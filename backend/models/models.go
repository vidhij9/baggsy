package models

// Bag struct defines the structure of a bag
// including bill and parent-child relationships
type Bag struct {
	ID          string `json:"id"`
	QRCode      string `json:"qr_code"`
	BagType     string `json:"bag_type"`
	Status      string `json:"status"`
	BillID      string `json:"bill_id"`
	ParentBagID string `json:"parent_bag_id"`
}

// Bill struct defines the structure of a bill
// including associated metadata
type Bill struct {
	ID          string `json:"id"`
	SAPBillID   string `json:"sap_bill_id"`
	Description string `json:"description"`
}
