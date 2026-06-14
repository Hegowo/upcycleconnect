package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LocaleHandler struct {
	DB    *gorm.DB
	Audit *services.AuditService
}

var localeCodeRe = regexp.MustCompile(`^[a-z]{2}(-[A-Za-z]{2,4})?$`)

func (h *LocaleHandler) PublicList(c *gin.Context) {
	var locales []models.Locale
	h.DB.Where("enabled = ?", true).Order("builtin DESC, name ASC").Find(&locales)
	out := make([]models.LocaleMeta, 0, len(locales))
	for i := range locales {
		out = append(out, models.ToLocaleMeta(&locales[i]))
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

func (h *LocaleHandler) PublicMessages(c *gin.Context) {
	code := c.Param("code")
	var loc models.Locale
	if err := h.DB.Where("code = ? AND enabled = ?", code, true).First(&loc).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Langue introuvable"})
		return
	}
	if strings.TrimSpace(loc.Messages) == "" {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	var tree map[string]interface{}
	if err := json.Unmarshal([]byte(loc.Messages), &tree); err != nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	c.JSON(http.StatusOK, tree)
}

func (h *LocaleHandler) AdminList(c *gin.Context) {
	var locales []models.Locale
	h.DB.Order("builtin DESC, name ASC").Find(&locales)
	out := make([]models.LocaleMeta, 0, len(locales))
	for i := range locales {
		out = append(out, models.ToLocaleMeta(&locales[i]))
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

func (h *LocaleHandler) AdminGet(c *gin.Context) {
	var loc models.Locale
	if err := h.DB.Where("code = ?", c.Param("code")).First(&loc).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Langue introuvable"})
		return
	}
	var tree map[string]interface{}
	if loc.Messages != "" {
		_ = json.Unmarshal([]byte(loc.Messages), &tree)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":     loc.Code,
		"name":     loc.Name,
		"enabled":  loc.Enabled,
		"builtin":  loc.Builtin,
		"messages": tree,
	})
}

type localePayload struct {
	Code     string                 `json:"code"`
	Name     string                 `json:"name"`
	Enabled  *bool                  `json:"enabled"`
	Messages map[string]interface{} `json:"messages"`
}

func (h *LocaleHandler) AdminCreate(c *gin.Context) {
	var req localePayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides"})
		return
	}
	req.Code = strings.TrimSpace(req.Code)
	if !localeCodeRe.MatchString(req.Code) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Code langue invalide (ex: es, de, pt-BR)"})
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Nom de la langue requis"})
		return
	}
	var existing models.Locale
	if h.DB.Where("code = ?", req.Code).First(&existing).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Cette langue existe déjà"})
		return
	}
	msgJSON := "{}"
	if req.Messages != nil {
		if b, err := json.Marshal(req.Messages); err == nil {
			msgJSON = string(b)
		}
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	loc := models.Locale{Code: req.Code, Name: req.Name, Enabled: enabled, Builtin: false, Messages: msgJSON}
	if err := h.DB.Create(&loc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création"})
		return
	}
	h.Audit.Log(c, "locale.created", "Locale", nil, nil, map[string]interface{}{"code": loc.Code, "name": loc.Name})
	c.JSON(http.StatusCreated, models.ToLocaleMeta(&loc))
}

func (h *LocaleHandler) AdminUpdate(c *gin.Context) {
	var loc models.Locale
	if err := h.DB.Where("code = ?", c.Param("code")).First(&loc).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Langue introuvable"})
		return
	}
	var req localePayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides"})
		return
	}
	updates := map[string]interface{}{}
	if strings.TrimSpace(req.Name) != "" {
		updates["name"] = req.Name
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}
	if req.Messages != nil {
		if b, err := json.Marshal(req.Messages); err == nil {
			updates["messages"] = string(b)
		}
	}
	h.DB.Model(&loc).Updates(updates)
	h.Audit.Log(c, "locale.updated", "Locale", nil, nil, map[string]interface{}{"code": loc.Code})
	c.JSON(http.StatusOK, models.ToLocaleMeta(&loc))
}

func (h *LocaleHandler) AdminDelete(c *gin.Context) {
	var loc models.Locale
	if err := h.DB.Where("code = ?", c.Param("code")).First(&loc).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Langue introuvable"})
		return
	}
	if loc.Builtin {
		c.JSON(http.StatusConflict, gin.H{"message": "Les langues intégrées (FR/EN) ne peuvent pas être supprimées"})
		return
	}
	h.DB.Delete(&loc)
	h.Audit.Log(c, "locale.deleted", "Locale", nil, map[string]interface{}{"code": loc.Code}, nil)
	c.JSON(http.StatusOK, gin.H{"message": "Langue supprimée"})
}
