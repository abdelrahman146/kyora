package accounting

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/bus"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
	"github.com/shopspring/decimal"
)

type accountingRequiredBusinessService interface {
	GetEffectivePaymentMethodFee(ctx context.Context, businessID string, descriptor business.PaymentMethodDescriptor) (enabled bool, feePercent decimal.Decimal, feeFixed decimal.Decimal, err error)
}

// BusHandler listens for background events and triggers accounting automation.
type BusHandler struct {
	svc         *Service
	businessSvc accountingRequiredBusinessService
}

// NewBusHandler registers accounting listeners on the event bus.
func NewBusHandler(b *bus.Bus, svc *Service, businessSvc accountingRequiredBusinessService) {
	h := &BusHandler{svc: svc, businessSvc: businessSvc}
	b.Listen(bus.OrderPaymentSucceededTopic, h.HandleOrderPaymentSucceeded)
}

func (h *BusHandler) HandleOrderPaymentSucceeded(event any) {
	e, ok := event.(*bus.OrderPaymentSucceededEvent)
	if !ok {
		logger.FromContext(context.Background()).Error("invalid event type for OrderPaymentSucceededEvent")
		return
	}
	if h.svc == nil || h.businessSvc == nil {
		logger.FromContext(e.Ctx).Error("missing dependencies for accounting bus handler")
		return
	}
	if e.BusinessID == "" || e.OrderID == "" {
		logger.FromContext(e.Ctx).Error("missing required fields in OrderPaymentSucceededEvent", "businessId", e.BusinessID, "orderId", e.OrderID)
		return
	}
	pm := business.PaymentMethodDescriptor(e.PaymentMethod)
	enabled, feePercent, feeFixed, err := h.businessSvc.GetEffectivePaymentMethodFee(e.Ctx, e.BusinessID, pm)
	if err != nil {
		logger.FromContext(e.Ctx).Error("failed to resolve payment method fee", "error", err, "businessId", e.BusinessID, "paymentMethod", e.PaymentMethod)
		return
	}
	if !enabled {
		return
	}

	fee := e.OrderTotal.Mul(feePercent).Add(feeFixed)
	// No fee => no expense.
	if fee.LessThanOrEqual(decimal.Zero) {
		return
	}
	fee = fee.Round(2)

	if err := h.svc.UpsertTransactionFeeExpenseForOrder(e.Ctx, e.BusinessID, e.OrderID, fee, e.Currency, e.PaidAt, e.PaymentMethod); err != nil {
		logger.FromContext(e.Ctx).Error("failed to upsert transaction fee expense", "error", err, "businessId", e.BusinessID, "orderId", e.OrderID)
		return
	}
}
