package business

import (
	"net/http"
	"strings"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/request"
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/gin-gonic/gin"
)

// HttpHandler handles HTTP requests for business domain operations.
//
// Business endpoints never accept workspaceId from the client.
// The workspace context is derived from the authenticated actor.
type HttpHandler struct {
	svc *Service
}

func NewHttpHandler(svc *Service) *HttpHandler {
	return &HttpHandler{svc: svc}
}

// ListBusinesses returns all businesses for the authenticated workspace.
//
// @Summary      List businesses
// @Description  Returns businesses for the authenticated workspace
// @Tags         business
// @Produce      json
// @Success      200 {object} map[string][]business.BusinessResponse
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses [get]
// @Security     BearerAuth
func (h *HttpHandler) ListBusinesses(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	items, err := h.svc.ListBusinesses(c.Request.Context(), actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	resp := make([]BusinessResponse, 0, len(items))
	for _, b := range items {
		resp = append(resp, ToBusinessResponse(b))
	}
	response.SuccessJSON(c, http.StatusOK, gin.H{"businesses": resp})
}

// GetBusinessByDescriptor returns a business by descriptor (scoped to workspace).
//
// @Summary      Get business
// @Description  Returns a business by descriptor for the authenticated workspace
// @Tags         business
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Success      200 {object} map[string]business.BusinessResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor} [get]
// @Security     BearerAuth
func (h *HttpHandler) GetBusinessByDescriptor(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	descriptor := strings.TrimSpace(c.Param("businessDescriptor"))
	if descriptor == "" {
		response.Error(c, problem.BadRequest("businessDescriptor is required"))
		return
	}

	biz, err := h.svc.GetBusinessByDescriptor(c.Request.Context(), actor, descriptor)
	if err != nil {
		response.Error(c, err)
		return
	}
	if biz == nil {
		response.Error(c, ErrBusinessNotFound(descriptor, nil))
		return
	}
	response.SuccessJSON(c, http.StatusOK, gin.H{"business": ToBusinessResponse(biz)})
}

// CreateBusiness creates a business within the authenticated workspace.
//
// @Summary      Create business
// @Description  Creates a business in the authenticated workspace
// @Tags         business
// @Accept       json
// @Produce      json
// @Param        request body business.CreateBusinessInput true "Create business"
// @Success      201 {object} map[string]business.BusinessResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      409 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses [post]
// @Security     BearerAuth
func (h *HttpHandler) CreateBusiness(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var input CreateBusinessInput
	if err := request.ValidBody(c, &input); err != nil {
		return
	}

	biz, err := h.svc.CreateBusiness(c.Request.Context(), actor, &input)
	if err != nil {
		// Return a more specific conflict error for descriptor collisions.
		if database.IsUniqueViolation(err) {
			response.Error(c, ErrBusinessDescriptorAlreadyTaken(strings.TrimSpace(strings.ToLower(input.Descriptor)), err))
			return
		}
		response.Error(c, err)
		return
	}

	response.SuccessJSON(c, http.StatusCreated, gin.H{"business": ToBusinessResponse(biz)})
}

// ListShippingZones returns all shipping zones for a business.
//
// @Summary      List shipping zones
// @Description  Returns all shipping zones for the authenticated business
// @Tags         business
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Success      200 {array} business.ShippingZoneResponse
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/shipping-zones [get]
// @Security     BearerAuth
func (h *HttpHandler) ListShippingZones(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := BusinessFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	items, err := h.svc.ListShippingZones(c.Request.Context(), actor, biz)
	if err != nil {
		response.Error(c, err)
		return
	}
	resp := make([]ShippingZoneResponse, 0, len(items))
	for i := range items {
		resp = append(resp, ToShippingZoneResponse(items[i]))
	}
	response.SuccessJSON(c, http.StatusOK, resp)
}

// ListPaymentMethods returns all payment methods (global catalog + business overrides).
//
// @Summary      List payment methods
// @Description  Returns business payment methods configuration (enabled + fees) based on a global catalog with per-business overrides
// @Tags         business
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Success      200 {array} business.PaymentMethodResponse
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/payment-methods [get]
// @Security     BearerAuth
func (h *HttpHandler) ListPaymentMethods(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := BusinessFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	items, err := h.svc.ListPaymentMethods(c.Request.Context(), actor, biz)
	if err != nil {
		response.Error(c, err)
		return
	}
	resp := make([]PaymentMethodResponse, 0, len(items))
	for _, it := range items {
		resp = append(resp, ToPaymentMethodResponse(it))
	}
	response.SuccessJSON(c, http.StatusOK, resp)
}

// UpdatePaymentMethod updates a business payment method (enabled + fees override).
//
// @Summary      Update payment method
// @Description  Enables/disables a payment method for a business and updates its fee configuration
// @Tags         business
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        descriptor path string true "Payment method descriptor"
// @Param        request body business.UpdateBusinessPaymentMethodRequest true "Payment method update"
// @Success      200 {object} business.PaymentMethodResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/payment-methods/{descriptor} [patch]
// @Security     BearerAuth
func (h *HttpHandler) UpdatePaymentMethod(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := BusinessFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	descriptor := strings.TrimSpace(c.Param("descriptor"))
	if descriptor == "" {
		response.Error(c, problem.BadRequest("descriptor is required"))
		return
	}

	var req UpdateBusinessPaymentMethodRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	updated, err := h.svc.UpdatePaymentMethod(c.Request.Context(), actor, biz, descriptor, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, ToPaymentMethodResponse(*updated))
}

// GetShippingZone returns a shipping zone by ID.
//
// @Summary      Get shipping zone
// @Description  Returns a shipping zone by ID
// @Tags         business
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        zoneId path string true "Shipping zone ID"
// @Success      200 {object} business.ShippingZoneResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/shipping-zones/{zoneId} [get]
// @Security     BearerAuth
func (h *HttpHandler) GetShippingZone(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := BusinessFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	zoneID := strings.TrimSpace(c.Param("zoneId"))
	if zoneID == "" {
		response.Error(c, problem.BadRequest("zoneId is required"))
		return
	}
	z, err := h.svc.GetShippingZoneByID(c.Request.Context(), actor, biz, zoneID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, ToShippingZoneResponse(z))
}

// CreateShippingZone creates a shipping zone.
//
// @Summary      Create shipping zone
// @Description  Creates a shipping zone for the business
// @Tags         business
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        request body business.CreateShippingZoneRequest true "Create shipping zone"
// @Success      201 {object} business.ShippingZoneResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      409 {object} problem.Problem
// @Failure      429 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/shipping-zones [post]
// @Security     BearerAuth
func (h *HttpHandler) CreateShippingZone(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := BusinessFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var req CreateShippingZoneRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}
	z, err := h.svc.CreateShippingZone(c.Request.Context(), actor, biz, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusCreated, ToShippingZoneResponse(z))
}

// UpdateShippingZone updates a shipping zone.
//
// @Summary      Update shipping zone
// @Description  Updates a shipping zone for the business
// @Tags         business
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        zoneId path string true "Shipping zone ID"
// @Param        request body business.UpdateShippingZoneRequest true "Update shipping zone"
// @Success      200 {object} business.ShippingZoneResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      409 {object} problem.Problem
// @Failure      429 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/shipping-zones/{zoneId} [patch]
// @Security     BearerAuth
func (h *HttpHandler) UpdateShippingZone(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := BusinessFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	zoneID := strings.TrimSpace(c.Param("zoneId"))
	if zoneID == "" {
		response.Error(c, problem.BadRequest("zoneId is required"))
		return
	}
	var req UpdateShippingZoneRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}
	z, err := h.svc.UpdateShippingZone(c.Request.Context(), actor, biz, zoneID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, ToShippingZoneResponse(z))
}

// DeleteShippingZone deletes a shipping zone.
//
// @Summary      Delete shipping zone
// @Description  Deletes a shipping zone
// @Tags         business
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        zoneId path string true "Shipping zone ID"
// @Success      204
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      429 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/shipping-zones/{zoneId} [delete]
// @Security     BearerAuth
func (h *HttpHandler) DeleteShippingZone(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := BusinessFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	zoneID := strings.TrimSpace(c.Param("zoneId"))
	if zoneID == "" {
		response.Error(c, problem.BadRequest("zoneId is required"))
		return
	}
	if err := h.svc.DeleteShippingZone(c.Request.Context(), actor, biz, zoneID); err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessEmpty(c, http.StatusNoContent)
}

// UpdateBusiness updates a business (scoped to workspace).
//
// @Summary      Update business
// @Description  Updates a business by descriptor in the authenticated workspace
// @Tags         business
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        request body business.UpdateBusinessInput true "Update business"
// @Success      200 {object} map[string]business.BusinessResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      409 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor} [patch]
// @Security     BearerAuth
func (h *HttpHandler) UpdateBusiness(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	descriptor := strings.TrimSpace(c.Param("businessDescriptor"))
	if descriptor == "" {
		response.Error(c, problem.BadRequest("businessDescriptor is required"))
		return
	}

	var input UpdateBusinessInput
	if err := request.ValidBody(c, &input); err != nil {
		return
	}

	ctx := c.Request.Context()
	current, err := h.svc.GetBusinessByDescriptor(ctx, actor, descriptor)
	if err != nil {
		response.Error(c, err)
		return
	}
	if current == nil {
		response.Error(c, ErrBusinessNotFound(descriptor, nil))
		return
	}

	biz, err := h.svc.UpdateBusiness(ctx, actor, current.ID, &input)
	if err != nil {
		if database.IsUniqueViolation(err) {
			// descriptor uniqueness
			var descriptor string
			if input.Descriptor != nil {
				descriptor = *input.Descriptor
			}
			response.Error(c, ErrBusinessDescriptorAlreadyTaken(strings.TrimSpace(strings.ToLower(descriptor)), err))
			return
		}
		response.Error(c, err)
		return
	}
	if biz == nil {
		response.Error(c, ErrBusinessNotFound(descriptor, nil))
		return
	}
	response.SuccessJSON(c, http.StatusOK, gin.H{"business": ToBusinessResponse(biz)})
}

// ArchiveBusiness marks a business as archived.
//
// @Summary      Archive business
// @Description  Archives a business by descriptor in the authenticated workspace
// @Tags         business
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Success      204
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/archive [post]
// @Security     BearerAuth
func (h *HttpHandler) ArchiveBusiness(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	descriptor := strings.TrimSpace(c.Param("businessDescriptor"))
	if descriptor == "" {
		response.Error(c, problem.BadRequest("businessDescriptor is required"))
		return
	}

	ctx := c.Request.Context()
	current, err := h.svc.GetBusinessByDescriptor(ctx, actor, descriptor)
	if err != nil {
		response.Error(c, err)
		return
	}
	if current == nil {
		response.Error(c, ErrBusinessNotFound(descriptor, nil))
		return
	}

	if err := h.svc.ArchiveBusiness(ctx, actor, current.ID); err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessEmpty(c, http.StatusNoContent)
}

// UnarchiveBusiness removes archive mark for a business.
//
// @Summary      Unarchive business
// @Description  Unarchives a business by descriptor in the authenticated workspace
// @Tags         business
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Success      204
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/unarchive [post]
// @Security     BearerAuth
func (h *HttpHandler) UnarchiveBusiness(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	descriptor := strings.TrimSpace(c.Param("businessDescriptor"))
	if descriptor == "" {
		response.Error(c, problem.BadRequest("businessDescriptor is required"))
		return
	}

	ctx := c.Request.Context()
	current, err := h.svc.GetBusinessByDescriptor(ctx, actor, descriptor)
	if err != nil {
		response.Error(c, err)
		return
	}
	if current == nil {
		response.Error(c, ErrBusinessNotFound(descriptor, nil))
		return
	}

	if err := h.svc.UnarchiveBusiness(ctx, actor, current.ID); err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessEmpty(c, http.StatusNoContent)
}

// DeleteBusiness deletes a business by ID (scoped to workspace).
//
// @Summary      Delete business
// @Description  Deletes a business by descriptor in the authenticated workspace
// @Tags         business
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Success      204
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor} [delete]
// @Security     BearerAuth
func (h *HttpHandler) DeleteBusiness(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	descriptor := strings.TrimSpace(c.Param("businessDescriptor"))
	if descriptor == "" {
		response.Error(c, problem.BadRequest("businessDescriptor is required"))
		return
	}

	ctx := c.Request.Context()
	current, err := h.svc.GetBusinessByDescriptor(ctx, actor, descriptor)
	if err != nil {
		response.Error(c, err)
		return
	}
	if current == nil {
		response.Error(c, ErrBusinessNotFound(descriptor, nil))
		return
	}

	if err := h.svc.DeleteBusiness(ctx, actor, current.ID); err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessEmpty(c, http.StatusNoContent)
}

type descriptorAvailabilityQuery struct {
	Descriptor string `form:"descriptor" binding:"required"`
}

// CheckDescriptorAvailability checks whether a descriptor is available in the authenticated workspace.
//
// @Summary      Check business descriptor availability
// @Description  Returns whether a business descriptor is available in the authenticated workspace
// @Tags         business
// @Produce      json
// @Param        descriptor query string true "Business descriptor"
// @Success      200 {object} map[string]bool
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/descriptor/availability [get]
// @Security     BearerAuth
func (h *HttpHandler) CheckDescriptorAvailability(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var q descriptorAvailabilityQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Error(c, problem.BadRequest("invalid query parameters").WithError(err))
		return
	}
	available, err := h.svc.IsBusinessDescriptorAvailable(c.Request.Context(), actor, q.Descriptor)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, gin.H{"available": available})
}
