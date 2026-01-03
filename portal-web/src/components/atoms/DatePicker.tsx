/**
 * DatePicker - Production-Grade Date Selection Component
 *
 * Mobile-first date picker with bottom sheet on mobile, dropdown on desktop.
 * Fully aligned with Kyora Design System (KDS) - matches FormInput/FormSelect patterns.
 *
 * Features:
 * - Mobile-first: Bottom sheet (<768px), dropdown (≥768px)
 * - RTL-native: Proper Arabic calendar layout
 * - daisyUI styled: Consistent with form system (input/button classes)
 * - Keyboard accessible: Arrow keys, Enter, Escape, Tab
 * - Touch-optimized: 50px minimum touch targets
 * - React Day Picker 9.x: No default CSS, custom KDS styling
 * - Locale-aware: DD/MM/YYYY (Arabic), MM/DD/YYYY (English)
 * - Validation: Min/max dates, disabled dates/weekends
 * - Clear button: Optional reset functionality
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

import {
  forwardRef,
  useCallback,
  useEffect,
  useId,
  useRef,
  useState,
} from 'react'
import { createPortal } from 'react-dom'
import { DayPicker } from 'react-day-picker'
import { format, isValid, parse } from 'date-fns'
import { ar, enUS } from 'date-fns/locale'
import { Calendar, ChevronLeft, ChevronRight, X } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import type { InputHTMLAttributes } from 'react'
import { cn } from '@/lib/utils'
import { useLanguage } from '@/hooks/useLanguage'
import { useMediaQuery } from '@/hooks/useMediaQuery'

export interface DatePickerProps extends Omit<
  InputHTMLAttributes<HTMLInputElement>,
  'value' | 'onChange' | 'size'
> {
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
    ref,
  ) => {
    const generatedId = useId()
    const inputId = id ?? generatedId
    const popupId = `${inputId}-popup`
    const { isRTL } = useLanguage()
    const { t } = useTranslation()
    const [isOpen, setIsOpen] = useState(false)
    const [isAnimating, setIsAnimating] = useState(false)
    const [inputValue, setInputValue] = useState('')
    const containerRef = useRef<HTMLDivElement>(null)
    const inputRef = useRef<HTMLInputElement>(null)
    const popupRef = useRef<HTMLDivElement>(null)
    const isMobile = useMediaQuery('(max-width: 768px)')

    const dateFormat = isRTL ? 'dd/MM/yyyy' : 'MM/dd/yyyy'
    const locale = isRTL ? ar : enUS

    useEffect(() => {
      if (value && isValid(value)) {
        setInputValue(format(value, dateFormat, { locale }))
      } else {
        setInputValue('')
      }
    }, [value, dateFormat, locale])

    const sizeClasses = {
      sm: 'h-[44px] text-sm',
      md: 'h-[50px] text-base',
      lg: 'h-[56px] text-lg',
    }

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      const newValue = e.target.value
      setInputValue(newValue)

      const parsedDate = parse(newValue, dateFormat, new Date(), { locale })
      if (isValid(parsedDate)) {
        onChange?.(parsedDate)
      }
    }

    const handleDateSelect = (date: Date | undefined) => {
      if (date && isValid(date)) {
        setInputValue(format(date, dateFormat, { locale }))
        onChange?.(date)
        handleClose()
      }
    }

    const handleClear = (e: React.MouseEvent) => {
      e.stopPropagation()
      setInputValue('')
      onChange?.(undefined)
      inputRef.current?.focus()
    }

    const handleBlur = () => {
      onBlur?.()
      const parsedDate = parse(inputValue, dateFormat, new Date(), { locale })
      if (inputValue && !isValid(parsedDate)) {
        setInputValue('')
        onChange?.(undefined)
      }
    }

    const handleOpen = useCallback(() => {
      setIsOpen(true)
      setIsAnimating(true)
    }, [])

    const handleClose = useCallback(() => {
      setIsAnimating(false)
      setTimeout(() => {
        setIsOpen(false)
        inputRef.current?.focus()
      }, 200)
    }, [])

    const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
      if (e.key === 'Escape') {
        if (isOpen) {
          e.preventDefault()
          handleClose()
        }
      } else if (e.key === 'ArrowDown' && !isOpen) {
        e.preventDefault()
        handleOpen()
      }
    }

    useEffect(() => {
      const handleClickOutside = (event: MouseEvent) => {
        if (
          containerRef.current &&
          !containerRef.current.contains(event.target as Node) &&
          popupRef.current &&
          !popupRef.current.contains(event.target as Node)
        ) {
          handleClose()
        }
      }

      if (isOpen && !isMobile) {
        document.addEventListener('mousedown', handleClickOutside)
        return () =>
          document.removeEventListener('mousedown', handleClickOutside)
      }
    }, [isOpen, isMobile, handleClose])

    useEffect(() => {
      if (isOpen && isMobile) {
        const originalOverflow = document.body.style.overflow
        const scrollbarWidth =
          window.innerWidth - document.documentElement.clientWidth

        document.body.style.overflow = 'hidden'
        if (scrollbarWidth > 0) {
          document.body.style.paddingInlineEnd = `${scrollbarWidth}px`
        }

        return () => {
          document.body.style.overflow = originalOverflow
          document.body.style.paddingInlineEnd = ''
        }
      }
    }, [isOpen, isMobile])

    const disabledMatcher = (date: Date) => {
      const isBeforeMin = minDate && date < minDate
      const isAfterMax = maxDate && date > maxDate
      const isInDisabledList = disabledDates.some(
        (d: Date) => format(d, 'yyyy-MM-dd') === format(date, 'yyyy-MM-dd'),
      )
      const isWeekend =
        disableWeekends && (date.getDay() === 0 || date.getDay() === 6)

      return !!(isBeforeMin || isAfterMax || isInDisabledList || isWeekend)
    }

    const calendarContent = (
      <DayPicker
        mode="single"
        selected={value}
        onSelect={handleDateSelect}
        disabled={disabledMatcher}
        locale={locale}
        dir={isRTL ? 'rtl' : 'ltr'}
        className="kyora-datepicker"
        classNames={{
          months: 'flex flex-col gap-4',
          month: 'space-y-4',
          month_caption: 'flex justify-center pt-1 relative items-center h-10',
          caption_label: 'text-base font-semibold text-base-content',
          nav: 'flex items-center gap-1',
          button_previous: cn(
            'btn btn-ghost btn-sm btn-circle absolute',
            'hover:bg-base-200 transition-colors',
            isRTL ? 'end-0' : 'start-0',
          ),
          button_next: cn(
            'btn btn-ghost btn-sm btn-circle absolute',
            'hover:bg-base-200 transition-colors',
            isRTL ? 'start-0' : 'end-0',
          ),
          month_grid: 'w-full border-collapse',
          weekdays: 'flex',
          weekday: cn(
            'text-base-content/60 rounded-md w-10 h-10',
            'font-medium text-sm flex items-center justify-center',
          ),
          week: 'flex w-full mt-2',
          day: cn(
            'relative p-0 text-center text-sm',
            'h-10 w-10 flex items-center justify-center',
            'font-normal transition-colors',
          ),
          day_button: cn(
            'btn btn-ghost btn-sm h-10 w-10 p-0 font-normal',
            'hover:bg-base-200 transition-colors',
          ),
          selected: 'bg-primary text-primary-content hover:bg-primary/90',
          today: 'border-2 border-primary',
          outside: 'text-base-content/30 opacity-50',
          disabled:
            'text-base-content/30 opacity-50 cursor-not-allowed hover:bg-transparent',
          hidden: 'invisible',
        }}
        components={{
          Chevron: ({ orientation }) =>
            orientation === 'left' ? (
              <ChevronLeft className="h-4 w-4" />
            ) : (
              <ChevronRight className="h-4 w-4" />
            ),
        }}
      />
    )

    return (
      <div
        ref={containerRef}
        className={cn('form-control relative', fullWidth && 'w-full')}
      >
        {label && (
          <label htmlFor={inputId} className="label pb-2">
            <span className="label-text text-base-content/70 font-medium">
              {label}
              {required && <span className="text-error ms-1">*</span>}
            </span>
          </label>
        )}

        <div className="relative">
          <div className="absolute inset-y-0 start-0 z-10 flex items-center ps-3 pointer-events-none text-base-content/50">
            <Calendar size={20} aria-hidden="true" />
          </div>

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
            onClick={() => !disabled && handleOpen()}
            disabled={disabled}
            required={required}
            placeholder={placeholder || dateFormat.toLowerCase()}
            aria-invalid={!!error}
            aria-describedby={
              error
                ? `${inputId}-error`
                : helperText
                  ? `${inputId}-helper`
                  : undefined
            }
            aria-expanded={isOpen}
            aria-controls={popupId}
            className={cn(
              'input input-bordered relative z-0 w-full transition-all duration-200',
              'cursor-pointer',
              sizeClasses[size],
              'ps-10',
              clearable && value && 'pe-10',
              error &&
                'input-error border-error focus:border-error focus:ring-error/20',
              disabled && 'opacity-50 cursor-not-allowed',
              isOpen && 'border-primary ring-2 ring-primary/20',
              className,
            )}
            {...props}
          />

          {clearable && value && !disabled && (
            <span
              role="button"
              tabIndex={0}
              onClick={handleClear}
              onKeyDown={(e) => {
                if (e.key === 'Enter' || e.key === ' ') {
                  e.preventDefault()
                  handleClear(e as unknown as React.MouseEvent)
                }
              }}
              className="absolute inset-y-0 end-0 z-10 flex items-center pe-3 text-base-content/50 hover:text-base-content transition-colors cursor-pointer"
              aria-label={t('common:clear')}
            >
              <X size={18} aria-hidden="true" />
            </span>
          )}
        </div>

        {helperText && !error && (
          <label className="label pt-1">
            <span
              className="label-text-alt text-base-content/60"
              id={`${inputId}-helper`}
            >
              {helperText}
            </span>
          </label>
        )}

        {error && (
          <label className="label pt-1">
            <span
              className="label-text-alt text-error"
              id={`${inputId}-error`}
              role="alert"
            >
              {error}
            </span>
          </label>
        )}

        {isOpen &&
          !disabled &&
          (isMobile ? (
            createPortal(
              <div
                className={cn(
                  'fixed inset-0 z-[9999] flex items-end justify-center',
                  'transition-opacity duration-200',
                  isAnimating ? 'opacity-100' : 'opacity-0',
                )}
                role="dialog"
                aria-modal="true"
                aria-label={label || t('common:select_date')}
                onClick={(e) => {
                  if (e.target === e.currentTarget) {
                    handleClose()
                  }
                }}
              >
                <div
                  className="absolute inset-0 bg-base-content/50 backdrop-blur-sm"
                  aria-hidden="true"
                />

                <div
                  ref={popupRef}
                  className={cn(
                    'relative bg-base-100 rounded-t-xl shadow-xl',
                    'w-full max-h-[85vh] overflow-y-auto',
                    'transition-transform duration-200',
                    isAnimating ? 'translate-y-0' : 'translate-y-full',
                  )}
                >
                  <div className="sticky top-0 z-10 bg-base-100 border-b border-base-300 px-4 py-3">
                    <div className="flex justify-between items-center">
                      <h3 className="text-lg font-semibold">
                        {label || t('common:select_date')}
                      </h3>
                      <button
                        type="button"
                        onClick={handleClose}
                        className="btn btn-ghost btn-sm btn-circle"
                        aria-label={t('common:close')}
                      >
                        <X size={20} />
                      </button>
                    </div>
                  </div>

                  <div className="p-4">{calendarContent}</div>

                  <div className="sticky bottom-0 z-10 bg-base-100 border-t border-base-300 px-4 py-3">
                    <button
                      type="button"
                      onClick={handleClose}
                      className="btn btn-primary w-full"
                    >
                      {t('common:done')}
                    </button>
                  </div>
                </div>
              </div>,
              document.body,
            )
          ) : (
            <div
              ref={popupRef}
              id={popupId}
              role="dialog"
              aria-modal="false"
              aria-label={label || t('common:select_date')}
              className={cn(
                'absolute top-full mt-2 z-50',
                'bg-base-100 rounded-box shadow-lg border border-base-300 p-4',
                isRTL ? 'end-0' : 'start-0',
              )}
            >
              {calendarContent}
            </div>
          ))}
      </div>
    )
  },
)

DatePicker.displayName = 'DatePicker'
