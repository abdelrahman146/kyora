/**
 * FormError Component - Form-Level Error Display
 *
 * Displays form-level errors that aren't specific to a single field.
 * Subscribes to form.state.errorMap for form-level validation errors.
 *
 * Usage:
 * ```tsx
 * <form.FormRoot>
 *   <form.FormError />
 *   <form.Field name="email">...</form.Field>
 * </form.FormRoot>
 * ```
 */

import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { useFormContext } from '../contexts'

export function FormError() {
  const form = useFormContext()
  const { t } = useTranslation('errors')

  return (
    <form.Subscribe selector={(state) => state.errorMap}>
      {(errorMap) => {
        // Get form-level error from onSubmit validator
        const formError = errorMap.onSubmit

        // Translate error if it exists
        const translatedError = useMemo(() => {
          if (!formError) return null

          if (typeof formError === 'string') {
            return t(formError)
          }

          if (typeof formError === 'object' && 'message' in formError) {
            const errorObj = formError as { message: string; code?: number }
            return t(errorObj.message)
          }

          return null
        }, [formError, t])

        if (!translatedError) {
          return null
        }

        return (
          <div role="alert" className="alert alert-error">
            <span>{translatedError}</span>
          </div>
        )
      }}
    </form.Subscribe>
  )
}
