package models

import "time"

type Bag struct {
	ID           uint      `gorm:"primaryKey"`
	QRCode       string    `gorm:"uniqueIndex"`    // Unique QR code for each bag (indexed for fast lookup)
	BagType      string    `gorm:"not null"`       // "Parent" or "Child" to distinguish bag type
	ChildCount   int       `gorm:"default:0"`      // (if used) number of child bags if this is a parent
	Linked       bool      `gorm:"default:false"`  // If Child bag, whether it is linked to a Parent (true once linked)
	ParentBagID  *uint     `gorm:"index"`          // Parent bag reference if this is a child bag (nullable) [oai_citation_attribution:20‡kapc.org.uk](https://kapc.org.uk/post/optional-relations-in-database-model-in-gorm#:~:text=performance.%20,or%20enforce%20unique%20tag%20names)
	LinkedToBill bool      `gorm:"default:false"`  // Flag to indicate if currently linked to a bill
	BillID       *uint     `gorm:"index"`          // Linked bill reference if any (nullable when not linked)
	CreatedAt    time.Time `gorm:"autoCreateTime"` // timestamp when the bag is created
	UpdatedAt    time.Time `gorm:"autoUpdateTime"` // timestamp for updates (if needed)
}

type Link struct {
	ID        uint   `gorm:"primaryKey"`
	ParentBag string `gorm:"not null"`
	BillID    string `gorm:"not null"`
	CreatedAt string `gorm:"default:CURRENT_TIMESTAMP"`
}
