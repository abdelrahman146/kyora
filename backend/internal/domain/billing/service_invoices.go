package billing

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
	"github.com/abdelrahman146/kyora/internal/platform/types/list"
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

// ListInvoices returns a ListResponse using limit/offset emulated over Stripe's cursor pagination.
// page & pageSize come from ListRequest; offset = (page-1)*pageSize. We fetch up to offset+pageSize+1 to compute hasMore.
func (s *Service) ListInvoices(ctx context.Context, ws *account.Workspace, status string, req *list.ListRequest) *list.ListResponse[InvoiceSummary] {
	custID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		return list.NewListResponse[InvoiceSummary](nil, req.Page(), req.PageSize(), 0, false)
	}
	pageSize := req.PageSize()
	if pageSize <= 0 {
		pageSize = 30
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := req.Offset()
	target := offset + pageSize + 1 // +1 to detect hasMore
	collected := make([]InvoiceSummary, 0)
	var lastID string
	var startingAfter *string
	remaining := target
	for remaining > 0 {
		// Stripe limit per request capped at 100
		batchLimit := remaining
		if batchLimit > 100 {
			batchLimit = 100
		}
		params := &stripelib.InvoiceListParams{Customer: stripelib.String(custID)}
		switch status {
		case string(stripelib.InvoiceStatusOpen):
			params.Status = stripelib.String(string(stripelib.InvoiceStatusOpen))
		case string(stripelib.InvoiceStatusPaid):
			params.Status = stripelib.String(string(stripelib.InvoiceStatusPaid))
		}
		params.Limit = stripelib.Int64(int64(batchLimit))
		if startingAfter != nil {
			params.StartingAfter = startingAfter
		}
		it := invoice.List(params)
		batchCount := 0
		for it.Next() {
			inv := it.Invoice()
			lastID = inv.ID
			batchCount++
			var due *time.Time
			if inv.DueDate != 0 {
				t := time.Unix(inv.DueDate, 0)
				due = &t
			}
			collected = append(collected, InvoiceSummary{
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
			if len(collected) >= target {
				break
			}
		}
		if err := it.Err(); err != nil {
			logger.FromContext(ctx).Error("invoice list iteration error", "error", err, "workspaceId", ws.ID)
			break
		}
		if len(collected) >= target || batchCount == 0 {
			break
		}
		// prepare next cursor
		if lastID != "" {
			startingAfter = stripelib.String(lastID)
		}
		remaining = target - len(collected)
	}
	// Slice to page window
	var window []InvoiceSummary
	if offset >= len(collected) {
		window = []InvoiceSummary{}
	} else {
		end := offset + pageSize
		if end > len(collected) {
			end = len(collected)
		}
		window = collected[offset:end]
	}
	// hasMore if we fetched at least offset+pageSize+1
	hasMore := len(collected) > offset+pageSize
	return list.NewListResponse(window, req.Page(), pageSize, 0, hasMore)
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
		if _, err := withStripeRetry[*stripelib.Invoice](ctx, 3, func() (*stripelib.Invoice, error) {
			return invoice.FinalizeInvoice(invoiceID, nil)
		}); err != nil {
			return err
		}
	}
	_, err = withStripeRetry[*stripelib.Invoice](ctx, 3, func() (*stripelib.Invoice, error) {
		return invoice.Pay(invoiceID, &stripelib.InvoicePayParams{})
	})
	return err
}

// CreateInvoice creates a new invoice for the workspace
func (s *Service) CreateInvoice(ctx context.Context, ws *account.Workspace, description string, amount int64, currency string, dueDate *string) (*stripelib.Invoice, error) {
	l := logger.FromContext(ctx).With("workspaceId", ws.ID, "amount", amount, "currency", currency)
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

	_, err = withStripeRetry[*stripelib.InvoiceItem](ctx, 3, func() (*stripelib.InvoiceItem, error) {
		return invoiceitem.New(invItemParams)
	})
	if err != nil {
		l.Error("failed to create invoice item", "error", err)
		return nil, ErrStripeOperationFailed(err, "create_invoice_item")
	}

	// Create the invoice
	inv, err := withStripeRetry[*stripelib.Invoice](ctx, 3, func() (*stripelib.Invoice, error) {
		return invoice.New(params)
	})
	if err != nil {
		l.Error("failed to create invoice", "error", err)
		return nil, ErrStripeOperationFailed(err, "create_invoice")
	}

	l.Info("invoice created successfully", "invoiceId", inv.ID)
	return inv, nil
}
