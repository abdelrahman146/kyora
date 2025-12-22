package billing

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/abdelrahman146/kyora/internal/platform/types/list"
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
func (h *HttpHandler) ListPlans(c *gin.Context) {
	plans, err := h.service.ListPlans(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, plans)
}

func (h *HttpHandler) GetPlan(c *gin.Context) {
	descriptor := c.Param("descriptor")
	plan, err := h.service.GetPlanByDescriptor(c.Request.Context(), descriptor)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, plan)
}

// Subscription Operations
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
	response.SuccessJSON(c, http.StatusOK, subscription)
}

func (h *HttpHandler) CreateSubscription(c *gin.Context) {
	var req subRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
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
	response.SuccessJSON(c, http.StatusOK, subscription)
}

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
	response.SuccessJSON(c, http.StatusOK, sub)
}

func (h *HttpHandler) ScheduleSubscriptionChange(c *gin.Context) {
	var req scheduleChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
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

func (h *HttpHandler) EstimateProration(c *gin.Context) {
	var req prorationEstimateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
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
func (h *HttpHandler) AttachPaymentMethod(c *gin.Context) {
	var req attachPMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
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

func (h *HttpHandler) CreateInvoice(c *gin.Context) {
	var req manualInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
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
func (h *HttpHandler) CreateCheckoutSession(c *gin.Context) {
	var req checkoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
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
	response.SuccessJSON(c, http.StatusOK, gin.H{"checkoutUrl": url})
}

func (h *HttpHandler) CreateBillingPortalSession(c *gin.Context) {
	var req billingPortalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
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
	response.SuccessJSON(c, http.StatusOK, gin.H{"portalUrl": url})
}

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

func (h *HttpHandler) CalculateTax(c *gin.Context) {
	var req taxCalculateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
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
func (h *HttpHandler) ExtendTrial(c *gin.Context) {
	var req trialExtendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
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

// Webhook handler
func (h *HttpHandler) HandleWebhook(c *gin.Context) {
	// This is a public endpoint that doesn't require authentication
	// Stripe will call this directly

	body, err := c.GetRawData()
	if err != nil {
		response.Error(c, err)
		return
	}

	// Process the webhook in the service layer
	if err := h.service.ProcessWebhook(c.Request.Context(), body, c.GetHeader("Stripe-Signature")); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessEmpty(c, http.StatusOK)
}
