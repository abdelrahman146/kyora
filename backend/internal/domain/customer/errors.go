package customer

import "github.com/abdelrahman146/kyora/internal/platform/types/problem"

// Customer errors

// ErrCustomerNotFound returns a not found error when a customer doesn't exist
func ErrCustomerNotFound(err error) *problem.Problem {
	return problem.NotFound("customer not found").WithError(err)
}

// ErrCustomerDuplicateEmail returns a conflict error when email already exists
func ErrCustomerDuplicateEmail(err error) *problem.Problem {
	return problem.Conflict("customer with this email already exists").WithError(err)
}

// ErrCustomerInvalidData returns a validation error for invalid customer data
func ErrCustomerInvalidData(message string) *problem.Problem {
	return problem.BadRequest(message)
}

// Customer address errors

// ErrCustomerAddressNotFound returns a not found error when an address doesn't exist
func ErrCustomerAddressNotFound(err error) *problem.Problem {
	return problem.NotFound("customer address not found").WithError(err)
}

// ErrCustomerAddressInvalidData returns a validation error for invalid address data
func ErrCustomerAddressInvalidData(message string) *problem.Problem {
	return problem.BadRequest(message)
}

// Customer note errors

// ErrCustomerNoteNotFound returns a not found error when a note doesn't exist
func ErrCustomerNoteNotFound(err error) *problem.Problem {
	return problem.NotFound("customer note not found").WithError(err)
}

// ErrCustomerNoteInvalidData returns a validation error for invalid note data
func ErrCustomerNoteInvalidData(message string) *problem.Problem {
	return problem.BadRequest(message)
}

// Authorization errors

// ErrCustomerUnauthorizedAccess returns a forbidden error for unauthorized customer access attempts
func ErrCustomerUnauthorizedAccess() *problem.Problem {
	return problem.Forbidden("you do not have permission to access this customer")
}

// ErrCustomerAddressUnauthorizedAccess returns a forbidden error for unauthorized address access attempts
func ErrCustomerAddressUnauthorizedAccess() *problem.Problem {
	return problem.Forbidden("you do not have permission to access this customer address")
}

// ErrCustomerNoteUnauthorizedAccess returns a forbidden error for unauthorized note access attempts
func ErrCustomerNoteUnauthorizedAccess() *problem.Problem {
	return problem.Forbidden("you do not have permission to access this customer note")
}
