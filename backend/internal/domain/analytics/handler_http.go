package analytics

import (
	"net/http"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/gin-gonic/gin"
)

// HttpHandler handles HTTP requests for analytics domain operations.
//
// Analytics endpoints must never accept workspaceId from the client.
// Workspace context is derived from the authenticated user via middleware.
type HttpHandler struct {
	service         *Service
	businessService *business.Service
}

// NewHttpHandler creates a new HTTP handler for analytics operations.
func NewHttpHandler(service *Service, businessService *business.Service) *HttpHandler {
	return &HttpHandler{
		service:         service,
		businessService: businessService,
	}
}

// getBusinessForWorkspace returns the first business for the authenticated workspace.
//
// Note: Multi-business support is planned; for now, Kyora assumes a single active business per workspace.
func (h *HttpHandler) getBusinessForWorkspace(c *gin.Context, actor *account.User) (*business.Business, error) {
	businesses, err := h.businessService.ListBusinesses(c.Request.Context(), actor)
	if err != nil {
		return nil, err
	}
	if len(businesses) == 0 {
		return nil, problem.NotFound("no business found for this workspace")
	}
	return businesses[0], nil
}

const dateLayout = "2006-01-02"

func parseDateParam(value string, field string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	t, err := time.Parse(dateLayout, value)
	if err != nil {
		return time.Time{}, problem.BadRequest("invalid " + field + " date format, use YYYY-MM-DD").WithError(err)
	}
	return t, nil
}

func defaultDateRange(from, to time.Time) (time.Time, time.Time) {
	now := time.Now().UTC()
	if to.IsZero() {
		to = now
	}
	if from.IsZero() {
		from = to.AddDate(0, 0, -30)
	}
	return from, to
}

// Dashboard endpoints

// GetDashboardMetrics returns dashboard metrics for the authenticated workspace.
//
// @Summary      Get dashboard metrics
// @Description  Returns dashboard metrics for the authenticated workspace
// @Tags         analytics
// @Produce      json
// @Success      200 {object} analytics.DashboardMetrics
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/analytics/dashboard [get]
// @Security     BearerAuth
func (h *HttpHandler) GetDashboardMetrics(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	metrics, err := h.service.ComputeDashboardMetrics(c.Request.Context(), actor, biz)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessJSON(c, http.StatusOK, metrics)
}

// Sales analytics

type salesAnalyticsQuery struct {
	From string `form:"from" binding:"omitempty"`
	To   string `form:"to" binding:"omitempty"`
}

// GetSalesAnalytics returns sales analytics for the authenticated workspace.
//
// @Summary      Get sales analytics
// @Description  Returns sales analytics for the authenticated workspace
// @Tags         analytics
// @Produce      json
// @Param        from query string false "Start date (YYYY-MM-DD)"
// @Param        to query string false "End date (YYYY-MM-DD)"
// @Success      200 {object} analytics.SalesAnalytics
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/analytics/sales [get]
// @Security     BearerAuth
func (h *HttpHandler) GetSalesAnalytics(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var query salesAnalyticsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, problem.BadRequest("invalid query parameters").WithError(err))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	from, err := parseDateParam(query.From, "from")
	if err != nil {
		response.Error(c, err)
		return
	}
	to, err := parseDateParam(query.To, "to")
	if err != nil {
		response.Error(c, err)
		return
	}
	from, to = defaultDateRange(from, to)
	if to.Before(from) {
		response.Error(c, problem.BadRequest("to must be on or after from").With("from", from.Format(dateLayout)).With("to", to.Format(dateLayout)))
		return
	}

	res, err := h.service.ComputeSalesAnalytics(c.Request.Context(), actor, biz, from, to)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessJSON(c, http.StatusOK, res)
}

// Inventory analytics

type inventoryAnalyticsQuery struct {
	From string `form:"from" binding:"omitempty"`
	To   string `form:"to" binding:"omitempty"`
}

// GetInventoryAnalytics returns inventory analytics for the authenticated workspace.
//
// @Summary      Get inventory analytics
// @Description  Returns inventory analytics for the authenticated workspace
// @Tags         analytics
// @Produce      json
// @Param        from query string false "Start date (YYYY-MM-DD)"
// @Param        to query string false "End date (YYYY-MM-DD)"
// @Success      200 {object} analytics.InventoryAnalytics
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/analytics/inventory [get]
// @Security     BearerAuth
func (h *HttpHandler) GetInventoryAnalytics(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var query inventoryAnalyticsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, problem.BadRequest("invalid query parameters").WithError(err))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	from, err := parseDateParam(query.From, "from")
	if err != nil {
		response.Error(c, err)
		return
	}
	to, err := parseDateParam(query.To, "to")
	if err != nil {
		response.Error(c, err)
		return
	}
	from, to = defaultDateRange(from, to)
	if to.Before(from) {
		response.Error(c, problem.BadRequest("to must be on or after from").With("from", from.Format(dateLayout)).With("to", to.Format(dateLayout)))
		return
	}

	res, err := h.service.ComputeInventoryAnalytics(c.Request.Context(), actor, biz, from, to)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}
	response.SuccessJSON(c, http.StatusOK, res)
}

// Customer analytics

type customerAnalyticsQuery struct {
	From string `form:"from" binding:"omitempty"`
	To   string `form:"to" binding:"omitempty"`
}

// GetCustomerAnalytics returns customer analytics for the authenticated workspace.
//
// @Summary      Get customer analytics
// @Description  Returns customer analytics for the authenticated workspace
// @Tags         analytics
// @Produce      json
// @Param        from query string false "Start date (YYYY-MM-DD)"
// @Param        to query string false "End date (YYYY-MM-DD)"
// @Success      200 {object} analytics.CustomerAnalytics
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/analytics/customers [get]
// @Security     BearerAuth
func (h *HttpHandler) GetCustomerAnalytics(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var query customerAnalyticsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, problem.BadRequest("invalid query parameters").WithError(err))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	from, err := parseDateParam(query.From, "from")
	if err != nil {
		response.Error(c, err)
		return
	}
	to, err := parseDateParam(query.To, "to")
	if err != nil {
		response.Error(c, err)
		return
	}
	from, to = defaultDateRange(from, to)
	if to.Before(from) {
		response.Error(c, problem.BadRequest("to must be on or after from").With("from", from.Format(dateLayout)).With("to", to.Format(dateLayout)))
		return
	}

	res, err := h.service.ComputeCustomerAnalytics(c.Request.Context(), actor, biz, from, to)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessJSON(c, http.StatusOK, res)
}

// Financial reports

type reportAsOfQuery struct {
	AsOf string `form:"asOf" binding:"omitempty"`
}

// GetFinancialPosition returns a balance-sheet style snapshot for the authenticated workspace.
//
// @Summary      Get financial position
// @Description  Returns a balance-sheet style snapshot as of a date for the authenticated workspace
// @Tags         analytics
// @Produce      json
// @Param        asOf query string false "As-of date (YYYY-MM-DD). Default: today"
// @Success      200 {object} analytics.FinancialPosition
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/analytics/reports/financial-position [get]
// @Security     BearerAuth
func (h *HttpHandler) GetFinancialPosition(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var query reportAsOfQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, problem.BadRequest("invalid query parameters").WithError(err))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	asOf, err := parseDateParam(query.AsOf, "asOf")
	if err != nil {
		response.Error(c, err)
		return
	}
	if asOf.IsZero() {
		asOf = time.Now().UTC()
	}

	res, err := h.service.ComputeFinancialPosition(c.Request.Context(), actor, biz, asOf)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}
	response.SuccessJSON(c, http.StatusOK, res)
}

// GetProfitAndLoss returns a profit and loss statement for the authenticated workspace.
//
// @Summary      Get profit and loss
// @Description  Returns a profit and loss statement as of a date for the authenticated workspace
// @Tags         analytics
// @Produce      json
// @Param        asOf query string false "As-of date (YYYY-MM-DD). Default: today"
// @Success      200 {object} analytics.ProfitAndLossStatement
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/analytics/reports/profit-and-loss [get]
// @Security     BearerAuth
func (h *HttpHandler) GetProfitAndLoss(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var query reportAsOfQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, problem.BadRequest("invalid query parameters").WithError(err))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	asOf, err := parseDateParam(query.AsOf, "asOf")
	if err != nil {
		response.Error(c, err)
		return
	}
	if asOf.IsZero() {
		asOf = time.Now().UTC()
	}

	res, err := h.service.ComputeProfitAndLoss(c.Request.Context(), actor, biz, asOf)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}
	response.SuccessJSON(c, http.StatusOK, res)
}

// GetCashFlow returns a cash flow statement for the authenticated workspace.
//
// @Summary      Get cash flow
// @Description  Returns a cash flow statement as of a date for the authenticated workspace
// @Tags         analytics
// @Produce      json
// @Param        asOf query string false "As-of date (YYYY-MM-DD). Default: today"
// @Success      200 {object} analytics.CashFlowStatement
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/analytics/reports/cash-flow [get]
// @Security     BearerAuth
func (h *HttpHandler) GetCashFlow(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var query reportAsOfQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, problem.BadRequest("invalid query parameters").WithError(err))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	asOf, err := parseDateParam(query.AsOf, "asOf")
	if err != nil {
		response.Error(c, err)
		return
	}
	if asOf.IsZero() {
		asOf = time.Now().UTC()
	}

	res, err := h.service.ComputeCashFlow(c.Request.Context(), actor, biz, asOf)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}
	response.SuccessJSON(c, http.StatusOK, res)
}
