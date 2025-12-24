package response

import (
	"errors"

	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, err error) {
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePrivate)
	}

	var p *problem.Problem
	if !errors.As(err, &p) {
		if database.IsRecordNotFound(err) {
			p = problem.NotFound("resource not found").WithError(err)
		} else if database.IsUniqueViolation(err) {
			p = problem.Conflict("resource already exists").WithError(err)
		} else {
			p = problem.InternalError().WithError(err)
		}
	}
	if p.Instance == "" && c.Request != nil && c.Request.URL != nil {
		p.Instance = c.Request.URL.Path
	}
	p.ServeJSON(c.Writer)
	c.Abort()
}

func SuccessJSON(c *gin.Context, status int, data any) {
	c.JSON(status, data)
}

func SuccessEmpty(c *gin.Context, status int) {
	c.Status(status)
}

func SuccessText(c *gin.Context, status int, data string) {
	c.String(status, data)
}
