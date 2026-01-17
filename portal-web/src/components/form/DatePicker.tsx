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
import { Calendar, X } from 'lucide-react'
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
  /** Button mode: Render as button instead of input (no clipping, modal-based) */
  buttonMode?: boolean
  /** Custom button content (only in buttonMode) */
  buttonContent?: React.ReactNode
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
      buttonMode = false,
      buttonContent,
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
    const { t } = useTranslation('common')
    const [isOpen, setIsOpen] = useState(false)
    const [isAnimating, setIsAnimating] = useState(false)
    const [displayMonth, setDisplayMonth] = useState<Date>(value || new Date())
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

    const currentMonth = displayMonth.getMonth()
    const currentYear = displayMonth.getFullYear()

    const months: Array<{ value: number; label: string }> = Array.from(
      { length: 12 },
      (_, monthIndex) => ({
        value: monthIndex,
        label: format(new Date(2020, monthIndex, 1), 'MMMM', { locale }),
      }),
    )

    const startYear = minDate ? minDate.getFullYear() : currentYear - 100
    const endYear = maxDate ? maxDate.getFullYear() : currentYear + 20
    const years: Array<number> = []
    for (
      let year = Math.min(startYear, endYear);
      year <= Math.max(startYear, endYear);
      year++
    ) {
      years.push(year)
    }

    const handleMonthChange = (month: number) => {
      setDisplayMonth(new Date(currentYear, month, 1))
    }

    const handleYearChange = (year: number) => {
      setDisplayMonth(new Date(year, currentMonth, 1))
    }

    const MonthYearSelector = () => (
      <div className="flex items-center justify-center gap-3 mb-4 px-1">
        <div className="relative flex-1">
          <select
            value={currentMonth}
            onChange={(e) => handleMonthChange(Number(e.target.value))}
            className={cn(
              'w-full appearance-none cursor-pointer transition-all duration-200',
              'bg-base-100 text-base-content border border-base-300 rounded-lg',
              'focus:outline-none focus:border-primary focus:ring-2 focus:ring-primary/20',
              'hover:border-base-400',
              isMobile
                ? 'h-[44px] text-base px-3 pe-8'
                : 'h-[38px] text-sm px-3 pe-8',
            )}
            aria-label={t('select_month', {
              defaultValue: 'Select month',
            })}
          >
            {months.map((month) => (
              <option key={month.value} value={month.value}>
                {month.label}
              </option>
            ))}
          </select>
          <div className="absolute inset-y-0 end-0 flex items-center pe-2 pointer-events-none text-base-content/50">
            <svg
              className="h-4 w-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              aria-hidden="true"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M19 9l-7 7-7-7"
              />
            </svg>
          </div>
        </div>
        <div
          className="relative"
          style={{ width: isMobile ? '110px' : '100px' }}
        >
          <select
            value={currentYear}
            onChange={(e) => handleYearChange(Number(e.target.value))}
            className={cn(
              'w-full appearance-none cursor-pointer transition-all duration-200',
              'bg-base-100 text-base-content border border-base-300 rounded-lg',
              'focus:outline-none focus:border-primary focus:ring-2 focus:ring-primary/20',
              'hover:border-base-400',
              isMobile
                ? 'h-[44px] text-base px-3 pe-8'
                : 'h-[38px] text-sm px-3 pe-8',
            )}
            aria-label={t('select_year', {
              defaultValue: 'Select year',
            })}
          >
            {years.map((year) => (
              <option key={year} value={year}>
                {year}
              </option>
            ))}
          </select>
          <div className="absolute inset-y-0 end-0 flex items-center pe-2 pointer-events-none text-base-content/50">
            <svg
              className="h-4 w-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              aria-hidden="true"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M19 9l-7 7-7-7"
              />
            </svg>
          </div>
        </div>
      </div>
    )

    const calendarContent = (
      <div className="w-full">
        <MonthYearSelector />
        <DayPicker
          mode="single"
          selected={value}
          onSelect={handleDateSelect}
          disabled={disabledMatcher}
          month={displayMonth}
          onMonthChange={setDisplayMonth}
          locale={locale}
          dir={isRTL ? 'rtl' : 'ltr'}
          className="kyora-datepicker w-full"
          classNames={{
            months: 'flex flex-col gap-4 w-full',
            month: 'space-y-3 w-full',
            month_caption: 'hidden',
            caption_label: 'hidden',
            nav: 'hidden',
            button_previous: 'hidden',
            button_next: 'hidden',
            month_grid: 'w-full border-collapse',
            weekdays: 'flex w-full',
            weekday: cn(
              'text-base-content/60 flex-1',
              'font-semibold flex items-center justify-center',
              'text-sm h-10',
            ),
            week: 'flex w-full mt-1',
            day: cn(
              'relative p-0 text-center flex-1',
              'flex items-center justify-center',
              'font-normal transition-colors',
              'h-10',
            ),
            day_button: cn(
              'w-full h-full p-0 font-normal rounded-lg',
              'hover:bg-base-200 transition-colors',
              'text-sm',
            ),
            selected: 'bg-primary text-primary-content hover:bg-primary/90',
            today: 'border-2 border-primary',
            outside: 'text-base-content/20 opacity-40',
            disabled:
              'text-base-content/20 opacity-40 cursor-not-allowed hover:bg-transparent',
            hidden: 'invisible',
          }}
        />
      </div>
    )

    // Button mode: renders as button with modal picker (no clipping)
    if (buttonMode) {
      return (
        <>
          <button
            ref={(node) => {
              if (typeof ref === 'function') {
                ref(node as any)
              } else if (ref) {
                ;(ref as any).current = node
              }
            }}
            type="button"
            onClick={() => !disabled && handleOpen()}
            disabled={disabled}
            aria-label={label}
            aria-expanded={isOpen}
            className={cn(
              'flex items-center gap-2 rounded-btn px-3 py-2 text-sm transition-colors cursor-pointer',
              'text-base-content/70 hover:bg-base-200 hover:text-base-content',
              disabled && 'opacity-50 cursor-not-allowed',
              className,
            )}
          >
            {buttonContent || (
              <>
                <Calendar className="h-4 w-4" aria-hidden="true" />
                <span className="whitespace-nowrap">
                  {value
                    ? format(value, dateFormat, { locale })
                    : label || t('date.selectDate')}
                </span>
              </>
            )}
          </button>

          {isOpen &&
            !disabled &&
            createPortal(
              <>
                {/* Backdrop */}
                <div
                  className="fixed inset-0 z-[100] bg-black/20 backdrop-blur-sm"
                  onClick={handleClose}
                  aria-hidden="true"
                />

                {/* Calendar Modal */}
                <div
                  ref={popupRef}
                  role="dialog"
                  aria-modal="true"
                  aria-label={label || t('date.selectDate')}
                  className={cn(
                    'fixed z-[101] bg-base-100 shadow-2xl',
                    isMobile
                      ? 'inset-x-0 bottom-0 rounded-t-3xl'
                      : 'left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 rounded-2xl',
                  )}
                >
                  {/* Header */}
                  <div className="flex items-center justify-between border-b border-base-300 px-4 py-3">
                    <h3 className="text-base font-semibold">
                      {label || t('date.selectDate')}
                    </h3>
                    <button
                      type="button"
                      onClick={handleClose}
                      className="btn btn-ghost btn-sm btn-square"
                      aria-label={t('close')}
                    >
                      <X className="h-4 w-4" />
                    </button>
                  </div>

                  {/* Calendar Content */}
                  <div className="p-4">{calendarContent}</div>
                </div>
              </>,
              document.body,
            )}
        </>
      )
    }

    // Input mode: standard input field with dropdown picker
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
              aria-label={t('clear')}
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
                aria-label={label || t('date.selectDate')}
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
                    'relative bg-base-100 rounded-t-xl',
                    'w-full max-h-[85vh] overflow-y-auto',
                    'transition-transform duration-200',
                    isAnimating ? 'translate-y-0' : 'translate-y-full',
                  )}
                >
                  <div className="sticky top-0 z-10 bg-base-100 border-b border-base-300 px-4 py-3">
                    <div className="flex justify-between items-center">
                      <h3 className="text-lg font-semibold">
                        {label || t('date.selectDate')}
                      </h3>
                      <button
                        type="button"
                        onClick={handleClose}
                        className="btn btn-ghost btn-sm btn-circle"
                        aria-label={t('close')}
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
                      {t('date.done')}
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
              aria-label={label || t('date.selectDate')}
              className={cn(
                'absolute top-full mt-2 z-50',
                'bg-base-100 rounded-box border border-base-300 p-4',
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
