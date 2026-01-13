package blob

import "github.com/abdelrahman146/kyora/internal/platform/types/problem"

func ErrProviderNotConfigured() error {
	return problem.InternalError().With("blob", "provider_not_configured").WithCode("blob.provider_not_configured")
}

func ErrBlobObjectNotFound(key string) error {
	return problem.NotFound("object not found").With("key", key).WithCode("blob.object_not_found")
}
