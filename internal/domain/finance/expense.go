package finance

import (
	"time"

	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/govalues/decimal"
	"gorm.io/gorm"
)

const (
	ExpenseTable  = "expenses"
	ExpenseAlias  = "exp"
	ExpenseStruct = "Expense"
)

type ExepenseType string

const (
	ExpenseTypeOneTime   ExepenseType = "one_time"
	ExpenseTypeRecurring ExepenseType = "recurring"
)

type Expense struct {
	gorm.Model
	ID              string          `gorm:"column:id;primaryKey;type:text" json:"id"`
	StoreID         string          `gorm:"column:store_id;type:text;not null;index" json:"storeId"`
	Name            string          `gorm:"column:name;type:text;not null" json:"name"`
	Amount          decimal.Decimal `gorm:"column:amount;type:numeric;not null;default:0" json:"amount"`
	Currency        string          `gorm:"column:currency;type:text;not null;default:'USD'" json:"currency"`
	Type            ExepenseType    `gorm:"column:type;type:text;not null;default:'one_time'" json:"type"`
	Notes           string          `gorm:"column:notes;type:text" json:"notes,omitempty"`
	RecurringPeriod time.Duration   `gorm:"column:recurring_period;type:interval" json:"recurringPeriod,omitempty"`
	RecurringUntil  *time.Time      `gorm:"column:recurring_until;type:timestamptz" json:"recurringUntil,omitempty"`
}

func (m *Expense) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = utils.ID.NewUlidWithPrefix(ExpenseAlias)
	}
	return
}

type CreateExpenseRequest struct {
	Name            string          `json:"name" binding:"required"`
	Amount          decimal.Decimal `json:"amount" binding:"required,gte=0"`
	Currency        string          `json:"currency" binding:"omitempty,len=3"`
	Type            ExepenseType    `json:"type" binding:"omitempty,oneof=one_time recurring"`
	Notes           string          `json:"notes" binding:"omitempty"`
	RecurringPeriod time.Duration   `json:"recurringPeriod" binding:"omitempty,gte=0,required_if=Type recurring"`
	RecurringUntil  *time.Time      `json:"recurringUntil" binding:"omitempty"`
}

type UpdateExpenseRequest struct {
	Name            string          `json:"name" binding:"omitempty"`
	Amount          decimal.Decimal `json:"amount" binding:"omitempty,gte=0"`
	Currency        string          `json:"currency" binding:"omitempty,len=3"`
	Type            ExepenseType    `json:"type" binding:"omitempty,oneof=one_time recurring"`
	Notes           string          `json:"notes" binding:"omitempty"`
	RecurringPeriod time.Duration   `json:"recurringPeriod" binding:"omitempty,gte=0,required_if=Type recurring"`
	RecurringUntil  *time.Time      `json:"recurringUntil" binding:"omitempty"`
}
