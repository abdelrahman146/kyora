/**
 * RFC 7807 ProblemDetails Error Parser
 *
 * Parses backend error responses that follow the RFC 7807 Problem Details format.
 * This is a stub that will be fully implemented in step 3.
 */

export interface ErrorResult {
  key: string
  params?: Record<string, unknown>
  fallback?: string
}

export function parseProblemDetails(error: unknown): ErrorResult {
  // Stub implementation - will be completed in step 3
  if (error instanceof Error) {
    return {
      key: 'errors.http.unknown',
      fallback: error.message,
    }
  }

  return {
    key: 'errors.http.unknown',
    fallback: 'An unknown error occurred',
  }
}

export function parseValidationErrors(
  _error: unknown,
): Record<string, string> | null {
  // Stub implementation - will be completed in step 3
  return null
}
