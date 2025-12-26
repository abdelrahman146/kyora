// Package middleware provides common HTTP middleware functions
package middleware

import (
	"slices"
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func corsConfigForOrigins(origins []string) cors.Config {
	// Kyora uses Bearer tokens (Authorization header). We do not rely on cookies,
	// so we keep AllowCredentials disabled to avoid the invalid "* + credentials"
	// combination and to simplify local development.
	cfg := cors.DefaultConfig()
	cfg.AllowCredentials = false
	cfg.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"}
	cfg.AllowHeaders = []string{"Authorization", "Content-Type", "X-Request-ID", "X-Trace-ID"}
	cfg.ExposeHeaders = []string{"X-Trace-ID"}
	cfg.MaxAge = 12 * time.Hour

	if len(origins) == 0 {
		origins = []string{"*"}
	}

	if slices.Contains(origins, "*") {
		cfg.AllowAllOrigins = true
		cfg.AllowOrigins = nil
		return cfg
	}

	cfg.AllowOrigins = origins
	return cfg
}

// NewCORSMiddleware creates CORS middleware from allowed_origins configuration.
// Just pass the origins you want to allow. Use "*" to allow all.
func NewCORSMiddleware() gin.HandlerFunc {
	allowedOrigins := viper.GetStringSlice(config.CORSAllowedOrigins)
	return cors.New(corsConfigForOrigins(allowedOrigins))
}

// NewPublicCORSMiddleware creates CORS middleware for public endpoints.
// Always allows all origins.
func NewPublicCORSMiddleware() gin.HandlerFunc {
	return cors.New(corsConfigForOrigins([]string{"*"}))
}
