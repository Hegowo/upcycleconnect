package models

import (
	"time"

	"gorm.io/gorm"
)

type Campaign struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	ProviderID      uint           `gorm:"index;not null" json:"provider_id"`
	Title           string         `gorm:"size:200;not null" json:"title"`
	Description     string         `gorm:"type:text" json:"description"`
	ImageURL        *string        `gorm:"size:500" json:"image_url"`
	BudgetCents     int64          `gorm:"not null" json:"budget_cents"`
	Status          string         `gorm:"size:20;default:draft;index" json:"status"`
	RejectionReason *string        `gorm:"type:text" json:"rejection_reason"`
	StartDate       *time.Time     `json:"start_date"`
	EndDate         *time.Time     `json:"end_date"`
	StripeSessionID *string        `gorm:"size:200" json:"stripe_session_id"`
	PaidAt          *time.Time     `json:"paid_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	Provider *User          `gorm:"foreignKey:ProviderID" json:"-"`
	Items    []CampaignItem `gorm:"foreignKey:CampaignID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

func (Campaign) TableName() string { return "campaigns" }

type CampaignItem struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	CampaignID uint   `gorm:"index;not null" json:"campaign_id"`
	ItemType   string `gorm:"size:20;not null" json:"item_type"`
	ItemID     uint   `gorm:"not null" json:"item_id"`
}

func (CampaignItem) TableName() string { return "campaign_items" }

type CampaignItemRef struct {
	Type string `json:"type"`
	ID   uint   `json:"id"`
}

type CampaignResponse struct {
	ID              uint              `json:"id"`
	Title           string            `json:"title"`
	Description     string            `json:"description"`
	ImageURL        *string           `json:"image_url"`
	BudgetCents     int64             `json:"budget_cents"`
	Status          string            `json:"status"`
	RejectionReason *string           `json:"rejection_reason"`
	StartDate       *string           `json:"start_date"`
	EndDate         *string           `json:"end_date"`
	PaidAt          *string           `json:"paid_at"`
	ProviderName    string            `json:"provider_name"`
	Items           []CampaignItemRef `json:"items"`
	CreatedAt       string            `json:"created_at"`
}

func ToCampaignResponse(c *Campaign) CampaignResponse {
	fmtDate := func(t *time.Time) *string {
		if t == nil {
			return nil
		}
		s := t.UTC().Format(time.RFC3339)
		return &s
	}
	name := ""
	if c.Provider != nil {
		name = c.Provider.FirstName + " " + c.Provider.LastName
	}
	items := make([]CampaignItemRef, 0, len(c.Items))
	for _, it := range c.Items {
		items = append(items, CampaignItemRef{Type: it.ItemType, ID: it.ItemID})
	}
	return CampaignResponse{
		ID:              c.ID,
		Title:           c.Title,
		Description:     c.Description,
		ImageURL:        c.ImageURL,
		BudgetCents:     c.BudgetCents,
		Status:          c.Status,
		RejectionReason: c.RejectionReason,
		StartDate:       fmtDate(c.StartDate),
		EndDate:         fmtDate(c.EndDate),
		PaidAt:          fmtDate(c.PaidAt),
		ProviderName:    name,
		Items:           items,
		CreatedAt:       c.CreatedAt.UTC().Format(time.RFC3339),
	}
}
