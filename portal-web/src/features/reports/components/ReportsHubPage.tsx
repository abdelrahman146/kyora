/**
 * Reports Hub Page
 *
 * Redesigned landing page for Financial Reports featuring:
 * - Financial Overview Hero: 3 prominent cards (Business Value, Net Profit, Cash Position)
 * - Your Financial Reports: 3 informative cards with descriptions and key metrics
 * - Mobile-first responsive design with elegant spacing
 * - RTL-compatible with proper semantic HTML
 * - All data fetched via TanStack Query with prefetching in route loader
 *
 * Design Principles:
 * - Generous white space and breathing room
 * - Clear visual hierarchy (hero metrics → report cards)
 * - Informative descriptions explaining what each report shows
 * - Color-coded metrics (green/red/blue for context)
 * - Professional financial dashboard feel
 *
 * Related routes:
 * - /business/$businessDescriptor/reports/health
 * - /business/$businessDescriptor/reports/profit
 * - /business/$businessDescriptor/reports/cashflow
 */
import {
  Link,
  useParams,
  useRouteContext,
  useSearch,
} from '@tanstack/react-router'
import {
  Activity,
  ArrowRight,
  BarChart3,
  TrendingDown,
  TrendingUp,
  Wallet,
} from 'lucide-react'
import { useTranslation } from 'react-i18next'

import {
  useAccountingSummaryQuery,
  useCashFlowQuery,
  useFinancialPositionQuery,
  useProfitAndLossQuery,
} from '@/api/accounting'
import { Button } from '@/components'
import { formatCurrency } from '@/lib/formatCurrency'
import { cn } from '@/lib/utils'
import { isRTL } from '@/lib/charts'

export function ReportsHubPage() {
  const { t } = useTranslation('reports')
  const { t: tCommon } = useTranslation('common')
  const { businessDescriptor } = useParams({
    from: '/business/$businessDescriptor/reports/',
  })
  const search = useSearch({
    from: '/business/$businessDescriptor/reports/',
  })
  const asOf = search.asOf
  const routeContext = useRouteContext({
    from: '/business/$businessDescriptor',
  })
  const business = routeContext.business

  // Data queries - all prefetched in route loader
  const {
    data: summary,
    isLoading: isSummaryLoading,
    error: summaryError,
  } = useAccountingSummaryQuery(businessDescriptor)

  const { data: financialPosition, isLoading: isFPLoading } =
    useFinancialPositionQuery(businessDescriptor, asOf)

  const { data: profitLoss, isLoading: isPLLoading } = useProfitAndLossQuery(
    businessDescriptor,
    asOf,
  )

  const { data: cashFlow, isLoading: isCFLoading } = useCashFlowQuery(
    businessDescriptor,
    asOf,
  )

  const isLoading =
    isSummaryLoading || isFPLoading || isPLLoading || isCFLoading
  const currency = summary?.currency ?? business.currency

  // Financial Position metrics
  const totalEquity = parseFloat(financialPosition?.totalEquity ?? '0')
  const totalAssets = parseFloat(financialPosition?.totalAssets ?? '0')

  // Profit metrics
  const netProfit = parseFloat(profitLoss?.netProfit ?? '0')
  const grossProfit = parseFloat(profitLoss?.grossProfit ?? '0')
  const revenue = parseFloat(profitLoss?.revenue ?? '0')
  const totalExpenses = parseFloat(profitLoss?.totalExpenses ?? '0')

  // Cash flow metrics
  const cashAtEnd = parseFloat(cashFlow?.cashAtEnd ?? '0')
  const totalCashIn = parseFloat(cashFlow?.totalCashIn ?? '0')
  const netCashFlow = parseFloat(cashFlow?.netCashFlow ?? '0')

  // Check if business has any data to show
  const hasData = revenue > 0 || totalExpenses > 0 || totalAssets > 0

  // Error state
  if (summaryError) {
    return (
      <div className="flex flex-col items-center justify-center px-4 py-16 text-center">
        <div className="mb-4 text-error">
          <Activity className="mx-auto h-12 w-12" aria-hidden="true" />
        </div>
        <h2 className="mb-2 text-lg font-semibold">{tCommon('error.title')}</h2>
        <p className="mb-4 max-w-md text-base-content/60">
          {tCommon('error.description')}
        </p>
        <Button variant="ghost" onClick={() => window.location.reload()}>
          {tCommon('retry')}
        </Button>
      </div>
    )
  }

  // Empty state (new business, no data yet)
  if (!isLoading && !hasData) {
    return (
      <div className="flex flex-col items-center justify-center px-4 py-16 text-center">
        <div className="mb-4 text-primary/40">
          <BarChart3 className="mx-auto h-16 w-16" aria-hidden="true" />
        </div>
        <h2 className="mb-2 text-lg font-semibold">
          {tCommon('empty.no_data')}
        </h2>
        <p className="mb-6 max-w-md text-base-content/60">
          {tCommon('empty.start_recording')}
        </p>
        <Link
          to="/business/$businessDescriptor/orders"
          params={{ businessDescriptor }}
          className="btn btn-primary"
        >
          {tCommon('cta.create_first_order')}
        </Link>
      </div>
    )
  }

  return (
    <div className="space-y-8 pb-20 md:pb-8">
      {/* Header with As-of Date */}
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex-1">
          <h1 className="text-2xl font-bold md:text-3xl">{t('hub.title')}</h1>
          <p className="mt-2 text-sm text-base-content/60 md:text-base">
            {t('hub.description')}
          </p>
        </div>
      </div>

      {/* Financial Overview Hero Section */}
      <section aria-labelledby="financial-overview-heading">
        <h2
          id="financial-overview-heading"
          className="mb-4 text-lg font-semibold text-base-content/80"
        >
          {t('hub.financial_overview')}
        </h2>

        {isLoading ? (
          <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
            {[1, 2, 3].map((i) => (
              <div
                key={i}
                className="rounded-box animate-pulse border border-base-300 bg-base-200/30 p-4"
              >
                <div className="mb-3 h-4 w-32 rounded bg-base-300" />
                <div className="mb-2 h-10 w-40 rounded bg-base-300" />
                <div className="h-3 w-24 rounded bg-base-300" />
              </div>
            ))}
          </div>
        ) : (
          <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
            {/* Business Value Card */}
            <div
              className={cn(
                'rounded-box border p-4',
                totalEquity >= 0
                  ? 'border-success/30 bg-success/5'
                  : 'border-error/30 bg-error/5',
              )}
            >
              <div className="flex items-start gap-3">
                <div
                  className={cn(
                    'rounded-box p-2',
                    totalEquity >= 0 ? 'bg-success/10' : 'bg-error/10',
                  )}
                >
                  <Activity
                    className={cn(
                      'h-5 w-5',
                      totalEquity >= 0 ? 'text-success' : 'text-error',
                    )}
                    aria-hidden="true"
                  />
                </div>
                <div className="flex-1">
                  <p className="text-xs font-medium text-base-content/60">
                    {t('metrics.business_value')}
                  </p>
                  <p
                    className={cn(
                      'mt-1 text-2xl font-bold tabular-nums',
                      totalEquity >= 0 ? 'text-success' : 'text-error',
                    )}
                  >
                    {formatCurrency(totalEquity, currency)}
                  </p>
                  <p className="mt-1 text-xs text-base-content/50">
                    {t('metrics.total_assets')}:{' '}
                    {formatCurrency(totalAssets, currency)}
                  </p>
                </div>
              </div>
            </div>

            {/* Net Profit Card */}
            <div
              className={cn(
                'rounded-box border p-4',
                netProfit >= 0
                  ? 'border-success/30 bg-success/5'
                  : 'border-error/30 bg-error/5',
              )}
            >
              <div className="flex items-start gap-3">
                <div
                  className={cn(
                    'rounded-box p-2',
                    netProfit >= 0 ? 'bg-success/10' : 'bg-error/10',
                  )}
                >
                  {netProfit >= 0 ? (
                    <TrendingUp
                      className="h-5 w-5 text-success"
                      aria-hidden="true"
                    />
                  ) : (
                    <TrendingDown
                      className="h-5 w-5 text-error"
                      aria-hidden="true"
                    />
                  )}
                </div>
                <div className="flex-1">
                  <p className="text-xs font-medium text-base-content/60">
                    {t('metrics.net_profit')}
                  </p>
                  <p
                    className={cn(
                      'mt-1 text-2xl font-bold tabular-nums',
                      netProfit >= 0 ? 'text-success' : 'text-error',
                    )}
                  >
                    {formatCurrency(netProfit, currency)}
                  </p>
                  <p className="mt-1 text-xs text-base-content/50">
                    {t('metrics.gross_profit')}:{' '}
                    {formatCurrency(grossProfit, currency)}
                  </p>
                </div>
              </div>
            </div>

            {/* Cash Position Card */}
            <div
              className={cn(
                'rounded-box border p-4',
                cashAtEnd >= 0
                  ? 'border-success/30 bg-success/5'
                  : 'border-error/30 bg-error/5',
              )}
            >
              <div className="flex items-start gap-3">
                <div
                  className={cn(
                    'rounded-box p-2',
                    cashAtEnd >= 0 ? 'bg-success/10' : 'bg-error/10',
                  )}
                >
                  <Wallet
                    className={cn(
                      'h-5 w-5',
                      cashAtEnd >= 0 ? 'text-success' : 'text-error',
                    )}
                    aria-hidden="true"
                  />
                </div>
                <div className="flex-1">
                  <p className="text-xs font-medium text-base-content/60">
                    {t('metrics.cash_position')}
                  </p>
                  <p
                    className={cn(
                      'mt-1 text-2xl font-bold tabular-nums',
                      cashAtEnd >= 0 ? 'text-success' : 'text-error',
                    )}
                  >
                    {formatCurrency(cashAtEnd, currency)}
                  </p>
                  <p className="mt-1 text-xs text-base-content/50">
                    {t('metrics.total_cash_in')}:{' '}
                    {formatCurrency(totalCashIn, currency)}
                  </p>
                </div>
              </div>
            </div>
          </div>
        )}
      </section>

      {/* Your Financial Reports Section */}
      <section aria-labelledby="reports-heading">
        <h2
          id="reports-heading"
          className="mb-4 text-lg font-semibold text-base-content/80"
        >
          {t('hub.your_reports')}
        </h2>

        {isLoading ? (
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {[1, 2, 3].map((i) => (
              <div
                key={i}
                className="rounded-box animate-pulse border border-base-300 bg-base-200/30 p-6"
              >
                <div className="mb-3 h-6 w-48 rounded bg-base-300" />
                <div className="mb-4 space-y-2">
                  <div className="h-4 w-full rounded bg-base-300" />
                  <div className="h-4 w-5/6 rounded bg-base-300" />
                </div>
                <div className="mb-4 space-y-3">
                  <div className="h-10 w-full rounded bg-base-300" />
                  <div className="h-10 w-full rounded bg-base-300" />
                </div>
                <div className="h-10 w-full rounded bg-base-300" />
              </div>
            ))}
          </div>
        ) : (
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {/* Business Health Report Card */}
            <article className="group rounded-box border border-base-300 bg-base-100 p-6 transition-all hover:border-primary/30">
              <div className="mb-4 flex items-start gap-3">
                <div className="rounded-box bg-primary/10 p-2.5">
                  <Activity
                    className="h-5 w-5 text-primary"
                    aria-hidden="true"
                  />
                </div>
                <div>
                  <h3 className="text-lg font-semibold">
                    {t('cards.business_health')}
                  </h3>
                </div>
              </div>

              <p className="mb-6 text-sm leading-relaxed text-base-content/70">
                {t('cards.business_health_desc')}
              </p>

              <div className="mb-6 space-y-3 rounded-lg bg-base-200/50 p-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-base-content/70">
                    {t('metrics.total_assets')}
                  </span>
                  <span className="font-semibold tabular-nums">
                    {formatCurrency(totalAssets, currency)}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-base-content/70">
                    {t('metrics.business_value')}
                  </span>
                  <span
                    className={cn(
                      'font-semibold tabular-nums',
                      totalEquity >= 0 ? 'text-success' : 'text-error',
                    )}
                  >
                    {formatCurrency(totalEquity, currency)}
                  </span>
                </div>
              </div>

              <Link
                to="/business/$businessDescriptor/reports/health"
                params={{ businessDescriptor }}
                search={asOf ? { asOf } : undefined}
                className="btn btn-ghost btn-block group-hover:btn-primary"
              >
                {t('hub.view_full_report')}
                <ArrowRight
                  className={`h-4 w-4 ${isRTL() ? 'rotate-180' : ''}`}
                  aria-hidden="true"
                />
              </Link>
            </article>
            {/* Profit & Loss Report Card */}
            <article className="group rounded-box border border-base-300 bg-base-100 p-6 transition-all hover:border-primary/30">
              <div className="mb-4 flex items-start gap-3">
                <div
                  className={cn(
                    'rounded-box p-2.5',
                    netProfit >= 0 ? 'bg-success/10' : 'bg-error/10',
                  )}
                >
                  {netProfit >= 0 ? (
                    <TrendingUp
                      className="h-5 w-5 text-success"
                      aria-hidden="true"
                    />
                  ) : (
                    <TrendingDown
                      className="h-5 w-5 text-error"
                      aria-hidden="true"
                    />
                  )}
                </div>
                <div>
                  <h3 className="text-lg font-semibold">
                    {t('cards.profit_loss')}
                  </h3>
                </div>
              </div>

              <p className="mb-6 text-sm leading-relaxed text-base-content/70">
                {t('cards.profit_loss_desc')}
              </p>

              <div className="mb-6 space-y-3 rounded-lg bg-base-200/50 p-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-base-content/70">
                    {t('metrics.total_revenue')}
                  </span>
                  <span className="font-semibold tabular-nums">
                    {formatCurrency(revenue, currency)}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-base-content/70">
                    {t('metrics.net_profit')}
                  </span>
                  <span
                    className={cn(
                      'font-semibold tabular-nums',
                      netProfit >= 0 ? 'text-success' : 'text-error',
                    )}
                  >
                    {formatCurrency(netProfit, currency)}
                  </span>
                </div>
              </div>

              <Link
                to="/business/$businessDescriptor/reports/profit"
                params={{ businessDescriptor }}
                search={asOf ? { asOf } : undefined}
                className="btn btn-ghost btn-block group-hover:btn-primary"
              >
                {t('hub.view_full_report')}
                <ArrowRight
                  className={`h-4 w-4 ${isRTL() ? 'rotate-180' : ''}`}
                  aria-hidden="true"
                />
              </Link>
            </article>

            {/* Cash Movement Report Card */}
            <article className="group rounded-box border border-base-300 bg-base-100 p-6 transition-all hover:border-primary/30">
              <div className="mb-4 flex items-start gap-3">
                <div className="rounded-box bg-info/10 p-2.5">
                  <Wallet className="h-5 w-5 text-info" aria-hidden="true" />
                </div>
                <div>
                  <h3 className="text-lg font-semibold">
                    {t('cards.cash_movement')}
                  </h3>
                </div>
              </div>

              <p className="mb-6 text-sm leading-relaxed text-base-content/70">
                {t('cards.cash_movement_desc')}
              </p>

              <div className="mb-6 space-y-3 rounded-lg bg-base-200/50 p-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-base-content/70">
                    {t('metrics.total_cash_in')}
                  </span>
                  <span className="font-semibold tabular-nums">
                    {formatCurrency(totalCashIn, currency)}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-base-content/70">
                    {t('metrics.net_cash_flow')}
                  </span>
                  <span
                    className={cn(
                      'font-semibold tabular-nums',
                      netCashFlow >= 0 ? 'text-success' : 'text-error',
                    )}
                  >
                    {netCashFlow >= 0 ? '+' : ''}
                    {formatCurrency(netCashFlow, currency)}
                  </span>
                </div>
              </div>

              <Link
                to="/business/$businessDescriptor/reports/cashflow"
                params={{ businessDescriptor }}
                search={asOf ? { asOf } : undefined}
                className="btn btn-ghost btn-block group-hover:btn-primary"
              >
                {t('hub.view_full_report')}
                <ArrowRight
                  className={`h-4 w-4 ${isRTL() ? 'rotate-180' : ''}`}
                  aria-hidden="true"
                />
              </Link>
            </article>
          </div>
        )}
      </section>
    </div>
  )
}

ReportsHubPage.displayName = 'ReportsHubPage'
