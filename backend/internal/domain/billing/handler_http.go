package billing

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/request"
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/abdelrahman146/kyora/internal/platform/types/list"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/gin-gonic/gin"
)

// HTTP Handler provides HTTP endpoints for billing operations
type HttpHandler struct {
	service    *Service
	accountSvc *account.Service
}

// NewHttpHandler configures the HTTP routes for billing operations
func NewHttpHandler(service *Service, accountSvc *account.Service) *HttpHandler {
	handler := &HttpHandler{
		service:    service,
		accountSvc: accountSvc,
	}
	return handler
}

// Plan Operations

// ListPlans returns all available billing plans.
//
// @Summary      List billing plans
// @Tags         billing
// @Produce      json
// @Success      200 {array} PlanResponse
// @Failure      500 {object} problem.Problem
// @Router       /v1/billing/plans [get]
func (h *HttpHandler) ListPlans(c *gin.Context) {
	plans, err := h.service.ListPlans(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, ToPlanResponses(plans))
}

// GetPlan returns a plan by descriptor.
//
// @Summary      Get billing plan
// @Tags         billing
// @Produce      json
// @Param        descriptor path string true "Plan descriptor"
// @Success      200 {object} PlanResponse
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/billing/plans/{descriptor} [get]
func (h *HttpHandler) GetPlan(c *gin.Context) {
	descriptor := c.Param("descriptor")
	plan, err := h.service.GetPlanByDescriptor(c.Request.Context(), descriptor)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, ToPlanResponse(plan))
}

// Subscription Operations

// GetSubscription returns the workspace subscription.
//
// @Summary      Get current subscription
// @Tags         billing
// @Produce      json
// @Success      200 {object} SubscriptionResponse
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Router       /v1/billing/subscription [get]
func (h *HttpHandler) GetSubscription(c *gin.Context) {
	ws, err := account.WorkspaceFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	subscription, err := h.service.GetSubscriptionByWorkspaceID(c.Request.Context(), ws.ID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, ToSubscriptionResponse(subscription))
}

// CreateSubscription creates or updates the workspace subscription.
//
// @Summary      Create or update subscription
// @Tags         billing
// @Accept       json
// @Produce      json
// @Param        request body subRequest true "Subscription request"
// @Success      200 {object} SubscriptionResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/billing/subscription [post]
func (h *HttpHandler) CreateSubscription(c *gin.Context) {
	var req subRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	ws, err := account.WorkspaceFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	plan, err := h.service.GetPlanByDescriptor(c.Request.Context(), req.PlanDescriptor)
	if err != nil {
		response.Error(c, err)
		return
	}

	subscription, err := h.service.CreateOrUpdateSubscription(c.Request.Context(), ws, plan)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, ToSubscriptionResponse(subscription))
}

// CancelSubscription cancels the workspace subscription immediately.
//
// @Summary      Cancel subscription
// @Tags         billing
// @Success      204
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/billing/subscription [delete]
func (h *HttpHandler) CancelSubscription(c *gin.Context) {
	ws, err := account.WorkspaceFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	if err := h.service.CancelSubscriptionImmediately(c.Request.Context(), ws); err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessEmpty(c, http.StatusNoContent)
}

// GetSubscriptionDetails returns subscription + plan + default payment method.
//
// @Summary      Get subscription details
// @Tags         billing
// @Produce      json
// @Success      200 {object} SubscriptionDetails
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/billing/subscription/details [get]
func (h *HttpHandler) GetSubscriptionDetails(c *gin.Context) {
	ws, err := account.WorkspaceFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	details, err := h.service.GetSubscriptionDetails(c.Request.Context(), ws)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, details)
}

// ResumeSubscription resumes a canceled/past_due/unpaid subscription when possible.
//
// @Summary      Resume subscription
// @Tags         billing
// @Produce      json
// @Success      200 {object} SubscriptionResponse
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/billing/subscription/resume [post]
func (h *HttpHandler) ResumeSubscription(c *gin.Context) {
	ws, err := account.WorkspaceFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	sub, err := h.service.ResumeSubscriptionIfNoDue(c.Request.Context(), ws)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, ToSubscriptionResponse(sub))
}

// ScheduleSubscriptionChange schedules a plan change.
//
// @Summary      Schedule subscription change
// @Tags         billing
// @Accept       json
// @Produce      json
// @Param        request body scheduleChangeRequest true "Schedule change"
// @Success      200 {object} stripe.SubscriptionSchedule
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/billing/subscription/schedule-change [post]
func (h *HttpHandler) ScheduleSubscriptionChange(c *gin.Context) {
	var req scheduleChangeRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}
	ws, err := account.WorkspaceFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	plan, err := h.service.GetPlanByDescriptor(c.Request.Context(), req.PlanDescriptor)
	if err != nil {
		response.Error(c, err)
		return
	}
	schedule, err := h.service.ScheduleSubscriptionChange(c.Request.Context(), ws, plan, req.EffectiveDate, req.ProrationMode)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, schedule)
}

// EstimateProration estimates proration amount when changing plans.
//
// @Summary      Estimate proration
// @Tags         billing
// @Accept       json
// @Produce      json
// @Param        request body prorationEstimateRequest true "Proration request"
// @Success      200 {object} map[string]int64
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/billing/subscription/estimate-proration [post]
func (h *HttpHandler) EstimateProration(c *gin.Context) {
	var req prorationEstimateRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}
	ws, err := account.WorkspaceFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	amount, err := h.service.EstimateProrationAmount(c.Request.Context(), ws, req.NewPlanDescriptor)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, gin.H{"amount": amount})
}

// Payment Method Operations

// AttachPaymentMethod attaches and sets a default payment method.
//
// @Summary      Attach payment method
// @Tags         billing
// @Accept       json
// @Produce      json
// @Param        request body attachPMRequest true "Attach payment method"
// @Success      200
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/billing/payment-methods/attach [post]
func (h *HttpHandler) AttachPaymentMethod(c *gin.Context) {
	var req attachPMRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	ws, err := account.WorkspaceFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	if err := h.service.AttachAndSetDefaultPaymentMethod(c.Request.Context(), ws, req.PaymentMethodID); err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessEmpty(c, http.StatusOK)
}

// CreateSetupIntent returns a Stripe SetupIntent client secret.
//
// @Summary      Create setup intent
// @Tags         billing
// @Produce      json
// @Success      200 {object} map[string]string
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/billing/payment-methods/setup-intent [post]
func (h *HttpHandler) CreateSetupIntent(c *gin.Context) {
	ws, err := account.WorkspaceFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	secret, err := h.service.CreateSetupIntent(c.Request.Context(), ws)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, gin.H{"clientSecret": secret})
}

// Invoice Operations

// ListInvoices lists Stripe invoices for the workspace.
//
// @Summary      List invoices
// @Tags         billing
// @Produce      json
// @Param        page query int false "Page" default(1)
// @Param        pageSize query int false "Page size" default(30)
// @Param        orderBy query string false "Order by (comma separated)"
// @Param        status query string false "Invoice status"
// @Success      200 {object} list.ListResponse[billing.InvoiceSummary]
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/billing/invoices [get]
func (h *HttpHandler) ListInvoices(c *gin.Context) {
	ws, err := account.WorkspaceFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	// Parse pagination inputs
	page := 1
	pageSize := 30
	if v := c.Query("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			page = p
		}
	}
	if v := c.Query("pageSize"); v != "" {
		if ps, err := strconv.Atoi(v); err == nil && ps > 0 {
			pageSize = ps
		}
	}
	var orderBy []string
	if v := c.Query("orderBy"); v != "" {
		orderBy = strings.Split(v, ",")
	}
	status := c.Query("status")
	req := list.NewListRequest(page, pageSize, orderBy, "")
	resp := h.service.ListInvoices(c.Request.Context(), ws, status, req)
	response.SuccessJSON(c, http.StatusOK, resp)
}

// DownloadInvoice redirects to an invoice PDF/url if invoice belongs to workspace.
//
// @Summary      Download invoice
// @Tags         billing
// @Param        id path string true "Invoice ID"
// @Success      302
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/billing/invoices/{id}/download [get]
func (h *HttpHandler) DownloadInvoice(c *gin.Context) {
	ws, err := account.WorkspaceFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	url, err := h.service.DownloadInvoiceURL(c.Request.Context(), ws, c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	c.Redirect(http.StatusFound, url)
}

// PayInvoice attempts to pay an open invoice.
//
// @Summary      Pay invoice
// @Tags         billing
// @Param        id path string true "Invoice ID"
// @Success      204
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/billing/invoices/{id}/pay [post]
func (h *HttpHandler) PayInvoice(c *gin.Context) {
	ws, err := account.WorkspaceFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	if err := h.service.PayInvoice(c.Request.Context(), ws, c.Param("id")); err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessEmpty(c, http.StatusNoContent)
}

// CreateInvoice creates a manual invoice for the workspace.
//
// @Summary      Create invoice
// @Tags         billing
// @Accept       json
// @Produce      json
// @Param        request body manualInvoiceRequest true "Invoice request"
// @Success      201
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/billing/invoices [post]
func (h *HttpHandler) CreateInvoice(c *gin.Context) {
	var req manualInvoiceRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}
	ws, err := account.WorkspaceFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	inv, err := h.service.CreateInvoice(c.Request.Context(), ws, req.Description, req.Amount, req.Currency, req.DueDate)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusCreated, inv)
}

// Checkout and Portal Operations

// CreateCheckoutSession creates a Stripe Checkout session.
//
// @Summary      Create checkout session
// @Tags         billing
// @Accept       json
// @Produce      json
// @Param        request body checkoutRequest true "Checkout request"
// @Success      200 {object} map[string]string
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/billing/checkout/session [post]
func (h *HttpHandler) CreateCheckoutSession(c *gin.Context) {
	var req checkoutRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	ws, err := account.WorkspaceFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	plan, err := h.service.GetPlanByDescriptor(c.Request.Context(), req.PlanDescriptor)
	if err != nil {
		response.Error(c, err)
		return
	}

	url, err := h.service.CreateCheckoutSession(c.Request.Context(), ws, plan, req.SuccessURL, req.CancelURL)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, gin.H{"url": url, "checkoutUrl": url})
}

// CreateBillingPortalSession creates a Stripe customer portal session.
//
// @Summary      Create billing portal session
// @Tags         billing
// @Accept       json
// @Produce      json
// @Param        request body billingPortalRequest true "Portal request"
// @Success      200 {object} map[string]string
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/billing/portal/session [post]
func (h *HttpHandler) CreateBillingPortalSession(c *gin.Context) {
	var req billingPortalRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	ws, err := account.WorkspaceFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	url, err := h.service.CreateBillingPortalSession(c.Request.Context(), ws, req.ReturnURL)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, gin.H{"url": url, "portalUrl": url})
}

// GetUsage returns current usage (best-effort).
//
// @Summary      Get usage
// @Tags         billing
// @Produce      json
// @Success      200 {object} map[string]int64
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/billing/usage [get]
func (h *HttpHandler) GetUsage(c *gin.Context) {
	ws, err := account.WorkspaceFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	usage, err := h.service.GetSubscriptionUsage(c.Request.Context(), ws)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, usage)
}

// GetUsageQuota returns used and limit for a specific quota type.
//
// @Summary      Get usage quota
// @Tags         billing
// @Produce      json
// @Param        type query string true "Quota type" Enums(orders_per_month,team_members,businesses)
// @Success      200 {object} map[string]any
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/billing/usage/quota [get]
func (h *HttpHandler) GetUsageQuota(c *gin.Context) {
	quotaType := c.Query("type")
	if quotaType == "" {
		response.Error(c, problem.BadRequest("missing required query parameter: type"))
		return
	}

	ws, err := account.WorkspaceFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	used, limit, err := h.service.GetUsageQuota(c.Request.Context(), ws, quotaType)
	if err != nil {
		response.Error(c, problem.BadRequest("invalid quota type").With("type", quotaType).WithError(err))
		return
	}

	response.SuccessJSON(c, http.StatusOK, gin.H{
		"type":  quotaType,
		"used":  used,
		"limit": limit,
	})
}

// CalculateTax calculates tax for the given amount.
//
// @Summary      Calculate tax
// @Tags         billing
// @Accept       json
// @Produce      json
// @Param        request body taxCalculateRequest true "Tax calculation request"
// @Success      200
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/billing/tax/calculate [post]
func (h *HttpHandler) CalculateTax(c *gin.Context) {
	var req taxCalculateRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}
	ws, err := account.WorkspaceFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	calc, err := h.service.CalculateTax(c.Request.Context(), ws, req.Amount, req.Currency)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, calc)
}

// Trial Operations

// ExtendTrial extends the Stripe trial for a trialing subscription.
//
// @Summary      Extend trial
// @Tags         billing
// @Accept       json
// @Success      204
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/billing/subscription/trial/extend [post]
func (h *HttpHandler) ExtendTrial(c *gin.Context) {
	var req trialExtendRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}
	ws, err := account.WorkspaceFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	if err := h.service.ExtendTrialPeriod(c.Request.Context(), ws, req.AdditionalDays); err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessEmpty(c, http.StatusNoContent)
}

// GetTrialStatus returns trial status details for the current workspace.
//
// @Summary      Get trial status
// @Tags         billing
// @Produce      json
// @Success      200 {object} TrialInfo
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/billing/subscription/trial [get]
func (h *HttpHandler) GetTrialStatus(c *gin.Context) {
	ws, err := account.WorkspaceFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	info, err := h.service.CheckTrialStatus(c.Request.Context(), ws)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, info)
}

type gracePeriodRequest struct {
	GraceDays int `json:"graceDays" binding:"required,min=1,max=30"`
}

// SetGracePeriod marks a grace period for past_due subscriptions.
//
// @Summary      Set grace period
// @Tags         billing
// @Accept       json
// @Success      204
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/billing/subscription/grace-period [post]
func (h *HttpHandler) SetGracePeriod(c *gin.Context) {
	var req gracePeriodRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}
	ws, err := account.WorkspaceFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	if err := h.service.HandleGracePeriod(c.Request.Context(), ws, req.GraceDays); err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessEmpty(c, http.StatusNoContent)
}

// Webhook handler

// HandleWebhook receives Stripe webhook events.
//
// @Summary      Stripe webhook
// @Tags         billing
// @Accept       json
// @Success      200
// @Failure      400 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /webhooks/stripe [post]
func (h *HttpHandler) HandleWebhook(c *gin.Context) {
	// This is a public endpoint that doesn't require authentication
	// Stripe will call this directly

	body, err := c.GetRawData()
	if err != nil {
		response.Error(c, problem.BadRequest("invalid webhook payload").WithError(err))
		return
	}

	// Process the webhook in the service layer
	if err := h.service.ProcessWebhook(c.Request.Context(), body, c.GetHeader("Stripe-Signature")); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessEmpty(c, http.StatusOK)
}
