package order

import (
	"fmt"
	"slices"

	"github.com/abdelrahman146/kyora/internal/utils"
)

type orderStateMachine struct {
	order *Order
}

func newOrderStateMachine(order *Order) *orderStateMachine {
	return &orderStateMachine{order: order}
}

func (sm *orderStateMachine) canTransitionStateTo(newState OrderStatus) bool {
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

func (sm *orderStateMachine) transitionStateTo(newState OrderStatus) error {
	if !sm.canTransitionStateTo(newState) {
		return utils.Problem.UnprocessableEntity(fmt.Sprintf("cannot transition order status from %s to %s", sm.order.Status, newState)).With("currentState", sm.order.Status)
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
		return utils.Problem.UnprocessableEntity(fmt.Sprintf("cannot transition payment status from %s to %s", sm.order.PaymentStatus, newStatus)).With("currentState", sm.order.PaymentStatus)
	}
	allowedOrderStatuses := []OrderStatus{OrderStatusPlaced, OrderStatusShipped, OrderStatusFulfilled}
	if !slices.Contains(allowedOrderStatuses, sm.order.Status) {
		return utils.Problem.UnprocessableEntity(fmt.Sprintf("cannot set payment status to %s when order status is %s", newStatus, sm.order.Status)).With("orderStatus", sm.order.Status)
	}
	sm.order.PaymentStatus = newStatus
	newStatus.UpdateTimestampField(sm.order)
	return nil
}
