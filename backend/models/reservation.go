package models

import (
	"time"

	"gorm.io/gorm"
)

type Reservation struct {
	ID                    uint           `gorm:"primaryKey" json:"id"`
	UserID                uint           `gorm:"index;not null" json:"user_id"`
	PrestationID          uint           `gorm:"index;not null" json:"prestation_id"`
	Status                string         `gorm:"size:20;default:pending" json:"status"`
	AmountCents           int64          `gorm:"not null;default:0" json:"amount_cents"`
	Currency              string         `gorm:"size:3;default:eur" json:"currency"`
	StripeSessionID       *string        `gorm:"size:255;index" json:"stripe_session_id"`
	StripePaymentIntentID *string        `gorm:"size:255;index" json:"stripe_payment_intent_id"`
	Notes                 *string        `gorm:"type:text" json:"notes"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`

	User       *User       `gorm:"foreignKey:UserID" json:"-"`
	Prestation *Prestation `gorm:"foreignKey:PrestationID" json:"-"`
}

func (Reservation) TableName() string {
	return "reservations"
}

type ReservationResponse struct {
	ID           uint                `json:"id"`
	Status       string              `json:"status"`
	AmountCents  int64               `json:"amount_cents"`
	Currency     string              `json:"currency"`
	Notes        *string             `json:"notes"`
	Prestation   *PrestationResponse `json:"prestation"`
	CreatedAt    string              `json:"created_at"`
	UpdatedAt    string              `json:"updated_at"`
}

func ToReservationResponse(r *Reservation) ReservationResponse {
	var prest *PrestationResponse
	if r.Prestation != nil {
		p := ToPrestationResponse(r.Prestation)
		prest = &p
	}
	return ReservationResponse{
		ID:          r.ID,
		Status:      r.Status,
		AmountCents: r.AmountCents,
		Currency:    r.Currency,
		Notes:       r.Notes,
		Prestation:  prest,
		CreatedAt:   r.CreatedAt.UTC().Format("2006-01-02T15:04:05.000000Z"),
		UpdatedAt:   r.UpdatedAt.UTC().Format("2006-01-02T15:04:05.000000Z"),
	}
}

func ToReservationResponses(items []Reservation) []ReservationResponse {
	result := make([]ReservationResponse, len(items))
	for i := range items {
		result[i] = ToReservationResponse(&items[i])
	}
	return result
}
