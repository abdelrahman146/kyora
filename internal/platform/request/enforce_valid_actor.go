package request

import (
	"context"
	"errors"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/abdelrahman146/kyora/internal/platform/types/ctxkey"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"

	"github.com/gin-gonic/gin"
)

type EnforceValidActorAuthService interface {
	GetUserByID(ctx context.Context, userID string) (*account.User, error)
}

var ActorKey = ctxkey.New("actor")

func EnforceValidActor(authService EnforceValidActorAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := ClaimsFromContext(c)
		if err != nil {
			response.Error(c, err)
			return
		}
		user, err := authService.GetUserByID(c.Request.Context(), claims.UserID)
		if err != nil {
			response.Error(c, err)
			return
		}
		c.Set(ActorKey, user)
		c.Next()
	}
}

func ActorFromContext(c *gin.Context) (*account.User, error) {
	user, exists := c.Get(ActorKey)
	if !exists {
		logger.FromContext(c.Request.Context()).Error("user not found in context, make sure EnforceValidActor middleware is applied")
		return nil, problem.InternalError().WithError(errors.New("user not found in context"))
	}
	if user, ok := user.(*account.User); ok {
		return user, nil
	}
	return nil, problem.InternalError().WithError(errors.New("unable to cast user from context"))
}
