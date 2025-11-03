package billing

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/bus"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/atomic"
	"github.com/shopspring/decimal"
	stripelib "github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/creditnote"
	"github.com/stripe/stripe-go/v83/customer"
	"github.com/stripe/stripe-go/v83/invoice"
	"github.com/stripe/stripe-go/v83/paymentmethod"
	"github.com/stripe/stripe-go/v83/price"
	"github.com/stripe/stripe-go/v83/product"
	"github.com/stripe/stripe-go/v83/setupintent"
	"github.com/stripe/stripe-go/v83/subscription"
)

type Service struct {
	storage         *Storage
	atomicProcessor atomic.AtomicProcessor
	bus             *bus.Bus
	account         *account.Service
}

func NewService(storage *Storage, atomicProcessor atomic.AtomicProcessor, bus *bus.Bus, accountSvc *account.Service) *Service {
	// Note: Stripe is used via package-level helpers using API key configured globally if needed in future.
	s := &Service{
		storage:         storage,
		atomicProcessor: atomicProcessor,
		bus:             bus,
		account:         accountSvc,
	}
	// Best-effort background plan sync on service creation
	go func() {
		// use background context to avoid blocking init
		if err := s.SyncPlansToStripe(context.Background()); err != nil {
			slog.Error("plan sync to stripe failed", "error", err)
		}
	}()
	return s
}

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

// PaymentMethodInfo represents the default card details associated with the customer's subscription
type PaymentMethodInfo struct {
	ID              string `json:"id,omitempty"`
	Brand           string `json:"brand,omitempty"`
	Last4           string `json:"last4,omitempty"`
	ExpMonth        int64  `json:"expMonth,omitempty"`
	ExpYear         int64  `json:"expYear,omitempty"`
	Expired         bool   `json:"expired"`
	ExpiringSoon    bool   `json:"expiringSoon"`
	DaysUntilExpiry int    `json:"daysUntilExpiry"`
}

// SubscriptionDetails aggregates subscription record and payment method details
type SubscriptionDetails struct {
	Subscription  *Subscription     `json:"subscription,omitempty"`
	Plan          *Plan             `json:"plan,omitempty"`
	PaymentMethod PaymentMethodInfo `json:"paymentMethod"`
}

func (s *Service) GetPlanByDescriptor(ctx context.Context, descriptor string) (*Plan, error) {
	return s.storage.plan.FindOne(ctx, s.storage.plan.ScopeEquals(PlanSchema.Descriptor, descriptor))
}

func (s *Service) GetPlanByID(ctx context.Context, id string) (*Plan, error) {
	return s.storage.plan.FindByID(ctx, id)
}

func (s *Service) ListPlans(ctx context.Context) ([]*Plan, error) {
	return s.storage.plan.FindMany(ctx)
}

func (s *Service) GetSubscriptionByID(ctx context.Context, id string) (*Subscription, error) {
	return s.storage.subscription.FindByID(ctx, id)
}

func (s *Service) GetSubscriptionByWorkspaceID(ctx context.Context, workspaceID string) (*Subscription, error) {
	return s.storage.subscription.FindOne(ctx, s.storage.subscription.ScopeWorkspaceID(workspaceID), s.storage.subscription.WithPreload(PlanStruct))
}

// EnsureCustomer makes sure the workspace has a Stripe customer and returns it
func (s *Service) EnsureCustomer(ctx context.Context, ws *account.Workspace) (string, error) {
	if ws.StripeCustomerID.Valid && ws.StripeCustomerID.String != "" {
		return ws.StripeCustomerID.String, nil
	}
	params := &stripelib.CustomerParams{
		Metadata: map[string]string{
			"workspace_id": ws.ID,
		},
	}
	c, err := customer.New(params)
	if err != nil {
		return "", err
	}
	if err := s.account.SetWorkspaceStripeCustomer(ctx, ws.ID, c.ID); err != nil {
		return "", err
	}
	return c.ID, nil
}

// AttachAndSetDefaultPaymentMethod attaches a payment method to the customer and sets it as default
func (s *Service) AttachAndSetDefaultPaymentMethod(ctx context.Context, ws *account.Workspace, pmID string) error {
	custID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		return err
	}
	// attach
	_, err = paymentmethod.Attach(pmID, &stripelib.PaymentMethodAttachParams{Customer: stripelib.String(custID)})
	if err != nil {
		return ErrInvalidPaymentMethod(err)
	}
	// set as default on customer
	_, err = customer.Update(custID, &stripelib.CustomerParams{
		InvoiceSettings: &stripelib.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripelib.String(pmID),
		},
	})
	if err != nil {
		return err
	}
	return s.account.SetWorkspaceDefaultPaymentMethod(ctx, ws.ID, pmID)
}

// CreateSetupIntent returns a client secret to collect and save a payment method for the workspace
func (s *Service) CreateSetupIntent(ctx context.Context, ws *account.Workspace) (string, error) {
	custID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		return "", err
	}
	si, err := setupintent.New(&stripelib.SetupIntentParams{
		Customer:           stripelib.String(custID),
		PaymentMethodTypes: []*string{stripelib.String("card")},
		Usage:              stripelib.String("off_session"),
	})
	if err != nil {
		return "", err
	}
	return si.ClientSecret, nil
}

// CreateOrUpdateSubscription creates a new subscription or updates existing to new plan with proration
func (s *Service) CreateOrUpdateSubscription(ctx context.Context, ws *account.Workspace, plan *Plan) (*Subscription, error) {
	// One active subscription per workspace
	existing, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil && !database.IsRecordNotFound(err) {
		return nil, err
	}
	if existing != nil && existing.PlanID == plan.ID && existing.Status == SubscriptionStatusActive {
		return nil, ErrCannotChangeToSamePlan(nil)
	}

	custID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		return nil, err
	}

	// If we have an existing subscription, update it; else create a new one
	var stripeSub *stripelib.Subscription
	if existing != nil {
		// basic downgrade protection: block moving to a cheaper plan without prior usage checks
		if currentPlan, err := s.GetPlanByID(ctx, existing.PlanID); err == nil {
			if plan.Price.LessThan(currentPlan.Price) && existing.Status == SubscriptionStatusActive {
				// 1) Feature-based compatibility: prevent downgrades that remove features currently available
				if err := s.ensureFeatureCompatibility(currentPlan, plan); err != nil {
					return nil, err
				}
				// Usage-aware checks across modules
				if err := s.ensureWithinNewPlanLimits(ctx, ws.ID, plan); err != nil {
					return nil, err
				}
			}
		}
		// proration on update
		up := &stripelib.SubscriptionParams{
			Items: []*stripelib.SubscriptionItemsParams{
				{Price: stripelib.String(plan.StripePlanID)},
			},
			ProrationBehavior: stripelib.String("create_prorations"),
			CancelAtPeriodEnd: stripelib.Bool(false),
		}
		stripeSub, err = subscription.Update(existing.StripeSubID, up)
		if err != nil {
			return nil, err
		}
		existing.PlanID = plan.ID
		existing.Status = mapStripeStatus(stripeSub.Status)
		existing.CurrentPeriodEnd = time.Now()
		if err := s.storage.subscription.UpdateOne(ctx, existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	// create new subscription (assumes a default payment method is set)
	sp := &stripelib.SubscriptionParams{
		Customer: stripelib.String(custID),
		Items: []*stripelib.SubscriptionItemsParams{
			{Price: stripelib.String(plan.StripePlanID)},
		},
		PaymentBehavior:   stripelib.String("default_incomplete"),
		ProrationBehavior: stripelib.String("create_prorations"),
		CancelAtPeriodEnd: stripelib.Bool(false),
	}
	// For free plan (price = 0), allow creation without requiring a payment method
	if plan.Price.IsZero() {
		sp.PaymentBehavior = stripelib.String("allow_incomplete")
	}
	stripeSub, err = subscription.New(sp)
	if err != nil {
		return nil, err
	}
	newSub := &Subscription{
		WorkspaceID:      ws.ID,
		PlanID:           plan.ID,
		StripeSubID:      stripeSub.ID,
		Status:           mapStripeStatus(stripeSub.Status),
		CurrentPeriodEnd: time.Now(),
	}
	if err := s.storage.subscription.CreateOne(ctx, newSub); err != nil {
		return nil, err
	}
	return newSub, nil
}

// CancelSubscriptionImmediately cancels subscription now (no proration here); prorated refund is handled by webhook
func (s *Service) CancelSubscriptionImmediately(ctx context.Context, ws *account.Workspace) error {
	subRec, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		return ErrSubscriptionNotFound(err, ws.ID)
	}
	// Cancel at Stripe immediately without proration; webhook will issue credit note refund
	_, err = subscription.Cancel(subRec.StripeSubID, &stripelib.SubscriptionCancelParams{InvoiceNow: stripelib.Bool(false), Prorate: stripelib.Bool(false)})
	if err != nil && !errors.Is(err, fmt.Errorf("%w", err)) { // keep best effort cancellation
		slog.Error("Failed to cancel Stripe subscription", "error", err)
	}
	// Update DB
	subRec.Status = SubscriptionStatusCanceled
	subRec.CurrentPeriodEnd = time.Now()
	return s.storage.subscription.UpdateOne(ctx, subRec)
}

func mapStripeStatus(s stripelib.SubscriptionStatus) SubscriptionStatus {
	switch s {
	case stripelib.SubscriptionStatusActive:
		return SubscriptionStatusActive
	case stripelib.SubscriptionStatusPastDue:
		return SubscriptionStatusPastDue
	case stripelib.SubscriptionStatusUnpaid:
		return SubscriptionStatusUnpaid
	case stripelib.SubscriptionStatusIncomplete:
		return SubscriptionStatusIncomplete
	case stripelib.SubscriptionStatusCanceled:
		return SubscriptionStatusCanceled
	default:
		return SubscriptionStatusIncomplete
	}
}

// ensureFeatureCompatibility prevents downgrades that remove features used by the current plan.
// Conservative rule: if a feature is enabled on the current plan but disabled on the new plan,
// block the downgrade. This can be relaxed later with usage-aware feature checks per module.
func (s *Service) ensureFeatureCompatibility(currentPlan, newPlan *Plan) error {
	// Compare all boolean feature flags
	cur := currentPlan.Features
	nxt := newPlan.Features
	// If current has a feature and next does not, disallow downgrade
	if cur.CustomerManagement && !nxt.CustomerManagement {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.InventoryManagement && !nxt.InventoryManagement {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.OrderManagement && !nxt.OrderManagement {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.ExpenseManagement && !nxt.ExpenseManagement {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.AssetsManagement && !nxt.AssetsManagement {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.Accounting && !nxt.Accounting {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.BasicAnalytics && !nxt.BasicAnalytics {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.FinancialReports && !nxt.FinancialReports {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.DataImport && !nxt.DataImport {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.DataExport && !nxt.DataExport {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.AdvancedAnalytics && !nxt.AdvancedAnalytics {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.AdvancedFinancialReports && !nxt.AdvancedFinancialReports {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.OrderPaymentLinks && !nxt.OrderPaymentLinks {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.InvoiceGeneration && !nxt.InvoiceGeneration {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.ExportAnalyticsData && !nxt.ExportAnalyticsData {
		return ErrCannotDowngradePlan(nil)
	}
	if cur.AIBusinessAssistant && !nxt.AIBusinessAssistant {
		return ErrCannotDowngradePlan(nil)
	}
	return nil
}

// SyncSubscriptionStatus updates the local record based on Stripe status
func (s *Service) SyncSubscriptionStatus(ctx context.Context, stripeSubID string, status string, periodStart, periodEnd int64) {
	rec, err := s.storage.subscription.FindOne(ctx, s.storage.subscription.ScopeEquals(SubscriptionSchema.StripeSubID, stripeSubID))
	if err != nil || rec == nil {
		return
	}
	rec.Status = mapStripeStatus(stripelib.SubscriptionStatus(status))
	if periodEnd > 0 {
		rec.CurrentPeriodEnd = time.Unix(periodEnd, 0)
	}
	_ = s.storage.subscription.UpdateOne(ctx, rec)
}

// MarkSubscriptionPastDue sets subscription status to past_due
func (s *Service) MarkSubscriptionPastDue(ctx context.Context, stripeSubID string) {
	rec, err := s.storage.subscription.FindOne(ctx, s.storage.subscription.ScopeEquals(SubscriptionSchema.StripeSubID, stripeSubID))
	if err != nil || rec == nil {
		return
	}
	rec.Status = SubscriptionStatusPastDue
	_ = s.storage.subscription.UpdateOne(ctx, rec)
}

// MarkSubscriptionActive sets subscription status to active
func (s *Service) MarkSubscriptionActive(ctx context.Context, stripeSubID string) {
	rec, err := s.storage.subscription.FindOne(ctx, s.storage.subscription.ScopeEquals(SubscriptionSchema.StripeSubID, stripeSubID))
	if err != nil || rec == nil {
		return
	}
	rec.Status = SubscriptionStatusActive
	_ = s.storage.subscription.UpdateOne(ctx, rec)
}

// RefundAndFinalizeCancellation computes prorated refund and cancels in Stripe, then updates local DB
func (s *Service) RefundAndFinalizeCancellation(ctx context.Context, stripeSubID string, periodStart, periodEnd int64) {
	rec, err := s.storage.subscription.FindOne(ctx, s.storage.subscription.ScopeEquals(SubscriptionSchema.StripeSubID, stripeSubID))
	if err != nil || rec == nil {
		return
	}
	// If period bounds are not provided by the event, skip refund calculation to avoid incorrect credits
	// (We still cancel immediately; future improvements can derive period from invoice line items.)
	// Get latest paid invoice for this subscription
	ip := &stripelib.InvoiceListParams{Subscription: stripelib.String(stripeSubID)}
	ip.Status = stripelib.String(string(stripelib.InvoiceStatusPaid))
	ip.Limit = stripelib.Int64(1)
	iter := invoice.List(ip)
	if iter.Next() {
		inv := iter.Invoice()
		// Only compute prorated refund if both periodStart and periodEnd are provided
		if periodStart > 0 && periodEnd > 0 {
			// Compute prorated refund based on remaining time
			now := time.Now()
			pEnd := time.Unix(periodEnd, 0)
			pStart := time.Unix(periodStart, 0)
			if pEnd.Before(now) {
				pEnd = rec.CurrentPeriodEnd
			}
			if !(pStart.After(now) || pEnd.Before(pStart)) {
				total := pEnd.Sub(pStart).Seconds()
				remaining := pEnd.Sub(now).Seconds()
				if total > 0 && remaining > 0 && inv.AmountPaid > 0 {
					refundAmount := int64(float64(inv.AmountPaid) * (remaining / total))
					if refundAmount > 0 {
						// Idempotency: skip if a cancel_prorated credit note already exists for this invoice
						lparams := &stripelib.CreditNoteListParams{Invoice: stripelib.String(inv.ID)}
						lparams.Limit = stripelib.Int64(10)
						alreadyRefunded := false
						itCN := creditnote.List(lparams)
						for itCN.Next() {
							cn := itCN.CreditNote()
							if cn != nil && cn.Metadata != nil {
								if v, ok := cn.Metadata["kyoraRefundKind"]; ok && v == "cancel_prorated" {
									alreadyRefunded = true
									break
								}
							}
						}
						if !alreadyRefunded {
							// Issue a credit note on the invoice with a refund for the prorated amount
							_, rerr := creditnote.New(&stripelib.CreditNoteParams{
								Invoice:      stripelib.String(inv.ID),
								Amount:       stripelib.Int64(refundAmount),
								RefundAmount: stripelib.Int64(refundAmount),
								Reason:       stripelib.String(string(stripelib.CreditNoteReasonOrderChange)),
								Memo:         stripelib.String("Prorated refund for immediate cancellation"),
								Metadata: map[string]string{
									"kyoraRefundKind": "cancel_prorated",
									"subscription":    stripeSubID,
									"workspaceId":     rec.WorkspaceID,
								},
							})
							if rerr != nil {
								slog.Error("failed to create credit note refund", "error", rerr, "invoice", inv.ID)
							}
						} else {
							slog.Info("skip credit note refund: already applied", "invoice", inv.ID, "subscription", stripeSubID)
						}
					}
				}
			}
		} else {
			slog.Info("Skipping refund calculation: missing period bounds", "subscription", stripeSubID)
		}
	}
	// Cancel at Stripe immediately without additional proration (credit note already applied)
	_, _ = subscription.Cancel(stripeSubID, &stripelib.SubscriptionCancelParams{InvoiceNow: stripelib.Bool(false), Prorate: stripelib.Bool(false)})
	// Update local record
	rec.Status = SubscriptionStatusCanceled
	rec.CurrentPeriodEnd = time.Now()
	_ = s.storage.subscription.UpdateOne(ctx, rec)
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
		return "", ErrSubscriptionNotFound(fmt.Errorf("invoice not owned by workspace"), ws.ID)
	}
	if inv.InvoicePDF == "" && inv.HostedInvoiceURL == "" {
		return "", ErrSubscriptionNotFound(fmt.Errorf("invoice has no downloadable link"), invoiceID)
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
		return ErrSubscriptionNotFound(fmt.Errorf("invoice not owned by workspace"), ws.ID)
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

// GetSubscriptionDetails returns current subscription plus default payment method details
func (s *Service) GetSubscriptionDetails(ctx context.Context, ws *account.Workspace) (*SubscriptionDetails, error) {
	rec, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil && !database.IsRecordNotFound(err) {
		return nil, err
	}
	var plan *Plan
	if rec != nil {
		plan, _ = s.GetPlanByID(ctx, rec.PlanID)
	}
	// Pull customer default payment method
	pmInfo := PaymentMethodInfo{}
	custID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		return &SubscriptionDetails{Subscription: rec, Plan: plan, PaymentMethod: pmInfo}, nil
	}
	c, err := customer.Get(custID, nil)
	if err != nil || c == nil || c.InvoiceSettings == nil || c.InvoiceSettings.DefaultPaymentMethod == nil {
		return &SubscriptionDetails{Subscription: rec, Plan: plan, PaymentMethod: pmInfo}, nil
	}
	dpm := c.InvoiceSettings.DefaultPaymentMethod
	pm, err := paymentmethod.Get(dpm.ID, nil)
	if err != nil || pm == nil || pm.Card == nil {
		return &SubscriptionDetails{Subscription: rec, Plan: plan, PaymentMethod: pmInfo}, nil
	}
	pmInfo.ID = pm.ID
	pmInfo.Brand = string(pm.Card.Brand)
	pmInfo.Last4 = pm.Card.Last4
	pmInfo.ExpMonth = int64(pm.Card.ExpMonth)
	pmInfo.ExpYear = int64(pm.Card.ExpYear)
	// compute expiry status
	now := time.Now()
	expTime := time.Date(int(pm.Card.ExpYear), time.Month(pm.Card.ExpMonth), 1, 0, 0, 0, 0, now.Location()).AddDate(0, 1, -1) // end of exp month
	days := int(expTime.Sub(now).Hours() / 24)
	pmInfo.DaysUntilExpiry = days
	pmInfo.Expired = days < 0
	pmInfo.ExpiringSoon = !pmInfo.Expired && days <= 30
	return &SubscriptionDetails{Subscription: rec, Plan: plan, PaymentMethod: pmInfo}, nil
}

// SyncPlansToStripe ensures all local plans exist in Stripe as products/prices and updates price ids when changed
func (s *Service) SyncPlansToStripe(ctx context.Context) error {
	plans, err := s.storage.plan.FindMany(ctx)
	if err != nil {
		return err
	}
	for _, p := range plans {
		// ensure product by descriptor
		var prod *stripelib.Product
		// Try to find by metadata is not directly supported; so use existing price to backtrack product when available
		if p.StripePlanID != "" {
			if pr, err := price.Get(p.StripePlanID, nil); err == nil && pr != nil {
				if pr.Product != nil {
					prod, _ = product.Get(pr.Product.ID, nil)
				}
			}
		}
		// If product nil, create a new one
		if prod == nil {
			prms := &stripelib.ProductParams{Name: stripelib.String(p.Name), Description: stripelib.String(p.Description)}
			prms.Metadata = map[string]string{"kyora_plan_id": p.ID, "descriptor": p.Descriptor}
			prod, err = product.New(prms)
			if err != nil {
				slog.Error("stripe product new failed", "err", err, "plan", p.ID)
				continue
			}
		} else {
			// Update name/description if changed
			_, _ = product.Update(prod.ID, &stripelib.ProductParams{Name: stripelib.String(p.Name), Description: stripelib.String(p.Description)})
		}
		// Ensure price matches
		var existing *stripelib.Price
		if p.StripePlanID != "" {
			if pr, err := price.Get(p.StripePlanID, nil); err == nil {
				existing = pr
			}
		}
		needNewPrice := true
		interval := "month"
		if p.BillingCycle == BillingCycleYearly {
			interval = "year"
		}
		unit := p.Price.Mul(decimal.NewFromInt(100)).IntPart()
		if existing != nil {
			if string(existing.Currency) == p.Currency && existing.Recurring != nil && string(existing.Recurring.Interval) == interval && existing.UnitAmount == unit {
				needNewPrice = false
			}
		}
		if needNewPrice {
			prms := &stripelib.PriceParams{
				Currency:   stripelib.String(p.Currency),
				UnitAmount: stripelib.Int64(unit),
				Recurring:  &stripelib.PriceRecurringParams{Interval: stripelib.String(interval)},
				Product:    stripelib.String(prod.ID),
			}
			// zero-amount price (free plan) is allowed
			newPrice, err := price.New(prms)
			if err != nil {
				slog.Error("stripe price new failed", "err", err, "plan", p.ID)
				continue
			}
			p.StripePlanID = newPrice.ID
			if uerr := s.storage.plan.UpdateOne(ctx, p); uerr != nil {
				slog.Error("update plan stripe id failed", "err", uerr, "plan", p.ID)
			}
		}
	}
	return nil
}

// ensureWithinNewPlanLimits enforces usage within new plan limits (users, businesses, monthly orders)
func (s *Service) ensureWithinNewPlanLimits(ctx context.Context, workspaceID string, newPlan *Plan) error {
	// Users
	users, err := s.account.CountWorkspaceUsers(ctx, workspaceID)
	if err != nil {
		return err
	}
	if users > newPlan.Limits.MaxTeamMembers {
		return ErrCannotDowngradePlan(nil)
	}
	// Businesses
	businesses, err := s.storage.CountBusinessesByWorkspace(ctx, workspaceID)
	if err != nil {
		return err
	}
	if businesses > newPlan.Limits.MaxBusinesses {
		return ErrCannotDowngradePlan(nil)
	}
	// Monthly orders
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthEnd := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location()).Add(-time.Nanosecond)
	orders, err := s.storage.CountMonthlyOrdersByWorkspace(ctx, workspaceID, monthStart, monthEnd)
	if err != nil {
		return err
	}
	if orders > newPlan.Limits.MaxOrdersPerMonth {
		return ErrCannotDowngradePlan(nil)
	}
	return nil
}

// ResumeSubscriptionIfNoDue attempts to pay open invoices then recreates a subscription with the same plan
func (s *Service) ResumeSubscriptionIfNoDue(ctx context.Context, ws *account.Workspace) (*Subscription, error) {
	// fetch current (canceled) record to get plan
	rec, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		return nil, ErrSubscriptionNotFound(err, ws.ID)
	}
	if rec.Status != SubscriptionStatusCanceled {
		return rec, nil
	}
	// attempt to pay open invoices
	custID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		return nil, err
	}
	ip := &stripelib.InvoiceListParams{Customer: stripelib.String(custID)}
	ip.Status = stripelib.String(string(stripelib.InvoiceStatusOpen))
	ip.Limit = stripelib.Int64(10)
	it := invoice.List(ip)
	for it.Next() {
		inv := it.Invoice()
		if _, err := invoice.Pay(inv.ID, &stripelib.InvoicePayParams{}); err != nil {
			return nil, ErrSubscriptionNotActive(err)
		}
	}
	// all due cleared — re-subscribe to same plan
	plan, err := s.GetPlanByID(ctx, rec.PlanID)
	if err != nil {
		return nil, err
	}
	return s.CreateOrUpdateSubscription(ctx, ws, plan)
}
