package core

import (
	"d2t_server/utils"
	"fmt"
	"log"
	// Import any database packages needed for executing SQL
)

// Store the database schema as a constant for NL to SQL conversion
const DatabaseSchema = `
-- Customers table
CREATE TABLE Customers
(
  cust_id      char(10)  NOT NULL ,
  cust_name    char(50)  NOT NULL ,
  cust_address char(50)  ,
  cust_city    char(50)  ,
  cust_state   char(5)   ,
  cust_zip     char(10)  ,
  cust_country char(50)  ,
  cust_contact char(50)  ,
  cust_email   char(255) 
);

-- OrderItems table
CREATE TABLE OrderItems
(
  order_num  int          NOT NULL ,
  order_item int          NOT NULL ,
  prod_id    char(10)     NOT NULL ,
  quantity   int          NOT NULL ,
  item_price decimal(8,2) NOT NULL 
);

-- Orders table
CREATE TABLE Orders
(
  order_num  int      NOT NULL ,
  order_date date     NOT NULL ,
  cust_id    char(10) NOT NULL 
);

-- Products table
CREATE TABLE Products
(
  prod_id    char(10)      NOT NULL ,
  vend_id    char(10)      NOT NULL ,
  prod_name  char(255)     NOT NULL ,
  prod_price decimal(8,2)  NOT NULL ,
  prod_desc  varchar(1000) NULL 
);

-- Vendors table
CREATE TABLE Vendors
(
  vend_id      char(10) NOT NULL ,
  vend_name    char(50) NOT NULL ,
  vend_address char(50) NULL ,
  vend_city    char(50) NULL ,
  vend_state   char(5)  NULL ,
  vend_zip     char(10) NULL ,
  vend_country char(50) NULL 
);

-- Primary keys
ALTER TABLE Customers ADD PRIMARY KEY (cust_id);
ALTER TABLE OrderItems ADD PRIMARY KEY (order_num, order_item);
ALTER TABLE Orders ADD PRIMARY KEY (order_num);
ALTER TABLE Products ADD PRIMARY KEY (prod_id);
ALTER TABLE Vendors ADD PRIMARY KEY (vend_id);

-- Foreign keys
ALTER TABLE OrderItems ADD CONSTRAINT FK_OrderItems_Orders FOREIGN KEY (order_num) REFERENCES Orders (order_num);
ALTER TABLE OrderItems ADD CONSTRAINT FK_OrderItems_Products FOREIGN KEY (prod_id) REFERENCES Products (prod_id);
ALTER TABLE Orders ADD CONSTRAINT FK_Orders_Customers FOREIGN KEY (cust_id) REFERENCES Customers (cust_id);
ALTER TABLE Products ADD CONSTRAINT FK_Products_Vendors FOREIGN KEY (vend_id) REFERENCES Vendors (vend_id);
`

// ProcessNaturalLanguageQuery takes a natural language query, converts it to SQL and executes it
// Returns both the generated SQL and the query result
func ProcessNaturalLanguageQuery(nlQuery string) (string, string, error) {
	// Step 1: Convert natural language to SQL using Deepseek with database schema
	sqlQuery, err := utils.DeepseekRequest(nlQuery, "nl2sql_with_schema", DatabaseSchema)
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
