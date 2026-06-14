package models

import (
	"time"

	"gorm.io/gorm"
)

type Tip struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:200;not null" json:"title"`
	Slug        string         `gorm:"size:220;uniqueIndex;not null" json:"slug"`
	Summary     string         `gorm:"size:500" json:"summary"`
	Content     string         `gorm:"type:longtext;not null" json:"content"`
	CoverImage  *string        `gorm:"size:255" json:"cover_image"`
	Category    string         `gorm:"size:60;index" json:"category"`
	AuthorID    uint           `gorm:"index;not null" json:"author_id"`
	Status      string         `gorm:"size:20;default:draft;index" json:"status"`
	PublishedAt *time.Time     `gorm:"index" json:"published_at"`
	ViewCount   uint           `gorm:"default:0" json:"view_count"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Author *User `gorm:"foreignKey:AuthorID" json:"-"`
}

func (Tip) TableName() string { return "tips" }

type TipResponse struct {
	ID          uint    `json:"id"`
	Title       string  `json:"title"`
	Slug        string  `json:"slug"`
	Summary     string  `json:"summary"`
	Content     string  `json:"content,omitempty"`
	CoverImage  *string `json:"cover_image"`
	Category    string  `json:"category"`
	Status      string  `json:"status"`
	PublishedAt *string `json:"published_at"`
	ViewCount   uint    `json:"view_count"`
	AuthorName  string  `json:"author_name"`
	CreatedAt   string  `json:"created_at"`
}

func ToTipResponse(t *Tip, includeContent bool) TipResponse {
	var pub *string
	if t.PublishedAt != nil {
		s := t.PublishedAt.UTC().Format(time.RFC3339)
		pub = &s
	}
	author := ""
	if t.Author != nil {
		author = t.Author.FirstName + " " + t.Author.LastName
	}
	r := TipResponse{
		ID:          t.ID,
		Title:       t.Title,
		Slug:        t.Slug,
		Summary:     t.Summary,
		CoverImage:  t.CoverImage,
		Category:    t.Category,
		Status:      t.Status,
		PublishedAt: pub,
		ViewCount:   t.ViewCount,
		AuthorName:  author,
		CreatedAt:   t.CreatedAt.UTC().Format(time.RFC3339),
	}
	if includeContent {
		r.Content = t.Content
	}
	return r
}
