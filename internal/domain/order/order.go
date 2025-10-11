package order

import (
	"database/sql"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/govalues/decimal"
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
	Subtotal           decimal.Decimal    `gorm:"column:subtotal;type:numeric;not null;default:0" json:"subtotal"`
	VAT                decimal.Decimal    `gorm:"column:vat;type:numeric;not null;default:0" json:"vat"`
	VATRate            decimal.Decimal    `gorm:"column:vat_rate;type:numeric;not null;default:0" json:"vatRate"`
	ShippingFee        decimal.Decimal    `gorm:"column:shipping_fee;type:numeric;not null;default:0" json:"shippingFee"`
	Discount           decimal.Decimal    `gorm:"column:discount;type:numeric;not null;default:0" json:"discount"`
	Total              decimal.Decimal    `gorm:"column:total;type:numeric;not null;default:0" json:"total"`
	Currency           string             `gorm:"column:currency;type:text;not null" json:"currency"`
	Status             OrderStatus        `gorm:"column:status;type:text;not null;default:'pending'" json:"status"`
	PaymentStatus      OrderPaymentStatus `gorm:"column:payment_status;type:text;not null;default:'pending'" json:"paymentStatus"`
	PaymentMethod      OrderPaymentMethod `gorm:"column:payment_method;type:text;not null;default:'credit_card'" json:"paymentMethod"`
	PaymentReference   string             `gorm:"column:payment_reference;type:text" json:"paymentReference,omitempty"`
	PlacedAt           sql.NullTime       `gorm:"column:placed_at" json:"placedAt"`
	ReadyForShipmentAt sql.NullTime       `gorm:"column:ready_for_shipment_at" json:"readyForShipmentAt"`
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
	ShippingAddressID string                    `json:"shippingAddressId" binding:"required"`
	ShippingFee       decimal.Decimal           `json:"shippingFee" binding:"omitempty,gte=0"`
	Discount          decimal.Decimal           `json:"discount" binding:"omitempty,gte=0"`
	PaymentMethod     OrderPaymentMethod        `json:"paymentMethod" binding:"omitempty,oneof=credit_card paypal bank_transfer cash_on_delivery"`
	PaymentReference  string                    `json:"paymentReference" binding:"omitempty"`
	Items             []*CreateOrderItemRequest `json:"items" binding:"required,dive,required"`
}

type UpdateOrderRequest struct {
	Status             OrderStatus               `json:"status" binding:"omitempty,oneof=pending placed ready_for_shipment shipped fulfilled cancelled returned"`
	PaymentStatus      OrderPaymentStatus        `json:"paymentStatus" binding:"omitempty,oneof=pending paid failed refunded"`
	PaymentMethod      OrderPaymentMethod        `json:"paymentMethod" binding:"omitempty,oneof=credit_card paypal bank_transfer cash_on_delivery"`
	PaymentReference   string                    `json:"paymentReference" binding:"omitempty,required_with=PaymentMethod"`
	ShippingFee        decimal.Decimal           `json:"shippingFee" binding:"omitempty,gte=0"`
	Discount           decimal.Decimal           `json:"discount" binding:"omitempty,gte=0"`
	PlacedAt           sql.NullTime              `json:"placedAt" binding:"omitempty,required_with=Status eq placed"`
	ReadyForShipmentAt sql.NullTime              `json:"readyForShipmentAt" binding:"omitempty,required_with=Status eq ready_for_shipment"`
	ShippedAt          sql.NullTime              `json:"shippedAt" binding:"omitempty,required_with=Status eq shipped"`
	FulfilledAt        sql.NullTime              `json:"fulfilledAt" binding:"omitempty,required_with=Status eq fulfilled"`
	CancelledAt        sql.NullTime              `json:"cancelledAt" binding:"omitempty,required_with=Status eq cancelled"`
	ReturnedAt         sql.NullTime              `json:"returnedAt" binding:"omitempty,required_with=Status eq returned"`
	PaidAt             sql.NullTime              `json:"paidAt" binding:"omitempty,required_with=PaymentStatus eq paid"`
	FailedAt           sql.NullTime              `json:"failedAt" binding:"omitempty,required_with=PaymentStatus eq failed"`
	RefundedAt         sql.NullTime              `json:"refundedAt" binding:"omitempty,required_with=PaymentStatus eq refunded"`
	Items              []*UpdateOrderItemRequest `json:"items" binding:"omitempty,dive,required"`
}

type OrderFilter struct {
	CustomerIDs     []string             `form:"customerIds" json:"customerIds" binding:"omitempty,dive,required"`
	Statuses        []OrderStatus        `form:"status" json:"status" binding:"omitempty,oneof=pending placed ready_for_shipment shipped fulfilled cancelled returned"`
	PaymentStatuses []OrderPaymentStatus `form:"paymentStatus" json:"paymentStatus" binding:"omitempty,oneof=pending paid failed refunded"`
	PaymentMethods  []OrderPaymentMethod `form:"paymentMethod" json:"paymentMethod" binding:"omitempty,oneof=credit_card paypal bank_transfer cash_on_delivery"`
	CountryCodes    []string             `form:"countryCode" json:"countryCode" binding:"omitempty,dive,len=2"`
}
