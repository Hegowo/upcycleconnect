package models

import "time"

type EmailVerification struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"index;not null" json:"user_id"`
	Token     string     `gorm:"uniqueIndex;size:64;not null" json:"token"`
	Type      string     `gorm:"size:20;not null" json:"type"`
	IPAddress *string    `gorm:"size:45" json:"ip_address"`
	UsedAt    *time.Time `json:"used_at"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
}

func (EmailVerification) TableName() string {
	return "email_verifications"
}

type UserKnownIP struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex:user_ip;not null" json:"user_id"`
	IPAddress string    `gorm:"uniqueIndex:user_ip;size:45;not null" json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}

func (UserKnownIP) TableName() string {
	return "user_known_ips"
}
