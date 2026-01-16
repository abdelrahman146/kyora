package accounting

import (
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/utils/transformer"
	"github.com/shopspring/decimal"
)

// AssetResponse is the API response for Asset entity
// No DeletedAt field (GORM leakage removed)
type AssetResponse struct {
	ID          string          `json:"id"`
	BusinessID  string          `json:"businessId"`
	Name        string          `json:"name"`
	Type        AssetType       `json:"type"`
	Value       decimal.Decimal `json:"value"`
	Currency    string          `json:"currency"`
	PurchasedAt time.Time       `json:"purchasedAt"`
	Note        string          `json:"note"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

// ToAssetResponse converts Asset model to AssetResponse
func ToAssetResponse(a *Asset) AssetResponse {
	if a == nil {
		return AssetResponse{}
	}

	return AssetResponse{
		ID:          a.ID,
		BusinessID:  a.BusinessID,
		Name:        a.Name,
		Type:        a.Type,
		Value:       a.Value,
		Currency:    a.Currency,
		PurchasedAt: a.PurchasedAt,
		Note:        a.Note,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

// ToAssetResponses converts a slice of Asset models to responses
func ToAssetResponses(assets []*Asset) []AssetResponse {
	responses := make([]AssetResponse, len(assets))
	for i, a := range assets {
		responses[i] = ToAssetResponse(a)
	}
	return responses
}

// InvestmentResponse is the API response for Investment entity
// No DeletedAt field (GORM leakage removed)
type InvestmentResponse struct {
	ID         string          `json:"id"`
	BusinessID string          `json:"businessId"`
	InvestorID string          `json:"investorId"`
	Amount     decimal.Decimal `json:"amount"`
	Currency   string          `json:"currency"`
	Note       string          `json:"note"`
	InvestedAt time.Time       `json:"investedAt"`
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt"`
}

// ToInvestmentResponse converts Investment model to InvestmentResponse
func ToInvestmentResponse(inv *Investment) InvestmentResponse {
	if inv == nil {
		return InvestmentResponse{}
	}

	return InvestmentResponse{
		ID:         inv.ID,
		BusinessID: inv.BusinessID,
		InvestorID: inv.InvestorID,
		Amount:     inv.Amount,
		Currency:   inv.Currency,
		Note:       inv.Note,
		InvestedAt: inv.InvestedAt,
		CreatedAt:  inv.CreatedAt,
		UpdatedAt:  inv.UpdatedAt,
	}
}

// ToInvestmentResponses converts a slice of Investment models to responses
func ToInvestmentResponses(investments []*Investment) []InvestmentResponse {
	responses := make([]InvestmentResponse, len(investments))
	for i, inv := range investments {
		responses[i] = ToInvestmentResponse(inv)
	}
	return responses
}

// WithdrawalResponse is the API response for Withdrawal entity
// No DeletedAt field (GORM leakage removed)
type WithdrawalResponse struct {
	ID           string          `json:"id"`
	BusinessID   string          `json:"businessId"`
	Amount       decimal.Decimal `json:"amount"`
	Currency     string          `json:"currency"`
	WithdrawerID string          `json:"withdrawerId"`
	Note         string          `json:"note"`
	WithdrawnAt  time.Time       `json:"withdrawnAt"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
}

// ToWithdrawalResponse converts Withdrawal model to WithdrawalResponse
func ToWithdrawalResponse(w *Withdrawal) WithdrawalResponse {
	if w == nil {
		return WithdrawalResponse{}
	}

	return WithdrawalResponse{
		ID:           w.ID,
		BusinessID:   w.BusinessID,
		Amount:       w.Amount,
		Currency:     w.Currency,
		WithdrawerID: w.WithdrawerID,
		Note:         w.Note,
		WithdrawnAt:  w.WithdrawnAt,
		CreatedAt:    w.CreatedAt,
		UpdatedAt:    w.UpdatedAt,
	}
}

// ToWithdrawalResponses converts a slice of Withdrawal models to responses
func ToWithdrawalResponses(withdrawals []*Withdrawal) []WithdrawalResponse {
	responses := make([]WithdrawalResponse, len(withdrawals))
	for i, w := range withdrawals {
		responses[i] = ToWithdrawalResponse(w)
	}
	return responses
}

// ExpenseResponse is the API response for Expense entity
// No DeletedAt field (GORM leakage removed)
// Optional fields use pointers (orderId, recurringExpenseId, note)
type ExpenseResponse struct {
	ID                 string          `json:"id"`
	BusinessID         string          `json:"businessId"`
	OrderID            *string         `json:"orderId,omitempty"`
	RecurringExpenseID *string         `json:"recurringExpenseId,omitempty"`
	Amount             decimal.Decimal `json:"amount"`
	Currency           string          `json:"currency"`
	OccurredOn         time.Time       `json:"occurredOn"`
	Category           ExpenseCategory `json:"category"`
	Type               ExpenseType     `json:"type"`
	Note               *string         `json:"note,omitempty"`
	CreatedAt          time.Time       `json:"createdAt"`
	UpdatedAt          time.Time       `json:"updatedAt"`
}

// ToExpenseResponse converts Expense model to ExpenseResponse
func ToExpenseResponse(exp *Expense) ExpenseResponse {
	if exp == nil {
		return ExpenseResponse{}
	}

	return ExpenseResponse{
		ID:                 exp.ID,
		BusinessID:         exp.BusinessID,
		OrderID:            transformer.NullStringPtr(exp.OrderID),
		RecurringExpenseID: transformer.NullStringPtr(exp.RecurringExpenseID),
		Amount:             exp.Amount,
		Currency:           exp.Currency,
		OccurredOn:         exp.OccurredOn,
		Category:           exp.Category,
		Type:               exp.Type,
		Note:               transformer.NullStringPtr(exp.Note),
		CreatedAt:          exp.CreatedAt,
		UpdatedAt:          exp.UpdatedAt,
	}
}

// ToExpenseResponses converts a slice of Expense models to responses
func ToExpenseResponses(expenses []*Expense) []ExpenseResponse {
	responses := make([]ExpenseResponse, len(expenses))
	for i, exp := range expenses {
		responses[i] = ToExpenseResponse(exp)
	}
	return responses
}

// RecurringExpenseResponse is the API response for RecurringExpense entity
// No DeletedAt field (GORM leakage removed)
// Optional fields use pointers (recurringEndDate, note)
// Expenses field intentionally omitted (accessed via separate endpoint)
type RecurringExpenseResponse struct {
	ID                 string                    `json:"id"`
	BusinessID         string                    `json:"businessId"`
	Frequency          RecurringExpenseFrequency `json:"frequency"`
	RecurringEndDate   *time.Time                `json:"recurringEndDate,omitempty"`
	RecurringStartDate time.Time                 `json:"recurringStartDate"`
	NextRecurringDate  time.Time                 `json:"nextRecurringDate"`
	Amount             decimal.Decimal           `json:"amount"`
	Currency           string                    `json:"currency"`
	Category           ExpenseCategory           `json:"category"`
	Status             RecurringExpenseStatus    `json:"status"`
	Note               *string                   `json:"note,omitempty"`
	CreatedAt          time.Time                 `json:"createdAt"`
	UpdatedAt          time.Time                 `json:"updatedAt"`
}

// ToRecurringExpenseResponse converts RecurringExpense model to RecurringExpenseResponse
func ToRecurringExpenseResponse(rexp *RecurringExpense) RecurringExpenseResponse {
	if rexp == nil {
		return RecurringExpenseResponse{}
	}

	return RecurringExpenseResponse{
		ID:                 rexp.ID,
		BusinessID:         rexp.BusinessID,
		Frequency:          rexp.Frequency,
		RecurringEndDate:   transformer.NullTimePtr(rexp.RecurringEndDate),
		RecurringStartDate: rexp.RecurringStartDate,
		NextRecurringDate:  rexp.NextRecurringDate,
		Amount:             rexp.Amount,
		Currency:           rexp.Currency,
		Category:           rexp.Category,
		Status:             rexp.Status,
		Note:               transformer.NullStringPtr(rexp.Note),
		CreatedAt:          rexp.CreatedAt,
		UpdatedAt:          rexp.UpdatedAt,
	}
}

// ToRecurringExpenseResponses converts a slice of RecurringExpense models to responses
func ToRecurringExpenseResponses(recurringExpenses []*RecurringExpense) []RecurringExpenseResponse {
	responses := make([]RecurringExpenseResponse, len(recurringExpenses))
	for i, rexp := range recurringExpenses {
		responses[i] = ToRecurringExpenseResponse(rexp)
	}
	return responses
}

// =============================================================================
// Recent Activity Response
// =============================================================================

// RecentActivityType represents the type of recent accounting activity
type RecentActivityType string

const (
	RecentActivityTypeExpense    RecentActivityType = "expense"
	RecentActivityTypeInvestment RecentActivityType = "investment"
	RecentActivityTypeWithdrawal RecentActivityType = "withdrawal"
)

// RecentActivityResponse is a unified response for recent accounting activities
// It provides a polymorphic view of expenses, investments, and withdrawals
type RecentActivityResponse struct {
	ID          string              `json:"id"`
	Type        RecentActivityType  `json:"type"`
	Amount      decimal.Decimal     `json:"amount"`
	Currency    string              `json:"currency"`
	Description string              `json:"description"`
	OccurredAt  time.Time           `json:"occurredAt"`
	CreatedAt   time.Time           `json:"createdAt"`
	Category    *ExpenseCategory    `json:"category,omitempty"`    // Only for expenses
	ExpenseType *ExpenseType        `json:"expenseType,omitempty"` // Only for expenses
	PersonID    *string             `json:"personId,omitempty"`    // InvestorID or WithdrawerID
}

// RecentActivitiesResponse is the response for the recent activities endpoint
type RecentActivitiesResponse struct {
	Items []RecentActivityResponse `json:"items"`
}

// ExpenseToRecentActivity converts an Expense to RecentActivityResponse
func ExpenseToRecentActivity(exp *Expense) RecentActivityResponse {
	var description string
	if exp.Note.Valid {
		description = exp.Note.String
	} else {
		description = string(exp.Category)
	}

	return RecentActivityResponse{
		ID:          exp.ID,
		Type:        RecentActivityTypeExpense,
		Amount:      exp.Amount,
		Currency:    exp.Currency,
		Description: description,
		OccurredAt:  exp.OccurredOn,
		CreatedAt:   exp.CreatedAt,
		Category:    &exp.Category,
		ExpenseType: &exp.Type,
	}
}

// InvestmentToRecentActivity converts an Investment to RecentActivityResponse
func InvestmentToRecentActivity(inv *Investment) RecentActivityResponse {
	description := inv.Note
	if description == "" {
		description = "Investment"
	}

	return RecentActivityResponse{
		ID:          inv.ID,
		Type:        RecentActivityTypeInvestment,
		Amount:      inv.Amount,
		Currency:    inv.Currency,
		Description: description,
		OccurredAt:  inv.InvestedAt,
		CreatedAt:   inv.CreatedAt,
		PersonID:    &inv.InvestorID,
	}
}

// WithdrawalToRecentActivity converts a Withdrawal to RecentActivityResponse
func WithdrawalToRecentActivity(w *Withdrawal) RecentActivityResponse {
	description := w.Note
	if description == "" {
		description = "Withdrawal"
	}

	return RecentActivityResponse{
		ID:          w.ID,
		Type:        RecentActivityTypeWithdrawal,
		Amount:      w.Amount,
		Currency:    w.Currency,
		Description: description,
		OccurredAt:  w.WithdrawnAt,
		CreatedAt:   w.CreatedAt,
		PersonID:    &w.WithdrawerID,
	}
}
