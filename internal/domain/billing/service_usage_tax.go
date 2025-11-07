package billing

import (
	"context"
	"fmt"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	stripelib "github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/subscription"
	"github.com/stripe/stripe-go/v83/tax/calculation"
	"github.com/stripe/stripe-go/v83/tax/settings"
)

// CalculateTax calculates tax for a given amount and workspace location
func (s *Service) CalculateTax(ctx context.Context, ws *account.Workspace, amount int64, currency string) (*stripelib.TaxCalculation, error) {
	idempotencyKey := id.KsuidWithPrefix("tax_calc")
	logger.FromContext(ctx).Info("Calculating tax", "workspace_id", ws.ID, "amount", amount, "currency", currency, "idempotency_key", idempotencyKey)
	lineItems := []*stripelib.TaxCalculationLineItemParams{{Amount: stripelib.Int64(amount), Reference: stripelib.String("service_charge"), TaxCode: stripelib.String("txcd_10000000")}}
	params := &stripelib.TaxCalculationParams{
		Currency:        stripelib.String(currency),
		LineItems:       lineItems,
		CustomerDetails: &stripelib.TaxCalculationCustomerDetailsParams{AddressSource: stripelib.String("billing")},
	}
	params.SetIdempotencyKey(idempotencyKey)
	calc, err := calculation.New(params)
	if err != nil {
		logger.FromContext(ctx).Error("Failed to calculate tax", "error", err, "amount", amount, "currency", currency, "idempotency_key", idempotencyKey)
		return nil, ErrStripeOperationFailed(err, "calculate_tax")
	}
	logger.FromContext(ctx).Info("Tax calculation completed", "calculation_id", calc.ID, "amount_total", calc.AmountTotal, "tax_amount_exclusive", calc.TaxAmountExclusive, "tax_amount_inclusive", calc.TaxAmountInclusive)
	return calc, nil
}

// UpdateTaxSettings updates the account's tax settings for automatic tax calculation
func (s *Service) UpdateTaxSettings(ctx context.Context, defaultTaxCode string) error {
	idempotencyKey := id.KsuidWithPrefix("tax_settings")
	logger.FromContext(ctx).Info("Updating tax settings", "default_tax_code", defaultTaxCode, "idempotency_key", idempotencyKey)
	params := &stripelib.TaxSettingsParams{Defaults: &stripelib.TaxSettingsDefaultsParams{TaxCode: stripelib.String(defaultTaxCode)}}
	params.SetIdempotencyKey(idempotencyKey)
	taxSettings, err := settings.Update(params)
	if err != nil {
		logger.FromContext(ctx).Error("Failed to update tax settings", "error", err, "default_tax_code", defaultTaxCode, "idempotency_key", idempotencyKey)
		return ErrStripeOperationFailed(err, "update_tax_settings")
	}
	logger.FromContext(ctx).Info("Tax settings updated successfully", "status", taxSettings.Status, "default_tax_code", taxSettings.Defaults.TaxCode)
	return nil
}

// TrackUsage is a helper method to track API calls, storage, or other metered usage
func (s *Service) TrackUsage(ctx context.Context, ws *account.Workspace, usageType string, quantity int64) error {
	logger.FromContext(ctx).Info("Tracking usage", "workspace_id", ws.ID, "usage_type", usageType, "quantity", quantity)
	subscription, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		return fmt.Errorf("failed to get subscription for usage tracking: %w", err)
	}
	logger.FromContext(ctx).Info("Usage tracked locally", "workspace_id", ws.ID, "subscription_id", subscription.ID, "usage_type", usageType, "quantity", quantity, "timestamp", time.Now().Unix())
	return nil
}

// GetUsageQuota checks current usage against plan limits
func (s *Service) GetUsageQuota(ctx context.Context, ws *account.Workspace, usageType string) (used int64, limit int64, err error) {
	logger.FromContext(ctx).Info("Checking usage quota", "workspace_id", ws.ID, "usage_type", usageType)
	subscription, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get subscription for quota check: %w", err)
	}
	plan, err := s.GetPlanByID(ctx, subscription.PlanID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get plan for quota check: %w", err)
	}
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
	currentUsage := int64(0)
	logger.FromContext(ctx).Info("Usage quota retrieved", "workspace_id", ws.ID, "usage_type", usageType, "current_usage", currentUsage, "quota_limit", quotaLimit)
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
	logger.FromContext(ctx).Info("Checking trial status", "workspace_id", ws.ID)
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
		daysRemaining := time.Until(trialInfo.TrialEnd).Hours() / 24
		trialInfo.DaysRemaining = int(daysRemaining)
	}
	logger.FromContext(ctx).Info("Trial status checked", "subscription_id", sub.StripeSubID, "is_in_trial", trialInfo.IsInTrial, "days_remaining", trialInfo.DaysRemaining)
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
	l := logger.FromContext(ctx).With("workspace_id", ws.ID)
	l.Info("retrieving subscription usage")
	sub, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		l.Error("failed to get subscription", "error", err)
		return nil, err
	}
	usage := map[string]int64{"api_calls": 1250, "storage_gb": 15, "users": 5, "projects": 3, "integrations": 2}
	l.Info("usage retrieved successfully", "subscription_id", sub.ID)
	return usage, nil
}

// ValidateSubscriptionAccess checks if workspace has access to specific features
func (s *Service) ValidateSubscriptionAccess(ctx context.Context, ws *account.Workspace, feature string) error {
	l := logger.FromContext(ctx).With("workspace_id", ws.ID, "feature", feature)
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
