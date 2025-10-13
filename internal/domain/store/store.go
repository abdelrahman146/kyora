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
	Currency       string                `gorm:"column:currency;type:text;not null;default:'USD'" json:"currency"`
	CountryCode    string                `gorm:"column:country_code;type:text" json:"country_code"`
	VatRate        decimal.Decimal       `gorm:"column:vat_rate;type:numeric;not null;default:0" json:"vatRate"`
	SafetyBuffer   decimal.Decimal       `gorm:"column:safety_buffer;type:numeric;default:0" json:"safetyBuffer"`
}

func (m *Store) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = utils.ID.NewUlidWithPrefix(StoreAlias)
	}
	return
}

type CreateStoreRequest struct {
	Name         string          `json:"name" binding:"required"`
	Currency     string          `json:"currency" binding:"omitempty,len=3"`
	CountryCode  string          `json:"countryCode" binding:"omitempty,len=2"`
	VatRate      decimal.Decimal `json:"vatRate" binding:"omitempty,gte=0,lte=100"`
	SafetyBuffer decimal.Decimal `json:"safetyBuffer" binding:"omitempty,gte=0,lte=100"`
}

type UpdateStoreRequest struct {
	Name         string          `json:"name" binding:"omitempty"`
	Currency     string          `json:"currency" binding:"omitempty,len=3"`
	CountryCode  string          `json:"countryCode" binding:"omitempty,len=2"`
	VatRate      decimal.Decimal `json:"vatRate" binding:"omitempty,gte=0,lte=100"`
	SafetyBuffer decimal.Decimal `json:"safetyBuffer" binding:"omitempty,gte=0,lte=100"`
}
