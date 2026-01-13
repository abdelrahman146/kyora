package inventory

import (
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
)

func ErrProductNotFound(err error) *problem.Problem {
	return problem.NotFound("product not found").WithError(err).WithCode("inventory.product_not_found")
}

// ErrVariantNotFound indicates that a variant could not be found.
func ErrVariantNotFound(err error) *problem.Problem {
	return problem.NotFound("variant not found").WithError(err).WithCode("inventory.variant_not_found")
}

// ErrCategoryNotFound indicates that a category could not be found.
func ErrCategoryNotFound(err error) *problem.Problem {
	return problem.NotFound("category not found").WithError(err).WithCode("inventory.category_not_found")
}
