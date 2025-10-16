package webutils

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func GetPaginationParams(c *gin.Context) (page int, pageSize int, orderBy string, isAscending bool) {
	page = cast.ToInt(c.DefaultQuery("page", "1"))
	pageSize = cast.ToInt(c.DefaultQuery("pageSize", "30"))
	orderBy = c.DefaultQuery("orderBy", "created_at")
	asc := c.DefaultQuery("asc", "false")
	isAscending = false
	if asc == "true" {
		isAscending = true
	}
	return
}
