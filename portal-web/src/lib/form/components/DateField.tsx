import type { DatePickerProps } from '@/components/atoms/DatePicker'
import { useFieldContext } from '@/lib/form/contexts'
import { DatePicker } from '@/components/atoms/DatePicker'

export interface DateFieldProps extends Omit<
  DatePickerProps,
  'value' | 'onChange' | 'error'
> {
  /**
   * If true, shows field label from the label prop
   * If false, no label is rendered (useful when label is in parent)
   * @default true
   */
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
 *   {(field) => <field.DateField label="Birth Date" minAge={18} />}
 * </form.AppField>
 */
export function DateField({
  showLabel = true,
  label,
  ...props
}: DateFieldProps) {
  const field = useFieldContext<Date>()

  return (
    <DatePicker
      {...props}
      label={showLabel ? label : undefined}
      value={field.state.value}
      onChange={(date) => field.handleChange(date ?? new Date())}
      onBlur={field.handleBlur}
      error={
        field.state.meta.isTouched
          ? field.state.meta.errors.join(', ')
          : undefined
      }
    />
  )
}
