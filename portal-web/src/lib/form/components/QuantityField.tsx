/**
 * QuantityField Component - Form Composition Layer
 *
 * Pre-bound quantity input with increment/decrement buttons.
 * Automatically wires to TanStack Form field context.
 * Follows the same UI patterns as TextField for consistency.
 *
 * Usage within form:
 * ```tsx
 * <form.Field name="quantity">
 *   {(field) => (
 *     <field.QuantityField
 *       label="Quantity"
 *       min={1}
 *       max={100}
 *     />
 *   )}
 * </form.Field>
 * ```
 */

import { useId, useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { useFieldContext } from '../contexts'
import type { BaseFieldProps } from '../types'
import { QuantityInput } from '@/components/atoms/QuantityInput'
import { cn } from '@/lib/utils'
import { getErrorText } from '@/lib/formErrors'

export interface QuantityFieldProps extends BaseFieldProps {
  /** Minimum allowed value (default: 0) */
  min?: number
  /** Maximum allowed value (default: Infinity) */
  max?: number
  /** Step increment/decrement value (default: 1) */
  step?: number
  /** Helper text shown below input */
  helperText?: string
  /** Full width (default: true) */
  fullWidth?: boolean
}

export function QuantityField({
  label,
  hint,
  helperText,
  required,
  disabled,
  min = 0,
  max = Infinity,
  step = 1,
  fullWidth = true,
}: QuantityFieldProps) {
  const field = useFieldContext<number>()
  const { t } = useTranslation('errors')
  const { t: tCommon } = useTranslation('common')
  const generatedId = useId()
  const inputId = generatedId

  // Extract error from field state and translate
  const error = useMemo(() => {
    const errors = field.state.meta.errors
    if (errors.length === 0) return undefined

    const firstError = errors[0]
    if (typeof firstError === 'string') {
      return t(firstError, { min, max })
    }

    if (
      typeof firstError === 'object' &&
      firstError &&
      'message' in firstError
    ) {
      const errorObj = firstError as { message: string; code?: number }
      return t(errorObj.message, { min, max })
    }

    return undefined
  }, [field.state.meta.errors, t, min, max])

  // Show error only after field has been touched (better UX)
  const showError = field.state.meta.isTouched && error
  const errorText = getErrorText(showError)
  const hasError = Boolean(errorText)

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

      <QuantityInput
        id={inputId}
        name={field.name}
        value={field.state.value}
        onChange={field.handleChange}
        onBlur={field.handleBlur}
        min={min}
        max={max}
        step={step}
        disabled={disabled || field.state.meta.isValidating}
        required={required}
        incrementLabel={tCommon('increment')}
        decrementLabel={tCommon('decrement')}
        error={showError}
      />

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

      {!hasError && (helperText || hint) && (
        <label className="label">
          <span
            id={`${inputId}-helper`}
            className="label-text-alt text-base-content/60"
          >
            {helperText || hint}
          </span>
        </label>
      )}
    </div>
  )
}
