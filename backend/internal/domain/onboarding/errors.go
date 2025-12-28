package onboarding

import (
	"math"
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
)

func ErrEmailAlreadyExists(err error) error {
	return problem.Conflict("email already registered").WithError(err)
}

func ErrSessionTokenRequired(err error) error {
	return problem.BadRequest("session token is required").WithError(err)
}

func ErrActiveSessionExists(err error) error {
	return problem.Conflict("an onboarding session already exists for this email").WithError(err)
}
func ErrSessionNotFound(err error) error {
	return problem.NotFound("onboarding session not found").WithError(err)
}
func ErrSessionExpired(err error) error {
	return problem.BadRequest("onboarding session expired").WithError(err)
}

func ErrSessionAlreadyCommitted(err error) error {
	return problem.BadRequest("onboarding session already committed").WithError(err)
}

func ErrInvalidStage(err error, expected string) error {
	return problem.BadRequest("invalid onboarding stage: expected " + expected).WithError(err)
}
func ErrInvalidOTP(err error) error {
	return problem.BadRequest("invalid or expired verification code").WithError(err)
}
func ErrPlanNotFound(err error) error {
	return problem.BadRequest("selected plan not found").WithError(err)
}

func ErrSessionUpdateFailed(err error) error {
	p := problem.InternalError()
	p.Detail = "failed to update onboarding session"
	return p.WithError(err)
}

func ErrSessionCleanupFailed(err error) error {
	p := problem.InternalError()
	p.Detail = "failed to cleanup onboarding sessions"
	return p.WithError(err)
}

// No special 402 helper; use BadRequest to indicate payment gating in the flow layer
func ErrPaymentRequired(err error) error {
	return problem.BadRequest("payment required").WithError(err)
}
func ErrStripeOperation(err error) error {
	return problem.BadRequest("payment initialization failed").WithError(err)
}

// throttle error helpers
func ErrRateLimitedRetryAfter(_ error, retryAfter time.Duration) error {
	p := problem.TooManyRequests("OTP request rate limit exceeded")
	// Rate limit error
	seconds := int(math.Ceil(retryAfter.Seconds()))
	p.Extensions = map[string]interface{}{"retryAfterSeconds": seconds}
	return p
}

func ErrRateLimited(_ error) error {
	p := problem.TooManyRequests("OTP request rate limit exceeded")
	p.Detail = "too many failed OTP attempts, please restart onboarding"
	return p
}
