package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"upcycleconnect/backend/middleware"
	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v76"
	"gorm.io/gorm"
)

type PaymentHandler struct {
	DB     *gorm.DB
	Audit  *services.AuditService
	Stripe *services.StripeService
	PDF    *services.PDFService
	Mailer *services.Mailer
}


type reserveRequest struct {
	Notes *string `json:"notes"`
}

func (h *PaymentHandler) Reserve(c *gin.Context) {
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

	var req reserveRequest
	_ = c.ShouldBindJSON(&req)

	if prestation.PriceType == "quote" {
		return_ := h.handleQuoteRequest(c, user, &prestation, req.Notes)
		_ = return_
		return
	}

	if prestation.Price == nil || *prestation.Price == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Cette prestation n'a pas de prix défini"})
		return
	}
	amountFloat, err := strconv.ParseFloat(*prestation.Price, 64)
	if err != nil || amountFloat <= 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Prix de la prestation invalide"})
		return
	}
	amountCents := int64(amountFloat * 100)

	reservation := models.Reservation{
		UserID:       user.ID,
		PrestationID: prestation.ID,
		Status:       "pending",
		AmountCents:  amountCents,
		Currency:     "eur",
		Notes:        req.Notes,
	}
	if err := h.DB.Create(&reservation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création de la réservation"})
		return
	}

	sess, err := h.Stripe.CreateCheckoutSession(services.CheckoutParams{
		ReservationID:   reservation.ID,
		UserEmail:       user.Email,
		PrestationTitle: prestation.Title,
		AmountCents:     amountCents,
		Currency:        "eur",
	})
	if err != nil {
		h.DB.Delete(&reservation)
		c.JSON(http.StatusBadGateway, gin.H{"message": "Stripe indisponible", "error": err.Error()})
		return
	}

	sessID := sess.ID
	reservation.StripeSessionID = &sessID
	h.DB.Save(&reservation)

	h.Audit.Log(c, "reservation.created", "Reservation", &reservation.ID, nil, map[string]interface{}{
		"prestation_id": prestation.ID,
		"amount_cents":  amountCents,
	})

	c.JSON(http.StatusOK, gin.H{
		"reservation_id": reservation.ID,
		"checkout_url":   sess.URL,
		"session_id":     sess.ID,
	})
}

func (h *PaymentHandler) handleQuoteRequest(c *gin.Context, user *models.User, prestation *models.Prestation, notes *string) error {
	reservation := models.Reservation{
		UserID:       user.ID,
		PrestationID: prestation.ID,
		Status:       "quote_requested",
		AmountCents:  0,
		Currency:     "eur",
		Notes:        notes,
	}
	if err := h.DB.Create(&reservation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création de la demande de devis"})
		return err
	}

	number := fmt.Sprintf("DEVIS-%d-%06d", time.Now().Year(), reservation.ID)
	var estimateCents int64
	if prestation.Price != nil {
		if v, err := strconv.ParseFloat(*prestation.Price, 64); err == nil {
			estimateCents = int64(v * 100)
		}
	}

	pdfPath, err := h.PDF.Generate(services.InvoiceData{
		Number:          number,
		Type:            "quote",
		IssuedAt:        time.Now(),
		CustomerName:    strings.TrimSpace(user.FirstName + " " + user.LastName),
		CustomerEmail:   user.Email,
		PrestationTitle: prestation.Title,
		AmountCents:     estimateCents,
		TVAPercent:      20.0,
		Currency:        "eur",
		Notes:           "Devis établi sur estimation. Le tarif final pourra être ajusté après analyse détaillée de votre besoin.",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la génération du devis"})
		return err
	}

	resID := reservation.ID
	invoice := models.Invoice{
		UserID:          user.ID,
		ReservationID:   &resID,
		Type:            "quote",
		Number:          number,
		PrestationTitle: prestation.Title,
		AmountCents:     estimateCents,
		TVAPercent:      20.0,
		Currency:        "eur",
		PDFPath:         &pdfPath,
		Status:          "issued",
		IssuedAt:        time.Now(),
	}
	if err := h.DB.Create(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de l'enregistrement du devis"})
		return err
	}

	go func() {
		_ = h.Mailer.SendQuote(user.Email, user.FirstName, prestation.Title, number, pdfPath)
	}()

	h.Audit.Log(c, "quote.issued", "Invoice", &invoice.ID, nil, map[string]interface{}{
		"number":        invoice.Number,
		"prestation_id": prestation.ID,
	})

	c.JSON(http.StatusOK, gin.H{
		"type":           "quote",
		"reservation_id": reservation.ID,
		"invoice_id":     invoice.ID,
		"number":         invoice.Number,
		"message":        "Votre devis a été envoyé par email et est disponible dans votre espace Mes factures.",
	})
	return nil
}

func (h *PaymentHandler) Webhook(c *gin.Context) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Payload invalide"})
		return
	}
	signature := c.GetHeader("Stripe-Signature")

	event, err := h.Stripe.VerifyWebhook(payload, signature)
	if err != nil {
		log.Printf("[webhook] verification failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Signature invalide", "error": err.Error()})
		return
	}
	log.Printf("[webhook] event received: type=%s id=%s", event.Type, event.ID)

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Payload session invalide"})
			return
		}
		if err := h.fulfillReservation(&session); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	case "checkout.session.expired", "checkout.session.async_payment_failed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err == nil {
			h.DB.Model(&models.Reservation{}).
				Where("stripe_session_id = ?", session.ID).
				Update("status", "failed")
		}
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}

func (h *PaymentHandler) fulfillReservation(session *stripe.CheckoutSession) error {
	var reservation models.Reservation
	if err := h.DB.Preload("User").Preload("Prestation").
		Where("stripe_session_id = ?", session.ID).
		First(&reservation).Error; err != nil {
		return fmt.Errorf("reservation introuvable: %w", err)
	}

	if reservation.Status == "paid" {
		return nil
	}

	reservation.Status = "paid"
	if session.PaymentIntent != nil {
		pi := session.PaymentIntent.ID
		reservation.StripePaymentIntentID = &pi
	}
	h.DB.Save(&reservation)

	if reservation.User == nil || reservation.Prestation == nil {
		return fmt.Errorf("données de réservation incomplètes")
	}

	number := fmt.Sprintf("FACT-%d-%06d", time.Now().Year(), reservation.ID)
	pdfPath, err := h.PDF.Generate(services.InvoiceData{
		Number:          number,
		Type:            "invoice",
		IssuedAt:        time.Now(),
		CustomerName:    strings.TrimSpace(reservation.User.FirstName + " " + reservation.User.LastName),
		CustomerEmail:   reservation.User.Email,
		PrestationTitle: reservation.Prestation.Title,
		AmountCents:     reservation.AmountCents,
		TVAPercent:      20.0,
		Currency:        reservation.Currency,
	})
	if err != nil {
		return fmt.Errorf("génération PDF: %w", err)
	}

	resID := reservation.ID
	invoice := models.Invoice{
		UserID:          reservation.UserID,
		ReservationID:   &resID,
		Type:            "invoice",
		Number:          number,
		PrestationTitle: reservation.Prestation.Title,
		AmountCents:     reservation.AmountCents,
		TVAPercent:      20.0,
		Currency:        reservation.Currency,
		PDFPath:         &pdfPath,
		Status:          "issued",
		IssuedAt:        time.Now(),
	}
	if err := h.DB.Create(&invoice).Error; err != nil {
		return fmt.Errorf("enregistrement facture: %w", err)
	}

	go func() {
		_ = h.Mailer.SendPaymentConfirmation(
			reservation.User.Email,
			reservation.User.FirstName,
			reservation.Prestation.Title,
			number,
			reservation.AmountCents,
			pdfPath,
		)
	}()

	return nil
}

func (h *PaymentHandler) SessionStatus(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	sessionID := c.Query("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "session_id requis"})
		return
	}

	var reservation models.Reservation
	if err := h.DB.Preload("Prestation").
		Where("user_id = ? AND stripe_session_id = ?", user.ID, sessionID).
		First(&reservation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Réservation introuvable"})
		return
	}

	c.JSON(http.StatusOK, models.ToReservationResponse(&reservation))
}
