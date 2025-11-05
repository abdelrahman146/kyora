package billing

import (
	"net/http"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/auth"
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/gin-gonic/gin"
)

// HTTP Handler provides HTTP endpoints for billing operations
type HttpHandler struct {
	service    *Service
	accountSvc *account.Service
}

// RegisterRoutes configures the HTTP routes for billing operations
func RegisterRoutes(router *gin.Engine, service *Service, accountSvc *account.Service) {
	handler := &HttpHandler{
		service:    service,
		accountSvc: accountSvc,
	}
	handler.RegisterRoutes(router)
}

func (h *HttpHandler) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/v1/billing")

	// Plan Operations (Public).
	group.GET("/plans", h.ListPlans)
	group.GET("/plans/:descriptor", h.GetPlan)

	// Subscription Operations
	subscriptionGroup := group.Group("/subscription")
	subscriptionGroup.Use(auth.EnforceAuthentication, account.EnforceValidActor(h.accountSvc), account.EnforceWorkspaceMembership(h.accountSvc))
	{
		subscriptionGroup.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceBilling), h.GetSubscription)
		subscriptionGroup.POST("", account.EnforceActorPermissions(role.ActionManage, role.ResourceBilling), h.CreateSubscription)
		subscriptionGroup.DELETE("", account.EnforceActorPermissions(role.ActionManage, role.ResourceBilling), h.CancelSubscription)
	}

	// Payment Method Operations
	paymentGroup := group.Group("/payment-methods")
	paymentGroup.Use(auth.EnforceAuthentication, account.EnforceValidActor(h.accountSvc), account.EnforceWorkspaceMembership(h.accountSvc))
	{
		paymentGroup.POST("/attach", account.EnforceActorPermissions(role.ActionManage, role.ResourceBilling), h.AttachPaymentMethod)
	}

	// Invoice Operations
	invoiceGroup := group.Group("/invoices")
	invoiceGroup.Use(auth.EnforceAuthentication, account.EnforceValidActor(h.accountSvc), account.EnforceWorkspaceMembership(h.accountSvc))
	{
		invoiceGroup.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceBilling), h.ListInvoices)
		invoiceGroup.GET("/:id/download", account.EnforceActorPermissions(role.ActionView, role.ResourceBilling), h.DownloadInvoice)
		invoiceGroup.POST("/:id/pay", account.EnforceActorPermissions(role.ActionManage, role.ResourceBilling), h.PayInvoice)
	}

	// Checkout and Portal Operations
	checkoutGroup := group.Group("/checkout")
	{
		checkoutGroup.POST("/session", h.CreateCheckoutSession)
	}

	portalGroup := group.Group("/portal")
	{
		portalGroup.POST("/session", h.CreateBillingPortalSession)
	}

	// Webhook endpoint (public - no auth required)
	router.POST("/webhooks/stripe", h.HandleWebhook)
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

// Invoice Operations
func (h *HttpHandler) ListInvoices(c *gin.Context) {
	ws, err := account.WorkspaceFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	status := c.Query("status")
	list, err := h.service.ListInvoices(c.Request.Context(), ws, status)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, list)
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
