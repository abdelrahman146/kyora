import { forwardRef, useId } from 'react'
import { cn } from '../../lib/utils'
import type { InputHTMLAttributes } from 'react'

export interface FormToggleProps
  extends Omit<InputHTMLAttributes<HTMLInputElement>, 'type' | 'size'> {
  label?: string
  description?: string
  error?: string
  size?: 'sm' | 'md' | 'lg'
  variant?: 'default' | 'primary' | 'secondary'
  labelPosition?: 'start' | 'end'
}

/**
 * FormToggle - Production-grade toggle/switch component
 *
 * Features:
 * - RTL-first design with logical properties
 * - Mobile-optimized touch target for easy interaction
 * - Accessible with ARIA attributes (role="switch", aria-checked)
 * - Flexible label positioning (start/end) for layout versatility
 * - Validation states with error messages
 * - Multiple sizes and color variants
 * - Keyboard navigation support
 * - Disabled state handling with visual feedback
 *
 * @example
 * ```tsx
 * <FormToggle
 *   label="Enable notifications"
 *   description="Receive email updates"
 *   labelPosition="start"
 *   error="Required field"
 *   variant="primary"
 * />
 * ```
 */
export const FormToggle = forwardRef<HTMLInputElement, FormToggleProps>(
  (
    {
      label,
      description,
      error,
      size = 'md',
      variant = 'primary',
      labelPosition = 'start',
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
      sm: 'toggle-sm',
      md: 'toggle-md',
      lg: 'toggle-lg',
    }

    const variantClasses = {
      default: '',
      primary: 'toggle-primary',
      secondary: 'toggle-secondary',
    }

    const labelContent = (label ?? description) && (
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
    )

    return (
      <div className="form-control">
        <label
          htmlFor={inputId}
          className={cn(
            'label cursor-pointer',
            labelPosition === 'start'
              ? 'justify-between'
              : 'justify-start gap-3',
            disabled && 'opacity-60 cursor-not-allowed'
          )}
        >
          {labelPosition === 'start' && labelContent}

          <input
            ref={ref}
            type="checkbox"
            id={inputId}
            disabled={disabled}
            className={cn(
              'toggle',
              sizeClasses[size],
              variantClasses[variant],
              error && 'toggle-error',
              className
            )}
            role="switch"
            aria-checked={props.checked}
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

          {labelPosition === 'end' && labelContent}
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

FormToggle.displayName = 'FormToggle'
