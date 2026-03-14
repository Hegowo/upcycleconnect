package models

import "time"

type ProviderProfile struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	UserID      uint       `gorm:"index" json:"user_id"`
	CompanyName string     `gorm:"size:255" json:"company_name"`
	Siret       *string    `gorm:"size:20" json:"siret"`
	Description *string    `gorm:"type:text" json:"description"`
	Website     *string    `gorm:"size:255" json:"website"`
	Status      string     `gorm:"size:20;default:pending" json:"status"`
	ApprovedAt  *time.Time `json:"approved_at"`
	ApprovedBy  *uint      `json:"approved_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (ProviderProfile) TableName() string {
	return "provider_profiles"
}

type ProviderProfileResponse struct {
	ID          uint    `json:"id"`
	CompanyName string  `json:"company_name"`
	Siret       *string `json:"siret"`
	Description *string `json:"description"`
	Website     *string `json:"website"`
	Status      string  `json:"status"`
	ApprovedAt  *string `json:"approved_at"`
}

func ToProviderProfileResponse(p *ProviderProfile) *ProviderProfileResponse {
	if p == nil {
		return nil
	}
	var approvedAt *string
	if p.ApprovedAt != nil {
		s := p.ApprovedAt.UTC().Format("2006-01-02T15:04:05.000000Z")
		approvedAt = &s
	}
	return &ProviderProfileResponse{
		ID:          p.ID,
		CompanyName: p.CompanyName,
		Siret:       p.Siret,
		Description: p.Description,
		Website:     p.Website,
		Status:      p.Status,
		ApprovedAt:  approvedAt,
	}
}

type ProviderResponse struct {
	ID        uint                     `json:"id"`
	Email     string                   `json:"email"`
	FirstName string                   `json:"first_name"`
	LastName  string                   `json:"last_name"`
	Phone     *string                  `json:"phone"`
	Status    string                   `json:"status"`
	CreatedAt string                   `json:"created_at"`
	Profile   *ProviderProfileResponse `json:"profile"`
}

func ToProviderResponse(u *User) ProviderResponse {
	return ProviderResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Phone:     u.Phone,
		Status:    u.Status,
		CreatedAt: u.CreatedAt.UTC().Format("2006-01-02T15:04:05.000000Z"),
		Profile:   ToProviderProfileResponse(u.ProviderProfile),
	}
}

func ToProviderResponses(users []User) []ProviderResponse {
	result := make([]ProviderResponse, len(users))
	for i := range users {
		result[i] = ToProviderResponse(&users[i])
	}
	return result
}
