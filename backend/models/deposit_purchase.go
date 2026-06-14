package models

import "time"

// DepositPurchase records a professional buying (vente) or claiming (don) a
// deposit on the marketplace. For "vente" it is created paid after Stripe;
// for "don" it is created paid immediately (amount 0).
type DepositPurchase struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	DepositID       uint       `gorm:"index;not null" json:"deposit_id"`
	ProviderID      uint       `gorm:"index;not null" json:"provider_id"`
	AmountCents     int64      `gorm:"not null;default:0" json:"amount_cents"`
	Currency        string     `gorm:"size:3;default:eur" json:"currency"`
	Status          string     `gorm:"size:20;default:pending" json:"status"` // pending | paid
	StripeSessionID *string    `gorm:"size:255;index" json:"-"`
	PaidAt          *time.Time `json:"paid_at"`
	CreatedAt       time.Time  `json:"created_at"`

	Deposit *DepositRequest `gorm:"foreignKey:DepositID" json:"deposit,omitempty"`
}

func (DepositPurchase) TableName() string { return "deposit_purchases" }

// DepositFavorite is a marketplace bookmark by a professional.
type DepositFavorite struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	DepositID  uint      `gorm:"uniqueIndex:idx_dep_fav;not null" json:"deposit_id"`
	ProviderID uint      `gorm:"uniqueIndex:idx_dep_fav;not null" json:"provider_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (DepositFavorite) TableName() string { return "deposit_favorites" }
