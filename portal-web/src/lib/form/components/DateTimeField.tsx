import { useFieldContext } from '@/lib/form/contexts'
import { DatePicker, type DatePickerProps } from '@/components/atoms/DatePicker'
import { TimePicker, type TimePickerProps } from '@/components/atoms/TimePicker'
import { useMediaQuery } from '@/hooks/useMediaQuery'
import { cn } from '@/lib/utils'

export interface DateTimeFieldProps {
  /**
   * Display mode for the field
   * - 'date': Only show date picker
   * - 'time': Only show time picker
   * - 'datetime': Show both date and time pickers
   * @default 'datetime'
   */
  mode?: 'date' | 'time' | 'datetime'

  /**
   * Label for the field
   */
  label?: string

  /**
   * If true, shows field label from the label prop
   * If false, no label is rendered (useful when label is in parent)
   * @default true
   */
  showLabel?: boolean

  /**
   * Helper text displayed below the field
   */
  helperText?: string

  /**
   * Props to pass to the DatePicker (when mode includes date)
   */
  datePickerProps?: Omit<DatePickerProps, 'value' | 'onChange' | 'onBlur' | 'error' | 'label'>

  /**
   * Props to pass to the TimePicker (when mode includes time)
   */
  timePickerProps?: Omit<TimePickerProps, 'value' | 'onChange' | 'onBlur' | 'error' | 'label'>

  /**
   * If true, field takes full width of container
   * @default true
   */
  fullWidth?: boolean

  /**
   * Size variant
   * @default 'md'
   */
  size?: 'sm' | 'md' | 'lg'

  /**
   * If true, field is disabled
   */
  disabled?: boolean

  /**
   * If true, field is required
   */
  required?: boolean
}

/**
 * DateTimeField - TanStack Form-integrated date/time picker with multiple modes
 *
 * Pre-bound to field context from useKyoraForm.
 * Supports date-only, time-only, or combined date+time input.
 * Responsive layout: side-by-side on desktop, stacked on mobile.
 *
 * @example
 * // Date only
 * <form.AppField name="birthdate" validators={{ onBlur: z.date() }}>
 *   {(field) => <field.DateTimeField mode="date" label="Birth Date" />}
 * </form.AppField>
 *
 * // Time only
 * <form.AppField name="appointmentTime" validators={{ onBlur: z.date() }}>
 *   {(field) => <field.DateTimeField mode="time" label="Appointment Time" />}
 * </form.AppField>
 *
 * // Date and time
 * <form.AppField name="eventDateTime" validators={{ onBlur: z.date() }}>
 *   {(field) => <field.DateTimeField mode="datetime" label="Event Date & Time" />}
 * </form.AppField>
 */
export function DateTimeField({
  mode = 'datetime',
  label,
  showLabel = true,
  helperText,
  datePickerProps,
  timePickerProps,
  fullWidth = true,
  size = 'md',
  disabled,
  required,
}: DateTimeFieldProps) {
  const field = useFieldContext<Date>()
  const isMobile = useMediaQuery('(max-width: 768px)')

  const error = field.state.meta.isTouched ? field.state.meta.errors.join(', ') : undefined

  // Date-only mode
  if (mode === 'date') {
    return (
      <DatePicker
        {...datePickerProps}
        label={showLabel ? label : undefined}
        helperText={helperText}
        value={field.state.value}
        onChange={(date: Date | undefined) => field.handleChange(date ?? field.state.value)}
        onBlur={field.handleBlur}
        error={error}
        fullWidth={fullWidth}
        size={size}
        disabled={disabled}
        required={required}
      />
    )
  }

  // Time-only mode
  if (mode === 'time') {
    return (
      <TimePicker
        {...timePickerProps}
        label={showLabel ? label : undefined}
        helperText={helperText}
        value={field.state.value}
        onChange={(date: Date | undefined) => field.handleChange(date ?? field.state.value)}
        onBlur={field.handleBlur}
        error={error}
        fullWidth={fullWidth}
        size={size}
        disabled={disabled}
        required={required}
      />
    )
  }

  // DateTime mode (both date and time)
  return (
    <div className={cn('form-control', fullWidth && 'w-full')}>
      {showLabel && label && (
        <label className="label">
          <span className="label-text text-base-content/70 font-medium">
            {label}
            {required && <span className="text-error ms-1">*</span>}
          </span>
        </label>
      )}

      <div
        className={cn(
          'flex gap-3',
          isMobile ? 'flex-col' : 'flex-row items-start'
        )}
      >
        {/* Date Picker */}
        <div className={cn(isMobile ? 'w-full' : 'flex-1')}>
          <DatePicker
            {...datePickerProps}
            value={field.state.value}
            onChange={(date: Date | undefined) => {
              // Preserve time when changing date
              if (date && field.state.value) {
                date.setHours(field.state.value.getHours())
                date.setMinutes(field.state.value.getMinutes())
              }
              if (date) field.handleChange(date)
            }}
            onBlur={field.handleBlur}
            error={undefined} // Only show error once at bottom
            fullWidth
            size={size}
            disabled={disabled}
            required={required}
          />
        </div>

        {/* Time Picker */}
        <div className={cn(isMobile ? 'w-full' : 'flex-1')}>
          <TimePicker
            {...timePickerProps}
            value={field.state.value}
            onChange={(date: Date | undefined) => {
              // Preserve date when changing time
              if (date && field.state.value) {
                const newDate = new Date(field.state.value)
                newDate.setHours(date.getHours())
                newDate.setMinutes(date.getMinutes())
                field.handleChange(newDate)
              }
            }}
            onBlur={field.handleBlur}
            error={undefined} // Only show error once at bottom
            fullWidth
            size={size}
            disabled={disabled}
            required={required}
          />
        </div>
      </div>

      {/* Error Message */}
      {error && (
        <label className="label">
          <span className="label-text-alt text-error" role="alert">
            {error}
          </span>
        </label>
      )}

      {/* Helper Text */}
      {helperText && !error && (
        <label className="label">
          <span className="label-text-alt text-base-content/60">
            {helperText}
          </span>
        </label>
      )}
    </div>
  )
}
