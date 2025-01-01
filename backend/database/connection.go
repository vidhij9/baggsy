package database

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func Connect() {
	var err error

	// PostgreSQL connection string
	// Replace the placeholders with your actual credentials

	// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
	// 	os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
	// 	os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	dsn := "host=localhost user=baggsy password=baggsy dbname=baggsy port=5432 sslmode=disable TimeZone=Asia/Kolkata"

	retryCount := 5
	for i := 0; i < retryCount; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("Database connection established")
			return
		}

		log.Printf("Failed to connect to database. Retrying... (%d/%d)\n", i+1, retryCount)
		time.Sleep(2 * time.Second) // Wait for 2 seconds before retrying
	}

	log.Fatal("Failed to connect to database after retries:", err)
}

func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
