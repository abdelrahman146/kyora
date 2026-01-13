package asset

import "github.com/abdelrahman146/kyora/internal/platform/types/problem"

func ErrAssetIdRequired() error {
	return problem.BadRequest("assetId is required").
		With("field", "assetId").
		WithCode("asset.id_required")
}

func ErrAssetNotAccessible() error {
	return problem.NotFound("asset not accessible").
		WithCode("asset.not_accessible")
}

func ErrAssetFileNotFound() error {
	return problem.NotFound("asset file not found").
		WithCode("asset.file_not_found")
}

func ErrAssetReadFailed(err error) error {
	return problem.InternalError().
		WithError(err).
		WithCode("asset.read_failed")
}

func ErrAssetWriteFailed(err error) error {
	return problem.InternalError().
		WithError(err).
		WithCode("asset.write_failed")
}

func ErrAssetDeleteFailed(err error) error {
	return problem.InternalError().
		WithError(err).
		WithCode("asset.delete_failed")
}

func ErrContentTypeRequired() error {
	return problem.BadRequest("Content-Type header is required").
		With("field", "Content-Type").
		WithCode("asset.content_type_required")
}
