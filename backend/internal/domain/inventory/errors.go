package inventory

import "github.com/abdelrahman146/kyora/internal/platform/types/problem"

// ErrProductNotFound returns a not found error when a product doesn't exist.
func ErrProductNotFound(err error) *problem.Problem {
	return problem.NotFound("product not found").WithError(err)
}

// ErrVariantNotFound returns a not found error when a variant doesn't exist.
func ErrVariantNotFound(err error) *problem.Problem {
	return problem.NotFound("variant not found").WithError(err)
}

// ErrCategoryNotFound returns a not found error when a category doesn't exist.
func ErrCategoryNotFound(err error) *problem.Problem {
	return problem.NotFound("category not found").WithError(err)
}
