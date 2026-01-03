import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'

import { useFieldContext } from '../contexts'
import type { PriceFieldProps } from '../types'
import { PriceInput } from '@/components/atoms/PriceInput'

export function PriceField(props: PriceFieldProps) {
  const field = useFieldContext<string>()
  const { t } = useTranslation('errors')

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

  const showError = field.state.meta.isTouched && error

  return (
    <PriceInput
      id={field.name}
      value={field.state.value}
      onChange={(e) => field.handleChange(e.target.value.replace(',', '.'))}
      onBlur={field.handleBlur}
      error={showError}
      disabled={props.disabled || field.state.meta.isValidating}
      aria-invalid={!field.state.meta.isValid && field.state.meta.isTouched}
      aria-describedby={showError ? `${field.name}-error` : undefined}
      currencyCode={props.currencyCode}
      helperText={props.hint}
      {...props}
    />
  )
}
