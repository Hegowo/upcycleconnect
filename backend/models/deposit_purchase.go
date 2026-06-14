package models

import "time"

type DepositPurchase struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	DepositID       uint       `gorm:"index;not null" json:"deposit_id"`
	ProviderID      uint       `gorm:"index;not null" json:"provider_id"`
	AmountCents     int64      `gorm:"not null;default:0" json:"amount_cents"`
	Currency        string     `gorm:"size:3;default:eur" json:"currency"`
	Status          string     `gorm:"size:20;default:pending" json:"status"`
	StripeSessionID *string    `gorm:"size:255;index" json:"-"`
	PaidAt          *time.Time `json:"paid_at"`
	CreatedAt       time.Time  `json:"created_at"`

	Deposit *DepositRequest `gorm:"foreignKey:DepositID" json:"deposit,omitempty"`
}

func (DepositPurchase) TableName() string { return "deposit_purchases" }

type DepositFavorite struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	DepositID  uint      `gorm:"uniqueIndex:idx_dep_fav;not null" json:"deposit_id"`
	ProviderID uint      `gorm:"uniqueIndex:idx_dep_fav;not null" json:"provider_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (DepositFavorite) TableName() string { return "deposit_favorites" }
