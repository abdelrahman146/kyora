package billing

import (
	"context"
	"fmt"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
	atomic "github.com/abdelrahman146/kyora/internal/platform/types/atomic"
	stripelib "github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/creditnote"
	"github.com/stripe/stripe-go/v83/invoice"
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

func (s *Service) GetSubscriptionByID(ctx context.Context, id string) (*Subscription, error) {
	return s.storage.subscription.FindByID(ctx, id)
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
	if ws == nil {
		return nil, fmt.Errorf("workspace cannot be nil")
	}
	if plan == nil {
		return nil, fmt.Errorf("plan cannot be nil")
	}
	existing, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil && !database.IsRecordNotFound(err) {
		return nil, fmt.Errorf("failed to check existing subscription: %w", err)
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
				Items:             []*stripelib.SubscriptionItemsParams{{Price: stripelib.String(plan.StripePlanID)}},
				ProrationBehavior: stripelib.String("create_prorations"),
				CancelAtPeriodEnd: stripelib.Bool(false),
				Metadata:          map[string]string{"workspace_id": ws.ID, "plan_id": plan.ID},
			}
			updateParams.SetIdempotencyKey(idempotencyKey)
			stripeSub, err = subscription.Update(existing.StripeSubID, updateParams)
			if err != nil {
				logger.FromContext(ctx).Error("Failed to update Stripe subscription", "error", err, "subscription_id", existing.StripeSubID, "plan_id", plan.ID)
				return fmt.Errorf("failed to update subscription: %w", err)
			}
			existing.PlanID = plan.ID
			existing.Status = mapStripeStatus(stripeSub.Status)
			if err := s.storage.subscription.UpdateOne(ctx, existing); err != nil {
				return fmt.Errorf("failed to update local subscription: %w", err)
			}
			result = existing
			logger.FromContext(ctx).Info("Updated subscription", "workspace_id", ws.ID, "subscription_id", existing.StripeSubID, "new_plan_id", plan.ID)
			return nil
		}
		idempotencyKey := fmt.Sprintf("sub_create_%s_%s", ws.ID, plan.ID)
		createParams := &stripelib.SubscriptionParams{
			Customer:          stripelib.String(custID),
			Items:             []*stripelib.SubscriptionItemsParams{{Price: stripelib.String(plan.StripePlanID)}},
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
		stripeSub, err = subscription.New(createParams)
		if err != nil {
			logger.FromContext(ctx).Error("Failed to create Stripe subscription", "error", err, "customer_id", custID, "plan_id", plan.ID)
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
		logger.FromContext(ctx).Info("Created new subscription", "workspace_id", ws.ID, "subscription_id", stripeSub.ID, "plan_id", plan.ID)
		return nil
	}, atomic.WithRetries(2))
	if err != nil {
		return nil, err
	}
	if result != nil && existing == nil && s.notification != nil {
		paymentMethodLastFour := ""
		if details, err := s.GetSubscriptionDetails(ctx, ws); err == nil && details.PaymentMethod.Last4 != "" {
			paymentMethodLastFour = details.PaymentMethod.Last4
		}
		go func() {
			if err := s.notification.SendSubscriptionWelcomeEmail(context.Background(), ws.ID, result, plan, paymentMethodLastFour); err != nil {
				logger.FromContext(ctx).Error("Failed to send subscription welcome email", "error", err, "workspace_id", ws.ID, "subscription_id", result.ID)
			}
		}()
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
		_, err = subscription.Cancel(subRec.StripeSubID, cancelParams)
		if err != nil {
			logger.FromContext(ctx).Error("Failed to cancel Stripe subscription", "error", err, "subscription_id", subRec.StripeSubID, "workspace_id", ws.ID)
		}
		subRec.Status = SubscriptionStatusCanceled
		subRec.CurrentPeriodEnd = time.Now()
		if updateErr := s.storage.subscription.UpdateOne(ctx, subRec); updateErr != nil {
			logger.FromContext(ctx).Error("Failed to update local subscription status", "error", updateErr, "subscription_id", subRec.StripeSubID)
			return fmt.Errorf("failed to update local subscription: %w", updateErr)
		}
		if err != nil {
			return fmt.Errorf("failed to cancel Stripe subscription: %w", err)
		}
		logger.FromContext(ctx).Info("Successfully canceled subscription", "workspace_id", ws.ID, "subscription_id", subRec.StripeSubID)
		if s.notification != nil {
			if plan, planErr := s.GetPlanByID(ctx, subRec.PlanID); planErr == nil {
				go func() {
					if err := s.notification.SendSubscriptionCanceledEmail(context.Background(), ws.ID, subRec, plan, time.Now(), ""); err != nil {
						logger.FromContext(ctx).Error("Failed to send subscription canceled email", "error", err, "workspace_id", ws.ID, "subscription_id", subRec.StripeSubID)
					}
				}()
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
	ip := &stripelib.InvoiceListParams{Subscription: stripelib.String(stripeSubID)}
	ip.Status = stripelib.String(string(stripelib.InvoiceStatusPaid))
	ip.Limit = stripelib.Int64(1)
	iter := invoice.List(ip)
	if iter.Next() {
		inv := iter.Invoice()
		if periodStart > 0 && periodEnd > 0 {
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
							_, rerr := creditnote.New(&stripelib.CreditNoteParams{
								Invoice:      stripelib.String(inv.ID),
								Amount:       stripelib.Int64(refundAmount),
								RefundAmount: stripelib.Int64(refundAmount),
								Reason:       stripelib.String(string(stripelib.CreditNoteReasonOrderChange)),
								Memo:         stripelib.String("Prorated refund for immediate cancellation"),
								Metadata:     map[string]string{"kyoraRefundKind": "cancel_prorated", "subscription": stripeSubID, "workspaceId": rec.WorkspaceID},
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
		} else {
			logger.FromContext(ctx).Info("Skipping refund calculation: missing period bounds", "subscription", stripeSubID)
		}
	}
	_, _ = subscription.Cancel(stripeSubID, &stripelib.SubscriptionCancelParams{InvoiceNow: stripelib.Bool(false), Prorate: stripelib.Bool(false)})
	rec.Status = SubscriptionStatusCanceled
	rec.CurrentPeriodEnd = time.Now()
	_ = s.storage.subscription.UpdateOne(ctx, rec)
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
func (s *Service) ResumeSubscriptionIfNoDue(ctx context.Context, ws *account.Workspace) (*Subscription, error) {
	rec, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		return nil, ErrSubscriptionNotFound(err, ws.ID)
	}
	if rec.Status != SubscriptionStatusCanceled {
		return rec, nil
	}
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
	plan, err := s.GetPlanByID(ctx, rec.PlanID)
	if err != nil {
		return nil, err
	}
	return s.CreateOrUpdateSubscription(ctx, ws, plan)
}

// ScheduleSubscriptionChange schedules a subscription change for a future date
func (s *Service) ScheduleSubscriptionChange(ctx context.Context, ws *account.Workspace, plan *Plan, effectiveDate, prorationMode string) (*stripelib.SubscriptionSchedule, error) {
	l := logger.FromContext(ctx).With("workspace_id", ws.ID, "plan_descriptor", plan.Descriptor, "effective_date", effectiveDate)
	l.Info("scheduling subscription change")
	sub, err := s.GetSubscriptionByWorkspaceID(ctx, ws.ID)
	if err != nil {
		l.Error("failed to get current subscription", "error", err)
		return nil, err
	}
	effectiveTime, err := time.Parse("2006-01-02T15:04:05Z", effectiveDate)
	if err != nil {
		if effectiveTime, err = time.Parse("2006-01-02", effectiveDate); err != nil {
			l.Error("invalid effective date format", "error", err)
			return nil, ErrStripeOperationFailed(err, "parse_date")
		}
	}
	if err := s.ensurePlanSynced(ctx, plan); err != nil {
		l.Error("failed to ensure plan synced", "error", err)
		return nil, err
	}
	scheduleParams := &stripelib.SubscriptionScheduleParams{FromSubscription: stripelib.String(sub.StripeSubID)}
	currentPhase := &stripelib.SubscriptionSchedulePhaseParams{
		Items:     []*stripelib.SubscriptionSchedulePhaseItemParams{{Price: stripelib.String(plan.StripePlanID)}},
		StartDate: stripelib.Int64(effectiveTime.Unix()),
	}
	if prorationMode != "" {
		currentPhase.ProrationBehavior = stripelib.String(prorationMode)
	}
	scheduleParams.Phases = []*stripelib.SubscriptionSchedulePhaseParams{currentPhase}
	schedule, err := subscriptionschedule.New(scheduleParams)
	if err != nil {
		l.Error("failed to create subscription schedule", "error", err)
		return nil, ErrStripeOperationFailed(err, "create_subscription_schedule")
	}
	l.Info("subscription schedule created successfully", "schedule_id", schedule.ID)
	return schedule, nil
}

// EstimateProrationAmount estimates the proration amount for plan changes
func (s *Service) EstimateProrationAmount(ctx context.Context, ws *account.Workspace, newPlanDescriptor string) (int64, error) {
	l := logger.FromContext(ctx).With("workspace_id", ws.ID, "new_plan", newPlanDescriptor)
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
	daysRemaining := int64(time.Until(currentSub.CurrentPeriodEnd).Hours() / 24)
	daysInMonth := int64(30)
	prorationAmount := ((newPricePerMonth - currentPricePerMonth) * daysRemaining) / daysInMonth
	l.Info("proration amount estimated", "amount", prorationAmount)
	return prorationAmount, nil
}

// CancelSubscription cancels a subscription immediately (alias)
func (s *Service) CancelSubscription(ctx context.Context, ws *account.Workspace) error {
	return s.CancelSubscriptionImmediately(ctx, ws)
}
