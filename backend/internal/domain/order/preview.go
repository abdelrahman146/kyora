package order

import "github.com/shopspring/decimal"

// OrderPreview represents a dry-run calculation for an order without persisting or mutating inventory.
type OrderPreview struct {
	Subtotal       decimal.Decimal    `json:"subtotal"`
	VAT            decimal.Decimal    `json:"vat"`
	VATRate        decimal.Decimal    `json:"vatRate"`
	ShippingFee    decimal.Decimal    `json:"shippingFee"`
	Discount       decimal.Decimal    `json:"discount"`
	COGS           decimal.Decimal    `json:"cogs"`
	Total          decimal.Decimal    `json:"total"`
	Currency       string             `json:"currency"`
	ShippingZoneID *string            `json:"shippingZoneId,omitempty"`
	PaymentMethod  OrderPaymentMethod `json:"paymentMethod"`
	Items          []OrderPreviewItem `json:"items"`
}

// OrderPreviewItem is the per-item breakdown used in preview responses.
type OrderPreviewItem struct {
	VariantID string          `json:"variantId"`
	ProductID string          `json:"productId"`
	Quantity  int             `json:"quantity"`
	UnitPrice decimal.Decimal `json:"unitPrice"`
	UnitCost  decimal.Decimal `json:"unitCost"`
	Total     decimal.Decimal `json:"total"`
	TotalCost decimal.Decimal `json:"totalCost"`
}
