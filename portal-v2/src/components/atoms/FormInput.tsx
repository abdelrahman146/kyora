import { forwardRef, useId } from 'react'
import { cn } from '../../lib/utils'
import type { InputHTMLAttributes, ReactNode } from 'react'

export interface FormInputProps extends Omit<
  InputHTMLAttributes<HTMLInputElement>,
  'size'
> {
  label?: string
  error?: string
  helperText?: string
  startIcon?: ReactNode
  endIcon?: ReactNode
  size?: 'sm' | 'md' | 'lg'
  variant?: 'default' | 'filled' | 'ghost'
  fullWidth?: boolean
}

/**
 * FormInput - Production-grade text input component
 *
 * Features:
 * - RTL-first design with logical properties
 * - Mobile-optimized (min-height 50px per KDS)
 * - Accessible with ARIA attributes
 * - Supports icons, validation states, and helper text
 * - Multiple variants and sizes
 * - Keyboard-specific inputs (email, tel, search)
 *
 * @example
 * <FormInput
 *   label="Email"
 *   type="email"
 *   error="Invalid email"
 *   startIcon={<Mail />}
 * />
 */
export const FormInput = forwardRef<HTMLInputElement, FormInputProps>(
  (
    {
      label,
      error,
      helperText,
      startIcon,
      endIcon,
      size = 'md',
      variant = 'default',
      fullWidth = true,
      className,
      id,
      disabled,
      required,
      ...props
    },
    ref,
  ) => {
    const generatedId = useId()
    const inputId = id ?? generatedId

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
          {startIcon && (
            <div className="absolute inset-y-0 start-0 z-10 flex items-center ps-3 pointer-events-none text-base-content/50">
              <span aria-hidden="true">{startIcon}</span>
            </div>
          )}

          <input
            ref={ref}
            id={inputId}
            disabled={disabled}
            required={required}
            className={cn(
              'input relative z-0 w-full transition-all duration-200',
              sizeClasses[size],
              variantClasses[variant],
              'text-base-content text-start placeholder:text-base-content/40',
              'focus:outline-none focus:border-primary focus:ring-2 focus:ring-primary/20',
              error &&
                'input-error border-error focus:border-error focus:ring-error/20',
              disabled && 'opacity-60 cursor-not-allowed',
              startIcon ? 'ps-10' : '',
              endIcon ? 'pe-10' : '',
              className,
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

          {endIcon && (
            <div className="absolute inset-y-0 end-0 z-10 flex items-center pe-3 text-base-content/50">
              <span aria-hidden="true">{endIcon}</span>
            </div>
          )}
        </div>

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
  },
)

FormInput.displayName = 'FormInput'
