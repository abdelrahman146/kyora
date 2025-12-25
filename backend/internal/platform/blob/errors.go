package blob

import "github.com/abdelrahman146/kyora/internal/platform/types/problem"

func ErrProviderNotConfigured() error {
	return problem.InternalError().With("blob", "provider_not_configured")
}

func ErrObjectNotFound(key string) error {
	return problem.NotFound("object not found").With("key", key)
}
