package server

import (
	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/accounting"
	"github.com/abdelrahman146/kyora/internal/domain/analytics"
	"github.com/abdelrahman146/kyora/internal/domain/billing"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/inventory"
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
	workspaceGroup.Use(auth.EnforceAuthentication, account.EnforceValidActor(accountService), account.EnforceWorkspaceMembership(accountService))
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
		subscriptionGroup.GET("/trial", account.EnforceActorPermissions(role.ActionView, role.ResourceBilling), h.GetTrialStatus)
		subscriptionGroup.POST("/trial/extend", account.EnforceActorPermissions(role.ActionManage, role.ResourceBilling), h.ExtendTrial)
		subscriptionGroup.POST("/grace-period", account.EnforceActorPermissions(role.ActionManage, role.ResourceBilling), h.SetGracePeriod)
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
	checkoutGroup.Use(auth.EnforceAuthentication, account.EnforceValidActor(accountService), account.EnforceWorkspaceMembership(accountService))
	{
		checkoutGroup.POST("/session", account.EnforceActorPermissions(role.ActionManage, role.ResourceBilling), h.CreateCheckoutSession)
	}

	portalGroup := group.Group("/portal")
	portalGroup.Use(auth.EnforceAuthentication, account.EnforceValidActor(accountService), account.EnforceWorkspaceMembership(accountService))
	{
		portalGroup.POST("/session", account.EnforceActorPermissions(role.ActionManage, role.ResourceBilling), h.CreateBillingPortalSession)
	}

	// Tax and Usage
	usageGroup := group.Group("/usage")
	usageGroup.Use(auth.EnforceAuthentication, account.EnforceValidActor(accountService), account.EnforceWorkspaceMembership(accountService))
	{
		usageGroup.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceBilling), h.GetUsage)
		usageGroup.GET("/quota", account.EnforceActorPermissions(role.ActionView, role.ResourceBilling), h.GetUsageQuota)
	}

	taxGroup := group.Group("/tax")
	taxGroup.Use(auth.EnforceAuthentication, account.EnforceValidActor(accountService), account.EnforceWorkspaceMembership(accountService))
	{
		taxGroup.POST("/calculate", account.EnforceActorPermissions(role.ActionView, role.ResourceBilling), h.CalculateTax)
	}

	// Webhook endpoint (public - no auth required)
	r.POST("/webhooks/stripe", h.HandleWebhook)
}

func registerAccountingRoutes(r *gin.Engine, h *accounting.HttpHandler, accountService *account.Service) {
	group := r.Group("/v1/accounting")
	group.Use(auth.EnforceAuthentication, account.EnforceValidActor(accountService), account.EnforceWorkspaceMembership(accountService))

	// Asset Operations
	assetGroup := group.Group("/assets")
	{
		assetGroup.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceAccounting), h.ListAssets)
		assetGroup.GET("/:assetId", account.EnforceActorPermissions(role.ActionView, role.ResourceAccounting), h.GetAsset)
		assetGroup.POST("", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), h.CreateAsset)
		assetGroup.PATCH("/:assetId", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), h.UpdateAsset)
		assetGroup.DELETE("/:assetId", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), h.DeleteAsset)
	}

	// Investment Operations
	investmentGroup := group.Group("/investments")
	{
		investmentGroup.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceAccounting), h.ListInvestments)
		investmentGroup.GET("/:investmentId", account.EnforceActorPermissions(role.ActionView, role.ResourceAccounting), h.GetInvestment)
		investmentGroup.POST("", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), h.CreateInvestment)
		investmentGroup.PATCH("/:investmentId", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), h.UpdateInvestment)
		investmentGroup.DELETE("/:investmentId", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), h.DeleteInvestment)
	}

	// Withdrawal Operations
	withdrawalGroup := group.Group("/withdrawals")
	{
		withdrawalGroup.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceAccounting), h.ListWithdrawals)
		withdrawalGroup.GET("/:withdrawalId", account.EnforceActorPermissions(role.ActionView, role.ResourceAccounting), h.GetWithdrawal)
		withdrawalGroup.POST("", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), h.CreateWithdrawal)
		withdrawalGroup.PATCH("/:withdrawalId", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), h.UpdateWithdrawal)
		withdrawalGroup.DELETE("/:withdrawalId", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), h.DeleteWithdrawal)
	}

	// Expense Operations
	expenseGroup := group.Group("/expenses")
	{
		expenseGroup.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceAccounting), h.ListExpenses)
		expenseGroup.GET("/:expenseId", account.EnforceActorPermissions(role.ActionView, role.ResourceAccounting), h.GetExpense)
		expenseGroup.POST("", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), h.CreateExpense)
		expenseGroup.PATCH("/:expenseId", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), h.UpdateExpense)
		expenseGroup.DELETE("/:expenseId", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), h.DeleteExpense)
	}

	// Recurring Expense Operations
	recurringExpenseGroup := group.Group("/recurring-expenses")
	{
		recurringExpenseGroup.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceAccounting), h.ListRecurringExpenses)
		recurringExpenseGroup.GET("/:recurringExpenseId", account.EnforceActorPermissions(role.ActionView, role.ResourceAccounting), h.GetRecurringExpense)
		recurringExpenseGroup.POST("", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), h.CreateRecurringExpense)
		recurringExpenseGroup.PATCH("/:recurringExpenseId", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), h.UpdateRecurringExpense)
		recurringExpenseGroup.DELETE("/:recurringExpenseId", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), h.DeleteRecurringExpense)
		recurringExpenseGroup.PATCH("/:recurringExpenseId/status", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), h.UpdateRecurringExpenseStatus)
		recurringExpenseGroup.GET("/:recurringExpenseId/occurrences", account.EnforceActorPermissions(role.ActionView, role.ResourceAccounting), h.GetRecurringExpenseOccurrences)
	}

	// Summary Operations
	group.GET("/summary", account.EnforceActorPermissions(role.ActionView, role.ResourceAccounting), h.GetAccountingSummary)
}

func registerAnalyticsRoutes(r *gin.Engine, h *analytics.HttpHandler, accountService *account.Service) {
	group := r.Group("/v1/analytics")
	group.Use(auth.EnforceAuthentication, account.EnforceValidActor(accountService), account.EnforceWorkspaceMembership(accountService))

	group.GET("/dashboard", account.EnforceActorPermissions(role.ActionView, role.ResourceBasicAnalytics), h.GetDashboardMetrics)
	group.GET("/sales", account.EnforceActorPermissions(role.ActionView, role.ResourceBasicAnalytics), h.GetSalesAnalytics)
	group.GET("/inventory", account.EnforceActorPermissions(role.ActionView, role.ResourceBasicAnalytics), h.GetInventoryAnalytics)
	group.GET("/customers", account.EnforceActorPermissions(role.ActionView, role.ResourceBasicAnalytics), h.GetCustomerAnalytics)

	reports := group.Group("/reports")
	reports.Use(account.EnforceActorPermissions(role.ActionView, role.ResourceFinancialReports))
	{
		reports.GET("/financial-position", h.GetFinancialPosition)
		reports.GET("/profit-and-loss", h.GetProfitAndLoss)
		reports.GET("/cash-flow", h.GetCashFlow)
	}
}

func registerBusinessRoutes(r *gin.Engine, h *business.HttpHandler, accountService *account.Service, billingService *billing.Service, businessService *business.Service) {
	group := r.Group("/v1/businesses")
	group.Use(auth.EnforceAuthentication, account.EnforceValidActor(accountService), account.EnforceWorkspaceMembership(accountService))

	group.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceBusiness), h.ListBusinesses)
	group.GET("/descriptor/availability", account.EnforceActorPermissions(role.ActionView, role.ResourceBusiness), h.CheckDescriptorAvailability)
	group.GET("/descriptor/:businessDescriptor", account.EnforceActorPermissions(role.ActionView, role.ResourceBusiness), h.GetBusinessByDescriptor)
	group.GET("/:businessId", account.EnforceActorPermissions(role.ActionView, role.ResourceBusiness), h.GetBusiness)

	group.POST("",
		account.EnforceActorPermissions(role.ActionManage, role.ResourceBusiness),
		billing.EnforceActiveSubscription(billingService),
		billing.EnforcePlanWorkspaceLimits(billing.PlanSchema.MaxBusinesses, businessService.MaxBusinessesEnforceFunc),
		h.CreateBusiness,
	)
	group.PATCH("/:businessId", account.EnforceActorPermissions(role.ActionManage, role.ResourceBusiness), h.UpdateBusiness)
	group.POST("/:businessId/archive", account.EnforceActorPermissions(role.ActionManage, role.ResourceBusiness), h.ArchiveBusiness)
	group.POST("/:businessId/unarchive", account.EnforceActorPermissions(role.ActionManage, role.ResourceBusiness), h.UnarchiveBusiness)
	group.DELETE("/:businessId", account.EnforceActorPermissions(role.ActionManage, role.ResourceBusiness), h.DeleteBusiness)
}

func registerCustomerRoutes(r *gin.Engine, h *customer.HttpHandler, accountService *account.Service) {
	group := r.Group("/v1/customers")
	group.Use(auth.EnforceAuthentication, account.EnforceValidActor(accountService), account.EnforceWorkspaceMembership(accountService))

	// Customer CRUD operations
	group.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceCustomer), h.ListCustomers)
	group.GET("/:customerId", account.EnforceActorPermissions(role.ActionView, role.ResourceCustomer), h.GetCustomer)
	group.POST("", account.EnforceActorPermissions(role.ActionManage, role.ResourceCustomer), h.CreateCustomer)
	group.PATCH("/:customerId", account.EnforceActorPermissions(role.ActionManage, role.ResourceCustomer), h.UpdateCustomer)
	group.DELETE("/:customerId", account.EnforceActorPermissions(role.ActionManage, role.ResourceCustomer), h.DeleteCustomer)

	// Customer address operations
	addressGroup := group.Group("/:customerId/addresses")
	{
		addressGroup.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceCustomer), h.ListCustomerAddresses)
		addressGroup.POST("", account.EnforceActorPermissions(role.ActionManage, role.ResourceCustomer), h.CreateCustomerAddress)
		addressGroup.PATCH("/:addressId", account.EnforceActorPermissions(role.ActionManage, role.ResourceCustomer), h.UpdateCustomerAddress)
		addressGroup.DELETE("/:addressId", account.EnforceActorPermissions(role.ActionManage, role.ResourceCustomer), h.DeleteCustomerAddress)
	}

	// Customer note operations
	noteGroup := group.Group("/:customerId/notes")
	{
		noteGroup.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceCustomer), h.ListCustomerNotes)
		noteGroup.POST("", account.EnforceActorPermissions(role.ActionManage, role.ResourceCustomer), h.CreateCustomerNote)
		noteGroup.DELETE("/:noteId", account.EnforceActorPermissions(role.ActionManage, role.ResourceCustomer), h.DeleteCustomerNote)
	}
}

func registerInventoryRoutes(r *gin.Engine, h *inventory.HttpHandler, accountService *account.Service) {
	group := r.Group("/v1/inventory")
	group.Use(auth.EnforceAuthentication, account.EnforceValidActor(accountService), account.EnforceWorkspaceMembership(accountService))

	group.GET("/summary", account.EnforceActorPermissions(role.ActionView, role.ResourceInventory), h.GetInventorySummary)
	group.GET("/top-products", account.EnforceActorPermissions(role.ActionView, role.ResourceInventory), h.GetTopProductsByInventoryValue)

	products := group.Group("/products")
	{
		products.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceInventory), h.ListProducts)
		products.GET("/:productId", account.EnforceActorPermissions(role.ActionView, role.ResourceInventory), h.GetProduct)
		products.GET("/:productId/variants", account.EnforceActorPermissions(role.ActionView, role.ResourceInventory), h.ListProductVariants)
		products.POST("", account.EnforceActorPermissions(role.ActionManage, role.ResourceInventory), h.CreateProduct)
		products.POST("/with-variants", account.EnforceActorPermissions(role.ActionManage, role.ResourceInventory), h.CreateProductWithVariants)
		products.PATCH("/:productId", account.EnforceActorPermissions(role.ActionManage, role.ResourceInventory), h.UpdateProduct)
		products.DELETE("/:productId", account.EnforceActorPermissions(role.ActionManage, role.ResourceInventory), h.DeleteProduct)
	}

	variants := group.Group("/variants")
	{
		variants.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceInventory), h.ListVariants)
		variants.GET("/:variantId", account.EnforceActorPermissions(role.ActionView, role.ResourceInventory), h.GetVariant)
		variants.POST("", account.EnforceActorPermissions(role.ActionManage, role.ResourceInventory), h.CreateVariant)
		variants.PATCH("/:variantId", account.EnforceActorPermissions(role.ActionManage, role.ResourceInventory), h.UpdateVariant)
		variants.DELETE("/:variantId", account.EnforceActorPermissions(role.ActionManage, role.ResourceInventory), h.DeleteVariant)
	}

	categories := group.Group("/categories")
	{
		categories.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceInventory), h.ListCategories)
		categories.GET("/:categoryId", account.EnforceActorPermissions(role.ActionView, role.ResourceInventory), h.GetCategory)
		categories.POST("", account.EnforceActorPermissions(role.ActionManage, role.ResourceInventory), h.CreateCategory)
		categories.PATCH("/:categoryId", account.EnforceActorPermissions(role.ActionManage, role.ResourceInventory), h.UpdateCategory)
		categories.DELETE("/:categoryId", account.EnforceActorPermissions(role.ActionManage, role.ResourceInventory), h.DeleteCategory)
	}
}
