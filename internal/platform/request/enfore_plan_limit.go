package request

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/abdelrahman146/kyora/internal/platform/types/schema"
	"github.com/gin-gonic/gin"
)

type EnforcePlanLimitFunc func(ctx context.Context, actor *account.User, businessID string) (int64, error)

func EnforcePlanLimit(feature schema.Field, enforceFunc EnforcePlanLimitFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		sub, err := SubscriptionFromContext(c)
		if err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}
		actor, err := ActorFromContext(c)
		if err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}
		business, err := BusinessFromContext(c)
		if err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}
		usage, err := enforceFunc(c.Request.Context(), actor, business.ID)
		if err != nil {
			response.Error(c, err)
		}
		if err := sub.Plan.Limits.CheckUsageLimit(feature, usage); err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}
		c.Next()
	}
}
