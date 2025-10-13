package handlers

import (
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/expense"
	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/domain/order"
	"github.com/abdelrahman146/kyora/internal/domain/owner"
	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/domain/supplier"
	"github.com/abdelrahman146/kyora/internal/web/middleware"
	"github.com/abdelrahman146/kyora/internal/web/views/pages"
	"github.com/abdelrahman146/kyora/internal/web/webcontext"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
)

type analyticsHandler struct {
	storeDomain     *store.StoreDomain
	orderDomain     *order.OrderDomain
	ownerDomain     *owner.OwnerDomain
	inventoryDomain *inventory.InventoryDomain
	expenseDomain   *expense.ExpenseDomain
	customerDomain  *customer.CustomerDomain
	supplierDomain  *supplier.SupplierDomain
}

func AddAnalyticsRoutes(r *gin.Engine, storeDomain *store.StoreDomain, orderDomain *order.OrderDomain, ownerDomain *owner.OwnerDomain, inventoryDomain *inventory.InventoryDomain, expenseDomain *expense.ExpenseDomain, customerDomain *customer.CustomerDomain, supplierDomain *supplier.SupplierDomain) {
	h := &analyticsHandler{
		storeDomain:     storeDomain,
		orderDomain:     orderDomain,
		ownerDomain:     ownerDomain,
		inventoryDomain: inventoryDomain,
		expenseDomain:   expenseDomain,
		customerDomain:  customerDomain,
		supplierDomain:  supplierDomain,
	}
	h.registerRoutes(r)
}

func (h *analyticsHandler) registerRoutes(r *gin.Engine) {
	r.Use(middleware.AuthRequired, middleware.UserRequired(nil))
	r.GET("/analytics/sales", h.getSalesReport)
	r.GET("/analytics/pnl", h.getPnlReport)
	r.GET("/analytics/customers", h.getCustomersReport)
	r.GET("/analytics/expenses", h.getExpensesReport)
}

func (h *analyticsHandler) getSalesReport(c *gin.Context) {
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Sales Report",
		Description: "View the sales report",
		Keywords:    "sales report, Kyora",
		Path:        "/analytics/sales",
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "/", Label: "Dashboard"},
			{Href: "/analytics", Label: "Analytics"},
			{Label: "Sales Report"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.SalesReport())
}

func (h *analyticsHandler) getPnlReport(c *gin.Context) {
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Profit and Loss Report",
		Description: "View the profit and loss report",
		Keywords:    "profit and loss report, Kyora",
		Path:        "/analytics/pnl",
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "/", Label: "Dashboard"},
			{Href: "/analytics", Label: "Analytics"},
			{Label: "Profit and Loss Report"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.ReportPNL())
}

func (h *analyticsHandler) getCustomersReport(c *gin.Context) {
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Customers Report",
		Description: "View the customers report",
		Keywords:    "customers report, Kyora",
		Path:        "/analytics/customers",
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "/", Label: "Dashboard"},
			{Href: "/analytics", Label: "Analytics"},
			{Label: "Customers Report"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.CustomersReport())
}

func (h *analyticsHandler) getExpensesReport(c *gin.Context) {
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Expenses Report",
		Description: "View the expenses report",
		Keywords:    "expenses report, Kyora",
		Path:        "/analytics/expenses",
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "/", Label: "Dashboard"},
			{Href: "/analytics", Label: "Analytics"},
			{Label: "Expenses Report"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.ExpensesReport())
}
