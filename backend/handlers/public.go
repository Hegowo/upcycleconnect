package handlers

import (
	"math"
	"net/http"
	"strconv"

	"upcycleconnect/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PublicHandler struct {
	DB *gorm.DB
}

func (h *PublicHandler) Categories(c *gin.Context) {
	type row struct {
		models.PrestationCategory
		PrestationsCount int64 `gorm:"column:prestations_count"`
	}
	var rows []row
	h.DB.Model(&models.PrestationCategory{}).
		Select("prestation_categories.*, (SELECT COUNT(*) FROM prestations WHERE prestations.category_id = prestation_categories.id AND prestations.status = 'published' AND prestations.deleted_at IS NULL) as prestations_count").
		Where("prestation_categories.deleted_at IS NULL AND prestation_categories.is_active = ?", true).
		Order("sort_order ASC, name ASC").
		Find(&rows)

	result := make([]models.CategoryResponse, len(rows))
	for i, r := range rows {
		result[i] = models.ToCategoryResponse(&r.PrestationCategory, r.PrestationsCount)
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *PublicHandler) Prestations(c *gin.Context) {
	query := h.DB.Model(&models.Prestation{}).
		Where("prestations.status = ? AND prestations.deleted_at IS NULL", "published")

	if catID := c.Query("category_id"); catID != "" {
		query = query.Where("category_id = ?", catID)
	}
	if search := c.Query("search"); search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if perPage < 1 || perPage > 100 {
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

	var prestations []models.Prestation
	fetch := h.DB.Preload("Category").Preload("Provider").
		Where("prestations.status = ? AND prestations.deleted_at IS NULL", "published")
	if catID := c.Query("category_id"); catID != "" {
		fetch = fetch.Where("category_id = ?", catID)
	}
	if search := c.Query("search"); search != "" {
		fetch = fetch.Where("title LIKE ?", "%"+search+"%")
	}
	fetch.Order("created_at DESC").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&prestations)

	c.JSON(http.StatusOK, gin.H{
		"data": models.ToPrestationResponses(prestations),
		"meta": gin.H{
			"current_page": page,
			"last_page":    lastPage,
			"per_page":     perPage,
			"total":        total,
		},
	})
}

func (h *PublicHandler) ShowPrestation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}
	var p models.Prestation
	if err := h.DB.Preload("Category").Preload("Provider").
		Where("id = ? AND status = ? AND deleted_at IS NULL", id, "published").
		First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}
	c.JSON(http.StatusOK, models.ToPrestationResponse(&p))
}

func (h *PublicHandler) Events(c *gin.Context) {
	query := h.DB.Model(&models.Event{}).
		Where("platform_events.status = ? AND platform_events.deleted_at IS NULL", "published")

	if from := c.Query("start_date_from"); from != "" {
		query = query.Where("start_date >= ?", from)
	}
	if to := c.Query("start_date_to"); to != "" {
		query = query.Where("start_date <= ?", to)
	}
	if upcoming := c.Query("upcoming"); upcoming == "true" || upcoming == "1" {
		query = query.Where("end_date >= NOW()")
	}

	var total int64
	query.Count(&total)

	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if perPage < 1 || perPage > 100 {
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

	var events []models.Event
	fetch := h.DB.Preload("Category").
		Where("platform_events.status = ? AND platform_events.deleted_at IS NULL", "published")
	if from := c.Query("start_date_from"); from != "" {
		fetch = fetch.Where("start_date >= ?", from)
	}
	if to := c.Query("start_date_to"); to != "" {
		fetch = fetch.Where("start_date <= ?", to)
	}
	if upcoming := c.Query("upcoming"); upcoming == "true" || upcoming == "1" {
		fetch = fetch.Where("end_date >= NOW()")
	}
	fetch.Order("start_date ASC").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&events)

	type regCount struct {
		EventID uint
		Count   int64
	}
	counts := make(map[uint]int64)
	if len(events) > 0 {
		ids := make([]uint, len(events))
		for i, e := range events {
			ids[i] = e.ID
		}
		var rows []regCount
		h.DB.Model(&models.EventRegistration{}).
			Select("event_id, COUNT(*) as count").
			Where("event_id IN ?", ids).
			Group("event_id").
			Scan(&rows)
		for _, r := range rows {
			counts[r.EventID] = r.Count
		}
	}

	responses := make([]models.EventResponse, len(events))
	for i := range events {
		responses[i] = models.ToEventResponseWithCount(&events[i], counts[events[i].ID])
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responses,
		"meta": gin.H{
			"current_page": page,
			"last_page":    lastPage,
			"per_page":     perPage,
			"total":        total,
		},
	})
}

func (h *PublicHandler) ShowEvent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}
	var e models.Event
	if err := h.DB.Preload("Category").Preload("Creator").
		Where("id = ? AND status = ? AND deleted_at IS NULL", id, "published").
		First(&e).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}
	var count int64
	h.DB.Model(&models.EventRegistration{}).Where("event_id = ?", e.ID).Count(&count)
	c.JSON(http.StatusOK, models.ToEventResponseWithCount(&e, count))
}

func (h *PublicHandler) Providers(c *gin.Context) {
	query := h.DB.Model(&models.User{}).
		Where("users.deleted_at IS NULL AND users.status = ?", "active").
		Where("users.id IN (SELECT user_id FROM provider_profiles WHERE status = 'approved')")

	if search := c.Query("search"); search != "" {
		like := "%" + search + "%"
		query = query.Where(
			"(users.id IN (SELECT user_id FROM provider_profiles WHERE company_name LIKE ?))",
			like,
		)
	}

	var total int64
	query.Count(&total)

	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if perPage < 1 || perPage > 100 {
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
	fetch := h.DB.Preload("ProviderProfile").
		Where("users.deleted_at IS NULL AND users.status = ?", "active").
		Where("users.id IN (SELECT user_id FROM provider_profiles WHERE status = 'approved')")
	if search := c.Query("search"); search != "" {
		like := "%" + search + "%"
		fetch = fetch.Where(
			"(users.id IN (SELECT user_id FROM provider_profiles WHERE company_name LIKE ?))",
			like,
		)
	}
	fetch.Order("users.created_at DESC").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&users)

	type publicProvider struct {
		ID          uint    `json:"id"`
		CompanyName string  `json:"company_name"`
		Description *string `json:"description"`
		Website     *string `json:"website"`
		FirstName   string  `json:"first_name"`
		LastName    string  `json:"last_name"`
	}
	out := make([]publicProvider, 0, len(users))
	for i := range users {
		u := &users[i]
		if u.ProviderProfile == nil {
			continue
		}
		out = append(out, publicProvider{
			ID:          u.ID,
			CompanyName: u.ProviderProfile.CompanyName,
			Description: u.ProviderProfile.Description,
			Website:     u.ProviderProfile.Website,
			FirstName:   u.FirstName,
			LastName:    u.LastName,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": out,
		"meta": gin.H{
			"current_page": page,
			"last_page":    lastPage,
			"per_page":     perPage,
			"total":        total,
		},
	})
}

func (h *PublicHandler) ShowProvider(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}
	var user models.User
	if err := h.DB.Preload("ProviderProfile").
		Where("id = ? AND deleted_at IS NULL AND status = ?", id, "active").
		First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}
	if user.ProviderProfile == nil || user.ProviderProfile.Status != "approved" {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var prestations []models.Prestation
	h.DB.Preload("Category").
		Where("provider_id = ? AND status = ? AND deleted_at IS NULL", user.ID, "published").
		Order("created_at DESC").
		Find(&prestations)

	c.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"company_name": user.ProviderProfile.CompanyName,
		"description":  user.ProviderProfile.Description,
		"website":      user.ProviderProfile.Website,
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"prestations":  models.ToPrestationResponses(prestations),
	})
}
