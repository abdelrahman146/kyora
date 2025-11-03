package accounting

import (
	"database/sql"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/types/schema"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const (
	AssetTable  = "assets"
	AssetStruct = "Asset"
	AssetPrefix = "ast"
)

type AssetType string

const (
	AssetTypeSoftware  AssetType = "software"
	AssetTypeEquipment AssetType = "equipment"
	AssetTypeVehicle   AssetType = "vehicle"
	AssetTypeFurniture AssetType = "furniture"
	AssetTypeOther     AssetType = "other"
)

var AllAssetTypes = []AssetType{
	AssetTypeSoftware,
	AssetTypeEquipment,
	AssetTypeVehicle,
	AssetTypeFurniture,
	AssetTypeOther,
}

type Asset struct {
	gorm.Model
	ID          string             `gorm:"column:id;primaryKey;type:text" json:"id"`
	BusinessID  string             `gorm:"column:business_id;type:text;not null;index" json:"businessId"`
	Business    *business.Business `gorm:"foreignKey:BusinessID;references:ID" json:"business,omitempty"`
	Name        string             `gorm:"column:name;type:text;not null" json:"name"`
	Type        AssetType          `gorm:"column:type;type:text;not null" json:"type"`
	Value       decimal.Decimal    `gorm:"column:value;type:numeric;not null" json:"value"`
	Currency    string             `gorm:"column:currency;type:text;not null;default:'USD'" json:"currency"`
	PurchasedAt time.Time          `gorm:"column:purchased_at;type:date;not null" json:"purchasedAt"`
	Note        string             `gorm:"column:note;type:text" json:"note"`
}

func (m *Asset) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(AssetPrefix)
	}
	return nil
}

type CreateAssetRequest struct {
	Name        string          `form:"name" json:"name" binding:"required"`
	Type        AssetType       `form:"type" json:"type" binding:"required"`
	Value       decimal.Decimal `form:"value" json:"value" binding:"required"`
	PurchasedAt time.Time       `form:"purchasedAt" json:"purchasedAt" binding:"omitempty"`
	Note        string          `form:"note" json:"note" binding:"omitempty"`
}

type UpdateAssetRequest struct {
	Name        string          `form:"name" json:"name" binding:"omitempty"`
	Type        AssetType       `form:"type" json:"type" binding:"omitempty"`
	Value       decimal.Decimal `form:"value" json:"value" binding:"omitempty"`
	PurchasedAt time.Time       `form:"purchasedAt" json:"purchasedAt" binding:"omitempty"`
	Note        string          `form:"note" json:"note" binding:"omitempty"`
}

var AssetSchema = struct {
	ID          schema.Field
	BusinessID  schema.Field
	Name        schema.Field
	Type        schema.Field
	Value       schema.Field
	Currency    schema.Field
	PurchasedAt schema.Field
	Note        schema.Field
	CreatedAt   schema.Field
	UpdatedAt   schema.Field
	DeletedAt   schema.Field
}{
	ID:          schema.NewField("id", "id"),
	BusinessID:  schema.NewField("business_id", "businessId"),
	Name:        schema.NewField("name", "name"),
	Type:        schema.NewField("type", "type"),
	Value:       schema.NewField("value", "value"),
	Currency:    schema.NewField("currency", "currency"),
	PurchasedAt: schema.NewField("purchased_at", "purchasedAt"),
	Note:        schema.NewField("note", "note"),
	CreatedAt:   schema.NewField("created_at", "createdAt"),
	UpdatedAt:   schema.NewField("updated_at", "updatedAt"),
	DeletedAt:   schema.NewField("deleted_at", "deletedAt"),
}

const (
	InvestmentTable  = "investments"
	InvestmentStruct = "Investment"
	InvestmentPrefix = "inv"
)

type Investment struct {
	gorm.Model
	ID         string             `json:"id" gorm:"column:id;primaryKey;type:text"`
	BusinessID string             `json:"businessId" gorm:"column:business_id;type:text;not null;index"`
	Business   *business.Business `json:"business,omitempty" gorm:"foreignKey:BusinessID;references:ID"`
	InvestorID string             `json:"investorId" gorm:"column:investor_id;type:text;not null;index"`
	Investor   *account.User      `json:"investor,omitempty" gorm:"foreignKey:InvestorID;references:ID"`
	Amount     decimal.Decimal    `json:"amount" gorm:"column:amount;type:numeric;not null"`
	Currency   string             `json:"currency" gorm:"column:currency;type:text;not null;default:'USD'"`
	Note       string             `json:"note" gorm:"column:note;type:text"`
	InvestedAt time.Time          `json:"investedAt" gorm:"column:invested_at;type:timestamptz;not null;default:now()"`
}

func (m *Investment) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(InvestmentPrefix)
	}
	return
}

type CreateInvestmentRequest struct {
	InvestorID string          `form:"investorId" json:"investorId" binding:"required"`
	Amount     decimal.Decimal `form:"amount" json:"amount" binding:"required"`
	Note       string          `form:"note" json:"note" binding:"omitempty"`
	InvestedAt time.Time       `form:"investedAt" json:"investedAt" binding:"omitempty"`
}

type UpdateInvestmentRequest struct {
	InvestorID string          `form:"investorId" json:"investorId" binding:"omitempty"`
	Amount     decimal.Decimal `form:"amount" json:"amount" binding:"omitempty"`
	Note       string          `form:"note" json:"note" binding:"omitempty"`
	InvestedAt time.Time       `form:"investedAt" json:"investedAt" binding:"omitempty"`
}

var InvestmentSchema = struct {
	ID         schema.Field
	BusinessID schema.Field
	InvestorID schema.Field
	Amount     schema.Field
	Currency   schema.Field
	Note       schema.Field
	InvestedAt schema.Field
	CreatedAt  schema.Field
	UpdatedAt  schema.Field
	DeletedAt  schema.Field
}{
	ID:         schema.NewField("id", "id"),
	BusinessID: schema.NewField("business_id", "businessId"),
	InvestorID: schema.NewField("investor_id", "investorId"),
	Amount:     schema.NewField("amount", "amount"),
	Currency:   schema.NewField("currency", "currency"),
	Note:       schema.NewField("note", "note"),
	InvestedAt: schema.NewField("invested_at", "investedAt"),
	CreatedAt:  schema.NewField("created_at", "createdAt"),
	UpdatedAt:  schema.NewField("updated_at", "updatedAt"),
	DeletedAt:  schema.NewField("deleted_at", "deletedAt"),
}

const (
	WithdrawalTable  = "withdrawals"
	WithdrawalStruct = "Withdrawal"
	WithdrawalPrefix = "wdl"
)

type Withdrawal struct {
	gorm.Model
	ID           string             `json:"id" gorm:"column:id;primaryKey;type:text"`
	BusinessID   string             `json:"businessId" gorm:"column:business_id;type:text;not null;index"`
	Business     *business.Business `json:"business,omitempty" gorm:"foreignKey:BusinessID;references:ID"`
	Amount       decimal.Decimal    `json:"amount" gorm:"column:amount;type:numeric;not null"`
	Currency     string             `json:"currency" gorm:"column:currency;type:text;not null;default:'USD'"`
	WithdrawerID string             `json:"withdrawerId" gorm:"column:withdrawer_id;type:text;not null;index"`
	Withdrawer   *account.User      `json:"withdrawer,omitempty" gorm:"foreignKey:WithdrawerID;references:ID"`
	Note         string             `json:"note" gorm:"column:note;type:text"`
	WithdrawnAt  time.Time          `json:"withdrawnAt" gorm:"column:withdrawn_at;type:timestamptz;not null;default:now()"`
}

func (m *Withdrawal) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(WithdrawalPrefix)
	}
	return
}

type CreateWithdrawalRequest struct {
	Amount       decimal.Decimal `form:"amount" json:"amount" binding:"required"`
	WithdrawerID string          `form:"withdrawerId" json:"withdrawerId" binding:"required"`
	Note         string          `form:"note" json:"note" binding:"omitempty"`
	WithdrawnAt  time.Time       `form:"withdrawnAt" json:"withdrawnAt" binding:"omitempty"`
}

type UpdateWithdrawalRequest struct {
	Amount       decimal.Decimal `form:"amount" json:"amount" binding:"omitempty"`
	WithdrawerID string          `form:"withdrawerId" json:"withdrawerId" binding:"omitempty"`
	Note         string          `form:"note" json:"note" binding:"omitempty"`
	WithdrawnAt  time.Time       `form:"withdrawnAt" json:"withdrawnAt" binding:"omitempty"`
}

var WithdrawalSchema = struct {
	ID           schema.Field
	BusinessID   schema.Field
	Amount       schema.Field
	Currency     schema.Field
	WithdrawerID schema.Field
	Note         schema.Field
	WithdrawnAt  schema.Field
	CreatedAt    schema.Field
	UpdatedAt    schema.Field
	DeletedAt    schema.Field
}{
	ID:           schema.NewField("id", "id"),
	BusinessID:   schema.NewField("business_id", "businessId"),
	Amount:       schema.NewField("amount", "amount"),
	Currency:     schema.NewField("currency", "currency"),
	WithdrawerID: schema.NewField("withdrawer_id", "withdrawerId"),
	Note:         schema.NewField("note", "note"),
	WithdrawnAt:  schema.NewField("withdrawn_at", "withdrawnAt"),
	CreatedAt:    schema.NewField("created_at", "createdAt"),
	UpdatedAt:    schema.NewField("updated_at", "updatedAt"),
	DeletedAt:    schema.NewField("deleted_at", "deletedAt"),
}

/* Expense Model */
//---------------*/

const (
	ExpenseTable  = "expenses"
	ExpenseStruct = "Expense"
	ExpensePrefix = "exp"
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

func ExpenseCategoriesList() []ExpenseCategory {
	return []ExpenseCategory{
		ExpenseCategoryOffice,
		ExpenseCategoryTravel,
		ExpenseCategorySupplies,
		ExpenseCategoryUtilities,
		ExpenseCategoryPayroll,
		ExpenseCategoryMarketing,
		ExpenseCategoryRent,
		ExpenseCategorySoftware,
		ExpenseCategoryMaintenance,
		ExpenseCategoryInsurance,
		ExpenseCategoryTaxes,
		ExpenseCategoryTraining,
		ExpenseCategoryConsulting,
		ExpenseCategoryMiscellaneous,
		ExpenseCategoryLegal,
		ExpenseCategoryResearch,
		ExpenseCategoryEquipment,
		ExpenseCategoryShipping,
		ExpenseCategoryOther,
	}
}

type Expense struct {
	gorm.Model
	ID                 string             `gorm:"column:id;primaryKey;type:text" json:"id"`
	BusinessID         string             `gorm:"column:business_id;type:text;not null;index" json:"business_id"`
	Business           *business.Business `gorm:"foreignKey:BusinessID;references:ID" json:"business,omitempty"`
	RecurringExpenseID sql.NullString     `gorm:"column:recurring_expense_id;type:text;index" json:"recurringExpenseId"`
	RecurringExpense   *RecurringExpense  `gorm:"foreignKey:RecurringExpenseID;references:ID" json:"recurringExpense,omitempty"`
	Amount             decimal.Decimal    `gorm:"column:amount;type:numeric;not null" json:"amount"`
	Currency           string             `gorm:"column:currency;type:text;not null;default:'USD'" json:"currency"`
	OccurredOn         time.Time          `gorm:"column:occurred_on;type:date;not null;default:now()" json:"occurredOn"`
	Category           ExpenseCategory    `gorm:"column:category;type:text;not null;index" json:"category"`
	Type               ExpenseType        `gorm:"column:type;type:text;not null;index" json:"type"`
	Note               sql.NullString     `gorm:"column:note;type:text" json:"note"`
}

func (m *Expense) TableName() string {
	return ExpenseTable
}

func (m *Expense) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(ExpensePrefix)
	}
	return
}

type CreateExpenseRequest struct {
	Amount             decimal.Decimal `form:"amount" json:"amount" binding:"required"`
	Category           ExpenseCategory `form:"category" json:"category" binding:"required"`
	Type               ExpenseType     `form:"type" json:"type" binding:"required"`
	RecurringExpenseID string          `form:"recurringExpenseId" json:"recurringExpenseId" binding:"omitempty,required_if=Type recurring"`
	Note               string          `form:"note" json:"note" binding:"omitempty"`
	OccurredOn         *time.Time      `form:"occurredOn" json:"occurredOn" binding:"omitempty"`
}

type UpdateExpenseRequest struct {
	Amount             decimal.Decimal `form:"amount" json:"amount" binding:"omitempty"`
	Category           ExpenseCategory `form:"category" json:"category" binding:"omitempty"`
	Type               ExpenseType     `form:"type" json:"type" binding:"omitempty"`
	RecurringExpenseID string          `form:"recurringExpenseId" json:"recurringExpenseId" binding:"omitempty,required_if=Type recurring"`
	Note               string          `form:"note" json:"note" binding:"omitempty"`
	OccurredOn         *time.Time      `form:"occurredOn" json:"occurredOn" binding:"omitempty"`
}

var ExpenseSchema = struct {
	ID                 schema.Field
	BusinessID         schema.Field
	RecurringExpenseID schema.Field
	Amount             schema.Field
	Currency           schema.Field
	OccurredOn         schema.Field
	Category           schema.Field
	Type               schema.Field
	Note               schema.Field
	CreatedAt          schema.Field
	UpdatedAt          schema.Field
	DeletedAt          schema.Field
}{
	ID:                 schema.NewField("id", "id"),
	BusinessID:         schema.NewField("business_id", "businessId"),
	RecurringExpenseID: schema.NewField("recurring_expense_id", "recurringExpenseId"),
	Amount:             schema.NewField("amount", "amount"),
	Currency:           schema.NewField("currency", "currency"),
	OccurredOn:         schema.NewField("occurred_on", "occurredOn"),
	Category:           schema.NewField("category", "category"),
	Type:               schema.NewField("type", "type"),
	Note:               schema.NewField("note", "note"),
	CreatedAt:          schema.NewField("created_at", "createdAt"),
	UpdatedAt:          schema.NewField("updated_at", "updatedAt"),
	DeletedAt:          schema.NewField("deleted_at", "deletedAt"),
}

/* Recurring Expense Model */
//--------------------------------*/

const (
	RecurringExpenseTable  = "recurring_expenses"
	RecurringExpenseStruct = "RecurringExpense"
	RecurringExpensePrefix = "rexp"
)

type RecurringExpenseFrequency string

var (
	RecurringExpenseFrequencyDaily   RecurringExpenseFrequency = "daily"
	RecurringExpenseFrequencyWeekly  RecurringExpenseFrequency = "weekly"
	RecurringExpenseFrequencyMonthly RecurringExpenseFrequency = "monthly"
	RecurringExpenseFrequencyYearly  RecurringExpenseFrequency = "yearly"
)

func (f RecurringExpenseFrequency) GetNextRecurrenceDate(from time.Time) time.Time {
	switch f {
	case RecurringExpenseFrequencyDaily:
		return from.AddDate(0, 0, 1)
	case RecurringExpenseFrequencyWeekly:
		return from.AddDate(0, 0, 7)
	case RecurringExpenseFrequencyMonthly:
		return from.AddDate(0, 1, 0)
	case RecurringExpenseFrequencyYearly:
		return from.AddDate(1, 0, 0)
	default:
		return from
	}
}

func RecurringExpenseFrequencies() []RecurringExpenseFrequency {
	return []RecurringExpenseFrequency{
		RecurringExpenseFrequencyDaily,
		RecurringExpenseFrequencyWeekly,
		RecurringExpenseFrequencyMonthly,
		RecurringExpenseFrequencyYearly,
	}
}

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
	BusinessID         string                    `gorm:"column:business_id;type:text;not null;index" json:"businessId"`
	Business           *business.Business        `gorm:"foreignKey:BusinessID;references:ID" json:"business,omitempty"`
	Frequency          RecurringExpenseFrequency `gorm:"column:frequency;type:text;not null" json:"frequency"`
	RecurringEndDate   sql.NullTime              `gorm:"column:recurring_end_date;type:timestamp" json:"recurringEndDate"`
	RecurringStartDate time.Time                 `gorm:"column:recurring_start_date;type:timestamp;not null" json:"recurringStartDate"`
	NextRecurringDate  time.Time                 `gorm:"column:next_recurring_date;type:date;index" json:"nextRecurringDate"`
	Amount             decimal.Decimal           `gorm:"column:amount;type:numeric;not null" json:"amount"`
	Currency           string                    `gorm:"column:currency;type:text;not null;default:'USD'" json:"currency"`
	Category           ExpenseCategory           `gorm:"column:category;type:text;not null" json:"category"`
	Status             RecurringExpenseStatus    `gorm:"column:status;type:text;not null;default:'active'" json:"status"`
	Note               sql.NullString            `gorm:"column:note;type:text" json:"note"`
	Expenses           []*Expense                `gorm:"foreignKey:RecurringExpenseID;references:ID" json:"expenses,omitempty"`
}

func (m *RecurringExpense) TableName() string {
	return RecurringExpenseTable
}

func (m *RecurringExpense) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(RecurringExpensePrefix)
	}
	return
}

type CreateRecurringExpenseRequest struct {
	Frequency                    RecurringExpenseFrequency `form:"frequency" json:"frequency" binding:"required,oneof=daily weekly monthly yearly"`
	RecurringEndDate             time.Time                 `form:"recurringEndDate" json:"recurringEndDate" binding:"omitempty,gtfield=RecurringStartDate"`
	RecurringStartDate           time.Time                 `form:"recurringStartDate" json:"recurringStartDate" binding:"required"`
	Amount                       decimal.Decimal           `form:"amount" json:"amount" binding:"required,gt=0"`
	Category                     ExpenseCategory           `form:"category" json:"category" binding:"required,oneof=office travel supplies utilities payroll marketing rent software maintenance insurance taxes training consulting miscellaneous legal research equipment shipping other"`
	Note                         string                    `form:"note" json:"note" binding:"omitempty"`
	AutoCreateHistoricalExpenses bool                      `form:"autoCreateHistoricalExpenses" json:"autoCreateHistoricalExpenses" binding:"omitempty"`
}

type UpdateRecurringExpenseRequest struct {
	Frequency          RecurringExpenseFrequency `form:"frequency" json:"frequency" binding:"omitempty,oneof=daily weekly monthly yearly"`
	RecurringEndDate   time.Time                 `form:"recurringEndDate" json:"recurringEndDate" binding:"omitempty,gtfield=RecurringStartDate"`
	RecurringStartDate time.Time                 `form:"recurringStartDate" json:"recurringStartDate" binding:"omitempty"`
	Amount             decimal.Decimal           `form:"amount" json:"amount" binding:"omitempty,gt=0"`
	Category           ExpenseCategory           `form:"category" json:"category" binding:"omitempty"`
	Note               string                    `form:"note" json:"note" binding:"omitempty"`
}

var RecurringExpenseSchema = struct {
	ID                 schema.Field
	BusinessID         schema.Field
	Frequency          schema.Field
	RecurringEndDate   schema.Field
	RecurringStartDate schema.Field
	NextRecurringDate  schema.Field
	Amount             schema.Field
	Currency           schema.Field
	Category           schema.Field
	Status             schema.Field
	Note               schema.Field
	CreatedAt          schema.Field
	UpdatedAt          schema.Field
	DeletedAt          schema.Field
}{
	ID:                 schema.NewField("id", "id"),
	BusinessID:         schema.NewField("business_id", "businessId"),
	Frequency:          schema.NewField("frequency", "frequency"),
	RecurringEndDate:   schema.NewField("recurring_end_date", "recurringEndDate"),
	RecurringStartDate: schema.NewField("recurring_start_date", "recurringStartDate"),
	NextRecurringDate:  schema.NewField("next_recurring_date", "nextRecurringDate"),
	Amount:             schema.NewField("amount", "amount"),
	Currency:           schema.NewField("currency", "currency"),
	Category:           schema.NewField("category", "category"),
	Status:             schema.NewField("status", "status"),
	Note:               schema.NewField("note", "note"),
	CreatedAt:          schema.NewField("created_at", "createdAt"),
	UpdatedAt:          schema.NewField("updated_at", "updatedAt"),
	DeletedAt:          schema.NewField("deleted_at", "deletedAt"),
}
