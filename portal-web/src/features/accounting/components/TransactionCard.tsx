/**
 * TransactionCard Component
 *
 * Mobile-first card for displaying capital transactions (investments & withdrawals).
 * Shows transaction type, amount, date, and optional note.
 *
 * Features:
 * - Color-coded by type (Green for investments, Red for withdrawals)
 * - Directional icons (arrow-down for money in, arrow-up for money out)
 * - Amount formatting with currency
 * - Date with localized formatting
 * - Quick actions dropdown for edit/delete
 * - RTL-compatible with logical properties
 * - Mobile-optimized touch targets
 */

import { useTranslation } from 'react-i18next'
import { ArrowDownCircle, ArrowUpCircle } from 'lucide-react'

import { TransactionQuickActions } from './TransactionQuickActions'
import type { Investment, Withdrawal } from '@/api/accounting'
import { formatCurrency } from '@/lib/formatCurrency'
import { formatDateShort } from '@/lib/formatDate'

export interface TransactionCardProps {
  transaction: Investment | Withdrawal
  type: 'investment' | 'withdrawal'
  currency: string
  businessDescriptor: string
  onActionComplete?: () => void
}

export function TransactionCard({
  transaction,
  type,
  currency,
  businessDescriptor,
  onActionComplete,
}: TransactionCardProps) {
  const { t } = useTranslation('accounting')

  const amount = parseFloat(transaction.amount)
  const date =
    type === 'investment'
      ? (transaction as Investment).investedAt
      : (transaction as Withdrawal).withdrawnAt

  const isInvestment = type === 'investment'
  const Icon = isInvestment ? ArrowDownCircle : ArrowUpCircle
  const colorClass = isInvestment
    ? 'bg-success/10 text-success'
    : 'bg-error/10 text-error'
  const amountColorClass = isInvestment ? 'text-success' : 'text-error'

  return (
    <div className="card bg-base-100 border border-base-300">
      <div className="card-body p-4 flex flex-row items-start gap-3">
        {/* Icon */}
        <div className={`rounded-full p-2 ${colorClass} shrink-0`}>
          <Icon className="h-5 w-5" />
        </div>

        {/* Content */}
        <div className="flex-1 min-w-0 space-y-1">
          {/* Type & Date */}
          <div className="flex items-center justify-between gap-2">
            <h3 className="font-medium text-base-content">
              {t(`type.${type}`)}
            </h3>
            <span className={`font-semibold ${amountColorClass} text-sm`}>
              {isInvestment ? '+' : '-'}
              {formatCurrency(amount, currency)}
            </span>
          </div>

          {/* Date */}
          <p className="text-xs text-base-content/60">
            {t(isInvestment ? 'list.invested_on' : 'list.withdrawn_on', {
              date: formatDateShort(date),
            })}
          </p>

          {/* Note (if present) */}
          {transaction.note && (
            <p className="text-sm text-base-content/70 line-clamp-2">
              {transaction.note}
            </p>
          )}
        </div>

        {/* Quick Actions Menu */}
        <TransactionQuickActions
          transaction={transaction}
          type={type}
          businessDescriptor={businessDescriptor}
          onActionComplete={onActionComplete}
        />
      </div>
    </div>
  )
}
