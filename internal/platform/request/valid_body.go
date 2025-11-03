package request

import (
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/gin-gonic/gin"
)

func ValidBody(c *gin.Context, obj any) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		response.Error(c, problem.BadRequest("invalid request body").WithError(err))
		return err
	}
	return nil
}
