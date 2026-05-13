package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"upcycleconnect/backend/middleware"
	"upcycleconnect/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CalendarHandler struct {
	DB *gorm.DB
}

func generateCalendarToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (h *CalendarHandler) GetToken(c *gin.Context) {
	user := middleware.GetAuthUser(c)

	var dbUser models.User
	if err := h.DB.First(&dbUser, user.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur serveur."})
		return
	}

	if dbUser.CalendarToken == nil {
		token, err := generateCalendarToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur serveur."})
			return
		}
		dbUser.CalendarToken = &token
		if err := h.DB.Model(&dbUser).Update("calendar_token", token).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur serveur."})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"calendar_token": *dbUser.CalendarToken})
}

func (h *CalendarHandler) RegenerateToken(c *gin.Context) {
	user := middleware.GetAuthUser(c)

	token, err := generateCalendarToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur serveur."})
		return
	}

	if err := h.DB.Model(&models.User{}).Where("id = ?", user.ID).
		Update("calendar_token", token).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur serveur."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"calendar_token": token})
}

func (h *CalendarHandler) Feed(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.String(http.StatusUnauthorized, "Token manquant")
		return
	}

	var user models.User
	if err := h.DB.Where("calendar_token = ? AND deleted_at IS NULL", token).
		First(&user).Error; err != nil {
		c.String(http.StatusUnauthorized, "Token invalide")
		return
	}

	var regs []models.EventRegistration
	h.DB.Preload("Event").
		Joins("JOIN platform_events e ON e.id = event_registrations.event_id AND e.deleted_at IS NULL").
		Where("event_registrations.user_id = ?", user.ID).
		Find(&regs)

	var reservations []models.Reservation
	h.DB.Preload("Prestation").
		Where("user_id = ? AND deleted_at IS NULL", user.ID).
		Find(&reservations)

	var sb strings.Builder
	sb.WriteString("BEGIN:VCALENDAR\r\n")
	sb.WriteString("VERSION:2.0\r\n")
	sb.WriteString("PRODID:-//UpcycleConnect//Planning//FR\r\n")
	sb.WriteString("CALSCALE:GREGORIAN\r\n")
	sb.WriteString("METHOD:PUBLISH\r\n")
	sb.WriteString("X-WR-CALNAME:Mon Planning UpcycleConnect\r\n")
	sb.WriteString("X-WR-CALDESC:Réservations et événements UpcycleConnect\r\n")
	sb.WriteString("X-WR-TIMEZONE:Europe/Paris\r\n")

	for _, reg := range regs {
		if reg.Event == nil {
			continue
		}
		e := reg.Event
		uid := fmt.Sprintf("event-%d-user-%d@upcycleconnect", e.ID, user.ID)
		desc := ""
		if e.Description != nil {
			desc = icsEscape(*e.Description)
		}
		loc := ""
		if e.Location != nil {
			loc = icsEscape(*e.Location)
		}
		sb.WriteString("BEGIN:VEVENT\r\n")
		sb.WriteString(fmt.Sprintf("UID:%s\r\n", uid))
		sb.WriteString(fmt.Sprintf("DTSTAMP:%s\r\n", icsTime(time.Now())))
		sb.WriteString(fmt.Sprintf("DTSTART:%s\r\n", icsTime(e.StartDate)))
		sb.WriteString(fmt.Sprintf("DTEND:%s\r\n", icsTime(e.EndDate)))
		sb.WriteString(fmt.Sprintf("SUMMARY:%s\r\n", icsEscape(e.Title)))
		if desc != "" {
			sb.WriteString(fmt.Sprintf("DESCRIPTION:%s\r\n", desc))
		}
		if loc != "" {
			sb.WriteString(fmt.Sprintf("LOCATION:%s\r\n", loc))
		}
		sb.WriteString("CATEGORIES:Événement\r\n")
		sb.WriteString("END:VEVENT\r\n")
	}

	for _, res := range reservations {
		if res.Prestation == nil {
			continue
		}
		uid := fmt.Sprintf("reservation-%d@upcycleconnect", res.ID)
		title := res.Prestation.Title
		desc := fmt.Sprintf("Prestation réservée - Statut: %s", res.Status)
		if res.AmountCents > 0 {
			desc += fmt.Sprintf(" - Montant: %.2f€", float64(res.AmountCents)/100)
		}
		start := res.CreatedAt
		end := start.Add(1 * time.Hour)
		sb.WriteString("BEGIN:VEVENT\r\n")
		sb.WriteString(fmt.Sprintf("UID:%s\r\n", uid))
		sb.WriteString(fmt.Sprintf("DTSTAMP:%s\r\n", icsTime(time.Now())))
		sb.WriteString(fmt.Sprintf("DTSTART;VALUE=DATE:%s\r\n", start.UTC().Format("20060102")))
		sb.WriteString(fmt.Sprintf("DTEND;VALUE=DATE:%s\r\n", end.UTC().Format("20060102")))
		sb.WriteString(fmt.Sprintf("SUMMARY:Prestation: %s\r\n", icsEscape(title)))
		sb.WriteString(fmt.Sprintf("DESCRIPTION:%s\r\n", icsEscape(desc)))
		sb.WriteString("CATEGORIES:Prestation\r\n")
		sb.WriteString("END:VEVENT\r\n")
	}

	sb.WriteString("END:VCALENDAR\r\n")

	c.Header("Content-Type", "text/calendar; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=planning.ics")
	c.String(http.StatusOK, sb.String())
}

func icsTime(t time.Time) string {
	return t.UTC().Format("20060102T150405Z")
}

func icsEscape(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, ";", "\\;")
	s = strings.ReplaceAll(s, ",", "\\,")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "")
	return s
}
