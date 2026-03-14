package handlers

import (
	"net/http"
	"strconv"

	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type CategoryHandler struct {
	DB    *gorm.DB
	Audit *services.AuditService
}

type categoryWithCount struct {
	models.PrestationCategory
	PrestationsCount int64 `gorm:"column:prestations_count"`
}

func (h *CategoryHandler) Index(c *gin.Context) {
	var categories []categoryWithCount
	h.DB.Model(&models.PrestationCategory{}).
		Select("prestation_categories.*, (SELECT COUNT(*) FROM prestations WHERE prestations.category_id = prestation_categories.id AND prestations.deleted_at IS NULL) as prestations_count").
		Where("prestation_categories.deleted_at IS NULL").
		Order("sort_order ASC, name ASC").
		Find(&categories)

	result := make([]models.CategoryResponse, len(categories))
	for i, cat := range categories {
		result[i] = models.ToCategoryResponse(&cat.PrestationCategory, cat.PrestationsCount)
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *CategoryHandler) Store(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Description *string `json:"description"`
		SortOrder   *int    `json:"sort_order"`
		IsActive    *bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return
	}

	category := models.PrestationCategory{
		Name: req.Name,
		Slug: slug.Make(req.Name),
	}
	if req.Description != nil {
		category.Description = req.Description
	}
	if req.SortOrder != nil {
		category.SortOrder = *req.SortOrder
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	} else {
		category.IsActive = true
	}

	if err := h.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création."})
		return
	}

	h.Audit.Log(c, "category.created", "PrestationCategory", &category.ID, nil, map[string]interface{}{
		"name": category.Name, "slug": category.Slug,
	})

	c.JSON(http.StatusCreated, models.ToCategoryResponse(&category, 0))
}

func (h *CategoryHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var cat categoryWithCount
	if err := h.DB.Model(&models.PrestationCategory{}).
		Select("prestation_categories.*, (SELECT COUNT(*) FROM prestations WHERE prestations.category_id = prestation_categories.id AND prestations.deleted_at IS NULL) as prestations_count").
		Where("prestation_categories.id = ? AND prestation_categories.deleted_at IS NULL", id).
		First(&cat).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	c.JSON(http.StatusOK, models.ToCategoryResponse(&cat.PrestationCategory, cat.PrestationsCount))
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var category models.PrestationCategory
	if err := h.DB.Where("id = ? AND deleted_at IS NULL", id).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var req struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		SortOrder   *int    `json:"sort_order"`
		IsActive    *bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return
	}

	old := map[string]interface{}{
		"name": category.Name, "slug": category.Slug,
		"sort_order": category.SortOrder, "is_active": category.IsActive,
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
		updates["slug"] = slug.Make(*req.Name)
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) > 0 {
		h.DB.Model(&category).Updates(updates)
		h.DB.Where("id = ? AND deleted_at IS NULL", id).First(&category)
	}

	h.Audit.Log(c, "category.updated", "PrestationCategory", &category.ID, old, map[string]interface{}{
		"name": category.Name, "slug": category.Slug,
	})

	// Get prestations count
	var count int64
	h.DB.Model(&models.Prestation{}).Where("category_id = ? AND deleted_at IS NULL", category.ID).Count(&count)

	c.JSON(http.StatusOK, models.ToCategoryResponse(&category, count))
}

func (h *CategoryHandler) Destroy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var category models.PrestationCategory
	if err := h.DB.Where("id = ? AND deleted_at IS NULL", id).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	h.Audit.Log(c, "category.deleted", "PrestationCategory", &category.ID, map[string]interface{}{
		"name": category.Name,
	}, nil)
	h.DB.Delete(&category)

	c.Status(http.StatusNoContent)
}

func (h *CategoryHandler) Toggle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var category models.PrestationCategory
	if err := h.DB.Where("id = ? AND deleted_at IS NULL", id).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	old := map[string]interface{}{"is_active": category.IsActive}
	newActive := !category.IsActive
	h.DB.Model(&category).Update("is_active", newActive)
	category.IsActive = newActive

	h.Audit.Log(c, "category.toggled", "PrestationCategory", &category.ID, old, map[string]interface{}{"is_active": category.IsActive})

	var count int64
	h.DB.Model(&models.Prestation{}).Where("category_id = ? AND deleted_at IS NULL", category.ID).Count(&count)

	c.JSON(http.StatusOK, models.ToCategoryResponse(&category, count))
}
