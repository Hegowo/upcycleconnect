package models

import (
	"time"

	"gorm.io/gorm"
)

type ForumThread struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	CategoryID uint           `gorm:"index;not null" json:"category_id"`
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	Title      string         `gorm:"size:255;not null" json:"title"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	Pinned     bool           `gorm:"default:false" json:"pinned"`
	Locked     bool           `gorm:"default:false" json:"locked"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	User     *User          `gorm:"foreignKey:UserID" json:"author,omitempty"`
	Category *ForumCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`

	ReplyCount int `gorm:"-" json:"reply_count"`
}

func (ForumThread) TableName() string { return "forum_threads" }
