# Kyora — Agent Orchestration Layer (SSOT)

These instructions are always-on context for working in this repository.

## 1) Product and Audience (Always Assume)

**Kyora** is a SaaS platform for social media commerce entrepreneurs. It acts like a “silent business partner” that manages **orders, inventory, accounting/finance, and insights**.

### Target audience

- **Primary region:** Middle East → **Arabic-first culture**.
- **Users:** solo entrepreneurs, side hustlers, home-based makers (artisans/bakers/fashion), micro-teams (2–5).
- **Tech literacy:** low–moderate; intimidated by Excel/ERPs.
- **Sales channels:** orders originate in **DMs** (Instagram, WhatsApp, TikTok, Facebook).
- **Usage pattern:** **mobile-heavy**; prefer mobile-first solutions.

### Core value proposition

- **Simplicity first:** zero accounting/business-management knowledge required.
- **Automatic heavy lifting:** revenue recognition, profitability, inventory levels, and records happen automatically.
- **Social-media native:** optimized for DM-driven commerce.
- **Peace of mind:** clear financial position, tax-ready records, actionable insights.

### Language and tone

- Use **plain language**. Avoid accounting jargon (ledger, accrual, EBITDA, COGS).
- Prefer: “Profit”, “Cash in hand”, “Best seller”, “Money in/out”, “What to do next”.
- Assume **Arabic/RTL-first** UX and i18n requirements for user-facing UI.

## 2) Product Concepts (What Kyora Does)

### Key modules (conceptual)

- **Order management:** quick order entry from any channel; track status; automatically recognize revenue.
- **Inventory management:** stock visibility, low-stock alerts, best sellers.
- **Customer management:** customer DB from orders; purchase history; best customers.
- **Expense management:** recurring + one-off direct/indirect expenses.
- **Owners management:** investments/withdrawals; safe amount to draw per owner.
- **Accounting & finance:** profit and cash flow in plain language; reporting in background.
- **Analytics:** dashboards without confusing charts/jargon.
- **Team management:** invite members; roles/permissions.
- **Multi-business:** multiple businesses per workspace.

### Storefront (important behavior)

Kyora does **not** do checkout. Sellers and customers prefer DM + payment links/COD/bank transfer.

Kyora’s storefront is a **public SPA per business** that:

1. Lets customers browse catalog and add to cart.
2. “Send order to WhatsApp” creates a **Pending** order in Kyora.
3. Redirects the customer to the seller’s WhatsApp chat with a **prepopulated message** of order details.
4. Seller confirms via quick lookup and updates order status.

## 3) Multi-Tenancy and Security (Non-Negotiable)

- **Workspace** is the top-level tenant.
- **Business** is a second-level scope inside a workspace.
- No cross-workspace and no cross-business data leaks.
- RBAC: `admin` / `member`.

## 4) Repo Reality and “Don’t Hallucinate” Rules

Kyora is **not live** and currently runs **locally**. Many components are incomplete or not yet implemented.

- Never claim a feature exists unless you can point to the code.
- Prefer: “Based on current code…” and verify using search/read before codifying instructions.
- If something is planned but not implemented, label it explicitly as planned.

## 5) Monorepo Structure (Where To Work)

```
kyora/
├── backend/          # Go API (source of truth for business logic)
├── portal-web/       # React dashboard (TanStack stack)
├── storefront-web/   # Customer storefront (React)
├── scripts/          # Repo scripts
└── .github/
    └── instructions/ # Specialized agent rules (SSOT)
```

**Path Prefix Rule:** Always include project prefix (`backend/`, `portal-web/`, `storefront-web/`).

## 6) Components and Status (High Level)

- **Backend:** Go monolith API (Gin, GORM/Postgres, Memcached). Integrations: Stripe, Resend, blob storage; Go HTML email templates. Architecture: domain/platform. Tests: heavy integration tests, minimal unit tests. Status: ~90%.
- **Portal Web:** SPA for clients (React 19 + TanStack Router/Query/Store/Form, i18n, `ky`, Zod, Chart.js, Tailwind v4 + daisyUI). Tests: none yet. Status: ~40%.
- **Storefront Web:** currently unmaintained/deprecated; intended to be replaced/migrated to Portal stack.
- **Planned (not implemented):** marketing website (SSG), docs app, admin dashboard + admin backend, mobile portal (React Native), infra deployment configs.

Breaking changes are acceptable (project under heavy development).

## 7) Instruction File Hierarchy (Which Rules Apply)

**Priority Order (resolve conflicts top-down):**

1. **Project-Specific Instructions** → Most authoritative for that domain
2. **Shared Technical Instructions** → Cross-cutting concerns (design, forms, HTTP)
3. **This File** → Meta-orchestration only

### Backend (Go)

**When:** Modifying `backend/**` or adding API endpoints.
**Read First:**

- `.github/instructions/go-backend-patterns.instructions.md` → Reusable Go backend patterns (Kyora-style)
- `.github/instructions/backend-core.instructions.md` → Architecture, patterns, conventions
- `.github/instructions/backend-testing.instructions.md` → Testing (E2E, unit, coverage) — only when writing tests

**Also Read (if relevant):**

- `.github/instructions/onboarding.instructions.md` → Onboarding flow SSOT (when touching onboarding)
- `.github/instructions/account-management.instructions.md` → Account management SSOT (auth, workspaces, team/invitations)
- `.github/instructions/business-management.instructions.md` → Business management SSOT (business CRUD, archive, shipping zones, payment methods)
- `.github/instructions/billing.instructions.md` → Billing workflow SSOT (plans/subscriptions/invoices/webhooks)
- `.github/instructions/inventory.instructions.md` → Inventory workflow SSOT (products/variants/categories/search/summary)
- `.github/instructions/orders.instructions.md` → Orders SSOT (orders lifecycle, payments, inventory adjustments, plan gates)
- `.github/instructions/customer.instructions.md` → Customer management SSOT (customers/addresses/notes/search/RBAC)
- `.github/instructions/analytics.instructions.md` → Analytics SSOT (dashboard metrics, analytics endpoints, financial reports)
- `.github/instructions/accounting.instructions.md` → Accounting SSOT (assets/investments/withdrawals/expenses/recurring/summary)
- `.github/instructions/resend.instructions.md` → Email functionality
- `.github/instructions/stripe.instructions.md` → Billing/payments
- `.github/instructions/asset_upload.instructions.md` → File uploads (backend contract)

### Portal Web (React Dashboard)

**When:** Modifying `portal-web/**` or building business dashboard features.
**Read First:**

- `.github/instructions/portal-web-architecture.instructions.md` → Tech stack, auth, routing, state management
- `.github/instructions/portal-web-development.instructions.md` → Development workflow, testing, deployment

**Also Read (if relevant):**

- `.github/instructions/onboarding.instructions.md` → Onboarding flow SSOT (when touching onboarding)
- `.github/instructions/account-management.instructions.md` → Account management SSOT (auth, workspaces, team/invitations)
- `.github/instructions/business-management.instructions.md` → Business management SSOT (business CRUD, archive, shipping zones, payment methods)
- `.github/instructions/billing.instructions.md` → Billing workflow SSOT (onboarding payment + future billing UI)
- `.github/instructions/inventory.instructions.md` → Inventory workflow SSOT (inventory list + CRUD sheets)
- `.github/instructions/orders.instructions.md` → Orders SSOT (orders list UI, query params, mutations)
- `.github/instructions/customer.instructions.md` → Customer management SSOT (customers/addresses/notes/search/RBAC)
- `.github/instructions/analytics.instructions.md` → Analytics SSOT (dashboard metrics, analytics endpoints, financial reports)
- `.github/instructions/accounting.instructions.md` → Accounting SSOT (assets/investments/withdrawals/expenses/recurring/summary)
- `.github/instructions/forms.instructions.md` → Form system (TanStack Form + all field components)
- `.github/instructions/portal-web-code-structure.instructions.md` → Code structure SSOT (routes/components/features/lib), no-legacy
- `.github/instructions/ui-implementation.instructions.md` → Components, RTL rules, daisyUI usage
- `.github/instructions/charts.instructions.md` → Chart.js visualizations, statistics
- `.github/instructions/design-tokens.instructions.md` → Colors, typography, spacing (SSOT)
- `.github/instructions/ky.instructions.md` → HTTP requests
- `.github/instructions/http-tanstack-query.instructions.md` → HTTP client + TanStack Query usage (no direct API calls, global error handling)
- `.github/instructions/state-management.instructions.md` → State ownership (URL vs Query vs Store vs Form), TanStack Store rules
- `.github/instructions/i18n-translations.instructions.md` → i18n keys, namespaces, no-duplication rules, locale parity (portal-web + storefront-web)
- `.github/instructions/errors-handling.instructions.md` → Errors & failure handling (backend + portal-web)
- `.github/instructions/responses-dtos-swagger.instructions.md` → Responses, DTOs & Swagger/OpenAPI (backend + portal-web)
- `.github/instructions/asset_upload.instructions.md` → File uploads (frontend flow)
- `.github/instructions/stripe.instructions.md` → Billing UI

### Storefront Web (Customer Portal)

**When:** Modifying `storefront-web/**` or building customer-facing features.
**Read First:** `storefront-web/DESIGN_SYSTEM.md` (storefront-specific patterns)
**Also Read (if relevant):**

- `.github/instructions/design-tokens.instructions.md` → Colors, typography, spacing (SSOT)
- `.github/instructions/ui-implementation.instructions.md` → Components, RTL rules, daisyUI usage

## 8) Execution Standards (Non-Negotiable)

**Code Quality Pillars:**

- **Robust:** Production-ready, handles edge cases, defensive coding
- **Secure:** No SQL injection, XSS, CSRF; validate all inputs
- **Maintainable:** Self-documenting code, clear naming, no TODOs/FIXMEs
- **Scalable:** Efficient queries, proper indexing, connection pooling

**Development Rules:**

- ✅ Complete implementations (never partial/example code)
- ✅ Delete deprecated code immediately (no "marked as deprecated")
- ✅ Refactor duplicates into shared utilities (`backend/internal/platform/utils/` or `portal-web/src/lib/`)
- ✅ Fix inefficiencies and anti-patterns when encountered
- ❌ No TODOs, FIXMEs, or placeholder comments
- ❌ No "brief implementations" or "scaffold for later"

**Breaking Changes:** Allowed (project under heavy development). Prioritize simplicity over backward compatibility.

## 9) Agent Decision Tree (How To Proceed)

```
Task received
    ↓
Does it modify backend/?
    YES → Read backend-core.instructions.md
        ↓
        Writing tests? → Read backend-testing.instructions.md
        Email/billing/assets? → Read resend/stripe/asset_upload instructions
    NO → Continue
        ↓
Does it modify portal-web/?
    YES → Read portal-web-architecture.instructions.md
        ↓
        Forms/HTTP/UI? → Read forms/ky/ui-implementation instructions
    NO → Continue
        ↓
Does it modify storefront-web/?
    YES → Read storefront-web/DESIGN_SYSTEM.md + design-tokens + ui-implementation
    NO → Error: unknown target
        ↓
Implement following SSOT rules
    ↓
Verify no TODOs, no duplication, production-ready
    ↓
Done
```

## 10) Anti-Patterns (Never Do This)

- ❌ **Cross-Domain References:** Don't copy-paste rules from one instruction file to another. Link to SSOT.
- ❌ **Vague Directives:** "Consider performance" is useless. Specify: "Use bulk inserts for >100 rows."
- ❌ **Conflicting Rules:** If backend says "use service pattern" and frontend says "inline logic," escalate.
- ❌ **Hallucinated Requirements:** If user says "add feature X" but X isn't in domain context, ask for clarification.
- ❌ **Token Waste:** Verbose explanations ("Now I will proceed to..."). Just execute.
- ❌ **Design Assumptions:** Never assume left/right (RTL-first). Never assume English labels (i18n required).

## 11) Conflict Resolution Protocol

**If instructions conflict:**

1. **Same-Level Conflict** (e.g., two instruction files disagree):

   - Prefer project-specific file (backend-core.instructions.md > asset_upload.instructions.md)
   - If still ambiguous, ask user

2. **Cross-Project Conflict** (e.g., backend pattern vs frontend pattern):

   - No conflict — each domain has its own rules
   - Backend is source of truth for API contracts
   - Frontend is source of truth for UI/UX patterns

3. **Legacy Code vs Instructions:**
   - Instructions win (refactor legacy code to match)
   - Exception: If refactor breaks production, ask user

## 12) Token Budget Guidance

**File Reading Strategy:**

- Read instruction files in full (they are optimized)
- Read source code selectively (use grep/semantic search first)
- Avoid reading entire `node_modules/` or `vendor/` directories

**Context Window Management:**

- Core instructions: ~300-400 lines each (token-optimized)
- Testing instructions: ~250 lines (loaded only when needed)
- Total instruction corpus: ~5,000 lines (~160K tokens after optimization)
- **Budget Per Task:** ~80K tokens context (40K instructions + 40K code)
- If nearing limit, prioritize: instruction file > domain models > handlers > tests

## 13) Meta-Instructions (For This File)

**Purpose:** Orchestrate agent behavior, not duplicate specialized rules.
**Scope:** What to build (domain), where to work (structure), which rules apply (hierarchy).
**Maintenance:** Update only when:

- New instruction file added
- Monorepo structure changes
- Conflict resolution rules change
- Business domain fundamentally shifts

**Never Add Here:**

- Code patterns (belongs in project-specific instructions)
- Tech stack details (belongs in project README or instructions)
- Color palettes, typography, components (belongs in design-tokens/ui-patterns)
