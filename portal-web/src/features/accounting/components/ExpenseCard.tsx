/**
 * ExpenseCard Component
 *
 * Mobile-first card for displaying expenses in list view.
 * Shows category icon, amount, date, and optional recurring badge.
 *
 * Features:
 * - Category-based icon mapping
 * - Amount highlighted in red (expense = money out)
 * - Date with localized formatting
 * - Recurring badge for recurring expense occurrences
 * - Quick actions dropdown for edit/delete
 * - RTL-compatible with logical properties
 * - Mobile-optimized touch targets
 */

import { useTranslation } from 'react-i18next'
import { Receipt, Repeat } from 'lucide-react'

import { categoryColors, categoryIcons } from '../schema/options'
import { ExpenseQuickActions } from './ExpenseQuickActions'
import type { Expense } from '@/api/accounting'
import { formatCurrency } from '@/lib/formatCurrency'
import { formatDateShort } from '@/lib/formatDate'

export interface ExpenseCardProps {
  expense: Expense
  currency: string
  businessDescriptor: string
  onActionComplete?: () => void
}

export function ExpenseCard({
  expense,
  currency,
  businessDescriptor,
  onActionComplete,
}: ExpenseCardProps) {
  const { t } = useTranslation('accounting')

  const CategoryIcon = categoryIcons[expense.category] ?? Receipt
  const colorClass =
    categoryColors[expense.category] ?? 'bg-base-200 text-base-content/70'

  const amount = parseFloat(expense.amount)
  const isRecurring = !!expense.recurringExpenseId

  return (
    <div className="w-full bg-base-100 border border-base-300 rounded-xl p-4">
      {/* Top row: Icon, Category, Badge, Amount, Actions */}
      <div className="flex items-center gap-3">
        {/* Category Icon */}
        <div
          className={`w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0 ${colorClass}`}
        >
          <CategoryIcon className="w-5 h-5" />
        </div>

        {/* Category Name + Recurring Badge */}
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2">
            <span className="font-medium text-base-content">
              {t(`category.${expense.category}`)}
            </span>
            {isRecurring && (
              <span className="badge badge-ghost badge-sm gap-1 text-secondary">
                <Repeat className="w-3 h-3" />
              </span>
            )}
          </div>
          {/* Date */}
          <p className="text-xs text-base-content/50">
            {t('list.occurred_on', {
              date: formatDateShort(expense.occurredOn),
            })}
          </p>
        </div>

        {/* Amount */}
        <div className="text-end flex-shrink-0">
          <span className="font-bold text-lg text-error">
            -{formatCurrency(amount, currency)}
          </span>
        </div>

        {/* Actions */}
        <div className="flex-shrink-0">
          <ExpenseQuickActions
            expense={expense}
            businessDescriptor={businessDescriptor}
            currency={currency}
            onActionComplete={onActionComplete}
          />
        </div>
      </div>

      {/* Note (full width, no truncation) */}
      {expense.note && (
        <p className="text-sm text-base-content/60 mt-2 pt-2 border-t border-base-200 whitespace-pre-wrap">
          {expense.note}
        </p>
      )}
    </div>
  )
}
