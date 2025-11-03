package request

import (
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/abdelrahman146/kyora/internal/platform/types/schema"
	"github.com/gin-gonic/gin"
)

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
