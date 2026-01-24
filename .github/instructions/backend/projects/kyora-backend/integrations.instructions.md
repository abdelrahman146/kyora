---
description: Kyora Backend Integrations — Resend Email, Stripe Billing, Blob Storage
applyTo: "backend/**"
---

# Kyora Backend Integrations

**External service integrations for Kyora backend.**

Use when: Working with email, payments, file storage in Kyora.

See also:

- `domain-modules.instructions.md` — Kyora domain overview
- `../general/architecture.instructions.md` — General patterns
- `.github/instructions/resend.instructions.md` — Full Resend docs
- `.github/instructions/stripe.instructions.md` — Full Stripe docs

---

## Integration Overview

Kyora backend integrates with:

| Service      | Purpose                | Config Prefix     | Provider Options            |
| ------------ | ---------------------- | ----------------- | --------------------------- |
| Resend       | Transactional emails   | `email.`          | `resend`, `mock`            |
| Stripe       | Billing, subscriptions | `billing.stripe.` | Real, mock (testcontainers) |
| Blob Storage | Asset uploads          | `storage.`        | `local`, `s3`               |
| Memcached    | Caching                | `cache.`          | Real, testcontainers        |
| PostgreSQL   | Database               | `database.`       | Real, testcontainers        |

---

## Resend Email Integration

### Configuration

```go
// Config keys
const (
    EmailProvider    = "email.provider"        // "resend" or "mock"
    ResendAPIKey     = "email.resend.api_key"
    ResendBaseURL    = "email.resend.base_url"
    EmailFromEmail   = "email.from_email"
    EmailFromName    = "email.from_name"
    EmailSupportEmail = "email.support_email"
)
```

### Initialization

```go
// In server.New()
emailClient := email.NewClient(
    cfg.GetString(config.EmailProvider),
    cfg.GetString(config.ResendAPIKey),
    cfg.GetString(config.ResendBaseURL),
)
```

### Available Templates

Located in `backend/internal/platform/email/templates/`:

| Template Constant               | Use Case               | Required Data                                    |
| ------------------------------- | ---------------------- | ------------------------------------------------ |
| `TemplateEmailVerification`     | Verify email on signup | `Name`, `VerificationLink`                       |
| `TemplateForgotPassword`        | Password reset         | `Name`, `ResetLink`                              |
| `TemplateWelcome`               | Welcome after signup   | `Name`, `WorkspaceName`                          |
| `TemplateWorkspaceInvitation`   | Workspace invite       | `InviterName`, `WorkspaceName`, `InvitationLink` |
| `TemplateSubscriptionConfirmed` | Subscription activated | `Name`, `PlanName`                               |
| `TemplateSubscriptionCancelled` | Subscription ended     | `Name`, `PlanName`                               |
| `TemplatePaymentFailed`         | Payment failure        | `Name`, `RetryLink`                              |

### Sending Emails

**Template-based** (recommended):

```go
err := emailClient.SendTemplate(ctx,
    email.TemplateEmailVerification,
    user.Email,
    email.From{
        Email: cfg.GetString(config.EmailFromEmail),
        Name:  cfg.GetString(config.EmailFromName),
    },
    "Verify your email - Kyora",
    map[string]interface{}{
        "Name": user.FirstName,
        "VerificationLink": verificationURL,
        "SupportEmail": cfg.GetString(config.EmailSupportEmail),
    },
)
```

**Event-driven** (via bus):

```go
// Emit event
bus.Emit(bus.VerifyEmailTopic, map[string]interface{}{
    "email": user.Email,
    "token": verificationToken,
})

// Listen in server setup
bus.Listen(bus.VerifyEmailTopic, func(payload any) {
    data := payload.(map[string]interface{})
    emailClient.SendTemplate(ctx, email.TemplateEmailVerification, ...)
})
```

### Mock Provider (Dev/Test)

Set `email.provider=mock`:

- Emails logged to console (not sent)
- No API calls made
- Useful for E2E tests without real email service

**Full docs**: `.github/instructions/resend.instructions.md`

---

## Stripe Billing Integration

### Configuration

```go
const (
    StripeAPIKey     = "billing.stripe.api_key"
    StripeWebhookSecret = "billing.stripe.webhook_secret"
    AutoSyncPlans    = "billing.auto_sync_plans" // Default: true
)
```

### Initialization

```go
// In server.New()
stripe.Key = cfg.GetString(config.StripeAPIKey)

// Auto-sync plans from Stripe on startup
if cfg.GetBool(config.AutoSyncPlans) {
    billingService.SyncPlansFromStripe(context.Background())
}
```

### Plan Management

**Plans stored in Kyora DB:**

- Synced from Stripe on startup (`billing.auto_sync_plans=true`)
- Manual sync: `kyora sync_plans`

**Plan structure:**

```go
type Plan struct {
    ID              string          `gorm:"primaryKey"`
    Name            string
    StripePriceID   string          // Stripe price ID
    Amount          decimal.Decimal
    Currency        string
    Features        types.JSONB     // Plan features/limits
}
```

**Plan limits enforced by middleware:**

```go
// In routes.go
businesses.Use(
    billing.EnforcePlanWorkspaceLimits(
        businessSvc,
        func(w *account.Workspace) int { return w.BusinessCount },
        func(p *billing.Plan) int { return p.Features.MaxBusinesses },
    ),
)
```

### Subscription Lifecycle

**Create subscription:**

```go
// In billing service
params := &stripe.SubscriptionParams{
    Customer: stripe.String(workspace.StripeCustomerID),
    Items: []*stripe.SubscriptionItemsParams{
        {Price: stripe.String(plan.StripePriceID)},
    },
    Metadata: map[string]string{
        "workspaceId": workspace.ID,
    },
}
sub, err := subscription.New(params)
```

**Update workspace:**

```go
workspace.StripeSubscriptionID = sub.ID
workspace.SubscriptionStatus = string(sub.Status)
workspace.PlanID = plan.ID
```

### Webhook Handling

**Endpoint:** `POST /webhooks/stripe` (public, signature verified)

**Handled events:**

- `customer.subscription.created` → Activate subscription
- `customer.subscription.updated` → Sync status (active, past_due, canceled)
- `customer.subscription.deleted` → Downgrade to free plan
- `invoice.payment_succeeded` → Record payment
- `invoice.payment_failed` → Notify user

**Signature verification:**

```go
event, err := webhook.ConstructEvent(
    bodyBytes,
    c.GetHeader("Stripe-Signature"),
    webhookSecret,
)
```

**Idempotency:**

- Store processed Stripe event IDs in cache
- Skip duplicate events

**Full docs**: `.github/instructions/stripe.instructions.md`

---

## Blob Storage Integration

### Configuration

```go
const (
    StorageProvider     = "storage.provider"      // "local" or "s3"
    StorageLocalPath    = "storage.local.path"    // "./tmp/assets"
    StorageS3Bucket     = "storage.s3.bucket"
    StorageS3Region     = "storage.s3.region"
    StorageS3Endpoint   = "storage.s3.endpoint"   // Optional (for S3-compatible)
)
```

### Initialization

```go
// In server.New()
blobClient := blob.FromConfig(cfg)
```

### Upload Flow

**1. Request presigned URL:**

```go
// Handler: POST /v1/assets/upload
func (h *Handler) RequestUpload(c *gin.Context) {
    var req RequestUploadRequest
    request.ValidBody(c, &req)

    uploadURL, assetID, err := h.blobClient.GenerateUploadURL(
        ctx,
        req.Filename,
        req.ContentType,
        req.Size,
    )

    response.SuccessJSON(c, http.StatusOK, map[string]interface{}{
        "uploadUrl": uploadURL,
        "assetId":   assetID,
    })
}
```

**2. Client uploads to URL:**

```bash
curl -X PUT "$uploadUrl" \
  -H "Content-Type: image/jpeg" \
  --upload-file photo.jpg
```

**3. Confirm upload:**

```go
// Handler: POST /v1/assets/:assetId/confirm
func (h *Handler) ConfirmUpload(c *gin.Context) {
    assetID := c.Param("assetId")

    asset, err := h.service.ConfirmAssetUpload(ctx, assetID)

    response.SuccessJSON(c, http.StatusOK, ToAssetResponse(asset))
}
```

### Local Provider (Dev/Test)

- Files stored in `./tmp/assets/`
- Direct file paths (no presigned URLs)
- Useful for local development

### S3 Provider (Production)

- Uses AWS S3 SDK
- Supports S3-compatible services (DigitalOcean Spaces, MinIO)
- Presigned URLs for direct client uploads

**Full docs**: `.github/instructions/asset_upload.instructions.md`

---

## Event Bus Integration

Kyora uses internal event bus for cross-domain automation.

### Bus Initialization

```go
// In server.New()
bus := bus.New()

// Subscribe domain handlers
accounting.NewBusHandler(bus, accountingSvc, businessSvc)
```

### Key Events

| Topic                        | Payload                      | Subscribers                         |
| ---------------------------- | ---------------------------- | ----------------------------------- |
| `OrderPaymentSucceededTopic` | `OrderPaymentSucceededEvent` | Accounting (upsert transaction fee) |
| `OrderCreatedTopic`          | `OrderCreatedEvent`          | Analytics (update metrics)          |
| `CustomerCreatedTopic`       | `CustomerCreatedEvent`       | Analytics (track acquisition)       |

### Publishing Events

```go
// In order service
s.bus.Emit(bus.OrderPaymentSucceededTopic, &bus.OrderPaymentSucceededEvent{
    Ctx:           ctx,
    OrderID:       order.ID,
    BusinessID:    order.BusinessID,
    PaymentMethod: string(order.PaymentMethod),
    Total:         order.Total,
})
```

### Subscribing to Events

```go
// In accounting/handler_bus.go
func NewBusHandler(b *bus.Bus, svc *Service, bizSvc accountingRequiredBusinessService) {
    h := &BusHandler{svc: svc, bizSvc: bizSvc}
    b.Listen(bus.OrderPaymentSucceededTopic, h.HandleOrderPaymentSucceeded)
}

func (h *BusHandler) HandleOrderPaymentSucceeded(payload interface{}) {
    event := payload.(*bus.OrderPaymentSucceededEvent)

    // Upsert transaction fee expense
    h.svc.UpsertTransactionFeeExpense(event.Ctx, event.BusinessID, event.OrderID, event.Total)
}
```

**Rules:**

- Handlers must be idempotent (events may be redelivered)
- Handler panics caught and logged (don't crash app)
- Events dispatched asynchronously (non-blocking)

---

## Testing Integrations

### Email Testing

```go
// Use mock provider in tests
cfg := &config.Config{EmailProvider: "mock"}
emailClient := email.NewClient("mock", "", "")

// Emails logged, not sent
emailClient.SendTemplate(ctx, email.TemplateWelcome, ...)
```

### Stripe Testing

**Test cards:**

- Success: `4242 4242 4242 4242`
- 3D Secure: `4000 0025 0000 3155`
- Decline: `4000 0000 0000 0002`

**Mock server (testcontainers):**

```go
// In E2E tests
stripeMock := testutils.CreateStripeMockCtx(ctx)
defer stripeMock.Terminate(ctx)

// Use mock endpoint
stripe.Key = "sk_test_mock"
stripe.API = stripeMock.Endpoint
```

**Webhook testing:**

```bash
# Use Stripe CLI
stripe listen --forward-to localhost:8080/webhooks/stripe
stripe trigger customer.subscription.created
```

### Blob Storage Testing

```go
// Use local provider in tests
cfg := &config.Config{
    StorageProvider: "local",
    StorageLocalPath: "./tmp/test_assets",
}
blobClient := blob.FromConfig(cfg)

// Files stored locally, no S3 calls
```

---

## Anti-Patterns

❌ Sending emails synchronously in HTTP handlers (use event bus)  
❌ Exposing Stripe API keys in frontend  
❌ Skipping Stripe webhook signature verification  
❌ Storing raw card numbers (use Stripe Payment Intents)  
❌ Processing webhooks without idempotency checks  
❌ Using production services in tests (use mocks)  
❌ Hardcoding email content (use templates)  
❌ Direct S3 uploads from backend (use presigned URLs)

---

## Quick Reference

**Resend**: Transactional emails, templates, mock provider for tests  
**Stripe**: Billing, subscriptions, webhooks, plan limits, test cards  
**Blob Storage**: Presigned uploads, local/S3 providers  
**Event Bus**: Cross-domain automation, idempotent handlers  
**Testing**: Mock providers, testcontainers, Stripe CLI
