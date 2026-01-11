import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import type { DateTimePickerProps } from '@/components/form/DateTimePicker'
import { useFieldContext } from '@/lib/form/contexts'
import { DateTimePicker } from '@/components/form/DateTimePicker'

export interface DateTimeFieldProps extends Omit<
  DateTimePickerProps,
  'value' | 'onChange' | 'error'
> {
  showLabel?: boolean
}

/**
 * DateTimeField - TanStack Form-integrated datetime picker
 *
 * Pre-bound to field context from useKyoraForm.
 * Automatically handles value, onChange, onBlur, and error display.
 * Combines date and time selection with tabbed interface.
 *
 * @example
 * <form.AppField name="eventDateTime" validators={{ onBlur: z.string() }}>
 *   {(field) => <field.DateTimeField label="Event Date & Time" />}
 * </form.AppField>
 */
export function DateTimeField({
  showLabel = true,
  label,
  ...props
}: DateTimeFieldProps) {
  const field = useFieldContext<string>()
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
    <DateTimePicker
      {...props}
      label={showLabel ? label : undefined}
      value={field.state.value}
      onChange={(datetime) => field.handleChange(datetime ?? '')}
      onBlur={field.handleBlur}
      error={showError ? error : undefined}
    />
  )
}
