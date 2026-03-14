package models

import (
	"time"

	"gorm.io/gorm"
)

type Prestation struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CategoryID  *uint          `gorm:"index" json:"category_id"`
	ProviderID  *uint          `gorm:"index" json:"provider_id"`
	Title       string         `gorm:"size:255" json:"title"`
	Description *string        `gorm:"type:text" json:"description"`
	Price       *string        `gorm:"type:decimal(10,2)" json:"price"`
	PriceType   string         `gorm:"size:20;default:fixed" json:"price_type"`
	Status      string         `gorm:"size:20;default:draft" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Category *PrestationCategory `gorm:"foreignKey:CategoryID" json:"-"`
	Provider *User               `gorm:"foreignKey:ProviderID" json:"-"`
}

func (Prestation) TableName() string {
	return "prestations"
}

type PrestationResponse struct {
	ID          uint               `json:"id"`
	Title       string             `json:"title"`
	Description *string            `json:"description"`
	Price       *string            `json:"price"`
	PriceType   string             `json:"price_type"`
	Status      string             `json:"status"`
	Category    *CategoryResponse  `json:"category"`
	Provider    *UserResponse      `json:"provider"`
	CreatedAt   string             `json:"created_at"`
	UpdatedAt   string             `json:"updated_at"`
}

func ToPrestationResponse(p *Prestation) PrestationResponse {
	var cat *CategoryResponse
	if p.Category != nil {
		c := ToCategoryResponse(p.Category, 0)
		cat = &c
	}
	var prov *UserResponse
	if p.Provider != nil {
		u := ToUserResponse(p.Provider)
		prov = &u
	}
	return PrestationResponse{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		Price:       p.Price,
		PriceType:   p.PriceType,
		Status:      p.Status,
		Category:    cat,
		Provider:    prov,
		CreatedAt:   p.CreatedAt.UTC().Format("2006-01-02T15:04:05.000000Z"),
		UpdatedAt:   p.UpdatedAt.UTC().Format("2006-01-02T15:04:05.000000Z"),
	}
}

func ToPrestationResponses(prestations []Prestation) []PrestationResponse {
	result := make([]PrestationResponse, len(prestations))
	for i := range prestations {
		result[i] = ToPrestationResponse(&prestations[i])
	}
	return result
}
