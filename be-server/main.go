package main

import (
	"d2t_server/routes"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	envFile := flag.String("env", "", "Path to .env file")
	flag.Parse()

	// 加载.env文件
	var err error

	if *envFile != "" {
		err = godotenv.Load(*envFile)
		if err != nil {
			log.Printf("Warning: Error loading specified .env file (%s): %v", *envFile, err)
		}
	} else {
		// 首先尝试直接加载当前目录下的.env
		err = godotenv.Load()

		// 如果加载失败，尝试从项目根目录加载
		if err != nil {
			execDir, err := os.Executable()
			if err == nil {
				execPath := filepath.Dir(execDir)
				err = godotenv.Load(filepath.Join(execPath, "../.env"))
				if err != nil {
					err = godotenv.Load(filepath.Join(execPath, "../../.env"))
				}
			}

			// 如果仍然失败，记录一条警告但不中断程序
			if err != nil {
				log.Printf("Warning: Error loading .env file: %v", err)
			}
		}
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser != "" {
		log.Printf("Database configuration loaded successfully. Connected as user: %s", dbUser)
	}

	// 设置Gin框架
	r := gin.Default()

	// 注册其他路由
	routes.RegisterRoutes(r)

	// 读取端口env或默认8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s ...", port)

	// 启动服务器
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
