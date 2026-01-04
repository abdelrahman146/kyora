package order

import (
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/platform/utils/transformer"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// OrderResponse normalizes nullable fields for API responses and preloads.
type OrderResponse struct {
	ID                 string                    `json:"id"`
	OrderNumber        string                    `json:"orderNumber"`
	BusinessID         string                    `json:"businessId"`
	CustomerID         string                    `json:"customerId"`
	Customer           *customer.Customer        `json:"customer,omitempty"`
	ShippingAddressID  string                    `json:"shippingAddressId"`
	ShippingAddress    *customer.CustomerAddress `json:"shippingAddress,omitempty"`
	ShippingZoneID     *string                   `json:"shippingZoneId,omitempty"`
	ShippingZone       *business.ShippingZone    `json:"shippingZone,omitempty"`
	Channel            string                    `json:"channel"`
	Subtotal           decimal.Decimal           `json:"subtotal"`
	VAT                decimal.Decimal           `json:"vat"`
	VATRate            decimal.Decimal           `json:"vatRate"`
	ShippingFee        decimal.Decimal           `json:"shippingFee"`
	Discount           decimal.Decimal           `json:"discount"`
	COGS               decimal.Decimal           `json:"cogs"`
	Total              decimal.Decimal           `json:"total"`
	Currency           string                    `json:"currency"`
	Status             OrderStatus               `json:"status"`
	PaymentStatus      OrderPaymentStatus        `json:"paymentStatus"`
	PaymentMethod      OrderPaymentMethod        `json:"paymentMethod"`
	PaymentReference   *string                   `json:"paymentReference"`
	PlacedAt           *time.Time                `json:"placedAt"`
	ReadyForShipmentAt *time.Time                `json:"readyForShipmentAt"`
	OrderedAt          time.Time                 `json:"orderedAt"`
	ShippedAt          *time.Time                `json:"shippedAt"`
	FulfilledAt        *time.Time                `json:"fulfilledAt"`
	CancelledAt        *time.Time                `json:"cancelledAt"`
	ReturnedAt         *time.Time                `json:"returnedAt"`
	PaidAt             *time.Time                `json:"paidAt"`
	FailedAt           *time.Time                `json:"failedAt"`
	RefundedAt         *time.Time                `json:"refundedAt"`
	Items              []*OrderItem              `json:"items,omitempty"`
	Notes              []*OrderNote              `json:"notes,omitempty"`
	CreatedAt          time.Time                 `json:"createdAt"`
	UpdatedAt          time.Time                 `json:"updatedAt"`
	DeletedAt          *time.Time                `json:"deletedAt"`
}

type OrderCustomerResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Email       *string `json:"email"`
	PhoneCode   *string `json:"phoneCode"`
	PhoneNumber *string `json:"phoneNumber"`
	AvatarURL   *string `json:"avatarUrl,omitempty"`
}

type OrderAddressResponse struct {
	ID          string  `json:"id"`
	Street      *string `json:"street"`
	City        string  `json:"city"`
	State       string  `json:"state"`
	ZipCode     *string `json:"zipCode"`
	CountryCode string  `json:"countryCode"`
	PhoneCode   string  `json:"phoneCode"`
	PhoneNumber string  `json:"phoneNumber"`
}

type OrderShippingZoneResponse struct {
	ID                    string          `json:"id"`
	Name                  string          `json:"name"`
	Countries             []string        `json:"countries"`
	Currency              string          `json:"currency"`
	ShippingCost          decimal.Decimal `json:"shippingCost"`
	FreeShippingThreshold decimal.Decimal `json:"freeShippingThreshold"`
}

type OrderItemResponse struct {
	ID        string                    `json:"id"`
	OrderID   string                    `json:"orderId"`
	ProductID string                    `json:"productId"`
	VariantID string                    `json:"variantId"`
	Quantity  int                       `json:"quantity"`
	Currency  string                    `json:"currency"`
	UnitPrice decimal.Decimal           `json:"unitPrice"`
	UnitCost  decimal.Decimal           `json:"unitCost"`
	TotalCost decimal.Decimal           `json:"totalCost"`
	Total     decimal.Decimal           `json:"total"`
	Product   *OrderItemProductResponse `json:"product,omitempty"`
	Variant   *OrderItemVariantResponse `json:"variant,omitempty"`
}

type OrderItemProductResponse struct {
	ID     string                       `json:"id"`
	Name   string                       `json:"name"`
	Photos inventory.AssetReferenceList `json:"photos"`
}

type OrderItemVariantResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
	SKU  string `json:"sku"`
}

type OrderNoteResponse struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"orderId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func toOrderResponses(orders []*Order) []*OrderResponse {
	out := make([]*OrderResponse, 0, len(orders))
	for _, ord := range orders {
		out = append(out, toOrderResponse(ord))
	}
	return out
}

func toOrderResponse(ord *Order) *OrderResponse {
	if ord == nil {
		return nil
	}

	return &OrderResponse{
		ID:                 ord.ID,
		OrderNumber:        ord.OrderNumber,
		BusinessID:         ord.BusinessID,
		CustomerID:         ord.CustomerID,
		Customer:           ord.Customer,
		ShippingAddressID:  ord.ShippingAddressID,
		ShippingAddress:    ord.ShippingAddress,
		ShippingZoneID:     ord.ShippingZoneID,
		ShippingZone:       ord.ShippingZone,
		Channel:            ord.Channel,
		Subtotal:           ord.Subtotal,
		VAT:                ord.VAT,
		VATRate:            ord.VATRate,
		ShippingFee:        ord.ShippingFee,
		Discount:           ord.Discount,
		COGS:               ord.COGS,
		Total:              ord.Total,
		Currency:           ord.Currency,
		Status:             ord.Status,
		PaymentStatus:      ord.PaymentStatus,
		PaymentMethod:      ord.PaymentMethod,
		PaymentReference:   transformer.NullStringPtr(ord.PaymentReference),
		PlacedAt:           transformer.NullTimePtr(ord.PlacedAt),
		ReadyForShipmentAt: transformer.NullTimePtr(ord.ReadyForShipmentAt),
		OrderedAt:          ord.OrderedAt,
		ShippedAt:          transformer.NullTimePtr(ord.ShippedAt),
		FulfilledAt:        transformer.NullTimePtr(ord.FulfilledAt),
		CancelledAt:        transformer.NullTimePtr(ord.CancelledAt),
		ReturnedAt:         transformer.NullTimePtr(ord.ReturnedAt),
		PaidAt:             transformer.NullTimePtr(ord.PaidAt),
		FailedAt:           transformer.NullTimePtr(ord.FailedAt),
		RefundedAt:         transformer.NullTimePtr(ord.RefundedAt),
		Items:              ord.Items,
		Notes:              ord.Notes,
		CreatedAt:          ord.CreatedAt,
		UpdatedAt:          ord.UpdatedAt,
		DeletedAt:          deletedAtPtr(ord.DeletedAt),
	}
}

func deletedAtPtr(da gorm.DeletedAt) *time.Time {
	if !da.Valid {
		return nil
	}
	t := da.Time
	return &t
}
