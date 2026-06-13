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
	DB            *gorm.DB
	Audit         *services.AuditService
	Stripe        *services.StripeService
	PDF           *services.PDFService
	Mailer        *services.Mailer
	Notifications *services.NotificationService
}


type reserveRequest struct {
	Notes     *string `json:"notes"`
	Signature string  `json:"signature"` // base64 PNG, required for paid prestations
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

	// Signature obligatoire pour toute prestation payante
	sigBytes, sigErr := decodeSignaturePNG(req.Signature)
	if sigErr != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Signature requise : " + sigErr.Error()})
		return
	}

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

	// Persist signed contract before redirecting to Stripe — proves engagement
	// even if the user abandons mid-payment.
	contract, cerr := createSignedContract(h.DB, h.PDF, user, &prestation,
		reservation.ID, amountCents, "eur", sigBytes, c.ClientIP())
	if cerr != nil {
		h.DB.Unscoped().Delete(&reservation)
		log.Printf("[reserve] contract creation failed: %v", cerr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de l'enregistrement du contrat"})
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
		"contract_id":   contract.ID,
	})

	// Notify the provider that a new reservation is incoming (pending payment).
	if prestation.ProviderID != nil {
		link := fmt.Sprintf("/profil/reservations/%d", reservation.ID)
		h.Notifications.MustNotify(*prestation.ProviderID, "reservation.received",
			"Nouvelle réservation en attente de paiement",
			fmt.Sprintf("%s a signé un contrat pour « %s ».", user.FirstName+" "+user.LastName, prestation.Title),
			link)
	}

	c.JSON(http.StatusOK, gin.H{
		"reservation_id": reservation.ID,
		"contract_id":    contract.ID,
		"contract_number": contract.Number,
		"checkout_url":   sess.URL,
		"session_id":     sess.ID,
	})
}

// AcceptQuote turns a quote_requested reservation into a signed contract + Stripe checkout.
// The client signs the same legal text as a direct purchase; the amount is locked to the issued quote.
func (h *PaymentHandler) AcceptQuote(c *gin.Context) {
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

	var req reserveRequest
	_ = c.ShouldBindJSON(&req)

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

	// Double-signature guard: refuse if a contract already exists for this reservation.
	var existing models.Contract
	if err := h.DB.Where("reservation_id = ?", reservation.ID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Un contrat existe déjà pour ce devis", "contract_id": existing.ID})
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
	amountCents := quote.AmountCents
	if amountCents <= 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Devis sans montant : impossible à signer"})
		return
	}

	sigBytes, sigErr := decodeSignaturePNG(req.Signature)
	if sigErr != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Signature requise : " + sigErr.Error()})
		return
	}

	contract, cerr := createSignedContract(h.DB, h.PDF, user, &prestation,
		reservation.ID, amountCents, quote.Currency, sigBytes, c.ClientIP())
	if cerr != nil {
		log.Printf("[accept-quote] contract creation failed: %v", cerr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de l'enregistrement du contrat"})
		return
	}

	sess, err := h.Stripe.CreateCheckoutSession(services.CheckoutParams{
		ReservationID:   reservation.ID,
		UserEmail:       user.Email,
		PrestationTitle: prestation.Title,
		AmountCents:     amountCents,
		Currency:        quote.Currency,
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "Stripe indisponible", "error": err.Error()})
		return
	}
	sessID := sess.ID
	reservation.StripeSessionID = &sessID
	reservation.Status = "pending"
	reservation.AmountCents = amountCents
	if err := h.DB.Save(&reservation).Error; err != nil {
		log.Printf("[accept-quote] failed to persist stripe_session_id for reservation %d: %v", reservation.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de l'enregistrement de la réservation"})
		return
	}

	h.Audit.Log(c, "quote.accepted", "Reservation", &reservation.ID, nil, map[string]interface{}{
		"quote_id":     quote.ID,
		"contract_id":  contract.ID,
		"amount_cents": amountCents,
	})

	if prestation.ProviderID != nil {
		h.Notifications.MustNotify(*prestation.ProviderID, "quote.accepted",
			"Devis accepté et signé",
			fmt.Sprintf("%s a accepté et signé le devis %s pour « %s ».",
				user.FirstName+" "+user.LastName, quote.Number, prestation.Title),
			fmt.Sprintf("/profil/reservations/%d", reservation.ID))
	}

	c.JSON(http.StatusOK, gin.H{
		"reservation_id":  reservation.ID,
		"contract_id":     contract.ID,
		"contract_number": contract.Number,
		"checkout_url":    sess.URL,
		"session_id":      sess.ID,
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

	h.Notifications.MustNotify(user.ID, "quote.issued",
		"Votre devis est disponible",
		fmt.Sprintf("Le devis %s pour « %s » a été émis et envoyé par email.", number, prestation.Title),
		"/profil/factures")
	if prestation.ProviderID != nil {
		h.Notifications.MustNotify(*prestation.ProviderID, "quote.requested",
			"Nouvelle demande de devis",
			fmt.Sprintf("%s a demandé un devis pour « %s ».", user.FirstName+" "+user.LastName, prestation.Title),
			"")
	}

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
		// Route based on the metadata.type set when the checkout was created.
		switch session.Metadata["type"] {
		case "subscription":
			if err := h.fulfillSubscription(&session); err != nil {
				log.Printf("[webhook] subscription fulfill failed for session=%s: %v", session.ID, err)
			}
		case "campaign":
			if err := h.fulfillCampaign(&session); err != nil {
				log.Printf("[webhook] campaign fulfill failed for session=%s: %v", session.ID, err)
			}
		default:
			if err := h.fulfillReservation(&session); err != nil {
				log.Printf("[webhook] fulfill failed for session=%s: %v", session.ID, err)
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
		}
	case "customer.subscription.deleted":
		// Subscription cancelled on Stripe's side — mark inactive locally.
		var sub stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &sub); err == nil {
			h.DB.Model(&models.Subscription{}).
				Where("stripe_subscription_id = ?", sub.ID).
				Update("status", "cancelled")
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

// fulfillSubscription activates a provider's subscription after a successful Stripe Checkout
// and generates the first monthly invoice PDF.
func (h *PaymentHandler) fulfillSubscription(session *stripe.CheckoutSession) error {
	if session.PaymentStatus != "" &&
		session.PaymentStatus != stripe.CheckoutSessionPaymentStatusPaid &&
		session.PaymentStatus != stripe.CheckoutSessionPaymentStatusNoPaymentRequired {
		return nil
	}
	userID, _ := strconv.ParseUint(session.Metadata["user_id"], 10, 64)
	plan := session.Metadata["plan"]
	if userID == 0 || plan == "" {
		return fmt.Errorf("metadata abonnement incomplète")
	}

	planDef, ok := models.SubscriptionPlans[plan]
	if !ok {
		return fmt.Errorf("plan inconnu: %s", plan)
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user introuvable: %w", err)
	}

	subID := ""
	if session.Subscription != nil {
		subID = session.Subscription.ID
	}
	custID := ""
	if session.Customer != nil {
		custID = session.Customer.ID
	}

	// Idempotency: if a subscription already exists for this stripe sub, do nothing.
	var existing models.Subscription
	if subID != "" && h.DB.Where("stripe_subscription_id = ?", subID).First(&existing).Error == nil {
		return nil
	}

	periodEnd := time.Now().AddDate(0, 1, 0)
	// Upsert: one subscription row per user.
	var sub models.Subscription
	if h.DB.Where("user_id = ?", userID).First(&sub).Error == nil {
		h.DB.Model(&sub).Updates(map[string]interface{}{
			"stripe_customer_id":     custID,
			"stripe_subscription_id": subID,
			"plan":                   plan,
			"status":                 "active",
			"current_period_end":     &periodEnd,
		})
	} else {
		sub = models.Subscription{
			UserID:               uint(userID),
			StripeCustomerID:     custID,
			StripeSubscriptionID: subID,
			Plan:                 plan,
			Status:               "active",
			CurrentPeriodEnd:     &periodEnd,
		}
		h.DB.Create(&sub)
	}

	// Generate the subscription invoice PDF.
	number := fmt.Sprintf("ABO-%d-%06d", time.Now().Year(), sub.ID)
	pdfPath, perr := h.PDF.Generate(services.InvoiceData{
		Number:          number,
		Type:            "invoice",
		IssuedAt:        time.Now(),
		CustomerName:    strings.TrimSpace(user.FirstName + " " + user.LastName),
		CustomerEmail:   user.Email,
		PrestationTitle: "Abonnement professionnel " + planDef.Label + " (mensuel)",
		AmountCents:     planDef.AmountCents,
		TVAPercent:      20.0,
		Currency:        "eur",
		Notes:           "Abonnement mensuel renouvelé automatiquement. Résiliable à tout moment depuis votre espace.",
	})
	if perr != nil {
		log.Printf("[webhook] subscription invoice PDF failed: %v", perr)
	} else {
		inv := models.Invoice{
			UserID:          uint(userID),
			Type:            "invoice",
			Number:          number,
			PrestationTitle: "Abonnement professionnel " + planDef.Label,
			AmountCents:     planDef.AmountCents,
			TVAPercent:      20.0,
			Currency:        "eur",
			PDFPath:         &pdfPath,
			Status:          "paid",
			IssuedAt:        time.Now(),
		}
		h.DB.Create(&inv)
	}

	if h.Notifications != nil {
		h.Notifications.MustNotify(uint(userID), "subscription.activated",
			"Abonnement activé",
			fmt.Sprintf("Votre abonnement %s est désormais actif. Profitez de toutes les fonctionnalités pro !", planDef.Label),
			"/profil/pro/abonnement")
	}
	return nil
}

// fulfillCampaign records the campaign payment. The campaign stays "pending"
// until an admin validates it (it then becomes "active").
func (h *PaymentHandler) fulfillCampaign(session *stripe.CheckoutSession) error {
	if session.PaymentStatus != "" &&
		session.PaymentStatus != stripe.CheckoutSessionPaymentStatusPaid &&
		session.PaymentStatus != stripe.CheckoutSessionPaymentStatusNoPaymentRequired {
		return nil
	}
	campaignID, _ := strconv.ParseUint(session.Metadata["campaign_id"], 10, 64)
	if campaignID == 0 {
		return fmt.Errorf("metadata campagne incomplète")
	}
	var camp models.Campaign
	if err := h.DB.First(&camp, campaignID).Error; err != nil {
		return fmt.Errorf("campagne introuvable: %w", err)
	}
	if camp.PaidAt != nil {
		return nil // already processed
	}
	now := time.Now()
	h.DB.Model(&camp).Updates(map[string]interface{}{
		"paid_at": &now,
		"status":  "pending", // awaiting admin validation
	})

	// Generate an invoice for the campaign payment.
	var user models.User
	h.DB.First(&user, camp.ProviderID)
	number := fmt.Sprintf("PUB-%d-%06d", time.Now().Year(), camp.ID)
	pdfPath, perr := h.PDF.Generate(services.InvoiceData{
		Number:          number,
		Type:            "invoice",
		IssuedAt:        time.Now(),
		CustomerName:    strings.TrimSpace(user.FirstName + " " + user.LastName),
		CustomerEmail:   user.Email,
		PrestationTitle: "Campagne publicitaire : " + camp.Title,
		AmountCents:     camp.BudgetCents,
		TVAPercent:      20.0,
		Currency:        "eur",
		Notes:           "Campagne soumise à validation par l'équipe UpcycleConnect avant diffusion.",
	})
	if perr == nil {
		inv := models.Invoice{
			UserID:          camp.ProviderID,
			Type:            "invoice",
			Number:          number,
			PrestationTitle: "Campagne publicitaire : " + camp.Title,
			AmountCents:     camp.BudgetCents,
			TVAPercent:      20.0,
			Currency:        "eur",
			PDFPath:         &pdfPath,
			Status:          "paid",
			IssuedAt:        time.Now(),
		}
		h.DB.Create(&inv)
	}

	if h.Notifications != nil {
		h.Notifications.NotifyAdmins("campaign.paid",
			"Campagne publicitaire à valider",
			fmt.Sprintf("%s a payé la campagne « %s ». À valider.", user.FirstName+" "+user.LastName, camp.Title),
			"/admin/monetization")
	}
	return nil
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

	// Activate the matching contract (best-effort — contract exists since Reserve created it).
	h.DB.Model(&models.Contract{}).
		Where("reservation_id = ?", reservation.ID).
		Update("status", "active")

	// Flip the reservation to paid only after the invoice is safely persisted.
	// If this Save fails, the next Stripe retry finds the invoice row and early-returns above.
	reservation.Status = "paid"
	// Freeze the platform commission (5-10%) and the provider's net share at payment time.
	commission, net := models.ComputeCommission(reservation.AmountCents)
	reservation.CommissionCents = commission
	reservation.NetCents = net
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

	// In-app notifications: client + provider
	if h.Notifications != nil {
		reservationLink := fmt.Sprintf("/profil/reservations/%d", reservation.ID)
		h.Notifications.MustNotify(reservation.UserID, "payment.confirmed",
			"Paiement confirmé",
			fmt.Sprintf("Votre paiement de %s pour « %s » a bien été reçu. Votre facture %s est disponible.",
				formatEUR(reservation.AmountCents), reservation.Prestation.Title, number),
			reservationLink)
		if reservation.Prestation.ProviderID != nil {
			h.Notifications.MustNotify(*reservation.Prestation.ProviderID, "payment.received",
				"Paiement reçu pour votre prestation",
				fmt.Sprintf("Le client %s a payé %s pour « %s ».",
					reservation.User.FirstName+" "+reservation.User.LastName,
					formatEUR(reservation.AmountCents),
					reservation.Prestation.Title),
				reservationLink)
		}
	}

	return nil
}

func formatEUR(cents int64) string {
	return fmt.Sprintf("%.2f €", float64(cents)/100)
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
