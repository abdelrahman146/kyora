package onboarding

import (
	"math"
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
)

func ErrEmailAlreadyExists(err error) error {
	return problem.Conflict("email already registered").WithError(err).WithCode("onboarding.email_already_registered")
}

func ErrSessionTokenRequired(err error) error {
	return problem.BadRequest("session token is required").WithError(err).WithCode("onboarding.session_token_required")
}

func ErrActiveSessionExists(err error) error {
	return problem.Conflict("an onboarding session already exists for this email").WithError(err).WithCode("onboarding.session_already_exists")
}
func ErrSessionNotFound(err error) error {
	return problem.NotFound("onboarding session not found").WithError(err).WithCode("onboarding.session_not_found")
}
func ErrSessionExpired(err error) error {
	return problem.BadRequest("onboarding session expired").WithError(err).WithCode("onboarding.session_expired")
}

func ErrSessionAlreadyCommitted(err error) error {
	return problem.BadRequest("onboarding session already committed").WithError(err).WithCode("onboarding.session_already_committed")
}

func ErrInvalidStage(err error, expected string) error {
	return problem.BadRequest("invalid onboarding stage: expected " + expected).WithError(err).WithCode("onboarding.invalid_stage")
}
func ErrInvalidOTP(err error) error {
	return problem.BadRequest("invalid or expired verification code").WithError(err).WithCode("onboarding.invalid_otp")
}
func ErrPlanNotFound(err error) error {
	return problem.BadRequest("selected plan not found").WithError(err).WithCode("onboarding.plan_not_found")
}

func ErrSessionUpdateFailed(err error) error {
	p := problem.InternalError()
	p.Detail = "failed to update onboarding session"
	return p.WithError(err).WithCode("onboarding.session_update_failed")
}

func ErrSessionCleanupFailed(err error) error {
	p := problem.InternalError()
	p.Detail = "failed to cleanup onboarding sessions"
	return p.WithError(err).WithCode("onboarding.session_cleanup_failed")
}

// No special 402 helper; use BadRequest to indicate payment gating in the flow layer
func ErrPaymentRequired(err error) error {
	return problem.BadRequest("payment required").WithError(err).WithCode("onboarding.payment_required")
}
func ErrStripeOperation(err error) error {
	return problem.BadRequest("payment initialization failed").WithError(err).WithCode("onboarding.stripe_operation_failed")
}

// throttle error helpers
func ErrRateLimitedRetryAfter(_ error, retryAfter time.Duration) error {
	p := problem.TooManyRequests("OTP request rate limit exceeded")
	// Rate limit error
	seconds := int(math.Ceil(retryAfter.Seconds()))
	p.Extensions = map[string]interface{}{
		"retryAfterSeconds": seconds,
		"code":              "onboarding.otp_rate_limited",
	}
	return p
}

func ErrRateLimited(_ error) error {
	p := problem.TooManyRequests("OTP request rate limit exceeded")
	p.Detail = "too many failed OTP attempts, please restart onboarding"
	return p.WithCode("onboarding.otp_attempts_exceeded")
}
