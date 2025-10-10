package handlers

import (
	"github.com/abdelrahman146/kyora/internal/web/views/pages"
	"github.com/abdelrahman146/kyora/internal/web/webcontext"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
)

type InvoiceHandler struct {
}

func NewInvoiceHandler() *InvoiceHandler {
	return &InvoiceHandler{}
}

func (h *InvoiceHandler) RegisterRoutes(r gin.IRoutes) {
	r.GET("/invoices", h.Index)
	r.GET("/invoices/:id", h.Show)
}

func (h *InvoiceHandler) Index(c *gin.Context) {
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

func (h *InvoiceHandler) Show(c *gin.Context) {
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
