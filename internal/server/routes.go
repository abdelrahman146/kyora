package server

import (
	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/billing"
	"github.com/abdelrahman146/kyora/internal/domain/onboarding"
	"github.com/abdelrahman146/kyora/internal/platform/auth"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/gin-gonic/gin"
)

func registerOnboardingRoutes(r *gin.Engine, h *onboarding.HttpHandler) {
	group := r.Group("/api/onboarding")
	group.POST("/start", h.Start)
	group.POST("/email/otp", h.SendEmailOTP)
	group.POST("/email/verify", h.VerifyEmail)
	group.POST("/oauth/google", h.OAuthGoogle)
	group.POST("/business", h.SetBusiness)
	group.POST("/payment/start", h.PaymentStart)
	group.POST("/complete", h.Complete)
}

func registerAccountRoutes(r *gin.Engine, h *account.HttpHandler, accountService *account.Service, billingService *billing.Service) {
	// Public authentication endpoints (no auth required)
	authGroup := r.Group("/v1/auth")
	{
		authGroup.POST("/login", h.Login)
		authGroup.POST("/google/login", h.LoginWithGoogle)
		authGroup.GET("/google/url", h.GetGoogleAuthURL)
		authGroup.POST("/forgot-password", h.ForgotPassword)
		authGroup.POST("/reset-password", h.ResetPassword)
		authGroup.POST("/verify-email/request", h.RequestEmailVerification)
		authGroup.POST("/verify-email", h.VerifyEmail)
	}

	// Public invitation acceptance endpoints (no auth required)
	invitationGroup := r.Group("/v1/invitations")
	{
		invitationGroup.POST("/accept", h.AcceptInvitation)
		invitationGroup.GET("/accept/google", h.AcceptInvitationWithGoogle)
	}

	// Protected user profile endpoints
	userGroup := r.Group("/v1/users")
	userGroup.Use(auth.EnforceAuthentication, account.EnforceValidActor(accountService))
	{
		userGroup.GET("/me", h.GetCurrentUser)
		userGroup.PATCH("/me", h.UpdateCurrentUser)
	}

	// Protected workspace endpoints
	workspaceGroup := r.Group("/v1/workspaces")
	workspaceGroup.Use(auth.EnforceAuthentication, account.EnforceValidActor(accountService))
	{
		// Workspace info (all authenticated users)
		workspaceGroup.GET("/me", h.GetCurrentWorkspace)

		// Workspace users (view permission required)
		workspaceGroup.GET("/users",
			account.EnforceActorPermissions(role.ActionView, role.ResourceAccount),
			h.GetWorkspaceUsers)
		workspaceGroup.GET("/users/:userId",
			account.EnforceActorPermissions(role.ActionView, role.ResourceAccount),
			h.GetWorkspaceUser)

		// User management (manage permission required)
		workspaceGroup.PATCH("/users/:userId/role",
			account.EnforceActorPermissions(role.ActionManage, role.ResourceAccount),
			h.UpdateUserRole)
		workspaceGroup.DELETE("/users/:userId",
			account.EnforceActorPermissions(role.ActionManage, role.ResourceAccount),
			h.RemoveUserFromWorkspace)

		// Invitation management (manage permission required)
		invitationsGroup := workspaceGroup.Group("/invitations")
		invitationsGroup.Use(account.EnforceActorPermissions(role.ActionManage, role.ResourceAccount))
		{
			invitationsGroup.POST("",
				billing.EnforceActiveSubscription(billingService),
				billing.EnforcePlanWorkspaceLimits(billing.PlanSchema.MaxTeamMembers, accountService.CountWorkspaceUsersForPlanLimit),
				h.InviteUserToWorkspace)
			invitationsGroup.GET("", h.GetWorkspaceInvitations)
			invitationsGroup.DELETE("/:invitationId", h.RevokeInvitation)
		}
	}
}

func registerBillingRoutes(r *gin.Engine, h *billing.HttpHandler, accountService *account.Service) {
	group := r.Group("/v1/billing")

	// Plan Operations (Public).
	group.GET("/plans", h.ListPlans)
	group.GET("/plans/:descriptor", h.GetPlan)

	// Subscription Operations
	subscriptionGroup := group.Group("/subscription")
	subscriptionGroup.Use(auth.EnforceAuthentication, account.EnforceValidActor(accountService), account.EnforceWorkspaceMembership(accountService))
	{
		subscriptionGroup.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceBilling), h.GetSubscription)
		subscriptionGroup.GET("/details", account.EnforceActorPermissions(role.ActionView, role.ResourceBilling), h.GetSubscriptionDetails)
		subscriptionGroup.POST("", account.EnforceActorPermissions(role.ActionManage, role.ResourceBilling), h.CreateSubscription)
		subscriptionGroup.DELETE("", account.EnforceActorPermissions(role.ActionManage, role.ResourceBilling), h.CancelSubscription)
		subscriptionGroup.POST("/resume", account.EnforceActorPermissions(role.ActionManage, role.ResourceBilling), h.ResumeSubscription)
		subscriptionGroup.POST("/schedule-change", account.EnforceActorPermissions(role.ActionManage, role.ResourceBilling), h.ScheduleSubscriptionChange)
		subscriptionGroup.POST("/estimate-proration", account.EnforceActorPermissions(role.ActionManage, role.ResourceBilling), h.EstimateProration)
		subscriptionGroup.POST("/trial/extend", account.EnforceActorPermissions(role.ActionManage, role.ResourceBilling), h.ExtendTrial)
	}

	// Payment Method Operations
	paymentGroup := group.Group("/payment-methods")
	paymentGroup.Use(auth.EnforceAuthentication, account.EnforceValidActor(accountService), account.EnforceWorkspaceMembership(accountService))
	{
		paymentGroup.POST("/attach", account.EnforceActorPermissions(role.ActionManage, role.ResourceBilling), h.AttachPaymentMethod)
		paymentGroup.POST("/setup-intent", account.EnforceActorPermissions(role.ActionManage, role.ResourceBilling), h.CreateSetupIntent)
	}

	// Invoice Operations
	invoiceGroup := group.Group("/invoices")
	invoiceGroup.Use(auth.EnforceAuthentication, account.EnforceValidActor(accountService), account.EnforceWorkspaceMembership(accountService))
	{
		invoiceGroup.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceBilling), h.ListInvoices)
		invoiceGroup.GET("/:id/download", account.EnforceActorPermissions(role.ActionView, role.ResourceBilling), h.DownloadInvoice)
		invoiceGroup.POST("/:id/pay", account.EnforceActorPermissions(role.ActionManage, role.ResourceBilling), h.PayInvoice)
		invoiceGroup.POST("", account.EnforceActorPermissions(role.ActionManage, role.ResourceBilling), h.CreateInvoice)
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

	// Tax and Usage
	usageGroup := group.Group("/usage")
	usageGroup.Use(auth.EnforceAuthentication, account.EnforceValidActor(accountService), account.EnforceWorkspaceMembership(accountService))
	{
		usageGroup.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceBilling), h.GetUsage)
	}

	taxGroup := group.Group("/tax")
	taxGroup.Use(auth.EnforceAuthentication, account.EnforceValidActor(accountService), account.EnforceWorkspaceMembership(accountService))
	{
		taxGroup.POST("/calculate", account.EnforceActorPermissions(role.ActionView, role.ResourceBilling), h.CalculateTax)
	}

	// Webhook endpoint (public - no auth required)
	r.POST("/webhooks/stripe", h.HandleWebhook)
}
