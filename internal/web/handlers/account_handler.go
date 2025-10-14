package handlers

import (
	"net/http"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/web/middleware"
	"github.com/abdelrahman146/kyora/internal/web/views/pages"
	"github.com/abdelrahman146/kyora/internal/web/webcontext"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
)

type accountHandler struct {
	accountDomain *account.AccountDomain
}

func AddAccountRoutes(r *gin.RouterGroup, accountDomain *account.AccountDomain) {
	h := &accountHandler{
		accountDomain: accountDomain,
	}
	h.registerRoutes(r, accountDomain)
}

func (h *accountHandler) registerRoutes(c *gin.RouterGroup, accountDomain *account.AccountDomain) {
	r := c.Group("/accounts")
	{
		r.Use(middleware.AuthRequired, middleware.UserRequired(accountDomain.AuthService))
		r.GET("/profile", h.profile)
		r.GET("/manage", h.manage)
		r.POST("/invite", h.invite)
	}
}

func (h *accountHandler) profile(c *gin.Context) {
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Edit Profile",
		Description: "Edit your profile information",
		Keywords:    "edit profile, Kyora",
		Path:        "/accounts/profile",
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "/accounts", Label: "Accounts"},
			{Label: "Profile"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.ProfileEdit())
}

func (h *accountHandler) manage(c *gin.Context) {
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Accounts Management",
		Description: "Manage your team accounts",
		Keywords:    "manage accounts, Kyora",
		Path:        "/accounts/manage",
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "/accounts", Label: "Accounts"},
			{Label: "Accounts Management"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Redirect(c, "/dashboard")
}

func (h *accountHandler) invite(c *gin.Context) {
	c.String(http.StatusOK, "Not implemented")
}
