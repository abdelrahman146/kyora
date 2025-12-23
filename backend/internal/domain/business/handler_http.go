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

// GetBusiness returns a business by ID (scoped to workspace).
//
// @Summary      Get business
// @Description  Returns a business by ID for the authenticated workspace
// @Tags         business
// @Produce      json
// @Param        businessId path string true "Business ID"
// @Success      200 {object} map[string]business.businessResponse
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessId} [get]
// @Security     BearerAuth
func (h *HttpHandler) GetBusiness(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	businessID := strings.TrimSpace(c.Param("businessId"))
	if businessID == "" {
		response.Error(c, problem.BadRequest("businessId is required"))
		return
	}

	biz, err := h.svc.GetBusinessByID(c.Request.Context(), actor, businessID)
	if err != nil {
		response.Error(c, err)
		return
	}
	if biz == nil {
		response.Error(c, ErrBusinessNotFound(businessID, nil))
		return
	}
	response.SuccessJSON(c, http.StatusOK, gin.H{"business": toBusinessResponse(biz)})
}

// GetBusinessByDescriptor returns a business by descriptor (scoped to workspace).
//
// @Summary      Get business by descriptor
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
// @Router       /v1/businesses/descriptor/{businessDescriptor} [get]
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
		response.Error(c, problem.NotFound("business not found"))
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

// UpdateBusiness updates a business (scoped to workspace).
//
// @Summary      Update business
// @Description  Updates a business by ID in the authenticated workspace
// @Tags         business
// @Accept       json
// @Produce      json
// @Param        businessId path string true "Business ID"
// @Param        request body business.UpdateBusinessInput true "Update business"
// @Success      200 {object} map[string]business.businessResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      409 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessId} [patch]
// @Security     BearerAuth
func (h *HttpHandler) UpdateBusiness(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	businessID := strings.TrimSpace(c.Param("businessId"))
	if businessID == "" {
		response.Error(c, problem.BadRequest("businessId is required"))
		return
	}

	var input UpdateBusinessInput
	if err := request.ValidBody(c, &input); err != nil {
		return
	}

	biz, err := h.svc.UpdateBusiness(c.Request.Context(), actor, businessID, &input)
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
		response.Error(c, ErrBusinessNotFound(businessID, nil))
		return
	}
	response.SuccessJSON(c, http.StatusOK, gin.H{"business": toBusinessResponse(biz)})
}

// ArchiveBusiness marks a business as archived.
//
// @Summary      Archive business
// @Description  Archives a business by ID in the authenticated workspace
// @Tags         business
// @Produce      json
// @Param        businessId path string true "Business ID"
// @Success      204
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessId}/archive [post]
// @Security     BearerAuth
func (h *HttpHandler) ArchiveBusiness(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	businessID := strings.TrimSpace(c.Param("businessId"))
	if businessID == "" {
		response.Error(c, problem.BadRequest("businessId is required"))
		return
	}
	if err := h.svc.ArchiveBusiness(c.Request.Context(), actor, businessID); err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessEmpty(c, http.StatusNoContent)
}

// UnarchiveBusiness removes archive mark for a business.
//
// @Summary      Unarchive business
// @Description  Unarchives a business by ID in the authenticated workspace
// @Tags         business
// @Produce      json
// @Param        businessId path string true "Business ID"
// @Success      204
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessId}/unarchive [post]
// @Security     BearerAuth
func (h *HttpHandler) UnarchiveBusiness(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	businessID := strings.TrimSpace(c.Param("businessId"))
	if businessID == "" {
		response.Error(c, problem.BadRequest("businessId is required"))
		return
	}
	if err := h.svc.UnarchiveBusiness(c.Request.Context(), actor, businessID); err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessEmpty(c, http.StatusNoContent)
}

// DeleteBusiness deletes a business by ID (scoped to workspace).
//
// @Summary      Delete business
// @Description  Deletes a business by ID in the authenticated workspace
// @Tags         business
// @Produce      json
// @Param        businessId path string true "Business ID"
// @Success      204
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessId} [delete]
// @Security     BearerAuth
func (h *HttpHandler) DeleteBusiness(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	businessID := strings.TrimSpace(c.Param("businessId"))
	if businessID == "" {
		response.Error(c, problem.BadRequest("businessId is required"))
		return
	}
	if err := h.svc.DeleteBusiness(c.Request.Context(), actor, businessID); err != nil {
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
