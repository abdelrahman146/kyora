package inventory

import (
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/types/asset"
	"github.com/shopspring/decimal"
)

// ProductResponse is the API response for Product entity
type ProductResponse struct {
	ID          string                 `json:"id"`
	BusinessID  string                 `json:"businessId"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Photos      []asset.AssetReference `json:"photos"`
	CategoryID  string                 `json:"categoryId"`
	Variants    []VariantResponse      `json:"variants,omitempty"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

// ToProductResponse converts Product model to ProductResponse
func ToProductResponse(p *Product) ProductResponse {
	photos := []asset.AssetReference{}
	if p.Photos != nil {
		photos = []asset.AssetReference(p.Photos)
	}

	var variants []VariantResponse
	if p.Variants != nil {
		variants = make([]VariantResponse, len(p.Variants))
		for i, v := range p.Variants {
			variants[i] = ToVariantResponse(v)
		}
	}

	return ProductResponse{
		ID:          p.ID,
		BusinessID:  p.BusinessID,
		Name:        p.Name,
		Description: p.Description,
		Photos:      photos,
		CategoryID:  p.CategoryID,
		Variants:    variants,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// ToProductResponses converts a slice of Product models to responses
func ToProductResponses(products []*Product) []ProductResponse {
	responses := make([]ProductResponse, len(products))
	for i, p := range products {
		responses[i] = ToProductResponse(p)
	}
	return responses
}

// VariantResponse is the API response for Variant entity
type VariantResponse struct {
	ID                 string                 `json:"id"`
	BusinessID         string                 `json:"businessId"`
	Name               string                 `json:"name"`
	Code               string                 `json:"code"`
	ProductID          string                 `json:"productId"`
	SKU                string                 `json:"sku"`
	CostPrice          decimal.Decimal        `json:"costPrice"`
	SalePrice          decimal.Decimal        `json:"salePrice"`
	Currency           string                 `json:"currency"`
	Photos             []asset.AssetReference `json:"photos"`
	StockQuantity      int                    `json:"stockQuantity"`
	StockQuantityAlert int                    `json:"stockQuantityAlert"`
	CreatedAt          time.Time              `json:"createdAt"`
	UpdatedAt          time.Time              `json:"updatedAt"`
}

// ToVariantResponse converts Variant model to VariantResponse
func ToVariantResponse(v *Variant) VariantResponse {
	photos := []asset.AssetReference{}
	if v.Photos != nil {
		photos = []asset.AssetReference(v.Photos)
	}

	return VariantResponse{
		ID:                 v.ID,
		BusinessID:         v.BusinessID,
		Name:               v.Name,
		Code:               v.Code,
		ProductID:          v.ProductID,
		SKU:                v.SKU,
		CostPrice:          v.CostPrice,
		SalePrice:          v.SalePrice,
		Currency:           v.Currency,
		Photos:             photos,
		StockQuantity:      v.StockQuantity,
		StockQuantityAlert: v.StockQuantityAlert,
		CreatedAt:          v.CreatedAt,
		UpdatedAt:          v.UpdatedAt,
	}
}

// ToVariantResponses converts a slice of Variant models to responses
func ToVariantResponses(variants []*Variant) []VariantResponse {
	responses := make([]VariantResponse, len(variants))
	for i, v := range variants {
		responses[i] = ToVariantResponse(v)
	}
	return responses
}

// CategoryResponse is the API response for Category entity
type CategoryResponse struct {
	ID         string    `json:"id"`
	BusinessID string    `json:"businessId"`
	Name       string    `json:"name"`
	Descriptor string    `json:"descriptor"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// ToCategoryResponse converts Category model to CategoryResponse
func ToCategoryResponse(c *Category) CategoryResponse {
	return CategoryResponse{
		ID:         c.ID,
		BusinessID: c.BusinessID,
		Name:       c.Name,
		Descriptor: c.Descriptor,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}
}

// ToCategoryResponses converts a slice of Category models to responses
func ToCategoryResponses(categories []*Category) []CategoryResponse {
	responses := make([]CategoryResponse, len(categories))
	for i, c := range categories {
		responses[i] = ToCategoryResponse(c)
	}
	return responses
}

// TopProductByInventoryValue represents a product with its computed inventory value
type TopProductByInventoryValue struct {
	Product        ProductResponse `json:"product"`
	InventoryValue decimal.Decimal `json:"inventoryValue"`
}

// ToTopProductByInventoryValue converts Product with inventory value to TopProductByInventoryValue response
func ToTopProductByInventoryValue(product *Product, inventoryValue decimal.Decimal) TopProductByInventoryValue {
	return TopProductByInventoryValue{
		Product:        ToProductResponse(product),
		InventoryValue: inventoryValue,
	}
}
