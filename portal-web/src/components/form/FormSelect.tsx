/**
 * FormSelect - Mobile-First Advanced Select Component
 *
 * Production-grade select/dropdown optimized for mobile with bottom sheet on small screens.
 *
 * Features:
 * - Mobile-first: Bottom sheet on mobile (< 768px), dropdown on desktop
 * - RTL-first design with logical properties
 * - Searchable with real-time filtering
 * - Multi-select support
 * - Full keyboard navigation
 * - Accessible with ARIA attributes (WCAG AA compliant)
 * - Optimized for use inside modals/bottom sheets (portal rendering)
 * - Smart scroll lock management (respects parent modals)
 * - Touch-optimized interactions
 * - Smooth animations and transitions
 * - Standard error prop integration
 */

import {
  forwardRef,
  useCallback,
  useEffect,
  useId,
  useMemo,
  useRef,
  useState,
} from 'react'
import { createPortal } from 'react-dom'
import { Check, ChevronDown, Plus, Search, X } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { cn } from '../../lib/utils'
import { useMediaQuery } from '../../hooks/useMediaQuery'
import { useSelectSearch } from './useSelectSearch'
import { useSelectKeyboard } from './useSelectKeyboard'
import { useClickOutside } from './useClickOutside'
import { FormInput } from './FormInput'
import type { ReactNode } from 'react'
import { getErrorText } from '@/lib/formErrors'

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
  /** Controlled search value (for external control) */
  searchValue?: string
  /** Callback when search value changes (for external control) */
  onSearchChange?: (value: string) => void
  /** Callback to create an option from search input */
  onCreateOption?: (label: string) => void | Promise<void>
  /** Custom label for the create option action */
  createOptionLabel?: (label: string) => string
  /** Loading state while creating a new option */
  isCreatingOption?: boolean
  /** Loading state for options (e.g., backend search in progress) */
  isLoading?: boolean
  id?: string
  disabled?: boolean
  required?: boolean
  className?: string
  /** Title for mobile bottom sheet (defaults to label) */
  mobileTitle?: string
  /** Callback when dropdown opens */
  onOpen?: () => void
  /** Callback when dropdown closes */
  onClose?: () => void
  /** Callback when clear button is clicked */
  onClear?: () => void
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
      mobileTitle,
      onOpen,
      onClose,
      onClear,
      searchValue: externalSearchValue,
      onSearchChange: externalOnSearchChange,
      onCreateOption,
      createOptionLabel,
      isCreatingOption,
      isLoading,
    }: FormSelectProps<T>,
    ref: React.ForwardedRef<HTMLDivElement>,
  ) => {
    const { t: tCommon } = useTranslation('common')
    const generatedId = useId()
    const inputId = id ?? generatedId
    const errorText = getErrorText(error)
    const hasError = Boolean(errorText)
    const placeholder = placeholderProp ?? tCommon('select')
    const [isOpen, setIsOpen] = useState(false)
    const [isAnimating, setIsAnimating] = useState(false)
    const searchInputRef = useRef<HTMLInputElement>(null)
    const mobileContentRef = useRef<HTMLDivElement>(null)
    const triggerRef = useRef<HTMLButtonElement>(null)
    const isMobile = useMediaQuery('(max-width: 768px)')

    // Selected values normalization
    const selectedValues = (() => {
      if (multiSelect && Array.isArray(value)) return value
      if (!multiSelect && value !== undefined && !Array.isArray(value))
        return [value]
      return []
    })()

    // Determine if using external (backend) search
    const isExternalSearch = Boolean(externalOnSearchChange)

    // Search management (only used for internal/client-side filtering)
    const {
      searchQuery: internalSearchQuery,
      setSearchQuery: setInternalSearchQuery,
      filteredOptions: clientFilteredOptions,
      clearSearch,
    } = useSelectSearch({
      options,
      // Only enable client-side filtering if NOT using external search
      searchable: searchable && !isExternalSearch,
    })

    // Use external search if provided, otherwise use internal
    const searchQuery = externalSearchValue ?? internalSearchQuery
    const setSearchQuery = externalOnSearchChange ?? setInternalSearchQuery

    // For external search, use options directly (already filtered by backend)
    // For internal search, use client-side filtered options
    const filteredOptions = isExternalSearch ? options : clientFilteredOptions

    const normalizedSearchQuery = searchQuery.trim()
    const hasExactOption = normalizedSearchQuery
      ? options.some(
          (opt) =>
            opt.label.trim().toLowerCase() ===
            normalizedSearchQuery.toLowerCase(),
        )
      : false
    const canCreateOption = Boolean(
      onCreateOption &&
      normalizedSearchQuery.length > 0 &&
      !hasExactOption &&
      !disabled,
    )

    const createActionLabel = createOptionLabel
      ? createOptionLabel(normalizedSearchQuery)
      : tCommon('create_option', { option: normalizedSearchQuery })

    // Close handler with animation
    const handleClose = useCallback(() => {
      setIsAnimating(false)

      // Wait for animation before unmounting
      setTimeout(() => {
        setIsOpen(false)
        clearSearch()
        onClose?.()

        // Restore focus to trigger
        triggerRef.current?.focus()
      }, 200)
    }, [clearSearch, onClose])

    // Open handler
    const handleOpen = useCallback(() => {
      setIsOpen(true)
      setIsAnimating(true)
      onOpen?.()
    }, [onOpen])

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

    // Remove last chip (for Backspace key on trigger)
    const handleRemoveLast = useCallback(() => {
      if (multiSelect && selectedValues.length > 0) {
        const newValues = selectedValues.slice(0, -1)
        onChange?.(newValues as T | Array<T>)
      }
    }, [multiSelect, selectedValues, onChange])

    // Keyboard navigation
    const { focusedIndex, setFocusedIndex, handleKeyDown } = useSelectKeyboard({
      isOpen,
      setIsOpen,
      filteredOptions,
      onSelectOption: handleToggleOption,
      onClose: handleClose,
      disabled,
      multiSelect,
      selectedValues,
      onRemoveLast: handleRemoveLast,
    })

    // Click outside detection (desktop only)
    const containerRef = useClickOutside<HTMLDivElement>({
      isActive: isOpen && !isMobile,
      onClickOutside: handleClose,
    })

    // Auto-focus search input when dropdown opens
    useEffect(() => {
      if (isOpen && searchable && searchInputRef.current) {
        // Add delay for mobile to allow animation to complete
        const delay = isMobile ? 300 : 50
        const timer = setTimeout(() => {
          searchInputRef.current?.focus()
        }, delay)
        return () => clearTimeout(timer)
      }
    }, [isOpen, searchable, isMobile])

    // Smart scroll lock for mobile (only if not already locked by parent)
    useEffect(() => {
      if (isOpen && isMobile && typeof window !== 'undefined') {
        // Check if body scroll is already locked by a parent modal
        const isAlreadyLocked = document.body.style.overflow === 'hidden'

        if (!isAlreadyLocked) {
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
      }
    }, [isOpen, isMobile])

    // Clear button handler
    const handleClear = useCallback(
      (e: React.MouseEvent) => {
        e.stopPropagation()
        if (onClear) {
          onClear()
        } else {
          onChange?.(
            multiSelect
              ? ([] as T | Array<T>)
              : (null as unknown as T | Array<T>),
          )
        }
      },
      [multiSelect, onChange, onClear],
    )

    // Remove chip handler
    const handleRemoveChip = useCallback(
      (valueToRemove: T, e?: React.MouseEvent | React.KeyboardEvent) => {
        e?.stopPropagation()
        if (disabled) return

        const newValues = selectedValues.filter((v) => v !== valueToRemove)
        onChange?.(newValues as T | Array<T>)
      },
      [disabled, selectedValues, onChange],
    )

    // Keyboard handler for chips
    const handleChipKeyDown = useCallback(
      (valueToRemove: T, e: React.KeyboardEvent) => {
        if (e.key === 'Backspace' || e.key === 'Delete') {
          e.preventDefault()
          handleRemoveChip(valueToRemove, e)
        }
      },
      [handleRemoveChip],
    )

    const handleCreateOption = useCallback(async () => {
      if (!canCreateOption || !onCreateOption || isCreatingOption) return

      await onCreateOption(normalizedSearchQuery)
      setSearchQuery('')
      setFocusedIndex(-1)
      if (!multiSelect) {
        handleClose()
      }
    }, [
      canCreateOption,
      onCreateOption,
      isCreatingOption,
      normalizedSearchQuery,
      setSearchQuery,
      setFocusedIndex,
      multiSelect,
      handleClose,
    ])

    // Display text (for single-select trigger)
    const getDisplayText = () => {
      if (selectedValues.length === 0) return placeholder

      if (!multiSelect) {
        const selectedOption = options.find(
          (opt) => opt.value === selectedValues[0],
        )
        return selectedOption?.label ?? placeholder
      }

      // Multi-select: return placeholder to be replaced by chips
      return placeholder
    }

    // Get selected option objects for chip display
    const selectedOptions = useMemo(() => {
      return selectedValues
        .map((val) => options.find((opt) => opt.value === val))
        .filter((opt): opt is FormSelectOption<T> => opt !== undefined)
    }, [selectedValues, options])

    const sizeClasses = {
      sm: 'h-[44px] text-sm',
      md: 'h-[50px] text-base',
      lg: 'h-[56px] text-lg',
    }

    const variantClasses = {
      default: 'input-bordered bg-base-100',
      filled:
        'input-bordered bg-base-200/50 border-transparent focus:bg-base-100',
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
            ref={triggerRef}
            type="button"
            id={inputId}
            onClick={() => {
              if (disabled) return
              if (isOpen) {
                handleClose()
              } else {
                handleOpen()
              }
            }}
            onKeyDown={handleKeyDown}
            disabled={disabled}
            className={cn(
              'input w-full flex items-center justify-between gap-2 transition-all duration-200',
              sizeClasses[size],
              variantClasses[variant],
              'text-start cursor-pointer',
              'focus:outline-none focus:border-primary focus:ring-2 focus:ring-primary/20',
              hasError &&
                'input-error border-error focus:border-error focus:ring-error/20',
              disabled && 'opacity-60 cursor-not-allowed',
              isOpen && 'border-primary ring-2 ring-primary/20',
              multiSelect && selectedValues.length > 0 && 'min-h-fit py-2',
              className,
            )}
            aria-haspopup="listbox"
            aria-expanded={isOpen}
            aria-labelledby={label ? `${inputId}-label` : undefined}
            aria-required={required}
            aria-invalid={hasError}
          >
            {multiSelect && selectedValues.length > 0 ? (
              // Multi-select: Show chips/badges
              <div className="flex-1 flex flex-wrap gap-1.5">
                {selectedOptions.map((option) => (
                  <span
                    key={option.value}
                    className={cn(
                      'badge badge-neutral gap-1.5 pe-1 ps-2.5 h-7 text-sm font-medium',
                      'transition-colors duration-200',
                    )}
                    tabIndex={0}
                    role="button"
                    aria-label={tCommon('remove_option', {
                      option: option.label,
                    })}
                    onKeyDown={(e) => handleChipKeyDown(option.value, e)}
                  >
                    {option.icon && (
                      <span className="shrink-0">{option.icon}</span>
                    )}
                    <span className="truncate max-w-[150px]">
                      {option.label}
                    </span>
                    <button
                      type="button"
                      onClick={(e) => handleRemoveChip(option.value, e)}
                      className={cn(
                        'btn btn-ghost btn-circle btn-xs',
                        'hover:bg-base-content/10 transition-colors',
                      )}
                      aria-label={tCommon('remove', { item: option.label })}
                      tabIndex={-1}
                    >
                      <X className="w-3.5 h-3.5" />
                    </button>
                  </span>
                ))}
              </div>
            ) : (
              // Single-select or empty: Show text
              <span
                className={cn(
                  'flex-1 truncate',
                  selectedValues.length === 0 && 'text-base-content/40',
                )}
              >
                {getDisplayText()}
              </span>
            )}

            <div className="flex items-center gap-1 shrink-0">
              {clearable &&
                selectedValues.length > 0 &&
                !disabled &&
                // Only show clear button if there's actually a selected value (not empty string)
                selectedValues.some((v) => v !== '') && (
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
                    className="p-1 hover:bg-base-200 rounded-md transition-colors cursor-pointer"
                    aria-label={tCommon('clear_selection')}
                  >
                    <X className="w-4 h-4" />
                  </span>
                )}
              <ChevronDown
                className={cn(
                  'w-5 h-5 transition-transform duration-200',
                  isOpen && 'rotate-180',
                )}
              />
            </div>
          </button>

          {/* Dropdown Panel - Mobile (Bottom Sheet) or Desktop (Dropdown) */}
          {isOpen &&
            !disabled &&
            (isMobile ? (
              // Mobile: Bottom Sheet with Portal
              createPortal(
                <div
                  className={cn(
                    'fixed inset-0 z-[9999] flex items-end justify-center',
                    'transition-opacity duration-200',
                    isAnimating ? 'opacity-100' : 'opacity-0',
                  )}
                  role="dialog"
                  aria-modal="true"
                  aria-label={mobileTitle || label || tCommon('select')}
                  onClick={(e) => {
                    // Close on backdrop click
                    if (e.target === e.currentTarget) {
                      handleClose()
                    }
                  }}
                >
                  {/* Backdrop */}
                  <div
                    className="absolute inset-0 bg-base-content/50 backdrop-blur-sm"
                    aria-hidden="true"
                  />

                  {/* Bottom Sheet Content */}
                  <div
                    ref={mobileContentRef}
                    className={cn(
                      'relative w-full max-h-[85vh]',
                      'bg-base-100 rounded-t-xl',
                      'overflow-hidden flex flex-col',
                      'transition-transform duration-200 ease-out',
                      isAnimating ? 'translate-y-0' : 'translate-y-full',
                    )}
                    onClick={(e) => e.stopPropagation()}
                  >
                    {/* Drag Handle */}
                    <div className="flex justify-center py-2 px-4 shrink-0">
                      <div
                        className="w-12 h-1 bg-base-300 rounded-full"
                        aria-hidden="true"
                      />
                    </div>

                    {/* Header */}
                    <div className="flex items-center justify-between px-4 pb-3 shrink-0 border-b border-base-300">
                      <h2 className="text-lg font-semibold text-base-content">
                        {mobileTitle || label || tCommon('select')}
                      </h2>
                      <button
                        type="button"
                        onClick={handleClose}
                        className="btn btn-ghost btn-sm btn-square"
                        aria-label={tCommon('close')}
                      >
                        <X className="w-5 h-5" />
                      </button>
                    </div>

                    {/* Search Input */}
                    {searchable && (
                      <div className="px-4 pt-4 pb-2 border-b border-base-300 shrink-0">
                        <FormInput
                          ref={searchInputRef}
                          type="text"
                          value={searchQuery}
                          onChange={(e) => {
                            setSearchQuery(e.target.value)
                            setFocusedIndex(-1)
                          }}
                          onKeyDown={(e) => {
                            if (e.key === 'Enter' && canCreateOption) {
                              e.preventDefault()
                              handleCreateOption()
                            }
                          }}
                          placeholder={tCommon('search_placeholder_generic')}
                          startIcon={<Search className="w-5 h-5" />}
                          size="lg"
                          variant="filled"
                          fullWidth
                          aria-label={tCommon('search_options')}
                        />
                      </div>
                    )}

                    {/* Options List */}
                    <ul
                      role="listbox"
                      aria-multiselectable={multiSelect}
                      className="overflow-y-auto flex-1 overscroll-contain"
                    >
                      {canCreateOption && (
                        <li
                          role="option"
                          aria-selected={false}
                          onClick={handleCreateOption}
                          className={cn(
                            'flex items-center gap-3 px-4 py-4 cursor-pointer transition-colors',
                            'min-h-[56px] border-b border-base-200',
                            'hover:bg-base-200 active:bg-base-300',
                          )}
                        >
                          <div className="p-2 rounded-lg bg-primary/10 text-primary">
                            <Plus className="w-5 h-5" />
                          </div>
                          <div className="flex-1 min-w-0">
                            <div className="font-medium text-base">
                              {createActionLabel}
                            </div>
                            <div className="text-sm text-base-content/60 mt-0.5">
                              {tCommon('create')}
                            </div>
                          </div>
                          {isCreatingOption && (
                            <span className="loading loading-spinner loading-sm text-primary" />
                          )}
                        </li>
                      )}

                      {isLoading ? (
                        <li className="p-6 flex justify-center">
                          <span className="loading loading-spinner loading-md text-primary" />
                        </li>
                      ) : filteredOptions.length === 0 && !canCreateOption ? (
                        <li className="p-6 text-center text-base-content/50">
                          {tCommon('no_options_found')}
                        </li>
                      ) : (
                        filteredOptions.map((option) => {
                          const isSelected = selectedValues.includes(
                            option.value,
                          )

                          return (
                            <li
                              key={option.value}
                              role="option"
                              aria-selected={isSelected}
                              aria-disabled={option.disabled}
                              onClick={() => {
                                if (!option.disabled) {
                                  handleToggleOption(option.value)
                                  // Close immediately on single select
                                  if (!multiSelect) {
                                    handleClose()
                                  }
                                }
                              }}
                              className={cn(
                                'flex items-center gap-3 px-4 py-4 cursor-pointer transition-colors',
                                'min-h-[56px]',
                                'active:bg-base-300',
                                isSelected && 'bg-primary/10',
                                option.disabled &&
                                  'opacity-50 cursor-not-allowed pointer-events-none',
                              )}
                            >
                              {option.renderCustom ? (
                                option.renderCustom()
                              ) : (
                                <>
                                  {option.icon && (
                                    <span className="shrink-0 text-xl">
                                      {option.icon}
                                    </span>
                                  )}
                                  <div className="flex-1 min-w-0">
                                    <div className="font-medium text-base">
                                      {option.label}
                                    </div>
                                    {option.description && (
                                      <div className="text-sm text-base-content/60 mt-0.5">
                                        {option.description}
                                      </div>
                                    )}
                                  </div>
                                  {isSelected && (
                                    <Check className="w-6 h-6 text-primary shrink-0" />
                                  )}
                                </>
                              )}
                            </li>
                          )
                        })
                      )}
                    </ul>
                  </div>
                </div>,
                document.body,
              )
            ) : (
              // Desktop: Dropdown
              <div
                className={cn(
                  'absolute z-50 mt-2 w-full',
                  'bg-base-100 border border-base-300 rounded-lg',
                  'overflow-hidden',
                  'transition-all duration-200 origin-top',
                  isAnimating ? 'opacity-100 scale-100' : 'opacity-0 scale-95',
                )}
                style={{ maxHeight }}
                role="presentation"
                onClick={(e) => e.stopPropagation()}
              >
                {/* Search Input */}
                {searchable && (
                  <div className="p-2 border-b border-base-300">
                    <FormInput
                      ref={searchInputRef}
                      type="text"
                      value={searchQuery}
                      onChange={(e) => {
                        setSearchQuery(e.target.value)
                        setFocusedIndex(-1)
                      }}
                      onKeyDown={(e) => {
                        if (e.key === 'Enter' && canCreateOption) {
                          e.preventDefault()
                          handleCreateOption()
                        }
                      }}
                      placeholder={tCommon('search_placeholder_generic')}
                      startIcon={<Search className="w-4 h-4" />}
                      size="sm"
                      variant="filled"
                      fullWidth
                      aria-label={tCommon('search_options')}
                    />
                  </div>
                )}

                {/* Options List */}
                <ul
                  role="listbox"
                  aria-multiselectable={multiSelect}
                  className="overflow-y-auto"
                  style={{ maxHeight: maxHeight - (searchable ? 60 : 0) }}
                >
                  {canCreateOption && (
                    <li
                      role="option"
                      aria-selected={false}
                      onClick={handleCreateOption}
                      className={cn(
                        'flex items-center gap-3 px-4 py-3 cursor-pointer transition-colors',
                        'min-h-[48px] border-b border-base-200',
                        'hover:bg-base-200 focus:bg-base-200 focus:outline-none',
                      )}
                    >
                      <div className="p-2 rounded-lg bg-primary/10 text-primary">
                        <Plus className="w-4 h-4" />
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="font-medium truncate">
                          {createActionLabel}
                        </div>
                        <div className="text-sm text-base-content/60 truncate">
                          {tCommon('create')}
                        </div>
                      </div>
                      {isCreatingOption && (
                        <span className="loading loading-spinner loading-xs text-primary" />
                      )}
                    </li>
                  )}

                  {isLoading ? (
                    <li className="p-4 flex justify-center">
                      <span className="loading loading-spinner loading-sm text-primary" />
                    </li>
                  ) : filteredOptions.length === 0 && !canCreateOption ? (
                    <li className="p-4 text-center text-base-content/50">
                      {tCommon('no_options_found')}
                    </li>
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
                          onClick={() =>
                            !option.disabled && handleToggleOption(option.value)
                          }
                          tabIndex={isFocused ? 0 : -1}
                          ref={(el) => {
                            if (isFocused && el) {
                              el.scrollIntoView({
                                block: 'nearest',
                                behavior: 'smooth',
                              })
                            }
                          }}
                          className={cn(
                            'flex items-center gap-3 px-4 py-3 cursor-pointer transition-colors',
                            'min-h-[48px]',
                            'hover:bg-base-200 focus:bg-base-200 focus:outline-none',
                            'active:bg-base-300',
                            isSelected && 'bg-primary/10',
                            isFocused &&
                              'bg-base-200 ring-2 ring-inset ring-primary/30',
                            option.disabled &&
                              'opacity-50 cursor-not-allowed pointer-events-none',
                          )}
                        >
                          {option.renderCustom ? (
                            option.renderCustom()
                          ) : (
                            <>
                              {option.icon && (
                                <span className="shrink-0">{option.icon}</span>
                              )}
                              <div className="flex-1 min-w-0">
                                <div className="font-medium truncate">
                                  {option.label}
                                </div>
                                {option.description && (
                                  <div className="text-sm text-base-content/60 truncate">
                                    {option.description}
                                  </div>
                                )}
                              </div>
                              {isSelected && (
                                <Check className="w-5 h-5 text-primary shrink-0" />
                              )}
                            </>
                          )}
                        </li>
                      )
                    })
                  )}
                </ul>
              </div>
            ))}
        </div>

        {/* Error Message */}
        {hasError && (
          <label className="label">
            <span
              id={`${inputId}-error`}
              className="label-text-alt text-error"
              role="alert"
            >
              {errorText}
            </span>
          </label>
        )}

        {/* Helper Text */}
        {!hasError && helperText && (
          <label className="label">
            <span
              id={`${inputId}-helper`}
              className="label-text-alt text-base-content/60"
            >
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
