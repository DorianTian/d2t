package model

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL驱动，使用下划线导入表示仅注册驱动而不直接使用
)

func GetPGDBConnection() (*sql.DB, error) {
	// Get database connection parameters from environment variables
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost" // Default value
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		// 如果环境变量不存在，使用默认值
		user = "Dorian"
		// 不再返回错误，而是使用默认值
		// return nil, fmt.Errorf("DB_USER environment variable not set")
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		// 如果环境变量不存在，使用默认值
		password = "tqy4468"
		// 不再返回错误，而是使用默认值
		// return nil, fmt.Errorf("DB_PASSWORD environment variable not set")
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		// 如果环境变量不存在，使用默认值
		dbname = "d2t_db"
		// 不再返回错误，而是使用默认值
		// return nil, fmt.Errorf("DB_NAME environment variable not set")
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432" // Default PostgreSQL port
	}

	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable" // Default to disable for local development
	}

	// Construct the DSN string from environment variables
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func ClosePGDB(db *sql.DB) error {
	return db.Close()
}

func ExecuteSQL(db *sql.DB, sqlStr string) ([]map[string]interface{}, error) {
	defer ClosePGDB(db)

	log.Printf("Executing SQL: %s", sqlStr)
	rows, err := db.Query(sqlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	columns, _ := rows.Columns()

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		ptrs := make([]interface{}, len(columns))
		for i := range values {
			ptrs[i] = &values[i]
		}

		if err := rows.Scan(ptrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		rowMap := make(map[string]interface{})
		for i, col := range columns {
			// 如果值是[]byte类型，将其转换为字符串
			if byteArray, ok := values[i].([]byte); ok {
				rowMap[col] = string(byteArray)
			} else {
				rowMap[col] = values[i]
			}
		}

		results = append(results, rowMap)
	}

	return results, nil
}
