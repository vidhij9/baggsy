package database

import (
	"log"
	"os",
	"strings"
)

func RunMigrations() {
    // Read the SQL file
    sqlFile, err := os.ReadFile("database/init.sql")
    if err != nil {
        log.Fatalf("Failed to read SQL file: %v", err)
    }

    // Convert the file content to string and split into individual statements
    statements := strings.Split(string(sqlFile), ";")

    // Execute each statement
    for _, statement := range statements {
        // Skip empty statements
        if strings.TrimSpace(statement) == "" {
            continue
        }

        _, err = DB.Exec(statement)
        if err != nil {
            log.Fatalf("Migration failed: %v", err)
        }
    }
    
    log.Println("Migrations ran successfully.")
}

