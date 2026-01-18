/**
 * Cash Movement Page
 *
 * Professional Statement of Cash Flows following proper accounting format:
 * - Hero: Cash Now (ending balance)
 * - Single comprehensive cash flow statement with three activity categories:
 *   - Operating Activities (core business operations)
 *   - Investing Activities (equipment/assets)
 *   - Financing Activities (owner capital)
 * - Proper indentation, subtotals, and visual hierarchy
 * - Insights with contextual advice
 * - Quick links to capital and expenses
 *
 * Responsive design: Mobile-first with elegant layout.
 * RTL-compatible.
 * Route: /business/$businessDescriptor/reports/cashflow
 */
import {
  Link,
  useNavigate,
  useParams,
  useRouteContext,
  useSearch,
} from '@tanstack/react-router'
import {
  ArrowLeft,
  Briefcase,
  Building2,
  ChevronLeft,
  ChevronRight,
  DollarSign,
  TrendingDown,
  TrendingUp,
} from 'lucide-react'
import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'

import { AsOfDatePicker } from './AsOfDatePicker'

import type { AdvisorInsight } from '@/components'
import { useCashFlowQuery } from '@/api/accounting'
import { AdvisorPanel, Button, Skeleton } from '@/components'
import { useLanguage } from '@/hooks/useLanguage'
import { cn } from '@/lib/utils'
import { formatCurrency } from '@/lib/formatCurrency'

export function CashMovementPage() {
  const { t } = useTranslation('reports')
  const { t: tCommon } = useTranslation('common')
  const { isRTL } = useLanguage()
  const navigate = useNavigate()
  const { businessDescriptor } = useParams({
    from: '/business/$businessDescriptor/reports/cashflow',
  })
  const search = useSearch({
    from: '/business/$businessDescriptor/reports/cashflow',
  })
  const asOf = search.asOf
  const routeContext = useRouteContext({
    from: '/business/$businessDescriptor',
  })
  const business = routeContext.business

  const ChevronIcon = isRTL ? ChevronLeft : ChevronRight

  // Data query
  const { data, isLoading, error } = useCashFlowQuery(businessDescriptor, asOf)

  const currency = business.currency

  // Parse amounts from decimal strings to numbers
  const cashAtStart = parseFloat(data?.cashAtStart ?? '0')
  const cashAtEnd = parseFloat(data?.cashAtEnd ?? '0')
  const cashFromCustomers = parseFloat(data?.cashFromCustomers ?? '0')
  const cashFromOwner = parseFloat(data?.cashFromOwner ?? '0')
  const inventoryPurchases = parseFloat(data?.inventoryPurchases ?? '0')
  const operatingExpenses = parseFloat(data?.operatingExpenses ?? '0')
  const businessInvestments = parseFloat(data?.businessInvestments ?? '0')
  const ownerDraws = parseFloat(data?.ownerDraws ?? '0')
  const netCashFlow = parseFloat(data?.netCashFlow ?? '0')

  // Calculate net cash flows by category for the comprehensive statement
  const netOperatingCash =
    cashFromCustomers - (inventoryPurchases + operatingExpenses)
  const netInvestingCash = -businessInvestments
  const netFinancingCash = cashFromOwner - ownerDraws

  // Build insights for advisor panel
  const advisorInsights = useMemo<Array<AdvisorInsight>>(() => {
    const insights: Array<AdvisorInsight> = []

    if (netCashFlow >= 0) {
      insights.push({ type: 'positive', message: t('insights.cash_healthy') })
    } else {
      insights.push({ type: 'alert', message: t('insights.cash_alert') })
    }

    // If owner withdrew more than what came in from customers
    if (ownerDraws > cashFromCustomers && cashFromCustomers > 0) {
      insights.push({
        type: 'suggestion',
        message: t('insights.withdrawal_tip'),
      })
    }

    return insights
  }, [netCashFlow, ownerDraws, cashFromCustomers, t])

  // Navigate back to Reports Hub
  const handleBack = () => {
    void navigate({
      to: '/business/$businessDescriptor/reports',
      params: { businessDescriptor },
      search: asOf ? { asOf } : undefined,
    })
  }

  // Error state
  if (error) {
    return (
      <div className="flex flex-col items-center justify-center px-4 py-16 text-center">
        <h2 className="mb-2 text-lg font-semibold">{t('error.title')}</h2>
        <p className="mb-4 text-base-content/60">{t('error.body')}</p>
        <Button variant="ghost" onClick={() => window.location.reload()}>
          {tCommon('retry')}
        </Button>
      </div>
    )
  }

  return (
    <div className="space-y-6 pb-20 md:pb-6">
      {/* Header */}
      <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
        <div className="flex items-center gap-3">
          <button
            type="button"
            onClick={handleBack}
            className="btn btn-ghost btn-sm btn-square"
            aria-label={tCommon('back')}
          >
            <ArrowLeft className={cn('h-5 w-5', isRTL && 'rotate-180')} />
          </button>
          <div>
            <h1 className="text-xl font-bold md:text-2xl">
              {t('cashflow.title')}
            </h1>
            <p className="text-sm text-base-content/60">
              {t('cashflow.subtitle')}
            </p>
          </div>
        </div>
        <AsOfDatePicker
          asOf={asOf}
          routeTo="/business/$businessDescriptor/reports/cashflow"
          businessDescriptor={businessDescriptor}
          className="self-end sm:self-auto"
        />
      </div>

      {/* Hero: Cash Now (Minimal) */}
      <section
        className="rounded-box bg-base-200/50 p-6 text-center"
        aria-labelledby="cash-now-heading"
      >
        <h2 id="cash-now-heading" className="sr-only">
          {t('cashflow.cash_now')}
        </h2>
        {isLoading ? (
          <>
            <Skeleton className="mx-auto mb-2 h-4 w-32" />
            <Skeleton className="mx-auto mb-1 h-10 w-48" />
            <Skeleton className="mx-auto h-3 w-40" />
          </>
        ) : (
          <>
            <p className="mb-1 text-sm text-base-content/60">
              {t('cashflow.cash_now')}
            </p>
            <p
              className={cn(
                'text-4xl font-bold tabular-nums md:text-5xl',
                cashAtEnd >= 0 ? 'text-success' : 'text-error',
              )}
            >
              {cashAtEnd >= 0 ? '' : '-'}
              {formatCurrency(Math.abs(cashAtEnd), currency)}
            </p>
            <p className="mt-1 text-xs text-base-content/50">
              {t('cashflow.cash_runway')}
            </p>
          </>
        )}
      </section>

      {/* Advisor Panel - Mobile (after hero) */}
      {!isLoading && advisorInsights.length > 0 && (
        <div className="lg:hidden">
          <AdvisorPanel
            title={t('insights.title')}
            insights={advisorInsights}
          />
        </div>
      )}

      {/* Responsive Grid: Statement + Advisor Panel sidebar on desktop */}
      <div className="grid gap-6 lg:grid-cols-[1fr_320px]">
        {/* Main Content: Statement of Cash Flows */}
        <div className="space-y-6">
          {/* Comprehensive Cash Flow Statement */}
          <section
            aria-labelledby="statement-heading"
            className="rounded-box border border-base-300 bg-base-100"
          >
            <div className="border-b border-base-300 bg-base-200/50 px-6 py-4">
              <h2 id="statement-heading" className="text-lg font-semibold">
                {t('cashflow.statement_title')}
              </h2>
            </div>

            {isLoading ? (
              <div className="space-y-3 p-6">
                {[1, 2, 3, 4, 5, 6, 7, 8, 9, 10].map((i) => (
                  <div key={i} className="flex items-center justify-between">
                    <Skeleton className="h-4 w-32" />
                    <Skeleton className="h-4 w-24" />
                  </div>
                ))}
              </div>
            ) : (
              <div className="divide-y divide-base-300">
                {/* Opening Balance */}
                <div className="flex items-center justify-between px-6 py-4">
                  <span className="font-medium">
                    {t('cashflow.cash_start')}
                  </span>
                  <span className="font-medium tabular-nums">
                    {formatCurrency(cashAtStart, currency)}
                  </span>
                </div>

                {/* OPERATING ACTIVITIES */}
                <div className="space-y-2 px-6 py-4">
                  <div className="mb-3 flex items-center gap-2">
                    <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10">
                      <Briefcase
                        className="h-4 w-4 text-primary"
                        aria-hidden="true"
                      />
                    </div>
                    <h3 className="font-semibold uppercase text-primary">
                      {t('cashflow.operating_activities')}
                    </h3>
                  </div>

                  {/* Operating line items with indentation */}
                  <div className="space-y-2 ps-10">
                    <div className="flex items-center justify-between text-sm">
                      <span className="text-base-content/80">
                        {t('cashflow.from_customers')}
                      </span>
                      <span
                        className={cn(
                          'tabular-nums',
                          cashFromCustomers > 0 && 'text-success',
                        )}
                      >
                        {cashFromCustomers > 0 && '+'}
                        {formatCurrency(cashFromCustomers, currency)}
                      </span>
                    </div>
                    <div className="flex items-center justify-between text-sm">
                      <span className="text-base-content/80">
                        {t('cashflow.inventory_purchases')}
                      </span>
                      <span
                        className={cn(
                          'tabular-nums',
                          inventoryPurchases > 0 && 'text-error',
                        )}
                      >
                        {inventoryPurchases > 0 && '-'}
                        {formatCurrency(inventoryPurchases, currency)}
                      </span>
                    </div>
                    <div className="flex items-center justify-between text-sm">
                      <span className="text-base-content/80">
                        {t('cashflow.running_costs')}
                      </span>
                      <span
                        className={cn(
                          'tabular-nums',
                          operatingExpenses > 0 && 'text-error',
                        )}
                      >
                        {operatingExpenses > 0 && '-'}
                        {formatCurrency(operatingExpenses, currency)}
                      </span>
                    </div>
                  </div>

                  {/* Operating subtotal */}
                  <div className="mt-2 flex items-center justify-between border-t border-base-300 pt-2 font-semibold">
                    <span>{t('cashflow.net_cash_from_operations')}</span>
                    <span
                      className={cn(
                        'tabular-nums',
                        netOperatingCash >= 0 ? 'text-success' : 'text-error',
                      )}
                    >
                      {netOperatingCash >= 0 ? '+' : ''}
                      {formatCurrency(netOperatingCash, currency)}
                    </span>
                  </div>
                </div>

                {/* INVESTING ACTIVITIES */}
                <div className="space-y-2 px-6 py-4">
                  <div className="mb-3 flex items-center gap-2">
                    <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-purple-500/10">
                      <Building2
                        className="h-4 w-4 text-purple-500"
                        aria-hidden="true"
                      />
                    </div>
                    <h3 className="font-semibold uppercase text-purple-500">
                      {t('cashflow.investing_activities')}
                    </h3>
                  </div>

                  {/* Investing line items with indentation */}
                  <div className="space-y-2 ps-10">
                    <div className="flex items-center justify-between text-sm">
                      <span className="text-base-content/80">
                        {t('cashflow.equipment_assets')}
                      </span>
                      <span
                        className={cn(
                          'tabular-nums',
                          businessInvestments > 0 && 'text-error',
                        )}
                      >
                        {businessInvestments > 0 && '-'}
                        {formatCurrency(businessInvestments, currency)}
                      </span>
                    </div>
                  </div>

                  {/* Investing subtotal */}
                  <div className="mt-2 flex items-center justify-between border-t border-base-300 pt-2 font-semibold">
                    <span>{t('cashflow.net_cash_from_investing')}</span>
                    <span
                      className={cn(
                        'tabular-nums',
                        netInvestingCash >= 0 ? 'text-success' : 'text-error',
                      )}
                    >
                      {netInvestingCash >= 0 ? '+' : ''}
                      {formatCurrency(netInvestingCash, currency)}
                    </span>
                  </div>
                </div>

                {/* FINANCING ACTIVITIES */}
                <div className="space-y-2 px-6 py-4">
                  <div className="mb-3 flex items-center gap-2">
                    <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-success/10">
                      <DollarSign
                        className="h-4 w-4 text-success"
                        aria-hidden="true"
                      />
                    </div>
                    <h3 className="font-semibold uppercase text-success">
                      {t('cashflow.financing_activities')}
                    </h3>
                  </div>

                  {/* Financing line items with indentation */}
                  <div className="space-y-2 ps-10">
                    <div className="flex items-center justify-between text-sm">
                      <span className="text-base-content/80">
                        {t('cashflow.from_owner')}
                      </span>
                      <span
                        className={cn(
                          'tabular-nums',
                          cashFromOwner > 0 && 'text-success',
                        )}
                      >
                        {cashFromOwner > 0 && '+'}
                        {formatCurrency(cashFromOwner, currency)}
                      </span>
                    </div>
                    <div className="flex items-center justify-between text-sm">
                      <span className="text-base-content/80">
                        {t('cashflow.to_owner')}
                      </span>
                      <span
                        className={cn(
                          'tabular-nums',
                          ownerDraws > 0 && 'text-error',
                        )}
                      >
                        {ownerDraws > 0 && '-'}
                        {formatCurrency(ownerDraws, currency)}
                      </span>
                    </div>
                  </div>

                  {/* Financing subtotal */}
                  <div className="mt-2 flex items-center justify-between border-t border-base-300 pt-2 font-semibold">
                    <span>{t('cashflow.net_cash_from_financing')}</span>
                    <span
                      className={cn(
                        'tabular-nums',
                        netFinancingCash >= 0 ? 'text-success' : 'text-error',
                      )}
                    >
                      {netFinancingCash >= 0 ? '+' : ''}
                      {formatCurrency(netFinancingCash, currency)}
                    </span>
                  </div>
                </div>

                {/* NET CHANGE IN CASH */}
                <div
                  className={cn(
                    'flex items-center justify-between px-6 py-4 font-bold',
                    netCashFlow >= 0
                      ? 'bg-success/5 text-success'
                      : 'bg-error/5 text-error',
                  )}
                >
                  <div className="flex items-center gap-2">
                    <span className="uppercase">
                      {t('cashflow.net_change')}
                    </span>
                    {netCashFlow >= 0 ? (
                      <TrendingUp className="h-5 w-5" aria-hidden="true" />
                    ) : (
                      <TrendingDown className="h-5 w-5" aria-hidden="true" />
                    )}
                  </div>
                  <span className="text-lg tabular-nums">
                    {netCashFlow >= 0 ? '+' : ''}
                    {formatCurrency(netCashFlow, currency)}
                  </span>
                </div>

                {/* Closing Balance */}
                <div className="flex items-center justify-between bg-base-200/50 px-6 py-4 font-bold">
                  <span>{t('cashflow.cash_end')}</span>
                  <span
                    className={cn(
                      'text-lg tabular-nums',
                      cashAtEnd >= 0 ? 'text-success' : 'text-error',
                    )}
                  >
                    {formatCurrency(cashAtEnd, currency)}
                  </span>
                </div>
              </div>
            )}
          </section>

          {/* Quick Links */}
          <section aria-label={t('quick_links.view_capital')}>
            <div className="flex flex-wrap gap-4">
              <Link
                to="/business/$businessDescriptor/accounting/capital"
                params={{ businessDescriptor }}
                className="flex items-center gap-1 text-sm font-medium text-primary hover:underline"
              >
                {t('quick_links.view_capital')}
                <ChevronIcon className="h-4 w-4" aria-hidden="true" />
              </Link>
              <Link
                to="/business/$businessDescriptor/accounting/expenses"
                params={{ businessDescriptor }}
                search={{
                  page: 1,
                  pageSize: 20,
                  sortBy: 'occurredOn',
                  sortOrder: 'desc',
                }}
                className="flex items-center gap-1 text-sm font-medium text-primary hover:underline"
              >
                {t('quick_links.view_expenses')}
                <ChevronIcon className="h-4 w-4" aria-hidden="true" />
              </Link>
            </div>
          </section>
        </div>

        {/* Advisor Panel - Desktop Sidebar (sticky with header clearance) */}
        {!isLoading && advisorInsights.length > 0 && (
          <aside className="hidden lg:block">
            <div className="sticky top-20">
              <AdvisorPanel
                title={t('insights.title')}
                insights={advisorInsights}
              />
            </div>
          </aside>
        )}
      </div>
    </div>
  )
}

CashMovementPage.displayName = 'CashMovementPage'
