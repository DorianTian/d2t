package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config 结构体包含应用程序的所有配置
type Config struct {
	Server ServerConfig
	DB     DBConfig
	API    APIConfig
}

// ServerConfig 服务器相关配置
type ServerConfig struct {
	Port string
}

// DBConfig 数据库相关配置
type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

// APIConfig API相关配置
type APIConfig struct {
	Timeout time.Duration
}

// LoadConfig 加载配置信息
func LoadConfig(envFile string) (*Config, error) {
	// 加载环境变量
	if envFile != "" {
		err := godotenv.Load(envFile)
		if err != nil {
			log.Printf("Warning: Error loading specified .env file (%s): %v", envFile, err)
		}
	} else {
		// 首先尝试直接加载当前目录下的.env
		err := godotenv.Load()

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

	// 读取API超时设置（默认30秒）
	apiTimeoutStr := getEnv("API_TIMEOUT_SECONDS", "300")
	apiTimeout, err := strconv.Atoi(apiTimeoutStr)
	if err != nil {
		log.Printf("Warning: Invalid API_TIMEOUT_SECONDS value: %v, using default 300 seconds", err)
		apiTimeout = 300
	}

	// 读取配置
	config := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
		DB: DBConfig{
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Name:     os.Getenv("DB_NAME"),
		},
		API: APIConfig{
			Timeout: time.Duration(apiTimeout) * time.Second,
		},
	}

	return config, nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetAPIConfig returns the API configuration
func GetAPIConfig() APIConfig {
	// Load config if not already loaded
	config, err := LoadConfig("")
	if err != nil {
		log.Printf("Warning: Failed to load config, using default API timeout (30s)")
		return APIConfig{
			Timeout: 30 * time.Second,
		}
	}
	return config.API
}
