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

// Notification encapsulates email sending for billing domain
type Notification struct {
	client     email.Client
	info       email.EmailInfo
	accountSvc *account.Service
}

// NewNotification wires the email client, defaults and account service
func NewNotification(client email.Client, info email.EmailInfo, accountSvc *account.Service) *Notification {
	return &Notification{client: client, info: info, accountSvc: accountSvc}
}

// SendSubscriptionWelcomeEmail sends a welcome email after successful subscription creation
func (n *Notification) SendSubscriptionWelcomeEmail(ctx context.Context, workspaceID string, subscription *Subscription, plan *Plan, paymentMethodLastFour string) error {
	if n.client == nil {
		return fmt.Errorf("email client not available")
	}
	logger := slog.With("action", "send_subscription_welcome", "workspace_id", workspaceID, "subscription_id", subscription.ID)
	logger.InfoContext(ctx, "Sending subscription welcome email")

	// Resolve primary user for workspace
	ws, err := n.accountSvc.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get workspace", "error", err)
		return fmt.Errorf("failed to get workspace: %w", err)
	}
	user, err := n.accountSvc.GetUserByID(ctx, ws.OwnerID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get workspace owner", "error", err)
		return fmt.Errorf("failed to get workspace owner: %w", err)
	}

	billingCycle := "month"
	if plan.BillingCycle == BillingCycleYearly {
		billingCycle = "year"
	}
	nextBillingDate := subscription.CurrentPeriodEnd.Format("January 2, 2006")

	data := map[string]any{
		"userName":        n.getUserDisplayName(user),
		"planName":        plan.Name,
		"amount":          n.formatAmount(plan.Price.String(), plan.Currency),
		"billingCycle":    billingCycle,
		"nextBillingDate": nextBillingDate,
		"lastFour":        paymentMethodLastFour,
		"dashboardURL":    fmt.Sprintf("%s/dashboard", n.info.BaseURL),
		"billingURL":      fmt.Sprintf("%s/billing", n.info.BaseURL),
		"productName":     n.info.ProductName,
		"supportEmail":    n.info.SupportEmail,
		"helpURL":         n.info.HelpURL,
		"currentYear":     fmt.Sprintf("%d", time.Now().Year()),
	}
	from := n.info.FormattedFrom()
	if _, err := n.client.SendTemplate(ctx, email.TemplateSubscriptionWelcome, []string{user.Email}, from, "", data); err != nil {
		logger.ErrorContext(ctx, "Failed to send subscription welcome email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	logger.InfoContext(ctx, "Subscription welcome email sent successfully")
	return nil
}

// SendPaymentFailedEmail sends an email notification when payment fails
func (n *Notification) SendPaymentFailedEmail(ctx context.Context, workspaceID string, subscription *Subscription, plan *Plan, paymentMethodLastFour string, attemptDate time.Time, nextAttemptDate *time.Time) error {
	if n.client == nil {
		return fmt.Errorf("email client not available")
	}
	logger := slog.With("action", "send_payment_failed", "workspace_id", workspaceID, "subscription_id", subscription.ID)
	logger.InfoContext(ctx, "Sending payment failed email")

	ws, err := n.accountSvc.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get workspace", "error", err)
		return fmt.Errorf("failed to get workspace: %w", err)
	}
	user, err := n.accountSvc.GetUserByID(ctx, ws.OwnerID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get workspace owner", "error", err)
		return fmt.Errorf("failed to get workspace owner: %w", err)
	}

	nextAttemptStr := "We will try again in a few days"
	if nextAttemptDate != nil {
		nextAttemptStr = nextAttemptDate.Format("January 2, 2006")
	}

	data := map[string]any{
		"userName":         n.getUserDisplayName(user),
		"amount":           n.formatAmount(plan.Price.String(), plan.Currency),
		"planName":         plan.Name,
		"lastFour":         paymentMethodLastFour,
		"attemptDate":      attemptDate.Format("January 2, 2006"),
		"nextAttemptDate":  nextAttemptStr,
		"gracePeriod":      "7 days",
		"updatePaymentURL": fmt.Sprintf("%s/billing/payment-methods", n.info.BaseURL),
		"retryPaymentURL":  fmt.Sprintf("%s/billing/retry-payment", n.info.BaseURL),
		"productName":      n.info.ProductName,
		"supportEmail":     n.info.SupportEmail,
		"helpURL":          n.info.HelpURL,
		"currentYear":      fmt.Sprintf("%d", time.Now().Year()),
	}
	from := n.info.FormattedFrom()
	if _, err := n.client.SendTemplate(ctx, email.TemplatePaymentFailed, []string{user.Email}, from, "", data); err != nil {
		logger.ErrorContext(ctx, "Failed to send payment failed email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	logger.InfoContext(ctx, "Payment failed email sent successfully")
	return nil
}

// SendSubscriptionCanceledEmail sends an email notification when subscription is canceled
func (n *Notification) SendSubscriptionCanceledEmail(ctx context.Context, workspaceID string, subscription *Subscription, plan *Plan, cancelDate time.Time, refundAmount string) error {
	if n.client == nil {
		return fmt.Errorf("email client not available")
	}
	logger := slog.With("action", "send_subscription_canceled", "workspace_id", workspaceID, "subscription_id", subscription.ID)
	logger.InfoContext(ctx, "Sending subscription canceled email")

	ws, err := n.accountSvc.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get workspace", "error", err)
		return fmt.Errorf("failed to get workspace: %w", err)
	}
	user, err := n.accountSvc.GetUserByID(ctx, ws.OwnerID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get workspace owner", "error", err)
		return fmt.Errorf("failed to get workspace owner: %w", err)
	}

	accessUntilDate := subscription.CurrentPeriodEnd.Format("January 2, 2006")
	data := map[string]any{
		"userName":        n.getUserDisplayName(user),
		"planName":        plan.Name,
		"cancelDate":      cancelDate.Format("January 2, 2006"),
		"accessUntilDate": accessUntilDate,
		"refundAmount":    refundAmount,
		"reactivateURL":   fmt.Sprintf("%s/billing/reactivate", n.info.BaseURL),
		"dashboardURL":    fmt.Sprintf("%s/dashboard", n.info.BaseURL),
		"feedbackURL":     fmt.Sprintf("%s/feedback?reason=cancellation", n.info.BaseURL),
		"productName":     n.info.ProductName,
		"supportEmail":    n.info.SupportEmail,
		"helpURL":         n.info.HelpURL,
		"currentYear":     fmt.Sprintf("%d", time.Now().Year()),
	}
	from := n.info.FormattedFrom()
	if _, err := n.client.SendTemplate(ctx, email.TemplateSubscriptionCanceled, []string{user.Email}, from, "", data); err != nil {
		logger.ErrorContext(ctx, "Failed to send subscription canceled email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	logger.InfoContext(ctx, "Subscription canceled email sent successfully")
	return nil
}

// SendTrialEndingEmail sends an email notification when trial is about to end
func (n *Notification) SendTrialEndingEmail(ctx context.Context, workspaceID string, subscription *Subscription, plan *Plan, trialInfo *TrialInfo) error {
	if n.client == nil {
		return fmt.Errorf("email client not available")
	}
	logger := slog.With("action", "send_trial_ending", "workspace_id", workspaceID, "subscription_id", subscription.ID)
	logger.InfoContext(ctx, "Sending trial ending email")

	ws, err := n.accountSvc.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get workspace", "error", err)
		return fmt.Errorf("failed to get workspace: %w", err)
	}
	user, err := n.accountSvc.GetUserByID(ctx, ws.OwnerID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get workspace owner", "error", err)
		return fmt.Errorf("failed to get workspace owner: %w", err)
	}

	trialPeriodDays := 14
	featuresUsed := []string{"Dashboard", "Basic Analytics", "Customer Management"}
	projectsCreated := 3

	data := map[string]any{
		"userName":        n.getUserDisplayName(user),
		"planName":        plan.Name,
		"trialPeriod":     fmt.Sprintf("%d days", trialPeriodDays),
		"daysRemaining":   trialInfo.DaysRemaining,
		"trialEndDate":    trialInfo.TrialEnd.Format("January 2, 2006"),
		"trialStartDate":  trialInfo.TrialEnd.AddDate(0, 0, -trialPeriodDays).Format("January 2, 2006"),
		"featuresUsed":    strings.Join(featuresUsed, ", "),
		"projectsCreated": fmt.Sprintf("%d", projectsCreated),
		"monthlyPrice":    n.formatAmount(plan.Price.String(), plan.Currency),
		"subscribeURL":    fmt.Sprintf("%s/billing/subscribe?plan=%s", n.info.BaseURL, plan.Descriptor),
		"plansURL":        fmt.Sprintf("%s/billing/plans", n.info.BaseURL),
		"productName":     n.info.ProductName,
		"supportEmail":    n.info.SupportEmail,
		"helpURL":         n.info.HelpURL,
		"currentYear":     fmt.Sprintf("%d", time.Now().Year()),
	}
	from := n.info.FormattedFrom()
	if _, err := n.client.SendTemplate(ctx, email.TemplateTrialEnding, []string{user.Email}, from, "", data); err != nil {
		logger.ErrorContext(ctx, "Failed to send trial ending email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	logger.InfoContext(ctx, "Trial ending email sent successfully")
	return nil
}

// SendSubscriptionConfirmedEmail sends email when first-time subscription is confirmed and paid
func (n *Notification) SendSubscriptionConfirmedEmail(ctx context.Context, workspaceID string, subscription *Subscription, plan *Plan, paymentMethodLastFour, invoiceURL string) error {
	if n.client == nil {
		return fmt.Errorf("email client not available")
	}
	logger := slog.With("action", "send_subscription_confirmed", "workspace_id", workspaceID, "subscription_id", subscription.ID)
	logger.InfoContext(ctx, "Sending subscription confirmed email")

	ws, err := n.accountSvc.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get workspace", "error", err)
		return fmt.Errorf("failed to get workspace: %w", err)
	}
	user, err := n.accountSvc.GetUserByID(ctx, ws.OwnerID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get workspace owner", "error", err)
		return fmt.Errorf("failed to get workspace owner: %w", err)
	}

	nextBillingDate := subscription.CurrentPeriodEnd.Format("January 2, 2006")
	invoiceNumber := fmt.Sprintf("INV-%d", time.Now().Unix())
	data := map[string]any{
		"userName":        n.getUserDisplayName(user),
		"planName":        plan.Name,
		"amount":          n.formatAmount(plan.Price.String(), plan.Currency),
		"paymentDate":     time.Now().Format("January 2, 2006"),
		"invoiceNumber":   invoiceNumber,
		"invoiceURL":      invoiceURL,
		"nextBillingDate": nextBillingDate,
		"dashboardURL":    fmt.Sprintf("%s/dashboard", n.info.BaseURL),
		"billingURL":      fmt.Sprintf("%s/billing", n.info.BaseURL),
		"supportEmail":    n.info.SupportEmail,
		"helpURL":         n.info.HelpURL,
		"productName":     n.info.ProductName,
		"currentYear":     fmt.Sprintf("%d", time.Now().Year()),
	}
	from := n.info.FormattedFrom()
	if _, err := n.client.SendTemplate(ctx, email.TemplateSubscriptionConfirmed, []string{user.Email}, from, "", data); err != nil {
		logger.ErrorContext(ctx, "Failed to send subscription confirmed email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	logger.InfoContext(ctx, "Subscription confirmed email sent successfully")
	return nil
}

// SendPaymentSucceededEmail sends an email notification for renewal payments
func (n *Notification) SendPaymentSucceededEmail(ctx context.Context, workspaceID string, subscription *Subscription, plan *Plan, paymentMethodLastFour, invoiceURL string, paymentDate time.Time) error {
	if n.client == nil {
		return fmt.Errorf("email client not available")
	}
	logger := slog.With("action", "send_payment_succeeded", "workspace_id", workspaceID, "subscription_id", subscription.ID)
	logger.InfoContext(ctx, "Sending payment succeeded email")

	ws, err := n.accountSvc.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get workspace", "error", err)
		return fmt.Errorf("failed to get workspace: %w", err)
	}
	user, err := n.accountSvc.GetUserByID(ctx, ws.OwnerID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get workspace owner", "error", err)
		return fmt.Errorf("failed to get workspace owner: %w", err)
	}

	nextBillingDate := subscription.CurrentPeriodEnd.Format("January 2, 2006")
	invoiceNumber := fmt.Sprintf("INV-%d", paymentDate.Unix())
	data := map[string]any{
		"userName":        n.getUserDisplayName(user),
		"planName":        plan.Name,
		"amount":          n.formatAmount(plan.Price.String(), plan.Currency),
		"paymentDate":     paymentDate.Format("January 2, 2006"),
		"invoiceNumber":   invoiceNumber,
		"nextBillingDate": nextBillingDate,
		"invoiceURL":      invoiceURL,
		"dashboardURL":    fmt.Sprintf("%s/dashboard", n.info.BaseURL),
		"billingURL":      fmt.Sprintf("%s/billing", n.info.BaseURL),
		"productName":     n.info.ProductName,
		"supportEmail":    n.info.SupportEmail,
		"helpURL":         n.info.HelpURL,
		"currentYear":     fmt.Sprintf("%d", time.Now().Year()),
	}
	from := n.info.FormattedFrom()
	if _, err := n.client.SendTemplate(ctx, email.TemplatePaymentSucceeded, []string{user.Email}, from, "", data); err != nil {
		logger.ErrorContext(ctx, "Failed to send payment succeeded email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	logger.InfoContext(ctx, "Payment succeeded email sent successfully")
	return nil
}

// Helpers
func (n *Notification) getUserDisplayName(user *account.User) string {
	if user.FirstName != "" {
		if user.LastName != "" {
			return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
		}
		return user.FirstName
	}
	if user.LastName != "" {
		return user.LastName
	}
	// Fallback
	parts := strings.Split(user.Email, "@")
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}
	return "User"
}

func (n *Notification) formatAmount(amount, currency string) string {
	switch strings.ToLower(currency) {
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
