package order

import (
	"database/sql"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/platform/types/schema"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPending          OrderStatus = "pending"
	OrderStatusPlaced           OrderStatus = "placed"
	OrderStatusReadyForShipment OrderStatus = "ready_for_shipment"
	OrderStatusShipped          OrderStatus = "shipped"
	OrderStatusFulfilled        OrderStatus = "fulfilled"
	OrderStatusCancelled        OrderStatus = "cancelled"
	OrderStatusReturned         OrderStatus = "returned"
)

func (s OrderStatus) UpdateTimestampField(order *Order) {
	switch s {
	case OrderStatusPlaced:
		order.PlacedAt = sql.NullTime{Time: time.Now(), Valid: true}
	case OrderStatusReadyForShipment:
		order.ReadyForShipmentAt = sql.NullTime{Time: time.Now(), Valid: true}
	case OrderStatusShipped:
		order.ShippedAt = sql.NullTime{Time: time.Now(), Valid: true}
	case OrderStatusFulfilled:
		order.FulfilledAt = sql.NullTime{Time: time.Now(), Valid: true}
	case OrderStatusCancelled:
		order.CancelledAt = sql.NullTime{Time: time.Now(), Valid: true}
	}
}

type OrderPaymentStatus string

const (
	OrderPaymentStatusPending  OrderPaymentStatus = "pending"
	OrderPaymentStatusPaid     OrderPaymentStatus = "paid"
	OrderPaymentStatusFailed   OrderPaymentStatus = "failed"
	OrderPaymentStatusRefunded OrderPaymentStatus = "refunded"
)

func (s OrderPaymentStatus) UpdateTimestampField(order *Order) {
	switch s {
	case OrderPaymentStatusPaid:
		order.PaidAt = sql.NullTime{Time: time.Now(), Valid: true}
	case OrderPaymentStatusFailed:
		order.FailedAt = sql.NullTime{Time: time.Now(), Valid: true}
	case OrderPaymentStatusRefunded:
		order.RefundedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}
}

type OrderPaymentMethod string

const (
	OrderPaymentMethodCreditCard     OrderPaymentMethod = "credit_card"
	OrderPaymentMethodPayPal         OrderPaymentMethod = "paypal"
	OrderPaymentMethodBankTransfer   OrderPaymentMethod = "bank_transfer"
	OrderPaymentMethodCashOnDelivery OrderPaymentMethod = "cash_on_delivery"
	OrderPaymentMethodTamara         OrderPaymentMethod = "tamara"
	OrderPaymentMethodTabby          OrderPaymentMethod = "tabby"
)

const (
	OrderTable  = "orders"
	OrderStruct = "Order"
	OrderPrefix = "ord"
)

type Order struct {
	gorm.Model
	ID                 string                    `gorm:"column:id;primaryKey;type:text" json:"id"`
	OrderNumber        string                    `gorm:"column:order_number;type:text;not null;uniqueIndex:order_number_business_id_idx" json:"orderNumber"`
	BusinessID         string                    `gorm:"column:business_id;type:text;not null;index;uniqueIndex:order_number_business_id_idx" json:"businessId"`
	Business           *business.Business        `gorm:"foreignKey:BusinessID;references:ID" json:"business,omitempty"`
	CustomerID         string                    `gorm:"column:customer_id;type:text;index" json:"customerId,omitempty"`
	Customer           *customer.Customer        `gorm:"foreignKey:CustomerID;references:ID" json:"customer,omitempty"`
	ShippingAddressID  string                    `gorm:"column:shipping_address_id;type:text" json:"shippingAddressId,omitempty"`
	ShippingAddress    *customer.CustomerAddress `gorm:"foreignKey:ShippingAddressID;references:ID" json:"shippingAddress,omitempty"`
	ShippingZoneID     *string                   `gorm:"column:shipping_zone_id;type:text;index" json:"shippingZoneId,omitempty"`
	ShippingZone       *business.ShippingZone    `gorm:"foreignKey:ShippingZoneID;references:ID" json:"shippingZone,omitempty"`
	Channel            string                    `gorm:"column:channel;type:text;not null" json:"channel"`
	Subtotal           decimal.Decimal           `gorm:"column:subtotal;type:numeric;not null;default:0" json:"subtotal"`
	VAT                decimal.Decimal           `gorm:"column:vat;type:numeric;not null;default:0" json:"vat"`
	VATRate            decimal.Decimal           `gorm:"column:vat_rate;type:numeric;not null;default:0" json:"vatRate"`
	ShippingFee        decimal.Decimal           `gorm:"column:shipping_fee;type:numeric;not null;default:0" json:"shippingFee"`
	Discount           decimal.Decimal           `gorm:"column:discount;type:numeric;not null;default:0" json:"discount"`
	DiscountType       DiscountType              `gorm:"column:discount_type;type:text" json:"discountType,omitempty"`
	DiscountValue      decimal.Decimal           `gorm:"column:discount_value;type:numeric;default:0" json:"discountValue,omitempty"`
	COGS               decimal.Decimal           `gorm:"column:cogs;type:numeric;not null;default:0" json:"cogs"`
	Total              decimal.Decimal           `gorm:"column:total;type:numeric;not null;default:0" json:"total"`
	Currency           string                    `gorm:"column:currency;type:text;not null" json:"currency"`
	Status             OrderStatus               `gorm:"column:status;type:text;not null;default:'pending'" json:"status"`
	PaymentStatus      OrderPaymentStatus        `gorm:"column:payment_status;type:text;not null;default:'pending'" json:"paymentStatus"`
	PaymentMethod      OrderPaymentMethod        `gorm:"column:payment_method;type:text;not null;default:'bank_transfer'" json:"paymentMethod"`
	PaymentReference   sql.NullString            `gorm:"column:payment_reference;type:text" json:"paymentReference,omitempty"`
	PlacedAt           sql.NullTime              `gorm:"column:placed_at" json:"placedAt"`
	ReadyForShipmentAt sql.NullTime              `gorm:"column:ready_for_shipment_at" json:"readyForShipmentAt"`
	OrderedAt          time.Time                 `gorm:"column:ordered_at;type:timestamptz;not null;default:now()" json:"orderedAt"`
	ShippedAt          sql.NullTime              `gorm:"column:shipped_at" json:"shippedAt"`
	FulfilledAt        sql.NullTime              `gorm:"column:fulfilled_at" json:"fulfilledAt"`
	CancelledAt        sql.NullTime              `gorm:"column:cancelled_at" json:"cancelledAt"`
	ReturnedAt         sql.NullTime              `gorm:"column:returned_at" json:"returnedAt"`
	PaidAt             sql.NullTime              `gorm:"column:paid_at" json:"paidAt"`
	FailedAt           sql.NullTime              `gorm:"column:failed_at" json:"failedAt"`
	RefundedAt         sql.NullTime              `gorm:"column:refunded_at" json:"refundedAt"`
	Items              []*OrderItem              `gorm:"foreignKey:OrderID;references:ID" json:"items"`
	Notes              []*OrderNote              `gorm:"foreignKey:OrderID;references:ID" json:"notes,omitempty"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == "" {
		o.ID = id.KsuidWithPrefix(OrderPrefix)
	}
	return
}

type DiscountType string

const (
	DiscountTypeAmount  DiscountType = "amount"
	DiscountTypePercent DiscountType = "percent"
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

var OrderSchema = struct {
	ID                 schema.Field
	OrderNumber        schema.Field
	BusinessID         schema.Field
	CustomerID         schema.Field
	ShippingAddressID  schema.Field
	ShippingZoneID     schema.Field
	Channel            schema.Field
	Subtotal           schema.Field
	VAT                schema.Field
	VATRate            schema.Field
	ShippingFee        schema.Field
	Discount           schema.Field
	COGS               schema.Field
	Total              schema.Field
	Currency           schema.Field
	Status             schema.Field
	PaymentStatus      schema.Field
	PaymentMethod      schema.Field
	PaymentReference   schema.Field
	PlacedAt           schema.Field
	ReadyForShipmentAt schema.Field
	OrderedAt          schema.Field
	ShippedAt          schema.Field
	FulfilledAt        schema.Field
	CancelledAt        schema.Field
	ReturnedAt         schema.Field
	PaidAt             schema.Field
	FailedAt           schema.Field
	RefundedAt         schema.Field
	CreatedAt          schema.Field
	UpdatedAt          schema.Field
	DeletedAt          schema.Field
}{
	ID:                 schema.NewField("id", "id"),
	OrderNumber:        schema.NewField("order_number", "orderNumber"),
	BusinessID:         schema.NewField("business_id", "businessId"),
	CustomerID:         schema.NewField("customer_id", "customerId"),
	ShippingAddressID:  schema.NewField("shipping_address_id", "shippingAddressId"),
	ShippingZoneID:     schema.NewField("shipping_zone_id", "shippingZoneId"),
	Channel:            schema.NewField("channel", "channel"),
	Subtotal:           schema.NewField("subtotal", "subtotal"),
	VAT:                schema.NewField("vat", "vat"),
	VATRate:            schema.NewField("vat_rate", "vatRate"),
	ShippingFee:        schema.NewField("shipping_fee", "shippingFee"),
	Discount:           schema.NewField("discount", "discount"),
	COGS:               schema.NewField("cogs", "cogs"),
	Total:              schema.NewField("total", "total"),
	Currency:           schema.NewField("currency", "currency"),
	Status:             schema.NewField("status", "status"),
	PaymentStatus:      schema.NewField("payment_status", "paymentStatus"),
	PaymentMethod:      schema.NewField("payment_method", "paymentMethod"),
	PaymentReference:   schema.NewField("payment_reference", "paymentReference"),
	PlacedAt:           schema.NewField("placed_at", "placedAt"),
	ReadyForShipmentAt: schema.NewField("ready_for_shipment_at", "readyForShipmentAt"),
	OrderedAt:          schema.NewField("ordered_at", "orderedAt"),
	ShippedAt:          schema.NewField("shipped_at", "shippedAt"),
	FulfilledAt:        schema.NewField("fulfilled_at", "fulfilledAt"),
	CancelledAt:        schema.NewField("cancelled_at", "cancelledAt"),
	ReturnedAt:         schema.NewField("returned_at", "returnedAt"),
	PaidAt:             schema.NewField("paid_at", "paidAt"),
	FailedAt:           schema.NewField("failed_at", "failedAt"),
	RefundedAt:         schema.NewField("refunded_at", "refundedAt"),
	CreatedAt:          schema.NewField("created_at", "createdAt"),
	UpdatedAt:          schema.NewField("updated_at", "updatedAt"),
	DeletedAt:          schema.NewField("deleted_at", "deletedAt"),
}

const (
	OrderItemTable  = "order_items"
	OrderItemStruct = "Items"
	OrderItemPrefix = "oitm"
)

type OrderItem struct {
	gorm.Model
	ID        string             `gorm:"column:id;primaryKey;type:text" json:"id"`
	OrderID   string             `gorm:"column:order_id;type:text;not null;index" json:"orderId"`
	Order     *Order             `gorm:"foreignKey:OrderID;references:ID;OnDelete:CASCADE" json:"order,omitempty"`
	ProductID string             `gorm:"column:product_id;type:text;not null;index" json:"productId"`
	Product   *inventory.Product `gorm:"foreignKey:ProductID;references:ID" json:"product,omitempty"`
	VariantID string             `gorm:"column:variant_id;type:text;not null;index" json:"variantId"`
	Variant   *inventory.Variant `gorm:"foreignKey:VariantID;references:ID" json:"variant,omitempty"`
	Quantity  int                `gorm:"column:quantity;type:int;not null;default:1" json:"quantity"`
	Currency  string             `gorm:"column:currency;type:text;not null" json:"currency"`
	UnitPrice decimal.Decimal    `gorm:"column:unit_price;type:numeric;not null;default:0" json:"unitPrice"`
	UnitCost  decimal.Decimal    `gorm:"column:unit_cost;type:numeric;not null;default:0" json:"unitCost"`
	TotalCost decimal.Decimal    `gorm:"column:total_cost;type:numeric;not null;default:0" json:"totalCost"`
	Total     decimal.Decimal    `gorm:"column:total;type:numeric;not null;default:0" json:"total"`
}

func (m *OrderItem) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(OrderItemPrefix)
	}
	return
}

type CreateOrderItemRequest struct {
	VariantID string          `json:"variantId" binding:"required"`
	Quantity  int             `json:"quantity" binding:"required,min=1"`
	UnitPrice decimal.Decimal `json:"unitPrice" binding:"required"`
	UnitCost  decimal.Decimal `json:"unitCost" binding:"omitempty"`
}

var OrderItemSchema = struct {
	ID        schema.Field
	OrderID   schema.Field
	ProductID schema.Field
	VariantID schema.Field
	Quantity  schema.Field
	Currency  schema.Field
	UnitPrice schema.Field
	UnitCost  schema.Field
	TotalCost schema.Field
	Total     schema.Field
	CreatedAt schema.Field
	UpdatedAt schema.Field
	DeletedAt schema.Field
}{
	ID:        schema.NewField("id", "id"),
	OrderID:   schema.NewField("order_id", "orderId"),
	ProductID: schema.NewField("product_id", "productId"),
	VariantID: schema.NewField("variant_id", "variantId"),
	Quantity:  schema.NewField("quantity", "quantity"),
	Currency:  schema.NewField("currency", "currency"),
	UnitPrice: schema.NewField("unit_price", "unitPrice"),
	UnitCost:  schema.NewField("unit_cost", "unitCost"),
	TotalCost: schema.NewField("total_cost", "totalCost"),
	Total:     schema.NewField("total", "total"),
	CreatedAt: schema.NewField("created_at", "createdAt"),
	UpdatedAt: schema.NewField("updated_at", "updatedAt"),
	DeletedAt: schema.NewField("deleted_at", "deletedAt"),
}

const (
	OrderNoteTable  = "order_notes"
	OrderNoteStruct = "Notes"
	OrderNotePrefix = "onot"
)

type OrderNote struct {
	gorm.Model
	ID      string `gorm:"column:id;primaryKey;type:text" json:"id"`
	OrderID string `gorm:"column:order_id;type:text;not null;index" json:"orderId"`
	Order   *Order `gorm:"foreignKey:OrderID;references:ID;OnDelete:CASCADE" json:"order,omitempty"`
	Content string `gorm:"column:content;type:text;not null" json:"content"`
}

func (m *OrderNote) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(OrderNotePrefix)
	}
	return
}

type CreateOrderNoteRequest struct {
	OrderID string `json:"orderId" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdateOrderNoteRequest struct {
	Content string `json:"content" binding:"omitempty"`
}

var OrderNoteSchema = struct {
	ID        schema.Field
	OrderID   schema.Field
	Content   schema.Field
	CreatedAt schema.Field
	UpdatedAt schema.Field
	DeletedAt schema.Field
}{
	ID:        schema.NewField("id", "id"),
	OrderID:   schema.NewField("order_id", "orderId"),
	Content:   schema.NewField("content", "content"),
	CreatedAt: schema.NewField("created_at", "createdAt"),
	UpdatedAt: schema.NewField("updated_at", "updatedAt"),
	DeletedAt: schema.NewField("deleted_at", "deletedAt"),
}
