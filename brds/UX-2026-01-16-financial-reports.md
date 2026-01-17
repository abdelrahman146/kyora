---
status: draft
created_at: 2026-01-16
updated_at: 2026-01-16
brd_ref: "brds/BRD-2026-01-16-financial-reports.md"
owners:
  - area: portal-web
    agent: UI/UX Designer
stakeholders:
  - name: Product Lead
    role: Approver
  - name: Engineering Manager
    role: Implementation Lead
areas:
  - portal-web
---

# UX Spec: Financial Reports — Business Health at a Glance

## 0) Inputs & Scope

- **BRD**: [BRD-2026-01-16-financial-reports.md](./BRD-2026-01-16-financial-reports.md)
- **Goals** (what should feel better for the user):
  - Instant financial clarity: user understands business health within 5 seconds
  - Zero accounting jargon: everything in plain language users already understand
  - Mobile-first excellence: beautiful, fast, fully functional on mobile
  - Confidence building: users feel informed and in control after viewing reports
- **Non-goals**:
  - Advanced accounting features (double-entry, ledgers)
  - Historical period comparisons (future iteration)
  - PDF/export functionality (future iteration)
  - Multi-currency reports
- **Assumptions**:
  - Backend APIs exist: `financial-position`, `profit-and-loss`, `cash-flow`
  - Reports are read-only (no write operations)
  - Both admin and member roles can view reports
  - As-of date defaults to today, editable by user

## 1) Reuse Map (Evidence)

### Existing routes/pages
- `/business/$businessDescriptor/accounting/` — [AccountingDashboard.tsx](../portal-web/src/features/accounting/components/AccountingDashboard.tsx) — Pattern for module landing page with hero stat, summary cards, quick nav cards
- `/business/$businessDescriptor/accounting/expenses/` — Pattern for list pages within a module

### Existing feature components
- `features/accounting/components/AccountingDashboard.tsx` — Reference implementation for:
  - "Safe to Draw" hero stat display with color variants
  - `StatCardGroup` with 3-column layout
  - `QuickNavCard` pattern for sub-page navigation
  - Recent activity list pattern
- `features/dashboard/components/BusinessDashboardPage.tsx` — Simple dashboard card layout pattern

### Existing shared components (components/lib)
- **Stat Cards**: `components/atoms/StatCard.tsx`, `StatCardSkeleton.tsx`
  - Props: `label`, `value`, `icon`, `trend`, `trendValue`, `variant` (success/warning/error/info)
- **Stat Card Group**: `components/molecules/StatCardGroup.tsx`
  - Responsive grid: 1 col mobile, 2 col sm, up to `cols` prop on lg
- **Charts**: 
  - `components/charts/ChartCard.tsx` — Wrapper with loading/empty/error states
  - `components/charts/BarChart.tsx` — Horizontal support, stacked support
  - `components/charts/DoughnutChart.tsx` — Center label support
  - `components/charts/LineChart.tsx` — Area fill, time series
- **BottomSheet**: `components/molecules/BottomSheet.tsx` — Mobile drawer, desktop side panel
- **Skeleton**: `components/atoms/Skeleton.tsx`
- **Button**: `components/atoms/Button.tsx`
- **Currency Formatting**: `lib/formatCurrency.ts`
- **DatePicker**: `components/forms/DatePicker.ts`
- **Date Formatting**: `lib/formatDate.ts`
- **Chart Theme**: `lib/charts/chartTheme.ts`, `useChartTheme` hook

### Do-not-duplicate list (MUST reuse these if applicable)
- **BottomSheet / sheet patterns**: Use existing `BottomSheet` for date picker overlay
- **StatCard / StatCardGroup**: Use for all summary metrics; do NOT create "ReportStatCard"
- **ChartCard**: Use as wrapper for all chart visualizations; do NOT create "ReportChartCard"
- **Empty/loading/error state patterns**: Follow `ChartCard` patterns exactly
- **Layout spacing**: Use `space-y-6` for page sections (matches AccountingDashboard)
- **Currency formatting**: Use `formatCurrency(amount, currency)` everywhere
- **Date formatting**: Use existing `formatDateShort` / `formatDateLong` helpers

## 2) IA + Surfaces

### Page Routes (4 pages)
1. **Reports Hub** — `/business/$businessDescriptor/reports/`
2. **Business Health** — `/business/$businessDescriptor/reports/health`
3. **Profit & Earnings** — `/business/$businessDescriptor/reports/profit`
4. **Cash Movement** — `/business/$businessDescriptor/reports/cashflow`

### Sheets (1 sheet)
1. **Date Picker Sheet** — For selecting "as-of" date on all report pages

### Components (new)
1. **ReportCard** — Clickable card showing key metric + "View Details" CTA (hub page)
2. **AssetBreakdownBar** — Horizontal segmented bar for asset distribution
3. **ProfitWaterfall** — Stepped bar visualization for P&L flow
4. **CashFlowDiagram** — Visual flow for cash in/out
5. **InsightCard** — Actionable insight with icon, message, and optional link

## 3) User Flows (Step-by-Step)

### Flow A — Quick Health Check (Most Common, Mobile-First)

1. User is on any business page (e.g., Dashboard, Orders)
2. User taps "Reports" in **sidebar**
3. User lands on **Reports Hub** (`/reports/`)
   - Sees "Safe to Draw" hero metric immediately (large, prominent)
   - Sees 3 report cards with key metrics visible without scrolling (mobile: vertical stack)
4. User glances at Safe to Draw amount, understands if they can withdraw money
5. User taps "Business Health" card
6. User sees full financial position in plain language
7. User reads insight card, feels informed, navigates back or closes app

### Flow B — Understanding Profit

1. User wonders "Am I making money?"
2. User navigates to Reports Hub → taps "Profit & Earnings" card
3. User sees waterfall visualization: Revenue → Product Costs → Gross Profit → Running Costs → **What You Keep**
4. User scrolls to see expense breakdown by category (doughnut/bar chart)
5. User reads insight: "Your biggest expense is [category]"
6. User taps "View All Expenses" quick link → navigates to Expenses page

### Flow C — Cash Visibility

1. User wonders "Where did my cash go?"
2. User navigates to Reports Hub → taps "Cash Movement" card
3. User sees cash flow diagram: Start → Money In → Money Out → **Cash Now**
4. User scrolls to see detailed breakdown tables
5. User reads insight: "You withdrew more than your sales revenue"
6. User taps "View Capital" quick link → navigates to Capital page

### Flow D — Change Report Date

1. User is on any report page (Hub, Health, Profit, or CashFlow)
2. User taps "As of [date]" button in header area
3. **Mobile**: BottomSheet slides up with date picker calendar
   **Desktop**: Dropdown/popover appears with date picker
4. User selects a different date
5. All metrics on current page refresh with data as of new date
6. If user navigates to another report, selected date persists

## 4) Per-Surface Specification (Implementation-Ready)

---

### Surface 1: Reports Hub (`/business/$businessDescriptor/reports/`)

#### Entry Points
- **Mobile**: Bottom navigation "Reports" item (new nav item needed)
- **Desktop**: Sidebar "Reports" item (new nav item needed)
- Dashboard quick-link (optional, future)

#### Layout Structure (Mobile-First)

```
┌─────────────────────────────────────────┐
│ Header: "Reports" + As-of Date Button   │
├─────────────────────────────────────────┤
│ Hero: Safe to Draw                      │
│ ┌─────────────────────────────────────┐ │
│ │ [Icon] Safe to Draw                 │ │
│ │        $12,500                      │ │
│ │ The amount you can safely take out  │ │
│ └─────────────────────────────────────┘ │
├─────────────────────────────────────────┤
│ Report Cards (vertical stack mobile)    │
│ ┌─────────────────────────────────────┐ │
│ │ Business Health                     │ │
│ │ What Your Business Is Worth: $25K   │ │
│ │ Cash: $8K | Inventory: $15K         │ │
│ │                    [View Details →] │ │
│ └─────────────────────────────────────┘ │
│ ┌─────────────────────────────────────┐ │
│ │ Profit & Earnings                   │ │
│ │ What You Keep: $5,200               │ │
│ │ Revenue: $18K | Costs: $12.8K       │ │
│ │                    [View Details →] │ │
│ └─────────────────────────────────────┘ │
│ ┌─────────────────────────────────────┐ │
│ │ Cash Movement                       │ │
│ │ Cash Now: $8,000                    │ │
│ │ In: $20K | Out: $12K                │ │
│ │                    [View Details →] │ │
│ └─────────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

#### Primary CTA
- **Each ReportCard** is tappable → navigates to detail page
- CTA label: "View Details" / "عرض التفاصيل"
- Entire card is clickable (not just the CTA text)

#### Secondary Actions
- **As-of Date Button**: Opens date picker sheet
- Date format: "As of Jan 16, 2026" / "حتى ١٦ يناير ٢٠٢٦"

#### Hero Section: Safe to Draw
- Use `StatCard` with `variant="success"` (positive) or `variant="error"` (negative)
- Icon: `PiggyBank` from lucide-react
- Large font size for amount: `text-3xl font-bold`
- Subtitle below: "The amount you can safely take out"
- Color coding:
  - Green (`variant="success"`): positive amount
  - Red (`variant="error"`): negative amount or zero

#### Report Cards
- **New component needed**: `ReportCard`
- Structure:
  ```
  ┌──────────────────────────────────┐
  │ [Icon]  Title                    │
  │ Key Metric Label: $Amount        │
  │ Secondary: Val1 | Val2           │
  │                  View Details →  │
  └──────────────────────────────────┘
  ```
- Touch target: entire card is tappable (min 44px height, more realistically ~120px)
- Hover state (desktop): subtle background change (`hover:bg-base-200/50`)
- Active state: `active:scale-[0.98]`
- Border: `border border-base-300 rounded-box`

#### Empty State
- **Trigger**: Business has no orders, expenses, or capital transactions
- **Illustration**: Friendly shop/store graphic (use existing empty state illustration pattern)
- **Headline**: "Your financial picture starts here"
- **Body**: "Once you start recording orders and expenses, you'll see a complete view of your business finances."
- **CTA**: Primary button → "Create Your First Order" → `/orders/`

#### Loading State
- `StatCard` for hero: Use `StatCardSkeleton`
- Report cards: 3 skeleton cards with:
  - Skeleton line for title (40% width)
  - Skeleton line for key metric (60% width)
  - Skeleton line for secondary values (80% width)
  - Pulse animation matching existing patterns

#### Error State
- **Trigger**: API fetch fails
- **Icon**: `AlertCircle` from lucide-react
- **Headline**: "Couldn't load your reports"
- **Body**: "We're having trouble connecting. Please check your internet and try again."
- **CTA**: Ghost button → "Retry" → refetch queries

#### Accessibility Notes
- Page title in `<title>`: "Reports — [Business Name]"
- Hero stat has `aria-label` describing the full context
- Report cards use `<article>` with `aria-labelledby` for title
- Focus visible on all interactive elements
- Color is not the only indicator (icons + text accompany colors)

---

### Surface 2: Business Health (`/business/$businessDescriptor/reports/health`)

#### Entry Points
- Reports Hub → "Business Health" card tap
- Direct URL access

#### Layout Structure (Mobile-First)

```
┌─────────────────────────────────────────┐
│ ← Back   Business Health   As of [date] │
├─────────────────────────────────────────┤
│ A snapshot of what your business owns   │
│ and is worth                            │
├─────────────────────────────────────────┤
│ Hero: What Your Business Is Worth       │
│ ┌─────────────────────────────────────┐ │
│ │        $25,000                      │ │
│ │ Your business value                  │ │
│ │ [?] info tooltip                     │ │
│ └─────────────────────────────────────┘ │
├─────────────────────────────────────────┤
│ Section: What You Own                   │
│ ┌─────────────────────────────────────┐ │
│ │ ████████░░░░░░░░░ $25,000           │ │
│ │ [Green]Cash: $8,000                 │ │
│ │ [Blue]Inventory: $15,000            │ │
│ │ [Purple]Equipment: $2,000           │ │
│ └─────────────────────────────────────┘ │
├─────────────────────────────────────────┤
│ Section: What You Owe                   │
│ ┌─────────────────────────────────────┐ │
│ │ $0                                  │ │
│ │ Kyora doesn't track loans yet.      │ │
│ └─────────────────────────────────────┘ │
├─────────────────────────────────────────┤
│ Section: Owner's Stake                  │
│ ┌─────────────────────────────────────┐ │
│ │ Money You Put In     $10,000        │ │
│ │ Money You Took Out   -$3,000        │ │
│ │ Profit Kept          $18,000        │ │
│ │ ─────────────────────────────       │ │
│ │ Your Business Value  $25,000        │ │
│ └─────────────────────────────────────┘ │
├─────────────────────────────────────────┤
│ Insight Card                            │
│ ┌─────────────────────────────────────┐ │
│ │ ✅ Your business is in a healthy    │ │
│ │ position. You own more than you owe.│ │
│ └─────────────────────────────────────┘ │
├─────────────────────────────────────────┤
│ Quick Links                             │
│ [View Assets →]  [View Capital →]       │
└─────────────────────────────────────────┘
```

#### Header
- **Back button**: `ChevronLeft` (or `ChevronRight` in RTL) → navigates to Reports Hub
- **Title**: "Business Health" / "صحة عملك"
- **As-of date button**: Same as Hub

#### Hero Section
- Large amount: `text-4xl font-bold` (larger than Hub hero)
- Label: "What Your Business Is Worth" / "قيمة عملك"
- Info tooltip icon: `HelpCircle` → explains "This is your business value: everything you own minus what you owe"
- Tooltip: Use daisyUI `tooltip` or custom popover

#### Section: What You Own
- **New component**: `AssetBreakdownBar`
- Horizontal stacked bar showing proportions:
  - Cash on Hand (green/success)
  - Inventory Value (blue/info)
  - Equipment & Assets (purple/secondary)
- Below bar: list with colored dots + labels + amounts
- Total row at bottom: bold, slightly larger

#### Section: What You Owe
- Simple display: "$0" or amount
- Note text: "Kyora doesn't track loans or credit yet. This will be available in a future update."
- Muted styling: `text-base-content/60`

#### Section: Owner's Stake
- Table-like layout (no actual table element):
  ```
  Label              Amount
  ──────────────────────────
  Money You Put In   $10,000
  Money You Took Out -$3,000
  Profit Kept        $18,000
  ══════════════════════════
  Your Business Value $25,000
  ```
- Negative amounts: show with minus sign, red color (`text-error`)
- Final row: divider above, bold styling

#### Insight Card
- **New component**: `InsightCard`
- Structure:
  ```
  ┌────────────────────────────────┐
  │ [Icon] Message text            │
  │        Optional link →         │
  └────────────────────────────────┘
  ```
- Variants:
  - **Positive** (green background): `bg-success/10`, `border-success/20`, icon `CheckCircle`
  - **Warning** (yellow background): `bg-warning/10`, `border-warning/20`, icon `AlertTriangle`
  - **Negative** (red background): `bg-error/10`, `border-error/20`, icon `AlertCircle`
- Logic:
  - TotalEquity > 0: Positive variant, "Your business is in a healthy position."
  - TotalEquity < 0: Negative variant, "Your business value is negative."
  - CashOnHand < 0: Warning variant, "Your cash position is estimated as negative."

#### Quick Links
- Horizontal row of text links (mobile: may wrap to 2 rows)
- Pattern: "[Label] →" with `ChevronRight` (or `ChevronLeft` in RTL)
- Links:
  - "View Assets" → `/accounting/assets`
  - "View Capital" → `/accounting/capital`

#### Empty State
- Same pattern as Hub, context-specific copy

#### Loading State
- Skeleton for hero amount
- Skeleton for asset breakdown bar
- Skeleton rows for tables

#### Error State
- Same pattern as Hub

#### Accessibility Notes
- Back button has `aria-label="Back to Reports"`
- Sections use semantic headings (`<h2>`)
- Info tooltip is keyboard accessible
- Amounts use `aria-label` for screen reader clarity on negative values

---

### Surface 3: Profit & Earnings (`/business/$businessDescriptor/reports/profit`)

#### Entry Points
- Reports Hub → "Profit & Earnings" card tap
- Direct URL access

#### Layout Structure (Mobile-First)

```
┌─────────────────────────────────────────┐
│ ← Back   Profit & Earnings  As of [date]│
├─────────────────────────────────────────┤
│ Where your money comes from and where   │
│ it goes                                 │
├─────────────────────────────────────────┤
│ Hero: What You Keep                     │
│ ┌─────────────────────────────────────┐ │
│ │        $5,200                       │ │
│ │ Your final profit after all costs   │ │
│ └─────────────────────────────────────┘ │
├─────────────────────────────────────────┤
│ Profit Flow (Waterfall)                 │
│ ┌─────────────────────────────────────┐ │
│ │ Revenue         ████████████ $18,000│ │
│ │      ↓                              │ │
│ │ - Product Costs ████      -$6,000   │ │
│ │      ↓                              │ │
│ │ = Gross Profit  ████████   $12,000  │ │
│ │      ↓                              │ │
│ │ - Running Costs ████       -$6,800  │ │
│ │      ↓                              │ │
│ │ = What You Keep █████       $5,200  │ │
│ └─────────────────────────────────────┘ │
├─────────────────────────────────────────┤
│ Margin Stats                            │
│ ┌───────────┐ ┌───────────┐            │
│ │Gross      │ │Profit     │            │
│ │Margin 67% │ │Margin 29% │            │
│ └───────────┘ └───────────┘            │
├─────────────────────────────────────────┤
│ Running Costs Breakdown                 │
│ ┌─────────────────────────────────────┐ │
│ │      [Doughnut Chart]               │ │
│ │    $6,800 total                     │ │
│ │ ○ Marketing    $2,500  37%          │ │
│ │ ○ Operations   $2,000  29%          │ │
│ │ ○ Shipping     $1,500  22%          │ │
│ │ ○ Other        $800    12%          │ │
│ └─────────────────────────────────────┘ │
├─────────────────────────────────────────┤
│ Insight Card                            │
│ ┌─────────────────────────────────────┐ │
│ │ 🎉 You're profitable! For every $1  │ │
│ │ of sales, you keep $0.29            │ │
│ │ 💡 Biggest expense: Marketing $2,500│ │
│ └─────────────────────────────────────┘ │
├─────────────────────────────────────────┤
│ Quick Links                             │
│ [View All Expenses →]  [View Orders →]  │
└─────────────────────────────────────────┘
```

#### Header
- Same pattern as Business Health page

#### Hero Section
- Amount: `text-4xl font-bold`
- Color: Green if positive, Red if negative
- Label: "What You Keep" / "ما تحتفظ به"
- Subtitle: "Your final profit after all costs" / "ربحك النهائي بعد كل التكاليف"

#### Profit Flow (Waterfall)
- **New component**: `ProfitWaterfall`
- Visual representation:
  ```
  Revenue          ████████████████████  $18,000
       ↓
  - Product Costs  ██████               -$6,000
       ↓
  = Gross Profit   ██████████████       $12,000
       ↓
  - Running Costs  ████████             -$6,800
       ↓
  = What You Keep  ██████                $5,200
  ```
- Bar colors:
  - Revenue: primary color
  - Subtracted items: muted/gray
  - Results (Gross Profit, What You Keep): success color
- Arrow indicators between steps: `ArrowDown` icon, muted color
- Mobile: Full width, vertically stacked
- Tablet/Desktop: Can remain same layout, more horizontal space

#### Margin Stats
- Use `StatCardGroup` with 2 columns
- Cards:
  - Gross Margin: percentage, info variant
  - Profit Margin: percentage, success/error variant based on positive/negative
- Calculation (client-side):
  - Gross Margin: `(GrossProfit / Revenue) * 100`
  - Profit Margin: `(NetProfit / Revenue) * 100`

#### Running Costs Breakdown
- Use `ChartCard` wrapper
- Use `DoughnutChart` with center label showing total
- Below chart: legend as list with colored dots
- Each row: `○ [Category] [Amount] [Percentage]%`
- Categories from backend expense data, grouped by category

#### Insight Card
- Multiple insights can stack:
  1. Profitability insight (positive or negative)
  2. Biggest expense insight
- Use `InsightCard` component (same as Health page)

#### Quick Links
- "View All Expenses" → `/accounting/expenses`
- "View Orders" → `/orders`

#### Empty State
- Specific copy: "Start recording orders to see your profit breakdown. Your earnings will appear here automatically."

#### Loading State
- Skeleton for hero
- Skeleton bars for waterfall (5 horizontal lines)
- Skeleton for doughnut chart (circle placeholder)

#### Error State
- Same pattern as Hub

---

### Surface 4: Cash Movement (`/business/$businessDescriptor/reports/cashflow`)

#### Entry Points
- Reports Hub → "Cash Movement" card tap
- Direct URL access

#### Layout Structure (Mobile-First)

```
┌─────────────────────────────────────────┐
│ ← Back   Cash Movement      As of [date]│
├─────────────────────────────────────────┤
│ How cash flows through your business    │
├─────────────────────────────────────────┤
│ Hero: Cash Now                          │
│ ┌─────────────────────────────────────┐ │
│ │        $8,000                       │ │
│ │ Your current cash runway            │ │
│ └─────────────────────────────────────┘ │
├─────────────────────────────────────────┤
│ Cash Flow Diagram                       │
│ ┌─────────────────────────────────────┐ │
│ │ Start                       $2,000  │ │
│ │     │                               │ │
│ │     ▼                               │ │
│ │ ┌─────────────────────────┐         │ │
│ │ │ Money Coming In         │         │ │
│ │ │ + From Customers $18,000│         │ │
│ │ │ + From Owner     $2,000 │         │ │
│ │ │ ─────────────────────── │         │ │
│ │ │ Total In        $20,000 │         │ │
│ │ └─────────────────────────┘         │ │
│ │     │                               │ │
│ │     ▼                               │ │
│ │ ┌─────────────────────────┐         │ │
│ │ │ Money Going Out         │         │ │
│ │ │ - Inventory      $5,000 │         │ │
│ │ │ - Running Costs  $6,000 │         │ │
│ │ │ - Equipment      $1,000 │         │ │
│ │ │ - To Owner       $2,000 │         │ │
│ │ │ ─────────────────────── │         │ │
│ │ │ Total Out       $14,000 │         │ │
│ │ └─────────────────────────┘         │ │
│ │     │                               │ │
│ │     ▼                               │ │
│ │ Net Change              +$6,000 ↑   │ │
│ │     │                               │ │
│ │     ▼                               │ │
│ │ End: Cash Now            $8,000     │ │
│ └─────────────────────────────────────┘ │
├─────────────────────────────────────────┤
│ Insight Card                            │
│ ┌─────────────────────────────────────┐ │
│ │ ✅ Cash Healthy: More cash came in  │ │
│ │ than went out.                      │ │
│ └─────────────────────────────────────┘ │
├─────────────────────────────────────────┤
│ Quick Links                             │
│ [View Capital →]  [View Expenses →]     │
└─────────────────────────────────────────┘
```

#### Header
- Same pattern as other detail pages

#### Hero Section
- Amount: `text-4xl font-bold`
- Color: Green if positive, Red if negative
- Label: "Cash Now" / "النقد الحالي"
- Subtitle: "Your current cash runway" / "المدى النقدي الحالي"

#### Cash Flow Diagram
- **New component**: `CashFlowDiagram`
- Visual flow representation with:
  - Start node (small, muted)
  - "Money Coming In" box (success border/accent)
  - "Money Going Out" box (warning border/accent)
  - Net Change indicator with arrow (up green, down red)
  - End node (emphasized)
- Inside each box: itemized list with amounts
- Connector lines between nodes: vertical line with arrow
- Mobile: full width, stacked vertically
- Desktop: can remain vertical or shift to horizontal flow

#### Money Coming In Box
- Title: "Money Coming In" / "الأموال الواردة"
- Items:
  - From Customers (Sales): `[CashFromCustomers]`
  - From Owner (Investment): `[CashFromOwner]`
  - Divider line
  - Total Cash In: `[TotalCashIn]` (bold)
- Background: `bg-success/5`
- Border: `border-success/20`

#### Money Going Out Box
- Title: "Money Going Out" / "الأموال الصادرة"
- Items:
  - Inventory Purchases: `[InventoryPurchases]`
  - Running Costs (Expenses): `[OperatingExpenses]`
  - Equipment & Assets: `[BusinessInvestments]`
  - To Owner (Withdrawals): `[OwnerDraws]`
  - Divider line
  - Total Cash Out: `[TotalCashOut]` (bold)
- Background: `bg-warning/5`
- Border: `border-warning/20`

#### Net Change Indicator
- Text: "Net Change: [amount]"
- Positive: green text, up arrow (`TrendingUp`)
- Negative: red text, down arrow (`TrendingDown`)
- Accompanying text: "Your cash [increased/decreased] by [amount]"

#### Insight Card
- Same `InsightCard` component
- Logic:
  - NetCashFlow > 0: Positive, "Cash Healthy: More cash came in than went out."
  - NetCashFlow < 0: Warning, "Cash Alert: You spent more than you received."
  - OwnerDraws > CashFromCustomers: Info tip, "You withdrew more than your sales revenue."

#### Quick Links
- "View Capital" → `/accounting/capital`
- "View Expenses" → `/accounting/expenses`

#### Empty State
- "Record your first order to see how cash flows through your business."

#### Loading State
- Skeleton for hero
- Skeleton for flow diagram (boxes with skeleton lines)

#### Error State
- Same pattern

---

### Sheet: Date Picker

#### Trigger
- "As of [date]" button on any report page header

#### Behavior by Device
- **Mobile** (< 768px): BottomSheet slides up from bottom
- **Desktop** (≥ 768px): Dropdown/popover anchored to button

#### Content
- Date picker calendar (use daisyUI or custom calendar component)
- Default selection: Today
- Max date: Today (cannot select future dates)
- Format display: Locale-aware (e.g., "Jan 16, 2026" or "١٦ يناير ٢٠٢٦")

#### Actions
- **Cancel**: Close sheet without changing date
- **Apply/Select**: Close sheet and update all queries with new date

#### State Persistence
- Selected date should be stored in URL search params: `?asOf=2026-01-16`
- Persists across report page navigation
- Resets to today on fresh visit (no `asOf` param)

---

## 5) Responsiveness + RTL Rules

### Mobile-First Layout Rules (< 768px)

1. **Single column layout**: All content stacks vertically
2. **Full-width cards**: Report cards span 100% width with proper padding
3. **Hero section**: Centered, prominent, no horizontal scrolling
4. **Charts**: Full width, minimum height 240px
5. **Touch targets**: Minimum 44x44px for all interactive elements
6. **Bottom navigation**: Reports icon added to existing 4-item nav
7. **Date picker**: BottomSheet overlay
8. **Spacing**: `space-y-6` between major sections, `space-y-4` within sections

### Tablet/Desktop Layout Rules (≥ 768px)

1. **Hub page**: 
   - Hero stays full width
   - Report cards: 3-column grid (`grid-cols-3`) on lg+, 2-column on md
2. **Detail pages**:
   - Max content width: `max-w-4xl mx-auto`
   - Charts can be larger (height 320px)
   - Side-by-side layouts for stat groups
3. **Date picker**: Dropdown popover instead of BottomSheet
4. **Sidebar navigation**: Reports icon in sidebar nav

### RTL Rules (Arabic)

1. **Layout direction**: Entire page uses `dir="rtl"` (inherited from root)
2. **Text alignment**: Natural start alignment (becomes right-aligned)
3. **Icons that need rotation**:
   - `ChevronRight` → `ChevronLeft` for "View Details" arrows
   - `ChevronLeft` → `ChevronRight` for back navigation
   - `ArrowRight` → `ArrowLeft` for links
4. **Logical properties** (MUST use):
   - `ps-*` / `pe-*` instead of `pl-*` / `pr-*`
   - `ms-*` / `me-*` instead of `ml-*` / `mr-*`
   - `start` / `end` instead of `left` / `right`
   - `border-s` / `border-e` instead of `border-l` / `border-r`
5. **Chart RTL**:
   - Use `useChartTheme` hook which handles RTL automatically
   - Legend position: `bottom` works in both directions
6. **Numbers**: Use `tabular-nums` for alignment; numbers remain LTR visually but flow is RTL
7. **Currency**: Format using `formatCurrency()` which handles locale

### Fields with `dir="ltr"` Override
- None on report pages (read-only, no input fields)
- Date picker may need `dir="ltr"` for calendar grid if using numeric date display

---

## 6) Copy & i18n Keys Inventory

### Namespaces to Use
- **Primary**: Create new namespace `reports` for all report-specific content
- **Reuse**: 
  - `common` for shared terms (Retry, Close, etc.)
  - `accounting` for existing terms (Safe to Draw already exists)

### New Keys Needed (en/ar parity required)

#### Namespace: `dashboard` (additions for navigation)

```json
{
  "reports": "Reports",
}
```

---

## 7) Component Gaps / Enhancements

### New Components Needed

#### 1. ReportCard
- **Why needed**: Standardized card for hub page with key metric, secondary info, and CTA
- **Where it should live**: `portal-web/src/features/reports/components/ReportCard.tsx`
- **Reuse candidates**:
  - Dashboard overview cards (can convert to use ReportCard)
  - Future module landing pages (similar pattern)
- **Props**:
  ```ts
  interface ReportCardProps {
    title: string
    icon: LucideIcon
    keyMetric: { label: string; value: string }
    secondaryValues: { label: string; value: string }[]
    href: string
    isLoading?: boolean
  }
  ```

#### 2. InsightCard
- **Why needed**: Reusable insight display with icon, variant, message, and optional link
- **Where it should live**: `portal-web/src/components/molecules/InsightCard.tsx` (shared)
- **Reuse candidates**:
  - Accounting dashboard insights
  - Dashboard alerts/tips
  - Order status warnings
- **Props**:
  ```ts
  interface InsightCardProps {
    variant: 'success' | 'warning' | 'error' | 'info'
    icon?: LucideIcon
    message: string
    link?: { label: string; href: string }
  }
  ```

#### 3. AssetBreakdownBar
- **Why needed**: Horizontal stacked bar showing asset distribution
- **Where it should live**: `portal-web/src/features/reports/components/AssetBreakdownBar.tsx`
- **Reuse candidates**:
  - Inventory category distribution
  - Any proportional breakdown visualization
- **Props**:
  ```ts
  interface AssetBreakdownBarProps {
    segments: Array<{
      label: string
      value: number
      color: 'success' | 'info' | 'secondary' | 'warning'
    }>
    total: number
    currency: string
    showLegend?: boolean
  }
  ```

#### 4. ProfitWaterfall
- **Why needed**: Stepped visualization for P&L flow
- **Where it should live**: `portal-web/src/features/reports/components/ProfitWaterfall.tsx`
- **Reuse candidates**:
  - Could be generalized for any "funnel" visualization
- **Props**:
  ```ts
  interface ProfitWaterfallProps {
    steps: Array<{
      label: string
      value: number
      type: 'total' | 'add' | 'subtract' | 'result'
    }>
    currency: string
  }
  ```

#### 5. CashFlowDiagram
- **Why needed**: Visual flow representation for cash movement
- **Where it should live**: `portal-web/src/features/reports/components/CashFlowDiagram.tsx`
- **Reuse candidates**:
  - Future detailed cash flow analysis
- **Props**:
  ```ts
  interface CashFlowDiagramProps {
    startAmount: number
    endAmount: number
    cashIn: Array<{ label: string; value: number }>
    cashOut: Array<{ label: string; value: number }>
    currency: string
  }
  ```

### Enhancements to Existing Components

#### 2. Sidebar (portal-web/src/components/organisms/Sidebar.tsx)
- **Enhancement**: Add "Reports" nav item
- **Change**: Add to `navItems` array:
  ```ts
  { key: 'reports', icon: FileBarChart, path: '/reports' }
  ```
- **Call sites to update**: None (array drives rendering)

#### 3. StatCard (portal-web/src/components/atoms/StatCard.tsx)
- **Enhancement**: Support larger "hero" variant
- **Change**: Add `size?: 'default' | 'lg'` prop
- **Behavior**: `lg` size uses `text-3xl` or `text-4xl` for value
- **Call sites to update**: None (backward compatible)

---

## 8) Acceptance Checklist (for Engineering Manager)

### Pattern Adherence
- [ ] Reuses existing Kyora UI patterns; no parallel "second versions" of flows
- [ ] Uses `StatCard`, `StatCardGroup`, `ChartCard` from existing components
- [ ] Uses `BottomSheet` for mobile date picker
- [ ] Empty/loading/error states follow existing patterns exactly
- [ ] Uses `formatCurrency()` and `formatDate*()` helpers

### Calm UI Principles
- [ ] Safe defaults: Today's date selected by default
- [ ] No overwhelming data: key metrics visible immediately, details on scroll
- [ ] Positive language for good metrics, gentle warnings for concerns
- [ ] One primary action per view (view details, back, retry)

### Mobile-First Verification
- [ ] Tested on iPhone SE (smallest common screen)
- [ ] Single-column layout on mobile, no horizontal scroll
- [ ] Touch targets minimum 44x44px
- [ ] Report cards fully tappable
- [ ] Date picker as BottomSheet on mobile
- [ ] Back navigation easily reachable

### RTL-First Verification
- [ ] All text naturally right-aligned in Arabic
- [ ] Directional icons flip correctly (ChevronRight ↔ ChevronLeft)
- [ ] Logical properties used (ps/pe, ms/me, start/end)
- [ ] Charts render correctly in RTL via `useChartTheme`
- [ ] Numbers maintain LTR visual order within RTL context

### States Specification
- [ ] Empty state for each page with actionable CTA
- [ ] Loading skeletons match final layout
- [ ] Error state with retry functionality
- [ ] Negative value states with appropriate color coding

### i18n Verification
- [ ] All new keys added to `reports` namespace
- [ ] en/ar parity for all keys
- [ ] Navigation keys added to `dashboard` namespace
- [ ] No hardcoded strings in components

### Component Gaps Addressed
- [ ] ReportCard component created in `features/reports/components/`
- [ ] InsightCard component created in `components/molecules/`
- [ ] AssetBreakdownBar created in `features/reports/components/`
- [ ] ProfitWaterfall created in `features/reports/components/`
- [ ] CashFlowDiagram created in `features/reports/components/`
- [ ] BottomNav updated with Reports item (or documented navigation decision)
- [ ] Sidebar updated with Reports item

### Accessibility
- [ ] Page titles set appropriately
- [ ] Semantic headings for sections
- [ ] Focus management on navigation
- [ ] Color not sole indicator of meaning
- [ ] Screen reader announcements for dynamic content

### Performance
- [ ] Data cached client-side (TanStack Query)
- [ ] Skeleton appears within 100ms
- [ ] Full data loads within 3s on 4G
- [ ] No unnecessary re-fetches when switching between report pages (if date unchanged)
