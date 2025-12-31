import { forwardRef, useEffect, useId, useRef, useState } from 'react'
import { ChevronDown, ChevronUp, Clock } from 'lucide-react'
import { format, isValid } from 'date-fns'
import { useTranslation } from 'react-i18next'
import type { InputHTMLAttributes } from 'react'
import { cn } from '@/lib/utils'

export interface TimePickerProps
  extends Omit<InputHTMLAttributes<HTMLInputElement>, 'value' | 'onChange' | 'size'> {
  label?: string
  value?: Date
  onChange?: (date: Date | undefined) => void
  onBlur?: () => void
  error?: string
  helperText?: string
  use24Hour?: boolean
  minuteStep?: number
  clearable?: boolean
  fullWidth?: boolean
  size?: 'sm' | 'md' | 'lg'
}

/**
 * TimePicker - Production-grade time input with hour/minute controls
 *
 * Features:
 * - Two numeric inputs for hours and minutes
 * - 12-hour format with AM/PM toggle (or 24-hour based on locale)
 * - Arrow buttons for increment/decrement
 * - Keyboard navigation (Arrow Up/Down, Tab, Type to set)
 * - Auto-advance from hours to minutes after valid input
 * - Mobile-optimized with large touch targets
 * - Accessible with ARIA attributes
 * - Validation for valid time ranges
 * - RTL support
 *
 * @example
 * <TimePicker
 *   label="Appointment Time"
 *   value={appointmentTime}
 *   onChange={setAppointmentTime}
 *   minuteStep={15}
 * />
 */
export const TimePicker = forwardRef<HTMLInputElement, TimePickerProps>(
  (
    {
      label,
      value,
      onChange,
      onBlur,
      error,
      helperText,
      use24Hour = false,
      minuteStep = 1,
      clearable = true,
      fullWidth = true,
      size = 'md',
      className,
      id,
      disabled,
      required,
      ...props
    },
    ref
  ) => {
    const generatedId = useId()
    const inputId = id ?? generatedId
    const { t } = useTranslation()

    // Extract hours and minutes from Date
    const [hours, setHours] = useState<string>('')
    const [minutes, setMinutes] = useState<string>('')
    const [period, setPeriod] = useState<'AM' | 'PM'>('AM')

    const hoursInputRef = useRef<HTMLInputElement>(null)
    const minutesInputRef = useRef<HTMLInputElement>(null)

    // Initialize from value prop
    useEffect(() => {
      if (value && isValid(value)) {
        const h = value.getHours()
        const m = value.getMinutes()

        if (use24Hour) {
          setHours(h.toString().padStart(2, '0'))
        } else {
          const hour12 = h === 0 ? 12 : h > 12 ? h - 12 : h
          setHours(hour12.toString().padStart(2, '0'))
          setPeriod(h >= 12 ? 'PM' : 'AM')
        }
        setMinutes(m.toString().padStart(2, '0'))
      } else {
        setHours('')
        setMinutes('')
        setPeriod('AM')
      }
    }, [value, use24Hour])

    // Size classes
    const sizeClasses = {
      sm: 'h-[44px] text-sm',
      md: 'h-[50px] text-base',
      lg: 'h-[56px] text-lg',
    }

    // Build Date from current inputs
    const buildDate = (h: string, m: string, p: 'AM' | 'PM'): Date | undefined => {
      const hourNum = parseInt(h, 10)
      const minuteNum = parseInt(m, 10)

      if (isNaN(hourNum) || isNaN(minuteNum)) {
        return undefined
      }

      // Validate ranges
      if (use24Hour) {
        if (hourNum < 0 || hourNum > 23) return undefined
      } else {
        if (hourNum < 1 || hourNum > 12) return undefined
      }
      if (minuteNum < 0 || minuteNum > 59) return undefined

      // Convert to 24-hour format if needed
      let hour24 = hourNum
      if (!use24Hour) {
        if (p === 'PM' && hourNum !== 12) {
          hour24 = hourNum + 12
        } else if (p === 'AM' && hourNum === 12) {
          hour24 = 0
        }
      }

      // Create date with current date and selected time
      const date = new Date()
      date.setHours(hour24, minuteNum, 0, 0)
      return date
    }

    // Handle hours change
    const handleHoursChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      const val = e.target.value.replace(/\D/g, '') // Only digits
      if (val.length <= 2) {
        setHours(val)
        const newDate = buildDate(val, minutes, period)
        if (newDate) {
          onChange?.(newDate)
        }

        // Auto-advance to minutes after valid hour
        if (val.length === 2) {
          const hourNum = parseInt(val, 10)
          const maxHour = use24Hour ? 23 : 12
          if (hourNum >= 0 && hourNum <= maxHour) {
            minutesInputRef.current?.focus()
          }
        }
      }
    }

    // Handle minutes change
    const handleMinutesChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      const val = e.target.value.replace(/\D/g, '')
      if (val.length <= 2) {
        setMinutes(val)
        const newDate = buildDate(hours, val, period)
        if (newDate) {
          onChange?.(newDate)
        }
      }
    }

    // Increment/Decrement hours
    const adjustHours = (delta: number) => {
      const currentHour = parseInt(hours, 10) || 0
      const maxHour = use24Hour ? 23 : 12
      const minHour = use24Hour ? 0 : 1
      let newHour = currentHour + delta

      if (newHour > maxHour) newHour = minHour
      if (newHour < minHour) newHour = maxHour

      const newHours = newHour.toString().padStart(2, '0')
      setHours(newHours)
      const newDate = buildDate(newHours, minutes, period)
      if (newDate) onChange?.(newDate)
    }

    // Increment/Decrement minutes
    const adjustMinutes = (delta: number) => {
      const currentMinute = parseInt(minutes, 10) || 0
      let newMinute = currentMinute + delta * minuteStep

      if (newMinute >= 60) newMinute = 0
      if (newMinute < 0) newMinute = 60 - minuteStep

      const newMinutes = newMinute.toString().padStart(2, '0')
      setMinutes(newMinutes)
      const newDate = buildDate(hours, newMinutes, period)
      if (newDate) onChange?.(newDate)
    }

    // Toggle AM/PM
    const togglePeriod = () => {
      const newPeriod = period === 'AM' ? 'PM' : 'AM'
      setPeriod(newPeriod)
      const newDate = buildDate(hours, minutes, newPeriod)
      if (newDate) onChange?.(newDate)
    }

    // Handle keyboard navigation
    const handleHoursKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
      if (e.key === 'ArrowUp') {
        e.preventDefault()
        adjustHours(1)
      } else if (e.key === 'ArrowDown') {
        e.preventDefault()
        adjustHours(-1)
      }
    }

    const handleMinutesKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
      if (e.key === 'ArrowUp') {
        e.preventDefault()
        adjustMinutes(1)
      } else if (e.key === 'ArrowDown') {
        e.preventDefault()
        adjustMinutes(-1)
      }
    }

    // Handle blur
    const handleBlur = () => {
      onBlur?.()
      // Pad with zeros if valid
      if (hours && hours.length === 1) {
        setHours(hours.padStart(2, '0'))
      }
      if (minutes && minutes.length === 1) {
        setMinutes(minutes.padStart(2, '0'))
      }
    }

    return (
      <div className={cn('form-control', fullWidth && 'w-full')}>
        {label && (
          <label htmlFor={inputId} className="label">
            <span className="label-text text-base-content/70 font-medium">
              {label}
              {required && <span className="text-error ms-1">*</span>}
            </span>
          </label>
        )}

        <div
          className={cn(
            'flex items-center gap-1 input input-bordered transition-all duration-200',
            sizeClasses[size],
            error && 'input-error border-error focus-within:border-error focus-within:ring-error/20',
            disabled && 'opacity-50 cursor-not-allowed',
            className
          )}
        >
          {/* Clock Icon */}
          <div className="flex-shrink-0 text-base-content/50">
            <Clock size={20} aria-hidden="true" />
          </div>

          {/* Hours Input */}
          <div className="relative flex flex-col items-center">
            <button
              type="button"
              onClick={() => adjustHours(1)}
              disabled={disabled}
              className="btn btn-ghost btn-xs p-0 h-4 min-h-0"
              aria-label={t('common.incrementHours') || 'Increment hours'}
              tabIndex={-1}
            >
              <ChevronUp size={12} aria-hidden="true" />
            </button>
            <input
              ref={hoursInputRef}
              type="text"
              inputMode="numeric"
              value={hours}
              onChange={handleHoursChange}
              onKeyDown={handleHoursKeyDown}
              onBlur={handleBlur}
              disabled={disabled}
              placeholder="HH"
              aria-label={t('common.hours') || 'Hours'}
              className="w-10 text-center bg-transparent border-0 focus:outline-none focus:ring-0 p-0"
              maxLength={2}
            />
            <button
              type="button"
              onClick={() => adjustHours(-1)}
              disabled={disabled}
              className="btn btn-ghost btn-xs p-0 h-4 min-h-0"
              aria-label={t('common.decrementHours') || 'Decrement hours'}
              tabIndex={-1}
            >
              <ChevronDown size={12} aria-hidden="true" />
            </button>
          </div>

          {/* Colon Separator */}
          <span className="text-base-content/50 font-semibold">:</span>

          {/* Minutes Input */}
          <div className="relative flex flex-col items-center">
            <button
              type="button"
              onClick={() => adjustMinutes(1)}
              disabled={disabled}
              className="btn btn-ghost btn-xs p-0 h-4 min-h-0"
              aria-label={t('common.incrementMinutes') || 'Increment minutes'}
              tabIndex={-1}
            >
              <ChevronUp size={12} aria-hidden="true" />
            </button>
            <input
              ref={minutesInputRef}
              type="text"
              inputMode="numeric"
              value={minutes}
              onChange={handleMinutesChange}
              onKeyDown={handleMinutesKeyDown}
              onBlur={handleBlur}
              disabled={disabled}
              placeholder="MM"
              aria-label={t('common.minutes') || 'Minutes'}
              className="w-10 text-center bg-transparent border-0 focus:outline-none focus:ring-0 p-0"
              maxLength={2}
            />
            <button
              type="button"
              onClick={() => adjustMinutes(-1)}
              disabled={disabled}
              className="btn btn-ghost btn-xs p-0 h-4 min-h-0"
              aria-label={t('common.decrementMinutes') || 'Decrement minutes'}
              tabIndex={-1}
            >
              <ChevronDown size={12} aria-hidden="true" />
            </button>
          </div>

          {/* AM/PM Toggle (12-hour format only) */}
          {!use24Hour && (
            <>
              <div className="w-px h-6 bg-base-300 mx-1" />
              <button
                type="button"
                onClick={togglePeriod}
                disabled={disabled}
                className="btn btn-ghost btn-sm min-h-0 h-auto px-2"
                aria-label={`${t('common.period')}: ${period}`}
              >
                <span className="font-semibold">{period}</span>
              </button>
            </>
          )}

          {/* Hidden input for form submission */}
          <input
            ref={ref}
            type="hidden"
            id={inputId}
            value={value ? format(value, 'HH:mm') : ''}
            aria-invalid={!!error}
            aria-describedby={error ? `${inputId}-error` : helperText ? `${inputId}-helper` : undefined}
            {...props}
          />
        </div>

        {/* Error Message */}
        {error && (
          <label className="label">
            <span
              id={`${inputId}-error`}
              className="label-text-alt text-error"
              role="alert"
            >
              {error}
            </span>
          </label>
        )}

        {/* Helper Text */}
        {helperText && !error && (
          <label className="label">
            <span id={`${inputId}-helper`} className="label-text-alt text-base-content/60">
              {helperText}
            </span>
          </label>
        )}
      </div>
    )
  }
)

TimePicker.displayName = 'TimePicker'
