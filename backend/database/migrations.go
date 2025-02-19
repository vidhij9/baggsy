package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// runMigrations opens and executes the init.sql schema script.
// It returns an error if something fails, ensuring the caller can handle it.
func RunMigrations(db *sql.DB) error {
	// Read the SQL schema file
	schema, err := os.ReadFile("init.sql")
	if err != nil {
		log.Printf("Failed to read schema file: %v\n", err)
		return err
	}

	// Execute the schema within a transaction for atomicity
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Failed to start transaction: %v\n", err)
		return err
	}
	if _, err := tx.Exec(string(schema)); err != nil {
		tx.Rollback()
		log.Printf("Schema migration failed: %v\n", err)
		return err
	}
	// Commit if all statements executed successfully
	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v\n", err)
		return err
	}

	log.Println("Database schema migrated successfully.")
	return nil
}
