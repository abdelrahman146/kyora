// Package config provides configuration utilities
// implements viper
package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Configuration keys
const (
	// environment
	Env = "env"
	// app configuration
	AppName               = "app.name"
	AppPort               = "app.port"
	AppDomain             = "app.domain"
	AppNotificationsEmail = "app.notifications_email"
	// log configuration
	LogFormat = "log.format"
	LogLevel  = "log.level"
	// http configuration
	HTTPPort          = "http.port"
	HTTPBaseURL       = "http.base_url"
	HTTPTraceIDHeader = "http.trace_id_header"
	HTTPMaxBodyBytes  = "http.max_body_bytes"
	// database configuration
	DatabaseDSN          = "database.dsn"
	DatabaseMaxOpenConns = "database.max_open_conns"
	DatabaseMaxIdleConns = "database.max_idle_conns"
	DatabaseMaxIdleTime  = "database.max_idle_time"
	DatabaseLogLevel     = "database.log_level"
	// cache configuration
	CacheHosts = "cache.hosts"
	// jwt configuration
	JWTSecret        = "auth.jwt.secret"
	JWTExpirySeconds = "auth.jwt.expiry_seconds"
	JWTIssuer        = "auth.jwt.issuer"
	JWTAudience      = "auth.jwt.audience"
	// Password reset configuration
	PasswordResetTokenExpirySeconds = "auth.password_reset_ttl_seconds"
	// Email verification configuration
	VerifyEmailTokenExpirySeconds = "auth.verify_email_ttl_seconds"
	// Workspace invitation configuration
	WorkspaceInvitationTokenExpirySeconds = "auth.invitation_token_ttl_seconds"
	// Google OAuth configuration
	GoogleOAuthClientID     = "auth.google_oauth.client_id"
	GoogleOAuthClientSecret = "auth.google_oauth.client_secret"
	GoogleOAuthRedirectURL  = "auth.google_oauth.redirect_url"
	// stripe configuration
	StripeAPIKey        = "billing.stripe.api_key"
	StripeWebhookSecret = "billing.stripe.webhook_secret"
	// billing configuration
	BillingAutoSyncPlans = "billing.auto_sync_plans" // bool - automatically sync plans on startup (default: true)
	// email configuration
	EmailProvider     = "email.provider"        // values: resend, mock
	EmailMockEnabled  = "email.mock.enabled"    // bool toggle to force mock
	ResendAPIKey      = "email.resend.api_key"  // API key for Resend
	ResendAPIBaseURL  = "email.resend.base_url" // optional, defaults to https://api.resend.com
	EmailFromEmail    = "email.from_email"      // default From email address
	EmailFromName     = "email.from_name"       // default From display name
	EmailSupportEmail = "email.support_email"   // support contact email
	EmailHelpURL      = "email.help_url"        // help/knowledge base URL
)

func init() {
	viper.SetConfigName(".kyora")                 // name of config file (without extension)
	viper.SetConfigType("yaml")                   // REQUIRED if the config file does not have the extension in the name
	viper.SetDefault(HTTPMaxBodyBytes, 1024*1024) // 1 MiB default max request body
	// Add current directory first
	viper.AddConfigPath(".")
	// Attempt to discover project root (where .kyora.yaml resides) by walking parent dirs
	if wd, err := os.Getwd(); err == nil {
		dir := wd
		for i := 0; i < 6; i++ { // limit depth to avoid infinite loops
			candidate := filepath.Join(dir, ".kyora.yaml")
			if _, statErr := os.Stat(candidate); statErr == nil {
				viper.AddConfigPath(dir)
				break
			}
			parent := filepath.Dir(dir)
			if parent == dir { // reached filesystem root
				break
			}
			dir = parent
		}
	}
	viper.AutomaticEnv() // read in environment variables that match
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}
