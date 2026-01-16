---
status: draft
created_at: 2026-01-15
updated_at: 2026-01-15
brd_ref: "brds/BRD-2026-01-15-accounting-management.md"
owners:
  - area: backend
    agent: Feature Builder
  - area: portal-web
    agent: Feature Builder
  - area: tests
    agent: Feature Builder
risk_level: low
---

# Engineering Plan: Accounting & Expense Management

## 0) Inputs

- BRD: [BRD-2026-01-15-accounting-management.md](brds/BRD-2026-01-15-accounting-management.md)
- UX Spec: [UX-2026-01-15-accounting-management.md](brds/UX-2026-01-15-accounting-management.md)
- Assumptions:
  - Backend endpoints are already fully implemented and available under `/v1/businesses/:businessDescriptor/accounting/...` (verified in `routes.go`).
  - No new backend work is required except potentially fixing minor bugs if discovered.
  - "Safe to Draw" logic resides 100% on the backend (`GET /summary`).
  - Search is NOT supported by backend; list views will not have a search bar.

## 1) Confirmation Gate

No confirmation-gated changes proposed.
- No new dependencies.
- No database migrations (backend exists).
- No new infrastructure.

## 2) Architecture summary

- **Backend:** Existing `accounting` domain. Routes are business-scoped. RBAC is standard (`ActionView`/`ActionManage` on `ResourceAccounting`).
- **Portal-web:**
  - **Feature Module:** `src/features/accounting/` (New).
  - **API Client:** `src/api/accounting.ts` + `src/api/types/accounting.ts` (Zod schemas).
  - **State:** URL-driven state (filters, page) + TanStack Query (`accountingQueries`).
  - **UI Patterns:** `ResourceListLayout` for lists, `DashboardLayout` for the main view, `BottomSheet` for mutations.
- **Data model:** No changes.
- **Security:** Standard RBAC. Member role can VIEW, Admin role can MANAGE. UI must hide "Add" buttons for non-admins (or show disabled with tooltip).

### Code Structure & Reuse

- **New Files (Feature Local):**
  - `src/features/accounting/components/AccountingDashboard.tsx`
  - `src/features/accounting/components/ExpenseList.tsx`
  - `src/features/accounting/components/CapitalList.tsx`
  - `src/features/accounting/components/AssetList.tsx`
  - `src/features/accounting/components/cards/*.tsx` (Resource cards)
  - `src/features/accounting/components/sheets/*.tsx` (Mutation forms)
- **Shared Libs (Reuse):**
  - `src/components/templates/ResourceListLayout.tsx` (Reuse for all lists)
  - `src/components/form/*` (All inputs)
  - `src/features/dashboard-layout/components/DashboardLayout.tsx`
  - `formatCurrency`
- **Do Not Duplicate:**
  - Date pickers, Modal shells, Query client setup.

### Repo Recon (Evidence)
- **Backend:** `registerBusinessScopedRoutes` in `routes.go` lines 428-485 confirms all accounting routes exist.
- **Frontend:** `portal-web/src/features/accounting` does not exist. This is a greenfield frontend feature.

## 3) Step-based execution plan

Execution protocol:
- Feature Builder will implement **one step per request**.
- Do not start Step N+1 before Step N is merged/verified.

### Step Index

- **Step 1:** API Client & i18n Foundation
- **Step 2:** Routing & Dashboard UI
- **Step 3:** Expenses Management (List + Create)
- **Step 4:** Recurring Expenses Support
- **Step 5:** Capital Management (Transactions)
- **Step 6:** Asset Management

---

### Step 1: API Client & i18n Foundation

**Goal:** Establish the data layer and translation strings required for the UI.

**Scope:**
- Create typed API client matching backend JSON.
- Create Zod schemas for responses.
- Add English translation keys.

**Tasks:**

- [ ] **Portal:** Create `src/api/types/accounting.ts`.
  - Define Zod schemas for `Expense`, `RecurringExpense`, `Investment`, `Withdrawal`, `Asset`, `AccountingSummary`.
  - **Drift Note:** Backend returns camelCase; ensure Zod matches. Dates are strings (RFC3339).
- [ ] **Portal:** Create `src/api/accounting.ts`.
  - Implement `accountingApi` object with methods for all endpoints documented in `accounting.instructions.md`.
  - Implement `accountingQueries` factory for TanStack Query options (keys: `['accounting', 'expenses', descriptor, filters]`, etc.).
- [ ] **Portal:** Create `src/i18n/locales/en/accounting.json`.
  - Populate with keys from UX Spec Section 7.
- [ ] **Portal:** Register the new namespace in `src/i18n/config.ts` (if manual registration needed) or ensure it loads.

**Verification:**
- Run a script or small test to verify `accountingApi.getSummary` calls the correct URL.
- Verify TypeScript types compile.

---

### Step 2: Routing & Dashboard UI

**Goal:** Users can navigate to `/accounting` and see high-level stats ("Safe to Draw").

**Scope:**
- Create routes.
- Build the Dashboard index page.
- "Safe to Draw" Hero card + Stats row.
- "Recent Activity" list (read-only view).

**Tasks:**

- [ ] **Portal:** Create route file `src/routes/business/$businessDescriptor/accounting/index.tsx`.
- [ ] **Portal:** Create `src/features/accounting/components/AccountingDashboard.tsx`.
- [ ] **Portal:** Implement "Safe to Draw" Hero card.
  - Use `accountingQueries.summary`.
  - Note: `safeToDraw` can be negative; handle UI color (Red).
- [ ] **Portal:** Implement "Recent Activity" list.
  - Fetch latest 5 items (might need to fetch expenses/transactions separately and merge client-side if backend doesn't have a "unified feed" endpoint yet. *Correction:* Backend instructions don't show a unified feed. We will fetch Expenses (limit 5) and Capital (limit 5) and merge/sort client-side for V1, or just show two small lists). *Decision: Just show "Recent Expenses" for now to keep it simple, or parallel queries.*
- [ ] **Portal:** Add "Accounting" link to the Sidebar (in `src/features/dashboard-layout/utils/navigation.ts` or similar).

**Verification:**
- Navigate to `/accounting`.
- Verify "Safe to Draw" statistic loads from backend.
- Verify navigation bar highlights "Accounting".

---

### Step 3: Expenses Management (List + Create)

**Goal:** Users can view full list of expenses and create new ones.

**Scope:**
- `/accounting/expenses` route.
- `ResourceListLayout` integration.
- `CreateExpenseSheet` form.

**Tasks:**

- [ ] **Portal:** Create route `src/routes/business/$businessDescriptor/accounting/expenses.tsx`.
- [ ] **Portal:** Create `src/features/accounting/components/ExpenseList.tsx`.
  - Use `ResourceListLayout`.
  - Filters: Date Range, Category.
  - List Item: `ExpenseCard` (Icon, Category, Amount, Date).
- [ ] **Portal:** Create `src/features/accounting/components/sheets/CreateExpenseSheet.tsx`.
  - Form Fields: Amount, Category (select), Date, Note.
  - Toggle: "Recurring" (logic in next step, visuals now).
  - Validation: Amount > 0.
- [ ] **Portal:** Wire up "Add Expense" button to open the sheet.
- [ ] **Portal:** Implement Create mutation (`accountingApi.createExpense`).
  - Invalidate queries on success.

**Verification:**
- Go to Expenses list.
- Click Add, fill form, save.
- Verify list updates.
- Verify Dashboard "Safe to Draw" updates.

---

### Step 4: Recurring Expenses Support

**Goal:** Handle the "Recurring" toggle logic in the Create form and list recurring templates.

**Scope:**
- Enhance `CreateExpenseSheet` to handle `isRecurring` toggle.
- Switch mutation API call based on toggle.
- Visualize "Recurring" badge in list.

**Tasks:**

- [ ] **Portal:** Update `CreateExpenseSheet` logic.
  - If `isRecurring` is true:
    - Show `frequency` select.
    - Show `autoBackfill` checkbox.
    - Submit to `createRecurringExpense` instead of `createExpense`.
- [ ] **Portal:** Update `ExpenseList` to optionally show recurring templates?
  - *Decision Check:* BRD says "List active templates" is separate. For V1, we will mix them or add a Tab "Recurring Templates" if easy.
  - *Decision:* Stick to simple Expenses List (occurrences) for now. The *creation* flow is the priority.
- [ ] **Portal:** Ensure `ExpenseCard` shows a recurring badge if the expense record has `recurringExpenseId` (if backend returns it) or just trust the create flow works.

**Verification:**
- Create a recurring expense (e.g., Weekly).
- Verify backend creates the template AND the backfilled history (if checked).

---

### Step 5: Capital Management (Transactions)

**Goal:** Manage Investments and Withdrawals.

**Scope:**
- `/accounting/capital` route.
- List view for combined Investments/Withdrawals.
- `CreateTransactionSheet`.

**Tasks:**

- [ ] **Portal:** Create route `src/routes/business/$businessDescriptor/accounting/capital.tsx`.
- [ ] **Portal:** Create `src/features/accounting/components/CapitalList.tsx`.
  - Queries: Fetch `investments` and `withdrawals` in parallel (or separate tabs). *Decision: Separate tabs "Investments" vs "Withdrawals" inside the page is safer for data types.*
  - List Item: `TransactionCard`.
- [ ] **Portal:** Create `src/features/accounting/components/sheets/CreateTransactionSheet.tsx`.
  - Type Switcher: Investment vs Withdrawal.
  - Logic: Submit to correct endpoint.
  - Helper: Show "Safe to Draw" value when "Withdrawal" is selected.

**Verification:**
- Record an investment. Verify it appears.
- Record a withdrawal. Verify it appears and Safe to Draw decreases.

---

### Step 6: Asset Management

**Goal:** Track fixed assets.

**Scope:**
- `/accounting/assets` route.
- List view.
- `CreateAssetSheet`.

**Tasks:**

- [ ] **Portal:** Create route `src/routes/business/$businessDescriptor/accounting/assets.tsx`.
- [ ] **Portal:** Create `src/features/accounting/components/AssetList.tsx`.
- [ ] **Portal:** Create `src/features/accounting/components/sheets/CreateAssetSheet.tsx`.
  - Fields: Name, Value, Date, Note.

**Verification:**
- Add an asset.
- Verify it lists correctly.
- Verify "Total Assets" on Dashboard updates.
