/**
 * ToggleField Component - Form Composition Layer
 *
 * Pre-bound toggle/switch that automatically wires to TanStack Form field context.
 * Eliminates boilerplate by auto-handling value, onChange, onBlur, errors,
 * and ARIA attributes.
 *
 * Usage within form:
 * ```tsx
 * <form.Field name="notifications">
 *   {(field) => (
 *     <field.ToggleField
 *       label="Enable notifications"
 *       description="Receive email updates"
 *     />
 *   )}
 * </form.Field>
 * ```
 */

import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { FormToggle } from '@/components/atoms/FormToggle'
import { useFieldContext } from '../contexts'
import type { ToggleFieldProps } from '../types'

export function ToggleField(props: ToggleFieldProps) {
  const field = useFieldContext<boolean>()
  const { t } = useTranslation('errors')

  // Extract error from field state and translate
  const error = useMemo(() => {
    const errors = field.state.meta.errors
    if (!errors || errors.length === 0) return undefined

    const firstError = errors[0]
    if (typeof firstError === 'string') {
      return t(firstError)
    }

    if (typeof firstError === 'object' && firstError && 'message' in firstError) {
      const errorObj = firstError as { message: string; code?: number }
      return t(errorObj.message)
    }

    return undefined
  }, [field.state.meta.errors, t])

  // Show error only after field has been touched (better UX)
  const showError = field.state.meta.isTouched && error

  return (
    <FormToggle
      id={field.name}
      checked={field.state.value}
      onChange={(e) => field.handleChange(e.target.checked)}
      onBlur={field.handleBlur}
      error={showError}
      disabled={props.disabled || field.state.meta.isValidating}
      aria-invalid={!field.state.meta.isValid && field.state.meta.isTouched}
      aria-describedby={showError ? `${field.name}-error` : undefined}
      {...props}
    />
  )
}
