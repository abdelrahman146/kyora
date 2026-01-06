---
description: Stripe Integration — Billing & Payments for Kyora
applyTo: "portal-web/**,backend/**"
---

# Stripe Integration

**Purpose**: Stripe fundamentals for Kyora billing and payments  
**Scope**: Backend (Stripe SDK integration) + Frontend (Stripe-hosted redirect flows / Stripe.js when needed)  
**Kyora workflow SSOT**: [.github/instructions/billing.instructions.md](./billing.instructions.md) (endpoints, webhook semantics, onboarding↔billing bridge)
**Architecture reference**: [backend-core.instructions.md](./backend-core.instructions.md)

## Backend Integration

**Config**:

- `billing.stripe.api_key` — Stripe secret key
- `billing.stripe.webhook_secret` — Webhook signing secret

**Initialization** (`internal/server/server.go`):

```go
stripe.Key = cfg.GetString(config.StripeAPIKey)
```

**Customer Management**:

- Create Stripe customer on workspace creation
- Store `stripeCustomerID` in `Workspace` model
- Link payment methods to customer

**Subscription Lifecycle**:

1. Create subscription with plan price ID
2. Handle payment intent if required (3D Secure, etc.)
3. Store `stripeSubscriptionID` in `Workspace`
4. Update subscription status via webhooks

**Webhooks**:

- Endpoint: `POST /webhooks/stripe` (see [.github/instructions/billing.instructions.md](./billing.instructions.md) for the authoritative route map)
- Signature verification is implemented in Kyora (HMAC verification of `Stripe-Signature` header with timestamp tolerance)
- Webhooks must be idempotent (Kyora persists processed Stripe event IDs)

**Error Handling**:

- Use `stripe.Error` type checking
- Return proper RFC 7807 problems
- Log stripe request IDs for debugging

## Frontend Integration (Portal Web)

**API Client**: Ky HTTP client (see [ky.instructions.md](./ky.instructions.md))  
**Billing UI**: not fully implemented yet; onboarding includes a paid-plan Checkout redirect flow.

**Key Flows**:

1. **Display Plans**: Fetch available plans from backend
2. **Select Plan**: Show payment form if upgrade required
3. **Payment Method**: Stripe Payment Element or saved payment methods
4. **Subscription Status**: Display current plan, next billing date, invoices

**Stripe Elements**:

- Load Stripe.js: `<script src="https://js.stripe.com/v3/"></script>`
- Initialize: `const stripe = Stripe(publishableKey)`
- Payment Element for new payment methods
- Display confirmation after successful payment

Note: If Kyora uses Stripe-hosted flows (Checkout / Billing Portal), portal-web should prefer backend-created session URLs and redirect the browser, then refresh server state (don’t assume success from redirects alone).

## Stripe API Patterns

**Use Payment Intents API** (not Charges):

- Create PaymentIntent on backend
- Return `clientSecret` to frontend
- Frontend confirms payment with Stripe.js
- Webhook processes result

**Customer Management**:

- One Stripe customer per workspace
- Attach payment methods to customer
- Set default payment method for subscriptions

**Subscription Management**:

- Use `billing_cycle_anchor` for consistent billing dates
- Handle prorations automatically with `proration_behavior: 'create_prorations'`
- Cancel at period end vs immediate cancellation

**Metadata**:

- Store workspace ID in Stripe metadata: `metadata: { workspaceId: "wks_xxx" }`
- Store user ID when relevant
- Use for webhook event processing

## Testing

**Test Cards** (development):

- Success: `4242 4242 4242 4242`
- 3D Secure: `4000 0025 0000 3155`
- Decline: `4000 0000 0000 0002`
- Insufficient funds: `4000 0000 0000 9995`

**Test Webhooks**:

- Use Stripe CLI: `stripe listen --forward-to localhost:8080/v1/webhooks/stripe`
- Trigger events: `stripe trigger customer.subscription.created`

**Mock Mode** (E2E tests):

- Use testcontainers stripe-mock
- Reference: `backend/internal/tests/e2e/main_test.go`

## Security Best Practices

**API Keys**:

- Use secret key on backend only (never frontend)
- Use publishable key on frontend
- Rotate keys if compromised

**Webhook Security**:

- ALWAYS verify webhook signatures
- Use rolling webhook secrets
- Reject unsigned events

**Payment Methods**:

- Use Payment Intents API (PCI compliant)
- Never handle raw card numbers on backend
- Let Stripe.js tokenize payment details

**Idempotency**:

- Use idempotency keys for create operations: `stripe.WithIdempotencyKey(key)`
- Prevent duplicate charges on retries

## Common Patterns

**Check Subscription Status**:

```go
if workspace.SubscriptionStatus != "active" {
    return problem.Forbidden("subscription required")
}
```

**Enforce Plan Limits**:

```go
// middleware_http.go pattern
func EnforcePlanWorkspaceLimits(limit int, counterFunc func) gin.HandlerFunc {
    return func(c *gin.Context) {
        workspace := account.WorkspaceFromContext(c)
        count := counterFunc(workspace.ID)
        if count >= limit {
            response.Error(c, problem.Forbidden("plan limit reached"))
            c.Abort()
            return
        }
        c.Next()
    }
}
```

**Create Subscription**:

```go
params := &stripe.SubscriptionParams{
    Customer: stripe.String(workspace.StripeCustomerID),
    Items: []*stripe.SubscriptionItemsParams{
        {Price: stripe.String(priceID)},
    },
    Metadata: map[string]string{
        "workspaceId": workspace.ID,
    },
}
sub, err := subscription.New(params)
```

**Process Webhook**:

- Verify `Stripe-Signature`
- Enforce timestamp tolerance
- Enforce idempotency by Stripe event ID
- Dispatch by `event.type`
- Persist side effects (DB updates) before acknowledging

## Documentation References

**Stripe Docs** (external):

- [Payment Intents API](https://docs.stripe.com/api/payment_intents)
- [Subscriptions API](https://docs.stripe.com/api/subscriptions)
- [Webhooks Guide](https://docs.stripe.com/webhooks)
- [Testing](https://docs.stripe.com/testing)
- [Error Codes](https://docs.stripe.com/error-codes)

**Kyora Internal**:

- [backend-core.instructions.md](./backend-core.instructions.md) — Architecture, error handling, HTTP patterns
- [backend-testing.instructions.md](./backend-testing.instructions.md) — Testing Stripe integration
- [ky.instructions.md](./ky.instructions.md) — Frontend HTTP client patterns

## Troubleshooting

**Webhook not receiving events**:

- Check webhook secret in config
- Verify signature verification logic
- Test with Stripe CLI
- Check Stripe Dashboard webhook logs

**Payment fails**:

- Check Stripe Dashboard for decline reason
- Verify payment method is valid
- Check for 3D Secure requirements
- Review error logs with Stripe request ID

**Subscription status mismatch**:

- Check webhook processing logs
- Verify database updates are atomic
- Manually sync from Stripe Dashboard if needed

## Anti-Patterns

❌ Never use Charges API (deprecated) — use Payment Intents  
❌ Never skip webhook signature verification  
❌ Never store raw card numbers  
❌ Never hardcode API keys  
❌ Never process webhooks without idempotency checks  
❌ Never update subscription without checking current status
