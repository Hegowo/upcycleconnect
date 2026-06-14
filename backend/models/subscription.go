package models

import "time"

type Subscription struct {
	ID                   uint       `gorm:"primaryKey" json:"id"`
	UserID               uint       `gorm:"uniqueIndex;not null" json:"user_id"`
	StripeCustomerID     string     `gorm:"size:100" json:"stripe_customer_id"`
	StripeSubscriptionID string     `gorm:"size:100;uniqueIndex" json:"stripe_subscription_id"`
	Plan                 string     `gorm:"size:20;default:basic" json:"plan"`
	Status               string     `gorm:"size:30;default:inactive" json:"status"`
	CurrentPeriodEnd     *time.Time `json:"current_period_end"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`

	User *User `gorm:"foreignKey:UserID" json:"-"`
}

func (Subscription) TableName() string { return "subscriptions" }

var SubscriptionPlans = map[string]struct {
	Label       string
	AmountCents int64
	Features    []string
}{
	"basic": {
		Label:       "Basic",
		AmountCents: 1500,
		Features: []string{
			"Tableaux de bord avancés (chiffre d'affaires, réservations)",
			"Statistiques sur les matériaux disponibles",
		},
	},
	"premium": {
		Label:       "Premium",
		AmountCents: 3000,
		Features: []string{
			"Tout ce qui est inclus dans Basic",
			"Analyse d'impact écologique détaillée",
			"Alertes priorisées pour la collecte",
			"Gestion de campagnes publicitaires (mise en avant)",
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
