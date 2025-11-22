package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB
var err error

func DbConnect() {
	// Get environment variables (with defaults for local development)
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "root")
	dbName := getEnv("DB_NAME", "testdb")

	// Build connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Retry logic for Docker container startup
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		DB, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("Failed to open database connection (attempt %d/%d): %v", i+1, maxRetries, err)
			time.Sleep(2 * time.Second)
			continue
		}

		err = DB.Ping()
		if err == nil {
			log.Print("Database connected successfully")
			createTables()
			return
		}

		log.Printf("Failed to ping database (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatal("Could not connect to database after retries")
	panic(err.Error())
}

func createTables() {
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS articles (
			article_id INT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			author VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)
	`

	_, err := DB.Exec(createTableQuery)
	if err != nil {
		log.Printf("Failed to create articles table: %v", err)
	} else {
		log.Print("Articles table ready")
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
