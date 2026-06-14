package handlers

import (
	"net/http"
	"time"

	"upcycleconnect/backend/middleware"
	"upcycleconnect/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProviderDashboardHandler struct {
	DB *gorm.DB
}

func (h *ProviderDashboardHandler) activeSubscriptionPlan(userID uint) string {
	var sub models.Subscription
	if err := h.DB.Where("user_id = ? AND status = ?", userID, "active").First(&sub).Error; err != nil {
		return ""
	}
	return sub.Plan
}

func (h *ProviderDashboardHandler) Dashboard(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}

	plan := h.activeSubscriptionPlan(user.ID)
	hasAdvanced := plan == "basic" || plan == "premium"
	hasPremium := plan == "premium"

	var prestationsCount, activePrestations, projectsCount, campaignsCount int64
	h.DB.Model(&models.Prestation{}).Where("provider_id = ? AND deleted_at IS NULL", user.ID).Count(&prestationsCount)
	h.DB.Model(&models.Prestation{}).Where("provider_id = ? AND status = ? AND deleted_at IS NULL", user.ID, "published").Count(&activePrestations)
	h.DB.Model(&models.UpcyclingProject{}).Where("provider_id = ?", user.ID).Count(&projectsCount)
	h.DB.Model(&models.Campaign{}).Where("provider_id = ? AND deleted_at IS NULL", user.ID).Count(&campaignsCount)

	resp := gin.H{
		"subscription": gin.H{
			"active": plan != "",
			"plan":   plan,
		},
		"basic": gin.H{
			"prestations_count":  prestationsCount,
			"active_prestations": activePrestations,
			"projects_count":     projectsCount,
			"campaigns_count":    campaignsCount,
		},
		"advanced": nil,
		"premium":  nil,
	}

	if hasAdvanced {
		type revRow struct {
			Total int64
			Count int64
		}
		var rev revRow
		h.DB.Model(&models.Reservation{}).
			Joins("JOIN prestations ON prestations.id = reservations.prestation_id").
			Where("prestations.provider_id = ? AND reservations.status = ?", user.ID, "paid").
			Select("COALESCE(SUM(reservations.amount_cents),0) AS total, COUNT(*) AS count").
			Scan(&rev)

		var reservationsTotal int64
		h.DB.Model(&models.Reservation{}).
			Joins("JOIN prestations ON prestations.id = reservations.prestation_id").
			Where("prestations.provider_id = ?", user.ID).
			Count(&reservationsTotal)

		type monthRow struct {
			Month string
			Cents int64
		}
		var months []monthRow
		h.DB.Model(&models.Reservation{}).
			Joins("JOIN prestations ON prestations.id = reservations.prestation_id").
			Where("prestations.provider_id = ? AND reservations.status = ? AND reservations.created_at >= ?",
				user.ID, "paid", time.Now().AddDate(0, -6, 0)).
			Select("DATE_FORMAT(reservations.created_at, '%Y-%m') AS month, COALESCE(SUM(reservations.amount_cents),0) AS cents").
			Group("month").Order("month ASC").Scan(&months)

		monthly := make([]gin.H, 0, len(months))
		for _, m := range months {
			monthly = append(monthly, gin.H{"month": m.Month, "cents": m.Cents})
		}

		type matRow struct {
			Category string
			Count    int64
			Weight   float64
		}
		var mats []matRow
		h.DB.Model(&models.DepositRequest{}).
			Joins("LEFT JOIN prestation_categories ON prestation_categories.id = deposit_requests.category_id").
			Where("deposit_requests.status = ? AND deposit_requests.collected_at IS NULL", "approved").
			Select("COALESCE(prestation_categories.name, 'Non catégorisé') AS category, COUNT(*) AS count, COALESCE(SUM(deposit_requests.estimated_weight),0) AS weight").
			Group("category").Order("count DESC").Scan(&mats)

		materials := make([]gin.H, 0, len(mats))
		for _, m := range mats {
			materials = append(materials, gin.H{"category": m.Category, "count": m.Count, "weight_kg": m.Weight})
		}

		resp["advanced"] = gin.H{
			"revenue_total_cents": rev.Total,
			"reservations_paid":   rev.Count,
			"reservations_total":  reservationsTotal,
			"monthly_revenue":     monthly,
			"material_stats":      materials,
		}
	}

	if hasPremium {
		type ecoRow struct {
			CO2    float64
			Weight float64
			Items  int64
		}
		var eco ecoRow
		h.DB.Model(&models.DepositRequest{}).
			Where("collected_by = ?", user.ID).
			Select("COALESCE(SUM(carbon_savings),0) AS co2, COALESCE(SUM(estimated_weight),0) AS weight, COUNT(*) AS items").
			Scan(&eco)

		var alerts []models.DepositRequest
		h.DB.Preload("Category").Preload("CollectionPoint").
			Where("status = ? AND collected_at IS NULL", "approved").
			Order("validated_at DESC").
			Limit(10).
			Find(&alerts)

		alertList := make([]gin.H, 0, len(alerts))
		for i := range alerts {
			d := &alerts[i]
			cat := ""
			if d.Category != nil {
				cat = d.Category.Name
			}
			point := ""
			if d.CollectionPoint != nil {
				point = d.CollectionPoint.Name
			}
			alertList = append(alertList, gin.H{
				"id":               d.ID,
				"title":            d.Title,
				"category":         cat,
				"collection_point": point,
				"weight_kg":        d.EstimatedWeight,
				"qr_code":          d.QRCode,
			})
		}

		resp["premium"] = gin.H{
			"eco_impact": gin.H{
				"co2_saved_kg":    eco.CO2,
				"weight_saved_kg": eco.Weight,
				"items_collected": eco.Items,
			},
			"collection_alerts": alertList,
		}
	}

	c.JSON(http.StatusOK, resp)
}
