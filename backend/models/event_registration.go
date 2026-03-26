package models

import "time"

type EventRegistration struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex:user_event;index;not null" json:"user_id"`
	EventID   uint      `gorm:"uniqueIndex:user_event;index;not null" json:"event_id"`
	CreatedAt time.Time `json:"created_at"`

	User  *User  `gorm:"foreignKey:UserID" json:"-"`
	Event *Event `gorm:"foreignKey:EventID" json:"-"`
}

func (EventRegistration) TableName() string {
	return "event_registrations"
}
