package logger

import (
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	whitelistPaths = []string{
		"/health",
		"/static",
	}
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if slices.Contains(whitelistPaths, c.Request.URL.Path) {
			c.Next()
			return
		}
		start := time.Now()
		traceID := id.Ksuid()
		traceIdKey := viper.GetString(config.HTTPTraceIDHeader)
		if traceIdKey == "" {
			traceIdKey = "X-Trace-ID"
		}
		if c.GetHeader(traceIdKey) != "" {
			traceID = c.GetHeader(traceIdKey)
		}
		c.Writer.Header().Set(traceIdKey, traceID)
		requestLogger := slog.Default().With("traceId", traceID)
		c.Request = c.Request.WithContext(WithContext(c.Request.Context(), requestLogger))
		requestLogger.Info("Request Started",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"clientIP", c.ClientIP(),
		)
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		switch {
		case status >= 500:
			requestLogger.Error("Request Completed With Failure",
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"status", c.Writer.Status(),
				"latency", fmt.Sprintf("%v", latency),
				"clientIP", c.ClientIP(),
				"error", c.Errors.ByType(gin.ErrorTypePrivate).String(),
			)
		case status >= 400:
			requestLogger.Warn("Request Completed",
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"status", c.Writer.Status(),
				"latency", fmt.Sprintf("%v", latency),
				"clientIP", c.ClientIP(),
				"error", c.Errors.ByType(gin.ErrorTypePrivate).String(),
			)
		default:
			requestLogger.Info("Request Completed",
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"status", c.Writer.Status(),
				"latency", fmt.Sprintf("%v", latency),
				"clientIP", c.ClientIP(),
			)
		}
	}
}
