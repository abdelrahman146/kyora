package billing

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/bus"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/atomic"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"github.com/shopspring/decimal"
	stripelib "github.com/stripe/stripe-go/v83"
	portalsession "github.com/stripe/stripe-go/v83/billingportal/session"
	checkoutsession "github.com/stripe/stripe-go/v83/checkout/session"
	"github.com/stripe/stripe-go/v83/creditnote"
	"github.com/stripe/stripe-go/v83/customer"
	"github.com/stripe/stripe-go/v83/invoice"
	"github.com/stripe/stripe-go/v83/invoiceitem"
	"github.com/stripe/stripe-go/v83/paymentmethod"
	"github.com/stripe/stripe-go/v83/price"
	"github.com/stripe/stripe-go/v83/product"
	"github.com/stripe/stripe-go/v83/setupintent"
	"github.com/stripe/stripe-go/v83/subscription"
	"github.com/stripe/stripe-go/v83/subscriptionschedule"
	"github.com/stripe/stripe-go/v83/tax/calculation"
	"github.com/stripe/stripe-go/v83/tax/settings"
)

type Service struct {
	storage         *Storage
	atomicProcessor atomic.AtomicProcessor
	bus             *bus.Bus
	account         *account.Service
}

func NewService(storage *Storage, atomicProcessor atomic.AtomicProcessor, bus *bus.Bus, accountSvc *account.Service) *Service {
	// Note: Stripe is used via package-level helpers using API key configured globally if needed in future.
	s := &Service{
		storage:         storage,
		atomicProcessor: atomicProcessor,
		bus:             bus,
		account:         accountSvc,
	}
	// Best-effort background plan sync on service creation
	go func() {
		// use background context to avoid blocking init
		if err := s.SyncPlansToStripe(context.Background()); err != nil {
			slog.Error("plan sync to stripe failed", "error", err)
		}
	}()
	return s
}

// InvoiceSummary is a lightweight view of a Stripe invoice for UI consumption
type InvoiceSummary struct {
	ID               string     `json:"id"`
	Number           string     `json:"number"`
	Status           string     `json:"status"`
	Currency         string     `json:"currency"`
	AmountDue        int64      `json:"amountDue"`
	AmountPaid       int64      `json:"amountPaid"`
	CreatedAt        time.Time  `json:"createdAt"`
	DueDate          *time.Time `json:"dueDate,omitempty"`
	HostedInvoiceURL string     `json:"hostedInvoiceUrl,omitempty"`
	InvoicePDF       string     `json:"invoicePdf,omitempty"`
}

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

func (s *Service) GetPlanByDescriptor(ctx context.Context, descriptor string) (*Plan, error) {
	return s.storage.plan.FindOne(ctx, s.storage.plan.ScopeEquals(PlanSchema.Descriptor, descriptor))
}

func (s *Service) GetPlanByID(ctx context.Context, id string) (*Plan, error) {
	return s.storage.plan.FindByID(ctx, id)
}

func (s *Service) ListPlans(ctx context.Context) ([]*Plan, error) {
	return s.storage.plan.FindMany(ctx)
}

func (s *Service) GetSubscriptionByID(ctx context.Context, id string) (*Subscription, error) {
	return s.storage.subscription.FindByID(ctx, id)
}

func (s *Service) GetSubscriptionByWorkspaceID(ctx context.Context, workspaceID string) (*Subscription, error) {
	return s.storage.subscription.FindOne(ctx, s.storage.subscription.ScopeWorkspaceID(workspaceID), s.storage.subscription.WithPreload(PlanStruct))
}

// EnsureCustomer makes sure the workspace has a Stripe customer and returns it
func (s *Service) EnsureCustomer(ctx context.Context, ws *account.Workspace) (string, error) {
	if ws.StripeCustomerID.Valid && ws.StripeCustomerID.String != "" {
		// Verify the customer still exists in Stripe
		if _, err := customer.Get(ws.StripeCustomerID.String, nil); err == nil {
			return ws.StripeCustomerID.String, nil
		}
		// Customer doesn't exist in Stripe anymore, need to create a new one
		slog.Warn("Stripe customer not found, creating new one", "workspace_id", ws.ID, "old_customer_id", ws.StripeCustomerID.String)
	}

	// Use workspace ID as idempotency key to prevent duplicate customers
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
		slog.Error("Failed to create Stripe customer", "error", err, "workspace_id", ws.ID)
		return "", fmt.Errorf("failed to create customer: %w", err)
	}

	if err := s.account.SetWorkspaceStripeCustomer(ctx, ws.ID, c.ID); err != nil {
		slog.Error("Failed to save Stripe customer ID to workspace", "error", err, "workspace_id", ws.ID, "customer_id", c.ID)
		return "", fmt.Errorf("failed to save customer ID: %w", err)
	}

	slog.Info("Created new Stripe customer", "workspace_id", ws.ID, "customer_id", c.ID)
	return c.ID, nil
}

// AttachAndSetDefaultPaymentMethod attaches a payment method to the customer and sets it as default
func (s *Service) AttachAndSetDefaultPaymentMethod(ctx context.Context, ws *account.Workspace, pmID string) error {
	if pmID == "" {
		return ErrInvalidPaymentMethod(fmt.Errorf("payment method ID cannot be empty"))
	}

	custID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		return fmt.Errorf("failed to ensure customer: %w", err)
	}

	// First verify the payment method exists and is valid
	pm, err := paymentmethod.Get(pmID, nil)
	if err != nil {
		slog.Error("Failed to retrieve payment method", "error", err, "payment_method_id", pmID)
		return ErrInvalidPaymentMethod(fmt.Errorf("payment method not found: %w", err))
	}

	// Validate payment method type (only allow cards for now)
	if pm.Type != stripelib.PaymentMethodTypeCard {
		return ErrInvalidPaymentMethod(fmt.Errorf("unsupported payment method type: %s", pm.Type))
	}

	// Attach payment method to customer with idempotency
	idempotencyKey := fmt.Sprintf("attach_pm_%s_%s", pmID, custID)
	attachParams := &stripelib.PaymentMethodAttachParams{
		Customer: stripelib.String(custID),
	}
	attachParams.SetIdempotencyKey(idempotencyKey)

	_, err = paymentmethod.Attach(pmID, attachParams)
	if err != nil {
		slog.Error("Failed to attach payment method", "error", err, "payment_method_id", pmID, "customer_id", custID)
		return ErrInvalidPaymentMethod(fmt.Errorf("failed to attach payment method: %w", err))
	}

	// Set as default payment method on customer
	updateIdempotencyKey := fmt.Sprintf("set_default_pm_%s_%s", pmID, custID)
	updateParams := &stripelib.CustomerParams{
		InvoiceSettings: &stripelib.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripelib.String(pmID),
		},
	}
	updateParams.SetIdempotencyKey(updateIdempotencyKey)

	_, err = customer.Update(custID, updateParams)
	if err != nil {
		slog.Error("Failed to set default payment method", "error", err, "payment_method_id", pmID, "customer_id", custID)
		return fmt.Errorf("failed to set default payment method: %w", err)
	}

	// Save to local database
	if err := s.account.SetWorkspaceDefaultPaymentMethod(ctx, ws.ID, pmID); err != nil {
		slog.Error("Failed to save payment method to workspace", "error", err, "workspace_id", ws.ID, "payment_method_id", pmID)
		return fmt.Errorf("failed to save payment method: %w", err)
	}

	slog.Info("Successfully attached and set default payment method", "workspace_id", ws.ID, "payment_method_id", pmID, "customer_id", custID)
	return nil
}

// CreateSetupIntent returns a client secret to collect and save a payment method for the workspace
func (s *Service) CreateSetupIntent(ctx context.Context, ws *account.Workspace) (string, error) {
	custID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		return "", fmt.Errorf("failed to ensure customer: %w", err)
	}

	// Use workspace ID as idempotency key to prevent duplicate setup intents
	idempotencyKey := fmt.Sprintf("setup_intent_%s_%d", ws.ID, time.Now().Unix())

	params := &stripelib.SetupIntentParams{
		Customer:           stripelib.String(custID),
		PaymentMethodTypes: []*string{stripelib.String("card")},
		Usage:              stripelib.String("off_session"),
		Metadata: map[string]string{
			"workspace_id": ws.ID,
		},
	}
	params.SetIdempotencyKey(idempotencyKey)

	si, err := setupintent.New(params)
	if err != nil {
		slog.Error("Failed to create setup intent", "error", err, "workspace_id", ws.ID, "customer_id", custID)
		return "", fmt.Errorf("failed to create setup intent: %w", err)
	}

	slog.Info("Created setup intent", "workspace_id", ws.ID, "customer_id", custID, "setup_intent_id", si.ID)
	return si.ClientSecret, nil
}

// CreateCheckoutSession creates a Stripe Checkout Session for subscription signup or changes
// This is the recommended approach for payment collection as per Stripe best practices
func (s *Service) CreateCheckoutSession(ctx context.Context, ws *account.Workspace, plan *Plan, successURL, cancelURL string) (string, error) {
	custID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		return "", fmt.Errorf("failed to ensure customer: %w", err)
	}

	// Check if customer already has an active subscription
	existing, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil && !database.IsRecordNotFound(err) {
		return "", fmt.Errorf("failed to check existing subscription: %w", err)
	}

	// Determine checkout mode based on whether it's new or update
	mode := stripelib.CheckoutSessionModeSubscription
	var lineItems []*stripelib.CheckoutSessionLineItemParams

	if plan.Price.IsZero() {
		// For free plans, use setup mode to save payment method for future use
		mode = stripelib.CheckoutSessionModeSetup
	} else {
		// For paid plans, create subscription
		lineItems = []*stripelib.CheckoutSessionLineItemParams{
			{
				Price:    stripelib.String(plan.StripePlanID),
				Quantity: stripelib.Int64(1),
			},
		}
	}

	// Use workspace and plan as idempotency key
	idempotencyKey := fmt.Sprintf("checkout_%s_%s_%d", ws.ID, plan.ID, time.Now().Unix())

	params := &stripelib.CheckoutSessionParams{
		Customer:   stripelib.String(custID),
		Mode:       stripelib.String(string(mode)),
		LineItems:  lineItems,
		SuccessURL: stripelib.String(successURL),
		CancelURL:  stripelib.String(cancelURL),
		Metadata: map[string]string{
			"workspace_id": ws.ID,
			"plan_id":      plan.ID,
		},
		PaymentMethodTypes:       []*string{stripelib.String("card")},
		BillingAddressCollection: stripelib.String("auto"),
	}

	// For subscription mode, handle existing subscription updates
	if mode == stripelib.CheckoutSessionModeSubscription && existing != nil && existing.Status == SubscriptionStatusActive {
		// This is an update to existing subscription
		params.SubscriptionData = &stripelib.CheckoutSessionSubscriptionDataParams{
			Metadata: map[string]string{
				"workspace_id": ws.ID,
				"plan_id":      plan.ID,
			},
		}
	}

	params.SetIdempotencyKey(idempotencyKey)

	session, err := checkoutsession.New(params)
	if err != nil {
		slog.Error("Failed to create checkout session", "error", err, "workspace_id", ws.ID, "plan_id", plan.ID)
		return "", fmt.Errorf("failed to create checkout session: %w", err)
	}

	slog.Info("Created checkout session", "workspace_id", ws.ID, "plan_id", plan.ID, "session_id", session.ID)
	return session.URL, nil
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
		slog.Error("Failed to create billing portal session", "error", err, "workspace_id", ws.ID, "customer_id", custID)
		return "", fmt.Errorf("failed to create billing portal session: %w", err)
	}

	slog.Info("Created billing portal session", "workspace_id", ws.ID, "customer_id", custID)
	return session.URL, nil
}

// CreateOrUpdateSubscription creates a new subscription or updates existing to new plan with proration
// This method now includes proper error handling, validation, and follows Stripe best practices
func (s *Service) CreateOrUpdateSubscription(ctx context.Context, ws *account.Workspace, plan *Plan) (*Subscription, error) {
	// Validate inputs
	if ws == nil {
		return nil, fmt.Errorf("workspace cannot be nil")
	}
	if plan == nil {
		return nil, fmt.Errorf("plan cannot be nil")
	}

	// One active subscription per workspace
	existing, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil && !database.IsRecordNotFound(err) {
		return nil, fmt.Errorf("failed to check existing subscription: %w", err)
	}

	// Don't allow changing to the same plan if already active
	if existing != nil && existing.PlanID == plan.ID && existing.Status == SubscriptionStatusActive {
		return existing, nil // Return existing instead of error
	}

	custID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure customer: %w", err)
	}

	var result *Subscription

	err = s.atomicProcessor.Exec(ctx, func(ctx context.Context) error {
		// If we have an existing subscription, update it; else create a new one
		var stripeSub *stripelib.Subscription
		if existing != nil {
			// Validate downgrade protection
			if currentPlan, err := s.GetPlanByID(ctx, existing.PlanID); err == nil {
				if plan.Price.LessThan(currentPlan.Price) && existing.Status == SubscriptionStatusActive {
					// Feature-based compatibility check
					if err := s.ensureFeatureCompatibility(currentPlan, plan); err != nil {
						return err
					}
					// Usage-aware checks across modules
					if err := s.ensureWithinNewPlanLimits(ctx, ws.ID, plan); err != nil {
						return err
					}
				}
			}

			// Update existing subscription with proper idempotency
			idempotencyKey := fmt.Sprintf("sub_update_%s_%s", existing.StripeSubID, plan.ID)
			updateParams := &stripelib.SubscriptionParams{
				Items: []*stripelib.SubscriptionItemsParams{
					{
						Price: stripelib.String(plan.StripePlanID),
					},
				},
				ProrationBehavior: stripelib.String("create_prorations"),
				CancelAtPeriodEnd: stripelib.Bool(false),
				Metadata: map[string]string{
					"workspace_id": ws.ID,
					"plan_id":      plan.ID,
				},
			}
			updateParams.SetIdempotencyKey(idempotencyKey)

			stripeSub, err = subscription.Update(existing.StripeSubID, updateParams)
			if err != nil {
				slog.Error("Failed to update Stripe subscription", "error", err, "subscription_id", existing.StripeSubID, "plan_id", plan.ID)
				return fmt.Errorf("failed to update subscription: %w", err)
			}

			// Update local record
			existing.PlanID = plan.ID
			existing.Status = mapStripeStatus(stripeSub.Status)
			// CurrentPeriodEnd will be updated via webhook events for accuracy
			if err := s.storage.subscription.UpdateOne(ctx, existing); err != nil {
				return fmt.Errorf("failed to update local subscription: %w", err)
			}

			result = existing
			slog.Info("Updated subscription", "workspace_id", ws.ID, "subscription_id", existing.StripeSubID, "new_plan_id", plan.ID)
			return nil
		}

		// Create new subscription with proper configuration
		idempotencyKey := fmt.Sprintf("sub_create_%s_%s", ws.ID, plan.ID)
		createParams := &stripelib.SubscriptionParams{
			Customer: stripelib.String(custID),
			Items: []*stripelib.SubscriptionItemsParams{
				{Price: stripelib.String(plan.StripePlanID)},
			},
			ProrationBehavior: stripelib.String("create_prorations"),
			CancelAtPeriodEnd: stripelib.Bool(false),
			Metadata: map[string]string{
				"workspace_id": ws.ID,
				"plan_id":      plan.ID,
			},
		}

		// Configure payment behavior based on plan type
		if plan.Price.IsZero() {
			// For free plans, allow creation without payment method
			createParams.PaymentBehavior = stripelib.String("allow_incomplete")
		} else {
			// For paid plans, require valid payment method
			createParams.PaymentBehavior = stripelib.String("default_incomplete")
			// Set collection method to charge automatically
			createParams.CollectionMethod = stripelib.String("charge_automatically")
		}

		createParams.SetIdempotencyKey(idempotencyKey)

		stripeSub, err = subscription.New(createParams)
		if err != nil {
			slog.Error("Failed to create Stripe subscription", "error", err, "customer_id", custID, "plan_id", plan.ID)
			return fmt.Errorf("failed to create subscription: %w", err)
		}

		// Create local subscription record
		newSub := &Subscription{
			WorkspaceID:      ws.ID,
			PlanID:           plan.ID,
			StripeSubID:      stripeSub.ID,
			Status:           mapStripeStatus(stripeSub.Status),
			CurrentPeriodEnd: time.Now(), // Will be updated via webhook events for accuracy
		}

		if err := s.storage.subscription.CreateOne(ctx, newSub); err != nil {
			return fmt.Errorf("failed to create local subscription: %w", err)
		}

		result = newSub
		slog.Info("Created new subscription", "workspace_id", ws.ID, "subscription_id", stripeSub.ID, "plan_id", plan.ID)
		return nil
	}, atomic.WithRetries(2))

	if err != nil {
		return nil, err
	}

	return result, nil
}

// CancelSubscriptionImmediately cancels subscription now with proper error handling and atomic updates
func (s *Service) CancelSubscriptionImmediately(ctx context.Context, ws *account.Workspace) error {
	if ws == nil {
		return fmt.Errorf("workspace cannot be nil")
	}

	subRec, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		return ErrSubscriptionNotFound(err, ws.ID)
	}

	if subRec.Status == SubscriptionStatusCanceled {
		return nil // Already canceled
	}

	return s.atomicProcessor.Exec(ctx, func(ctx context.Context) error {
		// Cancel at Stripe with idempotency
		idempotencyKey := fmt.Sprintf("cancel_%s", subRec.StripeSubID)
		cancelParams := &stripelib.SubscriptionCancelParams{
			InvoiceNow: stripelib.Bool(false),
			Prorate:    stripelib.Bool(false), // Proration handled by webhook
		}
		cancelParams.SetIdempotencyKey(idempotencyKey)

		_, err = subscription.Cancel(subRec.StripeSubID, cancelParams)
		if err != nil {
			slog.Error("Failed to cancel Stripe subscription", "error", err, "subscription_id", subRec.StripeSubID, "workspace_id", ws.ID)
			// Don't return error immediately - still update local state
		}

		// Update local record
		subRec.Status = SubscriptionStatusCanceled
		subRec.CurrentPeriodEnd = time.Now()

		if updateErr := s.storage.subscription.UpdateOne(ctx, subRec); updateErr != nil {
			slog.Error("Failed to update local subscription status", "error", updateErr, "subscription_id", subRec.StripeSubID)
			return fmt.Errorf("failed to update local subscription: %w", updateErr)
		}

		if err != nil {
			return fmt.Errorf("failed to cancel Stripe subscription: %w", err)
		}

		slog.Info("Successfully canceled subscription", "workspace_id", ws.ID, "subscription_id", subRec.StripeSubID)
		return nil
	}, atomic.WithRetries(2))
}

func mapStripeStatus(s stripelib.SubscriptionStatus) SubscriptionStatus {
	switch s {
	case stripelib.SubscriptionStatusActive:
		return SubscriptionStatusActive
	case stripelib.SubscriptionStatusPastDue:
		return SubscriptionStatusPastDue
	case stripelib.SubscriptionStatusUnpaid:
		return SubscriptionStatusUnpaid
	case stripelib.SubscriptionStatusIncomplete:
		return SubscriptionStatusIncomplete
	case stripelib.SubscriptionStatusCanceled:
		return SubscriptionStatusCanceled
	default:
		return SubscriptionStatusIncomplete
	}
}

// ensureFeatureCompatibility prevents downgrades that remove features used by the current plan.
// Conservative rule: if a feature is enabled on the current plan but disabled on the new plan,
// block the downgrade. This can be relaxed later with usage-aware feature checks per module.
func (s *Service) ensureFeatureCompatibility(currentPlan, newPlan *Plan) error {
	// Compare all boolean feature flags
	cur := currentPlan.Features
	nxt := newPlan.Features
	// If current has a feature and next does not, disallow downgrade
	if cur.CustomerManagement && !nxt.CustomerManagement {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.InventoryManagement && !nxt.InventoryManagement {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.OrderManagement && !nxt.OrderManagement {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.ExpenseManagement && !nxt.ExpenseManagement {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.AssetsManagement && !nxt.AssetsManagement {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.Accounting && !nxt.Accounting {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.BasicAnalytics && !nxt.BasicAnalytics {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.FinancialReports && !nxt.FinancialReports {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.DataImport && !nxt.DataImport {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.DataExport && !nxt.DataExport {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.AdvancedAnalytics && !nxt.AdvancedAnalytics {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.AdvancedFinancialReports && !nxt.AdvancedFinancialReports {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.OrderPaymentLinks && !nxt.OrderPaymentLinks {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.InvoiceGeneration && !nxt.InvoiceGeneration {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.ExportAnalyticsData && !nxt.ExportAnalyticsData {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.AIBusinessAssistant && !nxt.AIBusinessAssistant {
		return ErrCannotDowngradePlan(nil)
	}
	return nil
}

// SyncSubscriptionStatus updates the local record based on Stripe status
func (s *Service) SyncSubscriptionStatus(ctx context.Context, stripeSubID string, status string, periodStart, periodEnd int64) {
	rec, err := s.storage.subscription.FindOne(ctx, s.storage.subscription.ScopeEquals(SubscriptionSchema.StripeSubID, stripeSubID))
	if err != nil || rec == nil {
		return
	}
	rec.Status = mapStripeStatus(stripelib.SubscriptionStatus(status))
	if periodEnd > 0 {
		rec.CurrentPeriodEnd = time.Unix(periodEnd, 0)
	}
	_ = s.storage.subscription.UpdateOne(ctx, rec)
}

// MarkSubscriptionPastDue sets subscription status to past_due
func (s *Service) MarkSubscriptionPastDue(ctx context.Context, stripeSubID string) {
	rec, err := s.storage.subscription.FindOne(ctx, s.storage.subscription.ScopeEquals(SubscriptionSchema.StripeSubID, stripeSubID))
	if err != nil || rec == nil {
		return
	}
	rec.Status = SubscriptionStatusPastDue
	_ = s.storage.subscription.UpdateOne(ctx, rec)
}

// MarkSubscriptionActive sets subscription status to active
func (s *Service) MarkSubscriptionActive(ctx context.Context, stripeSubID string) {
	rec, err := s.storage.subscription.FindOne(ctx, s.storage.subscription.ScopeEquals(SubscriptionSchema.StripeSubID, stripeSubID))
	if err != nil || rec == nil {
		return
	}
	rec.Status = SubscriptionStatusActive
	_ = s.storage.subscription.UpdateOne(ctx, rec)
}

// RefundAndFinalizeCancellation computes prorated refund and cancels in Stripe, then updates local DB
func (s *Service) RefundAndFinalizeCancellation(ctx context.Context, stripeSubID string, periodStart, periodEnd int64) {
	rec, err := s.storage.subscription.FindOne(ctx, s.storage.subscription.ScopeEquals(SubscriptionSchema.StripeSubID, stripeSubID))
	if err != nil || rec == nil {
		return
	}
	// If period bounds are not provided by the event, skip refund calculation to avoid incorrect credits
	// (We still cancel immediately; future improvements can derive period from invoice line items.)
	// Get latest paid invoice for this subscription
	ip := &stripelib.InvoiceListParams{Subscription: stripelib.String(stripeSubID)}
	ip.Status = stripelib.String(string(stripelib.InvoiceStatusPaid))
	ip.Limit = stripelib.Int64(1)
	iter := invoice.List(ip)
	if iter.Next() {
		inv := iter.Invoice()
		// Only compute prorated refund if both periodStart and periodEnd are provided
		if periodStart > 0 && periodEnd > 0 {
			// Compute prorated refund based on remaining time
			now := time.Now()
			pEnd := time.Unix(periodEnd, 0)
			pStart := time.Unix(periodStart, 0)
			if pEnd.Before(now) {
				pEnd = rec.CurrentPeriodEnd
			}
			if !(pStart.After(now) || pEnd.Before(pStart)) {
				total := pEnd.Sub(pStart).Seconds()
				remaining := pEnd.Sub(now).Seconds()
				if total > 0 && remaining > 0 && inv.AmountPaid > 0 {
					refundAmount := int64(float64(inv.AmountPaid) * (remaining / total))
					if refundAmount > 0 {
						// Idempotency: skip if a cancel_prorated credit note already exists for this invoice
						lparams := &stripelib.CreditNoteListParams{Invoice: stripelib.String(inv.ID)}
						lparams.Limit = stripelib.Int64(10)
						alreadyRefunded := false
						itCN := creditnote.List(lparams)
						for itCN.Next() {
							cn := itCN.CreditNote()
							if cn != nil && cn.Metadata != nil {
								if v, ok := cn.Metadata["kyoraRefundKind"]; ok && v == "cancel_prorated" {
									alreadyRefunded = true
									break
								}
							}
						}
						if !alreadyRefunded {
							// Issue a credit note on the invoice with a refund for the prorated amount
							_, rerr := creditnote.New(&stripelib.CreditNoteParams{
								Invoice:      stripelib.String(inv.ID),
								Amount:       stripelib.Int64(refundAmount),
								RefundAmount: stripelib.Int64(refundAmount),
								Reason:       stripelib.String(string(stripelib.CreditNoteReasonOrderChange)),
								Memo:         stripelib.String("Prorated refund for immediate cancellation"),
								Metadata: map[string]string{
									"kyoraRefundKind": "cancel_prorated",
									"subscription":    stripeSubID,
									"workspaceId":     rec.WorkspaceID,
								},
							})
							if rerr != nil {
								slog.Error("failed to create credit note refund", "error", rerr, "invoice", inv.ID)
							}
						} else {
							slog.Info("skip credit note refund: already applied", "invoice", inv.ID, "subscription", stripeSubID)
						}
					}
				}
			}
		} else {
			slog.Info("Skipping refund calculation: missing period bounds", "subscription", stripeSubID)
		}
	}
	// Cancel at Stripe immediately without additional proration (credit note already applied)
	_, _ = subscription.Cancel(stripeSubID, &stripelib.SubscriptionCancelParams{InvoiceNow: stripelib.Bool(false), Prorate: stripelib.Bool(false)})
	// Update local record
	rec.Status = SubscriptionStatusCanceled
	rec.CurrentPeriodEnd = time.Now()
	_ = s.storage.subscription.UpdateOne(ctx, rec)
}

// ListInvoices returns invoice summaries for the workspace's customer
func (s *Service) ListInvoices(ctx context.Context, ws *account.Workspace, status string) ([]InvoiceSummary, error) {
	custID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		return nil, err
	}
	params := &stripelib.InvoiceListParams{Customer: stripelib.String(custID)}
	switch status {
	case string(stripelib.InvoiceStatusOpen):
		params.Status = stripelib.String(string(stripelib.InvoiceStatusOpen))
	case string(stripelib.InvoiceStatusPaid):
		params.Status = stripelib.String(string(stripelib.InvoiceStatusPaid))
	default:
		// all - no status filter
	}
	params.Limit = stripelib.Int64(50)
	it := invoice.List(params)
	res := make([]InvoiceSummary, 0)
	for it.Next() {
		inv := it.Invoice()
		var due *time.Time
		if inv.DueDate != 0 {
			t := time.Unix(inv.DueDate, 0)
			due = &t
		}
		res = append(res, InvoiceSummary{
			ID:               inv.ID,
			Number:           inv.Number,
			Status:           string(inv.Status),
			Currency:         string(inv.Currency),
			AmountDue:        inv.AmountDue,
			AmountPaid:       inv.AmountPaid,
			CreatedAt:        time.Unix(inv.Created, 0),
			DueDate:          due,
			HostedInvoiceURL: inv.HostedInvoiceURL,
			InvoicePDF:       inv.InvoicePDF,
		})
	}
	return res, nil
}

// DownloadInvoiceURL returns the downloadable PDF URL for an invoice if it belongs to the customer's workspace
func (s *Service) DownloadInvoiceURL(ctx context.Context, ws *account.Workspace, invoiceID string) (string, error) {
	custID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		return "", err
	}
	inv, err := invoice.Get(invoiceID, nil)
	if err != nil {
		return "", err
	}
	if inv.Customer == nil || inv.Customer.ID != custID {
		return "", ErrSubscriptionNotFound(fmt.Errorf("invoice not owned by workspace"), ws.ID)
	}
	if inv.InvoicePDF == "" && inv.HostedInvoiceURL == "" {
		return "", ErrSubscriptionNotFound(fmt.Errorf("invoice has no downloadable link"), invoiceID)
	}
	if inv.InvoicePDF != "" {
		return inv.InvoicePDF, nil
	}
	return inv.HostedInvoiceURL, nil
}

// PayInvoice attempts to pay an open invoice for the workspace's customer
func (s *Service) PayInvoice(ctx context.Context, ws *account.Workspace, invoiceID string) error {
	custID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		return err
	}
	inv, err := invoice.Get(invoiceID, nil)
	if err != nil {
		return err
	}
	if inv.Customer == nil || inv.Customer.ID != custID {
		return ErrSubscriptionNotFound(fmt.Errorf("invoice not owned by workspace"), ws.ID)
	}
	// If invoice is draft, finalize first
	if inv.Status == stripelib.InvoiceStatusDraft {
		if _, err := invoice.FinalizeInvoice(invoiceID, nil); err != nil {
			return err
		}
	}
	_, err = invoice.Pay(invoiceID, &stripelib.InvoicePayParams{})
	return err
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
	// Pull customer default payment method
	pmInfo := PaymentMethodInfo{}
	custID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		return &SubscriptionDetails{Subscription: rec, Plan: plan, PaymentMethod: pmInfo}, nil
	}
	c, err := customer.Get(custID, nil)
	if err != nil || c == nil || c.InvoiceSettings == nil || c.InvoiceSettings.DefaultPaymentMethod == nil {
		return &SubscriptionDetails{Subscription: rec, Plan: plan, PaymentMethod: pmInfo}, nil
	}
	dpm := c.InvoiceSettings.DefaultPaymentMethod
	pm, err := paymentmethod.Get(dpm.ID, nil)
	if err != nil || pm == nil || pm.Card == nil {
		return &SubscriptionDetails{Subscription: rec, Plan: plan, PaymentMethod: pmInfo}, nil
	}
	pmInfo.ID = pm.ID
	pmInfo.Brand = string(pm.Card.Brand)
	pmInfo.Last4 = pm.Card.Last4
	pmInfo.ExpMonth = int64(pm.Card.ExpMonth)
	pmInfo.ExpYear = int64(pm.Card.ExpYear)
	// compute expiry status
	now := time.Now()
	expTime := time.Date(int(pm.Card.ExpYear), time.Month(pm.Card.ExpMonth), 1, 0, 0, 0, 0, now.Location()).AddDate(0, 1, -1) // end of exp month
	days := int(expTime.Sub(now).Hours() / 24)
	pmInfo.DaysUntilExpiry = days
	pmInfo.Expired = days < 0
	pmInfo.ExpiringSoon = !pmInfo.Expired && days <= 30
	return &SubscriptionDetails{Subscription: rec, Plan: plan, PaymentMethod: pmInfo}, nil
}

// SyncPlansToStripe ensures all local plans exist in Stripe as products/prices with proper conflict resolution
func (s *Service) SyncPlansToStripe(ctx context.Context) error {
	plans, err := s.storage.plan.FindMany(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch plans: %w", err)
	}

	slog.Info("Starting plan sync to Stripe", "plan_count", len(plans))

	for _, p := range plans {
		if err := s.syncSinglePlanToStripe(ctx, p); err != nil {
			slog.Error("Failed to sync plan", "error", err, "plan_id", p.ID, "descriptor", p.Descriptor)
			// Continue with other plans instead of failing entire sync
			continue
		}
	}

	slog.Info("Completed plan sync to Stripe")
	return nil
}

// syncSinglePlanToStripe handles syncing a single plan with proper error handling and conflict resolution
func (s *Service) syncSinglePlanToStripe(ctx context.Context, p *Plan) error {
	// Try to find existing product by metadata
	prod, err := s.findOrCreateProduct(ctx, p)
	if err != nil {
		return fmt.Errorf("failed to find or create product: %w", err)
	}

	// Ensure price exists and is correct
	needNewPrice, err := s.validateExistingPrice(p, prod.ID)
	if err != nil {
		return fmt.Errorf("failed to validate existing price: %w", err)
	}

	if needNewPrice {
		newPrice, err := s.createPrice(p, prod.ID)
		if err != nil {
			return fmt.Errorf("failed to create new price: %w", err)
		}

		// Update local plan with new price ID
		p.StripePlanID = newPrice.ID
		if err := s.storage.plan.UpdateOne(ctx, p); err != nil {
			slog.Error("Failed to update plan with new Stripe price ID", "error", err, "plan_id", p.ID, "price_id", newPrice.ID)
			return fmt.Errorf("failed to update plan: %w", err)
		}

		slog.Info("Created new price for plan", "plan_id", p.ID, "price_id", newPrice.ID, "amount", p.Price)
	}

	return nil
}

// findOrCreateProduct finds existing product by metadata or creates new one
func (s *Service) findOrCreateProduct(ctx context.Context, p *Plan) (*stripelib.Product, error) {
	// First try to find by existing price if we have one
	if p.StripePlanID != "" {
		if pr, err := price.Get(p.StripePlanID, nil); err == nil && pr != nil && pr.Product != nil {
			if prod, err := product.Get(pr.Product.ID, nil); err == nil {
				// Verify metadata matches
				if prod.Metadata != nil {
					if kyoraID, exists := prod.Metadata["kyora_plan_id"]; exists && kyoraID == p.ID {
						return prod, nil
					}
					if descriptor, exists := prod.Metadata["descriptor"]; exists && descriptor == p.Descriptor {
						return prod, nil
					}
				}
				// Product exists but metadata doesn't match - update it
				return s.updateProductMetadata(prod, p)
			}
		}
	}

	// Search for existing product by metadata (this requires listing all products)
	// Note: This is expensive for large numbers of products, consider caching in production
	existingProd, err := s.findProductByMetadata(p)
	if err != nil {
		return nil, fmt.Errorf("failed to search for existing product: %w", err)
	}

	if existingProd != nil {
		return existingProd, nil
	}

	// Create new product
	return s.createProduct(p)
}

// findProductByMetadata searches for an existing product with matching metadata
func (s *Service) findProductByMetadata(p *Plan) (*stripelib.Product, error) {
	// List products to find by metadata (limit to reasonable number)
	params := &stripelib.ProductListParams{
		Active: stripelib.Bool(true),
	}
	params.Limit = stripelib.Int64(100) // Adjust based on expected product count

	iter := product.List(params)
	for iter.Next() {
		prod := iter.Product()
		if prod.Metadata != nil {
			if kyoraID, exists := prod.Metadata["kyora_plan_id"]; exists && kyoraID == p.ID {
				return prod, nil
			}
			if descriptor, exists := prod.Metadata["descriptor"]; exists && descriptor == p.Descriptor {
				return prod, nil
			}
		}
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return nil, nil // Not found
}

// createProduct creates a new Stripe product
func (s *Service) createProduct(p *Plan) (*stripelib.Product, error) {
	idempotencyKey := fmt.Sprintf("product_%s", p.ID)
	params := &stripelib.ProductParams{
		Name:        stripelib.String(p.Name),
		Description: stripelib.String(p.Description),
		Metadata: map[string]string{
			"kyora_plan_id": p.ID,
			"descriptor":    p.Descriptor,
		},
	}
	params.SetIdempotencyKey(idempotencyKey)

	prod, err := product.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe product: %w", err)
	}

	slog.Info("Created new Stripe product", "plan_id", p.ID, "product_id", prod.ID, "name", p.Name)
	return prod, nil
}

// updateProductMetadata updates an existing product's metadata
func (s *Service) updateProductMetadata(prod *stripelib.Product, p *Plan) (*stripelib.Product, error) {
	params := &stripelib.ProductParams{
		Name:        stripelib.String(p.Name),
		Description: stripelib.String(p.Description),
		Metadata: map[string]string{
			"kyora_plan_id": p.ID,
			"descriptor":    p.Descriptor,
		},
	}

	updatedProd, err := product.Update(prod.ID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update product metadata: %w", err)
	}

	slog.Info("Updated product metadata", "plan_id", p.ID, "product_id", prod.ID)
	return updatedProd, nil
}

// validateExistingPrice checks if the current price matches the plan requirements
func (s *Service) validateExistingPrice(p *Plan, productID string) (bool, error) {
	if p.StripePlanID == "" {
		return true, nil // Need new price
	}

	existingPrice, err := price.Get(p.StripePlanID, nil)
	if err != nil {
		slog.Warn("Failed to fetch existing price, will create new one", "price_id", p.StripePlanID, "error", err)
		return true, nil
	}

	// Check if price properties match
	interval := "month"
	if p.BillingCycle == BillingCycleYearly {
		interval = "year"
	}

	expectedAmount := p.Price.Mul(decimal.NewFromInt(100)).IntPart()

	// Validate all price properties
	if string(existingPrice.Currency) != p.Currency {
		slog.Info("Price currency mismatch", "expected", p.Currency, "actual", existingPrice.Currency)
		return true, nil
	}

	if existingPrice.Recurring == nil {
		slog.Info("Price missing recurring configuration")
		return true, nil
	}

	if string(existingPrice.Recurring.Interval) != interval {
		slog.Info("Price interval mismatch", "expected", interval, "actual", existingPrice.Recurring.Interval)
		return true, nil
	}

	if existingPrice.UnitAmount != expectedAmount {
		slog.Info("Price amount mismatch", "expected", expectedAmount, "actual", existingPrice.UnitAmount)
		return true, nil
	}

	// Verify product association
	if existingPrice.Product == nil || existingPrice.Product.ID != productID {
		slog.Info("Price associated with wrong product", "expected_product", productID, "actual_product", existingPrice.Product)
		return true, nil
	}

	return false, nil // Price is valid
}

// createPrice creates a new Stripe price
func (s *Service) createPrice(p *Plan, productID string) (*stripelib.Price, error) {
	interval := "month"
	if p.BillingCycle == BillingCycleYearly {
		interval = "year"
	}

	unitAmount := p.Price.Mul(decimal.NewFromInt(100)).IntPart()
	idempotencyKey := fmt.Sprintf("price_%s_%s_%d", p.ID, interval, unitAmount)

	params := &stripelib.PriceParams{
		Currency:   stripelib.String(p.Currency),
		UnitAmount: stripelib.Int64(unitAmount),
		Recurring:  &stripelib.PriceRecurringParams{Interval: stripelib.String(interval)},
		Product:    stripelib.String(productID),
		Metadata: map[string]string{
			"kyora_plan_id": p.ID,
			"descriptor":    p.Descriptor,
		},
	}
	params.SetIdempotencyKey(idempotencyKey)

	newPrice, err := price.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe price: %w", err)
	}

	return newPrice, nil
}

// ensureWithinNewPlanLimits enforces usage within new plan limits (users, businesses, monthly orders)
func (s *Service) ensureWithinNewPlanLimits(ctx context.Context, workspaceID string, newPlan *Plan) error {
	// Users
	users, err := s.account.CountWorkspaceUsers(ctx, workspaceID)
	if err != nil {
		return err
	}
	if users > newPlan.Limits.MaxTeamMembers {
		return ErrCannotDowngradePlan(nil)
	}
	// Businesses
	businesses, err := s.storage.CountBusinessesByWorkspace(ctx, workspaceID)
	if err != nil {
		return err
	}
	if businesses > newPlan.Limits.MaxBusinesses {
		return ErrCannotDowngradePlan(nil)
	}
	// Monthly orders
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthEnd := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location()).Add(-time.Nanosecond)
	orders, err := s.storage.CountMonthlyOrdersByWorkspace(ctx, workspaceID, monthStart, monthEnd)
	if err != nil {
		return err
	}
	if orders > newPlan.Limits.MaxOrdersPerMonth {
		return ErrCannotDowngradePlan(nil)
	}
	return nil
}

// ResumeSubscriptionIfNoDue attempts to pay open invoices then recreates a subscription with the same plan
func (s *Service) ResumeSubscriptionIfNoDue(ctx context.Context, ws *account.Workspace) (*Subscription, error) {
	// fetch current (canceled) record to get plan
	rec, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		return nil, ErrSubscriptionNotFound(err, ws.ID)
	}
	if rec.Status != SubscriptionStatusCanceled {
		return rec, nil
	}
	// attempt to pay open invoices
	custID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		return nil, err
	}
	ip := &stripelib.InvoiceListParams{Customer: stripelib.String(custID)}
	ip.Status = stripelib.String(string(stripelib.InvoiceStatusOpen))
	ip.Limit = stripelib.Int64(10)
	it := invoice.List(ip)
	for it.Next() {
		inv := it.Invoice()
		if _, err := invoice.Pay(inv.ID, &stripelib.InvoicePayParams{}); err != nil {
			return nil, ErrSubscriptionNotActive(err)
		}
	}
	// all due cleared — re-subscribe to same plan
	plan, err := s.GetPlanByID(ctx, rec.PlanID)
	if err != nil {
		return nil, err
	}
	return s.CreateOrUpdateSubscription(ctx, ws, plan)
}

// CalculateTax calculates tax for a given amount and workspace location
func (s *Service) CalculateTax(ctx context.Context, ws *account.Workspace, amount int64, currency string) (*stripelib.TaxCalculation, error) {
	idempotencyKey := id.KsuidWithPrefix("tax_calc")

	slog.InfoContext(ctx, "Calculating tax",
		"workspace_id", ws.ID,
		"amount", amount,
		"currency", currency,
		"idempotency_key", idempotencyKey,
	)

	// Prepare line items for tax calculation
	lineItems := []*stripelib.TaxCalculationLineItemParams{
		{
			Amount:    stripelib.Int64(amount),
			Reference: stripelib.String("service_charge"),
			TaxCode:   stripelib.String("txcd_10000000"), // Generic service tax code
		},
	}

	params := &stripelib.TaxCalculationParams{
		Currency:  stripelib.String(currency),
		LineItems: lineItems,
		CustomerDetails: &stripelib.TaxCalculationCustomerDetailsParams{
			AddressSource: stripelib.String("billing"),
		},
	}

	// Set idempotency key for safe retry
	params.SetIdempotencyKey(idempotencyKey)

	calc, err := calculation.New(params)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to calculate tax",
			"error", err.Error(),
			"amount", amount,
			"currency", currency,
			"idempotency_key", idempotencyKey,
		)
		return nil, ErrStripeOperationFailed(err, "calculate_tax")
	}

	slog.InfoContext(ctx, "Tax calculation completed",
		"calculation_id", calc.ID,
		"amount_total", calc.AmountTotal,
		"tax_amount_exclusive", calc.TaxAmountExclusive,
		"tax_amount_inclusive", calc.TaxAmountInclusive,
	)

	return calc, nil
}

// UpdateTaxSettings updates the account's tax settings for automatic tax calculation
func (s *Service) UpdateTaxSettings(ctx context.Context, defaultTaxCode string) error {
	idempotencyKey := id.KsuidWithPrefix("tax_settings")

	slog.InfoContext(ctx, "Updating tax settings",
		"default_tax_code", defaultTaxCode,
		"idempotency_key", idempotencyKey,
	)

	params := &stripelib.TaxSettingsParams{
		Defaults: &stripelib.TaxSettingsDefaultsParams{
			TaxCode: stripelib.String(defaultTaxCode),
		},
	}

	// Set idempotency key for safe retry
	params.SetIdempotencyKey(idempotencyKey)

	taxSettings, err := settings.Update(params)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to update tax settings",
			"error", err.Error(),
			"default_tax_code", defaultTaxCode,
			"idempotency_key", idempotencyKey,
		)
		return ErrStripeOperationFailed(err, "update_tax_settings")
	}

	slog.InfoContext(ctx, "Tax settings updated successfully",
		"status", taxSettings.Status,
		"default_tax_code", taxSettings.Defaults.TaxCode,
	)

	return nil
}

// TrackUsage is a helper method to track API calls, storage, or other metered usage
func (s *Service) TrackUsage(ctx context.Context, ws *account.Workspace, usageType string, quantity int64) error {
	slog.InfoContext(ctx, "Tracking usage",
		"workspace_id", ws.ID,
		"usage_type", usageType,
		"quantity", quantity,
	)

	// Get current subscription to find the appropriate metered item
	subscription, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		return fmt.Errorf("failed to get subscription for usage tracking: %w", err)
	}

	// Since our Subscription model doesn't have Stripe metadata,
	// we'll track usage locally and sync with Stripe periodically
	// For now, log the usage for manual tracking or batch processing
	slog.InfoContext(ctx, "Usage tracked locally",
		"workspace_id", ws.ID,
		"subscription_id", subscription.ID,
		"usage_type", usageType,
		"quantity", quantity,
		"timestamp", time.Now().Unix(),
	)

	// TODO: Store usage in local database for batch processing
	// This would typically involve:
	// 1. Store usage record in local usage_records table
	// 2. Batch process usage records to Stripe hourly/daily
	// 3. Handle usage-based pricing tiers and quotas

	return nil
}

// GetUsageQuota checks current usage against plan limits
func (s *Service) GetUsageQuota(ctx context.Context, ws *account.Workspace, usageType string) (used int64, limit int64, err error) {
	slog.InfoContext(ctx, "Checking usage quota",
		"workspace_id", ws.ID,
		"usage_type", usageType,
	)

	subscription, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get subscription for quota check: %w", err)
	}

	plan, err := s.GetPlanByID(ctx, subscription.PlanID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get plan for quota check: %w", err)
	}

	// Get quota limits from plan limits
	var quotaLimit int64
	switch usageType {
	case "orders_per_month":
		quotaLimit = plan.Limits.MaxOrdersPerMonth
	case "team_members":
		quotaLimit = plan.Limits.MaxTeamMembers
	case "businesses":
		quotaLimit = plan.Limits.MaxBusinesses
	default:
		return 0, 0, fmt.Errorf("unsupported usage type: %s", usageType)
	}

	// TODO: Calculate current usage from local usage_records table
	// For now, return placeholder values
	currentUsage := int64(0) // This would be calculated from stored usage records

	slog.InfoContext(ctx, "Usage quota retrieved",
		"workspace_id", ws.ID,
		"usage_type", usageType,
		"current_usage", currentUsage,
		"quota_limit", quotaLimit,
	)

	return currentUsage, quotaLimit, nil
}

// CheckUsageLimit verifies if usage is within plan limits
func (s *Service) CheckUsageLimit(ctx context.Context, ws *account.Workspace, usageType string, additionalUsage int64) error {
	current, limit, err := s.GetUsageQuota(ctx, ws, usageType)
	if err != nil {
		return fmt.Errorf("failed to check usage quota: %w", err)
	}

	// If limit is -1, it means unlimited
	if limit == -1 {
		return nil
	}

	if current+additionalUsage > limit {
		return fmt.Errorf("usage limit exceeded for %s: current %d + additional %d > limit %d",
			usageType, current, additionalUsage, limit)
	}

	return nil
}

// CreateTrialSubscription creates a subscription with a trial period
func (s *Service) CreateTrialSubscription(ctx context.Context, ws *account.Workspace, plan *Plan, trialDays int) (*stripelib.Subscription, error) {
	idempotencyKey := id.KsuidWithPrefix("trial_sub")

	slog.InfoContext(ctx, "Creating trial subscription",
		"workspace_id", ws.ID,
		"plan_id", plan.ID,
		"trial_days", trialDays,
		"idempotency_key", idempotencyKey,
	)

	// Ensure customer exists in Stripe
	customerID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure customer exists: %w", err)
	}

	// Calculate trial end date
	trialEnd := time.Now().AddDate(0, 0, trialDays).Unix()

	params := &stripelib.SubscriptionParams{
		Customer: stripelib.String(customerID),
		Items: []*stripelib.SubscriptionItemsParams{
			{
				Price: stripelib.String(plan.StripePlanID),
			},
		},
		TrialEnd:          stripelib.Int64(trialEnd),
		CollectionMethod:  stripelib.String("charge_automatically"),
		PaymentBehavior:   stripelib.String("default_incomplete"),
		ProrationBehavior: stripelib.String("none"),
	}

	// Expand latest invoice for payment intent details
	params.Expand = []*string{
		stripelib.String("latest_invoice.payment_intent"),
	}

	// Add metadata for tracking
	params.Metadata = map[string]string{
		"workspace_id": ws.ID,
		"plan_id":      plan.ID,
		"trial":        "true",
		"trial_days":   fmt.Sprintf("%d", trialDays),
	}

	// Set idempotency key for safe retry
	params.SetIdempotencyKey(idempotencyKey)

	sub, err := subscription.New(params)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create trial subscription",
			"error", err.Error(),
			"workspace_id", ws.ID,
			"plan_id", plan.ID,
			"trial_days", trialDays,
			"idempotency_key", idempotencyKey,
		)
		return nil, ErrStripeOperationFailed(err, "create_trial_subscription")
	}

	slog.InfoContext(ctx, "Trial subscription created successfully",
		"subscription_id", sub.ID,
		"customer_id", customerID,
		"trial_end", trialEnd,
		"status", sub.Status,
	)

	return sub, nil
}

// ExtendTrialPeriod extends the trial period for an existing subscription
func (s *Service) ExtendTrialPeriod(ctx context.Context, ws *account.Workspace, additionalDays int) error {
	idempotencyKey := id.KsuidWithPrefix("extend_trial")

	slog.InfoContext(ctx, "Extending trial period",
		"workspace_id", ws.ID,
		"additional_days", additionalDays,
		"idempotency_key", idempotencyKey,
	)

	// Get current subscription
	sub, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	// Check if subscription is still in trial
	if sub.Status != SubscriptionStatusTrialing {
		return fmt.Errorf("subscription is not in trial period, current status: %s", sub.Status)
	}

	// Get Stripe subscription to get current trial end
	stripeSub, err := subscription.Get(sub.StripeSubID, nil)
	if err != nil {
		return ErrStripeOperationFailed(err, "get_subscription")
	}

	// Calculate new trial end
	currentTrialEnd := time.Unix(stripeSub.TrialEnd, 0)
	newTrialEnd := currentTrialEnd.AddDate(0, 0, additionalDays).Unix()

	// Update subscription trial end
	params := &stripelib.SubscriptionParams{
		TrialEnd: stripelib.Int64(newTrialEnd),
	}
	params.SetIdempotencyKey(idempotencyKey)

	_, err = subscription.Update(sub.StripeSubID, params)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to extend trial period",
			"error", err.Error(),
			"subscription_id", sub.StripeSubID,
			"additional_days", additionalDays,
			"idempotency_key", idempotencyKey,
		)
		return ErrStripeOperationFailed(err, "extend_trial")
	}

	slog.InfoContext(ctx, "Trial period extended successfully",
		"subscription_id", sub.StripeSubID,
		"new_trial_end", newTrialEnd,
		"additional_days", additionalDays,
	)

	return nil
}

// HandleGracePeriod manages grace periods for failed payments
func (s *Service) HandleGracePeriod(ctx context.Context, ws *account.Workspace, graceDays int) error {
	idempotencyKey := id.KsuidWithPrefix("grace_period")

	slog.InfoContext(ctx, "Handling grace period",
		"workspace_id", ws.ID,
		"grace_days", graceDays,
		"idempotency_key", idempotencyKey,
	)

	// Get current subscription
	sub, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	// Check if subscription has unpaid invoices
	if sub.Status != SubscriptionStatusPastDue {
		return fmt.Errorf("subscription is not past due, current status: %s", sub.Status)
	}

	// Calculate grace period end
	gracePeriodEnd := time.Now().AddDate(0, 0, graceDays).Unix()

	// Update subscription to extend the collection period
	params := &stripelib.SubscriptionParams{
		CollectionMethod: stripelib.String("charge_automatically"),
	}

	// Add metadata to track grace period
	params.Metadata = map[string]string{
		"workspace_id":      ws.ID,
		"grace_period":      "true",
		"grace_period_end":  fmt.Sprintf("%d", gracePeriodEnd),
		"grace_period_days": fmt.Sprintf("%d", graceDays),
	}

	params.SetIdempotencyKey(idempotencyKey)

	_, err = subscription.Update(sub.StripeSubID, params)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to set grace period",
			"error", err.Error(),
			"subscription_id", sub.StripeSubID,
			"grace_days", graceDays,
			"idempotency_key", idempotencyKey,
		)
		return ErrStripeOperationFailed(err, "set_grace_period")
	}

	slog.InfoContext(ctx, "Grace period set successfully",
		"subscription_id", sub.StripeSubID,
		"grace_period_end", gracePeriodEnd,
		"grace_days", graceDays,
	)

	return nil
}

// CheckTrialStatus checks if a subscription is in trial and returns trial information
func (s *Service) CheckTrialStatus(ctx context.Context, ws *account.Workspace) (*TrialInfo, error) {
	slog.InfoContext(ctx, "Checking trial status",
		"workspace_id", ws.ID,
	)

	// Get current subscription
	sub, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	// Get Stripe subscription for detailed trial info
	stripeSub, err := subscription.Get(sub.StripeSubID, nil)
	if err != nil {
		return nil, ErrStripeOperationFailed(err, "get_subscription")
	}

	trialInfo := &TrialInfo{
		IsInTrial:     stripeSub.Status == "trialing",
		TrialEnd:      time.Unix(stripeSub.TrialEnd, 0),
		DaysRemaining: 0,
	}

	if trialInfo.IsInTrial {
		daysRemaining := time.Until(trialInfo.TrialEnd).Hours() / 24
		trialInfo.DaysRemaining = int(daysRemaining)
	}

	slog.InfoContext(ctx, "Trial status checked",
		"subscription_id", sub.StripeSubID,
		"is_in_trial", trialInfo.IsInTrial,
		"days_remaining", trialInfo.DaysRemaining,
	)

	return trialInfo, nil
}

// TrialInfo contains information about subscription trial status
type TrialInfo struct {
	IsInTrial     bool      `json:"isInTrial"`
	TrialEnd      time.Time `json:"trialEnd"`
	DaysRemaining int       `json:"daysRemaining"`
}

// CreateInvoice creates a new invoice for the workspace
func (s *Service) CreateInvoice(ctx context.Context, ws *account.Workspace, description string, amount int64, currency string, dueDate *string) (*stripelib.Invoice, error) {
	logger := slog.With("workspace_id", ws.ID, "amount", amount, "currency", currency)
	logger.Info("creating invoice")

	customerID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		logger.Error("failed to ensure customer", "error", err)
		return nil, ErrCustomerCreationFailed(ws.ID, err)
	}

	params := &stripelib.InvoiceParams{
		Customer:    stripelib.String(customerID),
		Currency:    stripelib.String(currency),
		Description: stripelib.String(description),
		AutoAdvance: stripelib.Bool(true),
	}

	if dueDate != nil {
		// Parse due date and set it
		if dueTime, err := time.Parse("2006-01-02", *dueDate); err == nil {
			params.DueDate = stripelib.Int64(dueTime.Unix())
		}
	}

	// Create invoice item first
	invItemParams := &stripelib.InvoiceItemParams{
		Customer:    stripelib.String(customerID),
		Amount:      stripelib.Int64(amount),
		Currency:    stripelib.String(currency),
		Description: stripelib.String(description),
	}

	_, err = invoiceitem.New(invItemParams)
	if err != nil {
		logger.Error("failed to create invoice item", "error", err)
		return nil, ErrStripeOperationFailed(err, "create_invoice_item")
	}

	// Create the invoice
	inv, err := invoice.New(params)
	if err != nil {
		logger.Error("failed to create invoice", "error", err)
		return nil, ErrStripeOperationFailed(err, "create_invoice")
	}

	logger.Info("invoice created successfully", "invoice_id", inv.ID)
	return inv, nil
}

// ScheduleSubscriptionChange schedules a subscription change for a future date
func (s *Service) ScheduleSubscriptionChange(ctx context.Context, ws *account.Workspace, plan *Plan, effectiveDate, prorationMode string) (*stripelib.SubscriptionSchedule, error) {
	logger := slog.With("workspace_id", ws.ID, "plan_descriptor", plan.Descriptor, "effective_date", effectiveDate)
	logger.Info("scheduling subscription change")

	// Get current subscription
	sub, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		logger.Error("failed to get current subscription", "error", err)
		return nil, err
	}

	// Parse effective date
	effectiveTime, err := time.Parse("2006-01-02T15:04:05Z", effectiveDate)
	if err != nil {
		// Try alternative format
		if effectiveTime, err = time.Parse("2006-01-02", effectiveDate); err != nil {
			logger.Error("invalid effective date format", "error", err)
			return nil, ErrStripeOperationFailed(err, "parse_date")
		}
	}

	// Ensure the plan exists in Stripe - for now we assume it exists
	// TODO: Add plan sync validation if needed

	// Create subscription schedule
	scheduleParams := &stripelib.SubscriptionScheduleParams{
		FromSubscription: stripelib.String(sub.StripeSubID),
	}

	// Add phases - use existing plan price in Stripe
	currentPhase := &stripelib.SubscriptionSchedulePhaseParams{
		Items: []*stripelib.SubscriptionSchedulePhaseItemParams{
			{
				Price: stripelib.String(plan.StripePlanID), // Use existing Stripe price ID
			},
		},
		StartDate: stripelib.Int64(effectiveTime.Unix()),
	}

	if prorationMode != "" {
		currentPhase.ProrationBehavior = stripelib.String(prorationMode)
	}

	scheduleParams.Phases = []*stripelib.SubscriptionSchedulePhaseParams{currentPhase}

	schedule, err := subscriptionschedule.New(scheduleParams)
	if err != nil {
		logger.Error("failed to create subscription schedule", "error", err)
		return nil, ErrStripeOperationFailed(err, "create_subscription_schedule")
	}

	logger.Info("subscription schedule created successfully", "schedule_id", schedule.ID)
	return schedule, nil
}

// Additional enhanced service methods for complete billing functionality

// CancelSubscription cancels a subscription immediately (alias for existing method)
func (s *Service) CancelSubscription(ctx context.Context, ws *account.Workspace) error {
	return s.CancelSubscriptionImmediately(ctx, ws)
}

// GetSubscriptionUsage retrieves usage data for metered billing
func (s *Service) GetSubscriptionUsage(ctx context.Context, ws *account.Workspace) (map[string]int64, error) {
	logger := slog.With("workspace_id", ws.ID)
	logger.Info("retrieving subscription usage")

	// Get current subscription
	sub, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		logger.Error("failed to get subscription", "error", err)
		return nil, err
	}

	// Mock usage data - in a real implementation, this would query usage records
	usage := map[string]int64{
		"api_calls":    1250,
		"storage_gb":   15,
		"users":        5,
		"projects":     3,
		"integrations": 2,
	}

	logger.Info("usage retrieved successfully", "subscription_id", sub.ID)
	return usage, nil
}

// ValidateSubscriptionAccess checks if workspace has access to specific features
func (s *Service) ValidateSubscriptionAccess(ctx context.Context, ws *account.Workspace, feature string) error {
	logger := slog.With("workspace_id", ws.ID, "feature", feature)
	logger.Info("validating subscription access")

	// Get current subscription and plan
	sub, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		logger.Error("failed to get subscription", "error", err)
		return err
	}

	if err := sub.IsActive(); err != nil {
		logger.Error("subscription not active", "error", err)
		return err
	}

	// Load plan details
	plan, err := s.GetPlanByID(ctx, sub.PlanID)
	if err != nil {
		logger.Error("failed to get plan", "error", err)
		return err
	}

	// Check feature availability based on plan
	// This is a simplified check - in a real implementation,
	// you'd have a more sophisticated feature matrix
	switch feature {
	case "advanced_analytics":
		if !plan.Features.AdvancedAnalytics {
			return ErrFeatureNotAvailable(nil, PlanSchema.AdvancedAnalytics)
		}
	case "advanced_reports":
		if !plan.Features.AdvancedFinancialReports {
			return ErrFeatureNotAvailable(nil, PlanSchema.AdvancedFinancialReports)
		}
	case "data_export":
		if !plan.Features.DataExport {
			return ErrFeatureNotAvailable(nil, PlanSchema.DataExport)
		}
	case "ai_assistant":
		if !plan.Features.AIBusinessAssistant {
			return ErrFeatureNotAvailable(nil, PlanSchema.AIBusinessAssistant)
		}
	}

	logger.Info("subscription access validated successfully")
	return nil
}

// EstimateProrationAmount estimates the proration amount for plan changes
func (s *Service) EstimateProrationAmount(ctx context.Context, ws *account.Workspace, newPlanDescriptor string) (int64, error) {
	logger := slog.With("workspace_id", ws.ID, "new_plan", newPlanDescriptor)
	logger.Info("estimating proration amount")

	// Get current subscription
	currentSub, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		logger.Error("failed to get current subscription", "error", err)
		return 0, err
	}

	// Get current and new plans
	currentPlan, err := s.GetPlanByID(ctx, currentSub.PlanID)
	if err != nil {
		logger.Error("failed to get current plan", "error", err)
		return 0, err
	}

	newPlan, err := s.GetPlanByDescriptor(ctx, newPlanDescriptor)
	if err != nil {
		logger.Error("failed to get new plan", "error", err)
		return 0, err
	}

	// Calculate proration (simplified calculation)
	// In a real implementation, you'd use Stripe's preview invoice API
	currentPricePerMonth := currentPlan.Price.IntPart()
	newPricePerMonth := newPlan.Price.IntPart()

	// Estimate based on remaining days in current period
	daysRemaining := int64(time.Until(currentSub.CurrentPeriodEnd).Hours() / 24)
	daysInMonth := int64(30) // Simplified

	prorationAmount := ((newPricePerMonth - currentPricePerMonth) * daysRemaining) / daysInMonth

	logger.Info("proration amount estimated", "amount", prorationAmount)
	return prorationAmount, nil
}

// ProcessWebhook processes incoming Stripe webhooks
func (s *Service) ProcessWebhook(ctx context.Context, payload []byte, signature string) error {
	logger := slog.With("event_type", "webhook")
	logger.Info("processing stripe webhook")

	// Note: In a real implementation, you'd need the webhook endpoint secret
	// For now, we'll skip signature verification in development
	// webhookSecret := viper.GetString("stripe.webhook_secret")

	// Parse the webhook event
	var event map[string]interface{}
	if err := json.Unmarshal(payload, &event); err != nil {
		logger.Error("failed to parse webhook payload", "error", err)
		return ErrWebhookProcessingFailed(err, "parse_payload")
	}

	eventType, ok := event["type"].(string)
	if !ok {
		logger.Error("webhook event missing type field")
		return ErrWebhookProcessingFailed(nil, "missing_event_type")
	}

	logger = logger.With("event_type", eventType)
	logger.Info("processing webhook event")

	// Handle different webhook events
	switch eventType {
	case "customer.subscription.created":
		return s.handleSubscriptionCreated(ctx, event)
	case "customer.subscription.updated":
		return s.handleSubscriptionUpdated(ctx, event)
	case "customer.subscription.deleted":
		return s.handleSubscriptionDeleted(ctx, event)
	case "invoice.payment_succeeded":
		return s.handleInvoicePaymentSucceeded(ctx, event)
	case "invoice.payment_failed":
		return s.handleInvoicePaymentFailed(ctx, event)
	case "customer.subscription.trial_will_end":
		return s.handleTrialWillEnd(ctx, event)
	default:
		logger.Info("unhandled webhook event type", "event_type", eventType)
		return nil // Not an error, just unhandled
	}
}

func (s *Service) handleSubscriptionCreated(ctx context.Context, event map[string]interface{}) error {
	logger := slog.With("event", "subscription_created")
	logger.Info("handling subscription created webhook")

	// Extract subscription data
	data, ok := event["data"].(map[string]interface{})
	if !ok {
		return ErrWebhookProcessingFailed(nil, "invalid_data_format")
	}

	object, ok := data["object"].(map[string]interface{})
	if !ok {
		return ErrWebhookProcessingFailed(nil, "invalid_object_format")
	}

	subscriptionID, ok := object["id"].(string)
	if !ok {
		return ErrWebhookProcessingFailed(nil, "missing_subscription_id")
	}

	status, ok := object["status"].(string)
	if !ok {
		return ErrWebhookProcessingFailed(nil, "missing_subscription_status")
	}

	// Extract period information
	currentPeriodStart := int64(0)
	currentPeriodEnd := int64(0)
	if period, ok := object["current_period_start"].(float64); ok {
		currentPeriodStart = int64(period)
	}
	if period, ok := object["current_period_end"].(float64); ok {
		currentPeriodEnd = int64(period)
	}

	s.SyncSubscriptionStatus(ctx, subscriptionID, status, currentPeriodStart, currentPeriodEnd)
	logger.Info("subscription created webhook processed successfully")
	return nil
}

func (s *Service) handleSubscriptionUpdated(ctx context.Context, event map[string]interface{}) error {
	logger := slog.With("event", "subscription_updated")
	logger.Info("handling subscription updated webhook")

	// Similar to created, but may need to handle plan changes
	data, ok := event["data"].(map[string]interface{})
	if !ok {
		return ErrWebhookProcessingFailed(nil, "invalid_data_format")
	}

	object, ok := data["object"].(map[string]interface{})
	if !ok {
		return ErrWebhookProcessingFailed(nil, "invalid_object_format")
	}

	subscriptionID, ok := object["id"].(string)
	if !ok {
		return ErrWebhookProcessingFailed(nil, "missing_subscription_id")
	}

	status, ok := object["status"].(string)
	if !ok {
		return ErrWebhookProcessingFailed(nil, "missing_subscription_status")
	}

	// Extract period information
	currentPeriodStart := int64(0)
	currentPeriodEnd := int64(0)
	if period, ok := object["current_period_start"].(float64); ok {
		currentPeriodStart = int64(period)
	}
	if period, ok := object["current_period_end"].(float64); ok {
		currentPeriodEnd = int64(period)
	}

	s.SyncSubscriptionStatus(ctx, subscriptionID, status, currentPeriodStart, currentPeriodEnd)
	logger.Info("subscription updated webhook processed successfully")
	return nil
}

func (s *Service) handleSubscriptionDeleted(ctx context.Context, event map[string]interface{}) error {
	logger := slog.With("event", "subscription_deleted")
	logger.Info("handling subscription deleted webhook")

	data, ok := event["data"].(map[string]interface{})
	if !ok {
		return ErrWebhookProcessingFailed(nil, "invalid_data_format")
	}

	object, ok := data["object"].(map[string]interface{})
	if !ok {
		return ErrWebhookProcessingFailed(nil, "invalid_object_format")
	}

	subscriptionID, ok := object["id"].(string)
	if !ok {
		return ErrWebhookProcessingFailed(nil, "missing_subscription_id")
	}

	// Extract period information for final processing
	currentPeriodStart := int64(0)
	currentPeriodEnd := int64(0)
	if period, ok := object["current_period_start"].(float64); ok {
		currentPeriodStart = int64(period)
	}
	if period, ok := object["current_period_end"].(float64); ok {
		currentPeriodEnd = int64(period)
	}

	s.RefundAndFinalizeCancellation(ctx, subscriptionID, currentPeriodStart, currentPeriodEnd)
	logger.Info("subscription deleted webhook processed successfully")
	return nil
}

func (s *Service) handleInvoicePaymentSucceeded(ctx context.Context, event map[string]interface{}) error {
	logger := slog.With("event", "invoice_payment_succeeded")
	logger.Info("handling invoice payment succeeded webhook")

	data, ok := event["data"].(map[string]interface{})
	if !ok {
		return ErrWebhookProcessingFailed(nil, "invalid_data_format")
	}

	object, ok := data["object"].(map[string]interface{})
	if !ok {
		return ErrWebhookProcessingFailed(nil, "invalid_object_format")
	}

	// Extract subscription ID if this is a subscription invoice
	if subscriptionID, ok := object["subscription"].(string); ok && subscriptionID != "" {
		s.MarkSubscriptionActive(ctx, subscriptionID)
		logger.Info("subscription marked as active after successful payment")
	}

	logger.Info("invoice payment succeeded webhook processed successfully")
	return nil
}

func (s *Service) handleInvoicePaymentFailed(ctx context.Context, event map[string]interface{}) error {
	logger := slog.With("event", "invoice_payment_failed")
	logger.Info("handling invoice payment failed webhook")

	data, ok := event["data"].(map[string]interface{})
	if !ok {
		return ErrWebhookProcessingFailed(nil, "invalid_data_format")
	}

	object, ok := data["object"].(map[string]interface{})
	if !ok {
		return ErrWebhookProcessingFailed(nil, "invalid_object_format")
	}

	// Extract subscription ID if this is a subscription invoice
	if subscriptionID, ok := object["subscription"].(string); ok && subscriptionID != "" {
		s.MarkSubscriptionPastDue(ctx, subscriptionID)
		logger.Info("subscription marked as past due after failed payment")
	}

	logger.Info("invoice payment failed webhook processed successfully")
	return nil
}

func (s *Service) handleTrialWillEnd(ctx context.Context, event map[string]interface{}) error {
	logger := slog.With("event", "trial_will_end")
	logger.Info("handling trial will end webhook")

	// This is a good place to send notification emails or trigger
	// other business logic when a trial is about to end

	data, ok := event["data"].(map[string]interface{})
	if !ok {
		return ErrWebhookProcessingFailed(nil, "invalid_data_format")
	}

	object, ok := data["object"].(map[string]interface{})
	if !ok {
		return ErrWebhookProcessingFailed(nil, "invalid_object_format")
	}

	subscriptionID, ok := object["id"].(string)
	if !ok {
		return ErrWebhookProcessingFailed(nil, "missing_subscription_id")
	}

	// Here you could:
	// 1. Send trial ending notification email
	// 2. Update subscription metadata
	// 3. Trigger business logic for trial conversion

	logger.Info("trial will end webhook processed successfully", "subscription_id", subscriptionID)
	return nil
}

// Enhanced middleware support methods

// CanUseFeature checks if a workspace's subscription allows a specific feature
// This method is designed to work with the enforce_plan_feature middleware
func (s *Service) CanUseFeature(ctx context.Context, workspaceID string, feature string) error {
	logger := slog.With("workspace_id", workspaceID, "feature", feature)
	logger.Info("checking feature availability")

	// Get subscription for workspace
	sub, err := s.GetSubscriptionByWorkspaceID(ctx, workspaceID)
	if err != nil {
		logger.Error("failed to get subscription", "error", err)
		return err
	}

	// Check if subscription is active
	if err := sub.IsActive(); err != nil {
		logger.Error("subscription not active", "error", err)
		return err
	}

	// Get plan details
	plan, err := s.GetPlanByID(ctx, sub.PlanID)
	if err != nil {
		logger.Error("failed to get plan", "error", err)
		return err
	}

	// Check feature availability
	switch feature {
	case "customer_management":
		if !plan.Features.CustomerManagement {
			return ErrFeatureNotAvailable(nil, PlanSchema.CustomerManagement)
		}
	case "inventory_management":
		if !plan.Features.InventoryManagement {
			return ErrFeatureNotAvailable(nil, PlanSchema.InventoryManagement)
		}
	case "order_management":
		if !plan.Features.OrderManagement {
			return ErrFeatureNotAvailable(nil, PlanSchema.OrderManagement)
		}
	case "expense_management":
		if !plan.Features.ExpenseManagement {
			return ErrFeatureNotAvailable(nil, PlanSchema.ExpenseManagement)
		}
	case "assets_management":
		if !plan.Features.AssetsManagement {
			return ErrFeatureNotAvailable(nil, PlanSchema.AssetsManagement)
		}
	case "accounting":
		if !plan.Features.Accounting {
			return ErrFeatureNotAvailable(nil, PlanSchema.Accounting)
		}
	case "basic_analytics":
		if !plan.Features.BasicAnalytics {
			return ErrFeatureNotAvailable(nil, PlanSchema.BasicAnalytics)
		}
	case "financial_reports":
		if !plan.Features.FinancialReports {
			return ErrFeatureNotAvailable(nil, PlanSchema.FinancialReports)
		}
	case "data_import":
		if !plan.Features.DataImport {
			return ErrFeatureNotAvailable(nil, PlanSchema.DataImport)
		}
	case "data_export":
		if !plan.Features.DataExport {
			return ErrFeatureNotAvailable(nil, PlanSchema.DataExport)
		}
	case "advanced_analytics":
		if !plan.Features.AdvancedAnalytics {
			return ErrFeatureNotAvailable(nil, PlanSchema.AdvancedAnalytics)
		}
	case "advanced_financial_reports":
		if !plan.Features.AdvancedFinancialReports {
			return ErrFeatureNotAvailable(nil, PlanSchema.AdvancedFinancialReports)
		}
	case "order_payment_links":
		if !plan.Features.OrderPaymentLinks {
			return ErrFeatureNotAvailable(nil, PlanSchema.OrderPaymentLinks)
		}
	case "invoice_generation":
		if !plan.Features.InvoiceGeneration {
			return ErrFeatureNotAvailable(nil, PlanSchema.InvoiceGeneration)
		}
	case "export_analytics_data":
		if !plan.Features.ExportAnalyticsData {
			return ErrFeatureNotAvailable(nil, PlanSchema.ExportAnalyticsData)
		}
	case "ai_business_assistant":
		if !plan.Features.AIBusinessAssistant {
			return ErrFeatureNotAvailable(nil, PlanSchema.AIBusinessAssistant)
		}
	default:
		logger.Warn("unknown feature requested", "feature", feature)
		return ErrUnknownFeature(nil, feature)
	}

	logger.Info("feature access granted")
	return nil
}

// CheckUsageLimitWithCallback checks if additional usage would exceed plan limits
// This method is designed to work with the enforce_plan_limit middleware
func (s *Service) CheckUsageLimitWithCallback(ctx context.Context, workspaceID, limitType string, additionalUsage int64, checkFunc func(used, limit int64) error) error {
	logger := slog.With("workspace_id", workspaceID, "limit_type", limitType, "additional_usage", additionalUsage)
	logger.Info("checking usage limit with callback")

	// Get current usage and limits
	current, limit, err := s.GetUsageQuota(ctx, &account.Workspace{ID: workspaceID}, limitType)
	if err != nil {
		logger.Error("failed to get usage quota", "error", err)
		return err
	}

	// If unlimited (-1), allow usage
	if limit == -1 {
		logger.Info("unlimited usage allowed")
		return nil
	}

	// Use custom check function if provided
	if checkFunc != nil {
		if err := checkFunc(current+additionalUsage, limit); err != nil {
			logger.Error("usage limit check failed", "error", err)
			return err
		}
	} else {
		// Default check
		if current+additionalUsage > limit {
			logger.Error("usage limit would be exceeded", "current", current, "additional", additionalUsage, "limit", limit)
			return ErrUsageLimitExceeded(nil, limitType, current+additionalUsage, limit)
		}
	}

	logger.Info("usage limit check passed")
	return nil
}

// GetSubscriptionByWorkspaceIDSafe safely retrieves subscription info for middleware
func (s *Service) GetSubscriptionByWorkspaceIDSafe(ctx context.Context, workspaceID string) (*Subscription, error) {
	logger := slog.With("workspace_id", workspaceID)
	logger.Info("getting subscription for middleware check")

	sub, err := s.GetSubscriptionByWorkspaceID(ctx, workspaceID)
	if err != nil {
		logger.Error("failed to get subscription", "error", err)
		return nil, ErrSubscriptionNotFound(err, workspaceID)
	}

	return sub, nil
}

// Enhanced middleware integration methods

// ValidateActiveSubscription validates that a workspace has an active subscription
// This method integrates with the enforce_active_sub middleware
func (s *Service) ValidateActiveSubscription(ctx context.Context, workspaceID string) (*Subscription, error) {
	sub, err := s.GetSubscriptionByWorkspaceIDSafe(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	if err := sub.IsActive(); err != nil {
		return nil, err
	}

	return sub, nil
}

// GetWorkspaceSubscriptionInfo returns comprehensive subscription information for a workspace
func (s *Service) GetWorkspaceSubscriptionInfo(ctx context.Context, workspaceID string) (*WorkspaceSubscriptionInfo, error) {
	logger := slog.With("workspace_id", workspaceID)
	logger.Info("getting workspace subscription info")

	sub, err := s.GetSubscriptionByWorkspaceID(ctx, workspaceID)
	if err != nil {
		return &WorkspaceSubscriptionInfo{
			HasSubscription: false,
			IsActive:        false,
		}, nil // Return default info instead of error for non-existent subscriptions
	}

	plan, err := s.GetPlanByID(ctx, sub.PlanID)
	if err != nil {
		logger.Error("failed to get plan details", "error", err)
		return nil, err
	}

	// Calculate trial information if applicable
	trialInfo, _ := s.CheckTrialStatus(ctx, &account.Workspace{ID: workspaceID})

	info := &WorkspaceSubscriptionInfo{
		HasSubscription:  true,
		IsActive:         sub.Status == SubscriptionStatusActive,
		SubscriptionID:   sub.ID,
		PlanDescriptor:   plan.Descriptor,
		PlanName:         plan.Name,
		Status:           string(sub.Status),
		CurrentPeriodEnd: sub.CurrentPeriodEnd,
		IsInTrial:        trialInfo != nil && trialInfo.IsInTrial,
		TrialEndsAt:      sub.CurrentPeriodEnd,
		Features:         plan.Features,
		Limits:           plan.Limits,
	}

	logger.Info("workspace subscription info retrieved successfully")
	return info, nil
}

// WorkspaceSubscriptionInfo contains comprehensive subscription information
type WorkspaceSubscriptionInfo struct {
	HasSubscription  bool        `json:"hasSubscription"`
	IsActive         bool        `json:"isActive"`
	SubscriptionID   string      `json:"subscriptionId,omitempty"`
	PlanDescriptor   string      `json:"planDescriptor,omitempty"`
	PlanName         string      `json:"planName,omitempty"`
	Status           string      `json:"status,omitempty"`
	CurrentPeriodEnd time.Time   `json:"currentPeriodEnd,omitempty"`
	IsInTrial        bool        `json:"isInTrial"`
	TrialEndsAt      time.Time   `json:"trialEndsAt,omitempty"`
	Features         PlanFeature `json:"features,omitempty"`
	Limits           PlanLimit   `json:"limits,omitempty"`
}
