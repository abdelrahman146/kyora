package owner

import (
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const (
	OwnerDrawTable  = "owner_draws"
	OwnerDrawStruct = "OwnerDraw"
	OwnerDrawAlias  = "odr"
)

type OwnerDraw struct {
	gorm.Model
	ID          string          `json:"id" gorm:"column:id;primaryKey;type:text"`
	StoreID     string          `json:"storeId" gorm:"column:store_id;type:text;not null;index"`
	Store       *store.Store    `json:"store,omitempty" gorm:"foreignKey:StoreID;references:ID"`
	OwnerID     string          `json:"ownerId" gorm:"column:owner_id;type:text;not null;index"`
	Owner       *Owner          `json:"owner,omitempty" gorm:"foreignKey:OwnerID;references:ID"`
	Amount      decimal.Decimal `json:"amount" gorm:"column:amount;type:numeric;not null"`
	Currency    string          `json:"currency" gorm:"column:currency;type:text;not null;default:'USD'"`
	WithdrawnAt time.Time       `json:"withdrawnAt" gorm:"column:withdrawn_at;type:timestamptz;not null;default:now()"`
	Note        string          `json:"note" gorm:"column:note;type:text"`
}

func (m *OwnerDraw) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = utils.ID.NewUlidWithPrefix(OwnerDrawAlias)
	}
	return
}

type CreateOwnerDrawRequest struct {
	Name        string          `form:"name" json:"name" binding:"required"`
	OwnerID     string          `form:"ownerId" json:"ownerId" binding:"required"`
	Amount      decimal.Decimal `form:"amount" json:"amount" binding:"required"`
	Note        string          `form:"note" json:"note" binding:"omitempty"`
	WithdrawnAt time.Time       `form:"withdrawnAt" json:"withdrawnAt" binding:"omitempty"`
}

type UpdateOwnerDrawRequest struct {
	Name        string          `form:"name" json:"name" binding:"omitempty"`
	OwnerID     string          `form:"ownerId" json:"ownerId" binding:"omitempty"`
	Amount      decimal.Decimal `form:"amount" json:"amount" binding:"omitempty"`
	Note        string          `form:"note" json:"note" binding:"omitempty"`
	WithdrawnAt time.Time       `form:"withdrawnAt" json:"withdrawnAt" binding:"omitempty"`
}
