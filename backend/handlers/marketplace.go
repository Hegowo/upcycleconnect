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

type MarketplaceHandler struct {
	DB            *gorm.DB
	Stripe        *services.StripeService
	Audit         *services.AuditService
	Notifications *services.NotificationService
}

func (h *MarketplaceHandler) approvedProvider(c *gin.Context) (*models.User, bool) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return nil, false
	}
	var p models.ProviderProfile
	if err := h.DB.Where("user_id = ? AND status = ?", user.ID, "approved").First(&p).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Profil prestataire approuvé requis."})
		return nil, false
	}
	return user, true
}

func (h *MarketplaceHandler) List(c *gin.Context) {
	user, ok := h.approvedProvider(c)
	if !ok {
		return
	}
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	base := h.DB.Model(&models.DepositRequest{}).
		Where("status = ? AND purchased_by IS NULL", "approved")
	if s := c.Query("search"); s != "" {
		like := "%" + s + "%"
		base = base.Where("title LIKE ? OR description LIKE ?", like, like)
	}
	if v := c.Query("category_id"); v != "" {
		base = base.Where("category_id = ?", v)
	}
	if v := c.Query("collection_point_id"); v != "" {
		base = base.Where("collection_point_id = ?", v)
	}
	if v := c.Query("sale_type"); v == "don" || v == "vente" {
		base = base.Where("sale_type = ?", v)
	}
	if v := c.Query("condition"); v == "good" || v == "fair" || v == "poor" {
		base = base.Where("`condition` = ?", v)
	}

	var total int64
	base.Count(&total)
	lastPage := int(math.Ceil(float64(total) / float64(perPage)))
	if lastPage < 1 {
		lastPage = 1
	}

	var items []models.DepositRequest
	base.Preload("User").Preload("Category").Preload("CollectionPoint").
		Order("validated_at DESC").
		Offset((page - 1) * perPage).Limit(perPage).
		Find(&items)

	favSet := map[uint]bool{}
	var favIDs []uint
	h.DB.Model(&models.DepositFavorite{}).Where("provider_id = ?", user.ID).Pluck("deposit_id", &favIDs)
	for _, id := range favIDs {
		favSet[id] = true
	}
	resp := models.ToDepositResponses(items)
	for i := range resp {
		resp[i].IsFavorited = favSet[resp[i].ID]
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp,
		"meta": gin.H{"current_page": page, "last_page": lastPage, "per_page": perPage, "total": total},
	})
}

func (h *MarketplaceHandler) Purchase(c *gin.Context) {
	user, ok := h.approvedProvider(c)
	if !ok {
		return
	}
	var dep models.DepositRequest
	if err := h.DB.First(&dep, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Annonce introuvable"})
		return
	}
	if dep.Status != "approved" || dep.PurchasedBy != nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Cet objet n'est plus disponible."})
		return
	}
	if dep.UserID == user.ID {
		c.JSON(http.StatusForbidden, gin.H{"message": "Vous ne pouvez pas acheter votre propre annonce."})
		return
	}

	if dep.SaleType == "vente" && dep.PriceCents > 0 {
		sess, err := h.Stripe.CreateDepositCheckout(user.ID, dep.ID, user.Email, dep.Title, dep.PriceCents)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		sid := sess.ID
		h.DB.Create(&models.DepositPurchase{
			DepositID: dep.ID, ProviderID: user.ID, AmountCents: dep.PriceCents,
			Status: "pending", StripeSessionID: &sid,
		})
		c.JSON(http.StatusOK, gin.H{"type": "vente", "checkout_url": sess.URL})
		return
	}

	now := time.Now()
	purchase := models.DepositPurchase{DepositID: dep.ID, ProviderID: user.ID, AmountCents: 0, Status: "paid", PaidAt: &now}
	h.DB.Create(&purchase)
	h.DB.Model(&dep).Updates(map[string]interface{}{"purchased_by": user.ID, "purchased_at": now})
	if h.Notifications != nil {
		h.Notifications.MustNotify(dep.UserID, "deposit.purchased",
			"Objet réservé",
			fmt.Sprintf("Un professionnel a réservé votre objet « %s ».", dep.Title),
			"/profil")
	}
	h.Audit.Log(c, "deposit.purchased", "DepositRequest", &dep.ID, nil, map[string]interface{}{"provider_id": user.ID})
	h.DB.Preload("Deposit").First(&purchase, purchase.ID)
	c.JSON(http.StatusCreated, gin.H{"type": "don", "data": purchase, "qr_code": dep.QRCode})
}

func (h *MarketplaceHandler) Purchases(c *gin.Context) {
	user, ok := h.approvedProvider(c)
	if !ok {
		return
	}
	var items []models.DepositPurchase
	h.DB.Preload("Deposit").Preload("Deposit.Category").Preload("Deposit.CollectionPoint").
		Where("provider_id = ?", user.ID).Order("created_at DESC").Find(&items)
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *MarketplaceHandler) PurchaseShow(c *gin.Context) {
	user, ok := h.approvedProvider(c)
	if !ok {
		return
	}
	var p models.DepositPurchase
	if err := h.DB.Preload("Deposit").Preload("Deposit.Category").Preload("Deposit.CollectionPoint").
		Where("id = ? AND provider_id = ?", c.Param("id"), user.ID).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Achat introuvable"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

func (h *MarketplaceHandler) FavoriteAdd(c *gin.Context) {
	user, ok := h.approvedProvider(c)
	if !ok {
		return
	}
	depID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var dep models.DepositRequest
	if err := h.DB.First(&dep, depID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Annonce introuvable"})
		return
	}
	fav := models.DepositFavorite{DepositID: uint(depID), ProviderID: user.ID}
	h.DB.Where("deposit_id = ? AND provider_id = ?", depID, user.ID).FirstOrCreate(&fav)
	c.JSON(http.StatusOK, gin.H{"message": "Ajouté aux favoris.", "is_favorited": true})
}

func (h *MarketplaceHandler) FavoriteRemove(c *gin.Context) {
	user, ok := h.approvedProvider(c)
	if !ok {
		return
	}
	h.DB.Where("deposit_id = ? AND provider_id = ?", c.Param("id"), user.ID).Delete(&models.DepositFavorite{})
	c.JSON(http.StatusOK, gin.H{"message": "Retiré des favoris.", "is_favorited": false})
}

func (h *MarketplaceHandler) Favorites(c *gin.Context) {
	user, ok := h.approvedProvider(c)
	if !ok {
		return
	}
	var favIDs []uint
	h.DB.Model(&models.DepositFavorite{}).Where("provider_id = ?", user.ID).Pluck("deposit_id", &favIDs)
	var items []models.DepositRequest
	if len(favIDs) > 0 {
		h.DB.Preload("User").Preload("Category").Preload("CollectionPoint").
			Where("id IN ?", favIDs).Find(&items)
	}
	resp := models.ToDepositResponses(items)
	for i := range resp {
		resp[i].IsFavorited = true
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}
