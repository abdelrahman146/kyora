package metadata

import (
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/abdelrahman146/kyora/internal/platform/utils/country"
	"github.com/gin-gonic/gin"
)

type HttpHandler struct{}

func NewHttpHandler() *HttpHandler {
	return &HttpHandler{}
}

func (h *HttpHandler) ListCountries(c *gin.Context) {
	response.SuccessJSON(c, 200, gin.H{
		"countries": country.Countries(),
	})
}
