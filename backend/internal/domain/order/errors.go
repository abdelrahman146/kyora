package order

import (
	"fmt"

	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
)

func ErrProductOutOfStock(variant *inventory.Variant) error {
	return problem.Conflict(fmt.Sprintf("product %q is out of stock", variant.Name)).
		With("variantId", variant.ID).
		With("stockQuantity", variant.StockQuantity).
		WithCode("order.product_out_of_stock")
}

func ErrInsufficientStock(variant *inventory.Variant, requestedQty int) error {
	return problem.Conflict(fmt.Sprintf("insufficient stock for product %q", variant.Name)).
		With("productId", variant.ProductID).
		With("variantId", variant.ID).
		With("requestedQuantity", requestedQty).
		With("availableQuantity", variant.StockQuantity).
		WithCode("order.insufficient_stock")
}

func ErrShippingZoneCountryMismatch(zoneID, countryCode string) error {
	return problem.BadRequest("shipping zone does not include destination country").
		With("countryCode", countryCode).
		With("zoneId", zoneID).
		WithCode("order.shipping_zone_country_mismatch")
}

// ErrEmptyOrderItems indicates that no items were provided in the order request
func ErrEmptyOrderItems() error {
	return problem.BadRequest("order must include at least one item").WithCode("order.empty_items")
}

// ErrVariantNotFound indicates that a requested variant could not be found
func ErrVariantNotFound(variantID string, err error) error {
	return problem.NotFound("variant not found").WithError(err).With("variantId", variantID).WithCode("order.variant_not_found")
}

// ErrInvalidOrderItemQuantity indicates that a provided item quantity is invalid
func ErrInvalidOrderItemQuantity(productID string, quantity int) error {
	return problem.BadRequest("invalid order item quantity").With("productId", productID).With("quantity", quantity).WithCode("order.invalid_item_quantity")
}

// ErrOrderNumberGenerationFailed indicates that generating a unique order number failed after retries
func ErrOrderNumberGenerationFailed(err error) error {
	return problem.InternalError().WithError(err).With("reason", "failed to generate a unique order number").WithCode("order.number_generation_failed")
}

// ErrOrderNotFound indicates that an order with the given id doesn't exist (in this business)
func ErrOrderNotFound(orderID string, err error) error {
	return problem.NotFound("order not found").WithError(err).With("orderId", orderID).WithCode("order.not_found")
}

// ErrOrderItemsCountMismatch indicates request items length doesn't match existing items length
func ErrOrderItemsCountMismatch(expected, got int) error {
	return problem.BadRequest("items count mismatch").With("expected", expected).With("got", got).WithCode("order.items_count_mismatch")
}

// ErrOrderItemsUpdateNotAllowed indicates items cannot be updated for orders in immutable statuses
func ErrOrderItemsUpdateNotAllowed(orderID string, status OrderStatus) error {
	return problem.Conflict("cannot update items for this order status").With("orderId", orderID).With("status", string(status)).WithCode("order.items_update_not_allowed")
}

func ErrOrderStatusUpdateNotAllowed(orderID string, from, to OrderStatus) error {
	return problem.Conflict(fmt.Sprintf("cannot update order status from %s to %s", from, to)).
		With("orderId", orderID).
		With("fromStatus", string(from)).
		With("toStatus", string(to)).
		WithCode("order.status_update_not_allowed")
}

func ErrOrderPaymentStatusUpdateNotAllowed(orderID string, from, to OrderPaymentStatus) error {
	return problem.Conflict(fmt.Sprintf("cannot update order payment status from %s to %s", from, to)).
		With("orderId", orderID).
		With("fromStatus", string(from)).
		With("toStatus", string(to)).
		WithCode("order.payment_status_update_not_allowed")
}

func ErrOrderPaymentStatusUpdateNotAllowedForOrderStatus(orderID string, orderStatus OrderStatus, paymentStatus OrderPaymentStatus) error {
	return problem.Conflict(fmt.Sprintf("cannot set payment status to %s when order status is %s", paymentStatus, orderStatus)).
		With("orderId", orderID).
		With("orderStatus", string(orderStatus)).
		With("paymentStatus", string(paymentStatus)).
		WithCode("order.payment_status_invalid_for_order_status")
}

func ErrOrderCannotBeDeleted(orderID string, status OrderStatus) error {
	return problem.Conflict("cannot delete order in its current status").
		With("orderId", orderID).
		With("status", string(status)).
		WithCode("order.cannot_delete")
}

func ErrOrderCannotBeCancelled(orderID string, status OrderStatus) error {
	return problem.Conflict("cannot cancel order in its current status").
		With("orderId", orderID).
		With("status", string(status)).
		WithCode("order.cannot_cancel")
}

func ErrUpdateOrderItemNotAllowed(orderID, itemID string, status OrderStatus) error {
	return problem.Conflict("cannot update order item in its current order status").
		With("orderId", orderID).
		With("itemId", itemID).
		With("status", string(status)).
		WithCode("order.item_update_not_allowed")
}

func ErrOrderNoteNotFound(orderNoteID string, err error) error {
	return problem.NotFound("order note not found").WithError(err).With("orderNoteId", orderNoteID).WithCode("order.note_not_found")
}

// ErrOrderRateLimited indicates the client exceeded rate limits.
func ErrOrderRateLimited() error {
	return problem.TooManyRequests("too many requests").WithCode("order.rate_limited")
}
