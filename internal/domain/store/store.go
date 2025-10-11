package store

import (
	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/govalues/decimal"
	"gorm.io/gorm"
)

const (
	StoreTable  = "stores"
	StoreStruct = "Store"
	StoreAlias  = "str"
)

type Store struct {
	gorm.Model
	ID             string                `gorm:"column:id;primaryKey;type:text" json:"id"`
	Code           string                `gorm:"column:code;type:text;not null;uniqueIndex:idx_code_organization_id" json:"code"`
	OrganizationID string                `gorm:"column:organization_id;type:text;not null;index;uniqueIndex:idx_name_organization_id;uniqueIndex:idx_code_organization_id" json:"organizationId"`
	Organization   *account.Organization `gorm:"foreignKey:OrganizationID;references:ID" json:"organization,omitempty"`
	Name           string                `gorm:"column:name;type:text;not null;uniqueIndex:idx_name_organization_id" json:"name"`
	Locale         string                `gorm:"column:locale;type:text;not null;default:'en'" json:"locale"`
	Currency       string                `gorm:"column:currency;type:text;not null;default:'USD'" json:"currency"`
	Timezone       string                `gorm:"column:timezone;type:text;not null;default:'UTC'" json:"timezone"`
	VATRate        decimal.Decimal       `gorm:"column:vat_rate;type:numeric;not null;default:0" json:"vatRate"`
}

func (m *Store) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = utils.ID.NewUlidWithPrefix(StoreAlias)
	}
	return
}

type CreateStoreRequest struct {
	Name     string          `json:"name" binding:"required"`
	Code     string          `json:"code" binding:"required,alphanum"`
	Locale   string          `json:"locale" binding:"omitempty"`
	Currency string          `json:"currency" binding:"omitempty,len=3"`
	Timezone string          `json:"timezone" binding:"omitempty"`
	VATRate  decimal.Decimal `json:"vatRate" binding:"omitempty,gte=0,lte=100"`
}

type UpdateStoreRequest struct {
	Name     string          `json:"name" binding:"omitempty"`
	Locale   string          `json:"locale" binding:"omitempty"`
	Currency string          `json:"currency" binding:"omitempty,len=3"`
	Timezone string          `json:"timezone" binding:"omitempty"`
	VATRate  decimal.Decimal `json:"vatRate" binding:"omitempty,gte=0,lte=100"`
}
