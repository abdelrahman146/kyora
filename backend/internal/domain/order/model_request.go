package order

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type CreateOrderRequest struct {
	CustomerID        string          `json:"customerId" binding:"required"`
	Channel           string          `json:"channel" binding:"required"`
	ShippingAddressID string          `json:"shippingAddressId" binding:"required"`
	ShippingZoneID    *string         `json:"shippingZoneId" binding:"omitempty"`
	ShippingFee       decimal.Decimal `json:"shippingFee" binding:"omitempty"`
	// Legacy discount field (amount-based). Still supported for backward compatibility.
	Discount decimal.Decimal `json:"discount" binding:"omitempty"`
	// New discount fields (preferred). When provided, these take precedence over Discount.
	DiscountType  DiscountType    `json:"discountType" binding:"omitempty,oneof=amount percent"`
	DiscountValue decimal.Decimal `json:"discountValue" binding:"omitempty"`
	// Optional target order status (advanced). If provided, backend will attempt to apply it atomically.
	Status *OrderStatus `json:"status" binding:"omitempty,oneof=pending placed ready_for_shipment shipped fulfilled cancelled returned"`
	// Optional target payment status (advanced). If provided, backend will attempt to apply it atomically.
	PaymentStatus    *OrderPaymentStatus `json:"paymentStatus" binding:"omitempty,oneof=pending paid failed refunded"`
	PaymentMethod    OrderPaymentMethod  `json:"paymentMethod" binding:"omitempty,oneof=credit_card paypal bank_transfer cash_on_delivery tamara tabby"`
	PaymentReference sql.NullString      `json:"paymentReference" binding:"omitempty"`
	OrderedAt        time.Time           `json:"orderedAt" binding:"omitempty"`
	// Optional single note content. If provided, a note will be created as part of order creation.
	Note  string                    `json:"note" binding:"omitempty"`
	Items []*CreateOrderItemRequest `json:"items" binding:"required,dive,required"`
}

type UpdateOrderRequest struct {
	ShippingAddressID *string             `json:"shippingAddressId" binding:"omitempty"`
	ShippingZoneID    *string             `json:"shippingZoneId" binding:"omitempty"`
	ShippingFee       decimal.NullDecimal `json:"shippingFee" binding:"omitempty"`
	Channel           string              `json:"channel" binding:"omitempty"`
	// Legacy discount field (amount-based). Still supported for backward compatibility.
	Discount decimal.NullDecimal `json:"discount" binding:"omitempty"`
	// New discount fields (preferred). When provided, these take precedence over Discount.
	DiscountType  DiscountType              `json:"discountType" binding:"omitempty,oneof=amount percent"`
	DiscountValue decimal.NullDecimal       `json:"discountValue" binding:"omitempty"`
	OrderedAt     time.Time                 `json:"orderedAt" binding:"omitempty"`
	Items         []*CreateOrderItemRequest `json:"items,omitempty" binding:"omitempty,dive,required"`
}

type AddOrderPaymentDetailsRequest struct {
	PaymentMethod    OrderPaymentMethod `json:"paymentMethod" binding:"required,oneof=credit_card paypal bank_transfer cash_on_delivery tamara tabby"`
	PaymentReference sql.NullString     `json:"paymentReference" binding:"omitempty"`
}

type CreateOrderItemRequest struct {
	VariantID string          `json:"variantId" binding:"required"`
	Quantity  int             `json:"quantity" binding:"required,min=1"`
	UnitPrice decimal.Decimal `json:"unitPrice" binding:"required"`
	UnitCost  decimal.Decimal `json:"unitCost" binding:"omitempty"`
}

type CreateOrderNoteRequest struct {
	OrderID string `json:"orderId" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdateOrderNoteRequest struct {
	Content string `json:"content" binding:"omitempty"`
}

// Query and handler request types

// listOrdersQuery represents the query parameters for listing orders.
type listOrdersQuery struct {
	Page            int       `form:"page" binding:"omitempty,min=1"`
	PageSize        int       `form:"pageSize" binding:"omitempty,min=1,max=100"`
	OrderBy         []string  `form:"orderBy" binding:"omitempty"`
	SearchTerm      string    `form:"search" binding:"omitempty"`
	Status          []string  `form:"status" binding:"omitempty"`
	PaymentStatus   []string  `form:"paymentStatus" binding:"omitempty"`
	SocialPlatforms []string  `form:"socialPlatforms" binding:"omitempty"`
	CustomerID      string    `form:"customerId" binding:"omitempty"`
	OrderNumber     string    `form:"orderNumber" binding:"omitempty"`
	From            time.Time `form:"from" time_format:"2006-01-02T15:04:05Z07:00" binding:"omitempty"`
	To              time.Time `form:"to" time_format:"2006-01-02T15:04:05Z07:00" binding:"omitempty"`
}

// updateOrderStatusRequest represents the request to update order status.
type updateOrderStatusRequest struct {
	Status OrderStatus `json:"status" binding:"required,oneof=pending placed ready_for_shipment shipped fulfilled cancelled returned"`
}

// updateOrderPaymentStatusRequest represents the request to update payment status.
type updateOrderPaymentStatusRequest struct {
	PaymentStatus OrderPaymentStatus `json:"paymentStatus" binding:"required,oneof=pending paid failed refunded"`
}

// addOrderPaymentDetailsRequest represents the request to add payment details (used in handler).
type addOrderPaymentDetailsRequest struct {
	PaymentMethod    OrderPaymentMethod `json:"paymentMethod" binding:"required,oneof=credit_card paypal bank_transfer cash_on_delivery tamara tabby"`
	PaymentReference string             `json:"paymentReference" binding:"omitempty"`
}

// createOrderNoteRequest represents the request to create an order note (used in handler).
type createOrderNoteRequest struct {
	Content string `json:"content" binding:"required"`
}

// updateOrderNoteRequest represents the request to update an order note (used in handler).
type updateOrderNoteRequest struct {
	Content string `json:"content" binding:"required"`
}
