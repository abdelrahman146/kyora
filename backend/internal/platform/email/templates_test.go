package email_test

import (
	"testing"

	"github.com/abdelrahman146/kyora/internal/platform/email"
	"github.com/stretchr/testify/require"
)

func TestRenderTemplate_OnboardingEmailOTP_RendersOTPCode(t *testing.T) {
	html, err := email.RenderTemplate(email.TemplateOnboardingEmailOTP, map[string]any{
		"currentYear":  "2025",
		"expiryTime":   "15 minutes",
		"productName":  "Kyora",
		"supportEmail": "support@kyora.com",
		"helpURL":      "https://help.kyora.com",
		"userName":     "Test User",
		"otpCode":      "123456",
	})
	require.NoError(t, err)
	require.Contains(t, html, "123456")
}

func TestRenderTemplate_OnboardingEmailOTP_MissingCode_DoesNotError(t *testing.T) {
	_, err := email.RenderTemplate(email.TemplateOnboardingEmailOTP, map[string]any{
		"currentYear":  "2025",
		"expiryTime":   "15 minutes",
		"productName":  "Kyora",
		"supportEmail": "support@kyora.com",
		"helpURL":      "https://help.kyora.com",
		"userName":     "Test User",
	})
	require.NoError(t, err)
}
