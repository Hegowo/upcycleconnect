package handlers

import (
	"net/http"
	"strconv"
	"time"

	"upcycleconnect/backend/middleware"
	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UpcyclingProjectHandler struct {
	DB    *gorm.DB
	Audit *services.AuditService
}

type projectPayload struct {
	Title       string  `json:"title" binding:"required,min=3,max=200"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	CoverImage  *string `json:"cover_image"`
	ImpactLabel *string `json:"impact_label"`
	Status      string  `json:"status"`
}

type stepPayload struct {
	Title       string  `json:"title" binding:"required,min=2,max=200"`
	Description string  `json:"description"`
	ImageURL    *string `json:"image_url"`
	StepOrder   int     `json:"step_order"`
	Completed   bool    `json:"completed"`
}

func (h *UpcyclingProjectHandler) PublicIndex(c *gin.Context) {
	q := h.DB.Preload("Provider").Preload("Steps", func(db *gorm.DB) *gorm.DB {
		return db.Order("step_order ASC")
	}).Where("status IN ?", []string{"showcased", "completed"}).Order("created_at DESC")

	if providerID := c.Query("provider_id"); providerID != "" {
		q = q.Where("provider_id = ?", providerID)
	}
	if category := c.Query("category"); category != "" {
		q = q.Where("category = ?", category)
	}

	var projects []models.UpcyclingProject
	q.Find(&projects)
	resp := make([]models.ProjectResponse, 0, len(projects))
	for i := range projects {
		resp = append(resp, models.ToProjectResponse(&projects[i], true))
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *UpcyclingProjectHandler) PublicShow(c *gin.Context) {
	var p models.UpcyclingProject
	if err := h.DB.Preload("Provider").
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Order("step_order ASC")
		}).
		Where("status IN ?", []string{"showcased", "completed"}).
		First(&p, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Projet introuvable"})
		return
	}
	c.JSON(http.StatusOK, models.ToProjectResponse(&p, true))
}

func (h *UpcyclingProjectHandler) MyProjects(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	var projects []models.UpcyclingProject
	h.DB.Preload("Steps", func(db *gorm.DB) *gorm.DB {
		return db.Order("step_order ASC")
	}).Where("provider_id = ?", user.ID).Order("created_at DESC").Find(&projects)
	resp := make([]models.ProjectResponse, 0, len(projects))
	for i := range projects {
		resp = append(resp, models.ToProjectResponse(&projects[i], true))
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *UpcyclingProjectHandler) CreateProject(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	var req projectPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	status := req.Status
	if status == "" {
		status = "in_progress"
	}
	p := models.UpcyclingProject{
		ProviderID:  user.ID,
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		CoverImage:  req.CoverImage,
		ImpactLabel: req.ImpactLabel,
		Status:      status,
	}
	if err := h.DB.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création"})
		return
	}
	h.Audit.Log(c, "project.created", "UpcyclingProject", &p.ID, nil, map[string]interface{}{"title": p.Title})
	h.DB.Preload("Provider").Preload("Steps").First(&p, p.ID)
	c.JSON(http.StatusCreated, models.ToProjectResponse(&p, true))
}

func (h *UpcyclingProjectHandler) UpdateProject(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var p models.UpcyclingProject
	if err := h.DB.Where("id = ? AND provider_id = ?", id, user.ID).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Projet introuvable"})
		return
	}
	var req projectPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	p.Title = req.Title
	p.Description = req.Description
	p.Category = req.Category
	p.CoverImage = req.CoverImage
	p.ImpactLabel = req.ImpactLabel
	if req.Status != "" {
		p.Status = req.Status
	}
	h.DB.Save(&p)
	h.DB.Preload("Provider").Preload("Steps", func(db *gorm.DB) *gorm.DB {
		return db.Order("step_order ASC")
	}).First(&p, p.ID)
	c.JSON(http.StatusOK, models.ToProjectResponse(&p, true))
}

func (h *UpcyclingProjectHandler) DeleteProject(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var p models.UpcyclingProject
	if err := h.DB.Where("id = ? AND provider_id = ?", id, user.ID).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Projet introuvable"})
		return
	}
	h.DB.Delete(&p)
	c.JSON(http.StatusOK, gin.H{"message": "Projet supprimé"})
}

func (h *UpcyclingProjectHandler) AddStep(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var p models.UpcyclingProject
	if err := h.DB.Where("id = ? AND provider_id = ?", id, user.ID).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Projet introuvable"})
		return
	}
	var req stepPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	var completedAt *time.Time
	if req.Completed {
		t := time.Now()
		completedAt = &t
	}
	step := models.ProjectStep{
		ProjectID:   p.ID,
		Title:       req.Title,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		StepOrder:   req.StepOrder,
		CompletedAt: completedAt,
	}
	h.DB.Create(&step)
	c.JSON(http.StatusCreated, step)
}

func (h *UpcyclingProjectHandler) UpdateStep(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	projectID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	stepID, _ := strconv.ParseUint(c.Param("step_id"), 10, 64)

	var p models.UpcyclingProject
	if err := h.DB.Where("id = ? AND provider_id = ?", projectID, user.ID).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Projet introuvable"})
		return
	}
	var step models.ProjectStep
	if err := h.DB.Where("id = ? AND project_id = ?", stepID, projectID).First(&step).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Étape introuvable"})
		return
	}
	var req stepPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	step.Title = req.Title
	step.Description = req.Description
	step.ImageURL = req.ImageURL
	step.StepOrder = req.StepOrder
	if req.Completed && step.CompletedAt == nil {
		t := time.Now()
		step.CompletedAt = &t
	} else if !req.Completed {
		step.CompletedAt = nil
	}
	h.DB.Save(&step)
	c.JSON(http.StatusOK, step)
}

func (h *UpcyclingProjectHandler) DeleteStep(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	projectID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	stepID, _ := strconv.ParseUint(c.Param("step_id"), 10, 64)
	var p models.UpcyclingProject
	if err := h.DB.Where("id = ? AND provider_id = ?", projectID, user.ID).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Projet introuvable"})
		return
	}
	h.DB.Where("id = ? AND project_id = ?", stepID, projectID).Delete(&models.ProjectStep{})
	c.JSON(http.StatusOK, gin.H{"message": "Étape supprimée"})
}

func (h *UpcyclingProjectHandler) AdminIndex(c *gin.Context) {
	var projects []models.UpcyclingProject
	h.DB.Preload("Provider").Preload("Steps").Order("created_at DESC").Find(&projects)
	resp := make([]models.ProjectResponse, 0, len(projects))
	for i := range projects {
		resp = append(resp, models.ToProjectResponse(&projects[i], false))
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}
