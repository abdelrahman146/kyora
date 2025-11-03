package business

import (
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
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
	VatRate       decimal.Decimal `form:"vatRate" json:"vatRate" binding:"omitempty,decimal"`
	SafetyBuffer  decimal.Decimal `form:"safetyBuffer" json:"safetyBuffer" binding:"omitempty,decimal"`
	EstablishedAt time.Time       `form:"establishedAt" json:"establishedAt" binding:"omitempty,date"`
}

type UpdateBusinessInput struct {
	Name          string              `form:"name" json:"name" binding:"required"`
	Descriptor    string              `form:"descriptor" json:"descriptor" binding:"required"`
	CountryCode   string              `form:"countryCode" json:"countryCode" binding:"required,len=2"`
	Currency      string              `form:"currency" json:"currency" binding:"required,len=3"`
	VatRate       decimal.NullDecimal `form:"vatRate" json:"vatRate" binding:"omitempty,decimal"`
	SafetyBuffer  decimal.NullDecimal `form:"safetyBuffer" json:"safetyBuffer" binding:"omitempty,decimal"`
	EstablishedAt time.Time           `form:"establishedAt" json:"establishedAt" binding:"omitempty,date"`
}
