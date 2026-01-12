package customer

import (
	"net/http"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/request"
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/abdelrahman146/kyora/internal/platform/types/list"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/gin-gonic/gin"
)

// HttpHandler handles HTTP requests for customer domain operations
type HttpHandler struct {
	service *Service
}

// NewHttpHandler creates a new HTTP handler for customer operations.
// Business context is derived from the request path via middleware.
func NewHttpHandler(service *Service) *HttpHandler {
	return &HttpHandler{service: service}
}

// getBusinessForWorkspace retrieves the first business for a workspace
func (h *HttpHandler) getBusinessForWorkspace(c *gin.Context, actor *account.User) (*business.Business, error) {
	_ = actor
	return business.BusinessFromContext(c)
}

// Customer endpoints

type listCustomersQuery struct {
	Page            int      `form:"page" binding:"omitempty,min=1"`
	PageSize        int      `form:"pageSize" binding:"omitempty,min=1,max=100"`
	OrderBy         []string `form:"orderBy" binding:"omitempty"`
	SearchTerm      string   `form:"search" binding:"omitempty"`
	CountryCode     string   `form:"countryCode" binding:"omitempty"`
	HasOrders       *bool    `form:"hasOrders" binding:"omitempty"`
	SocialPlatforms []string `form:"socialPlatforms" binding:"omitempty"`
}

// ListCustomers returns a paginated list of customers
//
// @Summary      List customers
// @Description  Returns a paginated list of all customers for the authenticated workspace with optional filters
// @Tags         customer
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        page query int false "Page number (default: 1)"
// @Param        pageSize query int false "Page size (default: 20, max: 100)"
// @Param        orderBy query []string false "Sort order (e.g., -name, email)"
// @Param        search query string false "Search term for customer name/email/phone/social handles"
// @Param        countryCode query string false "Filter by country code (e.g., US, AE)"
// @Param        hasOrders query bool false "Filter by customers with or without orders"
// @Param        socialPlatforms query []string false "Filter by social media platforms (instagram, tiktok, facebook, x, snapchat, whatsapp)"
// @Success      200 {object} list.ListResponse[customer.CustomerResponse]
// @Failure      401 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/customers [get]
// @Security     BearerAuth
func (h *HttpHandler) ListCustomers(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var query listCustomersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, problem.BadRequest("invalid query parameters").WithError(err))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}
	if query.SearchTerm != "" {
		term, err := list.NormalizeSearchTerm(query.SearchTerm)
		if err != nil {
			response.Error(c, problem.BadRequest("invalid search term"))
			return
		}
		query.SearchTerm = term
	}

	listReq := list.NewListRequest(query.Page, query.PageSize, query.OrderBy, query.SearchTerm)

	filters := &ListCustomersFilters{
		CountryCode:     query.CountryCode,
		HasOrders:       query.HasOrders,
		SocialPlatforms: query.SocialPlatforms,
	}

	customers, totalCount, err := h.service.ListCustomers(c.Request.Context(), actor, biz, listReq, filters)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	hasMore := int64(query.Page*query.PageSize) < totalCount
	listResp := list.NewListResponse(customers, query.Page, query.PageSize, totalCount, hasMore)
	response.SuccessJSON(c, http.StatusOK, listResp)
}

// GetCustomer returns a specific customer by ID
//
// @Summary      Get customer
// @Description  Returns a specific customer by ID with addresses and notes
// @Tags         customer
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        customerId path string true "Customer ID"
// @Success      200 {object} customer.CustomerResponse
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/customers/{customerId} [get]
// @Security     BearerAuth
func (h *HttpHandler) GetCustomer(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	customerID := c.Param("customerId")
	if customerID == "" {
		response.Error(c, problem.BadRequest("customerId is required"))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	customer, err := h.service.GetCustomerByID(c.Request.Context(), actor, biz, customerID)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrCustomerNotFound(err))
			return
		}
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	// Get customer aggregations (orders count, total spent)
	aggregations, err := h.service.storage.GetCustomerAggregations(c.Request.Context(), biz.ID, []string{customer.ID})
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	var ordersCount int
	var totalSpent float64
	if agg, ok := aggregations[customer.ID]; ok {
		ordersCount = int(agg.OrdersCount)
		totalSpent = agg.TotalSpent
	}

	customerResponse := ToCustomerResponse(customer, ordersCount, totalSpent)
	response.SuccessJSON(c, http.StatusOK, customerResponse)
}

// CreateCustomer creates a new customer
//
// @Summary      Create customer
// @Description  Creates a new customer for the authenticated workspace
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        request body CreateCustomerRequest true "Customer data"
// @Success      201 {object} customer.CustomerResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      409 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/customers [post]
// @Security     BearerAuth
func (h *HttpHandler) CreateCustomer(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var req CreateCustomerRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	customer, err := h.service.CreateCustomer(c.Request.Context(), actor, biz, &req)
	if err != nil {
		if database.IsUniqueViolation(err) {
			response.Error(c, ErrCustomerDuplicateEmail(err))
			return
		}
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	// New customers have no orders yet
	customerResponse := ToCustomerResponse(customer, 0, 0)
	response.SuccessJSON(c, http.StatusCreated, customerResponse)
}

// UpdateCustomer updates an existing customer
//
// @Summary      Update customer
// @Description  Updates an existing customer's information
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        customerId path string true "Customer ID"
// @Param        request body UpdateCustomerRequest true "Updated customer data"
// @Success      200 {object} customer.CustomerResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      409 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/customers/{customerId} [patch]
// @Security     BearerAuth
func (h *HttpHandler) UpdateCustomer(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	customerID := c.Param("customerId")
	if customerID == "" {
		response.Error(c, problem.BadRequest("customerId is required"))
		return
	}

	var req UpdateCustomerRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	customer, err := h.service.UpdateCustomer(c.Request.Context(), actor, biz, customerID, &req)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrCustomerNotFound(err))
			return
		}
		if database.IsUniqueViolation(err) {
			response.Error(c, ErrCustomerDuplicateEmail(err))
			return
		}
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	// Get customer aggregations (orders count, total spent)
	aggregations, err := h.service.storage.GetCustomerAggregations(c.Request.Context(), biz.ID, []string{customer.ID})
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	var ordersCount int
	var totalSpent float64
	if agg, ok := aggregations[customer.ID]; ok {
		ordersCount = int(agg.OrdersCount)
		totalSpent = agg.TotalSpent
	}

	customerResponse := ToCustomerResponse(customer, ordersCount, totalSpent)
	response.SuccessJSON(c, http.StatusOK, customerResponse)
}

// DeleteCustomer soft deletes a customer
//
// @Summary      Delete customer
// @Description  Soft deletes a customer (can be restored later)
// @Tags         customer
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        customerId path string true "Customer ID"
// @Success      204
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/customers/{customerId} [delete]
// @Security     BearerAuth
func (h *HttpHandler) DeleteCustomer(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	customerID := c.Param("customerId")
	if customerID == "" {
		response.Error(c, problem.BadRequest("customerId is required"))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	err = h.service.DeleteCustomer(c.Request.Context(), actor, biz, customerID)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrCustomerNotFound(err))
			return
		}
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessEmpty(c, http.StatusNoContent)
}

// Customer Address endpoints

// ListCustomerAddresses returns all addresses for a customer
//
// @Summary      List customer addresses
// @Description  Returns all addresses for a specific customer
// @Tags         customer
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        customerId path string true "Customer ID"
// @Success      200 {array} customer.CustomerAddressResponse
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/customers/{customerId}/addresses [get]
// @Security     BearerAuth
func (h *HttpHandler) ListCustomerAddresses(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	customerID := c.Param("customerId")
	if customerID == "" {
		response.Error(c, problem.BadRequest("customerId is required"))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	addresses, err := h.service.ListCustomerAddresses(c.Request.Context(), actor, biz, customerID)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrCustomerNotFound(err))
			return
		}
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	// Convert to response types
	addressResponses := make([]CustomerAddressResponse, len(addresses))
	for i, addr := range addresses {
		addressResponses[i] = ToCustomerAddressResponse(addr)
	}

	response.SuccessJSON(c, http.StatusOK, addressResponses)
}

// CreateCustomerAddress creates a new address for a customer
//
// @Summary      Create customer address
// @Description  Creates a new address for a specific customer
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        customerId path string true "Customer ID"
// @Param        request body CreateCustomerAddressRequest true "Address data"
// @Success      201 {object} customer.CustomerAddressResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/customers/{customerId}/addresses [post]
// @Security     BearerAuth
func (h *HttpHandler) CreateCustomerAddress(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	customerID := c.Param("customerId")
	if customerID == "" {
		response.Error(c, problem.BadRequest("customerId is required"))
		return
	}

	var req CreateCustomerAddressRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	address, err := h.service.CreateCustomerAddress(c.Request.Context(), actor, biz, customerID, &req)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrCustomerNotFound(err))
			return
		}
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	addressResponse := ToCustomerAddressResponse(address)
	response.SuccessJSON(c, http.StatusCreated, addressResponse)
}

// UpdateCustomerAddress updates an existing address
//
// @Summary      Update customer address
// @Description  Updates an existing address for a customer
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        customerId path string true "Customer ID"
// @Param        addressId path string true "Address ID"
// @Param        request body UpdateCustomerAddressRequest true "Updated address data"
// @Success      200 {object} customer.CustomerAddressResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/customers/{customerId}/addresses/{addressId} [patch]
// @Security     BearerAuth
func (h *HttpHandler) UpdateCustomerAddress(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	customerID := c.Param("customerId")
	addressID := c.Param("addressId")
	if customerID == "" || addressID == "" {
		response.Error(c, problem.BadRequest("customerId and addressId are required"))
		return
	}

	var req UpdateCustomerAddressRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	address, err := h.service.UpdateCustomerAddress(c.Request.Context(), actor, biz, customerID, addressID, &req)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrCustomerAddressNotFound(err))
			return
		}
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	addressResponse := ToCustomerAddressResponse(address)
	response.SuccessJSON(c, http.StatusOK, addressResponse)
}

// DeleteCustomerAddress deletes an address
//
// @Summary      Delete customer address
// @Description  Deletes an address for a customer
// @Tags         customer
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        customerId path string true "Customer ID"
// @Param        addressId path string true "Address ID"
// @Success      204
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/customers/{customerId}/addresses/{addressId} [delete]
// @Security     BearerAuth
func (h *HttpHandler) DeleteCustomerAddress(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	customerID := c.Param("customerId")
	addressID := c.Param("addressId")
	if customerID == "" || addressID == "" {
		response.Error(c, problem.BadRequest("customerId and addressId are required"))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	err = h.service.DeleteCustomerAddress(c.Request.Context(), actor, biz, customerID, addressID)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrCustomerAddressNotFound(err))
			return
		}
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessEmpty(c, http.StatusNoContent)
}

// Customer Note endpoints

// ListCustomerNotes returns all notes for a customer
//
// @Summary      List customer notes
// @Description  Returns all notes for a specific customer
// @Tags         customer
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        customerId path string true "Customer ID"
// @Success      200 {array} customer.CustomerNoteResponse
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/customers/{customerId}/notes [get]
// @Security     BearerAuth
func (h *HttpHandler) ListCustomerNotes(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	customerID := c.Param("customerId")
	if customerID == "" {
		response.Error(c, problem.BadRequest("customerId is required"))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	notes, err := h.service.ListCustomerNotes(c.Request.Context(), actor, biz, customerID)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrCustomerNotFound(err))
			return
		}
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	// Convert to response types
	noteResponses := make([]CustomerNoteResponse, len(notes))
	for i, note := range notes {
		noteResponses[i] = ToCustomerNoteResponse(note)
	}
	response.SuccessJSON(c, http.StatusOK, noteResponses)
}

// CreateCustomerNote creates a new note for a customer
//
// @Summary      Create customer note
// @Description  Creates a new note for a specific customer
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        customerId path string true "Customer ID"
// @Param        request body CreateCustomerNoteRequest true "Note data"
// @Success      201 {object} customer.CustomerNoteResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/customers/{customerId}/notes [post]
// @Security     BearerAuth
func (h *HttpHandler) CreateCustomerNote(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	customerID := c.Param("customerId")
	if customerID == "" {
		response.Error(c, problem.BadRequest("customerId is required"))
		return
	}

	var req CreateCustomerNoteRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	note, err := h.service.CreateCustomerNote(c.Request.Context(), actor, biz, customerID, req.Content)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrCustomerNotFound(err))
			return
		}
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	noteResponse := ToCustomerNoteResponse(note)
	response.SuccessJSON(c, http.StatusCreated, noteResponse)
}

// DeleteCustomerNote deletes a note
//
// @Summary      Delete customer note
// @Description  Deletes a note for a customer
// @Tags         customer
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        customerId path string true "Customer ID"
// @Param        noteId path string true "Note ID"
// @Success      204
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/customers/{customerId}/notes/{noteId} [delete]
// @Security     BearerAuth
func (h *HttpHandler) DeleteCustomerNote(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	customerID := c.Param("customerId")
	noteID := c.Param("noteId")
	if customerID == "" || noteID == "" {
		response.Error(c, problem.BadRequest("customerId and noteId are required"))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	err = h.service.DeleteCustomerNote(c.Request.Context(), actor, biz, customerID, noteID)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrCustomerNoteNotFound(err))
			return
		}
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessEmpty(c, http.StatusNoContent)
}
