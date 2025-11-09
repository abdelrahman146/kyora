package bus

import (
	"context"
)

type Topic string

const OnboardingPaymentSucceededTopic Topic = "onboarding_payment_succeeded"

type OnboardingPaymentSucceededEvent struct {
	Ctx                  context.Context `json:"-"`
	OnboardingSessionID  string          `json:"onboardingSessionId"`
	StripeCheckoutID     string          `json:"stripeCheckoutId"`
	StripeSubscriptionID string          `json:"stripeSubscriptionId"`
}
