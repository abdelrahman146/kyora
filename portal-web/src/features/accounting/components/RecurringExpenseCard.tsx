/**
 * RecurringExpenseCard Component
 *
 * Mobile-first card for displaying recurring expense templates.
 * Shows category icon, amount, frequency, status, and next occurrence.
 *
 * Features:
 * - Category-based icon mapping
 * - Amount highlighted in red (expense = money out)
 * - Status badge (active/paused/ended/canceled)
 * - Frequency badge
 * - Next occurrence date
 * - Quick actions dropdown for edit/delete/status change
 * - RTL-compatible with logical properties
 * - Mobile-optimized touch targets
 */

import { useTranslation } from 'react-i18next'
import { Calendar, Receipt, Repeat } from 'lucide-react'

import { categoryColors, categoryIcons } from '../schema/options'
import { RecurringExpenseQuickActions } from './RecurringExpenseQuickActions'
import type { RecurringExpense } from '@/api/accounting'
import { formatCurrency } from '@/lib/formatCurrency'
import { formatDateShort } from '@/lib/formatDate'

export interface RecurringExpenseCardProps {
  recurringExpense: RecurringExpense
  currency: string
  businessDescriptor: string
  onActionComplete?: () => void
}

const statusColors: Record<string, string> = {
  active: 'badge-success',
  paused: 'badge-warning',
  ended: 'badge-ghost',
  canceled: 'badge-error',
}

export function RecurringExpenseCard({
  recurringExpense,
  currency,
  businessDescriptor,
  onActionComplete,
}: RecurringExpenseCardProps) {
  const { t } = useTranslation('accounting')

  // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
  const CategoryIcon = categoryIcons[recurringExpense.category] ?? Receipt
  const colorClass =
    // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
    categoryColors[recurringExpense.category] ??
    'bg-base-200 text-base-content/70'

  const amount = parseFloat(recurringExpense.amount)
  const isActive = recurringExpense.status === 'active'

  return (
    <div
      className={`w-full bg-base-100 border border-base-300 rounded-xl overflow-hidden ${
        !isActive ? 'opacity-60' : ''
      }`}
    >
      {/* Header: Category Icon, Name, Status, Actions */}
      <div className="flex items-start gap-3 p-4 bg-base-200/30">
        {/* Category Icon */}
        <div
          className={`w-12 h-12 rounded-xl flex items-center justify-center flex-shrink-0 ${colorClass}`}
        >
          <CategoryIcon className="w-6 h-6" />
        </div>

        {/* Category Name + Status */}
        <div className="flex-1 min-w-0">
          <h3 className="font-semibold text-base text-base-content truncate">
            {t(`category.${recurringExpense.category}`)}
          </h3>
          <div className="flex items-center gap-2 mt-1">
            {/* Status badge */}
            <span
              className={`badge badge-sm ${statusColors[recurringExpense.status]}`}
            >
              {t(`status.${recurringExpense.status}`)}
            </span>
          </div>
        </div>

        {/* Actions */}
        <div className="flex-shrink-0 -mt-1">
          <RecurringExpenseQuickActions
            recurringExpense={recurringExpense}
            businessDescriptor={businessDescriptor}
            onActionComplete={onActionComplete}
          />
        </div>
      </div>

      {/* Amount Section */}
      <div className="px-4 py-3 bg-error/5 border-y border-base-200">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Repeat className="w-4 h-4 text-base-content/50" />
            <span className="text-sm text-base-content/70 font-medium">
              {t(`frequency.${recurringExpense.frequency}`)}
            </span>
          </div>
          <div className="text-end">
            <span className="font-bold text-xl text-error">
              -{formatCurrency(amount, currency)}
            </span>
            <span className="text-xs text-base-content/50 ms-1">
              /{t(`frequency.${recurringExpense.frequency}`).toLowerCase()}
            </span>
          </div>
        </div>
      </div>

      {/* Details Section */}
      <div className="p-4 space-y-2">
        {/* Next occurrence */}
        {recurringExpense.nextRecurringDate && isActive && (
          <div className="flex items-center gap-2 text-sm">
            <Calendar className="w-4 h-4 text-base-content/40 flex-shrink-0" />
            <span className="text-base-content/70">
              {t('recurring.next_on', {
                date: formatDateShort(recurringExpense.nextRecurringDate),
              })}
            </span>
          </div>
        )}

        {/* Start date */}
        <div className="flex items-center gap-2 text-xs text-base-content/50">
          <span className="w-4 h-4 flex-shrink-0" />
          <span>
            {t('recurring.started_on', {
              date: formatDateShort(recurringExpense.recurringStartDate),
            })}
          </span>
        </div>

        {/* End date (if exists) */}
        {recurringExpense.recurringEndDate && (
          <div className="flex items-center gap-2 text-xs text-base-content/50">
            <span className="w-4 h-4 flex-shrink-0" />
            <span>
              {t('recurring.ends_on', {
                date: formatDateShort(recurringExpense.recurringEndDate),
              })}
            </span>
          </div>
        )}

        {/* Note */}
        {recurringExpense.note && (
          <div className="pt-2 mt-2 border-t border-base-200">
            <p className="text-sm text-base-content/60 whitespace-pre-wrap leading-relaxed">
              {recurringExpense.note}
            </p>
          </div>
        )}
      </div>
    </div>
  )
}
