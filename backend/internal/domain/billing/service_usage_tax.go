package billing

import (
	"context"
	"fmt"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
	"github.com/abdelrahman146/kyora/internal/platform/utils/helpers"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	stripelib "github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/subscription"
	"github.com/stripe/stripe-go/v83/tax/calculation"
	"github.com/stripe/stripe-go/v83/tax/settings"
)

// CalculateTax calculates tax for a given amount and workspace location
func (s *Service) CalculateTax(ctx context.Context, ws *account.Workspace, amount int64, currency string) (*stripelib.TaxCalculation, error) {
	idempotencyKey := id.KsuidWithPrefix("tax_calc")
	logger.FromContext(ctx).Info("calculating tax", "workspaceId", ws.ID, "amount", amount, "currency", currency, "idempotencyKey", idempotencyKey)
	lineItems := []*stripelib.TaxCalculationLineItemParams{{Amount: stripelib.Int64(amount), Reference: stripelib.String("service_charge"), TaxCode: stripelib.String("txcd_10000000")}}
	params := &stripelib.TaxCalculationParams{
		Currency:        stripelib.String(currency),
		LineItems:       lineItems,
		CustomerDetails: &stripelib.TaxCalculationCustomerDetailsParams{AddressSource: stripelib.String("billing")},
	}
	params.SetIdempotencyKey(idempotencyKey)
	calc, err := withStripeRetry[*stripelib.TaxCalculation](ctx, 3, func() (*stripelib.TaxCalculation, error) {
		return calculation.New(params)
	})
	if err != nil {
		logger.FromContext(ctx).Error("failed to calculate tax", "error", err, "amount", amount, "currency", currency, "idempotencyKey", idempotencyKey)
		return nil, ErrStripeOperationFailed(err, "calculate_tax")
	}
	logger.FromContext(ctx).Info("tax calculation completed", "calculationId", calc.ID, "amountTotal", calc.AmountTotal, "taxAmountExclusive", calc.TaxAmountExclusive, "taxAmountInclusive", calc.TaxAmountInclusive)
	return calc, nil
}

// UpdateTaxSettings updates the account's tax settings for automatic tax calculation
func (s *Service) UpdateTaxSettings(ctx context.Context, defaultTaxCode string) error {
	idempotencyKey := id.KsuidWithPrefix("tax_settings")
	logger.FromContext(ctx).Info("updating tax settings", "defaultTaxCode", defaultTaxCode, "idempotencyKey", idempotencyKey)
	params := &stripelib.TaxSettingsParams{Defaults: &stripelib.TaxSettingsDefaultsParams{TaxCode: stripelib.String(defaultTaxCode)}}
	params.SetIdempotencyKey(idempotencyKey)
	taxSettings, err := withStripeRetry[*stripelib.TaxSettings](ctx, 3, func() (*stripelib.TaxSettings, error) {
		return settings.Update(params)
	})
	if err != nil {
		logger.FromContext(ctx).Error("failed to update tax settings", "error", err, "defaultTaxCode", defaultTaxCode, "idempotencyKey", idempotencyKey)
		return ErrStripeOperationFailed(err, "update_tax_settings")
	}
	logger.FromContext(ctx).Info("tax settings updated successfully", "status", taxSettings.Status, "defaultTaxCode", taxSettings.Defaults.TaxCode)
	return nil
}

// TrackUsage is a helper method to track API calls, storage, or other metered usage
func (s *Service) TrackUsage(ctx context.Context, ws *account.Workspace, usageType string, quantity int64) error {
	logger.FromContext(ctx).Info("tracking usage", "workspaceId", ws.ID, "usageType", usageType, "quantity", quantity)
	subscription, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		return fmt.Errorf("failed to get subscription for usage tracking: %w", err)
	}
	logger.FromContext(ctx).Info("usage tracked locally", "workspaceId", ws.ID, "subscriptionId", subscription.ID, "usageType", usageType, "quantity", quantity, "timestamp", time.Now().Unix())
	return nil
}

// GetUsageQuota checks current usage against plan limits
func (s *Service) GetUsageQuota(ctx context.Context, ws *account.Workspace, usageType string) (used int64, limit int64, err error) {
	logger.FromContext(ctx).Info("checking usage quota", "workspaceId", ws.ID, "usageType", usageType)
	subscription, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get subscription for quota check: %w", err)
	}
	plan, err := s.GetPlanByID(ctx, subscription.PlanID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get plan for quota check: %w", err)
	}
	var quotaLimit int64
	var currentUsage int64
	switch usageType {
	case "orders_per_month":
		quotaLimit = plan.Limits.MaxOrdersPerMonth
		now := time.Now().UTC()
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Nanosecond)
		currentUsage, err = s.storage.CountMonthlyOrdersByWorkspace(ctx, ws.ID, monthStart, monthEnd)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to count monthly orders: %w", err)
		}
	case "team_members":
		quotaLimit = plan.Limits.MaxTeamMembers
		currentUsage, err = s.account.CountWorkspaceUsers(ctx, ws.ID)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to count workspace users: %w", err)
		}
	case "businesses":
		quotaLimit = plan.Limits.MaxBusinesses
		currentUsage, err = s.storage.CountBusinessesByWorkspace(ctx, ws.ID)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to count workspace businesses: %w", err)
		}
	default:
		return 0, 0, fmt.Errorf("unsupported usage type: %s", usageType)
	}
	logger.FromContext(ctx).Info("usage quota retrieved", "workspaceId", ws.ID, "usageType", usageType, "currentUsage", currentUsage, "quotaLimit", quotaLimit)
	return currentUsage, quotaLimit, nil
}

// CheckUsageLimit verifies if usage is within plan limits
func (s *Service) CheckUsageLimit(ctx context.Context, ws *account.Workspace, usageType string, additionalUsage int64) error {
	current, limit, err := s.GetUsageQuota(ctx, ws, usageType)
	if err != nil {
		return fmt.Errorf("failed to check usage quota: %w", err)
	}
	if limit == -1 {
		return nil
	}
	if current+additionalUsage > limit {
		return fmt.Errorf("usage limit exceeded for %s: current %d + additional %d > limit %d", usageType, current, additionalUsage, limit)
	}
	return nil
}

// CreateTrialSubscription creates a subscription with a trial period
func (s *Service) CreateTrialSubscription(ctx context.Context, ws *account.Workspace, plan *Plan, trialDays int) (*stripelib.Subscription, error) {
	idempotencyKey := id.KsuidWithPrefix("trial_sub")
	logger.FromContext(ctx).Info("Creating trial subscription", "workspace_id", ws.ID, "plan_id", plan.ID, "trial_days", trialDays, "idempotency_key", idempotencyKey)
	customerID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure customer exists: %w", err)
	}
	if err := s.ensurePlanSynced(ctx, plan); err != nil {
		logger.FromContext(ctx).Error("failed to ensure plan synced before trial subscription", "error", err, "plan_id", plan.ID)
		return nil, fmt.Errorf("failed to ensure plan in stripe: %w", err)
	}
	trialEnd := time.Now().AddDate(0, 0, trialDays).Unix()
	params := &stripelib.SubscriptionParams{
		Customer:          stripelib.String(customerID),
		Items:             []*stripelib.SubscriptionItemsParams{{Price: stripelib.String(plan.StripePlanID)}},
		TrialEnd:          stripelib.Int64(trialEnd),
		CollectionMethod:  stripelib.String("charge_automatically"),
		PaymentBehavior:   stripelib.String("default_incomplete"),
		ProrationBehavior: stripelib.String("none"),
	}
	params.Expand = []*string{stripelib.String("latest_invoice.payment_intent")}
	params.Metadata = map[string]string{"workspace_id": ws.ID, "plan_id": plan.ID, "trial": "true", "trial_days": fmt.Sprintf("%d", trialDays)}
	params.SetIdempotencyKey(idempotencyKey)
	sub, err := subscription.New(params)
	if err != nil {
		logger.FromContext(ctx).Error("Failed to create trial subscription", "error", err, "workspace_id", ws.ID, "plan_id", plan.ID, "trial_days", trialDays, "idempotency_key", idempotencyKey)
		return nil, ErrStripeOperationFailed(err, "create_trial_subscription")
	}
	logger.FromContext(ctx).Info("Trial subscription created successfully", "subscription_id", sub.ID, "customer_id", customerID, "trial_end", trialEnd, "status", sub.Status)
	return sub, nil
}

// ExtendTrialPeriod extends the trial period for an existing subscription
func (s *Service) ExtendTrialPeriod(ctx context.Context, ws *account.Workspace, additionalDays int) error {
	idempotencyKey := id.KsuidWithPrefix("extend_trial")
	logger.FromContext(ctx).Info("Extending trial period", "workspace_id", ws.ID, "additional_days", additionalDays, "idempotency_key", idempotencyKey)
	sub, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}
	if sub.Status != SubscriptionStatusTrialing {
		return fmt.Errorf("subscription is not in trial period, current status: %s", sub.Status)
	}
	stripeSub, err := subscription.Get(sub.StripeSubID, nil)
	if err != nil {
		return ErrStripeOperationFailed(err, "get_subscription")
	}
	currentTrialEnd := time.Unix(stripeSub.TrialEnd, 0)
	newTrialEnd := currentTrialEnd.AddDate(0, 0, additionalDays).Unix()
	params := &stripelib.SubscriptionParams{TrialEnd: stripelib.Int64(newTrialEnd)}
	params.SetIdempotencyKey(idempotencyKey)
	_, err = subscription.Update(sub.StripeSubID, params)
	if err != nil {
		logger.FromContext(ctx).Error("Failed to extend trial period", "error", err, "subscription_id", sub.StripeSubID, "additional_days", additionalDays, "idempotency_key", idempotencyKey)
		return ErrStripeOperationFailed(err, "extend_trial")
	}
	logger.FromContext(ctx).Info("Trial period extended successfully", "subscription_id", sub.StripeSubID, "new_trial_end", newTrialEnd, "additional_days", additionalDays)
	return nil
}

// HandleGracePeriod manages grace periods for failed payments
func (s *Service) HandleGracePeriod(ctx context.Context, ws *account.Workspace, graceDays int) error {
	idempotencyKey := id.KsuidWithPrefix("grace_period")
	logger.FromContext(ctx).Info("Handling grace period", "workspace_id", ws.ID, "grace_days", graceDays, "idempotency_key", idempotencyKey)
	sub, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}
	if sub.Status != SubscriptionStatusPastDue {
		return fmt.Errorf("subscription is not past due, current status: %s", sub.Status)
	}
	gracePeriodEnd := time.Now().AddDate(0, 0, graceDays).Unix()
	params := &stripelib.SubscriptionParams{CollectionMethod: stripelib.String("charge_automatically")}
	params.Metadata = map[string]string{"workspace_id": ws.ID, "grace_period": "true", "grace_period_end": fmt.Sprintf("%d", gracePeriodEnd), "grace_period_days": fmt.Sprintf("%d", graceDays)}
	params.SetIdempotencyKey(idempotencyKey)
	_, err = subscription.Update(sub.StripeSubID, params)
	if err != nil {
		logger.FromContext(ctx).Error("Failed to set grace period", "error", err, "subscription_id", sub.StripeSubID, "grace_days", graceDays, "idempotency_key", idempotencyKey)
		return ErrStripeOperationFailed(err, "set_grace_period")
	}
	logger.FromContext(ctx).Info("Grace period set successfully", "subscription_id", sub.StripeSubID, "grace_period_end", gracePeriodEnd, "grace_days", graceDays)
	return nil
}

// CheckTrialStatus checks if a subscription is in trial and returns trial information
func (s *Service) CheckTrialStatus(ctx context.Context, ws *account.Workspace) (*TrialInfo, error) {
	logger.FromContext(ctx).Info("checking trial status", "workspaceId", ws.ID)
	sub, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}
	stripeSub, err := subscription.Get(sub.StripeSubID, nil)
	if err != nil {
		return nil, ErrStripeOperationFailed(err, "get_subscription")
	}
	trialInfo := &TrialInfo{IsInTrial: stripeSub.Status == "trialing", TrialEnd: time.Unix(stripeSub.TrialEnd, 0), DaysRemaining: 0}
	if trialInfo.IsInTrial {
		trialInfo.DaysRemaining = helpers.CeilPositiveDaysUntil(trialInfo.TrialEnd)
	}
	logger.FromContext(ctx).Info("trial status checked", "subscriptionId", sub.StripeSubID, "isInTrial", trialInfo.IsInTrial, "daysRemaining", trialInfo.DaysRemaining)
	return trialInfo, nil
}

// TrialInfo contains information about subscription trial status
type TrialInfo struct {
	IsInTrial     bool      `json:"isInTrial"`
	TrialEnd      time.Time `json:"trialEnd"`
	DaysRemaining int       `json:"daysRemaining"`
}

// GetSubscriptionUsage retrieves usage data for metered billing (placeholder)
func (s *Service) GetSubscriptionUsage(ctx context.Context, ws *account.Workspace) (map[string]int64, error) {
	l := logger.FromContext(ctx).With("workspaceId", ws.ID)
	l.Info("retrieving subscription usage")
	_, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		l.Error("failed to get subscription", "error", err)
		return nil, err
	}
	ordersUsed, _, err := s.GetUsageQuota(ctx, ws, "orders_per_month")
	if err != nil {
		return nil, err
	}
	usersUsed, _, err := s.GetUsageQuota(ctx, ws, "team_members")
	if err != nil {
		return nil, err
	}
	businessesUsed, _, err := s.GetUsageQuota(ctx, ws, "businesses")
	if err != nil {
		return nil, err
	}
	usage := map[string]int64{
		"ordersPerMonth": ordersUsed,
		"teamMembers":   usersUsed,
		"businesses":    businessesUsed,
	}
	l.Info("usage retrieved successfully")
	return usage, nil
}

// ValidateSubscriptionAccess checks if workspace has access to specific features
func (s *Service) ValidateSubscriptionAccess(ctx context.Context, ws *account.Workspace, feature string) error {
	l := logger.FromContext(ctx).With("workspaceId", ws.ID, "feature", feature)
	l.Info("validating subscription access")
	sub, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		l.Error("failed to get subscription", "error", err)
		return err
	}
	if err := sub.IsActive(); err != nil {
		l.Error("subscription not active", "error", err)
		return err
	}
	plan, err := s.GetPlanByID(ctx, sub.PlanID)
	if err != nil {
		l.Error("failed to get plan", "error", err)
		return err
	}
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
	l.Info("subscription access validated successfully")
	return nil
}
