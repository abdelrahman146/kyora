package handlers

import (
	"github.com/abdelrahman146/kyora/internal/web/views/pages"
	"github.com/abdelrahman146/kyora/internal/web/webcontext"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
)

type SupplierHandler struct {
}

func NewSupplierHandler() *SupplierHandler {
	return &SupplierHandler{}
}

func (h *SupplierHandler) RegisterRoutes(r gin.IRoutes) {
	r.GET("/suppliers", h.Index)
	r.GET("/suppliers/new", h.New)
	r.GET("/suppliers/:id/edit", h.Edit)
}

func (h *SupplierHandler) Index(c *gin.Context) {
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Suppliers",
		Description: "Manage your suppliers",
		Keywords:    "suppliers, Kyora",
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "/", Label: "Dashboard"},
			{Label: "Suppliers"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.SuppliersList())
}

func (h *SupplierHandler) New(c *gin.Context) {
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "New Supplier",
		Description: "Create a new supplier",
		Keywords:    "new supplier, Kyora",
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "/", Label: "Dashboard"},
			{Href: "/suppliers", Label: "Suppliers"},
			{Label: "New Supplier"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.SupplierForm(pages.SupplierFormProps{IsEdit: false}))
}

func (h *SupplierHandler) Edit(c *gin.Context) {
	supplierID := c.Param("id")
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Edit Supplier",
		Description: "Edit supplier " + supplierID,
		Keywords:    "edit supplier, Kyora",
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "/", Label: "Dashboard"},
			{Href: "/suppliers", Label: "Suppliers"},
			{Label: "Edit Supplier"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.SupplierForm(pages.SupplierFormProps{IsEdit: true}))
}
