/**
 * ExpenseQuickActions Component
 *
 * Dropdown actions for expense items.
 * Supports both regular and recurring expenses.
 *
 * For recurring expense occurrences:
 * - Edit/Delete this occurrence
 * - Edit/Delete the recurring template
 */

import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useQueryClient } from '@tanstack/react-query'
import { toast } from 'react-hot-toast'
import { Edit2, MoreVertical, Repeat, Trash2 } from 'lucide-react'

import { EditExpenseSheet } from './sheets/EditExpenseSheet'
import type { Expense } from '@/api/accounting'
import {
  accountingQueries,
  useDeleteExpenseMutation,
  useDeleteRecurringExpenseMutation,
} from '@/api/accounting'
import { ConfirmDialog } from '@/components/molecules/ConfirmDialog'
import { translateErrorAsync } from '@/lib/translateError'
import { formatCurrency } from '@/lib/formatCurrency'

interface ExpenseQuickActionsProps {
  expense: Expense
  businessDescriptor: string
  currency: string
  onActionComplete?: () => void
}

type DeleteTarget = 'occurrence' | 'template' | null

export function ExpenseQuickActions({
  expense,
  businessDescriptor,
  currency,
  onActionComplete,
}: ExpenseQuickActionsProps) {
  const { t } = useTranslation('accounting')
  const { t: tCommon } = useTranslation('common')
  const queryClient = useQueryClient()

  const [showEditSheet, setShowEditSheet] = useState(false)
  const [deleteTarget, setDeleteTarget] = useState<DeleteTarget>(null)
  const [showRecurringDeleteSelect, setShowRecurringDeleteSelect] =
    useState(false)

  const isRecurring = expense.recurringExpenseId !== null

  const deleteExpenseMutation = useDeleteExpenseMutation(
    businessDescriptor,
    expense.id,
  )

  const deleteRecurringMutation = useDeleteRecurringExpenseMutation(
    businessDescriptor,
    expense.recurringExpenseId ?? '',
  )

  const handleEditClick = (e: React.MouseEvent) => {
    e.stopPropagation()
    setShowEditSheet(true)
  }

  const handleDeleteClick = (e: React.MouseEvent) => {
    e.stopPropagation()
    if (isRecurring) {
      // Show selection for recurring expenses
      setShowRecurringDeleteSelect(true)
    } else {
      // Direct delete for one-time expenses
      setDeleteTarget('occurrence')
    }
  }

  const handleDeleteOccurrence = (e: React.MouseEvent) => {
    e.stopPropagation()
    setShowRecurringDeleteSelect(false)
    setDeleteTarget('occurrence')
  }

  const handleDeleteTemplate = (e: React.MouseEvent) => {
    e.stopPropagation()
    setShowRecurringDeleteSelect(false)
    setDeleteTarget('template')
  }

  const handleConfirmDelete = async () => {
    try {
      if (deleteTarget === 'occurrence') {
        await deleteExpenseMutation.mutateAsync()
        toast.success(t('toast.expense_deleted'))
      } else if (deleteTarget === 'template') {
        await deleteRecurringMutation.mutateAsync()
        toast.success(t('toast.recurring_deleted'))
      }
      await queryClient.invalidateQueries({
        queryKey: accountingQueries.expenses(),
      })
      await queryClient.invalidateQueries({
        queryKey: accountingQueries.recurringExpenses(),
      })
      onActionComplete?.()
    } catch (error) {
      const message = await translateErrorAsync(error, t)
      toast.error(message)
    }
    setDeleteTarget(null)
  }

  const handleActionComplete = async () => {
    await queryClient.invalidateQueries({
      queryKey: accountingQueries.expenses(),
    })
    onActionComplete?.()
  }

  const amount = parseFloat(expense.amount)

  return (
    <>
      <div className="dropdown dropdown-end">
        <button
          type="button"
          tabIndex={0}
          role="button"
          className="btn btn-ghost btn-sm btn-square"
          aria-label={tCommon('actionsLabel')}
          onClick={(e) => e.stopPropagation()}
        >
          <MoreVertical className="w-4 h-4" />
        </button>
        <ul
          tabIndex={0}
          className="dropdown-content menu bg-base-100 rounded-box z-[100] w-56 p-2 border border-base-300 mt-2"
        >
          {/* Edit Actions */}
          <li>
            <button type="button" onClick={handleEditClick}>
              <Edit2 className="w-4 h-4" />
              {isRecurring ? t('edit.edit_expense') : tCommon('edit')}
            </button>
          </li>

          <div className="divider my-1" />

          {/* Delete Actions */}
          {isRecurring && !showRecurringDeleteSelect ? (
            <li>
              <button
                type="button"
                onClick={handleDeleteClick}
                className="text-error"
              >
                <Trash2 className="w-4 h-4" />
                {tCommon('delete')}
              </button>
            </li>
          ) : isRecurring && showRecurringDeleteSelect ? (
            <>
              <li className="menu-title text-xs px-2 py-1">
                {t('delete.choose_what_to_delete')}
              </li>
              <li>
                <button type="button" onClick={handleDeleteOccurrence}>
                  <Trash2 className="w-4 h-4 text-error" />
                  <div className="flex flex-col items-start">
                    <span className="text-error">
                      {t('delete.this_occurrence')}
                    </span>
                    <span className="text-xs text-base-content/60">
                      {t('delete.this_occurrence_hint')}
                    </span>
                  </div>
                </button>
              </li>
              <li>
                <button type="button" onClick={handleDeleteTemplate}>
                  <Repeat className="w-4 h-4 text-error" />
                  <div className="flex flex-col items-start">
                    <span className="text-error">
                      {t('delete.recurring_template')}
                    </span>
                    <span className="text-xs text-base-content/60">
                      {t('delete.recurring_template_hint')}
                    </span>
                  </div>
                </button>
              </li>
            </>
          ) : (
            <li>
              <button
                type="button"
                onClick={handleDeleteClick}
                className="text-error"
              >
                <Trash2 className="w-4 h-4" />
                {tCommon('delete')}
              </button>
            </li>
          )}
        </ul>
      </div>

      {/* Edit Sheet */}
      <EditExpenseSheet
        isOpen={showEditSheet}
        onClose={() => setShowEditSheet(false)}
        expense={expense}
        businessDescriptor={businessDescriptor}
        onUpdated={handleActionComplete}
      />

      {/* Delete Confirmation */}
      <ConfirmDialog
        isOpen={deleteTarget !== null}
        onClose={() => setDeleteTarget(null)}
        onConfirm={handleConfirmDelete}
        title={
          deleteTarget === 'template'
            ? t('delete.recurring_title')
            : t('delete.expense_title')
        }
        message={
          deleteTarget === 'template'
            ? t('delete.recurring_warning')
            : t('delete.expense_warning', {
                category: t(`category.${expense.category}`),
                amount: formatCurrency(amount, currency),
              })
        }
        confirmText={tCommon('delete')}
        cancelText={tCommon('cancel')}
        variant="error"
      />
    </>
  )
}
