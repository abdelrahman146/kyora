package inventory

import (
	"github.com/abdelrahman146/kyora/internal/platform/types/asset"
	"github.com/shopspring/decimal"
)

// CreateProductRequest is the request DTO for creating a product.
type CreateProductRequest struct {
	Name        string                 `json:"name" binding:"required"`
	Description string                 `json:"description" binding:"omitempty"`
	Photos      []asset.AssetReference `json:"photos" binding:"omitempty,max=10,dive"`
	CategoryID  string                 `json:"categoryId" binding:"required"`
}

// UpdateProductRequest is the request DTO for updating a product.
type UpdateProductRequest struct {
	Name        string                 `json:"name" binding:"omitempty"`
	Description string                 `json:"description" binding:"omitempty"`
	Photos      []asset.AssetReference `json:"photos" binding:"omitempty,max=10,dive"`
	CategoryID  string                 `json:"categoryId" binding:"omitempty"`
}

// CreateVariantRequest is the request DTO for creating a variant.
type CreateVariantRequest struct {
	ProductID          string                 `form:"productId" json:"productId" binding:"required"`
	Code               string                 `form:"code" json:"code" binding:"required"`
	SKU                string                 `form:"sku" json:"sku" binding:"omitempty"`
	Photos             []asset.AssetReference `form:"photos" json:"photos" binding:"omitempty,max=10,dive"`
	CostPrice          *decimal.Decimal       `form:"costPrice" json:"costPrice" binding:"required"`
	SalePrice          *decimal.Decimal       `form:"salePrice" json:"salePrice" binding:"required"`
	StockQuantity      *int                   `form:"stockQuantity" json:"stockQuantity" binding:"required,gte=0"`
	StockQuantityAlert *int                   `form:"stockQuantityAlert" json:"stockQuantityAlert" binding:"required,gte=0"`
}

// UpdateVariantRequest is the request DTO for updating a variant.
type UpdateVariantRequest struct {
	Code               *string                `form:"code" json:"code" binding:"omitempty"`
	SKU                *string                `form:"sku" json:"sku" binding:"omitempty"`
	Photos             []asset.AssetReference `form:"photos" json:"photos" binding:"omitempty,max=10,dive"`
	CostPrice          *decimal.Decimal       `form:"costPrice" json:"costPrice" binding:"omitempty"`
	SalePrice          *decimal.Decimal       `form:"salePrice" json:"salePrice" binding:"omitempty"`
	Currency           *string                `form:"currency" json:"currency" binding:"omitempty,len=3"`
	StockQuantity      *int                   `form:"stockQuantity" json:"stockQuantity" binding:"omitempty,gte=0"`
	StockQuantityAlert *int                   `form:"stockQuantityAlert" json:"stockQuantityAlert" binding:"omitempty,gte=0"`
}

// CreateCategoryRequest is the request DTO for creating a category.
type CreateCategoryRequest struct {
	Name       string `json:"name" binding:"required"`
	Descriptor string `json:"descriptor" binding:"required"`
}

// UpdateCategoryRequest is the request DTO for updating a category.
type UpdateCategoryRequest struct {
	Name       string `json:"name" binding:"omitempty"`
	Descriptor string `json:"descriptor" binding:"omitempty"`
}
