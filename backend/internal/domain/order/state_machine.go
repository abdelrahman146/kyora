package order

import (
	"slices"
)

type orderStateMachine struct {
	order *Order
}

func newOrderStateMachine(order *Order) *orderStateMachine {
	return &orderStateMachine{order: order}
}

func (sm *orderStateMachine) canTransitionStateTo(newState OrderStatus) bool {
	allowedTransitions := map[OrderStatus][]OrderStatus{
		OrderStatusPending: {OrderStatusPlaced, OrderStatusCancelled},
		// ReadyForShipment is an optional intermediate state between placed and shipped.
		// We keep placed->shipped for backward compatibility.
		OrderStatusPlaced:           {OrderStatusReadyForShipment, OrderStatusShipped, OrderStatusCancelled},
		OrderStatusReadyForShipment: {OrderStatusShipped, OrderStatusCancelled},
		OrderStatusShipped:          {OrderStatusFulfilled},
		OrderStatusFulfilled:        {OrderStatusReturned},
		OrderStatusCancelled:        {},
		OrderStatusReturned:         {},
	}
	if allowed, ok := allowedTransitions[sm.order.Status]; ok {
		if slices.Contains(allowed, newState) {
			return true
		}
	}
	return false
}

func (sm *orderStateMachine) transitionStateTo(newState OrderStatus) error {
	if !sm.canTransitionStateTo(newState) {
		return ErrOrderStatusUpdateNotAllowed(sm.order.ID, sm.order.Status, newState)
	}
	sm.order.Status = newState
	newState.UpdateTimestampField(sm.order)
	return nil
}

func (sm *orderStateMachine) canTransitionPaymentStatusTo(newStatus OrderPaymentStatus) bool {
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

func (sm *orderStateMachine) transitionPaymentStatusTo(newStatus OrderPaymentStatus) error {
	if !sm.canTransitionPaymentStatusTo(newStatus) {
		return ErrOrderPaymentStatusUpdateNotAllowed(sm.order.ID, sm.order.PaymentStatus, newStatus)
	}
	allowedOrderStatuses := []OrderStatus{OrderStatusPlaced, OrderStatusReadyForShipment, OrderStatusShipped, OrderStatusFulfilled}
	if !slices.Contains(allowedOrderStatuses, sm.order.Status) {
		return ErrOrderPaymentStatusUpdateNotAllowedForOrderStatus(sm.order.ID, sm.order.Status, newStatus)
	}
	sm.order.PaymentStatus = newStatus
	newStatus.UpdateTimestampField(sm.order)
	return nil
}
