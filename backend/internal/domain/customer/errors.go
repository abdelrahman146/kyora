package customer

import (
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
)

// Customer errors
func ErrCustomerNotFound(err error) *problem.Problem {
	return problem.NotFound("customer not found").WithError(err).WithCode("customer.not_found")
}

func ErrCustomerDuplicateEmail(err error) *problem.Problem {
	return problem.Conflict("customer with this email already exists").WithError(err).WithCode("customer.duplicate_email")
}

func ErrCustomerInvalidData(message string) *problem.Problem {
	return problem.BadRequest(message).WithCode("customer.invalid_data")
}

// Customer address errors

func ErrCustomerAddressNotFound(err error) *problem.Problem {
	return problem.NotFound("customer address not found").WithError(err).WithCode("customer.address_not_found")
}

func ErrCustomerAddressInvalidData(message string) *problem.Problem {
	return problem.BadRequest(message).WithCode("customer.address_invalid_data")
}

// Customer note errors

func ErrCustomerNoteNotFound(err error) *problem.Problem {
	return problem.NotFound("customer note not found").WithError(err).WithCode("customer.note_not_found")
}

func ErrCustomerNoteInvalidData(message string) *problem.Problem {
	return problem.BadRequest(message).WithCode("customer.note_invalid_data")
}

// Authorization errors

func ErrCustomerUnauthorizedAccess() *problem.Problem {
	return problem.Forbidden("you do not have permission to access this customer").WithCode("customer.unauthorized")
}

func ErrCustomerAddressUnauthorizedAccess() *problem.Problem {
	return problem.Forbidden("you do not have permission to access this customer address").WithCode("customer.address_unauthorized")
}

func ErrCustomerNoteUnauthorizedAccess() *problem.Problem {
	return problem.Forbidden("you do not have permission to access this customer note").WithCode("customer.note_unauthorized")
}

// Handler/service inline errors

func ErrCustomerIdRequired() *problem.Problem {
	return problem.BadRequest("customerId is required").With("field", "customerId").WithCode("customer.id_required")
}

func ErrCustomerInvalidQueryParams(err error) *problem.Problem {
	return problem.BadRequest("invalid query parameters").WithError(err).WithCode("customer.invalid_query_params")
}

func ErrCustomerInvalidSearchTerm() *problem.Problem {
	return problem.BadRequest("invalid search term").WithCode("customer.invalid_search_term")
}

func ErrCustomerQueryFailed(err error) *problem.Problem {
	return problem.InternalError().WithError(err).WithCode("customer.query_failed")
}

func ErrCustomerAddressIdRequired() *problem.Problem {
	return problem.BadRequest("customerId and addressId are required").With("field", "addressId").WithCode("customer.address_id_required")
}

func ErrCustomerNoteIdRequired() *problem.Problem {
	return problem.BadRequest("customerId and noteId are required").With("field", "noteId").WithCode("customer.note_id_required")
}

func ErrBusinessRequired() *problem.Problem {
	return problem.InternalError().With("reason", "business is required").WithCode("customer.business_required")
}

func ErrCustomerDataRequired() *problem.Problem {
	return problem.BadRequest("customer is required").WithCode("customer.data_required")
}

func ErrCustomerEmailRequired() *problem.Problem {
	return problem.BadRequest("email is required").With("field", "email").WithCode("customer.email_required")
}

func ErrCustomerNameRequired() *problem.Problem {
	return problem.BadRequest("name is required").With("field", "name").WithCode("customer.name_required")
}
