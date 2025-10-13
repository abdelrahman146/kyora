package handlers

import (
	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/web/views/pages"
	"github.com/abdelrahman146/kyora/internal/web/webcontext"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
)

type storeHandler struct {
	storeDomain *store.StoreDomain
}

func AddStoreRoutes(r *gin.Engine, storeDomain *store.StoreDomain) {
	h := &storeHandler{storeDomain}
	h.registerRoutes(r)
}

func (h *storeHandler) registerRoutes(r *gin.Engine) {
	r.GET("/", h.business)
}

func (h *storeHandler) business(c *gin.Context) {
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Business Settings",
		Description: "Manage your business settings",
		Keywords:    "business, settings, Kyora",
		Path:        "/settings/business",
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "/", Label: "Dashboard"},
			{Label: "Business Settings"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.SettingsBusiness())
}
