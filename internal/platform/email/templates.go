package email

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"strings"
)

//go:embed templates/*.html
var templatesFS embed.FS

const (
	// Authentication and Account Management Templates
	TemplateForgotPassword       TemplateID = "forgot_password"
	TemplatePasswordResetConfirm TemplateID = "password_reset_confirmation"
	TemplateEmailVerification    TemplateID = "email_verification"
	TemplateWelcome              TemplateID = "welcome"
	TemplateLoginNotification    TemplateID = "login_notification"

	// Billing and Subscription Templates
	TemplateSubscriptionWelcome   TemplateID = "subscription_welcome"
	TemplatePaymentFailed         TemplateID = "payment_failed"
	TemplateSubscriptionCanceled  TemplateID = "subscription_canceled"
	TemplateTrialEnding           TemplateID = "trial_ending"
	TemplatePaymentSucceeded      TemplateID = "payment_succeeded"
	TemplateSubscriptionUpdated   TemplateID = "subscription_updated"
	TemplateInvoiceGenerated      TemplateID = "invoice_generated"
	TemplateSubscriptionConfirmed TemplateID = "subscription_confirmed"
)

// registry maps TemplateID to file paths within the embedded FS
var templateFiles = map[TemplateID]string{
	// Authentication and Account Management Templates
	TemplateForgotPassword:       "templates/forgot_password.html",
	TemplatePasswordResetConfirm: "templates/password_reset_confirmation.html",
	TemplateEmailVerification:    "templates/email_verification.html",
	TemplateWelcome:              "templates/welcome.html",
	TemplateLoginNotification:    "templates/login_notification.html",

	// Billing and Subscription Templates
	TemplateSubscriptionWelcome:   "templates/subscription_welcome.html",
	TemplatePaymentFailed:         "templates/payment_failed.html",
	TemplateSubscriptionCanceled:  "templates/subscription_canceled.html",
	TemplateTrialEnding:           "templates/trial_ending.html",
	TemplatePaymentSucceeded:      "templates/payment_succeeded.html",
	TemplateSubscriptionConfirmed: "templates/subscription_confirmed.html",
}

// subjects maps TemplateID to a default subject line
var subjects = map[TemplateID]string{
	// Authentication and Account Management Templates
	TemplateForgotPassword:       "Reset your password",
	TemplatePasswordResetConfirm: "Password successfully reset",
	TemplateEmailVerification:    "Verify your email address",
	TemplateWelcome:              "Welcome to Kyora!",
	TemplateLoginNotification:    "New login to your account",

	// Billing and Subscription Templates
	TemplateSubscriptionWelcome:   "Welcome to your subscription!",
	TemplatePaymentFailed:         "Payment Failed - Action Required",
	TemplateSubscriptionCanceled:  "Subscription canceled - We're sorry to see you go",
	TemplateTrialEnding:           "Your trial is ending soon",
	TemplatePaymentSucceeded:      "Payment received - Thank you!",
	TemplateSubscriptionUpdated:   "Your subscription has been updated",
	TemplateInvoiceGenerated:      "Your invoice is ready",
	TemplateSubscriptionConfirmed: "Subscription confirmed - You're all set!",
}

// RenderTemplate renders the embedded HTML template with provided data.
// Missing keys render as empty strings.
func RenderTemplate(id TemplateID, data map[string]any) (string, error) {
	path, ok := templateFiles[id]
	if !ok {
		return "", ErrTemplateNotFound(string(id))
	}
	contentBytes, err := templatesFS.ReadFile(path)
	if err != nil {
		return "", err
	}
	// prepare funcs: default returns the first non-empty string
	funcMap := template.FuncMap{
		"default": func(def string, v any) string {
			// best-effort convert to string
			switch vv := v.(type) {
			case string:
				if strings.TrimSpace(vv) != "" {
					return vv
				}
			case []byte:
				s := strings.TrimSpace(string(vv))
				if s != "" {
					return s
				}
			}
			return def
		},
	}
	t, err := template.New(string(id)).Funcs(funcMap).Option("missingkey=zero").Parse(string(contentBytes))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// SubjectFor returns a sensible default subject for a template id.
func SubjectFor(id TemplateID) string {
	if s, ok := subjects[id]; ok {
		return s
	}
	return string(id)
}

// ErrTemplateNotFound returns an error when a template ID is not registered.
func ErrTemplateNotFound(id string) error { return fmt.Errorf("email template not found: %s", id) }
