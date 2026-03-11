package models

import "time"

type Asset struct {
	ID            string     `json:"id" db:"id"`
	Brand         string     `json:"brand" db:"brand"`
	Model         string     `json:"model" db:"model"`
	SerialNumber  string     `json:"serialNumber" db:"serial_number"`
	Type          string     `json:"assetType" db:"asset_type"`
	Status        string     `json:"status" db:"status"`
	Owner         string     `json:"owner" db:"owner_type"`
	AssignedByID  *string    `json:"assignedByID" db:"assigned_by_id"`
	AssignedToID  *string    `json:"assignedTo" db:"assigned_to_id"`
	AssignedAt    *time.Time `json:"assignedAt" db:"assigned_at"`
	WarrantyStart time.Time  `json:"warrantyStart" db:"warranty_start"`
	WarrantyEnd   time.Time  `json:"warrantyEnd" db:"warranty_end"`
	ServiceStart  *time.Time `json:"serviceStart" db:"service_start"`
	ServiceEnd    *time.Time `json:"serviceEnd" db:"service_end"`
	ReturnedAt    *time.Time `json:"returnedAt" db:"returned_at"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
	ArchivedAt    *time.Time `db:"archived_at"`
	ArchivedByID  *string    `db:"archived_by_id"`
}

type Laptop struct {
	ID              string  `json:"id" db:"id"`
	AssetID         string  `json:"assetID" db:"asset_id"`
	Processor       string  `json:"processor" db:"processor"`
	RAM             string  `json:"ram" db:"ram"`
	Storage         string  `json:"storage" db:"storage"`
	OperatingSystem string  `json:"operatingSystem" db:"operating_system"`
	Charger         *string `json:"charger" db:"charger"`
	DevicePassword  string  `json:"devicePassword" db:"device_password"`
}

type Keyboard struct {
	ID           string  `json:"id" db:"id"`
	AssetID      string  `json:"assetID" db:"asset_id"`
	Layout       *string `json:"layout" db:"layout"`
	Connectivity string  `json:"connectivity" db:"connectivity"`
}

type Mouse struct {
	ID           string `json:"id" db:"id"`
	AssetID      string `json:"assetID" db:"asset_id"`
	DPI          *int   `json:"dpi" db:"dpi"`
	Connectivity string `json:"connectivity" db:"connectivity"`
}

type Mobile struct {
	ID              string  `json:"id" db:"id"`
	AssetID         string  `json:"assetID" db:"asset_id"`
	OperatingSystem string  `json:"operatingSystem" db:"operating_system"`
	RAM             string  `json:"ram" db:"ram"`
	Storage         string  `json:"storage" db:"storage"`
	Charger         *string `json:"charger" db:"charger"`
	DevicePassword  string  `json:"devicePassword" db:"device_password"`
}
