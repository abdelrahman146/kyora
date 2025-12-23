package billing

import (
	"github.com/abdelrahman146/kyora/internal/platform/types/schema"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"gorm.io/gorm"
)

const (
	InvoiceRecordTable  = "billing_invoice_records"
	InvoiceRecordStruct = "InvoiceRecord"
	InvoiceRecordPrefix = "binv"
)

// InvoiceRecord stores ownership mapping for manually-created invoices.
//
// This enables secure (BOLA-safe) access checks for invoice operations
// even when the billing provider responses are missing customer fields
// (e.g., in certain test/mocked environments).
type InvoiceRecord struct {
	gorm.Model
	ID               string `json:"id" gorm:"column:id;primaryKey;type:text"`
	WorkspaceID      string `json:"workspaceId" gorm:"column:workspace_id;type:text;index;not null"`
	StripeInvoiceID  string `json:"stripeInvoiceId" gorm:"column:stripe_invoice_id;type:text;uniqueIndex;not null"`
	HostedInvoiceURL string `json:"hostedInvoiceUrl" gorm:"column:hosted_invoice_url;type:text"`
	InvoicePDF       string `json:"invoicePdf" gorm:"column:invoice_pdf;type:text"`
}

func (m *InvoiceRecord) TableName() string {
	return InvoiceRecordTable
}

func (m *InvoiceRecord) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(InvoiceRecordPrefix)
	}
	return nil
}

var InvoiceRecordSchema = struct {
	ID               schema.Field
	WorkspaceID      schema.Field
	StripeInvoiceID  schema.Field
	HostedInvoiceURL schema.Field
	InvoicePDF       schema.Field
	CreatedAt        schema.Field
	UpdatedAt        schema.Field
	DeletedAt        schema.Field
}{
	ID:               schema.NewField("id", "id"),
	WorkspaceID:      schema.NewField("workspace_id", "workspaceId"),
	StripeInvoiceID:  schema.NewField("stripe_invoice_id", "stripeInvoiceId"),
	HostedInvoiceURL: schema.NewField("hosted_invoice_url", "hostedInvoiceUrl"),
	InvoicePDF:       schema.NewField("invoice_pdf", "invoicePdf"),
	CreatedAt:        schema.NewField("created_at", "createdAt"),
	UpdatedAt:        schema.NewField("updated_at", "updatedAt"),
	DeletedAt:        schema.NewField("deleted_at", "deletedAt"),
}
