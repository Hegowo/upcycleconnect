package models

import "time"

type PrestationImage struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	PrestationID uint      `gorm:"index;not null" json:"prestation_id"`
	URL          string    `gorm:"size:512;not null" json:"url"`
	Position     int       `gorm:"default:0" json:"order"`
	CreatedAt    time.Time `json:"created_at"`
}

func (PrestationImage) TableName() string {
	return "prestation_images"
}

type PrestationImageResponse struct {
	ID    uint   `json:"id"`
	URL   string `json:"url"`
	Order int    `json:"order"`
}
