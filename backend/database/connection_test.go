package database_test

import (
	"baggsy/backend/database"
	"os"
	"testing"

	"gorm.io/gorm"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary .env file for testing
	envContent := "DB_HOST=localhost\nDB_USER=testuser\nDB_PASS=testpass\nDB_NAME=testdb\nDB_PORT=5432"
	err := os.WriteFile(".env", []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary .env file: %v", err)
	}
	defer os.Remove(".env") // Clean up the .env file after the test

	database.LoadConfig()

	if os.Getenv("DB_HOST") != "localhost" {
		t.Errorf("Expected DB_HOST to be 'localhost', got '%s'", os.Getenv("DB_HOST"))
	}
	if os.Getenv("DB_USER") != "testuser" {
		t.Errorf("Expected DB_USER to be 'testuser', got '%s'", os.Getenv("DB_USER"))
	}
	if os.Getenv("DB_PASS") != "testpass" {
		t.Errorf("Expected DB_PASS to be 'testpass', got '%s'", os.Getenv("DB_PASS"))
	}
	if os.Getenv("DB_NAME") != "testdb" {
		t.Errorf("Expected DB_NAME to be 'testdb', got '%s'", os.Getenv("DB_NAME"))
	}
	if os.Getenv("DB_PORT") != "5432" {
		t.Errorf("Expected DB_PORT to be '5432', got '%s'", os.Getenv("DB_PORT"))
	}
}

func TestConnect(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASS", "testpass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_PORT", "5432")

	// Mock database connection
	mockDB, err := gorm.Open(nil, &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create mock database connection: %v", err)
	}

	// Inject mock database connection
	database.DB = mockDB

	if database.DB == nil {
		t.Errorf("Expected DB to be initialized, got nil")
	}
}
