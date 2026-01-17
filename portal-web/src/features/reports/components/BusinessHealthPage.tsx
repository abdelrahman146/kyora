/**
 * Business Health Page
 *
 * Shows financial position (balance sheet) in plain language:
 * - Hero: What Your Business Is Worth (Total Equity = Assets - Liabilities)
 * - Section: What You Own (Assets breakdown with visual bar)
 * - Section: What You Owe (Liabilities - currently $0 note)
 * - Section: Owner's Stake (Equity breakdown)
 * - Current Assets breakdown (Cash + Inventory)
 * - Insight Card with contextual advice
 * - Quick Links to related pages
 *
 * Responsive design: Mobile-first with 2-column desktop layout.
 * RTL-compatible.
 * Route: /business/$businessDescriptor/reports/health
 *
 * **Important**: totalEquity is the "business worth" (assets - liabilities).
 * cashOnHand is an approximation that can be NEGATIVE (see accounting.instructions.md).
 */
import {
  Link,
  useNavigate,
  useParams,
  useRouteContext,
  useSearch,
} from '@tanstack/react-router'
import {
  AlertTriangle,
  ArrowLeft,
  Building,
  ChevronLeft,
  ChevronRight,
  Package,
  Wallet,
} from 'lucide-react'
import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'

import { AsOfDatePicker } from './AsOfDatePicker'
import { AssetBreakdownBar } from './AssetBreakdownBar'

import type { AdvisorInsight } from '@/components'
import { useFinancialPositionQuery } from '@/api/accounting'
import { AdvisorPanel, Button, Skeleton } from '@/components'
import { useLanguage } from '@/hooks/useLanguage'
import { cn } from '@/lib/utils'
import { formatCurrency } from '@/lib/formatCurrency'

export function BusinessHealthPage() {
  const { t } = useTranslation('reports')
  const { t: tCommon } = useTranslation('common')
  const { isRTL } = useLanguage()
  const navigate = useNavigate()
  const { businessDescriptor } = useParams({
    from: '/business/$businessDescriptor/reports/health',
  })
  const search = useSearch({
    from: '/business/$businessDescriptor/reports/health',
  })
  const asOf = search.asOf
  const routeContext = useRouteContext({
    from: '/business/$businessDescriptor',
  })
  const business = routeContext.business

  const ChevronIcon = isRTL ? ChevronLeft : ChevronRight

  // Data query
  const { data, isLoading, error } = useFinancialPositionQuery(
    businessDescriptor,
    asOf,
  )

  const currency = business.currency

  // Parse amounts from decimal strings to numbers
  // totalEquity = Assets - Liabilities = "What Your Business Is Worth"
  const totalEquity = parseFloat(data?.totalEquity ?? '0')
  // totalAssets = currentAssets + fixedAssets = "Everything the business owns"
  const totalAssets = parseFloat(data?.totalAssets ?? '0')
  const totalLiabilities = parseFloat(data?.totalLiabilities ?? '0')
  // cashOnHand is an APPROXIMATION that can be NEGATIVE
  // Formula: (Revenue + OwnerInvestment) - (Expenses + OwnerDraws + FixedAssets + Inventory)
  const cashOnHand = parseFloat(data?.cashOnHand ?? '0')
  const isCashNegative = cashOnHand < 0
  // currentAssets = Cash + Inventory (liquid assets) - not displayed separately
  const inventoryValue = parseFloat(data?.totalInventoryValue ?? '0')
  const fixedAssets = parseFloat(data?.fixedAssets ?? '0')
  const ownerInvestment = parseFloat(data?.ownerInvestment ?? '0')
  const ownerDraws = parseFloat(data?.ownerDraws ?? '0')
  const retainedEarnings = parseFloat(data?.retainedEarnings ?? '0')

  // Build insights for advisor panel
  const advisorInsights = useMemo<Array<AdvisorInsight>>(() => {
    if (totalEquity < 0) {
      return [
        {
          type: 'alert',
          message: t('insights.negative_equity'),
        },
      ]
    }
    if (isCashNegative) {
      return [
        {
          type: 'suggestion',
          message: t('insights.negative_cash'),
        },
      ]
    }
    return [{ type: 'positive', message: t('insights.healthy') }]
  }, [totalEquity, isCashNegative, t])

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
              {t('health.title')}
            </h1>
            <p className="text-sm text-base-content/60">
              {t('health.subtitle')}
            </p>
          </div>
        </div>
        <AsOfDatePicker
          asOf={asOf}
          routeTo="/business/$businessDescriptor/reports/health"
          businessDescriptor={businessDescriptor}
          className="self-end sm:self-auto"
        />
      </div>

      {/* Responsive Grid: Main content + Advisor Panel sidebar on desktop */}
      <div className="grid gap-6 lg:grid-cols-[1fr_320px]">
        {/* Main Content */}
        <div className="space-y-6">
          {/* Hero: Business Worth - Full width on all screens */}
          <section
            className="rounded-box bg-base-200/50 p-6 text-center"
            aria-labelledby="business-worth-heading"
          >
            <h2 id="business-worth-heading" className="sr-only">
              {t('metrics.business_worth')}
            </h2>
            {isLoading ? (
              <>
                <Skeleton className="mx-auto mb-2 h-4 w-32" />
                <Skeleton className="mx-auto h-10 w-48" />
              </>
            ) : (
              <>
                <p className="mb-1 text-sm text-base-content/60">
                  {t('metrics.business_worth')}
                </p>
                <p
                  className={cn(
                    'text-4xl font-bold tabular-nums md:text-5xl',
                    totalEquity < 0 ? 'text-error' : 'text-success',
                  )}
                >
                  {formatCurrency(totalEquity, currency)}
                </p>
                <p className="mt-2 text-xs text-base-content/50">
                  {t('health.business_value_formula')}
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

          {/* Desktop: Two-column layout for key metrics */}
          <div className="grid gap-6 lg:grid-cols-2">
            {/* Section: What You Own */}
            <section aria-labelledby="what-you-own-heading">
              <h2
                id="what-you-own-heading"
                className="mb-4 text-lg font-semibold"
              >
                {t('health.what_you_own')}
              </h2>
              {isLoading ? (
                <div className="space-y-3 rounded-box border border-base-300 p-4">
                  <Skeleton className="h-4 w-full rounded-full" />
                  <Skeleton className="h-4 w-3/4" />
                  <Skeleton className="h-4 w-2/3" />
                  <Skeleton className="h-4 w-1/2" />
                </div>
              ) : (
                <div className="rounded-box border border-base-300 p-4">
                  <AssetBreakdownBar
                    segments={[
                      {
                        label: t('health.cash_on_hand'),
                        value: Math.max(0, cashOnHand), // Bar only shows positive portion
                        color: 'success',
                      },
                      {
                        label: t('health.inventory_value'),
                        value: inventoryValue,
                        color: 'info',
                      },
                      {
                        label: t('health.fixed_assets'),
                        value: fixedAssets,
                        color: 'secondary',
                      },
                    ]}
                    total={totalAssets}
                    currency={currency}
                    showLegend={false}
                  />
                  {/* Detailed breakdown with actual values */}
                  <div className="mt-4 space-y-2 border-t border-base-300 pt-3">
                    <div className="flex items-center justify-between text-sm">
                      <div className="flex items-center gap-2">
                        <Wallet
                          className={cn(
                            'h-4 w-4',
                            isCashNegative ? 'text-error' : 'text-success',
                          )}
                        />
                        <span>{t('health.cash_on_hand')}</span>
                      </div>
                      <span
                        className={cn(
                          'font-medium tabular-nums',
                          isCashNegative && 'text-error',
                        )}
                      >
                        {formatCurrency(cashOnHand, currency)}
                      </span>
                    </div>
                    <div className="flex items-center justify-between text-sm">
                      <div className="flex items-center gap-2">
                        <Package className="h-4 w-4 text-info" />
                        <span>{t('health.inventory_value')}</span>
                      </div>
                      <span className="font-medium tabular-nums">
                        {formatCurrency(inventoryValue, currency)}
                      </span>
                    </div>
                    <div className="flex items-center justify-between text-sm">
                      <div className="flex items-center gap-2">
                        <Building className="h-4 w-4 text-secondary" />
                        <span>{t('health.fixed_assets')}</span>
                      </div>
                      <span className="font-medium tabular-nums">
                        {formatCurrency(fixedAssets, currency)}
                      </span>
                    </div>
                  </div>
                  <div className="mt-4 flex justify-between border-t border-base-300 pt-3 font-semibold">
                    <span>{t('health.total_assets')}</span>
                    <span className="tabular-nums">
                      {formatCurrency(totalAssets, currency)}
                    </span>
                  </div>
                  {/* Warning for negative cash */}
                  {isCashNegative && (
                    <div className="mt-4 flex items-start gap-2 rounded-box bg-warning/10 p-3">
                      <AlertTriangle className="h-4 w-4 shrink-0 text-warning" />
                      <p className="text-xs text-base-content/70">
                        {t('health.note_negative_cash')}
                      </p>
                    </div>
                  )}
                </div>
              )}
            </section>

            {/* Right column on desktop: Liabilities + Owner's Stake */}
            <div className="space-y-6">
              {/* Section: What You Owe */}
              <section aria-labelledby="what-you-owe-heading">
                <h2
                  id="what-you-owe-heading"
                  className="mb-4 text-lg font-semibold"
                >
                  {t('health.what_you_owe')}
                </h2>
                <div className="rounded-box border border-base-300 p-4">
                  <div className="mb-2 flex items-center justify-between">
                    <span>{t('health.total_liabilities')}</span>
                    <span className="text-xl font-bold tabular-nums">
                      {isLoading ? (
                        <Skeleton className="h-6 w-24" />
                      ) : (
                        formatCurrency(totalLiabilities, currency)
                      )}
                    </span>
                  </div>
                  <p className="text-sm text-base-content/60">
                    {t('health.liabilities_note')}
                  </p>
                </div>
              </section>

              {/* Section: Owner's Stake */}
              <section aria-labelledby="owners-stake-heading">
                <h2
                  id="owners-stake-heading"
                  className="mb-4 text-lg font-semibold"
                >
                  {t('health.owners_stake')}
                </h2>
                {isLoading ? (
                  <div className="space-y-3 rounded-box border border-base-300 p-4">
                    <Skeleton className="h-6 w-full" />
                    <Skeleton className="h-6 w-full" />
                    <Skeleton className="h-6 w-full" />
                    <Skeleton className="h-8 w-full" />
                  </div>
                ) : (
                  <div className="space-y-3 rounded-box border border-base-300 p-4">
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
                          retainedEarnings >= 0 ? 'text-success' : 'text-error',
                        )}
                      >
                        {formatCurrency(retainedEarnings, currency)}
                      </span>
                    </div>
                    {/* Explanatory note - equity components don't directly sum to business value */}
                    <div className="mt-3 rounded-box bg-base-200/50 p-3 text-xs text-base-content/60">
                      <p>{t('health.equity_explanation')}</p>
                    </div>
                  </div>
                )}
              </section>
            </div>
          </div>

          {/* Quick Links */}
          <section aria-label={t('quick_links.view_assets')}>
            <div className="flex flex-wrap gap-4">
              <Link
                to="/business/$businessDescriptor/accounting/assets"
                params={{ businessDescriptor }}
                className="flex items-center gap-1 text-sm font-medium text-primary hover:underline"
              >
                {t('quick_links.view_assets')}
                <ChevronIcon className="h-4 w-4" aria-hidden="true" />
              </Link>
              <Link
                to="/business/$businessDescriptor/accounting/capital"
                params={{ businessDescriptor }}
                className="flex items-center gap-1 text-sm font-medium text-primary hover:underline"
              >
                {t('quick_links.view_capital')}
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

BusinessHealthPage.displayName = 'BusinessHealthPage'
