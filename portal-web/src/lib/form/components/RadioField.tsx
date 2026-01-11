/**
 * RadioField Component - Form Composition Layer
 *
 * Pre-bound radio button group that automatically wires to TanStack Form field context.
 * Eliminates boilerplate by auto-handling value, onChange, onBlur, errors,
 * and ARIA attributes.
 *
 * Usage within form:
 * ```tsx
 * <form.Field name="plan">
 *   {(field) => (
 *     <field.RadioField
 *       label="Select a plan"
 *       options={[
 *         { value: "free", label: "Free", description: "$0/month" },
 *         { value: "pro", label: "Pro", description: "$10/month" }
 *       ]}
 *       orientation="vertical"
 *     />
 *   )}
 * </form.Field>
 * ```
 */

import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { useFieldContext } from '../contexts'
import type { FormRadioProps } from '@/components/form/FormRadio'
import { FormRadio } from '@/components/form/FormRadio'

export interface RadioFieldProps extends Omit<
  FormRadioProps,
  'name' | 'value' | 'onChange'
> {
  disabled?: boolean
}

export function RadioField(props: RadioFieldProps) {
  const field = useFieldContext<string>()
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
    <FormRadio
      name={field.name}
      value={field.state.value}
      onChange={(e) => field.handleChange(e.target.value)}
      onBlur={field.handleBlur}
      error={showError}
      disabled={props.disabled || field.state.meta.isValidating}
      aria-invalid={!field.state.meta.isValid && field.state.meta.isTouched}
      aria-describedby={showError ? `${field.name}-error` : undefined}
      {...props}
    />
  )
}
