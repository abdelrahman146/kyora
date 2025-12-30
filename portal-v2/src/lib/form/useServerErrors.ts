/**
 * Server Error Injection for Kyora Forms
 *
 * Provides utilities to inject backend validation errors into form fields.
 * Integrates with TanStack Form's form-level validators to map server
 * errors to specific fields.
 *
 * Features:
 * - Parse RFC7807 ProblemDetails responses
 * - Map backend field errors to form fields
 * - Automatic error translation via i18n
 * - Type-safe error injection
 */

import { useCallback } from 'react'
import { parseProblemDetails } from '../errorParser'
import type { ServerErrors } from './types'
import type { HTTPError } from 'ky'

/**
 * Server Error Hook
 *
 * Provides utilities to handle server-side validation errors and inject
 * them into form fields.
 *
 * @returns Functions to parse and inject server errors
 *
 * @example
 * ```tsx
 * const { parseServerError, createServerErrorValidator } = useServerErrors()
 *
 * const form = useKyoraForm({
 *   defaultValues: { email: '', password: '' },
 *   onSubmit: async ({ value }) => {
 *     try {
 *       await loginMutation.mutateAsync(value)
 *     } catch (error) {
 *       const serverErrors = await parseServerError(error)
 *       // Handle display if needed
 *     }
 *   },
 * })
 * ```
 */
export function useServerErrors() {
  /**
   * Parse server error response into field-specific errors
   *
   * Extracts field validation errors from backend responses.
   * Returns a map of field names to error translation keys.
   */
  const parseServerError = useCallback(
    async (
      error: unknown,
    ): Promise<{ fields?: ServerErrors; form?: string }> => {
      if (!error) {
        return { form: 'generic.unexpected' }
      }

      // Handle HTTPError from ky
      if (
        typeof error === 'object' &&
        'response' in error &&
        error.response instanceof Response
      ) {
        const httpError = error as HTTPError

        try {
          // Try to parse JSON response body
          const body: unknown = await httpError.response.clone().json()

          if (body && typeof body === 'object') {
            // Check for field-specific validation errors
            // Backend may return: { fields: { email: 'validation.invalid_email' } }
            if ('fields' in body && typeof body.fields === 'object') {
              const fields = body.fields as Record<string, unknown>
              const serverErrors: ServerErrors = {}

              for (const [key, value] of Object.entries(fields)) {
                if (typeof value === 'string') {
                  serverErrors[key] = value
                }
              }

              return { fields: serverErrors }
            }

            // Check for single error message
            if ('detail' in body && typeof body.detail === 'string') {
              return { form: body.detail }
            }
          }

          // Fallback to generic HTTP error based on status
          const result = await parseProblemDetails(error)
          return { form: result.key }
        } catch {
          // Failed to parse response
          const result = await parseProblemDetails(error)
          return { form: result.key }
        }
      }

      // Fallback for non-HTTP errors
      const result = await parseProblemDetails(error)
      return { form: result.key }
    },
    [],
  )

  /**
   * Create a form-level validator that injects server errors
   *
   * Use this as the onSubmitAsync validator in your form to handle
   * server-side validation errors.
   *
   * @example
   * ```tsx
   * const form = useKyoraForm({
   *   validators: {
   *     onSubmitAsync: createServerErrorValidator(async ({ value }) => {
   *       return await loginMutation.mutateAsync(value)
   *     }),
   *   },
   * })
   * ```
   */
  const createServerErrorValidator = useCallback(
    (submitFn: (params: { value: any }) => Promise<void>) => {
      return async ({ value }: { value: any }) => {
        try {
          await submitFn({ value })
          return undefined // No errors
        } catch (error) {
          const serverErrors = await parseServerError(error)

          if (serverErrors.fields || serverErrors.form) {
            return {
              form: serverErrors.form,
              fields: serverErrors.fields,
            }
          }

          // Re-throw if not a validation error
          throw error
        }
      }
    },
    [parseServerError],
  )

  return {
    parseServerError,
    createServerErrorValidator,
  }
}

/**
 * Translate a server error key
 *
 * @param errorKey - Translation key from server
 * @param t - i18next translation function
 * @returns Translated error message
 */
export function translateServerError(
  errorKey: string | undefined,
  t: (key: string) => string,
): string | undefined {
  if (!errorKey) return undefined
  return t(errorKey)
}
