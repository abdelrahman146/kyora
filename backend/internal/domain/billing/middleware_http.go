package billing

import (
	"context"
	"errors"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/abdelrahman146/kyora/internal/platform/types/ctxkey"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/abdelrahman146/kyora/internal/platform/types/schema"
	"github.com/gin-gonic/gin"
)

var (
	SubscriptionKey = ctxkey.New("subscription")
)

type EnforceActiveSubscriptionBillingService interface {
	GetSubscriptionByWorkspaceID(ctx context.Context, workspaceID string) (*Subscription, error)
}

func EnforceActiveSubscription(billingService EnforceActiveSubscriptionBillingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := account.ActorFromContext(c)
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
		l := logger.FromContext(c.Request.Context())
		l.With("subscriptionID", subscription.ID)
		ctx := logger.WithContext(c.Request.Context(), l)
		c.Request = c.Request.WithContext(ctx)
		c.Set(SubscriptionKey, subscription)
		c.Next()
	}
}

func SubscriptionFromContext(c *gin.Context) (*Subscription, error) {
	subscription, exists := c.Get(SubscriptionKey)
	if !exists {
		logger.FromContext(c.Request.Context()).Error("subscription not found in context, make sure EnforceActiveSubscription middleware is applied")
		return nil, problem.InternalError().WithError(errors.New("subscription not found in context"))
	}
	if subscription, ok := subscription.(*Subscription); ok {
		return subscription, nil
	}
	logger.FromContext(c.Request.Context()).Error("unable to cast subscription from context, make sure EnforceActiveSubscription middleware is applied")
	return nil, problem.InternalError().WithError(errors.New("unable to cast subscription from context"))
}

type EnforcePlanLimitFunc func(ctx context.Context, actor *account.User, id string) (int64, error)

func EnforcePlanBusinessLimits(feature schema.Field, enforceFunc EnforcePlanLimitFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		sub, err := SubscriptionFromContext(c)
		if err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}
		actor, err := account.ActorFromContext(c)
		if err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}
		business, err := business.BusinessFromContext(c)
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

func EnforcePlanWorkspaceLimits(feature schema.Field, enforceFunc EnforcePlanLimitFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		sub, err := SubscriptionFromContext(c)
		if err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}
		actor, err := account.ActorFromContext(c)
		if err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}
		workspace, err := account.WorkspaceFromContext(c)
		if err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}
		usage, err := enforceFunc(c.Request.Context(), actor, workspace.ID)
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

func EnforcePlanFeatureRestriction(feature schema.Field) gin.HandlerFunc {
	return func(c *gin.Context) {
		sub, err := SubscriptionFromContext(c)
		if err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}
		if err := sub.Plan.Features.CanUseFeature(feature); err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}
		c.Next()
	}
}
