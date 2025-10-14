package handlers

import (
	"fmt"

	"github.com/abdelrahman146/kyora/internal/domain/expense"
	"github.com/abdelrahman146/kyora/internal/web/views/pages"
	"github.com/abdelrahman146/kyora/internal/web/webcontext"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
)

type expenseHandler struct {
	expenseDomain *expense.ExpenseDomain
}

func AddExpenseRoutes(r *gin.RouterGroup, expenseDomain *expense.ExpenseDomain) {
	h := &expenseHandler{expenseDomain: expenseDomain}
	h.registerRoutes(r)
}

func (h *expenseHandler) registerRoutes(c *gin.RouterGroup) {
	r := c.Group("/expenses")
	{
		r.GET("/", h.index)
		r.GET("/new", h.new)
		r.GET("/:id/edit", h.edit)
	}
}

func (h *expenseHandler) index(c *gin.Context) {
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

func (h *expenseHandler) new(c *gin.Context) {
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

func (h *expenseHandler) edit(c *gin.Context) {
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
