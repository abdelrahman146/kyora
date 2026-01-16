/**
 * EditRecurringExpenseSheet Component
 *
 * Bottom sheet form for editing a recurring expense template.
 *
 * Editable fields:
 * - Amount
 * - Category
 * - Frequency
 * - End date (optional)
 * - Note
 *
 * Note: Start date is not editable after creation.
 * Status changes are handled via a separate action.
 */

import { useTranslation } from 'react-i18next'
import { toast } from 'react-hot-toast'
import { useEffect } from 'react'
import { format, parseISO } from 'date-fns'
import { useQueryClient } from '@tanstack/react-query'

import { CATEGORY_OPTIONS, FREQUENCY_OPTIONS } from '../../schema/options'
import type { ExpenseCategory, RecurringExpense } from '@/api/accounting'
import {
  accountingQueries,
  useRecurringExpenseQuery,
  useUpdateRecurringExpenseMutation,
} from '@/api/accounting'
import { BottomSheet } from '@/components/molecules/BottomSheet'
import { Button } from '@/components/atoms/Button'
import { useKyoraForm } from '@/lib/form'
import { getSelectedBusiness } from '@/stores/businessStore'
import { FormInput } from '@/components/form/FormInput'

export interface EditRecurringExpenseSheetProps {
  isOpen: boolean
  onClose: () => void
  /** Either pass the full recurring expense or just the ID */
  recurringExpense?: RecurringExpense | null
  recurringExpenseId?: string | null
  businessDescriptor: string
  onUpdated?: () => void | Promise<void>
}

export function EditRecurringExpenseSheet({
  isOpen,
  onClose,
  recurringExpense: providedRecurringExpense,
  recurringExpenseId,
  businessDescriptor,
  onUpdated,
}: EditRecurringExpenseSheetProps) {
  const { t } = useTranslation('accounting')
  const { t: tCommon } = useTranslation('common')
  const queryClient = useQueryClient()

  const business = getSelectedBusiness()
  const currency = business?.currency ?? 'USD'

  // Determine which ID to use for fetching
  const effectiveId = providedRecurringExpense?.id ?? recurringExpenseId

  // Fetch recurring expense if only ID was provided
  const { data: fetchedRecurringExpense, isLoading: isFetching } =
    useRecurringExpenseQuery(businessDescriptor, effectiveId ?? '', {
      enabled: isOpen && !providedRecurringExpense && !!effectiveId,
    })

  // Use provided data or fetched data
  const recurringExpense = providedRecurringExpense ?? fetchedRecurringExpense

  const updateMutation = useUpdateRecurringExpenseMutation(
    businessDescriptor,
    recurringExpense?.id ?? '',
  )

  // Parse the start date to display (read-only)
  const getStartDate = (): string => {
    if (recurringExpense?.recurringStartDate) {
      try {
        return format(parseISO(recurringExpense.recurringStartDate), 'PP')
      } catch {
        return recurringExpense.recurringStartDate
      }
    }
    return ''
  }

  // Parse the end date to a Date object (editable)
  const getInitialEndDate = (): Date | null => {
    if (recurringExpense?.recurringEndDate) {
      try {
        return parseISO(recurringExpense.recurringEndDate)
      } catch {
        return null
      }
    }
    return null
  }

  const form = useKyoraForm({
    defaultValues: {
      amount: recurringExpense?.amount ?? '',
      category: recurringExpense?.category ?? '',
      frequency: recurringExpense?.frequency ?? 'monthly',
      endDate: getInitialEndDate(),
      note: recurringExpense?.note ?? '',
    },
    onSubmit: async ({ value }) => {
      if (!recurringExpense) return

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

      await updateMutation.mutateAsync({
        amount: value.amount,
        category: value.category as ExpenseCategory,
        frequency: value.frequency,
        recurringEndDate: value.endDate
          ? format(value.endDate, 'yyyy-MM-dd')
          : null,
        note: value.note || null,
      })

      toast.success(t('toast.recurring_updated'))

      // Invalidate queries
      await queryClient.invalidateQueries({
        queryKey: accountingQueries.recurringExpenses(),
      })
      await queryClient.invalidateQueries({
        queryKey: accountingQueries.expenses(),
      })

      // Notify parent and close
      await onUpdated?.()
      onClose()
    },
  })

  // Sync form values when recurring expense changes
  useEffect(() => {
    if (recurringExpense) {
      form.setFieldValue('amount', recurringExpense.amount)
      form.setFieldValue('category', recurringExpense.category)
      form.setFieldValue('frequency', recurringExpense.frequency)
      form.setFieldValue('endDate', getInitialEndDate())
      form.setFieldValue('note', recurringExpense.note ?? '')
    }
  }, [recurringExpense?.id])

  // Build options for selects
  const categoryOptions = CATEGORY_OPTIONS.map((opt) => ({
    value: opt.value,
    label: t(opt.labelKey),
  }))

  const frequencyOptions = FREQUENCY_OPTIONS.map((opt) => ({
    value: opt.value,
    label: t(opt.labelKey),
  }))

  const isSubmitting = updateMutation.isPending
  const isLoadingData = isFetching && !recurringExpense

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
        title={t('edit.recurring_template_title')}
        closeOnOverlayClick={!isSubmitting}
        closeOnEscape={!isSubmitting}
        footer={
          <div className="flex gap-2">
            <Button
              type="button"
              variant="ghost"
              className="flex-1"
              onClick={handleClose}
              disabled={isSubmitting || isLoadingData}
            >
              {tCommon('cancel')}
            </Button>
            <form.SubmitButton
              form="edit-recurring-expense-form"
              variant="primary"
              className="flex-1"
              disabled={isSubmitting || isLoadingData}
            >
              {isSubmitting ? tCommon('saving') : tCommon('save')}
            </form.SubmitButton>
          </div>
        }
      >
        {isLoadingData ? (
          <div className="space-y-4">
            <div className="skeleton h-12 w-full rounded-lg" />
            <div className="skeleton h-12 w-full rounded-lg" />
            <div className="skeleton h-12 w-full rounded-lg" />
            <div className="skeleton h-12 w-full rounded-lg" />
          </div>
        ) : (
          <form.FormRoot id="edit-recurring-expense-form" className="space-y-4">
            {/* Amount Field */}
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

            {/* Frequency Field */}
            <form.AppField name="frequency">
              {(field) => (
                <field.SelectField
                  label={t('frequency.label')}
                  options={frequencyOptions}
                  required
                />
              )}
            </form.AppField>

            {/* Start Date (read-only display) */}
            <div className="form-control">
              <label className="label pb-1">
                <span className="label-text font-medium">
                  {t('form.start_date')}
                </span>
              </label>
              <FormInput value={getStartDate()} readOnly />
              <label className="label pt-1">
                <span className="label-text-alt text-base-content/50">
                  {t('edit.start_date_not_editable')}
                </span>
              </label>
            </div>

            {/* End Date Field (editable) */}
            <form.AppField name="endDate">
              {(field) => <field.DateField label={t('form.end_date')} />}
            </form.AppField>

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

            {/* Status indicator */}
            {recurringExpense?.status && (
              <div className="flex items-center gap-2 p-3 bg-base-200/50 rounded-lg">
                <span className="text-sm text-base-content/70">
                  {t('status.label')}:
                </span>
                <span
                  className={`badge badge-sm ${
                    recurringExpense.status === 'active'
                      ? 'badge-success'
                      : recurringExpense.status === 'paused'
                        ? 'badge-warning'
                        : 'badge-ghost'
                  }`}
                >
                  {t(`status.${recurringExpense.status}`)}
                </span>
              </div>
            )}
          </form.FormRoot>
        )}
      </BottomSheet>
    </form.AppForm>
  )
}
