package middleware

import (
	"fmt"
	"slices"
	"time"

	"github.com/abdelrahman146/kyora/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	whitelistPaths = []string{
		"/health",
		"/static",
	}
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if slices.Contains(whitelistPaths, c.Request.URL.Path) {
			c.Next()
			return
		}
		start := time.Now()
		traceID := utils.ID.NewKsuid()
		traceIdKey := viper.GetString("server.trace_id_header")
		if traceIdKey == "" {
			traceIdKey = "X-Trace-ID"
		}
		if c.GetHeader(traceIdKey) != "" {
			traceID = c.GetHeader(traceIdKey)
		}
		c.Writer.Header().Set(traceIdKey, traceID)
		requestLogger := utils.Log.FromContext(c.Request.Context()).With("traceId", traceID)
		c.Request = c.Request.WithContext(utils.Log.WithContext(c.Request.Context(), requestLogger))
		requestLogger.Info("Request Started",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"clientIP", c.ClientIP(),
		)
		c.Next()
		latency := time.Since(start)
		requestLogger.Info("Request Completed",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency", fmt.Sprintf("%v", latency),
			"clientIP", c.ClientIP(),
		)
	}
}
