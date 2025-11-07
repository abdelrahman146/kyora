package billing

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/email"
)

// EmailIntegration provides methods to integrate billing operations with email notifications
type EmailIntegration struct {
	notificationService *email.NotificationService
	accountService      *account.Service
	baseURL             string
}

// NewEmailIntegration creates a new billing email integration helper
func NewEmailIntegration(notificationService *email.NotificationService, accountService *account.Service, baseURL string) *EmailIntegration {
	if baseURL == "" {
		baseURL = "https://app.kyora.com"
	}
	return &EmailIntegration{
		notificationService: notificationService,
		accountService:      accountService,
		baseURL:             baseURL,
	}
}

// SendSubscriptionWelcomeEmail sends a welcome email after successful subscription creation
func (e *EmailIntegration) SendSubscriptionWelcomeEmail(ctx context.Context, workspaceID string, subscription *Subscription, plan *Plan, paymentMethodLastFour string) error {
	if e.notificationService == nil {
		return fmt.Errorf("notification service not available")
	}

	logger := slog.With("action", "send_subscription_welcome", "workspace_id", workspaceID, "subscription_id", subscription.ID)
	logger.InfoContext(ctx, "Sending subscription welcome email")

	// Get workspace to find the primary user
	workspace, err := e.accountService.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get workspace", "error", err)
		return fmt.Errorf("failed to get workspace: %w", err)
	}

	// Get the owner/primary user by ID
	user, err := e.accountService.GetUserByID(ctx, workspace.OwnerID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get workspace owner", "error", err)
		return fmt.Errorf("failed to get workspace owner: %w", err)
	}

	// Calculate next billing date
	nextBillingDate := subscription.CurrentPeriodEnd.Format("January 2, 2006")

	// Format amount based on billing cycle
	billingCycle := "monthly"
	if plan.BillingCycle == BillingCycleYearly {
		billingCycle = "yearly"
	}

	params := email.SubscriptionWelcomeParams{
		Email:           user.Email,
		UserName:        e.getUserDisplayName(user),
		PlanName:        plan.Name,
		Amount:          e.formatAmount(plan.Price.String(), plan.Currency),
		BillingCycle:    billingCycle,
		NextBillingDate: nextBillingDate,
		LastFour:        paymentMethodLastFour,
		DashboardURL:    fmt.Sprintf("%s/dashboard", e.baseURL),
		BillingURL:      fmt.Sprintf("%s/billing", e.baseURL),
	}

	err = e.notificationService.SendSubscriptionWelcomeEmail(ctx, params)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to send subscription welcome email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.InfoContext(ctx, "Subscription welcome email sent successfully")
	return nil
}

// SendPaymentFailedEmail sends an email notification when payment fails
func (e *EmailIntegration) SendPaymentFailedEmail(ctx context.Context, workspaceID string, subscription *Subscription, plan *Plan, paymentMethodLastFour string, attemptDate time.Time, nextAttemptDate *time.Time) error {
	if e.notificationService == nil {
		return fmt.Errorf("notification service not available")
	}

	logger := slog.With("action", "send_payment_failed", "workspace_id", workspaceID, "subscription_id", subscription.ID)
	logger.InfoContext(ctx, "Sending payment failed email")

	// Get workspace to find the primary user
	workspace, err := e.accountService.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get workspace", "error", err)
		return fmt.Errorf("failed to get workspace: %w", err)
	}

	// Get the owner/primary user by ID
	user, err := e.accountService.GetUserByID(ctx, workspace.OwnerID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get workspace owner", "error", err)
		return fmt.Errorf("failed to get workspace owner: %w", err)
	}

	// Format next attempt date
	nextAttemptStr := "We will try again in a few days"
	if nextAttemptDate != nil {
		nextAttemptStr = nextAttemptDate.Format("January 2, 2006")
	}

	params := email.PaymentFailedParams{
		Email:            user.Email,
		UserName:         e.getUserDisplayName(user),
		Amount:           e.formatAmount(plan.Price.String(), plan.Currency),
		PlanName:         plan.Name,
		LastFour:         paymentMethodLastFour,
		AttemptDate:      attemptDate.Format("January 2, 2006"),
		NextAttemptDate:  nextAttemptStr,
		GracePeriod:      "7 days", // Default grace period
		UpdatePaymentURL: fmt.Sprintf("%s/billing/payment-methods", e.baseURL),
		RetryPaymentURL:  fmt.Sprintf("%s/billing/retry-payment", e.baseURL),
	}

	err = e.notificationService.SendPaymentFailedEmail(ctx, params)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to send payment failed email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.InfoContext(ctx, "Payment failed email sent successfully")
	return nil
}

// SendSubscriptionCanceledEmail sends an email notification when subscription is canceled
func (e *EmailIntegration) SendSubscriptionCanceledEmail(ctx context.Context, workspaceID string, subscription *Subscription, plan *Plan, cancelDate time.Time, refundAmount string) error {
	if e.notificationService == nil {
		return fmt.Errorf("notification service not available")
	}

	logger := slog.With("action", "send_subscription_canceled", "workspace_id", workspaceID, "subscription_id", subscription.ID)
	logger.InfoContext(ctx, "Sending subscription canceled email")

	// Get workspace to find the primary user
	workspace, err := e.accountService.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get workspace", "error", err)
		return fmt.Errorf("failed to get workspace: %w", err)
	}

	// Get the owner/primary user by ID
	user, err := e.accountService.GetUserByID(ctx, workspace.OwnerID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get workspace owner", "error", err)
		return fmt.Errorf("failed to get workspace owner: %w", err)
	}

	// Calculate access until date (typically end of current billing period)
	accessUntilDate := subscription.CurrentPeriodEnd.Format("January 2, 2006")

	params := email.SubscriptionCanceledParams{
		Email:           user.Email,
		UserName:        e.getUserDisplayName(user),
		PlanName:        plan.Name,
		CancelDate:      cancelDate.Format("January 2, 2006"),
		AccessUntilDate: accessUntilDate,
		RefundAmount:    refundAmount,
		ReactivateURL:   fmt.Sprintf("%s/billing/reactivate", e.baseURL),
		DashboardURL:    fmt.Sprintf("%s/dashboard", e.baseURL),
		FeedbackURL:     fmt.Sprintf("%s/feedback?reason=cancellation", e.baseURL),
	}

	err = e.notificationService.SendSubscriptionCanceledEmail(ctx, params)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to send subscription canceled email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.InfoContext(ctx, "Subscription canceled email sent successfully")
	return nil
}

// SendTrialEndingEmail sends an email notification when trial is about to end
func (e *EmailIntegration) SendTrialEndingEmail(ctx context.Context, workspaceID string, subscription *Subscription, plan *Plan, trialInfo *TrialInfo) error {
	if e.notificationService == nil {
		return fmt.Errorf("notification service not available")
	}

	logger := slog.With("action", "send_trial_ending", "workspace_id", workspaceID, "subscription_id", subscription.ID)
	logger.InfoContext(ctx, "Sending trial ending email")

	// Get workspace to find the primary user
	workspace, err := e.accountService.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get workspace", "error", err)
		return fmt.Errorf("failed to get workspace: %w", err)
	}

	// Get the owner/primary user by ID
	user, err := e.accountService.GetUserByID(ctx, workspace.OwnerID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get workspace owner", "error", err)
		return fmt.Errorf("failed to get workspace owner: %w", err)
	}

	// Calculate trial period duration (approximation)
	trialPeriodDays := 14 // Default trial period, could be calculated from subscription data

	// Get usage statistics (mock data for now)
	featuresUsed := []string{"Dashboard", "Basic Analytics", "Customer Management"}
	projectsCreated := 3

	params := email.TrialEndingParams{
		Email:           user.Email,
		UserName:        e.getUserDisplayName(user),
		PlanName:        plan.Name,
		TrialPeriod:     fmt.Sprintf("%d days", trialPeriodDays),
		DaysRemaining:   trialInfo.DaysRemaining,
		TrialEndDate:    trialInfo.TrialEnd.Format("January 2, 2006"),
		TrialStartDate:  trialInfo.TrialEnd.AddDate(0, 0, -trialPeriodDays).Format("January 2, 2006"),
		FeaturesUsed:    strings.Join(featuresUsed, ", "),
		ProjectsCreated: fmt.Sprintf("%d", projectsCreated),
		MonthlyPrice:    e.formatAmount(plan.Price.String(), plan.Currency),
		SubscribeURL:    fmt.Sprintf("%s/billing/subscribe?plan=%s", e.baseURL, plan.Descriptor),
		PlansURL:        fmt.Sprintf("%s/billing/plans", e.baseURL),
	}

	err = e.notificationService.SendTrialEndingEmail(ctx, params)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to send trial ending email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.InfoContext(ctx, "Trial ending email sent successfully")
	return nil
}

// getUserDisplayName returns a user-friendly display name
func (e *EmailIntegration) getUserDisplayName(user *account.User) string {
	if user.FirstName != "" {
		if user.LastName != "" {
			return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
		}
		return user.FirstName
	}
	if user.LastName != "" {
		return user.LastName
	}
	// Fallback to email prefix if no names are available
	return "User"
}

// formatAmount formats a price amount with currency symbol
func (e *EmailIntegration) formatAmount(amount, currency string) string {
	switch currency {
	case "usd":
		return fmt.Sprintf("$%s", amount)
	case "eur":
		return fmt.Sprintf("€%s", amount)
	case "gbp":
		return fmt.Sprintf("£%s", amount)
	default:
		return fmt.Sprintf("%s %s", amount, currency)
	}
}

// Helper methods for different subscription events

// NotifySubscriptionCreated sends appropriate notifications when a subscription is created
func (e *EmailIntegration) NotifySubscriptionCreated(ctx context.Context, workspaceID string, subscription *Subscription, plan *Plan, paymentMethodLastFour string) error {
	// For paid plans, send welcome email
	if !plan.Price.IsZero() {
		return e.SendSubscriptionWelcomeEmail(ctx, workspaceID, subscription, plan, paymentMethodLastFour)
	}
	// For free plans, we might send a different type of welcome email
	// For now, still send the subscription welcome but with $0 amount
	return e.SendSubscriptionWelcomeEmail(ctx, workspaceID, subscription, plan, paymentMethodLastFour)
}

// NotifyPaymentFailed sends appropriate notifications when a payment fails
func (e *EmailIntegration) NotifyPaymentFailed(ctx context.Context, workspaceID string, subscription *Subscription, plan *Plan, paymentMethodLastFour string) error {
	// Calculate next attempt date (typically 3-7 days later)
	nextAttempt := time.Now().AddDate(0, 0, 3)

	return e.SendPaymentFailedEmail(ctx, workspaceID, subscription, plan, paymentMethodLastFour, time.Now(), &nextAttempt)
}

// NotifySubscriptionCanceled sends appropriate notifications when a subscription is canceled
func (e *EmailIntegration) NotifySubscriptionCanceled(ctx context.Context, workspaceID string, subscription *Subscription, plan *Plan, refundAmount string) error {
	return e.SendSubscriptionCanceledEmail(ctx, workspaceID, subscription, plan, time.Now(), refundAmount)
}

// NotifyTrialEnding sends appropriate notifications when a trial is about to end
func (e *EmailIntegration) NotifyTrialEnding(ctx context.Context, workspaceID string, subscription *Subscription, plan *Plan, trialInfo *TrialInfo) error {
	// Only send if trial is ending within the next 3 days
	if trialInfo.DaysRemaining <= 3 && trialInfo.DaysRemaining > 0 {
		return e.SendTrialEndingEmail(ctx, workspaceID, subscription, plan, trialInfo)
	}
	return nil
}

// SendSubscriptionConfirmedEmail sends an email notification when first-time subscription is confirmed and paid
func (e *EmailIntegration) SendSubscriptionConfirmedEmail(ctx context.Context, workspaceID string, subscription *Subscription, plan *Plan, paymentMethodLastFour, invoiceURL string) error {
	if e.notificationService == nil {
		return fmt.Errorf("notification service not available")
	}

	logger := slog.With("action", "send_subscription_confirmed", "workspace_id", workspaceID, "subscription_id", subscription.ID)
	logger.InfoContext(ctx, "Sending subscription confirmed email")

	// Get workspace to find the primary user
	workspace, err := e.accountService.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get workspace", "error", err)
		return fmt.Errorf("failed to get workspace: %w", err)
	}

	// Get the owner/primary user by ID
	user, err := e.accountService.GetUserByID(ctx, workspace.OwnerID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get workspace owner", "error", err)
		return fmt.Errorf("failed to get workspace owner: %w", err)
	}

	// Calculate next billing date
	nextBillingDate := subscription.CurrentPeriodEnd.Format("January 2, 2006")

	params := email.SubscriptionConfirmedParams{
		Email:           user.Email,
		UserName:        e.getUserDisplayName(user),
		PlanName:        plan.Name,
		Amount:          e.formatAmount(plan.Price.String(), plan.Currency),
		PaymentDate:     time.Now().Format("January 2, 2006"),
		InvoiceNumber:   fmt.Sprintf("INV-%d", time.Now().Unix()),
		InvoiceURL:      invoiceURL,
		NextBillingDate: nextBillingDate,
		DashboardURL:    fmt.Sprintf("%s/dashboard", e.baseURL),
		BillingURL:      fmt.Sprintf("%s/billing", e.baseURL),
		SupportEmail:    "support@kyora.com",
	}

	err = e.notificationService.SendSubscriptionConfirmedEmail(ctx, params)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to send subscription confirmed email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.InfoContext(ctx, "Subscription confirmed email sent successfully")
	return nil
}

// SendPaymentSucceededEmail sends an email notification for successful automatic renewal payments with invoice
func (e *EmailIntegration) SendPaymentSucceededEmail(ctx context.Context, workspaceID string, subscription *Subscription, plan *Plan, paymentMethodLastFour, invoiceURL string, paymentDate time.Time) error {
	if e.notificationService == nil {
		return fmt.Errorf("notification service not available")
	}

	logger := slog.With("action", "send_payment_succeeded", "workspace_id", workspaceID, "subscription_id", subscription.ID)
	logger.InfoContext(ctx, "Sending payment succeeded email")

	// Get workspace to find the primary user
	workspace, err := e.accountService.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get workspace", "error", err)
		return fmt.Errorf("failed to get workspace: %w", err)
	}

	// Get the owner/primary user by ID
	user, err := e.accountService.GetUserByID(ctx, workspace.OwnerID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get workspace owner", "error", err)
		return fmt.Errorf("failed to get workspace owner: %w", err)
	}

	// Calculate next billing date
	nextBillingDate := subscription.CurrentPeriodEnd.Format("January 2, 2006")

	params := email.PaymentSucceededParams{
		Email:           user.Email,
		UserName:        e.getUserDisplayName(user),
		PlanName:        plan.Name,
		Amount:          e.formatAmount(plan.Price.String(), plan.Currency),
		PaymentDate:     paymentDate.Format("January 2, 2006"),
		InvoiceNumber:   fmt.Sprintf("INV-%d", paymentDate.Unix()),
		NextBillingDate: nextBillingDate,
		InvoiceURL:      invoiceURL,
		DashboardURL:    fmt.Sprintf("%s/dashboard", e.baseURL),
	}

	err = e.notificationService.SendPaymentSucceededEmail(ctx, params)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to send payment succeeded email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.InfoContext(ctx, "Payment succeeded email sent successfully")
	return nil
}
