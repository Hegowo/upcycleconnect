package models

import "time"

type ForumCategory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Slug        string    `gorm:"size:100;uniqueIndex;not null" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	Icon        string    `gorm:"size:50" json:"icon"`
	Color       string    `gorm:"size:20" json:"color"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	ThreadCount int `gorm:"-" json:"thread_count"`
	ReplyCount  int `gorm:"-" json:"reply_count"`
}

func (ForumCategory) TableName() string { return "forum_categories" }
