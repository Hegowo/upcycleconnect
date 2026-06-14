package models

import "time"

type ForumReport struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Type       string    `gorm:"size:20;not null" json:"type"`
	TargetID   uint      `gorm:"index;not null" json:"target_id"`
	ReporterID uint      `gorm:"index;not null" json:"reporter_id"`
	Reason     string    `gorm:"type:text" json:"reason"`
	Status     string    `gorm:"size:20;default:pending" json:"status"`
	CreatedAt  time.Time `json:"created_at"`

	Reporter *User `gorm:"foreignKey:ReporterID" json:"reporter,omitempty"`
}

func (ForumReport) TableName() string { return "forum_reports" }
