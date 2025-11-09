package account

import (
	"errors"
	"fmt"

	"github.com/abdelrahman146/kyora/internal/platform/auth"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/abdelrahman146/kyora/internal/platform/types/ctxkey"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/gin-gonic/gin"
)

func EnforceActorPermissions(action role.Action, resource role.Resource) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := ActorFromContext(c)
		if err != nil {
			response.Error(c, err)
			return
		}
		role := user.Role
		if err := role.HasPermission(action, resource); err != nil {
			response.Error(c, err)
			return
		}
		c.Next()
	}
}

var ActorKey = ctxkey.New("actor")

func EnforceValidActor(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := auth.ClaimsFromContext(c)
		if err != nil {
			response.Error(c, err)
			return
		}
		user, err := service.GetUserByID(c.Request.Context(), claims.UserID)
		if err != nil {
			response.Error(c, err)
			return
		}
		l := logger.FromContext(c.Request.Context())
		l.With("actorID", user.ID, "actorEmail", user.Email, "actorName", fmt.Sprintf("%s %s", user.FirstName, user.LastName), "actorRole", user.Role)
		ctx := logger.WithContext(c.Request.Context(), l)
		c.Request = c.Request.WithContext(ctx)
		c.Set(ActorKey, user)
		c.Next()
	}
}

func ActorFromContext(c *gin.Context) (*User, error) {
	user, exists := c.Get(ActorKey)
	if !exists {
		logger.FromContext(c.Request.Context()).Error("user not found in context, make sure EnforceValidActor middleware is applied")
		return nil, problem.InternalError().WithError(errors.New("user not found in context"))
	}
	if user, ok := user.(*User); ok {
		return user, nil
	}
	return nil, problem.InternalError().WithError(errors.New("unable to cast user from context"))
}

var WorkspaceKey = ctxkey.New("workspace")

func EnforceWorkspaceMembership(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := ActorFromContext(c)
		if err != nil {
			response.Error(c, err)
			return
		}
		workspaceID := c.Param("workspaceId")

		if user.WorkspaceID != workspaceID {
			response.Error(c, problem.Forbidden("user is not a member of the workspace"))
			return
		}
		workspace, err := service.GetWorkspaceByID(c.Request.Context(), workspaceID)
		if err != nil {
			response.Error(c, err)
			return
		}
		l := logger.FromContext(c.Request.Context())
		l.With("workspaceID", workspace.ID)
		ctx := logger.WithContext(c.Request.Context(), l)
		c.Request = c.Request.WithContext(ctx)
		c.Set(WorkspaceKey, workspace)
		c.Next()
	}
}

func WorkspaceFromContext(c *gin.Context) (*Workspace, error) {
	workspace, exists := c.Get(WorkspaceKey)
	if !exists {
		logger.FromContext(c.Request.Context()).Error("workspace not found in context, make sure EnforceWorkspaceMembership middleware is applied")
		return nil, problem.InternalError().WithError(errors.New("workspace not found in context"))
	}
	if workspace, ok := workspace.(*Workspace); ok {
		return workspace, nil
	}
	logger.FromContext(c.Request.Context()).Error("unable to cast workspace from context, make sure EnforceWorkspaceMembership middleware is applied")
	return nil, problem.InternalError().WithError(errors.New("unable to cast workspace from context"))
}
