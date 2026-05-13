package models

import "time"

type CollectionPoint struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Name         string     `gorm:"size:255;not null" json:"name"`
	Address      string     `gorm:"size:512;not null" json:"address"`
	City         string     `gorm:"size:100;not null" json:"city"`
	PostalCode   string     `gorm:"size:10;not null" json:"postal_code"`
	Latitude     float64    `gorm:"type:decimal(10,7);not null" json:"latitude"`
	Longitude    float64    `gorm:"type:decimal(10,7);not null" json:"longitude"`
	Phone        *string    `gorm:"size:30" json:"phone"`
	OpeningHours *string    `gorm:"size:512" json:"opening_hours"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (CollectionPoint) TableName() string {
	return "collection_points"
}

type CollectionPointResponse struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Address      string  `json:"address"`
	City         string  `json:"city"`
	PostalCode   string  `json:"postal_code"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Phone        *string `json:"phone"`
	OpeningHours *string `json:"opening_hours"`
	IsActive     bool    `json:"is_active"`
	CreatedAt    string  `json:"created_at"`
}

func ToCollectionPointResponse(cp *CollectionPoint) CollectionPointResponse {
	return CollectionPointResponse{
		ID:           cp.ID,
		Name:         cp.Name,
		Address:      cp.Address,
		City:         cp.City,
		PostalCode:   cp.PostalCode,
		Latitude:     cp.Latitude,
		Longitude:    cp.Longitude,
		Phone:        cp.Phone,
		OpeningHours: cp.OpeningHours,
		IsActive:     cp.IsActive,
		CreatedAt:    cp.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}
}

func ToCollectionPointResponses(items []CollectionPoint) []CollectionPointResponse {
	result := make([]CollectionPointResponse, len(items))
	for i := range items {
		result[i] = ToCollectionPointResponse(&items[i])
	}
	return result
}
