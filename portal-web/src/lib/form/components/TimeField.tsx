import { useFieldContext } from '@/lib/form/contexts'
import { TimePicker, type TimePickerProps } from '@/components/atoms/TimePicker'

export interface TimeFieldProps extends Omit<TimePickerProps, 'value' | 'onChange' | 'error'> {
  /**
   * If true, shows field label from the label prop
   * If false, no label is rendered (useful when label is in parent)
   * @default true
   */
  showLabel?: boolean
}

/**
 * TimeField - TanStack Form-integrated time picker
 *
 * Pre-bound to field context from useKyoraForm.
 * Automatically handles value, onChange, onBlur, and error display.
 *
 * @example
 * <form.AppField name="appointmentTime" validators={{ onBlur: z.date() }}>
 *   {(field) => <field.TimeField label="Appointment Time" minuteStep={15} />}
 * </form.AppField>
 */
export function TimeField({ showLabel = true, label, ...props }: TimeFieldProps) {
  const field = useFieldContext<Date>()

  return (
    <TimePicker
      {...props}
      label={showLabel ? label : undefined}
      value={field.state.value}
      onChange={(date) => field.handleChange(date)}
      onBlur={field.handleBlur}
      error={field.state.meta.isTouched ? field.state.meta.errors.join(', ') : undefined}
    />
  )
}
