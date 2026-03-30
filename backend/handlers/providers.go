package handlers

import (
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

type ProviderHandler struct {
	DB    *gorm.DB
	Audit *services.AuditService
}

func (h *ProviderHandler) Index(c *gin.Context) {
	query := h.DB.Model(&models.User{}).
		Where("deleted_at IS NULL").
		Where("id IN (SELECT user_id FROM provider_profiles)")

	if status := c.Query("status"); status != "" {
		query = query.Where("id IN (SELECT user_id FROM provider_profiles WHERE status = ?)", status)
	}

	if search := c.Query("search"); search != "" {
		like := "%" + search + "%"
		query = query.Where(
			"(email LIKE ? OR first_name LIKE ? OR last_name LIKE ? OR id IN (SELECT user_id FROM provider_profiles WHERE company_name LIKE ?))",
			like, like, like, like,
		)
	}

	var total int64
	query.Count(&total)

	var pendingCount int64
	var activeCount int64

	h.DB.Table("provider_profiles").Where("status = ?", "pending").Count(&pendingCount)
	h.DB.Table("provider_profiles").Where("status = ?", "approved").Count(&activeCount)

	resPending := h.DB.Table("provider_profiles").Where("status = ?", "pending").Count(&pendingCount)
	if resPending.Error != nil {
		fmt.Println("ERREUR SQL PENDING:", resPending.Error)
	}
	fmt.Printf("COMPTEUR TROUVÉ : %d\n", pendingCount)

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

	fmt.Println("--- DEBUG PROVIDERS INDEX ---")
	fmt.Printf("Status Query: %s\n", c.Query("status"))
	fmt.Printf("Pending Count calculated: %d\n", pendingCount)
	fmt.Printf("Active Count calculated: %d\n", activeCount)
	fmt.Println("-----------------------------")

	var users []models.User
	query.Preload("ProviderProfile").
		Order("created_at DESC").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"data": models.ToProviderResponses(users),
		"meta": gin.H{
			"current_page":  page,
			"last_page":     lastPage,
			"per_page":      perPage,
			"total":         total,
			"pending_count": pendingCount,
			"active_count":  activeCount,
		},
	})
}

func (h *ProviderHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var user models.User
	if err := h.DB.Preload("ProviderProfile").Where("id = ? AND deleted_at IS NULL", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	c.JSON(http.StatusOK, models.ToProviderResponse(&user))
}

func (h *ProviderHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var user models.User
	if err := h.DB.Preload("ProviderProfile").Where("id = ? AND deleted_at IS NULL", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	if user.ProviderProfile == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Profil prestataire introuvable."})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=pending approved suspended"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return
	}

	profile := user.ProviderProfile
	old := map[string]string{"status": profile.Status}

	updates := map[string]interface{}{
		"status": req.Status,
	}

	if req.Status == "approved" && profile.ApprovedAt == nil {
		now := time.Now()
		updates["approved_at"] = now
		authUser := middleware.GetAuthUser(c)
		if authUser != nil {
			updates["approved_by"] = authUser.ID
		}
	}

	h.DB.Model(profile).Updates(updates)

	h.Audit.Log(c, "provider.status_changed", "ProviderProfile", &profile.ID, old, map[string]string{"status": req.Status})

	h.DB.Preload("ProviderProfile").First(&user, user.ID)

	c.JSON(http.StatusOK, models.ToProviderResponse(&user))
}
