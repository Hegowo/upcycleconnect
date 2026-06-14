package database

import (
	"log"
	"os"
	"time"

	"upcycleconnect/backend/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	seedRoles(db)
	seedLocales(db)
	seedDefaultAdmin(db)
}

func seedRoles(db *gorm.DB) {
	for _, name := range []string{"super_admin", "admin", "user"} {
		var count int64
		db.Model(&models.Role{}).Where("name = ?", name).Count(&count)
		if count == 0 {
			db.Create(&models.Role{Name: name})
		}
	}
}

func seedLocales(db *gorm.DB) {
	builtins := []models.Locale{
		{Code: "fr", Name: "Français", Enabled: true, Builtin: true, Messages: "{}"},
		{Code: "en", Name: "English", Enabled: true, Builtin: true, Messages: "{}"},
	}
	for _, l := range builtins {
		var count int64
		db.Model(&models.Locale{}).Where("code = ?", l.Code).Count(&count)
		if count == 0 {
			db.Create(&l)
		}
	}
}

func seedDefaultAdmin(db *gorm.DB) {
	email := os.Getenv("SEED_ADMIN_EMAIL")
	password := os.Getenv("SEED_ADMIN_PASSWORD")
	if email == "" || password == "" {
		return
	}

	var superCount int64
	db.Table("user_roles").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("roles.name = ?", "super_admin").
		Count(&superCount)
	if superCount > 0 {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[seed] cannot hash default admin password: %v", err)
		return
	}

	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		now := time.Now()
		user = models.User{
			Email:           email,
			Password:        string(hash),
			FirstName:       "Super",
			LastName:        "Admin",
			Status:          "active",
			EmailVerifiedAt: &now,
		}
		if err := db.Create(&user).Error; err != nil {
			log.Printf("[seed] cannot create default admin: %v", err)
			return
		}
	}

	var superRole models.Role
	if db.Where("name = ?", "super_admin").First(&superRole).Error == nil {
		db.Exec("INSERT IGNORE INTO user_roles (user_id, role_id) VALUES (?, ?)", user.ID, superRole.ID)
		log.Printf("[seed] default super-admin ensured for %s", email)
	}
}
