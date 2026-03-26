package models

import "time"

type EmailChangeRequest struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Step      string    `gorm:"size:20;not null" json:"step"` // verify_current | verify_new
	Code      string    `gorm:"size:6;not null" json:"code"`
	NewEmail  *string   `gorm:"size:255" json:"new_email"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (EmailChangeRequest) TableName() string {
	return "email_change_requests"
}
