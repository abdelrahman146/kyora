import { forwardRef, useId, useRef } from 'react'
import type { InputHTMLAttributes } from 'react'
import { cn } from '@/lib/utils'

export interface FormFileInputProps extends Omit<
  InputHTMLAttributes<HTMLInputElement>,
  'type' | 'size'
> {
  onFilesSelected?: (files: Array<File>) => void
  label?: string
  helperText?: string
  error?: unknown
  size?: 'sm' | 'md' | 'lg'
  fullWidth?: boolean
}

/**
 * FormFileInput - Hidden file input with accessible trigger
 *
 * Features:
 * - Hidden native file input
 * - Keyboard accessible (Space/Enter to trigger)
 * - Mobile camera/gallery support via capture attribute
 * - Multiple file selection support
 * - ARIA labels for screen readers
 * - RTL-first design
 *
 * @example
 * <FormFileInput
 *   accept="image/*"
 *   multiple
 *   capture="environment"
 *   onFilesSelected={(files) => console.log(files)}
 * />
 */
export const FormFileInput = forwardRef<HTMLInputElement, FormFileInputProps>(
  (
    {
      onFilesSelected,
      label,
      helperText,
      error,
      size = 'md',
      fullWidth = true,
      className,
      disabled,
      id,
      ...props
    },
    ref,
  ) => {
    const generatedId = useId()
    const inputId = id ?? generatedId
    const internalRef = useRef<HTMLInputElement>(null)

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      const files = Array.from(e.target.files || [])
      if (files.length > 0) {
        onFilesSelected?.(files)
      }
      props.onChange?.(e)
    }

    const handleTriggerClick = () => {
      if (ref && 'current' in ref && ref.current) {
        ref.current.click()
      } else if (internalRef.current) {
        internalRef.current.click()
      }
    }

    const handleKeyDown = (e: React.KeyboardEvent) => {
      if (e.key === ' ' || e.key === 'Enter') {
        e.preventDefault()
        handleTriggerClick()
      }
    }

    const sizeClasses = {
      sm: 'h-[44px]',
      md: 'h-[50px]',
      lg: 'h-[56px]',
    }

    const hasError = Boolean(error)

    return (
      <div className={cn('form-control', fullWidth && 'w-full')}>
        {label && (
          <label htmlFor={inputId} className="label">
            <span className="label-text text-base-content/70 font-medium">
              {label}
              {props.required && <span className="text-error ms-1">*</span>}
            </span>
          </label>
        )}

        {/* Hidden native input */}
        <input
          ref={ref || internalRef}
          type="file"
          id={inputId}
          className="hidden"
          disabled={disabled}
          {...props}
          onChange={handleChange}
        />

        {/* Visible trigger button */}
        <div
          role="button"
          tabIndex={disabled ? -1 : 0}
          onClick={disabled ? undefined : handleTriggerClick}
          onKeyDown={disabled ? undefined : handleKeyDown}
          className={cn(
            'btn btn-outline',
            sizeClasses[size],
            'w-full cursor-pointer',
            hasError && 'btn-error',
            disabled && 'btn-disabled cursor-not-allowed',
            className,
          )}
          aria-label={label || 'Select files'}
          aria-disabled={disabled}
        >
          <span className="text-base-content/70">
            {props.accept?.includes('image') ? 'Choose Images' : 'Choose Files'}
          </span>
        </div>

        {helperText && !hasError && (
          <label className="label">
            <span className="label-text-alt text-base-content/60">
              {helperText}
            </span>
          </label>
        )}

        {hasError && typeof error === 'string' && (
          <label className="label">
            <span className="label-text-alt text-error" role="alert">
              {error}
            </span>
          </label>
        )}
      </div>
    )
  },
)

FormFileInput.displayName = 'FormFileInput'
