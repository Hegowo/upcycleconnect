package handlers

import (
	"math"
	"net/http"
	"strconv"

	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PrestationHandler struct {
	DB    *gorm.DB
	Audit *services.AuditService
}

func (h *PrestationHandler) Index(c *gin.Context) {
	query := h.DB.Model(&models.Prestation{}).Where("prestations.deleted_at IS NULL")

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if catID := c.Query("category_id"); catID != "" {
		query = query.Where("category_id = ?", catID)
	}
	if provID := c.Query("provider_id"); provID != "" {
		query = query.Where("provider_id = ?", provID)
	}
	if search := c.Query("search"); search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if perPage < 1 {
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

	var prestations []models.Prestation
	h.DB.Preload("Category").Preload("Provider.Roles").
		Where("prestations.deleted_at IS NULL").
		Scopes(func(db *gorm.DB) *gorm.DB {
			if status := c.Query("status"); status != "" {
				db = db.Where("status = ?", status)
			}
			if catID := c.Query("category_id"); catID != "" {
				db = db.Where("category_id = ?", catID)
			}
			if provID := c.Query("provider_id"); provID != "" {
				db = db.Where("provider_id = ?", provID)
			}
			if search := c.Query("search"); search != "" {
				db = db.Where("title LIKE ?", "%"+search+"%")
			}
			return db
		}).
		Order("created_at DESC").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&prestations)

	c.JSON(http.StatusOK, gin.H{
		"data": models.ToPrestationResponses(prestations),
		"meta": gin.H{
			"current_page": page,
			"last_page":    lastPage,
			"per_page":     perPage,
			"total":        total,
		},
	})
}

func (h *PrestationHandler) Store(c *gin.Context) {
	var req struct {
		CategoryID  *uint   `json:"category_id"`
		ProviderID  *uint   `json:"provider_id"`
		Title       string  `json:"title" binding:"required"`
		Description *string `json:"description"`
		Price       *string `json:"price"`
		PriceType   string  `json:"price_type"`
		Status      string  `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return
	}

	prestation := models.Prestation{
		CategoryID:  req.CategoryID,
		ProviderID:  req.ProviderID,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		PriceType:   req.PriceType,
		Status:      req.Status,
	}
	if prestation.PriceType == "" {
		prestation.PriceType = "fixed"
	}
	if prestation.Status == "" {
		prestation.Status = "draft"
	}

	if err := h.DB.Create(&prestation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création."})
		return
	}

	h.DB.Preload("Category").Preload("Provider.Roles").First(&prestation, prestation.ID)

	h.Audit.Log(c, "prestation.created", "Prestation", &prestation.ID, nil, map[string]interface{}{
		"title": prestation.Title,
	})

	c.JSON(http.StatusCreated, models.ToPrestationResponse(&prestation))
}

func (h *PrestationHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var prestation models.Prestation
	if err := h.DB.Preload("Category").Preload("Provider.Roles").
		Where("id = ? AND deleted_at IS NULL", id).First(&prestation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	c.JSON(http.StatusOK, models.ToPrestationResponse(&prestation))
}

func (h *PrestationHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var prestation models.Prestation
	if err := h.DB.Where("id = ? AND deleted_at IS NULL", id).First(&prestation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var req struct {
		CategoryID  *uint   `json:"category_id"`
		ProviderID  *uint   `json:"provider_id"`
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Price       *string `json:"price"`
		PriceType   *string `json:"price_type"`
		Status      *string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return
	}

	old := map[string]interface{}{
		"title":      prestation.Title,
		"status":     prestation.Status,
		"price":      prestation.Price,
		"price_type": prestation.PriceType,
	}

	updates := map[string]interface{}{}
	after := map[string]interface{}{"title": prestation.Title}

	if req.CategoryID != nil {
		updates["category_id"] = *req.CategoryID
	}
	if req.ProviderID != nil {
		updates["provider_id"] = *req.ProviderID
	}
	if req.Title != nil {
		updates["title"] = *req.Title
		after["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
		after["description"] = *req.Description
	}
	if req.Price != nil {
		updates["price"] = *req.Price
		after["price"] = *req.Price
	}
	if req.PriceType != nil {
		updates["price_type"] = *req.PriceType
		after["price_type"] = *req.PriceType
	}
	if req.Status != nil {
		updates["status"] = *req.Status
		after["status"] = *req.Status
	}

	if len(updates) > 0 {
		h.DB.Model(&prestation).Updates(updates)
	}

	h.DB.Preload("Category").Preload("Provider.Roles").First(&prestation, prestation.ID)

	h.Audit.Log(c, "prestation.updated", "Prestation", &prestation.ID, old, after)

	c.JSON(http.StatusOK, models.ToPrestationResponse(&prestation))
}

func (h *PrestationHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var prestation models.Prestation
	if err := h.DB.Where("id = ? AND deleted_at IS NULL", id).First(&prestation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=draft published suspended archived"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return
	}

	old := map[string]string{"status": prestation.Status}
	h.DB.Model(&prestation).Update("status", req.Status)

	h.Audit.Log(c, "prestation.status_changed", "Prestation", &prestation.ID, old, map[string]string{"status": req.Status})

	h.DB.Preload("Category").Preload("Provider.Roles").First(&prestation, prestation.ID)

	c.JSON(http.StatusOK, models.ToPrestationResponse(&prestation))
}

func (h *PrestationHandler) Destroy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var prestation models.Prestation
	if err := h.DB.Where("id = ? AND deleted_at IS NULL", id).First(&prestation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	h.Audit.Log(c, "prestation.deleted", "Prestation", &prestation.ID, map[string]interface{}{
		"title": prestation.Title,
	}, nil)
	h.DB.Delete(&prestation)

	c.Status(http.StatusNoContent)
}
