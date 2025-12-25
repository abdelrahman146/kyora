package asset

import "github.com/abdelrahman146/kyora/internal/platform/types/problem"

func ErrAssetNotFound(id string, err error) error {
	return problem.NotFound("asset not found").WithError(err).With("assetId", id)
}

func ErrIdempotencyConflict() error {
	return problem.Conflict("idempotency key conflict")
}

func ErrIdempotencyInProgress() error {
	return problem.Conflict("request in progress")
}

func ErrUploadNotAllowed(field, reason string) error {
	return problem.BadRequest("upload not allowed").With("field", field).With("reason", reason)
}

func ErrUploadNotReady() error {
	return problem.Conflict("upload not completed")
}

func ErrRateLimited() error {
	return problem.TooManyRequests("too many requests")
}

func ErrAssetInUse(id string) error {
	return problem.Conflict("asset is referenced by other resources").With("assetId", id)
}
