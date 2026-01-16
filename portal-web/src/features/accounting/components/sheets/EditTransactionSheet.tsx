/**
 * EditTransactionSheet Component
 *
 * Bottom sheet form for editing capital transactions (investments/withdrawals).
 */

import { useTranslation } from 'react-i18next'
import { toast } from 'react-hot-toast'
import { parseISO } from 'date-fns'
import { useStore } from '@tanstack/react-store'
import { z } from 'zod'

import type { Investment, Withdrawal } from '@/api/accounting'
import {
  useUpdateInvestmentMutation,
  useUpdateWithdrawalMutation,
} from '@/api/accounting'
import { BottomSheet } from '@/components/molecules/BottomSheet'
import { Button } from '@/components/atoms/Button'
import { useKyoraForm } from '@/lib/form'
import { getSelectedBusiness } from '@/stores/businessStore'
import { authStore } from '@/stores/authStore'

interface EditTransactionSheetProps {
  isOpen: boolean
  onClose: () => void
  transaction: Investment | Withdrawal
  type: 'investment' | 'withdrawal'
  businessDescriptor: string
  onUpdated?: () => void | Promise<void>
}

export function EditTransactionSheet({
  isOpen,
  onClose,
  transaction,
  type,
  businessDescriptor,
  onUpdated,
}: EditTransactionSheetProps) {
  const { t } = useTranslation('accounting')
  const { t: tCommon } = useTranslation('common')
  const business = getSelectedBusiness()
  const currency = business?.currency ?? 'USD'
  const user = useStore(authStore, (state) => state.user)

  const updateInvestmentMutation =
    useUpdateInvestmentMutation(businessDescriptor)
  const updateWithdrawalMutation =
    useUpdateWithdrawalMutation(businessDescriptor)

  const dateField =
    type === 'investment'
      ? (transaction as Investment).investedAt
      : (transaction as Withdrawal).withdrawnAt

  const form = useKyoraForm({
    defaultValues: {
      amount: transaction.amount,
      date: parseISO(dateField),
      note: transaction.note ?? '',
    },
    onSubmit: async ({ value }) => {
      if (!user) {
        toast.error(tCommon('errors.authentication_required'))
        return
      }

      // Let global error handler catch errors - don't manually toast them
      if (type === 'investment') {
        await updateInvestmentMutation.mutateAsync({
          id: transaction.id,
          data: {
            investorId: user.id,
            amount: value.amount,
            investedAt: value.date.toISOString(),
            note: value.note || undefined,
          },
        })
        toast.success(t('toast.investment_updated'))
      } else {
        await updateWithdrawalMutation.mutateAsync({
          id: transaction.id,
          data: {
            withdrawerId: user.id,
            amount: value.amount,
            withdrawnAt: value.date.toISOString(),
            note: value.note || undefined,
          },
        })
        toast.success(t('toast.withdrawal_updated'))
      }

      await onUpdated?.()
      onClose()
    },
  })

  const isSubmitting =
    updateInvestmentMutation.isPending || updateWithdrawalMutation.isPending

  const handleClose = () => {
    if (isSubmitting) return
    onClose()
  }

  return (
    <form.AppForm>
      <BottomSheet
        isOpen={isOpen}
        onClose={handleClose}
        title={t(`edit.${type}_title`)}
        closeOnOverlayClick={!isSubmitting}
        closeOnEscape={!isSubmitting}
        footer={
          <div className="flex gap-2">
            <Button
              type="button"
              variant="ghost"
              className="flex-1"
              onClick={handleClose}
              disabled={isSubmitting}
            >
              {tCommon('cancel')}
            </Button>
            <form.SubmitButton
              form="edit-transaction-form"
              variant="primary"
              className="flex-1"
              disabled={isSubmitting}
            >
              {isSubmitting ? tCommon('saving') : tCommon('actions.save')}
            </form.SubmitButton>
          </div>
        }
      >
        <form.FormRoot id="edit-transaction-form" className="space-y-4">
          {/* Amount */}
          <form.AppField
            name="amount"
            validators={{
              onChange: z
                .string()
                .min(1, 'accounting:validation.amount_required')
                .refine(
                  (val) => {
                    const num = parseFloat(val)
                    return !isNaN(num) && num > 0
                  },
                  { message: 'accounting:validation.amount_positive' },
                ),
            }}
          >
            {(field) => (
              <field.PriceField
                label={t('form.amount')}
                placeholder={t('form.amount_placeholder')}
                currencyCode={currency}
                required
              />
            )}
          </form.AppField>

          {/* Date */}
          <form.AppField
            name="date"
            validators={{
              onChange: z
                .date({
                  required_error: 'accounting:validation.date_required',
                  invalid_type_error: 'accounting:validation.date_required',
                })
                .max(new Date(), 'accounting:validation.date_not_future'),
            }}
          >
            {(field) => (
              <field.DateField
                label={t('form.date')}
                maxDate={new Date()}
                required
              />
            )}
          </form.AppField>

          {/* Note */}
          <form.AppField name="note">
            {(field) => (
              <field.TextareaField
                label={t('form.note')}
                placeholder={t('form.note_placeholder')}
              />
            )}
          </form.AppField>
        </form.FormRoot>
      </BottomSheet>
    </form.AppForm>
  )
}
