/**
 * Accounting Dashboard
 *
 * Landing page for the Accounting module showing:
 * - "Safe to Draw" hero statistic
 * - Summary stats (Total Assets, Total Expenses)
 * - Recent Activity list (mixed: expenses, investments, withdrawals)
 * - Quick navigation to sub-pages (Expenses, Capital, Assets)
 */

import { Link, useParams, useRouteContext } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import {
  ArrowDownRight,
  ArrowUpRight,
  Box,
  Boxes,
  ChevronLeft,
  ChevronRight,
  PiggyBank,
  Plus,
  Receipt,
  Repeat,
  Wallet,
} from 'lucide-react'

import { categoryIcons } from '../schema/options'

import type { RecentActivity } from '@/api/accounting'

import {
  useAccountingSummaryQuery,
  useRecentActivitiesQuery,
} from '@/api/accounting'
import { StatCard, StatCardSkeleton } from '@/components'
import { StatCardGroup } from '@/components/molecules/StatCardGroup'
import { formatCurrency } from '@/lib/formatCurrency'
import { formatDateShort } from '@/lib/formatDate'
import { useLanguage } from '@/hooks/useLanguage'

export function AccountingDashboard() {
  const { t } = useTranslation('accounting')
  const { businessDescriptor } = useParams({
    from: '/business/$businessDescriptor/accounting/',
  })
  const { business } = useRouteContext({
    from: '/business/$businessDescriptor',
  })
  const { isRTL } = useLanguage()

  // Fetch accounting summary
  const { data: summary, isLoading: isSummaryLoading } =
    useAccountingSummaryQuery(businessDescriptor)

  // Fetch recent activities (mixed: expenses, investments, withdrawals)
  const { data: recentActivitiesData, isLoading: isActivitiesLoading } =
    useRecentActivitiesQuery(businessDescriptor, { limit: 5 })

  const recentActivities = recentActivitiesData?.items ?? []
  const currency = summary?.currency ?? business.currency

  // Parse amounts from string to number
  const safeToDrawAmount = parseFloat(summary?.safeToDrawAmount ?? '0')
  const totalAssets = parseFloat(summary?.totalAssetValue ?? '0')
  const totalExpenses = parseFloat(summary?.totalExpenses ?? '0')

  // Determine if safe to draw is negative (warning)
  const isSafeToDrawNegative = safeToDrawAmount < 0

  return (
    <div className="space-y-6">
      {/* Summary Stats Section */}
      <section aria-labelledby="summary-heading">
        <h2 id="summary-heading" className="sr-only">
          {t('header.dashboard')}
        </h2>

        <StatCardGroup cols={3}>
          {/* Safe to Draw - Hero Stat */}
          {isSummaryLoading ? (
            <StatCardSkeleton />
          ) : (
            <StatCard
              label={t('stats.safe_to_draw')}
              value={formatCurrency(Math.abs(safeToDrawAmount), currency)}
              icon={<PiggyBank className="h-5 w-5 text-primary" />}
              variant={isSafeToDrawNegative ? 'error' : 'success'}
              trend={isSafeToDrawNegative ? 'down' : undefined}
              trendValue={
                isSafeToDrawNegative
                  ? t('helper.exceeds_safe_amount')
                  : undefined
              }
            />
          )}

          {/* Total Assets */}
          {isSummaryLoading ? (
            <StatCardSkeleton />
          ) : (
            <StatCard
              label={t('stats.total_assets')}
              value={formatCurrency(totalAssets, currency)}
              icon={<Box className="h-5 w-5 text-info" />}
              variant="info"
            />
          )}

          {/* Total Expenses */}
          {isSummaryLoading ? (
            <StatCardSkeleton />
          ) : (
            <StatCard
              label={t('stats.total_expenses')}
              value={formatCurrency(totalExpenses, currency)}
              icon={<Receipt className="h-5 w-5 text-warning" />}
              variant="warning"
            />
          )}
        </StatCardGroup>
      </section>

      {/* Quick Navigation Cards */}
      <section aria-labelledby="quick-nav-heading">
        <h2 id="quick-nav-heading" className="text-lg font-semibold mb-4">
          {t('nav.title')}
        </h2>

        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <QuickNavCard
            to="/business/$businessDescriptor/accounting/expenses"
            params={{ businessDescriptor }}
            search={{
              page: 1,
              pageSize: 20,
              sortBy: 'occurredOn',
              sortOrder: 'desc',
            }}
            icon={<Receipt className="h-6 w-6" />}
            title={t('header.expenses')}
            isRTL={isRTL}
          />
          <QuickNavCard
            to="/business/$businessDescriptor/accounting/capital"
            params={{ businessDescriptor }}
            icon={<Wallet className="h-6 w-6" />}
            title={t('header.capital')}
            isRTL={isRTL}
          />
          <QuickNavCard
            to="/business/$businessDescriptor/accounting/assets"
            params={{ businessDescriptor }}
            icon={<Box className="h-6 w-6" />}
            title={t('header.assets')}
            isRTL={isRTL}
          />
        </div>
      </section>

      {/* Recent Activity Section */}
      <section aria-labelledby="recent-activity-heading">
        <div className="flex items-center justify-between mb-4">
          <h2 id="recent-activity-heading" className="text-lg font-semibold">
            {t('list.recent_activity')}
          </h2>
          {/* <Link
            to="/business/$businessDescriptor/accounting/expenses"
            params={{ businessDescriptor }}
            className="btn btn-ghost btn-sm gap-1"
          >
            {t('list.view_all')}
            {isRTL ? (
              <ChevronLeft className="h-4 w-4" />
            ) : (
              <ChevronRight className="h-4 w-4" />
            )}
          </Link> */}
        </div>

        {isActivitiesLoading ? (
          <RecentActivitySkeleton />
        ) : recentActivities.length === 0 ? (
          <EmptyRecentActivity businessDescriptor={businessDescriptor} />
        ) : (
          <div className="card bg-base-100 border border-base-300">
            <ul className="divide-y divide-base-300">
              {recentActivities.map((activity) => (
                <ActivityListItem
                  key={activity.id}
                  activity={activity}
                  currency={currency}
                />
              ))}
            </ul>
          </div>
        )}
      </section>
    </div>
  )
}

// =============================================================================
// Sub-components
// =============================================================================

interface QuickNavCardProps {
  to: string
  params: { businessDescriptor: string }
  search?: Record<string, unknown>
  icon: React.ReactNode
  title: string
  isRTL: boolean
}

function QuickNavCard({
  to,
  params,
  search,
  icon,
  title,
  isRTL,
}: QuickNavCardProps) {
  return (
    <Link
      to={to}
      params={params}
      search={search}
      className="card bg-base-100 border border-base-300 hover:border-primary hover:bg-base-200 transition-colors"
    >
      <div className="card-body p-4 flex-row items-center gap-4">
        <div className="flex items-center justify-center h-12 w-12 rounded-full bg-primary/10 text-primary">
          {icon}
        </div>
        <span className="font-medium flex-1">{title}</span>
        {isRTL ? (
          <ChevronLeft className="h-5 w-5 text-base-content/40" />
        ) : (
          <ChevronRight className="h-5 w-5 text-base-content/40" />
        )}
      </div>
    </Link>
  )
}

interface ActivityListItemProps {
  activity: RecentActivity
  currency: string
}

function ActivityListItem({ activity, currency }: ActivityListItemProps) {
  const { t } = useTranslation('accounting')
  const amount = parseFloat(activity.amount)

  // Determine icon, color, and label based on activity type
  const getActivityDetails = () => {
    switch (activity.type) {
      case 'expense': {
        const CategoryIcon = activity.category
          ? categoryIcons[activity.category]
          : Receipt
        return {
          icon: <CategoryIcon className="h-5 w-5 text-error" />,
          bgColor: 'bg-error/10',
          amountColor: 'text-error',
          amountPrefix: '-',
          label: activity.category
            ? t(`category.${activity.category}`)
            : t('activity.expense'),
          badgeText:
            activity.expenseType === 'recurring' ? (
              <Repeat className="w-3 h-3" />
            ) : null,
        }
      }
      case 'investment':
        return {
          icon: <ArrowDownRight className="h-5 w-5 text-success" />,
          bgColor: 'bg-success/10',
          amountColor: 'text-success',
          amountPrefix: '+',
          label: t('activity.investment'),
          badgeText: null,
        }
      case 'withdrawal':
        return {
          icon: <ArrowUpRight className="h-5 w-5 text-warning" />,
          bgColor: 'bg-warning/10',
          amountColor: 'text-warning',
          amountPrefix: '-',
          label: t('activity.withdrawal'),
          badgeText: null,
        }
      case 'asset':
        return {
          icon: <Boxes className="h-5 w-5 text-success" />,
          bgColor: 'bg-success/10',
          amountColor: 'text-success',
          amountPrefix: '+',
          label: t('activity.asset'),
          badgeText: null,
        }
      default:
        return {
          icon: <Receipt className="h-5 w-5 text-base-content/60" />,
          bgColor: 'bg-base-300',
          amountColor: 'text-base-content',
          amountPrefix: '',
          label: t('activity.unknown'),
          badgeText: null,
        }
    }
  }

  const details = getActivityDetails()

  return (
    <li className="flex items-center gap-4 p-4">
      <div
        className={`flex items-center justify-center h-10 w-10 rounded-full ${details.bgColor}`}
      >
        {details.icon}
      </div>
      <div className="flex-1 min-w-0">
        <p className="font-medium truncate">{details.label}</p>
        <p className="text-sm text-base-content/60 truncate">
          {activity.description !== details.label
            ? activity.description
            : t('list.occurred_on', {
                date: formatDateShort(activity.occurredAt),
              })}
        </p>
      </div>
      {details.badgeText && (
        <span className="badge badge-ghost badge-sm gap-1 text-secondary">
          {details.badgeText}
        </span>
      )}
      <div className="text-end">
        <p className={`font-semibold tabular-nums ${details.amountColor}`}>
          {details.amountPrefix}
          {formatCurrency(amount, currency)}
        </p>
        <p className="text-xs text-base-content/60">
          {formatDateShort(activity.occurredAt)}
        </p>
      </div>
    </li>
  )
}

function RecentActivitySkeleton() {
  return (
    <div className="card bg-base-100 border border-base-300 animate-pulse">
      <div className="divide-y divide-base-300">
        {[1, 2, 3].map((i) => (
          <div key={i} className="flex items-center gap-4 p-4">
            <div className="h-10 w-10 rounded-full bg-base-300" />
            <div className="flex-1 space-y-2">
              <div className="h-4 w-24 rounded bg-base-300" />
              <div className="h-3 w-32 rounded bg-base-300" />
            </div>
            <div className="space-y-2 text-end">
              <div className="h-4 w-20 rounded bg-base-300 ms-auto" />
              <div className="h-3 w-16 rounded bg-base-300 ms-auto" />
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}

interface EmptyRecentActivityProps {
  businessDescriptor: string
}

function EmptyRecentActivity({ businessDescriptor }: EmptyRecentActivityProps) {
  const { t } = useTranslation('accounting')

  return (
    <div className="card bg-base-200/50 border border-base-300">
      <div className="card-body items-center text-center py-12">
        <div className="flex items-center justify-center h-16 w-16 rounded-full bg-base-300/50 mb-4">
          <Receipt className="h-8 w-8 text-base-content/40" />
        </div>
        <h3 className="font-semibold text-lg">
          {t('empty.recent_activity_title')}
        </h3>
        <p className="text-base-content/60 max-w-sm">
          {t('empty.recent_activity_description')}
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
          className="btn btn-primary btn-sm mt-4 gap-2"
        >
          <Plus className="h-4 w-4" />
          {t('actions.add_first_expense')}
        </Link>
      </div>
    </div>
  )
}
