/**
 * CreateExpenseSheet Component
 *
 * Bottom sheet form for creating a new expense.
 *
 * Supports two modes:
 * 1. One-time expense (default)
 * 2. Recurring expense (when toggle is enabled)
 *
 * For recurring expenses, additional fields appear:
 * - Frequency (daily, weekly, monthly, yearly)
 * - Start date
 * - End date (optional)
 * - Auto-backfill historical expenses toggle
 */

import { useTranslation } from 'react-i18next'
import { toast } from 'react-hot-toast'
import { format } from 'date-fns'

import { CATEGORY_OPTIONS, FREQUENCY_OPTIONS } from '../../schema/options'
import type { RecurringExpenseFrequency } from '@/api/types/accounting'
import type { ExpenseCategory } from '@/api/accounting'
import {
  useCreateExpenseMutation,
  useCreateRecurringExpenseMutation,
} from '@/api/accounting'
import { BottomSheet } from '@/components/molecules/BottomSheet'
import { Button } from '@/components/atoms/Button'
import { useKyoraForm } from '@/lib/form'
import { getSelectedBusiness } from '@/stores/businessStore'

interface CreateExpenseSheetProps {
  isOpen: boolean
  onClose: () => void
  businessDescriptor: string
  onCreated?: () => void | Promise<void>
}

export function CreateExpenseSheet({
  isOpen,
  onClose,
  businessDescriptor,
  onCreated,
}: CreateExpenseSheetProps) {
  const { t } = useTranslation('accounting')
  const { t: tCommon } = useTranslation('common')

  const business = getSelectedBusiness()
  const currency = business?.currency ?? 'USD'

  const createExpenseMutation = useCreateExpenseMutation(businessDescriptor)
  const createRecurringExpenseMutation =
    useCreateRecurringExpenseMutation(businessDescriptor)

  const form = useKyoraForm({
    defaultValues: {
      amount: '',
      category: '' as ExpenseCategory | '',
      date: new Date(),
      note: '',
      isRecurring: false,
      frequency: 'monthly' as RecurringExpenseFrequency,
      endDate: null as Date | null,
      autoBackfill: false,
    },
    onSubmit: async ({ value }) => {
      // Validate required fields
      if (!value.category) {
        toast.error(t('validation.category_required'))
        return
      }

      const amount = parseFloat(value.amount)
      if (isNaN(amount) || amount <= 0) {
        toast.error(t('validation.amount_positive'))
        return
      }

      // Format date to YYYY-MM-DD
      const dateString = format(value.date, 'yyyy-MM-dd')

      if (value.isRecurring) {
        // Create recurring expense
        await createRecurringExpenseMutation.mutateAsync({
          amount: value.amount,
          category: value.category,
          note: value.note || undefined,
          frequency: value.frequency,
          recurringStartDate: dateString,
          recurringEndDate: value.endDate
            ? format(value.endDate, 'yyyy-MM-dd')
            : undefined,
          autoCreateHistoricalExpenses: value.autoBackfill,
        })
        toast.success(t('toast.recurring_expense_created'))
      } else {
        // Create one-time expense (backend requires type: 'one_time')
        await createExpenseMutation.mutateAsync({
          amount: value.amount,
          category: value.category,
          type: 'one_time',
          occurredOn: dateString,
          note: value.note || undefined,
        })
        toast.success(t('toast.expense_created'))
      }

      // Reset form and close
      form.reset()
      await onCreated?.()
      onClose()
    },
  })

  // Build options for selects
  const categoryOptions = [
    { value: '', label: t('form.select_category') },
    ...CATEGORY_OPTIONS.map((opt) => ({
      value: opt.value,
      label: t(opt.labelKey),
    })),
  ]

  const frequencyOptions = FREQUENCY_OPTIONS.map((opt) => ({
    value: opt.value,
    label: t(opt.labelKey),
  }))

  const isSubmitting =
    createExpenseMutation.isPending || createRecurringExpenseMutation.isPending

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
        title={t('form.add_expense_title')}
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
            <form.Subscribe selector={(state) => state.values.isRecurring}>
              {(isRecurring) => (
                <form.SubmitButton
                  form="create-expense-form"
                  variant="primary"
                  className="flex-1"
                  disabled={isSubmitting}
                >
                  {isSubmitting
                    ? tCommon('saving')
                    : isRecurring
                      ? t('actions.create_recurring')
                      : t('actions.create_expense')}
                </form.SubmitButton>
              )}
            </form.Subscribe>
          </div>
        }
      >
        <form.FormRoot id="create-expense-form" className="space-y-4">
          {/* Amount Field - Use PriceField for money input */}
          <form.AppField
            name="amount"
            validators={{
              onChange: ({ value }) => {
                if (!value) return t('validation.amount_required')
                const num = parseFloat(value)
                if (isNaN(num) || num <= 0)
                  return t('validation.amount_positive')
                return undefined
              },
            }}
          >
            {(field) => (
              <field.PriceField
                label={t('form.amount')}
                currencyCode={currency}
                placeholder="0.00"
                required
              />
            )}
          </form.AppField>

          {/* Category Field */}
          <form.AppField
            name="category"
            validators={{
              onChange: ({ value }) => {
                if (!value) return t('validation.category_required')
                return undefined
              },
            }}
          >
            {(field) => (
              <field.SelectField
                label={t('category.label')}
                options={categoryOptions}
                required
              />
            )}
          </form.AppField>

          {/* Date Field - changes label based on isRecurring */}
          <form.Subscribe selector={(state) => state.values.isRecurring}>
            {(isRecurring) => (
              <form.AppField name="date">
                {(field) => (
                  <field.DateField
                    label={isRecurring ? t('form.start_date') : t('form.date')}
                    required
                  />
                )}
              </form.AppField>
            )}
          </form.Subscribe>

          {/* Note Field */}
          <form.AppField name="note">
            {(field) => (
              <field.TextareaField
                label={t('form.note')}
                placeholder={t('form.note_placeholder')}
                rows={2}
              />
            )}
          </form.AppField>

          {/* Recurring Toggle */}
          <form.AppField name="isRecurring">
            {(field) => <field.ToggleField label={t('form.is_recurring')} />}
          </form.AppField>

          {/* Recurring-specific fields */}
          <form.Subscribe selector={(state) => state.values.isRecurring}>
            {(isRecurring) =>
              isRecurring ? (
                <div className="space-y-4 border-t border-base-200 pt-4">
                  {/* Frequency */}
                  <form.AppField name="frequency">
                    {(field) => (
                      <field.SelectField
                        label={t('form.frequency')}
                        options={frequencyOptions}
                        required
                      />
                    )}
                  </form.AppField>

                  {/* End Date (optional) */}
                  <form.AppField name="endDate">
                    {(field) => <field.DateField label={t('form.end_date')} />}
                  </form.AppField>

                  {/* Auto-backfill toggle */}
                  <form.AppField name="autoBackfill">
                    {(field) => (
                      <field.ToggleField
                        label={t('form.auto_backfill')}
                        hint={t('form.auto_backfill_hint')}
                      />
                    )}
                  </form.AppField>
                </div>
              ) : null
            }
          </form.Subscribe>
        </form.FormRoot>
      </BottomSheet>
    </form.AppForm>
  )
}
