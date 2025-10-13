package handlers

import (
	"github.com/abdelrahman146/kyora/internal/domain/order"
	"github.com/abdelrahman146/kyora/internal/web/views/pages"
	"github.com/abdelrahman146/kyora/internal/web/webcontext"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
)

type invoiceHandler struct {
	orderDomain *order.OrderDomain
}

func AddInvoiceRoutes(r *gin.Engine, orderDomain *order.OrderDomain) {
	h := &invoiceHandler{orderDomain: orderDomain}
	h.registerRoutes(r)
}

func (h *invoiceHandler) registerRoutes(r *gin.Engine) {
	r.Group("/invoices")
	{
		r.GET("/", h.index)
		r.GET("/:id", h.show)
	}
}

func (h *invoiceHandler) index(c *gin.Context) {
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Payments & Invoices",
		Description: "Manage customer invoices",
		Keywords:    "invoices, Kyora",
		Path:        "/invoices",
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "/", Label: "Dashboard"},
			{Label: "Payments & Invoices"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.InvoicesList())
}

func (h *invoiceHandler) show(c *gin.Context) {
	invoiceID := c.Param("id")
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Invoice " + invoiceID,
		Description: "View and manage invoice " + invoiceID,
		Keywords:    "Invoice, Kyora",
		Path:        "/invoices/" + invoiceID,
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "/", Label: "Dashboard"},
			{Href: "/invoices", Label: "Payments & Invoices"},
			{Label: "Invoice " + invoiceID},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.InvoiceView(invoiceID))
}
