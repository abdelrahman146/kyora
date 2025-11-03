// Package auth provides authentication utilities
// JWT Tokens
// Session Management
// OAuth2 Integration
// TOTP Generation and Validation
package request

import (
	"errors"

	"github.com/abdelrahman146/kyora/internal/platform/auth"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
	"github.com/abdelrahman146/kyora/internal/platform/types/ctxkey"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/gin-gonic/gin"
)

var (
	ClaimsKey = ctxkey.New("claims")
)

func EnforceAuthentication(c *gin.Context) {
	jwtToken := auth.JwtFromContext(c)
	if jwtToken == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		c.Abort()
		return
	}
	// verify jwtToken
	claims, err := auth.ParseJwtToken(jwtToken)
	if err != nil {
		c.JSON(401, gin.H{"error": "unauthorized"})
		c.Abort()
		return
	}
	c.Set(ClaimsKey, claims)
	c.Next()
}

func ClaimsFromContext(c *gin.Context) (*auth.CustomClaims, error) {
	claimsVal, exists := c.Get(ClaimsKey)
	if !exists {
		logger.FromContext(c.Request.Context()).Error("claims not found in context, make sure EnforceAuthentication middleware is applied")
		return nil, problem.InternalError().WithError(errors.New("claims not found in context"))
	}
	claims, ok := claimsVal.(*auth.CustomClaims)
	if !ok {
		logger.FromContext(c.Request.Context()).Error("unable to cast claims from context")
		return nil, problem.InternalError().WithError(errors.New("unable to cast claims from context"))
	}
	return claims, nil
}
