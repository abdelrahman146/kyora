# UX Spec: Accounting & Expense Management

**References:**
- BRD: [BRD-2026-01-15-accounting-management.md](BRD-2026-01-15-accounting-management.md)

## 1. Inputs & Scope

**Goal:** Implement the "Accounting" module to help users track expenses, capital, and assets, focusing on the "Safe to Draw" metric.
**Context:** Mobile-first, Arabic-first. Users are non-accountants.
**Constraints:** Must use existing backend API endpoints (`/accounting/expenses`, `/summary`, etc.).
**Backend Limitations:** No search endpoints exist for accounting resources today. Lists must rely on date-range filtering, paging, and client-side scanning if needed (though paginated).

---

## 2. Reuse Map (Evidence)

We will leverage existing portal patterns to avoid fragmentation.

| Component | Path | Usage |
|-----------|------|-------|
| **Layout** | `src/features/dashboard-layout/components/DashboardLayout.tsx` | Main shell for all pages. |
| **List Template** | `src/components/templates/ResourceListLayout.tsx` | For Expenses/Capital/Assets list pages (handles search, filter, pagination). |
| **Data Cards** | `src/features/orders/components/OrderCard.tsx` (pattern) | Reference for creating `ExpenseCard`, `TransactionCard`, `AssetCard`. |
| **Stats** | `src/components/molecules/StatCardGroup.tsx` | For the top-level summary stats. |
| **Sheets** | `src/components/molecules/BottomSheet.tsx` | For "Add Expense" and other forms. |
| **Forms** | `src/components/form/*` | `FormInput`, `FormSelect`, `DatePicker`, `FormTextarea`, `FormToggle`. |
| **Formatting** | `src/lib/formatCurrency.ts` | for prices. |

**Do Not Build:**
- A new table component (use `ResourceListLayout`'s slot or `Table` if desktop view needed).
- Custom date pickers (use `DatePicker`).
- Custom toast notifications (use `src/lib/toast.ts`).

---

## 3. IA + Surfaces

### Route Structure
All under `/business/$businessDescriptor/accounting/`:

1.  **Dashboard (Index)**: `/`
    - Summary stats ("Safe to Draw" + others).
    - Shortcuts to lists.
    - Recent activity feed.
2.  **Expenses List**: `/expenses`
    - Filterable list of operating expenses.
    - Primary action: "Add Expense".
3.  **Capital List**: `/capital`
    - Mixed list of Investments & Withdrawals.
    - Primary action: "Record Transaction".
4.  **Assets List**: `/assets`
    - List of fixed assets.
    - Primary action: "Add Asset".

### Sheets (Modals/Drawers)
- `CreateExpenseSheet` (and Edit)
- `CreateTransactionSheet` (Invest/Withdraw)
- `CreateAssetSheet`

---

## 4. User Flows

### Flow 1: Check Financial Health
1.  User taps "Accounting" in sidebar.
2.  Lands on **Accounting Dashboard**.
3.  Sees "Safe to Draw" summary prominently (Hero on mobile, Grid on desktop).
4.  Scans **"Recent Activity"** list below it.

### Flow 2: Log a One-Time Expense (Dashboard Entry)
1.  From Dashboard, user taps **"+ Expense"** (FAB on mobile, Header action on desktop).
2.  **CreateExpenseSheet** opens.
3.  User enters **Amount**, selects **Category** (e.g., Packaging), adds **Note**.
4.  Taps **"Save Expense"**.
5.  Sheet closes, toast appears ("Expense saved").
6.  Dashboard refreshes, "Safe to Draw" decreases.

### Flow 3: Log a One-Time Expense (List Entry)
1.  Only on **Expenses List** page.
2.  User taps **"+ Add Expense"** (Primary button in header).
3.  **CreateExpenseSheet** opens.
4.  User fills and saves.
5.  List refreshes to show the new item at the top.

### Flow 4: Log a Withdrawal (Smart Helper)
1.  From Capital List (or Dashboard), tap **"Record Transaction"** -> Select **"Withdrawal"**.
2.  **CreateTransactionSheet** opens (Mode: Withdrawal).
3.  Helper text shows: *"Safe to withdraw: SAR 5,000"*.
4.  User enters Amount (e.g., 500).
5.  Taps **"Withdraw"**.
6.  User is returned to the list, which now shows the withdrawal in Red.

### Flow 5: Add a Fixed Asset
1.  User navigates to **Assets List**.
2.  Taps **"+ Add Asset"**.
3.  **CreateAssetSheet** opens.
4.  User enters Name ("Toyota Van"), Value (40000), Purchase Date.
5.  Taps **"Save Asset"**.
6.  List refreshes.

---

## 5. Per-Surface Spec

### 5.1 Accounting Dashboard
**Route:** `/business/$businessDescriptor/accounting/`

**Layout Strategy:**
- **Mobile:** Vertical stack. Hero card is full width.
- **Desktop:** 2-column grid.
    - Left/Top: Summary Stats grid (Safe to Draw, Assets, Expenses). Safe to Draw is prominent but constrained size (stat card), not a banner.
    - Right/Bottom: Recent Activity list.

**Components:**
- **Header:** Title: "Accounting".
- **Summary Section:**
    - `StatCardGroup` containing:
        - **Safe to Draw** (Highlighted/Large).
        - **Expenses (This Month)**.
        - **Total Assets**.
- **Quick Links:**
    - Button group or Cards navigating to sub-pages.
- **Recent Activity:**
    - List: Last 5 items (mixed `Expense` and `Withdrawal/Investment`).
    - **Empty State:** "No activity yet. Add an expense or investment."

**Actions:**
- **Mobile FAB:** "+ Expense".
- **Desktop Header Action:** Split button "Add Expense" / "Record Transaction".

### 5.2 Expenses List
**Route:** `/expenses`

**Template:** `ResourceListLayout`
- **Heading:** "Expenses".
- **Primary Action (Top Right):** "+ Add Expense" (opens `CreateExpenseSheet`).
- **Filters:**
    - Date Range (`date-range`).
    - Category (`select`).
    - Type: `One-time` vs `Recurring`.
- **Search:** *Disabled* (Backend does not support search param).
- **List Item (`ExpenseCard`):**
    - **Icon:** Category icon (Lucide react icons mapped to category).
    - **Main:** Category Name + Note (truncated).
    - **Right:** Amount (Red color).
    - **Sub:** Date.
    - **Badge:** "Recurring" (if applicable).

### 5.3 Capital List
**Route:** `/capital`

**Template:** `ResourceListLayout`
- **Heading:** "Capital & Draws".
- **Primary Action (Top Right):** "Record Transaction" (opens `CreateTransactionSheet`).
- **Filters:** Type (Investment / Withdrawal), Date Range.
- **Search:** *Disabled*.
- **List Item (`TransactionCard`):**
    - **Icon:** arrow-down-circle (Green/Invest), arrow-up-circle (Red/Withdraw).
    - **Main:** "Investment" or "Withdrawal".
    - **Right:** Amount (Green for In, Red for Out).

### 5.4 Assets List
**Route:** `/assets`

**Template:** `ResourceListLayout`
- **Heading:** "Assets".
- **Primary Action (Top Right):** "+ Add Asset" (opens `CreateAssetSheet`).
- **Filters:** Date Range.
- **Search:** *Disabled*.
- **List Item (`AssetCard`):**
    - **Icon:** box/package icon.
    - **Main:** Asset Name.
    - **Right:** Value.
    - **Sub:** Purchased on [Date].

### 5.5 Forms (Sheets)

#### Create Expense Sheet
- **Components:** `BottomSheet`, `useKyoraForm`.
- **Fields:**
    1.  `amount` (`FormInput`, type="number", min="0").
        - *Validation:* Required, > 0.
    2.  `category` (`FormSelect`).
        - *Options (i18n):* rent, marketing, salaries, packaging, software, logistics, transaction_fee, other.
    3.  `date` (`DatePicker`). Default: Today.
    4.  `note` (`FormTextarea`). Optional.
    5.  `isRecurring` (`FormToggle` or `FormCheckbox`). "Repeat this expense?".
        - *Conditional:* If true, show `frequency` (`FormSelect`: weekly, monthly, yearly).
        - *Conditional:* If true, show `autoBackfill` (`FormCheckbox`). "Log past expenses from start date?".

#### Create Transaction Sheet
- **Fields:**
    1.  `type` (`FormRadio` / Segmented Control): "Investment" vs "Withdrawal".
    2.  `amount` (`FormInput`).
        - *Helper (Withdrawal only):* "Safe to draw: [Amount]".
    3.  `date` (`DatePicker`).
    4.  `note` (`FormTextarea`).

#### Create Asset Sheet
- **Fields:**
    1.  `name` (`FormInput`).
    2.  `value` (`FormInput`, type="number").
    3.  `purchaseDate` (`DatePicker`).
    4.  `note` (`FormTextarea`).

---

## 6. Responsiveness & RTL

- **Mobile First:** All lists render as Cards (`ExpenseCard`), not Tables.
- **Desktop:** `ResourceListLayout` constrains width. "Safe to Draw" card should not stretch awkwardly; use `max-w` or grid system.
- **RTL:** Prices (`SAR 500` -> `٥٠٠ ر.س`), directional icons flip.

---

## 7. i18n Keys Inventory

**Namespace:** `accounting`

| Key | English | Context |
|-----|---------|---------|
| `nav.title` | Accounting | Sidebar |
| `stats.safe_to_draw` | Safe to Draw | Dashboard Hero |
| `stats.total_assets` | Total Assets | Dashboard Stat |
| `stats.expenses_month` | Expenses (Month) | Dashboard Stat |
| `actions.add_expense` | Add Expense | Button |
| `actions.record_transaction` | Record Transaction | Button |
| `header.expenses` | Expenses | Page Title |
| `header.capital` | Capital & Draws | Page Title |
| `category.rent` | Rent | |
| `category.marketing` | Marketing | |
| ... | (Full category list) | |
| `helper.safe_amount` | Safe limit: {{amount}} | Withdrawal form |

---

## 8. Component Gaps / Enhancements

1.  **`ResourceListLayout` usage:** Ensure it supports "mixed" lists (for Dashboard Recent Activity) or just standard pagination. *Decision:* Use `ResourceListLayout` only for full pages. For Dashboard "Recent Activity", build a simple `<div>` list since it's limited to 5 items.
2.  **`ExpenseCard` & `TransactionCard`:** Need to create these. Recommend creating a generic `AccountingResourceCard` that takes `icon`, `title`, `subtitle`, `amount`, `variant (positive/negative)` to maximize reuse.

## 9. Acceptance Checklist

- [ ] "Safe to Draw" calculation matches backend summary exactly.
- [ ] Adding an expense immediately updates the underlying query cache (optimistic or invalidation).
- [ ] Recurring expense toggle correctly switches the API payload (to `/recurring-expenses`).
- [ ] Withdrawal form shows the "Safe limit" helper text.
- [ ] Empty states are helpful ("Add your rent...").
- [ ] Mobile view is fully functional (no horizontal scrolling tables).
- [ ] RTL layout is verified.
