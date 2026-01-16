/**
 * TransactionQuickActions Component
 *
 * Dropdown menu for transaction (investment/withdrawal) actions.
 * Provides edit and delete operations.
 */

import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { toast } from 'react-hot-toast'
import { MoreVertical, Pencil, Trash2 } from 'lucide-react'

import { EditTransactionSheet } from './sheets/EditTransactionSheet'
import type { Investment, Withdrawal } from '@/api/accounting'
import {
  useDeleteInvestmentMutation,
  useDeleteWithdrawalMutation,
} from '@/api/accounting'
import { Button } from '@/components/atoms/Button'
import { ConfirmDialog } from '@/components/molecules/ConfirmDialog'
import { translateErrorAsync } from '@/lib/translateError'

interface TransactionQuickActionsProps {
  transaction: Investment | Withdrawal
  type: 'investment' | 'withdrawal'
  businessDescriptor: string
  onActionComplete?: () => void
}

export function TransactionQuickActions({
  transaction,
  type,
  businessDescriptor,
  onActionComplete,
}: TransactionQuickActionsProps) {
  const { t } = useTranslation('accounting')
  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const [isEditOpen, setIsEditOpen] = useState(false)
  const [isDeleteOpen, setIsDeleteOpen] = useState(false)

  const deleteInvestmentMutation =
    useDeleteInvestmentMutation(businessDescriptor)
  const deleteWithdrawalMutation =
    useDeleteWithdrawalMutation(businessDescriptor)

  const handleDelete = async () => {
    try {
      if (type === 'investment') {
        await deleteInvestmentMutation.mutateAsync(transaction.id)
        toast.success(t('toast.investment_deleted'))
      } else {
        await deleteWithdrawalMutation.mutateAsync(transaction.id)
        toast.success(t('toast.withdrawal_deleted'))
      }
      await onActionComplete?.()
      setIsDeleteOpen(false)
    } catch (error) {
      const message = await translateErrorAsync(error, t)
      toast.error(message)
    }
  }

  return (
    <>
      {/* Dropdown Menu */}
      <div className="dropdown dropdown-end">
        <Button
          type="button"
          variant="ghost"
          size="sm"
          className="btn-square"
          aria-label="More actions"
          onClick={() => setIsMenuOpen(!isMenuOpen)}
        >
          <MoreVertical className="h-4 w-4" />
        </Button>

        {isMenuOpen && (
          <ul className="dropdown-content z-10 menu p-2 shadow bg-base-100 border border-base-300 rounded-box w-52 mt-2">
            <li>
              <button
                type="button"
                onClick={() => {
                  setIsMenuOpen(false)
                  setIsEditOpen(true)
                }}
                className="flex items-center gap-2"
              >
                <Pencil className="h-4 w-4" />
                {t(`edit.${type}_title`)}
              </button>
            </li>
            <li>
              <button
                type="button"
                onClick={() => {
                  setIsMenuOpen(false)
                  setIsDeleteOpen(true)
                }}
                className="flex items-center gap-2 text-error"
              >
                <Trash2 className="h-4 w-4" />
                {t(`delete.${type}_title`)}
              </button>
            </li>
          </ul>
        )}
      </div>

      {/* Edit Sheet */}
      <EditTransactionSheet
        isOpen={isEditOpen}
        onClose={() => setIsEditOpen(false)}
        transaction={transaction}
        type={type}
        businessDescriptor={businessDescriptor}
        onUpdated={onActionComplete}
      />

      {/* Delete Confirmation */}
      <ConfirmDialog
        isOpen={isDeleteOpen}
        onClose={() => setIsDeleteOpen(false)}
        onConfirm={handleDelete}
        title={t(`delete.${type}_title`)}
        message={t(`delete.${type}_message`)}
        confirmText={t('common:actions.delete')}
        cancelText={t('common:actions.cancel')}
        variant="error"
      />
    </>
  )
}
