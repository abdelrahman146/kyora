package order

import (
	"fmt"

	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
)

func ErrProductOutOfStock(variant *inventory.Variant) error {
	return problem.Conflict(fmt.Sprintf("product %q is out of stock", variant.Name)).
		With("variantId", variant.ID).
		With("stockQuantity", variant.StockQuantity)
}

func ErrInsufficientStock(variant *inventory.Variant, requestedQty int) error {
	return problem.Conflict(fmt.Sprintf("insufficient stock for product %q", variant.Name)).
		With("productId", variant.ProductID).
		With("variantId", variant.ID).
		With("requestedQuantity", requestedQty).
		With("availableQuantity", variant.StockQuantity)
}

// ErrEmptyOrderItems indicates that no items were provided in the order request
func ErrEmptyOrderItems() error {
	return problem.BadRequest("order must include at least one item")
}

// ErrVariantNotFound indicates that a requested variant could not be found
func ErrVariantNotFound(variantID string, err error) error {
	return problem.NotFound("variant not found").WithError(err).With("variantId", variantID)
}

// ErrInvalidOrderItemQuantity indicates that a provided item quantity is invalid
func ErrInvalidOrderItemQuantity(productID string, quantity int) error {
	return problem.BadRequest("invalid order item quantity").With("productId", productID).With("quantity", quantity)
}

// ErrOrderNumberGenerationFailed indicates that generating a unique order number failed after retries
func ErrOrderNumberGenerationFailed(err error) error {
	return problem.InternalError().WithError(err).With("reason", "failed to generate a unique order number")
}

// ErrOrderNotFound indicates that an order with the given id doesn't exist (in this business)
func ErrOrderNotFound(orderID string, err error) error {
	return problem.NotFound("order not found").WithError(err).With("orderId", orderID)
}

// ErrOrderItemsCountMismatch indicates request items length doesn't match existing items length
func ErrOrderItemsCountMismatch(expected, got int) error {
	return problem.BadRequest("items count mismatch").With("expected", expected).With("got", got)
}

// ErrOrderItemsUpdateNotAllowed indicates items cannot be updated for orders in immutable statuses
func ErrOrderItemsUpdateNotAllowed(orderID string, status OrderStatus) error {
	return problem.Conflict("cannot update items for this order status").With("orderId", orderID).With("status", string(status))
}

func ErrOrderStatusUpdateNotAllowed(orderID string, from, to OrderStatus) error {
	return problem.Conflict(fmt.Sprintf("cannot update order status from %s to %s", from, to)).
		With("orderId", orderID).
		With("fromStatus", string(from)).
		With("toStatus", string(to))
}

func ErrOrderPaymentStatusUpdateNotAllowed(orderID string, from, to OrderPaymentStatus) error {
	return problem.Conflict(fmt.Sprintf("cannot update order payment status from %s to %s", from, to)).
		With("orderId", orderID).
		With("fromStatus", string(from)).
		With("toStatus", string(to))
}

func ErrOrderPaymentStatusUpdateNotAllowedForOrderStatus(orderID string, orderStatus OrderStatus, paymentStatus OrderPaymentStatus) error {
	return problem.Conflict(fmt.Sprintf("cannot set payment status to %s when order status is %s", paymentStatus, orderStatus)).
		With("orderId", orderID).
		With("orderStatus", string(orderStatus)).
		With("paymentStatus", string(paymentStatus))
}

func ErrOrderCannotBeDeleted(orderID string, status OrderStatus) error {
	return problem.Conflict("cannot delete order in its current status").
		With("orderId", orderID).
		With("status", string(status))
}

func ErrOrderCannotBeCancelled(orderID string, status OrderStatus) error {
	return problem.Conflict("cannot cancel order in its current status").
		With("orderId", orderID).
		With("status", string(status))
}

func ErrUpdateOrderItemNotAllowed(orderID, itemID string, status OrderStatus) error {
	return problem.Conflict("cannot update order item in its current order status").
		With("orderId", orderID).
		With("itemId", itemID).
		With("status", string(status))
}

func ErrOrderNoteNotFound(orderNoteID string, err error) error {
	return problem.NotFound("order note not found").WithError(err).With("orderNoteId", orderNoteID)
}
