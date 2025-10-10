package handlers

import (
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
