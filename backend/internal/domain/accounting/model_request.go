package accounting

import (
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/types/date"
	"github.com/shopspring/decimal"
)

// CreateAssetRequest is the request DTO for creating an asset.
type CreateAssetRequest struct {
	Name        string          `form:"name" json:"name" binding:"required"`
	Type        AssetType       `form:"type" json:"type" binding:"required"`
	Value       decimal.Decimal `form:"value" json:"value" binding:"required"`
	PurchasedAt time.Time       `form:"purchasedAt" json:"purchasedAt" binding:"omitempty"`
	Note        string          `form:"note" json:"note" binding:"omitempty"`
}

// UpdateAssetRequest is the request DTO for updating an asset.
type UpdateAssetRequest struct {
	Name        string          `form:"name" json:"name" binding:"omitempty"`
	Type        AssetType       `form:"type" json:"type" binding:"omitempty"`
	Value       decimal.Decimal `form:"value" json:"value" binding:"omitempty"`
	PurchasedAt time.Time       `form:"purchasedAt" json:"purchasedAt" binding:"omitempty"`
	Note        string          `form:"note" json:"note" binding:"omitempty"`
}

// CreateInvestmentRequest is the request DTO for creating an investment.
type CreateInvestmentRequest struct {
	InvestorID string          `form:"investorId" json:"investorId" binding:"required"`
	Amount     decimal.Decimal `form:"amount" json:"amount" binding:"required"`
	Note       string          `form:"note" json:"note" binding:"omitempty"`
	InvestedAt time.Time       `form:"investedAt" json:"investedAt" binding:"omitempty"`
}

// UpdateInvestmentRequest is the request DTO for updating an investment.
type UpdateInvestmentRequest struct {
	InvestorID string          `form:"investorId" json:"investorId" binding:"omitempty"`
	Amount     decimal.Decimal `form:"amount" json:"amount" binding:"omitempty"`
	Note       string          `form:"note" json:"note" binding:"omitempty"`
	InvestedAt time.Time       `form:"investedAt" json:"investedAt" binding:"omitempty"`
}

// CreateWithdrawalRequest is the request DTO for creating a withdrawal.
type CreateWithdrawalRequest struct {
	Amount       decimal.Decimal `form:"amount" json:"amount" binding:"required"`
	WithdrawerID string          `form:"withdrawerId" json:"withdrawerId" binding:"required"`
	Note         string          `form:"note" json:"note" binding:"omitempty"`
	WithdrawnAt  time.Time       `form:"withdrawnAt" json:"withdrawnAt" binding:"omitempty"`
}

// UpdateWithdrawalRequest is the request DTO for updating a withdrawal.
type UpdateWithdrawalRequest struct {
	Amount       decimal.Decimal `form:"amount" json:"amount" binding:"omitempty"`
	WithdrawerID string          `form:"withdrawerId" json:"withdrawerId" binding:"omitempty"`
	Note         string          `form:"note" json:"note" binding:"omitempty"`
	WithdrawnAt  time.Time       `form:"withdrawnAt" json:"withdrawnAt" binding:"omitempty"`
}

// CreateExpenseRequest is the request DTO for creating an expense.
type CreateExpenseRequest struct {
	Amount             decimal.Decimal `form:"amount" json:"amount" binding:"required"`
	Category           ExpenseCategory `form:"category" json:"category" binding:"required"`
	Type               ExpenseType     `form:"type" json:"type" binding:"required"`
	RecurringExpenseID string          `form:"recurringExpenseId" json:"recurringExpenseId" binding:"omitempty,required_if=Type recurring"`
	Note               string          `form:"note" json:"note" binding:"omitempty"`
	OccurredOn         *date.Date      `form:"occurredOn" json:"occurredOn" binding:"omitempty"`
}

// UpdateExpenseRequest is the request DTO for updating an expense.
type UpdateExpenseRequest struct {
	Amount             decimal.Decimal `form:"amount" json:"amount" binding:"omitempty"`
	Category           ExpenseCategory `form:"category" json:"category" binding:"omitempty"`
	Type               ExpenseType     `form:"type" json:"type" binding:"omitempty"`
	RecurringExpenseID string          `form:"recurringExpenseId" json:"recurringExpenseId" binding:"omitempty,required_if=Type recurring"`
	Note               string          `form:"note" json:"note" binding:"omitempty"`
	OccurredOn         *date.Date      `form:"occurredOn" json:"occurredOn" binding:"omitempty"`
}

// CreateRecurringExpenseRequest is the request DTO for creating a recurring expense.
type CreateRecurringExpenseRequest struct {
	Frequency                    RecurringExpenseFrequency `form:"frequency" json:"frequency" binding:"required,oneof=daily weekly monthly yearly"`
	RecurringEndDate             *date.Date                `form:"recurringEndDate" json:"recurringEndDate" binding:"omitempty"`
	RecurringStartDate           date.Date                 `form:"recurringStartDate" json:"recurringStartDate" binding:"required"`
	Amount                       decimal.Decimal           `form:"amount" json:"amount" binding:"required"`
	Category                     ExpenseCategory           `form:"category" json:"category" binding:"required,oneof=office travel supplies utilities payroll marketing rent software maintenance insurance taxes training consulting miscellaneous legal research equipment shipping transaction_fee other"`
	Note                         string                    `form:"note" json:"note" binding:"omitempty"`
	AutoCreateHistoricalExpenses bool                      `form:"autoCreateHistoricalExpenses" json:"autoCreateHistoricalExpenses" binding:"omitempty"`
}

// UpdateRecurringExpenseRequest is the request DTO for updating a recurring expense.
type UpdateRecurringExpenseRequest struct {
	Frequency          RecurringExpenseFrequency `form:"frequency" json:"frequency" binding:"omitempty,oneof=daily weekly monthly yearly"`
	RecurringEndDate   *date.Date                `form:"recurringEndDate" json:"recurringEndDate" binding:"omitempty"`
	RecurringStartDate *date.Date                `form:"recurringStartDate" json:"recurringStartDate" binding:"omitempty"`
	Amount             decimal.Decimal           `form:"amount" json:"amount" binding:"omitempty"`
	Category           ExpenseCategory           `form:"category" json:"category" binding:"omitempty"`
	Note               string                    `form:"note" json:"note" binding:"omitempty"`
}

// Query and handler request types

// listAssetsQuery represents the query parameters for listing assets.
type listAssetsQuery struct {
	Page     int      `form:"page" binding:"omitempty,min=1"`
	PageSize int      `form:"pageSize" binding:"omitempty,min=1,max=100"`
	OrderBy  []string `form:"orderBy" binding:"omitempty"`
}

// listExpensesQuery represents the query parameters for listing expenses.
type listExpensesQuery struct {
	Page     int             `form:"page" binding:"omitempty,min=1"`
	PageSize int             `form:"pageSize" binding:"omitempty,min=1,max=100"`
	OrderBy  []string        `form:"orderBy" binding:"omitempty"`
	Category ExpenseCategory `form:"category" binding:"omitempty"`
	From     *time.Time      `form:"from" binding:"omitempty" time_format:"2006-01-02"`
	To       *time.Time      `form:"to" binding:"omitempty" time_format:"2006-01-02"`
}

// updateRecurringExpenseStatusRequest represents the request to update recurring expense status.
type updateRecurringExpenseStatusRequest struct {
	Status RecurringExpenseStatus `json:"status" binding:"required,oneof=active paused ended canceled"`
}

// summaryQuery represents the query parameters for accounting summary.
type summaryQuery struct {
	From string `form:"from" binding:"omitempty"`
	To   string `form:"to" binding:"omitempty"`
}

// recentActivitiesQuery represents the query parameters for recent activities.
type recentActivitiesQuery struct {
	Limit int `form:"limit" binding:"omitempty,min=1,max=50"`
}
