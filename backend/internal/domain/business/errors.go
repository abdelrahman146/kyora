package business

import (
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
)

// ErrBusinessNotFound indicates that a business could not be found in the authenticated workspace.
func ErrBusinessNotFound(businessDescriptor string, err error) error {
	return problem.NotFound("business not found").WithError(err).With("businessDescriptor", businessDescriptor)
}

// ErrBusinessDescriptorAlreadyTaken indicates the requested descriptor is already used within the workspace.
func ErrBusinessDescriptorAlreadyTaken(descriptor string, err error) error {
	return problem.Conflict("business descriptor is already taken").WithError(err).With("descriptor", descriptor)
}

// ErrInvalidBusinessDescriptor indicates the descriptor doesn't match allowed format.
func ErrInvalidBusinessDescriptor(descriptor string) error {
	return problem.BadRequest("invalid descriptor").With("descriptor", descriptor).With("hint", "use lowercase letters, numbers, and hyphens")
}

// ErrBusinessRateLimited indicates the client exceeded rate limits.
func ErrBusinessRateLimited() error {
	return problem.TooManyRequests("too many requests")
}
