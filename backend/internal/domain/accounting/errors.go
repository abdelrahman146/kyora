package accounting

import "github.com/abdelrahman146/kyora/internal/platform/types/problem"

// Asset errors

// ErrAssetNotFound returns a not found error for an asset
func ErrAssetNotFound(err error) *problem.Problem {
	return problem.NotFound("asset not found").WithError(err).WithCode("accounting.asset_not_found")
}

// ErrAssetInvalidType returns a validation error for invalid asset type
func ErrAssetInvalidType(assetType string) *problem.Problem {
	return problem.BadRequest("invalid asset type").With("type", assetType).WithCode("accounting.asset_invalid_type")
}

// Investment errors

// ErrInvestmentNotFound returns a not found error for an investment
func ErrInvestmentNotFound(err error) *problem.Problem {
	return problem.NotFound("investment not found").WithError(err).WithCode("accounting.investment_not_found")
}

// ErrInvestmentInvalidAmount returns a validation error for invalid investment amount
func ErrInvestmentInvalidAmount() *problem.Problem {
	return problem.BadRequest("investment amount must be greater than zero").WithCode("accounting.investment_invalid_amount")
}

// ErrInvestorNotFound returns a not found error when investor user doesn't exist
func ErrInvestorNotFound(investorID string) *problem.Problem {
	return problem.NotFound("investor not found").With("investorId", investorID).WithCode("accounting.investor_not_found")
}

// Withdrawal errors

// ErrWithdrawalNotFound returns a not found error for a withdrawal
func ErrWithdrawalNotFound(err error) *problem.Problem {
	return problem.NotFound("withdrawal not found").WithError(err).WithCode("accounting.withdrawal_not_found")
}

// ErrWithdrawalInvalidAmount returns a validation error for invalid withdrawal amount
func ErrWithdrawalInvalidAmount() *problem.Problem {
	return problem.BadRequest("withdrawal amount must be greater than zero").WithCode("accounting.withdrawal_invalid_amount")
}

// ErrWithdrawerNotFound returns a not found error when withdrawer user doesn't exist
func ErrWithdrawerNotFound(withdrawerID string) *problem.Problem {
	return problem.NotFound("withdrawer not found").With("withdrawerId", withdrawerID).WithCode("accounting.withdrawer_not_found")
}

// ErrInsufficientFunds returns a validation error when there are insufficient funds for withdrawal
func ErrInsufficientFunds(requested, available string) *problem.Problem {
	return problem.BadRequest("insufficient funds for withdrawal").
		With("requested", requested).
		With("available", available).
		WithCode("accounting.insufficient_funds")
}

// Expense errors

// ErrExpenseNotFound returns a not found error for an expense
func ErrExpenseNotFound(err error) *problem.Problem {
	return problem.NotFound("expense not found").WithError(err).WithCode("accounting.expense_not_found")
}

// ErrExpenseInvalidCategory returns a validation error for invalid expense category
func ErrExpenseInvalidCategory(category string) *problem.Problem {
	return problem.BadRequest("invalid expense category").With("category", category).WithCode("accounting.expense_invalid_category")
}

// ErrExpenseInvalidAmount returns a validation error for invalid expense amount
func ErrExpenseInvalidAmount() *problem.Problem {
	return problem.BadRequest("expense amount must be greater than zero").WithCode("accounting.expense_invalid_amount")
}

// Recurring Expense errors

// ErrRecurringExpenseNotFound returns a not found error for a recurring expense
func ErrRecurringExpenseNotFound(err error) *problem.Problem {
	return problem.NotFound("recurring expense not found").WithError(err).WithCode("accounting.recurring_expense_not_found")
}

// ErrRecurringExpenseInvalidFrequency returns a validation error for invalid frequency
func ErrRecurringExpenseInvalidFrequency(frequency string) *problem.Problem {
	return problem.BadRequest("invalid recurring expense frequency").With("frequency", frequency).WithCode("accounting.recurring_expense_invalid_frequency")
}

// ErrRecurringExpenseInvalidStatus returns a validation error for invalid status
func ErrRecurringExpenseInvalidStatus(status string) *problem.Problem {
	return problem.BadRequest("invalid recurring expense status").With("status", status).WithCode("accounting.recurring_expense_invalid_status")
}

// ErrRecurringExpenseInvalidDateRange returns a validation error for invalid date range
func ErrRecurringExpenseInvalidDateRange() *problem.Problem {
	return problem.BadRequest("recurring end date must be after start date").WithCode("accounting.recurring_expense_invalid_date_range")
}

// ErrRecurringExpenseInvalidAmount returns a validation error for invalid recurring expense amount.
func ErrRecurringExpenseInvalidAmount() *problem.Problem {
	return problem.BadRequest("recurring expense amount must be greater than zero").WithCode("accounting.recurring_expense_invalid_amount")
}

// ErrRecurringExpenseInvalidTransition returns a conflict error for invalid status transition
func ErrRecurringExpenseInvalidTransition(from, to string) *problem.Problem {
	return problem.Conflict("invalid status transition").
		With("from", from).
		With("to", to).
		WithCode("accounting.recurring_expense_invalid_transition")
}
