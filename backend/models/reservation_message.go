package models

import "time"

type ReservationMessage struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ReservationID uint      `gorm:"index;not null" json:"reservation_id"`
	SenderID      uint      `gorm:"index;not null" json:"sender_id"`
	Content       string    `gorm:"type:text" json:"content"`
	CreatedAt     time.Time `json:"created_at"`

	Sender *User `gorm:"foreignKey:SenderID" json:"-"`
}

func (ReservationMessage) TableName() string {
	return "reservation_messages"
}

type ReservationMessageResponse struct {
	ID         uint      `json:"id"`
	SenderID   uint      `json:"sender_id"`
	SenderName string    `json:"sender_name"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

func ToReservationMessageResponse(m *ReservationMessage) ReservationMessageResponse {
	name := ""
	if m.Sender != nil {
		name = m.Sender.FirstName + " " + m.Sender.LastName
	}
	return ReservationMessageResponse{
		ID:         m.ID,
		SenderID:   m.SenderID,
		SenderName: name,
		Content:    m.Content,
		CreatedAt:  m.CreatedAt,
	}
}
