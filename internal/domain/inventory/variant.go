package inventory

import (
	"time"

	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/govalues/decimal"
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
	Name          string          `gorm:"column:name;type:text;not null;index:variant_trgm_idx,type:gin,option:gin_trgm_ops" json:"name"`
	ProductID     string          `gorm:"column:product_id;type:text;not null;index" json:"productId"`
	Product       *Product        `gorm:"foreignKey:ProductID;references:ID" json:"product,omitempty"`
	SKU           string          `gorm:"column:sku;type:text;not null;unique" json:"sku"`
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
	Name          string          `json:"name" binding:"required"`
	SKU           string          `json:"sku" binding:"required"`
	CostPrice     decimal.Decimal `json:"costPrice" binding:"required,gte=0"`
	SalePrice     decimal.Decimal `json:"salePrice" binding:"required,gte=0"`
	StockQuantity int             `json:"stockQuantity" binding:"required,gte=0"`
	StockAlert    int             `json:"stockAlert" binding:"required,gte=0"`
}

type UpdateVariantRequest struct {
	Name          string          `json:"name" binding:"omitempty"`
	SKU           string          `json:"sku" binding:"omitempty"`
	CostPrice     decimal.Decimal `json:"costPrice" binding:"omitempty,gte=0"`
	SalePrice     decimal.Decimal `json:"salePrice" binding:"omitempty,gte=0"`
	Currency      string          `json:"currency" binding:"omitempty,len=3"`
	StockQuantity int             `json:"stockQuantity" binding:"omitempty,gte=0"`
	StockAlert    int             `json:"stockAlert" binding:"omitempty,gte=0"`
}

type VariantFilter struct {
	IDs         []string  `json:"ids" binding:"omitempty,dive,required"`
	SKUs        []string  `json:"skus" binding:"omitempty,dive,required"`
	ProductIDs  []string  `json:"productIds" binding:"omitempty,dive,required"`
	ProductTags []string  `json:"productTags" binding:"omitempty,dive,required"`
	SearchQuery string    `json:"searchQuery" binding:"omitempty"`
	From        time.Time `json:"from" binding:"omitempty"`
	To          time.Time `json:"to" binding:"omitempty"`
}
