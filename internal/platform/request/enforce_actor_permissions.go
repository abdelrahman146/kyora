package request

import (
	"github.com/abdelrahman146/kyora/internal/platform/response"
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
