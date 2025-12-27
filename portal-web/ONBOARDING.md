# Onboarding Flow Implementation

## Overview

The onboarding flow is a multi-step process that guides new users through:
1. **Plan Selection** - Choose a billing plan
2. **Email Verification** - Verify identity via OTP or Google OAuth
3. **Business Setup** - Configure business details
4. **Payment** - Complete Stripe checkout (paid plans only)
5. **Completion** - Finalize workspace creation and login

## Architecture

### State Management
- **OnboardingContext** - Manages flow state with sessionStorage persistence
- **URL-based navigation** - Each step is a separate route
- **Session token** - Backend tracks progress via token

### Backend Integration
- **API Endpoints**: `/api/onboarding/*`
- **Session-based** - All staged data stored in `onboarding_sessions` table
- **Atomic commit** - Final step commits everything in a transaction
- **Stripe integration** - Checkout sessions for paid plans

### Routes Structure

```
/onboarding
  ├── /plan                     - Step 1: Plan selection
  ├── /verify                   - Step 2: Email verification
  ├── /business                 - Step 3: Business setup
  ├── /payment                  - Step 4: Payment (paid plans)
  ├── /complete                 - Step 5: Finalization
  └── /oauth-callback           - Google OAuth callback handler
```

## Components

### OnboardingLayout
Simple, focused layout without sidebar:
- Progress indicator
- Language switcher
- Mobile-first responsive
- RTL support

### Route Guards
- **OnboardingRoot** - Redirects if already authenticated
- **Step-specific** - Each step validates prerequisites before rendering

## API Integration

### Flow Sequence

```typescript
// 1. Start session
POST /api/onboarding/start
→ { sessionToken, stage, isPaid }

// 2. Verify identity (Email)
POST /api/onboarding/email/otp
POST /api/onboarding/email/verify
→ { stage: "identity_verified" }

// 2. Verify identity (Google OAuth)
POST /api/onboarding/oauth/google
→ { stage: "identity_verified" }

// 3. Set business details
POST /api/onboarding/business
→ { stage: "business_staged" | "ready_to_commit" }

// 4. Payment (paid plans only)
POST /api/onboarding/payment/start
→ { checkoutUrl }
// User completes Stripe checkout
// Webhook updates stage to "ready_to_commit"

// 5. Complete onboarding
POST /api/onboarding/complete
→ { user, token, refreshToken }
```

### Error Handling
- All API errors are translated via `translateErrorAsync`
- Inline error messages in forms
- Retry buttons on critical failures
- Graceful degradation (e.g., payment cancellation)

## Stripe Integration

### Payment Flow
1. Backend creates Stripe Customer
2. Backend creates Checkout Session
3. User redirects to Stripe
4. Stripe processes payment
5. Webhook confirms payment
6. User returns to `/onboarding/payment?status=success`
7. Frontend polls status and proceeds to completion

### Success/Cancel Handling
- **Success**: Auto-redirect to complete step after 2s
- **Cancelled**: Option to retry or change plan

## i18n Support

### Translation Files
- `locales/en/onboarding.json`
- `locales/ar/onboarding.json`

### Key Namespaces
```typescript
t("onboarding.plan.title")
t("onboarding.verify.otpSent")
t("onboarding.business.nameRequired")
t("onboarding.payment.successTitle")
t("onboarding.complete.welcomeMessage", { businessName })
```

## Usage Example

### Linking to Onboarding
```typescript
// From login page
<Link to="/onboarding/plan?plan=starter">
  Sign up for Starter Plan
</Link>

// From marketing page
<Link to="/onboarding/plan?plan=professional&email=user@example.com">
  Get Started
</Link>
```

### Pre-selecting Plan
```typescript
// URL parameter automatically selects plan
/onboarding/plan?plan=starter
```

## State Persistence

### sessionStorage Schema
```typescript
interface OnboardingState {
  sessionToken: string | null;
  stage: SessionStage | null;
  selectedPlan: Plan | null;
  isPaidPlan: boolean;
  email: string | null;
  isEmailVerified: boolean;
  businessName: string | null;
  businessDescriptor: string | null;
  businessCountry: string | null;
  businessCurrency: string | null;
  isPaymentComplete: boolean;
  checkoutUrl: string | null;
}
```

### Session Recovery
- State persists across page refreshes
- Cleared on completion or logout
- Expires after 24 hours (backend)

## Testing Checklist

- [ ] Free plan flow (no payment)
- [ ] Paid plan flow with Stripe
- [ ] Email OTP verification
- [ ] Google OAuth sign-up
- [ ] Payment cancellation
- [ ] Payment success
- [ ] Invalid descriptor validation
- [ ] Country/currency selection
- [ ] Progress bar updates
- [ ] Mobile responsiveness
- [ ] RTL layout (Arabic)
- [ ] Session recovery after refresh
- [ ] Error handling and retries

## Future Enhancements

1. **Email verification link** as alternative to OTP
2. **Skip business setup** for quick trials
3. **Plan comparison modal** with detailed features
4. **Payment method preview** before Stripe redirect
5. **Onboarding analytics** for conversion tracking
6. **Resume incomplete sessions** via email link
7. **Multi-currency pricing** based on IP geolocation
