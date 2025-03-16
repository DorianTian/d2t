package utils

import (
	"bytes"
	"d2t_server/internal/config"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status    string      `json:"status"`
	Timestamp string      `json:"timestamp"`
	Data      interface{} `json:"data,omitempty"`
}

func getAskURL() string {
	return "https://api.deepseek.com/chat/completions"
}

// DeepseekRequest handles all interactions with the Deepseek API
// mode: "nl2sql" for natural language to SQL conversion, "analyze" for SQL analysis, "nl2sql_with_schema" for NL to SQL with DB schema
func DeepseekRequest(input string, mode string, schema ...string) (string, error) {
	url := getAskURL()

	// Define request structure
	type Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	type RequestBody struct {
		Model    string    `json:"model"`
		Messages []Message `json:"messages"`
	}

	// Build request body based on mode
	var reqBody RequestBody
	reqBody.Model = "deepseek-chat"

	switch mode {
	case "nl2sql":
		reqBody.Messages = []Message{
			{
				Role:    "user",
				Content: "我想把如下问题转换成sql语句: " + input,
			},
		}
	case "nl2sql_with_schema":
		if len(schema) == 0 {
			return "", fmt.Errorf("schema is required for nl2sql_with_schema mode")
		}

		prettyJSON, _ := json.MarshalIndent(reqBody, "", "  ")
		log.Printf("请求体: %s", string(prettyJSON))
		log.Printf("DeepseekRequest: 模式=%s, 输入=%s", mode, input)
		reqBody.Messages = []Message{
			{
				Role:    "system",
				Content: "You are an SQL expert. Convert natural language questions to SQL queries. Use the provided database schema to create accurate queries. Only respond with valid SQL queries, no explanations.",
			},
			{
				Role:    "user",
				Content: "Database Schema:\n" + schema[0] + "\n\nConvert this question to SQL: " + input,
			},
		}
	case "analyze":
		reqBody.Messages = []Message{
			{
				Role:    "system",
				Content: "You are an SQL expert. Analyze the provided SQL query and explain its purpose, potential optimizations, and any issues it might have.",
			},
			{
				Role:    "user",
				Content: input,
			},
		}
	default:
		return "", fmt.Errorf("未知的操作模式: %s", mode)
	}

	// Serialize request body to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("json序列化失败: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Add API key
	apiKey := "sk-05560ff5cfcf472188f269aff3fb053b"
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	// Get API timeout config
	apiConfig := config.GetAPIConfig()

	// Send request with timeout
	client := &http.Client{
		Timeout: apiConfig.Timeout,
	}

	log.Printf("Making API request with timeout of %v seconds", apiConfig.Timeout.Seconds())
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	log.Printf("Deepseek响应: %v", response["choices"])
	if choices, ok := response["choices"].([]interface{}); ok && len(choices) > 0 {
		if firstChoice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := firstChoice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					return content, nil
				}
			}
		}
	}

	return "", fmt.Errorf("无法从响应中提取有效内容: %s", string(body))
}

// CleanSQLFromMarkdown removes markdown formatting from SQL strings
// It handles cases like "### Revised Query Example: ```sql SELECT * FROM table ```"
func CleanSQLFromMarkdown(sqlStr string) string {
	// 清除常见的markdown标记
	re := regexp.MustCompile("(?s).*?```sql(.+?)```.*")
	matches := re.FindStringSubmatch(sqlStr)

	// 如果找到了SQL代码块，提取其中内容
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	// 如果没有找到SQL代码块，尝试清除其他可能的markdown格式
	result := sqlStr

	// 去除标题（###, ##, #）
	result = regexp.MustCompile("(?m)^#+\\s+.*$").ReplaceAllString(result, "")

	// 去除"SQL Query:"等前缀
	result = regexp.MustCompile("(?i)(SQL Query:|Query:|Revised Query Example:)").ReplaceAllString(result, "")

	// 去除反引号
	result = strings.ReplaceAll(result, "`", "")

	// 去除多余空行并整理格式
	result = regexp.MustCompile(`(?m)^\s*$[\r\n]*`).ReplaceAllString(result, "")
	result = strings.TrimSpace(result)

	return result
}

// ExtractSQLFromMarkdown 是一个更复杂的版本，可以处理多种情况
func ExtractSQLFromMarkdown(markdown string) string {
	// 首先尝试提取代码块
	codeBlockPattern := regexp.MustCompile("(?s)```sql(.+?)```")
	matches := codeBlockPattern.FindStringSubmatch(markdown)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	// 如果没有代码块标记，尝试其他常见模式
	patterns := []string{
		// 尝试匹配 "SQL Query:" 或类似前缀后的内容
		"(?i)(?:SQL Query:|Query:|Revised Query:)\\s*(.+)",
		// 匹配任何看起来像SQL的内容 (SELECT, UPDATE, INSERT等开头)
		"(?i)(SELECT|UPDATE|INSERT|DELETE|CREATE|ALTER|DROP).+",
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(markdown)
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		} else if len(matches) == 1 {
			return strings.TrimSpace(matches[0])
		}
	}

	// 如果所有模式都失败，返回清理后的原始文本
	// 删除markdown标记
	cleaned := markdown
	cleaned = regexp.MustCompile(`#+ .*`).ReplaceAllString(cleaned, "")      // 删除标题
	cleaned = regexp.MustCompile("`").ReplaceAllString(cleaned, "")          // 删除反引号
	cleaned = regexp.MustCompile(`\*\*|\*|__`).ReplaceAllString(cleaned, "") // 删除粗体和斜体标记

	return strings.TrimSpace(cleaned)
}

// TrimStringValues processes the results from database queries and trims
// trailing whitespace from string values
func TrimStringValues(results []map[string]interface{}) []map[string]interface{} {
	trimmedResults := make([]map[string]interface{}, len(results))

	for i, row := range results {
		trimmedRow := make(map[string]interface{})

		for key, value := range row {
			// 如果是字符串类型，去除尾部空白
			if strValue, ok := value.(string); ok {
				trimmedRow[key] = strings.TrimSpace(strValue)
			} else {
				trimmedRow[key] = value
			}
		}

		trimmedResults[i] = trimmedRow
	}

	return trimmedResults
}

// 健康检查处理函数 - 使用Gin框架格式
func HealthCheckHandler(c *gin.Context) {
	// 构造响应
	response := Response{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)
}
