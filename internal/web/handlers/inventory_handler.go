package handlers

import (
	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/web/views/pages"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
)

type inventoryHandler struct {
	inventoryDomain *inventory.InventoryDomain
}

func AddInventoryRoutes(r *gin.RouterGroup, inventoryDomain *inventory.InventoryDomain) {
	h := &inventoryHandler{
		inventoryDomain: inventoryDomain,
	}
	h.registerRoutes(r)
}

func (h *inventoryHandler) registerRoutes(c *gin.RouterGroup) {
	r := c.Group("/inventory/products")
	{
		r.GET("/", h.index)
		r.POST("/", h.create)
		r.GET("/:id", h.show)
		r.PUT("/:id", h.update)
		r.DELETE("/:id", h.delete)
	}
	rv := c.Group("/inventory/variants")
	{
		rv.GET("/", h.index)
	}
}

func (h *inventoryHandler) index(c *gin.Context) {
	storeId := c.Param("storeId")
	page, pageSize, orderBy, isAscending := webutils.GetPaginationParams(c)
	_, err := h.inventoryDomain.InventoryService.ListProducts(c.Request.Context(), storeId, page, pageSize, orderBy, isAscending)
	if err != nil {
		webutils.Render(c, 500, pages.ErrorPage(500, "Failed to load products"))
		return
	}
	webutils.Render(c, 200, pages.NotImplemented("Products List"))
}

func (h *inventoryHandler) create(c *gin.Context) {
	_ = c.Param("storeId")
	// receive form data and validate inventory.CreateProductRequest
	c.String(200, "not implemented")
}

func (h *inventoryHandler) show(c *gin.Context) {
	storeId := c.Param("storeId")
	id := c.Param("id")
	_, err := h.inventoryDomain.InventoryService.GetProductByID(c.Request.Context(), storeId, id)
	if err != nil {
		webutils.Render(c, 500, pages.ErrorPage(500, "Failed to load product"))
		return
	}
	webutils.Render(c, 200, pages.NotImplemented("Product Details"))
}

func (h *inventoryHandler) update(c *gin.Context) {
	storeId := c.Param("storeId")
	id := c.Param("id")
	_ = storeId
	_ = id
	// receive form data and validate inventory.UpdateProductRequest
	c.String(200, "not implemented")
}

func (h *inventoryHandler) delete(c *gin.Context) {
	storeId := c.Param("storeId")
	id := c.Param("id")
	_ = storeId
	_ = id
	// perform delete operation
	c.String(200, "not implemented")
}

func (h *inventoryHandler) variantsIndex(c *gin.Context) {
	storeId := c.Param("storeId")
	page, pageSize, orderBy, isAscending := webutils.GetPaginationParams(c)
	_, err := h.inventoryDomain.InventoryService.ListVariants(c.Request.Context(), storeId, page, pageSize, orderBy, isAscending)
	if err != nil {
		webutils.Render(c, 500, pages.ErrorPage(500, "Failed to load variants"))
		return
	}
	webutils.Render(c, 200, pages.NotImplemented("Variants List"))
}
