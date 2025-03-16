package main

import (
	"d2t_server/internal/api"
	"d2t_server/internal/config"
	"flag"
	"log"
)

func main() {
	// 解析命令行参数
	envFile := flag.String("env", "", "Path to .env file")
	flag.Parse()

	// 加载配置
	cfg, err := config.LoadConfig(*envFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 输出数据库连接信息
	if cfg.DB.User != "" {
		log.Printf("Database configuration loaded successfully. Connected as user: %s", cfg.DB.User)
	}

	// 创建服务器
	server := api.NewServer(cfg)

	// 启动服务器
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
