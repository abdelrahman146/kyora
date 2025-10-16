package inventory

import (
	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

const (
	ProductTable  = "products"
	ProductStruct = "Product"
	ProductAlias  = "prod"
)

type Product struct {
	gorm.Model
	ID          string         `gorm:"column:id;primaryKey;type:text" json:"id"`
	StoreID     string         `gorm:"column:store_id;type:text;not null;index" json:"storeId"`
	Store       *store.Store   `gorm:"foreignKey:StoreID;references:ID" json:"store,omitempty"`
	Name        string         `gorm:"column:name;type:text;not null;index:product_trgm_idx,type:gin,option:gin_trgm_ops" json:"name"`
	Description string         `gorm:"column:description;type:text" json:"description"`
	Tags        pq.StringArray `gorm:"column:tags;type:text[]" json:"tags"`
	Variants    []*Variant     `gorm:"foreignKey:ProductID;references:ID" json:"variants,omitempty"`
}

func (m *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = utils.ID.NewUlidWithPrefix(ProductAlias)
	}
	return
}

type CreateProductRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"omitempty"`
	Tags        []string `json:"tags" binding:"omitempty,dive,required"`
}

type UpdateProductRequest struct {
	Name        string   `json:"name" binding:"omitempty"`
	Description string   `json:"description" binding:"omitempty"`
	Tags        []string `json:"tags" binding:"omitempty,dive,required"`
}
