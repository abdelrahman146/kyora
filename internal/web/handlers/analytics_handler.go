package handlers

import (
	"github.com/abdelrahman146/kyora/internal/web/views/pages"
	"github.com/abdelrahman146/kyora/internal/web/webcontext"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
}

func NewAnalyticsHandler() *AnalyticsHandler {
	return &AnalyticsHandler{}
}

func (h *AnalyticsHandler) RegisterRoutes(r gin.IRoutes) {
	r.GET("/analytics/sales", h.GetSalesReport)
	r.GET("/analytics/pnl", h.GetPnlReport)
	r.GET("/analytics/customers", h.GetCustomersReport)
	r.GET("/analytics/expenses", h.GetExpensesReport)
}

func (h *AnalyticsHandler) GetSalesReport(c *gin.Context) {
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

func (h *AnalyticsHandler) GetPnlReport(c *gin.Context) {
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

func (h *AnalyticsHandler) GetCustomersReport(c *gin.Context) {
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

func (h *AnalyticsHandler) GetExpensesReport(c *gin.Context) {
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
