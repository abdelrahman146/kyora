package webrouter

import (
	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/abdelrahman146/kyora/internal/web/webutils"

	"github.com/gin-gonic/gin"
)

func registerRoutes(r gin.IRoutes) {
	r.Static("/static", "./public")
	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.POST("/logout", func(c *gin.Context) {
		utils.JWT.ClearJwtCookie(c)
		webutils.Redirect(c, "/login")
	})
}
