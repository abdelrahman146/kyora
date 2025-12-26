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
 * Parses backend ProblemDetails error responses into user-friendly messages
 * Falls back to generic messages if parsing fails
 */
export async function parseProblemDetails(error: unknown): Promise<string> {
  if (!error) {
    return "An unexpected error occurred. Please try again.";
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

      // Prioritize detail, then title, then status code message
      if (
        body &&
        typeof body === "object" &&
        "detail" in body &&
        typeof body.detail === "string"
      ) {
        return body.detail;
      }

      if (
        body &&
        typeof body === "object" &&
        "title" in body &&
        typeof body.title === "string"
      ) {
        return body.title;
      }

      // Generate message from status code
      return getStatusMessage(httpError.response.status);
    } catch {
      // Response body is not JSON or failed to parse
      return getStatusMessage(httpError.response.status);
    }
  }

  // Handle TimeoutError
  if (error instanceof Error && error.name === "TimeoutError") {
    return "Request timed out. Please check your connection and try again.";
  }

  // Handle network errors
  if (error instanceof Error && error.name === "TypeError") {
    return "Network error. Please check your internet connection.";
  }

  // Handle generic Error objects
  if (error instanceof Error && error.message) {
    return error.message;
  }

  // Fallback for unknown error types
  return "An unexpected error occurred. Please try again.";
}

/**
 * Generates user-friendly messages for HTTP status codes
 */
function getStatusMessage(status: number): string {
  switch (status) {
    case 400:
      return "Invalid request. Please check your input and try again.";
    case 401:
      return "You are not authorized. Please log in again.";
    case 403:
      return "You don't have permission to perform this action.";
    case 404:
      return "The requested resource was not found.";
    case 409:
      return "This operation conflicts with existing data.";
    case 422:
      return "The provided data is invalid. Please check and try again.";
    case 429:
      return "Too many requests. Please wait a moment and try again.";
    case 500:
      return "Server error. Please try again later.";
    case 502:
      return "Service temporarily unavailable. Please try again later.";
    case 503:
      return "Service is currently under maintenance. Please try again later.";
    case 504:
      return "Request timed out. Please try again.";
    default:
      if (status >= 400 && status < 500) {
        return "Client error. Please check your request and try again.";
      }
      if (status >= 500) {
        return "Server error. Please try again later.";
      }
      return `Request failed with status ${String(status)}.`;
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
