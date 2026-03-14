package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	CategoryID      *uint          `gorm:"index" json:"category_id"`
	Title           string         `gorm:"size:255" json:"title"`
	Description     *string        `gorm:"type:text" json:"description"`
	Location        *string        `gorm:"size:255" json:"location"`
	StartDate       time.Time      `json:"start_date"`
	EndDate         time.Time      `json:"end_date"`
	MaxParticipants *int           `json:"max_participants"`
	Status          string         `gorm:"size:20;default:draft" json:"status"`
	CreatedBy       *uint          `gorm:"index" json:"created_by"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	Category *PrestationCategory `gorm:"foreignKey:CategoryID" json:"-"`
	Creator  *User               `gorm:"foreignKey:CreatedBy" json:"-"`
}

func (Event) TableName() string {
	return "platform_events"
}

type EventResponse struct {
	ID              uint              `json:"id"`
	Title           string            `json:"title"`
	Description     *string           `json:"description"`
	Location        *string           `json:"location"`
	StartDate       string            `json:"start_date"`
	EndDate         string            `json:"end_date"`
	MaxParticipants *int              `json:"max_participants"`
	Status          string            `json:"status"`
	Category        *CategoryResponse `json:"category"`
	Creator         *UserResponse     `json:"creator"`
	CreatedAt       string            `json:"created_at"`
	UpdatedAt       string            `json:"updated_at"`
}

func ToEventResponse(e *Event) EventResponse {
	var cat *CategoryResponse
	if e.Category != nil {
		c := ToCategoryResponse(e.Category, 0)
		cat = &c
	}
	var creator *UserResponse
	if e.Creator != nil {
		u := ToUserResponse(e.Creator)
		creator = &u
	}
	return EventResponse{
		ID:              e.ID,
		Title:           e.Title,
		Description:     e.Description,
		Location:        e.Location,
		StartDate:       e.StartDate.UTC().Format("2006-01-02T15:04:05.000000Z"),
		EndDate:         e.EndDate.UTC().Format("2006-01-02T15:04:05.000000Z"),
		MaxParticipants: e.MaxParticipants,
		Status:          e.Status,
		Category:        cat,
		Creator:         creator,
		CreatedAt:       e.CreatedAt.UTC().Format("2006-01-02T15:04:05.000000Z"),
		UpdatedAt:       e.UpdatedAt.UTC().Format("2006-01-02T15:04:05.000000Z"),
	}
}

func ToEventResponses(events []Event) []EventResponse {
	result := make([]EventResponse, len(events))
	for i := range events {
		result[i] = ToEventResponse(&events[i])
	}
	return result
}
