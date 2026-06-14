package handlers

import (
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"upcycleconnect/backend/middleware"
	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func saveKbisFile(header *multipart.FileHeader, src multipart.File, userID uint) (string, error) {
	ext := filepath.Ext(header.Filename)
	allowed := map[string]bool{".pdf": true, ".jpg": true, ".jpeg": true, ".png": true}
	if !allowed[ext] {
		return "", fmt.Errorf("format invalide")
	}
	if header.Size > 10*1024*1024 {
		return "", fmt.Errorf("fichier trop volumineux")
	}
	dir := "/uploads/kbis"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	destPath := filepath.Join(dir, fmt.Sprintf("kbis-%d%s", userID, ext))
	dst, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}
	return destPath, nil
}

type UserProviderHandler struct {
	DB            *gorm.DB
	Audit         *services.AuditService
	Notifications *services.NotificationService
}

var siretRegex = regexp.MustCompile(`^\d{14}$`)

func (h *UserProviderHandler) Apply(c *gin.Context) {
	user := middleware.GetAuthUser(c)

	var existing models.ProviderProfile
	if err := h.DB.Where("user_id = ?", user.ID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"message": "Une candidature prestataire existe déjà.",
			"profile": models.ToProviderProfileResponse(&existing),
		})
		return
	}

	var req struct {
		CompanyName string  `json:"company_name" binding:"required"`
		Siret       *string `json:"siret"`
		Description *string `json:"description"`
		Website     *string `json:"website"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides. Le nom de l'entreprise est requis."})
		return
	}

	if req.Siret != nil && *req.Siret != "" && !siretRegex.MatchString(*req.Siret) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "SIRET invalide (14 chiffres requis)."})
		return
	}

	profile := models.ProviderProfile{
		UserID:      user.ID,
		CompanyName: req.CompanyName,
		Siret:       req.Siret,
		Description: req.Description,
		Website:     req.Website,
		Status:      "pending",
	}
	if err := h.DB.Create(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création."})
		return
	}

	h.Audit.Log(c, "provider.applied", "ProviderProfile", &profile.ID, nil, map[string]interface{}{
		"company_name": profile.CompanyName,
	})

	if h.Notifications != nil {
		h.Notifications.NotifyAdmins("provider.applied",
			"Nouvelle candidature prestataire",
			profile.CompanyName+" a soumis une candidature à valider.",
			"/admin/providers")
		h.Notifications.MustNotify(user.ID, "provider.application_received",
			"Candidature enregistrée",
			"Votre candidature prestataire « "+profile.CompanyName+" » a bien été reçue. Elle sera examinée par un administrateur.",
			"/profil/parametres")
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Candidature prestataire enregistrée. Elle sera examinée par un administrateur.",
		"profile": models.ToProviderProfileResponse(&profile),
	})
}

func (h *UserProviderHandler) GetProfile(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	var profile models.ProviderProfile
	if err := h.DB.Where("user_id = ?", user.ID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Aucun profil prestataire trouvé."})
		return
	}
	c.JSON(http.StatusOK, models.ToProviderProfileResponse(&profile))
}

func (h *UserProviderHandler) UpdateProfile(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	var profile models.ProviderProfile
	if err := h.DB.Where("user_id = ?", user.ID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Aucun profil prestataire trouvé."})
		return
	}

	var req struct {
		CompanyName *string `json:"company_name"`
		Description *string `json:"description"`
		Website     *string `json:"website"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return
	}

	updates := map[string]interface{}{}
	old := map[string]interface{}{"company_name": profile.CompanyName}

	if req.CompanyName != nil && *req.CompanyName != "" {
		updates["company_name"] = *req.CompanyName
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Website != nil {
		updates["website"] = *req.Website
	}

	if len(updates) > 0 {
		h.DB.Model(&profile).Updates(updates)
		h.DB.Where("user_id = ?", user.ID).First(&profile)
	}

	h.Audit.Log(c, "provider.profile_updated", "ProviderProfile", &profile.ID, old, updates)

	c.JSON(http.StatusOK, models.ToProviderProfileResponse(&profile))
}

func (h *UserProviderHandler) ensureApproved(c *gin.Context, user *models.User) *models.ProviderProfile {
	var profile models.ProviderProfile
	if err := h.DB.Where("user_id = ?", user.ID).First(&profile).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Profil prestataire requis."})
		return nil
	}
	if profile.Status != "approved" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Votre profil prestataire n'est pas encore approuvé."})
		return nil
	}
	return &profile
}

func (h *UserProviderHandler) ListPrestations(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if h.ensureApproved(c, user) == nil {
		return
	}

	query := h.DB.Model(&models.Prestation{}).
		Where("provider_id = ? AND deleted_at IS NULL", user.ID)
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	lastPage := int(math.Ceil(float64(total) / float64(perPage)))
	if lastPage < 1 {
		lastPage = 1
	}

	var items []models.Prestation
	fetch := h.DB.Preload("Category").Preload("Images").
		Where("provider_id = ? AND deleted_at IS NULL", user.ID)
	if status := c.Query("status"); status != "" {
		fetch = fetch.Where("status = ?", status)
	}
	fetch.Order("created_at DESC").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&items)

	c.JSON(http.StatusOK, gin.H{
		"data": models.ToPrestationResponses(items),
		"meta": gin.H{
			"current_page": page,
			"last_page":    lastPage,
			"per_page":     perPage,
			"total":        total,
		},
	})
}

func (h *UserProviderHandler) CreatePrestation(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if h.ensureApproved(c, user) == nil {
		return
	}

	var req struct {
		CategoryID  *uint   `json:"category_id"`
		Title       string  `json:"title" binding:"required"`
		Description *string `json:"description"`
		Price       *string `json:"price"`
		PriceType   string  `json:"price_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return
	}

	allowed := map[string]bool{"": true, "fixed": true, "hourly": true, "quote": true}
	if !allowed[req.PriceType] {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Type de prix invalide (fixed, hourly, quote)."})
		return
	}

	providerID := user.ID
	p := models.Prestation{
		CategoryID:  req.CategoryID,
		ProviderID:  &providerID,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		PriceType:   req.PriceType,
		Status:      "draft",
	}
	if p.PriceType == "" {
		p.PriceType = "fixed"
	}

	if err := h.DB.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création."})
		return
	}
	h.DB.Preload("Category").First(&p, p.ID)

	h.Audit.Log(c, "provider.prestation_created", "Prestation", &p.ID, nil, map[string]interface{}{
		"title": p.Title,
	})

	c.JSON(http.StatusCreated, models.ToPrestationResponse(&p))
}

func (h *UserProviderHandler) UpdatePrestation(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if h.ensureApproved(c, user) == nil {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var p models.Prestation
	if err := h.DB.Where("id = ? AND provider_id = ? AND deleted_at IS NULL", id, user.ID).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var req struct {
		CategoryID  *uint   `json:"category_id"`
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Price       *string `json:"price"`
		PriceType   *string `json:"price_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return
	}

	updates := map[string]interface{}{}
	old := map[string]interface{}{"title": p.Title, "status": p.Status}

	if req.CategoryID != nil {
		updates["category_id"] = *req.CategoryID
	}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.PriceType != nil {
		allowed := map[string]bool{"fixed": true, "hourly": true, "quote": true}
		if !allowed[*req.PriceType] {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Type de prix invalide."})
			return
		}
		updates["price_type"] = *req.PriceType
	}

	if p.Status == "published" {
		updates["status"] = "draft"
	}

	if len(updates) > 0 {
		h.DB.Model(&p).Updates(updates)
	}
	h.DB.Preload("Category").First(&p, p.ID)

	h.Audit.Log(c, "provider.prestation_updated", "Prestation", &p.ID, old, updates)

	c.JSON(http.StatusOK, models.ToPrestationResponse(&p))
}

func (h *UserProviderHandler) SubmitPrestation(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if h.ensureApproved(c, user) == nil {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var p models.Prestation
	if err := h.DB.Where("id = ? AND provider_id = ? AND deleted_at IS NULL", id, user.ID).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	if p.Status != "draft" && p.Status != "archived" {
		c.JSON(http.StatusConflict, gin.H{"message": "La prestation ne peut pas être soumise dans son état actuel."})
		return
	}

	old := map[string]string{"status": p.Status}

	h.DB.Model(&p).Update("status", "pending")
	h.DB.Preload("Category").First(&p, p.ID)

	h.Audit.Log(c, "provider.prestation_submitted", "Prestation", &p.ID, old, map[string]string{"status": "pending"})

	if h.Notifications != nil {
		h.Notifications.NotifyAdmins("prestation.pending_validation",
			"Annonce à valider",
			fmt.Sprintf("%s a soumis l'annonce « %s » pour validation.",
				user.FirstName+" "+user.LastName, p.Title),
			"/admin/prestations")

		h.Notifications.MustNotify(user.ID, "prestation.submitted",
			"Annonce soumise",
			"Votre prestation « "+p.Title+" » a été soumise et sera examinée par un administrateur.",
			"/profil/pro")
	}

	c.JSON(http.StatusOK, models.ToPrestationResponse(&p))
}

func (h *UserProviderHandler) DestroyPrestation(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if h.ensureApproved(c, user) == nil {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var p models.Prestation
	if err := h.DB.Where("id = ? AND provider_id = ? AND deleted_at IS NULL", id, user.ID).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	h.Audit.Log(c, "provider.prestation_deleted", "Prestation", &p.ID, map[string]interface{}{
		"title": p.Title,
	}, nil)
	h.DB.Delete(&p)
	c.Status(http.StatusNoContent)
}

func (h *UserProviderHandler) UploadKbis(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}

	var profile models.ProviderProfile
	if err := h.DB.Where("user_id = ?", user.ID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Aucun profil prestataire trouvé. Soumets d'abord ta candidature."})
		return
	}

	file, header, err := c.Request.FormFile("kbis")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Fichier Kbis requis (champ 'kbis')"})
		return
	}
	defer file.Close()

	destPath, err := saveKbisFile(header, file, user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Format invalide (PDF/JPG/PNG, max 10 Mo)"})
		return
	}

	if err := h.DB.Model(&profile).Update("kbis_path", destPath).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la mise à jour du profil"})
		return
	}
	profile.KbisPath = &destPath

	h.Audit.Log(c, "provider.kbis_uploaded", "ProviderProfile", &profile.ID, nil, map[string]interface{}{
		"filename": filepath.Base(destPath),
	})

	c.JSON(http.StatusOK, gin.H{
		"message":  "Kbis enregistré avec succès.",
		"has_kbis": true,
		"profile":  models.ToProviderProfileResponse(&profile),
	})
}

func (h *UserProviderHandler) DownloadKbis(c *gin.Context) {
	providerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID invalide"})
		return
	}

	var profile models.ProviderProfile
	if err := h.DB.Where("user_id = ?", providerID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Profil prestataire introuvable"})
		return
	}

	if profile.KbisPath == nil || *profile.KbisPath == "" {
		c.JSON(http.StatusNotFound, gin.H{"message": "Aucun Kbis déposé pour ce prestataire"})
		return
	}

	path := *profile.KbisPath
	if _, err := os.Stat(path); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Fichier introuvable sur le serveur"})
		return
	}

	ext := filepath.Ext(path)
	contentType := "application/octet-stream"
	switch ext {
	case ".pdf":
		contentType = "application/pdf"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	}

	filename := fmt.Sprintf("kbis-provider-%d%s", providerID, ext)
	c.Header("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, filename))
	c.Header("Content-Type", contentType)
	c.File(path)
}

func (h *UserProviderHandler) CompleteOnboarding(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	var profile models.ProviderProfile
	if err := h.DB.Where("user_id = ?", user.ID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Profil prestataire introuvable"})
		return
	}
	if profile.OnboardingCompletedAt != nil {
		c.JSON(http.StatusOK, models.ToProviderProfileResponse(&profile))
		return
	}
	now := time.Now()
	if err := h.DB.Model(&profile).Update("onboarding_completed_at", now).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de l'enregistrement"})
		return
	}
	profile.OnboardingCompletedAt = &now
	h.Audit.Log(c, "provider.onboarded", "ProviderProfile", &profile.ID, nil, nil)
	c.JSON(http.StatusOK, models.ToProviderProfileResponse(&profile))
}

func (h *UserProviderHandler) ownedPrestation(c *gin.Context) (*models.Prestation, bool) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return nil, false
	}
	var p models.Prestation
	if err := h.DB.Where("id = ? AND deleted_at IS NULL", c.Param("id")).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Prestation introuvable"})
		return nil, false
	}
	if p.ProviderID == nil || *p.ProviderID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"message": "Accès refusé"})
		return nil, false
	}
	return &p, true
}

func (h *UserProviderHandler) ShowPrestation(c *gin.Context) {
	p, ok := h.ownedPrestation(c)
	if !ok {
		return
	}
	h.DB.Preload("Category").Preload("Images").First(p, p.ID)
	c.JSON(http.StatusOK, models.ToPrestationResponse(p))
}

func (h *UserProviderHandler) ListParticipants(c *gin.Context) {
	p, ok := h.ownedPrestation(c)
	if !ok {
		return
	}
	type row struct {
		ID           uint
		FirstName    string
		LastName     string
		Email        string
		AvatarURL    *string
		RegisteredAt time.Time
	}
	var rows []row
	h.DB.Table("reservations r").
		Select("u.id, u.first_name, u.last_name, u.email, u.avatar_url, MIN(r.created_at) AS registered_at").
		Joins("JOIN users u ON u.id = r.user_id").
		Where("r.prestation_id = ? AND r.deleted_at IS NULL", p.ID).
		Group("u.id").
		Order("registered_at DESC").
		Scan(&rows)

	out := make([]gin.H, 0, len(rows))
	for _, r := range rows {
		out = append(out, gin.H{
			"id":            r.ID,
			"first_name":    r.FirstName,
			"last_name":     r.LastName,
			"email":         r.Email,
			"avatar_url":    r.AvatarURL,
			"registered_at": r.RegisteredAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

func (h *UserProviderHandler) UploadPrestationImage(c *gin.Context) {
	p, ok := h.ownedPrestation(c)
	if !ok {
		return
	}
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Image requise (champ 'image')"})
		return
	}
	defer file.Close()
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowed[ext] {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Format invalide (JPG, PNG, WEBP)"})
		return
	}
	if header.Size > 5*1024*1024 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Image trop volumineuse (max 5 Mo)"})
		return
	}
	dir := "/uploads/prestations"
	if err := os.MkdirAll(dir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur de stockage"})
		return
	}
	filename := fmt.Sprintf("prestation-%d-%d%s", p.ID, time.Now().UnixNano(), ext)
	dst, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur de stockage"})
		return
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de l'enregistrement"})
		return
	}

	var count int64
	h.DB.Model(&models.PrestationImage{}).Where("prestation_id = ?", p.ID).Count(&count)
	img := models.PrestationImage{PrestationID: p.ID, URL: "/uploads/prestations/" + filename, Position: int(count)}
	if err := h.DB.Create(&img).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de l'enregistrement"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": img.ID, "url": img.URL, "order": img.Position})
}

func (h *UserProviderHandler) DeletePrestationImage(c *gin.Context) {
	p, ok := h.ownedPrestation(c)
	if !ok {
		return
	}
	res := h.DB.Where("id = ? AND prestation_id = ?", c.Param("imageId"), p.ID).Delete(&models.PrestationImage{})
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Image introuvable"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Image supprimée."})
}
