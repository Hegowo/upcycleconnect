package models

import "time"

type Passkey struct {
	ID           uint       `gorm:"primarykey" json:"id"`
	UserID       uint       `gorm:"not null;index" json:"-"`
	CredentialID string     `gorm:"not null;uniqueIndex;size:512" json:"-"`
	Credential   []byte     `gorm:"not null;type:blob" json:"-"`
	Name         string     `gorm:"size:100" json:"name"`
	CreatedAt    time.Time  `json:"created_at"`
	LastUsedAt   *time.Time `gorm:"default:null" json:"last_used_at"`
}
