package order

import (
	"database/sql"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/utils"
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
)

const (
	OrderTable  = "orders"
	OrderAlias  = "ord"
	OrderStruct = "Order"
)

type Order struct {
	gorm.Model
	ID                 string             `gorm:"column:id;primaryKey;type:text" json:"id"`
	OrderNumber        string             `gorm:"column:order_number;type:text;not null;uniqueIndex:order_number_store_id_idx" json:"orderNumber"`
	StoreID            string             `gorm:"column:store_id;type:text;not null;index;uniqueIndex:order_number_store_id_idx" json:"storeId"`
	Store              *store.Store       `gorm:"foreignKey:StoreID;references:ID" json:"store,omitempty"`
	CustomerID         string             `gorm:"column:customer_id;type:text;index" json:"customerId,omitempty"`
	Customer           *customer.Customer `gorm:"foreignKey:CustomerID;references:ID" json:"customer,omitempty"`
	ShippingAddressID  string             `gorm:"column:shipping_address_id;type:text" json:"shippingAddressId,omitempty"`
	ShippingAddress    *customer.Address  `gorm:"foreignKey:ShippingAddressID;references:ID" json:"shippingAddress,omitempty"`
	Channel            string             `gorm:"column:channel;type:text;not null" json:"channel"`
	Subtotal           decimal.Decimal    `gorm:"column:subtotal;type:numeric;not null;default:0" json:"subtotal"`
	VAT                decimal.Decimal    `gorm:"column:vat;type:numeric;not null;default:0" json:"vat"`
	VATRate            decimal.Decimal    `gorm:"column:vat_rate;type:numeric;not null;default:0" json:"vatRate"`
	ShippingFee        decimal.Decimal    `gorm:"column:shipping_fee;type:numeric;not null;default:0" json:"shippingFee"`
	Discount           decimal.Decimal    `gorm:"column:discount;type:numeric;not null;default:0" json:"discount"`
	COGS               decimal.Decimal    `gorm:"column:cogs;type:numeric;not null;default:0" json:"cogs"`
	Total              decimal.Decimal    `gorm:"column:total;type:numeric;not null;default:0" json:"total"`
	Currency           string             `gorm:"column:currency;type:text;not null" json:"currency"`
	Status             OrderStatus        `gorm:"column:status;type:text;not null;default:'pending'" json:"status"`
	PaymentStatus      OrderPaymentStatus `gorm:"column:payment_status;type:text;not null;default:'pending'" json:"paymentStatus"`
	PaymentMethod      OrderPaymentMethod `gorm:"column:payment_method;type:text;not null;default:'bank_transfer'" json:"paymentMethod"`
	PaymentReference   sql.NullString     `gorm:"column:payment_reference;type:text" json:"paymentReference,omitempty"`
	PlacedAt           sql.NullTime       `gorm:"column:placed_at" json:"placedAt"`
	ReadyForShipmentAt sql.NullTime       `gorm:"column:ready_for_shipment_at" json:"readyForShipmentAt"`
	OrderedAt          time.Time          `gorm:"column:ordered_at;type:timestamptz;not null;default:now()" json:"orderedAt"`
	ShippedAt          sql.NullTime       `gorm:"column:shipped_at" json:"shippedAt"`
	FulfilledAt        sql.NullTime       `gorm:"column:fulfilled_at" json:"fulfilledAt"`
	CancelledAt        sql.NullTime       `gorm:"column:cancelled_at" json:"cancelledAt"`
	ReturnedAt         sql.NullTime       `gorm:"column:returned_at" json:"returnedAt"`
	PaidAt             sql.NullTime       `gorm:"column:paid_at" json:"paidAt"`
	FailedAt           sql.NullTime       `gorm:"column:failed_at" json:"failedAt"`
	RefundedAt         sql.NullTime       `gorm:"column:refunded_at" json:"refundedAt"`
	Items              []*OrderItem       `gorm:"foreignKey:OrderID;references:ID" json:"items"`
	Notes              []*OrderNote       `gorm:"foreignKey:OrderID;references:ID" json:"notes,omitempty"`
}

func (m *Order) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = utils.ID.NewUlidWithPrefix(OrderAlias)
	}
	return
}

type CreateOrderRequest struct {
	CustomerID        string                    `json:"customerId" binding:"required"`
	Channel           string                    `json:"channel" binding:"required"`
	ShippingAddressID string                    `json:"shippingAddressId" binding:"required"`
	ShippingFee       decimal.Decimal           `json:"shippingFee" binding:"omitempty,gte=0"`
	Discount          decimal.Decimal           `json:"discount" binding:"omitempty,gte=0"`
	PaymentMethod     OrderPaymentMethod        `json:"paymentMethod" binding:"omitempty,oneof=credit_card paypal bank_transfer cash_on_delivery"`
	PaymentReference  sql.NullString            `json:"paymentReference" binding:"omitempty"`
	OrderedAt         time.Time                 `json:"orderedAt" binding:"omitempty"`
	Items             []*CreateOrderItemRequest `json:"items" binding:"required,dive,required"`
}

type UpdateOrderRequest struct {
	ShippingFee decimal.Decimal           `json:"shippingFee" binding:"omitempty,gte=0"`
	Channel     string                    `json:"channel" binding:"omitempty"`
	Discount    decimal.Decimal           `json:"discount" binding:"omitempty,gte=0"`
	OrderedAt   time.Time                 `json:"orderedAt" binding:"omitempty"`
	Items       []*UpdateOrderItemRequest `json:"items" binding:"omitempty,dive,required"`
}

type AddOrderPaymentDetailsRequest struct {
	PaymentMethod    OrderPaymentMethod `json:"paymentMethod" binding:"required,oneof=credit_card paypal bank_transfer cash_on_delivery"`
	PaymentReference sql.NullString     `json:"paymentReference" binding:"omitempty"`
}

type OrderFilter struct {
	CustomerIDs     []string             `form:"customerIds" json:"customerIds" binding:"omitempty,dive,required"`
	Statuses        []OrderStatus        `form:"status" json:"status" binding:"omitempty,oneof=pending placed ready_for_shipment shipped fulfilled cancelled returned"`
	PaymentStatuses []OrderPaymentStatus `form:"paymentStatus" json:"paymentStatus" binding:"omitempty,oneof=pending paid failed refunded"`
	PaymentMethods  []OrderPaymentMethod `form:"paymentMethod" json:"paymentMethod" binding:"omitempty,oneof=credit_card paypal bank_transfer cash_on_delivery"`
	CountryCodes    []string             `form:"countryCode" json:"countryCode" binding:"omitempty,dive,len=2"`
}
