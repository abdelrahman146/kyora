// Package middleware provides common HTTP middleware functions
package middleware

import (
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// NewCORSMiddleware creates CORS middleware based on configuration.
// It supports:
//   - Configurable allowed origins (defaults to localhost:5173, localhost:5174 for dev)
//   - Configurable allowed methods, headers, and exposed headers
//   - Optional allow-all-origins mode (use with caution, only for public APIs)
//   - Credentials support (required for authenticated endpoints with cookies/tokens)
//
// Configuration keys:
//   - cors.allowed_origins: array of allowed origin URLs
//   - cors.allowed_methods: array of allowed HTTP methods
//   - cors.allowed_headers: array of allowed request headers
//   - cors.expose_headers: array of headers to expose to clients
//   - cors.max_age: preflight cache duration in seconds
//   - cors.allow_all_origins: bool to allow all origins (disables credentials)
func NewCORSMiddleware() gin.HandlerFunc {
	allowedOrigins := viper.GetStringSlice(config.CORSAllowedOrigins)
	allowedMethods := viper.GetStringSlice(config.CORSAllowedMethods)
	allowedHeaders := viper.GetStringSlice(config.CORSAllowedHeaders)
	exposeHeaders := viper.GetStringSlice(config.CORSExposeHeaders)
	maxAge := viper.GetInt64(config.CORSMaxAge)
	allowAllOrigins := viper.GetBool(config.CORSAllowAllOrigins)

	corsConfig := cors.Config{
		AllowAllOrigins: allowAllOrigins,
		AllowOrigins:    allowedOrigins,
		AllowMethods:    allowedMethods,
		AllowHeaders:    allowedHeaders,
		ExposeHeaders:   exposeHeaders,
		MaxAge:          time.Duration(maxAge) * time.Second,
	}

	// Enable credentials support unless allow-all-origins is enabled
	// (credentials cannot be used with wildcard origins for security)
	if !allowAllOrigins {
		corsConfig.AllowCredentials = true
	}

	return cors.New(corsConfig)
}

// NewPublicCORSMiddleware creates CORS middleware for public endpoints
// that should be accessible from any origin (e.g., storefront, webhooks).
// This uses AllowAllOrigins=true and does NOT support credentials.
func NewPublicCORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Idempotency-Key",
			"X-Trace-ID",
		},
		ExposeHeaders: []string{"X-Trace-ID"},
		MaxAge:        12 * time.Hour,
	})
}
