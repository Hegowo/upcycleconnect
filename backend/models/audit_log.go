package models

import "time"

type AuditLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       *uint     `gorm:"index" json:"user_id"`
	Action       string    `gorm:"size:100" json:"action"`
	ResourceType string    `gorm:"size:100" json:"resource_type"`
	ResourceID   *uint     `json:"resource_id"`
	OldValues    *string   `gorm:"type:json" json:"old_values"`
	NewValues    *string   `gorm:"type:json" json:"new_values"`
	IPAddress    *string   `gorm:"size:45" json:"ip_address"`
	CreatedAt    time.Time `json:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
