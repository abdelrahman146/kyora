/**
 * FormRoot Component - Form Wrapper with Submit Handling
 *
 * Wraps form element with proper submit event handling.
 * Prevents default browser submission and handles form validation.
 *
 * Usage:
 * ```tsx
 * <form.FormRoot>
 *   <form.Field name="email">...</form.Field>
 *   <form.SubmitButton>Submit</form.SubmitButton>
 * </form.FormRoot>
 * ```
 */

/**
 * FormRoot Component - Form Wrapper with Submit Handling
 *
 * Wraps form element with proper submit event handling.
 * Prevents default browser submission and handles form validation.
 *
 * Note: This component must be used as `form.FormRoot` from useKyoraForm,
 * not imported directly, to ensure proper form context.
 *
 * Usage:
 * ```tsx
 * const form = useKyoraForm({...})
 * <form.FormRoot>
 *   <form.Field name="email">...</form.Field>
 *   <form.SubmitButton>Submit</form.SubmitButton>
 * </form.FormRoot>
 * ```
 */

import { useFormContext } from '../contexts'
import type { FormHTMLAttributes } from 'react'

interface FormRootProps extends Omit<FormHTMLAttributes<HTMLFormElement>, 'onSubmit'> {
  children: React.ReactNode
}

export function FormRoot({ children, ...props }: FormRootProps) {
  // This will throw if not used as form.FormRoot - which is correct!
  // The error message will guide users to use it through the composition API
  const form = useFormContext()

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    e.stopPropagation()
    void form.handleSubmit()
  }

  return (
    <form onSubmit={handleSubmit} noValidate {...props}>
      {children}
    </form>
  )
}


