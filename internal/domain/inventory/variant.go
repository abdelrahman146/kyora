package inventory

import (
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const (
	VariantTable  = "variants"
	VariantStruct = "Variant"
	VariantAlias  = "var"
)

type Variant struct {
	gorm.Model
	ID            string          `gorm:"column:id;primaryKey;type:text" json:"id"`
	StoreID       string          `gorm:"column:store_id;type:text;not null;index;uniqueIndex:sku_store_idx" json:"storeId"`
	Store         *store.Store    `gorm:"foreignKey:StoreID;references:ID" json:"store,omitempty"`
	Name          string          `gorm:"column:name;type:text;not null;index:variant_trgm_idx,type:gin,option:gin_trgm_ops" json:"name"`
	Code          string          `gorm:"column:code;type:text;not null;uniqueIndex:code_product_idx" json:"code"`
	ProductID     string          `gorm:"column:product_id;type:text;not null;index;uniqueIndex:code_product_idx" json:"productId"`
	Product       *Product        `gorm:"foreignKey:ProductID;references:ID" json:"product,omitempty"`
	SKU           string          `gorm:"column:sku;type:text;not null;uniqueIndex:sku_store_idx" json:"sku"`
	CostPrice     decimal.Decimal `gorm:"column:cost_price;type:numeric;not null;default:0" json:"costPrice"`
	SalePrice     decimal.Decimal `gorm:"column:sale_price;type:numeric;not null;default:0" json:"salePrice"`
	Currency      string          `gorm:"column:currency;type:text;not null;default:'USD'" json:"currency"`
	StockQuantity int             `gorm:"column:stock_quantity;type:int;not null;default:0" json:"stockQuantity"`
	StockAlert    int             `gorm:"column:stock_alert;type:int;not null;default:0" json:"stockAlert"`
}

func (m *Variant) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = utils.ID.NewUlidWithPrefix(VariantAlias)
	}
	return
}

type CreateVariantRequest struct {
	Code          string          `form:"code" json:"code" binding:"required"`
	SKU           string          `form:"sku" json:"sku" binding:"omitempty"`
	CostPrice     decimal.Decimal `form:"costPrice" json:"costPrice" binding:"required,gte=0"`
	SalePrice     decimal.Decimal `form:"salePrice" json:"salePrice" binding:"required,gte=0"`
	StockQuantity int             `form:"stockQuantity" json:"stockQuantity" binding:"required,gte=0"`
	StockAlert    int             `form:"stockAlert" json:"stockAlert" binding:"required,gte=0"`
}

type UpdateVariantRequest struct {
	Code          string          `form:"code" json:"code" binding:"omitempty"`
	SKU           string          `form:"sku" json:"sku" binding:"omitempty"`
	CostPrice     decimal.Decimal `form:"costPrice" json:"costPrice" binding:"omitempty,gte=0"`
	SalePrice     decimal.Decimal `form:"salePrice" json:"salePrice" binding:"omitempty,gte=0"`
	Currency      string          `form:"currency" json:"currency" binding:"omitempty,len=3"`
	StockQuantity int             `form:"stockQuantity" json:"stockQuantity" binding:"omitempty,gte=0"`
	StockAlert    int             `form:"stockAlert" json:"stockAlert" binding:"omitempty,gte=0"`
}

type VariantFilter struct {
	IDs         []string  `form:"ids" json:"ids" binding:"omitempty,dive,required"`
	SKUs        []string  `form:"skus" json:"skus" binding:"omitempty,dive,required"`
	ProductIDs  []string  `form:"productIds" json:"productIds" binding:"omitempty,dive,required"`
	ProductTags []string  `form:"productTags" json:"productTags" binding:"omitempty,dive,required"`
	SearchQuery string    `form:"searchQuery" json:"searchQuery" binding:"omitempty"`
	From        time.Time `form:"from" json:"from" binding:"omitempty"`
	To          time.Time `form:"to" json:"to" binding:"omitempty"`
}
