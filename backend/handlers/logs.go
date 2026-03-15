package handlers

import (
	"net/http"
	"time"

	"upcycleconnect/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LogsHandler struct {
	DB *gorm.DB
}

type auditLogRow struct {
	ID           uint      `json:"id"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resource_type"`
	ResourceID   *uint     `json:"resource_id"`
	AdminName    string    `json:"admin_name"`
	IPAddress    *string   `json:"ip_address"`
	CreatedAt    time.Time `json:"created_at"`
}

func (h *LogsHandler) Index(c *gin.Context) {
	var rows []struct {
		models.AuditLog
		AdminName string
	}

	h.DB.Table("audit_logs al").
		Select("al.*, CONCAT(COALESCE(u.first_name,''), ' ', COALESCE(u.last_name,'')) as admin_name").
		Joins("LEFT JOIN users u ON u.id = al.user_id").
		Order("al.created_at DESC").
		Limit(50).
		Find(&rows)

	result := make([]auditLogRow, len(rows))
	for i, r := range rows {
		result[i] = auditLogRow{
			ID:           r.ID,
			Action:       r.Action,
			ResourceType: r.ResourceType,
			ResourceID:   r.ResourceID,
			AdminName:    r.AdminName,
			IPAddress:    r.IPAddress,
			CreatedAt:    r.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

type activityEntry struct {
	ID        uint      `json:"id"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Company   string    `json:"company,omitempty"`
	Status    string    `json:"status,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *LogsHandler) Activity(c *gin.Context) {
	type userRow struct {
		ID        uint
		FirstName string
		LastName  string
		Email     string
		CreatedAt time.Time
	}
	var users []userRow
	h.DB.Raw(`
		SELECT id, first_name, last_name, email, created_at
		FROM users
		WHERE deleted_at IS NULL
		  AND id NOT IN (SELECT user_id FROM user_roles)
		ORDER BY created_at DESC
		LIMIT 50
	`).Scan(&users)

	type providerRow struct {
		ID          uint
		FirstName   string
		LastName    string
		Email       string
		CompanyName string
		Status      string
		CreatedAt   time.Time
	}
	var providers []providerRow
	h.DB.Raw(`
		SELECT u.id, u.first_name, u.last_name, u.email, pp.company_name, pp.status, u.created_at
		FROM users u
		JOIN provider_profiles pp ON pp.user_id = u.id
		WHERE u.deleted_at IS NULL
		ORDER BY u.created_at DESC
		LIMIT 50
	`).Scan(&providers)

	var entries []activityEntry
	for _, u := range users {
		entries = append(entries, activityEntry{
			ID:        u.ID,
			Type:      "user",
			Name:      u.FirstName + " " + u.LastName,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
		})
	}
	for _, p := range providers {
		entries = append(entries, activityEntry{
			ID:        p.ID,
			Type:      "provider",
			Name:      p.FirstName + " " + p.LastName,
			Email:     p.Email,
			Company:   p.CompanyName,
			Status:    p.Status,
			CreatedAt: p.CreatedAt,
		})
	}

	for i := 0; i < len(entries)-1; i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[j].CreatedAt.After(entries[i].CreatedAt) {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}
	if len(entries) > 50 {
		entries = entries[:50]
	}

	c.JSON(http.StatusOK, gin.H{"data": entries})
}
