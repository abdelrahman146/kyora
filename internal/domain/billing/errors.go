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
	return problem.BadRequest("cannot downgrade to a plan with fewer features").WithError(err)
}
