package inventory

import (
	"time"

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

type ProductFilter struct {
	IDs         []string  `json:"ids" binding:"omitempty,dive,required"`
	Tags        []string  `json:"tags" binding:"omitempty,dive,required"`
	SearchQuery string    `json:"searchQuery" binding:"omitempty"`
	From        time.Time `json:"from" binding:"omitempty"`
	To          time.Time `json:"to" binding:"omitempty"`
}
