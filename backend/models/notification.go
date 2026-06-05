package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	Type      string         `gorm:"size:50;not null;index" json:"type"`
	Title     string         `gorm:"size:200;not null" json:"title"`
	Body      string         `gorm:"type:text" json:"body"`
	Link      *string        `gorm:"size:500" json:"link"`
	ReadAt    *time.Time     `gorm:"index" json:"read_at"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Notification) TableName() string {
	return "notifications"
}

type NotificationResponse struct {
	ID        uint    `json:"id"`
	Type      string  `json:"type"`
	Title     string  `json:"title"`
	Body      string  `json:"body"`
	Link      *string `json:"link"`
	Read      bool    `json:"read"`
	CreatedAt string  `json:"created_at"`
}

func ToNotificationResponse(n *Notification) NotificationResponse {
	return NotificationResponse{
		ID:        n.ID,
		Type:      n.Type,
		Title:     n.Title,
		Body:      n.Body,
		Link:      n.Link,
		Read:      n.ReadAt != nil,
		CreatedAt: n.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func ToNotificationResponses(items []Notification) []NotificationResponse {
	out := make([]NotificationResponse, len(items))
	for i := range items {
		out[i] = ToNotificationResponse(&items[i])
	}
	return out
}
