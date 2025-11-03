package inventory

import (
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/types/schema"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

/* Product Model */
//---------------*/

const (
	ProductTable  = "products"
	ProductStruct = "Product"
	ProductPrefix = "prod"
)

type Product struct {
	gorm.Model
	ID          string             `gorm:"column:id;primaryKey;type:text" json:"id"`
	BusinessID  string             `gorm:"column:business_id;type:text;not null;index" json:"businessId"`
	Business    *business.Business `gorm:"foreignKey:BusinessID;references:ID" json:"business,omitempty"`
	Name        string             `gorm:"column:name;type:text;not null;index:product_trgm_idx,type:gin,option:gin_trgm_ops" json:"name"`
	Description string             `gorm:"column:description;type:text" json:"description"`
	CategoryID  string             `gorm:"column:category_id;type:text;index" json:"categoryId"`
	Category    *Category          `gorm:"foreignKey:CategoryID;references:ID" json:"category,omitempty"`
	Variants    []*Variant         `gorm:"foreignKey:ProductID;references:ID" json:"variants,omitempty"`
}

func (m *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(ProductPrefix)
	}
	return
}

type CreateProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"omitempty"`
	CategoryID  string `json:"categoryId" binding:"required"`
}

type UpdateProductRequest struct {
	Name        string `json:"name" binding:"omitempty"`
	Description string `json:"description" binding:"omitempty"`
	CategoryID  string `json:"categoryId" binding:"omitempty"`
}

var ProductSchema = struct {
	ID          schema.Field
	BusinessID  schema.Field
	Name        schema.Field
	Description schema.Field
	CategoryID  schema.Field
	CreatedAt   schema.Field
	UpdatedAt   schema.Field
	DeletedAt   schema.Field
}{
	ID:          schema.NewField("id", "id"),
	BusinessID:  schema.NewField("business_id", "businessId"),
	Name:        schema.NewField("name", "name"),
	Description: schema.NewField("description", "description"),
	CategoryID:  schema.NewField("category_id", "categoryId"),
}

func CreateProductSKU(businessDescriptor, productName, variantCode string) string {
	s := id.NewCodeFromString(businessDescriptor, 3)
	s += "-" + id.NewCodeFromString(productName, 3)
	s += "-" + id.NewCodeFromString(variantCode, 3)
	rand, _ := id.RandomString(4)
	s += "-" + rand
	return s
}

/* Variant Model */
//---------------*/

const (
	VariantTable  = "variants"
	VariantStruct = "Variant"
	VariantPrefix = "var"
)

type Variant struct {
	gorm.Model
	ID                 string             `gorm:"column:id;primaryKey;type:text" json:"id"`
	BusinessID         string             `gorm:"column:business_id;type:text;not null;index;uniqueIndex:sku_business_idx" json:"businessId"`
	Business           *business.Business `gorm:"foreignKey:BusinessID;references:ID" json:"business,omitempty"`
	Name               string             `gorm:"column:name;type:text;not null;index:variant_trgm_idx,type:gin,option:gin_trgm_ops" json:"name"`
	Code               string             `gorm:"column:code;type:text;not null;uniqueIndex:code_product_idx" json:"code"`
	ProductID          string             `gorm:"column:product_id;type:text;not null;index;uniqueIndex:code_product_idx" json:"productId"`
	Product            *Product           `gorm:"foreignKey:ProductID;references:ID;OnDelete:CASCADE;" json:"product,omitempty"`
	SKU                string             `gorm:"column:sku;type:text;not null;uniqueIndex:sku_store_idx" json:"sku"`
	CostPrice          decimal.Decimal    `gorm:"column:cost_price;type:numeric;not null;default:0" json:"costPrice"`
	SalePrice          decimal.Decimal    `gorm:"column:sale_price;type:numeric;not null;default:0" json:"salePrice"`
	Currency           string             `gorm:"column:currency;type:text;not null;default:'USD'" json:"currency"`
	StockQuantity      int                `gorm:"column:stock_quantity;type:int;not null;default:0" json:"stockQuantity"`
	StockQuantityAlert int                `gorm:"column:stock_alert;type:int;not null;default:0" json:"stockQuantityAlert"`
}

func (m *Variant) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(VariantPrefix)
	}
	return
}

type CreateVariantRequest struct {
	Code               string          `form:"code" json:"code" binding:"required"`
	SKU                string          `form:"sku" json:"sku" binding:"omitempty"`
	CostPrice          decimal.Decimal `form:"costPrice" json:"costPrice" binding:"required,gte=0"`
	SalePrice          decimal.Decimal `form:"salePrice" json:"salePrice" binding:"required,gte=0"`
	StockQuantity      int             `form:"stockQuantity" json:"stockQuantity" binding:"required,gte=0"`
	StockQuantityAlert int             `form:"stockQuantityAlert" json:"stockQuantityAlert" binding:"required,gte=0"`
}

type UpdateVariantRequest struct {
	Code               string          `form:"code" json:"code" binding:"omitempty"`
	SKU                string          `form:"sku" json:"sku" binding:"omitempty"`
	CostPrice          decimal.Decimal `form:"costPrice" json:"costPrice" binding:"omitempty,gte=0"`
	SalePrice          decimal.Decimal `form:"salePrice" json:"salePrice" binding:"omitempty,gte=0"`
	Currency           string          `form:"currency" json:"currency" binding:"omitempty,len=3"`
	StockQuantity      int             `form:"stockQuantity" json:"stockQuantity" binding:"omitempty,gte=0"`
	StockQuantityAlert int             `form:"stockQuantityAlert" json:"stockQuantityAlert" binding:"omitempty,gte=0"`
}

var VariantSchema = struct {
	ID                 schema.Field
	BusinessID         schema.Field
	Name               schema.Field
	Code               schema.Field
	ProductID          schema.Field
	SKU                schema.Field
	CostPrice          schema.Field
	SalePrice          schema.Field
	Currency           schema.Field
	StockQuantity      schema.Field
	StockQuantityAlert schema.Field
	CreatedAt          schema.Field
	UpdatedAt          schema.Field
	DeletedAt          schema.Field
}{
	ID:                 schema.NewField("id", "id"),
	BusinessID:         schema.NewField("business_id", "businessId"),
	Name:               schema.NewField("name", "name"),
	Code:               schema.NewField("code", "code"),
	ProductID:          schema.NewField("product_id", "productId"),
	SKU:                schema.NewField("sku", "sku"),
	CostPrice:          schema.NewField("cost_price", "costPrice"),
	SalePrice:          schema.NewField("sale_price", "salePrice"),
	Currency:           schema.NewField("currency", "currency"),
	StockQuantity:      schema.NewField("stock_quantity", "stockQuantity"),
	StockQuantityAlert: schema.NewField("stock_alert", "stockQuantityAlert"),
	CreatedAt:          schema.NewField("created_at", "createdAt"),
	UpdatedAt:          schema.NewField("updated_at", "updatedAt"),
	DeletedAt:          schema.NewField("deleted_at", "deletedAt"),
}

/* Category Model */
//-----------------*/

const (
	CategoryTable  = "categories"
	CategoryStruct = "Category"
	CategoryPrefix = "cat"
)

type Category struct {
	gorm.Model
	ID         string             `gorm:"column:id;primaryKey;type:text" json:"id"`
	BusinessID string             `gorm:"column:business_id;type:text;not null;index" json:"businessId"`
	Business   *business.Business `gorm:"foreignKey:BusinessID;references:ID" json:"business,omitempty"`
	Name       string             `gorm:"column:name;type:text;not null;index:category_trgm_idx,type:gin,option:gin_trgm_ops" json:"name"`
	Descriptor string             `gorm:"column:descriptor;type:text;not null;uniqueIndex:descriptor_business_idx" json:"descriptor"`
}

func (m *Category) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(CategoryPrefix)
	}
	return
}

type CreateCategoryRequest struct {
	Name       string `json:"name" binding:"required"`
	Descriptor string `json:"descriptor" binding:"required"`
}

type UpdateCategoryRequest struct {
	Name       string `json:"name" binding:"omitempty"`
	Descriptor string `json:"descriptor" binding:"omitempty"`
}

var CategorySchema = struct {
	ID         schema.Field
	BusinessID schema.Field
	Name       schema.Field
	Descriptor schema.Field
	CreatedAt  schema.Field
	UpdatedAt  schema.Field
	DeletedAt  schema.Field
}{
	ID:         schema.NewField("id", "id"),
	BusinessID: schema.NewField("business_id", "businessId"),
	Name:       schema.NewField("name", "name"),
	Descriptor: schema.NewField("descriptor", "descriptor"),
	CreatedAt:  schema.NewField("created_at", "createdAt"),
	UpdatedAt:  schema.NewField("updated_at", "updatedAt"),
	DeletedAt:  schema.NewField("deleted_at", "deletedAt"),
}
