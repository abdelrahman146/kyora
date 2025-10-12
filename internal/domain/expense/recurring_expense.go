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
	RecurringExpenseTable  = "recurring_expenses"
	RecurringExpenseStruct = "RecurringExpense"
	RecurringExpenseAlias  = "rexp"
)

type RecurringExpenseFrequency = sql.NullString

var (
	RecurringExpenseFrequencyDaily   RecurringExpenseFrequency = sql.NullString{String: "daily", Valid: true}
	RecurringExpenseFrequencyWeekly  RecurringExpenseFrequency = sql.NullString{String: "weekly", Valid: true}
	RecurringExpenseFrequencyMonthly RecurringExpenseFrequency = sql.NullString{String: "monthly", Valid: true}
	RecurringExpenseFrequencyYearly  RecurringExpenseFrequency = sql.NullString{String: "yearly", Valid: true}
)

type RecurringExpenseStatus string

const (
	RecurringExpenseStatusActive   RecurringExpenseStatus = "active"
	RecurringExpenseStatusPaused   RecurringExpenseStatus = "paused"
	RecurringExpenseStatusEnded    RecurringExpenseStatus = "ended"
	RecurringExpenseStatusCanceled RecurringExpenseStatus = "canceled"
)

type RecurringExpense struct {
	gorm.Model
	ID                 string                    `gorm:"column:id;primaryKey;type:text" json:"id"`
	StoreID            string                    `gorm:"column:store_id;type:text;not null;index" json:"storeId"`
	Store              *store.Store              `gorm:"foreignKey:StoreID;references:ID" json:"store,omitempty"`
	Frequency          RecurringExpenseFrequency `gorm:"column:frequency;type:text;not null" json:"frequency"`
	RecurringEndDate   sql.NullTime              `gorm:"column:recurring_end_date;type:timestamp" json:"recurringEndDate"`
	RecurringStartDate time.Time                 `gorm:"column:recurring_start_date;type:timestamp;not null" json:"recurringStartDate"`
	Amount             decimal.Decimal           `gorm:"column:amount;type:numeric;not null" json:"amount"`
	Category           ExpenseCategory           `gorm:"column:category;type:text;not null" json:"category"`
	Status             RecurringExpenseStatus    `gorm:"column:status;type:text;not null;default:'active'" json:"status"`
	Note               string                    `gorm:"column:note;type:text" json:"note"`
}

func (m *RecurringExpense) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = utils.ID.NewUlidWithPrefix(RecurringExpenseAlias)
	}
	return
}

type CreateRecurringExpenseRequest struct {
	Frequency          RecurringExpenseFrequency `form:"frequency" json:"frequency" binding:"required,oneof=daily weekly monthly yearly"`
	RecurringEndDate   sql.NullTime              `form:"recurringEndDate" json:"recurringEndDate" binding:"omitempty,gtfield=RecurringStartDate"`
	RecurringStartDate time.Time                 `form:"recurringStartDate" json:"recurringStartDate" binding:"required"`
	Amount             decimal.Decimal           `form:"amount" json:"amount" binding:"required,gt=0"`
	Category           ExpenseCategory           `form:"category" json:"category" binding:"required,oneof=office travel supplies utilities payroll marketing rent software maintenance insurance taxes training consulting miscellaneous legal research equipment shipping other"`
	Note               string                    `form:"note" json:"note" binding:"omitempty"`
}

type UpdateRecurringExpenseRequest struct {
	Frequency          RecurringExpenseFrequency `form:"frequency" json:"frequency" binding:"omitempty,oneof=daily weekly monthly yearly"`
	RecurringEndDate   sql.NullTime              `form:"recurringEndDate" json:"recurringEndDate" binding:"omitempty,gtfield=RecurringStartDate"`
	RecurringStartDate time.Time                 `form:"recurringStartDate" json:"recurringStartDate" binding:"omitempty"`
	Amount             decimal.Decimal           `form:"amount" json:"amount" binding:"omitempty,gt=0"`
	Category           ExpenseCategory           `form:"category" json:"category" binding:"omitempty"`
	Note               string                    `form:"note" json:"note" binding:"omitempty"`
}

type RecurringExpenseFilter struct {
	IDs         []string `form:"ids" json:"ids" binding:"omitempty,dive,required"`
	Categories  []string `form:"categories" json:"categories" binding:"omitempty,dive,required"`
	Frequencies []string `form:"frequencies" json:"frequencies" binding:"omitempty,dive,required,oneof=daily weekly monthly yearly"`
}
