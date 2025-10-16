package handlers

import (
	"github.com/abdelrahman146/kyora/internal/domain/expense"
	"github.com/abdelrahman146/kyora/internal/web/views/pages"
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
		r.POST("/", h.create)
		r.GET("/:id", h.show)      // View expense details
		r.PUT("/:id", h.update)    // Update expense details
		r.DELETE("/:id", h.delete) // Delete expense
	}
	rc := c.Group("/recurring-expenses")
	{
		rc.GET("/", h.recurringIndex)
		rc.POST("/", h.recurringCreate)
		rc.GET("/:id", h.recurringShow)
		rc.PUT("/:id", h.recurringUpdate)
		rc.DELETE("/:id", h.recurringDelete)
	}
}

func (h *expenseHandler) index(c *gin.Context) {
	storeId := c.Param("storeId")
	page, pageSize, orderBy, isAscending := webutils.GetPaginationParams(c)
	_, err := h.expenseDomain.ExpenseService.ListExpenses(c.Request.Context(), storeId, page, pageSize, orderBy, isAscending)
	if err != nil {
		webutils.Render(c, 500, pages.ErrorPage(500, "Failed to load expenses"))
		return
	}
	webutils.Render(c, 200, pages.NotImplemented("Expenses List"))
}

func (h *expenseHandler) create(c *gin.Context) {
	_ = c.Param("storeId")
	// receive form data and validate expense.CreateExpenseRequest
	c.String(200, "not implemented")
}

func (h *expenseHandler) show(c *gin.Context) {
	storeId := c.Param("storeId")
	id := c.Param("id")
	_, err := h.expenseDomain.ExpenseService.GetExpenseByID(c.Request.Context(), storeId, id)
	if err != nil {
		webutils.Render(c, 500, pages.ErrorPage(500, "Failed to load expense"))
		return
	}
	webutils.Render(c, 200, pages.NotImplemented("Expense Details"))
}

func (h *expenseHandler) update(c *gin.Context) {
	_ = c.Param("storeId")
	_ = c.Param("id")
	// receive form data and validate expense.UpdateExpenseRequest
	c.String(200, "not implemented")
}

func (h *expenseHandler) delete(c *gin.Context) {
	_ = c.Param("storeId")
	_ = c.Param("id")
	// perform delete operation
	c.String(200, "not implemented")
}

func (h *expenseHandler) recurringIndex(c *gin.Context) {
	storeId := c.Param("storeId")
	page, pageSize, orderBy, isAscending := webutils.GetPaginationParams(c)
	_, err := h.expenseDomain.ExpenseService.ListRecurringExpenses(c.Request.Context(), storeId, page, pageSize, orderBy, isAscending)
	if err != nil {
		webutils.Render(c, 500, pages.ErrorPage(500, "Failed to load recurring expenses"))
		return
	}
	webutils.Render(c, 200, pages.NotImplemented("Recurring Expenses List"))
}

func (h *expenseHandler) recurringCreate(c *gin.Context) {
	_ = c.Param("storeId")
	// receive form data and validate expense.CreateRecurringExpenseRequest
	c.String(200, "not implemented")
}

func (h *expenseHandler) recurringShow(c *gin.Context) {
	storeId := c.Param("storeId")
	id := c.Param("id")
	_, err := h.expenseDomain.ExpenseService.GetRecurringExpenseByID(c.Request.Context(), storeId, id)
	if err != nil {
		webutils.Render(c, 500, pages.ErrorPage(500, "Failed to load recurring expense"))
		return
	}
	webutils.Render(c, 200, pages.NotImplemented("Recurring Expense Details"))
}

func (h *expenseHandler) recurringUpdate(c *gin.Context) {
	_ = c.Param("storeId")
	_ = c.Param("id")
	// receive form data and validate expense.UpdateRecurringExpenseRequest
	c.String(200, "not implemented")
}

func (h *expenseHandler) recurringDelete(c *gin.Context) {
	_ = c.Param("storeId")
	_ = c.Param("id")
	// perform delete operation
	c.String(200, "not implemented")
}
