import ky from "ky";
import { parseProblemDetails } from "../lib/errorParser";
import { getCookie, setCookie, deleteCookie } from "../lib/cookies";

const API_BASE_URL: string =
  (import.meta.env.VITE_API_BASE_URL as string | undefined) ??
  (import.meta.env.DEV ? window.location.origin : "http://localhost:8080");

const REFRESH_TOKEN_COOKIE_NAME = "kyora_refresh_token";

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

async function refreshAccessToken(): Promise<string | null> {
  const refresh = getRefreshToken();

  if (!refresh) {
    return null;
  }

  try {
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

// API Client Configuration

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
    // Before Request Hook: Add Authentication & Headers
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
      },
    ],

    // After Response Hook: Handle Authentication & Token Refresh
    afterResponse: [
      async (request, options, response) => {
        // Handle 401 Unauthorized - Try to refresh token
        // IMPORTANT: Auth endpoints intentionally return 401 for invalid credentials.
        // We must not treat those as "session expired" or we'll cause redirects/retries
        // on the login page.
        if (response.status === 401) {
          let pathname = "";
          try {
            pathname = new URL(request.url).pathname;
          } catch {
            // ignore
          }

          if (pathname.startsWith("/v1/auth/")) {
            return response;
          }

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

    // Before Retry Hook
    beforeRetry: [],

    // Before Error Hook: Parse ProblemDetails and enhance error messages
    // Note: Error messages are now translation keys. Use translateError()
    // in components to get localized messages.
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

        return error;
      },
    ],
  },
});

// Typed API Client Methods (Optional Convenience Wrappers)

type KyOptions = Parameters<typeof apiClient.get>[1];

const inFlight = new Map<string, Promise<unknown>>();

function stableStringify(value: unknown): string {
  if (value === null) return "null";
  if (value === undefined) return "undefined";

  const t = typeof value;
  if (t === "string") return JSON.stringify(value);
  if (t === "number" || t === "boolean") return JSON.stringify(value);
  if (t === "function") return '"[function]"';

  if (Array.isArray(value)) {
    return `[${value.map(stableStringify).join(",")}]`;
  }

  if (value instanceof URLSearchParams) {
    const entries = Array.from(value.entries()).sort(([a], [b]) =>
      a.localeCompare(b)
    );
    return `URLSearchParams(${stableStringify(entries)})`;
  }

  if (value instanceof Headers) {
    const entries = Array.from(value.entries()).sort(([a], [b]) =>
      a.localeCompare(b)
    );
    return `Headers(${stableStringify(entries)})`;
  }

  if (t === "object") {
    const obj = value as Record<string, unknown>;
    const keys = Object.keys(obj).sort();
    const parts = keys.map(
      (k) => `${JSON.stringify(k)}:${stableStringify(obj[k])}`
    );
    return `{${parts.join(",")}}`;
  }

  try {
    return JSON.stringify(value);
  } catch {
    return "[unstringifiable]";
  }
}

function buildDedupeKey(
  method: string,
  url: string,
  options?: KyOptions
): string {
  const token = getAccessToken() ?? "";

  const searchParams = (options as { searchParams?: unknown } | undefined)
    ?.searchParams;
  const jsonPayload = (options as { json?: unknown } | undefined)?.json;
  const body = (options as { body?: unknown } | undefined)?.body;

  return [
    method.toUpperCase(),
    url,
    `token:${token}`,
    searchParams !== undefined ? `sp:${stableStringify(searchParams)}` : "",
    jsonPayload !== undefined ? `json:${stableStringify(jsonPayload)}` : "",
    body !== undefined ? `body:${stableStringify(body)}` : "",
  ]
    .filter(Boolean)
    .join("|");
}

async function deduped<T>(
  method: string,
  url: string,
  options: KyOptions | undefined,
  fn: () => Promise<T>
): Promise<T> {
  const key = buildDedupeKey(method, url, options);
  const existing = inFlight.get(key);
  if (existing) {
    return existing as Promise<T>;
  }

  const promise = fn().finally(() => {
    inFlight.delete(key);
  });
  inFlight.set(key, promise as Promise<unknown>);
  return promise;
}

/**
 * Type-safe GET request
 */
export async function get<T>(
  url: string,
  options?: Parameters<typeof apiClient.get>[1]
): Promise<T> {
  return deduped("get", url, options, () =>
    apiClient.get(url, options).json<T>()
  );
}

/**
 * Type-safe POST request
 */
export async function post<T>(
  url: string,
  options?: Parameters<typeof apiClient.post>[1]
): Promise<T> {
  return deduped("post", url, options, () =>
    apiClient.post(url, options).json<T>()
  );
}

/**
 * POST request with no JSON body expected (e.g., 204 No Content)
 */
export async function postVoid(
  url: string,
  options?: Parameters<typeof apiClient.post>[1]
): Promise<void> {
  await deduped("post", url, options, async () => {
    await apiClient.post(url, options);
  });
}

/**
 * Type-safe PUT request
 */
export async function put<T>(
  url: string,
  options?: Parameters<typeof apiClient.put>[1]
): Promise<T> {
  return deduped("put", url, options, () =>
    apiClient.put(url, options).json<T>()
  );
}

/**
 * Type-safe PATCH request
 */
export async function patch<T>(
  url: string,
  options?: Parameters<typeof apiClient.patch>[1]
): Promise<T> {
  return deduped("patch", url, options, () =>
    apiClient.patch(url, options).json<T>()
  );
}

/**
 * Type-safe DELETE request
 */
export async function del<T>(
  url: string,
  options?: Parameters<typeof apiClient.delete>[1]
): Promise<T> {
  return deduped("delete", url, options, () =>
    apiClient.delete(url, options).json<T>()
  );
}

/**
 * DELETE request with no JSON body expected
 */
export async function delVoid(
  url: string,
  options?: Parameters<typeof apiClient.delete>[1]
): Promise<void> {
  await deduped("delete", url, options, async () => {
    await apiClient.delete(url, options);
  });
}

export default apiClient;
