import type { HTTPError } from "ky";

/**
 * ProblemDetails interface based on RFC 7807
 * Matches the backend problem.Problem type from swagger.json
 */
export interface ProblemDetails {
  type?: string;
  title?: string;
  status?: number;
  detail?: string;
  instance?: string;
  extensions?: Record<string, unknown>;
}

/**
 * Result of parsing an error with i18n translation key and optional interpolation params
 */
export interface ErrorResult {
  key: string;
  params?: Record<string, string | number>;
  fallback?: string;
}

/**
 * Parses backend ProblemDetails error responses into translation keys
 * Returns an ErrorResult with the i18n key and optional interpolation params
 * Falls back to generic error keys if parsing fails
 */
export async function parseProblemDetails(
  error: unknown
): Promise<ErrorResult> {
  if (!error) {
    return { key: "errors.generic.unexpected" };
  }

  // Handle HTTPError from ky
  if (
    typeof error === "object" &&
    "response" in error &&
    error.response instanceof Response
  ) {
    const httpError = error as HTTPError;

    try {
      // Try to parse JSON response body
      const body = await httpError.response.json();

      // If backend provides a detail message, use it as fallback but return status-based key
      const fallback =
        body &&
        typeof body === "object" &&
        "detail" in body &&
        typeof body.detail === "string"
          ? body.detail
          : undefined;

      // Return translation key based on status code
      return {
        ...getStatusErrorKey(httpError.response.status),
        fallback,
      };
    } catch {
      // Response body is not JSON or failed to parse
      return getStatusErrorKey(httpError.response.status);
    }
  }

  // Handle TimeoutError
  if (error instanceof Error && error.name === "TimeoutError") {
    return { key: "errors.network.timeout" };
  }

  // Handle network errors
  if (error instanceof Error && error.name === "TypeError") {
    return { key: "errors.network.connection" };
  }

  // Handle generic Error objects
  if (error instanceof Error && error.message) {
    return {
      key: "errors.generic.message",
      params: { message: error.message },
      fallback: error.message,
    };
  }

  // Fallback for unknown error types
  return { key: "errors.generic.unexpected" };
}

/**
 * Returns translation key and params for HTTP status codes
 */
function getStatusErrorKey(status: number): ErrorResult {
  switch (status) {
    case 400:
      return { key: "errors.http.400" };
    case 401:
      return { key: "errors.http.401" };
    case 403:
      return { key: "errors.http.403" };
    case 404:
      return { key: "errors.http.404" };
    case 409:
      return { key: "errors.http.409" };
    case 422:
      return { key: "errors.http.422" };
    case 429:
      return { key: "errors.http.429" };
    case 500:
      return { key: "errors.http.500" };
    case 502:
      return { key: "errors.http.502" };
    case 503:
      return { key: "errors.http.503" };
    case 504:
      return { key: "errors.http.504" };
    default:
      if (status >= 400 && status < 500) {
        return { key: "errors.http.4xx", params: { status } };
      }
      if (status >= 500) {
        return { key: "errors.http.5xx", params: { status } };
      }
      return { key: "errors.http.unknown", params: { status } };
  }
}

/**
 * Extracts validation errors from ProblemDetails extensions
 * Returns a map of field names to error messages
 */
export async function parseValidationErrors(
  error: unknown
): Promise<Record<string, string> | null> {
  if (!error) {
    return null;
  }

  // Handle HTTPError from ky
  if (
    typeof error === "object" &&
    "response" in error &&
    error.response instanceof Response
  ) {
    const httpError = error as HTTPError;

    try {
      const body = await httpError.response.json();

      // Check if extensions contains validation errors
      if (
        body &&
        typeof body === "object" &&
        "extensions" in body &&
        body.extensions &&
        typeof body.extensions === "object"
      ) {
        // Common patterns: errors, validationErrors, fieldErrors
        const validationData =
          (body.extensions as Record<string, unknown>).errors ??
          (body.extensions as Record<string, unknown>).validationErrors ??
          (body.extensions as Record<string, unknown>).fieldErrors;

        if (validationData && typeof validationData === "object") {
          return validationData as Record<string, string>;
        }
      }
    } catch {
      // Failed to parse validation errors
    }
  }

  return null;
}

/**
 * Checks if an error is a specific HTTP status code
 */
export function isHttpError(error: unknown, status?: number): boolean {
  if (
    typeof error === "object" &&
    error !== null &&
    "response" in error &&
    error.response instanceof Response
  ) {
    if (status !== undefined) {
      return error.response.status === status;
    }
    return true;
  }
  return false;
}

/**
 * Type guard for HTTPError
 */
export function isHTTPError(error: unknown): error is HTTPError {
  return (
    typeof error === "object" &&
    error !== null &&
    "response" in error &&
    error.response instanceof Response
  );
}
