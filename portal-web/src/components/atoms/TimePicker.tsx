/**
 * TimePicker - Production-Grade Time Selection Component
 *
 * Mobile-first time picker with bottom sheet on mobile, dropdown on desktop.
 * Fully aligned with Kyora Design System (KDS).
 *
 * Features:
 * - Mobile-first: Bottom sheet (<768px), dropdown (≥768px)
 * - 12/24-hour format support
 * - Touch-optimized dropdown selectors for hours/minutes
 * - Keyboard accessible
 * - daisyUI styled: Consistent with form system
 * - AM/PM toggle for 12-hour format
 * - Clear button: Optional reset functionality
 *
 * @example
 * <TimePicker
 *   label="Appointment Time"
 *   value="14:30"
 *   onChange={setTime}
 *   format={24}
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
import { Clock, X } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import type { InputHTMLAttributes } from 'react'
import { cn } from '@/lib/utils'
import { useLanguage } from '@/hooks/useLanguage'
import { useMediaQuery } from '@/hooks/useMediaQuery'

export interface TimePickerProps extends Omit<
  InputHTMLAttributes<HTMLInputElement>,
  'value' | 'onChange' | 'size'
> {
  label?: string
  value?: string
  onChange?: (time: string | undefined) => void
  onBlur?: () => void
  error?: string
  helperText?: string
  format?: 12 | 24
  minuteStep?: number
  clearable?: boolean
  fullWidth?: boolean
  size?: 'sm' | 'md' | 'lg'
}

export const TimePicker = forwardRef<HTMLInputElement, TimePickerProps>(
  (
    {
      label,
      value,
      onChange,
      onBlur,
      error,
      helperText,
      format = 24,
      minuteStep = 1,
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
    const { t } = useTranslation('common')
    const [isOpen, setIsOpen] = useState(false)
    const [isAnimating, setIsAnimating] = useState(false)
    const [dropdownWidth, setDropdownWidth] = useState<number | undefined>()
    const [selectedHour, setSelectedHour] = useState<number>(12)
    const [selectedMinute, setSelectedMinute] = useState<number>(0)
    const [selectedPeriod, setSelectedPeriod] = useState<'AM' | 'PM'>('AM')
    const containerRef = useRef<HTMLDivElement>(null)
    const inputRef = useRef<HTMLInputElement>(null)
    const popupRef = useRef<HTMLDivElement>(null)
    const isMobile = useMediaQuery('(max-width: 768px)')

    const sizeClasses = {
      sm: 'h-[44px] text-sm',
      md: 'h-[50px] text-base',
      lg: 'h-[56px] text-lg',
    }

    const hours =
      format === 12
        ? Array.from({ length: 12 }, (_, i) => i + 1)
        : Array.from({ length: 24 }, (_, i) => i)
    const minutes = Array.from(
      { length: 60 / minuteStep },
      (_, i) => i * minuteStep,
    )

    useEffect(() => {
      if (value) {
        const [hourStr, minuteStr] = value.split(':')
        const hour = parseInt(hourStr, 10)
        const minute = parseInt(minuteStr, 10)

        if (format === 12) {
          if (hour === 0) {
            setSelectedHour(12)
            setSelectedPeriod('AM')
          } else if (hour === 12) {
            setSelectedHour(12)
            setSelectedPeriod('PM')
          } else if (hour > 12) {
            setSelectedHour(hour - 12)
            setSelectedPeriod('PM')
          } else {
            setSelectedHour(hour)
            setSelectedPeriod('AM')
          }
        } else {
          setSelectedHour(hour)
        }
        setSelectedMinute(minute)
      }
    }, [value, format])

    const formatTime = (
      hour: number,
      minute: number,
      period?: 'AM' | 'PM',
    ): string => {
      if (format === 12 && period) {
        return `${hour.toString().padStart(2, '0')}:${minute.toString().padStart(2, '0')} ${period}`
      }
      return `${hour.toString().padStart(2, '0')}:${minute.toString().padStart(2, '0')}`
    }

    const formatTimeFor24 = (
      hour: number,
      minute: number,
      period: 'AM' | 'PM',
    ): string => {
      let hour24 = hour
      if (format === 12) {
        if (period === 'AM' && hour === 12) {
          hour24 = 0
        } else if (period === 'PM' && hour !== 12) {
          hour24 = hour + 12
        }
      }
      return `${hour24.toString().padStart(2, '0')}:${minute.toString().padStart(2, '0')}`
    }

    const handleApply = () => {
      const time24 = formatTimeFor24(
        selectedHour,
        selectedMinute,
        selectedPeriod,
      )
      onChange?.(time24)
      handleClose()
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

    const timePickerContent = (
      <div className="flex gap-2 justify-center items-center">
        <div className="flex flex-col items-center gap-2 min-w-[80px]">
          <label className="text-sm font-medium text-base-content/70">
            {t('date.hours')}
          </label>
          <select
            value={selectedHour}
            onChange={(e) => setSelectedHour(Number(e.target.value))}
            className="select select-bordered w-full"
          >
            {hours.map((hour) => (
              <option key={hour} value={hour}>
                {hour.toString().padStart(2, '0')}
              </option>
            ))}
          </select>
        </div>

        <span className="text-2xl font-bold text-base-content mt-8">:</span>

        <div className="flex flex-col items-center gap-2 min-w-[80px]">
          <label className="text-sm font-medium text-base-content/70">
            {t('date.minutes')}
          </label>
          <select
            value={selectedMinute}
            onChange={(e) => setSelectedMinute(Number(e.target.value))}
            className="select select-bordered w-full"
          >
            {minutes.map((minute) => (
              <option key={minute} value={minute}>
                {minute.toString().padStart(2, '0')}
              </option>
            ))}
          </select>
        </div>

        {format === 12 && (
          <div className="flex flex-col items-center gap-2 min-w-[80px]">
            <label className="text-sm font-medium text-base-content/70">
              {t('date.period')}
            </label>
            <div className="join">
              <button
                type="button"
                onClick={() => setSelectedPeriod('AM')}
                className={cn(
                  'btn join-item',
                  selectedPeriod === 'AM' ? 'btn-primary' : 'btn-ghost',
                )}
              >
                AM
              </button>
              <button
                type="button"
                onClick={() => setSelectedPeriod('PM')}
                className={cn(
                  'btn join-item',
                  selectedPeriod === 'PM' ? 'btn-primary' : 'btn-ghost',
                )}
              >
                PM
              </button>
            </div>
          </div>
        )}
      </div>
    )

    const displayValue = value
      ? format === 12
        ? formatTime(selectedHour, selectedMinute, selectedPeriod)
        : value
      : ''

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
            <Clock size={20} aria-hidden="true" />
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
            value={displayValue}
            readOnly
            onClick={() => !disabled && handleOpen()}
            onKeyDown={handleKeyDown}
            disabled={disabled}
            required={required}
            placeholder={
              placeholder || (format === 12 ? 'hh:mm AM/PM' : 'HH:mm')
            }
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
                aria-label={label || t('date.selectTime')}
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
                    'w-full',
                    'transition-transform duration-200',
                    isAnimating ? 'translate-y-0' : 'translate-y-full',
                  )}
                >
                  <div className="sticky top-0 z-10 bg-base-100 border-b border-base-300 px-4 py-3">
                    <div className="flex justify-between items-center">
                      <h3 className="text-lg font-semibold">
                        {label || t('date.selectTime')}
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

                  <div className="p-4">{timePickerContent}</div>

                  <div className="sticky bottom-0 z-10 bg-base-100 border-t border-base-300 px-4 py-3">
                    <div className="flex gap-2">
                      <button
                        type="button"
                        onClick={() => {
                          onChange?.(undefined)
                          handleClose()
                        }}
                        className="btn btn-ghost flex-1"
                      >
                        {t('clear')}
                      </button>
                      <button
                        type="button"
                        onClick={handleApply}
                        className="btn btn-primary flex-1"
                      >
                        {t('apply')}
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
              aria-label={label || t('date.selectTime')}
              className={cn(
                'absolute top-full mt-2 z-50',
                'bg-base-100 rounded-box shadow-lg border border-base-300 p-4',
                isRTL ? 'end-0' : 'start-0',
                'min-w-[320px]',
              )}
              style={{
                width: dropdownWidth ? `${dropdownWidth}px` : undefined,
              }}
            >
              {timePickerContent}
              <div className="flex gap-2 mt-4 pt-4 border-t border-base-300">
                <button
                  type="button"
                  className="btn btn-sm btn-ghost flex-1"
                  onClick={() => {
                    onChange?.(undefined)
                    handleClose()
                  }}
                >
                  {t('clear')}
                </button>
                <button
                  type="button"
                  className="btn btn-sm btn-primary flex-1"
                  onClick={handleApply}
                >
                  {t('apply')}
                </button>
              </div>
            </div>
          ))}
      </div>
    )
  },
)

TimePicker.displayName = 'TimePicker'
