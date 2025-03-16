package core

import (
	"d2t_server/internal/config"
	"d2t_server/utils"
	"fmt"
	"log"
	// Import any database packages needed for executing SQL
)

// ProcessNaturalLanguageQuery takes a natural language query, converts it to SQL and executes it
// Returns both the generated SQL and the query result
func ProcessNaturalLanguageQuery(nlQuery string) (string, string, error) {
	// Step 1: Convert natural language to SQL using Deepseek with database schema
	sqlQuery, err := utils.DeepseekRequest(nlQuery, "nl2sql_with_schema", config.DatabaseSchema)
	if err != nil {
		return "", "", fmt.Errorf("failed to convert natural language to SQL: %w", err)
	}

	if sqlQuery == "" {
		return "", "", fmt.Errorf("empty SQL query returned from Deepseek")
	} else {
		sqlQuery = utils.CleanSQLFromMarkdown(sqlQuery)

	}

	log.Printf("Natural language query converted to SQL: '%s'", sqlQuery)

	// Step 2: Execute the SQL query against the database
	// Note: This should be implemented based on your database setup
	// For now, we're returning the SQL without executing it

	// Step 3 (Optional): Analyze the SQL query for additional insights
	analysisResult, err := utils.DeepseekRequest(sqlQuery, "analyze")
	if err != nil {
		// We don't fail the whole process if analysis fails
		log.Printf("Warning: failed to analyze SQL query: %v", err)
		analysisResult = "No analysis available"
	}

	return sqlQuery, analysisResult, nil
}
