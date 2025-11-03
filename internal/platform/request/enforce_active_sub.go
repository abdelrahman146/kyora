package request

import (
	"context"
	"errors"

	"github.com/abdelrahman146/kyora/internal/domain/billing"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/abdelrahman146/kyora/internal/platform/types/ctxkey"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/gin-gonic/gin"
)

var (
	SubscriptionKey = ctxkey.New("subscription")
)

type EnforceActiveSubscriptionBillingService interface {
	GetSubscriptionByWorkspaceID(ctx context.Context, workspaceID string) (*billing.Subscription, error)
}

func EnforceActiveSubscription(billingService EnforceActiveSubscriptionBillingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := ActorFromContext(c)
		if err != nil {
			response.Error(c, err)
			return
		}
		subscription, err := billingService.GetSubscriptionByWorkspaceID(c.Request.Context(), user.WorkspaceID)
		if err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}
		if err := subscription.IsActive(); err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}
		c.Set(SubscriptionKey, subscription)
		c.Next()
	}
}

func SubscriptionFromContext(c *gin.Context) (*billing.Subscription, error) {
	subscription, exists := c.Get(SubscriptionKey)
	if !exists {
		logger.FromContext(c.Request.Context()).Error("subscription not found in context, make sure EnforceActiveSubscription middleware is applied")
		return nil, problem.InternalError().WithError(errors.New("subscription not found in context"))
	}
	if subscription, ok := subscription.(*billing.Subscription); ok {
		return subscription, nil
	}
	logger.FromContext(c.Request.Context()).Error("unable to cast subscription from context, make sure EnforceActiveSubscription middleware is applied")
	return nil, problem.InternalError().WithError(errors.New("unable to cast subscription from context"))
}
