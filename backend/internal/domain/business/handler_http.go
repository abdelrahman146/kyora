package business

import (
	"net/http"
	"strings"
	"time"

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

type businessResponse struct {
	ID            string     `json:"id"`
	WorkspaceID   string     `json:"workspaceId"`
	Descriptor    string     `json:"descriptor"`
	Name          string     `json:"name"`
	CountryCode   string     `json:"countryCode"`
	Currency      string     `json:"currency"`
	VatRate       string     `json:"vatRate"`
	SafetyBuffer  string     `json:"safetyBuffer"`
	EstablishedAt time.Time  `json:"establishedAt"`
	ArchivedAt    *time.Time `json:"archivedAt,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

type shippingZoneResponse struct {
	ID                    string    `json:"id"`
	BusinessID            string    `json:"businessId"`
	Name                  string    `json:"name"`
	Countries             []string  `json:"countries"`
	Currency              string    `json:"currency"`
	ShippingCost          string    `json:"shippingCost"`
	FreeShippingThreshold string    `json:"freeShippingThreshold"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

func toShippingZoneResponse(z *ShippingZone) shippingZoneResponse {
	resp := shippingZoneResponse{
		ID:                    z.ID,
		BusinessID:            z.BusinessID,
		Name:                  z.Name,
		Countries:             []string(z.Countries),
		Currency:              z.Currency,
		ShippingCost:          z.ShippingCost.String(),
		FreeShippingThreshold: z.FreeShippingThreshold.String(),
		CreatedAt:             z.CreatedAt,
		UpdatedAt:             z.UpdatedAt,
	}
	if resp.Countries == nil {
		resp.Countries = []string{}
	}
	return resp
}

func toBusinessResponse(b *Business) businessResponse {
	return businessResponse{
		ID:            b.ID,
		WorkspaceID:   b.WorkspaceID,
		Descriptor:    b.Descriptor,
		Name:          b.Name,
		CountryCode:   b.CountryCode,
		Currency:      b.Currency,
		VatRate:       b.VatRate.String(),
		SafetyBuffer:  b.SafetyBuffer.StringFixed(2),
		EstablishedAt: b.EstablishedAt,
		ArchivedAt:    b.ArchivedAt,
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
	}
}

// ListBusinesses returns all businesses for the authenticated workspace.
//
// @Summary      List businesses
// @Description  Returns businesses for the authenticated workspace
// @Tags         business
// @Produce      json
// @Success      200 {object} map[string][]business.businessResponse
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

	resp := make([]businessResponse, 0, len(items))
	for _, b := range items {
		resp = append(resp, toBusinessResponse(b))
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
// @Success      200 {object} map[string]business.businessResponse
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
	response.SuccessJSON(c, http.StatusOK, gin.H{"business": toBusinessResponse(biz)})
}

// CreateBusiness creates a business within the authenticated workspace.
//
// @Summary      Create business
// @Description  Creates a business in the authenticated workspace
// @Tags         business
// @Accept       json
// @Produce      json
// @Param        request body business.CreateBusinessInput true "Create business"
// @Success      201 {object} map[string]business.businessResponse
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

	response.SuccessJSON(c, http.StatusCreated, gin.H{"business": toBusinessResponse(biz)})
}

// ListShippingZones returns all shipping zones for a business.
//
// @Summary      List shipping zones
// @Description  Returns all shipping zones for the authenticated business
// @Tags         business
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Success      200 {array} business.shippingZoneResponse
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
	resp := make([]shippingZoneResponse, 0, len(items))
	for _, z := range items {
		resp = append(resp, toShippingZoneResponse(z))
	}
	response.SuccessJSON(c, http.StatusOK, resp)
}

// GetShippingZone returns a shipping zone by ID.
//
// @Summary      Get shipping zone
// @Description  Returns a shipping zone by ID
// @Tags         business
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        zoneId path string true "Shipping zone ID"
// @Success      200 {object} business.shippingZoneResponse
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
	response.SuccessJSON(c, http.StatusOK, toShippingZoneResponse(z))
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
// @Success      201 {object} business.shippingZoneResponse
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
	response.SuccessJSON(c, http.StatusCreated, toShippingZoneResponse(z))
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
// @Success      200 {object} business.shippingZoneResponse
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
	response.SuccessJSON(c, http.StatusOK, toShippingZoneResponse(z))
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
// @Success      200 {object} map[string]business.businessResponse
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
	response.SuccessJSON(c, http.StatusOK, gin.H{"business": toBusinessResponse(biz)})
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
