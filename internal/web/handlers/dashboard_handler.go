package handlers

import (
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/expense"
	"github.com/abdelrahman146/kyora/internal/domain/order"
	"github.com/abdelrahman146/kyora/internal/domain/owner"
	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/domain/supplier"
	"github.com/abdelrahman146/kyora/internal/web/views/pages"
	"github.com/abdelrahman146/kyora/internal/web/webcontext"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
)

type dashboardHandler struct {
	storeDomain    *store.StoreDomain
	orderDomain    *order.OrderDomain
	ownerDomain    *owner.OwnerDomain
	expenseDomain  *expense.ExpenseDomain
	customerDomain *customer.CustomerDomain
	supplierDomain *supplier.SupplierDomain
}

func AddDashboardRoutes(
	r *gin.RouterGroup,
	storeDomain *store.StoreDomain,
	orderDomain *order.OrderDomain,
	ownerDomain *owner.OwnerDomain,
	expenseDomain *expense.ExpenseDomain,
	customerDomain *customer.CustomerDomain,
	supplierDomain *supplier.SupplierDomain,
) {
	h := &dashboardHandler{
		storeDomain:    storeDomain,
		orderDomain:    orderDomain,
		ownerDomain:    ownerDomain,
		expenseDomain:  expenseDomain,
		customerDomain: customerDomain,
		supplierDomain: supplierDomain,
	}
	h.registerRoutes(r)
}

func (h *dashboardHandler) registerRoutes(c *gin.RouterGroup) {
	r := c.Group("/dashboard")
	{
		r.GET("/", h.index)
	}
}

func (h *dashboardHandler) index(c *gin.Context) {
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Dashboard",
		Description: "Dashboard page",
		Keywords:    "dashboard, Kyora",
		Path:        "/",
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "Home", Label: "/"},
			{Href: "Dashboard", Label: ""},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.Dashboard())
}
