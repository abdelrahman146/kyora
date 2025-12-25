// Package config provides configuration utilities
// implements viper
package config

import (
	"errors"
	"fmt"
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
	// refresh token configuration
	RefreshTokenExpirySeconds = "auth.refresh_token_ttl_seconds"
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
	// database configuration (advanced)
	DatabaseAutoMigrate = "database.auto_migrate" // bool - auto-migrate models on startup (default: true)
	// email configuration
	EmailProvider     = "email.provider"        // values: resend, mock
	EmailMockEnabled  = "email.mock.enabled"    // bool toggle to force mock
	ResendAPIKey      = "email.resend.api_key"  // API key for Resend
	ResendAPIBaseURL  = "email.resend.base_url" // optional, defaults to https://api.resend.com
	EmailFromEmail    = "email.from_email"      // default From email address
	EmailFromName     = "email.from_name"       // default From display name
	EmailSupportEmail = "email.support_email"   // support contact email
	EmailHelpURL      = "email.help_url"        // help/knowledge base URL

	// blob storage / uploads
	StorageProvider        = "storage.provider"          // values: local, s3
	StorageBucket          = "storage.bucket"            // bucket/container name
	StorageRegion          = "storage.region"            // region (e.g., nyc3 for DO Spaces)
	StorageEndpoint        = "storage.endpoint"          // optional (e.g., https://nyc3.digitaloceanspaces.com)
	StorageAccessKeyID     = "storage.access_key_id"     // S3-compatible access key
	StorageSecretAccessKey = "storage.secret_access_key" // S3-compatible secret
	StoragePublicBaseURL   = "storage.public_base_url"   // optional public base URL (e.g., https://bucket.endpoint)

	UploadsMaxBytes = "uploads.max_bytes" // max file size in bytes for direct uploads
)

var configured bool

// Configure prepares Viper defaults and search paths.
// It is safe to call multiple times.
func Configure() {
	if configured {
		return
	}
	configured = true

	viper.SetConfigName(".kyora")
	viper.SetConfigType("yaml")

	// Defaults
	viper.SetDefault(HTTPMaxBodyBytes, int64(1024*1024)) // 1 MiB default max request body
	viper.SetDefault(BillingAutoSyncPlans, true)
	viper.SetDefault(DatabaseAutoMigrate, true)
	viper.SetDefault(StorageProvider, "local")
	viper.SetDefault(UploadsMaxBytes, int64(5*1024*1024)) // 5 MiB per upload default
	// Auth defaults
	// Refresh tokens are long-lived and rotated; keep configurable.
	viper.SetDefault(RefreshTokenExpirySeconds, int64(30*24*60*60)) // 30 days

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

	viper.AutomaticEnv()
}

type LoadOption func(*LoadOptions)

type LoadOptions struct {
	RequireConfigFile bool
}

// WithRequiredConfigFile makes Load return an error if the config file cannot be found.
func WithRequiredConfigFile(required bool) LoadOption {
	return func(o *LoadOptions) { o.RequireConfigFile = required }
}

// Load reads the config file if present.
// It never terminates the process; callers should handle returned errors.
func Load(opts ...LoadOption) error {
	Configure()
	options := &LoadOptions{}
	for _, opt := range opts {
		opt(options)
	}

	if err := viper.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if errors.As(err, &notFound) {
			if options.RequireConfigFile {
				return fmt.Errorf("config file not found: %w", err)
			}
			return nil
		}
		return fmt.Errorf("failed to read config: %w", err)
	}

	// Keep a log line for visibility in CLI usage; avoid fatal exits.
	if file := viper.ConfigFileUsed(); file != "" {
		log.Printf("Loaded config file: %s", file)
	}

	return nil
}
