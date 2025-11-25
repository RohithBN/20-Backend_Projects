package lib

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
}

func createTables() {
    // Create notes table
    createNoteTable := `
    CREATE TABLE IF NOT EXISTS notes (
        id INT AUTO_INCREMENT PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        markdown_content TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        INDEX idx_created_at (created_at),
        INDEX idx_title (title)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
    `

    // Create attachments table
    createAttachmentTable := `
    CREATE TABLE IF NOT EXISTS attachments (
        id INT AUTO_INCREMENT PRIMARY KEY,
        note_id INT NOT NULL,
        original_file_name VARCHAR(255) NOT NULL,
        stored_file_name VARCHAR(255) NOT NULL UNIQUE,
        file_url VARCHAR(512) NOT NULL,
        uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        mime_type VARCHAR(100),
        size BIGINT DEFAULT 0,
        FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
        INDEX idx_note_id (note_id),
        INDEX idx_uploaded_at (uploaded_at)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
    `

    // Execute notes table creation
    _, err := DB.Exec(createNoteTable)
    if err != nil {
        log.Printf("Error creating notes table: %v", err)
    } else {
        log.Println("Notes table created or already exists")
    }

    // Execute attachments table creation
    _, err = DB.Exec(createAttachmentTable)
    if err != nil {
        log.Printf("Error creating attachments table: %v", err)
    } else {
        log.Println("Attachments table created or already exists")
    }
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}