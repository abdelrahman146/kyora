import { forwardRef, useEffect, useId, useRef, useState } from 'react'
import { DayPicker } from 'react-day-picker'
import 'react-day-picker/style.css'
import { format, isValid, parse } from 'date-fns'
import { ar, enUS } from 'date-fns/locale'
import { Calendar, X } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import type { InputHTMLAttributes } from 'react'
import { cn } from '@/lib/utils'
import { useLanguage } from '@/hooks/useLanguage'

export interface DatePickerProps
  extends Omit<InputHTMLAttributes<HTMLInputElement>, 'value' | 'onChange' | 'size'> {
  label?: string
  value?: Date
  onChange?: (date: Date | undefined) => void
  onBlur?: () => void
  error?: string
  helperText?: string
  minDate?: Date
  maxDate?: Date
  disabledDates?: Array<Date>
  disableWeekends?: boolean
  clearable?: boolean
  fullWidth?: boolean
  size?: 'sm' | 'md' | 'lg'
}

/**
 * DatePicker - Production-grade date input with calendar popup
 *
 * Features:
 * - React Day Picker integration with daisyUI styling
 * - RTL support for Arabic (calendar flows right-to-left)
 * - Locale-aware date formatting (DD/MM/YYYY for Arabic, MM/DD/YYYY for English)
 * - Keyboard navigation (Arrow keys, Page Up/Down, Home/End, Enter, Escape)
 * - Mobile-optimized with full-screen modal on small screens
 * - Accessible with ARIA attributes and focus management
 * - Clear button to reset value
 * - Min/max date validation
 * - Disabled dates/weekends support
 *
 * @example
 * <DatePicker
 *   label="Birth Date"
 *   value={birthDate}
 *   onChange={setBirthDate}
 *   maxDate={new Date()}
 *   clearable
 * />
 */
export const DatePicker = forwardRef<HTMLInputElement, DatePickerProps>(
  (
    {
      label,
      value,
      onChange,
      onBlur,
      error,
      helperText,
      minDate,
      maxDate,
      disabledDates = [],
      disableWeekends = false,
      clearable = true,
      fullWidth = true,
      size = 'md',
      className,
      id,
      disabled,
      required,
      placeholder,
      ...props
    },
    ref
  ) => {
    const generatedId = useId()
    const inputId = id ?? generatedId
    const popupId = `${inputId}-popup`
    const { isRTL } = useLanguage()
    const { t } = useTranslation()
    const [isOpen, setIsOpen] = useState(false)
    const [inputValue, setInputValue] = useState('')
    const containerRef = useRef<HTMLDivElement>(null)
    const inputRef = useRef<HTMLInputElement>(null)
    const popupRef = useRef<HTMLDivElement>(null)

    // Date format based on locale
    const dateFormat = isRTL ? 'dd/MM/yyyy' : 'MM/dd/yyyy'
    const locale = isRTL ? ar : enUS

    // Initialize input value from prop
    useEffect(() => {
      if (value && isValid(value)) {
        setInputValue(format(value, dateFormat, { locale }))
      } else {
        setInputValue('')
      }
    }, [value, dateFormat, locale])

    // Size classes
    const sizeClasses = {
      sm: 'h-[44px] text-sm',
      md: 'h-[50px] text-base',
      lg: 'h-[56px] text-lg',
    }

    // Handle input change (manual typing)
    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      const newValue = e.target.value
      setInputValue(newValue)

      // Try to parse the date
      const parsedDate = parse(newValue, dateFormat, new Date(), { locale })
      if (isValid(parsedDate)) {
        onChange?.(parsedDate)
      }
    }

    // Handle date selection from calendar
    const handleDateSelect = (date: Date | undefined) => {
      if (date && isValid(date)) {
        setInputValue(format(date, dateFormat, { locale }))
        onChange?.(date)
        setIsOpen(false)
        inputRef.current?.focus()
      }
    }

    // Handle clear
    const handleClear = () => {
      setInputValue('')
      onChange?.(undefined)
      inputRef.current?.focus()
    }

    // Handle input blur
    const handleBlur = () => {
      onBlur?.()
      // If input is invalid, clear it
      const parsedDate = parse(inputValue, dateFormat, new Date(), { locale })
      if (inputValue && !isValid(parsedDate)) {
        setInputValue('')
        onChange?.(undefined)
      }
    }

    // Handle keyboard navigation
    const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
      if (e.key === 'Escape') {
        setIsOpen(false)
        inputRef.current?.focus()
      } else if (e.key === 'ArrowDown' && !isOpen) {
        setIsOpen(true)
      }
    }

    // Close popup on click outside
    useEffect(() => {
      const handleClickOutside = (event: MouseEvent) => {
        if (
          containerRef.current &&
          !containerRef.current.contains(event.target as Node)
        ) {
          setIsOpen(false)
        }
      }

      if (isOpen) {
        document.addEventListener('mousedown', handleClickOutside)
        return () => document.removeEventListener('mousedown', handleClickOutside)
      }
    }, [isOpen])

    // Focus trap for popup
    useEffect(() => {
      if (isOpen && popupRef.current) {
        const focusableElements = popupRef.current.querySelectorAll(
          'button, [tabindex]:not([tabindex="-1"])'
        )
        const firstElement = focusableElements[0] as HTMLElement
        firstElement.focus()
      }
    }, [isOpen])

    // Disabled dates matcher
    const disabledMatcher = (date: Date) => {
      const isBeforeMin = minDate && date < minDate
      const isAfterMax = maxDate && date > maxDate
      const isInDisabledList = disabledDates.some(
        (d: Date) => format(d, 'yyyy-MM-dd') === format(date, 'yyyy-MM-dd')
      )
      const isWeekend = disableWeekends && (date.getDay() === 0 || date.getDay() === 6)

      return !!(isBeforeMin || isAfterMax || isInDisabledList || isWeekend)
    }

    return (
      <div
        ref={containerRef}
        className={cn('form-control relative', fullWidth && 'w-full')}
      >
        {label && (
          <label htmlFor={inputId} className="label">
            <span className="label-text text-base-content/70 font-medium">
              {label}
              {required && <span className="text-error ms-1">*</span>}
            </span>
          </label>
        )}

        <div className="relative">
          {/* Calendar Icon */}
          <div className="absolute inset-y-0 start-0 z-10 flex items-center ps-3 pointer-events-none text-base-content/50">
            <Calendar size={20} aria-hidden="true" />
          </div>

          {/* Input Field */}
          <input
            ref={(node) => {
              inputRef.current = node
              if (typeof ref === 'function') {
                ref(node)
              } else if (ref) {
                ref.current = node
              }
            }}
            type="text"
            id={inputId}
            value={inputValue}
            onChange={handleInputChange}
            onBlur={handleBlur}
            onKeyDown={handleKeyDown}
            onFocus={() => setIsOpen(true)}
            disabled={disabled}
            required={required}
            placeholder={placeholder || dateFormat.toLowerCase()}
            aria-invalid={!!error}
            aria-describedby={error ? `${inputId}-error` : helperText ? `${inputId}-helper` : undefined}
            aria-expanded={isOpen}
            aria-controls={popupId}
            className={cn(
              'input input-bordered relative z-0 w-full transition-all duration-200',
              sizeClasses[size],
              'ps-10', // Space for calendar icon
              clearable && value && 'pe-10', // Space for clear button
              error && 'input-error border-error focus:border-error focus:ring-error/20',
              disabled && 'opacity-50 cursor-not-allowed',
              className
            )}
            {...props}
          />

          {/* Clear Button */}
          {clearable && value && !disabled && (
            <button
              type="button"
              onClick={handleClear}
              className="absolute inset-y-0 end-0 z-10 flex items-center pe-3 text-base-content/50 hover:text-base-content transition-colors"
              aria-label={t('common.clear')}
            >
              <X size={18} aria-hidden="true" />
            </button>
          )}
        </div>

        {/* Calendar Popup */}
        {isOpen && !disabled && (
          <>
            {/* Mobile: Full-screen modal */}
            <div
              className="fixed inset-0 bg-black/50 z-40 md:hidden"
              onClick={() => setIsOpen(false)}
            />
            <div
              ref={popupRef}
              id={popupId}
              role="dialog"
              aria-modal="true"
              aria-label={t('common.selectDate') || 'Select date'}
              className={cn(
                // Mobile: Full-screen modal
                'fixed inset-x-0 bottom-0 z-50 bg-base-100 rounded-t-lg shadow-xl p-4 md:hidden',
                // Desktop: Dropdown popup
                'md:absolute md:top-full md:mt-2 md:z-50 md:bg-base-100 md:rounded-box md:shadow-lg md:p-4 md:block',
                // RTL positioning
                isRTL ? 'md:end-0' : 'md:start-0'
              )}
            >
              {/* Mobile header */}
              <div className="flex justify-between items-center mb-4 md:hidden">
                <h3 className="text-lg font-semibold">{label || t('common.selectDate')}</h3>
                <button
                  type="button"
                  onClick={() => setIsOpen(false)}
                  className="btn btn-ghost btn-sm btn-circle"
                  aria-label={t('common.close')}
                >
                  <X size={20} />
                </button>
              </div>

              {/* Calendar */}
              <DayPicker
                mode="single"
                selected={value}
                onSelect={handleDateSelect}
                disabled={disabledMatcher}
                locale={locale}
                dir={isRTL ? 'rtl' : 'ltr'}
                className={cn(
                  'date-picker',
                  // Custom styling for daisyUI integration
                  '[&_.rdp-day_button]:h-10 [&_.rdp-day_button]:w-10',
                  '[&_.rdp-day_button]:rounded-lg',
                  '[&_.rdp-day_button]:transition-colors',
                  '[&_.rdp-day_button:hover]:bg-base-200',
                  '[&_.rdp-day_button.rdp-selected]:bg-primary [&_.rdp-day_button.rdp-selected]:text-primary-content',
                  '[&_.rdp-day_button.rdp-today]:border [&_.rdp-day_button.rdp-today]:border-primary',
                  '[&_.rdp-day_button:disabled]:text-base-300 [&_.rdp-day_button:disabled]:cursor-not-allowed',
                  '[&_.rdp-nav_button]:h-10 [&_.rdp-nav_button]:w-10',
                  '[&_.rdp-nav_button]:rounded-lg',
                  '[&_.rdp-nav_button]:transition-colors',
                  '[&_.rdp-nav_button:hover]:bg-base-200',
                  '[&_.rdp-month_caption]:font-semibold [&_.rdp-month_caption]:text-base-content',
                  // Mobile optimization
                  'md:[&_.rdp-day_button]:h-12 md:[&_.rdp-day_button]:w-12'
                )}
                modifiersClassNames={{
                  selected: 'rdp-selected',
                  today: 'rdp-today',
                  disabled: 'rdp-disabled',
                }}
              />

              {/* Mobile footer with action buttons */}
              <div className="flex gap-2 mt-4 md:hidden">
                <button
                  type="button"
                  onClick={() => {
                    const today = new Date()
                    if (!disabledMatcher(today)) {
                      handleDateSelect(today)
                    }
                  }}
                  className="btn btn-ghost flex-1"
                >
                  {t('common.today')}
                </button>
                {clearable && (
                  <button
                    type="button"
                    onClick={() => {
                      handleClear()
                      setIsOpen(false)
                    }}
                    className="btn btn-ghost flex-1"
                  >
                    {t('common.clear')}
                  </button>
                )}
              </div>
            </div>
          </>
        )}

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

DatePicker.displayName = 'DatePicker'
