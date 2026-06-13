package handlers

import (
	"net/http"
	"time"

	"upcycleconnect/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FinanceHandler struct {
	DB *gorm.DB
}

// Overview returns the platform's revenue breakdown for the admin financial dashboard.
// Sources (per Annexe 1): commissions on transactions, pro subscriptions, ad campaigns,
// and formations/ateliers (sold as priced prestations, hence captured in commissions+GMV).
func (h *FinanceHandler) Overview(c *gin.Context) {
	// ── Transactions (paid prestation reservations) ──────────────────────
	type txAgg struct {
		Gross      int64
		Commission int64
		Net        int64
		Count      int64
	}
	var tx txAgg
	h.DB.Model(&models.Reservation{}).
		Where("status = ?", "paid").
		Select("COALESCE(SUM(amount_cents),0) AS gross, COALESCE(SUM(commission_cents),0) AS commission, COALESCE(SUM(net_cents),0) AS net, COUNT(*) AS count").
		Scan(&tx)

	// ── Subscriptions (active) ───────────────────────────────────────────
	var subs []models.Subscription
	h.DB.Where("status = ?", "active").Find(&subs)
	var subMRR int64
	subCount := int64(len(subs))
	for _, s := range subs {
		if plan, ok := models.SubscriptionPlans[s.Plan]; ok {
			subMRR += plan.AmountCents
		}
	}

	// ── Campaigns (paid) ─────────────────────────────────────────────────
	type campAgg struct {
		Total int64
		Count int64
	}
	var camp campAgg
	h.DB.Model(&models.Campaign{}).
		Where("paid_at IS NOT NULL AND deleted_at IS NULL").
		Select("COALESCE(SUM(budget_cents),0) AS total, COUNT(*) AS count").
		Scan(&camp)

	// Platform revenue = commissions + subscription MRR + campaign payments.
	platformRevenue := tx.Commission + subMRR + camp.Total

	// ── Monthly commissions (last 6 months) ──────────────────────────────
	type monthRow struct {
		Month string
		Cents int64
	}
	var months []monthRow
	h.DB.Model(&models.Reservation{}).
		Where("status = ? AND updated_at >= ?", "paid", time.Now().AddDate(0, -6, 0)).
		Select("DATE_FORMAT(updated_at, '%Y-%m') AS month, COALESCE(SUM(commission_cents),0) AS cents").
		Group("month").Order("month ASC").Scan(&months)
	monthly := make([]gin.H, 0, len(months))
	for _, m := range months {
		monthly = append(monthly, gin.H{"month": m.Month, "cents": m.Cents})
	}

	c.JSON(http.StatusOK, gin.H{
		"commission_rate_percent": int(models.CommissionRate * 100),
		"platform_revenue_cents":  platformRevenue,
		"transactions": gin.H{
			"gross_cents":      tx.Gross,      // total GMV processed
			"commission_cents": tx.Commission, // platform's cut
			"net_to_providers_cents": tx.Net,  // owed to providers
			"count":            tx.Count,
		},
		"subscriptions": gin.H{
			"active_count": subCount,
			"mrr_cents":    subMRR, // monthly recurring revenue
		},
		"campaigns": gin.H{
			"paid_count":  camp.Count,
			"total_cents": camp.Total,
		},
		"monthly_commissions": monthly,
	})
}

// Transactions lists recent paid transactions with their commission split (admin view).
func (h *FinanceHandler) Transactions(c *gin.Context) {
	var reservations []models.Reservation
	h.DB.Preload("User").Preload("Prestation").
		Where("status = ?", "paid").
		Order("updated_at DESC").
		Limit(100).
		Find(&reservations)

	rows := make([]gin.H, 0, len(reservations))
	for i := range reservations {
		r := &reservations[i]
		title := ""
		if r.Prestation != nil {
			title = r.Prestation.Title
		}
		customer := ""
		if r.User != nil {
			customer = r.User.FirstName + " " + r.User.LastName
		}
		rows = append(rows, gin.H{
			"id":               r.ID,
			"title":            title,
			"customer":         customer,
			"amount_cents":     r.AmountCents,
			"commission_cents": r.CommissionCents,
			"net_cents":        r.NetCents,
			"date":             r.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": rows})
}
