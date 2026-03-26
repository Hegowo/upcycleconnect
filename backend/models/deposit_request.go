package models

import "time"

type DepositRequest struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	UserID          uint       `gorm:"index;not null" json:"user_id"`
	CategoryID      *uint      `gorm:"index" json:"category_id"`
	Title           string     `gorm:"size:255;not null" json:"title"`
	Description     string     `gorm:"type:text;not null" json:"description"`
	Condition       string     `gorm:"size:20;default:good" json:"condition"`
	History         *string    `gorm:"size:255" json:"history"`
	EstimatedWeight *float64   `gorm:"type:decimal(6,2)" json:"estimated_weight"`
	CarbonSavings   *float64   `gorm:"type:decimal(8,2)" json:"carbon_savings"`
	Status          string     `gorm:"size:20;default:pending" json:"status"`
	ValidationNote  *string    `gorm:"type:text" json:"validation_note"`
	ValidatedBy     *uint      `gorm:"index" json:"validated_by"`
	ValidatedAt     *time.Time `json:"validated_at"`
	QRCode          *string    `gorm:"size:100" json:"qr_code"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	User     *User               `gorm:"foreignKey:UserID" json:"-"`
	Category *PrestationCategory `gorm:"foreignKey:CategoryID" json:"-"`
}

func (DepositRequest) TableName() string {
	return "deposit_requests"
}

type DepositRequestResponse struct {
	ID              uint              `json:"id"`
	Title           string            `json:"title"`
	Description     string            `json:"description"`
	Condition       string            `json:"condition"`
	History         *string           `json:"history"`
	EstimatedWeight *float64          `json:"estimated_weight"`
	CarbonSavings   *float64          `json:"carbon_savings"`
	Status          string            `json:"status"`
	ValidationNote  *string           `json:"validation_note"`
	QRCode          *string           `json:"qr_code"`
	User            *UserResponse     `json:"user"`
	Category        *CategoryResponse `json:"category"`
	CreatedAt       string            `json:"created_at"`
	ValidatedAt     *string           `json:"validated_at"`
}

func ToDepositResponse(d *DepositRequest) DepositRequestResponse {
	var user *UserResponse
	if d.User != nil {
		u := ToUserResponse(d.User)
		user = &u
	}
	var cat *CategoryResponse
	if d.Category != nil {
		c := ToCategoryResponse(d.Category, 0)
		cat = &c
	}
	var validatedAt *string
	if d.ValidatedAt != nil {
		s := d.ValidatedAt.UTC().Format("2006-01-02T15:04:05.000000Z")
		validatedAt = &s
	}
	return DepositRequestResponse{
		ID:              d.ID,
		Title:           d.Title,
		Description:     d.Description,
		Condition:       d.Condition,
		History:         d.History,
		EstimatedWeight: d.EstimatedWeight,
		CarbonSavings:   d.CarbonSavings,
		Status:          d.Status,
		ValidationNote:  d.ValidationNote,
		QRCode:          d.QRCode,
		User:            user,
		Category:        cat,
		CreatedAt:       d.CreatedAt.UTC().Format("2006-01-02T15:04:05.000000Z"),
		ValidatedAt:     validatedAt,
	}
}

func ToDepositResponses(items []DepositRequest) []DepositRequestResponse {
	result := make([]DepositRequestResponse, len(items))
	for i := range items {
		result[i] = ToDepositResponse(&items[i])
	}
	return result
}
