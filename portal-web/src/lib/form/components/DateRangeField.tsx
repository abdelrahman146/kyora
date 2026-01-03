import { useTranslation } from 'react-i18next'
import { useMemo } from 'react'
import type { DateRange } from 'react-day-picker'

import { DateRangePicker } from '@/components/atoms/DateRangePicker'
import { useFieldContext } from '@/lib/form'

interface DateRangeFieldProps {
  label?: string
  placeholder?: string
  minDate?: Date
  maxDate?: Date
  disabledDates?: Array<Date>
  required?: boolean
  disabled?: boolean
  helperText?: string
  size?: 'sm' | 'md' | 'lg'
  numberOfMonths?: number
}

export function DateRangeField({
  label,
  placeholder,
  minDate,
  maxDate,
  disabledDates,
  required,
  disabled,
  helperText,
  size = 'md',
  numberOfMonths = 2,
}: DateRangeFieldProps) {
  const field = useFieldContext<DateRange | undefined>()
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
    <DateRangePicker
      id={field.name}
      label={label}
      value={field.state.value}
      onChange={(range: DateRange | undefined) => field.handleChange(range)}
      onBlur={field.handleBlur}
      minDate={minDate}
      maxDate={maxDate}
      disabledDates={disabledDates}
      placeholder={placeholder}
      required={required}
      disabled={disabled}
      error={showError ? error : undefined}
      helperText={helperText}
      size={size}
      numberOfMonths={numberOfMonths}
    />
  )
}
