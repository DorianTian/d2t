package routes

import (
	"d2t_server/core"
	"d2t_server/model"
	"d2t_server/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.POST("/api/askQA", func(c *gin.Context) {
		fmt.Println("==== /api/askQA was called ====")

		// 定义请求结构体
		var req struct {
			Question string `json:"question"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		sqlStr, analysisResult, err := core.ProcessNaturalLanguageQuery(req.Question)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		db, err := model.GetPGDBConnection()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		results, err := model.ExecuteSQL(db, sqlStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"results":  utils.TrimStringValues(results),
			"sql":      sqlStr,
			"analysis": analysisResult,
		})
	})
}
