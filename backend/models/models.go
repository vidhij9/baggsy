package models

import (
	"gorm.io/gorm"
)

type Bag struct {
	ID         uint           `gorm:"primaryKey"`
	QRCode     string         `gorm:"uniqueIndex;not null"` // Unique QR Code for the bag
	BagType    string         `gorm:"not null"`             // "Parent" or "Child"
	ChildCount int            `gorm:"default:0"`            // Number of child bags for parent bags
	Linked     bool           `gorm:"default:false"`        // Indicates if the bag is linked to a bill
	ParentBag  string         `gorm:"default:null"`         // QR Code of the parent bag for child bags
	DeletedAt  gorm.DeletedAt `gorm:"index"`                // For soft delete functionality
}

type BagRequest struct {
	ParentBag string `json:"parentBag" binding:"required"` // Ensure this matches the JSON payload
	ChildBag  string `json:"childBag" binding:"required"`  // Ensure this matches the JSON payload
}

// type BagMap struct {
// 	ID        uint   `gorm:"primaryKey"`
// 	ParentBag string `gorm:"not null"`
// 	ChildBag  string `gorm:"not null"`
// 	CreatedAt string `gorm:"default:current_timestamp"`
// }

type Link struct {
	ID        uint   `gorm:"primaryKey"`
	ParentBag string `gorm:"not null"`
	BillID    string `gorm:"not null"`
	CreatedAt string `gorm:"default:CURRENT_TIMESTAMP"`
}
