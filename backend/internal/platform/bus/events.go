package bus

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type Topic string

const OnboardingPaymentSucceededTopic Topic = "onboarding_payment_succeeded"

// OrderPaymentSucceededTopic is emitted when an order transitions to payment status "paid".
// Consumers should treat it as best-effort and must be idempotent.
const OrderPaymentSucceededTopic Topic = "order_payment_succeeded"

type OnboardingPaymentSucceededEvent struct {
	Ctx                  context.Context `json:"-"`
	OnboardingSessionID  string          `json:"onboardingSessionId"`
	StripeCheckoutID     string          `json:"stripeCheckoutId"`
	StripeSubscriptionID string          `json:"stripeSubscriptionId"`
}

// OrderPaymentSucceededEvent is emitted when an order becomes paid.
// It includes the minimal order snapshot needed for downstream automation without requiring additional DB reads.
type OrderPaymentSucceededEvent struct {
	Ctx           context.Context `json:"-"`
	BusinessID    string          `json:"businessId"`
	OrderID       string          `json:"orderId"`
	PaymentMethod string          `json:"paymentMethod"`
	OrderTotal    decimal.Decimal `json:"orderTotal"`
	Currency      string          `json:"currency"`
	PaidAt        time.Time       `json:"paidAt"`
}
