import { forwardRef, useId, useState } from 'react'
import { Eye, EyeOff, Lock } from 'lucide-react'
import { cn } from '../../lib/utils'
import type { InputHTMLAttributes } from 'react'

export interface PasswordInputProps
  extends Omit<InputHTMLAttributes<HTMLInputElement>, 'type'> {
  label?: string
  error?: string
  helperText?: string
  fullWidth?: boolean
  showPasswordToggle?: boolean
  showDefaultIcon?: boolean
}

/**
 * PasswordInput - Production-grade password input component
 *
 * Features:
 * - RTL-first design with logical properties (start/end instead of left/right)
 * - Mobile-optimized (min-height 50px per KDS)
 * - Show/hide password toggle with keyboard accessibility
 * - Lock icon by default (can be disabled)
 * - Accessible with ARIA attributes
 * - Validation states with error messages
 * - Helper text for password requirements
 * - Supports all standard input props
 * - Focus management and keyboard navigation
 *
 * @example
 * ```tsx
 * <PasswordInput
 *   label="Password"
 *   error="Password too short"
 *   helperText="At least 8 characters"
 *   showPasswordToggle
 * />
 * ```
 */
export const PasswordInput = forwardRef<HTMLInputElement, PasswordInputProps>(
  (
    {
      label,
      error,
      helperText,
      fullWidth = true,
      showPasswordToggle = true,
      showDefaultIcon = true,
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
    const [isPasswordVisible, setIsPasswordVisible] = useState(false)

    const togglePasswordVisibility = () => {
      setIsPasswordVisible((prev) => !prev)
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

        <div className="relative">
          {/* Lock Icon (start) */}
          {showDefaultIcon && (
            <div className="absolute inset-y-0 start-0 z-10 flex items-center ps-3 pointer-events-none text-base-content/50">
              <Lock size={20} aria-hidden="true" />
            </div>
          )}

          {/* Password Input */}
          <input
            ref={ref}
            id={inputId}
            type={isPasswordVisible ? 'text' : 'password'}
            disabled={disabled}
            required={required}
            className={cn(
              'input input-bordered relative z-0 w-full h-[50px] text-base transition-all duration-200',
              'bg-base-100 text-base-content',
              'text-start placeholder:text-base-content/40',
              'focus:outline-none focus:border-primary focus:ring-2 focus:ring-primary/20',
              error &&
                'input-error border-error focus:border-error focus:ring-error/20',
              disabled && 'opacity-60 cursor-not-allowed',
              showDefaultIcon ? 'ps-10' : '',
              showPasswordToggle ? 'pe-10' : '',
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

          {/* Show/Hide Password Toggle (end) */}
          {showPasswordToggle && (
            <button
              type="button"
              onClick={togglePasswordVisibility}
              disabled={disabled}
              className={cn(
                'absolute inset-y-0 end-0 z-10 flex items-center pe-3',
                'text-base-content/50 hover:text-base-content/70',
                'transition-colors duration-200',
                'focus:outline-none focus-visible:ring-2 focus-visible:ring-primary/50 rounded',
                disabled && 'opacity-60 cursor-not-allowed'
              )}
              aria-label={isPasswordVisible ? 'Hide password' : 'Show password'}
              tabIndex={disabled ? -1 : 0}
            >
              {isPasswordVisible ? (
                <EyeOff size={20} aria-hidden="true" />
              ) : (
                <Eye size={20} aria-hidden="true" />
              )}
            </button>
          )}
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

PasswordInput.displayName = 'PasswordInput'
