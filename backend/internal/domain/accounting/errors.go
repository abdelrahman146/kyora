package accounting

import "github.com/abdelrahman146/kyora/internal/platform/types/problem"

// Asset errors

// ErrAssetNotFound returns a not found error for an asset
func ErrAssetNotFound(err error) *problem.Problem {
	return problem.NotFound("asset not found").WithError(err)
}

// ErrAssetInvalidType returns a validation error for invalid asset type
func ErrAssetInvalidType(assetType string) *problem.Problem {
	return problem.BadRequest("invalid asset type").With("type", assetType)
}

// Investment errors

// ErrInvestmentNotFound returns a not found error for an investment
func ErrInvestmentNotFound(err error) *problem.Problem {
	return problem.NotFound("investment not found").WithError(err)
}

// ErrInvestmentInvalidAmount returns a validation error for invalid investment amount
func ErrInvestmentInvalidAmount() *problem.Problem {
	return problem.BadRequest("investment amount must be greater than zero")
}

// ErrInvestorNotFound returns a not found error when investor user doesn't exist
func ErrInvestorNotFound(investorID string) *problem.Problem {
	return problem.NotFound("investor not found").With("investorId", investorID)
}

// Withdrawal errors

// ErrWithdrawalNotFound returns a not found error for a withdrawal
func ErrWithdrawalNotFound(err error) *problem.Problem {
	return problem.NotFound("withdrawal not found").WithError(err)
}

// ErrWithdrawalInvalidAmount returns a validation error for invalid withdrawal amount
func ErrWithdrawalInvalidAmount() *problem.Problem {
	return problem.BadRequest("withdrawal amount must be greater than zero")
}

// ErrWithdrawerNotFound returns a not found error when withdrawer user doesn't exist
func ErrWithdrawerNotFound(withdrawerID string) *problem.Problem {
	return problem.NotFound("withdrawer not found").With("withdrawerId", withdrawerID)
}

// ErrInsufficientFunds returns a validation error when there are insufficient funds for withdrawal
func ErrInsufficientFunds(requested, available string) *problem.Problem {
	return problem.BadRequest("insufficient funds for withdrawal").
		With("requested", requested).
		With("available", available)
}

// Expense errors

// ErrExpenseNotFound returns a not found error for an expense
func ErrExpenseNotFound(err error) *problem.Problem {
	return problem.NotFound("expense not found").WithError(err)
}

// ErrExpenseInvalidCategory returns a validation error for invalid expense category
func ErrExpenseInvalidCategory(category string) *problem.Problem {
	return problem.BadRequest("invalid expense category").With("category", category)
}

// ErrExpenseInvalidAmount returns a validation error for invalid expense amount
func ErrExpenseInvalidAmount() *problem.Problem {
	return problem.BadRequest("expense amount must be greater than zero")
}

// Recurring Expense errors

// ErrRecurringExpenseNotFound returns a not found error for a recurring expense
func ErrRecurringExpenseNotFound(err error) *problem.Problem {
	return problem.NotFound("recurring expense not found").WithError(err)
}

// ErrRecurringExpenseInvalidFrequency returns a validation error for invalid frequency
func ErrRecurringExpenseInvalidFrequency(frequency string) *problem.Problem {
	return problem.BadRequest("invalid recurring expense frequency").With("frequency", frequency)
}

// ErrRecurringExpenseInvalidStatus returns a validation error for invalid status
func ErrRecurringExpenseInvalidStatus(status string) *problem.Problem {
	return problem.BadRequest("invalid recurring expense status").With("status", status)
}

// ErrRecurringExpenseInvalidDateRange returns a validation error for invalid date range
func ErrRecurringExpenseInvalidDateRange() *problem.Problem {
	return problem.BadRequest("recurring end date must be after start date")
}

// ErrRecurringExpenseInvalidAmount returns a validation error for invalid recurring expense amount.
func ErrRecurringExpenseInvalidAmount() *problem.Problem {
	return problem.BadRequest("recurring expense amount must be greater than zero")
}

// ErrRecurringExpenseInvalidTransition returns a conflict error for invalid status transition
func ErrRecurringExpenseInvalidTransition(from, to string) *problem.Problem {
	return problem.Conflict("invalid status transition").
		With("from", from).
		With("to", to)
}
