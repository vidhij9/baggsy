package database

import (
	"log"

	"baggsy/backend/models"
)

func RunMigrations() error {
	// Check for the existence of the constraint
	var exists int
	constraintCheckQuery := `
		SELECT COUNT(*)
		FROM pg_constraint
		WHERE conname = 'uni_bags_qr_code'
		AND conrelid = 'bags'::regclass;
	`

	if err := DB.Raw(constraintCheckQuery).Scan(&exists).Error; err != nil {
		log.Printf("Error checking for constraint existence: %v", err)
		return err
	}

	if exists > 0 {
		log.Println("Dropping existing constraint: uni_bags_qr_code")
		if err := DB.Exec(`ALTER TABLE "bags" DROP CONSTRAINT "uni_bags_qr_code";`).Error; err != nil {
			log.Printf("Error dropping constraint: %v", err)
			return err
		}
	} else {
		log.Println("Constraint uni_bags_qr_code does not exist. Skipping drop.")
	}

	// Migrate tables: Bags, BagMap, and Links
	if err := DB.AutoMigrate(&models.BagMap{}, &models.Link{}, &models.Bag{}); err != nil {
		log.Printf("Error during AutoMigrate: %v", err)
		return err
	}

	log.Println("Database migration completed successfully!")
	return nil
}
