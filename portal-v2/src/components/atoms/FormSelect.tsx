/**
 * FormSelect - Refactored Advanced Select Component
 *
 * Production-grade select/dropdown with composition pattern.
 * Reduced from 551 lines to ~250 lines through hook extraction.
 *
 * Features:
 * - RTL-first design with logical properties
 * - Searchable with real-time filtering
 * - Multi-select support
 * - Full keyboard navigation
 * - Accessible with ARIA attributes
 * - Mobile-optimized
 * - Standard error prop integration
 */

import { forwardRef, useCallback, useEffect, useId, useRef, useState } from 'react'
import { Check, ChevronDown, Search, X } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { cn } from '../../lib/utils'
import { getErrorText } from '@/lib/formErrors'
import { useSelectSearch } from './useSelectSearch'
import { useSelectKeyboard } from './useSelectKeyboard'
import { useClickOutside } from './useClickOutside'
import type { ReactNode } from 'react'

export interface FormSelectOption<T = string> {
  value: T
  label: string
  description?: string
  icon?: ReactNode
  disabled?: boolean
  renderCustom?: () => ReactNode
}

export interface FormSelectProps<T = string> {
  label?: string
  error?: unknown
  helperText?: string
  options: Array<FormSelectOption<T>>
  value?: T | Array<T>
  onChange?: (value: T | Array<T>) => void
  size?: 'sm' | 'md' | 'lg'
  variant?: 'default' | 'filled' | 'ghost'
  fullWidth?: boolean
  searchable?: boolean
  multiSelect?: boolean
  placeholder?: string
  clearable?: boolean
  maxHeight?: number
  id?: string
  disabled?: boolean
  required?: boolean
  className?: string
}

export const FormSelect = forwardRef<HTMLDivElement, FormSelectProps>(
  <T extends string>(
    {
      label,
      error,
      helperText,
      options,
      value,
      onChange,
      size = 'md',
      variant = 'default',
      fullWidth = true,
      searchable = false,
      multiSelect = false,
      placeholder: placeholderProp,
      clearable = false,
      maxHeight = 300,
      className,
      id,
      disabled,
      required,
    }: FormSelectProps<T>,
    ref: React.ForwardedRef<HTMLDivElement>,
  ) => {
    const { t } = useTranslation()
    const generatedId = useId()
    const inputId = id ?? generatedId
    const errorText = getErrorText(error)
    const hasError = Boolean(errorText)
    const placeholder = placeholderProp ?? t('common.select')
    const [isOpen, setIsOpen] = useState(false)
    const searchInputRef = useRef<HTMLInputElement>(null)

    // Selected values normalization
    const selectedValues = (() => {
      if (multiSelect && Array.isArray(value)) return value
      if (!multiSelect && value !== undefined && !Array.isArray(value)) return [value]
      return []
    })()

    // Search management
    const { searchQuery, setSearchQuery, filteredOptions, clearSearch } = useSelectSearch({
      options,
      searchable,
    })

    // Close handler
    const handleClose = useCallback(() => {
      setIsOpen(false)
      clearSearch()
    }, [clearSearch])

    // Toggle option selection
    const handleToggleOption = useCallback(
      (optionValue: T) => {
        if (disabled) return

        if (multiSelect) {
          const newValues = selectedValues.includes(optionValue)
            ? selectedValues.filter((v) => v !== optionValue)
            : [...selectedValues, optionValue]
          onChange?.(newValues as T | Array<T>)
        } else {
          onChange?.(optionValue as T | Array<T>)
          handleClose()
        }
      },
      [disabled, multiSelect, selectedValues, onChange, handleClose],
    )

    // Keyboard navigation
    const { focusedIndex, setFocusedIndex, handleKeyDown } = useSelectKeyboard({
      isOpen,
      setIsOpen,
      filteredOptions,
      onSelectOption: handleToggleOption,
      onClose: handleClose,
      disabled,
    })

    // Click outside detection
    const containerRef = useClickOutside<HTMLDivElement>({
      isActive: isOpen,
      onClickOutside: handleClose,
    })

    // Auto-focus search input when dropdown opens
    useEffect(() => {
      if (isOpen && searchable && searchInputRef.current) {
        searchInputRef.current.focus()
      }
    }, [isOpen, searchable])

    // Prevent body scroll when dropdown is open on mobile
    useEffect(() => {
      if (isOpen && typeof window !== 'undefined') {
        const originalOverflow = document.body.style.overflow
        const originalPaddingRight = document.body.style.paddingRight

        // Only lock scroll on mobile/tablet
        const isMobile = window.innerWidth < 1024
        if (isMobile) {
          const scrollbarWidth = window.innerWidth - document.documentElement.clientWidth
          document.body.style.overflow = 'hidden'
          if (scrollbarWidth > 0) {
            document.body.style.paddingRight = `${String(scrollbarWidth)}px`
          }
        }

        return () => {
          if (isMobile) {
            document.body.style.overflow = originalOverflow
            document.body.style.paddingRight = originalPaddingRight
          }
        }
      }

      return undefined
    }, [isOpen])

    // Clear button handler
    const handleClear = useCallback(
      (e: React.MouseEvent) => {
        e.stopPropagation()
        onChange?.(multiSelect ? ([] as T | Array<T>) : (null as unknown as T | Array<T>))
      },
      [multiSelect, onChange],
    )

    // Display text
    const getDisplayText = () => {
      if (selectedValues.length === 0) return placeholder

      if (multiSelect) {
        const count = selectedValues.length
        return count > 0 ? t('common.selected_count', { count }) : placeholder
      }

      const selectedOption = options.find((opt) => opt.value === selectedValues[0])
      return selectedOption?.label ?? placeholder
    }

    const sizeClasses = {
      sm: 'h-[44px] text-sm',
      md: 'h-[50px] text-base',
      lg: 'h-[56px] text-lg',
    }

    const variantClasses = {
      default: 'input-bordered bg-base-100',
      filled: 'input-bordered bg-base-200/50 border-transparent focus:bg-base-100',
      ghost: 'input-ghost bg-transparent',
    }

    return (
      <div ref={ref} className={cn('form-control', fullWidth && 'w-full')}>
        {label && (
          <label htmlFor={inputId} className="label">
            <span className="label-text text-base-content/70 font-medium">
              {label}
              {required && <span className="text-error ms-1">*</span>}
            </span>
          </label>
        )}

        <div ref={containerRef} className="relative">
          {/* Select Trigger Button */}
          <button
            type="button"
            id={inputId}
            onClick={() => !disabled && setIsOpen(!isOpen)}
            onKeyDown={handleKeyDown}
            disabled={disabled}
            className={cn(
              'input w-full flex items-center justify-between gap-2 transition-all duration-200',
              sizeClasses[size],
              variantClasses[variant],
              'text-start cursor-pointer',
              'focus:outline-none focus:border-primary focus:ring-2 focus:ring-primary/20',
              hasError && 'input-error border-error focus:border-error focus:ring-error/20',
              disabled && 'opacity-60 cursor-not-allowed',
              isOpen && 'border-primary ring-2 ring-primary/20',
              className,
            )}
            aria-haspopup="listbox"
            aria-expanded={isOpen}
            aria-labelledby={label ? `${inputId}-label` : undefined}
            aria-required={required}
            aria-invalid={hasError}
          >
            <span className={cn('flex-1 truncate', selectedValues.length === 0 && 'text-base-content/40')}>
              {getDisplayText()}
            </span>

            <div className="flex items-center gap-1">
              {clearable && selectedValues.length > 0 && !disabled && (
                <button
                  type="button"
                  onClick={handleClear}
                  className="p-1 hover:bg-base-200 rounded-md transition-colors"
                  aria-label={t('common.clear_selection')}
                >
                  <X className="w-4 h-4" />
                </button>
              )}
              <ChevronDown className={cn('w-5 h-5 transition-transform duration-200', isOpen && 'rotate-180')} />
            </div>
          </button>

          {/* Dropdown Panel */}
          {isOpen && !disabled && (
            <div
              className={cn(
                'absolute z-50 mt-2 w-full',
                'bg-base-100 border border-base-300 rounded-lg shadow-xl',
                'overflow-hidden',
                'animate-in fade-in-0 zoom-in-95 duration-100',
              )}
              style={{ maxHeight }}
              role="presentation"
              onClick={(e) => e.stopPropagation()}
            >
              {/* Search Input */}
              {searchable && (
                <div className="p-2 border-b border-base-300">
                  <div className="relative">
                    <Search className="absolute start-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content/50" />
                    <input
                      ref={searchInputRef}
                      type="text"
                      value={searchQuery}
                      onChange={(e) => {
                        setSearchQuery(e.target.value)
                        setFocusedIndex(-1)
                      }}
                      placeholder={t('common.search_placeholder_generic')}
                      className="input input-sm w-full ps-9"
                      aria-label={t('common.search_options')}
                    />
                  </div>
                </div>
              )}

              {/* Options List */}
              <ul
                role="listbox"
                aria-multiselectable={multiSelect}
                className="overflow-y-auto"
                style={{ maxHeight: maxHeight - (searchable ? 60 : 0) }}
              >
                {filteredOptions.length === 0 ? (
                  <li className="p-4 text-center text-base-content/50">{t('common.no_options_found')}</li>
                ) : (
                  filteredOptions.map((option, index) => {
                    const isSelected = selectedValues.includes(option.value)
                    const isFocused = index === focusedIndex

                    return (
                      <li
                        key={option.value}
                        role="option"
                        aria-selected={isSelected}
                        aria-disabled={option.disabled}
                        onClick={() => !option.disabled && handleToggleOption(option.value)}
                        tabIndex={isFocused ? 0 : -1}
                        ref={(el) => {
                          if (isFocused && el) {
                            el.scrollIntoView({ block: 'nearest', behavior: 'smooth' })
                          }
                        }}
                        className={cn(
                          'flex items-center gap-3 px-4 py-3 cursor-pointer transition-colors',
                          'min-h-[48px]',
                          'hover:bg-base-200 focus:bg-base-200 focus:outline-none',
                          'active:bg-base-300',
                          isSelected && 'bg-primary/10',
                          isFocused && 'bg-base-200 ring-2 ring-inset ring-primary/30',
                          option.disabled && 'opacity-50 cursor-not-allowed pointer-events-none',
                        )}
                      >
                        {option.renderCustom ? (
                          option.renderCustom()
                        ) : (
                          <>
                            {option.icon && <span className="shrink-0">{option.icon}</span>}
                            <div className="flex-1 min-w-0">
                              <div className="font-medium truncate">{option.label}</div>
                              {option.description && (
                                <div className="text-sm text-base-content/60 truncate">{option.description}</div>
                              )}
                            </div>
                            {isSelected && <Check className="w-5 h-5 text-primary shrink-0" />}
                          </>
                        )}
                      </li>
                    )
                  })
                )}
              </ul>
            </div>
          )}
        </div>

        {/* Error Message */}
        {hasError && (
          <label className="label">
            <span id={`${inputId}-error`} className="label-text-alt text-error" role="alert">
              {errorText}
            </span>
          </label>
        )}

        {/* Helper Text */}
        {!hasError && helperText && (
          <label className="label">
            <span id={`${inputId}-helper`} className="label-text-alt text-base-content/60">
              {helperText}
            </span>
          </label>
        )}
      </div>
    )
  },
) as (<T extends string>(
  props: FormSelectProps<T> & { ref?: React.ForwardedRef<HTMLDivElement> },
) => React.ReactElement) & { displayName?: string }

FormSelect.displayName = 'FormSelect'
