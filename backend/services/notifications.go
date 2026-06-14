package services

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"upcycleconnect/backend/config"
	"upcycleconnect/backend/models"

	"gorm.io/gorm"
)

type NotificationService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewNotificationService(db *gorm.DB, cfg *config.Config) *NotificationService {
	return &NotificationService{db: db, cfg: cfg}
}

type NotifyInput struct {
	UserID uint
	Type   string
	Title  string
	Body   string
	Link   string
}

func (s *NotificationService) Notify(in NotifyInput) {
	if in.UserID == 0 || in.Title == "" {
		return
	}

	link := in.Link
	var linkPtr *string
	if link != "" {
		linkPtr = &link
	}

	n := models.Notification{
		UserID: in.UserID,
		Type:   in.Type,
		Title:  in.Title,
		Body:   in.Body,
		Link:   linkPtr,
	}
	if err := s.db.Create(&n).Error; err != nil {
		log.Printf("[notify] db create failed for user=%d type=%s: %v", in.UserID, in.Type, err)
		return
	}

	go s.sendPush(in)
}

func (s *NotificationService) sendPush(in NotifyInput) {
	if s.cfg.OneSignalAppID == "" || s.cfg.OneSignalAPIKey == "" {
		return
	}

	var user models.User
	if err := s.db.Select("id, one_signal_player_id").First(&user, in.UserID).Error; err != nil {
		return
	}
	if user.OneSignalPlayerID == nil || *user.OneSignalPlayerID == "" {
		return
	}

	payload := map[string]interface{}{
		"app_id":             s.cfg.OneSignalAppID,
		"include_player_ids": []string{*user.OneSignalPlayerID},
		"headings":           map[string]string{"en": in.Title, "fr": in.Title},
		"contents":           map[string]string{"en": in.Body, "fr": in.Body},
	}
	if in.Link != "" {
		url := in.Link
		if strings.HasPrefix(url, "/") {
			url = strings.TrimRight(s.cfg.AppURL, "/") + url
		}
		payload["url"] = url
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "https://onesignal.com/api/v1/notifications", bytes.NewReader(body))
	if err != nil {
		log.Printf("[notify] build push request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Basic "+s.cfg.OneSignalAPIKey)

	client := &http.Client{Timeout: 6 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[notify] push request failed: %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		log.Printf("[notify] push returned status %d", resp.StatusCode)
	}
}

func (s *NotificationService) NotifyMany(userIDs []uint, typ, title, body, link string) {
	for _, id := range userIDs {
		s.Notify(NotifyInput{UserID: id, Type: typ, Title: title, Body: body, Link: link})
	}
}

func (s *NotificationService) NotifyAdmins(typ, title, body, link string) {
	var ids []uint
	s.db.
		Table("user_roles").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("roles.name IN ?", []string{"admin", "super_admin"}).
		Pluck("user_roles.user_id", &ids)

	seen := make(map[uint]bool, len(ids))
	uniq := ids[:0]
	for _, id := range ids {
		if !seen[id] {
			seen[id] = true
			uniq = append(uniq, id)
		}
	}
	s.NotifyMany(uniq, typ, title, body, link)
}

func (s *NotificationService) MustNotify(userID uint, typ, title, body, link string) {
	s.Notify(NotifyInput{UserID: userID, Type: typ, Title: title, Body: body, Link: link})
}
