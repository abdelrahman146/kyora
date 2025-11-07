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

// EmailIntegration provides methods to integrate account operations with email notifications
type EmailIntegration struct {
	notificationService *email.NotificationService
	baseURL             string
}

// NewEmailIntegration creates a new email integration helper
func NewEmailIntegration(notificationService *email.NotificationService, baseURL string) *EmailIntegration {
	if baseURL == "" {
		baseURL = "https://app.kyora.com"
	}
	return &EmailIntegration{
		notificationService: notificationService,
		baseURL:             baseURL,
	}
}

// SendForgotPasswordEmail sends a forgot password email using the notification service
func (e *EmailIntegration) SendForgotPasswordEmail(ctx context.Context, user *User, token string, expiryTime time.Time) error {
	if e.notificationService == nil {
		return fmt.Errorf("notification service not available")
	}

	logger := slog.With("action", "send_forgot_password", "user_id", user.ID, "email", user.Email)
	logger.InfoContext(ctx, "Sending forgot password email")

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", e.baseURL, token)
	expiryHours := int(time.Until(expiryTime).Hours())
	if expiryHours < 1 {
		expiryHours = 1 // Minimum 1 hour for display
	}

	params := email.ForgotPasswordParams{
		Email:      user.Email,
		UserName:   e.getUserDisplayName(user),
		ResetURL:   resetURL,
		ExpiryTime: fmt.Sprintf("%d hours", expiryHours),
	}

	err := e.notificationService.SendForgotPasswordEmail(ctx, params)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to send forgot password email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.InfoContext(ctx, "Forgot password email sent successfully")
	return nil
}

// SendPasswordResetConfirmationEmail sends a password reset confirmation email
func (e *EmailIntegration) SendPasswordResetConfirmationEmail(ctx context.Context, user *User, clientIP string) error {
	if e.notificationService == nil {
		return fmt.Errorf("notification service not available")
	}

	logger := slog.With("action", "send_password_reset_confirmation", "user_id", user.ID, "email", user.Email)
	logger.InfoContext(ctx, "Sending password reset confirmation email")

	resetTime := time.Now()
	location := e.getLocationFromIP(clientIP)

	params := email.PasswordResetConfirmationParams{
		Email:         user.Email,
		UserName:      e.getUserDisplayName(user),
		LoginURL:      fmt.Sprintf("%s/login", e.baseURL),
		ResetDate:     resetTime.Format("January 2, 2006"),
		ResetTime:     resetTime.Format("3:04 PM MST"),
		ResetLocation: location,
		ResetIP:       clientIP,
	}

	err := e.notificationService.SendPasswordResetConfirmationEmail(ctx, params)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to send password reset confirmation email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.InfoContext(ctx, "Password reset confirmation email sent successfully")
	return nil
}

// SendEmailVerificationEmail sends an email verification email
func (e *EmailIntegration) SendEmailVerificationEmail(ctx context.Context, user *User, token string, expiryTime time.Time) error {
	if e.notificationService == nil {
		return fmt.Errorf("notification service not available")
	}

	logger := slog.With("action", "send_email_verification", "user_id", user.ID, "email", user.Email)
	logger.InfoContext(ctx, "Sending email verification email")

	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", e.baseURL, token)
	expiryHours := int(time.Until(expiryTime).Hours())
	if expiryHours < 1 {
		expiryHours = 1 // Minimum 1 hour for display
	}

	params := email.EmailVerificationParams{
		Email:      user.Email,
		UserName:   e.getUserDisplayName(user),
		VerifyURL:  verifyURL,
		ExpiryTime: fmt.Sprintf("%d hours", expiryHours),
	}

	err := e.notificationService.SendEmailVerificationEmail(ctx, params)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to send email verification email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.InfoContext(ctx, "Email verification email sent successfully")
	return nil
}

// SendWelcomeEmail sends a welcome email to a new user
func (e *EmailIntegration) SendWelcomeEmail(ctx context.Context, user *User) error {
	if e.notificationService == nil {
		return fmt.Errorf("notification service not available")
	}

	logger := slog.With("action", "send_welcome_email", "user_id", user.ID, "email", user.Email)
	logger.InfoContext(ctx, "Sending welcome email")

	params := email.WelcomeParams{
		Email:        user.Email,
		UserName:     e.getUserDisplayName(user),
		DashboardURL: fmt.Sprintf("%s/dashboard", e.baseURL),
		GuideURL:     fmt.Sprintf("%s/getting-started", e.baseURL),
	}

	err := e.notificationService.SendWelcomeEmail(ctx, params)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to send welcome email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.InfoContext(ctx, "Welcome email sent successfully")
	return nil
}

// SendLoginNotificationEmail sends a security notification email when user logs in
func (e *EmailIntegration) SendLoginNotificationEmail(ctx context.Context, user *User, clientIP, userAgent string) error {
	if e.notificationService == nil {
		return fmt.Errorf("notification service not available")
	}

	logger := slog.With("action", "send_login_notification", "user_id", user.ID, "email", user.Email)
	logger.InfoContext(ctx, "Sending login notification email")

	loginTime := time.Now()
	location := e.getLocationFromIP(clientIP)
	deviceInfo := e.parseUserAgent(userAgent)
	resetURL := fmt.Sprintf("%s/reset-password", e.baseURL)

	params := email.LoginNotificationParams{
		Email:         user.Email,
		UserName:      e.getUserDisplayName(user),
		LoginDate:     loginTime.Format("January 2, 2006"),
		LoginTime:     loginTime.Format("3:04 PM MST"),
		LoginLocation: location,
		LoginIP:       clientIP,
		DeviceInfo:    deviceInfo,
		ResetURL:      resetURL,
		SupportEmail:  "support@kyora.com", // Could be made configurable
	}

	err := e.notificationService.SendLoginNotificationEmail(ctx, params)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to send login notification email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.InfoContext(ctx, "Login notification email sent successfully")
	return nil
}

// parseUserAgent extracts device/browser info from user agent string
func (e *EmailIntegration) parseUserAgent(userAgent string) string {
	if userAgent == "" {
		return "Unknown device"
	}

	// Simple user agent parsing - in production you might use a proper library
	ua := strings.ToLower(userAgent)

	// Detect operating system
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

	// Detect browser
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

// getUserDisplayName returns a user-friendly display name
func (e *EmailIntegration) getUserDisplayName(user *User) string {
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
	parts := strings.Split(user.Email, "@")
	if len(parts) > 0 {
		return parts[0]
	}
	return "User"
}

// getLocationFromIP attempts to get a user-friendly location from IP address
func (e *EmailIntegration) getLocationFromIP(ip string) string {
	if ip == "" {
		return "Unknown location"
	}

	// Basic IP validation
	if net.ParseIP(ip) == nil {
		return "Unknown location"
	}

	// Check for local/private IPs
	if e.isPrivateIP(ip) {
		return "Local network"
	}

	// For now, just return the IP. In a production system, you would:
	// 1. Use a geolocation service like MaxMind GeoIP or IPinfo
	// 2. Cache results to avoid repeated lookups
	// 3. Handle rate limiting and errors gracefully
	return fmt.Sprintf("IP: %s", ip)
}

// isPrivateIP checks if an IP address is private/local
func (e *EmailIntegration) isPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	// Check for IPv4 private ranges
	if ip4 := ip.To4(); ip4 != nil {
		// 127.0.0.0/8 (localhost)
		if ip4[0] == 127 {
			return true
		}
		// 10.0.0.0/8
		if ip4[0] == 10 {
			return true
		}
		// 172.16.0.0/12
		if ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31 {
			return true
		}
		// 192.168.0.0/16
		if ip4[0] == 192 && ip4[1] == 168 {
			return true
		}
	}

	// Check for IPv6 loopback
	if ip.IsLoopback() {
		return true
	}

	return false
}
