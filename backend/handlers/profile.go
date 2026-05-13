package handlers

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"upcycleconnect/backend/middleware"
	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ProfileHandler struct {
	DB     *gorm.DB
	Mailer *services.Mailer
}

func (h *ProfileHandler) Stats(c *gin.Context) {
	user := middleware.GetAuthUser(c)

	var deposits []models.DepositRequest
	h.DB.Preload("Category").Where("user_id = ?", user.ID).Order("created_at DESC").Find(&deposits)

	var co2Total float64
	var acceptedCount int
	for _, d := range deposits {
		if d.Status == "accepted" || d.Status == "validated" {
			acceptedCount++
			if d.CarbonSavings != nil {
				co2Total += *d.CarbonSavings
			}
		}
	}

	var regs []models.EventRegistration
	h.DB.Preload("Event").Preload("Event.Category").
		Joins("JOIN platform_events e ON e.id = event_registrations.event_id AND e.deleted_at IS NULL").
		Where("event_registrations.user_id = ?", user.ID).
		Order("e.start_date DESC").
		Find(&regs)

	var prestaReservations []models.Reservation
	h.DB.Preload("Prestation").
		Where("user_id = ? AND deleted_at IS NULL", user.ID).
		Order("created_at DESC").
		Find(&prestaReservations)

	score := int(float64(acceptedCount)*5 + co2Total*1.5)
	if score > 100 {
		score = 100
	}

	type depSummary struct {
		ID        uint    `json:"id"`
		Title     string  `json:"title"`
		Status    string  `json:"status"`
		Category  *string `json:"category"`
		Photo1    *string `json:"photo1,omitempty"`
		CreatedAt string  `json:"created_at"`
	}
	type regSummary struct {
		ID          uint    `json:"id"`
		Type        string  `json:"type"`
		EventID     uint    `json:"event_id,omitempty"`
		Title       string  `json:"title"`
		StartDate   string  `json:"start_date"`
		Location    *string `json:"location,omitempty"`
		Past        bool    `json:"past"`
		Status      string  `json:"status,omitempty"`
		AmountCents int64   `json:"amount_cents,omitempty"`
		Currency    string  `json:"currency,omitempty"`
	}
	type badgeDef struct {
		Key    string `json:"key"`
		Label  string `json:"label"`
		Desc   string `json:"desc"`
		Earned bool   `json:"earned"`
	}

	depSummaries := make([]depSummary, len(deposits))
	for i, d := range deposits {
		var cat *string
		if d.Category != nil {
			cat = &d.Category.Name
		}
		depSummaries[i] = depSummary{
			ID:        d.ID,
			Title:     d.Title,
			Status:    d.Status,
			Category:  cat,
			Photo1:    d.Photo1,
			CreatedAt: d.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		}
	}

	now := time.Now()
	regSummaries := make([]regSummary, 0, len(regs)+len(prestaReservations))
	for _, r := range regs {
		if r.Event == nil {
			continue
		}
		regSummaries = append(regSummaries, regSummary{
			ID:        r.ID,
			Type:      "event",
			EventID:   r.EventID,
			Title:     r.Event.Title,
			StartDate: r.Event.StartDate.UTC().Format("2006-01-02T15:04:05Z"),
			Location:  r.Event.Location,
			Past:      r.Event.EndDate.Before(now),
		})
	}
	for _, pr := range prestaReservations {
		title := "Prestation"
		if pr.Prestation != nil {
			title = pr.Prestation.Title
		}
		regSummaries = append(regSummaries, regSummary{
			ID:          pr.ID,
			Type:        "prestation",
			Title:       title,
			StartDate:   pr.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
			Past:        pr.Status == "paid" || pr.Status == "failed",
			Status:      pr.Status,
			AmountCents: pr.AmountCents,
			Currency:    pr.Currency,
		})
	}

	badges := []badgeDef{
		{Key: "premier_pas", Label: "Premier Pas", Desc: "Effectuer son premier dépôt", Earned: len(deposits) >= 1},
		{Key: "donneur", Label: "Donneur", Desc: "5 dépôts acceptés", Earned: acceptedCount >= 5},
		{Key: "reboiseur", Label: "Reboiseur", Desc: "10 kg de CO₂ économisés", Earned: co2Total >= 10},
		{Key: "activateur", Label: "Activateur", Desc: "Participer à un événement", Earned: len(regs) >= 1},
		{Key: "expert", Label: "Expert", Desc: "10 dépôts effectués", Earned: len(deposits) >= 10},
		{Key: "eco_champion", Label: "Éco-Champion", Desc: "Atteindre un score de 50", Earned: score >= 50},
	}

	c.JSON(http.StatusOK, gin.H{
		"score":         score,
		"co2_saved":     co2Total,
		"objects_saved": acceptedCount,
		"deposits":      depSummaries,
		"reservations":  regSummaries,
		"badges":        badges,
	})
}

func (h *ProfileHandler) UpdateInfo(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	var req struct {
		FirstName string  `json:"first_name" binding:"required"`
		LastName  string  `json:"last_name" binding:"required"`
		Phone     *string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return
	}
	if err := h.DB.Model(user).Updates(map[string]interface{}{
		"first_name": req.FirstName,
		"last_name":  req.LastName,
		"phone":      req.Phone,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la mise à jour."})
		return
	}
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Phone = req.Phone
	c.JSON(http.StatusOK, models.ToUserResponse(user))
}

func (h *ProfileHandler) UploadAvatar(c *gin.Context) {
	user := middleware.GetAuthUser(c)

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Fichier manquant."})
		return
	}
	if file.Size > 2*1024*1024 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Fichier trop volumineux (max 2 Mo)."})
		return
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowed[ext] {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Format non supporté (jpg, png, webp)."})
		return
	}

	dir := "/uploads/avatars"
	if err := os.MkdirAll(dir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur serveur."})
		return
	}
	filename := fmt.Sprintf("%d%s", user.ID, ext)
	savePath := filepath.Join(dir, filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de l'enregistrement."})
		return
	}

	url := fmt.Sprintf("/uploads/avatars/%s", filename)
	if err := h.DB.Model(user).Update("avatar_url", url).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la mise à jour."})
		return
	}
	user.AvatarURL = &url
	c.JSON(http.StatusOK, gin.H{"avatar_url": url})
}

func (h *ProfileHandler) ChangePassword(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	var req struct {
		Current string `json:"current" binding:"required"`
		New     string `json:"new" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides. Le nouveau mot de passe doit faire au moins 8 caractères."})
		return
	}

	var dbUser models.User
	h.DB.Where("id = ?", user.ID).First(&dbUser)
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(req.Current)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Mot de passe actuel incorrect."})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.New), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur serveur."})
		return
	}
	h.DB.Model(&dbUser).Update("password", string(hashed))
	c.JSON(http.StatusOK, gin.H{"message": "Mot de passe modifié avec succès."})
}

func (h *ProfileHandler) EmailChangeStart(c *gin.Context) {
	user := middleware.GetAuthUser(c)

	h.DB.Where("user_id = ?", user.ID).Delete(&models.EmailChangeRequest{})

	code := randCode()
	ecr := models.EmailChangeRequest{
		UserID:    user.ID,
		Step:      "verify_current",
		Code:      code,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	h.DB.Create(&ecr)

	if err := h.Mailer.SendEmailChangeCode(user.Email, user.FirstName, code); err != nil {
		log.Printf("[mailer] email change start to %s: %v", user.Email, err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Code envoyé à votre adresse actuelle."})
}

func (h *ProfileHandler) EmailVerifyCurrent(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	var req struct {
		Code     string `json:"code" binding:"required"`
		NewEmail string `json:"new_email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return
	}

	var ecr models.EmailChangeRequest
	if err := h.DB.Where("user_id = ? AND step = ? AND code = ?", user.ID, "verify_current", req.Code).
		First(&ecr).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Code incorrect."})
		return
	}
	if time.Now().After(ecr.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{"message": "Code expiré. Recommencez."})
		return
	}

	var count int64
	h.DB.Model(&models.User{}).
		Where("email = ? AND deleted_at IS NULL AND id != ?", req.NewEmail, user.ID).
		Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"message": "Cette adresse email est déjà utilisée."})
		return
	}

	newCode := randCode()
	newEmail := req.NewEmail
	h.DB.Model(&ecr).Updates(map[string]interface{}{
		"step":       "verify_new",
		"code":       newCode,
		"new_email":  newEmail,
		"expires_at": time.Now().Add(10 * time.Minute),
	})

	if err := h.Mailer.SendEmailChangeCode(newEmail, user.FirstName, newCode); err != nil {
		log.Printf("[mailer] email change new to %s: %v", newEmail, err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Code envoyé à votre nouvelle adresse."})
}

func (h *ProfileHandler) EmailVerifyNew(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return
	}

	var ecr models.EmailChangeRequest
	if err := h.DB.Where("user_id = ? AND step = ? AND code = ?", user.ID, "verify_new", req.Code).
		First(&ecr).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Code incorrect."})
		return
	}
	if time.Now().After(ecr.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{"message": "Code expiré. Recommencez."})
		return
	}
	if ecr.NewEmail == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur interne."})
		return
	}

	h.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("email", *ecr.NewEmail)
	h.DB.Delete(&ecr)

	c.JSON(http.StatusOK, gin.H{
		"message": "Email modifié avec succès.",
		"email":   *ecr.NewEmail,
	})
}

func (h *ProfileHandler) RegisterForEvent(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	eventID := c.Param("id")

	var event models.Event
	if err := h.DB.Where("id = ? AND status = 'published' AND deleted_at IS NULL", eventID).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Événement introuvable."})
		return
	}
	if time.Now().After(event.EndDate) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Cet événement est terminé."})
		return
	}
	if event.MaxParticipants != nil {
		var count int64
		h.DB.Model(&models.EventRegistration{}).Where("event_id = ?", event.ID).Count(&count)
		if int(count) >= *event.MaxParticipants {
			c.JSON(http.StatusConflict, gin.H{"message": "Cet événement est complet."})
			return
		}
	}

	reg := models.EventRegistration{UserID: user.ID, EventID: event.ID}
	if err := h.DB.Create(&reg).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Vous êtes déjà inscrit à cet événement."})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Inscription confirmée.", "id": reg.ID})
}

func (h *ProfileHandler) UnregisterFromEvent(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	eventID := c.Param("id")

	result := h.DB.Where("user_id = ? AND event_id = ?", user.ID, eventID).
		Delete(&models.EventRegistration{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Inscription introuvable."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Désinscription effectuée."})
}

func randCode() string {
	b := make([]byte, 4)
	rand.Read(b)
	v := (uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])) % 1000000
	return fmt.Sprintf("%06d", v)
}
