package models

import "time"

// Bill represents a billing record in the database.
type Bill struct {
	ID        int       `gorm:"primaryKey"`  // Primary Key
	BillCode  string    `gorm:"uniqueIndex"` // Unique bill identifier
	CreatedAt time.Time // Timestamp of creation
	UpdatedAt time.Time // Timestamp of last update
}
