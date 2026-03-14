package handlers

import (
	"net/http"

	"upcycleconnect/backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardHandler struct {
	DB *gorm.DB
}

func (h *DashboardHandler) Stats(c *gin.Context) {
	var usersCount int64
	h.DB.Model(&models.User{}).
		Where("deleted_at IS NULL").
		Where("id NOT IN (SELECT user_id FROM user_roles)").
		Count(&usersCount)

	var providersCount int64
	h.DB.Model(&models.ProviderProfile{}).Where("status = ?", "approved").Count(&providersCount)

	var providersPending int64
	h.DB.Model(&models.ProviderProfile{}).Where("status = ?", "pending").Count(&providersPending)

	var prestationsCount int64
	h.DB.Model(&models.Prestation{}).Where("deleted_at IS NULL").Count(&prestationsCount)

	var eventsCount int64
	h.DB.Model(&models.Event{}).
		Where("deleted_at IS NULL").
		Where("status IN ?", []string{"draft", "published"}).
		Count(&eventsCount)

	type statusCount struct {
		Status string
		Count  int64
	}
	var statusCounts []statusCount
	h.DB.Model(&models.Prestation{}).
		Select("status, COUNT(*) as count").
		Where("deleted_at IS NULL").
		Group("status").
		Find(&statusCounts)

	prestationsByStatus := make(map[string]int64)
	for _, sc := range statusCounts {
		prestationsByStatus[sc.Status] = sc.Count
	}

	c.JSON(http.StatusOK, gin.H{
		"users_count":           usersCount,
		"providers_count":       providersCount,
		"providers_pending":     providersPending,
		"prestations_count":     prestationsCount,
		"events_count":          eventsCount,
		"prestations_by_status": prestationsByStatus,
	})
}
