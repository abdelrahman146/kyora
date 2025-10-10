package handlers

import (
	"github.com/abdelrahman146/kyora/internal/web/views/layouts"
	"github.com/abdelrahman146/kyora/internal/web/views/pages"
	"github.com/abdelrahman146/kyora/internal/web/webcontext"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

func (h *DashboardHandler) RegisterRoutes(r gin.IRoutes) {
	r.GET("/", h.Index)
}

func (h *DashboardHandler) Index(c *gin.Context) {
	info := webcontext.PageInfo{
		Title:       "Dashboard",
		Description: "Dashboard page",
		Keywords:    "dashboard, Kyora",
		Path:        "/",
		Breadcrumbs: []webcontext.Breadcrumb{
			{Title: "Home", Link: "/"},
			{Title: "Dashboard", Link: ""},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	stats := []pages.Stat{
		{Label: "Revenue (7d)", Value: "$2,430", Delta: "+8%"},
		{Label: "Gross Profit (7d)", Value: "$1,120", Delta: "+5%"},
		{Label: "Orders", Value: "34", Delta: "+2"},
		{Label: "AOV", Value: "$71.47", Delta: "+3%"},
		{Label: "New Customers", Value: "12", Delta: "+4"},
		{Label: "Unpaid Invoices", Value: "3", Delta: ""},
	}
	webutils.Render(c, 200, layouts.AppLayout(pages.Dashboard(stats)))
}
