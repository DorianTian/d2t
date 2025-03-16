package routes

import (
	"d2t_server/internal/services"
	"d2t_server/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	// 健康检查路由
	r.GET("/health", HealthCheckHandler)
	r.GET("/ping", PingHandler)

	// API 路由
	api := r.Group("/api")
	{
		api.POST("/askQA", AskQAHandler)
		// 其他API路由可以添加在这里
	}
}

// PingHandler 处理Ping请求
func PingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// HealthCheckHandler 处理健康检查请求
func HealthCheckHandler(c *gin.Context) {
	utils.HealthCheckHandler(c)
}

// AskQAHandler 处理问答请求
func AskQAHandler(c *gin.Context) {
	fmt.Println("==== /api/askQA was called ====")

	// 定义请求结构体
	var req struct {
		Question string `json:"question"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 使用服务层处理问题
	qaService := services.NewQAService()
	sqlStr, analysisResult, results, err := qaService.ProcessQuestion(req.Question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results":  utils.TrimStringValues(results),
		"sql":      sqlStr,
		"analysis": analysisResult,
	})
}
