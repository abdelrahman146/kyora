package storefront

import (
	"bytes"
	"io"
	"net/http"

	"github.com/abdelrahman146/kyora/internal/platform/request"
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	service *Service
}

func NewHttpHandler(service *Service) *HttpHandler {
	return &HttpHandler{service: service}
}

// GetCatalog godoc
// @Summary Get storefront catalog
// @Description Returns public business info, categories, products, and variants for a storefront
// @Tags storefront
// @Param storefrontPublicId path string true "Storefront Public ID"
// @Success 200 {object} CatalogResponse
// @Failure 404 {object} problem.Problem
// @Router /v1/storefront/{storefrontPublicId}/catalog [get]
func (h *HttpHandler) GetCatalog(c *gin.Context) {
	storefrontID := c.Param("storefrontPublicId")
	data, err := h.service.GetCatalog(c.Request.Context(), storefrontID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, data)
}

// ListShippingZones godoc
// @Summary List storefront shipping zones
// @Description Returns shipping zones for the storefront (public)
// @Tags storefront
// @Param storefrontPublicId path string true "Storefront Public ID"
// @Success 200 {array} storefront.PublicShippingZone
// @Failure 404 {object} problem.Problem
// @Router /v1/storefront/{storefrontPublicId}/shipping-zones [get]
func (h *HttpHandler) ListShippingZones(c *gin.Context) {
	storefrontID := c.Param("storefrontPublicId")
	data, err := h.service.ListShippingZones(c.Request.Context(), storefrontID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, data)
}

// CreateOrder godoc
// @Summary Create storefront order
// @Description Creates a pending, unpaid order from the public storefront (idempotent via Idempotency-Key)
// @Tags storefront
// @Param storefrontPublicId path string true "Storefront Public ID"
// @Param Idempotency-Key header string true "Idempotency key"
// @Accept json
// @Produce json
// @Success 201 {object} CreateOrderResponse
// @Failure 400 {object} problem.Problem
// @Failure 409 {object} problem.Problem
// @Failure 429 {object} problem.Problem
// @Router /v1/storefront/{storefrontPublicId}/orders [post]
func (h *HttpHandler) CreateOrder(c *gin.Context) {
	storefrontID := c.Param("storefrontPublicId")
	idempotencyKey := c.GetHeader("Idempotency-Key")
	clientIP := c.ClientIP()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.Error(c, err)
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	var req CreateOrderRequest
	if err := request.ValidBody(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	out, err := h.service.CreatePendingOrder(c.Request.Context(), storefrontID, idempotencyKey, body, clientIP, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusCreated, out)
}
