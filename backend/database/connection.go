package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" // Import PostgreSQL driver
)

var DB *sql.DB

func InitDBWithRetry() {
	var err error
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	retries := 5
	for i := 0; i < retries; i++ {
		DB, err = sql.Open("postgres", connStr)
		if err == nil && DB.Ping() == nil {
			log.Println("Database connection established")
			return
		}
		log.Printf("Database connection failed (attempt %d/%d): %v", i+1, retries, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("Failed to connect to the database after %d attempts: %v", retries, err)
}
