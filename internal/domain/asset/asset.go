package asset

import (
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/govalues/decimal"
	"gorm.io/gorm"
)

const (
	AssetTable  = "assets"
	AssetStruct = "Asset"
	AssetAlias  = "ast"
)

type AssetType string

const (
	AssetTypeSoftware  AssetType = "software"
	AssetTypeEquipment AssetType = "equipment"
	AssetTypeVehicle   AssetType = "vehicle"
	AssetTypeFurniture AssetType = "furniture"
	AssetTypeOther     AssetType = "other"
)

type Asset struct {
	gorm.Model
	ID          string          `gorm:"column:id;primaryKey;type:text" json:"id"`
	StoreID     string          `gorm:"column:store_id;type:text;not null;index" json:"storeId"`
	Store       *store.Store    `gorm:"foreignKey:StoreID;references:ID" json:"store,omitempty"`
	Name        string          `gorm:"column:name;type:text;not null" json:"name"`
	Type        AssetType       `gorm:"column:type;type:text;not null" json:"type"`
	Value       decimal.Decimal `gorm:"column:value;type:numeric;not null" json:"value"`
	Currency    string          `gorm:"column:currency;type:text;not null;default:'USD'" json:"currency"`
	PurchasedAt time.Time       `gorm:"column:purchased_at;type:date;not null" json:"purchasedAt"`
	Note        string          `gorm:"column:note;type:text" json:"note"`
}

func (m *Asset) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = utils.ID.NewUlidWithPrefix(AssetAlias)
	}
	return
}

type CreateAssetRequest struct {
	Name        string          `form:"name" json:"name" binding:"required"`
	Type        AssetType       `form:"type" json:"type" binding:"required"`
	Value       decimal.Decimal `form:"value" json:"value" binding:"required"`
	Currency    string          `form:"currency" json:"currency" binding:"omitempty,len=3"`
	PurchasedAt *time.Time      `form:"purchasedAt" json:"purchasedAt" binding:"omitempty"`
	Note        string          `form:"note" json:"note" binding:"omitempty"`
}

type UpdateAssetRequest struct {
	Name        string          `form:"name" json:"name" binding:"omitempty"`
	Type        AssetType       `form:"type" json:"type" binding:"omitempty"`
	Value       decimal.Decimal `form:"value" json:"value" binding:"omitempty"`
	Currency    string          `form:"currency" json:"currency" binding:"omitempty,len=3"`
	PurchasedAt *time.Time      `form:"purchasedAt" json:"purchasedAt" binding:"omitempty"`
	Note        string          `form:"note" json:"note" binding:"omitempty"`
}

type AssetFilter struct {
	Types []AssetType `form:"types" json:"types" binding:"omitempty"`
}
