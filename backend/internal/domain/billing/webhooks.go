package billing

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
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
)

type webhookEvent struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Data struct {
		Object json.RawMessage `json:"object"`
	} `json:"data"`
}

func verifyStripeSignature(payload []byte, signatureHeader, secret string, tolerance time.Duration) error {
	if strings.TrimSpace(signatureHeader) == "" {
		return errors.New("webhook has no Stripe-Signature header")
	}
	if secret == "" {
		return errors.New("webhook secret is not configured")
	}

	var tsStr string
	var sigHex string
	for _, part := range strings.Split(signatureHeader, ",") {
		p := strings.TrimSpace(part)
		if strings.HasPrefix(p, "t=") {
			tsStr = strings.TrimPrefix(p, "t=")
			continue
		}
		if strings.HasPrefix(p, "v1=") {
			sigHex = strings.TrimPrefix(p, "v1=")
			continue
		}
	}
	if tsStr == "" || sigHex == "" {
		return errors.New("invalid Stripe-Signature header format")
	}

	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid Stripe-Signature timestamp: %w", err)
	}
	now := time.Now().Unix()
	if tolerance > 0 {
		delta := now - ts
		if delta < 0 {
			delta = -delta
		}
		if time.Duration(delta)*time.Second > tolerance {
			return errors.New("timestamp wasn't within tolerance")
		}
	}

	signedPayload := append([]byte(tsStr+"."), payload...)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(signedPayload)
	expected := mac.Sum(nil)

	given, err := hex.DecodeString(sigHex)
	if err != nil {
		return fmt.Errorf("invalid Stripe-Signature v1 value: %w", err)
	}
	if subtle.ConstantTimeCompare(expected, given) != 1 {
		return errors.New("no signatures found matching the expected signature")
	}

	return nil
}

// ProcessWebhook verifies signature, ensures idempotency, and dispatches typed handlers
func (s *Service) ProcessWebhook(ctx context.Context, payload []byte, signature string) error {
	log := logger.FromContext(ctx)
	secret := viper.GetString(config.StripeWebhookSecret)
	if secret == "" {
		log.Warn("Stripe webhook secret not configured; rejecting webhook for security")
		return ErrWebhookProcessingFailed(nil, "missing_webhook_secret")
	}
	if err := verifyStripeSignature(payload, signature, secret, 5*time.Minute); err != nil {
		log.Error("webhook signature verification failed", "error", err)
		return ErrWebhookSignatureInvalid(err)
	}

	var evt webhookEvent
	if err := json.Unmarshal(payload, &evt); err != nil {
		log.Error("failed to parse webhook payload", "error", err)
		return ErrWebhookPayloadInvalid(err)
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
		if err := s.handleSubscriptionCreated(ctx, evt.Data.Object); err != nil {
			return err
		}
	case "customer.subscription.updated":
		if err := s.handleSubscriptionUpdated(ctx, evt.Data.Object); err != nil {
			return err
		}
	case "customer.subscription.deleted":
		if err := s.handleSubscriptionDeleted(ctx, evt.Data.Object); err != nil {
			return err
		}
	case "invoice.payment_succeeded":
		if err := s.handleInvoicePaymentSucceeded(ctx, evt.Data.Object); err != nil {
			return err
		}
	case "invoice.payment_failed":
		if err := s.handleInvoicePaymentFailed(ctx, evt.Data.Object); err != nil {
			return err
		}
	case "invoice.finalized":
		if err := s.handleInvoiceFinalized(ctx, evt.Data.Object); err != nil {
			return err
		}
	case "invoice.marked_uncollectible":
		if err := s.handleInvoiceMarkedUncollectible(ctx, evt.Data.Object); err != nil {
			return err
		}
	case "invoice.voided":
		if err := s.handleInvoiceVoided(ctx, evt.Data.Object); err != nil {
			return err
		}
	case "customer.subscription.trial_will_end":
		if err := s.handleTrialWillEnd(ctx, evt.Data.Object); err != nil {
			return err
		}
	case "payment_method.automatically_updated":
		if err := s.handlePaymentMethodAutomaticallyUpdated(ctx, evt.Data.Object); err != nil {
			return err
		}
	case "checkout.session.completed":
		if err := s.handleCheckoutSessionCompleted(ctx, evt.Data.Object); err != nil {
			return err
		}
	default:
		// unhandled types are ignored
	}

	// Mark processed
	if evt.ID != "" {
		se := &StripeEvent{EventID: evt.ID, Type: evt.Type, ProcessedAt: time.Now()}
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
	if err := s.SyncSubscriptionStatus(ctx, id, status, start, end); err != nil {
		logger.FromContext(ctx).Error("failed to sync subscription status", "error", err, "stripeSubId", id)
		return ErrWebhookProcessingFailed(err, "sync_subscription")
	}
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
	if err := s.RefundAndFinalizeCancellation(ctx, id, start, end); err != nil {
		logger.FromContext(ctx).Error("failed to finalize cancellation", "error", err, "stripeSubId", id)
		return ErrWebhookProcessingFailed(err, "finalize_cancellation")
	}
	return nil
}

func (s *Service) handleInvoicePaymentSucceeded(ctx context.Context, raw json.RawMessage) error {
	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		return ErrWebhookProcessingFailed(err, "parse_payload")
	}
	subID := cast.ToString(obj["subscription"])
	if subID != "" {
		if err := s.MarkSubscriptionActive(ctx, subID); err != nil {
			logger.FromContext(ctx).Error("failed to mark subscription active", "error", err, "stripeSubId", subID)
			return ErrWebhookProcessingFailed(err, "mark_subscription_active")
		}
		if s.Notification != nil {
			l := logger.FromContext(ctx)
			go func(subscriptionID string, createdAt int64) {
				bg, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()
				subscription, err := s.storage.subscription.FindOne(bg, s.storage.subscription.ScopeEquals(SubscriptionSchema.StripeSubID, subscriptionID))
				if err != nil || subscription == nil {
					if err != nil {
						l.Warn("failed to load subscription for payment notification", "error", err, "stripeSubId", subscriptionID)
					}
					return
				}
				plan, err := s.GetPlanByID(bg, subscription.PlanID)
				if err != nil || plan == nil {
					if err != nil {
						l.Warn("failed to load plan for payment notification", "error", err, "planId", subscription.PlanID)
					}
					return
				}
				invoiceURL := cast.ToString(obj["hosted_invoice_url"])
				if invoiceURL == "" {
					invoiceURL = cast.ToString(obj["invoice_pdf"])
				}
				payAt := time.Now().UTC()
				if createdAt != 0 {
					payAt = time.Unix(createdAt, 0).UTC()
				}
				last4 := "****" // masked; extraction omitted
				isFirst := time.Since(subscription.CreatedAt) < 24*time.Hour
				if isFirst {
					if err := s.Notification.SendSubscriptionConfirmedEmail(bg, subscription.WorkspaceID, subscription, plan, last4, invoiceURL); err != nil {
						l.Warn("failed to send subscription confirmed email", "error", err, "workspaceId", subscription.WorkspaceID)
					}
				} else {
					if err := s.Notification.SendPaymentSucceededEmail(bg, subscription.WorkspaceID, subscription, plan, last4, invoiceURL, payAt); err != nil {
						l.Warn("failed to send payment succeeded email", "error", err, "workspaceId", subscription.WorkspaceID)
					}
				}
			}(subID, cast.ToInt64(obj["created"]))
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
		if err := s.MarkSubscriptionPastDue(ctx, subID); err != nil {
			logger.FromContext(ctx).Error("failed to mark subscription past due", "error", err, "stripeSubId", subID)
			return ErrWebhookProcessingFailed(err, "mark_subscription_past_due")
		}
		if s.Notification != nil {
			l := logger.FromContext(ctx)
			go func(subscriptionID string) {
				bg, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()
				subscription, err := s.storage.subscription.FindOne(bg, s.storage.subscription.ScopeEquals(SubscriptionSchema.StripeSubID, subscriptionID))
				if err != nil || subscription == nil {
					if err != nil {
						l.Warn("failed to load subscription for payment failed notification", "error", err, "stripeSubId", subscriptionID)
					}
					return
				}
				plan, err := s.GetPlanByID(bg, subscription.PlanID)
				if err != nil || plan == nil {
					if err != nil {
						l.Warn("failed to load plan for payment failed notification", "error", err, "planId", subscription.PlanID)
					}
					return
				}
				if err := s.Notification.SendPaymentFailedEmail(bg, subscription.WorkspaceID, subscription, plan, "****", time.Now().UTC(), nil); err != nil {
					l.Warn("failed to send payment failed email", "error", err, "workspaceId", subscription.WorkspaceID)
				}
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
			l := logger.FromContext(ctx)
			go func(subscriptionID string) {
				bg, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()
				subscription, err := s.storage.subscription.FindOne(bg, s.storage.subscription.ScopeEquals(SubscriptionSchema.StripeSubID, subscriptionID))
				if err != nil || subscription == nil {
					if err != nil {
						l.Warn("failed to load subscription for trial ending notification", "error", err, "stripeSubId", subscriptionID)
					}
					return
				}
				plan, err := s.GetPlanByID(bg, subscription.PlanID)
				if err != nil || plan == nil {
					if err != nil {
						l.Warn("failed to load plan for trial ending notification", "error", err, "planId", subscription.PlanID)
					}
					return
				}
				trialInfo, err := s.CheckTrialStatus(bg, &account.Workspace{ID: subscription.WorkspaceID})
				if err != nil {
					l.Warn("failed to check trial status", "error", err, "workspaceId", subscription.WorkspaceID)
					return
				}
				if err := s.Notification.SendTrialEndingEmail(bg, subscription.WorkspaceID, subscription, plan, trialInfo); err != nil {
					l.Warn("failed to send trial ending email", "error", err, "workspaceId", subscription.WorkspaceID)
				}
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
		if err := s.MarkSubscriptionPastDue(ctx, subID); err != nil {
			logger.FromContext(ctx).Error("failed to mark subscription past due", "error", err, "stripeSubId", subID)
			return ErrWebhookProcessingFailed(err, "mark_subscription_past_due")
		}
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
				if _, err := withStripeRetry[*stripelib.Customer](ctx, 3, func() (*stripelib.Customer, error) {
					return customer.Update(cid, &stripelib.CustomerParams{InvoiceSettings: &stripelib.CustomerInvoiceSettingsParams{DefaultPaymentMethod: stripelib.String(pm.ID)}})
				}); err != nil {
					logger.FromContext(ctx).Warn("failed to update customer default payment method", "error", err, "stripeCustomerId", cid, "stripePaymentMethodId", pm.ID)
				}
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
			if _, err := withStripeRetry[*stripelib.Customer](ctx, 3, func() (*stripelib.Customer, error) {
				return customer.Update(sess.Customer.ID, &stripelib.CustomerParams{InvoiceSettings: &stripelib.CustomerInvoiceSettingsParams{DefaultPaymentMethod: stripelib.String(pmID)}})
			}); err != nil {
				logger.FromContext(ctx).Warn("failed to set customer invoice settings default payment method", "error", err, "stripeCustomerId", sess.Customer.ID, "stripePaymentMethodId", pmID)
			}
		}
	}
	if sess.Subscription != nil {
		if err := s.MarkSubscriptionActive(ctx, sess.Subscription.ID); err != nil {
			logger.FromContext(ctx).Error("failed to mark subscription active", "error", err, "stripeSubId", sess.Subscription.ID)
			return ErrWebhookProcessingFailed(err, "mark_subscription_active")
		}
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
