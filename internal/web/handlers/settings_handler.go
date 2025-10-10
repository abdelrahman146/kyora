package handlers

import (
	"github.com/abdelrahman146/kyora/internal/web/views/pages"
	"github.com/abdelrahman146/kyora/internal/web/webcontext"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
)

type SettingsHandler struct {
}

func NewSettingsHandler() *SettingsHandler {
	return &SettingsHandler{}
}

func (h *SettingsHandler) RegisterRoutes(r gin.IRoutes) {
	r.GET("/settings/business", h.Business)
	r.GET("/profile", h.Profile)
}

func (h *SettingsHandler) Business(c *gin.Context) {
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

func (h *SettingsHandler) Profile(c *gin.Context) {
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Edit Profile",
		Description: "Edit your profile information",
		Keywords:    "edit profile, Kyora",
		Path:        "/profile",
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "/", Label: "Dashboard"},
			{Label: "Profile"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.ProfileEdit())
}
