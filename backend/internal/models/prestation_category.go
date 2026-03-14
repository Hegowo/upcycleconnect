package models

import (
	"time"

	"gorm.io/gorm"
)

type PrestationCategory struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:255" json:"name"`
	Slug        string         `gorm:"size:255;uniqueIndex" json:"slug"`
	Description *string        `gorm:"type:text" json:"description"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Prestations []Prestation `gorm:"foreignKey:CategoryID" json:"-"`
}

func (PrestationCategory) TableName() string {
	return "prestation_categories"
}

type CategoryResponse struct {
	ID               uint    `json:"id"`
	Name             string  `json:"name"`
	Slug             string  `json:"slug"`
	Description      *string `json:"description"`
	IsActive         bool    `json:"is_active"`
	SortOrder        int     `json:"sort_order"`
	PrestationsCount int64   `json:"prestations_count"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

func ToCategoryResponse(c *PrestationCategory, count int64) CategoryResponse {
	return CategoryResponse{
		ID:               c.ID,
		Name:             c.Name,
		Slug:             c.Slug,
		Description:      c.Description,
		IsActive:         c.IsActive,
		SortOrder:        c.SortOrder,
		PrestationsCount: count,
		CreatedAt:        c.CreatedAt.UTC().Format("2006-01-02T15:04:05.000000Z"),
		UpdatedAt:        c.UpdatedAt.UTC().Format("2006-01-02T15:04:05.000000Z"),
	}
}
