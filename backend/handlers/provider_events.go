package handlers

import (
	"math"
	"net/http"
	"strconv"

	"upcycleconnect/backend/middleware"
	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProviderEventHandler struct {
	DB    *gorm.DB
	Audit *services.AuditService
}

func (h *ProviderEventHandler) ensureApproved(c *gin.Context, user *models.User) bool {
	var count int64
	h.DB.Model(&models.ProviderProfile{}).
		Where("user_id = ? AND status = ?", user.ID, "approved").
		Count(&count)
	if count == 0 {
		c.JSON(http.StatusForbidden, gin.H{"message": "Profil prestataire approuvé requis pour gérer des événements."})
		return false
	}
	return true
}

func (h *ProviderEventHandler) List(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if !h.ensureApproved(c, user) {
		return
	}

	query := h.DB.Model(&models.Event{}).Where("created_by = ? AND deleted_at IS NULL", user.ID)
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
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

	var items []models.Event
	fetch := h.DB.Preload("Category").Where("created_by = ? AND deleted_at IS NULL", user.ID)
	if status := c.Query("status"); status != "" {
		fetch = fetch.Where("status = ?", status)
	}
	fetch.Order("created_at DESC").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&items)

	type regCount struct {
		EventID uint
		Count   int64
	}
	counts := make(map[uint]int64)
	if len(items) > 0 {
		ids := make([]uint, len(items))
		for i, e := range items {
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

	responses := make([]models.EventResponse, len(items))
	for i := range items {
		responses[i] = models.ToEventResponseWithCount(&items[i], counts[items[i].ID])
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

func (h *ProviderEventHandler) Store(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if !h.ensureApproved(c, user) {
		return
	}

	var req struct {
		CategoryID      *uint   `json:"category_id"`
		Title           string  `json:"title" binding:"required"`
		Description     *string `json:"description"`
		Location        *string `json:"location"`
		StartDate       string  `json:"start_date" binding:"required"`
		EndDate         string  `json:"end_date" binding:"required"`
		MaxParticipants *int    `json:"max_participants"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides. Titre, dates début et fin requis."})
		return
	}

	start, err := parseDateTime(req.StartDate)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Format de date de début invalide."})
		return
	}
	end, err := parseDateTime(req.EndDate)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Format de date de fin invalide."})
		return
	}
	if !end.After(start) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "La date de fin doit être après la date de début."})
		return
	}

	createdBy := user.ID
	event := models.Event{
		CategoryID:      req.CategoryID,
		Title:           req.Title,
		Description:     req.Description,
		Location:        req.Location,
		StartDate:       start,
		EndDate:         end,
		MaxParticipants: req.MaxParticipants,
		Status:          "draft",
		CreatedBy:       &createdBy,
	}
	if err := h.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création."})
		return
	}
	h.DB.Preload("Category").First(&event, event.ID)

	h.Audit.Log(c, "provider.event_created", "Event", &event.ID, nil, map[string]interface{}{
		"title": event.Title,
	})

	c.JSON(http.StatusCreated, models.ToEventResponseWithCount(&event, 0))
}

func (h *ProviderEventHandler) Update(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if !h.ensureApproved(c, user) {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var event models.Event
	if err := h.DB.Where("id = ? AND created_by = ? AND deleted_at IS NULL", id, user.ID).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}
	if event.Status == "published" || event.Status == "cancelled" {
		c.JSON(http.StatusConflict, gin.H{"message": "Un événement publié ou annulé ne peut pas être modifié."})
		return
	}

	var req struct {
		CategoryID      *uint   `json:"category_id"`
		Title           *string `json:"title"`
		Description     *string `json:"description"`
		Location        *string `json:"location"`
		StartDate       *string `json:"start_date"`
		EndDate         *string `json:"end_date"`
		MaxParticipants *int    `json:"max_participants"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return
	}

	updates := map[string]interface{}{}
	if req.CategoryID != nil {
		updates["category_id"] = *req.CategoryID
	}
	if req.Title != nil && *req.Title != "" {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Location != nil {
		updates["location"] = *req.Location
	}
	if req.MaxParticipants != nil {
		updates["max_participants"] = *req.MaxParticipants
	}
	if req.StartDate != nil {
		start, err := parseDateTime(*req.StartDate)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Format de date de début invalide."})
			return
		}
		updates["start_date"] = start
	}
	if req.EndDate != nil {
		end, err := parseDateTime(*req.EndDate)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Format de date de fin invalide."})
			return
		}
		updates["end_date"] = end
	}
	if event.Status == "pending" {
		updates["status"] = "draft"
	}

	if len(updates) > 0 {
		h.DB.Model(&event).Updates(updates)
	}
	h.DB.Preload("Category").First(&event, event.ID)

	var count int64
	h.DB.Model(&models.EventRegistration{}).Where("event_id = ?", event.ID).Count(&count)

	h.Audit.Log(c, "provider.event_updated", "Event", &event.ID, nil, updates)
	c.JSON(http.StatusOK, models.ToEventResponseWithCount(&event, count))
}

func (h *ProviderEventHandler) Submit(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if !h.ensureApproved(c, user) {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var event models.Event
	if err := h.DB.Where("id = ? AND created_by = ? AND deleted_at IS NULL", id, user.ID).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}
	if event.Status != "draft" {
		c.JSON(http.StatusConflict, gin.H{"message": "Seuls les événements en brouillon peuvent être soumis."})
		return
	}

	h.DB.Model(&event).Update("status", "pending")
	h.Audit.Log(c, "provider.event_submitted", "Event", &event.ID,
		map[string]string{"status": "draft"}, map[string]string{"status": "pending"})

	c.JSON(http.StatusOK, gin.H{"message": "Événement soumis pour validation par un administrateur.", "status": "pending"})
}

func (h *ProviderEventHandler) Destroy(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if !h.ensureApproved(c, user) {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var event models.Event
	if err := h.DB.Where("id = ? AND created_by = ? AND deleted_at IS NULL", id, user.ID).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}
	if event.Status == "published" {
		c.JSON(http.StatusConflict, gin.H{"message": "Un événement publié ne peut pas être supprimé."})
		return
	}

	h.Audit.Log(c, "provider.event_deleted", "Event", &event.ID, map[string]interface{}{"title": event.Title}, nil)
	h.DB.Delete(&event)
	c.Status(http.StatusNoContent)
}
