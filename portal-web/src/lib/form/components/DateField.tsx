import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import type { DatePickerProps } from '@/components/atoms/DatePicker'
import { useFieldContext } from '@/lib/form/contexts'
import { DatePicker } from '@/components/atoms/DatePicker'

export interface DateFieldProps extends Omit<
  DatePickerProps,
  'value' | 'onChange' | 'error'
> {
  showLabel?: boolean
}

/**
 * DateField - TanStack Form-integrated date picker
 *
 * Pre-bound to field context from useKyoraForm.
 * Automatically handles value, onChange, onBlur, and error display.
 *
 * @example
 * <form.AppField name="birthdate" validators={{ onBlur: z.date() }}>
 *   {(field) => <field.DateField label="Birth Date" />}
 * </form.AppField>
 */
export function DateField({
  showLabel = true,
  label,
  ...props
}: DateFieldProps) {
  const field = useFieldContext<Date>()
  const { t } = useTranslation('errors')

  const error = useMemo(() => {
    const errors = field.state.meta.errors
    if (errors.length === 0) return undefined

    const firstError = errors[0]
    if (typeof firstError === 'string') {
      return t(firstError)
    }

    if (
      typeof firstError === 'object' &&
      firstError &&
      'message' in firstError
    ) {
      const errorObj = firstError as { message: string; code?: number }
      return t(errorObj.message)
    }

    return undefined
  }, [field.state.meta.errors, t])

  const showError = field.state.meta.isTouched && error

  return (
    <DatePicker
      {...props}
      label={showLabel ? label : undefined}
      value={field.state.value}
      onChange={(date) => field.handleChange(date ?? new Date())}
      onBlur={field.handleBlur}
      error={showError ? error : undefined}
    />
  )
}
