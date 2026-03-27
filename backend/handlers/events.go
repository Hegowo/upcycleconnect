package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"upcycleconnect/backend/middleware"
	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EventHandler struct {
	DB    *gorm.DB
	Audit *services.AuditService
}

func (h *EventHandler) Index(c *gin.Context) {
	query := h.DB.Model(&models.Event{}).Where("platform_events.deleted_at IS NULL")

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if from := c.Query("start_date_from"); from != "" {
		query = query.Where("start_date >= ?", from)
	}
	if to := c.Query("start_date_to"); to != "" {
		query = query.Where("start_date <= ?", to)
	}

	var total int64
	query.Count(&total)

	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if perPage < 1 {
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
	fetchQuery := h.DB.Preload("Category").Preload("Creator.Roles").
		Where("platform_events.deleted_at IS NULL")

	if status := c.Query("status"); status != "" {
		fetchQuery = fetchQuery.Where("status = ?", status)
	}
	if from := c.Query("start_date_from"); from != "" {
		fetchQuery = fetchQuery.Where("start_date >= ?", from)
	}
	if to := c.Query("start_date_to"); to != "" {
		fetchQuery = fetchQuery.Where("start_date <= ?", to)
	}

	fetchQuery.Order("start_date ASC").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&events)

	// Charger les compteurs de participants pour tous les événements récupérés
	type regCount struct {
		EventID uint
		Count   int64
	}
	var regCounts []regCount
	if len(events) > 0 {
		ids := make([]uint, len(events))
		for i, e := range events {
			ids[i] = e.ID
		}
		h.DB.Model(&models.EventRegistration{}).
			Select("event_id, COUNT(*) as count").
			Where("event_id IN ?", ids).
			Group("event_id").
			Scan(&regCounts)
	}
	countByID := make(map[uint]int64)
	for _, rc := range regCounts {
		countByID[rc.EventID] = rc.Count
	}

	responses := make([]models.EventResponse, len(events))
	for i := range events {
		responses[i] = models.ToEventResponseWithCount(&events[i], countByID[events[i].ID])
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

func (h *EventHandler) Store(c *gin.Context) {
	var req struct {
		CategoryID      *uint   `json:"category_id"`
		Title           string  `json:"title" binding:"required"`
		Description     *string `json:"description"`
		Location        *string `json:"location"`
		StartDate       string  `json:"start_date" binding:"required"`
		EndDate         string  `json:"end_date" binding:"required"`
		MaxParticipants *int    `json:"max_participants"`
		Status          string  `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return
	}

	authUser := middleware.GetAuthUser(c)
	var createdBy *uint
	if authUser != nil {
		createdBy = &authUser.ID
	}

	event := models.Event{
		CategoryID:      req.CategoryID,
		Title:           req.Title,
		Description:     req.Description,
		Location:        req.Location,
		MaxParticipants: req.MaxParticipants,
		Status:          req.Status,
		CreatedBy:       createdBy,
	}
	if event.Status == "" {
		event.Status = "draft"
	}

	startDate, err := parseDateTime(req.StartDate)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Format de date de début invalide."})
		return
	}
	endDate, err := parseDateTime(req.EndDate)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Format de date de fin invalide."})
		return
	}
	event.StartDate = startDate
	event.EndDate = endDate

	if err := h.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création."})
		return
	}

	h.DB.Preload("Category").Preload("Creator.Roles").First(&event, event.ID)

	h.Audit.Log(c, "event.created", "Event", &event.ID, nil, map[string]interface{}{
		"title": event.Title,
	})

	c.JSON(http.StatusCreated, models.ToEventResponse(&event))
}

func (h *EventHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var event models.Event
	if err := h.DB.Preload("Category").Preload("Creator.Roles").
		Where("id = ? AND deleted_at IS NULL", id).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var count int64
	h.DB.Model(&models.EventRegistration{}).Where("event_id = ?", event.ID).Count(&count)

	c.JSON(http.StatusOK, models.ToEventResponseWithCount(&event, count))
}

func (h *EventHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var event models.Event
	if err := h.DB.Where("id = ? AND deleted_at IS NULL", id).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
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
		Status          *string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return
	}

	// Capture full before state for the diff
	oldLocation := ""
	if event.Location != nil {
		oldLocation = *event.Location
	}
	old := map[string]interface{}{
		"title":            event.Title,
		"status":           event.Status,
		"location":         oldLocation,
		"start_date":       event.StartDate.Format("2006-01-02T15:04"),
		"end_date":         event.EndDate.Format("2006-01-02T15:04"),
		"max_participants": event.MaxParticipants,
	}

	updates := map[string]interface{}{}
	after := map[string]interface{}{"title": event.Title}

	if req.CategoryID != nil {
		updates["category_id"] = *req.CategoryID
	}
	if req.Title != nil {
		updates["title"] = *req.Title
		after["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
		after["description"] = *req.Description
	}
	if req.Location != nil {
		updates["location"] = *req.Location
		after["location"] = *req.Location
	}
	if req.StartDate != nil {
		d, err := parseDateTime(*req.StartDate)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Format de date de début invalide."})
			return
		}
		updates["start_date"] = d
		after["start_date"] = d.Format("2006-01-02T15:04")
	}
	if req.EndDate != nil {
		d, err := parseDateTime(*req.EndDate)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Format de date de fin invalide."})
			return
		}
		updates["end_date"] = d
		after["end_date"] = d.Format("2006-01-02T15:04")
	}
	if req.MaxParticipants != nil {
		updates["max_participants"] = *req.MaxParticipants
		after["max_participants"] = *req.MaxParticipants
	}
	if req.Status != nil {
		updates["status"] = *req.Status
		after["status"] = *req.Status
	}

	if len(updates) > 0 {
		h.DB.Model(&event).Updates(updates)
	}

	h.DB.Preload("Category").Preload("Creator.Roles").First(&event, event.ID)

	h.Audit.Log(c, "event.updated", "Event", &event.ID, old, after)

	c.JSON(http.StatusOK, models.ToEventResponse(&event))
}

func (h *EventHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var event models.Event
	if err := h.DB.Where("id = ? AND deleted_at IS NULL", id).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=draft published cancelled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return
	}

	old := map[string]string{"status": event.Status}
	h.DB.Model(&event).Update("status", req.Status)

	h.Audit.Log(c, "event.status_changed", "Event", &event.ID, old, map[string]string{"status": req.Status})

	h.DB.Preload("Category").Preload("Creator.Roles").First(&event, event.ID)

	c.JSON(http.StatusOK, models.ToEventResponse(&event))
}

func (h *EventHandler) Destroy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var event models.Event
	if err := h.DB.Where("id = ? AND deleted_at IS NULL", id).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	h.Audit.Log(c, "event.deleted", "Event", &event.ID, map[string]interface{}{
		"title": event.Title,
	}, nil)
	h.DB.Delete(&event)

	c.Status(http.StatusNoContent)
}

func parseDateTime(s string) (t time.Time, err error) {
	formats := []string{
		"2006-01-02T15:04:05.000000Z",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	}
	for _, f := range formats {
		t, err = time.Parse(f, s)
		if err == nil {
			return t, nil
		}
	}
	return t, fmt.Errorf("unable to parse date: %s", s)
}
