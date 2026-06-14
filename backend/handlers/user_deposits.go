package handlers

import (
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

type UserDepositHandler struct {
	DB    *gorm.DB
	Audit *services.AuditService
}

const depositPhotosDir = "/uploads/deposits"

func saveDepositPhoto(p *string) *string {
	if p == nil || *p == "" {
		return p
	}
	url, _, err := services.SaveDataURIImage(*p, depositPhotosDir, "/uploads/deposits")
	if err != nil {
		return p
	}
	return &url
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
	if err := h.DB.Preload("Category").Preload("CollectionPoint").
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
		CategoryID        *uint    `json:"category_id"`
		CollectionPointID *uint    `json:"collection_point_id"`
		Title             string   `json:"title" binding:"required"`
		Description       string   `json:"description" binding:"required"`
		Condition         string   `json:"condition"`
		History           *string  `json:"history"`
		EstimatedWeight   *float64 `json:"estimated_weight"`
		SaleType          string   `json:"sale_type"`
		PriceCents        int64    `json:"price_cents"`
		Photo1            *string  `json:"photo1"`
		Photo2            *string  `json:"photo2"`
		Photo3            *string  `json:"photo3"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides. Titre et description requis."})
		return
	}

	if req.Condition == "" {
		req.Condition = "good"
	}
	if req.SaleType != "vente" {
		req.SaleType = "don"
		req.PriceCents = 0
	} else if req.PriceCents <= 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Un prix (price_cents > 0) est requis pour une vente."})
		return
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

	req.Photo1 = saveDepositPhoto(req.Photo1)
	req.Photo2 = saveDepositPhoto(req.Photo2)
	req.Photo3 = saveDepositPhoto(req.Photo3)

	deposit := models.DepositRequest{
		UserID:            user.ID,
		CategoryID:        req.CategoryID,
		CollectionPointID: req.CollectionPointID,
		Title:             req.Title,
		Description:       req.Description,
		Condition:         req.Condition,
		History:           req.History,
		EstimatedWeight:   req.EstimatedWeight,
		CarbonSavings:     carbon,
		Status:            "pending",
		SaleType:          req.SaleType,
		PriceCents:        req.PriceCents,
		Photo1:            req.Photo1,
		Photo2:            req.Photo2,
		Photo3:            req.Photo3,
	}
	if err := h.DB.Create(&deposit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création du dépôt."})
		return
	}

	h.Audit.Log(c, "deposit.created", "DepositRequest", &deposit.ID, nil, map[string]interface{}{
		"title": deposit.Title,
	})

	h.DB.Preload("Category").Preload("CollectionPoint").First(&deposit, deposit.ID)
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
		"score":             score,
		"level":             level,
		"deposits_total":    totalCount,
		"deposits_accepted": acceptedCount,
		"deposits_pending":  pendingCount,
		"deposits_rejected": rejectedCount,
		"weight_saved_kg":   totalWeight,
		"co2_saved_kg":      totalCarbon,
		"events_attended":   regCount,
	})
}

func (h *UserDepositHandler) CollectDeposit(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}

	var req struct {
		Barcode string `json:"barcode" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Code-barres requis"})
		return
	}

	var deposit models.DepositRequest
	if err := h.DB.Preload("User").Preload("Category").Preload("CollectionPoint").
		Where("qr_code = ? AND status = ?", req.Barcode, "approved").
		First(&deposit).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Aucun objet validé trouvé pour ce code-barres. Vérifie que l'objet a bien été approuvé et que le code est correct."})
		return
	}

	if deposit.CollectedAt != nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Cet objet a déjà été collecté.", "deposit": models.ToDepositResponse(&deposit)})
		return
	}

	now := time.Now()
	collectedBy := user.ID
	h.DB.Model(&deposit).Updates(map[string]interface{}{
		"collected_at": &now,
		"collected_by": &collectedBy,
		"status":       "collected",
	})
	deposit.CollectedAt = &now
	deposit.CollectedBy = &collectedBy
	deposit.Status = "collected"

	h.Audit.Log(c, "deposit.collected", "DepositRequest", &deposit.ID, nil, map[string]interface{}{
		"barcode":      req.Barcode,
		"collected_by": user.ID,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Objet collecté avec succès.",
		"deposit": models.ToDepositResponse(&deposit),
	})
}
