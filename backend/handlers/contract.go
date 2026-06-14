package handlers

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"upcycleconnect/backend/middleware"
	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ContractHandler struct {
	DB    *gorm.DB
	PDF   *services.PDFService
	Audit *services.AuditService
}

func (h *ContractHandler) Preview(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Prestation introuvable"})
		return
	}

	var prestation models.Prestation
	if err := h.DB.Preload("Provider").
		Where("id = ? AND status = ? AND deleted_at IS NULL", id, "published").
		First(&prestation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Prestation introuvable"})
		return
	}

	amountCents := int64(0)
	if prestation.Price != nil {
		if v, perr := strconv.ParseFloat(*prestation.Price, 64); perr == nil {
			amountCents = int64(v * 100)
		}
	}

	providerName := "UpcycleConnect"
	providerEmail := "contact@upcycleconnect.xyz"
	if prestation.Provider != nil {
		providerName = strings.TrimSpace(prestation.Provider.FirstName + " " + prestation.Provider.LastName)
		providerEmail = prestation.Provider.Email
	}

	c.JSON(http.StatusOK, gin.H{
		"prestation": gin.H{
			"id":          prestation.ID,
			"title":       prestation.Title,
			"description": prestation.Description,
			"price_type":  prestation.PriceType,
		},
		"amount_cents": amountCents,
		"currency":     "eur",
		"customer": gin.H{
			"name":    strings.TrimSpace(user.FirstName + " " + user.LastName),
			"email":   user.Email,
			"phone":   user.Phone,
			"address": user.Address,
		},
		"provider": gin.H{
			"name":  providerName,
			"email": providerEmail,
		},
		"generated_at": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ContractHandler) QuotePreview(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}

	resID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Réservation introuvable"})
		return
	}

	var reservation models.Reservation
	if err := h.DB.First(&reservation, resID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Réservation introuvable"})
		return
	}
	if reservation.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"message": "Accès refusé"})
		return
	}
	if reservation.Status != "quote_requested" && reservation.Status != "quote_issued" {
		c.JSON(http.StatusConflict, gin.H{"message": "Ce devis ne peut plus être accepté"})
		return
	}

	var prestation models.Prestation
	if err := h.DB.Preload("Provider").First(&prestation, reservation.PrestationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Prestation introuvable"})
		return
	}

	var quote models.Invoice
	if err := h.DB.Where("reservation_id = ? AND type = ?", reservation.ID, "quote").
		Order("created_at DESC").First(&quote).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Aucun devis émis pour cette réservation"})
		return
	}
	if quote.AmountCents <= 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Devis sans montant : impossible à signer"})
		return
	}

	providerName := "UpcycleConnect"
	providerEmail := "contact@upcycleconnect.xyz"
	if prestation.Provider != nil {
		providerName = strings.TrimSpace(prestation.Provider.FirstName + " " + prestation.Provider.LastName)
		providerEmail = prestation.Provider.Email
	}

	c.JSON(http.StatusOK, gin.H{
		"prestation": gin.H{
			"id":          prestation.ID,
			"title":       prestation.Title,
			"description": prestation.Description,
			"price_type":  prestation.PriceType,
		},
		"amount_cents": quote.AmountCents,
		"currency":     quote.Currency,
		"quote_number": quote.Number,
		"customer": gin.H{
			"name":    strings.TrimSpace(user.FirstName + " " + user.LastName),
			"email":   user.Email,
			"phone":   user.Phone,
			"address": user.Address,
		},
		"provider": gin.H{
			"name":  providerName,
			"email": providerEmail,
		},
		"generated_at": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ContractHandler) Download(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Contrat introuvable"})
		return
	}

	var contract models.Contract
	if err := h.DB.First(&contract, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Contrat introuvable"})
		return
	}
	if !canAccessContract(h.DB, user, &contract) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Accès refusé"})
		return
	}
	if contract.PDFPath == nil || *contract.PDFPath == "" {
		c.JSON(http.StatusNotFound, gin.H{"message": "PDF introuvable"})
		return
	}
	if _, err := os.Stat(*contract.PDFPath); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "PDF introuvable sur le serveur"})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf(`inline; filename="%s.pdf"`, contract.Number))
	c.Header("Content-Type", "application/pdf")
	c.File(*contract.PDFPath)
}

func (h *ContractHandler) ByReservation(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	resID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Réservation introuvable"})
		return
	}

	var contract models.Contract
	if err := h.DB.Where("reservation_id = ?", resID).First(&contract).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Aucun contrat pour cette réservation"})
		return
	}
	if !canAccessContract(h.DB, user, &contract) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Accès refusé"})
		return
	}
	c.JSON(http.StatusOK, models.ToContractResponse(&contract))
}

func (h *ContractHandler) ProviderContracts(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}

	var contracts []models.Contract
	h.DB.
		Joins("JOIN reservations ON reservations.id = contracts.reservation_id").
		Joins("JOIN prestations ON prestations.id = reservations.prestation_id").
		Where("prestations.provider_id = ?", user.ID).
		Where("contracts.deleted_at IS NULL").
		Order("contracts.signed_at DESC").
		Find(&contracts)

	resp := make([]models.ContractResponse, 0, len(contracts))
	for i := range contracts {
		resp = append(resp, models.ToContractResponse(&contracts[i]))
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func canAccessContract(db *gorm.DB, user *models.User, contract *models.Contract) bool {
	if user == nil {
		return false
	}
	if user.IsAdmin() {
		return true
	}
	if contract.UserID == user.ID {
		return true
	}

	var providerID *uint
	row := db.
		Table("reservations").
		Select("prestations.provider_id").
		Joins("JOIN prestations ON prestations.id = reservations.prestation_id").
		Where("reservations.id = ?", contract.ReservationID).
		Row()
	if row != nil {
		_ = row.Scan(&providerID)
	}
	return providerID != nil && *providerID == user.ID
}

func decodeSignaturePNG(raw string) ([]byte, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("signature vide")
	}
	if idx := strings.Index(raw, "base64,"); idx >= 0 {
		raw = raw[idx+len("base64,"):]
	}
	data, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf("signature invalide: %w", err)
	}

	if len(data) < 8 || data[0] != 0x89 || data[1] != 0x50 || data[2] != 0x4E || data[3] != 0x47 {
		return nil, fmt.Errorf("signature : format PNG attendu")
	}

	if len(data) > 1024*1024 {
		return nil, fmt.Errorf("signature trop volumineuse")
	}
	return data, nil
}

func createSignedContract(
	db *gorm.DB,
	pdfSvc *services.PDFService,
	user *models.User,
	prestation *models.Prestation,
	reservationID uint,
	amountCents int64,
	currency string,
	signaturePNG []byte,
	clientIP string,
) (*models.Contract, error) {
	number := fmt.Sprintf("CONTRAT-%d-%06d", time.Now().Year(), reservationID)
	now := time.Now()

	providerName := "UpcycleConnect"
	providerEmail := "contact@upcycleconnect.xyz"
	if prestation.Provider != nil {
		providerName = strings.TrimSpace(prestation.Provider.FirstName + " " + prestation.Provider.LastName)
		providerEmail = prestation.Provider.Email
	}

	contract := models.Contract{
		Number:          number,
		UserID:          user.ID,
		ReservationID:   reservationID,
		Signature:       signaturePNG,
		SignedAt:        now,
		SignedIP:        &clientIP,
		CustomerName:    strings.TrimSpace(user.FirstName + " " + user.LastName),
		CustomerEmail:   user.Email,
		CustomerAddress: user.Address,
		CustomerPhone:   user.Phone,
		PrestationTitle: prestation.Title,
		PrestationDesc:  prestation.Description,
		ProviderName:    providerName,
		ProviderEmail:   providerEmail,
		AmountCents:     amountCents,
		Currency:        currency,
		Status:          "signed_pending_payment",
	}
	if err := db.Create(&contract).Error; err != nil {
		return nil, fmt.Errorf("save contract: %w", err)
	}

	descStr := ""
	if prestation.Description != nil {
		descStr = *prestation.Description
	}
	addrStr := ""
	if user.Address != nil {
		addrStr = *user.Address
	}
	phoneStr := ""
	if user.Phone != nil {
		phoneStr = *user.Phone
	}
	pdfPath, err := pdfSvc.GenerateContract(services.ContractData{
		Number:          number,
		SignedAt:        now,
		CustomerName:    contract.CustomerName,
		CustomerEmail:   contract.CustomerEmail,
		CustomerAddress: addrStr,
		CustomerPhone:   phoneStr,
		PrestationTitle: prestation.Title,
		PrestationDesc:  descStr,
		ProviderName:    providerName,
		ProviderEmail:   providerEmail,
		AmountCents:     amountCents,
		Currency:        currency,
		SignedIP:        clientIP,
		SignaturePNG:    signaturePNG,
	})
	if err != nil {
		log.Printf("[contract] PDF generation failed for %s: %v", number, err)

		return &contract, nil
	}

	contract.PDFPath = &pdfPath
	if err := db.Save(&contract).Error; err != nil {
		log.Printf("[contract] failed to persist pdf path for %s: %v", number, err)
	}
	return &contract, nil
}
