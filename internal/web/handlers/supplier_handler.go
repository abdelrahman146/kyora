package handlers

import (
	"github.com/abdelrahman146/kyora/internal/domain/supplier"
	"github.com/abdelrahman146/kyora/internal/web/views/pages"
	"github.com/abdelrahman146/kyora/internal/web/webcontext"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
)

type supplierHandler struct {
	supplierDomain *supplier.SupplierDomain
}

func AddSupplierRoutes(r *gin.RouterGroup, supplierDomain *supplier.SupplierDomain) {
	h := &supplierHandler{supplierDomain: supplierDomain}
	h.RegisterRoutes(r)
}

func (h *supplierHandler) RegisterRoutes(c *gin.RouterGroup) {
	r := c.Group("/suppliers")
	{
		r.GET("/", h.Index)
		r.GET("/new", h.New)
		r.GET("/:id/edit", h.Edit)
	}
}

func (h *supplierHandler) Index(c *gin.Context) {
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

func (h *supplierHandler) New(c *gin.Context) {
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

func (h *supplierHandler) Edit(c *gin.Context) {
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
