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

type UserHandler struct {
	DB    *gorm.DB
	Audit *services.AuditService
}

func (h *UserHandler) Index(c *gin.Context) {
	query := h.DB.Model(&models.User{}).
		Where("deleted_at IS NULL").
		Where("id NOT IN (SELECT user_id FROM user_roles)")

	if search := c.Query("search"); search != "" {
		like := "%" + search + "%"
		query = query.Where("(email LIKE ? OR first_name LIKE ? OR last_name LIKE ?)", like, like, like)
	}

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var pendingCount int64
	h.DB.Model(&models.ProviderProfile{}).Where("status = ?", "pending").Count(&pendingCount)

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

	var users []models.User
	query.Preload("Roles").
		Order("created_at DESC").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"data": models.ToUserResponses(users),
		"meta": gin.H{
			"current_page":  page,
			"last_page":     lastPage,
			"per_page":      perPage,
			"total":         total,
			"pending_count": pendingCount,
		},
	})
}

func (h *UserHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var user models.User
	if err := h.DB.Preload("Roles").Preload("ProviderProfile").Where("id = ? AND deleted_at IS NULL", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	c.JSON(http.StatusOK, models.ToUserResponse(&user))
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var user models.User
	if err := h.DB.Preload("Roles").Where("id = ? AND deleted_at IS NULL", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var req struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Phone     *string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return
	}

	updates := map[string]interface{}{}
	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}

	if len(updates) > 0 {
		h.DB.Model(&user).Updates(updates)
		h.DB.Preload("Roles").First(&user, user.ID)
	}

	c.JSON(http.StatusOK, models.ToUserResponse(&user))
}

func (h *UserHandler) Ban(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var user models.User
	if err := h.DB.Preload("Roles").Where("id = ? AND deleted_at IS NULL", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	old := map[string]string{"status": user.Status}
	h.DB.Model(&user).Update("status", "banned")
	user.Status = "banned"

	h.Audit.Log(c, "user.banned", "User", &user.ID, old, map[string]string{"status": "banned"})

	c.JSON(http.StatusOK, models.ToUserResponse(&user))
}

func (h *UserHandler) Activate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var user models.User
	if err := h.DB.Preload("Roles").Where("id = ? AND deleted_at IS NULL", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	old := map[string]string{"status": user.Status}
	h.DB.Model(&user).Update("status", "active")
	user.Status = "active"

	h.Audit.Log(c, "user.activated", "User", &user.ID, old, map[string]string{"status": "active"})

	c.JSON(http.StatusOK, models.ToUserResponse(&user))
}

func (h *UserHandler) Destroy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var user models.User
	if err := h.DB.Where("id = ? AND deleted_at IS NULL", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	h.Audit.Log(c, "user.deleted", "User", &user.ID, map[string]interface{}{"email": user.Email}, nil)
	h.DB.Delete(&user)

	c.Status(http.StatusNoContent)
}
