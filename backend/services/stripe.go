package services

import (
	"fmt"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	stripeSub "github.com/stripe/stripe-go/v76/subscription"
	"github.com/stripe/stripe-go/v76/webhook"

	"upcycleconnect/backend/config"
)

type StripeService struct {
	cfg *config.Config
}

func NewStripeService(cfg *config.Config) *StripeService {
	stripe.Key = cfg.StripeSecret
	return &StripeService{cfg: cfg}
}

type CheckoutParams struct {
	ReservationID   uint
	UserEmail       string
	PrestationTitle string
	AmountCents     int64
	Currency        string
}

func (s *StripeService) CreateCheckoutSession(p CheckoutParams) (*stripe.CheckoutSession, error) {
	if s.cfg.StripeSecret == "" {
		return nil, fmt.Errorf("Clé secrète stripe pas configurée")
	}
	stripe.Key = s.cfg.StripeSecret

	currency := p.Currency
	if currency == "" {
		currency = "eur"
	}

	successURL := fmt.Sprintf("%s/paiement/succes?session_id={CHECKOUT_SESSION_ID}", s.cfg.AppURL)
	cancelURL := fmt.Sprintf("%s/paiement/annule", s.cfg.AppURL)

	params := &stripe.CheckoutSessionParams{
		Mode:              stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:        stripe.String(successURL),
		CancelURL:         stripe.String(cancelURL),
		CustomerEmail:     stripe.String(p.UserEmail),
		ClientReferenceID: stripe.String(fmt.Sprintf("%d", p.ReservationID)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(currency),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(p.PrestationTitle),
					},
					UnitAmount: stripe.Int64(p.AmountCents),
				},
				Quantity: stripe.Int64(1),
			},
		},
	}
	params.AddMetadata("reservation_id", fmt.Sprintf("%d", p.ReservationID))

	return session.New(params)
}

type SubscriptionCheckoutParams struct {
	UserID      uint
	UserEmail   string
	Plan        string
	Label       string
	AmountCents int64
}

func (s *StripeService) CreateSubscriptionCheckout(p SubscriptionCheckoutParams) (*stripe.CheckoutSession, error) {
	if s.cfg.StripeSecret == "" || s.cfg.StripeSecret == "sk_test_local_dummy" {
		return nil, fmt.Errorf("Stripe non configuré — ajoute ta clé dans le .env pour activer les abonnements")
	}
	stripe.Key = s.cfg.StripeSecret
	successURL := fmt.Sprintf("%s/profil/pro/abonnement?success=1&session_id={CHECKOUT_SESSION_ID}", s.cfg.AppURL)
	cancelURL := fmt.Sprintf("%s/profil/pro/abonnement?cancelled=1", s.cfg.AppURL)
	params := &stripe.CheckoutSessionParams{
		Mode:          stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL:    stripe.String(successURL),
		CancelURL:     stripe.String(cancelURL),
		CustomerEmail: stripe.String(p.UserEmail),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("eur"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name:        stripe.String("UpcycleConnect " + p.Label),
						Description: stripe.String("Abonnement mensuel professionnel/artisan"),
					},
					Recurring: &stripe.CheckoutSessionLineItemPriceDataRecurringParams{
						Interval: stripe.String("month"),
					},
					UnitAmount: stripe.Int64(p.AmountCents),
				},
				Quantity: stripe.Int64(1),
			},
		},
	}
	params.AddMetadata("user_id", fmt.Sprintf("%d", p.UserID))
	params.AddMetadata("plan", p.Plan)
	params.AddMetadata("type", "subscription")
	return session.New(params)
}

func (s *StripeService) CancelSubscription(stripeSubID string) error {
	if s.cfg.StripeSecret == "" || s.cfg.StripeSecret == "sk_test_local_dummy" {
		return fmt.Errorf("Stripe non configuré")
	}
	stripe.Key = s.cfg.StripeSecret
	params := &stripe.SubscriptionParams{}
	params.AddExpand("latest_invoice")
	_, err := stripeSub.Cancel(stripeSubID, nil)
	return err
}

func (s *StripeService) CreateCampaignCheckout(userID, campaignID uint, userEmail, title string, amountCents int64) (*stripe.CheckoutSession, error) {
	if s.cfg.StripeSecret == "" || s.cfg.StripeSecret == "sk_test_local_dummy" {
		return nil, fmt.Errorf("Stripe non configuré — ajoute ta clé dans le .env")
	}
	stripe.Key = s.cfg.StripeSecret
	successURL := fmt.Sprintf("%s/profil/pro/campagnes?success=1&session_id={CHECKOUT_SESSION_ID}", s.cfg.AppURL)
	cancelURL := fmt.Sprintf("%s/profil/pro/campagnes?cancelled=1", s.cfg.AppURL)
	params := &stripe.CheckoutSessionParams{
		Mode:          stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:    stripe.String(successURL),
		CancelURL:     stripe.String(cancelURL),
		CustomerEmail: stripe.String(userEmail),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("eur"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Campagne publicitaire : " + title),
					},
					UnitAmount: stripe.Int64(amountCents),
				},
				Quantity: stripe.Int64(1),
			},
		},
	}
	params.AddMetadata("user_id", fmt.Sprintf("%d", userID))
	params.AddMetadata("campaign_id", fmt.Sprintf("%d", campaignID))
	params.AddMetadata("type", "campaign")
	return session.New(params)
}

func (s *StripeService) CreateDepositCheckout(providerID, depositID uint, userEmail, title string, amountCents int64) (*stripe.CheckoutSession, error) {
	if s.cfg.StripeSecret == "" || s.cfg.StripeSecret == "sk_test_local_dummy" {
		return nil, fmt.Errorf("Stripe non configuré — ajoute ta clé dans le .env")
	}
	stripe.Key = s.cfg.StripeSecret
	successURL := fmt.Sprintf("%s/profil/pro/achats?success=1&session_id={CHECKOUT_SESSION_ID}", s.cfg.AppURL)
	cancelURL := fmt.Sprintf("%s/profil/pro/achats?cancelled=1", s.cfg.AppURL)
	params := &stripe.CheckoutSessionParams{
		Mode:          stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:    stripe.String(successURL),
		CancelURL:     stripe.String(cancelURL),
		CustomerEmail: stripe.String(userEmail),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("eur"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Achat objet : " + title),
					},
					UnitAmount: stripe.Int64(amountCents),
				},
				Quantity: stripe.Int64(1),
			},
		},
	}
	params.AddMetadata("provider_id", fmt.Sprintf("%d", providerID))
	params.AddMetadata("deposit_id", fmt.Sprintf("%d", depositID))
	params.AddMetadata("type", "deposit_purchase")
	return session.New(params)
}

func (s *StripeService) VerifyWebhook(payload []byte, signature string) (stripe.Event, error) {
	if s.cfg.StripeWebhookSecret == "" {
		return stripe.Event{}, fmt.Errorf("stripe webhook secret is not configured")
	}
	return webhook.ConstructEventWithOptions(payload, signature, s.cfg.StripeWebhookSecret, webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true,
	})
}
