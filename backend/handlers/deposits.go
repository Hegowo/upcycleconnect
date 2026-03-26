package handlers

import (
	"crypto/rand"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"upcycleconnect/backend/middleware"
	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DepositHandler struct {
	DB    *gorm.DB
	Audit *services.AuditService
}

func (h *DepositHandler) Index(c *gin.Context) {
	query := h.DB.Model(&models.DepositRequest{})

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if search := c.Query("search"); search != "" {
		like := "%" + search + "%"
		query = query.Where("title LIKE ? OR description LIKE ?", like, like)
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

	var items []models.DepositRequest
	h.DB.Preload("User").Preload("Category").
		Where(query).
		Order("created_at DESC").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&items)

	c.JSON(http.StatusOK, gin.H{
		"data": models.ToDepositResponses(items),
		"meta": gin.H{
			"current_page": page,
			"last_page":    lastPage,
			"per_page":     perPage,
			"total":        total,
		},
	})
}

func (h *DepositHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var item models.DepositRequest
	if err := h.DB.Preload("User").Preload("Category").First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	c.JSON(http.StatusOK, models.ToDepositResponse(&item))
}

func (h *DepositHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var item models.DepositRequest
	if err := h.DB.Preload("User").Preload("Category").First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var req struct {
		Status         string  `json:"status" binding:"required,oneof=approved rejected"`
		ValidationNote *string `json:"validation_note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return
	}

	old := item.Status
	now := time.Now()
	updates := map[string]interface{}{
		"status":          req.Status,
		"validation_note": req.ValidationNote,
		"validated_at":    now,
	}

	authUser := middleware.GetAuthUser(c)
	if authUser != nil {
		updates["validated_by"] = authUser.ID
	}

	if req.Status == "approved" {
		qr := generateQRCode(item.ID)
		updates["qr_code"] = qr
	}

	h.DB.Model(&item).Updates(updates)

	action := "deposit.approved"
	if req.Status == "rejected" {
		action = "deposit.rejected"
	}
	h.Audit.Log(c, action, "DepositRequest", &item.ID,
		map[string]string{"status": old},
		map[string]string{"status": req.Status},
	)

	h.DB.Preload("User").Preload("Category").First(&item, item.ID)
	c.JSON(http.StatusOK, models.ToDepositResponse(&item))
}

func generateQRCode(id uint) string {
	b := make([]byte, 3)
	rand.Read(b)
	return fmt.Sprintf("UP-%d-%X", id, b)
}
