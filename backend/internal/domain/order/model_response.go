package order

import (
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/platform/types/asset"
	"github.com/abdelrahman146/kyora/internal/platform/utils/transformer"
	"github.com/shopspring/decimal"
)

// OrderResponse is the API response for Order entity
// All date fields use time.Time (not sql.NullTime)
// No DeletedAt field (GORM leakage removed)
type OrderResponse struct {
	ID                 string                            `json:"id"`
	OrderNumber        string                            `json:"orderNumber"`
	BusinessID         string                            `json:"businessId"`
	CustomerID         string                            `json:"customerId"`
	Customer           *customer.CustomerResponse        `json:"customer,omitempty"`
	ShippingAddressID  string                            `json:"shippingAddressId"`
	ShippingAddress    *customer.CustomerAddressResponse `json:"shippingAddress,omitempty"`
	ShippingZoneID     *string                           `json:"shippingZoneId,omitempty"`
	ShippingZone       *business.ShippingZoneResponse    `json:"shippingZone,omitempty"`
	Channel            string                            `json:"channel"`
	Subtotal           decimal.Decimal                   `json:"subtotal"`
	VAT                decimal.Decimal                   `json:"vat"`
	VATRate            decimal.Decimal                   `json:"vatRate"`
	ShippingFee        decimal.Decimal                   `json:"shippingFee"`
	Discount           decimal.Decimal                   `json:"discount"`
	DiscountType       DiscountType                      `json:"discountType,omitempty"`
	DiscountValue      decimal.Decimal                   `json:"discountValue,omitempty"`
	COGS               decimal.Decimal                   `json:"cogs"`
	Total              decimal.Decimal                   `json:"total"`
	Currency           string                            `json:"currency"`
	Status             OrderStatus                       `json:"status"`
	PaymentStatus      OrderPaymentStatus                `json:"paymentStatus"`
	PaymentMethod      OrderPaymentMethod                `json:"paymentMethod"`
	PaymentReference   *string                           `json:"paymentReference,omitempty"`
	PlacedAt           *time.Time                        `json:"placedAt,omitempty"`
	ReadyForShipmentAt *time.Time                        `json:"readyForShipmentAt,omitempty"`
	OrderedAt          time.Time                         `json:"orderedAt"`
	ShippedAt          *time.Time                        `json:"shippedAt,omitempty"`
	FulfilledAt        *time.Time                        `json:"fulfilledAt,omitempty"`
	CancelledAt        *time.Time                        `json:"cancelledAt,omitempty"`
	ReturnedAt         *time.Time                        `json:"returnedAt,omitempty"`
	PaidAt             *time.Time                        `json:"paidAt,omitempty"`
	FailedAt           *time.Time                        `json:"failedAt,omitempty"`
	RefundedAt         *time.Time                        `json:"refundedAt,omitempty"`
	Items              []OrderItemResponse               `json:"items,omitempty"`
	Notes              []OrderNoteResponse               `json:"notes,omitempty"`
	CreatedAt          time.Time                         `json:"createdAt"`
	UpdatedAt          time.Time                         `json:"updatedAt"`
}

// OrderItemResponse is the API response for OrderItem entity
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
	CreatedAt time.Time                 `json:"createdAt"`
	UpdatedAt time.Time                 `json:"updatedAt"`
}

// OrderItemProductResponse is a simplified product representation in order items
type OrderItemProductResponse struct {
	ID     string                 `json:"id"`
	Name   string                 `json:"name"`
	Photos []asset.AssetReference `json:"photos"`
}

// OrderItemVariantResponse is a simplified variant representation in order items
type OrderItemVariantResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
	SKU  string `json:"sku"`
}

// OrderNoteResponse is the API response for OrderNote entity
type OrderNoteResponse struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"orderId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// OrderPreviewResponse is the API response for order preview calculations.
type OrderPreviewResponse struct {
	Subtotal       decimal.Decimal            `json:"subtotal"`
	VAT            decimal.Decimal            `json:"vat"`
	VATRate        decimal.Decimal            `json:"vatRate"`
	ShippingFee    decimal.Decimal            `json:"shippingFee"`
	Discount       decimal.Decimal            `json:"discount"`
	COGS           decimal.Decimal            `json:"cogs"`
	Total          decimal.Decimal            `json:"total"`
	Currency       string                     `json:"currency"`
	ShippingZoneID *string                    `json:"shippingZoneId,omitempty"`
	PaymentMethod  OrderPaymentMethod         `json:"paymentMethod"`
	Items          []OrderPreviewItemResponse `json:"items"`
}

// OrderPreviewItemResponse is the per-item breakdown in preview responses.
type OrderPreviewItemResponse struct {
	VariantID string          `json:"variantId"`
	ProductID string          `json:"productId"`
	Quantity  int             `json:"quantity"`
	UnitPrice decimal.Decimal `json:"unitPrice"`
	UnitCost  decimal.Decimal `json:"unitCost"`
	Total     decimal.Decimal `json:"total"`
	TotalCost decimal.Decimal `json:"totalCost"`
}

// ToOrderResponse converts Order model to OrderResponse
func ToOrderResponse(ord *Order) OrderResponse {
	if ord == nil {
		return OrderResponse{}
	}

	var customerResp *customer.CustomerResponse
	if ord.Customer != nil {
		// Note: ordersCount and totalSpent are not computed here
		// They should be computed separately for customer detail views
		resp := customer.ToCustomerResponse(ord.Customer, 0, 0.0)
		customerResp = &resp
	}

	var shippingAddressResp *customer.CustomerAddressResponse
	if ord.ShippingAddress != nil {
		resp := customer.ToCustomerAddressResponse(ord.ShippingAddress)
		shippingAddressResp = &resp
	}

	var shippingZoneResp *business.ShippingZoneResponse
	if ord.ShippingZone != nil {
		resp := business.ToShippingZoneResponse(ord.ShippingZone)
		shippingZoneResp = &resp
	}

	items := ToOrderItemResponses(ord.Items)
	notes := ToOrderNoteResponses(ord.Notes)

	return OrderResponse{
		ID:                 ord.ID,
		OrderNumber:        ord.OrderNumber,
		BusinessID:         ord.BusinessID,
		CustomerID:         ord.CustomerID,
		Customer:           customerResp,
		ShippingAddressID:  ord.ShippingAddressID,
		ShippingAddress:    shippingAddressResp,
		ShippingZoneID:     ord.ShippingZoneID,
		ShippingZone:       shippingZoneResp,
		Channel:            ord.Channel,
		Subtotal:           ord.Subtotal,
		VAT:                ord.VAT,
		VATRate:            ord.VATRate,
		ShippingFee:        ord.ShippingFee,
		Discount:           ord.Discount,
		DiscountType:       ord.DiscountType,
		DiscountValue:      ord.DiscountValue,
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
		Items:              items,
		Notes:              notes,
		CreatedAt:          ord.CreatedAt,
		UpdatedAt:          ord.UpdatedAt,
	}
}

// ToOrderResponses converts a slice of Order models to responses
func ToOrderResponses(orders []*Order) []OrderResponse {
	responses := make([]OrderResponse, len(orders))
	for i, ord := range orders {
		responses[i] = ToOrderResponse(ord)
	}
	return responses
}

// ToOrderItemResponse converts OrderItem model to OrderItemResponse
func ToOrderItemResponse(item *OrderItem) OrderItemResponse {
	if item == nil {
		return OrderItemResponse{}
	}

	var productResp *OrderItemProductResponse
	if item.Product != nil {
		photos := []asset.AssetReference{}
		if item.Product.Photos != nil {
			photos = []asset.AssetReference(item.Product.Photos)
		}
		productResp = &OrderItemProductResponse{
			ID:     item.Product.ID,
			Name:   item.Product.Name,
			Photos: photos,
		}
	}

	var variantResp *OrderItemVariantResponse
	if item.Variant != nil {
		variantResp = &OrderItemVariantResponse{
			ID:   item.Variant.ID,
			Name: item.Variant.Name,
			Code: item.Variant.Code,
			SKU:  item.Variant.SKU,
		}
	}

	return OrderItemResponse{
		ID:        item.ID,
		OrderID:   item.OrderID,
		ProductID: item.ProductID,
		VariantID: item.VariantID,
		Quantity:  item.Quantity,
		Currency:  item.Currency,
		UnitPrice: item.UnitPrice,
		UnitCost:  item.UnitCost,
		TotalCost: item.TotalCost,
		Total:     item.Total,
		Product:   productResp,
		Variant:   variantResp,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

// ToOrderItemResponses converts a slice of OrderItem models to responses
func ToOrderItemResponses(items []*OrderItem) []OrderItemResponse {
	responses := make([]OrderItemResponse, len(items))
	for i, item := range items {
		responses[i] = ToOrderItemResponse(item)
	}
	return responses
}

// ToOrderNoteResponse converts OrderNote model to OrderNoteResponse
func ToOrderNoteResponse(note *OrderNote) OrderNoteResponse {
	if note == nil {
		return OrderNoteResponse{}
	}

	return OrderNoteResponse{
		ID:        note.ID,
		OrderID:   note.OrderID,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}

// ToOrderNoteResponses converts a slice of OrderNote models to responses
func ToOrderNoteResponses(notes []*OrderNote) []OrderNoteResponse {
	responses := make([]OrderNoteResponse, len(notes))
	for i, note := range notes {
		responses[i] = ToOrderNoteResponse(note)
	}
	return responses
}

// ToOrderPreviewResponse maps an OrderPreview to its API response shape.
func ToOrderPreviewResponse(preview *OrderPreview) OrderPreviewResponse {
	if preview == nil {
		return OrderPreviewResponse{}
	}

	return OrderPreviewResponse{
		Subtotal:       preview.Subtotal,
		VAT:            preview.VAT,
		VATRate:        preview.VATRate,
		ShippingFee:    preview.ShippingFee,
		Discount:       preview.Discount,
		COGS:           preview.COGS,
		Total:          preview.Total,
		Currency:       preview.Currency,
		ShippingZoneID: preview.ShippingZoneID,
		PaymentMethod:  preview.PaymentMethod,
		Items:          ToOrderPreviewItemResponses(preview.Items),
	}
}

// ToOrderPreviewItemResponses maps preview items to API response items.
func ToOrderPreviewItemResponses(items []OrderPreviewItem) []OrderPreviewItemResponse {
	responses := make([]OrderPreviewItemResponse, len(items))
	for i, item := range items {
		responses[i] = OrderPreviewItemResponse{
			VariantID: item.VariantID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			UnitCost:  item.UnitCost,
			Total:     item.Total,
			TotalCost: item.TotalCost,
		}
	}
	return responses
}
