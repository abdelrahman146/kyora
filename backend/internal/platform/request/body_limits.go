package request

import (
	"net/http"

	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/gin-gonic/gin"
)

// LimitBodySize enforces a maximum request body size.
//
// It rejects requests early when Content-Length is known and exceeds the limit,
// otherwise it wraps the request body with http.MaxBytesReader so downstream
// JSON decoding fails with *http.MaxBytesError.
func LimitBodySize(maxBytes int64) gin.HandlerFunc {
	if maxBytes <= 0 {
		return func(c *gin.Context) { c.Next() }
	}

	return func(c *gin.Context) {
		if c.Request.Body == nil {
			c.Next()
			return
		}

		if c.Request.ContentLength > maxBytes {
			response.Error(c, problem.PayloadTooLarge("request body too large").WithCode("request.body_too_large"))
			c.Abort()
			return
		}

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}
