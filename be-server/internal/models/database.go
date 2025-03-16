package models

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// GetPGDBConnection 创建并返回一个PostgreSQL数据库连接
func GetPGDBConnection() (*sql.DB, error) {
	// 从环境变量获取数据库连接信息
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// 如果没有设置端口，使用默认端口5432
	if dbPort == "" {
		dbPort = "5432"
	}

	// 构建连接字符串
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// 打开数据库连接
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("无法连接到数据库: %w", err)
	}

	// 设置连接池参数
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	// 测试连接
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("无法ping数据库: %w", err)
	}

	return db, nil
}

// ExecuteSQL 执行SQL语句并返回结果
func ExecuteSQL(db *sql.DB, sqlQuery string) ([]map[string]interface{}, error) {
	rows, err := db.Query(sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("执行SQL错误: %w", err)
	}
	defer rows.Close()

	// 获取列名
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("获取列名错误: %w", err)
	}

	// 结果集
	var results []map[string]interface{}

	// 为每行数据创建一个值的容器
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	// 遍历结果集
	for rows.Next() {
		// 扫描行数据到值容器
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, fmt.Errorf("扫描行数据错误: %w", err)
		}

		// 创建结果映射
		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}

		results = append(results, row)
	}

	// 检查遍历时是否有错误
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历结果集错误: %w", err)
	}

	return results, nil
}
