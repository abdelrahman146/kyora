package expense

import (
	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/store"
)

type ExpenseDomain struct {
	ExpenseService *ExpenseService
}

func SetupExpenseDomain(postgres *db.Postgres, atomicProcess *db.AtomicProcess, storeService *store.StoreService) *ExpenseDomain {
	expenseRepo := NewExpenseRepository(postgres)
	recurringExpenseRepo := NewRecurringExpenseRepository(postgres)
	postgres.AutoMigrate(&Expense{}, &RecurringExpense{})
	return &ExpenseDomain{
		ExpenseService: NewExpenseService(expenseRepo, recurringExpenseRepo, storeService, atomicProcess),
	}
}
