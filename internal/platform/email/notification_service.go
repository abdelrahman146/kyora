package email

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// NotificationService provides high-level methods for sending transactional emails
// with proper data validation, error handling, and logging
type NotificationService struct {
	client Client
	config NotificationConfig
}

// NotificationConfig holds configuration for the notification service
type NotificationConfig struct {
	FromEmail    string
	FromName     string
	ProductName  string
	SupportEmail string
	HelpURL      string
	BaseURL      string
	NoReply      bool
}

// NewNotificationService creates a new notification service with the given client and config
func NewNotificationService(client Client, config NotificationConfig) *NotificationService {
	// Set defaults if not provided
	if config.FromEmail == "" {
		config.FromEmail = viper.GetString("email.from_email")
		if config.FromEmail == "" {
			config.FromEmail = "no-reply@kyora.com"
		}
	}
	if config.FromName == "" {
		config.FromName = viper.GetString("email.from_name")
		if config.FromName == "" {
			config.FromName = config.ProductName
		}
	}
	if config.ProductName == "" {
		config.ProductName = viper.GetString("app.name")
		if config.ProductName == "" {
			config.ProductName = "Kyora"
		}
	}
	if config.SupportEmail == "" {
		config.SupportEmail = viper.GetString("email.support_email")
		if config.SupportEmail == "" {
			config.SupportEmail = "support@kyora.com"
		}
	}
	if config.HelpURL == "" {
		config.HelpURL = viper.GetString("email.help_url")
		if config.HelpURL == "" {
			config.HelpURL = "https://help.kyora.com"
		}
	}
	if config.BaseURL == "" {
		config.BaseURL = viper.GetString("http.base_url")
		if config.BaseURL == "" {
			config.BaseURL = "https://app.kyora.com"
		}
	}

	return &NotificationService{
		client: client,
		config: config,
	}
}

// SendForgotPasswordEmail sends a password reset email to the user
func (s *NotificationService) SendForgotPasswordEmail(ctx context.Context, params ForgotPasswordParams) error {
	logger := slog.With("action", "send_forgot_password_email", "email", params.Email)
	logger.InfoContext(ctx, "Sending forgot password email")

	if err := s.validateForgotPasswordParams(params); err != nil {
		logger.ErrorContext(ctx, "Invalid forgot password parameters", "error", err)
		return fmt.Errorf("invalid parameters: %w", err)
	}

	data := map[string]any{
		"userName":     params.UserName,
		"resetURL":     params.ResetURL,
		"productName":  s.config.ProductName,
		"supportEmail": s.config.SupportEmail,
		"helpURL":      s.config.HelpURL,
		"expiryTime":   params.ExpiryTime,
		"currentYear":  strconv.Itoa(time.Now().Year()),
	}

	from := s.formatFromAddress()
	result, err := s.client.SendTemplate(ctx, TemplateForgotPassword, []string{params.Email}, from, "", data)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to send forgot password email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.InfoContext(ctx, "Forgot password email sent successfully", "provider", result.Provider, "id", result.ID)
	return nil
}

// SendPasswordResetConfirmationEmail sends a confirmation email after password reset
func (s *NotificationService) SendPasswordResetConfirmationEmail(ctx context.Context, params PasswordResetConfirmationParams) error {
	logger := slog.With("action", "send_password_reset_confirmation", "email", params.Email)
	logger.InfoContext(ctx, "Sending password reset confirmation email")

	if err := s.validatePasswordResetConfirmationParams(params); err != nil {
		logger.ErrorContext(ctx, "Invalid password reset confirmation parameters", "error", err)
		return fmt.Errorf("invalid parameters: %w", err)
	}

	data := map[string]any{
		"userName":      params.UserName,
		"loginURL":      params.LoginURL,
		"productName":   s.config.ProductName,
		"supportEmail":  s.config.SupportEmail,
		"helpURL":       s.config.HelpURL,
		"resetDate":     params.ResetDate,
		"resetTime":     params.ResetTime,
		"resetLocation": params.ResetLocation,
		"resetIP":       params.ResetIP,
		"currentYear":   strconv.Itoa(time.Now().Year()),
	}

	from := s.formatFromAddress()
	result, err := s.client.SendTemplate(ctx, TemplatePasswordResetConfirm, []string{params.Email}, from, "", data)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to send password reset confirmation email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.InfoContext(ctx, "Password reset confirmation email sent successfully", "provider", result.Provider, "id", result.ID)
	return nil
}

// SendEmailVerificationEmail sends an email verification email to the user
func (s *NotificationService) SendEmailVerificationEmail(ctx context.Context, params EmailVerificationParams) error {
	logger := slog.With("action", "send_email_verification", "email", params.Email)
	logger.InfoContext(ctx, "Sending email verification email")

	if err := s.validateEmailVerificationParams(params); err != nil {
		logger.ErrorContext(ctx, "Invalid email verification parameters", "error", err)
		return fmt.Errorf("invalid parameters: %w", err)
	}

	data := map[string]any{
		"userName":     params.UserName,
		"verifyURL":    params.VerifyURL,
		"productName":  s.config.ProductName,
		"supportEmail": s.config.SupportEmail,
		"helpURL":      s.config.HelpURL,
		"expiryTime":   params.ExpiryTime,
		"currentYear":  strconv.Itoa(time.Now().Year()),
	}

	from := s.formatFromAddress()
	result, err := s.client.SendTemplate(ctx, TemplateEmailVerification, []string{params.Email}, from, "", data)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to send email verification email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.InfoContext(ctx, "Email verification email sent successfully", "provider", result.Provider, "id", result.ID)
	return nil
}

// SendWelcomeEmail sends a welcome email to new users
func (s *NotificationService) SendWelcomeEmail(ctx context.Context, params WelcomeParams) error {
	logger := slog.With("action", "send_welcome_email", "email", params.Email)
	logger.InfoContext(ctx, "Sending welcome email")

	if err := s.validateWelcomeParams(params); err != nil {
		logger.ErrorContext(ctx, "Invalid welcome parameters", "error", err)
		return fmt.Errorf("invalid parameters: %w", err)
	}

	data := map[string]any{
		"userName":     params.UserName,
		"dashboardURL": params.DashboardURL,
		"guideURL":     params.GuideURL,
		"productName":  s.config.ProductName,
		"supportEmail": s.config.SupportEmail,
		"helpURL":      s.config.HelpURL,
		"currentYear":  strconv.Itoa(time.Now().Year()),
	}

	from := s.formatFromAddress()
	result, err := s.client.SendTemplate(ctx, TemplateWelcome, []string{params.Email}, from, "", data)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to send welcome email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.InfoContext(ctx, "Welcome email sent successfully", "provider", result.Provider, "id", result.ID)
	return nil
}

// SendSubscriptionWelcomeEmail sends a welcome email for new subscriptions
func (s *NotificationService) SendSubscriptionWelcomeEmail(ctx context.Context, params SubscriptionWelcomeParams) error {
	logger := slog.With("action", "send_subscription_welcome", "email", params.Email)
	logger.InfoContext(ctx, "Sending subscription welcome email")

	if err := s.validateSubscriptionWelcomeParams(params); err != nil {
		logger.ErrorContext(ctx, "Invalid subscription welcome parameters", "error", err)
		return fmt.Errorf("invalid parameters: %w", err)
	}

	data := map[string]any{
		"userName":        params.UserName,
		"planName":        params.PlanName,
		"amount":          params.Amount,
		"billingCycle":    params.BillingCycle,
		"nextBillingDate": params.NextBillingDate,
		"lastFour":        params.LastFour,
		"dashboardURL":    params.DashboardURL,
		"billingURL":      params.BillingURL,
		"productName":     s.config.ProductName,
		"supportEmail":    s.config.SupportEmail,
		"helpURL":         s.config.HelpURL,
		"currentYear":     strconv.Itoa(time.Now().Year()),
	}

	from := s.formatFromAddress()
	result, err := s.client.SendTemplate(ctx, TemplateSubscriptionWelcome, []string{params.Email}, from, "", data)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to send subscription welcome email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.InfoContext(ctx, "Subscription welcome email sent successfully", "provider", result.Provider, "id", result.ID)
	return nil
}

// SendPaymentFailedEmail sends an email notification for failed payments
func (s *NotificationService) SendPaymentFailedEmail(ctx context.Context, params PaymentFailedParams) error {
	logger := slog.With("action", "send_payment_failed", "email", params.Email)
	logger.InfoContext(ctx, "Sending payment failed email")

	if err := s.validatePaymentFailedParams(params); err != nil {
		logger.ErrorContext(ctx, "Invalid payment failed parameters", "error", err)
		return fmt.Errorf("invalid parameters: %w", err)
	}

	data := map[string]any{
		"userName":         params.UserName,
		"amount":           params.Amount,
		"planName":         params.PlanName,
		"lastFour":         params.LastFour,
		"attemptDate":      params.AttemptDate,
		"nextAttemptDate":  params.NextAttemptDate,
		"gracePeriod":      params.GracePeriod,
		"updatePaymentURL": params.UpdatePaymentURL,
		"retryPaymentURL":  params.RetryPaymentURL,
		"productName":      s.config.ProductName,
		"supportEmail":     s.config.SupportEmail,
		"helpURL":          s.config.HelpURL,
		"currentYear":      strconv.Itoa(time.Now().Year()),
	}

	from := s.formatFromAddress()
	result, err := s.client.SendTemplate(ctx, TemplatePaymentFailed, []string{params.Email}, from, "", data)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to send payment failed email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.InfoContext(ctx, "Payment failed email sent successfully", "provider", result.Provider, "id", result.ID)
	return nil
}

// SendSubscriptionCanceledEmail sends an email notification for canceled subscriptions
func (s *NotificationService) SendSubscriptionCanceledEmail(ctx context.Context, params SubscriptionCanceledParams) error {
	logger := slog.With("action", "send_subscription_canceled", "email", params.Email)
	logger.InfoContext(ctx, "Sending subscription canceled email")

	if err := s.validateSubscriptionCanceledParams(params); err != nil {
		logger.ErrorContext(ctx, "Invalid subscription canceled parameters", "error", err)
		return fmt.Errorf("invalid parameters: %w", err)
	}

	data := map[string]any{
		"userName":        params.UserName,
		"planName":        params.PlanName,
		"cancelDate":      params.CancelDate,
		"accessUntilDate": params.AccessUntilDate,
		"refundAmount":    params.RefundAmount,
		"reactivateURL":   params.ReactivateURL,
		"dashboardURL":    params.DashboardURL,
		"feedbackURL":     params.FeedbackURL,
		"productName":     s.config.ProductName,
		"supportEmail":    s.config.SupportEmail,
		"helpURL":         s.config.HelpURL,
		"currentYear":     strconv.Itoa(time.Now().Year()),
	}

	from := s.formatFromAddress()
	result, err := s.client.SendTemplate(ctx, TemplateSubscriptionCanceled, []string{params.Email}, from, "", data)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to send subscription canceled email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.InfoContext(ctx, "Subscription canceled email sent successfully", "provider", result.Provider, "id", result.ID)
	return nil
}

// SendTrialEndingEmail sends an email notification when trial is ending
func (s *NotificationService) SendTrialEndingEmail(ctx context.Context, params TrialEndingParams) error {
	logger := slog.With("action", "send_trial_ending", "email", params.Email)
	logger.InfoContext(ctx, "Sending trial ending email")

	if err := s.validateTrialEndingParams(params); err != nil {
		logger.ErrorContext(ctx, "Invalid trial ending parameters", "error", err)
		return fmt.Errorf("invalid parameters: %w", err)
	}

	data := map[string]any{
		"userName":        params.UserName,
		"planName":        params.PlanName,
		"trialPeriod":     params.TrialPeriod,
		"daysRemaining":   params.DaysRemaining,
		"trialEndDate":    params.TrialEndDate,
		"trialStartDate":  params.TrialStartDate,
		"featuresUsed":    params.FeaturesUsed,
		"projectsCreated": params.ProjectsCreated,
		"monthlyPrice":    params.MonthlyPrice,
		"subscribeURL":    params.SubscribeURL,
		"plansURL":        params.PlansURL,
		"productName":     s.config.ProductName,
		"supportEmail":    s.config.SupportEmail,
		"helpURL":         s.config.HelpURL,
		"currentYear":     strconv.Itoa(time.Now().Year()),
	}

	from := s.formatFromAddress()
	result, err := s.client.SendTemplate(ctx, TemplateTrialEnding, []string{params.Email}, from, "", data)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to send trial ending email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.InfoContext(ctx, "Trial ending email sent successfully", "provider", result.Provider, "id", result.ID)
	return nil
}

// SendLoginNotificationEmail sends a security notification email when user logs in
func (s *NotificationService) SendLoginNotificationEmail(ctx context.Context, params LoginNotificationParams) error {
	logger := slog.With("action", "send_login_notification", "email", params.Email)
	logger.InfoContext(ctx, "Sending login notification email")

	if err := s.validateLoginNotificationParams(params); err != nil {
		logger.ErrorContext(ctx, "Invalid login notification parameters", "error", err)
		return fmt.Errorf("invalid parameters: %w", err)
	}

	data := map[string]any{
		"userName":      params.UserName,
		"loginDate":     params.LoginDate,
		"loginTime":     params.LoginTime,
		"loginLocation": params.LoginLocation,
		"loginIP":       params.LoginIP,
		"deviceInfo":    params.DeviceInfo,
		"resetURL":      params.ResetURL,
		"productName":   s.config.ProductName,
		"supportEmail":  params.SupportEmail,
		"helpURL":       s.config.HelpURL,
		"currentYear":   strconv.Itoa(time.Now().Year()),
	}

	from := s.formatFromAddress()
	result, err := s.client.SendTemplate(ctx, TemplateLoginNotification, []string{params.Email}, from, "", data)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to send login notification email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.InfoContext(ctx, "Login notification email sent successfully", "provider", result.Provider, "id", result.ID)
	return nil
}

// SendSubscriptionConfirmedEmail sends a confirmation email after first payment is processed
func (s *NotificationService) SendSubscriptionConfirmedEmail(ctx context.Context, params SubscriptionConfirmedParams) error {
	logger := slog.With("action", "send_subscription_confirmed", "email", params.Email)
	logger.InfoContext(ctx, "Sending subscription confirmed email")

	if err := s.validateSubscriptionConfirmedParams(params); err != nil {
		logger.ErrorContext(ctx, "Invalid subscription confirmed parameters", "error", err)
		return fmt.Errorf("invalid parameters: %w", err)
	}

	data := map[string]any{
		"userName":        params.UserName,
		"planName":        params.PlanName,
		"amount":          params.Amount,
		"paymentDate":     params.PaymentDate,
		"invoiceNumber":   params.InvoiceNumber,
		"invoiceURL":      params.InvoiceURL,
		"nextBillingDate": params.NextBillingDate,
		"dashboardURL":    params.DashboardURL,
		"billingURL":      params.BillingURL,
		"productName":     s.config.ProductName,
		"supportEmail":    params.SupportEmail,
		"helpURL":         s.config.HelpURL,
		"currentYear":     strconv.Itoa(time.Now().Year()),
	}

	from := s.formatFromAddress()
	result, err := s.client.SendTemplate(ctx, TemplateSubscriptionConfirmed, []string{params.Email}, from, "", data)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to send subscription confirmed email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.InfoContext(ctx, "Subscription confirmed email sent successfully", "provider", result.Provider, "id", result.ID)
	return nil
}

// SendPaymentSucceededEmail sends a confirmation email for successful automatic renewals
func (s *NotificationService) SendPaymentSucceededEmail(ctx context.Context, params PaymentSucceededParams) error {
	logger := slog.With("action", "send_payment_succeeded", "email", params.Email)
	logger.InfoContext(ctx, "Sending payment succeeded email")

	if err := s.validatePaymentSucceededParams(params); err != nil {
		logger.ErrorContext(ctx, "Invalid payment succeeded parameters", "error", err)
		return fmt.Errorf("invalid parameters: %w", err)
	}

	data := map[string]any{
		"userName":        params.UserName,
		"amount":          params.Amount,
		"planName":        params.PlanName,
		"paymentDate":     params.PaymentDate,
		"invoiceNumber":   params.InvoiceNumber,
		"nextBillingDate": params.NextBillingDate,
		"invoiceURL":      params.InvoiceURL,
		"dashboardURL":    params.DashboardURL,
		"productName":     s.config.ProductName,
		"supportEmail":    s.config.SupportEmail,
		"helpURL":         s.config.HelpURL,
		"currentYear":     strconv.Itoa(time.Now().Year()),
	}

	from := s.formatFromAddress()
	result, err := s.client.SendTemplate(ctx, TemplatePaymentSucceeded, []string{params.Email}, from, "", data)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to send payment succeeded email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.InfoContext(ctx, "Payment succeeded email sent successfully", "provider", result.Provider, "id", result.ID)
	return nil
}

// formatFromAddress returns a properly formatted "From" address
func (s *NotificationService) formatFromAddress() string {
	if s.config.FromName != "" {
		return fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	}
	return s.config.FromEmail
}

// validateEmail checks if an email address is valid
func (s *NotificationService) validateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if !strings.Contains(email, "@") {
		return fmt.Errorf("email must contain @ symbol")
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return fmt.Errorf("email format is invalid")
	}
	if parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("email format is invalid")
	}
	// Basic domain validation
	if _, err := net.LookupMX(parts[1]); err != nil {
		// If MX lookup fails, try A record as fallback
		if _, err := net.LookupHost(parts[1]); err != nil {
			slog.Warn("Email domain validation failed", "email", email, "error", err)
			// Don't fail completely as this could be a temporary DNS issue
		}
	}
	return nil
}

// validateURL checks if a URL is valid
func (s *NotificationService) validateURL(url string) error {
	if url == "" {
		return fmt.Errorf("URL is required")
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("URL must start with http:// or https://")
	}
	return nil
}

// Validation methods for each email type
func (s *NotificationService) validateForgotPasswordParams(params ForgotPasswordParams) error {
	if err := s.validateEmail(params.Email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}
	if err := s.validateURL(params.ResetURL); err != nil {
		return fmt.Errorf("invalid reset URL: %w", err)
	}
	if params.UserName == "" {
		return fmt.Errorf("user name is required")
	}
	return nil
}

func (s *NotificationService) validatePasswordResetConfirmationParams(params PasswordResetConfirmationParams) error {
	if err := s.validateEmail(params.Email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}
	if err := s.validateURL(params.LoginURL); err != nil {
		return fmt.Errorf("invalid login URL: %w", err)
	}
	if params.UserName == "" {
		return fmt.Errorf("user name is required")
	}
	return nil
}

func (s *NotificationService) validateEmailVerificationParams(params EmailVerificationParams) error {
	if err := s.validateEmail(params.Email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}
	if err := s.validateURL(params.VerifyURL); err != nil {
		return fmt.Errorf("invalid verify URL: %w", err)
	}
	if params.UserName == "" {
		return fmt.Errorf("user name is required")
	}
	return nil
}

func (s *NotificationService) validateWelcomeParams(params WelcomeParams) error {
	if err := s.validateEmail(params.Email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}
	if err := s.validateURL(params.DashboardURL); err != nil {
		return fmt.Errorf("invalid dashboard URL: %w", err)
	}
	if params.UserName == "" {
		return fmt.Errorf("user name is required")
	}
	return nil
}

func (s *NotificationService) validateSubscriptionWelcomeParams(params SubscriptionWelcomeParams) error {
	if err := s.validateEmail(params.Email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}
	if err := s.validateURL(params.DashboardURL); err != nil {
		return fmt.Errorf("invalid dashboard URL: %w", err)
	}
	if err := s.validateURL(params.BillingURL); err != nil {
		return fmt.Errorf("invalid billing URL: %w", err)
	}
	if params.UserName == "" {
		return fmt.Errorf("user name is required")
	}
	if params.PlanName == "" {
		return fmt.Errorf("plan name is required")
	}
	return nil
}

func (s *NotificationService) validatePaymentFailedParams(params PaymentFailedParams) error {
	if err := s.validateEmail(params.Email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}
	if err := s.validateURL(params.UpdatePaymentURL); err != nil {
		return fmt.Errorf("invalid update payment URL: %w", err)
	}
	if params.UserName == "" {
		return fmt.Errorf("user name is required")
	}
	if params.Amount == "" {
		return fmt.Errorf("amount is required")
	}
	return nil
}

func (s *NotificationService) validateSubscriptionCanceledParams(params SubscriptionCanceledParams) error {
	if err := s.validateEmail(params.Email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}
	if err := s.validateURL(params.DashboardURL); err != nil {
		return fmt.Errorf("invalid dashboard URL: %w", err)
	}
	if params.UserName == "" {
		return fmt.Errorf("user name is required")
	}
	if params.PlanName == "" {
		return fmt.Errorf("plan name is required")
	}
	return nil
}

func (s *NotificationService) validateTrialEndingParams(params TrialEndingParams) error {
	if err := s.validateEmail(params.Email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}
	if err := s.validateURL(params.SubscribeURL); err != nil {
		return fmt.Errorf("invalid subscribe URL: %w", err)
	}
	if params.UserName == "" {
		return fmt.Errorf("user name is required")
	}
	if params.DaysRemaining < 0 {
		return fmt.Errorf("days remaining must be non-negative")
	}
	return nil
}

func (s *NotificationService) validateLoginNotificationParams(params LoginNotificationParams) error {
	if err := s.validateEmail(params.Email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}
	if err := s.validateURL(params.ResetURL); err != nil {
		return fmt.Errorf("invalid reset URL: %w", err)
	}
	if params.UserName == "" {
		return fmt.Errorf("user name is required")
	}
	return nil
}

func (s *NotificationService) validateSubscriptionConfirmedParams(params SubscriptionConfirmedParams) error {
	if err := s.validateEmail(params.Email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}
	if err := s.validateURL(params.DashboardURL); err != nil {
		return fmt.Errorf("invalid dashboard URL: %w", err)
	}
	if err := s.validateURL(params.BillingURL); err != nil {
		return fmt.Errorf("invalid billing URL: %w", err)
	}
	if params.InvoiceURL != "" {
		if err := s.validateURL(params.InvoiceURL); err != nil {
			return fmt.Errorf("invalid invoice URL: %w", err)
		}
	}
	if params.UserName == "" {
		return fmt.Errorf("user name is required")
	}
	if params.PlanName == "" {
		return fmt.Errorf("plan name is required")
	}
	return nil
}

func (s *NotificationService) validatePaymentSucceededParams(params PaymentSucceededParams) error {
	if err := s.validateEmail(params.Email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}
	if err := s.validateURL(params.DashboardURL); err != nil {
		return fmt.Errorf("invalid dashboard URL: %w", err)
	}
	if params.InvoiceURL != "" {
		if err := s.validateURL(params.InvoiceURL); err != nil {
			return fmt.Errorf("invalid invoice URL: %w", err)
		}
	}
	if params.UserName == "" {
		return fmt.Errorf("user name is required")
	}
	if params.Amount == "" {
		return fmt.Errorf("amount is required")
	}
	return nil
}
