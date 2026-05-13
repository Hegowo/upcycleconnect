package handlers

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"upcycleconnect/backend/middleware"
	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EventMessageHandler struct {
	DB    *gorm.DB
	Audit *services.AuditService
}

// checkAccess returns true if the user is the event organizer or a registered participant.
func (h *EventMessageHandler) checkAccess(user *models.User, eventID uint) bool {
	var event models.Event
	if err := h.DB.Where("id = ? AND deleted_at IS NULL", eventID).First(&event).Error; err != nil {
		return false
	}
	if event.CreatedBy != nil && user.ID == *event.CreatedBy {
		return true
	}
	var reg models.EventRegistration
	err := h.DB.Where("user_id = ? AND event_id = ?", user.ID, eventID).First(&reg).Error
	return err == nil
}

// Index handles GET /api/v1/events/:id/messages
func (h *EventMessageHandler) Index(c *gin.Context) {
	user := middleware.GetAuthUser(c)

	var eventID uint
	if _, err := fmt.Sscanf(c.Param("id"), "%d", &eventID); err != nil || eventID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var event models.Event
	if err := h.DB.Where("id = ? AND deleted_at IS NULL", eventID).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	if !h.checkAccess(user, eventID) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Accès refusé."})
		return
	}

	var messages []models.EventMessage
	h.DB.Preload("User").Where("event_id = ?", eventID).Order("created_at ASC").Find(&messages)

	type messageResponse struct {
		ID          uint    `json:"id"`
		UserID      uint    `json:"user_id"`
		UserName    string  `json:"user_name"`
		Avatar      *string `json:"avatar"`
		Content     string  `json:"content"`
		ImageURL    *string `json:"image_url"`
		IsMe        bool    `json:"is_me"`
		IsOrganizer bool    `json:"is_organizer"`
		CreatedAt   string  `json:"created_at"`
	}

	data := make([]messageResponse, 0, len(messages))
	for _, msg := range messages {
		var userName string
		var avatar *string
		if msg.User != nil {
			userName = strings.TrimSpace(msg.User.FirstName + " " + msg.User.LastName)
			avatar = msg.User.AvatarURL
		}

		isOrganizer := false
		if event.CreatedBy != nil {
			isOrganizer = msg.UserID == *event.CreatedBy
		}

		data = append(data, messageResponse{
			ID:          msg.ID,
			UserID:      msg.UserID,
			UserName:    userName,
			Avatar:      avatar,
			Content:     msg.Content,
			ImageURL:    msg.ImageURL,
			IsMe:        msg.UserID == user.ID,
			IsOrganizer: isOrganizer,
			CreatedAt:   msg.CreatedAt.UTC().Format("2006-01-02T15:04:05.000000Z"),
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

// Store handles POST /api/v1/events/:id/messages
func (h *EventMessageHandler) Store(c *gin.Context) {
	user := middleware.GetAuthUser(c)

	var eventID uint
	if _, err := fmt.Sscanf(c.Param("id"), "%d", &eventID); err != nil || eventID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var event models.Event
	if err := h.DB.Where("id = ? AND deleted_at IS NULL", eventID).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	if !h.checkAccess(user, eventID) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Accès refusé."})
		return
	}

	var req struct {
		Content  string  `json:"content"`
		ImageURL *string `json:"image_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return
	}

	if strings.TrimSpace(req.Content) == "" && req.ImageURL == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Un message ou une image est requis."})
		return
	}

	msg := models.EventMessage{
		EventID:  eventID,
		UserID:   user.ID,
		Content:  req.Content,
		ImageURL: req.ImageURL,
	}
	if err := h.DB.Create(&msg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de l'envoi du message."})
		return
	}

	h.DB.Preload("User").First(&msg, msg.ID)

	var userName string
	var avatar *string
	if msg.User != nil {
		userName = strings.TrimSpace(msg.User.FirstName + " " + msg.User.LastName)
		avatar = msg.User.AvatarURL
	}

	isOrganizer := false
	if event.CreatedBy != nil {
		isOrganizer = msg.UserID == *event.CreatedBy
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":           msg.ID,
		"user_id":      msg.UserID,
		"user_name":    userName,
		"avatar":       avatar,
		"content":      msg.Content,
		"image_url":    msg.ImageURL,
		"is_me":        true,
		"is_organizer": isOrganizer,
		"created_at":   msg.CreatedAt.UTC().Format("2006-01-02T15:04:05.000000Z"),
	})
}

// UploadImage handles POST /api/v1/events/:id/messages/image
func (h *EventMessageHandler) UploadImage(c *gin.Context) {
	user := middleware.GetAuthUser(c)

	var eventID uint
	if _, err := fmt.Sscanf(c.Param("id"), "%d", &eventID); err != nil || eventID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var event models.Event
	if err := h.DB.Where("id = ? AND deleted_at IS NULL", eventID).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	if !h.checkAccess(user, eventID) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Accès refusé."})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Fichier image requis."})
		return
	}

	const maxSize = 5 << 20 // 5 MB
	if file.Size > maxSize {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "L'image ne doit pas dépasser 5 Mo."})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true}
	if !allowed[ext] {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Format non supporté (jpg, jpeg, png, webp, gif)."})
		return
	}

	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur serveur."})
		return
	}
	filename := fmt.Sprintf("%x%s", b, ext)

	dir := "/uploads/chat"
	if err := os.MkdirAll(dir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur serveur."})
		return
	}

	savePath := filepath.Join(dir, filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de l'enregistrement."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": "/uploads/chat/" + filename})
}
