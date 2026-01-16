/**
 * CreateTransactionSheet Component
 *
 * Bottom sheet form for recording capital transactions.
 *
 * Supports two transaction types:
 * 1. Investment (Owner capital in)
 * 2. Withdrawal (Owner capital out)
 *
 * Features:
 * - Type switcher (Investment/Withdrawal)
 * - Smart helper showing "Safe to Draw" amount for withdrawals
 * - Amount validation
 * - Optional note field
 * - Date picker
 */

import { useTranslation } from 'react-i18next'
import { toast } from 'react-hot-toast'
import { format } from 'date-fns'
import { ArrowDownCircle, ArrowUpCircle, Info } from 'lucide-react'

import {
  useAccountingSummaryQuery,
  useCreateInvestmentMutation,
  useCreateWithdrawalMutation,
} from '@/api/accounting'
import { BottomSheet } from '@/components/molecules/BottomSheet'
import { Button } from '@/components/atoms/Button'
import { useKyoraForm } from '@/lib/form'
import { translateErrorAsync } from '@/lib/translateError'
import { getSelectedBusiness } from '@/stores/businessStore'
import { formatCurrency } from '@/lib/formatCurrency'

interface CreateTransactionSheetProps {
  isOpen: boolean
  onClose: () => void
  businessDescriptor: string
  defaultType?: 'investment' | 'withdrawal'
  onCreated?: () => void | Promise<void>
}

export function CreateTransactionSheet({
  isOpen,
  onClose,
  businessDescriptor,
  defaultType = 'investment',
  onCreated,
}: CreateTransactionSheetProps) {
  const { t } = useTranslation('accounting')
  const { t: tCommon } = useTranslation('common')

  const business = getSelectedBusiness()
  const currency = business?.currency ?? 'USD'

  // Fetch safe to draw amount for withdrawal helper
  const { data: summary } = useAccountingSummaryQuery(businessDescriptor)
  const safeToDrawAmount = parseFloat(summary?.safeToDrawAmount ?? '0')

  const createInvestmentMutation =
    useCreateInvestmentMutation(businessDescriptor)
  const createWithdrawalMutation =
    useCreateWithdrawalMutation(businessDescriptor)

  const form = useKyoraForm({
    defaultValues: {
      type: defaultType,
      amount: '',
      date: new Date(),
      note: '',
    },
    onSubmit: async ({ value }) => {
      const amount = parseFloat(value.amount)
      if (isNaN(amount) || amount <= 0) {
        toast.error(t('validation.amount_positive'))
        return
      }

      // Warn if withdrawal exceeds safe amount (soft block)
      if (value.type === 'withdrawal' && amount > safeToDrawAmount) {
        const proceed = window.confirm(
          `${t('helper.exceeds_safe_amount')}\n\n${t('helper.safe_amount', { amount: formatCurrency(safeToDrawAmount, currency) })}\n\n${tCommon('actions.continue')}?`,
        )
        if (!proceed) return
      }

      try {
        const dateString = format(value.date, 'yyyy-MM-dd')

        if (value.type === 'investment') {
          await createInvestmentMutation.mutateAsync({
            amount: value.amount,
            investedAt: dateString,
            note: value.note || undefined,
          })
          toast.success(t('toast.investment_created'))
        } else {
          await createWithdrawalMutation.mutateAsync({
            amount: value.amount,
            withdrawnAt: dateString,
            note: value.note || undefined,
          })
          toast.success(t('toast.withdrawal_created'))
        }

        // Reset and close
        form.reset()
        await onCreated?.()
        onClose()
      } catch (error) {
        const message = await translateErrorAsync(error, t)
        toast.error(message)
      }
    },
  })

  const isSubmitting =
    createInvestmentMutation.isPending || createWithdrawalMutation.isPending

  const handleClose = () => {
    if (isSubmitting) return
    form.reset()
    onClose()
  }

  return (
    <form.AppForm>
      <BottomSheet
        isOpen={isOpen}
        onClose={handleClose}
        title={t('actions.record_transaction')}
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
            <form.Subscribe selector={(state) => state.values.type}>
              {(transactionType) => (
                <form.SubmitButton
                  form="create-transaction-form"
                  variant="primary"
                  className="flex-1"
                  disabled={isSubmitting}
                >
                  {isSubmitting
                    ? tCommon('saving')
                    : transactionType === 'investment'
                      ? t('actions.invest')
                      : t('actions.withdraw')}
                </form.SubmitButton>
              )}
            </form.Subscribe>
          </div>
        }
      >
        <form.FormRoot id="create-transaction-form" className="space-y-4">
          {/* Transaction Type Switcher */}
          <form.Subscribe selector={(state) => state.values.type}>
            {(transactionType) => (
              <div className="flex gap-2 p-1 bg-base-200 rounded-lg">
                <button
                  type="button"
                  onClick={() => form.setFieldValue('type', 'investment')}
                  className={`flex-1 flex items-center justify-center gap-2 px-4 py-2 rounded-md transition-colors ${
                    transactionType === 'investment'
                      ? 'bg-base-100 shadow text-success font-medium'
                      : 'text-base-content/60 hover:text-base-content'
                  }`}
                >
                  <ArrowDownCircle className="h-4 w-4" />
                  {t('type.investment')}
                </button>
                <button
                  type="button"
                  onClick={() => form.setFieldValue('type', 'withdrawal')}
                  className={`flex-1 flex items-center justify-center gap-2 px-4 py-2 rounded-md transition-colors ${
                    transactionType === 'withdrawal'
                      ? 'bg-base-100 shadow text-error font-medium'
                      : 'text-base-content/60 hover:text-base-content'
                  }`}
                >
                  <ArrowUpCircle className="h-4 w-4" />
                  {t('type.withdrawal')}
                </button>
              </div>
            )}
          </form.Subscribe>

          {/* Amount Field */}
          <form.AppField name="amount">
            {(field) => (
              <field.PriceField
                label={t('form.amount')}
                placeholder={t('form.amount_placeholder')}
                currencyCode={currency}
                required
              />
            )}
          </form.AppField>

          {/* Safe to Draw Helper (Withdrawal only) */}
          <form.Subscribe selector={(state) => state.values.type}>
            {(transactionType) =>
              transactionType === 'withdrawal' && safeToDrawAmount > 0 ? (
                <div className="alert alert-info py-2 px-3">
                  <Info className="h-4 w-4 shrink-0" />
                  <span className="text-sm">
                    {t('helper.safe_amount', {
                      amount: formatCurrency(safeToDrawAmount, currency),
                    })}
                  </span>
                </div>
              ) : null
            }
          </form.Subscribe>

          {/* Date Field */}
          <form.AppField name="date">
            {(field) => (
              <field.DateField
                label={t('form.date')}
                maxDate={new Date()}
                required
              />
            )}
          </form.AppField>

          {/* Note Field */}
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
