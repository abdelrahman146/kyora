package expense

import (
	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/store"
)

type ExpenseDomain struct {
	ExpenseService *ExpenseService
}

func NewDomain(postgres *db.Postgres, atomicProcess *db.AtomicProcess, cache *db.Memcache, storeDomain *store.StoreDomain) *ExpenseDomain {
	expenseRepo := NewExpenseRepository(postgres)
	recurringExpenseRepo := NewRecurringExpenseRepository(postgres)
	postgres.AutoMigrate(&Expense{}, &RecurringExpense{})
	return &ExpenseDomain{
		ExpenseService: NewExpenseService(expenseRepo, recurringExpenseRepo, storeDomain.StoreService, atomicProcess),
	}
}
