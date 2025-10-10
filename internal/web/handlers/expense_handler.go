package handlers

import (
	"fmt"

	"github.com/abdelrahman146/kyora/internal/web/views/pages"
	"github.com/abdelrahman146/kyora/internal/web/webcontext"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
)

type ExpenseHandler struct {
}

func NewExpenseHandler() *ExpenseHandler {
	return &ExpenseHandler{}
}

func (h *ExpenseHandler) RegisterRoutes(r gin.IRoutes) {
	r.GET("/expenses", h.Index)
	r.GET("/expenses/new", h.New)
	r.GET("/expenses/:id/edit", h.Edit)
}

func (h *ExpenseHandler) Index(c *gin.Context) {
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Expenses",
		Description: "Manage your expenses",
		Keywords:    "expenses, Kyora",
		Path:        "/expenses",
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "/", Label: "Dashboard"},
			{Label: "Expenses"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.ExpensesList())
}

func (h *ExpenseHandler) New(c *gin.Context) {
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "New Expense",
		Description: "Create a new expense",
		Keywords:    "new expense, Kyora",
		Path:        "/expenses/new",
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "/", Label: "Dashboard"},
			{Href: "/expenses", Label: "Expenses"},
			{Label: "New Expense"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.ExpenseForm(pages.ExpenseFormProps{IsEdit: false}))
}

func (h *ExpenseHandler) Edit(c *gin.Context) {
	expenseID := c.Param("id")
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Edit Expense",
		Description: "Edit expense " + expenseID,
		Keywords:    "edit expense, Kyora",
		Path:        fmt.Sprintf("/expenses/%s/edit", expenseID),
		Breadcrumbs: []webcontext.Breadcrumb{
			{Href: "/", Label: "Dashboard"},
			{Href: "/expenses", Label: "Expenses"},
			{Label: "Edit Expense"},
		},
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.ExpenseForm(pages.ExpenseFormProps{IsEdit: true}))
}
