import { forwardRef, useId } from 'react'
import { cn } from '../../lib/utils'
import type { InputHTMLAttributes } from 'react'

export interface FormRadioOption {
  value: string
  label: string
  description?: string
  disabled?: boolean
}

export interface FormRadioProps
  extends Omit<InputHTMLAttributes<HTMLInputElement>, 'type' | 'size'> {
  label?: string
  options: Array<FormRadioOption>
  error?: string
  size?: 'sm' | 'md' | 'lg'
  variant?: 'default' | 'primary' | 'secondary'
  orientation?: 'vertical' | 'horizontal'
  name: string
}

/**
 * FormRadio - Production-grade radio button group component
 *
 * Features:
 * - RTL-first design with logical properties
 * - Mobile-optimized touch targets for easy selection
 * - Accessible with ARIA attributes (role="radiogroup")
 * - Supports descriptions for each option for better UX
 * - Validation states with error messages
 * - Flexible layout: vertical or horizontal orientation
 * - Multiple sizes and color variants
 * - Individual option disabled states
 * - Keyboard navigation support (arrow keys)
 *
 * @example
 * ```tsx
 * <FormRadio
 *   name="plan"
 *   label="Select a plan"
 *   options={[
 *     { value: "free", label: "Free", description: "$0/month" },
 *     { value: "pro", label: "Pro", description: "$10/month" }
 *   ]}
 *   error="Please select a plan"
 *   orientation="vertical"
 * />
 * ```
 */
export const FormRadio = forwardRef<HTMLInputElement, FormRadioProps>(
  (
    {
      label,
      options,
      error,
      size = 'md',
      variant = 'primary',
      orientation = 'vertical',
      name,
      className,
      value: selectedValue,
      onChange,
      ...props
    },
    ref
  ) => {
    const generatedId = useId()
    const groupId = `${generatedId}-group`

    const sizeClasses = {
      sm: 'radio-sm',
      md: 'radio-md',
      lg: 'radio-lg',
    }

    const variantClasses = {
      default: '',
      primary: 'radio-primary',
      secondary: 'radio-secondary',
    }

    return (
      <div className="form-control">
        {label && (
          <label className="label">
            <span className="label-text text-base-content/70 font-medium">
              {label}
            </span>
          </label>
        )}

        <div
          role="radiogroup"
          aria-labelledby={label ? groupId : undefined}
          aria-invalid={error ? 'true' : 'false'}
          aria-describedby={error ? `${groupId}-error` : undefined}
          className={cn(
            'flex gap-4',
            orientation === 'vertical' ? 'flex-col' : 'flex-row flex-wrap'
          )}
        >
          {options.map((option) => {
            const optionId = `${groupId}-${option.value}`
            const isChecked = selectedValue === option.value

            return (
              <label
                key={option.value}
                htmlFor={optionId}
                className={cn(
                  'label cursor-pointer justify-start gap-3',
                  (option.disabled ?? props.disabled) &&
                    'opacity-60 cursor-not-allowed'
                )}
              >
                <input
                  ref={isChecked ? ref : undefined}
                  type="radio"
                  id={optionId}
                  name={name}
                  value={option.value}
                  checked={isChecked}
                  onChange={onChange}
                  disabled={option.disabled ?? props.disabled}
                  className={cn(
                    'radio',
                    sizeClasses[size],
                    variantClasses[variant],
                    error && 'radio-error',
                    className
                  )}
                  {...props}
                />

                <div className="flex flex-col gap-1">
                  <span className="label-text text-base-content font-medium">
                    {option.label}
                  </span>
                  {option.description && (
                    <span className="label-text-alt text-base-content/60">
                      {option.description}
                    </span>
                  )}
                </div>
              </label>
            )
          })}
        </div>

        {error && (
          <label className="label">
            <span
              id={`${groupId}-error`}
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

FormRadio.displayName = 'FormRadio'
