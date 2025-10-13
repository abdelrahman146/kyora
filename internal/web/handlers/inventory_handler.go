package handlers

import (
	"fmt"

	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/web/views/pages"
	"github.com/abdelrahman146/kyora/internal/web/webcontext"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
)

type inventoryHandler struct {
	inventoryDomain *inventory.InventoryDomain
}

func AddInventoryRoutes(r *gin.Engine, inventoryDomain *inventory.InventoryDomain) {
	h := &inventoryHandler{
		inventoryDomain: inventoryDomain,
	}
	h.registerRoutes(r)
}

func (h *inventoryHandler) registerRoutes(r *gin.Engine) {
	r.Group("/inventory")
	{
		r.GET("/", h.index)
		r.GET("/new", h.new)
		r.GET("/:id/edit", h.edit)
	}
}

func (h *inventoryHandler) index(c *gin.Context) {
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

func (h *inventoryHandler) new(c *gin.Context) {
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

func (h *inventoryHandler) edit(c *gin.Context) {
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
