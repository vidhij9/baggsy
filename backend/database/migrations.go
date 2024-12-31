package database

import (
	"baggsy/backend/models"
	"log"
)

func RunMigrations() error {
	log.Println("Starting database migrations...")

	// AutoMigrate ensures the schema is created or updated
	if err := DB.AutoMigrate(&models.Bag{}, &models.BagMap{}, &models.Link{}); err != nil {
		log.Printf("Error during AutoMigrate: %v", err)
		return err
	}

	log.Println("Database migrations completed successfully!")
	return nil
}
