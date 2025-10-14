package owner

import (
	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const (
	InvestmentTable  = "investments"
	InvestmentStruct = "Investment"
	InvestmentAlias  = "inv"
)

type Investment struct {
	gorm.Model
	ID       string          `json:"id" gorm:"column:id;primaryKey;type:text"`
	StoreID  string          `json:"storeId" gorm:"column:store_id;type:text;not null;index"`
	Store    *store.Store    `json:"store,omitempty" gorm:"foreignKey:StoreID;references:ID"`
	OwnerID  string          `json:"ownerId" gorm:"column:owner_id;type:text;not null;index"`
	Owner    *Owner          `json:"owner,omitempty" gorm:"foreignKey:OwnerID;references:ID"`
	Name     string          `json:"name" gorm:"column:name;type:text;not null"`
	Amount   decimal.Decimal `json:"amount" gorm:"column:amount;type:numeric;not null"`
	Currency string          `json:"currency" gorm:"column:currency;type:text;not null;default:'USD'"`
	Note     string          `json:"note" gorm:"column:note;type:text"`
}

func (m *Investment) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = utils.ID.NewUlidWithPrefix(InvestmentAlias)
	}
	return
}

type CreateInvestmentRequest struct {
	Name    string          `form:"name" json:"name" binding:"required"`
	OwnerID string          `form:"ownerId" json:"ownerId" binding:"required"`
	Amount  decimal.Decimal `form:"amount" json:"amount" binding:"required"`
	Note    string          `form:"note" json:"note" binding:"omitempty"`
}

type UpdateInvestmentRequest struct {
	Name    string          `form:"name" json:"name" binding:"omitempty"`
	OwnerID string          `form:"ownerId" json:"ownerId" binding:"omitempty"`
	Amount  decimal.Decimal `form:"amount" json:"amount" binding:"omitempty"`
	Note    string          `form:"note" json:"note" binding:"omitempty"`
}
