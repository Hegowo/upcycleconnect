package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"upcycleconnect/backend/config"
	"upcycleconnect/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type OAuthHandler struct {
	DB  *gorm.DB
	Cfg *config.Config
}

type googleTokenInfo struct {
	Sub              string `json:"sub"`
	Email            string `json:"email"`
	EmailVerified    string `json:"email_verified"`
	GivenName        string `json:"given_name"`
	FamilyName       string `json:"family_name"`
	Picture          string `json:"picture"`
	Aud              string `json:"aud"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (h *OAuthHandler) GoogleAuth(c *gin.Context) {
	var req struct {
		Credential string `json:"credential" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Token Google manquant."})
		return
	}

	resp, err := http.Get(fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", req.Credential))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "Impossible de vérifier le token Google."})
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var info googleTokenInfo
	if err := json.Unmarshal(body, &info); err != nil || info.Error != "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Token Google invalide."})
		return
	}
	if info.Email == "" || info.EmailVerified != "true" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Email Google non vérifié."})
		return
	}
	if info.Aud != h.Cfg.GoogleClientID {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Token Google invalide."})
		return
	}

	email := strings.ToLower(strings.TrimSpace(info.Email))

	var user models.User
	err = h.DB.Where("email = ?", email).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		now := time.Now()
		user = models.User{
			Email:           email,
			Password:        "",
			FirstName:       info.GivenName,
			LastName:        info.FamilyName,
			Status:          "active",
			EmailVerifiedAt: &now,
			GoogleID:        &info.Sub,
		}
		if info.Picture != "" {
			user.AvatarURL = &info.Picture
		}
		if err := h.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création du compte."})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur interne."})
		return
	} else {
		if user.Status == "banned" {
			c.JSON(http.StatusForbidden, gin.H{"message": "Ce compte a été suspendu."})
			return
		}
		// Lier le Google ID si pas encore fait
		if user.GoogleID == nil {
			googleID := info.Sub
			h.DB.Model(&user).Update("google_id", googleID)
		}
	}

	tokenStr, err := h.oauthIssueJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur de génération du token."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenStr,
		"user":  models.ToUserResponse(&user),
	})
}

func (h *OAuthHandler) oauthIssueJWT(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": fmt.Sprintf("%d", userID),
		"exp": time.Now().Add(720 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	})
	return token.SignedString([]byte(h.Cfg.AppKey))
}
