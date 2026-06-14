package handlers

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"upcycleconnect/backend/config"
	"upcycleconnect/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type OAuthHandler struct {
	DB  *gorm.DB
	Cfg *config.Config
}

func (h *OAuthHandler) buildUserWithProfile(user *models.User) map[string]interface{} {
	h.DB.Preload("ProviderProfile").First(user, user.ID)
	resp := models.ToUserResponse(user)
	return map[string]interface{}{
		"id": resp.ID, "email": resp.Email,
		"first_name": resp.FirstName, "last_name": resp.LastName,
		"phone": resp.Phone, "address": resp.Address,
		"avatar_url": resp.AvatarURL, "status": resp.Status,
		"role": resp.Role, "email_verified_at": resp.EmailVerifiedAt,
		"onboarding_completed_at": resp.OnboardingCompletedAt,
		"created_at":              resp.CreatedAt, "updated_at": resp.UpdatedAt,
		"provider_profile": models.ToProviderProfileResponse(user.ProviderProfile),
	}
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
		if user.GoogleID == nil {
			if user.Password == "" {
				c.JSON(http.StatusConflict, gin.H{
					"status":  "oauth_conflict",
					"message": "Un compte existe déjà avec cet email via un autre fournisseur.",
				})
				return
			}
			linkToken, err := h.generateLinkToken("google", info.Sub, email)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur interne."})
				return
			}
			c.JSON(http.StatusConflict, gin.H{
				"status":     "email_conflict",
				"provider":   "google",
				"email":      email,
				"link_token": linkToken,
			})
			return
		}
	}

	tokenStr, err := h.oauthIssueJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur de génération du token."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenStr, "user": h.buildUserWithProfile(&user)})
}

type appleJWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}
type appleJWKS struct {
	Keys []appleJWK `json:"keys"`
}

var (
	appleJWKSCache    *appleJWKS
	appleJWKSCachedAt time.Time
	appleJWKSMu       sync.Mutex
)

func fetchAppleJWKS() (*appleJWKS, error) {
	appleJWKSMu.Lock()
	defer appleJWKSMu.Unlock()
	if appleJWKSCache != nil && time.Since(appleJWKSCachedAt) < time.Hour {
		return appleJWKSCache, nil
	}
	resp, err := http.Get("https://appleid.apple.com/auth/keys")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var jwks appleJWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}
	appleJWKSCache = &jwks
	appleJWKSCachedAt = time.Now()
	return &jwks, nil
}

func verifyAppleIDToken(idToken, servicesID string) (jwt.MapClaims, error) {
	unverified, _, err := new(jwt.Parser).ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}
	kid, _ := unverified.Header["kid"].(string)

	jwks, err := fetchAppleJWKS()
	if err != nil {
		return nil, fmt.Errorf("impossible de récupérer les clés Apple: %w", err)
	}

	var matched *appleJWK
	for i := range jwks.Keys {
		if jwks.Keys[i].Kid == kid {
			matched = &jwks.Keys[i]
			break
		}
	}
	if matched == nil {
		appleJWKSMu.Lock()
		appleJWKSCache = nil
		appleJWKSMu.Unlock()
		jwks, err = fetchAppleJWKS()
		if err != nil {
			return nil, err
		}
		for i := range jwks.Keys {
			if jwks.Keys[i].Kid == kid {
				matched = &jwks.Keys[i]
				break
			}
		}
		if matched == nil {
			return nil, fmt.Errorf("clé Apple introuvable: %s", kid)
		}
	}

	nBytes, err := base64.RawURLEncoding.DecodeString(matched.N)
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(matched.E)
	if err != nil {
		return nil, err
	}
	pubKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: int(new(big.Int).SetBytes(eBytes).Int64()),
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(idToken, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("méthode de signature inattendue: %v", t.Header["alg"])
		}
		return pubKey, nil
	},
		jwt.WithValidMethods([]string{"RS256"}),
	)
	if err != nil || !token.Valid {
		log.Printf("[apple] ParseWithClaims error: %v", err)
		return nil, fmt.Errorf("signature invalide: %w", err)
	}

	iss, _ := claims.GetIssuer()
	if iss != "https://appleid.apple.com" {
		log.Printf("[apple] unexpected issuer: %q", iss)
		return nil, fmt.Errorf("issuer invalide: %q", iss)
	}

	audOK := false
	switch v := claims["aud"].(type) {
	case string:
		audOK = v == servicesID
	case []interface{}:
		for _, a := range v {
			if s, ok := a.(string); ok && s == servicesID {
				audOK = true
				break
			}
		}
	}
	if !audOK {
		log.Printf("[apple] unexpected aud: %v (expected %q)", claims["aud"], servicesID)
		return nil, fmt.Errorf("audience invalide")
	}

	return claims, nil
}

func (h *OAuthHandler) AppleAuth(c *gin.Context) {
	var req struct {
		IDToken   string `json:"id_token" binding:"required"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Token Apple manquant."})
		return
	}

	claims, err := verifyAppleIDToken(req.IDToken, h.Cfg.AppleServicesID)
	if err != nil {
		log.Printf("[apple] verifyAppleIDToken failed (servicesID=%q): %v", h.Cfg.AppleServicesID, err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Token Apple invalide."})
		return
	}

	sub, _ := claims["sub"].(string)
	email, _ := claims["email"].(string)
	if sub == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Informations Apple insuffisantes."})
		return
	}
	email = strings.ToLower(strings.TrimSpace(email))

	var user models.User

	err = h.DB.Where("apple_id = ?", sub).First(&user).Error
	if err == gorm.ErrRecordNotFound && email != "" {
		err = h.DB.Where("email = ?", email).First(&user).Error
	}

	if err == gorm.ErrRecordNotFound {
		now := time.Now()
		firstName := strings.TrimSpace(req.FirstName)
		lastName := strings.TrimSpace(req.LastName)
		if firstName == "" && email != "" {
			firstName = strings.Split(email, "@")[0]
		}
		user = models.User{
			Email:           email,
			Password:        "",
			FirstName:       firstName,
			LastName:        lastName,
			Status:          "active",
			EmailVerifiedAt: &now,
			AppleID:         &sub,
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
		if user.AppleID == nil {
			if user.Password == "" {
				c.JSON(http.StatusConflict, gin.H{
					"status":  "oauth_conflict",
					"message": "Un compte existe déjà avec cet email via un autre fournisseur.",
				})
				return
			}
			linkToken, err := h.generateLinkToken("apple", sub, email)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur interne."})
				return
			}
			c.JSON(http.StatusConflict, gin.H{
				"status":     "email_conflict",
				"provider":   "apple",
				"email":      email,
				"link_token": linkToken,
			})
			return
		}
	}

	tokenStr, err := h.oauthIssueJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur de génération du token."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenStr, "user": h.buildUserWithProfile(&user)})
}

func (h *OAuthHandler) generateLinkToken(provider, sub, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"purpose":  "oauth_link",
		"provider": provider,
		"sub":      sub,
		"email":    email,
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
	})
	return token.SignedString([]byte(h.Cfg.AppKey))
}

func (h *OAuthHandler) OAuthLink(c *gin.Context) {
	var req struct {
		LinkToken string `json:"link_token" binding:"required"`
		Password  string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Données invalides."})
		return
	}

	token, err := jwt.Parse(req.LinkToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("méthode inattendue")
		}
		return []byte(h.Cfg.AppKey), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Lien invalide ou expiré."})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["purpose"] != "oauth_link" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Lien invalide."})
		return
	}

	provider, _ := claims["provider"].(string)
	sub, _ := claims["sub"].(string)
	email, _ := claims["email"].(string)

	var user models.User
	if err := h.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Utilisateur introuvable."})
		return
	}

	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Ce compte ne possède pas de mot de passe."})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Mot de passe incorrect."})
		return
	}

	switch provider {
	case "google":
		if err := h.DB.Model(&user).Update("google_id", sub).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la liaison."})
			return
		}
	case "apple":
		if err := h.DB.Model(&user).Update("apple_id", sub).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la liaison."})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"message": "Fournisseur inconnu."})
		return
	}

	tokenStr, err := h.oauthIssueJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur de génération du token."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenStr, "user": h.buildUserWithProfile(&user)})
}

func (h *OAuthHandler) oauthIssueJWT(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": fmt.Sprintf("%d", userID),
		"exp": time.Now().Add(720 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	})
	return token.SignedString([]byte(h.Cfg.AppKey))
}
