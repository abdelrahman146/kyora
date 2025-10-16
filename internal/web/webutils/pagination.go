package webutils

import (
	"github.com/abdelrahman146/kyora/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func GetPaginationParams(c *gin.Context) *types.ListRequest {
	page := cast.ToInt(c.DefaultQuery("page", "1"))
	pageSize := cast.ToInt(c.DefaultQuery("pageSize", "30"))
	orderBy := c.QueryArray("orderBy")
	search := c.Query("search")
	return &types.ListRequest{
		Page:     page,
		PageSize: pageSize,
		OrderBy:  orderBy,
		Search:   search,
	}
}
