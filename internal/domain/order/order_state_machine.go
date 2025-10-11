package order

import (
	"fmt"
	"slices"

	"github.com/abdelrahman146/kyora/internal/utils"
)

type OrderStateMachine struct {
	order *Order
}

func NewOrderStateMachine(order *Order) *OrderStateMachine {
	return &OrderStateMachine{order: order}
}

func (sm *OrderStateMachine) CanTransitionStateTo(newState OrderStatus) bool {
	allowedTransitions := map[OrderStatus][]OrderStatus{
		OrderStatusPending:   {OrderStatusPlaced, OrderStatusCancelled},
		OrderStatusPlaced:    {OrderStatusShipped, OrderStatusCancelled},
		OrderStatusShipped:   {OrderStatusFulfilled},
		OrderStatusFulfilled: {OrderStatusReturned},
		OrderStatusCancelled: {},
	}
	if allowed, ok := allowedTransitions[sm.order.Status]; ok {
		if slices.Contains(allowed, newState) {
			return true
		}
	}
	return false
}

func (sm *OrderStateMachine) TransitionStateTo(newState OrderStatus) error {
	if !sm.CanTransitionStateTo(newState) {
		return utils.Problem.UnprocessableEntity(fmt.Sprintf("cannot transition order status from %s to %s", sm.order.Status, newState)).With("currentState", sm.order.Status)
	}
	sm.order.Status = newState
	newState.UpdateTimestampField(sm.order)
	return nil
}

func (sm *OrderStateMachine) CanTransitionPaymentStatusTo(newStatus OrderPaymentStatus) bool {
	allowedTransitions := map[OrderPaymentStatus][]OrderPaymentStatus{
		OrderPaymentStatusPending:  {OrderPaymentStatusPaid, OrderPaymentStatusFailed},
		OrderPaymentStatusPaid:     {OrderPaymentStatusRefunded},
		OrderPaymentStatusFailed:   {OrderPaymentStatusPending},
		OrderPaymentStatusRefunded: {},
	}
	if allowed, ok := allowedTransitions[sm.order.PaymentStatus]; ok {
		if slices.Contains(allowed, newStatus) {
			return true
		}
	}
	return false
}

func (sm *OrderStateMachine) TransitionPaymentStatusTo(newStatus OrderPaymentStatus) error {
	if !sm.CanTransitionPaymentStatusTo(newStatus) {
		return utils.Problem.UnprocessableEntity(fmt.Sprintf("cannot transition payment status from %s to %s", sm.order.PaymentStatus, newStatus)).With("currentState", sm.order.PaymentStatus)
	}
	sm.order.PaymentStatus = newStatus
	newStatus.UpdateTimestampField(sm.order)
	return nil
}
