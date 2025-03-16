package services

import (
	"d2t_server/core"
	"d2t_server/internal/models"
	"fmt"
)

// QAService 处理问答相关的业务逻辑
type QAService struct {
}

// NewQAService 创建一个新的QAService实例
func NewQAService() *QAService {
	return &QAService{}
}

// ProcessQuestion 处理问题并返回结果
func (s *QAService) ProcessQuestion(question string) (string, string, []map[string]interface{}, error) {
	// 处理自然语言查询
	sqlStr, analysisResult, err := core.ProcessNaturalLanguageQuery(question)
	if err != nil {
		return "", "", nil, fmt.Errorf("处理查询失败: %w", err)
	}

	// 获取数据库连接
	db, err := models.GetPGDBConnection()
	if err != nil {
		return "", "", nil, fmt.Errorf("数据库连接失败: %w", err)
	}

	// 执行SQL查询
	results, err := models.ExecuteSQL(db, sqlStr)
	if err != nil {
		return "", "", nil, fmt.Errorf("执行SQL失败: %w", err)
	}

	return sqlStr, analysisResult, results, nil
}
