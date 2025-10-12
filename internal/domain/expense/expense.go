package expense

import (
	"database/sql"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/govalues/decimal"
	"gorm.io/gorm"
)

const (
	ExpenseTable  = "expenses"
	ExpenseStruct = "Expense"
	ExpenseAlias  = "exp"
)

type ExpenseType string

const (
	ExpenseTypeOneTime   ExpenseType = "one_time"
	ExpenseTypeRecurring ExpenseType = "recurring"
)

type ExpenseCategory string

const (
	ExpenseCategoryOffice        ExpenseCategory = "office"
	ExpenseCategoryTravel        ExpenseCategory = "travel"
	ExpenseCategorySupplies      ExpenseCategory = "supplies"
	ExpenseCategoryUtilities     ExpenseCategory = "utilities"
	ExpenseCategoryPayroll       ExpenseCategory = "payroll"
	ExpenseCategoryMarketing     ExpenseCategory = "marketing"
	ExpenseCategoryRent          ExpenseCategory = "rent"
	ExpenseCategorySoftware      ExpenseCategory = "software"
	ExpenseCategoryMaintenance   ExpenseCategory = "maintenance"
	ExpenseCategoryInsurance     ExpenseCategory = "insurance"
	ExpenseCategoryTaxes         ExpenseCategory = "taxes"
	ExpenseCategoryTraining      ExpenseCategory = "training"
	ExpenseCategoryConsulting    ExpenseCategory = "consulting"
	ExpenseCategoryMiscellaneous ExpenseCategory = "miscellaneous"
	ExpenseCategoryLegal         ExpenseCategory = "legal"
	ExpenseCategoryResearch      ExpenseCategory = "research"
	ExpenseCategoryEquipment     ExpenseCategory = "equipment"
	ExpenseCategoryShipping      ExpenseCategory = "shipping"
	ExpenseCategoryOther         ExpenseCategory = "other"
)

type Expense struct {
	gorm.Model
	ID                 string            `gorm:"column:id;primaryKey;type:text" json:"id"`
	StoreID            string            `gorm:"column:store_id;type:text;not null;index" json:"store_id"`
	Store              *store.Store      `gorm:"foreignKey:StoreID;references:ID" json:"store,omitempty"`
	RecurringExpenseID sql.NullString    `gorm:"column:recurring_expense_id;type:text;index" json:"recurringExpenseId"`
	RecurringExpense   *RecurringExpense `gorm:"foreignKey:RecurringExpenseID;references:ID" json:"recurringExpense,omitempty"`
	Amount             decimal.Decimal   `gorm:"column:amount;type:numeric;not null" json:"amount"`
	Currency           string            `gorm:"column:currency;type:text;not null;default:'USD'" json:"currency"`
	OccurredOn         time.Time         `gorm:"column:occurred_on;type:date;not null" json:"occurredOn"`
	Category           ExpenseCategory   `gorm:"column:category;type:text;not null;index" json:"category"`
	Type               ExpenseType       `gorm:"column:type;type:text;not null;index" json:"type"`
	Note               string            `gorm:"column:note;type:text" json:"note"`
}

func (m *Expense) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = utils.ID.NewUlidWithPrefix(ExpenseAlias)
	}
	return
}

type CreateExpenseRequest struct {
	Amount             decimal.Decimal `form:"amount" json:"amount" binding:"required"`
	Category           ExpenseCategory `form:"category" json:"category" binding:"required"`
	Type               ExpenseType     `form:"type" json:"type" binding:"required"`
	RecurringExpenseID sql.NullString  `form:"recurringExpenseId" json:"recurringExpenseId" binding:"omitempty,required_if=Type recurring"`
	Note               string          `form:"note" json:"note" binding:"omitempty"`
	OccurredOn         *time.Time      `form:"occurredOn" json:"occurredOn" binding:"omitempty"`
}

type UpdateExpenseRequest struct {
	Amount             decimal.Decimal `form:"amount" json:"amount" binding:"omitempty"`
	Category           ExpenseCategory `form:"category" json:"category" binding:"omitempty"`
	Type               ExpenseType     `form:"type" json:"type" binding:"omitempty"`
	RecurringExpenseID sql.NullString  `form:"recurringExpenseId" json:"recurringExpenseId" binding:"omitempty,required_if=Type recurring"`
	Note               string          `form:"note" json:"note" binding:"omitempty"`
	OccurredOn         *time.Time      `form:"occurredOn" json:"occurredOn" binding:"omitempty"`
}

type ExpenseFilter struct {
	IDs        []string          `form:"ids" json:"ids" binding:"omitempty,dive,required"`
	Categories []ExpenseCategory `form:"categories" json:"categories" binding:"omitempty"`
	Types      []ExpenseType     `form:"types" json:"types" binding:"omitempty"`
}
