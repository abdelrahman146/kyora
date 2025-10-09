package middleware

import (
	"github.com/abdelrahman146/kyora/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthRequiredMiddleware(c *gin.Context) {
	jwtToken := utils.JWT.GetJwtFromContext(c)
	if jwtToken == "" {
		c.Redirect(302, loginPath)
		return
	}
	// verify jwtToken
	claims, err := utils.JWT.ParseToken(jwtToken)
	if err != nil {
		c.Redirect(302, loginPath)
		return
	}
	c.Set(ClaimsKey, claims)
	c.Next()
}
