/**
 * RecurringExpenseQuickActions Component
 *
 * Dropdown actions for recurring expense templates.
 * Supports:
 * - Status changes (pause, resume, end, cancel)
 * - Edit recurring template
 * - Delete recurring template
 */

import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useQueryClient } from '@tanstack/react-query'
import { toast } from 'react-hot-toast'
import {
  Ban,
  CheckCircle,
  Edit2,
  MoreVertical,
  Pause,
  Play,
  Trash2,
} from 'lucide-react'

import { EditRecurringExpenseSheet } from './sheets/EditRecurringExpenseSheet'
import type { RecurringExpense } from '@/api/accounting'
import type { RecurringExpenseStatus } from '@/api/types/accounting'
import {
  accountingQueries,
  useDeleteRecurringExpenseMutation,
  useUpdateRecurringExpenseStatusMutation,
} from '@/api/accounting'
import { ConfirmDialog } from '@/components/molecules/ConfirmDialog'
import { translateErrorAsync } from '@/lib/translateError'

interface RecurringExpenseQuickActionsProps {
  recurringExpense: RecurringExpense
  businessDescriptor: string
  onActionComplete?: () => void
}

export function RecurringExpenseQuickActions({
  recurringExpense,
  businessDescriptor,
  onActionComplete,
}: RecurringExpenseQuickActionsProps) {
  const { t } = useTranslation('accounting')
  const { t: tCommon } = useTranslation('common')
  const queryClient = useQueryClient()

  const [showEditSheet, setShowEditSheet] = useState(false)
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false)
  const [statusToUpdate, setStatusToUpdate] =
    useState<RecurringExpenseStatus | null>(null)

  const deleteMutation = useDeleteRecurringExpenseMutation(
    businessDescriptor,
    recurringExpense.id,
  )

  const updateStatusMutation = useUpdateRecurringExpenseStatusMutation(
    businessDescriptor,
    recurringExpense.id,
  )

  const invalidateQueries = async () => {
    await queryClient.invalidateQueries({
      queryKey: accountingQueries.recurringExpenses(),
    })
    await queryClient.invalidateQueries({
      queryKey: accountingQueries.expenses(),
    })
    await queryClient.invalidateQueries({
      queryKey: accountingQueries.summaries(),
    })
    await queryClient.invalidateQueries({
      queryKey: accountingQueries.recentActivitiesKey(),
    })
    onActionComplete?.()
  }

  const handleEditClick = (e: React.MouseEvent) => {
    e.stopPropagation()
    setShowEditSheet(true)
  }

  const handleDeleteClick = (e: React.MouseEvent) => {
    e.stopPropagation()
    setShowDeleteConfirm(true)
  }

  const handleStatusClick = (
    e: React.MouseEvent,
    status: RecurringExpenseStatus,
  ) => {
    e.stopPropagation()
    setStatusToUpdate(status)
  }

  const handleConfirmStatusChange = async () => {
    if (!statusToUpdate) return

    try {
      await updateStatusMutation.mutateAsync({ status: statusToUpdate })
      toast.success(t('toast.recurring_status_updated'))
      await invalidateQueries()
    } catch (error) {
      const message = await translateErrorAsync(error, t)
      toast.error(message)
    }
    setStatusToUpdate(null)
  }

  const handleConfirmDelete = async () => {
    try {
      await deleteMutation.mutateAsync()
      toast.success(t('toast.recurring_deleted'))
      await invalidateQueries()
    } catch (error) {
      const message = await translateErrorAsync(error, t)
      toast.error(message)
    }
    setShowDeleteConfirm(false)
  }

  const handleActionComplete = async () => {
    await invalidateQueries()
  }

  const status = recurringExpense.status

  // Determine available status transitions
  // Based on backend state machine:
  // active → paused|ended|canceled
  // paused → active|ended|canceled
  // ended → active|canceled
  // canceled → active
  const statusActions: Array<{
    status: RecurringExpenseStatus
    icon: typeof Play
    labelKey: string
  }> = []

  if (status === 'active') {
    statusActions.push(
      { status: 'paused', icon: Pause, labelKey: 'recurring.pause' },
      { status: 'ended', icon: CheckCircle, labelKey: 'recurring.end' },
      { status: 'canceled', icon: Ban, labelKey: 'recurring.cancel' },
    )
  } else if (status === 'paused') {
    statusActions.push(
      { status: 'active', icon: Play, labelKey: 'recurring.resume' },
      { status: 'ended', icon: CheckCircle, labelKey: 'recurring.end' },
      { status: 'canceled', icon: Ban, labelKey: 'recurring.cancel' },
    )
  } else if (status === 'ended') {
    statusActions.push(
      {
        status: 'active',
        icon: Play,
        labelKey: 'recurring.reactivate',
      },
      { status: 'canceled', icon: Ban, labelKey: 'recurring.cancel' },
    )
  } else {
    // status === 'canceled'
    statusActions.push({
      status: 'active',
      icon: Play,
      labelKey: 'recurring.reactivate',
    })
  }

  return (
    <>
      <div className="dropdown dropdown-end ">
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
          {/* Status Actions */}
          {statusActions.length > 0 && (
            <>
              <li className="menu-title text-xs px-2 py-1">
                {t('recurring.change_status')}
              </li>
              {statusActions.map(
                ({ status: newStatus, icon: Icon, labelKey }) => (
                  <li key={newStatus}>
                    <button
                      type="button"
                      onClick={(e) => handleStatusClick(e, newStatus)}
                    >
                      <Icon className="w-4 h-4" />
                      {t(labelKey)}
                    </button>
                  </li>
                ),
              )}
              <div className="divider my-1" />
            </>
          )}

          {/* Edit */}
          <li>
            <button type="button" onClick={handleEditClick}>
              <Edit2 className="w-4 h-4" />
              {t('edit.recurring_template_title')}
            </button>
          </li>

          <div className="divider my-1" />

          {/* Delete */}
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
        </ul>
      </div>

      {/* Edit Sheet */}
      <EditRecurringExpenseSheet
        isOpen={showEditSheet}
        onClose={() => setShowEditSheet(false)}
        recurringExpense={recurringExpense}
        businessDescriptor={businessDescriptor}
        onUpdated={handleActionComplete}
      />

      {/* Status Change Confirmation */}
      <ConfirmDialog
        isOpen={statusToUpdate !== null}
        onClose={() => setStatusToUpdate(null)}
        onConfirm={handleConfirmStatusChange}
        title={t('recurring.confirm_status_title')}
        message={t('recurring.confirm_status_message', {
          status: statusToUpdate ? t(`status.${statusToUpdate}`) : '',
        })}
        confirmText={tCommon('confirm')}
        cancelText={tCommon('cancel')}
        variant="warning"
      />

      {/* Delete Confirmation */}
      <ConfirmDialog
        isOpen={showDeleteConfirm}
        onClose={() => setShowDeleteConfirm(false)}
        onConfirm={handleConfirmDelete}
        title={t('delete.recurring_title')}
        message={t('delete.recurring_warning')}
        confirmText={tCommon('delete')}
        cancelText={tCommon('cancel')}
        variant="error"
      />
    </>
  )
}
