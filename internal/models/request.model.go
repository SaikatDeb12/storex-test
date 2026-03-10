package models

type RegisterRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=50"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8,max=20"`
	PhoneNumber string `json:"phoneNumber" validate:"required,len=10"`
	Role        string `json:"role" validate:"required,oneof=admin employee project_manager asset_manager employee_manager"`
	Employment  string `json:"employment" validate:"required,oneof=full_time intern freelancer"`
}

type LoginRequest struct {
	Email    string `json:"email" db:"email" validate:"required,email"`
	Password string `json:"password" db:"password" validate:"required,min=8,max=20"`
}

type RequestContext struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	Role      string `json:"role"`
}

type UserInfoRequest struct {
	ID           string             `json:"id" db:"id"`
	Name         string             `json:"name" db:"name" validate:"required,min=3,max=50"`
	Email        string             `json:"email" db:"email" validate:"required,email"`
	PhoneNumber  string             `json:"phoneNumber" db:"phone_number" validate:"required,len=10"`
	Role         string             `json:"role" db:"role" validate:"required"`
	Employment   string             `json:"employment" db:"employment" validate:"required"`
	CreatedAt    string             `json:"createdAt" db:"created_at" validate:"required"`
	AssetDetails []AssetInfoRequest `json:"assetDetails"`
}

type CreateAssetRequest struct {
	Brand         string `json:"brand" db:"brand" validate:"required"`
	Model         string `json:"model" db:"model" validate:"required"`
	SerialNumber  string `json:"serialNumber" db:"serial_number" validate:"required"`
	Type          string `json:"assetType" db:"asset_type" validate:"required,oneof=laptop keyboard mouse mobile"`
	Owner         string `json:"owner" db:"owner_type" validate:"required,oneof=client remotestate"`
	WarrantyStart string `json:"warrantyStart" db:"warranty_start" validate:"required"`
	WarrantyEnd   string `json:"warrantyEnd" db:"warranty_end" validate:"required"`

	Laptop   *LaptopRequest   `json:"laptop"`
	Keyboard *KeyboardRequest `json:"keyboard"`
	Mouse    *MouseRequest    `json:"mouse"`
	Mobile   *MobileRequest   `json:"mobile"`
}

type AssetInfoRequest struct {
	ID     string `json:"id" db:"id"`
	Brand  string `json:"brand" db:"brand"`
	Model  string `json:"model" db:"model"`
	Status string `json:"status" db:"status"`
	Type   string `json:"assetType" db:"asset_type"`
}

type UpdateAssetRequest struct {
	Brand         string `json:"brand" db:"brand" validate:"required"`
	Model         string `json:"model" db:"model" validate:"required"`
	SerialNumber  string `json:"serialNumber" db:"serial_number" validate:"required"`
	Type          string `json:"assetType" db:"asset_type" validate:"required,oneof=laptop keyboard mouse mobile"`
	Status        string `json:"status" db:"status" validate:"required,oneof=available assigned in_service under_repair damaged"`
	Owner         string `json:"owner" db:"owner_type" validate:"required,oneof=client remotestate"`
	WarrantyStart string `json:"warrantyStart" db:"warranty_start" validate:"required"`
	WarrantyEnd   string `json:"warrantyEnd" db:"warranty_end" validate:"required"`

	Laptop   *LaptopRequest   `json:"laptop"`
	Keyboard *KeyboardRequest `json:"keyboard"`
	Mouse    *MouseRequest    `json:"mouse"`
	Mobile   *MobileRequest   `json:"mobile"`
}

type LaptopRequest struct {
	AssetID         string  `json:"assetID" db:"asset_id"`
	Processor       string  `json:"processor" db:"processor" validate:"required"`
	RAM             string  `json:"ram" db:"ram" validate:"required"`
	Storage         string  `json:"storage" db:"storage" validate:"required"`
	OperatingSystem string  `json:"operatingSystem" db:"operating_system" validate:"required"`
	Charger         *string `json:"charger" db:"charger"`
	DevicePassword  string  `json:"devicePassword" db:"device_password" validate:"required"`
}

type KeyboardRequest struct {
	AssetID      string  `json:"assetID" db:"asset_id"`
	Layout       *string `json:"layout" db:"layout" validate:"required"`
	Connectivity string  `json:"connectivity" db:"connectivity" validate:"required"`
}

type MouseRequest struct {
	AssetID      string `json:"assetID" db:"asset_id"`
	DPI          *int   `json:"dpi" db:"dpi" validate:"required"`
	Connectivity string `json:"connectivity" db:"connectivity" validate:"required"`
}

type MobileRequest struct {
	AssetID         string  `json:"assetID" db:"asset_id"`
	OperatingSystem string  `json:"operatingSystem" db:"operating_system" validate:"required"`
	RAM             string  `json:"ram" db:"ram" validate:"required"`
	Storage         string  `json:"storage" db:"storage" validate:"required"`
	Charger         *string `json:"charger" db:"charger" validate:"required"`
	DevicePassword  string  `json:"devicePassword" db:"device_password" validate:"required"`
}

type AllAssetsInfoRequest struct {
	ID            string  `json:"id" db:"id"`
	Brand         string  `json:"brand" db:"brand" validate:"required"`
	Model         string  `json:"model" db:"model" validate:"required"`
	SerialNumber  string  `json:"serialNumber" db:"serial_number" validate:"required"`
	Type          string  `json:"assetType" db:"asset_type" validate:"required, oneof=laptop keyboard mouse mobile"`
	Status        string  `json:"status" db:"status" validate:"required,oneof=available assigned in_sercice under_repair damaged"`
	Owner         string  `json:"owner" db:"owner_type" validate:"required,oneof=client remotestate"`
	AssignedByID  *string `json:"assignedByID" db:"assigned_by_id"`
	AssignedToID  *string `json:"assignedTo" db:"assigned_to_id"`
	AssignedAt    *string `json:"assignedAt" db:"assigned_at"`
	WarrantyStart string  `json:"warrantyStart" db:"warranty_start" validate:"required"`
	WarrantyEnd   string  `json:"warrantyEnd" db:"warranty_end" validate:"required"`
	ServiceStart  *string `json:"serviceStart" db:"service_start"`
	ServiceEnd    *string `json:"serviceEnd" db:"service_end"`
	ReturnedAt    *string `json:"returnedAt" db:"returned_at"`
	CreatedAt     string  `db:"created_at"`
	UpdatedAt     *string `db:"updated_at"`
}

type AssetAssignRequest struct {
	AssetID string `json:"assetID" db:"asset_id"`
	UserID  string `json:"userID" db:"user_id"`
}

type DashboardSummaryRequest struct {
	Total            int `json:"total"`
	Available        int `json:"available"`
	Assigned         int `json:"assigned"`
	WaitingForRepair int `json:"waitingForRepair"`
	InService        int `json:"inService"`
	Damaged          int `json:"damaged"`
}

type DashboardData struct {
	Summary DashboardSummaryRequest `json:"summary"`
	Assets  []AllAssetsInfoRequest  `json:"assetInfo"`
}

type SentServiceRequest struct {
	StartDate string `json:"start_date" db:"service_start" validate:"required"`
	EndDate   string `json:"end_date" db:"service_end" validate:"required"`
}

type ErrorModel struct {
	Error      string
	Message    string
	StatusCode int
}
