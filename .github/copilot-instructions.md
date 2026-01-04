# Kyora — Agent Orchestration Layer

## 1. Domain Context (30-Second Brief)

**Product:** B2B SaaS for Middle East social commerce entrepreneurs. Automates accounting, inventory, and revenue recognition.
**Tech:** Go monolith (backend), React + TanStack (portal-web), React (storefront-web).
**User:** Non-technical business owners selling via Instagram/WhatsApp/TikTok DMs.
**Philosophy:** "Professional tools that feel effortless" — zero accounting knowledge required.
**Architecture:** Workspace-based multi-tenancy. RBAC: admin/member. Billing via Stripe.

## 2. Business Logic (What To Build)

**Core Pain Solved:** "Financial blindness" — users don't know profit vs revenue, inventory levels, or cash position.

**Key Flows:**

- Onboarding → workspace setup → add first order in <60s
- Order entry (DM source) → auto revenue recognition → inventory update → profit calculation
- Dashboard → plain-language insights ("You made $X profit this month", "Top seller: Y")

**Multi-Tenancy Rule:** All data scoped to `workspaceId`. No cross-workspace leaks.

**Avoid:** Accounting jargon (EBITDA, ledgers, accruals). Use: "Profit", "Cash in hand", "Best seller".

## 3. Monorepo Structure (Where To Work)

```
kyora/
├── backend/          # Go API (source of truth for business logic)
├── portal-web/       # React dashboard (TanStack stack)
├── storefront-web/   # Customer storefront (React)
└── .github/
    └── instructions/ # Specialized agent rules (SSOT)
```

**Path Prefix Rule:** Always include project prefix (`backend/`, `portal-web/`, `storefront-web/`).

## 4. Instruction File Hierarchy (Which Rules Apply)

**Priority Order (resolve conflicts top-down):**

1. **Project-Specific Instructions** → Most authoritative for that domain
2. **Shared Technical Instructions** → Cross-cutting concerns (design, forms, HTTP)
3. **This File** → Meta-orchestration only

### Backend (Go)

**When:** Modifying `backend/**` or adding API endpoints.
**Read First:**

- `.github/instructions/backend-core.instructions.md` → Architecture, patterns, conventions
- `.github/instructions/backend-testing.instructions.md` → Testing (E2E, unit, coverage) — only when writing tests

**Also Read (if relevant):**

- `.github/instructions/resend.instructions.md` → Email functionality
- `.github/instructions/stripe.instructions.md` → Billing/payments
- `.github/instructions/asset_upload.instructions.md` → File uploads (backend contract)

### Portal Web (React Dashboard)

**When:** Modifying `portal-web/**` or building business dashboard features.
**Read First:**

- `.github/instructions/portal-web-architecture.instructions.md` → Tech stack, auth, routing, state management
- `.github/instructions/portal-web-development.instructions.md` → Development workflow, testing, deployment

**Also Read (if relevant):**

- `.github/instructions/forms.instructions.md` → Form system (TanStack Form + all field components)
- `.github/instructions/ui-implementation.instructions.md` → Components, RTL rules, daisyUI usage
- `.github/instructions/charts.instructions.md` → Chart.js visualizations, statistics
- `.github/instructions/design-tokens.instructions.md` → Colors, typography, spacing (SSOT)
- `.github/instructions/ky.instructions.md` → HTTP requests
- `.github/instructions/asset_upload.instructions.md` → File uploads (frontend flow)
- `.github/instructions/stripe.instructions.md` → Billing UI

### Storefront Web (Customer Portal)

**When:** Modifying `storefront-web/**` or building customer-facing features.
**Read First:** `storefront-web/DESIGN_SYSTEM.md` (storefront-specific patterns)
**Also Read (if relevant):**

- `.github/instructions/design-tokens.instructions.md` → Colors, typography, spacing (SSOT)
- `.github/instructions/ui-implementation.instructions.md` → Components, RTL rules, daisyUI usage

## 5. Execution Standards (Non-Negotiable)

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

## 6. Agent Decision Tree (How To Proceed)

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

## 7. Anti-Patterns (Never Do This)

- ❌ **Cross-Domain References:** Don't copy-paste rules from one instruction file to another. Link to SSOT.
- ❌ **Vague Directives:** "Consider performance" is useless. Specify: "Use bulk inserts for >100 rows."
- ❌ **Conflicting Rules:** If backend says "use service pattern" and frontend says "inline logic," escalate.
- ❌ **Hallucinated Requirements:** If user says "add feature X" but X isn't in domain context, ask for clarification.
- ❌ **Token Waste:** Verbose explanations ("Now I will proceed to..."). Just execute.
- ❌ **Design Assumptions:** Never assume left/right (RTL-first). Never assume English labels (i18n required).

## 8. Conflict Resolution Protocol

-core
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

## 9. Token Budget Guidance

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

## 10. Meta-Instructions (For This File)

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
