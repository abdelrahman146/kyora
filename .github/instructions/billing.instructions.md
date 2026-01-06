---
description: Billing SSOT — Plans, Subscriptions, Invoices, Webhooks (Backend + Portal-Web guidance)
applyTo: "backend/internal/domain/billing/**,backend/internal/server/routes.go,backend/internal/domain/onboarding/service.go,backend/internal/domain/onboarding/handler_http.go,backend/internal/domain/onboarding/handler_bus.go,backend/internal/platform/bus/events.go,portal-web/src/api/onboarding.ts,portal-web/src/routes/onboarding/payment.tsx,portal-web/src/routes/onboarding/plan.tsx,portal-web/src/routes/billing/**,portal-web/src/api/billing.ts,portal-web/src/api/types/billing.ts"
---

# Billing — Single Source of Truth (SSOT)

This file documents **how billing works today** in Kyora, based on the current backend + portal-web implementation.

- Stripe details (SDK usage, keys, security fundamentals) live in: [.github/instructions/stripe.instructions.md](./stripe.instructions.md)
- This file owns **Kyora’s billing workflow**, endpoint contracts, webhook semantics, and the onboarding↔billing bridge.

## Scope and Mental Model

- **Billing is workspace-scoped**: subscriptions, Stripe customer, invoices, and plan enforcement all attach to the workspace.
- **Billing “gates” product actions** via middleware in `backend/internal/domain/billing/middleware_http.go`.
- **Onboarding paid-plan checkout is a special flow** implemented in the onboarding domain (it creates a Stripe Checkout session and relies on the billing webhook handler to confirm payment).

## Backend HTTP Surface (authoritative)

Routes are wired in `backend/internal/server/routes.go`.

### Public (no auth)

- `GET /v1/billing/plans` — list plans
- `GET /v1/billing/plans/:descriptor` — get plan by descriptor
- `POST /webhooks/stripe` — Stripe webhooks (public; signature verified)

### Protected (auth + workspace membership + RBAC)

All routes below require:

- `auth.EnforceAuthentication`
- `account.EnforceValidActor`
- `account.EnforceWorkspaceMembership`
- plus `account.EnforceActorPermissions(..., role.ResourceBilling)` per endpoint

Key groups:

- **Subscription**: `GET/POST/DELETE /v1/billing/subscription` + detail/estimate/schedule/resume + trial/grace
- **Payment methods**: `POST /v1/billing/payment-methods/setup-intent`, `POST /v1/billing/payment-methods/attach`
- **Invoices**: `GET /v1/billing/invoices`, `GET /v1/billing/invoices/:id/download`, `POST /v1/billing/invoices/:id/pay`, `POST /v1/billing/invoices`
- **Checkout**: `POST /v1/billing/checkout/session`
- **Billing portal**: `POST /v1/billing/portal/session`
- **Usage/tax**: `GET /v1/billing/usage`, `GET /v1/billing/usage/quota`, `POST /v1/billing/tax/calculate`

Implementation reference: `backend/internal/domain/billing/handler_http.go`.

## Plan Enforcement (middleware)

Billing enforcement middleware lives in `backend/internal/domain/billing/middleware_http.go`.

### `EnforceActiveSubscription`

- Loads the workspace subscription and attaches it to request context.
- Calls `subscription.IsActive()` which currently means `status == "active"` only.
- Intended usage pattern (example in routing): apply it **before** plan limit checks (e.g. inviting users, creating businesses).

### `EnforcePlanWorkspaceLimits` / `EnforcePlanBusinessLimits`

- Reads the subscription from context (requires `EnforceActiveSubscription` earlier).
- Computes usage via injected function (DB count) and enforces plan limit via `sub.Plan.Limits.CheckUsageLimit(feature, usage)`.

### `EnforcePlanFeatureRestriction`

- Feature flag gate: `sub.Plan.Features.CanUseFeature(feature)`.

## Webhooks (security + idempotency + side effects)

Webhook handler entry point: `POST /webhooks/stripe` -> `billing.HttpHandler.HandleWebhook` -> `billing.Service.ProcessWebhook`.

### Signature verification

- Signature verification is implemented manually in `backend/internal/domain/billing/webhooks.go` (`verifyStripeSignature`).
- It parses `Stripe-Signature` header (`t=...`, `v1=...`) and checks an HMAC SHA-256 signature with a tolerance (currently 5 minutes).

### Idempotency

- Stripe event IDs are persisted to `stripe_events` (see `backend/internal/domain/billing/stripe_events.go`).
- If an event ID was already stored, processing is skipped.

### Event types handled

The backend currently handles (and ignores unhandled event types):

- `customer.subscription.created/updated/deleted`
- `invoice.payment_succeeded`, `invoice.payment_failed`, `invoice.finalized`, `invoice.marked_uncollectible`, `invoice.voided`
- `customer.subscription.trial_will_end`
- `payment_method.automatically_updated`
- `checkout.session.completed`

## Onboarding ↔ Billing bridge (critical)

Paid-plan onboarding is implemented in the onboarding domain:

- Backend endpoint: `POST /v1/onboarding/payment/start` in `backend/internal/domain/onboarding/handler_http.go`
- It creates (or reuses) a Stripe Checkout Session in `backend/internal/domain/onboarding/service.go`.
- The Stripe customer + checkout session include metadata:
  - `onboarding_session_id`
  - `plan_id`

**Payment confirmation is webhook-driven**:

- When Stripe emits `checkout.session.completed`, the billing webhook handler:
  - fetches the Checkout Session from Stripe (server-side)
  - best-effort attaches / sets default payment method
  - if metadata contains `onboarding_session_id`, emits `bus.OnboardingPaymentSucceededTopic` with:
    - `OnboardingSessionID`
    - `StripeCheckoutID`
    - `StripeSubscriptionID`
- The onboarding domain listens to this bus topic in `backend/internal/domain/onboarding/handler_bus.go` and marks the onboarding session as payment succeeded.

**Implication**: frontend should treat “paid-plan onboarding payment success” as **asynchronous** and confirmed only when the onboarding session state changes.

## Portal-web (current reality + guidance)

### What exists today

- Plan listing during onboarding uses:

  - `GET /v1/billing/plans`
  - `GET /v1/billing/plans/:descriptor`
  - implemented in `portal-web/src/api/onboarding.ts`

- Paid-plan payment step in onboarding is implemented in:
  - `portal-web/src/routes/onboarding/payment.tsx`
  - It calls `POST /v1/onboarding/payment/start` to get `checkoutUrl`, then redirects the browser to Stripe.
  - It polls `GET /v1/onboarding/session` every 3 seconds (React Query `refetchInterval`) to observe server-confirmed payment progression.

### Rules for implementing billing UI later

If/when building a real `/billing` area in portal-web:

- Prefer backend as the source of truth: use the billing endpoints in `backend/internal/server/routes.go` and the DTOs in `backend/internal/domain/billing/model.go` / `handler_http.go`.
- Never assume payment success from a frontend redirect alone. Always fetch:
  - subscription status via `GET /v1/billing/subscription` (workspace billing)
  - or onboarding session via `GET /v1/onboarding/session` (onboarding)
- Treat webhook-confirmed state transitions as eventual-consistency events.

## Don’ts (SSOT / correctness)

- Don’t duplicate Stripe fundamentals here; link to [.github/instructions/stripe.instructions.md](./stripe.instructions.md).
- Don’t hardcode a portal-web billing route structure that doesn’t exist yet.
- Don’t add new webhook event types without:
  - signature verification,
  - idempotency protection,
  - and a clearly defined side effect (DB write or bus event).
