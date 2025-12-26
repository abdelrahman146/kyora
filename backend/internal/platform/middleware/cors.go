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
//
// Configuration is intentionally simple:
//   - cors.enabled: if false, CORS is completely disabled (useful for local development)
//   - cors.allowed_origins: array of specific origins to allow, or empty to allow all
//
// When CORS is enabled:
//   - If allowed_origins is empty: allows all origins with credentials
//   - If allowed_origins has values: only allows those specific origins with credentials
//   - All HTTP methods are allowed
//   - All headers are allowed
//   - Credentials (cookies, auth headers) are always supported when enabled
//
// Examples:
//
//	# Disable CORS entirely (local development)
//	cors:
//	  enabled: false
//
//	# Allow all origins (staging/testing)
//	cors:
//	  enabled: true
//	  allowed_origins: []
//
//	# Allow specific origins (production)
//	cors:
//	  enabled: true
//	  allowed_origins:
//	    - "https://portal.kyora.io"
//	    - "https://app.kyora.io"
func NewCORSMiddleware() gin.HandlerFunc {
	enabled := viper.GetBool(config.CORSEnabled)

	// If CORS is disabled, return a no-op middleware
	if !enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	allowedOrigins := viper.GetStringSlice(config.CORSAllowedOrigins)

	// If no specific origins configured, allow all with credentials
	if len(allowedOrigins) == 0 {
		return cors.New(cors.Config{
			AllowAllOrigins: false,
			AllowOriginFunc: func(origin string) bool {
				return true // Allow any origin
			},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
			AllowHeaders:     []string{"*"},
			ExposeHeaders:    []string{"*"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		})
	}

	// Allow only specific origins
	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// NewPublicCORSMiddleware creates CORS middleware for public endpoints
// that should be accessible from any origin (e.g., storefront, webhooks).
// This always allows all origins and supports credentials.
func NewPublicCORSMiddleware() gin.HandlerFunc {
	enabled := viper.GetBool(config.CORSEnabled)

	// If CORS is disabled, return a no-op middleware
	if !enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return cors.New(cors.Config{
		AllowAllOrigins: false,
		AllowOriginFunc: func(origin string) bool {
			return true // Allow any origin
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
