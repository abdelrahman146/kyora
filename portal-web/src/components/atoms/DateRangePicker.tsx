/**
 * DateRangePicker - Production-Grade Date Range Selection Component
 *
 * Mobile-first date range picker with bottom sheet on mobile, dropdown on desktop.
 * Fully aligned with Kyora Design System (KDS) - matches DatePicker patterns.
 *
 * Features:
 * - Mobile-first: Bottom sheet with single month (<768px), dropdown with two months (≥768px)
 * - RTL-native: Proper Arabic calendar layout
 * - daisyUI styled: Consistent with form system
 * - Range selection: Intuitive from/to date picking
 * - Keyboard accessible: Arrow keys, Enter, Escape
 * - Touch-optimized: 50px minimum touch targets
 * - Quick actions: Clear and Apply buttons
 * - Validation: Min/max dates, disabled dates
 *
 * @example
 * <DateRangePicker
 *   label="Order Date Range"
 *   value={dateRange}
 *   onChange={setDateRange}
 *   minDate={new Date(2020, 0, 1)}
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
import { format } from 'date-fns'
import { ar, enUS } from 'date-fns/locale'
import { Calendar, X } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import type { InputHTMLAttributes } from 'react'
import type { DateRange } from 'react-day-picker'
import { cn } from '@/lib/utils'
import { useLanguage } from '@/hooks/useLanguage'
import { useMediaQuery } from '@/hooks/useMediaQuery'

export interface DateRangePickerProps extends Omit<
  InputHTMLAttributes<HTMLInputElement>,
  'value' | 'onChange' | 'size'
> {
  label?: string
  value?: DateRange
  onChange?: (range: DateRange | undefined) => void
  onBlur?: () => void
  error?: string
  helperText?: string
  minDate?: Date
  maxDate?: Date
  disabledDates?: Array<Date>
  clearable?: boolean
  fullWidth?: boolean
  size?: 'sm' | 'md' | 'lg'
  numberOfMonths?: number
}

export const DateRangePicker = forwardRef<
  HTMLInputElement,
  DateRangePickerProps
>(
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
      clearable = true,
      fullWidth = true,
      size = 'md',
      numberOfMonths,
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
    const [dropdownWidth, setDropdownWidth] = useState<number | undefined>()
    const [displayMonth, setDisplayMonth] = useState<Date>(
      value?.from || new Date(),
    )
    const containerRef = useRef<HTMLDivElement>(null)
    const inputRef = useRef<HTMLInputElement>(null)
    const popupRef = useRef<HTMLDivElement>(null)
    const isMobile = useMediaQuery('(max-width: 768px)')

    const locale = isRTL ? ar : enUS
    const monthCount = numberOfMonths ?? 1

    const currentYear = displayMonth.getFullYear()
    const currentMonth = displayMonth.getMonth()

    const years = Array.from({ length: 100 }, (_, i) => currentYear - 50 + i)
    const months = Array.from({ length: 12 }, (_, i) => {
      const date = new Date(2000, i, 1)
      return {
        value: i,
        label: format(date, 'MMMM', { locale }),
      }
    })

    const sizeClasses = {
      sm: 'h-[44px] text-sm',
      md: 'h-[50px] text-base',
      lg: 'h-[56px] text-lg',
    }

    const formatDateRange = (range?: DateRange): string => {
      if (!range?.from) return ''
      if (!range.to) {
        return format(range.from, 'PP', { locale })
      }
      const fromStr = format(range.from, 'PP', { locale })
      const toStr = format(range.to, 'PP', { locale })
      return `${fromStr} - ${toStr}`
    }

    const handleDateSelect = (range: DateRange | undefined) => {
      onChange?.(range)
      // Don't auto-close - let user manually click Apply to confirm selection
      // This provides better UX especially on mobile where user picks both dates in one session
    }

    const handleClear = (e: React.MouseEvent) => {
      e.stopPropagation()
      onChange?.(undefined)
      inputRef.current?.focus()
    }

    const handleOpen = useCallback(() => {
      if (containerRef.current) {
        setDropdownWidth(containerRef.current.offsetWidth)
      }
      setIsOpen(true)
      setIsAnimating(true)
    }, [])

    const handleClose = useCallback(() => {
      setIsAnimating(false)
      setTimeout(() => {
        setIsOpen(false)
        onBlur?.()
        inputRef.current?.focus()
      }, 200)
    }, [onBlur])

    const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
      if (e.key === 'Escape') {
        if (isOpen) {
          e.preventDefault()
          handleClose()
        }
      } else if (e.key === 'ArrowDown' && !isOpen) {
        e.preventDefault()
        handleOpen()
      } else if (e.key === 'Enter' && !isOpen) {
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

      return !!(isBeforeMin || isAfterMax || isInDisabledList)
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
            aria-label={t('common:select_month', {
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
            aria-label={t('common:select_year', {
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
          mode="range"
          selected={value}
          onSelect={handleDateSelect}
          disabled={disabledMatcher}
          month={displayMonth}
          onMonthChange={setDisplayMonth}
          numberOfMonths={monthCount}
          locale={locale}
          dir={isRTL ? 'rtl' : 'ltr'}
          className="kyora-datepicker w-full"
          classNames={{
            months: cn(
              'flex gap-4 w-full',
              isMobile ? 'flex-col' : 'flex-col sm:flex-row',
            ),
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
              'font-medium transition-colors',
              'h-10',
            ),
            day_button: cn(
              'w-full h-full p-0 font-normal rounded-lg',
              'hover:bg-base-200 transition-colors',
              'text-sm',
            ),
            selected: 'bg-primary text-primary-content hover:bg-primary/90',
            range_start:
              'bg-primary text-primary-content hover:bg-primary/90 rounded-s-lg',
            range_end:
              'bg-primary text-primary-content hover:bg-primary/90 rounded-e-lg',
            range_middle:
              'bg-primary/20 text-base-content hover:bg-primary/30 rounded-none',
            today: 'border-2 border-primary',
            outside: 'text-base-content/20 opacity-40',
            disabled:
              'text-base-content/20 opacity-40 cursor-not-allowed hover:bg-transparent',
            hidden: 'invisible',
          }}
        />
      </div>
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
            value={formatDateRange(value)}
            readOnly
            onClick={() => !disabled && handleOpen()}
            onKeyDown={handleKeyDown}
            disabled={disabled}
            required={required}
            placeholder={placeholder || t('common:select_date_range')}
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
              'focus:outline-none focus:border-primary focus:ring-2 focus:ring-primary/20',
              sizeClasses[size],
              'ps-10',
              clearable && value?.from && 'pe-10',
              error &&
                'input-error border-error focus:border-error focus:ring-error/20',
              disabled && 'opacity-50 cursor-not-allowed',
              isOpen && 'border-primary ring-2 ring-primary/20',
              className,
            )}
            {...props}
          />

          {clearable && value?.from && !disabled && (
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
                aria-label={label || t('common:select_date_range')}
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
                        {label || t('common:select_date_range')}
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
                    <div className="flex gap-3">
                      {clearable && (
                        <button
                          type="button"
                          onClick={() => {
                            onChange?.(undefined)
                            handleClose()
                          }}
                          className="btn btn-ghost flex-1"
                        >
                          {t('common:clear')}
                        </button>
                      )}
                      <button
                        type="button"
                        onClick={handleClose}
                        disabled={!value?.from || !value.to}
                        className="btn btn-primary flex-1"
                      >
                        {t('common:apply')}
                      </button>
                    </div>
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
              aria-label={label || t('common:select_date_range')}
              className={cn(
                'absolute top-full mt-2 z-50',
                'bg-base-100 rounded-box shadow-xl border border-base-300 p-4',
                'max-h-[70vh] overflow-auto',
                isRTL ? 'end-0' : 'start-0',
              )}
              style={{ width: dropdownWidth ? `${dropdownWidth}px` : '100%' }}
            >
              {calendarContent}
              <div className="flex gap-3 mt-4 pt-4 border-t border-base-300">
                {clearable && (
                  <button
                    type="button"
                    className="btn btn-sm btn-ghost flex-1"
                    onClick={() => {
                      onChange?.(undefined)
                      handleClose()
                    }}
                  >
                    {t('common:clear')}
                  </button>
                )}
                <button
                  type="button"
                  className="btn btn-sm btn-primary flex-1"
                  onClick={handleClose}
                  disabled={!value?.from || !value.to}
                >
                  {t('common:apply')}
                </button>
              </div>
            </div>
          ))}
      </div>
    )
  },
)

DateRangePicker.displayName = 'DateRangePicker'
