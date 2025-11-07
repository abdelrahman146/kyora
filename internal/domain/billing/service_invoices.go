package billing

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
	stripelib "github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/invoice"
	"github.com/stripe/stripe-go/v83/invoiceitem"
)

// InvoiceSummary is a lightweight view of a Stripe invoice for UI consumption
type InvoiceSummary struct {
	ID               string     `json:"id"`
	Number           string     `json:"number"`
	Status           string     `json:"status"`
	Currency         string     `json:"currency"`
	AmountDue        int64      `json:"amountDue"`
	AmountPaid       int64      `json:"amountPaid"`
	CreatedAt        time.Time  `json:"createdAt"`
	DueDate          *time.Time `json:"dueDate,omitempty"`
	HostedInvoiceURL string     `json:"hostedInvoiceUrl,omitempty"`
	InvoicePDF       string     `json:"invoicePdf,omitempty"`
}

// ListInvoices returns invoice summaries for the workspace's customer
func (s *Service) ListInvoices(ctx context.Context, ws *account.Workspace, status string) ([]InvoiceSummary, error) {
	custID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		return nil, err
	}
	params := &stripelib.InvoiceListParams{Customer: stripelib.String(custID)}
	switch status {
	case string(stripelib.InvoiceStatusOpen):
		params.Status = stripelib.String(string(stripelib.InvoiceStatusOpen))
	case string(stripelib.InvoiceStatusPaid):
		params.Status = stripelib.String(string(stripelib.InvoiceStatusPaid))
	default:
		// all - no status filter
	}
	params.Limit = stripelib.Int64(50)
	it := invoice.List(params)
	res := make([]InvoiceSummary, 0)
	for it.Next() {
		inv := it.Invoice()
		var due *time.Time
		if inv.DueDate != 0 {
			t := time.Unix(inv.DueDate, 0)
			due = &t
		}
		res = append(res, InvoiceSummary{
			ID:               inv.ID,
			Number:           inv.Number,
			Status:           string(inv.Status),
			Currency:         string(inv.Currency),
			AmountDue:        inv.AmountDue,
			AmountPaid:       inv.AmountPaid,
			CreatedAt:        time.Unix(inv.Created, 0),
			DueDate:          due,
			HostedInvoiceURL: inv.HostedInvoiceURL,
			InvoicePDF:       inv.InvoicePDF,
		})
	}
	return res, nil
}

// DownloadInvoiceURL returns the downloadable PDF URL for an invoice if it belongs to the customer's workspace
func (s *Service) DownloadInvoiceURL(ctx context.Context, ws *account.Workspace, invoiceID string) (string, error) {
	custID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		return "", err
	}
	inv, err := invoice.Get(invoiceID, nil)
	if err != nil {
		return "", err
	}
	if inv.Customer == nil || inv.Customer.ID != custID {
		return "", ErrSubscriptionNotFound(err, ws.ID)
	}
	if inv.InvoicePDF == "" && inv.HostedInvoiceURL == "" {
		return "", ErrSubscriptionNotFound(err, invoiceID)
	}
	if inv.InvoicePDF != "" {
		return inv.InvoicePDF, nil
	}
	return inv.HostedInvoiceURL, nil
}

// PayInvoice attempts to pay an open invoice for the workspace's customer
func (s *Service) PayInvoice(ctx context.Context, ws *account.Workspace, invoiceID string) error {
	custID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		return err
	}
	inv, err := invoice.Get(invoiceID, nil)
	if err != nil {
		return err
	}
	if inv.Customer == nil || inv.Customer.ID != custID {
		return ErrSubscriptionNotFound(err, ws.ID)
	}
	// If invoice is draft, finalize first
	if inv.Status == stripelib.InvoiceStatusDraft {
		if _, err := invoice.FinalizeInvoice(invoiceID, nil); err != nil {
			return err
		}
	}
	_, err = invoice.Pay(invoiceID, &stripelib.InvoicePayParams{})
	return err
}

// CreateInvoice creates a new invoice for the workspace
func (s *Service) CreateInvoice(ctx context.Context, ws *account.Workspace, description string, amount int64, currency string, dueDate *string) (*stripelib.Invoice, error) {
	l := logger.FromContext(ctx).With("workspace_id", ws.ID, "amount", amount, "currency", currency)
	l.Info("creating invoice")

	customerID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		l.Error("failed to ensure customer", "error", err)
		return nil, ErrCustomerCreationFailed(ws.ID, err)
	}

	params := &stripelib.InvoiceParams{
		Customer:    stripelib.String(customerID),
		Currency:    stripelib.String(currency),
		Description: stripelib.String(description),
		AutoAdvance: stripelib.Bool(true),
	}

	if dueDate != nil {
		// Parse due date and set it
		if dueTime, err := time.Parse("2006-01-02", *dueDate); err == nil {
			params.DueDate = stripelib.Int64(dueTime.Unix())
		}
	}

	// Create invoice item first
	invItemParams := &stripelib.InvoiceItemParams{
		Customer:    stripelib.String(customerID),
		Amount:      stripelib.Int64(amount),
		Currency:    stripelib.String(currency),
		Description: stripelib.String(description),
	}

	_, err = invoiceitem.New(invItemParams)
	if err != nil {
		l.Error("failed to create invoice item", "error", err)
		return nil, ErrStripeOperationFailed(err, "create_invoice_item")
	}

	// Create the invoice
	inv, err := invoice.New(params)
	if err != nil {
		l.Error("failed to create invoice", "error", err)
		return nil, ErrStripeOperationFailed(err, "create_invoice")
	}

	l.Info("invoice created successfully", "invoice_id", inv.ID)
	return inv, nil
}
