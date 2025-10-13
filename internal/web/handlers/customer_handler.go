package handlers

import (
	"fmt"

	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/web/views/pages"
	"github.com/abdelrahman146/kyora/internal/web/webcontext"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
)

type customerHandler struct {
	customerDomain *customer.CustomerDomain
}

func AddCustomerRoutes(r *gin.Engine, customerDomain *customer.CustomerDomain) {
	h := &customerHandler{customerDomain: customerDomain}
	h.registerRoutes(r)
}

func (h *customerHandler) registerRoutes(r *gin.Engine) {
	r.Group("/customers")
	{
		r.GET("/", h.index)
		r.GET("/new", h.new)
		r.GET("/:id/edit", h.edit)
	}
}

func (h *customerHandler) index(c *gin.Context) {
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Customers",
		Description: "Manage your customers",
		Keywords:    "customers, Kyora",
		Path:        "/customers",
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "/", Label: "Dashboard"},
			{Label: "Customers"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.CustomersList())
}

func (h *customerHandler) new(c *gin.Context) {
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "New Customer",
		Description: "Create a new customer",
		Keywords:    "new customer, Kyora",
		Path:        "/customers/new",
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "/", Label: "Dashboard"},
			{Href: "/customers", Label: "Customers"},
			{Label: "New Customer"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.CustomerForm(pages.CustomerFormProps{IsEdit: false}))
}

func (h *customerHandler) edit(c *gin.Context) {
	customerID := c.Param("id")
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Edit Customer",
		Description: "Edit customer " + customerID,
		Keywords:    "edit customer, Kyora",
		Path:        fmt.Sprintf("/customers/%s/edit", customerID),
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "/", Label: "Dashboard"},
			{Href: "/customers", Label: "Customers"},
			{Label: "Edit Customer"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.CustomerForm(pages.CustomerFormProps{IsEdit: true}))
}
