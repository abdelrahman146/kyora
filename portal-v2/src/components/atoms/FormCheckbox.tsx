import { forwardRef, useId } from 'react'
import { cn } from '../../lib/utils'
import type { InputHTMLAttributes } from 'react'

export interface FormCheckboxProps
  extends Omit<InputHTMLAttributes<HTMLInputElement>, 'type' | 'size'> {
  label?: string
  description?: string
  error?: string
  size?: 'sm' | 'md' | 'lg'
  variant?: 'default' | 'primary' | 'secondary'
}

/**
 * FormCheckbox - Production-grade checkbox component
 *
 * Features:
 * - RTL-first design with logical properties
 * - Mobile-optimized touch target (44x44px minimum per WCAG)
 * - Accessible with ARIA attributes and semantic HTML
 * - Label and description support for clarity
 * - Validation states with error messages
 * - Multiple sizes and color variants
 * - Keyboard navigation support
 * - Disabled state handling with visual feedback
 *
 * @example
 * ```tsx
 * <FormCheckbox
 *   label="Accept terms"
 *   description="I agree to the terms and conditions"
 *   error="You must accept the terms"
 *   variant="primary"
 * />
 * ```
 */
export const FormCheckbox = forwardRef<HTMLInputElement, FormCheckboxProps>(
  (
    {
      label,
      description,
      error,
      size = 'md',
      variant = 'primary',
      className,
      id,
      disabled,
      ...props
    },
    ref
  ) => {
    const generatedId = useId()
    const inputId = id ?? generatedId

    const sizeClasses = {
      sm: 'checkbox-sm',
      md: 'checkbox-md',
      lg: 'checkbox-lg',
    }

    const variantClasses = {
      default: '',
      primary: 'checkbox-primary',
      secondary: 'checkbox-secondary',
    }

    return (
      <div className="form-control">
        <label
          htmlFor={inputId}
          className={cn(
            'label cursor-pointer justify-start gap-3',
            disabled && 'opacity-60 cursor-not-allowed'
          )}
        >
          <input
            ref={ref}
            type="checkbox"
            id={inputId}
            disabled={disabled}
            className={cn(
              'checkbox',
              sizeClasses[size],
              variantClasses[variant],
              error && 'checkbox-error',
              className
            )}
            aria-invalid={error ? 'true' : 'false'}
            aria-describedby={
              error
                ? `${inputId}-error`
                : description
                  ? `${inputId}-description`
                  : undefined
            }
            {...props}
          />

          <div className="flex flex-col gap-1">
            {label && (
              <span className="label-text text-base-content font-medium">
                {label}
              </span>
            )}
            {description && (
              <span
                id={`${inputId}-description`}
                className="label-text-alt text-base-content/60"
              >
                {description}
              </span>
            )}
          </div>
        </label>

        {error && (
          <label className="label pt-0">
            <span
              id={`${inputId}-error`}
              className="label-text-alt text-error"
              role="alert"
            >
              {error}
            </span>
          </label>
        )}
      </div>
    )
  }
)

FormCheckbox.displayName = 'FormCheckbox'
