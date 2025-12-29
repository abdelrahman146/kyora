import { forwardRef, useId } from 'react'
import { cn } from '../../lib/utils'
import type { TextareaHTMLAttributes } from 'react'

export interface FormTextareaProps
  extends Omit<TextareaHTMLAttributes<HTMLTextAreaElement>, 'size'> {
  label?: string
  error?: string
  helperText?: string
  size?: 'sm' | 'md' | 'lg'
  variant?: 'default' | 'filled' | 'ghost'
  fullWidth?: boolean
  maxLength?: number
  showCount?: boolean
}

/**
 * FormTextarea - Production-grade textarea component
 *
 * Features:
 * - RTL-first design with logical properties
 * - Mobile-optimized with proper touch targets (min-height per KDS)
 * - Auto-resize capability with resize-y
 * - Character counter with live updates (optional)
 * - Accessible with ARIA attributes
 * - Validation states with error messages
 * - Multiple variants and sizes
 * - Full keyboard navigation support
 *
 * @example
 * ```tsx
 * <FormTextarea
 *   label="Description"
 *   maxLength={500}
 *   showCount
 *   rows={4}
 *   error="Description too short"
 * />
 * ```
 */
export const FormTextarea = forwardRef<
  HTMLTextAreaElement,
  FormTextareaProps
>(
  (
    {
      label,
      error,
      helperText,
      size = 'md',
      variant = 'default',
      fullWidth = true,
      className,
      id,
      disabled,
      required,
      maxLength,
      showCount,
      value,
      ...props
    },
    ref
  ) => {
    const generatedId = useId()
    const inputId = id ?? generatedId

    const currentLength = typeof value === 'string' ? value.length : 0

    const sizeClasses = {
      sm: 'min-h-[88px] text-sm',
      md: 'min-h-[120px] text-base',
      lg: 'min-h-[160px] text-lg',
    }

    const variantClasses = {
      default: 'textarea-bordered bg-base-100',
      filled:
        'textarea-bordered bg-base-200/50 border-transparent focus:bg-base-100',
      ghost: 'textarea-ghost bg-transparent',
    }

    return (
      <div className={cn('form-control', fullWidth && 'w-full')}>
        {label && (
          <label htmlFor={inputId} className="label">
            <span className="label-text text-base-content/70 font-medium">
              {label}
              {required && <span className="text-error ms-1">*</span>}
            </span>
            {showCount && maxLength && (
              <span className="label-text-alt text-base-content/50">
                {currentLength}/{maxLength}
              </span>
            )}
          </label>
        )}

        <textarea
          ref={ref}
          id={inputId}
          disabled={disabled}
          required={required}
          maxLength={maxLength}
          value={value}
          className={cn(
            'textarea w-full transition-all duration-200',
            sizeClasses[size],
            variantClasses[variant],
            'text-base-content text-start placeholder:text-base-content/40',
            'focus:outline-none focus:border-primary focus:ring-2 focus:ring-primary/20',
            'resize-y',
            error &&
              'textarea-error border-error focus:border-error focus:ring-error/20',
            disabled && 'opacity-60 cursor-not-allowed',
            className
          )}
          aria-invalid={error ? 'true' : 'false'}
          aria-describedby={
            error
              ? `${inputId}-error`
              : helperText
                ? `${inputId}-helper`
                : undefined
          }
          aria-required={required}
          {...props}
        />

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

        {!error && helperText && (
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
  }
)

FormTextarea.displayName = 'FormTextarea'
