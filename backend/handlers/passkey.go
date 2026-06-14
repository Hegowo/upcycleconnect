package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"upcycleconnect/backend/config"
	"upcycleconnect/backend/middleware"
	"upcycleconnect/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type PasskeyHandler struct {
	DB  *gorm.DB
	Cfg *config.Config
	WA  *webauthn.WebAuthn
}

func NewPasskeyHandler(db *gorm.DB, cfg *config.Config) (*PasskeyHandler, error) {
	rpID := cfg.WebAuthnRPID
	if rpID == "" {
		rpID = "upcycleconnect.xyz"
	}
	origin := cfg.WebAuthnOrigin
	if origin == "" {
		origin = "https://upcycleconnect.xyz"
	}
	wa, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "UpcycleConnect",
		RPID:          rpID,
		RPOrigins:     []string{origin},
	})
	if err != nil {
		return nil, err
	}
	return &PasskeyHandler{DB: db, Cfg: cfg, WA: wa}, nil
}

type sessionEntry struct {
	data      *webauthn.SessionData
	expiresAt time.Time
}

var (
	sessions   = make(map[string]sessionEntry)
	sessionsMu sync.Mutex
)

func saveSession(key string, data *webauthn.SessionData) {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()

	for k, v := range sessions {
		if time.Now().After(v.expiresAt) {
			delete(sessions, k)
		}
	}
	sessions[key] = sessionEntry{data: data, expiresAt: time.Now().Add(5 * time.Minute)}
}

func loadSession(key string) (*webauthn.SessionData, bool) {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	e, ok := sessions[key]
	if !ok || time.Now().After(e.expiresAt) {
		delete(sessions, key)
		return nil, false
	}
	delete(sessions, key)
	return e.data, true
}

type waUser struct {
	user        *models.User
	credentials []webauthn.Credential
}

func (u *waUser) WebAuthnID() []byte                         { return []byte(fmt.Sprintf("%d", u.user.ID)) }
func (u *waUser) WebAuthnName() string                       { return u.user.Email }
func (u *waUser) WebAuthnDisplayName() string                { return u.user.FirstName + " " + u.user.LastName }
func (u *waUser) WebAuthnCredentials() []webauthn.Credential { return u.credentials }

func (h *PasskeyHandler) loadWAUser(userID uint) (*waUser, error) {
	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	var passkeys []models.Passkey
	h.DB.Where("user_id = ?", userID).Find(&passkeys)

	creds := make([]webauthn.Credential, 0, len(passkeys))
	for _, pk := range passkeys {
		var cred webauthn.Credential
		if err := json.Unmarshal(pk.Credential, &cred); err == nil {
			creds = append(creds, cred)
		}
	}
	return &waUser{user: &user, credentials: creds}, nil
}

func (h *PasskeyHandler) RegisterBegin(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié."})
		return
	}
	wu, err := h.loadWAUser(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Utilisateur introuvable."})
		return
	}

	options, sessionData, err := h.WA.BeginRegistration(wu,
		webauthn.WithAuthenticatorSelection(protocol.AuthenticatorSelection{
			ResidentKey:      protocol.ResidentKeyRequirementPreferred,
			UserVerification: protocol.VerificationPreferred,
		}),
	)
	if err != nil {
		log.Printf("[passkey] BeginRegistration error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la génération du challenge."})
		return
	}

	sessionKey := fmt.Sprintf("reg:%d", user.ID)
	saveSession(sessionKey, sessionData)
	c.JSON(http.StatusOK, options)
}

func (h *PasskeyHandler) RegisterComplete(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié."})
		return
	}

	name := c.Query("name")
	if name == "" {
		name = "Ma clé d'accès"
	}

	wu, err := h.loadWAUser(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Utilisateur introuvable."})
		return
	}

	sessionKey := fmt.Sprintf("reg:%d", user.ID)
	sessionData, ok := loadSession(sessionKey)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Session expirée, veuillez réessayer."})
		return
	}

	credential, err := h.WA.FinishRegistration(wu, *sessionData, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Vérification échouée : " + err.Error()})
		return
	}

	credJSON, err := json.Marshal(credential)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur interne."})
		return
	}

	passkey := models.Passkey{
		UserID:       user.ID,
		CredentialID: base64.RawURLEncoding.EncodeToString(credential.ID),
		Credential:   credJSON,
		Name:         name,
	}
	if err := h.DB.Create(&passkey).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de l'enregistrement."})
		return
	}

	c.JSON(http.StatusCreated, passkey)
}

func (h *PasskeyHandler) List(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié."})
		return
	}
	var passkeys []models.Passkey
	h.DB.Where("user_id = ?", user.ID).Order("created_at desc").Find(&passkeys)
	c.JSON(http.StatusOK, passkeys)
}

func (h *PasskeyHandler) Delete(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié."})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID invalide."})
		return
	}

	result := h.DB.Where("id = ? AND user_id = ?", id, user.ID).Delete(&models.Passkey{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Clé introuvable."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Clé supprimée."})
}

func (h *PasskeyHandler) Rename(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié."})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID invalide."})
		return
	}
	var body struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Nom requis."})
		return
	}
	result := h.DB.Model(&models.Passkey{}).Where("id = ? AND user_id = ?", id, user.ID).Update("name", body.Name)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Clé introuvable."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Clé renommée."})
}

const authSessionKey = "auth:login"

func (h *PasskeyHandler) AuthBegin(c *gin.Context) {
	options, sessionData, err := h.WA.BeginDiscoverableLogin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la génération du challenge."})
		return
	}

	challenge := base64.RawURLEncoding.EncodeToString([]byte(sessionData.Challenge))
	saveSession("login:"+challenge, sessionData)
	c.JSON(http.StatusOK, gin.H{"options": options, "session_key": challenge})
}

func (h *PasskeyHandler) AuthComplete(c *gin.Context) {
	var body struct {
		SessionKey string `json:"session_key" binding:"required"`
	}

	sessionKey := c.Query("session_key")
	if sessionKey == "" {
		if err := c.ShouldBindJSON(&body); err != nil || body.SessionKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "session_key manquant."})
			return
		}
		sessionKey = body.SessionKey
	}

	sessionData, ok := loadSession("login:" + sessionKey)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Session expirée. Veuillez réessayer."})
		return
	}

	userHandler := func(rawID, userHandle []byte) (webauthn.User, error) {
		rawIDStr := base64.RawURLEncoding.EncodeToString(rawID)
		var passkey models.Passkey
		if err := h.DB.Where("credential_id = ?", rawIDStr).First(&passkey).Error; err != nil {
			return nil, fmt.Errorf("clé inconnue")
		}
		wu, err := h.loadWAUser(passkey.UserID)
		if err != nil {
			return nil, err
		}
		return wu, nil
	}

	credential, err := h.WA.FinishDiscoverableLogin(userHandler, *sessionData, c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Vérification échouée."})
		return
	}

	rawIDStr := base64.RawURLEncoding.EncodeToString(credential.ID)
	var passkey models.Passkey
	if err := h.DB.Where("credential_id = ?", rawIDStr).First(&passkey).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Clé inconnue."})
		return
	}

	now := time.Now()
	credJSON, _ := json.Marshal(credential)
	h.DB.Model(&passkey).Updates(map[string]interface{}{
		"credential":   credJSON,
		"last_used_at": now,
	})

	var user models.User
	if err := h.DB.First(&user, passkey.UserID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Utilisateur introuvable."})
		return
	}
	if user.Status == "banned" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Ce compte a été suspendu."})
		return
	}

	tokenStr, err := h.issueJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur de génération du token."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenStr, "user": models.ToUserResponse(&user)})
}

func (h *PasskeyHandler) issueJWT(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": fmt.Sprintf("%d", userID),
		"exp": time.Now().Add(720 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	})
	return token.SignedString([]byte(h.Cfg.AppKey))
}
