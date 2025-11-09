package billing

import (
	"context"
	"fmt"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/bus"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/email"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
	"github.com/abdelrahman146/kyora/internal/platform/types/atomic"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	stripelib "github.com/stripe/stripe-go/v83"
	checkoutsession "github.com/stripe/stripe-go/v83/checkout/session"
)

type Service struct {
	storage         *Storage
	atomicProcessor atomic.AtomicProcessor
	bus             *bus.Bus
	account         *account.Service
	Notification    *Notification
}

func NewService(storage *Storage, atomicProcessor atomic.AtomicProcessor, bus *bus.Bus, accountSvc *account.Service, emailClient email.Client) *Service {
	// Note: Stripe is used via package-level helpers using API key configured globally if needed in future.
	emailInfo := email.NewEmail()
	notification := NewNotification(emailClient, emailInfo, accountSvc)
	s := &Service{
		storage:         storage,
		atomicProcessor: atomicProcessor,
		bus:             bus,
		account:         accountSvc,
		Notification:    notification,
	}
	// Best-effort background plan sync on service creation
	go func() {
		// use background context to avoid blocking init
		if err := s.SyncPlansToStripe(context.Background()); err != nil {
			logger.FromContext(context.Background()).Error("plan sync to stripe failed", "error", err)
		}
	}()
	return s
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
		// Free plan: create or update subscription directly without collecting payment method.
		// Returns empty URL indicating no redirect needed.
		_, subErr := s.CreateOrUpdateSubscription(ctx, ws, plan)
		if subErr != nil {
			return "", subErr
		}
		logger.FromContext(ctx).Info("free plan subscription created without checkout session", "workspaceId", ws.ID, "planId", plan.ID)
		return "", nil
	} else {
		// Ensure plan exists in Stripe
		if err := s.ensurePlanSynced(ctx, plan); err != nil {
			logger.FromContext(ctx).Error("failed to ensure plan synced before checkout", "error", err, "planId", plan.ID)
			return "", fmt.Errorf("failed to ensure plan in stripe: %w", err)
		}
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
		logger.FromContext(ctx).Error("failed to create checkout session", "error", err, "workspaceId", ws.ID, "planId", plan.ID)
		return "", fmt.Errorf("failed to create checkout session: %w", err)
	}

	logger.FromContext(ctx).Info("created checkout session", "workspaceId", ws.ID, "planId", plan.ID, "sessionId", session.ID)
	return session.URL, nil
}

// CanUseFeature checks if a workspace's subscription allows a specific feature
// This method is designed to work with the enforce_plan_feature middleware
func (s *Service) CanUseFeature(ctx context.Context, workspaceID string, feature role.Resource) error {
	logger := logger.FromContext(ctx).With("workspaceId", workspaceID, "feature", feature)
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
	case role.ResourceCustomer:
		if !plan.Features.CustomerManagement {
			return ErrFeatureNotAvailable(nil, PlanSchema.CustomerManagement)
		}
	case role.ResourceInventory:
		if !plan.Features.InventoryManagement {
			return ErrFeatureNotAvailable(nil, PlanSchema.InventoryManagement)
		}
	case role.ResourceOrder:
		if !plan.Features.OrderManagement {
			return ErrFeatureNotAvailable(nil, PlanSchema.OrderManagement)
		}
	case role.ResourceExpense:
		if !plan.Features.ExpenseManagement {
			return ErrFeatureNotAvailable(nil, PlanSchema.ExpenseManagement)
		}
	case role.ResourceAccounting:
		if !plan.Features.Accounting {
			return ErrFeatureNotAvailable(nil, PlanSchema.Accounting)
		}
	case role.ResourceBasicAnalytics:
		if !plan.Features.BasicAnalytics {
			return ErrFeatureNotAvailable(nil, PlanSchema.BasicAnalytics)
		}
	case role.ResourceFinancialReports:
		if !plan.Features.FinancialReports {
			return ErrFeatureNotAvailable(nil, PlanSchema.FinancialReports)
		}
	case role.ResourceDataImport:
		if !plan.Features.DataImport {
			return ErrFeatureNotAvailable(nil, PlanSchema.DataImport)
		}
	case role.ResourceDataExport:
		if !plan.Features.DataExport {
			return ErrFeatureNotAvailable(nil, PlanSchema.DataExport)
		}
	case role.ResourceAdvancedAnalytics:
		if !plan.Features.AdvancedAnalytics {
			return ErrFeatureNotAvailable(nil, PlanSchema.AdvancedAnalytics)
		}
	case role.ResourceAdvancedFinancialReports:
		if !plan.Features.AdvancedFinancialReports {
			return ErrFeatureNotAvailable(nil, PlanSchema.AdvancedFinancialReports)
		}
	case role.ResourceOrderPaymentLinks:
		if !plan.Features.OrderPaymentLinks {
			return ErrFeatureNotAvailable(nil, PlanSchema.OrderPaymentLinks)
		}
	case role.ResourceOrderInvoiceGeneration:
		if !plan.Features.InvoiceGeneration {
			return ErrFeatureNotAvailable(nil, PlanSchema.InvoiceGeneration)
		}
	case role.ResourceExportAnalyticsData:
		if !plan.Features.ExportAnalyticsData {
			return ErrFeatureNotAvailable(nil, PlanSchema.ExportAnalyticsData)
		}
	case role.ResourceAIBusinessAssistant:
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
	logger := logger.FromContext(ctx).With("workspaceId", workspaceID, "limitType", limitType, "additionalUsage", additionalUsage)
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
	logger := logger.FromContext(ctx).With("workspaceId", workspaceID)
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
	logger := logger.FromContext(ctx).With("workspaceId", workspaceID)
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
