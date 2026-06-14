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
	"time"

	"upcycleconnect/backend/middleware"
	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// saveKbisFile validates and writes a Kbis upload to /uploads/kbis, returning the
// stored path. Shared by the authenticated UploadKbis endpoint and by registration
// (where the provider isn't active yet and can't authenticate). Accepts PDF/JPG/PNG
// up to 10 MB and names the file after the user id so a re-upload replaces it.
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
	fetch := h.DB.Preload("Category").
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
	// Annonces require validation by the administrative service (per spec p.5):
	// submission moves the prestation to "pending", an admin then publishes it.
	h.DB.Model(&p).Update("status", "pending")
	h.DB.Preload("Category").First(&p, p.ID)

	h.Audit.Log(c, "provider.prestation_submitted", "Prestation", &p.ID, old, map[string]string{"status": "pending"})

	if h.Notifications != nil {
		// Notify admins that an annonce awaits validation.
		h.Notifications.NotifyAdmins("prestation.pending_validation",
			"Annonce à valider",
			fmt.Sprintf("%s a soumis l'annonce « %s » pour validation.",
				user.FirstName+" "+user.LastName, p.Title),
			"/admin/prestations")
		// Confirm reception to the provider.
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

// UploadKbis accepts a multipart PDF or image file as the provider's official Kbis document.
// Replaces any previously uploaded file. The file is stored in /uploads/kbis/ and is never
// deleted — admins can always retrieve it even after the provider is approved.
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

// DownloadKbis streams the Kbis file to the admin. Only admins can access this route.
func (h *UserProviderHandler) DownloadKbis(c *gin.Context) {
	providerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID invalide"})
		return
	}

	// Load the provider's profile
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

// CompleteOnboarding marks the provider's onboarding wizard as completed.
// Called by the frontend when the pro finishes the interactive setup on /profil/pro.
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
