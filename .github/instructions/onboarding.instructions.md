---
description: "Kyora onboarding SSOT (backend + portal-web): stages, endpoint contracts, routing, polling, and do/don't rules"
applyTo: "backend/internal/domain/onboarding/**,portal-web/src/routes/onboarding/**,portal-web/src/api/onboarding.ts,portal-web/src/api/types/onboarding.ts,portal-web/src/lib/onboarding.ts"
---

# Kyora Onboarding SSOT (Backend + Portal Web)

This file is the **single source of truth** for how onboarding works end-to-end across:

- Backend: `backend/internal/domain/onboarding/**`
- Portal Web: `portal-web/src/routes/onboarding/**`, `portal-web/src/api/onboarding.ts`, `portal-web/src/lib/onboarding.ts`

It documents what is **actually implemented today**. If you change the flow, you must update both backend and portal-web and keep them consistent.

## Non-negotiables

- Backend is the **source of truth** for onboarding stage and payment state.
- Portal-web must treat `GET /v1/onboarding/session` as truth and be able to **resume** from any step.
- Never invent a new stage or “skip” behavior in portal-web without adding it to backend and updating:
  - Portal route gating + redirects
  - `portal-web/src/lib/onboarding.ts` stage→route mapping
  - Zod enums in `portal-web/src/api/types/onboarding.ts`

## Backend: state machine (authoritative)

Backend stage enum (see `backend/internal/domain/onboarding/model.go`):

- `plan_selected`
- `identity_pending`
- `identity_verified`
- `business_staged`
- `payment_pending` (paid plans)
- `payment_confirmed` (defined, but not currently set by backend logic)
- `ready_to_commit`
- `committed`

Backend payment status enum:

- `skipped` (free plan)
- `pending` (paid plan, awaiting Stripe)
- `succeeded` (Stripe confirmed via webhook)

### Canonical transitions (today)

- Start session: `plan_selected`
- Send OTP: `plan_selected` → `identity_pending`
- Verify OTP: `identity_pending` → `identity_verified`
- Google OAuth: (any non-expired session) → `identity_verified` (sets `method=google`, `emailVerified=true`)
- Set business:
  - free plan: `identity_verified|business_staged` → `ready_to_commit` and `paymentStatus=skipped`
  - paid plan: `identity_verified|business_staged` → `payment_pending` and `paymentStatus=pending`
- Stripe webhook: `payment_pending` → `ready_to_commit` and `paymentStatus=succeeded`
- Complete onboarding: `ready_to_commit` → `committed` (atomic commit)

Important: `payment_confirmed` exists in the enum, and portal-web supports it in routing, but backend currently sets `ready_to_commit` directly when payment succeeds.

## Backend: endpoints and contracts

All onboarding endpoints live under `/v1/onboarding`.

### `POST /v1/onboarding/start`

Request:

- `email` (required, email)
- `planDescriptor` (required)

Response:

- `sessionToken`
- `stage`
- `isPaid`

Behavior:

- Rejects if the email already belongs to an existing user (`409 Conflict`).
- If an **active session** already exists for the email, it is **resumed** (plan may be updated).
- Session expiry is **24 hours** from creation.

### `POST /v1/onboarding/email/otp`

Request:

- `sessionToken` (required)

Response:

- `retryAfterSeconds` (int)

Behavior:

- Allowed stages: `plan_selected` or `identity_pending`.
- Sets stage to `identity_pending`.
- OTP expiry is **15 minutes**.
- Cooldown is **2 minutes** per email (returns `429 Too Many Requests` with Problem JSON extensions: `retryAfterSeconds`).

### `POST /v1/onboarding/email/verify`

Request:

- `sessionToken`
- `code` (6 digits)
- `firstName`
- `lastName`
- `password` (min 8)

Response:

- `stage`

Behavior:

- Allowed stage: `identity_pending`.
- On success: `identity_verified`.
- On invalid/expired code: `400 Bad Request`.

### `POST /v1/onboarding/oauth/google`

Request:

- `sessionToken`
- `code` (OAuth authorization code)

Response:

- `stage`

Behavior:

- Sets `method=google`, `emailVerified=true`, names from Google profile, random password hash.
- Sets stage to `identity_verified`.

### `POST /v1/onboarding/business`

Request:

- `sessionToken`
- `name`
- `descriptor`
- `country` (len=2)
- `currency` (len=3)

Response:

- `stage`

Behavior:

- Allowed stages: `identity_verified` or `business_staged`.
- Free plan: sets `ready_to_commit`.
- Paid plan: sets `payment_pending`.

### `POST /v1/onboarding/payment/start`

Request:

- `sessionToken`
- `successUrl` (absolute URL)
- `cancelUrl` (absolute URL)

Response:

- `checkoutUrl` (string)

Behavior:

- For free plans: returns empty string (portal-web should not call it).
- Rate limit: max 3 attempts per 10 minutes and minimum 30 seconds apart (errors as `429` with generic detail).
- Reuses existing pending Stripe checkout session URL when possible.
- Uses Stripe idempotency key: `onboarding_checkout_<sessionID>`.

### `POST /v1/onboarding/complete`

Request:

- `sessionToken`

Response:

- `user`, `token`, `refreshToken`

Behavior:

- Allowed stage: `ready_to_commit` only.
- Performs an atomic commit:
  - bootstraps workspace + owner user
  - creates business
  - creates subscription for paid plan
  - marks session committed (`stage=committed`, `committedAt=now`)
- Issues JWT + refresh token.
- Sends welcome email best-effort.

### `GET /v1/onboarding/session?sessionToken=...`

Response includes:

- `sessionToken`, `email`, `stage`, `planId`, `planDescriptor`, `isPaidPlan`
- optional `firstName`, `lastName`
- optional `businessName`, `businessDescriptor`, `businessCountry`, `businessCurrency`
- `paymentStatus` (`skipped|pending|succeeded`)
- optional `checkoutSessionId`, `otpExpiry`
- `expiresAt`

Behavior:

- Rejects missing token with `400`.
- Rejects expired sessions (`400`).
- Treats committed sessions as expired (current code returns `ErrSessionExpired`).

### `DELETE /v1/onboarding/session`

Request:

- `sessionToken`

Response:

- `204 No Content`

Behavior:

- Allows cancel/restart.
- Rejects deletion if already committed.

## Portal Web: routes, redirects, and source-of-truth

Routes:

- `/onboarding` → redirects to `/onboarding/plan` (unless already authenticated, then `/`)
- `/onboarding/plan`
- `/onboarding/email?plan=<descriptor>`
- `/onboarding/verify?session=<token>`
- `/onboarding/oauth-callback?session=<token>&code=...`
- `/onboarding/business?session=<token>`
- `/onboarding/payment?session=<token>[&status=success|cancelled]`
- `/onboarding/complete?session=<token>`

### Stage→route mapping

`portal-web/src/lib/onboarding.ts` maps backend stages to a route:

- `plan_selected` → `/onboarding/email`
- `identity_pending` → `/onboarding/verify`
- `identity_verified` → `/onboarding/business`
- `business_staged` → `/onboarding/payment`
- `payment_pending` → `/onboarding/payment`
- `payment_confirmed` → `/onboarding/complete`
- `ready_to_commit` → `/onboarding/complete`
- `committed` → `/onboarding/complete`

Most onboarding route loaders call `redirectToCorrectStage(currentPath, session.stage, sessionToken)`.

### Polling behavior

`portal-web/src/api/onboarding.ts` session query options:

- `staleTime: 0`, `gcTime: 0`
- `refetchInterval: 3000`

This is required so the portal can observe payment status transitions driven by Stripe webhook → backend `ready_to_commit`.

## Portal Web: step behavior and required invariants

### Plan (`/onboarding/plan`)

- Loads plans via `GET /v1/billing/plans`.
- Selecting a plan navigates to `/onboarding/email?plan=<descriptor>`.

### Email (`/onboarding/email?plan=...`)

- Requires plan descriptor in search params.
- Starts session via `POST /v1/onboarding/start`.
- On success navigates to `/onboarding/verify?session=<token>`.
- "Continue with Google" uses `/v1/auth/google/url` (auth module), stores plan descriptor in session storage.

### Verify (`/onboarding/verify?session=...`)

- Loader fetches session and redirects to correct stage.
- Auto-sends OTP on mount when session stage is `plan_selected`.
- Resend uses backend `retryAfterSeconds` (from success body or Problem JSON extensions on 429).
- Two UI steps:
  - OTP entry step (collects code only)
  - profile step (first/last/password) then calls `POST /v1/onboarding/email/verify`

### OAuth callback (`/onboarding/oauth-callback`)

- Calls `POST /v1/onboarding/oauth/google`.
- On success navigates to business.
- On error navigates back to verify after a delay.

### Business (`/onboarding/business`)

- Loader fetches session and redirects to correct stage.
- Sets business via `POST /v1/onboarding/business`.
- Navigation decision:
  - if response.stage is `ready_to_commit` → `/onboarding/complete`
  - else if response.stage is `business_staged` → `/onboarding/payment`
  - else → `/onboarding/complete`

Note: backend currently returns `payment_pending` for paid plans (not `business_staged`). If you change backend stages, you must also change this client logic.

### Payment (`/onboarding/payment`)

- Route guard sends free plans to `/onboarding/complete`.
- Route guard currently only allows stages: `business_staged|payment_pending|payment_confirmed`.
  - Backend sets `payment_pending`, so `business_staged` is generally not expected.
- On `status=success`, it navigates to `/onboarding/complete`.
- Starts payment by calling `POST /v1/onboarding/payment/start` and then `window.location.href = checkoutUrl`.

### Complete (`/onboarding/complete`)

- Route guard currently allows only `ready_to_commit` or `payment_confirmed`.
- Calls `POST /v1/onboarding/complete`, stores tokens, sets user, and deletes session.
- Redirects to `/`.

## Consistency checklist (when changing onboarding)

When you change onboarding behavior, you must update all of these together:

- Backend:
  - stage transitions in `backend/internal/domain/onboarding/service.go`
  - DTOs in `backend/internal/domain/onboarding/handler_http.go`
  - stage/payment enums in `backend/internal/domain/onboarding/model.go`
- Portal:
  - Zod enums in `portal-web/src/api/types/onboarding.ts`
  - Stage routing in `portal-web/src/lib/onboarding.ts`
  - Route guards/loaders in `portal-web/src/routes/onboarding/*.tsx`
  - Any step-specific assumptions (e.g., `business.tsx` expecting `business_staged` vs `payment_pending`)

## Known mismatches / sharp edges (documented so you don’t reintroduce them)

- Backend sets `payment_pending` after business for paid plans; portal business step’s navigation checks `business_staged`.
- Backend does not set stage `payment_confirmed`; portal supports it in routing and guards.
- Backend `DELETE /v1/onboarding/session` uses JSON body and returns `400` if the session is committed.

Keep this file updated if you fix or intentionally change any of the above.
