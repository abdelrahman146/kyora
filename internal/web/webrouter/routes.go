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
	handlers.AddAuthRoutes(&r.RouterGroup, accountDomain)
	handlers.AddOnboardingRoutes(&r.RouterGroup, accountDomain)
	handlers.AddAccountRoutes(&r.RouterGroup, accountDomain)
	rs := r.Group("/:storeId", middleware.AuthRequired, middleware.UserRequired(accountDomain.AuthService), middleware.StoreRequired(storeDomain.StoreService))
	{
		handlers.AddStoreRoutes(rs, storeDomain)
		handlers.AddDashboardRoutes(rs, storeDomain, orderDomain, ownerDomain, expenseDomain, customerDomain, supplierDomain)
		handlers.AddInventoryRoutes(rs, inventoryDomain)
		handlers.AddOrderRoutes(rs, orderDomain)
		handlers.AddSupplierRoutes(rs, supplierDomain)
		handlers.AddCustomerRoutes(rs, customerDomain)
		handlers.AddExpenseRoutes(rs, expenseDomain)
		handlers.AddAnalyticsRoutes(rs, storeDomain, orderDomain, ownerDomain, inventoryDomain, expenseDomain, customerDomain, supplierDomain)

	}
}
