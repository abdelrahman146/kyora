---
description: Kyora Backend Domain Modules — Account, Inventory, Orders, Customers, Accounting, Analytics
applyTo: "backend/**"
---

# Kyora Backend Domain Modules

**Kyora-specific domain module overview and integration patterns.**

Use when: Working with Kyora backend domains, understanding domain interactions.

See also:

- `../general/architecture.instructions.md` — General backend architecture
- `integrations.instructions.md` — External service integrations
- `testing-specifics.instructions.md` — Kyora testing patterns
- Project-specific SSOT: `.github/instructions/*-management.instructions.md` files

---

## Domain Map

Kyora backend includes these domain modules:

| Domain       | Purpose                                | Key Models                                    |
| ------------ | -------------------------------------- | --------------------------------------------- |
| `account`    | Users, workspaces, sessions, RBAC      | User, Workspace, Session, Invitation          |
| `business`   | Business profiles, descriptors, zones  | Business, ShippingZone, PaymentMethod         |
| `inventory`  | Products, variants, categories, stock  | Product, Variant, Category                    |
| `order`      | Orders, items, notes, status tracking  | Order, OrderItem, OrderNote                   |
| `customer`   | Customer profiles, addresses, notes    | Customer, Address, CustomerNote               |
| `accounting` | Assets, expenses, withdrawals, summary | Asset, Expense, Withdrawal, AccountingSummary |
| `analytics`  | Dashboards, reports, metrics           | (aggregation queries)                         |
| `billing`    | Plans, subscriptions, invoices         | Plan, Subscription, Invoice                   |
| `onboarding` | Multi-stage onboarding flow            | OnboardingSession                             |
| `storefront` | Public storefront (separate access)    | (uses business data)                          |
| `metadata`   | Categories, tags, custom fields        | Category, Tag                                 |
| `asset`      | File uploads, blob storage             | Asset                                         |

---

## Multi-Tenancy Hierarchy

**Two-level isolation:**

1. **Workspace**: Top-level tenant (one per registered account)
2. **Business**: Second-level tenant (multiple per workspace)

**Scoping rules:**

- Workspace-scoped: Users, sessions, invitations, subscriptions
- Business-scoped: Products, orders, customers, expenses, analytics

**Critical invariants:**

- User can only access their own workspace
- Business must belong to user's workspace
- All business-scoped queries MUST filter by `business_id`
- Never trust workspace/business IDs from URL params for authorization

---

## Account Domain

**Purpose**: Authentication, authorization, workspace/user management

**Models:**

- `User`: Email, password hash, auth version, workspace membership
- `Workspace`: Top-level tenant, Stripe customer ID, subscription
- `Session`: Refresh tokens (hashed), expiration tracking
- `Invitation`: Workspace invites, email-based registration
- `Role`: RBAC permissions per user

**Key flows:**

- Login/logout → JWT access token + rotating refresh token
- Token refresh → rotate session, return new tokens
- Logout all → increment `auth_version` (invalidates all tokens)
- Invite user → send email, create pending invitation
- Accept invite → create user, join workspace

**Middleware:**

- `EnforceAuthentication` → validate JWT
- `EnforceValidActor` → load user, check auth version
- `EnforceWorkspaceMembership` → load workspace (never trust URL)
- `EnforceActorPermissions` → RBAC checks

**SSOT**: `.github/instructions/account-management.instructions.md`

---

## Business Domain

**Purpose**: Business profiles, descriptors, payment/shipping config

**Models:**

- `Business`: Profile, descriptor (unique slug), workspace ownership
- `ShippingZone`: Delivery areas, pricing
- `PaymentMethod`: Accepted payment types

**Key rules:**

- Business descriptor is unique per workspace (not globally)
- Archived businesses hidden from UI but not deleted
- Business-scoped middleware loads business via descriptor
- Business must belong to actor's workspace

**Middleware:**

- `EnforceBusinessValidity` → load business by descriptor, verify workspace ownership

**SSOT**: `.github/instructions/business-management.instructions.md`

---

## Inventory Domain

**Purpose**: Products, variants, categories, stock tracking

**Models:**

- `Product`: Name, description, category, images
- `Variant`: SKU, price, cost, stock quantity, stock alert
- `Category`: Hierarchical product organization

**Key rules:**

- Products can have multiple variants (size, color, etc.)
- Stock tracked at variant level
- Low stock alerts when `stock_quantity <= stock_alert`
- Search uses PostgreSQL full-text search (`search_vector`)

**Stock semantics:**

- `stock_quantity`: Current available stock
- `stock_alert`: Threshold for low stock warning
- Track stock? → variant flag for inventory management

**SSOT**: `.github/instructions/inventory.instructions.md`

---

## Order Domain

**Purpose**: Order lifecycle, items, payments, status tracking

**Models:**

- `Order`: Business-scoped, customer link, totals, status, payment
- `OrderItem`: Product/variant reference, quantity, price snapshot
- `OrderNote`: Internal notes, timeline tracking

**Status machine:**

- `pending` → `confirmed` → `processing` → `shipped` → `delivered`
- `cancelled` at any stage
- Transitions validated in service layer

**Payment states:**

- `unpaid`, `partial`, `paid`, `refunded`
- Independent from order status

**Inventory integration:**

- Creating order reserves stock (decrements variant stock)
- Cancelling order releases stock (increments variant stock)

**SSOT**: `.github/instructions/orders.instructions.md`

---

## Customer Domain

**Purpose**: Customer profiles, addresses, purchase history

**Models:**

- `Customer`: Name, email, phone, social handles
- `Address`: Shipping/billing addresses (multiple per customer)
- `CustomerNote`: Timeline notes, interaction tracking

**Key features:**

- Link orders to customers
- Track lifetime value (calculated)
- Social platform tracking (Instagram, WhatsApp, TikTok)
- Search by name, email, phone

**SSOT**: `.github/instructions/customer.instructions.md`

---

## Accounting Domain

**Purpose**: Financial tracking, expenses, withdrawals, summary

**Models:**

- `Asset`: Money in business (initial capital, investments)
- `Expense`: Money out (purchases, fees, salaries)
- `Withdrawal`: Owner withdrawals
- `AccountingSummary`: Cached aggregates (profit, cash in hand)

**Automation (via event bus):**

- Order payment succeeded → upsert transaction fee expense
- Auto-categorize Stripe fees as "Transaction Fees"

**Key calculations:**

- Profit = Revenue - Expenses - COGS
- Cash in hand = Assets + Revenue - Expenses - Withdrawals

**SSOT**: `.github/instructions/accounting.instructions.md`

---

## Analytics Domain

**Purpose**: Dashboards, sales reports, inventory reports, customer reports

**Key features:**

- Sales over time (revenue, orders, avg order value)
- Inventory analytics (stock levels, best sellers, low stock alerts)
- Customer analytics (top customers, lifetime value, acquisition)
- Date range filters (today, week, month, year, custom)

**Implementation pattern:**

- Use raw SQL aggregations for performance
- Cache frequently accessed metrics
- Use PostgreSQL window functions for trends

**SSOT**: `.github/instructions/analytics.instructions.md`

---

## Billing Domain

**Purpose**: Stripe integration, plans, subscriptions, invoices

**Models:**

- `Plan`: Available subscription plans (Free, Starter, Professional)
- `Subscription`: Active workspace subscription, Stripe ID
- `Invoice`: Billing history

**Plan gates:**

- Workspace limits (users, businesses, products, orders)
- Feature flags (advanced analytics, API access)
- Middleware enforces plan limits before mutations

**Webhook events:**

- `customer.subscription.created` → activate subscription
- `customer.subscription.updated` → sync status
- `customer.subscription.deleted` → downgrade to free
- `invoice.payment_succeeded` → record payment

**SSOT**: `.github/instructions/billing.instructions.md`

---

## Onboarding Domain

**Purpose**: Multi-stage onboarding flow for new users

**Models:**

- `OnboardingSession`: Token-based session, stage tracking

**Stages:**

1. `select_plan` → Choose Free/Starter/Professional
2. `create_workspace` → Set workspace name
3. `create_business` → Set business descriptor
4. `setup_stripe` (paid only) → Stripe Checkout
5. `complete` → Finalize, create workspace/business

**Key rules:**

- Sessions expire after 1 hour
- Stage progression validated (can't skip stages)
- Workspace/business creation atomic
- Cleanup job removes expired sessions

**SSOT**: `.github/instructions/onboarding.instructions.md`

---

## Storefront Domain

**Purpose**: Public-facing storefront for customers (separate from portal)

**Key features:**

- Public CORS (no auth required)
- Read-only access to products, categories
- Order placement (optional auth for registered customers)

**Access pattern:**

- Public routes under `/v1/storefront`
- Business lookup by descriptor
- Rate limiting to prevent abuse

---

## Metadata Domain

**Purpose**: Shared categories, tags, custom fields

**Models:**

- `Category`: Product/expense categories (hierarchical)
- `Tag`: Flexible tagging system

**Used by:**

- Inventory (product categories)
- Accounting (expense categories)
- Orders (custom tags)

---

## Asset Domain

**Purpose**: File uploads, blob storage integration

**Models:**

- `Asset`: Upload metadata, storage path, MIME type

**Storage providers:**

- Local filesystem (dev/test)
- S3-compatible (production)

**Upload flow:**

1. Request presigned upload URL
2. Client uploads to storage
3. Backend records asset metadata

**SSOT**: `.github/instructions/asset_upload.instructions.md`

---

## Cross-Domain Interactions

### Order → Inventory

- Creating order: Reserve stock (decrement variant stock)
- Cancelling order: Release stock (increment variant stock)
- Validation: Check sufficient stock before order creation

### Order → Customer

- Link order to customer profile
- Track customer purchase history
- Calculate customer lifetime value

### Order → Accounting

- Payment succeeded → create asset/revenue record
- Payment fees → create expense record (via bus event)

### Billing → Account

- Subscription status affects workspace access
- Plan limits enforced by middleware
- Downgrade when subscription expires

### Onboarding → Account + Business + Billing

- Creates workspace + admin user
- Creates first business
- Activates paid subscription (if selected)

---

## Testing Kyora Domains

### E2E Test Patterns

```go
// Create workspace + user (via account helper)
user := s.accountHelper.CreateTestUser("test@example.com", "Pass123!")
token := s.accountHelper.LoginUser("test@example.com", "Pass123!")

// Create business (via business helper)
biz := s.businessHelper.CreateBusiness(user.WorkspaceID, "demo")

// Make business-scoped request
resp, _ := s.client.Post(
    "/v1/businesses/"+biz.Descriptor+"/orders",
    orderPayload,
    testutils.WithAuth(token),
)
```

### Isolation Testing

```go
// Verify workspace isolation
user1 := s.createUser("user1@example.com")
user2 := s.createUser("user2@example.com")

biz1 := s.createBusiness(user1.WorkspaceID, "biz1")

// User2 tries to access biz1 (should fail)
token2 := s.loginUser(user2.Email, "Pass123!")
resp := s.client.Get("/v1/businesses/"+biz1.Descriptor+"/orders",
    testutils.WithAuth(token2))

s.Equal(http.StatusForbidden, resp.StatusCode)
```

---

## Quick Reference

**Account**: Users, workspaces, sessions, RBAC  
**Business**: Profiles, descriptors, zones, payment methods  
**Inventory**: Products, variants, stock tracking  
**Order**: Order lifecycle, items, status machine  
**Customer**: Profiles, addresses, notes  
**Accounting**: Assets, expenses, withdrawals, profit calc  
**Analytics**: Dashboards, reports, aggregations  
**Billing**: Stripe, plans, subscriptions, plan gates  
**Onboarding**: Multi-stage new user flow  
**Storefront**: Public product display  
**Metadata**: Categories, tags  
**Asset**: File uploads, blob storage

**Multi-tenancy**: Workspace (top) → Business (second)  
**Critical**: Always scope by workspace/business, never trust URL params
