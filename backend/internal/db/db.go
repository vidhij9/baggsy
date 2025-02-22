package db

import (
	"fmt"
	"time"

	"baggsy/backend/internal/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	var err error
	DB, err = gorm.Open("postgres", "host=localhost user=baggsy password=baggsy dbname=baggsy sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	sqlDB := DB.DB()
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Explicitly create ENUM types and tables
	fmt.Println("Creating ENUM types and schema...")
	err = DB.Exec(`
		DO $$ BEGIN
			CREATE TYPE user_role AS ENUM ('employee', 'admin');
		EXCEPTION
			WHEN duplicate_object THEN NULL;
		END $$;

		DO $$ BEGIN
			CREATE TYPE bag_type AS ENUM ('parent', 'child');
		EXCEPTION
			WHEN duplicate_object THEN NULL;
		END $$;

		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			role user_role NOT NULL
		);

		CREATE TABLE IF NOT EXISTS bags (
			id SERIAL PRIMARY KEY,
			qr_code VARCHAR(255) UNIQUE NOT NULL,
			type bag_type NOT NULL,
			child_count INT DEFAULT 0 CHECK (child_count >= 0),
			parent_id INT,
			linked BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES bags(id),
			CONSTRAINT check_parent_child CHECK ((type = 'parent' AND parent_id IS NULL) OR (type = 'child' AND parent_id IS NOT NULL))
		);
		CREATE INDEX IF NOT EXISTS idx_bags_qr_code ON bags(lower(qr_code));
		CREATE INDEX IF NOT EXISTS idx_bags_parent_id ON bags(parent_id);

		CREATE TABLE IF NOT EXISTS links (
			id SERIAL PRIMARY KEY,
			parent_id INT UNIQUE NOT NULL,
			bill_id VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			CONSTRAINT fk_link_parent FOREIGN KEY (parent_id) REFERENCES bags(id)
		);
		CREATE INDEX IF NOT EXISTS idx_links_bill_id ON links(bill_id);
	`).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create schema: %v", err)
	}

	// AutoMigrate to ensure GORM models align with schema
	fmt.Println("Running AutoMigrate for model consistency...")
	err = DB.AutoMigrate(&models.User{}, &models.Bag{}, &models.Link{}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to auto-migrate: %v", err)
	}

	// Create initial admin user
	var admin models.User
	if err := DB.Where("username = ?", "admin").First(&admin).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			fmt.Println("Creating initial admin user...")
			hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
			if err != nil {
				return nil, fmt.Errorf("failed to hash password: %v", err)
			}
			adminUser := models.User{
				Username:     "admin",
				PasswordHash: string(hash),
				Role:         "admin",
			}
			if err := DB.Create(&adminUser).Error; err != nil {
				return nil, fmt.Errorf("failed to create admin user: %v", err)
			}
			fmt.Println("Admin user created successfully.")
		} else {
			return nil, fmt.Errorf("error checking for admin user: %v", err)
		}
	} else {
		fmt.Println("Admin user already exists.")
	}

	return DB, nil
}
