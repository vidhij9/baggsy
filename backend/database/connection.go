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
	// Build DSN from environment variables
	// host := os.Getenv("DB_HOST")
	// port := os.Getenv("DB_PORT")
	// user := os.Getenv("DB_USER")
	// password := os.Getenv("DB_PASSWORD")
	// dbName := os.Getenv("DB_NAME")

	// dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	// 	host, port, user, password, dbName)

	dsn := "host=localhost user=baggsy password=baggsy dbname=baggsy port=5432 sslmode=disable TimeZone=Asia/Kolkata"

	retryCount := 5
	for i := 0; i < retryCount; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("Database connection established")
			return
		}

		// Setup connection pooling
		sqlDB, err := DB.DB()
		if err != nil {
			log.Fatalf("Failed to get sql.DB from GORM: %v", err)
		}
		sqlDB.SetMaxOpenConns(100) // max number of open connections
		sqlDB.SetMaxIdleConns(50)  // max number of idle connections
		sqlDB.SetConnMaxLifetime(5 * time.Minute)

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
