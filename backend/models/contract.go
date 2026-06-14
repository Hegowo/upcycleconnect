package models

import (
	"time"

	"gorm.io/gorm"
)

type Contract struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	Number        string `gorm:"size:50;uniqueIndex;not null" json:"number"`
	UserID        uint   `gorm:"index;not null" json:"user_id"`
	ReservationID uint   `gorm:"uniqueIndex;not null" json:"reservation_id"`

	Signature []byte    `gorm:"type:mediumblob;not null" json:"-"`
	SignedAt  time.Time `gorm:"not null" json:"signed_at"`
	SignedIP  *string   `gorm:"size:45" json:"signed_ip"`

	CustomerName    string  `gorm:"size:200;not null" json:"customer_name"`
	CustomerEmail   string  `gorm:"size:255;not null" json:"customer_email"`
	CustomerAddress *string `gorm:"size:512" json:"customer_address"`
	CustomerPhone   *string `gorm:"size:50" json:"customer_phone"`
	PrestationTitle string  `gorm:"size:255;not null" json:"prestation_title"`
	PrestationDesc  *string `gorm:"type:text" json:"prestation_description"`
	ProviderName    string  `gorm:"size:200" json:"provider_name"`
	ProviderEmail   string  `gorm:"size:255" json:"provider_email"`
	AmountCents     int64   `gorm:"not null;default:0" json:"amount_cents"`
	Currency        string  `gorm:"size:3;default:eur" json:"currency"`

	Status  string  `gorm:"size:30;default:signed_pending_payment" json:"status"`
	PDFPath *string `gorm:"size:500" json:"-"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User        *User        `gorm:"foreignKey:UserID" json:"-"`
	Reservation *Reservation `gorm:"foreignKey:ReservationID" json:"-"`
}

func (Contract) TableName() string {
	return "contracts"
}

type ContractResponse struct {
	ID              uint   `json:"id"`
	Number          string `json:"number"`
	Status          string `json:"status"`
	CustomerName    string `json:"customer_name"`
	CustomerEmail   string `json:"customer_email"`
	PrestationTitle string `json:"prestation_title"`
	AmountCents     int64  `json:"amount_cents"`
	Currency        string `json:"currency"`
	SignedAt        string `json:"signed_at"`
	HasPDF          bool   `json:"has_pdf"`
	ReservationID   uint   `json:"reservation_id"`
}

func ToContractResponse(c *Contract) ContractResponse {
	return ContractResponse{
		ID:              c.ID,
		Number:          c.Number,
		Status:          c.Status,
		CustomerName:    c.CustomerName,
		CustomerEmail:   c.CustomerEmail,
		PrestationTitle: c.PrestationTitle,
		AmountCents:     c.AmountCents,
		Currency:        c.Currency,
		SignedAt:        c.SignedAt.UTC().Format(time.RFC3339),
		HasPDF:          c.PDFPath != nil && *c.PDFPath != "",
		ReservationID:   c.ReservationID,
	}
}
