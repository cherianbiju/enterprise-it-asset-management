package models

import "time"

type AssetStatus string

const (
	StatusAvailable   AssetStatus = "AVAILABLE"
	StatusAssigned    AssetStatus = "ASSIGNED"
	StatusUnderRepair AssetStatus = "UNDER_REPAIR"
	StatusRetired     AssetStatus = "RETIRED"
)

func (s AssetStatus) IsValid() bool {
	switch s {
	case StatusAvailable, StatusAssigned, StatusUnderRepair, StatusRetired:
		return true
	default:
		return false
	}
}

type Asset struct {
	ID           int         `json:"id"`
	AssetCode    string      `json:"asset_code"`
	AssetName    string      `json:"asset_name"`
	AssetType    string      `json:"asset_type"`
	SerialNumber string      `json:"serial_number"`
	PurchaseDate string      `json:"purchase_date"`
	Status       AssetStatus `json:"status"`
	AssignedTo   *string     `json:"assigned_to"`
	Location     string      `json:"location"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

type AssetHistory struct {
	ID          int       `json:"id"`
	AssetID     int       `json:"asset_id"`
	Action      string    `json:"action"`
	PerformedBy string    `json:"performed_by"`
	Details     string    `json:"details"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateAssetRequest struct {
	AssetCode    string      `json:"asset_code"`
	AssetName    string      `json:"asset_name"`
	AssetType    string      `json:"asset_type"`
	SerialNumber string      `json:"serial_number"`
	PurchaseDate string      `json:"purchase_date"`
	Status       AssetStatus `json:"status,omitempty"`
	AssignedTo   *string     `json:"assigned_to,omitempty"`
	Location     string      `json:"location"`
}

type UpdateAssetRequest struct {
	AssetName    string `json:"asset_name"`
	AssetType    string `json:"asset_type"`
	PurchaseDate string `json:"purchase_date"`
	Location     string `json:"location"`
}

type AssignAssetRequest struct {
	AssignedTo  string `json:"assigned_to"`
	Location    string `json:"location,omitempty"`
	Notes       string `json:"notes,omitempty"`
	PerformedBy string `json:"performed_by,omitempty"`
}

type ReturnAssetRequest struct {
	Location    string `json:"location,omitempty"`
	Notes       string `json:"notes,omitempty"`
	PerformedBy string `json:"performed_by,omitempty"`
}

type StatusChangeRequest struct {
	Location    string `json:"location,omitempty"`
	Notes       string `json:"notes,omitempty"`
	PerformedBy string `json:"performed_by,omitempty"`
}

type PaginatedResponse struct {
	Data       []Asset `json:"data"`
	Total      int     `json:"total"`
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
	TotalPages int     `json:"total_pages"`
}

type DashboardStats struct {
	Total       int `json:"total"`
	Available   int `json:"available"`
	Assigned    int `json:"assigned"`
	UnderRepair int `json:"under_repair"`
	Retired     int `json:"retired"`
}
