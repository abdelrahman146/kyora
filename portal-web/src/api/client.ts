import ky from "ky";
import { parseProblemDetails } from "../lib/errorParser";
import { getCookie, setCookie, deleteCookie } from "../lib/cookies";

/**
 * Centralized API Client for Kyora Portal Web
 *
 * Features:
 * - JWT Bearer token authentication (Access token in memory, Refresh token in secure cookie)
 * - Automatic token refresh on 401 responses
 * - Retry logic with exponential backoff
 * - Comprehensive error handling with ProblemDetails parsing
 * - Request/response logging in development
 */

const API_BASE_URL: string =
  (import.meta.env.VITE_API_BASE_URL as string | undefined) ??
  (import.meta.env.DEV ? window.location.origin : "http://localhost:8080");

const REFRESH_TOKEN_COOKIE_NAME = "kyora_refresh_token";

// ============================================================================
// Token Management
// Access Token: In-Memory (more secure, cleared on tab close/refresh)
// Refresh Token: Secure Cookie (persistent, HttpOnly-like with Secure flag in production)
// ============================================================================

let accessToken: string | null = null;
let isRefreshing = false;
let refreshPromise: Promise<string | null> | null = null;

export function getAccessToken(): string | null {
  return accessToken;
}

export function getRefreshToken(): string | null {
  return getCookie(REFRESH_TOKEN_COOKIE_NAME);
}

export function setTokens(access: string, refresh: string): void {
  accessToken = access;
  // Store refresh token in secure cookie (365 days)
  // Note: For true HttpOnly, this should be set by backend via Set-Cookie header
  // This is client-side cookie with Secure flag in production
  setCookie(REFRESH_TOKEN_COOKIE_NAME, refresh, 365);
}

export function clearTokens(): void {
  accessToken = null;
  deleteCookie(REFRESH_TOKEN_COOKIE_NAME);
  isRefreshing = false;
  refreshPromise = null;
}

export function hasValidToken(): boolean {
  return accessToken !== null;
}

// ============================================================================
// Token Refresh Logic
// ============================================================================

async function refreshAccessToken(): Promise<string | null> {
  const refresh = getRefreshToken();

  if (!refresh) {
    return null;
  }

  try {
    // Create a new ky instance without auth hooks to avoid infinite loop
    const refreshClient = ky.create({
      prefixUrl: API_BASE_URL,
      timeout: 10000,
    });

    const response: { token: string; refreshToken: string } =
      await refreshClient
        .post("v1/auth/refresh", {
          json: { refreshToken: refresh },
        })
        .json<{ token: string; refreshToken: string }>();

    // Update tokens
    setTokens(response.token, response.refreshToken);
    return response.token;
  } catch {
    // Refresh failed - clear tokens and redirect to login
    clearTokens();
    return null;
  }
}

// ============================================================================
// API Client Configuration
// ============================================================================

export const apiClient = ky.create({
  prefixUrl: API_BASE_URL,
  timeout: 30000, // 30 seconds default timeout

  // Retry configuration with exponential backoff
  retry: {
    limit: 2,
    methods: ["get", "put", "head", "delete", "options", "trace"],
    statusCodes: [408, 413, 429, 500, 502, 503, 504],
    backoffLimit: 3000, // Max 3 seconds between retries
  },

  hooks: {
    // ========================================================================
    // Before Request Hook: Add Authentication & Headers
    // ========================================================================
    beforeRequest: [
      (request) => {
        // Add JWT Bearer token if available
        const token = getAccessToken();
        if (token) {
          request.headers.set("Authorization", `Bearer ${token}`);
        }

        // Add request ID for tracing (optional but useful for debugging)
        if (!request.headers.has("X-Request-ID")) {
          request.headers.set("X-Request-ID", crypto.randomUUID());
        }

        // Log request in development
        if (import.meta.env.DEV) {
          console.log(`[API] → ${request.method} ${request.url}`);
        }
      },
    ],

    // ========================================================================
    // After Response Hook: Handle Authentication & Token Refresh
    // ========================================================================
    afterResponse: [
      async (request, options, response) => {
        // Log response in development
        if (import.meta.env.DEV) {
          console.log(
            `[API] ← ${request.method} ${request.url} - ${String(
              response.status
            )}`
          );
        }

        // Handle 401 Unauthorized - Try to refresh token
        if (response.status === 401) {
          // Prevent multiple simultaneous refresh requests
          if (!isRefreshing) {
            isRefreshing = true;
            refreshPromise = refreshAccessToken().finally(() => {
              isRefreshing = false;
              refreshPromise = null;
            });
          }

          try {
            // Wait for token refresh to complete
            const refreshPromiseResolved = refreshPromise;
            if (!refreshPromiseResolved) {
              throw new Error("No refresh promise available");
            }
            const newToken = await refreshPromiseResolved;

            if (newToken) {
              // Retry original request with new token
              const headers = new Headers(request.headers);
              headers.set("Authorization", `Bearer ${newToken}`);

              return await ky(request, { ...options, headers });
            } else {
              // Refresh failed - redirect to login
              window.location.href = "/login";
              throw new Error("Authentication required");
            }
          } catch (err) {
            // Refresh failed - redirect to login
            clearTokens();
            window.location.href = "/login";
            throw err;
          }
        }

        return response;
      },
    ],

    // ========================================================================
    // Before Retry Hook: Log retry attempts
    // ========================================================================
    beforeRetry: [
      ({ request, retryCount }) => {
        if (import.meta.env.DEV) {
          console.log(
            `[API] ⟳ Retry #${String(retryCount)} for ${request.method} ${
              request.url
            }`
          );
        }
      },
    ],

    // ========================================================================
    // Before Error Hook: Parse ProblemDetails and enhance error messages
    // Note: Error messages are now translation keys. Use translateError()
    // in components to get localized messages.
    // ========================================================================
    beforeError: [
      async (error) => {
        // Parse backend ProblemDetails error format
        try {
          const errorResult = await parseProblemDetails(error);
          // Store translation key in error message for now
          // Components should use translateError(errorResult, t) for proper i18n
          error.message = errorResult.fallback ?? errorResult.key;
        } catch {
          // If parsing fails, keep original error message
        }

        // Log error in development
        if (import.meta.env.DEV) {
          console.error(`[API] ✗ ${error.message}`, error);
        }

        return error;
      },
    ],
  },
});

// ============================================================================
// Typed API Client Methods (Optional Convenience Wrappers)
// ============================================================================

/**
 * Type-safe GET request
 */
export async function get<T>(
  url: string,
  options?: Parameters<typeof apiClient.get>[1]
): Promise<T> {
  return apiClient.get(url, options).json<T>();
}

/**
 * Type-safe POST request
 */
export async function post<T>(
  url: string,
  options?: Parameters<typeof apiClient.post>[1]
): Promise<T> {
  return apiClient.post(url, options).json<T>();
}

/**
 * Type-safe PUT request
 */
export async function put<T>(
  url: string,
  options?: Parameters<typeof apiClient.put>[1]
): Promise<T> {
  return apiClient.put(url, options).json<T>();
}

/**
 * Type-safe PATCH request
 */
export async function patch<T>(
  url: string,
  options?: Parameters<typeof apiClient.patch>[1]
): Promise<T> {
  return apiClient.patch(url, options).json<T>();
}

/**
 * Type-safe DELETE request
 */
export async function del<T>(
  url: string,
  options?: Parameters<typeof apiClient.delete>[1]
): Promise<T> {
  return apiClient.delete(url, options).json<T>();
}

export default apiClient;
