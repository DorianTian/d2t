package api

import (
	"d2t_server/internal/config"
	"d2t_server/internal/middleware"
	"d2t_server/internal/routes"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// Server 代表API服务器
type Server struct {
	router *gin.Engine
	config *config.Config
}

// NewServer 创建一个新的服务器实例
func NewServer(config *config.Config) *Server {
	router := gin.Default()

	// 添加中间件
	middleware.RegisterMiddleware(router)

	// 注册路由
	routes.RegisterRoutes(router)

	return &Server{
		router: router,
		config: config,
	}
}

// Start 启动服务器
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.config.Server.Port)
	log.Printf("Starting server on port %s ...", s.config.Server.Port)
	return s.router.Run(addr)
}
