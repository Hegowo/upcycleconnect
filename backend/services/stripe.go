package services

import (
	"fmt"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
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

func (s *StripeService) VerifyWebhook(payload []byte, signature string) (stripe.Event, error) {
	if s.cfg.StripeWebhookSecret == "" {
		return stripe.Event{}, fmt.Errorf("stripe webhook secret is not configured")
	}
	return webhook.ConstructEventWithOptions(payload, signature, s.cfg.StripeWebhookSecret, webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true,
	})
}
