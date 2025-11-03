package request

import (
	"context"
	"errors"

	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/abdelrahman146/kyora/internal/platform/types/ctxkey"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/gin-gonic/gin"
)

var (
	BusinessKey = ctxkey.New("business")
)

type businessRequiredBusinessService interface {
	GetBusinessByDescriptor(ctx context.Context, workspaceID string, descriptor string) (*business.Business, error)
}

func EnforceBusinessValidity(businessService businessRequiredBusinessService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := ActorFromContext(c)
		if err != nil {
			response.Error(c, err)
			return
		}
		descriptor := c.Param("businessDescriptor")
		biz, err := businessService.GetBusinessByDescriptor(c.Request.Context(), user.WorkspaceID, descriptor)
		if err != nil || biz == nil {
			response.Error(c, problem.NotFound("business not found"))
			return
		}
		c.Set(BusinessKey, biz)
		c.Next()
	}
}

func BusinessFromContext(c *gin.Context) (*business.Business, error) {
	biz, exists := c.Get(BusinessKey)
	if !exists {
		logger.FromContext(c.Request.Context()).Error("business not found in context, make sure EnforceBusinessValidity middleware is applied")
		return nil, problem.InternalError().WithError(errors.New("business not found in context"))
	}
	if b, ok := biz.(*business.Business); ok {
		return b, nil
	}
	return nil, problem.InternalError().WithError(errors.New("business not found in context"))
}
