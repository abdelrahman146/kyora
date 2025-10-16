package handlers

import (
	"fmt"

	"github.com/abdelrahman146/kyora/internal/domain/analytics"
	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/abdelrahman146/kyora/internal/web/views/pages"
	"github.com/abdelrahman146/kyora/internal/web/webcontext"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type analyticsHandler struct {
	storeDomain     *store.StoreDomain
	analyticsDomain *analytics.AnalyticsDomain
}

func AddAnalyticsRoutes(c *gin.RouterGroup, storeDomain *store.StoreDomain, analyticsDomain *analytics.AnalyticsDomain) {
	h := &analyticsHandler{
		storeDomain:     storeDomain,
		analyticsDomain: analyticsDomain,
	}
	h.registerRoutes(c)
}

func (h *analyticsHandler) registerRoutes(c *gin.RouterGroup) {
	r := c.Group("/analytics")
	{
		r.GET("/sales", h.getSalesReport)
		r.GET("/inventory", h.getInventoryReport)
		r.GET("/customers", h.getCustomersReport)
		r.GET("/assets", h.getAssetsReport)
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
	from := cast.ToTime(c.Query("from"))
	to := cast.ToTime(c.Query("to"))
	_, err := h.analyticsDomain.Service.GenerateSalesAnalytics(c.Request.Context(), storeId, from, to)
	if err != nil {
		utils.Log.FromContext(c.Request.Context()).Error("failed to generate sales report", "error", err, "storeId", storeId)
		webutils.Render(c, 500, pages.ErrorPage(500, "Failed to load sales report"))
		return
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.NotImplemented("Sales Report"))
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
	from := cast.ToTime(c.Query("from"))
	to := cast.ToTime(c.Query("to"))
	_, err := h.analyticsDomain.Service.GenerateCustomerAnalytics(c.Request.Context(), storeId, from, to)
	if err != nil {
		utils.Log.FromContext(c.Request.Context()).Error("failed to generate customers report", "error", err, "storeId", storeId)
		webutils.Render(c, 500, pages.ErrorPage(500, "Failed to load customers report"))
		return
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.NotImplemented("Customers Report"))
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
	from := cast.ToTime(c.Query("from"))
	to := cast.ToTime(c.Query("to"))
	_, err := h.analyticsDomain.Service.GenerateExpenseAnalytics(c.Request.Context(), storeId, from, to)
	if err != nil {
		utils.Log.FromContext(c.Request.Context()).Error("failed to generate expenses report", "error", err, "storeId", storeId)
		webutils.Render(c, 500, pages.ErrorPage(500, "Failed to load expenses report"))
		return
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.NotImplemented("Expenses Report"))
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
	from := cast.ToTime(c.Query("from"))
	to := cast.ToTime(c.Query("to"))
	_, err := h.analyticsDomain.Service.GenerateInventoryAnalytics(c.Request.Context(), storeId, from, to)
	if err != nil {
		utils.Log.FromContext(c.Request.Context()).Error("failed to generate inventory report", "error", err, "storeId", storeId)
		webutils.Render(c, 500, pages.ErrorPage(500, "Failed to load inventory report"))
		return
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.NotImplemented("Inventory Report"))
}

func (h *analyticsHandler) getAssetsReport(c *gin.Context) {
	path := c.Request.URL.RequestURI()
	storeId := c.Param("storeId")
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Assets Report",
		Description: "View the assets report",
		Keywords:    "assets report, Kyora",
		Path:        path,
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: fmt.Sprintf("/%s/dashboard", storeId), Label: "Dashboard"},
			{Href: fmt.Sprintf("/%s/analytics", storeId), Label: "Analytics"},
			{Label: "Assets Report"},
		},
	}
	from := cast.ToTime(c.Query("from"))
	to := cast.ToTime(c.Query("to"))
	_, err := h.analyticsDomain.Service.GenerateAssetAnalytics(c.Request.Context(), storeId, from, to)
	if err != nil {
		utils.Log.FromContext(c.Request.Context()).Error("failed to generate assets report", "error", err, "storeId", storeId)
		webutils.Render(c, 500, pages.ErrorPage(500, "Failed to load assets report"))
		return
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.NotImplemented("Assets Report"))
}
