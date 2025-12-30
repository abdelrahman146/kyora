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

import type { FormHTMLAttributes } from 'react'
import { useFormContext } from '../contexts'

interface FormRootProps extends Omit<FormHTMLAttributes<HTMLFormElement>, 'onSubmit'> {
  children: React.ReactNode
}

export function FormRoot({ children, ...props }: FormRootProps) {
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
