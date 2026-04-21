package models

import (
	"time"

	"gorm.io/gorm"
)

type Invoice struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	UserID          uint           `gorm:"index;not null" json:"user_id"`
	ReservationID   *uint          `gorm:"index" json:"reservation_id"`
	Type            string         `gorm:"size:20;default:invoice" json:"type"`
	Number          string         `gorm:"size:50;uniqueIndex;not null" json:"number"`
	PrestationTitle string         `gorm:"size:255;not null" json:"prestation_title"`
	AmountCents     int64          `gorm:"not null;default:0" json:"amount_cents"`
	TVAPercent      float64        `gorm:"type:decimal(5,2);default:20.00" json:"tva_percent"`
	Currency        string         `gorm:"size:3;default:eur" json:"currency"`
	PDFPath         *string        `gorm:"size:500" json:"pdf_path"`
	Status          string         `gorm:"size:20;default:issued" json:"status"`
	IssuedAt        time.Time      `gorm:"not null" json:"issued_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	User        *User        `gorm:"foreignKey:UserID" json:"-"`
	Reservation *Reservation `gorm:"foreignKey:ReservationID" json:"-"`
}

func (Invoice) TableName() string {
	return "invoices"
}

type InvoiceResponse struct {
	ID              uint    `json:"id"`
	Type            string  `json:"type"`
	Number          string  `json:"number"`
	PrestationTitle string  `json:"prestation_title"`
	AmountCents     int64   `json:"amount_cents"`
	TVAPercent      float64 `json:"tva_percent"`
	Currency        string  `json:"currency"`
	Status          string  `json:"status"`
	HasPDF          bool    `json:"has_pdf"`
	IssuedAt        string  `json:"issued_at"`
	CreatedAt       string  `json:"created_at"`
}

func ToInvoiceResponse(i *Invoice) InvoiceResponse {
	return InvoiceResponse{
		ID:              i.ID,
		Type:            i.Type,
		Number:          i.Number,
		PrestationTitle: i.PrestationTitle,
		AmountCents:     i.AmountCents,
		TVAPercent:      i.TVAPercent,
		Currency:        i.Currency,
		Status:          i.Status,
		HasPDF:          i.PDFPath != nil && *i.PDFPath != "",
		IssuedAt:        i.IssuedAt.UTC().Format("2006-01-02T15:04:05.000000Z"),
		CreatedAt:       i.CreatedAt.UTC().Format("2006-01-02T15:04:05.000000Z"),
	}
}

func ToInvoiceResponses(items []Invoice) []InvoiceResponse {
	result := make([]InvoiceResponse, len(items))
	for i := range items {
		result[i] = ToInvoiceResponse(&items[i])
	}
	return result
}
