package models

import "time"

type Locale struct {
	Code      string    `gorm:"primaryKey;size:10" json:"code"`
	Name      string    `gorm:"size:80;not null" json:"name"`
	Enabled   bool      `gorm:"default:true" json:"enabled"`
	Builtin   bool      `gorm:"default:false" json:"builtin"`
	Messages  string    `gorm:"type:longtext" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Locale) TableName() string { return "locales" }

type LocaleMeta struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Builtin bool   `json:"builtin"`
}

func ToLocaleMeta(l *Locale) LocaleMeta {
	return LocaleMeta{Code: l.Code, Name: l.Name, Enabled: l.Enabled, Builtin: l.Builtin}
}
