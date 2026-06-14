package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"upcycleconnect/backend/middleware"
	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReservationMessageHandler struct {
	DB            *gorm.DB
	Notifications *services.NotificationService
}

func (h *ReservationMessageHandler) access(c *gin.Context) (*models.Reservation, *models.Prestation, bool, bool) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return nil, nil, false, false
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Réservation introuvable"})
		return nil, nil, false, false
	}
	var reservation models.Reservation
	if err := h.DB.First(&reservation, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Réservation introuvable"})
		return nil, nil, false, false
	}
	var prestation models.Prestation
	if err := h.DB.First(&prestation, reservation.PrestationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Prestation introuvable"})
		return nil, nil, false, false
	}
	isClient := reservation.UserID == user.ID
	isProvider := prestation.ProviderID != nil && *prestation.ProviderID == user.ID
	if !isClient && !isProvider {
		c.JSON(http.StatusForbidden, gin.H{"message": "Accès refusé"})
		return nil, nil, false, false
	}
	return &reservation, &prestation, isProvider, true
}

func (h *ReservationMessageHandler) Index(c *gin.Context) {
	reservation, _, _, ok := h.access(c)
	if !ok {
		return
	}
	var msgs []models.ReservationMessage
	h.DB.Preload("Sender").
		Where("reservation_id = ?", reservation.ID).
		Order("created_at ASC").
		Find(&msgs)

	out := make([]models.ReservationMessageResponse, 0, len(msgs))
	for i := range msgs {
		out = append(out, models.ToReservationMessageResponse(&msgs[i]))
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

func (h *ReservationMessageHandler) Store(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	reservation, prestation, isProvider, ok := h.access(c)
	if !ok {
		return
	}
	var req struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Content) == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Message vide."})
		return
	}
	if len(req.Content) > 4000 {
		req.Content = req.Content[:4000]
	}

	msg := models.ReservationMessage{
		ReservationID: reservation.ID,
		SenderID:      user.ID,
		Content:       strings.TrimSpace(req.Content),
	}
	if err := h.DB.Create(&msg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de l'envoi du message."})
		return
	}
	msg.Sender = user

	if h.Notifications != nil {
		var recipient uint
		if isProvider {
			recipient = reservation.UserID
		} else if prestation.ProviderID != nil {
			recipient = *prestation.ProviderID
		}
		if recipient != 0 {
			h.Notifications.MustNotify(recipient, "reservation.message",
				"Nouveau message",
				fmt.Sprintf("%s : %s", strings.TrimSpace(user.FirstName+" "+user.LastName), truncateMsg(msg.Content, 80)),
				fmt.Sprintf("/profil/reservations/%d", reservation.ID))
		}
	}

	c.JSON(http.StatusCreated, gin.H{"data": models.ToReservationMessageResponse(&msg)})
}

func truncateMsg(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
