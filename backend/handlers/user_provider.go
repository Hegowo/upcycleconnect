package handlers

import (
	"math"
	"net/http"
	"regexp"
	"strconv"

	"upcycleconnect/backend/middleware"
	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserProviderHandler struct {
	DB    *gorm.DB
	Audit *services.AuditService
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
	h.DB.Model(&p).Update("status", "published")
	h.DB.Preload("Category").First(&p, p.ID)

	h.Audit.Log(c, "provider.prestation_submitted", "Prestation", &p.ID, old, map[string]string{"status": "published"})

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
