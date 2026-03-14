package services

import (
	"encoding/json"
	"time"

	"upcycleconnect/backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditService struct {
	DB *gorm.DB
}

func NewAuditService(db *gorm.DB) *AuditService {
	return &AuditService{DB: db}
}

func (s *AuditService) Log(c *gin.Context, action string, resourceType string, resourceID *uint, oldValues interface{}, newValues interface{}) {
	var userID *uint
	if val, exists := c.Get("user"); exists {
		if user, ok := val.(*models.User); ok {
			userID = &user.ID
		}
	}

	var oldJSON *string
	if oldValues != nil {
		b, err := json.Marshal(oldValues)
		if err == nil {
			s := string(b)
			oldJSON = &s
		}
	}

	var newJSON *string
	if newValues != nil {
		b, err := json.Marshal(newValues)
		if err == nil {
			s := string(b)
			newJSON = &s
		}
	}

	ip := c.ClientIP()

	log := models.AuditLog{
		UserID:       userID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		OldValues:    oldJSON,
		NewValues:    newJSON,
		IPAddress:    &ip,
		CreatedAt:    time.Now(),
	}

	s.DB.Create(&log)
}
