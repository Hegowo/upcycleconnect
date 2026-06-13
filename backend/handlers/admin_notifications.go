package handlers

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminNotificationHandler struct {
	DB            *gorm.DB
	Notifications *services.NotificationService
	Audit         *services.AuditService
}

// SentList returns every notification dispatched to users, newest first,
// with the recipient's identity (admin oversight of "notifications envoyées").
func (h *AdminNotificationHandler) SentList(c *gin.Context) {
	q := h.DB.Model(&models.Notification{})
	if t := c.Query("type"); t != "" {
		q = q.Where("type = ?", t)
	}
	if search := c.Query("search"); search != "" {
		like := "%" + search + "%"
		q = q.Where("title LIKE ? OR body LIKE ?", like, like)
	}

	var total int64
	q.Count(&total)

	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	if perPage < 1 || perPage > 100 {
		perPage = 30
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	lastPage := int(math.Ceil(float64(total) / float64(perPage)))
	if lastPage < 1 {
		lastPage = 1
	}

	var rows []struct {
		models.Notification
		Email     string
		FirstName string
		LastName  string
	}
	q.Select("notifications.*, users.email, users.first_name, users.last_name").
		Joins("LEFT JOIN users ON users.id = notifications.user_id").
		Order("notifications.created_at DESC").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Scan(&rows)

	data := make([]gin.H, 0, len(rows))
	for _, r := range rows {
		data = append(data, gin.H{
			"id":         r.ID,
			"type":       r.Type,
			"title":      r.Title,
			"body":       r.Body,
			"link":       r.Link,
			"read":       r.ReadAt != nil,
			"recipient":  r.FirstName + " " + r.LastName,
			"email":      r.Email,
			"created_at": r.CreatedAt.UTC().Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
		"meta": gin.H{"current_page": page, "last_page": lastPage, "per_page": perPage, "total": total},
	})
}

// Broadcast sends a notification to a target audience (all, particuliers, pros).
func (h *AdminNotificationHandler) Broadcast(c *gin.Context) {
	var req struct {
		Title    string `json:"title" binding:"required,min=2"`
		Body     string `json:"body" binding:"required,min=2"`
		Link     string `json:"link"`
		Audience string `json:"audience" binding:"required,oneof=all particuliers pros"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides (titre, message et audience requis)."})
		return
	}

	// Resolve recipient ids based on audience.
	var ids []uint
	switch req.Audience {
	case "pros":
		h.DB.Model(&models.ProviderProfile{}).Where("status = ?", "approved").Pluck("user_id", &ids)
	case "particuliers":
		// Users without an approved provider profile.
		h.DB.Model(&models.User{}).
			Where("status = ? AND id NOT IN (SELECT user_id FROM provider_profiles WHERE status = 'approved')", "active").
			Pluck("id", &ids)
	default: // all
		h.DB.Model(&models.User{}).Where("status = ?", "active").Pluck("id", &ids)
	}

	if len(ids) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Aucun destinataire pour cette audience.", "sent": 0})
		return
	}

	h.Notifications.NotifyMany(ids, "admin.broadcast", req.Title, req.Body, req.Link)
	h.Audit.Log(c, "notification.broadcast", "Notification", nil, nil, map[string]interface{}{
		"audience": req.Audience, "count": len(ids), "title": req.Title,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Notification envoyée.", "sent": len(ids)})
}
