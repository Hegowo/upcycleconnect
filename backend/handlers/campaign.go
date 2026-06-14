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

type CampaignHandler struct {
	DB     *gorm.DB
	Stripe *services.StripeService
	Audit  *services.AuditService
}

type campaignPayload struct {
	Title       string  `json:"title" binding:"required,min=3,max=200"`
	Description string  `json:"description"`
	ImageURL    *string `json:"image_url"`

	BudgetEuros   float64 `json:"budget_euros" binding:"required,min=100,max=500"`
	StartDate     *string `json:"start_date"`
	EndDate       *string `json:"end_date"`
	PrestationIDs []uint  `json:"prestation_ids"`
	EventIDs      []uint  `json:"event_ids"`
}

func (h *CampaignHandler) syncCampaignItems(campaignID, providerID uint, prestationIDs, eventIDs []uint) {
	h.DB.Where("campaign_id = ?", campaignID).Delete(&models.CampaignItem{})

	for _, pid := range prestationIDs {
		var cnt int64
		h.DB.Model(&models.Prestation{}).Where("id = ? AND provider_id = ?", pid, providerID).Count(&cnt)
		if cnt > 0 {
			h.DB.Create(&models.CampaignItem{CampaignID: campaignID, ItemType: "prestation", ItemID: pid})
		}
	}
	for _, eid := range eventIDs {
		var cnt int64
		h.DB.Model(&models.Event{}).Where("id = ? AND created_by = ?", eid, providerID).Count(&cnt)
		if cnt > 0 {
			h.DB.Create(&models.CampaignItem{CampaignID: campaignID, ItemType: "event", ItemID: eid})
		}
	}
}

func parseCampaignDate(raw *string) *time.Time {
	if raw == nil || *raw == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, *raw)
	if err != nil {
		t2, err2 := time.Parse("2006-01-02", *raw)
		if err2 != nil {
			return nil
		}
		return &t2
	}
	return &t
}

func (h *CampaignHandler) MyCampaigns(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	var campaigns []models.Campaign
	h.DB.Preload("Provider").Preload("Items").Where("provider_id = ? AND deleted_at IS NULL", user.ID).
		Order("created_at DESC").Find(&campaigns)
	resp := make([]models.CampaignResponse, 0, len(campaigns))
	for i := range campaigns {
		resp = append(resp, models.ToCampaignResponse(&campaigns[i]))
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *CampaignHandler) CreateCampaign(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}

	if !h.hasPremiumSubscription(user.ID) {
		c.JSON(http.StatusPaymentRequired, gin.H{
			"message": "La gestion de campagnes publicitaires nécessite un abonnement Premium.",
			"upgrade": "/profil/pro/abonnement",
		})
		return
	}
	var req campaignPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	budget := int64(math.Round(req.BudgetEuros * 100))
	camp := models.Campaign{
		ProviderID:  user.ID,
		Title:       req.Title,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		BudgetCents: budget,
		Status:      "draft",
		StartDate:   parseCampaignDate(req.StartDate),
		EndDate:     parseCampaignDate(req.EndDate),
	}
	if err := h.DB.Create(&camp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création"})
		return
	}
	h.syncCampaignItems(camp.ID, user.ID, req.PrestationIDs, req.EventIDs)
	h.Audit.Log(c, "campaign.created", "Campaign", &camp.ID, nil, map[string]interface{}{"title": camp.Title})
	h.DB.Preload("Provider").Preload("Items").First(&camp, camp.ID)
	c.JSON(http.StatusCreated, models.ToCampaignResponse(&camp))
}

func (h *CampaignHandler) UpdateCampaign(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var camp models.Campaign
	if err := h.DB.Where("id = ? AND provider_id = ?", id, user.ID).First(&camp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Campagne introuvable"})
		return
	}
	if camp.Status != "draft" && camp.Status != "rejected" {
		c.JSON(http.StatusConflict, gin.H{"message": "Seules les campagnes en brouillon peuvent être modifiées"})
		return
	}
	var req campaignPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	camp.Title = req.Title
	camp.Description = req.Description
	camp.ImageURL = req.ImageURL
	camp.BudgetCents = int64(math.Round(req.BudgetEuros * 100))
	camp.StartDate = parseCampaignDate(req.StartDate)
	camp.EndDate = parseCampaignDate(req.EndDate)
	h.DB.Save(&camp)
	h.syncCampaignItems(camp.ID, user.ID, req.PrestationIDs, req.EventIDs)
	h.DB.Preload("Provider").Preload("Items").First(&camp, camp.ID)
	c.JSON(http.StatusOK, models.ToCampaignResponse(&camp))
}

func (h *CampaignHandler) DeleteCampaign(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var camp models.Campaign
	if err := h.DB.Where("id = ? AND provider_id = ?", id, user.ID).First(&camp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Campagne introuvable"})
		return
	}
	if camp.Status == "active" {
		c.JSON(http.StatusConflict, gin.H{"message": "Impossible de supprimer une campagne active"})
		return
	}
	h.DB.Delete(&camp)
	c.JSON(http.StatusOK, gin.H{"message": "Campagne supprimée"})
}

func (h *CampaignHandler) SubmitCampaign(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var camp models.Campaign
	if err := h.DB.Where("id = ? AND provider_id = ?", id, user.ID).First(&camp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Campagne introuvable"})
		return
	}
	if camp.Status != "draft" && camp.Status != "rejected" {
		c.JSON(http.StatusConflict, gin.H{"message": "Cette campagne a déjà été soumise"})
		return
	}
	sess, err := h.Stripe.CreateCampaignCheckout(user.ID, camp.ID, user.Email, camp.Title, camp.BudgetCents)
	if err != nil {
		h.DB.Model(&camp).Updates(map[string]interface{}{"status": "pending"})
		c.JSON(http.StatusOK, gin.H{
			"message":      fmt.Sprintf("Campagne soumise (mode hors-ligne : %s)", err.Error()),
			"checkout_url": nil,
		})
		return
	}
	h.DB.Model(&camp).Updates(map[string]interface{}{
		"status":            "pending",
		"stripe_session_id": sess.ID,
	})
	h.Audit.Log(c, "campaign.submitted", "Campaign", &camp.ID, nil, map[string]interface{}{"stripe_session": sess.ID})
	c.JSON(http.StatusOK, gin.H{"checkout_url": sess.URL, "session_id": sess.ID})
}

func (h *CampaignHandler) ActiveCampaigns(c *gin.Context) {
	var campaigns []models.Campaign
	now := time.Now()
	h.DB.Preload("Provider").Preload("Items").
		Where("status = ? AND (start_date IS NULL OR start_date <= ?) AND (end_date IS NULL OR end_date >= ?)",
			"active", now, now).
		Order("created_at DESC").
		Limit(6).
		Find(&campaigns)
	resp := make([]models.CampaignResponse, 0, len(campaigns))
	for i := range campaigns {
		resp = append(resp, models.ToCampaignResponse(&campaigns[i]))
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *CampaignHandler) PublicShow(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var camp models.Campaign
	if err := h.DB.Preload("Provider").Preload("Items").
		Where("status = ?", "active").First(&camp, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Campagne introuvable ou inactive"})
		return
	}

	prestationIDs := []uint{}
	eventIDs := []uint{}
	for _, it := range camp.Items {
		if it.ItemType == "prestation" {
			prestationIDs = append(prestationIDs, it.ItemID)
		} else if it.ItemType == "event" {
			eventIDs = append(eventIDs, it.ItemID)
		}
	}

	var prestations []models.Prestation
	if len(prestationIDs) > 0 {
		h.DB.Preload("Category").Where("id IN ? AND status = ?", prestationIDs, "published").Find(&prestations)
	}
	var events []models.Event
	if len(eventIDs) > 0 {
		h.DB.Where("id IN ? AND status = ?", eventIDs, "published").Find(&events)
	}

	c.JSON(http.StatusOK, gin.H{
		"campaign":    models.ToCampaignResponse(&camp),
		"prestations": models.ToPrestationResponses(prestations),
		"events":      models.ToEventResponses(events),
	})
}

func (h *CampaignHandler) AdminIndex(c *gin.Context) {
	var campaigns []models.Campaign
	q := h.DB.Preload("Provider").Preload("Items").Order("created_at DESC")
	if status := c.Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}
	q.Find(&campaigns)
	resp := make([]models.CampaignResponse, 0, len(campaigns))
	for i := range campaigns {
		resp = append(resp, models.ToCampaignResponse(&campaigns[i]))
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *CampaignHandler) AdminUpdateStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var camp models.Campaign
	if err := h.DB.First(&camp, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Campagne introuvable"})
		return
	}
	var req struct {
		Status          string  `json:"status" binding:"required,oneof=active rejected ended"`
		RejectionReason *string `json:"rejection_reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	updates := map[string]interface{}{"status": req.Status}
	if req.Status == "active" {
		now := time.Now()
		updates["paid_at"] = &now
	}
	if req.RejectionReason != nil {
		updates["rejection_reason"] = req.RejectionReason
	}
	h.DB.Model(&camp).Updates(updates)
	h.Audit.Log(c, "campaign.status_updated", "Campaign", &camp.ID, map[string]string{"status": camp.Status}, map[string]string{"status": req.Status})
	h.DB.Preload("Provider").First(&camp, camp.ID)
	c.JSON(http.StatusOK, models.ToCampaignResponse(&camp))
}

func (h *CampaignHandler) hasPremiumSubscription(userID uint) bool {
	var sub models.Subscription
	if err := h.DB.Where("user_id = ? AND status = ? AND plan = ?", userID, "active", "premium").
		First(&sub).Error; err != nil {
		return false
	}
	return true
}
