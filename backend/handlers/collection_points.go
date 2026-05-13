package handlers

import (
	"net/http"
	"strconv"

	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CollectionPointHandler struct {
	DB    *gorm.DB
	Audit *services.AuditService
}

func (h *CollectionPointHandler) Index(c *gin.Context) {
	var items []models.CollectionPoint
	h.DB.Order("name ASC").Find(&items)
	c.JSON(http.StatusOK, models.ToCollectionPointResponses(items))
}

func (h *CollectionPointHandler) PublicIndex(c *gin.Context) {
	var items []models.CollectionPoint
	h.DB.Where("is_active = true").Order("name ASC").Find(&items)
	c.JSON(http.StatusOK, models.ToCollectionPointResponses(items))
}

func (h *CollectionPointHandler) Store(c *gin.Context) {
	var req struct {
		Name         string  `json:"name" binding:"required"`
		Address      string  `json:"address" binding:"required"`
		City         string  `json:"city" binding:"required"`
		PostalCode   string  `json:"postal_code" binding:"required"`
		Latitude     float64 `json:"latitude" binding:"required"`
		Longitude    float64 `json:"longitude" binding:"required"`
		Phone        *string `json:"phone"`
		OpeningHours *string `json:"opening_hours"`
		IsActive     *bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return
	}

	active := true
	if req.IsActive != nil {
		active = *req.IsActive
	}

	cp := models.CollectionPoint{
		Name:         req.Name,
		Address:      req.Address,
		City:         req.City,
		PostalCode:   req.PostalCode,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Phone:        req.Phone,
		OpeningHours: req.OpeningHours,
		IsActive:     active,
	}
	if err := h.DB.Create(&cp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création."})
		return
	}

	h.Audit.Log(c, "collection_point.created", "CollectionPoint", &cp.ID, nil, nil)
	c.JSON(http.StatusCreated, models.ToCollectionPointResponse(&cp))
}

func (h *CollectionPointHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Introuvable"})
		return
	}

	var cp models.CollectionPoint
	if err := h.DB.First(&cp, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Introuvable"})
		return
	}

	var req struct {
		Name         string  `json:"name"`
		Address      string  `json:"address"`
		City         string  `json:"city"`
		PostalCode   string  `json:"postal_code"`
		Latitude     float64 `json:"latitude"`
		Longitude    float64 `json:"longitude"`
		Phone        *string `json:"phone"`
		OpeningHours *string `json:"opening_hours"`
		IsActive     *bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return
	}

	updates := map[string]interface{}{
		"name":          req.Name,
		"address":       req.Address,
		"city":          req.City,
		"postal_code":   req.PostalCode,
		"latitude":      req.Latitude,
		"longitude":     req.Longitude,
		"phone":         req.Phone,
		"opening_hours": req.OpeningHours,
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	h.DB.Model(&cp).Updates(updates)

	h.Audit.Log(c, "collection_point.updated", "CollectionPoint", &cp.ID, nil, nil)
	h.DB.First(&cp, cp.ID)
	c.JSON(http.StatusOK, models.ToCollectionPointResponse(&cp))
}

func (h *CollectionPointHandler) Destroy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Introuvable"})
		return
	}

	var cp models.CollectionPoint
	if err := h.DB.First(&cp, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Introuvable"})
		return
	}

	h.DB.Delete(&cp)
	h.Audit.Log(c, "collection_point.deleted", "CollectionPoint", &cp.ID, nil, nil)
	c.JSON(http.StatusOK, gin.H{"message": "Point de collecte supprimé."})
}
