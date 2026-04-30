package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID              uint            `gorm:"primaryKey" json:"id"`
	Email           string          `gorm:"uniqueIndex;size:255" json:"email"`
	Password        string          `gorm:"size:255" json:"-"`
	FirstName       string          `gorm:"size:100" json:"first_name"`
	LastName        string          `gorm:"size:100" json:"last_name"`
	Phone           *string         `gorm:"size:20" json:"phone"`
	Address         *string         `gorm:"size:512" json:"address"`
	AvatarURL       *string         `gorm:"size:255" json:"avatar_url"`
	Status          string          `gorm:"size:20;default:active" json:"status"`
	EmailVerifiedAt *time.Time      `json:"email_verified_at"`
	RememberToken   *string         `gorm:"size:100" json:"-"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	DeletedAt       gorm.DeletedAt  `gorm:"index" json:"-"`
	Roles           []Role          `gorm:"many2many:user_roles;" json:"-"`
	ProviderProfile *ProviderProfile `gorm:"foreignKey:UserID" json:"-"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) PrimaryRole() *string {
	if len(u.Roles) == 0 {
		return nil
	}
	for _, r := range u.Roles {
		if r.Name == "super_admin" {
			name := r.Name
			return &name
		}
	}
	for _, r := range u.Roles {
		if r.Name == "admin" {
			name := r.Name
			return &name
		}
	}
	return nil
}

func (u *User) IsAdmin() bool {
	for _, r := range u.Roles {
		if r.Name == "admin" || r.Name == "super_admin" {
			return true
		}
	}
	return false
}

func (u *User) IsSuperAdmin() bool {
	for _, r := range u.Roles {
		if r.Name == "super_admin" {
			return true
		}
	}
	return false
}

type UserResponse struct {
	ID              uint    `json:"id"`
	Email           string  `json:"email"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	Phone           *string `json:"phone"`
	Address         *string `json:"address"`
	AvatarURL       *string `json:"avatar_url"`
	Status          string  `json:"status"`
	Role            *string `json:"role"`
	EmailVerifiedAt *string `json:"email_verified_at"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

func ToUserResponse(u *User) UserResponse {
	var emailVerified *string
	if u.EmailVerifiedAt != nil {
		s := u.EmailVerifiedAt.UTC().Format(time.RFC3339Nano)
		emailVerified = &s
	}
	return UserResponse{
		ID:              u.ID,
		Email:           u.Email,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		Phone:           u.Phone,
		Address:         u.Address,
		AvatarURL:       u.AvatarURL,
		Status:          u.Status,
		Role:            u.PrimaryRole(),
		EmailVerifiedAt: emailVerified,
		CreatedAt:       u.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:       u.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}

func ToUserResponses(users []User) []UserResponse {
	result := make([]UserResponse, len(users))
	for i := range users {
		result[i] = ToUserResponse(&users[i])
	}
	return result
}
