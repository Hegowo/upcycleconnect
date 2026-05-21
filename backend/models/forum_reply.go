package models

import (
	"time"

	"gorm.io/gorm"
)

type ForumReply struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ThreadID  uint           `gorm:"index;not null" json:"thread_id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User *User `gorm:"foreignKey:UserID" json:"author,omitempty"`
}

func (ForumReply) TableName() string { return "forum_replies" }
