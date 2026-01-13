package storefront

import "github.com/abdelrahman146/kyora/internal/platform/types/problem"

func ErrStorefrontNotFound(storefrontPublicID string, err error) *problem.Problem {
	return problem.NotFound("storefront not found").WithError(err).With("storefrontPublicId", storefrontPublicID).WithCode("storefront.not_found")
}

func ErrStorefrontDisabled(storefrontPublicID string) *problem.Problem {
	return problem.Forbidden("storefront is disabled").With("storefrontPublicId", storefrontPublicID).WithCode("storefront.disabled")
}

func ErrIdempotencyKeyRequired() *problem.Problem {
	return problem.BadRequest("Idempotency-Key header is required").With("header", "Idempotency-Key").WithCode("storefront.idempotency_key_required")
}

func ErrIdempotencyConflict() *problem.Problem {
	return problem.Conflict("idempotency key already used with a different payload").WithCode("storefront.idempotency_conflict")
}

func ErrIdempotencyInProgress() *problem.Problem {
	return problem.Conflict("request already in progress").WithCode("storefront.idempotency_in_progress")
}
