package business

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/types/date"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/abdelrahman146/kyora/internal/platform/types/schema"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const (
	BusinessTable  = "businesses"
	BusinessStruct = "Business"
	BusinessPrefix = "bus"
)

type Business struct {
	gorm.Model
	ID            string             `gorm:"column:id;primaryKey;type:text" json:"id"`
	Descriptor    string             `gorm:"column:descriptor;type:text;uniqueIndex:idx_workspace_id_descriptor" json:"descriptor"`
	WorkspaceID   string             `gorm:"column:workspace_id;type:text;uniqueIndex:idx_workspace_id_descriptor" json:"workspaceId"`
	Workspace     *account.Workspace `gorm:"foreignKey:WorkspaceID;references:ID" json:"workspace,omitempty"`
	Name          string             `gorm:"column:name;type:text" json:"name"`
	CountryCode   string             `gorm:"column:country_code;type:text" json:"countryCode"`
	Currency      string             `gorm:"column:currency;type:text" json:"currency"`
	VatRate       decimal.Decimal    `gorm:"column:vat_rate;type:numeric;not null;default:0" json:"vatRate"`
	SafetyBuffer  decimal.Decimal    `gorm:"column:safety_buffer;type:numeric;not null;default:0" json:"safetyBuffer"`
	EstablishedAt time.Time          `gorm:"column:established_at;type:date;default:now()" json:"establishedAt,omitempty"`
	ArchivedAt    *time.Time         `gorm:"column:archived_at;type:timestamp with time zone" json:"archivedAt,omitempty"`
}

func (m *Business) TableName() string {
	return BusinessTable
}

func (m *Business) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(BusinessPrefix)
	}
	return nil
}

var BusinessSchema = struct {
	ID            schema.Field
	Descriptor    schema.Field
	WorkspaceID   schema.Field
	Name          schema.Field
	CountryCode   schema.Field
	Currency      schema.Field
	VatRate       schema.Field
	SafetyBuffer  schema.Field
	EstablishedAt schema.Field
	ArchivedAt    schema.Field
	CreatedAt     schema.Field
	UpdatedAt     schema.Field
	DeletedAt     schema.Field
}{
	ID:            schema.NewField("id", "id"),
	Descriptor:    schema.NewField("descriptor", "descriptor"),
	WorkspaceID:   schema.NewField("workspace_id", "workspaceId"),
	Name:          schema.NewField("name", "name"),
	CountryCode:   schema.NewField("country_code", "countryCode"),
	Currency:      schema.NewField("currency", "currency"),
	VatRate:       schema.NewField("vat_rate", "vatRate"),
	SafetyBuffer:  schema.NewField("safety_buffer", "safetyBuffer"),
	EstablishedAt: schema.NewField("established_at", "establishedAt"),
	ArchivedAt:    schema.NewField("archived_at", "archivedAt"),
	CreatedAt:     schema.NewField("created_at", "createdAt"),
	UpdatedAt:     schema.NewField("updated_at", "updatedAt"),
	DeletedAt:     schema.NewField("deleted_at", "deletedAt"),
}

type CreateBusinessInput struct {
	Name          string          `form:"name" json:"name" binding:"required"`
	Descriptor    string          `form:"descriptor" json:"descriptor" binding:"required"`
	CountryCode   string          `form:"countryCode" json:"countryCode" binding:"required,len=2"`
	Currency      string          `form:"currency" json:"currency" binding:"required,len=3"`
	VatRate       decimal.Decimal `form:"vatRate" json:"vatRate" binding:"omitempty"`
	SafetyBuffer  decimal.Decimal `form:"safetyBuffer" json:"safetyBuffer" binding:"omitempty"`
	EstablishedAt date.Date       `form:"establishedAt" json:"establishedAt" binding:"omitempty"`
}

type UpdateBusinessInput struct {
	Name          *string             `form:"name" json:"name" binding:"omitempty"`
	Descriptor    *string             `form:"descriptor" json:"descriptor" binding:"omitempty"`
	CountryCode   *string             `form:"countryCode" json:"countryCode" binding:"omitempty,len=2"`
	Currency      *string             `form:"currency" json:"currency" binding:"omitempty,len=3"`
	VatRate       decimal.NullDecimal `form:"vatRate" json:"vatRate" binding:"omitempty"`
	SafetyBuffer  decimal.NullDecimal `form:"safetyBuffer" json:"safetyBuffer" binding:"omitempty"`
	EstablishedAt *date.Date          `form:"establishedAt" json:"establishedAt" binding:"omitempty"`
}

const (
	ShippingZoneTable  = "shipping_zones"
	ShippingZoneStruct = "ShippingZone"
	ShippingZonePrefix = "sz"
)

// CountryCodeList is a JSONB-backed list of ISO 3166-1 alpha-2 country codes.
// Stored as JSONB to keep reads fast and avoid join overhead.
type CountryCodeList []string

func (c CountryCodeList) Value() (driver.Value, error) {
	if c == nil {
		c = []string{}
	}
	b, err := json.Marshal([]string(c))
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (c *CountryCodeList) Scan(value any) error {
	if c == nil {
		return problem.InternalError().WithError(errors.New("CountryCodeList scan into nil receiver"))
	}
	if value == nil {
		*c = CountryCodeList{}
		return nil
	}
	switch v := value.(type) {
	case []byte:
		var out []string
		if err := json.Unmarshal(v, &out); err != nil {
			return err
		}
		*c = CountryCodeList(out)
		return nil
	case string:
		var out []string
		if err := json.Unmarshal([]byte(v), &out); err != nil {
			return err
		}
		*c = CountryCodeList(out)
		return nil
	default:
		return problem.InternalError().WithError(errors.New("unexpected scan type for CountryCodeList"))
	}
}

func (c CountryCodeList) Contains(countryCode string) bool {
	for _, x := range c {
		if x == countryCode {
			return true
		}
	}
	return false
}

// ShippingZone maps destination countries to shipping pricing rules.
// This is part of the business domain (not a separate domain) to stay DRY/KISS.
type ShippingZone struct {
	gorm.Model
	ID                    string          `gorm:"column:id;primaryKey;type:text" json:"id"`
	BusinessID            string          `gorm:"column:business_id;type:text;not null;index;uniqueIndex:idx_shipping_zone_business_name" json:"businessId"`
	Business              *Business       `gorm:"foreignKey:BusinessID;references:ID" json:"business,omitempty"`
	Name                  string          `gorm:"column:name;type:text;not null;uniqueIndex:idx_shipping_zone_business_name" json:"name"`
	Countries             CountryCodeList `gorm:"column:countries;type:jsonb;not null;default:'[]'" json:"countries"`
	Currency              string          `gorm:"column:currency;type:text;not null" json:"currency"`
	ShippingCost          decimal.Decimal `gorm:"column:shipping_cost;type:numeric;not null;default:0" json:"shippingCost"`
	FreeShippingThreshold decimal.Decimal `gorm:"column:free_shipping_threshold;type:numeric;not null;default:0" json:"freeShippingThreshold"`
}

func (m *ShippingZone) TableName() string { return ShippingZoneTable }

func (m *ShippingZone) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(ShippingZonePrefix)
	}
	return nil
}

var ShippingZoneSchema = struct {
	ID                    schema.Field
	BusinessID            schema.Field
	Name                  schema.Field
	Countries             schema.Field
	Currency              schema.Field
	ShippingCost          schema.Field
	FreeShippingThreshold schema.Field
	CreatedAt             schema.Field
	UpdatedAt             schema.Field
	DeletedAt             schema.Field
}{
	ID:                    schema.NewField("id", "id"),
	BusinessID:            schema.NewField("business_id", "businessId"),
	Name:                  schema.NewField("name", "name"),
	Countries:             schema.NewField("countries", "countries"),
	Currency:              schema.NewField("currency", "currency"),
	ShippingCost:          schema.NewField("shipping_cost", "shippingCost"),
	FreeShippingThreshold: schema.NewField("free_shipping_threshold", "freeShippingThreshold"),
	CreatedAt:             schema.NewField("created_at", "createdAt"),
	UpdatedAt:             schema.NewField("updated_at", "updatedAt"),
	DeletedAt:             schema.NewField("deleted_at", "deletedAt"),
}

type CreateShippingZoneRequest struct {
	Name                  string          `json:"name" binding:"required"`
	Countries             []string        `json:"countries" binding:"required,min=1,max=50,dive,required,len=2"`
	ShippingCost          decimal.Decimal `json:"shippingCost" binding:"omitempty"`
	FreeShippingThreshold decimal.Decimal `json:"freeShippingThreshold" binding:"omitempty"`
}

type UpdateShippingZoneRequest struct {
	Name                  *string             `json:"name" binding:"omitempty"`
	Countries             []string            `json:"countries" binding:"omitempty,min=1,max=50,dive,required,len=2"`
	ShippingCost          decimal.NullDecimal `json:"shippingCost" binding:"omitempty"`
	FreeShippingThreshold decimal.NullDecimal `json:"freeShippingThreshold" binding:"omitempty"`
}
