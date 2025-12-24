package business

import (
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/abdelrahman146/kyora/internal/platform/types/schema"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const (
	BusinessPaymentMethodTable  = "business_payment_methods"
	BusinessPaymentMethodStruct = "BusinessPaymentMethod"
	BusinessPaymentMethodPrefix = "bpm"
)

// BusinessPaymentMethod stores per-business enablement and fee overrides for a global payment method.
// We keep it in the business domain (same as shipping zones) to stay DRY/KISS.
type BusinessPaymentMethod struct {
	gorm.Model
	ID         string                  `gorm:"column:id;primaryKey;type:text" json:"id"`
	BusinessID string                  `gorm:"column:business_id;type:text;not null;index;uniqueIndex:idx_business_payment_method" json:"businessId"`
	Business   *Business               `gorm:"foreignKey:BusinessID;references:ID" json:"business,omitempty"`
	Descriptor PaymentMethodDescriptor `gorm:"column:descriptor;type:text;not null;uniqueIndex:idx_business_payment_method" json:"descriptor"`
	Enabled    bool                    `gorm:"column:enabled;type:boolean;not null;default:false" json:"enabled"`
	FeePercent decimal.Decimal         `gorm:"column:fee_percent;type:numeric;not null;default:0" json:"feePercent"`
	FeeFixed   decimal.Decimal         `gorm:"column:fee_fixed;type:numeric;not null;default:0" json:"feeFixed"`
}

func (m *BusinessPaymentMethod) TableName() string { return BusinessPaymentMethodTable }

func (m *BusinessPaymentMethod) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(BusinessPaymentMethodPrefix)
	}
	return nil
}

var BusinessPaymentMethodSchema = struct {
	ID         schema.Field
	BusinessID schema.Field
	Descriptor schema.Field
	Enabled    schema.Field
	FeePercent schema.Field
	FeeFixed   schema.Field
	CreatedAt  schema.Field
	UpdatedAt  schema.Field
	DeletedAt  schema.Field
}{
	ID:         schema.NewField("id", "id"),
	BusinessID: schema.NewField("business_id", "businessId"),
	Descriptor: schema.NewField("descriptor", "descriptor"),
	Enabled:    schema.NewField("enabled", "enabled"),
	FeePercent: schema.NewField("fee_percent", "feePercent"),
	FeeFixed:   schema.NewField("fee_fixed", "feeFixed"),
	CreatedAt:  schema.NewField("created_at", "createdAt"),
	UpdatedAt:  schema.NewField("updated_at", "updatedAt"),
	DeletedAt:  schema.NewField("deleted_at", "deletedAt"),
}

type UpdateBusinessPaymentMethodRequest struct {
	Enabled    *bool               `json:"enabled" binding:"omitempty"`
	FeePercent decimal.NullDecimal `json:"feePercent" binding:"omitempty"`
	FeeFixed   decimal.NullDecimal `json:"feeFixed" binding:"omitempty"`
}

func (r *UpdateBusinessPaymentMethodRequest) Validate() error {
	if r == nil {
		return problem.BadRequest("request is required")
	}
	if r.FeePercent.Valid {
		if r.FeePercent.Decimal.LessThan(decimal.Zero) {
			return problem.BadRequest("feePercent cannot be negative").With("field", "feePercent")
		}
		if r.FeePercent.Decimal.GreaterThan(decimal.NewFromInt(1)) {
			return problem.BadRequest("feePercent must be between 0 and 1").With("field", "feePercent")
		}
	}
	if r.FeeFixed.Valid {
		if r.FeeFixed.Decimal.LessThan(decimal.Zero) {
			return problem.BadRequest("feeFixed cannot be negative").With("field", "feeFixed")
		}
	}
	return nil
}
