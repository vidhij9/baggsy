package models

import (
	"gorm.io/gorm"
)

type Bag struct {
	ID         uint           `gorm:"primaryKey"`
	QRCode     string         `gorm:"uniqueIndex;not null" json:"qrCode" binding:"required"`
	BagType    string         `gorm:"not null" json:"bagType" binding:"required"` // "Parent" or "Child"
	ChildCount int            `json:"childCount,omitempty"`                       // Optional: Only for Parent bags
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type BagMap struct {
	ID        uint   `gorm:"primaryKey"`
	ParentBag string `gorm:"not null"`
	ChildBag  string `gorm:"not null"`
	CreatedAt string `gorm:"default:current_timestamp"`
}

type Link struct {
	ID        uint   `gorm:"primaryKey"`
	ParentBag string `gorm:"not null"`
	BillID    string `gorm:"not null"`
	CreatedAt string `gorm:"default:current_timestamp"`
}
