/**
 * SelectField Component - Form Composition Layer
 *
 * Pre-bound select/dropdown that automatically wires to TanStack Form field context.
 * Eliminates boilerplate by auto-handling value, onChange, onBlur, errors,
 * and ARIA attributes. Supports both single and multi-select modes.
 *
 * Usage within form:
 * ```tsx
 * // Single select
 * <form.Field name="country">
 *   {(field) => (
 *     <field.SelectField
 *       label="Country"
 *       options={countries}
 *       searchable
 *     />
 *   )}
 * </form.Field>
 *
 * // Multi-select
 * <form.Field name="tags">
 *   {(field) => (
 *     <field.SelectField
 *       label="Tags"
 *       options={tagOptions}
 *       multiSelect
 *       searchable
 *     />
 *   )}
 * </form.Field>
 * ```
 */

import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { useFieldContext } from '../contexts'
import type { SelectFieldProps } from '../types'
import { FormSelect } from '@/components/form/FormSelect'

export function SelectField<T extends string = string>(
  props: SelectFieldProps<T>,
) {
  const field = useFieldContext<T | Array<T>>()
  const { t } = useTranslation('errors')

  // Extract error from field state and translate
  const error = useMemo(() => {
    const errors = field.state.meta.errors
    if (errors.length === 0) return undefined

    const firstError = errors[0]
    if (typeof firstError === 'string') {
      return t(firstError)
    }

    if (
      typeof firstError === 'object' &&
      firstError &&
      'message' in firstError
    ) {
      const errorObj = firstError as { message: string; code?: number }
      return t(errorObj.message)
    }

    return undefined
  }, [field.state.meta.errors, t])

  // Show error only after field has been touched (better UX)
  const showError = field.state.meta.isTouched && error

  return (
    <FormSelect<T>
      id={field.name}
      value={field.state.value}
      onChange={(value) => field.handleChange(value)}
      error={showError}
      disabled={props.disabled || field.state.meta.isValidating}
      aria-invalid={!field.state.meta.isValid && field.state.meta.isTouched}
      aria-describedby={showError ? `${field.name}-error` : undefined}
      {...props}
    />
  )
}
