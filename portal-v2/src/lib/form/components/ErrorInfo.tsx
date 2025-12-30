/**
 * ErrorInfo Component - Form Composition Layer
 *
 * Displays field validation errors with source-specific styling.
 * Distinguishes between real-time (onChange) and blur (onBlur) errors.
 *
 * Usage within form:
 * ```tsx
 * <form.Field name="email">
 *   {(field) => (
 *     <>
 *       <field.TextField label="Email" />
 *       <field.ErrorInfo />
 *     </>
 *   )}
 * </form.Field>
 * ```
 */

import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { useFieldContext } from '../contexts'

export function ErrorInfo() {
  const field = useFieldContext()
  const { t } = useTranslation('errors')

  const translatedError = useMemo(() => {
    const errors = field.state.meta.errors
    if (!errors || errors.length === 0) return null

    const firstError = errors[0]
    if (typeof firstError === 'string') {
      return t(firstError)
    }

    if (typeof firstError === 'object' && firstError && 'message' in firstError) {
      const errorObj = firstError as { message: string; code?: number }
      return t(errorObj.message)
    }

    return null
  }, [field.state.meta.errors, t])

  // Show error only after field has been touched
  if (!field.state.meta.isTouched || !translatedError) {
    return null
  }

  return (
    <div className="label">
      <span
        id={`${field.name}-error`}
        role="alert"
        className="label-text-alt text-error"
      >
        {translatedError}
      </span>
    </div>
  )
}
