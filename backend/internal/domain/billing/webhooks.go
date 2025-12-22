package billing

import (
	"context"
	"encoding/json"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/bus"
	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	stripelib "github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/checkout/session"
	"github.com/stripe/stripe-go/v83/customer"
	"github.com/stripe/stripe-go/v83/paymentmethod"
	"github.com/stripe/stripe-go/v83/setupintent"
	"github.com/stripe/stripe-go/v83/webhook"
)

// ProcessWebhook verifies signature, ensures idempotency, and dispatches typed handlers
func (s *Service) ProcessWebhook(ctx context.Context, payload []byte, signature string) error {
	log := logger.FromContext(ctx)
	secret := viper.GetString(config.StripeWebhookSecret)
	if secret == "" {
		log.Warn("Stripe webhook secret not configured; rejecting webhook for security")
		return ErrWebhookProcessingFailed(nil, "missing_webhook_secret")
	}

	// Verify signature and construct event
	evt, err := webhook.ConstructEvent(payload, signature, secret)
	if err != nil {
		log.Error("webhook signature verification failed", "error", err)
		return ErrWebhookProcessingFailed(err, "verify_signature")
	}

	// Idempotency: skip already processed events
	if evt.ID != "" {
		if _, err := s.storage.event.FindOne(ctx, s.storage.event.ScopeEquals(StripeEventSchema.EventID, evt.ID)); err == nil {
			// already processed
			return nil
		}
	}

	// Dispatch by type
	switch evt.Type {
	case "customer.subscription.created":
		if err := s.handleSubscriptionCreated(ctx, evt.Data.Raw); err != nil {
			return err
		}
	case "customer.subscription.updated":
		if err := s.handleSubscriptionUpdated(ctx, evt.Data.Raw); err != nil {
			return err
		}
	case "customer.subscription.deleted":
		if err := s.handleSubscriptionDeleted(ctx, evt.Data.Raw); err != nil {
			return err
		}
	case "invoice.payment_succeeded":
		if err := s.handleInvoicePaymentSucceeded(ctx, evt.Data.Raw); err != nil {
			return err
		}
	case "invoice.payment_failed":
		if err := s.handleInvoicePaymentFailed(ctx, evt.Data.Raw); err != nil {
			return err
		}
	case "invoice.finalized":
		if err := s.handleInvoiceFinalized(ctx, evt.Data.Raw); err != nil {
			return err
		}
	case "invoice.marked_uncollectible":
		if err := s.handleInvoiceMarkedUncollectible(ctx, evt.Data.Raw); err != nil {
			return err
		}
	case "invoice.voided":
		if err := s.handleInvoiceVoided(ctx, evt.Data.Raw); err != nil {
			return err
		}
	case "customer.subscription.trial_will_end":
		if err := s.handleTrialWillEnd(ctx, evt.Data.Raw); err != nil {
			return err
		}
	case "payment_method.automatically_updated":
		if err := s.handlePaymentMethodAutomaticallyUpdated(ctx, evt.Data.Raw); err != nil {
			return err
		}
	case "checkout.session.completed":
		if err := s.handleCheckoutSessionCompleted(ctx, evt.Data.Raw); err != nil {
			return err
		}
	default:
		// unhandled types are ignored
	}

	// Mark processed
	if evt.ID != "" {
		se := &StripeEvent{EventID: evt.ID, Type: string(evt.Type), ProcessedAt: time.Now()}
		if err := s.storage.event.CreateOne(ctx, se); err != nil && !database.IsUniqueViolation(err) {
			// Log but continue
			logger.FromContext(ctx).Warn("failed to persist webhook event id", "error", err, "eventId", evt.ID)
		}
	}
	return nil
}

// Old helpers rewritten to accept raw JSON object only (from verified event)
func (s *Service) handleSubscriptionCreated(ctx context.Context, raw json.RawMessage) error {
	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		return ErrWebhookProcessingFailed(err, "parse_payload")
	}
	id := cast.ToString(obj["id"])
	status := cast.ToString(obj["status"])
	start := cast.ToInt64(obj["current_period_start"])
	end := cast.ToInt64(obj["current_period_end"])
	s.SyncSubscriptionStatus(ctx, id, status, start, end)
	return nil
}

func (s *Service) handleSubscriptionUpdated(ctx context.Context, raw json.RawMessage) error {
	return s.handleSubscriptionCreated(ctx, raw)
}

func (s *Service) handleSubscriptionDeleted(ctx context.Context, raw json.RawMessage) error {
	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		return ErrWebhookProcessingFailed(err, "parse_payload")
	}
	id := cast.ToString(obj["id"])
	start := cast.ToInt64(obj["current_period_start"])
	end := cast.ToInt64(obj["current_period_end"])
	s.RefundAndFinalizeCancellation(ctx, id, start, end)
	return nil
}

func (s *Service) handleInvoicePaymentSucceeded(ctx context.Context, raw json.RawMessage) error {
	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		return ErrWebhookProcessingFailed(err, "parse_payload")
	}
	subID := cast.ToString(obj["subscription"])
	if subID != "" {
		s.MarkSubscriptionActive(ctx, subID)
		if s.Notification != nil {
			go func(subscriptionID string) {
				subscription, err := s.storage.subscription.FindOne(context.Background(), s.storage.subscription.ScopeEquals(SubscriptionSchema.StripeSubID, subscriptionID))
				if err != nil || subscription == nil {
					return
				}
				plan, err := s.GetPlanByID(context.Background(), subscription.PlanID)
				if err != nil || plan == nil {
					return
				}
				invoiceURL := cast.ToString(obj["hosted_invoice_url"])
				if invoiceURL == "" {
					invoiceURL = cast.ToString(obj["invoice_pdf"])
				}
				payAt := time.Now()
				if v := cast.ToInt64(obj["created"]); v != 0 {
					payAt = time.Unix(v, 0)
				}
				last4 := "****" // masked; extraction omitted
				isFirst := time.Since(subscription.CreatedAt) < 24*time.Hour
				if isFirst {
					_ = s.Notification.SendSubscriptionConfirmedEmail(context.Background(), subscription.WorkspaceID, subscription, plan, last4, invoiceURL)
				} else {
					_ = s.Notification.SendPaymentSucceededEmail(context.Background(), subscription.WorkspaceID, subscription, plan, last4, invoiceURL, payAt)
				}
			}(subID)
		}
	}
	return nil
}

func (s *Service) handleInvoicePaymentFailed(ctx context.Context, raw json.RawMessage) error {
	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		return ErrWebhookProcessingFailed(err, "parse_payload")
	}
	subID := cast.ToString(obj["subscription"])
	if subID != "" {
		s.MarkSubscriptionPastDue(ctx, subID)
		if s.Notification != nil {
			go func(subscriptionID string) {
				subscription, err := s.storage.subscription.FindOne(context.Background(), s.storage.subscription.ScopeEquals(SubscriptionSchema.StripeSubID, subscriptionID))
				if err != nil || subscription == nil {
					return
				}
				plan, err := s.GetPlanByID(context.Background(), subscription.PlanID)
				if err != nil || plan == nil {
					return
				}
				_ = s.Notification.SendPaymentFailedEmail(context.Background(), subscription.WorkspaceID, subscription, plan, "****", time.Now(), nil)
			}(subID)
		}
	}
	return nil
}

func (s *Service) handleTrialWillEnd(ctx context.Context, raw json.RawMessage) error {
	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		return ErrWebhookProcessingFailed(err, "parse_payload")
	}
	subID := cast.ToString(obj["id"])
	if subID != "" {
		if s.Notification != nil {
			go func(subscriptionID string) {
				subscription, err := s.storage.subscription.FindOne(context.Background(), s.storage.subscription.ScopeEquals(SubscriptionSchema.StripeSubID, subscriptionID))
				if err != nil || subscription == nil {
					return
				}
				plan, err := s.GetPlanByID(context.Background(), subscription.PlanID)
				if err != nil || plan == nil {
					return
				}
				trialInfo, err := s.CheckTrialStatus(context.Background(), &account.Workspace{ID: subscription.WorkspaceID})
				if err != nil {
					return
				}
				_ = s.Notification.SendTrialEndingEmail(context.Background(), subscription.WorkspaceID, subscription, plan, trialInfo)
			}(subID)
		}
	}
	return nil
}

func (s *Service) handleInvoiceFinalized(ctx context.Context, raw json.RawMessage) error {
	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		return ErrWebhookProcessingFailed(err, "parse_payload")
	}
	logger.FromContext(ctx).Info("invoice finalized", "invoiceId", obj["id"])
	return nil
}

func (s *Service) handleInvoiceMarkedUncollectible(ctx context.Context, raw json.RawMessage) error {
	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		return ErrWebhookProcessingFailed(err, "parse_payload")
	}
	if subID, ok := obj["subscription"].(string); ok && subID != "" {
		// Downgrade local status so UI can show attention needed
		s.MarkSubscriptionPastDue(ctx, subID)
	}
	return nil
}

func (s *Service) handleInvoiceVoided(ctx context.Context, raw json.RawMessage) error {
	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		return ErrWebhookProcessingFailed(err, "parse_payload")
	}
	logger.FromContext(ctx).Info("invoice voided", "invoiceId", obj["id"])
	return nil
}

func (s *Service) handlePaymentMethodAutomaticallyUpdated(ctx context.Context, raw json.RawMessage) error {
	var pm struct {
		ID       string `json:"id"`
		Customer any    `json:"customer"`
	}
	if err := json.Unmarshal(raw, &pm); err != nil {
		return ErrWebhookProcessingFailed(err, "parse_payload")
	}
	if pm.ID == "" {
		return nil
	}
	if c, ok := pm.Customer.(map[string]any); ok {
		cid := cast.ToString(c["id"])
		if cid != "" {
			cu, err := customer.Get(cid, nil)
			if err == nil && (cu.InvoiceSettings == nil || cu.InvoiceSettings.DefaultPaymentMethod == nil) {
				_, _ = withStripeRetry[*stripelib.Customer](ctx, 3, func() (*stripelib.Customer, error) {
					return customer.Update(cid, &stripelib.CustomerParams{InvoiceSettings: &stripelib.CustomerInvoiceSettingsParams{DefaultPaymentMethod: stripelib.String(pm.ID)}})
				})
			}
		}
	}
	return nil
}

func (s *Service) handleCheckoutSessionCompleted(ctx context.Context, raw json.RawMessage) error {
	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		return ErrWebhookProcessingFailed(err, "parse_payload")
	}
	sid := cast.ToString(obj["id"])
	if sid == "" {
		return nil
	}
	sess, err := session.Get(sid, nil)
	if err != nil {
		return ErrWebhookProcessingFailed(err, "get_session")
	}
	if sess.SetupIntent != nil && sess.Customer != nil {
		si, err := setupintent.Get(sess.SetupIntent.ID, nil)
		if err != nil {
			return ErrWebhookProcessingFailed(err, "get_setup_intent")
		}
		if si.PaymentMethod != nil {
			pmID := si.PaymentMethod.ID
			if _, err := withStripeRetry[*stripelib.PaymentMethod](ctx, 3, func() (*stripelib.PaymentMethod, error) {
				return paymentmethod.Attach(pmID, &stripelib.PaymentMethodAttachParams{Customer: stripelib.String(sess.Customer.ID)})
			}); err != nil {
				return ErrWebhookProcessingFailed(err, "attach_payment_method")
			}
			_, _ = withStripeRetry[*stripelib.Customer](ctx, 3, func() (*stripelib.Customer, error) {
				return customer.Update(sess.Customer.ID, &stripelib.CustomerParams{InvoiceSettings: &stripelib.CustomerInvoiceSettingsParams{DefaultPaymentMethod: stripelib.String(pmID)}})
			})
		}
	}
	if sess.Subscription != nil {
		s.MarkSubscriptionActive(ctx, sess.Subscription.ID)
		// If this checkout was part of onboarding, notify the onboarding flow via event
		if sess.Metadata != nil {
			if obID, ok := sess.Metadata["onboarding_session_id"]; ok && obID != "" {
				s.bus.Emit(bus.OnboardingPaymentSucceededTopic, &bus.OnboardingPaymentSucceededEvent{
					OnboardingSessionID:  obID,
					StripeCheckoutID:     sid,
					StripeSubscriptionID: sess.Subscription.ID,
				})
			}
		}
	}
	return nil
}
