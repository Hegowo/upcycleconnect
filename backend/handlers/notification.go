package handlers

import (
	"net/http"
	"strconv"
	"time"

	"upcycleconnect/backend/middleware"
	"upcycleconnect/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NotificationHandler struct {
	DB *gorm.DB
}

func (h *NotificationHandler) Index(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "30"))
	if limit < 1 || limit > 100 {
		limit = 30
	}

	var items []models.Notification
	h.DB.Where("user_id = ?", user.ID).
		Order("created_at DESC").
		Limit(limit).
		Find(&items)

	var unread int64
	h.DB.Model(&models.Notification{}).
		Where("user_id = ? AND read_at IS NULL", user.ID).
		Count(&unread)

	c.JSON(http.StatusOK, gin.H{
		"data":         models.ToNotificationResponses(items),
		"unread_count": unread,
	})
}

func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	var unread int64
	h.DB.Model(&models.Notification{}).
		Where("user_id = ? AND read_at IS NULL", user.ID).
		Count(&unread)
	c.JSON(http.StatusOK, gin.H{"unread_count": unread})
}

func (h *NotificationHandler) MarkRead(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Notification introuvable"})
		return
	}
	now := time.Now()
	result := h.DB.Model(&models.Notification{}).
		Where("id = ? AND user_id = ? AND read_at IS NULL", id, user.ID).
		Update("read_at", now)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur interne"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"updated": result.RowsAffected})
}

func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	now := time.Now()
	result := h.DB.Model(&models.Notification{}).
		Where("user_id = ? AND read_at IS NULL", user.ID).
		Update("read_at", now)
	c.JSON(http.StatusOK, gin.H{"updated": result.RowsAffected})
}

func (h *NotificationHandler) RegisterPushToken(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	var body struct {
		PlayerID string `json:"player_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.PlayerID == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "player_id requis"})
		return
	}

	if len(body.PlayerID) > 100 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "player_id invalide"})
		return
	}
	if err := h.DB.Model(&models.User{}).
		Where("id = ?", user.ID).
		Update("one_signal_player_id", body.PlayerID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur interne"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Token enregistré"})
}

func (h *NotificationHandler) UnregisterPushToken(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	h.DB.Model(&models.User{}).
		Where("id = ?", user.ID).
		Update("one_signal_player_id", nil)
	c.JSON(http.StatusOK, gin.H{"message": "Token supprimé"})
}
