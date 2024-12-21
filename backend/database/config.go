package database

import (
	"log"
	"os"

	"baggsy/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB establishes a connection to the database
func ConnectDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=password dbname=baggsy port=5432 sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	DB = db
	log.Println("Database connected successfully")
	return db
}

// RunMigrations runs database migrations
func RunMigrations(db *gorm.DB) {
	log.Println("Running database migrations...")
	if err := db.AutoMigrate(&models.Bag{}, &models.Bill{}, &models.Distributor{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migrations completed successfully")
}
