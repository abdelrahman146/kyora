package webrouter

import (
	"github.com/abdelrahman146/kyora/internal/web/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())
	r.SetTrustedProxies(nil)
	registerRoutes(r)
	return r
}
