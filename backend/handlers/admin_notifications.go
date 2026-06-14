package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
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

func (h *AdminNotificationHandler) Broadcast(c *gin.Context) {
	var req struct {
		Title    string `json:"title" binding:"required,min=2"`
		Body     string `json:"body" binding:"required,min=2"`
		Link     string `json:"link"`
		Audience string `json:"audience" binding:"required,oneof=all particuliers pros user"`
		UserID   *uint  `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides (titre, message et audience requis)."})
		return
	}

	var ids []uint
	switch req.Audience {
	case "user":
		if req.UserID == nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Destinataire requis."})
			return
		}
		var u models.User
		if err := h.DB.Select("id").First(&u, *req.UserID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Utilisateur introuvable."})
			return
		}
		ids = []uint{u.ID}
	case "pros":
		h.DB.Model(&models.ProviderProfile{}).Where("status = ?", "approved").Pluck("user_id", &ids)
	case "particuliers":

		h.DB.Model(&models.User{}).
			Where("status = ? AND id NOT IN (SELECT user_id FROM provider_profiles WHERE status = 'approved')", "active").
			Pluck("id", &ids)
	default:
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

func (h *AdminNotificationHandler) RecipientSearch(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))
	if len(q) < 2 {
		c.JSON(http.StatusOK, gin.H{"data": []any{}})
		return
	}
	like := "%" + q + "%"
	type row struct {
		ID        uint
		FirstName string
		LastName  string
		Email     string
		Company   *string
	}
	var rows []row
	h.DB.Table("users u").
		Select("u.id, u.first_name, u.last_name, u.email, pp.company_name AS company").
		Joins("LEFT JOIN provider_profiles pp ON pp.user_id = u.id").
		Where("u.deleted_at IS NULL AND u.status = ?", "active").
		Where("u.email LIKE ? OR u.first_name LIKE ? OR u.last_name LIKE ? OR CONCAT(u.first_name, ' ', u.last_name) LIKE ? OR pp.company_name LIKE ?",
			like, like, like, like, like).
		Group("u.id").
		Limit(10).
		Scan(&rows)

	out := make([]gin.H, 0, len(rows))
	for _, r := range rows {
		out = append(out, gin.H{
			"id":      r.ID,
			"name":    strings.TrimSpace(r.FirstName + " " + r.LastName),
			"email":   r.Email,
			"company": r.Company,
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

func (h *AdminNotificationHandler) LinkTargets(c *gin.Context) {
	typ := c.Query("type")
	q := strings.TrimSpace(c.Query("q"))
	like := "%" + q + "%"
	out := []gin.H{}

	type row struct {
		ID    uint
		Title string
	}
	switch typ {
	case "prestation":
		var rows []row
		qy := h.DB.Table("prestations").Select("id, title").Where("deleted_at IS NULL")
		if q != "" {
			qy = qy.Where("title LIKE ?", like)
		}
		qy.Order("created_at DESC").Limit(10).Scan(&rows)
		for _, r := range rows {
			out = append(out, gin.H{"label": r.Title, "path": fmt.Sprintf("/prestations/%d", r.ID)})
		}
	case "event":
		var rows []row
		qy := h.DB.Table("events").Select("id, title").Where("deleted_at IS NULL")
		if q != "" {
			qy = qy.Where("title LIKE ?", like)
		}
		qy.Order("created_at DESC").Limit(10).Scan(&rows)
		for _, r := range rows {
			out = append(out, gin.H{"label": r.Title, "path": fmt.Sprintf("/evenements/%d", r.ID)})
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}
