package handlers

import (
	"math"
	"net/http"
	"strconv"

	"upcycleconnect/backend/middleware"
	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserDepositHandler struct {
	DB    *gorm.DB
	Audit *services.AuditService
}

func (h *UserDepositHandler) Index(c *gin.Context) {
	user := middleware.GetAuthUser(c)

	query := h.DB.Model(&models.DepositRequest{}).Where("user_id = ?", user.ID)

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

	var items []models.DepositRequest
	fetch := h.DB.Preload("Category").Where("user_id = ?", user.ID)
	if status := c.Query("status"); status != "" {
		fetch = fetch.Where("status = ?", status)
	}
	fetch.Order("created_at DESC").
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

func (h *UserDepositHandler) Show(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var item models.DepositRequest
	if err := h.DB.Preload("Category").
		Where("id = ? AND user_id = ?", id, user.ID).
		First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}
	c.JSON(http.StatusOK, models.ToDepositResponse(&item))
}

func (h *UserDepositHandler) Store(c *gin.Context) {
	user := middleware.GetAuthUser(c)

	var req struct {
		CategoryID      *uint    `json:"category_id"`
		Title           string   `json:"title" binding:"required"`
		Description     string   `json:"description" binding:"required"`
		Condition       string   `json:"condition"`
		History         *string  `json:"history"`
		EstimatedWeight *float64 `json:"estimated_weight"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides. Titre et description requis."})
		return
	}

	if req.Condition == "" {
		req.Condition = "good"
	}
	allowed := map[string]bool{"good": true, "fair": true, "poor": true}
	if !allowed[req.Condition] {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "État invalide (good, fair, poor)."})
		return
	}

	if req.CategoryID != nil {
		var cnt int64
		h.DB.Model(&models.PrestationCategory{}).
			Where("id = ? AND is_active = ? AND deleted_at IS NULL", *req.CategoryID, true).
			Count(&cnt)
		if cnt == 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Catégorie invalide."})
			return
		}
	}

	var carbon *float64
	if req.EstimatedWeight != nil {
		v := *req.EstimatedWeight * 2.5
		carbon = &v
	}

	deposit := models.DepositRequest{
		UserID:          user.ID,
		CategoryID:      req.CategoryID,
		Title:           req.Title,
		Description:     req.Description,
		Condition:       req.Condition,
		History:         req.History,
		EstimatedWeight: req.EstimatedWeight,
		CarbonSavings:   carbon,
		Status:          "pending",
	}
	if err := h.DB.Create(&deposit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création du dépôt."})
		return
	}

	h.Audit.Log(c, "deposit.created", "DepositRequest", &deposit.ID, nil, map[string]interface{}{
		"title": deposit.Title,
	})

	h.DB.Preload("Category").First(&deposit, deposit.ID)
	c.JSON(http.StatusCreated, models.ToDepositResponse(&deposit))
}

func (h *UserDepositHandler) Destroy(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var item models.DepositRequest
	if err := h.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}
	if item.Status != "pending" {
		c.JSON(http.StatusConflict, gin.H{"message": "Seuls les dépôts en attente peuvent être annulés."})
		return
	}

	h.Audit.Log(c, "deposit.cancelled", "DepositRequest", &item.ID,
		map[string]string{"status": item.Status}, map[string]string{"status": "cancelled"})
	h.DB.Delete(&item)
	c.Status(http.StatusNoContent)
}

func (h *UserDepositHandler) Score(c *gin.Context) {
	user := middleware.GetAuthUser(c)

	var deposits []models.DepositRequest
	h.DB.Where("user_id = ?", user.ID).Find(&deposits)

	var totalCount, acceptedCount, pendingCount, rejectedCount int
	var totalWeight, totalCarbon float64
	for _, d := range deposits {
		totalCount++
		switch d.Status {
		case "accepted", "approved", "validated":
			acceptedCount++
			if d.EstimatedWeight != nil {
				totalWeight += *d.EstimatedWeight
			}
			if d.CarbonSavings != nil {
				totalCarbon += *d.CarbonSavings
			}
		case "pending":
			pendingCount++
		case "rejected":
			rejectedCount++
		}
	}

	var regCount int64
	h.DB.Model(&models.EventRegistration{}).Where("user_id = ?", user.ID).Count(&regCount)

	rawScore := float64(acceptedCount)*5 + totalCarbon*1.5 + float64(regCount)*3
	score := int(rawScore)
	if score > 100 {
		score = 100
	}

	level := "Débutant"
	switch {
	case score >= 80:
		level = "Éco-Champion"
	case score >= 50:
		level = "Confirmé"
	case score >= 20:
		level = "Engagé"
	}

	c.JSON(http.StatusOK, gin.H{
		"score":               score,
		"level":               level,
		"deposits_total":      totalCount,
		"deposits_accepted":   acceptedCount,
		"deposits_pending":    pendingCount,
		"deposits_rejected":   rejectedCount,
		"weight_saved_kg":     totalWeight,
		"co2_saved_kg":        totalCarbon,
		"events_attended":     regCount,
	})
}
