package handlers

import (
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

func AddAccountRoutes(r *gin.Engine, accountDomain *account.AccountDomain) {
	h := &accountHandler{
		accountDomain: accountDomain,
	}
	h.registerRoutes(r, accountDomain)
}

func (h *accountHandler) registerRoutes(r *gin.Engine, accountDomain *account.AccountDomain) {
	r.Group("/accounts")
	{
		r.Use(middleware.AuthRequired, middleware.UserRequired(accountDomain.AuthService))
		r.GET("/", h.profile)
		r.GET("/new", h.new)
		r.GET("/:id/edit", h.edit)
	}
}

func (h *accountHandler) new(c *gin.Context) {
	c.String(200, "New Account")
}

func (h *accountHandler) edit(c *gin.Context) {
	id := c.Param("id")
	c.String(200, "Edit Account %s", id)
}

func (h *accountHandler) profile(c *gin.Context) {
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
