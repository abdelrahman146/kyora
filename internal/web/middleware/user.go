package middleware

import (
	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/utils"

	"github.com/gin-gonic/gin"
)

func UserMiddleware(authService *account.AuthenticationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsVal, exists := c.Get(ClaimsKey)
		if !exists {
			c.Redirect(302, loginPath)
			return
		}
		customClaims, ok := claimsVal.(*utils.CustomClaims)
		if !ok {
			c.Redirect(302, loginPath)
			return
		}
		user, err := authService.GetUserByID(c.Request.Context(), customClaims.UserID)
		if err != nil {
			utils.Log.FromContext(c.Request.Context()).Error("failed to get user by id from claims", "claims", customClaims)
			c.Redirect(302, loginPath)
			return
		}
		c.Set(UserKey, user)
		c.Next()
	}
}
