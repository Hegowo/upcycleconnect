package handlers

import (
	"math"
	"net/http"
	"path/filepath"
	"strconv"

	"upcycleconnect/backend/middleware"
	"upcycleconnect/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InvoiceHandler struct {
	DB *gorm.DB
}

func (h *InvoiceHandler) Index(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}

	query := h.DB.Model(&models.Invoice{}).Where("user_id = ?", user.ID)
	if t := c.Query("type"); t == "invoice" || t == "quote" {
		query = query.Where("type = ?", t)
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

	var items []models.Invoice
	fetch := h.DB.Where("user_id = ?", user.ID)
	if t := c.Query("type"); t == "invoice" || t == "quote" {
		fetch = fetch.Where("type = ?", t)
	}
	fetch.Order("issued_at DESC").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&items)

	c.JSON(http.StatusOK, gin.H{
		"data": models.ToInvoiceResponses(items),
		"meta": gin.H{
			"current_page": page,
			"last_page":    lastPage,
			"per_page":     perPage,
			"total":        total,
		},
	})
}

func (h *InvoiceHandler) Show(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var invoice models.Invoice
	if err := h.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&invoice).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}
	c.JSON(http.StatusOK, models.ToInvoiceResponse(&invoice))
}

func (h *InvoiceHandler) Download(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var invoice models.Invoice
	if err := h.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&invoice).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}
	if invoice.PDFPath == nil || *invoice.PDFPath == "" {
		c.JSON(http.StatusNotFound, gin.H{"message": "PDF indisponible"})
		return
	}

	filename := filepath.Base(*invoice.PDFPath)
	c.FileAttachment(*invoice.PDFPath, filename)
}
