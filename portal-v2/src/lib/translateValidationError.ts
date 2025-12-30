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
export function translateValidationError(error: unknown, t: TFunction): string {
  if (!error) return ''

  // Handle error objects from TanStack Form/Zod
  let errorStr: string
  if (typeof error === 'string') {
    errorStr = error
  } else if (typeof error === 'object' && error !== null) {
    // Try to extract message from error object
    const errorObj = error as Record<string, unknown>
    if ('message' in errorObj && typeof errorObj.message === 'string') {
      errorStr = errorObj.message
    } else if (
      'toString' in errorObj &&
      typeof errorObj.toString === 'function'
    ) {
      errorStr = errorObj.toString()
    } else {
      errorStr = JSON.stringify(error)
    }
  } else {
    errorStr = String(error)
  }

  // Check if error is a translation key (starts with known namespace)
  if (errorStr.startsWith('validation.')) {
    // Try with 'common.validation.' prefix first
    const commonKey = `common.${errorStr}`
    if (t(commonKey) !== commonKey) {
      return t(commonKey)
    }
    // Fallback to original key
    return t(errorStr)
  }

  if (errorStr.startsWith('errors.') || errorStr.startsWith('onboarding.')) {
    return t(errorStr)
  }

  // Already translated or custom error message
  return errorStr
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
