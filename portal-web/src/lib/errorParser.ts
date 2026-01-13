import type { HTTPError } from 'ky'

/**
 * ProblemDetails interface based on RFC 7807
 * Matches the backend problem.Problem type from swagger.json
 *
 * All backend errors now include `extensions.code` for machine-readable error mapping.
 * Example: { extensions: { code: "account.invalid_credentials" } }
 */
export interface ProblemDetails {
  type?: string
  title?: string
  status?: number
  detail?: string
  instance?: string
  extensions?: Record<string, unknown>
}

/**
 * Result of parsing an error with i18n translation key and optional interpolation params
 */
export interface ErrorResult {
  key: string
  ns?: string
  params?: Record<string, string | number | Record<string, unknown>>
  fallback?: string
}

/**
 * Parses backend ProblemDetails error responses into translation keys
 * Returns an ErrorResult with the i18n key and optional interpolation params
 * Falls back to generic error keys if parsing fails
 *
 * Priority:
 * 1. Use extensions.code if available → "errors.backend.<code>"
 * 2. Fall back to status-based keys for graceful degradation
 */
export async function parseProblemDetails(
  error: unknown,
): Promise<ErrorResult> {
  if (!error) {
    return { ns: 'errors', key: 'generic.unexpected' }
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

      // If backend provides a detail message, use it as fallback.
      const fallback =
        body &&
        typeof body === 'object' &&
        'detail' in body &&
        typeof (body as { detail?: unknown }).detail === 'string'
          ? (body as { detail: string }).detail
          : undefined

      // Priority 1: Check for extensions.code (machine-readable error code)
      if (
        body &&
        typeof body === 'object' &&
        'extensions' in body &&
        (body as { extensions?: unknown }).extensions &&
        typeof (body as { extensions: unknown }).extensions === 'object'
      ) {
        const extensions = (body as { extensions: Record<string, unknown> })
          .extensions
        if (extensions.code && typeof extensions.code === 'string') {
          // Extract all extension properties except 'code' as details for translation
          const { code, ...details } = extensions

          // Return "errors.backend.<code>" for i18n lookup
          // Include details in params for flexible translations like:
          // "you cannot update order status from {details.orderStatusBefore} to {details.orderStatusAfter}"
          return {
            ns: 'errors',
            key: `backend.${code}`,
            params: Object.keys(details).length > 0 ? { details } : undefined,
            fallback,
          }
        }
      }

      // Priority 2: Fall back to status code mapping for graceful degradation
      return {
        ...getStatusErrorKey(httpError.response.status),
        fallback,
      }
    } catch {
      // Response body is not JSON or failed to parse
      return getStatusErrorKey(httpError.response.status)
    }
  }

  // Handle TimeoutError
  if (error instanceof Error && error.name === 'TimeoutError') {
    return { ns: 'errors', key: 'network.timeout' }
  }

  // Handle network errors
  if (error instanceof Error && error.name === 'TypeError') {
    return { ns: 'errors', key: 'network.connection' }
  }

  // Handle generic Error objects
  if (error instanceof Error && error.message) {
    return {
      ns: 'errors',
      key: 'generic.message',
      params: { message: error.message },
      fallback: error.message,
    }
  }

  // Fallback for unknown error types
  return { ns: 'errors', key: 'generic.unexpected' }
}

/**
 * Returns translation key and params for HTTP status codes
 */
function getStatusErrorKey(status: number): ErrorResult {
  switch (status) {
    case 400:
      return { ns: 'errors', key: 'http.400' }
    case 401:
      return { ns: 'errors', key: 'http.401' }
    case 403:
      return { ns: 'errors', key: 'http.403' }
    case 404:
      return { ns: 'errors', key: 'http.404' }
    case 409:
      return { ns: 'errors', key: 'http.409' }
    case 422:
      return { ns: 'errors', key: 'http.422' }
    case 429:
      return { ns: 'errors', key: 'http.429' }
    case 500:
      return { ns: 'errors', key: 'http.500' }
    case 502:
      return { ns: 'errors', key: 'http.502' }
    case 503:
      return { ns: 'errors', key: 'http.503' }
    case 504:
      return { ns: 'errors', key: 'http.504' }
    default:
      if (status >= 400 && status < 500) {
        return { ns: 'errors', key: 'http.4xx', params: { status } }
      }
      if (status >= 500) {
        return { ns: 'errors', key: 'http.5xx', params: { status } }
      }
      return { ns: 'errors', key: 'http.unknown', params: { status } }
  }
}

/**
 * Extracts validation errors from ProblemDetails extensions
 * Returns a map of field names to error messages
 */
export async function parseValidationErrors(
  error: unknown,
): Promise<Record<string, string> | null> {
  if (!error) {
    return null
  }

  // Handle HTTPError from ky
  if (
    typeof error === 'object' &&
    'response' in error &&
    error.response instanceof Response
  ) {
    const httpError = error as HTTPError

    try {
      const body: unknown = await httpError.response.clone().json()

      // Check if extensions contains validation errors
      if (
        body &&
        typeof body === 'object' &&
        'extensions' in body &&
        (body as { extensions?: unknown }).extensions &&
        typeof (body as { extensions: unknown }).extensions === 'object'
      ) {
        // Common patterns: errors, validationErrors, fieldErrors
        const extensions = (body as { extensions: Record<string, unknown> })
          .extensions
        const validationData =
          extensions.errors ??
          extensions.validationErrors ??
          extensions.fieldErrors

        if (validationData && typeof validationData === 'object') {
          return validationData as Record<string, string>
        }
      }
    } catch {
      // Failed to parse validation errors
    }
  }

  return null
}

/**
 * Checks if an error is a specific HTTP status code
 */
export function isHttpError(error: unknown, status?: number): boolean {
  if (
    typeof error === 'object' &&
    error !== null &&
    'response' in error &&
    error.response instanceof Response
  ) {
    if (status !== undefined) {
      return error.response.status === status
    }
    return true
  }
  return false
}

/**
 * Type guard for HTTPError
 */
export function isHTTPError(error: unknown): error is HTTPError {
  return (
    typeof error === 'object' &&
    error !== null &&
    'response' in error &&
    error.response instanceof Response
  )
}
