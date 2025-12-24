package billing

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
	atomic "github.com/abdelrahman146/kyora/internal/platform/types/atomic"
	"github.com/abdelrahman146/kyora/internal/platform/utils/helpers"
	"github.com/shopspring/decimal"
	stripelib "github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/creditnote"
	"github.com/stripe/stripe-go/v83/invoice"
	"github.com/stripe/stripe-go/v83/paymentmethod"
	"github.com/stripe/stripe-go/v83/subscription"
	"github.com/stripe/stripe-go/v83/subscriptionschedule"
)

func (s *Service) GetPlanByDescriptor(ctx context.Context, descriptor string) (*Plan, error) {
	return s.storage.plan.FindOne(ctx, s.storage.plan.ScopeEquals(PlanSchema.Descriptor, descriptor))
}

func (s *Service) GetPlanByID(ctx context.Context, id string) (*Plan, error) {
	return s.storage.plan.FindByID(ctx, id)
}

func (s *Service) ListPlans(ctx context.Context) ([]*Plan, error) {
	return s.storage.plan.FindMany(ctx)
}

func (s *Service) GetSubscriptionByWorkspaceID(ctx context.Context, workspaceID string) (*Subscription, error) {
	return s.storage.subscription.FindOne(ctx, s.storage.subscription.ScopeWorkspaceID(workspaceID), s.storage.subscription.WithPreload(PlanStruct))
}

// CreateOrUpdateSubscription creates a new subscription or updates existing to new plan with proration
func (s *Service) CreateOrUpdateSubscription(ctx context.Context, ws *account.Workspace, plan *Plan) (*Subscription, error) {
	if err := s.ensurePlanSynced(ctx, plan); err != nil {
		logger.FromContext(ctx).Error("Failed to ensure plan synced before subscription", "error", err, "plan_id", plan.ID)
		return nil, fmt.Errorf("failed to ensure plan in stripe: %w", err)
	}
	if plan == nil || plan.StripePlanID == nil || *plan.StripePlanID == "" {
		return nil, fmt.Errorf("plan missing stripe price id")
	}
	if ws == nil {
		return nil, fmt.Errorf("workspace cannot be nil")
	}
	if plan == nil {
		return nil, fmt.Errorf("plan cannot be nil")
	}
	existing, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		if database.IsRecordNotFound(err) {
			existing = nil
		} else {
			return nil, fmt.Errorf("failed to check existing subscription: %w", err)
		}
	}
	if existing != nil && existing.PlanID == plan.ID && existing.Status == SubscriptionStatusActive {
		return existing, nil
	}
	custID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure customer: %w", err)
	}
	var result *Subscription
	err = s.atomicProcessor.Exec(ctx, func(ctx context.Context) error {
		var stripeSub *stripelib.Subscription
		if existing != nil {
			if currentPlan, err := s.GetPlanByID(ctx, existing.PlanID); err == nil {
				if plan.Price.LessThan(currentPlan.Price) && existing.Status == SubscriptionStatusActive {
					if err := s.ensureFeatureCompatibility(currentPlan, plan); err != nil {
						return err
					}
					if err := s.ensureWithinNewPlanLimits(ctx, ws.ID, plan); err != nil {
						return err
					}
				}
			}
			idempotencyKey := fmt.Sprintf("sub_update_%s_%s", existing.StripeSubID, plan.ID)
			updateParams := &stripelib.SubscriptionParams{
				Items:             []*stripelib.SubscriptionItemsParams{{Price: stripelib.String(*plan.StripePlanID)}},
				ProrationBehavior: stripelib.String("create_prorations"),
				CancelAtPeriodEnd: stripelib.Bool(false),
				Metadata:          map[string]string{"workspace_id": ws.ID, "plan_id": plan.ID},
			}
			updateParams.SetIdempotencyKey(idempotencyKey)
			stripeSub, err = subscription.Update(existing.StripeSubID, updateParams)
			if err != nil {
				logger.FromContext(ctx).Error("Failed to update Stripe subscription", "error", err, "subscriptionId", existing.StripeSubID, "planId", plan.ID)
				return fmt.Errorf("failed to update subscription: %w", err)
			}
			existing.PlanID = plan.ID
			existing.Status = mapStripeStatus(stripeSub.Status)
			if err := s.storage.subscription.UpdateOne(ctx, existing); err != nil {
				return fmt.Errorf("failed to update local subscription: %w", err)
			}
			result = existing
			logger.FromContext(ctx).Info("Updated subscription", "workspaceId", ws.ID, "subscriptionId", existing.StripeSubID, "newPlanId", plan.ID)
			return nil
		}
		idempotencyKey := fmt.Sprintf("sub_create_%s_%s", ws.ID, plan.ID)
		createParams := &stripelib.SubscriptionParams{
			Customer:          stripelib.String(custID),
			Items:             []*stripelib.SubscriptionItemsParams{{Price: stripelib.String(*plan.StripePlanID)}},
			ProrationBehavior: stripelib.String("create_prorations"),
			CancelAtPeriodEnd: stripelib.Bool(false),
			Metadata:          map[string]string{"workspace_id": ws.ID, "plan_id": plan.ID},
		}
		if plan.Price.IsZero() {
			createParams.PaymentBehavior = stripelib.String("allow_incomplete")
		} else {
			createParams.PaymentBehavior = stripelib.String("default_incomplete")
			createParams.CollectionMethod = stripelib.String("charge_automatically")
		}
		createParams.SetIdempotencyKey(idempotencyKey)
		stripeSub, err = withStripeRetry(ctx, 3, func() (*stripelib.Subscription, error) { return subscription.New(createParams) })
		if err != nil {
			logger.FromContext(ctx).Error("Failed to create Stripe subscription", "error", err, "customerId", custID, "planId", plan.ID)
			return fmt.Errorf("failed to create subscription: %w", err)
		}
		newSub := &Subscription{
			WorkspaceID:      ws.ID,
			PlanID:           plan.ID,
			StripeSubID:      stripeSub.ID,
			Status:           mapStripeStatus(stripeSub.Status),
			CurrentPeriodEnd: time.Now(),
		}
		if err := s.storage.subscription.CreateOne(ctx, newSub); err != nil {
			return fmt.Errorf("failed to create local subscription: %w", err)
		}
		result = newSub
		logger.FromContext(ctx).Info("Created new subscription", "workspaceId", ws.ID, "subscriptionId", stripeSub.ID, "planId", plan.ID)
		return nil
	}, atomic.WithRetries(2))
	if err != nil {
		return nil, err
	}
	if result != nil && existing == nil && s.Notification != nil {
		paymentMethodLastFour := ""
		if details, err := s.GetSubscriptionDetails(ctx, ws); err == nil && details.PaymentMethod.Last4 != "" {
			paymentMethodLastFour = details.PaymentMethod.Last4
		}
		l := logger.FromContext(ctx)
		go func(workspaceID string, sub *Subscription, p *Plan, last4 string) {
			bg, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := s.Notification.SendSubscriptionWelcomeEmail(bg, workspaceID, sub, p, last4); err != nil {
				l.Warn("failed to send subscription welcome email", "error", err, "workspaceId", workspaceID, "subscriptionId", sub.ID)
			}
		}(ws.ID, result, plan, paymentMethodLastFour)
	}
	return result, nil
}

// CancelSubscriptionImmediately cancels subscription now
func (s *Service) CancelSubscriptionImmediately(ctx context.Context, ws *account.Workspace) error {
	if ws == nil {
		return fmt.Errorf("workspace cannot be nil")
	}
	subRec, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		return ErrSubscriptionNotFound(err, ws.ID)
	}
	if subRec.Status == SubscriptionStatusCanceled {
		return nil
	}
	return s.atomicProcessor.Exec(ctx, func(ctx context.Context) error {
		idempotencyKey := fmt.Sprintf("cancel_%s", subRec.StripeSubID)
		cancelParams := &stripelib.SubscriptionCancelParams{InvoiceNow: stripelib.Bool(false), Prorate: stripelib.Bool(false)}
		cancelParams.SetIdempotencyKey(idempotencyKey)
		_, err = withStripeRetry(ctx, 3, func() (*stripelib.Subscription, error) { return subscription.Cancel(subRec.StripeSubID, cancelParams) })
		if err != nil {
			logger.FromContext(ctx).Error("Failed to cancel Stripe subscription", "error", err, "subscriptionId", subRec.StripeSubID, "workspaceId", ws.ID)
		}
		subRec.Status = SubscriptionStatusCanceled
		subRec.CurrentPeriodEnd = time.Now()
		if updateErr := s.storage.subscription.UpdateOne(ctx, subRec); updateErr != nil {
			logger.FromContext(ctx).Error("Failed to update local subscription status", "error", updateErr, "subscriptionId", subRec.StripeSubID)
			return fmt.Errorf("failed to update local subscription: %w", updateErr)
		}
		if err != nil {
			return fmt.Errorf("failed to cancel Stripe subscription: %w", err)
		}
		logger.FromContext(ctx).Info("Successfully canceled subscription", "workspaceId", ws.ID, "subscriptionId", subRec.StripeSubID)
		if s.Notification != nil {
			if plan, planErr := s.GetPlanByID(ctx, subRec.PlanID); planErr == nil {
				l := logger.FromContext(ctx)
				go func(workspaceID string, sub *Subscription, p *Plan) {
					bg, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()
					if err := s.Notification.SendSubscriptionCanceledEmail(bg, workspaceID, sub, p, time.Now().UTC(), ""); err != nil {
						l.Warn("failed to send subscription canceled email", "error", err, "workspaceId", workspaceID, "subscriptionId", sub.StripeSubID)
					}
				}(ws.ID, subRec, plan)
			}
		}
		return nil
	}, atomic.WithRetries(2))
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
func (s *Service) ensureFeatureCompatibility(currentPlan, newPlan *Plan) error {
	cur := currentPlan.Features
	nxt := newPlan.Features
	curMap := map[string]bool{
		"customerManagement":       cur.CustomerManagement,
		"inventoryManagement":      cur.InventoryManagement,
		"orderManagement":          cur.OrderManagement,
		"expenseManagement":        cur.ExpenseManagement,
		"accounting":               cur.Accounting,
		"basicAnalytics":           cur.BasicAnalytics,
		"financialReports":         cur.FinancialReports,
		"dataImport":               cur.DataImport,
		"dataExport":               cur.DataExport,
		"advancedAnalytics":        cur.AdvancedAnalytics,
		"advancedFinancialReports": cur.AdvancedFinancialReports,
		"orderPaymentLinks":        cur.OrderPaymentLinks,
		"invoiceGeneration":        cur.InvoiceGeneration,
		"exportAnalyticsData":      cur.ExportAnalyticsData,
		"aiBusinessAssistant":      cur.AIBusinessAssistant,
	}
	nxtMap := map[string]bool{
		"customerManagement":       nxt.CustomerManagement,
		"inventoryManagement":      nxt.InventoryManagement,
		"orderManagement":          nxt.OrderManagement,
		"expenseManagement":        nxt.ExpenseManagement,
		"accounting":               nxt.Accounting,
		"basicAnalytics":           nxt.BasicAnalytics,
		"financialReports":         nxt.FinancialReports,
		"dataImport":               nxt.DataImport,
		"dataExport":               nxt.DataExport,
		"advancedAnalytics":        nxt.AdvancedAnalytics,
		"advancedFinancialReports": nxt.AdvancedFinancialReports,
		"orderPaymentLinks":        nxt.OrderPaymentLinks,
		"invoiceGeneration":        nxt.InvoiceGeneration,
		"exportAnalyticsData":      nxt.ExportAnalyticsData,
		"aiBusinessAssistant":      nxt.AIBusinessAssistant,
	}
	for k, v := range curMap {
		if v && !nxtMap[k] {
			return ErrCannotDowngradePlan(nil)
		}
	}
	return nil
}

// SyncSubscriptionStatus updates the local record based on Stripe status
func (s *Service) SyncSubscriptionStatus(ctx context.Context, stripeSubID string, status string, periodStart, periodEnd int64) error {
	rec, err := s.storage.subscription.FindOne(ctx, s.storage.subscription.ScopeEquals(SubscriptionSchema.StripeSubID, stripeSubID))
	if err != nil {
		return err
	}
	if rec == nil {
		logger.FromContext(ctx).Warn("subscription not found for stripeSubID", "stripeSubId", stripeSubID)
		return nil
	}
	rec.Status = mapStripeStatus(stripelib.SubscriptionStatus(status))
	if periodEnd > 0 {
		rec.CurrentPeriodEnd = time.Unix(periodEnd, 0)
	}
	if err := s.storage.subscription.UpdateOne(ctx, rec); err != nil {
		return err
	}
	return nil
}

// MarkSubscriptionPastDue sets subscription status to past_due
func (s *Service) MarkSubscriptionPastDue(ctx context.Context, stripeSubID string) error {
	rec, err := s.storage.subscription.FindOne(ctx, s.storage.subscription.ScopeEquals(SubscriptionSchema.StripeSubID, stripeSubID))
	if err != nil {
		return err
	}
	if rec == nil {
		logger.FromContext(ctx).Warn("subscription not found for stripeSubID", "stripeSubId", stripeSubID)
		return nil
	}
	rec.Status = SubscriptionStatusPastDue
	return s.storage.subscription.UpdateOne(ctx, rec)
}

// MarkSubscriptionActive sets subscription status to active
func (s *Service) MarkSubscriptionActive(ctx context.Context, stripeSubID string) error {
	rec, err := s.storage.subscription.FindOne(ctx, s.storage.subscription.ScopeEquals(SubscriptionSchema.StripeSubID, stripeSubID))
	if err != nil {
		return err
	}
	if rec == nil {
		logger.FromContext(ctx).Warn("subscription not found for stripeSubID", "stripeSubId", stripeSubID)
		return nil
	}
	rec.Status = SubscriptionStatusActive
	return s.storage.subscription.UpdateOne(ctx, rec)
}

// RefundAndFinalizeCancellation computes prorated refund and cancels in Stripe, then updates local DB
func (s *Service) RefundAndFinalizeCancellation(ctx context.Context, stripeSubID string, periodStart, periodEnd int64) error {
	rec, err := s.storage.subscription.FindOne(ctx, s.storage.subscription.ScopeEquals(SubscriptionSchema.StripeSubID, stripeSubID))
	if err != nil {
		return err
	}
	if rec == nil {
		logger.FromContext(ctx).Warn("subscription not found for stripeSubID", "stripeSubId", stripeSubID)
		return nil
	}
	if periodStart > 0 && periodEnd > 0 {
		ip := &stripelib.InvoiceListParams{Subscription: stripelib.String(stripeSubID)}
		ip.Status = stripelib.String(string(stripelib.InvoiceStatusPaid))
		ip.Limit = stripelib.Int64(1)
		iter := invoice.List(ip)
		if iter.Next() {
			inv := iter.Invoice()
			if inv != nil && inv.AmountPaid > 0 {
				now := time.Now().UTC()
				pStart := time.Unix(periodStart, 0).UTC()
				pEnd := time.Unix(periodEnd, 0).UTC()
				if pEnd.Before(now) {
					pEnd = rec.CurrentPeriodEnd
				}
				if !pStart.After(now) && !pEnd.Before(pStart) {
					totalNanos := pEnd.Sub(pStart).Nanoseconds()
					remainingNanos := pEnd.Sub(now).Nanoseconds()
					if totalNanos > 0 && remainingNanos > 0 {
						refundAmount := decimal.NewFromInt(inv.AmountPaid).
							Mul(decimal.NewFromInt(remainingNanos)).
							Div(decimal.NewFromInt(totalNanos)).
							IntPart()
						if refundAmount > 0 {
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
							if err := itCN.Err(); err != nil {
								logger.FromContext(ctx).Warn("credit note listing failed", "error", err, "invoice", inv.ID)
							}
							if !alreadyRefunded {
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
									logger.FromContext(ctx).Error("failed to create credit note refund", "error", rerr, "invoice", inv.ID)
								}
							} else {
								logger.FromContext(ctx).Info("skip credit note refund: already applied", "invoice", inv.ID, "subscription", stripeSubID)
							}
						}
					}
				}
			}
		}
		if err := iter.Err(); err != nil {
			logger.FromContext(ctx).Warn("invoice listing failed", "error", err, "subscription", stripeSubID)
		}
	} else {
		logger.FromContext(ctx).Info("Skipping refund calculation: missing period bounds", "subscription", stripeSubID)
	}
	cancelParams := &stripelib.SubscriptionCancelParams{InvoiceNow: stripelib.Bool(false), Prorate: stripelib.Bool(false)}
	if _, err := withStripeRetry(ctx, 3, func() (*stripelib.Subscription, error) { return subscription.Cancel(stripeSubID, cancelParams) }); err != nil {
		return err
	}
	rec.Status = SubscriptionStatusCanceled
	rec.CurrentPeriodEnd = time.Now().UTC()
	if err := s.storage.subscription.UpdateOne(ctx, rec); err != nil {
		return err
	}
	return nil
}

// ensureWithinNewPlanLimits enforces usage within new plan limits (users, businesses, monthly orders)
func (s *Service) ensureWithinNewPlanLimits(ctx context.Context, workspaceID string, newPlan *Plan) error {
	users, err := s.account.CountWorkspaceUsers(ctx, workspaceID)
	if err != nil {
		return err
	}
	if users > newPlan.Limits.MaxTeamMembers {
		return ErrCannotDowngradePlan(nil)
	}
	businesses, err := s.storage.CountBusinessesByWorkspace(ctx, workspaceID)
	if err != nil {
		return err
	}
	if businesses > newPlan.Limits.MaxBusinesses {
		return ErrCannotDowngradePlan(nil)
	}
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

// ResumeSubscriptionIfNoDue attempts to pay open invoices then recreates a subscription
// ResumeSubscriptionIfNoDue attempts to restore a canceled/past_due/unpaid/incomplete subscription by settling invoices
// and recreating or updating the subscription. It performs robust checks and returns a fully active local record if possible.
func (s *Service) ResumeSubscriptionIfNoDue(ctx context.Context, ws *account.Workspace) (*Subscription, error) {
	l := logger.FromContext(ctx).With("workspaceId", ws.ID)
	l.Info("attempting subscription resume")
	rec, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		return nil, ErrSubscriptionNotFound(err, ws.ID)
	}
	if rec.Status == SubscriptionStatusActive {
		l.Info("subscription already active")
		return rec, nil
	}
	custID, err := s.EnsureCustomer(ctx, ws)
	if err != nil {
		l.Error("ensure customer failed", "error", err)
		return nil, err
	}
	// Fetch Stripe subscription (may be deleted)
	stripeSub, subErr := subscription.Get(rec.StripeSubID, nil)
	if subErr != nil {
		// If not found (deleted), create new later
		l.Warn("stripe subscription fetch failed; will recreate", "error", subErr)
		stripeSub = nil
	}
	// Attempt to pay any open invoices
	payParams := &stripelib.InvoiceListParams{Customer: stripelib.String(custID)}
	payParams.Status = stripelib.String(string(stripelib.InvoiceStatusOpen))
	payParams.Limit = stripelib.Int64(25)
	it := invoice.List(payParams)
	var payFailures []error
	for it.Next() {
		inv := it.Invoice()
		if inv.Status == stripelib.InvoiceStatusOpen || inv.Status == stripelib.InvoiceStatusDraft {
			if inv.Status == stripelib.InvoiceStatusDraft {
				if _, err := invoice.FinalizeInvoice(inv.ID, nil); err != nil {
					l.Error("failed to finalize invoice", "invoice_id", inv.ID, "error", err)
					payFailures = append(payFailures, err)
					continue
				}
			}
			if _, err := invoice.Pay(inv.ID, &stripelib.InvoicePayParams{}); err != nil {
				l.Error("failed to pay invoice", "invoice_id", inv.ID, "error", err)
				payFailures = append(payFailures, err)
			} else {
				l.Info("invoice paid", "invoice_id", inv.ID)
			}
		}
	}
	if err := it.Err(); err != nil {
		l.Error("invoice listing error during resume", "error", err)
	}
	// If we still have failures, surface a composite error
	if len(payFailures) > 0 {
		return nil, ErrSubscriptionNotActive(errors.Join(payFailures...))
	}
	// Ensure plan exists
	plan, err := s.GetPlanByID(ctx, rec.PlanID)
	if err != nil {
		return nil, err
	}
	if stripeSub == nil || stripeSub.Status == stripelib.SubscriptionStatusCanceled {
		l.Info("recreating subscription in Stripe")
		newRec, createErr := s.CreateOrUpdateSubscription(ctx, ws, plan)
		if createErr != nil {
			return nil, createErr
		}
		l.Info("subscription recreated", "subscriptionId", newRec.StripeSubID)
		return newRec, nil
	}
	// If subscription exists but not active try to update payment behavior
	if stripeSub.Status != stripelib.SubscriptionStatusActive {
		l.Info("updating existing stripe subscription for resume", "stripeStatus", stripeSub.Status)
		updParams := &stripelib.SubscriptionParams{CancelAtPeriodEnd: stripelib.Bool(false)}
		updParams.PaymentBehavior = stripelib.String("default_incomplete")
		updParams.ProrationBehavior = stripelib.String("none")
		if _, err := withStripeRetry(ctx, 3, func() (*stripelib.Subscription, error) { return subscription.Update(stripeSub.ID, updParams) }); err != nil {
			l.Error("subscription update failed", "error", err)
			return nil, ErrSubscriptionNotActive(err)
		}
		// Attempt to refresh default payment method if automatically updated
		if stripeSub.DefaultPaymentMethod != nil {
			if pm, err := paymentmethod.Get(stripeSub.DefaultPaymentMethod.ID, nil); err == nil {
				l.Info("verified default payment method", "paymentMethodId", pm.ID)
			}
		}
	}
	// Sync local status (fetch again)
	stripeSub2, err := subscription.Get(rec.StripeSubID, nil)
	if err == nil {
		if err := s.SyncSubscriptionStatus(ctx, stripeSub2.ID, string(stripeSub2.Status), 0, 0); err != nil {
			l.Warn("failed to sync local subscription status", "error", err)
		} else {
			rec, _ = s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
			l.Info("subscription resume completed", "finalStatus", rec.Status)
		}
	}
	return rec, nil
}

// ScheduleSubscriptionChange schedules a subscription change for a future date
func (s *Service) ScheduleSubscriptionChange(ctx context.Context, ws *account.Workspace, plan *Plan, effectiveDate, prorationMode string) (*stripelib.SubscriptionSchedule, error) {
	l := logger.FromContext(ctx).With("workspaceId", ws.ID, "planDescriptor", plan.Descriptor, "effectiveDate", effectiveDate)
	l.Info("scheduling subscription change")
	sub, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		l.Error("failed to get current subscription", "error", err)
		return nil, err
	}
	if prorationMode != "" {
		switch prorationMode {
		case "none", "create_prorations", "always_invoice":
			// ok
		default:
			return nil, ErrInvalidProrationMode(nil, prorationMode)
		}
	}
	effectiveTime, err := time.Parse("2006-01-02T15:04:05Z", effectiveDate)
	if err != nil {
		if effectiveTime, err = time.Parse("2006-01-02", effectiveDate); err != nil {
			l.Error("invalid effective date format", "error", err)
			return nil, ErrInvalidEffectiveDate(err, effectiveDate)
		}
	}
	if err := s.ensurePlanSynced(ctx, plan); err != nil {
		l.Error("failed to ensure plan synced", "error", err)
		return nil, err
	}
	if plan.StripePlanID == nil || *plan.StripePlanID == "" {
		return nil, fmt.Errorf("plan missing stripe price id")
	}
	scheduleParams := &stripelib.SubscriptionScheduleParams{FromSubscription: stripelib.String(sub.StripeSubID)}
	currentPhase := &stripelib.SubscriptionSchedulePhaseParams{
		Items:     []*stripelib.SubscriptionSchedulePhaseItemParams{{Price: stripelib.String(*plan.StripePlanID)}},
		StartDate: stripelib.Int64(effectiveTime.Unix()),
	}
	if prorationMode != "" {
		currentPhase.ProrationBehavior = stripelib.String(prorationMode)
	}
	scheduleParams.Phases = []*stripelib.SubscriptionSchedulePhaseParams{currentPhase}
	schedule, err := withStripeRetry(ctx, 3, func() (*stripelib.SubscriptionSchedule, error) { return subscriptionschedule.New(scheduleParams) })
	if err != nil {
		l.Error("failed to create subscription schedule", "error", err)
		return nil, ErrStripeOperationFailed(err, "create_subscription_schedule")
	}
	l.Info("subscription schedule created successfully", "scheduleId", schedule.ID)
	return schedule, nil
}

// EstimateProrationAmount estimates the proration amount for plan changes
func (s *Service) EstimateProrationAmount(ctx context.Context, ws *account.Workspace, newPlanDescriptor string) (int64, error) {
	l := logger.FromContext(ctx).With("workspaceId", ws.ID, "newPlan", newPlanDescriptor)
	l.Info("estimating proration amount")
	currentSub, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		l.Error("failed to get current subscription", "error", err)
		return 0, err
	}
	currentPlan, err := s.GetPlanByID(ctx, currentSub.PlanID)
	if err != nil {
		l.Error("failed to get current plan", "error", err)
		return 0, err
	}
	newPlan, err := s.GetPlanByDescriptor(ctx, newPlanDescriptor)
	if err != nil {
		l.Error("failed to get new plan", "error", err)
		return 0, err
	}
	currentPricePerMonth := currentPlan.Price.IntPart()
	newPricePerMonth := newPlan.Price.IntPart()
	daysRemaining := int64(helpers.CeilPositiveDaysUntil(currentSub.CurrentPeriodEnd))
	daysInMonth := int64(30)
	prorationAmount := ((newPricePerMonth - currentPricePerMonth) * daysRemaining) / daysInMonth
	l.Info("proration amount estimated", "amount", prorationAmount)
	return prorationAmount, nil
}

// CancelSubscription cancels a subscription immediately (alias)
func (s *Service) CancelSubscription(ctx context.Context, ws *account.Workspace) error {
	return s.CancelSubscriptionImmediately(ctx, ws)
}
