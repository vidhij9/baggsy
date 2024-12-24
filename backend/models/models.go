package models

// Bag model for storing QR code and type
type Bag struct {
	ID      uint   `gorm:"primaryKey"`
	QRCode  string `gorm:"uniqueIndex;not null"`
	BagType string `gorm:"not null"` // Can be "Parent" or "Child"
}

// BagMap model for linking parent bags to bill IDs
type BagMap struct {
	ID         uint   `gorm:"primaryKey"`
	ParentBag  string `gorm:"not null"`             // QR Code of parent bag
	ChildBag   string `gorm:"not null"`             // QR Code of child bag
	BillID     string `gorm:"not null"`             // Associated bill ID
	UniqueLink string `gorm:"uniqueIndex;not null"` // Unique combination of ParentBag and ChildBag
}

// Link model for associating parent and child bags
type Link struct {
	ID        uint   `gorm:"primaryKey"`
	ParentBag string `gorm:"not null"`            // QR Code of the parent bag
	ChildBag  string `gorm:"not null"`            // QR Code of the child bag
	BillID    string `gorm:"not null"`            // Bill ID
	CreatedAt int64  `gorm:"autoCreateTime:nano"` // Timestamp of creation
}
