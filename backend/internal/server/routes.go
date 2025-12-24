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
	"github.com/abdelrahman146/kyora/internal/domain/order"
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
		authGroup.POST("/refresh", h.Refresh)
		authGroup.POST("/logout", h.Logout)
		authGroup.POST("/logout-all", h.LogoutAll)
		authGroup.POST("/logout-others", h.LogoutOtherDevices)
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

func registerBusinessRoutes(r *gin.Engine, h *business.HttpHandler, accountService *account.Service, billingService *billing.Service, businessService *business.Service) {
	group := r.Group("/v1/businesses")
	group.Use(auth.EnforceAuthentication, account.EnforceValidActor(accountService), account.EnforceWorkspaceMembership(accountService))

	group.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceBusiness), h.ListBusinesses)
	group.GET("/descriptor/availability", account.EnforceActorPermissions(role.ActionView, role.ResourceBusiness), h.CheckDescriptorAvailability)
	group.GET("/:businessDescriptor", account.EnforceActorPermissions(role.ActionView, role.ResourceBusiness), h.GetBusinessByDescriptor)

	group.POST("",
		account.EnforceActorPermissions(role.ActionManage, role.ResourceBusiness),
		billing.EnforceActiveSubscription(billingService),
		billing.EnforcePlanWorkspaceLimits(billing.PlanSchema.MaxBusinesses, businessService.MaxBusinessesEnforceFunc),
		h.CreateBusiness,
	)
	group.PATCH("/:businessDescriptor", account.EnforceActorPermissions(role.ActionManage, role.ResourceBusiness), h.UpdateBusiness)
	group.POST("/:businessDescriptor/archive", account.EnforceActorPermissions(role.ActionManage, role.ResourceBusiness), h.ArchiveBusiness)
	group.POST("/:businessDescriptor/unarchive", account.EnforceActorPermissions(role.ActionManage, role.ResourceBusiness), h.UnarchiveBusiness)
	group.DELETE("/:businessDescriptor", account.EnforceActorPermissions(role.ActionManage, role.ResourceBusiness), h.DeleteBusiness)
}

func registerBusinessScopedRoutes(
	r *gin.Engine,
	accountService *account.Service,
	billingService *billing.Service,
	businessService *business.Service,
	businessHandler *business.HttpHandler,
	accountingHandler *accounting.HttpHandler,
	analyticsHandler *analytics.HttpHandler,
	customerHandler *customer.HttpHandler,
	inventoryHandler *inventory.HttpHandler,
	orderHandler *order.HttpHandler,
) {
	group := r.Group("/v1/businesses/:businessDescriptor")
	group.Use(
		auth.EnforceAuthentication,
		account.EnforceValidActor(accountService),
		account.EnforceWorkspaceMembership(accountService),
		business.EnforceBusinessValidity(businessService),
	)

	// Customer routes
	customers := group.Group("/customers")
	{
		customers.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceCustomer), customerHandler.ListCustomers)
		customers.GET("/:customerId", account.EnforceActorPermissions(role.ActionView, role.ResourceCustomer), customerHandler.GetCustomer)
		customers.POST("", account.EnforceActorPermissions(role.ActionManage, role.ResourceCustomer), customerHandler.CreateCustomer)
		customers.PATCH("/:customerId", account.EnforceActorPermissions(role.ActionManage, role.ResourceCustomer), customerHandler.UpdateCustomer)
		customers.DELETE("/:customerId", account.EnforceActorPermissions(role.ActionManage, role.ResourceCustomer), customerHandler.DeleteCustomer)

		addressGroup := customers.Group("/:customerId/addresses")
		{
			addressGroup.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceCustomer), customerHandler.ListCustomerAddresses)
			addressGroup.POST("", account.EnforceActorPermissions(role.ActionManage, role.ResourceCustomer), customerHandler.CreateCustomerAddress)
			addressGroup.PATCH("/:addressId", account.EnforceActorPermissions(role.ActionManage, role.ResourceCustomer), customerHandler.UpdateCustomerAddress)
			addressGroup.DELETE("/:addressId", account.EnforceActorPermissions(role.ActionManage, role.ResourceCustomer), customerHandler.DeleteCustomerAddress)
		}

		noteGroup := customers.Group("/:customerId/notes")
		{
			noteGroup.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceCustomer), customerHandler.ListCustomerNotes)
			noteGroup.POST("", account.EnforceActorPermissions(role.ActionManage, role.ResourceCustomer), customerHandler.CreateCustomerNote)
			noteGroup.DELETE("/:noteId", account.EnforceActorPermissions(role.ActionManage, role.ResourceCustomer), customerHandler.DeleteCustomerNote)
		}
	}

	// Inventory routes
	inventoryGroup := group.Group("/inventory")
	{
		inventoryGroup.GET("/summary", account.EnforceActorPermissions(role.ActionView, role.ResourceInventory), inventoryHandler.GetInventorySummary)
		inventoryGroup.GET("/top-products", account.EnforceActorPermissions(role.ActionView, role.ResourceInventory), inventoryHandler.GetTopProductsByInventoryValue)

		products := inventoryGroup.Group("/products")
		{
			products.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceInventory), inventoryHandler.ListProducts)
			products.GET("/:productId", account.EnforceActorPermissions(role.ActionView, role.ResourceInventory), inventoryHandler.GetProduct)
			products.GET("/:productId/variants", account.EnforceActorPermissions(role.ActionView, role.ResourceInventory), inventoryHandler.ListProductVariants)
			products.POST("", account.EnforceActorPermissions(role.ActionManage, role.ResourceInventory), inventoryHandler.CreateProduct)
			products.POST("/with-variants", account.EnforceActorPermissions(role.ActionManage, role.ResourceInventory), inventoryHandler.CreateProductWithVariants)
			products.PATCH("/:productId", account.EnforceActorPermissions(role.ActionManage, role.ResourceInventory), inventoryHandler.UpdateProduct)
			products.DELETE("/:productId", account.EnforceActorPermissions(role.ActionManage, role.ResourceInventory), inventoryHandler.DeleteProduct)
		}

		variants := inventoryGroup.Group("/variants")
		{
			variants.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceInventory), inventoryHandler.ListVariants)
			variants.GET("/:variantId", account.EnforceActorPermissions(role.ActionView, role.ResourceInventory), inventoryHandler.GetVariant)
			variants.POST("", account.EnforceActorPermissions(role.ActionManage, role.ResourceInventory), inventoryHandler.CreateVariant)
			variants.PATCH("/:variantId", account.EnforceActorPermissions(role.ActionManage, role.ResourceInventory), inventoryHandler.UpdateVariant)
			variants.DELETE("/:variantId", account.EnforceActorPermissions(role.ActionManage, role.ResourceInventory), inventoryHandler.DeleteVariant)
		}

		categories := inventoryGroup.Group("/categories")
		{
			categories.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceInventory), inventoryHandler.ListCategories)
			categories.GET("/:categoryId", account.EnforceActorPermissions(role.ActionView, role.ResourceInventory), inventoryHandler.GetCategory)
			categories.POST("", account.EnforceActorPermissions(role.ActionManage, role.ResourceInventory), inventoryHandler.CreateCategory)
			categories.PATCH("/:categoryId", account.EnforceActorPermissions(role.ActionManage, role.ResourceInventory), inventoryHandler.UpdateCategory)
			categories.DELETE("/:categoryId", account.EnforceActorPermissions(role.ActionManage, role.ResourceInventory), inventoryHandler.DeleteCategory)
		}
	}

	// Shipping zones (business settings)
	shippingZones := group.Group("/shipping-zones")
	{
		shippingZones.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceBusiness), businessHandler.ListShippingZones)
		shippingZones.GET("/:zoneId", account.EnforceActorPermissions(role.ActionView, role.ResourceBusiness), businessHandler.GetShippingZone)

		manageZones := shippingZones.Group("")
		manageZones.Use(
			account.EnforceActorPermissions(role.ActionManage, role.ResourceBusiness),
			billing.EnforceActiveSubscription(billingService),
		)
		{
			manageZones.POST("", businessHandler.CreateShippingZone)
			manageZones.PATCH("/:zoneId", businessHandler.UpdateShippingZone)
			manageZones.DELETE("/:zoneId", businessHandler.DeleteShippingZone)
		}
	}

	// Payment methods (business settings)
	paymentMethods := group.Group("/payment-methods")
	{
		paymentMethods.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceBusiness), businessHandler.ListPaymentMethods)

		managePaymentMethods := paymentMethods.Group("")
		managePaymentMethods.Use(
			account.EnforceActorPermissions(role.ActionManage, role.ResourceBusiness),
			billing.EnforceActiveSubscription(billingService),
		)
		{
			managePaymentMethods.PATCH("/:descriptor", businessHandler.UpdatePaymentMethod)
		}
	}

	// Order routes
	orders := group.Group("/orders")
	{
		orders.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceOrder), orderHandler.ListOrders)
		orders.GET("/by-number/:orderNumber", account.EnforceActorPermissions(role.ActionView, role.ResourceOrder), orderHandler.GetOrderByNumber)
		orders.GET("/:orderId", account.EnforceActorPermissions(role.ActionView, role.ResourceOrder), orderHandler.GetOrder)

		manageOrders := orders.Group("")
		manageOrders.Use(
			account.EnforceActorPermissions(role.ActionManage, role.ResourceOrder),
			billing.EnforceActiveSubscription(billingService),
			billing.EnforcePlanFeatureRestriction(billing.PlanSchema.OrderManagement),
		)
		{
			manageOrders.POST("",
				billing.EnforcePlanWorkspaceLimits(billing.PlanSchema.MaxOrdersPerMonth, billingService.CountMonthlyOrdersForPlanLimit),
				orderHandler.CreateOrder,
			)
			manageOrders.PATCH("/:orderId", orderHandler.UpdateOrder)
			manageOrders.DELETE("/:orderId", orderHandler.DeleteOrder)
			manageOrders.PATCH("/:orderId/status", orderHandler.UpdateOrderStatus)
			manageOrders.PATCH("/:orderId/payment-status", orderHandler.UpdateOrderPaymentStatus)
			manageOrders.PATCH("/:orderId/payment-details", orderHandler.AddOrderPaymentDetails)

			notes := manageOrders.Group("/:orderId/notes")
			{
				notes.POST("", orderHandler.CreateOrderNote)
				notes.PATCH("/:noteId", orderHandler.UpdateOrderNote)
				notes.DELETE("/:noteId", orderHandler.DeleteOrderNote)
			}
		}
	}

	// Analytics routes
	analyticsGroup := group.Group("/analytics")
	{
		analyticsGroup.GET("/dashboard", account.EnforceActorPermissions(role.ActionView, role.ResourceBasicAnalytics), analyticsHandler.GetDashboardMetrics)
		analyticsGroup.GET("/sales", account.EnforceActorPermissions(role.ActionView, role.ResourceBasicAnalytics), analyticsHandler.GetSalesAnalytics)
		analyticsGroup.GET("/inventory", account.EnforceActorPermissions(role.ActionView, role.ResourceBasicAnalytics), analyticsHandler.GetInventoryAnalytics)
		analyticsGroup.GET("/customers", account.EnforceActorPermissions(role.ActionView, role.ResourceBasicAnalytics), analyticsHandler.GetCustomerAnalytics)

		reports := analyticsGroup.Group("/reports")
		reports.Use(account.EnforceActorPermissions(role.ActionView, role.ResourceFinancialReports))
		{
			reports.GET("/financial-position", analyticsHandler.GetFinancialPosition)
			reports.GET("/profit-and-loss", analyticsHandler.GetProfitAndLoss)
			reports.GET("/cash-flow", analyticsHandler.GetCashFlow)
		}
	}

	// Accounting routes
	accountingGroup := group.Group("/accounting")
	{
		assets := accountingGroup.Group("/assets")
		{
			assets.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceAccounting), accountingHandler.ListAssets)
			assets.GET("/:assetId", account.EnforceActorPermissions(role.ActionView, role.ResourceAccounting), accountingHandler.GetAsset)
			assets.POST("", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), accountingHandler.CreateAsset)
			assets.PATCH("/:assetId", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), accountingHandler.UpdateAsset)
			assets.DELETE("/:assetId", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), accountingHandler.DeleteAsset)
		}

		investments := accountingGroup.Group("/investments")
		{
			investments.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceAccounting), accountingHandler.ListInvestments)
			investments.GET("/:investmentId", account.EnforceActorPermissions(role.ActionView, role.ResourceAccounting), accountingHandler.GetInvestment)
			investments.POST("", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), accountingHandler.CreateInvestment)
			investments.PATCH("/:investmentId", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), accountingHandler.UpdateInvestment)
			investments.DELETE("/:investmentId", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), accountingHandler.DeleteInvestment)
		}

		withdrawals := accountingGroup.Group("/withdrawals")
		{
			withdrawals.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceAccounting), accountingHandler.ListWithdrawals)
			withdrawals.GET("/:withdrawalId", account.EnforceActorPermissions(role.ActionView, role.ResourceAccounting), accountingHandler.GetWithdrawal)
			withdrawals.POST("", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), accountingHandler.CreateWithdrawal)
			withdrawals.PATCH("/:withdrawalId", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), accountingHandler.UpdateWithdrawal)
			withdrawals.DELETE("/:withdrawalId", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), accountingHandler.DeleteWithdrawal)
		}

		expenses := accountingGroup.Group("/expenses")
		{
			expenses.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceAccounting), accountingHandler.ListExpenses)
			expenses.GET("/:expenseId", account.EnforceActorPermissions(role.ActionView, role.ResourceAccounting), accountingHandler.GetExpense)
			expenses.POST("", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), accountingHandler.CreateExpense)
			expenses.PATCH("/:expenseId", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), accountingHandler.UpdateExpense)
			expenses.DELETE("/:expenseId", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), accountingHandler.DeleteExpense)
		}

		recurring := accountingGroup.Group("/recurring-expenses")
		{
			recurring.GET("", account.EnforceActorPermissions(role.ActionView, role.ResourceAccounting), accountingHandler.ListRecurringExpenses)
			recurring.GET("/:recurringExpenseId", account.EnforceActorPermissions(role.ActionView, role.ResourceAccounting), accountingHandler.GetRecurringExpense)
			recurring.POST("", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), accountingHandler.CreateRecurringExpense)
			recurring.PATCH("/:recurringExpenseId", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), accountingHandler.UpdateRecurringExpense)
			recurring.DELETE("/:recurringExpenseId", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), accountingHandler.DeleteRecurringExpense)
			recurring.PATCH("/:recurringExpenseId/status", account.EnforceActorPermissions(role.ActionManage, role.ResourceAccounting), accountingHandler.UpdateRecurringExpenseStatus)
			recurring.GET("/:recurringExpenseId/occurrences", account.EnforceActorPermissions(role.ActionView, role.ResourceAccounting), accountingHandler.GetRecurringExpenseOccurrences)
		}

		accountingGroup.GET("/summary", account.EnforceActorPermissions(role.ActionView, role.ResourceAccounting), accountingHandler.GetAccountingSummary)
	}
}
