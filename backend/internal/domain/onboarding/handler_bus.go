package onboarding

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/platform/bus"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
)

type BusHandler struct {
	svc *Service
}

func NewBusHandler(b *bus.Bus, svc *Service) {
	h := &BusHandler{svc: svc}
	b.Listen(bus.OnboardingPaymentSucceededTopic, h.HandleOnboardingPaymentSucceeded)
}

func (h *BusHandler) HandleOnboardingPaymentSucceeded(event any) {
	e, ok := event.(*bus.OnboardingPaymentSucceededEvent)
	if !ok {
		logger.FromContext(context.Background()).Error("invalid event type for OnboardingPaymentSucceededEvent")
		return
	}
	if err := h.svc.MarkPaymentSucceeded(e.Ctx, e.OnboardingSessionID, e.StripeSubscriptionID); err != nil {
		logger.FromContext(e.Ctx).Error("failed to mark onboarding session as payment succeeded", "error", err)
	}
}
