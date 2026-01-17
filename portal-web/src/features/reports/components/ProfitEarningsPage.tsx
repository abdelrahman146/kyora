/**
 * Profit & Earnings Page
 *
 * Shows where money comes from and where it goes:
 * - Hero: What You Keep (Net Profit) with key metrics on desktop
 * - Waterfall: Revenue → COGS → Gross Profit → Expenses → Net Profit
 * - Margin stats (Gross Margin, Profit Margin)
 * - Running costs breakdown (doughnut chart + detailed list)
 * - Insights with contextual advice
 * - Quick links to expenses and orders
 *
 * Responsive design: Mobile-first with 2-column desktop layout.
 * RTL-compatible.
 * Route: /business/$businessDescriptor/reports/profit
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
  ChevronLeft,
  ChevronRight,
  DollarSign,
  Minus,
  Percent,
  ShoppingCart,
  TrendingDown,
  TrendingUp,
} from 'lucide-react'
import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'

import { AsOfDatePicker } from './AsOfDatePicker'
import { ProfitWaterfall } from './ProfitWaterfall'

import type { WaterfallStep } from './ProfitWaterfall'
import type { ExpenseCategory } from '@/api/accounting'

import type { AdvisorInsight } from '@/components'
import { useProfitAndLossQuery } from '@/api/accounting'
import { AdvisorPanel, Button, Skeleton, StatCard } from '@/components'
import { useLanguage } from '@/hooks/useLanguage'
import { cn } from '@/lib/utils'
import { formatCurrency } from '@/lib/formatCurrency'
import {
  categoryColors,
  categoryIcons,
} from '@/features/accounting/schema/options'

export function ProfitEarningsPage() {
  const { t } = useTranslation('reports')
  const { t: tCommon } = useTranslation('common')
  const { t: tAccounting } = useTranslation('accounting')
  const { isRTL } = useLanguage()
  const navigate = useNavigate()
  const { businessDescriptor } = useParams({
    from: '/business/$businessDescriptor/reports/profit',
  })
  const search = useSearch({
    from: '/business/$businessDescriptor/reports/profit',
  })
  const asOf = search.asOf
  const routeContext = useRouteContext({
    from: '/business/$businessDescriptor',
  })
  const business = routeContext.business

  const ChevronIcon = isRTL ? ChevronLeft : ChevronRight

  // Data query
  const { data, isLoading, error } = useProfitAndLossQuery(
    businessDescriptor,
    asOf,
  )

  const currency = business.currency

  // Parse amounts from decimal strings to numbers
  const revenue = parseFloat(data?.revenue ?? '0')
  const cogs = parseFloat(data?.cogs ?? '0')
  const grossProfit = parseFloat(data?.grossProfit ?? '0')
  const totalExpenses = parseFloat(data?.totalExpenses ?? '0')
  const netProfit = parseFloat(data?.netProfit ?? '0')

  // Calculate margins
  const grossMargin = revenue > 0 ? (grossProfit / revenue) * 100 : 0
  const profitMargin = revenue > 0 ? (netProfit / revenue) * 100 : 0

  // Prepare waterfall steps
  const waterfallSteps: Array<WaterfallStep> = useMemo(
    () => [
      { label: t('profit.revenue'), value: revenue, type: 'total' },
      { label: t('profit.cogs'), value: cogs, type: 'subtract' },
      { label: t('profit.gross_profit'), value: grossProfit, type: 'result' },
      {
        label: t('profit.operating_expenses'),
        value: totalExpenses,
        type: 'subtract',
      },
      {
        label: t('profit.net_profit_result'),
        value: netProfit,
        type: 'result',
      },
    ],
    [t, revenue, cogs, grossProfit, totalExpenses, netProfit],
  )

  // Expense breakdown list with icons and colors
  const expenseBreakdownList = useMemo(() => {
    const expensesByCategory = data?.expensesByCategory
    if (!expensesByCategory?.length) return []

    return expensesByCategory
      .map((exp) => {
        const category = exp.Key as ExpenseCategory
        const value = parseFloat(exp.Value)
        return {
          key: category,
          label: tAccounting(`category.${category}`, {
            defaultValue: category,
          }),
          value,
          Icon: categoryIcons[category],
          colorClass: categoryColors[category],
          percent: totalExpenses > 0 ? (value / totalExpenses) * 100 : 0,
        }
      })
      .filter((c) => c.value > 0)
      .sort((a, b) => b.value - a.value)
  }, [data?.expensesByCategory, totalExpenses, tAccounting])

  // Find biggest expense for insight
  const biggestExpense = useMemo(() => {
    if (expenseBreakdownList.length === 0) return null

    const biggest = expenseBreakdownList[0]
    return {
      category: biggest.label,
      amount: formatCurrency(biggest.value, currency),
      percent: Math.round(biggest.percent),
    }
  }, [expenseBreakdownList, currency])

  // Build insights for advisor panel
  const advisorInsights = useMemo<Array<AdvisorInsight>>(() => {
    const insights: Array<AdvisorInsight> = []

    // Profitability insight
    if (netProfit > 0) {
      insights.push({
        type: 'positive',
        message: t('insights.profitable', {
          currency,
          margin: Math.round(profitMargin),
        }),
      })
    } else if (netProfit < 0) {
      insights.push({
        type: 'alert',
        message: t('insights.losing', {
          amount: formatCurrency(Math.abs(netProfit), currency),
        }),
      })
    } else {
      insights.push({
        type: 'info',
        message: t('insights.profitable', { currency, margin: 0 }),
      })
    }

    // Biggest expense insight
    if (biggestExpense) {
      insights.push({
        type: 'suggestion',
        message: t('insights.biggest_expense', {
          category: biggestExpense.category,
          amount: biggestExpense.amount,
          percent: biggestExpense.percent,
        }),
      })
    }

    return insights
  }, [netProfit, profitMargin, currency, biggestExpense, t])

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
              {t('profit.title')}
            </h1>
            <p className="text-sm text-base-content/60">
              {t('profit.subtitle')}
            </p>
          </div>
        </div>
        <AsOfDatePicker
          asOf={asOf}
          routeTo="/business/$businessDescriptor/reports/profit"
          businessDescriptor={businessDescriptor}
          className="self-end sm:self-auto"
        />
      </div>

      {/* Hero + Key Stats Grid - Full Width */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* Hero: Net Profit */}
        <section
          className="rounded-box bg-base-200/50 p-6 text-center lg:flex lg:flex-col lg:justify-center"
          aria-labelledby="net-profit-heading"
        >
          <h2 id="net-profit-heading" className="sr-only">
            {t('profit.net_profit')}
          </h2>
          {isLoading ? (
            <>
              <Skeleton className="mx-auto mb-2 h-4 w-32" />
              <Skeleton className="mx-auto h-10 w-48" />
            </>
          ) : (
            <>
              <p className="mb-1 text-sm text-base-content/60">
                {t('profit.net_profit')}
              </p>
              <p
                className={cn(
                  'text-4xl font-bold tabular-nums md:text-5xl',
                  netProfit < 0 ? 'text-error' : 'text-success',
                )}
              >
                {netProfit >= 0 ? '' : '-'}
                {formatCurrency(Math.abs(netProfit), currency)}
              </p>
              <p className="mt-1 text-xs text-base-content/50">
                {t('profit.net_profit_subtitle')}
              </p>
            </>
          )}
        </section>

        {/* Key Stats - Desktop */}
        <div className="hidden lg:grid lg:grid-cols-2 lg:gap-4">
          <StatCard
            label={t('profit.revenue')}
            value={formatCurrency(revenue, currency)}
            icon={<DollarSign className="h-4 w-4" />}
            variant="default"
          />
          <StatCard
            label={t('profit.cogs')}
            value={formatCurrency(cogs, currency)}
            icon={<ShoppingCart className="h-4 w-4" />}
            variant={cogs > 0 ? 'warning' : 'default'}
          />
          <StatCard
            label={t('profit.gross_profit')}
            value={formatCurrency(grossProfit, currency)}
            icon={
              grossProfit >= 0 ? (
                <TrendingUp className="h-4 w-4" />
              ) : (
                <TrendingDown className="h-4 w-4" />
              )
            }
            variant={grossProfit >= 0 ? 'success' : 'error'}
          />
          <StatCard
            label={t('profit.operating_expenses')}
            value={formatCurrency(totalExpenses, currency)}
            icon={<Minus className="h-4 w-4" />}
            variant={totalExpenses > 0 ? 'warning' : 'default'}
          />
        </div>
      </div>

      {/* Advisor Panel - Mobile (after hero) */}
      {!isLoading && advisorInsights.length > 0 && (
        <div className="lg:hidden">
          <AdvisorPanel
            title={t('insights.title')}
            insights={advisorInsights}
          />
        </div>
      )}

      {/* Responsive Grid: Detailed Content + Advisor Panel sidebar on desktop */}
      <div className="grid gap-6 lg:grid-cols-[1fr_320px]">
        {/* Main Content */}
        <div className="space-y-6">
          {/* Waterfall Visualization */}
          <section aria-labelledby="profit-waterfall-heading">
            <h2
              id="profit-waterfall-heading"
              className="mb-4 text-lg font-semibold"
            >
              {t('profit.statement_title')}
            </h2>
            {isLoading ? (
              <div className="space-y-4 rounded-box border border-base-300 p-4">
                {[1, 2, 3, 4, 5].map((i) => (
                  <div key={i} className="space-y-2">
                    <div className="flex justify-between">
                      <Skeleton className="h-4 w-24" />
                      <Skeleton className="h-4 w-20" />
                    </div>
                    <Skeleton className="h-3 w-full rounded-full" />
                  </div>
                ))}
              </div>
            ) : (
              <div className="rounded-box border border-base-300 p-4">
                <ProfitWaterfall steps={waterfallSteps} currency={currency} />
              </div>
            )}
          </section>

          {/* Two-column layout for margins + expense breakdown on desktop */}
          <div className="grid gap-6 lg:grid-cols-2">
            {/* Margin Stats */}
            <section aria-label={t('profit.margins_title')}>
              <h3 className="mb-4 text-lg font-semibold">
                {t('profit.margins_title')}
              </h3>
              {isLoading ? (
                <div className="space-y-4">
                  <Skeleton className="h-20 w-full rounded-box" />
                  <Skeleton className="h-20 w-full rounded-box" />
                </div>
              ) : (
                <div className="space-y-4">
                  <div className="rounded-box border border-base-300 p-4">
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="text-sm font-medium text-base-content">
                          {t('profit.gross_margin')}
                        </p>
                        <p
                          className={cn(
                            'text-2xl font-bold tabular-nums',
                            grossMargin > 0
                              ? 'text-success'
                              : 'text-base-content',
                          )}
                        >
                          {grossMargin.toFixed(1)}%
                        </p>
                      </div>
                      <Percent
                        className={cn(
                          'h-8 w-8',
                          grossMargin > 0
                            ? 'text-success/30'
                            : 'text-base-content/20',
                        )}
                      />
                    </div>
                    <p className="mt-3 text-xs text-base-content/60">
                      {t('profit.gross_margin_help')}
                    </p>
                    {grossMargin > 0 && (
                      <div className="mt-2 rounded-lg bg-base-200 px-3 py-2">
                        <p className="text-xs text-base-content/70">
                          {t('profit.margin_example', {
                            percent: grossMargin.toFixed(0),
                            currency,
                            amount: grossMargin.toFixed(0),
                          })}
                        </p>
                      </div>
                    )}
                  </div>
                  <div className="rounded-box border border-base-300 p-4">
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="text-sm font-medium text-base-content">
                          {t('profit.profit_margin')}
                        </p>
                        <p
                          className={cn(
                            'text-2xl font-bold tabular-nums',
                            profitMargin > 0
                              ? 'text-success'
                              : profitMargin < 0
                                ? 'text-error'
                                : 'text-base-content',
                          )}
                        >
                          {profitMargin.toFixed(1)}%
                        </p>
                      </div>
                      <Percent
                        className={cn(
                          'h-8 w-8',
                          profitMargin > 0
                            ? 'text-success/30'
                            : profitMargin < 0
                              ? 'text-error/30'
                              : 'text-base-content/20',
                        )}
                      />
                    </div>
                    <p className="mt-3 text-xs text-base-content/60">
                      {t('profit.profit_margin_help')}
                    </p>
                    {profitMargin !== 0 && (
                      <div className="mt-2 rounded-lg bg-base-200 px-3 py-2">
                        <p className="text-xs text-base-content/70">
                          {t('profit.margin_example', {
                            percent: Math.abs(profitMargin).toFixed(0),
                            currency,
                            amount: Math.abs(profitMargin).toFixed(0),
                          })}
                        </p>
                      </div>
                    )}
                  </div>
                </div>
              )}
            </section>

            {/* Operating Expenses Breakdown */}
            <section aria-labelledby="expense-breakdown-heading">
              <h3
                id="expense-breakdown-heading"
                className="mb-4 text-lg font-semibold"
              >
                {t('profit.expense_detail')}
              </h3>
              {isLoading ? (
                <div className="space-y-2">
                  <Skeleton className="h-16 w-full rounded-box" />
                  <Skeleton className="h-16 w-full rounded-box" />
                  <Skeleton className="h-16 w-full rounded-box" />
                </div>
              ) : expenseBreakdownList.length > 0 ? (
                <div className="rounded-box border border-base-300 divide-y divide-base-300">
                  {expenseBreakdownList.map((expense) => {
                    const Icon = expense.Icon
                    return (
                      <div
                        key={expense.key}
                        className="flex items-center justify-between gap-3 p-4"
                      >
                        <div className="flex items-center gap-3 min-w-0">
                          <div
                            className={cn(
                              'flex h-9 w-9 shrink-0 items-center justify-center rounded-lg',
                              expense.colorClass,
                            )}
                          >
                            <Icon className="h-4 w-4" aria-hidden="true" />
                          </div>
                          <span className="text-sm truncate">
                            {expense.label}
                          </span>
                        </div>
                        <div className="flex items-center gap-3 shrink-0">
                          <span className="text-xs text-base-content/50 tabular-nums">
                            {expense.percent.toFixed(1)}%
                          </span>
                          <span className="font-medium tabular-nums">
                            {formatCurrency(expense.value, currency)}
                          </span>
                        </div>
                      </div>
                    )
                  })}
                </div>
              ) : (
                <div className="rounded-box border border-dashed border-base-300 bg-base-200/30 p-8 text-center">
                  <TrendingDown
                    className="mx-auto mb-3 h-12 w-12 text-base-content/20"
                    aria-hidden="true"
                  />
                  <p className="text-sm text-base-content/60">
                    {t('profit.no_expenses')}
                  </p>
                  <Link
                    to="/business/$businessDescriptor/accounting/expenses"
                    params={{ businessDescriptor }}
                    search={{
                      page: 1,
                      pageSize: 20,
                      sortBy: 'occurredOn',
                      sortOrder: 'desc',
                    }}
                    className="btn btn-sm btn-ghost mt-4"
                  >
                    {t('quick_links.view_expenses')}
                  </Link>
                </div>
              )}
            </section>
          </div>

          {/* Quick Links */}
          <section aria-label={t('quick_links.view_expenses')}>
            <div className="flex flex-wrap gap-4">
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
              <Link
                to="/business/$businessDescriptor/orders"
                params={{ businessDescriptor }}
                className="flex items-center gap-1 text-sm font-medium text-primary hover:underline"
              >
                {t('quick_links.view_orders')}
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

ProfitEarningsPage.displayName = 'ProfitEarningsPage'
