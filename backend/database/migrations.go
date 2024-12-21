package database

import "log"

func RunMigrations() {
	query := `
	CREATE TABLE IF NOT EXISTS bags (
		id SERIAL PRIMARY KEY,
		qr_code VARCHAR(255) UNIQUE NOT NULL,
		bag_type VARCHAR(50) NOT NULL,
		status VARCHAR(50)
	);`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migrations ran successfully.")
}
