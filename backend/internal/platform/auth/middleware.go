package auth

import (
	"errors"

	"github.com/abdelrahman146/kyora/internal/platform/logger"
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/abdelrahman146/kyora/internal/platform/types/ctxkey"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/gin-gonic/gin"
)

var (
	ClaimsKey = ctxkey.New("claims")
)

func EnforceAuthentication(c *gin.Context) {
	jwtToken := JwtFromContext(c)
	if jwtToken == "" {
		response.Error(c, problem.Unauthorized("unauthorized"))
		return
	}
	// verify jwtToken
	claims, err := ParseJwtToken(jwtToken)
	if err != nil {
		response.Error(c, problem.Unauthorized("unauthorized").WithError(err))
		return
	}
	c.Set(ClaimsKey, claims)
	c.Next()
}

func ClaimsFromContext(c *gin.Context) (*CustomClaims, error) {
	claimsVal, exists := c.Get(ClaimsKey)
	if !exists {
		logger.FromContext(c.Request.Context()).Error("claims not found in context, make sure EnforceAuthentication middleware is applied")
		return nil, problem.InternalError().WithError(errors.New("claims not found in context"))
	}
	claims, ok := claimsVal.(*CustomClaims)
	if !ok {
		logger.FromContext(c.Request.Context()).Error("unable to cast claims from context")
		return nil, problem.InternalError().WithError(errors.New("unable to cast claims from context"))
	}
	return claims, nil
}
