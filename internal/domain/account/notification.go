package account

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"strings"
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/email"
)

// Notification encapsulates email sending for account domain
type Notification struct {
	client email.Client
	info   email.EmailInfo
}

// NewNotification wires the email client and defaults. from/fromName can be overridden via email.WithFrom
func NewNotification(client email.Client, info email.EmailInfo) *Notification {
	return &Notification{client: client, info: info}
}

// SendForgotPasswordEmail sends a forgot password email using templates
func (n *Notification) SendForgotPasswordEmail(ctx context.Context, user *User, token string, expiryTime time.Time) error {
	if n.client == nil {
		return fmt.Errorf("email client not available")
	}
	logger := slog.With("action", "send_forgot_password", "user_id", user.ID, "email", user.Email)
	logger.InfoContext(ctx, "Sending forgot password email")

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", n.info.BaseURL, token)
	expiryHours := int(time.Until(expiryTime).Hours())
	if expiryHours < 1 {
		expiryHours = 1
	}

	data := map[string]any{
		"userName":     n.getUserDisplayName(user),
		"resetURL":     resetURL,
		"productName":  n.info.ProductName,
		"supportEmail": n.info.SupportEmail,
		"helpURL":      n.info.HelpURL,
		"expiryTime":   fmt.Sprintf("%d hours", expiryHours),
		"currentYear":  fmt.Sprintf("%d", time.Now().Year()),
	}
	from := n.info.FormattedFrom()
	if _, err := n.client.SendTemplate(ctx, email.TemplateForgotPassword, []string{user.Email}, from, "", data); err != nil {
		logger.ErrorContext(ctx, "Failed to send forgot password email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	logger.InfoContext(ctx, "Forgot password email sent successfully")
	return nil
}

// SendPasswordResetConfirmationEmail sends a password reset confirmation email
func (n *Notification) SendPasswordResetConfirmationEmail(ctx context.Context, user *User, clientIP string) error {
	if n.client == nil {
		return fmt.Errorf("email client not available")
	}
	logger := slog.With("action", "send_password_reset_confirmation", "user_id", user.ID, "email", user.Email)
	logger.InfoContext(ctx, "Sending password reset confirmation email")

	resetTime := time.Now()
	location := n.getLocationFromIP(clientIP)
	data := map[string]any{
		"userName":      n.getUserDisplayName(user),
		"loginURL":      fmt.Sprintf("%s/login", n.info.BaseURL),
		"productName":   n.info.ProductName,
		"supportEmail":  n.info.SupportEmail,
		"helpURL":       n.info.HelpURL,
		"resetDate":     resetTime.Format("January 2, 2006"),
		"resetTime":     resetTime.Format("3:04 PM MST"),
		"resetLocation": location,
		"resetIP":       clientIP,
		"currentYear":   fmt.Sprintf("%d", time.Now().Year()),
	}
	from := n.info.FormattedFrom()
	if _, err := n.client.SendTemplate(ctx, email.TemplatePasswordResetConfirm, []string{user.Email}, from, "", data); err != nil {
		logger.ErrorContext(ctx, "Failed to send password reset confirmation email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	logger.InfoContext(ctx, "Password reset confirmation email sent successfully")
	return nil
}

// SendEmailVerificationEmail sends an email verification email
func (n *Notification) SendEmailVerificationEmail(ctx context.Context, user *User, token string, expiryTime time.Time) error {
	if n.client == nil {
		return fmt.Errorf("email client not available")
	}
	logger := slog.With("action", "send_email_verification", "user_id", user.ID, "email", user.Email)
	logger.InfoContext(ctx, "Sending email verification email")

	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", n.info.BaseURL, token)
	expiryHours := int(time.Until(expiryTime).Hours())
	if expiryHours < 1 {
		expiryHours = 1
	}
	data := map[string]any{
		"userName":     n.getUserDisplayName(user),
		"verifyURL":    verifyURL,
		"productName":  n.info.ProductName,
		"supportEmail": n.info.SupportEmail,
		"helpURL":      n.info.HelpURL,
		"expiryTime":   fmt.Sprintf("%d hours", expiryHours),
		"currentYear":  fmt.Sprintf("%d", time.Now().Year()),
	}
	from := n.info.FormattedFrom()
	if _, err := n.client.SendTemplate(ctx, email.TemplateEmailVerification, []string{user.Email}, from, "", data); err != nil {
		logger.ErrorContext(ctx, "Failed to send email verification email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	logger.InfoContext(ctx, "Email verification email sent successfully")
	return nil
}

// SendWelcomeEmail sends a welcome email to a new user
func (n *Notification) SendWelcomeEmail(ctx context.Context, user *User) error {
	if n.client == nil {
		return fmt.Errorf("email client not available")
	}
	logger := slog.With("action", "send_welcome_email", "user_id", user.ID, "email", user.Email)
	logger.InfoContext(ctx, "Sending welcome email")

	data := map[string]any{
		"userName":     n.getUserDisplayName(user),
		"dashboardURL": fmt.Sprintf("%s/dashboard", n.info.BaseURL),
		"guideURL":     fmt.Sprintf("%s/getting-started", n.info.BaseURL),
		"productName":  n.info.ProductName,
		"supportEmail": n.info.SupportEmail,
		"helpURL":      n.info.HelpURL,
		"currentYear":  fmt.Sprintf("%d", time.Now().Year()),
	}
	from := n.info.FormattedFrom()
	if _, err := n.client.SendTemplate(ctx, email.TemplateWelcome, []string{user.Email}, from, "", data); err != nil {
		logger.ErrorContext(ctx, "Failed to send welcome email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	logger.InfoContext(ctx, "Welcome email sent successfully")
	return nil
}

// SendLoginNotificationEmail sends a security notification email when user logs in
func (n *Notification) SendLoginNotificationEmail(ctx context.Context, user *User, clientIP, userAgent string) error {
	if n.client == nil {
		return fmt.Errorf("email client not available")
	}
	logger := slog.With("action", "send_login_notification", "user_id", user.ID, "email", user.Email)
	logger.InfoContext(ctx, "Sending login notification email")

	loginTime := time.Now()
	location := n.getLocationFromIP(clientIP)
	deviceInfo := n.parseUserAgent(userAgent)
	resetURL := fmt.Sprintf("%s/reset-password", n.info.BaseURL)

	data := map[string]any{
		"userName":      n.getUserDisplayName(user),
		"loginDate":     loginTime.Format("January 2, 2006"),
		"loginTime":     loginTime.Format("3:04 PM MST"),
		"loginLocation": location,
		"loginIP":       clientIP,
		"deviceInfo":    deviceInfo,
		"resetURL":      resetURL,
		"productName":   n.info.ProductName,
		"supportEmail":  n.info.SupportEmail,
		"helpURL":       n.info.HelpURL,
		"currentYear":   fmt.Sprintf("%d", time.Now().Year()),
	}
	from := n.info.FormattedFrom()
	if _, err := n.client.SendTemplate(ctx, email.TemplateLoginNotification, []string{user.Email}, from, "", data); err != nil {
		logger.ErrorContext(ctx, "Failed to send login notification email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	logger.InfoContext(ctx, "Login notification email sent successfully")
	return nil
}

// Helpers (moved from old integration file)
func (n *Notification) parseUserAgent(userAgent string) string {
	if userAgent == "" {
		return "Unknown device"
	}
	ua := strings.ToLower(userAgent)
	var os string
	if strings.Contains(ua, "windows") {
		os = "Windows"
	} else if strings.Contains(ua, "macintosh") || strings.Contains(ua, "mac os") {
		os = "macOS"
	} else if strings.Contains(ua, "linux") {
		os = "Linux"
	} else if strings.Contains(ua, "android") {
		os = "Android"
	} else if strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") {
		os = "iOS"
	} else {
		os = "Unknown OS"
	}
	var browser string
	if strings.Contains(ua, "chrome") && !strings.Contains(ua, "edge") {
		browser = "Chrome"
	} else if strings.Contains(ua, "firefox") {
		browser = "Firefox"
	} else if strings.Contains(ua, "safari") && !strings.Contains(ua, "chrome") {
		browser = "Safari"
	} else if strings.Contains(ua, "edge") {
		browser = "Edge"
	} else if strings.Contains(ua, "opera") {
		browser = "Opera"
	} else {
		browser = "Unknown browser"
	}
	return fmt.Sprintf("%s on %s", browser, os)
}

func (n *Notification) getUserDisplayName(user *User) string {
	if user.FirstName != "" {
		if user.LastName != "" {
			return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
		}
		return user.FirstName
	}
	if user.LastName != "" {
		return user.LastName
	}
	parts := strings.Split(user.Email, "@")
	if len(parts) > 0 {
		return parts[0]
	}
	return "User"
}

func (n *Notification) getLocationFromIP(ip string) string {
	if ip == "" {
		return "Unknown location"
	}
	if net.ParseIP(ip) == nil {
		return "Unknown location"
	}
	if isPrivateIP(ip) {
		return "Local network"
	}
	return fmt.Sprintf("IP: %s", ip)
}

func isPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	if ip4 := ip.To4(); ip4 != nil {
		if ip4[0] == 127 {
			return true
		}
		if ip4[0] == 10 {
			return true
		}
		if ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31 {
			return true
		}
		if ip4[0] == 192 && ip4[1] == 168 {
			return true
		}
	}
	if ip.IsLoopback() {
		return true
	}
	return false
}

// Backward compatibility with previous type name used by services
type EmailIntegration = Notification

// NewEmailIntegration is a backward-compatible constructor
func NewEmailIntegration(client email.Client, info email.EmailInfo) *EmailIntegration {
	n := NewNotification(client, info)
	return (*EmailIntegration)(n)
}
