package handlers

import (
	"net/http"
	"strconv"
	"time"

	"upcycleconnect/backend/middleware"
	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminHandler struct {
	DB    *gorm.DB
	Audit *services.AuditService
}

func (h *AdminHandler) Index(c *gin.Context) {
	var users []models.User
	h.DB.Preload("Roles").
		Where("deleted_at IS NULL").
		Where("id IN (SELECT user_id FROM user_roles WHERE role_id IN (SELECT id FROM roles WHERE name IN ('admin','super_admin')))").
		Order("created_at DESC").
		Find(&users)

	c.JSON(http.StatusOK, gin.H{"data": models.ToUserResponses(users)})
}

func (h *AdminHandler) Show(c *gin.Context) {
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

	c.JSON(http.StatusOK, models.ToUserResponse(&user))
}

func (h *AdminHandler) Store(c *gin.Context) {
	var req struct {
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required,min=8"`
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
		Role      string `json:"role" binding:"required,oneof=admin super_admin"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return
	}

	var existing models.User
	if err := h.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Cet email est déjà utilisé."})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur interne."})
		return
	}

	now := time.Now()
	user := models.User{
		Email:           req.Email,
		Password:        string(hashedPassword),
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Status:          "active",
		EmailVerifiedAt: &now,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création."})
		return
	}

	var role models.Role
	if err := h.DB.Where("name = ?", req.Role).First(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Rôle introuvable."})
		return
	}

	h.DB.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", user.ID, role.ID)

	h.Audit.Log(c, "admin.created", "User", &user.ID, nil, map[string]interface{}{
		"email": user.Email, "role": req.Role,
	})

	h.DB.Preload("Roles").First(&user, user.ID)

	c.JSON(http.StatusCreated, models.ToUserResponse(&user))
}

func (h *AdminHandler) Update(c *gin.Context) {
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
		Password  *string `json:"password"`
		Role      *string `json:"role"`
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
	if req.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur interne."})
			return
		}
		updates["password"] = string(hashedPassword)
	}

	if len(updates) > 0 {
		h.DB.Model(&user).Updates(updates)
	}

	if req.Role != nil {
		var role models.Role
		if err := h.DB.Where("name = ?", *req.Role).First(&role).Error; err == nil {
			h.DB.Exec("DELETE FROM user_roles WHERE user_id = ?", user.ID)
			h.DB.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", user.ID, role.ID)
		}
	}

	h.DB.Preload("Roles").First(&user, user.ID)

	c.JSON(http.StatusOK, models.ToUserResponse(&user))
}

func (h *AdminHandler) Destroy(c *gin.Context) {
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

	authUser := middleware.GetAuthUser(c)
	if authUser != nil && authUser.ID == user.ID {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Vous ne pouvez pas supprimer votre propre compte."})
		return
	}

	h.Audit.Log(c, "admin.deleted", "User", &user.ID, map[string]interface{}{"email": user.Email}, nil)
	h.DB.Delete(&user)

	c.Status(http.StatusNoContent)
}
