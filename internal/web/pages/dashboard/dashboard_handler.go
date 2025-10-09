package dashboardPage

import (
	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/web/middleware"
	"github.com/abdelrahman146/kyora/internal/web/webcontext"
	"github.com/abdelrahman146/kyora/internal/web/webutils"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	authService *account.AuthenticationService
}

func NewDashboardHandler(authService *account.AuthenticationService) *DashboardHandler {
	return &DashboardHandler{authService}
}

func (h *DashboardHandler) RegisterRoutes(r gin.IRoutes) {
	r.GET("/", middleware.AuthRequiredMiddleware, middleware.UserMiddleware(h.authService), h.Index)
}

func (h *DashboardHandler) Index(c *gin.Context) {
	info := webcontext.PageInfo{
		Title:       "Dashboard",
		Description: "Dashboard page",
		Keywords:    "dashboard, zard",
		Breadcrumbs: []webcontext.Breadcrumb{
			{Title: "Home", Link: "/"},
			{Title: "Dashboard", Link: ""},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, Dashboard())
}
