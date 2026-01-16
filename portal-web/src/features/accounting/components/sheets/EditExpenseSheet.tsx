/**
 * EditExpenseSheet Component
 *
 * Bottom sheet form for editing an existing expense.
 *
 * For regular (one-time) expenses, allows editing:
 * - Amount
 * - Category
 * - Date
 * - Note
 *
 * For recurring expense occurrences, additional options are shown:
 * - Edit only this occurrence
 * - Edit the recurring template (opens EditRecurringExpenseSheet)
 */

import { useTranslation } from 'react-i18next'
import { toast } from 'react-hot-toast'
import { useEffect, useState } from 'react'
import { format, parseISO } from 'date-fns'
import { ChevronRight, FileText, Repeat } from 'lucide-react'

import { CATEGORY_OPTIONS } from '../../schema/options'
import { EditRecurringExpenseSheet } from './EditRecurringExpenseSheet'
import type { Expense, ExpenseCategory } from '@/api/accounting'
import { useUpdateExpenseMutation } from '@/api/accounting'
import { BottomSheet } from '@/components/molecules/BottomSheet'
import { Button } from '@/components/atoms/Button'
import { useKyoraForm } from '@/lib/form'
import { translateErrorAsync } from '@/lib/translateError'
import { getSelectedBusiness } from '@/stores/businessStore'

interface EditExpenseSheetProps {
  isOpen: boolean
  onClose: () => void
  expense: Expense | null
  businessDescriptor: string
  onUpdated?: () => void | Promise<void>
}

type EditMode = 'select' | 'occurrence'

/**
 * Mode selection component for recurring expense occurrences.
 * Allows user to choose between editing this occurrence or the recurring template.
 */
interface ModeSelectionProps {
  onEditOccurrence: () => void
  onEditTemplate: () => void
  t: (key: string) => string
}

function ModeSelection({
  onEditOccurrence,
  onEditTemplate,
  t,
}: ModeSelectionProps) {
  return (
    <div className="space-y-3">
      <p className="text-sm text-base-content/70 mb-4">
        {t('edit.recurring_choose_mode')}
      </p>

      {/* Edit This Occurrence */}
      <button
        type="button"
        className="w-full flex items-center gap-3 p-4 rounded-xl border border-base-300 hover:bg-base-200/50 transition-colors text-start cursor-pointer"
        onClick={onEditOccurrence}
      >
        <div className="w-10 h-10 rounded-lg bg-primary/10 text-primary flex items-center justify-center">
          <FileText className="w-5 h-5" />
        </div>
        <div className="flex-1">
          <p className="font-medium">{t('edit.edit_this_occurrence')}</p>
          <p className="text-sm text-base-content/60">
            {t('edit.edit_this_occurrence_hint')}
          </p>
        </div>
        <ChevronRight className="w-5 h-5 text-base-content/40" />
      </button>

      {/* Edit Recurring Template */}
      <button
        type="button"
        className="w-full flex items-center gap-3 p-4 rounded-xl border border-base-300 hover:bg-base-200/50 transition-colors text-start cursor-pointer"
        onClick={onEditTemplate}
      >
        <div className="w-10 h-10 rounded-lg bg-secondary/10 text-secondary flex items-center justify-center">
          <Repeat className="w-5 h-5" />
        </div>
        <div className="flex-1">
          <p className="font-medium">{t('edit.edit_recurring_template')}</p>
          <p className="text-sm text-base-content/60">
            {t('edit.edit_recurring_template_hint')}
          </p>
        </div>
        <ChevronRight className="w-5 h-5 text-base-content/40" />
      </button>
    </div>
  )
}

export function EditExpenseSheet({
  isOpen,
  onClose,
  expense,
  businessDescriptor,
  onUpdated,
}: EditExpenseSheetProps) {
  const { t } = useTranslation('accounting')
  const { t: tCommon } = useTranslation('common')

  const business = getSelectedBusiness()
  const currency = business?.currency ?? 'USD'

  // For recurring expenses, show mode selection first
  const [editMode, setEditMode] = useState<EditMode>('occurrence')
  // Track if we should show the recurring template sheet
  const [showRecurringSheet, setShowRecurringSheet] = useState(false)

  // Determine if this is a recurring expense occurrence
  const isRecurringOccurrence = expense?.recurringExpenseId != null

  // Reset mode when opening
  useEffect(() => {
    if (isOpen && expense) {
      // If recurring occurrence, show mode selection; otherwise go straight to edit
      setEditMode(isRecurringOccurrence ? 'select' : 'occurrence')
      setShowRecurringSheet(false)
    }
  }, [isOpen, expense, isRecurringOccurrence])

  const updateExpenseMutation = useUpdateExpenseMutation(
    businessDescriptor,
    expense?.id ?? '',
  )

  // Parse the expense date to a Date object
  const getInitialDate = (): Date => {
    if (expense?.occurredOn) {
      try {
        // occurredOn is YYYY-MM-DD
        return parseISO(expense.occurredOn)
      } catch {
        return new Date()
      }
    }
    return new Date()
  }

  const form = useKyoraForm({
    defaultValues: {
      amount: expense?.amount ?? '',
      category: expense?.category ?? '',
      date: getInitialDate(),
      note: expense?.note ?? '',
    },
    onSubmit: async ({ value }) => {
      if (!expense) return

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

      try {
        // Format date to YYYY-MM-DD
        const dateString = format(value.date, 'yyyy-MM-dd')

        await updateExpenseMutation.mutateAsync({
          amount: value.amount,
          category: value.category as ExpenseCategory,
          occurredOn: dateString,
          note: value.note || null,
        })

        toast.success(t('toast.expense_updated'))

        // Close and notify parent
        await onUpdated?.()
        onClose()
      } catch (error) {
        const message = await translateErrorAsync(error, t)
        toast.error(message)
      }
    },
  })

  // Sync form values when expense changes
  useEffect(() => {
    if (expense) {
      form.setFieldValue('amount', expense.amount)
      form.setFieldValue('category', expense.category)
      form.setFieldValue('date', getInitialDate())
      form.setFieldValue('note', expense.note ?? '')
    }
  }, [expense?.id])

  // Build options for selects
  const categoryOptions = CATEGORY_OPTIONS.map((opt) => ({
    value: opt.value,
    label: t(opt.labelKey),
  }))

  const isSubmitting = updateExpenseMutation.isPending

  const handleClose = () => {
    if (isSubmitting) return
    form.reset()
    setShowRecurringSheet(false)
    onClose()
  }

  const handleEditOccurrence = () => {
    setEditMode('occurrence')
  }

  const handleEditTemplate = () => {
    // Close this sheet and open the recurring expense sheet
    setShowRecurringSheet(true)
    onClose()
  }

  const handleRecurringSheetClose = () => {
    setShowRecurringSheet(false)
  }

  const handleRecurringUpdated = async () => {
    await onUpdated?.()
    setShowRecurringSheet(false)
  }

  // Determine title based on mode
  const getTitle = () => {
    if (editMode === 'select') {
      return t('edit.expense_title')
    }
    return t('edit.edit_expense')
  }

  // Determine footer based on mode
  const getFooter = () => {
    if (editMode === 'select') {
      // Mode selection - just a cancel button
      return (
        <Button variant="ghost" className="w-full" onClick={handleClose}>
          {tCommon('cancel')}
        </Button>
      )
    }

    // Edit form - save and cancel buttons
    return (
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
          form="edit-expense-form"
          variant="primary"
          className="flex-1"
          disabled={isSubmitting}
        >
          {isSubmitting ? tCommon('saving') : tCommon('save')}
        </form.SubmitButton>
      </div>
    )
  }

  // Render edit occurrence form
  const renderEditForm = () => (
    <form.FormRoot id="edit-expense-form" className="space-y-4">
      {/* Amount Field */}
      <form.AppField
        name="amount"
        validators={{
          onChange: ({ value }) => {
            if (!value) return t('validation.amount_required')
            const num = parseFloat(value)
            if (isNaN(num) || num <= 0) return t('validation.amount_positive')
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

      {/* Date Field */}
      <form.AppField name="date">
        {(field) => <field.DateField label={t('form.date')} required />}
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
    </form.FormRoot>
  )

  return (
    <>
      <form.AppForm>
        <BottomSheet
          isOpen={isOpen}
          onClose={handleClose}
          title={getTitle()}
          closeOnOverlayClick={!isSubmitting}
          closeOnEscape={!isSubmitting}
          footer={getFooter()}
        >
          {editMode === 'select' ? (
            <ModeSelection
              onEditOccurrence={handleEditOccurrence}
              onEditTemplate={handleEditTemplate}
              t={t}
            />
          ) : (
            renderEditForm()
          )}
        </BottomSheet>
      </form.AppForm>

      {/* Recurring Expense Template Sheet */}
      <EditRecurringExpenseSheet
        isOpen={showRecurringSheet}
        onClose={handleRecurringSheetClose}
        recurringExpenseId={expense?.recurringExpenseId}
        businessDescriptor={businessDescriptor}
        onUpdated={handleRecurringUpdated}
      />
    </>
  )
}
