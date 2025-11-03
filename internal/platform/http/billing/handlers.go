package billinghttp

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/billing"
	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/abdelrahman146/kyora/internal/platform/request"
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/stripe/stripe-go/v83/webhook"
)

type attachPMRequest struct {
	PaymentMethodID string `json:"paymentMethodId" binding:"required"`
}

type subRequest struct {
	PlanDescriptor string `json:"planDescriptor" binding:"required"`
}

// Register wires billing endpoints under /api/billing and webhook at /webhooks/stripe (root group must call webhook separately)
func Register(root *gin.Engine, auth *gin.RouterGroup, svc *billing.Service, accountSvc *account.Service) {
	grp := auth.Group("/billing")

	grp.GET("/plans", func(c *gin.Context) {
		plans, err := svc.ListPlans(c.Request.Context())
		if err != nil {
			response.Error(c, err)
			return
		}
		response.SuccessJSON(c, http.StatusOK, plans)
	})

	grp.POST("/payment-method/attach", func(c *gin.Context) {
		var req attachPMRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, err)
			return
		}
		actor, err := request.ActorFromContext(c)
		if err != nil {
			response.Error(c, err)
			return
		}
		ws, err := accountSvc.GetWorkspaceByID(c.Request.Context(), actor.WorkspaceID)
		if err != nil {
			response.Error(c, err)
			return
		}
		if err := svc.AttachAndSetDefaultPaymentMethod(c.Request.Context(), ws, req.PaymentMethodID); err != nil {
			response.Error(c, err)
			return
		}
		response.SuccessEmpty(c, http.StatusNoContent)
	})

	grp.GET("/payment-method/setup-intent", func(c *gin.Context) {
		actor, err := request.ActorFromContext(c)
		if err != nil {
			response.Error(c, err)
			return
		}
		ws, err := accountSvc.GetWorkspaceByID(c.Request.Context(), actor.WorkspaceID)
		if err != nil {
			response.Error(c, err)
			return
		}
		secret, err := svc.CreateSetupIntent(c.Request.Context(), ws)
		if err != nil {
			response.Error(c, err)
			return
		}
		response.SuccessJSON(c, http.StatusOK, gin.H{"clientSecret": secret})
	})

	grp.POST("/subscription", func(c *gin.Context) {
		var req subRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, err)
			return
		}
		actor, err := request.ActorFromContext(c)
		if err != nil {
			response.Error(c, err)
			return
		}
		ws, err := accountSvc.GetWorkspaceByID(c.Request.Context(), actor.WorkspaceID)
		if err != nil {
			response.Error(c, err)
			return
		}
		plan, err := svc.GetPlanByDescriptor(c.Request.Context(), req.PlanDescriptor)
		if err != nil {
			response.Error(c, err)
			return
		}
		subRec, err := svc.CreateOrUpdateSubscription(c.Request.Context(), ws, plan)
		if err != nil {
			response.Error(c, err)
			return
		}
		response.SuccessJSON(c, http.StatusCreated, subRec)
	})

	grp.POST("/subscription/change", func(c *gin.Context) {
		var req subRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, err)
			return
		}
		actor, err := request.ActorFromContext(c)
		if err != nil {
			response.Error(c, err)
			return
		}
		ws, err := accountSvc.GetWorkspaceByID(c.Request.Context(), actor.WorkspaceID)
		if err != nil {
			response.Error(c, err)
			return
		}
		plan, err := svc.GetPlanByDescriptor(c.Request.Context(), req.PlanDescriptor)
		if err != nil {
			response.Error(c, err)
			return
		}
		subRec, err := svc.CreateOrUpdateSubscription(c.Request.Context(), ws, plan)
		if err != nil {
			response.Error(c, err)
			return
		}
		response.SuccessJSON(c, http.StatusOK, subRec)
	})

	grp.POST("/subscription/cancel", func(c *gin.Context) {
		actor, err := request.ActorFromContext(c)
		if err != nil {
			response.Error(c, err)
			return
		}
		ws, err := accountSvc.GetWorkspaceByID(c.Request.Context(), actor.WorkspaceID)
		if err != nil {
			response.Error(c, err)
			return
		}
		if err := svc.CancelSubscriptionImmediately(c.Request.Context(), ws); err != nil {
			response.Error(c, err)
			return
		}
		response.SuccessEmpty(c, http.StatusNoContent)
	})

	grp.POST("/subscription/resume", func(c *gin.Context) {
		actor, err := request.ActorFromContext(c)
		if err != nil {
			response.Error(c, err)
			return
		}
		ws, err := accountSvc.GetWorkspaceByID(c.Request.Context(), actor.WorkspaceID)
		if err != nil {
			response.Error(c, err)
			return
		}
		subRec, err := svc.ResumeSubscriptionIfNoDue(c.Request.Context(), ws)
		if err != nil {
			response.Error(c, err)
			return
		}
		response.SuccessJSON(c, http.StatusOK, subRec)
	})

	// Subscription details (includes default payment method info)
	grp.GET("/subscription", func(c *gin.Context) {
		actor, err := request.ActorFromContext(c)
		if err != nil {
			response.Error(c, err)
			return
		}
		ws, err := accountSvc.GetWorkspaceByID(c.Request.Context(), actor.WorkspaceID)
		if err != nil {
			response.Error(c, err)
			return
		}
		details, err := svc.GetSubscriptionDetails(c.Request.Context(), ws)
		if err != nil {
			response.Error(c, err)
			return
		}
		response.SuccessJSON(c, http.StatusOK, details)
	})

	// List invoices: /billing/invoices?status=paid|open|all
	grp.GET("/invoices", func(c *gin.Context) {
		actor, err := request.ActorFromContext(c)
		if err != nil {
			response.Error(c, err)
			return
		}
		ws, err := accountSvc.GetWorkspaceByID(c.Request.Context(), actor.WorkspaceID)
		if err != nil {
			response.Error(c, err)
			return
		}
		status := c.Query("status")
		list, err := svc.ListInvoices(c.Request.Context(), ws, status)
		if err != nil {
			response.Error(c, err)
			return
		}
		response.SuccessJSON(c, http.StatusOK, list)
	})

	// Download invoice: redirect to Stripe hosted PDF
	grp.GET("/invoices/:id/download", func(c *gin.Context) {
		actor, err := request.ActorFromContext(c)
		if err != nil {
			response.Error(c, err)
			return
		}
		ws, err := accountSvc.GetWorkspaceByID(c.Request.Context(), actor.WorkspaceID)
		if err != nil {
			response.Error(c, err)
			return
		}
		url, err := svc.DownloadInvoiceURL(c.Request.Context(), ws, c.Param("id"))
		if err != nil {
			response.Error(c, err)
			return
		}
		c.Redirect(http.StatusFound, url)
	})

	// Pay invoice
	grp.POST("/invoices/:id/pay", func(c *gin.Context) {
		actor, err := request.ActorFromContext(c)
		if err != nil {
			response.Error(c, err)
			return
		}
		ws, err := accountSvc.GetWorkspaceByID(c.Request.Context(), actor.WorkspaceID)
		if err != nil {
			response.Error(c, err)
			return
		}
		if err := svc.PayInvoice(c.Request.Context(), ws, c.Param("id")); err != nil {
			response.Error(c, err)
			return
		}
		response.SuccessEmpty(c, http.StatusNoContent)
	})
}

// Webhook verifies signature and processes relevant events; delegates to service
func Webhook(c *gin.Context, svc *billing.Service) {
	secret := viper.GetString(config.StripeWebhookSecret)
	if secret == "" {
		response.Error(c, gin.Error{Err: io.EOF})
		return
	}
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.Error(c, err)
		return
	}
	sig := c.GetHeader("Stripe-Signature")
	evt, err := webhook.ConstructEvent(payload, sig, secret)
	if err != nil {
		response.Error(c, err)
		return
	}
	switch evt.Type {
	case "customer.subscription.created", "customer.subscription.updated":
		var sub struct {
			ID                 string `json:"id"`
			Status             string `json:"status"`
			CurrentPeriodStart int64  `json:"current_period_start"`
			CurrentPeriodEnd   int64  `json:"current_period_end"`
		}
		if err := json.Unmarshal(evt.Data.Raw, &sub); err == nil {
			go svc.SyncSubscriptionStatus(c.Request.Context(), sub.ID, sub.Status, sub.CurrentPeriodStart, sub.CurrentPeriodEnd)
		}
	case "customer.subscription.deleted":
		var sub struct {
			ID                 string `json:"id"`
			CurrentPeriodStart int64  `json:"current_period_start"`
			CurrentPeriodEnd   int64  `json:"current_period_end"`
		}
		if err := json.Unmarshal(evt.Data.Raw, &sub); err == nil {
			// Attempt prorated refund with provided period bounds
			go svc.RefundAndFinalizeCancellation(c.Request.Context(), sub.ID, sub.CurrentPeriodStart, sub.CurrentPeriodEnd)
		}
	case "invoice.payment_failed":
		var iv struct {
			Subscription string `json:"subscription"`
		}
		if err := json.Unmarshal(evt.Data.Raw, &iv); err == nil && iv.Subscription != "" {
			go svc.MarkSubscriptionPastDue(c.Request.Context(), iv.Subscription)
		}
	case "invoice.payment_succeeded":
		var iv struct {
			Subscription string `json:"subscription"`
		}
		if err := json.Unmarshal(evt.Data.Raw, &iv); err == nil && iv.Subscription != "" {
			go svc.MarkSubscriptionActive(c.Request.Context(), iv.Subscription)
		}
	default:
		// ignore other events
	}
	c.Status(http.StatusOK)
}
