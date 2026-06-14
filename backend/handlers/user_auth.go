package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"upcycleconnect/backend/config"
	"upcycleconnect/backend/middleware"
	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserAuthHandler struct {
	DB     *gorm.DB
	Cfg    *config.Config
	Mailer *services.Mailer
}

func (h *UserAuthHandler) Register(c *gin.Context) {
	var req struct {
		FirstName   string  `json:"first_name" binding:"required"`
		LastName    string  `json:"last_name" binding:"required"`
		Email       string  `json:"email" binding:"required,email"`
		Password    string  `json:"password" binding:"required,min=8"`
		Phone       *string `json:"phone"`
		AccountType string  `json:"account_type"`
		CompanyName *string `json:"company_name"`
		Siret       *string `json:"siret"`
		Activity    *string `json:"activity"`
		Address     *string `json:"address"`
	}

	var kbisHeader *multipart.FileHeader
	var kbisFile multipart.File
	if strings.HasPrefix(c.ContentType(), "multipart/form-data") {
		req.FirstName = c.PostForm("first_name")
		req.LastName = c.PostForm("last_name")
		req.Email = c.PostForm("email")
		req.Password = c.PostForm("password")
		req.AccountType = c.PostForm("account_type")
		if v := c.PostForm("phone"); v != "" {
			req.Phone = &v
		}
		if v := c.PostForm("company_name"); v != "" {
			req.CompanyName = &v
		}
		if v := c.PostForm("siret"); v != "" {
			req.Siret = &v
		}
		if v := c.PostForm("activity"); v != "" {
			req.Activity = &v
		}
		if v := c.PostForm("address"); v != "" {
			req.Address = &v
		}
		if req.FirstName == "" || req.LastName == "" || req.Email == "" || len(req.Password) < 8 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides. Vérifiez les champs obligatoires."})
			return
		}
		if f, header, err := c.Request.FormFile("kbis"); err == nil {
			kbisFile, kbisHeader = f, header
			defer kbisFile.Close()
		}
	} else {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides. Vérifiez les champs obligatoires."})
			return
		}
	}

	var count int64
	h.DB.Model(&models.User{}).Where("email = ? AND deleted_at IS NULL", req.Email).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"message": "Cette adresse email est déjà utilisée."})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur serveur."})
		return
	}

	user := models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  string(hashed),
		Phone:     req.Phone,
		Address:   req.Address,
		Status:    "pending",
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création du compte."})
		return
	}

	if req.AccountType == "provider" && req.CompanyName != nil && *req.CompanyName != "" {
		profile := models.ProviderProfile{
			UserID:      user.ID,
			CompanyName: *req.CompanyName,
			Siret:       req.Siret,
			Description: req.Activity,
			Status:      "pending",
		}
		h.DB.Create(&profile)
		if kbisHeader != nil {
			if path, err := saveKbisFile(kbisHeader, kbisFile, user.ID); err == nil {
				h.DB.Model(&profile).Update("kbis_path", path)
			} else {
				log.Printf("[register] kbis save failed for user %d: %v", user.ID, err)
			}
		}
	}

	token := newToken()
	expires := time.Now().Add(24 * time.Hour)
	ev := models.EmailVerification{
		UserID:    user.ID,
		Token:     token,
		Type:      "register",
		ExpiresAt: expires,
	}
	h.DB.Create(&ev)

	if err := h.Mailer.SendRegisterVerification(user.Email, user.FirstName, token); err != nil {
		log.Printf("[mailer] register verification to %s: %v", user.Email, err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Compte créé. Vérifiez votre email pour activer votre compte.",
	})
}

func (h *UserAuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return
	}

	var user models.User
	if err := h.DB.Where("email = ? AND deleted_at IS NULL", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Identifiants invalides."})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Identifiants invalides."})
		return
	}

	if user.Status == "pending" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Votre compte n'est pas encore activé. Vérifiez votre email."})
		return
	}

	if user.Status != "active" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Votre compte a été désactivé."})
		return
	}

	clientIP := c.ClientIP()

	var knownIP models.UserKnownIP
	isKnown := h.DB.Where("user_id = ? AND ip_address = ?", user.ID, clientIP).First(&knownIP).Error == nil

	if isKnown {
		tokenStr, err := h.issueJWT(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création du token."})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"token": tokenStr,
			"user":  h.userWithProfile(&user),
		})
		return
	}

	token := newToken()
	expires := time.Now().Add(15 * time.Minute)
	ip := clientIP
	ev := models.EmailVerification{
		UserID:    user.ID,
		Token:     token,
		Type:      "login",
		IPAddress: &ip,
		ExpiresAt: expires,
	}
	h.DB.Create(&ev)

	attemptTime := time.Now().Format("02/01/2006 à 15:04:05 (UTC)")
	if err := h.Mailer.SendLoginVerification(user.Email, user.FirstName, token, clientIP, attemptTime); err != nil {
		log.Printf("[mailer] login verification to %s: %v", user.Email, err)
	}

	c.JSON(http.StatusAccepted, gin.H{
		"pending_verification": true,
		"message":              "Un email de vérification a été envoyé à votre adresse.",
	})
}

func (h *UserAuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Token manquant."})
		return
	}

	var ev models.EmailVerification
	if err := h.DB.Where("token = ? AND type = ?", token, "register").First(&ev).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Lien invalide ou expiré."})
		return
	}
	if ev.UsedAt != nil {
		c.JSON(http.StatusGone, gin.H{"message": "Ce lien a déjà été utilisé."})
		return
	}
	if time.Now().After(ev.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{"message": "Ce lien a expiré."})
		return
	}

	now := time.Now()
	h.DB.Model(&models.User{}).Where("id = ?", ev.UserID).Updates(map[string]interface{}{
		"status":            "active",
		"email_verified_at": now,
	})

	h.DB.Model(&ev).Update("used_at", now)

	clientIP := c.ClientIP()
	h.DB.Exec(
		"INSERT IGNORE INTO user_known_ips (user_id, ip_address, created_at) VALUES (?, ?, ?)",
		ev.UserID, clientIP, now,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Email vérifié. Votre compte est activé."})
}

func (h *UserAuthHandler) VerifyLogin(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Token manquant."})
		return
	}

	var ev models.EmailVerification
	if err := h.DB.Where("token = ? AND type = ?", token, "login").First(&ev).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Lien invalide ou expiré."})
		return
	}
	if ev.UsedAt != nil {
		c.JSON(http.StatusGone, gin.H{"message": "Ce lien a déjà été utilisé."})
		return
	}
	if time.Now().After(ev.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{"message": "Ce lien a expiré (15 minutes). Reconnectez-vous."})
		return
	}

	var user models.User
	if err := h.DB.Where("id = ? AND status = 'active' AND deleted_at IS NULL", ev.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Compte introuvable ou désactivé."})
		return
	}

	now := time.Now()

	if ev.IPAddress != nil {
		h.DB.Exec(
			"INSERT IGNORE INTO user_known_ips (user_id, ip_address, created_at) VALUES (?, ?, ?)",
			ev.UserID, *ev.IPAddress, now,
		)
	}
	clientIP := c.ClientIP()
	h.DB.Exec(
		"INSERT IGNORE INTO user_known_ips (user_id, ip_address, created_at) VALUES (?, ?, ?)",
		ev.UserID, clientIP, now,
	)

	h.DB.Model(&ev).Update("used_at", now)

	tokenStr, err := h.issueJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création du token."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenStr,
		"user":  h.userWithProfile(&user),
	})
}

func (h *UserAuthHandler) userWithProfile(user *models.User) map[string]interface{} {
	h.DB.Preload("ProviderProfile").First(user, user.ID)
	resp := models.ToUserResponse(user)
	return map[string]interface{}{
		"id":                      resp.ID,
		"email":                   resp.Email,
		"first_name":              resp.FirstName,
		"last_name":               resp.LastName,
		"phone":                   resp.Phone,
		"address":                 resp.Address,
		"avatar_url":              resp.AvatarURL,
		"status":                  resp.Status,
		"role":                    resp.Role,
		"email_verified_at":       resp.EmailVerifiedAt,
		"onboarding_completed_at": resp.OnboardingCompletedAt,
		"created_at":              resp.CreatedAt,
		"updated_at":              resp.UpdatedAt,
		"provider_profile":        models.ToProviderProfileResponse(user.ProviderProfile),
	}
}

func (h *UserAuthHandler) Me(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	c.JSON(http.StatusOK, h.userWithProfile(user))
}

func (h *UserAuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Déconnexion réussie."})
}

func (h *UserAuthHandler) CompleteOnboarding(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	if user.OnboardingCompletedAt != nil {
		c.JSON(http.StatusOK, h.userWithProfile(user))
		return
	}
	now := time.Now()
	if err := h.DB.Model(user).Update("onboarding_completed_at", now).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de l'enregistrement"})
		return
	}
	user.OnboardingCompletedAt = &now
	c.JSON(http.StatusOK, h.userWithProfile(user))
}

func (h *UserAuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides (mot de passe de 8 caractères minimum)."})
		return
	}

	var ev models.EmailVerification
	if err := h.DB.Where("token = ? AND type = ?", req.Token, "password_reset").First(&ev).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Lien invalide ou expiré."})
		return
	}
	if ev.UsedAt != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Ce lien a déjà été utilisé."})
		return
	}
	if time.Now().After(ev.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Ce lien a expiré."})
		return
	}

	var user models.User
	if err := h.DB.First(&user, ev.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Utilisateur introuvable."})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur interne."})
		return
	}

	now := time.Now()
	h.DB.Model(&user).Updates(map[string]interface{}{"password": string(hash)})
	h.DB.Model(&ev).Update("used_at", now)

	c.JSON(http.StatusOK, gin.H{"message": "Mot de passe réinitialisé avec succès."})
}

func (h *UserAuthHandler) issueJWT(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": fmt.Sprintf("%d", userID),
		"exp": time.Now().Add(720 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	})
	return token.SignedString([]byte(h.Cfg.AppKey))
}

func newToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
