package handlers

import (
	"fmt"

	"github.com/abdelrahman146/kyora/internal/web/views/pages"
	"github.com/abdelrahman146/kyora/internal/web/webcontext"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

func (h *ProductHandler) RegisterRoutes(r gin.IRoutes) {
	r.GET("/products", h.Index)
	r.GET("/products/new", h.New)
	r.GET("/products/:id/edit", h.Edit)
}

func (h *ProductHandler) Index(c *gin.Context) {
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Products",
		Description: "Manage your products",
		Keywords:    "products, Kyora",
		Path:        "/products",
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "/", Label: "Dashboard"},
			{Label: "Products"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.ProductsList())
}

func (h *ProductHandler) New(c *gin.Context) {
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "New Product",
		Description: "Create a new product",
		Keywords:    "new product, Kyora",
		Path:        "/products/new",
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "/", Label: "Dashboard"},
			{Href: "/products", Label: "Products"},
			{Label: "New Product"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.ProductForm(pages.ProductFormProps{IsEdit: false}))
}

func (h *ProductHandler) Edit(c *gin.Context) {
	productID := c.Param("id")
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Edit Product " + productID,
		Description: "Edit product " + productID,
		Keywords:    "edit product, Kyora",
		Path:        fmt.Sprintf("/products/%s/edit", productID),
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "/", Label: "Dashboard"},
			{Href: "/products", Label: "Products"},
			{Label: "Edit Product"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.ProductForm(pages.ProductFormProps{IsEdit: true, Product: &pages.Product{}}))
}
