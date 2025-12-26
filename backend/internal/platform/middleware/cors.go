// Package middleware provides common HTTP middleware functions
package middleware

import (
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// NewCORSMiddleware creates CORS middleware from allowed_origins configuration.
// Just pass the origins you want to allow. Use "*" to allow all.
func NewCORSMiddleware() gin.HandlerFunc {
	allowedOrigins := viper.GetStringSlice(config.CORSAllowedOrigins)

	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// NewPublicCORSMiddleware creates CORS middleware for public endpoints.
// Always allows all origins.
func NewPublicCORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
