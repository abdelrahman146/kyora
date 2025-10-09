package middleware

import (
	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/gin-gonic/gin"
)

func GuestRequiredMiddleware(c *gin.Context) {
	jwtToken := utils.JWT.GetJwtFromContext(c)
	if jwtToken != "" {
		claims, err := utils.JWT.ParseToken(jwtToken[len("Bearer "):])
		if err == nil && claims != nil {
			c.Redirect(302, "/")
			return
		}
	}
	c.Next()
}
