package webrouter

import (
	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/asset"
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/expense"
	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/domain/order"
	"github.com/abdelrahman146/kyora/internal/domain/owner"
	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/domain/supplier"
	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/abdelrahman146/kyora/internal/web/handlers"
	"github.com/abdelrahman146/kyora/internal/web/middleware"
	"github.com/abdelrahman146/kyora/internal/web/webutils"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine,
	accountDomain *account.AccountDomain,
	storeDomain *store.StoreDomain,
	inventoryDomain *inventory.InventoryDomain,
	orderDomain *order.OrderDomain,
	customerDomain *customer.CustomerDomain,
	ownerDomain *owner.OwnerDomain,
	assetDomain *asset.AssetDomain,
	expenseDomain *expense.ExpenseDomain,
	supplierDomain *supplier.SupplierDomain,
) {
	r.Static("/static", "./public")
	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.POST("/logout", func(c *gin.Context) {
		utils.JWT.ClearJwtCookie(c)
		webutils.Redirect(c, "/login")
	})

	// handlers
	handlers.AddAuthRoutes(r, accountDomain)
	handlers.AddOnboardingRoutes(r, accountDomain)
	handlers.AddAccountRoutes(r, accountDomain)
	r.Group("/:storeId", middleware.AuthRequired, middleware.UserRequired(accountDomain.AuthService), middleware.StoreRequired(storeDomain.StoreService))
	{
		handlers.AddStoreRoutes(r, storeDomain)
		handlers.AddDashboardRoutes(r, storeDomain, orderDomain, ownerDomain, expenseDomain, customerDomain, supplierDomain)
		handlers.AddInventoryRoutes(r, inventoryDomain)
		handlers.AddOrderRoutes(r, orderDomain)
		handlers.AddSupplierRoutes(r, supplierDomain)
		handlers.AddCustomerRoutes(r, customerDomain)
		handlers.AddExpenseRoutes(r, expenseDomain)
		handlers.AddAnalyticsRoutes(r, storeDomain, orderDomain, ownerDomain, inventoryDomain, expenseDomain, customerDomain, supplierDomain)

	}
}
