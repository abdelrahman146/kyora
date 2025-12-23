package billing

import (
	"fmt"

	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/abdelrahman146/kyora/internal/platform/types/schema"
)

func ErrFeatureNotAvailable(err error, feature schema.Field) error {
	return problem.Forbidden(fmt.Sprintf("this feature is not available for your plan: %s", feature.JSONField())).With("feature", feature.JSONField()).WithError(err)
}

func ErrFeatureMaxLimitReached(err error, feature schema.Field, limit any) error {
	return problem.Forbidden(fmt.Sprintf("you have reached the maximum limit for this feature: %s", feature.JSONField())).With("feature", feature.JSONField()).With("limit", limit).WithError(err)
}

func ErrSubscriptionNotActive(err error) error {
	return problem.Forbidden("your subscription is not active. please renew your subscription").WithError(err)
}

func ErrSubscriptionCanceled(err error) error {
	return problem.Forbidden("your subscription has been canceled. please renew your subscription").WithError(err)
}

func ErrPlanNotFound(err error, descriptor string) error {
	return problem.NotFound("plan not found").With("descriptor", descriptor).WithError(err)
}

func ErrUnknownFeature(err error, feature any) error {
	return problem.BadRequest("the requested feature is unknown").With("feature", feature).WithError(err)
}

func ErrSubscriptionNotFound(err error, workspaceID string) error {
	return problem.NotFound("subscription not found for the given workspace").With("workspaceID", workspaceID).WithError(err)
}

func ErrInvalidPaymentMethod(err error) error {
	return problem.BadRequest("the provided payment method is invalid").WithError(err)
}

func ErrCannotCancelFreePlan(err error) error {
	return problem.BadRequest("cannot cancel a free plan subscription").WithError(err)
}

func ErrCannotChangeToSamePlan(err error) error {
	return problem.BadRequest("cannot change to the same plan").WithError(err)
}

func ErrCannotDowngradePlan(err error) error {
	return problem.BadRequest("cannot downgrade to a plan with fewer features or lower limits").WithError(err)
}

func ErrStripeOperationFailed(err error, operation string) error {
	return problem.InternalError().With("detail", "billing operation failed").With("operation", operation).WithError(err)
}

func ErrInvalidEffectiveDate(err error, effectiveDate string) error {
	return problem.BadRequest("invalid effective date").With("effectiveDate", effectiveDate).WithError(err)
}

func ErrInvalidProrationMode(err error, prorationMode string) error {
	return problem.BadRequest("invalid proration mode").With("prorationMode", prorationMode).WithError(err)
}

func ErrSubscriptionNotInTrial(err error, status string) error {
	return problem.BadRequest("subscription is not in trial").With("status", status).WithError(err)
}

func ErrSubscriptionNotPastDue(err error, status string) error {
	return problem.BadRequest("subscription is not past due").With("status", status).WithError(err)
}

// ErrCustomerCreationFailed indicates a failure creating a Stripe customer
func ErrCustomerCreationFailed(customerID string, err error) *problem.Problem {
	return problem.InternalError().With("detail",
		fmt.Sprintf("Failed to create Stripe customer for user %s", customerID),
	).WithError(err)
}

// ErrPlanSyncFailed indicates failure syncing plan with billing provider
func ErrPlanSyncFailed(planID string, err error) *problem.Problem {
	return problem.InternalError().With("detail", "failed to sync plan with billing provider").With("planId", planID).WithError(err)
}

func ErrWebhookProcessingFailed(err error, eventType string) error {
	return problem.InternalError().With("detail", "failed to process billing webhook").With("eventType", eventType).WithError(err)
}

func ErrWebhookSignatureInvalid(err error) error {
	return problem.BadRequest("invalid webhook signature").WithError(err)
}

func ErrWebhookPayloadInvalid(err error) error {
	return problem.BadRequest("invalid webhook payload").WithError(err)
}

func ErrCheckoutSessionFailed(err error) error {
	return problem.BadRequest("failed to create checkout session").WithError(err)
}

func ErrBillingPortalFailed(err error) error {
	return problem.BadRequest("failed to create billing portal session").WithError(err)
}

func ErrInvoiceNotReady(err error) error {
	return problem.Conflict("invoice is not ready for download").WithError(err)
}

func ErrUsageLimitExceeded(err error, feature string, current, limit int64) error {
	return problem.Forbidden("usage limit exceeded for feature").
		With("feature", feature).
		With("current", current).
		With("limit", limit).
		WithError(err)
}

func ErrUnauthorized() error {
	return problem.Unauthorized("authentication required")
}
