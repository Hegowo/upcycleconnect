package models

import "time"

type EventMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	EventID   uint      `gorm:"index;not null" json:"event_id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Content   string    `gorm:"type:text" json:"content"`
	ImageURL  *string   `gorm:"size:512" json:"image_url"`
	CreatedAt time.Time `json:"created_at"`

	User  *User  `gorm:"foreignKey:UserID" json:"-"`
	Event *Event `gorm:"foreignKey:EventID" json:"-"`
}

func (EventMessage) TableName() string {
	return "event_messages"
}
