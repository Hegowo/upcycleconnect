package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
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
	amountCents := int64(math.Round(amountFloat * 100))

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
		h.DB.Unscoped().Delete(&reservation)
		c.JSON(http.StatusBadGateway, gin.H{"message": "Stripe indisponible", "error": err.Error()})
		return
	}

	sessID := sess.ID
	reservation.StripeSessionID = &sessID
	if err := h.DB.Save(&reservation).Error; err != nil {
		log.Printf("[reserve] failed to persist stripe_session_id for reservation %d: %v", reservation.ID, err)
		h.DB.Unscoped().Delete(&reservation)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de l'enregistrement de la réservation"})
		return
	}

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
			estimateCents = int64(math.Round(v * 100))
		} else {
			log.Printf("[quote] prestation %d has unparseable price %q: %v", prestation.ID, *prestation.Price, err)
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
	case "checkout.session.completed", "checkout.session.async_payment_succeeded":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Payload session invalide"})
			return
		}
		if err := h.fulfillReservation(&session); err != nil {
			log.Printf("[webhook] fulfill failed for session=%s: %v", session.ID, err)
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
	// For async payment methods (SEPA Debit, Bancontact, bank redirects),
	// checkout.session.completed fires before funds clear. Wait for the
	// async_payment_succeeded event (which comes with payment_status=paid).
	if session.PaymentStatus != "" &&
		session.PaymentStatus != stripe.CheckoutSessionPaymentStatusPaid &&
		session.PaymentStatus != stripe.CheckoutSessionPaymentStatusNoPaymentRequired {
		log.Printf("[webhook] session %s not yet paid (status=%s); awaiting async settlement", session.ID, session.PaymentStatus)
		return nil
	}

	reservation, err := h.lookupReservation(session)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Permanent miss: webhook misconfigured or Save side-effect was lost.
			// Returning nil => 200 to Stripe, which stops retries.
			log.Printf("[webhook] no reservation for session=%s cref=%q — acking to stop retries", session.ID, session.ClientReferenceID)
			return nil
		}
		return fmt.Errorf("lookup réservation: %w", err)
	}

	// Idempotency: use invoice-row existence as the source of truth.
	// A retry that runs after the invoice is already stored returns cleanly;
	// a retry after a crash *before* the invoice was stored re-runs the full path.
	var existing models.Invoice
	switch err := h.DB.Where("reservation_id = ? AND type = ?", reservation.ID, "invoice").
		First(&existing).Error; {
	case err == nil:
		return nil
	case !errors.Is(err, gorm.ErrRecordNotFound):
		return fmt.Errorf("lookup facture existante: %w", err)
	}

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

	// Flip the reservation to paid only after the invoice is safely persisted.
	// If this Save fails, the next Stripe retry finds the invoice row and early-returns above.
	reservation.Status = "paid"
	if session.PaymentIntent != nil {
		pi := session.PaymentIntent.ID
		reservation.StripePaymentIntentID = &pi
	}
	if err := h.DB.Save(reservation).Error; err != nil {
		log.Printf("[webhook] invoice %s stored but reservation %d status save failed: %v", number, reservation.ID, err)
		return fmt.Errorf("mise à jour réservation: %w", err)
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

// lookupReservation finds the reservation for a Stripe checkout session. It
// first tries stripe_session_id, then falls back to ClientReferenceID — which
// protects against a silent failure of the Save that persists the session ID
// in Reserve().
func (h *PaymentHandler) lookupReservation(session *stripe.CheckoutSession) (*models.Reservation, error) {
	var reservation models.Reservation
	err := h.DB.Preload("User").Preload("Prestation").
		Where("stripe_session_id = ?", session.ID).
		First(&reservation).Error
	if err == nil {
		return &reservation, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if session.ClientReferenceID == "" {
		return nil, gorm.ErrRecordNotFound
	}
	resID, perr := strconv.ParseUint(session.ClientReferenceID, 10, 64)
	if perr != nil {
		return nil, gorm.ErrRecordNotFound
	}
	if err := h.DB.Preload("User").Preload("Prestation").
		First(&reservation, resID).Error; err != nil {
		return nil, err
	}
	if reservation.StripeSessionID == nil || *reservation.StripeSessionID == "" {
		sessID := session.ID
		reservation.StripeSessionID = &sessID
		if serr := h.DB.Save(&reservation).Error; serr != nil {
			log.Printf("[webhook] backfill stripe_session_id failed for reservation %d: %v", reservation.ID, serr)
		}
	}
	return &reservation, nil
}

func (h *PaymentHandler) ShowReservation(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Réservation introuvable"})
		return
	}

	var reservation models.Reservation
	if err := h.DB.
		Preload("Prestation").
		Preload("Prestation.Category").
		Preload("Prestation.Provider").
		Where("id = ? AND user_id = ?", id, user.ID).
		First(&reservation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Réservation introuvable"})
		return
	}

	resp := models.ToReservationResponse(&reservation)

	var invoice models.Invoice
	if err := h.DB.Where("reservation_id = ? AND user_id = ?", reservation.ID, user.ID).
		Order("issued_at DESC").
		First(&invoice).Error; err == nil {
		inv := models.ToInvoiceResponse(&invoice)
		c.JSON(http.StatusOK, gin.H{
			"reservation": resp,
			"invoice":     inv,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reservation": resp,
		"invoice":     nil,
	})
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
