package billing

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
	stripelib "github.com/stripe/stripe-go/v83"
	portalsession "github.com/stripe/stripe-go/v83/billingportal/session"
	"github.com/stripe/stripe-go/v83/customer"
	"github.com/stripe/stripe-go/v83/paymentmethod"
	"github.com/stripe/stripe-go/v83/setupintent"
)

// PaymentMethodInfo represents the default card details associated with the customer's subscription
type PaymentMethodInfo struct {
	ID              string `json:"id,omitempty"`
	Brand           string `json:"brand,omitempty"`
	Last4           string `json:"last4,omitempty"`
	ExpMonth        int64  `json:"expMonth,omitempty"`
	ExpYear         int64  `json:"expYear,omitempty"`
	Expired         bool   `json:"expired"`
	ExpiringSoon    bool   `json:"expiringSoon"`
	DaysUntilExpiry int    `json:"daysUntilExpiry"`
}

// SubscriptionDetails aggregates subscription record and payment method details
type SubscriptionDetails struct {
	Subscription  *Subscription     `json:"subscription,omitempty"`
	Plan          *Plan             `json:"plan,omitempty"`
	PaymentMethod PaymentMethodInfo `json:"paymentMethod"`
}

// GetSubscriptionDetails returns current subscription plus default payment method details
func (s *Service) GetSubscriptionDetails(ctx context.Context, ws *account.Workspace) (*SubscriptionDetails, error) {
	rec, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil && !database.IsRecordNotFound(err) {
		return nil, err
	}
	var plan *Plan
	if rec != nil {
		plan, _ = s.GetPlanByID(ctx, rec.PlanID)
	}
	pmInfo := PaymentMethodInfo{}
	custID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		return &SubscriptionDetails{Subscription: rec, Plan: plan, PaymentMethod: pmInfo}, nil
	}
	c, err := customer.Get(custID, nil)
	pmID := ""
	if err == nil && c != nil && c.InvoiceSettings != nil && c.InvoiceSettings.DefaultPaymentMethod != nil {
		pmID = c.InvoiceSettings.DefaultPaymentMethod.ID
	}
	// Stripe may not immediately reflect customer invoice settings in all environments.
	// Fall back to the workspace-stored default payment method when available.
	if pmID == "" && ws.StripePaymentMethodID.Valid && ws.StripePaymentMethodID.String != "" {
		pmID = ws.StripePaymentMethodID.String
	}
	if pmID == "" {
		return &SubscriptionDetails{Subscription: rec, Plan: plan, PaymentMethod: pmInfo}, nil
	}
	pm, err := paymentmethod.Get(pmID, nil)
	if err != nil || pm == nil || pm.Card == nil {
		return &SubscriptionDetails{Subscription: rec, Plan: plan, PaymentMethod: pmInfo}, nil
	}
	pmInfo.ID = pm.ID
	pmInfo.Brand = string(pm.Card.Brand)
	pmInfo.Last4 = pm.Card.Last4
	pmInfo.ExpMonth = int64(pm.Card.ExpMonth)
	pmInfo.ExpYear = int64(pm.Card.ExpYear)
	now := time.Now()
	expTime := time.Date(int(pm.Card.ExpYear), time.Month(pm.Card.ExpMonth), 1, 0, 0, 0, 0, now.Location()).AddDate(0, 1, -1)
	days := int(expTime.Sub(now).Hours() / 24)
	pmInfo.DaysUntilExpiry = days
	pmInfo.Expired = days < 0
	pmInfo.ExpiringSoon = !pmInfo.Expired && days <= 30
	return &SubscriptionDetails{Subscription: rec, Plan: plan, PaymentMethod: pmInfo}, nil
}

// EnsureCustomer makes sure the workspace has a Stripe customer and returns it
func (s *Service) EnsureCustomer(ctx context.Context, ws *account.Workspace) (string, error) {
	if ws.StripeCustomerID.Valid && ws.StripeCustomerID.String != "" {
		if _, err := customer.Get(ws.StripeCustomerID.String, nil); err == nil {
			return ws.StripeCustomerID.String, nil
		}
		logger.FromContext(ctx).Warn("Stripe customer not found, creating new one", "workspace_id", ws.ID, "old_customer_id", ws.StripeCustomerID.String)
	}
	idempotencyKey := fmt.Sprintf("customer_%s", ws.ID)
	params := &stripelib.CustomerParams{
		Description: stripelib.String(fmt.Sprintf("Customer for workspace %s", ws.ID)),
		Metadata: map[string]string{
			"workspace_id": ws.ID,
		},
	}
	params.SetIdempotencyKey(idempotencyKey)
	c, err := customer.New(params)
	if err != nil {
		logger.FromContext(ctx).Error("Failed to create Stripe customer", "error", err, "workspace_id", ws.ID)
		return "", fmt.Errorf("failed to create customer: %w", err)
	}
	if err := s.account.SetWorkspaceStripeCustomer(ctx, ws.ID, c.ID); err != nil {
		logger.FromContext(ctx).Error("Failed to save Stripe customer ID to workspace", "error", err, "workspace_id", ws.ID, "customer_id", c.ID)
		return "", fmt.Errorf("failed to save customer ID: %w", err)
	}
	// Keep the in-memory workspace in sync within the same request context.
	// This prevents creating multiple Stripe customers for the same workspace
	// when EnsureCustomer is called more than once in a single request.
	ws.StripeCustomerID = sql.NullString{String: c.ID, Valid: true}
	logger.FromContext(ctx).Info("Created new Stripe customer", "workspace_id", ws.ID, "customer_id", c.ID)
	return c.ID, nil
}

// AttachAndSetDefaultPaymentMethod attaches a payment method to the customer and sets it as default
func (s *Service) AttachAndSetDefaultPaymentMethod(ctx context.Context, ws *account.Workspace, pmID string) error {
	if pmID == "" {
		return ErrInvalidPaymentMethod(fmt.Errorf("payment method ID cannot be empty"))
	}
	if !isStripePaymentMethodID(pmID) {
		return ErrInvalidPaymentMethod(fmt.Errorf("invalid payment method ID format"))
	}
	custID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		return fmt.Errorf("failed to ensure customer: %w", err)
	}
	pm, err := paymentmethod.Get(pmID, nil)
	if err != nil {
		logger.FromContext(ctx).Error("Failed to retrieve payment method", "error", err, "payment_method_id", pmID)
		return ErrInvalidPaymentMethod(fmt.Errorf("payment method not found: %w", err))
	}
	if pm.Type != stripelib.PaymentMethodTypeCard {
		return ErrInvalidPaymentMethod(fmt.Errorf("unsupported payment method type: %s", pm.Type))
	}
	if pm.Card == nil || pm.Card.Last4 == "" || pm.Card.ExpMonth == 0 || pm.Card.ExpYear == 0 {
		return ErrInvalidPaymentMethod(fmt.Errorf("payment method does not have valid card details"))
	}
	idempotencyKey := fmt.Sprintf("attach_pm_%s_%s", pmID, custID)
	attachParams := &stripelib.PaymentMethodAttachParams{Customer: stripelib.String(custID)}
	attachParams.SetIdempotencyKey(idempotencyKey)
	_, err = paymentmethod.Attach(pmID, attachParams)
	if err != nil {
		logger.FromContext(ctx).Error("Failed to attach payment method", "error", err, "payment_method_id", pmID, "customer_id", custID)
		return ErrInvalidPaymentMethod(fmt.Errorf("failed to attach payment method: %w", err))
	}
	updateIdempotencyKey := fmt.Sprintf("set_default_pm_%s_%s", pmID, custID)
	updateParams := &stripelib.CustomerParams{InvoiceSettings: &stripelib.CustomerInvoiceSettingsParams{DefaultPaymentMethod: stripelib.String(pmID)}}
	updateParams.SetIdempotencyKey(updateIdempotencyKey)
	_, err = customer.Update(custID, updateParams)
	if err != nil {
		logger.FromContext(ctx).Error("Failed to set default payment method", "error", err, "payment_method_id", pmID, "customer_id", custID)
		return fmt.Errorf("failed to set default payment method: %w", err)
	}
	if err := s.account.SetWorkspaceDefaultPaymentMethod(ctx, ws.ID, pmID); err != nil {
		logger.FromContext(ctx).Error("Failed to save payment method to workspace", "error", err, "workspace_id", ws.ID, "payment_method_id", pmID)
		return fmt.Errorf("failed to save payment method: %w", err)
	}
	ws.StripePaymentMethodID = sql.NullString{String: pmID, Valid: true}
	logger.FromContext(ctx).Info("Successfully attached and set default payment method", "workspace_id", ws.ID, "payment_method_id", pmID, "customer_id", custID)
	return nil
}

func isStripePaymentMethodID(id string) bool {
	if !strings.HasPrefix(id, "pm_") {
		return false
	}
	suffix := strings.TrimPrefix(id, "pm_")
	if suffix == "" {
		return false
	}
	for _, r := range suffix {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// CreateSetupIntent returns a client secret to collect and save a payment method for the workspace
func (s *Service) CreateSetupIntent(ctx context.Context, ws *account.Workspace) (string, error) {
	custID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		return "", fmt.Errorf("failed to ensure customer: %w", err)
	}
	idempotencyKey := fmt.Sprintf("setup_intent_%s_%d", ws.ID, time.Now().Unix())
	params := &stripelib.SetupIntentParams{
		Customer:           stripelib.String(custID),
		PaymentMethodTypes: []*string{stripelib.String("card")},
		Usage:              stripelib.String("off_session"),
		Metadata:           map[string]string{"workspace_id": ws.ID},
	}
	params.SetIdempotencyKey(idempotencyKey)
	si, err := setupintent.New(params)
	if err != nil {
		logger.FromContext(ctx).Error("Failed to create setup intent", "error", err, "workspace_id", ws.ID, "customer_id", custID)
		return "", fmt.Errorf("failed to create setup intent: %w", err)
	}
	logger.FromContext(ctx).Info("Created setup intent", "workspace_id", ws.ID, "customer_id", custID, "setup_intent_id", si.ID)
	return si.ClientSecret, nil
}

// CreateBillingPortalSession creates a Stripe Customer Portal session for self-service billing management
func (s *Service) CreateBillingPortalSession(ctx context.Context, ws *account.Workspace, returnURL string) (string, error) {
	custID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		return "", fmt.Errorf("failed to ensure customer: %w", err)
	}

	params := &stripelib.BillingPortalSessionParams{
		Customer:  stripelib.String(custID),
		ReturnURL: stripelib.String(returnURL),
	}

	session, err := portalsession.New(params)
	if err != nil {
		logger.FromContext(ctx).Error("Failed to create billing portal session", "error", err, "workspace_id", ws.ID, "customer_id", custID)
		return "", fmt.Errorf("failed to create billing portal session: %w", err)
	}

	logger.FromContext(ctx).Info("Created billing portal session", "workspace_id", ws.ID, "customer_id", custID)
	return session.URL, nil
}
