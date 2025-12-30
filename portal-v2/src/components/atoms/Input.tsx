import { forwardRef, useId } from 'react'
import type { InputHTMLAttributes } from 'react'
import { cn } from '@/lib/utils'
import { getErrorText } from '@/lib/formErrors'

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string
  error?: unknown
  helperText?: string
  startIcon?: React.ReactNode
  endIcon?: React.ReactNode
}

/**
 * Input Component
 *
 * Form input with label, error states, and icon support.
 * Uses daisyUI input classes with RTL-ready logical properties.
 */
export const Input = forwardRef<HTMLInputElement, InputProps>(
  (
    { label, error, helperText, startIcon, endIcon, className, id, ...props },
    ref,
  ) => {
    const generatedId = useId()
    const inputId = id ?? generatedId
    const errorText = getErrorText(error)
    const hasError = Boolean(errorText)

    return (
      <div className="w-full">
        {label && (
          <label
            htmlFor={inputId}
            className="label mb-2 block text-sm font-medium text-base-content"
          >
            <span className="label-text">{label}</span>
          </label>
        )}
        <div className="relative">
          {startIcon && (
            <div className="pointer-events-none absolute inset-y-0 start-0 z-10 flex items-center ps-3 text-base-content/50">
              <span aria-hidden="true">{startIcon}</span>
            </div>
          )}
          <input
            ref={ref}
            id={inputId}
            className={cn(
              'input input-bordered relative z-0 h-[50px] w-full',
              'bg-base-100 text-base-content',
              'text-start placeholder:text-base-content/40',
              'transition-all duration-200',
              'focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/20',
              hasError
                ? 'input-error border-error focus:border-error focus:ring-error/20'
                : '',
              startIcon ? 'ps-10' : '',
              endIcon ? 'pe-10' : '',
              className,
            )}
            aria-invalid={hasError ? 'true' : 'false'}
            aria-describedby={
              hasError
                ? `${inputId}-error`
                : helperText
                  ? `${inputId}-helper`
                  : undefined
            }
            {...props}
          />
          {endIcon && (
            <div className="pointer-events-none absolute inset-y-0 end-0 z-10 flex items-center pe-3 text-base-content/50">
              <span aria-hidden="true">{endIcon}</span>
            </div>
          )}
        </div>
        {hasError && (
          <p
            id={`${inputId}-error`}
            className="mt-2 text-sm text-error"
            role="alert"
          >
            {errorText}
          </p>
        )}
        {!hasError && helperText && (
          <p
            id={`${inputId}-helper`}
            className="mt-2 text-sm text-base-content/60"
          >
            {helperText}
          </p>
        )}
      </div>
    )
  },
)

Input.displayName = 'Input'
