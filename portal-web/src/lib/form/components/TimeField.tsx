import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import type { TimePickerProps } from '@/components/atoms/TimePicker'
import { useFieldContext } from '@/lib/form/contexts'
import { TimePicker } from '@/components/atoms/TimePicker'

export interface TimeFieldProps extends Omit<
  TimePickerProps,
  'value' | 'onChange' | 'error'
> {
  showLabel?: boolean
}

/**
 * TimeField - TanStack Form-integrated time picker
 *
 * Pre-bound to field context from useKyoraForm.
 * Automatically handles value, onChange, onBlur, and error display.
 *
 * @example
 * <form.AppField name="appointmentTime" validators={{ onBlur: z.string() }}>
 *   {(field) => <field.TimeField label="Appointment Time" minuteStep={15} />}
 * </form.AppField>
 */
export function TimeField({
  showLabel = true,
  label,
  ...props
}: TimeFieldProps) {
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
    <TimePicker
      {...props}
      label={showLabel ? label : undefined}
      value={field.state.value}
      onChange={(time) => field.handleChange(time ?? '')}
      onBlur={field.handleBlur}
      error={showError ? error : undefined}
    />
  )
}
