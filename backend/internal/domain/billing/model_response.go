package billing

import (
	"time"

	"github.com/shopspring/decimal"
)

/* Plan Response */
//------------------*/

// PlanResponse represents the API response shape for a Plan.
// It excludes GORM metadata (gorm.Model, DeletedAt) and returns clean JSON.
type PlanResponse struct {
	ID           string          `json:"id"`
	Descriptor   string          `json:"descriptor"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	StripePlanID *string         `json:"stripePlanId,omitempty"`
	Price        decimal.Decimal `json:"price"`
	Currency     string          `json:"currency"`
	BillingCycle BillingCycle    `json:"billingCycle"`
	Features     PlanFeature     `json:"features"`
	Limits       PlanLimit       `json:"limits"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
}

// ToPlanResponse converts a Plan model to a PlanResponse
func ToPlanResponse(plan *Plan) *PlanResponse {
	if plan == nil {
		return nil
	}
	return &PlanResponse{
		ID:           plan.ID,
		Descriptor:   plan.Descriptor,
		Name:         plan.Name,
		Description:  plan.Description,
		StripePlanID: plan.StripePlanID,
		Price:        plan.Price,
		Currency:     plan.Currency,
		BillingCycle: plan.BillingCycle,
		Features:     plan.Features,
		Limits:       plan.Limits,
		CreatedAt:    plan.CreatedAt,
		UpdatedAt:    plan.UpdatedAt,
	}
}

// ToPlanResponses converts a slice of Plan models to PlanResponses
func ToPlanResponses(plans []*Plan) []*PlanResponse {
	if plans == nil {
		return nil
	}
	responses := make([]*PlanResponse, len(plans))
	for i, plan := range plans {
		responses[i] = ToPlanResponse(plan)
	}
	return responses
}

/* Subscription Response */
//--------------------------*/

// SubscriptionResponse represents the API response shape for a Subscription.
// It excludes GORM metadata and returns clean JSON with optional plan details.
type SubscriptionResponse struct {
	ID               string             `json:"id"`
	WorkspaceID      string             `json:"workspaceId"`
	PlanID           string             `json:"planId"`
	Plan             *PlanResponse      `json:"plan,omitempty"`
	StripeSubID      string             `json:"stripeSubId"`
	CurrentPeriodEnd time.Time          `json:"currentPeriodEnd"`
	Status           SubscriptionStatus `json:"status"`
	CreatedAt        time.Time          `json:"createdAt"`
	UpdatedAt        time.Time          `json:"updatedAt"`
}

// ToSubscriptionResponse converts a Subscription model to a SubscriptionResponse
func ToSubscriptionResponse(subscription *Subscription) *SubscriptionResponse {
	if subscription == nil {
		return nil
	}

	resp := &SubscriptionResponse{
		ID:               subscription.ID,
		WorkspaceID:      subscription.WorkspaceID,
		PlanID:           subscription.PlanID,
		StripeSubID:      subscription.StripeSubID,
		CurrentPeriodEnd: subscription.CurrentPeriodEnd,
		Status:           subscription.Status,
		CreatedAt:        subscription.CreatedAt,
		UpdatedAt:        subscription.UpdatedAt,
	}

	// Include plan details if loaded
	if subscription.Plan != nil {
		resp.Plan = ToPlanResponse(subscription.Plan)
	}

	return resp
}

/* Invoice Summary Response */
//-----------------------------*/

// InvoiceSummary is already a clean response type (not a GORM model)
// It's defined in model_invoice_record.go and doesn't need conversion
