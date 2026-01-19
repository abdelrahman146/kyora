package order

import (
	"net/http"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/request"
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/abdelrahman146/kyora/internal/platform/types/list"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/abdelrahman146/kyora/internal/platform/utils/transformer"
	"github.com/gin-gonic/gin"
)

// HttpHandler handles HTTP requests for order domain operations.
// Business context is derived from the request path via middleware.
type HttpHandler struct {
	service *Service
}

func NewHttpHandler(service *Service) *HttpHandler {
	return &HttpHandler{service: service}
}

func (h *HttpHandler) getBusinessForRequest(c *gin.Context, actor *account.User) (*business.Business, error) {
	_ = actor
	return business.BusinessFromContext(c)
}

// ListOrders returns a paginated list of orders.
//
// @Summary      List orders
// @Description  Returns a paginated list of orders for the authenticated business
// @Tags         order
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        page query int false "Page number (default: 1)"
// @Param        pageSize query int false "Page size (default: 20, max: 100)"
// @Param        orderBy query []string false "Sort order (e.g., -createdAt, -orderedAt, -total)"
// @Param        search query string false "Search term (matches orderNumber, channel, or customer name/email)"
// @Param        status query []string false "Filter by status (repeatable)"
// @Param        paymentStatus query []string false "Filter by payment status (repeatable)"
// @Param        socialPlatforms query []string false "Filter by platform/channel (instagram, tiktok, facebook, x, snapchat, whatsapp)"
// @Param        customerId query string false "Filter by customerId"
// @Param        orderNumber query string false "Filter by exact orderNumber"
// @Param        from query string false "Filter by orderedAt >= from (RFC3339)"
// @Param        to query string false "Filter by orderedAt <= to (RFC3339)"
// @Success      200 {object} list.ListResponse[order.OrderResponse]
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/orders [get]
// @Security     BearerAuth
func (h *HttpHandler) ListOrders(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var query listOrdersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, problem.BadRequest("invalid query parameters").WithError(err))
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
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

	filters := &ListOrdersFilters{
		Channels:    query.SocialPlatforms,
		CustomerID:  query.CustomerID,
		OrderNumber: query.OrderNumber,
		From:        query.From,
		To:          query.To,
	}
	for _, s := range query.Status {
		filters.Statuses = append(filters.Statuses, OrderStatus(s))
	}
	for _, ps := range query.PaymentStatus {
		filters.PaymentStatuses = append(filters.PaymentStatuses, OrderPaymentStatus(ps))
	}

	items, total, err := h.service.ListOrders(c.Request.Context(), actor, biz, listReq, filters)
	if err != nil {
		response.Error(c, err)
		return
	}
	respItems := ToOrderResponses(items)
	hasMore := int64(query.Page*query.PageSize) < total
	response.SuccessJSON(c, http.StatusOK, list.NewListResponse(respItems, query.Page, query.PageSize, total, hasMore))
}

// GetOrder returns an order by ID with items and notes.
//
// @Summary      Get order
// @Description  Returns an order by ID including items and notes
// @Tags         order
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        orderId path string true "Order ID"
// @Success      200 {object} order.OrderResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/orders/{orderId} [get]
// @Security     BearerAuth
func (h *HttpHandler) GetOrder(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	orderID := c.Param("orderId")
	if orderID == "" {
		response.Error(c, problem.BadRequest("orderId is required"))
		return
	}
	ord, err := h.service.GetOrderByID(c.Request.Context(), actor, biz, orderID)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrOrderNotFound(orderID, err))
			return
		}
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, ToOrderResponse(ord))
}

// GetOrderByNumber returns an order by its order number.
//
// @Summary      Get order by order number
// @Description  Returns an order by its order number (unique per business)
// @Tags         order
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        orderNumber path string true "Order number"
// @Success      200 {object} order.OrderResponse
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/orders/by-number/{orderNumber} [get]
// @Security     BearerAuth
func (h *HttpHandler) GetOrderByNumber(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	orderNumber := c.Param("orderNumber")
	if orderNumber == "" {
		response.Error(c, problem.BadRequest("orderNumber is required"))
		return
	}
	ord, err := h.service.GetOrderByOrderNumber(c.Request.Context(), actor, biz, orderNumber)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, problem.NotFound("order not found"))
			return
		}
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, ToOrderResponse(ord))
}

// PreviewOrder validates an order payload and returns computed totals without creating the order.
//
// @Summary      Preview order totals
// @Description  Validates the order payload and returns computed totals without creating an order or adjusting inventory
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        body body CreateOrderRequest true "Order preview"
// @Success      200 {object} order.OrderPreviewResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      409 {object} problem.Problem
// @Failure      429 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/orders/preview [post]
// @Security     BearerAuth
func (h *HttpHandler) PreviewOrder(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	var req CreateOrderRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}
	if len(req.Items) > 100 {
		response.Error(c, problem.BadRequest("too many order items").With("max", 100))
		return
	}
	preview, err := h.service.PreviewOrder(c.Request.Context(), actor, biz, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, ToOrderPreviewResponse(preview))
}

// CreateOrder creates a new order.
//
// @Summary      Create order
// @Description  Creates an order and allocates inventory
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        body body CreateOrderRequest true "Order"
// @Success      201 {object} order.OrderResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      409 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/orders [post]
// @Security     BearerAuth
func (h *HttpHandler) CreateOrder(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	var req CreateOrderRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}
	if len(req.Items) > 100 {
		response.Error(c, problem.BadRequest("too many order items").With("max", 100))
		return
	}
	ord, err := h.service.CreateOrder(c.Request.Context(), actor, biz, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	loaded, loadErr := h.service.GetOrderByID(c.Request.Context(), actor, biz, ord.ID)
	if loadErr != nil {
		response.Error(c, loadErr)
		return
	}
	response.SuccessJSON(c, http.StatusCreated, ToOrderResponse(loaded))
}

// UpdateOrder updates an order.
//
// @Summary      Update order
// @Description  Updates editable order fields; items updates are restricted by order status
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        orderId path string true "Order ID"
// @Param        body body UpdateOrderRequest true "Order updates"
// @Success      200 {object} order.OrderResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      409 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/orders/{orderId} [patch]
// @Security     BearerAuth
func (h *HttpHandler) UpdateOrder(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	orderID := c.Param("orderId")
	if orderID == "" {
		response.Error(c, problem.BadRequest("orderId is required"))
		return
	}
	var req UpdateOrderRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}
	if len(req.Items) > 100 {
		response.Error(c, problem.BadRequest("too many order items").With("max", 100))
		return
	}
	ord, err := h.service.UpdateOrder(c.Request.Context(), actor, biz, orderID, &req)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrOrderNotFound(orderID, err))
			return
		}
		response.Error(c, err)
		return
	}
	loaded, loadErr := h.service.GetOrderByID(c.Request.Context(), actor, biz, ord.ID)
	if loadErr != nil {
		response.Error(c, loadErr)
		return
	}
	response.SuccessJSON(c, http.StatusOK, ToOrderResponse(loaded))
}

// DeleteOrder deletes an order (restricted to safe statuses) and restocks inventory.
//
// @Summary      Delete order
// @Description  Deletes an order (only allowed for pending/cancelled) and restocks inventory
// @Tags         order
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        orderId path string true "Order ID"
// @Success      204
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      409 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/orders/{orderId} [delete]
// @Security     BearerAuth
func (h *HttpHandler) DeleteOrder(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	orderID := c.Param("orderId")
	if orderID == "" {
		response.Error(c, problem.BadRequest("orderId is required"))
		return
	}
	if err := h.service.DeleteOrder(c.Request.Context(), actor, biz, orderID); err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessEmpty(c, http.StatusNoContent)
}

// UpdateOrderStatus updates order lifecycle status.
//
// @Summary      Update order status
// @Description  Updates the order lifecycle status using the order state machine
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        orderId path string true "Order ID"
// @Param        body body updateOrderStatusRequest true "Status"
// @Success      200 {object} order.OrderResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      409 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/orders/{orderId}/status [patch]
// @Security     BearerAuth
func (h *HttpHandler) UpdateOrderStatus(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	orderID := c.Param("orderId")
	if orderID == "" {
		response.Error(c, problem.BadRequest("orderId is required"))
		return
	}
	var req updateOrderStatusRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}
	ord, err := h.service.UpdateOrderStatus(c.Request.Context(), actor, biz, orderID, req.Status)
	if err != nil {
		response.Error(c, err)
		return
	}
	loaded, loadErr := h.service.GetOrderByID(c.Request.Context(), actor, biz, ord.ID)
	if loadErr != nil {
		response.Error(c, loadErr)
		return
	}
	response.SuccessJSON(c, http.StatusOK, ToOrderResponse(loaded))
}

// UpdateOrderPaymentStatus updates order payment status.
//
// @Summary      Update order payment status
// @Description  Updates the order payment status using the payment state machine
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        orderId path string true "Order ID"
// @Param        body body updateOrderPaymentStatusRequest true "Payment status"
// @Success      200 {object} order.OrderResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      409 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/orders/{orderId}/payment-status [patch]
// @Security     BearerAuth
func (h *HttpHandler) UpdateOrderPaymentStatus(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	orderID := c.Param("orderId")
	if orderID == "" {
		response.Error(c, problem.BadRequest("orderId is required"))
		return
	}
	var req updateOrderPaymentStatusRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}
	ord, err := h.service.UpdateOrderPaymentStatus(c.Request.Context(), actor, biz, orderID, req.PaymentStatus)
	if err != nil {
		response.Error(c, err)
		return
	}
	loaded, loadErr := h.service.GetOrderByID(c.Request.Context(), actor, biz, ord.ID)
	if loadErr != nil {
		response.Error(c, loadErr)
		return
	}
	response.SuccessJSON(c, http.StatusOK, ToOrderResponse(loaded))
}

// AddOrderPaymentDetails sets payment method/reference without changing payment status.
//
// @Summary      Add order payment details
// @Description  Sets payment method/reference fields without changing payment status
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        orderId path string true "Order ID"
// @Param        body body addOrderPaymentDetailsRequest true "Payment details"
// @Success      200 {object} order.OrderResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      409 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/orders/{orderId}/payment-details [patch]
// @Security     BearerAuth
func (h *HttpHandler) AddOrderPaymentDetails(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	orderID := c.Param("orderId")
	if orderID == "" {
		response.Error(c, problem.BadRequest("orderId is required"))
		return
	}
	var body addOrderPaymentDetailsRequest
	if err := request.ValidBody(c, &body); err != nil {
		return
	}
	req := &AddOrderPaymentDetailsRequest{
		PaymentMethod:    body.PaymentMethod,
		PaymentReference: transformer.ToNullString(body.PaymentReference),
	}
	ord, err := h.service.AddOrderPaymentDetails(c.Request.Context(), actor, biz, orderID, req)
	if err != nil {
		response.Error(c, err)
		return
	}
	loaded, loadErr := h.service.GetOrderByID(c.Request.Context(), actor, biz, ord.ID)
	if loadErr != nil {
		response.Error(c, loadErr)
		return
	}
	response.SuccessJSON(c, http.StatusOK, ToOrderResponse(loaded))
}

// CreateOrderNote creates a note for an order.
//
// @Summary      Create order note
// @Description  Creates an order note
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        orderId path string true "Order ID"
// @Param        body body createOrderNoteRequest true "Note"
// @Success      201 {object} order.OrderNoteResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/orders/{orderId}/notes [post]
// @Security     BearerAuth
func (h *HttpHandler) CreateOrderNote(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	orderID := c.Param("orderId")
	if orderID == "" {
		response.Error(c, problem.BadRequest("orderId is required"))
		return
	}
	var body createOrderNoteRequest
	if err := request.ValidBody(c, &body); err != nil {
		return
	}
	if len(body.Content) > 2000 {
		response.Error(c, problem.BadRequest("note content too long").With("max", 2000))
		return
	}
	note, err := h.service.CreateOrderNote(c.Request.Context(), actor, biz, orderID, &CreateOrderNoteRequest{Content: body.Content})
	if err != nil {
		response.Error(c, err)
		return
	}
	noteResponse := ToOrderNoteResponse(note)
	response.SuccessJSON(c, http.StatusCreated, noteResponse)
}

// UpdateOrderNote updates an order note.
//
// @Summary      Update order note
// @Description  Updates an order note
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        orderId path string true "Order ID"
// @Param        noteId path string true "Note ID"
// @Param        body body updateOrderNoteRequest true "Note"
// @Success      200 {object} order.OrderNoteResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/orders/{orderId}/notes/{noteId} [patch]
// @Security     BearerAuth
func (h *HttpHandler) UpdateOrderNote(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	orderID := c.Param("orderId")
	noteID := c.Param("noteId")
	if orderID == "" || noteID == "" {
		response.Error(c, problem.BadRequest("orderId and noteId are required"))
		return
	}
	var body updateOrderNoteRequest
	if err := request.ValidBody(c, &body); err != nil {
		return
	}
	if len(body.Content) > 2000 {
		response.Error(c, problem.BadRequest("note content too long").With("max", 2000))
		return
	}
	note, err := h.service.UpdateOrderNote(c.Request.Context(), actor, biz, orderID, noteID, &UpdateOrderNoteRequest{Content: body.Content})
	if err != nil {
		response.Error(c, err)
		return
	}
	noteResponse := ToOrderNoteResponse(note)
	response.SuccessJSON(c, http.StatusOK, noteResponse)
}

// DeleteOrderNote deletes an order note.
//
// @Summary      Delete order note
// @Description  Deletes an order note
// @Tags         order
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        orderId path string true "Order ID"
// @Param        noteId path string true "Note ID"
// @Success      204
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/orders/{orderId}/notes/{noteId} [delete]
// @Security     BearerAuth
func (h *HttpHandler) DeleteOrderNote(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	orderID := c.Param("orderId")
	noteID := c.Param("noteId")
	if orderID == "" || noteID == "" {
		response.Error(c, problem.BadRequest("orderId and noteId are required"))
		return
	}
	if err := h.service.DeleteOrderNote(c.Request.Context(), actor, biz, orderID, noteID); err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessEmpty(c, http.StatusNoContent)
}
