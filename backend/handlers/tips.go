package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"upcycleconnect/backend/middleware"
	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TipHandler struct {
	DB    *gorm.DB
	Audit *services.AuditService
}

func (h *TipHandler) PublicIndex(c *gin.Context) {
	q := h.DB.Preload("Author").
		Where("status = ?", "published").
		Where("published_at IS NULL OR published_at <= ?", time.Now()).
		Order("published_at DESC, created_at DESC")

	if cat := c.Query("category"); cat != "" {
		q = q.Where("category = ?", cat)
	}
	if search := strings.TrimSpace(c.Query("search")); search != "" {
		like := "%" + search + "%"
		q = q.Where("title LIKE ? OR summary LIKE ?", like, like)
	}

	var tips []models.Tip
	q.Find(&tips)

	resp := make([]models.TipResponse, 0, len(tips))
	for i := range tips {
		resp = append(resp, models.ToTipResponse(&tips[i], false))
	}

	categorySet := map[string]struct{}{}
	for _, t := range tips {
		if t.Category != "" {
			categorySet[t.Category] = struct{}{}
		}
	}
	cats := make([]string, 0, len(categorySet))
	for c := range categorySet {
		cats = append(cats, c)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       resp,
		"categories": cats,
	})
}

func (h *TipHandler) PublicShow(c *gin.Context) {
	slugOrID := c.Param("slug")
	var tip models.Tip
	q := h.DB.Preload("Author").
		Where("status = ? AND (published_at IS NULL OR published_at <= ?)", "published", time.Now())
	if id, err := strconv.ParseUint(slugOrID, 10, 64); err == nil {
		q = q.Where("id = ?", id)
	} else {
		q = q.Where("slug = ?", slugOrID)
	}
	if err := q.First(&tip).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Conseil introuvable"})
		return
	}

	h.DB.Model(&tip).UpdateColumn("view_count", gorm.Expr("view_count + 1"))
	tip.ViewCount++
	c.JSON(http.StatusOK, models.ToTipResponse(&tip, true))
}

type tipPayload struct {
	Title       string  `json:"title" binding:"required,min=3,max=200"`
	Slug        string  `json:"slug"`
	Summary     string  `json:"summary"`
	Content     string  `json:"content" binding:"required,min=10"`
	CoverImage  *string `json:"cover_image"`
	Category    string  `json:"category"`
	Status      string  `json:"status"`
	PublishedAt *string `json:"published_at"`
}

func (h *TipHandler) AdminIndex(c *gin.Context) {
	var tips []models.Tip
	q := h.DB.Preload("Author").Order("created_at DESC")
	if status := c.Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}
	q.Find(&tips)
	resp := make([]models.TipResponse, 0, len(tips))
	for i := range tips {
		resp = append(resp, models.ToTipResponse(&tips[i], false))
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *TipHandler) AdminShow(c *gin.Context) {
	var tip models.Tip
	if err := h.DB.Preload("Author").First(&tip, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Conseil introuvable"})
		return
	}
	c.JSON(http.StatusOK, models.ToTipResponse(&tip, true))
}

func (h *TipHandler) AdminStore(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	var req tipPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides", "error": err.Error()})
		return
	}
	slug := normalizeSlug(req.Slug, req.Title)
	tip := models.Tip{
		Title:      req.Title,
		Slug:       uniqueSlug(h.DB, slug, 0),
		Summary:    req.Summary,
		Content:    req.Content,
		CoverImage: req.CoverImage,
		Category:   req.Category,
		Status:     defaultIfEmpty(req.Status, "draft"),
		AuthorID:   user.ID,
	}
	tip.PublishedAt = computePublishedAt(req.PublishedAt, tip.Status)

	if err := h.DB.Create(&tip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création"})
		return
	}
	h.Audit.Log(c, "tip.created", "Tip", &tip.ID, nil, map[string]interface{}{"title": tip.Title})
	h.DB.Preload("Author").First(&tip, tip.ID)
	c.JSON(http.StatusCreated, models.ToTipResponse(&tip, true))
}

func (h *TipHandler) AdminUpdate(c *gin.Context) {
	var tip models.Tip
	if err := h.DB.First(&tip, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Conseil introuvable"})
		return
	}
	var req tipPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides", "error": err.Error()})
		return
	}
	newSlug := normalizeSlug(req.Slug, req.Title)
	if newSlug != tip.Slug {
		tip.Slug = uniqueSlug(h.DB, newSlug, tip.ID)
	}
	tip.Title = req.Title
	tip.Summary = req.Summary
	tip.Content = req.Content
	tip.CoverImage = req.CoverImage
	tip.Category = req.Category
	tip.Status = defaultIfEmpty(req.Status, tip.Status)
	tip.PublishedAt = computePublishedAt(req.PublishedAt, tip.Status)

	if err := h.DB.Save(&tip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la mise à jour"})
		return
	}
	h.Audit.Log(c, "tip.updated", "Tip", &tip.ID, nil, map[string]interface{}{"title": tip.Title})
	h.DB.Preload("Author").First(&tip, tip.ID)
	c.JSON(http.StatusOK, models.ToTipResponse(&tip, true))
}

func (h *TipHandler) AdminDestroy(c *gin.Context) {
	var tip models.Tip
	if err := h.DB.First(&tip, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Conseil introuvable"})
		return
	}
	if err := h.DB.Delete(&tip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la suppression"})
		return
	}
	h.Audit.Log(c, "tip.deleted", "Tip", &tip.ID, map[string]string{"title": tip.Title}, nil)
	c.JSON(http.StatusOK, gin.H{"message": "Conseil supprimé"})
}

var slugRe = regexp.MustCompile(`[^a-z0-9]+`)

func normalizeSlug(slug, fallback string) string {
	if strings.TrimSpace(slug) == "" {
		slug = fallback
	}
	slug = strings.ToLower(slug)
	slug = slugRe.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = fmt.Sprintf("conseil-%d", time.Now().Unix())
	}
	if len(slug) > 200 {
		slug = slug[:200]
	}
	return slug
}

func uniqueSlug(db *gorm.DB, base string, ignoreID uint) string {
	candidate := base
	i := 2
	for {
		var count int64
		q := db.Model(&models.Tip{}).Where("slug = ?", candidate)
		if ignoreID > 0 {
			q = q.Where("id <> ?", ignoreID)
		}
		q.Count(&count)
		if count == 0 {
			return candidate
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
		i++
	}
}

func defaultIfEmpty(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return v
}

func computePublishedAt(raw *string, status string) *time.Time {
	if status != "published" {
		return nil
	}
	if raw != nil && strings.TrimSpace(*raw) != "" {
		if t, err := time.Parse(time.RFC3339, *raw); err == nil {
			return &t
		}
	}
	now := time.Now()
	return &now
}
