package models

import "time"

type ForumBan struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ThreadID  uint      `gorm:"uniqueIndex:idx_thread_user;index;not null" json:"thread_id"`
	UserID    uint      `gorm:"uniqueIndex:idx_thread_user;index;not null" json:"user_id"`
	BannedBy  uint      `gorm:"not null" json:"banned_by"`
	CreatedAt time.Time `json:"created_at"`

	User         *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	BannedByUser *User `gorm:"foreignKey:BannedBy" json:"banned_by_user,omitempty"`
}

func (ForumBan) TableName() string { return "forum_bans" }
