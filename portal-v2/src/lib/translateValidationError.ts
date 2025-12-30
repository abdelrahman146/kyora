import type { TFunction } from 'i18next'

/**
 * Translate validation error messages from Zod schemas
 *
 * This utility translates validation error keys (like 'validation.required')
 * into user-friendly localized messages.
 *
 * @param error - The error message (may be a translation key or already translated)
 * @param t - i18next translation function
 * @returns Translated error message
 *
 * @example
 * ```ts
 * // In Zod schema
 * const schema = z.object({
 *   email: z.string().min(1, 'validation.required').email('validation.invalid_email'),
 * })
 *
 * // In component
 * const errorMessage = translateValidationError(field.state.meta.errors[0], t)
 * ```
 */
export function translateValidationError(
  error: string | undefined,
  t: TFunction,
): string {
  if (!error) return ''

  // Check if error is a translation key (starts with known namespace)
  if (
    error.startsWith('validation.') ||
    error.startsWith('errors.') ||
    error.startsWith('onboarding.')
  ) {
    return t(error)
  }

  // Already translated or custom error message
  return error
}

/**
 * Translate an array of validation errors
 *
 * @param errors - Array of error messages
 * @param t - i18next translation function
 * @returns First translated error message, or empty string
 */
export function translateValidationErrors(
  errors: Array<string> | undefined,
  t: TFunction,
): string {
  if (!errors || errors.length === 0) return ''
  return translateValidationError(errors[0], t)
}
