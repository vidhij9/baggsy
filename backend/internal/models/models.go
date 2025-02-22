package models

import "time"

type User struct {
	ID           uint   `gorm:"primary_key" json:"id"`
	Username     string `gorm:"unique;not null" json:"username"`
	PasswordHash string `gorm:"not null" json:"passwordHash"`
	Role         string `gorm:"type:user_role;not null" json:"role"` // Use custom ENUM type
}

type Bag struct {
	ID         uint      `gorm:"primary_key" json:"id"`
	QRCode     string    `gorm:"unique;not null" json:"qrCode"`
	Type       string    `gorm:"type:bag_type;not null" json:"type"` // Use custom ENUM type
	ChildCount int       `gorm:"default:0" json:"childCount"`
	ParentID   *uint     `gorm:"index" json:"parentId,omitempty"`
	Linked     bool      `gorm:"default:false" json:"linked"`
	CreatedAt  time.Time `gorm:"not null" json:"createdAt"`
	Children   []Bag     `gorm:"-" json:"children,omitempty"` // For response only
}

type Link struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	ParentID  uint      `gorm:"unique;not null;index" json:"parentId"`
	BillID    string    `gorm:"not null;index" json:"billId"`
	CreatedAt time.Time `gorm:"not null" json:"createdAt"`
}
