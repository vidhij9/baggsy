package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// ConnectDB tries to open a database connection with retries and returns the *sql.DB.
func ConnectDB(dsn string) (*sql.DB, error) {
	var db *sql.DB
	var err error
	maxRetries := 5

	for i := 1; i <= maxRetries; i++ {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Printf("Database connection attempt %d failed: %v", i, err)
		} else {
			// Try pinging the database to ensure the connection is valid
			pingErr := db.Ping()
			if pingErr == nil {
				log.Printf("Successfully connected to database on attempt %d", i)
				break // exit loop on success
			}
			// If ping failed, close this db handle and prepare to retry
			log.Printf("Database ping attempt %d failed: %v", i, pingErr)
			_ = db.Close() // close the opened connection before retry
			err = pingErr  // treat ping error as the current error
		}

		if i < maxRetries {
			// Exponential backoff before next attempt
			time.Sleep(time.Duration(i) * time.Second)
		}
	}

	if err != nil {
		// All attempts failed
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
	}

	// Configure connection pool settings
	db.SetMaxOpenConns(10)                  // max open connections (in-use + idle)
	db.SetMaxIdleConns(5)                   // max idle connections to retain
	db.SetConnMaxLifetime(30 * time.Minute) // recycle connections periodically
	log.Println("Database connection pool configured (MaxOpenConns=10, MaxIdleConns=5).")

	return db, nil
}

// CloseDB closes the given database connection cleanly.
func CloseDB(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}
