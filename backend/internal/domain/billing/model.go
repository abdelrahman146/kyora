package billing

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/types/schema"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

/* Plan Model */
//-----------------*/

const (
	PlanTable  = "plans"
	PlanStruct = "Plan"
	PlanPrefix = "plan"
)

type BillingCycle string

const (
	BillingCycleMonthly BillingCycle = "monthly"
	BillingCycleYearly  BillingCycle = "yearly"
)

func (bc BillingCycle) GetMonths() int {
	switch bc {
	case BillingCycleMonthly:
		return 1
	case BillingCycleYearly:
		return 12
	default:
		return 0
	}
}

type PlanFeature struct {
	CustomerManagement       bool `json:"customerManagement"`
	InventoryManagement      bool `json:"inventoryManagement"`
	OrderManagement          bool `json:"orderManagement"`
	ExpenseManagement        bool `json:"expenseManagement"`
	Accounting               bool `json:"accounting"` // owner draws, investments, loans tracking
	BasicAnalytics           bool `json:"basicAnalytics"`
	FinancialReports         bool `json:"financialReports"`
	DataImport               bool `json:"dataImport"`
	DataExport               bool `json:"dataExport"`
	AdvancedAnalytics        bool `json:"advancedAnalytics"`
	AdvancedFinancialReports bool `json:"advancedFinancialReports"`
	OrderPaymentLinks        bool `json:"orderPaymentLinks"` // integrate with stripe and auto generate shareable payment links for orders
	InvoiceGeneration        bool `json:"invoiceGeneration"` // Up to here pro plan
	ExportAnalyticsData      bool `json:"exportAnalyticsData"`
	AIBusinessAssistant      bool `json:"aiBusinessAssistant"` // Up to here premium plan
}

// Scan implements the Scanner interface for JSONB deserialization
func (pf *PlanFeature) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan PlanFeature: value is not []byte")
	}
	return json.Unmarshal(bytes, pf)
}

// Value implements the Valuer interface for JSONB serialization
func (pf PlanFeature) Value() (driver.Value, error) {
	return json.Marshal(pf)
}

// CanUseFeature checks boolean (on/off) features.
func (pf *PlanFeature) CanUseFeature(feature schema.Field) error {
	// A map defines all boolean features and their current values.
	boolFeatures := map[schema.Field]bool{
		PlanSchema.CustomerManagement:       pf.CustomerManagement,
		PlanSchema.InventoryManagement:      pf.InventoryManagement,
		PlanSchema.OrderManagement:          pf.OrderManagement,
		PlanSchema.ExpenseManagement:        pf.ExpenseManagement,
		PlanSchema.Accounting:               pf.Accounting,
		PlanSchema.BasicAnalytics:           pf.BasicAnalytics,
		PlanSchema.FinancialReports:         pf.FinancialReports,
		PlanSchema.DataImport:               pf.DataImport,
		PlanSchema.DataExport:               pf.DataExport,
		PlanSchema.AdvancedAnalytics:        pf.AdvancedAnalytics,
		PlanSchema.AdvancedFinancialReports: pf.AdvancedFinancialReports,
		PlanSchema.OrderPaymentLinks:        pf.OrderPaymentLinks,
		PlanSchema.InvoiceGeneration:        pf.InvoiceGeneration,
		PlanSchema.ExportAnalyticsData:      pf.ExportAnalyticsData,
		PlanSchema.AIBusinessAssistant:      pf.AIBusinessAssistant,
	}

	enabled, ok := boolFeatures[feature]
	if !ok {
		return ErrUnknownFeature(nil, feature)
	}
	if !enabled {
		return ErrFeatureNotAvailable(nil, feature)
	}
	return nil
}

type PlanLimit struct {
	MaxOrdersPerMonth int64 `json:"maxOrdersPerMonth"`
	MaxTeamMembers    int64 `json:"maxTeamMembers"`
	MaxBusinesses     int64 `json:"maxBusinesses"`
}

// Scan implements the Scanner interface for JSONB deserialization
func (pl *PlanLimit) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan PlanLimit: value is not []byte")
	}
	return json.Unmarshal(bytes, pl)
}

// Value implements the Valuer interface for JSONB serialization
func (pl PlanLimit) Value() (driver.Value, error) {
	return json.Marshal(pl)
}

func (pl *PlanLimit) CheckUsageLimit(feature schema.Field, currentUsage int64) error {
	limitFeatures := map[schema.Field]int64{
		PlanSchema.MaxOrdersPerMonth: pl.MaxOrdersPerMonth,
		PlanSchema.MaxTeamMembers:    pl.MaxTeamMembers,
		PlanSchema.MaxBusinesses:     pl.MaxBusinesses,
	}
	limit, ok := limitFeatures[feature]
	if !ok {
		return ErrUnknownFeature(nil, feature)
	}
	newUsage := currentUsage + 1
	if newUsage > limit {
		return ErrFeatureMaxLimitReached(nil, feature, limit)
	}
	return nil
}

type Plan struct {
	gorm.Model
	ID           string          `json:"id" gorm:"column:id;primaryKey;type:text"`
	Descriptor   string          `json:"descriptor" gorm:"column:descriptor;type:text;not null;unique"`
	Name         string          `json:"name" gorm:"column:name;type:text;not null"`
	Description  string          `json:"description" gorm:"column:description;type:text;not null"`
	StripePlanID *string         `json:"stripePlanId" gorm:"column:stripe_plan_id;type:text;unique"`
	Price        decimal.Decimal `json:"price" gorm:"column:price;type:decimal(10,2);not null"`
	Currency     string          `json:"currency" gorm:"column:currency;type:text;not null"`
	BillingCycle BillingCycle    `json:"billingCycle" gorm:"column:billing_cycle;type:text;not null"`
	Features     PlanFeature     `json:"features" gorm:"column:features;type:jsonb;not null"`
	Limits       PlanLimit       `json:"limits" gorm:"column:limits;type:jsonb;not null"`
}

func (m *Plan) TableName() string {
	return PlanTable
}

func (m *Plan) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(PlanPrefix)
	}
	return nil
}

var PlanSchema = struct {
	ID           schema.Field
	Descriptor   schema.Field
	Name         schema.Field
	Description  schema.Field
	StripePlanID schema.Field
	Price        schema.Field
	Currency     schema.Field
	BillingCycle schema.Field
	CreatedAt    schema.Field
	UpdatedAt    schema.Field
	DeletedAt    schema.Field
	// Features (jsonb: features)
	CustomerManagement       schema.Field
	InventoryManagement      schema.Field
	OrderManagement          schema.Field
	ExpenseManagement        schema.Field
	Accounting               schema.Field
	BasicAnalytics           schema.Field
	FinancialReports         schema.Field
	DataImport               schema.Field
	DataExport               schema.Field
	AdvancedAnalytics        schema.Field
	AdvancedFinancialReports schema.Field
	OrderPaymentLinks        schema.Field
	InvoiceGeneration        schema.Field
	ExportAnalyticsData      schema.Field
	AIBusinessAssistant      schema.Field
	// Limits (jsonb: limits)
	MaxOrdersPerMonth schema.Field
	MaxTeamMembers    schema.Field
	MaxBusinesses     schema.Field
}{
	ID:           schema.NewField("id", "id"),
	Descriptor:   schema.NewField("descriptor", "descriptor"),
	Name:         schema.NewField("name", "name"),
	Description:  schema.NewField("description", "description"),
	StripePlanID: schema.NewField("stripe_plan_id", "stripePlanId"),
	Price:        schema.NewField("price", "price"),
	Currency:     schema.NewField("currency", "currency"),
	BillingCycle: schema.NewField("billing_cycle", "billingCycle"),
	CreatedAt:    schema.NewField("created_at", "createdAt"),
	UpdatedAt:    schema.NewField("updated_at", "updatedAt"),
	DeletedAt:    schema.NewField("deleted_at", "deletedAt"),
	// Features
	CustomerManagement:       schema.NewField("features->customerManagement", "features.customerManagement"),
	InventoryManagement:      schema.NewField("features->inventoryManagement", "features.inventoryManagement"),
	OrderManagement:          schema.NewField("features->orderManagement", "features.orderManagement"),
	ExpenseManagement:        schema.NewField("features->expenseManagement", "features.expenseManagement"),
	Accounting:               schema.NewField("features->accounting", "features.accounting"),
	BasicAnalytics:           schema.NewField("features->basicAnalytics", "features.basicAnalytics"),
	FinancialReports:         schema.NewField("features->financialReports", "features.financialReports"),
	DataImport:               schema.NewField("features->dataImport", "features.dataImport"),
	DataExport:               schema.NewField("features->dataExport", "features.dataExport"),
	AdvancedAnalytics:        schema.NewField("features->advancedAnalytics", "features.advancedAnalytics"),
	AdvancedFinancialReports: schema.NewField("features->advancedFinancialReports", "features.advancedFinancialReports"),
	OrderPaymentLinks:        schema.NewField("features->orderPaymentLinks", "features.orderPaymentLinks"),
	InvoiceGeneration:        schema.NewField("features->invoiceGeneration", "features.invoiceGeneration"),
	ExportAnalyticsData:      schema.NewField("features->exportAnalyticsData", "features.exportAnalyticsData"),
	AIBusinessAssistant:      schema.NewField("features->aiBusinessAssistant", "features.aiBusinessAssistant"),
	// Limits
	MaxOrdersPerMonth: schema.NewField("limits->maxOrdersPerMonth", "limits.maxOrdersPerMonth"),
	MaxTeamMembers:    schema.NewField("limits->maxTeamMembers", "limits.maxTeamMembers"),
	MaxBusinesses:     schema.NewField("limits->maxBusinesses", "limits.maxBusinesses"),
}

/* Subscription Model */
//-----------------*/

const (
	SubscriptionTable  = "subscriptions"
	SubscriptionStruct = "Subscription"
	SubscriptionPrefix = "sub"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive     SubscriptionStatus = "active"
	SubscriptionStatusPastDue    SubscriptionStatus = "past_due"
	SubscriptionStatusCanceled   SubscriptionStatus = "canceled"
	SubscriptionStatusUnpaid     SubscriptionStatus = "unpaid"
	SubscriptionStatusIncomplete SubscriptionStatus = "incomplete"
	SubscriptionStatusTrialing   SubscriptionStatus = "trialing"
)

type Subscription struct {
	gorm.Model
	ID               string             `json:"id" gorm:"column:id;primaryKey;type:text"`
	WorkspaceID      string             `json:"workspaceId" gorm:"column:workspace_id;type:text;not null;index"`
	PlanID           string             `json:"planId" gorm:"column:plan_id;type:text;not null;index"`
	Plan             *Plan              `json:"plan,omitempty" gorm:"foreignKey:PlanID;references:ID"`
	StripeSubID      string             `json:"stripeSubId" gorm:"column:stripe_sub_id;type:text;not null;unique"`
	CurrentPeriodEnd time.Time          `json:"currentPeriodEnd" gorm:"column:current_period_end;type:timestamp;not null"`
	Status           SubscriptionStatus `json:"status" gorm:"column:status;type:text;not null;index"`
}

func (m *Subscription) TableName() string {
	return SubscriptionTable
}

func (m *Subscription) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(SubscriptionPrefix)
	}
	return nil
}

func (m *Subscription) IsActive() error {
	if m.Status != SubscriptionStatusActive {
		return ErrSubscriptionNotActive(nil)
	}
	return nil
}

var SubscriptionSchema = struct {
	ID               schema.Field
	WorkspaceID      schema.Field
	PlanID           schema.Field
	StripeSubID      schema.Field
	StartedAt        schema.Field
	CurrentPeriodEnd schema.Field
	Status           schema.Field
	CreatedAt        schema.Field
	UpdatedAt        schema.Field
	DeletedAt        schema.Field
}{
	ID:               schema.NewField("id", "id"),
	WorkspaceID:      schema.NewField("workspace_id", "workspaceId"),
	PlanID:           schema.NewField("plan_id", "planId"),
	StripeSubID:      schema.NewField("stripe_sub_id", "stripeSubId"),
	StartedAt:        schema.NewField("started_at", "startedAt"),
	CurrentPeriodEnd: schema.NewField("current_period_end", "currentPeriodEnd"),
	Status:           schema.NewField("status", "status"),
	CreatedAt:        schema.NewField("created_at", "createdAt"),
	UpdatedAt:        schema.NewField("updated_at", "updatedAt"),
	DeletedAt:        schema.NewField("deleted_at", "deletedAt"),
}

// Request/Response DTOs
type attachPMRequest struct {
	PaymentMethodID string `json:"paymentMethodId" binding:"required"`
}

type subRequest struct {
	PlanDescriptor string `json:"planDescriptor" binding:"required"`
}

type checkoutRequest struct {
	PlanDescriptor string `json:"planDescriptor" binding:"required"`
	SuccessURL     string `json:"successUrl" binding:"required,url"`
	CancelURL      string `json:"cancelUrl" binding:"required,url"`
}

type billingPortalRequest struct {
	ReturnURL string `json:"returnUrl" binding:"required,url"`
}

// Additional request DTOs (new endpoints)
type scheduleChangeRequest struct {
	PlanDescriptor string `json:"planDescriptor" binding:"required"`
	EffectiveDate  string `json:"effectiveDate" binding:"required"`  // ISO8601 ("2006-01-02" or "2006-01-02T15:04:05Z")
	ProrationMode  string `json:"prorationMode" binding:"omitempty"` // Stripe proration behavior ("create_prorations" | "none")
}

type prorationEstimateRequest struct {
	NewPlanDescriptor string `json:"newPlanDescriptor" binding:"required"`
}

type resumeSubscriptionRequest struct {
	// currently no fields; placeholder for future (e.g. payment method enforcement)
}

type manualInvoiceRequest struct {
	Description string  `json:"description" binding:"required"`
	Amount      int64   `json:"amount" binding:"required,min=1"` // amount in minor units (e.g., cents)
	Currency    string  `json:"currency" binding:"required"`
	DueDate     *string `json:"dueDate" binding:"omitempty"` // YYYY-MM-DD
}

type trialExtendRequest struct {
	AdditionalDays int `json:"additionalDays" binding:"required,min=1,max=30"`
}

type taxCalculateRequest struct {
	Amount   int64  `json:"amount" binding:"required,min=1"`
	Currency string `json:"currency" binding:"required"`
}
