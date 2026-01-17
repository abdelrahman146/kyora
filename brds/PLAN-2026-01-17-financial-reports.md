---
status: draft
created_at: 2026-01-17
updated_at: 2026-01-17
brd_ref: "brds/BRD-2026-01-16-financial-reports.md"
ux_ref: "brds/UX-2026-01-16-financial-reports.md"
owners:
  - area: portal-web
    agent: Feature Builder
risk_level: low
---

# Engineering Plan: Financial Reports — Business Health at a Glance

## 0) Inputs

- **BRD**: [BRD-2026-01-16-financial-reports.md](./BRD-2026-01-16-financial-reports.md)
- **UX Spec**: [UX-2026-01-16-financial-reports.md](./UX-2026-01-16-financial-reports.md)
- **Assumptions**:
  - Backend APIs already exist and are stable (no backend changes required)
  - Frontend-only implementation (portal-web)
  - Both admin and member roles can view reports (read-only)
  - As-of date parameter defaults to today

## 1) Confirmation Gate (must be approved before implementation)

- ❌ New dependency/library? **No** — uses existing Chart.js, TanStack Router/Query/Store
- ❌ New project/app? **No** — portal-web only
- ❌ Breaking change? **No** — additive feature only
- ❌ Migration? **No** — no data model changes
- ❌ Data model change with customer impact? **No** — read-only feature

**No confirmation-gated changes proposed.**

## 2) Architecture Summary (high level)

- **Backend**: No changes required. Uses existing analytics API endpoints:
  - `GET /v1/businesses/:businessDescriptor/analytics/reports/financial-position`
  - `GET /v1/businesses/:businessDescriptor/analytics/reports/profit-and-loss`
  - `GET /v1/businesses/:businessDescriptor/analytics/reports/cash-flow`
  - `GET /v1/businesses/:businessDescriptor/accounting/summary` (for Safe to Draw)

- **Portal-web**:
  - 4 new routes under `/business/$businessDescriptor/reports/`
  - 1 new feature module: `features/reports/`
  - 1 new shared component: `components/molecules/InsightCard.tsx`
  - 5 feature-specific components: `ReportCard`, `AssetBreakdownBar`, `ProfitWaterfall`, `CashFlowDiagram`, + page components
  - Navigation updates: Sidebar + (optional) BottomNav
  - New i18n namespace: `reports.json` (en/ar)

- **Data model**: No changes — read-only access to existing backend models

- **Security/tenancy**: Existing middleware chain enforces business scoping. No new security considerations.

## 3) Step-based Execution Plan (handoff-ready)

**Execution protocol:**
- Feature Builder will implement **one step per request**.
- Each step below is sized to be completed "perfectly" in a single AI request.
- Do not start Step N+1 before Step N is merged/verified.
- After each step: TypeScript must compile without errors, ESLint must pass.

### Step Index

- **Step 0** — i18n namespace setup (en/ar parity)
- **Step 1** — API types + query functions for financial reports
- **Step 2** — InsightCard shared component (molecules)
- **Step 3** — Feature folder structure + ReportCard component
- **Step 4** — Reports Hub page (route + feature page)
- **Step 5** — AssetBreakdownBar + Business Health page
- **Step 6** — ProfitWaterfall + Profit & Earnings page
- **Step 7** — CashFlowDiagram + Cash Movement page
- **Step 8** — Navigation updates (Sidebar) + final polish

---

### Step 0 — i18n Namespace Setup (en/ar parity)

- **Goal (user-visible outcome)**: All report-related translation keys are available in both English and Arabic.

- **Scope**:
  - **In**: Create `reports.json` namespace with all keys from BRD/UX spec
  - **Out**: No UI components yet

- **Targets (files/symbols/endpoints)**:
  - Create: `portal-web/src/i18n/en/reports.json`
  - Create: `portal-web/src/i18n/ar/reports.json`
  - Modify: `portal-web/src/i18n/en/dashboard.json` (add `reports` nav key)
  - Modify: `portal-web/src/i18n/ar/dashboard.json` (add `reports` nav key)
  - Modify: `portal-web/src/i18n/index.ts` (import and register `reports` namespace)

- **Tasks (detailed checklist)**:

  **Portal-web i18n:**
  - [ ] Create `portal-web/src/i18n/en/reports.json` with all keys:
    ```json
    {
      "hub": {
        "title": "Reports",
        "safe_to_draw": "Safe to Draw",
        "safe_to_draw_subtitle": "The amount you can safely take out",
        "as_of": "As of {{date}}"
      },
      "cards": {
        "business_health": "Business Health",
        "profit_earnings": "Profit & Earnings",
        "cash_movement": "Cash Movement",
        "view_details": "View Details"
      },
      "metrics": {
        "business_worth": "What Your Business Is Worth",
        "what_you_keep": "What You Keep",
        "cash_now": "Cash Now"
      },
      "health": {
        "title": "Business Health",
        "subtitle": "A snapshot of what your business owns and is worth",
        "what_you_own": "What You Own",
        "what_you_owe": "What You Owe",
        "owners_stake": "Owner's Stake",
        "cash_on_hand": "Cash on Hand",
        "inventory_value": "Inventory Value",
        "fixed_assets": "Equipment & Assets",
        "total_assets": "Total You Own",
        "total_liabilities": "Total You Owe",
        "money_put_in": "Money You Put In",
        "money_took_out": "Money You Took Out",
        "profit_kept": "Profit Kept in Business",
        "business_value": "Your Business Value",
        "liabilities_note": "Kyora doesn't track loans or credit yet. This will be available in a future update."
      },
      "profit": {
        "title": "Profit & Earnings",
        "subtitle": "Where your money comes from and where it goes",
        "what_you_keep": "What You Keep",
        "what_you_keep_subtitle": "Your final profit after all costs",
        "money_in": "Money Coming In",
        "from_sales": "From your sales",
        "product_costs": "Product Costs",
        "product_costs_explanation": "The cost of making or buying what you sold",
        "gross_profit": "Gross Profit",
        "gross_profit_explanation": "What's left after product costs",
        "gross_margin": "Gross Margin",
        "running_costs": "Running Costs",
        "total_running_costs": "Total Running Costs",
        "final_profit": "Final Profit",
        "profit_margin": "Profit Margin"
      },
      "cashflow": {
        "title": "Cash Movement",
        "subtitle": "How cash flows through your business",
        "cash_now": "Cash Now",
        "cash_runway": "Your current cash runway",
        "cash_start": "Cash at Start",
        "cash_end": "Cash at End",
        "money_in": "Money Coming In",
        "money_out": "Money Going Out",
        "from_customers": "From Customers",
        "from_owner": "From Owner",
        "total_in": "Total Cash In",
        "inventory_purchases": "Inventory Purchases",
        "running_costs": "Running Costs",
        "equipment_assets": "Equipment & Assets",
        "to_owner": "To Owner",
        "total_out": "Total Cash Out",
        "net_change": "Net Change",
        "cash_increased": "Your cash increased by {{amount}}",
        "cash_decreased": "Your cash decreased by {{amount}}"
      },
      "insights": {
        "healthy": "Your business is in a healthy position. You own more than you owe.",
        "negative_equity": "Your business value is negative. This means expenses and withdrawals exceed profits.",
        "negative_cash": "Your cash position is estimated as negative. Review recent expenses or add capital.",
        "profitable": "You're profitable! For every {{currency}} of sales, you keep {{margin}}%.",
        "losing": "You're currently losing money. Your expenses exceed your gross profit by {{amount}}.",
        "biggest_expense": "Your biggest expense is {{category}} at {{amount}} ({{percent}}%).",
        "cash_healthy": "Cash Healthy: More cash came in than went out.",
        "cash_alert": "Cash Alert: You spent more than you received. Monitor your cash carefully.",
        "withdrawal_tip": "Tip: You withdrew more than your sales revenue. Consider aligning withdrawals with income."
      },
      "quick_links": {
        "view_assets": "View Assets",
        "view_capital": "View Capital",
        "view_expenses": "View All Expenses",
        "view_orders": "View Orders"
      },
      "empty": {
        "title": "Your financial picture starts here",
        "body": "Once you start recording orders and expenses, you'll see a complete view of your business finances.",
        "cta": "Create Your First Order"
      },
      "error": {
        "title": "Couldn't load your reports",
        "body": "We're having trouble connecting. Please check your internet and try again."
      }
    }
    ```
  - [ ] Create `portal-web/src/i18n/ar/reports.json` with Arabic translations (all keys must have parity):
    ```json
    {
      "hub": {
        "title": "التقارير",
        "safe_to_draw": "المبلغ الآمن للسحب",
        "safe_to_draw_subtitle": "المبلغ الذي يمكنك سحبه بأمان",
        "as_of": "حتى تاريخ {{date}}"
      },
      "cards": {
        "business_health": "صحة عملك",
        "profit_earnings": "أرباحك",
        "cash_movement": "حركة النقد",
        "view_details": "عرض التفاصيل"
      },
      "metrics": {
        "business_worth": "قيمة عملك",
        "what_you_keep": "ما تحتفظ به",
        "cash_now": "النقد الحالي"
      },
      "health": {
        "title": "صحة عملك",
        "subtitle": "لمحة عما يملكه عملك وقيمته",
        "what_you_own": "ما تملكه",
        "what_you_owe": "ما عليك",
        "owners_stake": "حصة المالك",
        "cash_on_hand": "النقد المتوفر",
        "inventory_value": "قيمة المخزون",
        "fixed_assets": "المعدات والأصول",
        "total_assets": "إجمالي ما تملكه",
        "total_liabilities": "إجمالي ما عليك",
        "money_put_in": "المال الذي وضعته",
        "money_took_out": "المال الذي سحبته",
        "profit_kept": "الأرباح المحتفظ بها",
        "business_value": "قيمة عملك",
        "liabilities_note": "كيورا لا تتتبع القروض أو الائتمان بعد. ستتوفر هذه الميزة في تحديث مستقبلي."
      },
      "profit": {
        "title": "أرباحك",
        "subtitle": "من أين يأتي مالك وأين يذهب",
        "what_you_keep": "ما تحتفظ به",
        "what_you_keep_subtitle": "ربحك النهائي بعد كل التكاليف",
        "money_in": "الأموال الواردة",
        "from_sales": "من مبيعاتك",
        "product_costs": "تكلفة المنتجات",
        "product_costs_explanation": "تكلفة صنع أو شراء ما بعته",
        "gross_profit": "الربح الإجمالي",
        "gross_profit_explanation": "ما يتبقى بعد تكلفة المنتجات",
        "gross_margin": "هامش الربح الإجمالي",
        "running_costs": "مصاريف التشغيل",
        "total_running_costs": "إجمالي مصاريف التشغيل",
        "final_profit": "الربح النهائي",
        "profit_margin": "هامش الربح"
      },
      "cashflow": {
        "title": "حركة النقد",
        "subtitle": "كيف ينتقل النقد في عملك",
        "cash_now": "النقد الحالي",
        "cash_runway": "المدى النقدي الحالي",
        "cash_start": "النقد في البداية",
        "cash_end": "النقد في النهاية",
        "money_in": "الأموال الواردة",
        "money_out": "الأموال الصادرة",
        "from_customers": "من العملاء",
        "from_owner": "من المالك",
        "total_in": "إجمالي النقد الوارد",
        "inventory_purchases": "مشتريات المخزون",
        "running_costs": "مصاريف التشغيل",
        "equipment_assets": "المعدات والأصول",
        "to_owner": "للمالك",
        "total_out": "إجمالي النقد الصادر",
        "net_change": "صافي التغير",
        "cash_increased": "ازداد نقدك بمقدار {{amount}}",
        "cash_decreased": "انخفض نقدك بمقدار {{amount}}"
      },
      "insights": {
        "healthy": "عملك في وضع صحي. أنت تملك أكثر مما عليك.",
        "negative_equity": "قيمة عملك سالبة. هذا يعني أن المصروفات والسحوبات تتجاوز الأرباح.",
        "negative_cash": "وضعك النقدي المقدر سالب. راجع المصروفات الأخيرة أو أضف رأس مال.",
        "profitable": "أنت تحقق ربحاً! مقابل كل {{currency}} من المبيعات، تحتفظ بـ {{margin}}%.",
        "losing": "أنت تخسر حالياً. مصروفاتك تتجاوز ربحك الإجمالي بمقدار {{amount}}.",
        "biggest_expense": "أكبر مصروفاتك هي {{category}} بمبلغ {{amount}} ({{percent}}%).",
        "cash_healthy": "النقد صحي: دخل نقد أكثر مما خرج.",
        "cash_alert": "تنبيه نقدي: أنفقت أكثر مما استلمت. راقب نقدك بعناية.",
        "withdrawal_tip": "نصيحة: سحبت أكثر من إيرادات مبيعاتك. فكر في مواءمة السحوبات مع الدخل."
      },
      "quick_links": {
        "view_assets": "عرض الأصول",
        "view_capital": "عرض رأس المال",
        "view_expenses": "عرض جميع المصروفات",
        "view_orders": "عرض الطلبات"
      },
      "empty": {
        "title": "صورتك المالية تبدأ هنا",
        "body": "بمجرد أن تبدأ بتسجيل الطلبات والمصروفات، ستظهر لك صورة كاملة عن ماليات عملك.",
        "cta": "أنشئ طلبك الأول"
      },
      "error": {
        "title": "تعذر تحميل تقاريرك",
        "body": "نواجه مشكلة في الاتصال. يرجى التحقق من الإنترنت والمحاولة مجدداً."
      }
    }
    ```
  - [ ] Update `portal-web/src/i18n/en/dashboard.json` — add nav key:
    ```json
    {
      "reports": "Reports"
    }
    ```
  - [ ] Update `portal-web/src/i18n/ar/dashboard.json` — add nav key:
    ```json
    {
      "reports": "التقارير"
    }
    ```
  - [ ] Update `portal-web/src/i18n/index.ts` — import and register `reports` namespace in both `enResources` and `arResources` objects

- **Edge cases + error handling**: N/A for i18n setup

- **Verification checklist**:
  - [ ] Run `pnpm tsc --noEmit` — no TypeScript errors
  - [ ] Run `pnpm lint` — no ESLint errors
  - [ ] Verify all keys exist in both `en/reports.json` and `ar/reports.json`
  - [ ] Verify `reports` namespace is imported in i18n index

- **Definition of done**:
  - [ ] TypeScript compiles without errors
  - [ ] ESLint passes without errors
  - [ ] Both en/ar i18n files exist with identical key structure
  - [ ] Namespace is registered and loadable

- **Relevant instruction files**:
  - `.github/instructions/i18n-translations.instructions.md` — namespace rules, key naming, parity requirements

---

### Step 1 — API Types + Query Functions for Financial Reports

- **Goal (user-visible outcome)**: API layer is ready to fetch financial report data from backend.

- **Scope**:
  - **In**: Add types for financial reports, query functions, query keys
  - **Out**: No UI yet

- **Targets (files/symbols/endpoints)**:
  - Create: `portal-web/src/api/types/analytics.ts` (if not exists, or extend)
  - Modify: `portal-web/src/api/accounting.ts` (add report query functions)
  - Modify: `portal-web/src/lib/queryKeys.ts` (add `reports` query key factory)

- **Tasks (detailed checklist)**:

  **Portal-web API types:**
  - [ ] Create/extend `portal-web/src/api/types/analytics.ts` with backend response types:
    ```typescript
    // Financial Position (Balance Sheet)
    export interface FinancialPosition {
      businessID: string
      asOf: string
      totalAssets: string      // decimal as string
      totalLiabilities: string
      totalEquity: string
      cashOnHand: string
      totalInventoryValue: string
      currentAssets: string
      fixedAssets: string
      ownerInvestment: string
      retainedEarnings: string
      ownerDraws: string
    }

    // Profit and Loss Statement
    export interface ProfitAndLossStatement {
      businessID: string
      asOf: string
      grossProfit: string
      totalExpenses: string
      netProfit: string
      revenue: string
      cogs: string
      expensesByCategory: Array<{ key: string; value: string }>
    }

    // Cash Flow Statement
    export interface CashFlowStatement {
      businessID: string
      asOf: string
      cashAtStart: string
      cashAtEnd: string
      cashFromCustomers: string
      cashFromOwner: string
      totalCashIn: string
      inventoryPurchases: string
      operatingExpenses: string
      totalBusinessOperation: string
      businessInvestments: string
      ownerDraws: string
      totalCashOut: string
      netCashFlow: string
    }
    ```

  **Portal-web API functions:**
  - [ ] Add to `portal-web/src/api/accounting.ts`:
    ```typescript
    // Query options for financial reports
    export const financialPositionQueryOptions = (
      businessDescriptor: string,
      asOf?: string
    ) => queryOptions({
      queryKey: reports.financialPosition(businessDescriptor, asOf),
      queryFn: () => get(`v1/businesses/${businessDescriptor}/analytics/reports/financial-position`, {
        searchParams: asOf ? { asOf } : undefined
      }).json<FinancialPosition>(),
      staleTime: STALE_TIME.FIVE_MINUTES,
    })

    export const profitAndLossQueryOptions = (
      businessDescriptor: string,
      asOf?: string
    ) => queryOptions({
      queryKey: reports.profitAndLoss(businessDescriptor, asOf),
      queryFn: () => get(`v1/businesses/${businessDescriptor}/analytics/reports/profit-and-loss`, {
        searchParams: asOf ? { asOf } : undefined
      }).json<ProfitAndLossStatement>(),
      staleTime: STALE_TIME.FIVE_MINUTES,
    })

    export const cashFlowQueryOptions = (
      businessDescriptor: string,
      asOf?: string
    ) => queryOptions({
      queryKey: reports.cashFlow(businessDescriptor, asOf),
      queryFn: () => get(`v1/businesses/${businessDescriptor}/analytics/reports/cash-flow`, {
        searchParams: asOf ? { asOf } : undefined
      }).json<CashFlowStatement>(),
      staleTime: STALE_TIME.FIVE_MINUTES,
    })

    // Hook wrappers
    export function useFinancialPositionQuery(businessDescriptor: string, asOf?: string) {
      return useQuery(financialPositionQueryOptions(businessDescriptor, asOf))
    }

    export function useProfitAndLossQuery(businessDescriptor: string, asOf?: string) {
      return useQuery(profitAndLossQueryOptions(businessDescriptor, asOf))
    }

    export function useCashFlowQuery(businessDescriptor: string, asOf?: string) {
      return useQuery(cashFlowQueryOptions(businessDescriptor, asOf))
    }
    ```

  **Portal-web query keys:**
  - [ ] Add to `portal-web/src/lib/queryKeys.ts`:
    ```typescript
    /**
     * Reports queries (business-scoped)
     * StaleTime: 5 minutes (computed data, acceptable staleness)
     * Invalidated on business switch
     */
    export const reports = {
      all: ['reports'] as const,
      businessScoped: true,
      financialPosition: (businessDescriptor: string, asOf?: string) =>
        [...reports.all, 'financial-position', businessDescriptor, asOf] as const,
      profitAndLoss: (businessDescriptor: string, asOf?: string) =>
        [...reports.all, 'profit-and-loss', businessDescriptor, asOf] as const,
      cashFlow: (businessDescriptor: string, asOf?: string) =>
        [...reports.all, 'cash-flow', businessDescriptor, asOf] as const,
    } as const
    ```

- **Edge cases + error handling**:
  - API errors are handled by global error handling (see `http-tanstack-query.instructions.md`)
  - `asOf` parameter is optional; backend defaults to today

- **Verification checklist**:
  - [ ] Run `pnpm tsc --noEmit` — no TypeScript errors
  - [ ] Run `pnpm lint` — no ESLint errors
  - [ ] Types match backend response structure in `backend/internal/domain/analytics/model.go`
  - [ ] Query keys follow existing pattern in `queryKeys.ts`

- **Definition of done**:
  - [ ] TypeScript compiles without errors
  - [ ] ESLint passes without errors
  - [ ] All three report query functions are exported
  - [ ] Query keys are properly structured for invalidation

- **Relevant instruction files**:
  - `.github/instructions/http-tanstack-query.instructions.md` — query options, hooks, error handling
  - `.github/instructions/ky.instructions.md` — HTTP client usage
  - `.github/instructions/analytics.instructions.md` — backend API contract reference

---

### Step 2 — InsightCard Shared Component (molecules)

- **Goal (user-visible outcome)**: Reusable InsightCard component available for all report pages.

- **Scope**:
  - **In**: Create `InsightCard` in `components/molecules/`
  - **Out**: No route integration yet

- **Targets (files/symbols/endpoints)**:
  - Create: `portal-web/src/components/molecules/InsightCard.tsx`
  - Modify: `portal-web/src/components/index.ts` (export InsightCard)

- **Tasks (detailed checklist)**:

  **Portal-web InsightCard component:**
  - [ ] Create `portal-web/src/components/molecules/InsightCard.tsx`:
    ```typescript
    /**
     * InsightCard Component
     *
     * Displays contextual insights with icon, variant styling, and optional action link.
     * Used in financial reports to show actionable suggestions based on metrics.
     *
     * Variants:
     * - success: Green background, CheckCircle icon (positive metrics)
     * - warning: Yellow background, AlertTriangle icon (concerning metrics)
     * - error: Red background, AlertCircle icon (critical issues)
     * - info: Blue background, Info icon (tips and information)
     */
    import { Link } from '@tanstack/react-router'
    import {
      AlertCircle,
      AlertTriangle,
      CheckCircle,
      ChevronLeft,
      ChevronRight,
      Info,
      type LucideIcon,
    } from 'lucide-react'
    import { cn } from '@/lib/utils'
    import { useLanguage } from '@/hooks/useLanguage'

    export interface InsightCardProps {
      variant: 'success' | 'warning' | 'error' | 'info'
      message: string
      icon?: LucideIcon
      link?: {
        label: string
        href: string
        params?: Record<string, string>
      }
      className?: string
    }

    const variantConfig = {
      success: {
        bg: 'bg-success/10',
        border: 'border-success/20',
        text: 'text-success',
        defaultIcon: CheckCircle,
      },
      warning: {
        bg: 'bg-warning/10',
        border: 'border-warning/20',
        text: 'text-warning',
        defaultIcon: AlertTriangle,
      },
      error: {
        bg: 'bg-error/10',
        border: 'border-error/20',
        text: 'text-error',
        defaultIcon: AlertCircle,
      },
      info: {
        bg: 'bg-info/10',
        border: 'border-info/20',
        text: 'text-info',
        defaultIcon: Info,
      },
    } as const

    export function InsightCard({
      variant,
      message,
      icon,
      link,
      className,
    }: InsightCardProps) {
      const { isRTL } = useLanguage()
      const config = variantConfig[variant]
      const Icon = icon ?? config.defaultIcon
      const ChevronIcon = isRTL ? ChevronLeft : ChevronRight

      return (
        <div
          className={cn(
            'rounded-box border p-4',
            config.bg,
            config.border,
            className
          )}
          role="note"
          aria-label={message}
        >
          <div className="flex items-start gap-3">
            <Icon className={cn('h-5 w-5 shrink-0 mt-0.5', config.text)} />
            <div className="flex-1 space-y-2">
              <p className="text-sm text-base-content">{message}</p>
              {link && (
                <Link
                  to={link.href}
                  params={link.params}
                  className={cn(
                    'inline-flex items-center gap-1 text-sm font-medium',
                    config.text,
                    'hover:underline'
                  )}
                >
                  {link.label}
                  <ChevronIcon className="h-4 w-4" />
                </Link>
              )}
            </div>
          </div>
        </div>
      )
    }
    ```

  - [ ] Update `portal-web/src/components/index.ts` — add export:
    ```typescript
    export { InsightCard } from './molecules/InsightCard'
    export type { InsightCardProps } from './molecules/InsightCard'
    ```

- **Edge cases + error handling**:
  - RTL: ChevronRight becomes ChevronLeft for link arrows
  - Link is optional; component works without it
  - Icon is optional; defaults to variant-appropriate icon

- **Verification checklist**:
  - [ ] Run `pnpm tsc --noEmit` — no TypeScript errors
  - [ ] Run `pnpm lint` — no ESLint errors
  - [ ] Component uses logical properties (no `ml-`/`mr-`, only `ms-`/`me-`)
  - [ ] Component is exported from index.ts

- **Definition of done**:
  - [ ] TypeScript compiles without errors
  - [ ] ESLint passes without errors
  - [ ] InsightCard renders correctly in all 4 variants
  - [ ] RTL directional icons work correctly

- **Relevant instruction files**:
  - `.github/instructions/ui-implementation.instructions.md` — RTL rules, daisyUI classes
  - `.github/instructions/design-tokens.instructions.md` — color tokens, spacing
  - `.github/instructions/portal-web-code-structure.instructions.md` — molecules placement

---

### Step 3 — Feature Folder Structure + ReportCard Component

- **Goal (user-visible outcome)**: Reports feature folder exists with ReportCard component ready for hub page.

- **Scope**:
  - **In**: Create feature folder, ReportCard, ReportCardSkeleton
  - **Out**: No routes yet

- **Targets (files/symbols/endpoints)**:
  - Create: `portal-web/src/features/reports/` folder structure
  - Create: `portal-web/src/features/reports/components/ReportCard.tsx`
  - Create: `portal-web/src/features/reports/components/index.ts`

- **Tasks (detailed checklist)**:

  **Portal-web feature structure:**
  - [ ] Create folder: `portal-web/src/features/reports/`
  - [ ] Create folder: `portal-web/src/features/reports/components/`
  - [ ] Create `portal-web/src/features/reports/components/ReportCard.tsx`:
    ```typescript
    /**
     * ReportCard Component
     *
     * Clickable card for the Reports Hub page displaying:
     * - Icon and title
     * - Key metric (large, prominent)
     * - Secondary metrics (smaller, supporting)
     * - "View Details" CTA
     *
     * Entire card is clickable with proper touch targets (min 44px).
     */
    import { Link } from '@tanstack/react-router'
    import { ChevronLeft, ChevronRight, type LucideIcon } from 'lucide-react'
    import { useTranslation } from 'react-i18next'
    import { cn } from '@/lib/utils'
    import { useLanguage } from '@/hooks/useLanguage'
    import { Skeleton } from '@/components/atoms/Skeleton'

    export interface ReportCardProps {
      title: string
      icon: LucideIcon
      keyMetric: {
        label: string
        value: string
      }
      secondaryValues?: Array<{
        label: string
        value: string
      }>
      href: string
      params?: Record<string, string>
      searchParams?: Record<string, string>
      className?: string
    }

    export function ReportCard({
      title,
      icon: Icon,
      keyMetric,
      secondaryValues = [],
      href,
      params,
      searchParams,
      className,
    }: ReportCardProps) {
      const { t } = useTranslation('reports')
      const { isRTL } = useLanguage()
      const ChevronIcon = isRTL ? ChevronLeft : ChevronRight

      return (
        <Link
          to={href}
          params={params}
          search={searchParams}
          className={cn(
            'block rounded-box border border-base-300 p-4',
            'transition-colors hover:bg-base-200/50 active:scale-[0.98]',
            'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary',
            className
          )}
        >
          {/* Header: Icon + Title */}
          <div className="flex items-center gap-2 mb-3">
            <Icon className="h-5 w-5 text-primary" />
            <h3 className="text-base font-semibold text-base-content">{title}</h3>
          </div>

          {/* Key Metric */}
          <div className="mb-2">
            <p className="text-xs text-base-content/60 mb-1">{keyMetric.label}</p>
            <p className="text-2xl font-bold text-base-content">{keyMetric.value}</p>
          </div>

          {/* Secondary Values */}
          {secondaryValues.length > 0 && (
            <div className="flex flex-wrap gap-x-4 gap-y-1 mb-3 text-sm text-base-content/70">
              {secondaryValues.map((item, index) => (
                <span key={index}>
                  {item.label}: <span className="font-medium">{item.value}</span>
                </span>
              ))}
            </div>
          )}

          {/* CTA */}
          <div className="flex items-center justify-end gap-1 text-sm font-medium text-primary">
            {t('cards.view_details')}
            <ChevronIcon className="h-4 w-4" />
          </div>
        </Link>
      )
    }

    export function ReportCardSkeleton({ className }: { className?: string }) {
      return (
        <div
          className={cn(
            'rounded-box border border-base-300 p-4',
            className
          )}
        >
          {/* Header skeleton */}
          <div className="flex items-center gap-2 mb-3">
            <Skeleton className="h-5 w-5 rounded" />
            <Skeleton className="h-5 w-32" />
          </div>

          {/* Key metric skeleton */}
          <div className="mb-2">
            <Skeleton className="h-3 w-24 mb-2" />
            <Skeleton className="h-8 w-40" />
          </div>

          {/* Secondary values skeleton */}
          <div className="flex gap-4 mb-3">
            <Skeleton className="h-4 w-20" />
            <Skeleton className="h-4 w-24" />
          </div>

          {/* CTA skeleton */}
          <div className="flex justify-end">
            <Skeleton className="h-4 w-24" />
          </div>
        </div>
      )
    }
    ```

  - [ ] Create `portal-web/src/features/reports/components/index.ts`:
    ```typescript
    export { ReportCard, ReportCardSkeleton } from './ReportCard'
    export type { ReportCardProps } from './ReportCard'
    ```

- **Edge cases + error handling**:
  - RTL: Chevron direction flips
  - Secondary values are optional
  - Entire card is clickable (Link wrapper)

- **Verification checklist**:
  - [ ] Run `pnpm tsc --noEmit` — no TypeScript errors
  - [ ] Run `pnpm lint` — no ESLint errors
  - [ ] Component uses logical properties for RTL
  - [ ] Touch target is adequate (card is fully tappable)

- **Definition of done**:
  - [ ] TypeScript compiles without errors
  - [ ] ESLint passes without errors
  - [ ] Feature folder structure matches code-structure.instructions.md
  - [ ] ReportCard and ReportCardSkeleton are exported

- **Relevant instruction files**:
  - `.github/instructions/portal-web-code-structure.instructions.md` — feature folder placement
  - `.github/instructions/ui-implementation.instructions.md` — RTL, touch targets, daisyUI

---

### Step 4 — Reports Hub Page (route + feature page)

- **Goal (user-visible outcome)**: Users can navigate to `/business/:descriptor/reports/` and see the Reports Hub with Safe to Draw hero and 3 report cards.

- **Scope**:
  - **In**: Route file, Hub page component, wire up data fetching
  - **Out**: Detail pages (Health, Profit, CashFlow) are separate steps

- **Targets (files/symbols/endpoints)**:
  - Create: `portal-web/src/routes/business/$businessDescriptor/reports/index.tsx`
  - Create: `portal-web/src/features/reports/components/ReportsHubPage.tsx`
  - Modify: `portal-web/src/features/reports/components/index.ts`

- **Tasks (detailed checklist)**:

  **Portal-web route:**
  - [ ] Create `portal-web/src/routes/business/$businessDescriptor/reports/index.tsx`:
    ```typescript
    /**
     * Reports Hub Route
     *
     * Landing page for financial reports showing:
     * - Safe to Draw hero metric
     * - 3 report cards (Business Health, Profit & Earnings, Cash Movement)
     *
     * URL: /business/:businessDescriptor/reports
     * Search params: asOf (optional date filter, YYYY-MM-DD)
     */
    import { createFileRoute } from '@tanstack/react-router'
    import { z } from 'zod'
    import {
      accountingSummaryQueryOptions,
      financialPositionQueryOptions,
      profitAndLossQueryOptions,
      cashFlowQueryOptions,
    } from '@/api/accounting'
    import { ReportsHubPage } from '@/features/reports/components/ReportsHubPage'

    const searchSchema = z.object({
      asOf: z.string().optional(),
    })

    export const Route = createFileRoute(
      '/business/$businessDescriptor/reports/'
    )({
      validateSearch: searchSchema,
      staticData: {
        titleKey: 'common.pages.reports',
      },
      loader: async ({ context, params, search }) => {
        const { businessDescriptor } = params
        const { asOf } = search

        // Prefetch all required data in parallel
        await Promise.all([
          context.queryClient.ensureQueryData(
            accountingSummaryQueryOptions(businessDescriptor)
          ),
          context.queryClient.ensureQueryData(
            financialPositionQueryOptions(businessDescriptor, asOf)
          ),
          context.queryClient.ensureQueryData(
            profitAndLossQueryOptions(businessDescriptor, asOf)
          ),
          context.queryClient.ensureQueryData(
            cashFlowQueryOptions(businessDescriptor, asOf)
          ),
        ])
      },
      component: ReportsHubPage,
    })
    ```

  **Portal-web feature page:**
  - [ ] Create `portal-web/src/features/reports/components/ReportsHubPage.tsx`:
    ```typescript
    /**
     * Reports Hub Page
     *
     * Displays:
     * - Page header with "As of" date button
     * - Safe to Draw hero stat (from accounting summary)
     * - 3 ReportCards linking to detail pages
     * - Empty/loading/error states
     */
    import { useNavigate, useParams, useSearch, useRouteContext } from '@tanstack/react-router'
    import { useTranslation } from 'react-i18next'
    import {
      Activity,
      ArrowDownRight,
      ArrowUpRight,
      BarChart3,
      Calendar,
      PiggyBank,
      TrendingDown,
      TrendingUp,
      Wallet,
    } from 'lucide-react'

    import {
      useAccountingSummaryQuery,
      useFinancialPositionQuery,
      useProfitAndLossQuery,
      useCashFlowQuery,
    } from '@/api/accounting'
    import { StatCard, StatCardSkeleton, Button } from '@/components'
    import { formatCurrency } from '@/lib/formatCurrency'
    import { formatDateShort } from '@/lib/formatDate'
    import { ReportCard, ReportCardSkeleton } from './ReportCard'

    export function ReportsHubPage() {
      const { t } = useTranslation('reports')
      const { t: tCommon } = useTranslation('common')
      const navigate = useNavigate()
      const { businessDescriptor } = useParams({
        from: '/business/$businessDescriptor/reports/',
      })
      const { asOf } = useSearch({
        from: '/business/$businessDescriptor/reports/',
      })
      const { business } = useRouteContext({
        from: '/business/$businessDescriptor',
      })

      // Data queries
      const { data: summary, isLoading: isSummaryLoading, error: summaryError } =
        useAccountingSummaryQuery(businessDescriptor)

      const { data: financialPosition, isLoading: isFPLoading } =
        useFinancialPositionQuery(businessDescriptor, asOf)

      const { data: profitLoss, isLoading: isPLLoading } =
        useProfitAndLossQuery(businessDescriptor, asOf)

      const { data: cashFlow, isLoading: isCFLoading } =
        useCashFlowQuery(businessDescriptor, asOf)

      const isLoading = isSummaryLoading || isFPLoading || isPLLoading || isCFLoading
      const currency = summary?.currency ?? business.currency

      // Parse amounts
      const safeToDrawAmount = parseFloat(summary?.safeToDrawAmount ?? '0')
      const isSafeToDrawNegative = safeToDrawAmount < 0

      // Financial Position metrics
      const totalEquity = parseFloat(financialPosition?.totalEquity ?? '0')
      const cashOnHand = parseFloat(financialPosition?.cashOnHand ?? '0')
      const inventoryValue = parseFloat(financialPosition?.totalInventoryValue ?? '0')

      // Profit metrics
      const netProfit = parseFloat(profitLoss?.netProfit ?? '0')
      const revenue = parseFloat(profitLoss?.revenue ?? '0')
      const totalExpenses = parseFloat(profitLoss?.totalExpenses ?? '0')

      // Cash flow metrics
      const cashAtEnd = parseFloat(cashFlow?.cashAtEnd ?? '0')
      const totalCashIn = parseFloat(cashFlow?.totalCashIn ?? '0')
      const totalCashOut = parseFloat(cashFlow?.totalCashOut ?? '0')

      // Date display
      const displayDate = asOf ? new Date(asOf) : new Date()

      // TODO: Implement date picker sheet in future iteration
      const handleDateClick = () => {
        // For now, do nothing. Date picker will be added later.
      }

      // Check if business has any data
      const hasData = revenue > 0 || totalExpenses > 0 || cashOnHand !== 0

      // Error state
      if (summaryError) {
        return (
          <div className="flex flex-col items-center justify-center py-16 px-4 text-center">
            <div className="text-error mb-4">
              <Activity className="h-12 w-12 mx-auto" />
            </div>
            <h2 className="text-lg font-semibold mb-2">{t('error.title')}</h2>
            <p className="text-base-content/60 mb-4 max-w-md">{t('error.body')}</p>
            <Button
              variant="ghost"
              onClick={() => window.location.reload()}
            >
              {tCommon('retry')}
            </Button>
          </div>
        )
      }

      // Empty state (new business, no data)
      if (!isLoading && !hasData) {
        return (
          <div className="flex flex-col items-center justify-center py-16 px-4 text-center">
            <div className="text-primary/40 mb-4">
              <BarChart3 className="h-16 w-16 mx-auto" />
            </div>
            <h2 className="text-lg font-semibold mb-2">{t('empty.title')}</h2>
            <p className="text-base-content/60 mb-6 max-w-md">{t('empty.body')}</p>
            <Button
              as="link"
              to="/business/$businessDescriptor/orders"
              params={{ businessDescriptor }}
            >
              {t('empty.cta')}
            </Button>
          </div>
        )
      }

      return (
        <div className="space-y-6 pb-20 md:pb-6">
          {/* Header with As-of Date */}
          <div className="flex items-center justify-between">
            <h1 className="text-xl font-bold">{t('hub.title')}</h1>
            <button
              onClick={handleDateClick}
              className="flex items-center gap-2 text-sm text-base-content/70 hover:text-base-content"
            >
              <Calendar className="h-4 w-4" />
              <span>{t('hub.as_of', { date: formatDateShort(displayDate) })}</span>
            </button>
          </div>

          {/* Safe to Draw Hero */}
          <section aria-labelledby="safe-to-draw-heading">
            <h2 id="safe-to-draw-heading" className="sr-only">
              {t('hub.safe_to_draw')}
            </h2>
            {isLoading ? (
              <StatCardSkeleton />
            ) : (
              <StatCard
                label={t('hub.safe_to_draw')}
                value={formatCurrency(Math.abs(safeToDrawAmount), currency)}
                icon={<PiggyBank className="h-5 w-5 text-primary" />}
                variant={isSafeToDrawNegative ? 'error' : 'success'}
                trend={isSafeToDrawNegative ? 'down' : undefined}
                trendValue={
                  isSafeToDrawNegative
                    ? t('hub.safe_to_draw_subtitle')
                    : t('hub.safe_to_draw_subtitle')
                }
              />
            )}
          </section>

          {/* Report Cards Grid */}
          <section aria-labelledby="reports-heading">
            <h2 id="reports-heading" className="sr-only">
              {t('hub.title')}
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {/* Business Health Card */}
              {isLoading ? (
                <ReportCardSkeleton />
              ) : (
                <ReportCard
                  title={t('cards.business_health')}
                  icon={Activity}
                  keyMetric={{
                    label: t('metrics.business_worth'),
                    value: formatCurrency(totalEquity, currency),
                  }}
                  secondaryValues={[
                    { label: t('health.cash_on_hand'), value: formatCurrency(cashOnHand, currency) },
                    { label: t('health.inventory_value'), value: formatCurrency(inventoryValue, currency) },
                  ]}
                  href="/business/$businessDescriptor/reports/health"
                  params={{ businessDescriptor }}
                  searchParams={asOf ? { asOf } : undefined}
                />
              )}

              {/* Profit & Earnings Card */}
              {isLoading ? (
                <ReportCardSkeleton />
              ) : (
                <ReportCard
                  title={t('cards.profit_earnings')}
                  icon={netProfit >= 0 ? TrendingUp : TrendingDown}
                  keyMetric={{
                    label: t('metrics.what_you_keep'),
                    value: formatCurrency(netProfit, currency),
                  }}
                  secondaryValues={[
                    { label: t('profit.money_in'), value: formatCurrency(revenue, currency) },
                    { label: t('profit.running_costs'), value: formatCurrency(totalExpenses, currency) },
                  ]}
                  href="/business/$businessDescriptor/reports/profit"
                  params={{ businessDescriptor }}
                  searchParams={asOf ? { asOf } : undefined}
                />
              )}

              {/* Cash Movement Card */}
              {isLoading ? (
                <ReportCardSkeleton />
              ) : (
                <ReportCard
                  title={t('cards.cash_movement')}
                  icon={Wallet}
                  keyMetric={{
                    label: t('metrics.cash_now'),
                    value: formatCurrency(cashAtEnd, currency),
                  }}
                  secondaryValues={[
                    { label: t('cashflow.total_in'), value: formatCurrency(totalCashIn, currency) },
                    { label: t('cashflow.total_out'), value: formatCurrency(totalCashOut, currency) },
                  ]}
                  href="/business/$businessDescriptor/reports/cashflow"
                  params={{ businessDescriptor }}
                  searchParams={asOf ? { asOf } : undefined}
                />
              )}
            </div>
          </section>
        </div>
      )
    }
    ```

  - [ ] Update `portal-web/src/features/reports/components/index.ts`:
    ```typescript
    export { ReportCard, ReportCardSkeleton } from './ReportCard'
    export type { ReportCardProps } from './ReportCard'
    export { ReportsHubPage } from './ReportsHubPage'
    ```

  - [ ] Add page title key to `portal-web/src/i18n/en/common.json` under `pages`:
    ```json
    "reports": "Reports"
    ```

  - [ ] Add page title key to `portal-web/src/i18n/ar/common.json` under `pages`:
    ```json
    "reports": "التقارير"
    ```

- **Edge cases + error handling**:
  - Empty state: shown when no revenue/expenses/cash data exists
  - Error state: shown when API fetch fails, with retry button
  - Loading state: skeleton cards displayed
  - Negative Safe to Draw: red variant with warning styling

- **Verification checklist**:
  - [ ] Run `pnpm tsc --noEmit` — no TypeScript errors
  - [ ] Run `pnpm lint` — no ESLint errors
  - [ ] Route is accessible at `/business/:descriptor/reports/`
  - [ ] Page displays Safe to Draw hero with correct color variant
  - [ ] 3 report cards display with correct metrics
  - [ ] Empty state renders for new businesses
  - [ ] Loading skeletons match final layout

- **Definition of done**:
  - [ ] TypeScript compiles without errors
  - [ ] ESLint passes without errors
  - [ ] Hub page renders correctly with all states (loading, empty, error, data)
  - [ ] Route file follows portal-web-code-structure.instructions.md

- **Relevant instruction files**:
  - `.github/instructions/portal-web-code-structure.instructions.md` — route vs feature page separation
  - `.github/instructions/http-tanstack-query.instructions.md` — loader prefetching, query usage
  - `.github/instructions/state-management.instructions.md` — URL-driven state (asOf param)
  - `.github/instructions/ui-implementation.instructions.md` — loading/empty/error states

---

### Step 5 — AssetBreakdownBar + Business Health Page

- **Goal (user-visible outcome)**: Users can view the Business Health page showing what they own, owe, and their business worth.

- **Scope**:
  - **In**: AssetBreakdownBar component, Business Health page, route
  - **Out**: Other detail pages

- **Targets (files/symbols/endpoints)**:
  - Create: `portal-web/src/features/reports/components/AssetBreakdownBar.tsx`
  - Create: `portal-web/src/features/reports/components/BusinessHealthPage.tsx`
  - Create: `portal-web/src/routes/business/$businessDescriptor/reports/health.tsx`
  - Modify: `portal-web/src/features/reports/components/index.ts`

- **Tasks (detailed checklist)**:

  **Portal-web AssetBreakdownBar:**
  - [ ] Create `portal-web/src/features/reports/components/AssetBreakdownBar.tsx`:
    ```typescript
    /**
     * AssetBreakdownBar Component
     *
     * Horizontal stacked bar showing asset distribution with legend.
     * Used on Business Health page to visualize What You Own breakdown.
     */
    import { cn } from '@/lib/utils'
    import { formatCurrency } from '@/lib/formatCurrency'

    export interface AssetSegment {
      label: string
      value: number
      color: 'success' | 'info' | 'secondary' | 'warning' | 'primary'
    }

    export interface AssetBreakdownBarProps {
      segments: Array<AssetSegment>
      total: number
      currency: string
      showLegend?: boolean
      className?: string
    }

    const colorClasses = {
      success: { bg: 'bg-success', dot: 'bg-success' },
      info: { bg: 'bg-info', dot: 'bg-info' },
      secondary: { bg: 'bg-secondary', dot: 'bg-secondary' },
      warning: { bg: 'bg-warning', dot: 'bg-warning' },
      primary: { bg: 'bg-primary', dot: 'bg-primary' },
    } as const

    export function AssetBreakdownBar({
      segments,
      total,
      currency,
      showLegend = true,
      className,
    }: AssetBreakdownBarProps) {
      // Filter out zero-value segments for the bar
      const nonZeroSegments = segments.filter((s) => s.value > 0)

      return (
        <div className={cn('space-y-3', className)}>
          {/* Stacked Bar */}
          <div className="h-4 w-full rounded-full bg-base-200 overflow-hidden flex">
            {nonZeroSegments.map((segment, index) => {
              const percentage = total > 0 ? (segment.value / total) * 100 : 0
              return (
                <div
                  key={index}
                  className={cn(
                    colorClasses[segment.color].bg,
                    'h-full transition-all'
                  )}
                  style={{ width: `${percentage}%` }}
                  role="img"
                  aria-label={`${segment.label}: ${formatCurrency(segment.value, currency)} (${percentage.toFixed(0)}%)`}
                />
              )
            })}
          </div>

          {/* Legend */}
          {showLegend && (
            <div className="space-y-2">
              {segments.map((segment, index) => (
                <div
                  key={index}
                  className="flex items-center justify-between text-sm"
                >
                  <div className="flex items-center gap-2">
                    <span
                      className={cn(
                        'h-3 w-3 rounded-full shrink-0',
                        colorClasses[segment.color].dot
                      )}
                    />
                    <span className="text-base-content/80">{segment.label}</span>
                  </div>
                  <span className="font-medium tabular-nums">
                    {formatCurrency(segment.value, currency)}
                  </span>
                </div>
              ))}
            </div>
          )}
        </div>
      )
    }
    ```

  **Portal-web Business Health Page:**
  - [ ] Create `portal-web/src/features/reports/components/BusinessHealthPage.tsx`:
    ```typescript
    /**
     * Business Health Page
     *
     * Shows financial position (balance sheet) in plain language:
     * - Hero: What Your Business Is Worth (Total Equity)
     * - Section: What You Own (Assets breakdown)
     * - Section: What You Owe (Liabilities - currently $0)
     * - Section: Owner's Stake (Equity breakdown)
     * - Insight Card with contextual advice
     * - Quick Links to related pages
     */
    import { Link, useNavigate, useParams, useSearch, useRouteContext } from '@tanstack/react-router'
    import { useTranslation } from 'react-i18next'
    import {
      ArrowLeft,
      ArrowRight,
      Calendar,
      ChevronLeft,
      ChevronRight,
      HelpCircle,
    } from 'lucide-react'

    import { useFinancialPositionQuery } from '@/api/accounting'
    import { Button, InsightCard, Skeleton } from '@/components'
    import { formatCurrency } from '@/lib/formatCurrency'
    import { formatDateShort } from '@/lib/formatDate'
    import { useLanguage } from '@/hooks/useLanguage'
    import { cn } from '@/lib/utils'
    import { AssetBreakdownBar } from './AssetBreakdownBar'

    export function BusinessHealthPage() {
      const { t } = useTranslation('reports')
      const { t: tCommon } = useTranslation('common')
      const { isRTL } = useLanguage()
      const navigate = useNavigate()
      const { businessDescriptor } = useParams({
        from: '/business/$businessDescriptor/reports/health',
      })
      const { asOf } = useSearch({
        from: '/business/$businessDescriptor/reports/health',
      })
      const { business } = useRouteContext({
        from: '/business/$businessDescriptor',
      })

      const BackIcon = isRTL ? ArrowRight : ArrowLeft
      const ChevronIcon = isRTL ? ChevronLeft : ChevronRight

      // Data query
      const { data, isLoading, error } = useFinancialPositionQuery(
        businessDescriptor,
        asOf
      )

      const currency = business.currency
      const displayDate = asOf ? new Date(asOf) : new Date()

      // Parse amounts
      const totalEquity = parseFloat(data?.totalEquity ?? '0')
      const totalAssets = parseFloat(data?.totalAssets ?? '0')
      const totalLiabilities = parseFloat(data?.totalLiabilities ?? '0')
      const cashOnHand = parseFloat(data?.cashOnHand ?? '0')
      const inventoryValue = parseFloat(data?.totalInventoryValue ?? '0')
      const fixedAssets = parseFloat(data?.fixedAssets ?? '0')
      const ownerInvestment = parseFloat(data?.ownerInvestment ?? '0')
      const ownerDraws = parseFloat(data?.ownerDraws ?? '0')
      const retainedEarnings = parseFloat(data?.retainedEarnings ?? '0')

      // Determine insight variant
      const getInsightConfig = () => {
        if (cashOnHand < 0) {
          return { variant: 'warning' as const, message: t('insights.negative_cash') }
        }
        if (totalEquity < 0) {
          return { variant: 'error' as const, message: t('insights.negative_equity') }
        }
        return { variant: 'success' as const, message: t('insights.healthy') }
      }

      const insight = getInsightConfig()

      // Error state
      if (error) {
        return (
          <div className="flex flex-col items-center justify-center py-16 px-4 text-center">
            <h2 className="text-lg font-semibold mb-2">{t('error.title')}</h2>
            <p className="text-base-content/60 mb-4">{t('error.body')}</p>
            <Button variant="ghost" onClick={() => window.location.reload()}>
              {tCommon('retry')}
            </Button>
          </div>
        )
      }

      return (
        <div className="space-y-6 pb-20 md:pb-6">
          {/* Header */}
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <Button
                variant="ghost"
                size="sm"
                onClick={() =>
                  navigate({
                    to: '/business/$businessDescriptor/reports',
                    params: { businessDescriptor },
                    search: asOf ? { asOf } : undefined,
                  })
                }
                aria-label={tCommon('back')}
              >
                <BackIcon className="h-5 w-5" />
              </Button>
              <div>
                <h1 className="text-xl font-bold">{t('health.title')}</h1>
                <p className="text-sm text-base-content/60">{t('health.subtitle')}</p>
              </div>
            </div>
            <button className="flex items-center gap-2 text-sm text-base-content/70">
              <Calendar className="h-4 w-4" />
              <span>{t('hub.as_of', { date: formatDateShort(displayDate) })}</span>
            </button>
          </div>

          {/* Hero: Business Worth */}
          <section className="bg-base-200/50 rounded-box p-6 text-center">
            {isLoading ? (
              <>
                <Skeleton className="h-4 w-32 mx-auto mb-2" />
                <Skeleton className="h-10 w-48 mx-auto" />
              </>
            ) : (
              <>
                <p className="text-sm text-base-content/60 mb-1">
                  {t('health.business_worth')}
                </p>
                <p
                  className={cn(
                    'text-4xl font-bold',
                    totalEquity < 0 ? 'text-error' : 'text-success'
                  )}
                >
                  {formatCurrency(totalEquity, currency)}
                </p>
              </>
            )}
          </section>

          {/* Section: What You Own */}
          <section>
            <h2 className="text-lg font-semibold mb-4">{t('health.what_you_own')}</h2>
            {isLoading ? (
              <div className="space-y-3">
                <Skeleton className="h-4 w-full rounded-full" />
                <Skeleton className="h-4 w-3/4" />
                <Skeleton className="h-4 w-2/3" />
              </div>
            ) : (
              <div className="rounded-box border border-base-300 p-4">
                <AssetBreakdownBar
                  segments={[
                    { label: t('health.cash_on_hand'), value: cashOnHand, color: 'success' },
                    { label: t('health.inventory_value'), value: inventoryValue, color: 'info' },
                    { label: t('health.fixed_assets'), value: fixedAssets, color: 'secondary' },
                  ]}
                  total={totalAssets}
                  currency={currency}
                />
                <div className="border-t border-base-300 mt-4 pt-3 flex justify-between font-semibold">
                  <span>{t('health.total_assets')}</span>
                  <span className="tabular-nums">{formatCurrency(totalAssets, currency)}</span>
                </div>
              </div>
            )}
          </section>

          {/* Section: What You Owe */}
          <section>
            <h2 className="text-lg font-semibold mb-4">{t('health.what_you_owe')}</h2>
            <div className="rounded-box border border-base-300 p-4">
              <div className="flex justify-between items-center mb-2">
                <span>{t('health.total_liabilities')}</span>
                <span className="text-xl font-bold tabular-nums">
                  {formatCurrency(totalLiabilities, currency)}
                </span>
              </div>
              <p className="text-sm text-base-content/60">{t('health.liabilities_note')}</p>
            </div>
          </section>

          {/* Section: Owner's Stake */}
          <section>
            <h2 className="text-lg font-semibold mb-4">{t('health.owners_stake')}</h2>
            {isLoading ? (
              <div className="space-y-3">
                <Skeleton className="h-6 w-full" />
                <Skeleton className="h-6 w-full" />
                <Skeleton className="h-6 w-full" />
              </div>
            ) : (
              <div className="rounded-box border border-base-300 p-4 space-y-3">
                <div className="flex justify-between text-sm">
                  <span>{t('health.money_put_in')}</span>
                  <span className="font-medium tabular-nums text-success">
                    {formatCurrency(ownerInvestment, currency)}
                  </span>
                </div>
                <div className="flex justify-between text-sm">
                  <span>{t('health.money_took_out')}</span>
                  <span className="font-medium tabular-nums text-error">
                    -{formatCurrency(ownerDraws, currency)}
                  </span>
                </div>
                <div className="flex justify-between text-sm">
                  <span>{t('health.profit_kept')}</span>
                  <span
                    className={cn(
                      'font-medium tabular-nums',
                      retainedEarnings >= 0 ? 'text-success' : 'text-error'
                    )}
                  >
                    {formatCurrency(retainedEarnings, currency)}
                  </span>
                </div>
                <div className="border-t border-base-300 pt-3 flex justify-between font-semibold">
                  <span>{t('health.business_value')}</span>
                  <span
                    className={cn(
                      'tabular-nums',
                      totalEquity >= 0 ? 'text-success' : 'text-error'
                    )}
                  >
                    {formatCurrency(totalEquity, currency)}
                  </span>
                </div>
              </div>
            )}
          </section>

          {/* Insight Card */}
          {!isLoading && (
            <InsightCard variant={insight.variant} message={insight.message} />
          )}

          {/* Quick Links */}
          <section>
            <div className="flex flex-wrap gap-4">
              <Link
                to="/business/$businessDescriptor/accounting/assets"
                params={{ businessDescriptor }}
                className="flex items-center gap-1 text-sm font-medium text-primary hover:underline"
              >
                {t('quick_links.view_assets')}
                <ChevronIcon className="h-4 w-4" />
              </Link>
              <Link
                to="/business/$businessDescriptor/accounting/capital"
                params={{ businessDescriptor }}
                className="flex items-center gap-1 text-sm font-medium text-primary hover:underline"
              >
                {t('quick_links.view_capital')}
                <ChevronIcon className="h-4 w-4" />
              </Link>
            </div>
          </section>
        </div>
      )
    }
    ```

  **Portal-web route:**
  - [ ] Create `portal-web/src/routes/business/$businessDescriptor/reports/health.tsx`:
    ```typescript
    import { createFileRoute } from '@tanstack/react-router'
    import { z } from 'zod'
    import { financialPositionQueryOptions } from '@/api/accounting'
    import { BusinessHealthPage } from '@/features/reports/components/BusinessHealthPage'

    const searchSchema = z.object({
      asOf: z.string().optional(),
    })

    export const Route = createFileRoute(
      '/business/$businessDescriptor/reports/health'
    )({
      validateSearch: searchSchema,
      staticData: {
        titleKey: 'common.pages.reports_health',
      },
      loader: async ({ context, params, search }) => {
        await context.queryClient.ensureQueryData(
          financialPositionQueryOptions(params.businessDescriptor, search.asOf)
        )
      },
      component: BusinessHealthPage,
    })
    ```

  - [ ] Update `portal-web/src/features/reports/components/index.ts` — add exports:
    ```typescript
    export { AssetBreakdownBar } from './AssetBreakdownBar'
    export type { AssetBreakdownBarProps, AssetSegment } from './AssetBreakdownBar'
    export { BusinessHealthPage } from './BusinessHealthPage'
    ```

  - [ ] Add page title to i18n common.json (en):
    ```json
    "reports_health": "Business Health"
    ```

  - [ ] Add page title to i18n common.json (ar):
    ```json
    "reports_health": "صحة عملك"
    ```

- **Edge cases + error handling**:
  - Negative equity: show error variant insight
  - Negative cash: show warning variant insight
  - Zero total assets: bar is empty, legend still shows all items
  - Loading: skeleton for all sections

- **Verification checklist**:
  - [ ] Run `pnpm tsc --noEmit` — no TypeScript errors
  - [ ] Run `pnpm lint` — no ESLint errors
  - [ ] Route accessible at `/business/:descriptor/reports/health`
  - [ ] Back button navigates to hub with preserved asOf param
  - [ ] AssetBreakdownBar renders proportional segments
  - [ ] Insight card shows appropriate variant based on data

- **Definition of done**:
  - [ ] TypeScript compiles without errors
  - [ ] ESLint passes without errors
  - [ ] Business Health page renders all sections correctly
  - [ ] RTL layout works (back arrow direction, chevrons)

- **Relevant instruction files**:
  - `.github/instructions/ui-implementation.instructions.md` — RTL, daisyUI
  - `.github/instructions/design-tokens.instructions.md` — color tokens
  - `.github/instructions/charts.instructions.md` — visualization patterns

---

### Step 6 — ProfitWaterfall + Profit & Earnings Page

- **Goal (user-visible outcome)**: Users can view the Profit & Earnings page showing revenue → costs → profit breakdown.

- **Scope**:
  - **In**: ProfitWaterfall component, Profit page, expenses doughnut chart, route
  - **Out**: Cash Movement page

- **Targets (files/symbols/endpoints)**:
  - Create: `portal-web/src/features/reports/components/ProfitWaterfall.tsx`
  - Create: `portal-web/src/features/reports/components/ProfitEarningsPage.tsx`
  - Create: `portal-web/src/routes/business/$businessDescriptor/reports/profit.tsx`
  - Modify: `portal-web/src/features/reports/components/index.ts`

- **Tasks (detailed checklist)**:

  **Portal-web ProfitWaterfall:**
  - [ ] Create `portal-web/src/features/reports/components/ProfitWaterfall.tsx`:
    ```typescript
    /**
     * ProfitWaterfall Component
     *
     * Stepped visualization showing P&L flow:
     * Revenue → Product Costs → Gross Profit → Running Costs → Net Profit
     */
    import { ArrowDown } from 'lucide-react'
    import { cn } from '@/lib/utils'
    import { formatCurrency } from '@/lib/formatCurrency'

    export interface WaterfallStep {
      label: string
      value: number
      type: 'total' | 'subtract' | 'result'
    }

    export interface ProfitWaterfallProps {
      steps: Array<WaterfallStep>
      currency: string
      className?: string
    }

    export function ProfitWaterfall({
      steps,
      currency,
      className,
    }: ProfitWaterfallProps) {
      // Find max value for bar scaling
      const maxValue = Math.max(...steps.map((s) => Math.abs(s.value)))

      return (
        <div className={cn('space-y-2', className)}>
          {steps.map((step, index) => {
            const percentage = maxValue > 0 ? (Math.abs(step.value) / maxValue) * 100 : 0
            const isSubtract = step.type === 'subtract'
            const isResult = step.type === 'result'
            const isLast = index === steps.length - 1

            return (
              <div key={index}>
                {/* Step Row */}
                <div className="flex items-center gap-3">
                  {/* Prefix */}
                  <span
                    className={cn(
                      'w-4 text-center text-sm font-medium',
                      isSubtract && 'text-error',
                      isResult && 'text-base-content'
                    )}
                  >
                    {isSubtract ? '-' : isResult ? '=' : ''}
                  </span>

                  {/* Label and Bar */}
                  <div className="flex-1">
                    <div className="flex justify-between items-center mb-1">
                      <span
                        className={cn(
                          'text-sm',
                          isResult && 'font-semibold',
                          isSubtract && 'text-base-content/70'
                        )}
                      >
                        {step.label}
                      </span>
                      <span
                        className={cn(
                          'font-medium tabular-nums',
                          isSubtract && 'text-error',
                          isResult && step.value >= 0 && 'text-success',
                          isResult && step.value < 0 && 'text-error'
                        )}
                      >
                        {isSubtract && '-'}
                        {formatCurrency(Math.abs(step.value), currency)}
                      </span>
                    </div>
                    <div className="h-3 w-full bg-base-200 rounded-full overflow-hidden">
                      <div
                        className={cn(
                          'h-full rounded-full transition-all',
                          isSubtract && 'bg-base-300',
                          !isSubtract && !isResult && 'bg-primary',
                          isResult && step.value >= 0 && 'bg-success',
                          isResult && step.value < 0 && 'bg-error'
                        )}
                        style={{ width: `${percentage}%` }}
                      />
                    </div>
                  </div>
                </div>

                {/* Arrow between steps */}
                {!isLast && (
                  <div className="flex justify-center py-1">
                    <ArrowDown className="h-4 w-4 text-base-content/30" />
                  </div>
                )}
              </div>
            )
          })}
        </div>
      )
    }
    ```

  **Portal-web Profit & Earnings Page:**
  - [ ] Create `portal-web/src/features/reports/components/ProfitEarningsPage.tsx`:
    ```typescript
    /**
     * Profit & Earnings Page
     *
     * Shows P&L in plain language:
     * - Hero: What You Keep (Net Profit)
     * - Waterfall: Revenue → Costs → Profit
     * - Margin stats
     * - Running costs breakdown (doughnut chart)
     * - Insights
     * - Quick links
     */
    import { Link, useNavigate, useParams, useSearch, useRouteContext } from '@tanstack/react-router'
    import { useTranslation } from 'react-i18next'
    import {
      ArrowLeft,
      ArrowRight,
      Calendar,
      ChevronLeft,
      ChevronRight,
    } from 'lucide-react'

    import { useProfitAndLossQuery } from '@/api/accounting'
    import { Button, InsightCard, Skeleton, StatCard, StatCardSkeleton } from '@/components'
    import { StatCardGroup } from '@/components/molecules/StatCardGroup'
    import { ChartCard } from '@/components/charts/ChartCard'
    import { DoughnutChart } from '@/components/charts/DoughnutChart'
    import { formatCurrency } from '@/lib/formatCurrency'
    import { formatDateShort } from '@/lib/formatDate'
    import { useLanguage } from '@/hooks/useLanguage'
    import { cn } from '@/lib/utils'
    import { ProfitWaterfall } from './ProfitWaterfall'

    export function ProfitEarningsPage() {
      const { t } = useTranslation('reports')
      const { t: tCommon } = useTranslation('common')
      const { isRTL } = useLanguage()
      const navigate = useNavigate()
      const { businessDescriptor } = useParams({
        from: '/business/$businessDescriptor/reports/profit',
      })
      const { asOf } = useSearch({
        from: '/business/$businessDescriptor/reports/profit',
      })
      const { business } = useRouteContext({
        from: '/business/$businessDescriptor',
      })

      const BackIcon = isRTL ? ArrowRight : ArrowLeft
      const ChevronIcon = isRTL ? ChevronLeft : ChevronRight

      // Data query
      const { data, isLoading, error } = useProfitAndLossQuery(
        businessDescriptor,
        asOf
      )

      const currency = business.currency
      const displayDate = asOf ? new Date(asOf) : new Date()

      // Parse amounts
      const revenue = parseFloat(data?.revenue ?? '0')
      const cogs = parseFloat(data?.cogs ?? '0')
      const grossProfit = parseFloat(data?.grossProfit ?? '0')
      const totalExpenses = parseFloat(data?.totalExpenses ?? '0')
      const netProfit = parseFloat(data?.netProfit ?? '0')

      // Calculate margins
      const grossMargin = revenue > 0 ? ((grossProfit / revenue) * 100).toFixed(1) : '0'
      const profitMargin = revenue > 0 ? ((netProfit / revenue) * 100).toFixed(1) : '0'

      // Expenses by category for chart
      const expensesByCategory = data?.expensesByCategory ?? []
      const chartLabels = expensesByCategory.map((e) => e.key)
      const chartValues = expensesByCategory.map((e) => parseFloat(e.value))

      // Find biggest expense
      const biggestExpense = expensesByCategory.reduce(
        (max, curr) => {
          const val = parseFloat(curr.value)
          return val > max.value ? { category: curr.key, value: val } : max
        },
        { category: '', value: 0 }
      )

      const biggestExpensePercent =
        totalExpenses > 0
          ? ((biggestExpense.value / totalExpenses) * 100).toFixed(0)
          : '0'

      // Error state
      if (error) {
        return (
          <div className="flex flex-col items-center justify-center py-16 px-4 text-center">
            <h2 className="text-lg font-semibold mb-2">{t('error.title')}</h2>
            <p className="text-base-content/60 mb-4">{t('error.body')}</p>
            <Button variant="ghost" onClick={() => window.location.reload()}>
              {tCommon('retry')}
            </Button>
          </div>
        )
      }

      return (
        <div className="space-y-6 pb-20 md:pb-6">
          {/* Header */}
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <Button
                variant="ghost"
                size="sm"
                onClick={() =>
                  navigate({
                    to: '/business/$businessDescriptor/reports',
                    params: { businessDescriptor },
                    search: asOf ? { asOf } : undefined,
                  })
                }
                aria-label={tCommon('back')}
              >
                <BackIcon className="h-5 w-5" />
              </Button>
              <div>
                <h1 className="text-xl font-bold">{t('profit.title')}</h1>
                <p className="text-sm text-base-content/60">{t('profit.subtitle')}</p>
              </div>
            </div>
            <button className="flex items-center gap-2 text-sm text-base-content/70">
              <Calendar className="h-4 w-4" />
              <span>{t('hub.as_of', { date: formatDateShort(displayDate) })}</span>
            </button>
          </div>

          {/* Hero: What You Keep */}
          <section className="bg-base-200/50 rounded-box p-6 text-center">
            {isLoading ? (
              <>
                <Skeleton className="h-4 w-32 mx-auto mb-2" />
                <Skeleton className="h-10 w-48 mx-auto mb-1" />
                <Skeleton className="h-3 w-48 mx-auto" />
              </>
            ) : (
              <>
                <p className="text-sm text-base-content/60 mb-1">
                  {t('profit.what_you_keep')}
                </p>
                <p
                  className={cn(
                    'text-4xl font-bold',
                    netProfit >= 0 ? 'text-success' : 'text-error'
                  )}
                >
                  {formatCurrency(netProfit, currency)}
                </p>
                <p className="text-sm text-base-content/60 mt-1">
                  {t('profit.what_you_keep_subtitle')}
                </p>
              </>
            )}
          </section>

          {/* Profit Flow Waterfall */}
          <section>
            <h2 className="text-lg font-semibold mb-4">{t('profit.money_in')}</h2>
            {isLoading ? (
              <div className="space-y-4">
                {[1, 2, 3, 4, 5].map((i) => (
                  <Skeleton key={i} className="h-8 w-full" />
                ))}
              </div>
            ) : (
              <div className="rounded-box border border-base-300 p-4">
                <ProfitWaterfall
                  steps={[
                    { label: t('profit.money_in'), value: revenue, type: 'total' },
                    { label: t('profit.product_costs'), value: cogs, type: 'subtract' },
                    { label: t('profit.gross_profit'), value: grossProfit, type: 'result' },
                    { label: t('profit.running_costs'), value: totalExpenses, type: 'subtract' },
                    { label: t('profit.what_you_keep'), value: netProfit, type: 'result' },
                  ]}
                  currency={currency}
                />
              </div>
            )}
          </section>

          {/* Margin Stats */}
          <section>
            <StatCardGroup cols={2}>
              {isLoading ? (
                <>
                  <StatCardSkeleton />
                  <StatCardSkeleton />
                </>
              ) : (
                <>
                  <StatCard
                    label={t('profit.gross_margin')}
                    value={`${grossMargin}%`}
                    variant="info"
                  />
                  <StatCard
                    label={t('profit.profit_margin')}
                    value={`${profitMargin}%`}
                    variant={parseFloat(profitMargin) >= 0 ? 'success' : 'error'}
                  />
                </>
              )}
            </StatCardGroup>
          </section>

          {/* Running Costs Breakdown */}
          {expensesByCategory.length > 0 && (
            <section>
              <h2 className="text-lg font-semibold mb-4">{t('profit.running_costs')}</h2>
              <ChartCard
                title={t('profit.total_running_costs')}
                subtitle={formatCurrency(totalExpenses, currency)}
                isLoading={isLoading}
                isEmpty={expensesByCategory.length === 0}
              >
                <div className="h-64">
                  <DoughnutChart
                    data={{
                      labels: chartLabels,
                      datasets: [
                        {
                          data: chartValues,
                          backgroundColor: [
                            'hsl(var(--p))',
                            'hsl(var(--s))',
                            'hsl(var(--a))',
                            'hsl(var(--in))',
                            'hsl(var(--wa))',
                            'hsl(var(--er))',
                          ],
                        },
                      ],
                    }}
                    centerLabel={{
                      value: formatCurrency(totalExpenses, currency),
                      label: t('profit.total_running_costs'),
                    }}
                  />
                </div>
              </ChartCard>
            </section>
          )}

          {/* Insights */}
          {!isLoading && (
            <div className="space-y-3">
              {netProfit >= 0 ? (
                <InsightCard
                  variant="success"
                  message={t('insights.profitable', {
                    currency,
                    margin: profitMargin,
                  })}
                />
              ) : (
                <InsightCard
                  variant="error"
                  message={t('insights.losing', {
                    amount: formatCurrency(Math.abs(netProfit), currency),
                  })}
                />
              )}
              {biggestExpense.category && (
                <InsightCard
                  variant="info"
                  message={t('insights.biggest_expense', {
                    category: biggestExpense.category,
                    amount: formatCurrency(biggestExpense.value, currency),
                    percent: biggestExpensePercent,
                  })}
                />
              )}
            </div>
          )}

          {/* Quick Links */}
          <section>
            <div className="flex flex-wrap gap-4">
              <Link
                to="/business/$businessDescriptor/accounting/expenses"
                params={{ businessDescriptor }}
                className="flex items-center gap-1 text-sm font-medium text-primary hover:underline"
              >
                {t('quick_links.view_expenses')}
                <ChevronIcon className="h-4 w-4" />
              </Link>
              <Link
                to="/business/$businessDescriptor/orders"
                params={{ businessDescriptor }}
                className="flex items-center gap-1 text-sm font-medium text-primary hover:underline"
              >
                {t('quick_links.view_orders')}
                <ChevronIcon className="h-4 w-4" />
              </Link>
            </div>
          </section>
        </div>
      )
    }
    ```

  **Portal-web route:**
  - [ ] Create `portal-web/src/routes/business/$businessDescriptor/reports/profit.tsx`:
    ```typescript
    import { createFileRoute } from '@tanstack/react-router'
    import { z } from 'zod'
    import { profitAndLossQueryOptions } from '@/api/accounting'
    import { ProfitEarningsPage } from '@/features/reports/components/ProfitEarningsPage'

    const searchSchema = z.object({
      asOf: z.string().optional(),
    })

    export const Route = createFileRoute(
      '/business/$businessDescriptor/reports/profit'
    )({
      validateSearch: searchSchema,
      staticData: {
        titleKey: 'common.pages.reports_profit',
      },
      loader: async ({ context, params, search }) => {
        await context.queryClient.ensureQueryData(
          profitAndLossQueryOptions(params.businessDescriptor, search.asOf)
        )
      },
      component: ProfitEarningsPage,
    })
    ```

  - [ ] Update `portal-web/src/features/reports/components/index.ts`:
    ```typescript
    export { ProfitWaterfall } from './ProfitWaterfall'
    export type { ProfitWaterfallProps, WaterfallStep } from './ProfitWaterfall'
    export { ProfitEarningsPage } from './ProfitEarningsPage'
    ```

  - [ ] Add page title to i18n common.json (en): `"reports_profit": "Profit & Earnings"`
  - [ ] Add page title to i18n common.json (ar): `"reports_profit": "أرباحك"`

- **Edge cases + error handling**:
  - Negative net profit: error variant insight
  - Zero revenue: margins show 0%
  - No expenses: doughnut chart section hidden
  - Loading: skeleton for all sections

- **Verification checklist**:
  - [ ] Run `pnpm tsc --noEmit` — no TypeScript errors
  - [ ] Run `pnpm lint` — no ESLint errors
  - [ ] Route accessible at `/business/:descriptor/reports/profit`
  - [ ] Waterfall shows correct flow with arrows
  - [ ] Doughnut chart renders expense breakdown
  - [ ] Insights show appropriate messages

- **Definition of done**:
  - [ ] TypeScript compiles without errors
  - [ ] ESLint passes without errors
  - [ ] Profit page renders all sections correctly
  - [ ] Charts use existing ChartCard/DoughnutChart components

- **Relevant instruction files**:
  - `.github/instructions/charts.instructions.md` — DoughnutChart usage, chart patterns
  - `.github/instructions/ui-implementation.instructions.md` — RTL, daisyUI

---

### Step 7 — CashFlowDiagram + Cash Movement Page

- **Goal (user-visible outcome)**: Users can view the Cash Movement page showing cash in/out flow.

- **Scope**:
  - **In**: CashFlowDiagram component, Cash Movement page, route
  - **Out**: Navigation updates

- **Targets (files/symbols/endpoints)**:
  - Create: `portal-web/src/features/reports/components/CashFlowDiagram.tsx`
  - Create: `portal-web/src/features/reports/components/CashMovementPage.tsx`
  - Create: `portal-web/src/routes/business/$businessDescriptor/reports/cashflow.tsx`
  - Modify: `portal-web/src/features/reports/components/index.ts`

- **Tasks (detailed checklist)**:

  **Portal-web CashFlowDiagram:**
  - [ ] Create `portal-web/src/features/reports/components/CashFlowDiagram.tsx`:
    ```typescript
    /**
     * CashFlowDiagram Component
     *
     * Visual flow representation for cash movement:
     * Start → Money In box → Money Out box → Net Change → End
     */
    import { ArrowDown, TrendingDown, TrendingUp } from 'lucide-react'
    import { useTranslation } from 'react-i18next'
    import { cn } from '@/lib/utils'
    import { formatCurrency } from '@/lib/formatCurrency'

    export interface CashFlowItem {
      label: string
      value: number
    }

    export interface CashFlowDiagramProps {
      startAmount: number
      endAmount: number
      cashIn: Array<CashFlowItem>
      cashOut: Array<CashFlowItem>
      currency: string
      className?: string
    }

    export function CashFlowDiagram({
      startAmount,
      endAmount,
      cashIn,
      cashOut,
      currency,
      className,
    }: CashFlowDiagramProps) {
      const { t } = useTranslation('reports')

      const totalIn = cashIn.reduce((sum, item) => sum + item.value, 0)
      const totalOut = cashOut.reduce((sum, item) => sum + item.value, 0)
      const netChange = totalIn - totalOut
      const isPositiveChange = netChange >= 0

      return (
        <div className={cn('space-y-4', className)}>
          {/* Start Node */}
          <div className="text-center">
            <span className="text-sm text-base-content/60">{t('cashflow.cash_start')}</span>
            <p className="text-lg font-semibold tabular-nums">
              {formatCurrency(startAmount, currency)}
            </p>
          </div>

          <div className="flex justify-center">
            <ArrowDown className="h-5 w-5 text-base-content/30" />
          </div>

          {/* Money Coming In Box */}
          <div className="rounded-box border border-success/20 bg-success/5 p-4">
            <h3 className="font-semibold text-success mb-3">{t('cashflow.money_in')}</h3>
            <div className="space-y-2">
              {cashIn.map((item, index) => (
                <div key={index} className="flex justify-between text-sm">
                  <span className="text-base-content/80">{item.label}</span>
                  <span className="font-medium tabular-nums text-success">
                    +{formatCurrency(item.value, currency)}
                  </span>
                </div>
              ))}
              <div className="border-t border-success/20 pt-2 mt-2 flex justify-between font-semibold">
                <span>{t('cashflow.total_in')}</span>
                <span className="tabular-nums text-success">
                  {formatCurrency(totalIn, currency)}
                </span>
              </div>
            </div>
          </div>

          <div className="flex justify-center">
            <ArrowDown className="h-5 w-5 text-base-content/30" />
          </div>

          {/* Money Going Out Box */}
          <div className="rounded-box border border-warning/20 bg-warning/5 p-4">
            <h3 className="font-semibold text-warning mb-3">{t('cashflow.money_out')}</h3>
            <div className="space-y-2">
              {cashOut.map((item, index) => (
                <div key={index} className="flex justify-between text-sm">
                  <span className="text-base-content/80">{item.label}</span>
                  <span className="font-medium tabular-nums text-error">
                    -{formatCurrency(item.value, currency)}
                  </span>
                </div>
              ))}
              <div className="border-t border-warning/20 pt-2 mt-2 flex justify-between font-semibold">
                <span>{t('cashflow.total_out')}</span>
                <span className="tabular-nums text-error">
                  {formatCurrency(totalOut, currency)}
                </span>
              </div>
            </div>
          </div>

          <div className="flex justify-center">
            <ArrowDown className="h-5 w-5 text-base-content/30" />
          </div>

          {/* Net Change Indicator */}
          <div className="text-center py-2">
            <div className="flex items-center justify-center gap-2">
              {isPositiveChange ? (
                <TrendingUp className="h-5 w-5 text-success" />
              ) : (
                <TrendingDown className="h-5 w-5 text-error" />
              )}
              <span className="text-sm text-base-content/60">{t('cashflow.net_change')}</span>
            </div>
            <p
              className={cn(
                'text-xl font-bold tabular-nums',
                isPositiveChange ? 'text-success' : 'text-error'
              )}
            >
              {isPositiveChange ? '+' : ''}
              {formatCurrency(netChange, currency)}
            </p>
          </div>

          <div className="flex justify-center">
            <ArrowDown className="h-5 w-5 text-base-content/30" />
          </div>

          {/* End Node */}
          <div className="text-center bg-base-200/50 rounded-box p-4">
            <span className="text-sm text-base-content/60">{t('cashflow.cash_end')}</span>
            <p
              className={cn(
                'text-2xl font-bold tabular-nums',
                endAmount >= 0 ? 'text-success' : 'text-error'
              )}
            >
              {formatCurrency(endAmount, currency)}
            </p>
          </div>
        </div>
      )
    }
    ```

  **Portal-web Cash Movement Page:**
  - [ ] Create `portal-web/src/features/reports/components/CashMovementPage.tsx`:
    ```typescript
    /**
     * Cash Movement Page
     *
     * Shows cash flow in plain language:
     * - Hero: Cash Now
     * - Flow Diagram: Start → In → Out → End
     * - Insights
     * - Quick links
     */
    import { Link, useNavigate, useParams, useSearch, useRouteContext } from '@tanstack/react-router'
    import { useTranslation } from 'react-i18next'
    import {
      ArrowLeft,
      ArrowRight,
      Calendar,
      ChevronLeft,
      ChevronRight,
    } from 'lucide-react'

    import { useCashFlowQuery } from '@/api/accounting'
    import { Button, InsightCard, Skeleton } from '@/components'
    import { formatCurrency } from '@/lib/formatCurrency'
    import { formatDateShort } from '@/lib/formatDate'
    import { useLanguage } from '@/hooks/useLanguage'
    import { cn } from '@/lib/utils'
    import { CashFlowDiagram } from './CashFlowDiagram'

    export function CashMovementPage() {
      const { t } = useTranslation('reports')
      const { t: tCommon } = useTranslation('common')
      const { isRTL } = useLanguage()
      const navigate = useNavigate()
      const { businessDescriptor } = useParams({
        from: '/business/$businessDescriptor/reports/cashflow',
      })
      const { asOf } = useSearch({
        from: '/business/$businessDescriptor/reports/cashflow',
      })
      const { business } = useRouteContext({
        from: '/business/$businessDescriptor',
      })

      const BackIcon = isRTL ? ArrowRight : ArrowLeft
      const ChevronIcon = isRTL ? ChevronLeft : ChevronRight

      // Data query
      const { data, isLoading, error } = useCashFlowQuery(businessDescriptor, asOf)

      const currency = business.currency
      const displayDate = asOf ? new Date(asOf) : new Date()

      // Parse amounts
      const cashAtStart = parseFloat(data?.cashAtStart ?? '0')
      const cashAtEnd = parseFloat(data?.cashAtEnd ?? '0')
      const cashFromCustomers = parseFloat(data?.cashFromCustomers ?? '0')
      const cashFromOwner = parseFloat(data?.cashFromOwner ?? '0')
      const inventoryPurchases = parseFloat(data?.inventoryPurchases ?? '0')
      const operatingExpenses = parseFloat(data?.operatingExpenses ?? '0')
      const businessInvestments = parseFloat(data?.businessInvestments ?? '0')
      const ownerDraws = parseFloat(data?.ownerDraws ?? '0')
      const netCashFlow = parseFloat(data?.netCashFlow ?? '0')
      const totalCashIn = parseFloat(data?.totalCashIn ?? '0')

      // Determine insights
      const getInsights = () => {
        const insights: Array<{ variant: 'success' | 'warning' | 'info'; message: string }> = []

        if (netCashFlow >= 0) {
          insights.push({ variant: 'success', message: t('insights.cash_healthy') })
        } else {
          insights.push({ variant: 'warning', message: t('insights.cash_alert') })
        }

        if (ownerDraws > cashFromCustomers && cashFromCustomers > 0) {
          insights.push({ variant: 'info', message: t('insights.withdrawal_tip') })
        }

        return insights
      }

      const insights = getInsights()

      // Error state
      if (error) {
        return (
          <div className="flex flex-col items-center justify-center py-16 px-4 text-center">
            <h2 className="text-lg font-semibold mb-2">{t('error.title')}</h2>
            <p className="text-base-content/60 mb-4">{t('error.body')}</p>
            <Button variant="ghost" onClick={() => window.location.reload()}>
              {tCommon('retry')}
            </Button>
          </div>
        )
      }

      return (
        <div className="space-y-6 pb-20 md:pb-6">
          {/* Header */}
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <Button
                variant="ghost"
                size="sm"
                onClick={() =>
                  navigate({
                    to: '/business/$businessDescriptor/reports',
                    params: { businessDescriptor },
                    search: asOf ? { asOf } : undefined,
                  })
                }
                aria-label={tCommon('back')}
              >
                <BackIcon className="h-5 w-5" />
              </Button>
              <div>
                <h1 className="text-xl font-bold">{t('cashflow.title')}</h1>
                <p className="text-sm text-base-content/60">{t('cashflow.subtitle')}</p>
              </div>
            </div>
            <button className="flex items-center gap-2 text-sm text-base-content/70">
              <Calendar className="h-4 w-4" />
              <span>{t('hub.as_of', { date: formatDateShort(displayDate) })}</span>
            </button>
          </div>

          {/* Hero: Cash Now */}
          <section className="bg-base-200/50 rounded-box p-6 text-center">
            {isLoading ? (
              <>
                <Skeleton className="h-4 w-32 mx-auto mb-2" />
                <Skeleton className="h-10 w-48 mx-auto mb-1" />
                <Skeleton className="h-3 w-40 mx-auto" />
              </>
            ) : (
              <>
                <p className="text-sm text-base-content/60 mb-1">
                  {t('cashflow.cash_now')}
                </p>
                <p
                  className={cn(
                    'text-4xl font-bold',
                    cashAtEnd >= 0 ? 'text-success' : 'text-error'
                  )}
                >
                  {formatCurrency(cashAtEnd, currency)}
                </p>
                <p className="text-sm text-base-content/60 mt-1">
                  {t('cashflow.cash_runway')}
                </p>
              </>
            )}
          </section>

          {/* Cash Flow Diagram */}
          <section>
            {isLoading ? (
              <div className="space-y-4">
                <Skeleton className="h-16 w-full" />
                <Skeleton className="h-32 w-full" />
                <Skeleton className="h-32 w-full" />
                <Skeleton className="h-16 w-full" />
              </div>
            ) : (
              <CashFlowDiagram
                startAmount={cashAtStart}
                endAmount={cashAtEnd}
                cashIn={[
                  { label: t('cashflow.from_customers'), value: cashFromCustomers },
                  { label: t('cashflow.from_owner'), value: cashFromOwner },
                ]}
                cashOut={[
                  { label: t('cashflow.inventory_purchases'), value: inventoryPurchases },
                  { label: t('cashflow.running_costs'), value: operatingExpenses },
                  { label: t('cashflow.equipment_assets'), value: businessInvestments },
                  { label: t('cashflow.to_owner'), value: ownerDraws },
                ]}
                currency={currency}
              />
            )}
          </section>

          {/* Insights */}
          {!isLoading && (
            <div className="space-y-3">
              {insights.map((insight, index) => (
                <InsightCard
                  key={index}
                  variant={insight.variant}
                  message={insight.message}
                />
              ))}
            </div>
          )}

          {/* Quick Links */}
          <section>
            <div className="flex flex-wrap gap-4">
              <Link
                to="/business/$businessDescriptor/accounting/capital"
                params={{ businessDescriptor }}
                className="flex items-center gap-1 text-sm font-medium text-primary hover:underline"
              >
                {t('quick_links.view_capital')}
                <ChevronIcon className="h-4 w-4" />
              </Link>
              <Link
                to="/business/$businessDescriptor/accounting/expenses"
                params={{ businessDescriptor }}
                className="flex items-center gap-1 text-sm font-medium text-primary hover:underline"
              >
                {t('quick_links.view_expenses')}
                <ChevronIcon className="h-4 w-4" />
              </Link>
            </div>
          </section>
        </div>
      )
    }
    ```

  **Portal-web route:**
  - [ ] Create `portal-web/src/routes/business/$businessDescriptor/reports/cashflow.tsx`:
    ```typescript
    import { createFileRoute } from '@tanstack/react-router'
    import { z } from 'zod'
    import { cashFlowQueryOptions } from '@/api/accounting'
    import { CashMovementPage } from '@/features/reports/components/CashMovementPage'

    const searchSchema = z.object({
      asOf: z.string().optional(),
    })

    export const Route = createFileRoute(
      '/business/$businessDescriptor/reports/cashflow'
    )({
      validateSearch: searchSchema,
      staticData: {
        titleKey: 'common.pages.reports_cashflow',
      },
      loader: async ({ context, params, search }) => {
        await context.queryClient.ensureQueryData(
          cashFlowQueryOptions(params.businessDescriptor, search.asOf)
        )
      },
      component: CashMovementPage,
    })
    ```

  - [ ] Update `portal-web/src/features/reports/components/index.ts`:
    ```typescript
    export { CashFlowDiagram } from './CashFlowDiagram'
    export type { CashFlowDiagramProps, CashFlowItem } from './CashFlowDiagram'
    export { CashMovementPage } from './CashMovementPage'
    ```

  - [ ] Add page title to i18n common.json (en): `"reports_cashflow": "Cash Movement"`
  - [ ] Add page title to i18n common.json (ar): `"reports_cashflow": "حركة النقد"`

- **Edge cases + error handling**:
  - Negative net change: warning insight
  - Withdrawals > sales: tip insight
  - Zero values: items still shown with $0
  - Loading: skeleton for all sections

- **Verification checklist**:
  - [ ] Run `pnpm tsc --noEmit` — no TypeScript errors
  - [ ] Run `pnpm lint` — no ESLint errors
  - [ ] Route accessible at `/business/:descriptor/reports/cashflow`
  - [ ] Cash flow diagram shows flow with arrows
  - [ ] Insights show appropriate messages
  - [ ] RTL layout works correctly

- **Definition of done**:
  - [ ] TypeScript compiles without errors
  - [ ] ESLint passes without errors
  - [ ] Cash Movement page renders all sections
  - [ ] All 4 report pages are now accessible

- **Relevant instruction files**:
  - `.github/instructions/ui-implementation.instructions.md` — RTL, daisyUI
  - `.github/instructions/design-tokens.instructions.md` — color tokens

---

### Step 8 — Navigation Updates (Sidebar) + Final Polish

- **Goal (user-visible outcome)**: Users can access Reports from the sidebar navigation.

- **Scope**:
  - **In**: Add Reports to Sidebar, verify all routes work, final polish
  - **Out**: BottomNav update (documented as future consideration due to 5-item limit)

- **Targets (files/symbols/endpoints)**:
  - Modify: `portal-web/src/components/organisms/Sidebar.tsx`
  - Verify: All route imports generate correctly (`routeTree.gen.ts`)

- **Tasks (detailed checklist)**:

  **Portal-web Sidebar update:**
  - [ ] Modify `portal-web/src/components/organisms/Sidebar.tsx` — add Reports to `navItems` array:
    ```typescript
    import {
      Calculator,
      ChevronLeft,
      ChevronRight,
      FileBarChart, // Add this import
      LayoutDashboard,
      Package,
      ShoppingCart,
      Users,
      X,
    } from 'lucide-react'

    const navItems: Array<NavItem> = [
      { key: 'dashboard', icon: LayoutDashboard, path: '' },
      { key: 'inventory', icon: Package, path: '/inventory' },
      { key: 'orders', icon: ShoppingCart, path: '/orders' },
      { key: 'customers', icon: Users, path: '/customers' },
      { key: 'accounting', icon: Calculator, path: '/accounting' },
      { key: 'reports', icon: FileBarChart, path: '/reports' }, // Add this line
    ]
    ```

  **Final verification tasks:**
  - [ ] Run TanStack Router codegen to regenerate `routeTree.gen.ts`
  - [ ] Verify all 4 routes are registered:
    - `/business/$businessDescriptor/reports/`
    - `/business/$businessDescriptor/reports/health`
    - `/business/$businessDescriptor/reports/profit`
    - `/business/$businessDescriptor/reports/cashflow`
  - [ ] Test navigation flow:
    - Sidebar → Reports → Hub
    - Hub → Business Health → Back to Hub
    - Hub → Profit & Earnings → Back to Hub
    - Hub → Cash Movement → Back to Hub
  - [ ] Test `asOf` search param persistence across pages
  - [ ] Test empty state on new business
  - [ ] Test error state (can simulate by blocking network)
  - [ ] Test RTL layout (switch to Arabic)

  **Documentation note for BottomNav:**
  - BottomNav currently has 4 items (inventory, orders, customers, accounting)
  - Adding a 5th item requires UX decision:
    - Option A: Replace one item with Reports
    - Option B: Add "More" menu with Reports + Accounting
    - Option C: Keep Reports desktop-only (sidebar)
  - For now, Reports is accessible via Sidebar on all screen sizes
  - Document this as future enhancement in `DRIFT_TODO.md` if needed

- **Edge cases + error handling**:
  - Sidebar collapsed state: icon-only should show Reports icon
  - Mobile drawer: Reports should appear in nav list
  - Active state: should highlight when on any `/reports/*` path

- **Verification checklist**:
  - [ ] Run `pnpm tsc --noEmit` — no TypeScript errors
  - [ ] Run `pnpm lint` — no ESLint errors
  - [ ] Run `pnpm build` — build succeeds
  - [ ] Sidebar shows Reports icon in both expanded and collapsed states
  - [ ] Reports nav item is active when on any reports page
  - [ ] All pages render correctly in LTR and RTL
  - [ ] All pages handle loading/empty/error states

- **Definition of done**:
  - [ ] TypeScript compiles without errors
  - [ ] ESLint passes without errors
  - [ ] Build succeeds without errors
  - [ ] Reports accessible from Sidebar navigation
  - [ ] All 4 report pages function correctly
  - [ ] Feature is complete per BRD acceptance criteria

- **Relevant instruction files**:
  - `.github/instructions/portal-web-code-structure.instructions.md` — component placement
  - `.github/instructions/ui-implementation.instructions.md` — navigation patterns, RTL

---

## 4) API Contracts (high level)

All endpoints are **existing** backend APIs. No changes required.

### Endpoints Used:

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/v1/businesses/:businessDescriptor/accounting/summary` | GET | Safe to Draw + summary metrics |
| `/v1/businesses/:businessDescriptor/analytics/reports/financial-position` | GET | Balance sheet data |
| `/v1/businesses/:businessDescriptor/analytics/reports/profit-and-loss` | GET | P&L statement data |
| `/v1/businesses/:businessDescriptor/analytics/reports/cash-flow` | GET | Cash flow statement data |

### Query Parameters:

- `asOf` (optional): `YYYY-MM-DD` format, defaults to today

### Response Types:

See `backend/internal/domain/analytics/model.go` for authoritative shapes:
- `FinancialPosition`
- `ProfitAndLossStatement`
- `CashFlowStatement`

## 5) Data Model & Migrations

**No changes required.** This is a read-only feature using existing backend data.

## 6) Security & Privacy

- **Tenant scoping**: All endpoints are business-scoped via `:businessDescriptor` route param
- **RBAC**: Both `admin` and `member` roles can view reports (`role.ActionView` on `role.ResourceFinancialReports`)
- **Abuse prevention**: N/A (read-only, no write operations)

## 7) Observability & KPIs

### Events to Track (future implementation):

| Event | Properties | Trigger |
|-------|------------|---------|
| `reports.hub.viewed` | businessId, hasData | Hub page loaded |
| `reports.health.viewed` | businessId, asOf | Health page loaded |
| `reports.profit.viewed` | businessId, asOf | Profit page loaded |
| `reports.cashflow.viewed` | businessId, asOf | CashFlow page loaded |

### KPI Expectations:

- 60%+ of active businesses view reports weekly
- 40%+ view multiple reports per session

## 8) Test Strategy

### What Should Be Tested:

- **Manual testing**: All 4 pages render correctly with real data
- **State testing**: Empty, loading, error states for each page
- **RTL testing**: All pages in Arabic locale
- **Navigation testing**: Sidebar links, back buttons, quick links
- **Query param persistence**: `asOf` param persists across navigation

### Future E2E Tests (not in scope for this plan):

- Report data accuracy against backend
- Date filter functionality
- Cross-browser rendering

## 9) Risks & Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| API response shape mismatch | Low | Medium | Types defined from backend model.go |
| Performance on large datasets | Low | Low | Queries use staleTime caching |
| RTL rendering issues | Low | Medium | Use logical properties, test in Arabic |
| Chart rendering issues | Low | Low | Use existing chart components |

## 10) Definition of Done

- [ ] Meets BRD acceptance criteria (Section 8 of BRD)
- [ ] Mobile-first UX verified on iPhone SE viewport
- [ ] RTL/i18n parity verified (all keys in en/ar)
- [ ] Multi-tenancy verified (business-scoped data only)
- [ ] Error handling complete (loading/empty/error states)
- [ ] No TODO/FIXME comments in code
- [ ] TypeScript compiles without errors
- [ ] ESLint passes without errors
- [ ] Build succeeds without errors
- [ ] Navigation updated (Sidebar includes Reports)
- [ ] All 4 pages accessible and functional
