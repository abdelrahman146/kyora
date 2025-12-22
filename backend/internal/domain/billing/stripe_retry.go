package billing

import (
	"context"
	"math/rand"
	"time"

	stripelib "github.com/stripe/stripe-go/v83"
)

// withStripeRetry retries fn on transient Stripe API connection errors with exponential backoff and jitter
func withStripeRetry[T any](ctx context.Context, attempts int, fn func() (T, error)) (T, error) {
	var zero T
	if attempts < 1 {
		attempts = 1
	}
	base := 200 * time.Millisecond
	for i := 0; i < attempts; i++ {
		v, err := fn()
		if err == nil {
			return v, nil
		}
		// Only retry on APIConnection errors (network issues)
		if se, ok := err.(*stripelib.Error); ok {
			if se.Type != "api_connection_error" && se.Type != "rate_limit_error" {
				return zero, err
			}
		} else {
			return zero, err
		}
		// backoff
		delay := base * time.Duration(1<<i)
		// add jitter up to 100ms
		delay += time.Duration(rand.Intn(100)) * time.Millisecond
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return zero, err
		case <-timer.C:
		}
	}
	return zero, nil
}
