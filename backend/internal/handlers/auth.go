package handlers

import (
	"fmt"
	"net/http"
	"time"

	"upcycleconnect/backend/internal/config"
	"upcycleconnect/backend/internal/middleware"
	"upcycleconnect/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB  *gorm.DB
	Cfg *config.Config
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return
	}

	var user models.User
	if err := h.DB.Preload("Roles").Where("email = ? AND deleted_at IS NULL", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Identifiants invalides."})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Identifiants invalides."})
		return
	}

	if !user.IsAdmin() {
		c.JSON(http.StatusForbidden, gin.H{"message": "Accès refusé : compte administrateur requis."})
		return
	}

	if user.Status != "active" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Compte désactivé. Contactez un super-administrateur."})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": fmt.Sprintf("%d", user.ID),
		"exp": time.Now().Add(480 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.Cfg.AppKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création du token."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user":  models.ToUserResponse(&user),
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Déconnexion réussie."})
}

func (h *AuthHandler) Me(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	c.JSON(http.StatusOK, models.ToUserResponse(user))
}
