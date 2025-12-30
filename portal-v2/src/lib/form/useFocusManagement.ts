/**
 * Focus Management Utilities for Kyora Forms
 *
 * Provides automatic focus management for form validation errors
 * and initial field autofocus behavior.
 *
 * Features:
 * - Auto-focus first invalid field on submit error
 * - Autofocus on component mount
 * - Keyboard-accessible focus handling
 * - Integration with TanStack Form's onSubmitInvalid
 */

import { useEffect, useRef } from 'react'
import type { FormApi } from '@tanstack/react-form'

/**
 * Focus Management Hook
 *
 * Returns an onSubmitInvalid handler that automatically focuses
 * the first invalid field when form submission fails validation.
 *
 * @example
 * ```tsx
 * const form = useKyoraForm({
 *   defaultValues: {...},
 *   onSubmit: async ({ value }) => {...},
 *   onSubmitInvalid: useFocusOnError(),
 * })
 * ```
 */
export function useFocusOnError() {
  return () => {
    // Query for the first invalid input using ARIA attributes
    // This works because our field components automatically set aria-invalid
    const invalidInput = document.querySelector('[aria-invalid="true"]') as
      | HTMLInputElement
      | HTMLTextAreaElement
      | HTMLSelectElement
      | null

    if (invalidInput) {
      // Focus the invalid field
      invalidInput.focus()

      // Scroll into view if needed (with smooth behavior)
      invalidInput.scrollIntoView({
        behavior: 'smooth',
        block: 'center',
      })
    }
  }
}

/**
 * Auto-Focus Hook
 *
 * Automatically focuses an input element on component mount.
 * Useful for first field in forms or modal inputs.
 *
 * @param enabled - Whether autofocus is enabled (default: true)
 * @returns Ref to attach to the input element
 *
 * @example
 * ```tsx
 * function LoginForm() {
 *   const emailRef = useAutoFocus()
 *
 *   return (
 *     <form.Field name="email">
 *       {(field) => (
 *         <field.TextField
 *           label="Email"
 *           ref={emailRef}
 *         />
 *       )}
 *     </form.Field>
 *   )
 * }
 * ```
 */
export function useAutoFocus<T extends HTMLElement = HTMLInputElement>(
  enabled: boolean = true,
) {
  const ref = useRef<T>(null)

  useEffect(() => {
    if (enabled && ref.current) {
      // Small delay to ensure rendering is complete
      const timeoutId = setTimeout(() => {
        ref.current?.focus()
      }, 100)

      return () => clearTimeout(timeoutId)
    }
  }, [enabled])

  return ref
}

/**
 * Create a focus management configuration for forms
 *
 * Returns form options that include onSubmitInvalid handler
 * with automatic focus management.
 *
 * @example
 * ```tsx
 * const form = useKyoraForm({
 *   ...createFocusManagement(),
 *   defaultValues: {...},
 *   onSubmit: async ({ value }) => {...},
 * })
 * ```
 */
export function createFocusManagement() {
  return {
    onSubmitInvalid: () => {
      const invalidInput = document.querySelector('[aria-invalid="true"]') as
        | HTMLInputElement
        | HTMLTextAreaElement
        | HTMLSelectElement
        | null

      if (invalidInput) {
        invalidInput.focus()
        invalidInput.scrollIntoView({
          behavior: 'smooth',
          block: 'center',
        })
      }
    },
  }
}
