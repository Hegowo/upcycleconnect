package database

import (
	"log"

	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"gorm.io/gorm"
)

func MigrateDepositPhotos(db *gorm.DB) {
	var deposits []models.DepositRequest
	if err := db.Where(
		"photo1 LIKE 'data:%' OR photo2 LIKE 'data:%' OR photo3 LIKE 'data:%'",
	).Find(&deposits).Error; err != nil {
		log.Printf("[migrate] deposit photos: query failed: %v", err)
		return
	}
	if len(deposits) == 0 {
		return
	}

	conv := func(p *string) (*string, bool) {
		if p == nil || *p == "" {
			return p, false
		}
		url, changed, err := services.SaveDataURIImage(*p, "/uploads/deposits", "/uploads/deposits")
		if err != nil {
			log.Printf("[migrate] deposit photo decode failed: %v", err)
			return p, false
		}
		if !changed {
			return p, false
		}
		return &url, true
	}

	migrated := 0
	for i := range deposits {
		d := &deposits[i]
		p1, c1 := conv(d.Photo1)
		p2, c2 := conv(d.Photo2)
		p3, c3 := conv(d.Photo3)
		if !c1 && !c2 && !c3 {
			continue
		}
		db.Model(&models.DepositRequest{}).Where("id = ?", d.ID).Updates(map[string]interface{}{
			"photo1": p1, "photo2": p2, "photo3": p3,
		})
		migrated++
	}
	log.Printf("[migrate] deposit photos: converted %d deposit(s) from base64 to files", migrated)
}
