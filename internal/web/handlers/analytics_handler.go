package handlers

import (
	"fmt"

	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/expense"
	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/domain/order"
	"github.com/abdelrahman146/kyora/internal/domain/owner"
	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/domain/supplier"
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

func AddAnalyticsRoutes(c *gin.RouterGroup, storeDomain *store.StoreDomain, orderDomain *order.OrderDomain, ownerDomain *owner.OwnerDomain, inventoryDomain *inventory.InventoryDomain, expenseDomain *expense.ExpenseDomain, customerDomain *customer.CustomerDomain, supplierDomain *supplier.SupplierDomain) {
	h := &analyticsHandler{
		storeDomain:     storeDomain,
		orderDomain:     orderDomain,
		ownerDomain:     ownerDomain,
		inventoryDomain: inventoryDomain,
		expenseDomain:   expenseDomain,
		customerDomain:  customerDomain,
		supplierDomain:  supplierDomain,
	}
	h.registerRoutes(c)
}

func (h *analyticsHandler) registerRoutes(c *gin.RouterGroup) {
	r := c.Group("/analytics")
	{
		r.GET("/inventory", h.getInventoryReport)
		r.GET("/sales", h.getSalesReport)
		r.GET("/customers", h.getCustomersReport)
		r.GET("/expenses", h.getExpensesReport)
		r.GET("/capital", h.getCapitalReport)
	}
}

func (h *analyticsHandler) getSalesReport(c *gin.Context) {
	path := c.Request.URL.RequestURI()
	storeId := c.Param("storeId")
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Sales Report",
		Description: "View the sales report",
		Keywords:    "sales report, Kyora",
		Path:        path,
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: fmt.Sprintf("/%s/dashboard", storeId), Label: "Dashboard"},
			{Href: fmt.Sprintf("/%s/analytics", storeId), Label: "Analytics"},
			{Label: "Sales Report"},
		},
	}
	_ = webutils.NewFromDate(c.Request.Context(), c.Query("from"))
	_ = webutils.NewToDate(c.Request.Context(), c.Query("to"))
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.SalesReport())
}

func (h *analyticsHandler) getCustomersReport(c *gin.Context) {
	path := c.Request.URL.RequestURI()
	storeId := c.Param("storeId")
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Customers Report",
		Description: "View the customers report",
		Keywords:    "customers report, Kyora",
		Path:        path,
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: fmt.Sprintf("/%s/dashboard", storeId), Label: "Dashboard"},
			{Href: fmt.Sprintf("/%s/analytics", storeId), Label: "Analytics"},
			{Label: "Customers Report"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.CustomersReport())
}

func (h *analyticsHandler) getExpensesReport(c *gin.Context) {
	path := c.Request.URL.RequestURI()
	storeId := c.Param("storeId")
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Expenses Report",
		Description: "View the expenses report",
		Keywords:    "expenses report, Kyora",
		Path:        path,
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: fmt.Sprintf("/%s/dashboard", storeId), Label: "Dashboard"},
			{Href: fmt.Sprintf("/%s/analytics", storeId), Label: "Analytics"},
			{Label: "Expenses Report"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.ExpensesReport())
}

func (h *analyticsHandler) getInventoryReport(c *gin.Context) {
	path := c.Request.URL.RequestURI()
	storeId := c.Param("storeId")
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Inventory Report",
		Description: "View the inventory report",
		Keywords:    "inventory report, Kyora",
		Path:        path,
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: fmt.Sprintf("/%s/dashboard", storeId), Label: "Dashboard"},
			{Href: fmt.Sprintf("/%s/analytics", storeId), Label: "Analytics"},
			{Label: "Inventory Report"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	c.String(200, "Not implemented")
}

func (h *analyticsHandler) getCapitalReport(c *gin.Context) {
	path := c.Request.URL.RequestURI()
	storeId := c.Param("storeId")
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Capital Report",
		Description: "View the capital report",
		Keywords:    "capital report, Kyora",
		Path:        path,
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: fmt.Sprintf("/%s/dashboard", storeId), Label: "Dashboard"},
			{Href: fmt.Sprintf("/%s/analytics", storeId), Label: "Analytics"},
			{Label: "Capital Report"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	c.String(200, "Not implemented")
}
