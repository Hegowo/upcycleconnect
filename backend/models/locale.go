package models

import "time"

// Locale is a translation language managed at runtime (no code change required, per spec p.7).
// Messages holds the full i18n message tree as a JSON string. Admins can add a new language
// and fill its translations from the back-office; the frontend loads enabled locales dynamically.
type Locale struct {
	Code      string    `gorm:"primaryKey;size:10" json:"code"` // e.g. "es", "de", "pt-BR"
	Name      string    `gorm:"size:80;not null" json:"name"`   // e.g. "Español"
	Enabled   bool      `gorm:"default:true" json:"enabled"`
	Builtin   bool      `gorm:"default:false" json:"builtin"` // fr/en shipped with the app
	Messages  string    `gorm:"type:longtext" json:"-"`       // JSON message tree (omitted from list responses)
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
