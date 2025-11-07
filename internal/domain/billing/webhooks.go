package billing

import (
	"context"
	"encoding/json"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
)

// ProcessWebhook processes incoming Stripe webhooks
func (s *Service) ProcessWebhook(ctx context.Context, payload []byte, signature string) error {
	logger := logger.FromContext(ctx).With("event_type", "webhook")
	logger.Info("processing stripe webhook")

	// Note: In a real implementation, you'd need the webhook endpoint secret
	// For now, we'll skip signature verification in development
	// webhookSecret := viper.GetString("stripe.webhook_secret")

	// Parse the webhook event
	var event map[string]interface{}
	if err := json.Unmarshal(payload, &event); err != nil {
		logger.Error("failed to parse webhook payload", "error", err)
		return ErrWebhookProcessingFailed(err, "parse_payload")
	}

	eventType, ok := event["type"].(string)
	if !ok {
		logger.Error("webhook event missing type field")
		return ErrWebhookProcessingFailed(nil, "missing_event_type")
	}

	logger = logger.With("event_type", eventType)
	logger.Info("processing webhook event")

	// Handle different webhook events
	switch eventType {
	case "customer.subscription.created":
		return s.handleSubscriptionCreated(ctx, event)
	case "customer.subscription.updated":
		return s.handleSubscriptionUpdated(ctx, event)
	case "customer.subscription.deleted":
		return s.handleSubscriptionDeleted(ctx, event)
	case "invoice.payment_succeeded":
		return s.handleInvoicePaymentSucceeded(ctx, event)
	case "invoice.payment_failed":
		return s.handleInvoicePaymentFailed(ctx, event)
	case "customer.subscription.trial_will_end":
		return s.handleTrialWillEnd(ctx, event)
	default:
		logger.Info("unhandled webhook event type", "event_type", eventType)
		return nil // Not an error, just unhandled
	}
}

func (s *Service) handleSubscriptionCreated(ctx context.Context, event map[string]interface{}) error {
	logger := logger.FromContext(ctx).With("event", "subscription_created")
	logger.Info("handling subscription created webhook")

	// Extract subscription data
	data, ok := event["data"].(map[string]interface{})
	if !ok {
		return ErrWebhookProcessingFailed(nil, "invalid_data_format")
	}

	object, ok := data["object"].(map[string]interface{})
	if !ok {
		return ErrWebhookProcessingFailed(nil, "invalid_object_format")
	}

	subscriptionID, ok := object["id"].(string)
	if !ok {
		return ErrWebhookProcessingFailed(nil, "missing_subscription_id")
	}

	status, ok := object["status"].(string)
	if !ok {
		return ErrWebhookProcessingFailed(nil, "missing_subscription_status")
	}

	// Extract period information
	currentPeriodStart := int64(0)
	currentPeriodEnd := int64(0)
	if period, ok := object["current_period_start"].(float64); ok {
		currentPeriodStart = int64(period)
	}
	if period, ok := object["current_period_end"].(float64); ok {
		currentPeriodEnd = int64(period)
	}

	s.SyncSubscriptionStatus(ctx, subscriptionID, status, currentPeriodStart, currentPeriodEnd)
	logger.Info("subscription created webhook processed successfully")
	return nil
}

func (s *Service) handleSubscriptionUpdated(ctx context.Context, event map[string]interface{}) error {
	logger := logger.FromContext(ctx).With("event", "subscription_updated")
	logger.Info("handling subscription updated webhook")

	// Similar to created, but may need to handle plan changes
	data, ok := event["data"].(map[string]interface{})
	if !ok {
		return ErrWebhookProcessingFailed(nil, "invalid_data_format")
	}

	object, ok := data["object"].(map[string]interface{})
	if !ok {
		return ErrWebhookProcessingFailed(nil, "invalid_object_format")
	}

	subscriptionID, ok := object["id"].(string)
	if !ok {
		return ErrWebhookProcessingFailed(nil, "missing_subscription_id")
	}

	status, ok := object["status"].(string)
	if !ok {
		return ErrWebhookProcessingFailed(nil, "missing_subscription_status")
	}

	// Extract period information
	currentPeriodStart := int64(0)
	currentPeriodEnd := int64(0)
	if period, ok := object["current_period_start"].(float64); ok {
		currentPeriodStart = int64(period)
	}
	if period, ok := object["current_period_end"].(float64); ok {
		currentPeriodEnd = int64(period)
	}

	s.SyncSubscriptionStatus(ctx, subscriptionID, status, currentPeriodStart, currentPeriodEnd)
	logger.Info("subscription updated webhook processed successfully")
	return nil
}

func (s *Service) handleSubscriptionDeleted(ctx context.Context, event map[string]interface{}) error {
	logger := logger.FromContext(ctx).With("event", "subscription_deleted")
	logger.Info("handling subscription deleted webhook")

	data, ok := event["data"].(map[string]interface{})
	if !ok {
		return ErrWebhookProcessingFailed(nil, "invalid_data_format")
	}

	object, ok := data["object"].(map[string]interface{})
	if !ok {
		return ErrWebhookProcessingFailed(nil, "invalid_object_format")
	}

	subscriptionID, ok := object["id"].(string)
	if !ok {
		return ErrWebhookProcessingFailed(nil, "missing_subscription_id")
	}

	// Extract period information for final processing
	currentPeriodStart := int64(0)
	currentPeriodEnd := int64(0)
	if period, ok := object["current_period_start"].(float64); ok {
		currentPeriodStart = int64(period)
	}
	if period, ok := object["current_period_end"].(float64); ok {
		currentPeriodEnd = int64(period)
	}

	s.RefundAndFinalizeCancellation(ctx, subscriptionID, currentPeriodStart, currentPeriodEnd)
	logger.Info("subscription deleted webhook processed successfully")
	return nil
}

func (s *Service) handleInvoicePaymentSucceeded(ctx context.Context, event map[string]interface{}) error {
	logger := logger.FromContext(ctx).With("event", "invoice_payment_succeeded")
	logger.Info("handling invoice payment succeeded webhook")

	data, ok := event["data"].(map[string]interface{})
	if !ok {
		return ErrWebhookProcessingFailed(nil, "invalid_data_format")
	}

	object, ok := data["object"].(map[string]interface{})
	if !ok {
		return ErrWebhookProcessingFailed(nil, "invalid_object_format")
	}

	// Extract subscription ID if this is a subscription invoice
	if subscriptionID, ok := object["subscription"].(string); ok && subscriptionID != "" {
		s.MarkSubscriptionActive(ctx, subscriptionID)
		logger.Info("subscription marked as active after successful payment")

		// Send payment confirmation email notification
		if s.notification != nil {
			go func() {
				// Get subscription details to send notification
				subscription, err := s.storage.subscription.FindOne(context.Background(), s.storage.subscription.ScopeEquals(SubscriptionSchema.StripeSubID, subscriptionID))
				if err != nil {
					logger.Error("Failed to get subscription for payment succeeded email", "error", err, "subscription_id", subscriptionID)
					return
				}

				// Get plan details
				plan, err := s.GetPlanByID(context.Background(), subscription.PlanID)
				if err != nil {
					logger.Error("Failed to get plan for payment succeeded email", "error", err, "plan_id", subscription.PlanID)
					return
				}

				// Try to get payment method info and invoice details from event
				paymentMethodLastFour := "****"
				invoiceURL := ""

				// Extract payment method from event (simplified approach)
				if charge, ok := object["charge"].(map[string]interface{}); ok {
					if pm, ok := charge["payment_method_details"].(map[string]interface{}); ok {
						if card, ok := pm["card"].(map[string]interface{}); ok {
							if last4, ok := card["last4"].(string); ok {
								paymentMethodLastFour = last4
							}
						}
					}
				}

				// Get invoice URL from object
				if hostedInvoiceURL, ok := object["hosted_invoice_url"].(string); ok {
					invoiceURL = hostedInvoiceURL
				} else if invoicePDF, ok := object["invoice_pdf"].(string); ok {
					invoiceURL = invoicePDF
				}

				// Determine if this is first payment (subscription creation) or renewal
				// Check if subscription was recently created (within last 24 hours)
				isFirstPayment := time.Since(subscription.CreatedAt) < 24*time.Hour

				paymentDate := time.Now()
				if created, ok := object["created"].(float64); ok {
					paymentDate = time.Unix(int64(created), 0)
				}

				var emailErr error
				if isFirstPayment {
					// Send subscription confirmed email for first payment
					emailErr = s.notification.SendSubscriptionConfirmedEmail(context.Background(), subscription.WorkspaceID, subscription, plan, paymentMethodLastFour, invoiceURL)
				} else {
					// Send payment succeeded email for renewals
					emailErr = s.notification.SendPaymentSucceededEmail(context.Background(), subscription.WorkspaceID, subscription, plan, paymentMethodLastFour, invoiceURL, paymentDate)
				}

				if emailErr != nil {
					logger.Error("Failed to send payment confirmation email", "error", emailErr, "workspace_id", subscription.WorkspaceID, "is_first_payment", isFirstPayment)
				}
			}()
		}
	}

	logger.Info("invoice payment succeeded webhook processed successfully")
	return nil
}

func (s *Service) handleInvoicePaymentFailed(ctx context.Context, event map[string]interface{}) error {
	logger := logger.FromContext(ctx).With("event", "invoice_payment_failed")
	logger.Info("handling invoice payment failed webhook")

	data, ok := event["data"].(map[string]interface{})
	if !ok {
		return ErrWebhookProcessingFailed(nil, "invalid_data_format")
	}

	object, ok := data["object"].(map[string]interface{})
	if !ok {
		return ErrWebhookProcessingFailed(nil, "invalid_object_format")
	}

	// Extract subscription ID if this is a subscription invoice
	if subscriptionID, ok := object["subscription"].(string); ok && subscriptionID != "" {
		s.MarkSubscriptionPastDue(ctx, subscriptionID)
		logger.Info("subscription marked as past due after failed payment")

		// Send payment failed email notification
		if s.notification != nil {
			go func() {
				// Get subscription details to send notification
				subscription, err := s.storage.subscription.FindOne(context.Background(), s.storage.subscription.ScopeEquals(SubscriptionSchema.StripeSubID, subscriptionID))
				if err != nil {
					logger.Error("Failed to get subscription for payment failed email", "error", err, "subscription_id", subscriptionID)
					return
				}

				// Get plan details
				plan, err := s.GetPlanByID(context.Background(), subscription.PlanID)
				if err != nil {
					logger.Error("Failed to get plan for payment failed email", "error", err, "plan_id", subscription.PlanID)
					return
				}

				// Try to get payment method info (simplified - would need to extract from Stripe)
				paymentMethodLastFour := "****"

				err = s.notification.SendPaymentFailedEmail(context.Background(), subscription.WorkspaceID, subscription, plan, paymentMethodLastFour, time.Now(), nil)
				if err != nil {
					logger.Error("Failed to send payment failed email", "error", err, "workspace_id", subscription.WorkspaceID)
				}
			}()
		}
	}

	logger.Info("invoice payment failed webhook processed successfully")
	return nil
}

func (s *Service) handleTrialWillEnd(ctx context.Context, event map[string]interface{}) error {
	logger := logger.FromContext(ctx).With("event", "trial_will_end")
	logger.Info("handling trial will end webhook")

	// This is a good place to send notification emails or trigger
	// other business logic when a trial is about to end

	data, ok := event["data"].(map[string]interface{})
	if !ok {
		return ErrWebhookProcessingFailed(nil, "invalid_data_format")
	}

	object, ok := data["object"].(map[string]interface{})
	if !ok {
		return ErrWebhookProcessingFailed(nil, "invalid_object_format")
	}

	subscriptionID, ok := object["id"].(string)
	if !ok {
		return ErrWebhookProcessingFailed(nil, "missing_subscription_id")
	}

	// Send trial ending notification email
	if s.notification != nil {
		go func() {
			// Get subscription details
			subscription, err := s.storage.subscription.FindOne(context.Background(), s.storage.subscription.ScopeEquals(SubscriptionSchema.StripeSubID, subscriptionID))
			if err != nil {
				logger.Error("Failed to get subscription for trial ending email", "error", err, "subscription_id", subscriptionID)
				return
			}

			// Get plan details
			plan, err := s.GetPlanByID(context.Background(), subscription.PlanID)
			if err != nil {
				logger.Error("Failed to get plan for trial ending email", "error", err, "plan_id", subscription.PlanID)
				return
			}

			// Get trial status
			trialInfo, err := s.CheckTrialStatus(context.Background(), &account.Workspace{ID: subscription.WorkspaceID})
			if err != nil {
				logger.Error("Failed to get trial status for trial ending email", "error", err, "workspace_id", subscription.WorkspaceID)
				return
			}

			err = s.notification.SendTrialEndingEmail(context.Background(), subscription.WorkspaceID, subscription, plan, trialInfo)
			if err != nil {
				logger.Error("Failed to send trial ending email", "error", err, "workspace_id", subscription.WorkspaceID)
			}
		}()
	}

	logger.Info("trial will end webhook processed successfully", "subscription_id", subscriptionID)
	return nil
}
