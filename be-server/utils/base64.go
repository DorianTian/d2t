package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

// DecodeBase64FromResults processes a result map or array of maps
// and decodes any Base64 encoded string values
func DecodeBase64FromResults(results interface{}) interface{} {
	// 处理结果为数组的情况
	if resultsList, ok := results.([]interface{}); ok {
		decodedList := make([]interface{}, len(resultsList))
		for i, item := range resultsList {
			decodedList[i] = DecodeBase64FromResults(item)
		}
		return decodedList
	}

	// 处理结果为单个对象的情况
	if resultsMap, ok := results.(map[string]interface{}); ok {
		decodedMap := make(map[string]interface{})
		for key, value := range resultsMap {
			// 递归处理嵌套对象
			if nestedMap, isMap := value.(map[string]interface{}); isMap {
				decodedMap[key] = DecodeBase64FromResults(nestedMap)
				continue
			}

			// 递归处理嵌套数组
			if nestedList, isList := value.([]interface{}); isList {
				decodedMap[key] = DecodeBase64FromResults(nestedList)
				continue
			}

			// 处理字符串值（可能是Base64编码）
			if strValue, isString := value.(string); isString {
				decodedMap[key] = DecodeBase64IfNeeded(strValue)
				continue
			}

			// 其他类型保持不变
			decodedMap[key] = value
		}
		return decodedMap
	}

	// 如果不是数组或对象，直接返回原值
	return results
}

// DecodeBase64IfNeeded determines if a string is Base64 encoded
// and decodes it if necessary
func DecodeBase64IfNeeded(input string) string {
	// 检查字符串是否可能是Base64编码
	if !isBase64Encoded(input) {
		return input
	}

	// 尝试解码
	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		// 如果解码失败，返回原始字符串
		return input
	}

	// 解码成功，清理结果（去除尾部空白）
	return strings.TrimRight(string(decoded), " \t\r\n\x00")
}

// isBase64Encoded checks if a string appears to be Base64 encoded
func isBase64Encoded(s string) bool {
	// 检查是否是空字符串
	if s == "" {
		return false
	}

	// 基本的Base64格式检查：长度是4的倍数，并且只包含Base64字符集
	if len(s)%4 != 0 {
		return false
	}

	// 使用正则表达式检查是否符合Base64编码模式
	matched, _ := regexp.MatchString("^[A-Za-z0-9+/]*={0,2}$", s)
	if !matched {
		return false
	}

	// 额外检查：确保解码后是可打印字符（避免误解码二进制数据）
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return false
	}

	// 简单检查解码后的字符串是否包含大量不可打印字符
	isPrintable := true
	nonPrintableCount := 0

	for _, b := range decoded {
		// 检查是否是ASCII可打印字符或常见控制字符
		if (b < 32 || b > 126) && b != 9 && b != 10 && b != 13 {
			nonPrintableCount++
		}
	}

	// 如果不可打印字符超过20%，可能不是文本内容的Base64
	if float64(nonPrintableCount)/float64(len(decoded)) > 0.2 {
		isPrintable = false
	}

	return isPrintable
}

// DecodeJSONWithBase64 unmarshals a JSON string and decodes any Base64 values
func DecodeJSONWithBase64(jsonStr string) (interface{}, error) {
	var result interface{}

	// 解析JSON
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	// 解码Base64值
	decoded := DecodeBase64FromResults(result)
	return decoded, nil
}

// ProcessAndDecodeResponse is a helper function that processes API response and automatically decodes any Base64 encoded values
func ProcessAndDecodeResponse(resp *http.Response) (interface{}, error) {
	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析并解码JSON
	result, err := DecodeJSONWithBase64(string(body))
	if err != nil {
		return nil, err
	}

	return result, nil
}
