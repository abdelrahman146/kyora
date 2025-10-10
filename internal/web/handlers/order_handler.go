package handlers

import (
	"github.com/abdelrahman146/kyora/internal/web/views/pages"
	"github.com/abdelrahman146/kyora/internal/web/webcontext"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

func (h *OrderHandler) RegisterRoutes(r gin.IRoutes) {
	r.GET("/orders", h.Index)
	r.GET("/orders/new", h.New)
	r.GET("/orders/:id", h.Show)
	r.GET("/orders/:id/edit", h.Edit)
}

func (h *OrderHandler) Index(c *gin.Context) {
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Orders",
		Description: "Manage customer orders",
		Keywords:    "orders, Kyora",
		Path:        "/orders",
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "/", Label: "Dashboard"},
			{Label: "Orders"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.OrdersList())
}

func (h *OrderHandler) New(c *gin.Context) {
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "New Order",
		Description: "Create a new order",
		Keywords:    "new order, Kyora",
		Path:        "/orders/new",
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "/", Label: "Dashboard"},
			{Href: "/orders", Label: "Orders"},
			{Label: "New Order"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.OrderForm())
}

func (h *OrderHandler) Show(c *gin.Context) {
	orderID := c.Param("id")
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Order " + orderID,
		Description: "View and manage order " + orderID,
		Keywords:    "order, Kyora",
		Path:        "/orders/" + orderID,
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "/", Label: "Dashboard"},
			{Href: "/orders", Label: "Orders"},
			{Label: "Order " + orderID},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.OrderView(orderID))
}

func (h *OrderHandler) Edit(c *gin.Context) {
	orderID := c.Param("id")
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Edit Order " + orderID,
		Description: "Edit order " + orderID,
		Keywords:    "edit order, Kyora",
		Path:        "/orders/" + orderID + "/edit",
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "/", Label: "Dashboard"},
			{Href: "/orders", Label: "Orders"},
			{Label: "Edit Order"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.OrderForm())
}
