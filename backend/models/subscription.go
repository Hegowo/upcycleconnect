package models

import "time"

// Subscription tracks a provider's Stripe subscription (15€ basic / 30€ premium/month).
// Created when the provider completes a Stripe Checkout in subscription mode.
type Subscription struct {
	ID                   uint       `gorm:"primaryKey" json:"id"`
	UserID               uint       `gorm:"uniqueIndex;not null" json:"user_id"`
	StripeCustomerID     string     `gorm:"size:100" json:"stripe_customer_id"`
	StripeSubscriptionID string     `gorm:"size:100;uniqueIndex" json:"stripe_subscription_id"`
	Plan                 string     `gorm:"size:20;default:basic" json:"plan"` // basic | premium
	Status               string     `gorm:"size:30;default:inactive" json:"status"` // active | inactive | cancelled | past_due
	CurrentPeriodEnd     *time.Time `json:"current_period_end"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`

	User *User `gorm:"foreignKey:UserID" json:"-"`
}

func (Subscription) TableName() string { return "subscriptions" }

var SubscriptionPlans = map[string]struct {
	Label      string
	AmountCents int64
	Features   []string
}{
	"basic": {
		Label:       "Basic",
		AmountCents: 1500,
		Features: []string{
			"Accès aux annonces avec achats",
			"Récupération d'objets en conteneurs",
			"Gestion de vos projets d'upcycling",
			"Tableau de bord basique",
		},
	},
	"premium": {
		Label:       "Premium",
		AmountCents: 3000,
		Features: []string{
			"Tout ce qui est inclus dans Basic",
			"Tableaux de bord avancés",
			"Analyse d'impact écologique détaillée",
			"Statistiques sur les matériaux disponibles",
			"Alertes priorisées pour la collecte",
			"Mise en avant de vos projets",
			"Gestion de campagnes publicitaires",
		},
	},
}

type SubscriptionResponse struct {
	ID               uint    `json:"id"`
	Plan             string  `json:"plan"`
	PlanLabel        string  `json:"plan_label"`
	Status           string  `json:"status"`
	CurrentPeriodEnd *string `json:"current_period_end"`
	AmountCents      int64   `json:"amount_cents"`
}

func ToSubscriptionResponse(s *Subscription) SubscriptionResponse {
	var end *string
	if s.CurrentPeriodEnd != nil {
		v := s.CurrentPeriodEnd.UTC().Format(time.RFC3339)
		end = &v
	}
	plan := SubscriptionPlans[s.Plan]
	return SubscriptionResponse{
		ID:               s.ID,
		Plan:             s.Plan,
		PlanLabel:        plan.Label,
		Status:           s.Status,
		CurrentPeriodEnd: end,
		AmountCents:      plan.AmountCents,
	}
}
