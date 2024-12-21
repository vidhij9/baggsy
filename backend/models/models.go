package models

type Bag struct {
	ID      string `json:"id"`
	QRCode  string `json:"qr_code"`
	BagType string `json:"bag_type"`
	Status  string `json:"status"`
}

type Bill struct {
	ID          string   `json:"id"`
	SAPBillID   string   `json:"sap_bill_id"`
	Description string   `json:"description"`
	ParentBags  []string `json:"parent_bags"`
}
