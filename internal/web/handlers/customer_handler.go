package handlers

import (
	"fmt"

	"github.com/abdelrahman146/kyora/internal/web/views/pages"
	"github.com/abdelrahman146/kyora/internal/web/webcontext"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
}

func NewCustomerHandler() *CustomerHandler {
	return &CustomerHandler{}
}

func (h *CustomerHandler) RegisterRoutes(r gin.IRoutes) {
	r.GET("/customers", h.Index)
	r.GET("/customers/new", h.New)
	r.GET("/customers/:id/edit", h.Edit)
}

func (h *CustomerHandler) Index(c *gin.Context) {
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

func (h *CustomerHandler) New(c *gin.Context) {
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

func (h *CustomerHandler) Edit(c *gin.Context) {
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
