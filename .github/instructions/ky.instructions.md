---
description: Ky HTTP Client Library - Production-Grade HTTP Requests
applyTo: "portal-web/**"
---

# Ky HTTP Client - Production-Grade HTTP Requests for Kyora Portal

Ky is a tiny, elegant HTTP client based on the Fetch API. This guide covers how to leverage its strengths for production-grade HTTP requests in the Kyora Portal Web App.

## Why Ky Over Plain Fetch?

1. **Simpler API**: Method shortcuts (`ky.post()`, `ky.get()`)
2. **Automatic Error Handling**: Treats non-2xx status codes as errors (after redirects)
3. **Built-in Retries**: Configurable retry logic with exponential backoff
4. **JSON Shortcuts**: `.json()` method with TypeScript generics support
5. **Timeout Support**: Built-in request timeout handling
6. **URL Prefix**: Set base URLs per instance
7. **Instances with Defaults**: Create custom instances with predefined options
8. **Hooks**: Powerful lifecycle hooks for intercepting requests/responses
9. **TypeScript First**: `.json()` defaults to `unknown`, not `any`

## Core Setup

### Basic API Client Structure

Create a centralized API client in `src/api/client.ts`:

```ts
import ky from "ky";
import type { KyInstance } from "ky";

// Base API client with default configuration
export const apiClient = ky.create({
  prefixUrl: import.meta.env.VITE_API_URL || "http://localhost:8080",
  timeout: 30000, // 30 seconds
  retry: {
    limit: 2,
    methods: ["get", "put", "head", "delete", "options", "trace"],
    statusCodes: [408, 413, 429, 500, 502, 503, 504],
    backoffLimit: 3000, // Max 3 seconds between retries
  },
  headers: {
    "Content-Type": "application/json",
  },
});
```

### Typed Responses with Generics

Always use TypeScript generics for type-safe responses:

```ts
// Define your API response types
interface User {
  id: string;
  email: string;
  name: string;
}

interface ApiResponse<T> {
  data: T;
  meta?: {
    page: number;
    total: number;
  };
}

// Use generics for type safety
const user = await apiClient.get<User>("api/v1/users/123").json<User>();
// user is typed as User, not any or unknown

// For paginated responses
const response = await apiClient
  .get<ApiResponse<User[]>>("api/v1/users")
  .json<ApiResponse<User[]>>();
```

## Authentication Implementation

### JWT with Refresh Token Pattern

```ts
import ky, { HTTPError } from "ky";

let accessToken: string | null = null;
let refreshToken: string | null = null;
let isRefreshing = false;
let refreshPromise: Promise<string> | null = null;

// Get stored tokens (from memory or secure storage)
function getAccessToken(): string | null {
  return accessToken;
}

function getRefreshToken(): string | null {
  return refreshToken;
}

function setTokens(access: string, refresh: string): void {
  accessToken = access;
  refreshToken = refresh;
}

function clearTokens(): void {
  accessToken = null;
  refreshToken = null;
}

// Refresh access token
async function refreshAccessToken(): Promise<string> {
  const refresh = getRefreshToken();

  if (!refresh) {
    throw new Error("No refresh token available");
  }

  const response = await ky
    .post("/api/v1/auth/refresh", {
      json: { refresh_token: refresh },
    })
    .json<{ access_token: string; refresh_token: string }>();

  setTokens(response.access_token, response.refresh_token);
  return response.access_token;
}

// Create authenticated API client
export const apiClient = ky.create({
  prefixUrl: import.meta.env.VITE_API_URL,
  timeout: 30000,
  hooks: {
    beforeRequest: [
      (request) => {
        const token = getAccessToken();
        if (token) {
          request.headers.set("Authorization", `Bearer ${token}`);
        }
      },
    ],
    afterResponse: [
      async (request, options, response) => {
        // Handle 401 - Unauthorized (token expired)
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
            // Wait for token refresh
            const newToken = await refreshPromise!;

            // Retry original request with new token
            const headers = new Headers(request.headers);
            headers.set("Authorization", `Bearer ${newToken}`);

            return ky(request, { ...options, headers });
          } catch (error) {
            // Refresh failed, clear tokens and redirect to login
            clearTokens();
            window.location.href = "/login";
            throw error;
          }
        }
      },
    ],
    beforeError: [
      async (error) => {
        const { response } = error;

        if (response) {
          // Parse backend error format (ProblemDetails)
          try {
            const body = await response.json();
            error.message = body.detail || body.message || error.message;
          } catch {
            // Response is not JSON, use default message
          }
        }

        return error;
      },
    ],
  },
});
```

## Advanced Retry Strategies

### Custom Retry Logic

```ts
export const apiClient = ky.create({
  prefixUrl: import.meta.env.VITE_API_URL,
  retry: {
    limit: 3,
    methods: ["get", "put", "head", "delete", "options", "trace"],
    statusCodes: [408, 413, 429, 500, 502, 503, 504],

    // Custom retry delay with exponential backoff
    // Default: 0.3 * (2 ** (attemptCount - 1)) * 1000
    delay: (attemptCount) => {
      return Math.min(1000 * 2 ** attemptCount, 10000); // Max 10 seconds
    },

    // Add jitter to prevent thundering herd
    jitter: true, // Randomizes delay between 0 and computed value

    // Retry on timeout
    retryOnTimeout: true,

    // Custom retry logic
    shouldRetry: ({ error, retryCount }) => {
      // Retry on specific business logic errors
      if (error instanceof HTTPError) {
        const status = error.response.status;

        // Always retry rate limits (429) but only for first 2 attempts
        if (status === 429 && retryCount <= 2) {
          return true;
        }

        // Don't retry on 4xx client errors except rate limits
        if (status >= 400 && status < 500) {
          return false;
        }
      }

      // Use default retry logic for other errors
      return undefined;
    },
  },
});
```

### Respecting Server Retry-After Headers

Ky automatically respects `Retry-After` headers for status codes 413, 429, and 503:

```ts
export const apiClient = ky.create({
  retry: {
    limit: 5,
    afterStatusCodes: [413, 429, 503], // Wait for Retry-After on these
    maxRetryAfter: 60000, // Max 60 seconds wait
  },
});
```

## Hook System - Request/Response Interceptors

### beforeRequest Hook

Modify requests before they are sent:

```ts
export const apiClient = ky.create({
  hooks: {
    beforeRequest: [
      (request, options, { retryCount }) => {
        // Add request ID for tracing
        request.headers.set("X-Request-ID", crypto.randomUUID());

        // Add business ID from context (if available)
        const businessId = getCurrentBusinessId();
        if (businessId) {
          request.headers.set("X-Business-ID", businessId);
        }

        // Log only on initial request, not retries
        if (retryCount === 0) {
          console.log(`[API] ${request.method} ${request.url}`);
        }
      },
    ],
  },
});
```

### afterResponse Hook

Process responses before they reach your code:

```ts
export const apiClient = ky.create({
  hooks: {
    afterResponse: [
      async (request, options, response, { retryCount }) => {
        // Log response time
        const duration = Date.now() - request.startTime;
        console.log(
          `[API] ${request.method} ${request.url} - ${response.status} (${duration}ms)`
        );

        // Handle specific business logic based on response
        if (response.status === 200) {
          const data = await response.clone().json();

          // Force retry based on response body content
          if (data.error?.code === "TEMPORARY_ERROR") {
            return ky.retry({ code: "TEMPORARY_ERROR" });
          }

          // Retry with custom delay from API response
          if (data.error?.code === "RATE_LIMIT") {
            return ky.retry({
              delay: data.error.retryAfter * 1000,
              code: "RATE_LIMIT",
            });
          }
        }
      },
    ],
  },
});
```

### beforeRetry Hook

Modify requests before retry attempts:

```ts
export const apiClient = ky.create({
  hooks: {
    beforeRetry: [
      async ({ request, options, error, retryCount }) => {
        console.log(`[API] Retry #${retryCount} for ${request.url}`);

        // Refresh token on 401 errors
        if (error instanceof HTTPError && error.response.status === 401) {
          const newToken = await refreshAccessToken();
          request.headers.set("Authorization", `Bearer ${newToken}`);
        }

        // Add backoff information to request
        request.headers.set("X-Retry-Count", String(retryCount));
      },
    ],
  },
});
```

### beforeError Hook

Enhance error messages before they are thrown:

```ts
export const apiClient = ky.create({
  hooks: {
    beforeError: [
      async (error, { retryCount }) => {
        const { response } = error;

        if (response) {
          // Parse backend error format
          try {
            const body = await response.json();

            // Enhance error with backend details
            error.name = body.type || "APIError";
            error.message = body.detail || body.message || error.message;

            // Add custom properties
            (error as any).code = body.code;
            (error as any).validationErrors = body.errors;
          } catch {
            // Response is not JSON
          }

          // Show different message on final retry
          if (retryCount === error.options.retry.limit) {
            error.message = `${error.message} (failed after ${retryCount} retries)`;
          }
        }

        return error;
      },
    ],
  },
});
```

## Timeout Management

```ts
// Global timeout
export const apiClient = ky.create({
  timeout: 30000, // 30 seconds for all requests
});

// Per-request timeout override
const data = await apiClient
  .get("api/v1/large-report", {
    timeout: 120000, // 2 minutes for large reports
  })
  .json();

// Disable timeout for specific requests
const stream = await apiClient.get("api/v1/stream", {
  timeout: false, // No timeout
});

// Retry on timeout
export const apiClient = ky.create({
  timeout: 10000,
  retry: {
    limit: 3,
    retryOnTimeout: true, // Retry when request times out
  },
});
```

## JSON Handling

### Sending JSON

```ts
// Automatic JSON serialization
const response = await apiClient
  .post("api/v1/users", {
    json: {
      name: "Ahmed",
      email: "ahmed@example.com",
    },
  })
  .json<User>();

// ky automatically:
// - Sets Content-Type: application/json
// - Stringifies the object
// - Sets Accept: application/json
```

### Custom JSON Parsing

```ts
import { parse as parseBourne } from "@hapi/bourne";
import { parse as parseDate } from "date-fns";

export const apiClient = ky.create({
  // Protect from prototype pollution
  parseJson: (text) => parseBourne(text),

  // Custom stringification (e.g., handle dates)
  stringifyJson: (data) =>
    JSON.stringify(data, (key, value) => {
      // Convert DateTime objects to ISO strings
      if (value instanceof Date) {
        return value.toISOString();
      }
      return value;
    }),
});
```

## Search Parameters

```ts
// Using object
const users = await apiClient
  .get("api/v1/users", {
    searchParams: {
      page: 1,
      limit: 20,
      status: "active",
      sort: "created_at",
    },
  })
  .json<User[]>();
// => GET /api/v1/users?page=1&limit=20&status=active&sort=created_at

// Using URLSearchParams
const params = new URLSearchParams();
params.set("q", "search query");
params.set("filter", "active");

const results = await apiClient
  .get("api/v1/search", {
    searchParams: params,
  })
  .json<SearchResult[]>();

// undefined values are filtered out
const users = await apiClient
  .get("api/v1/users", {
    searchParams: {
      page: 1,
      filter: undefined, // This will be omitted
    },
  })
  .json<User[]>();
// => GET /api/v1/users?page=1
```

## Creating Domain-Specific Instances

Create specialized instances for different API domains:

```ts
// Base authenticated client
const apiClient = ky.create({
  prefixUrl: import.meta.env.VITE_API_URL,
  timeout: 30000,
  hooks: {
    beforeRequest: [authHook],
  },
});

// Users API
export const usersApi = apiClient.extend({
  prefixUrl: `${apiClient.defaults.options.prefixUrl}/api/v1/users`,
});

// Orders API with longer timeout
export const ordersApi = apiClient.extend({
  prefixUrl: `${apiClient.defaults.options.prefixUrl}/api/v1/orders`,
  timeout: 60000, // Orders may take longer
});

// Public API (no auth)
export const publicApi = ky.create({
  prefixUrl: import.meta.env.VITE_API_URL,
  timeout: 10000,
});

// Usage
const user = await usersApi.get("123").json<User>();
// => GET /api/v1/users/123

const order = await ordersApi.post("", { json: orderData }).json<Order>();
// => POST /api/v1/orders
```

## Error Handling Best Practices

### Type Guards for Error Types

```ts
import { HTTPError, TimeoutError } from "ky";

export function isHTTPError(error: unknown): error is HTTPError {
  return error instanceof HTTPError;
}

export function isTimeoutError(error: unknown): error is TimeoutError {
  return error instanceof TimeoutError;
}

// Usage
try {
  const data = await apiClient.get("api/v1/users").json();
} catch (error) {
  if (isHTTPError(error)) {
    // HTTPError - has response property
    console.error(`HTTP ${error.response.status}: ${error.message}`);

    // Read error response body
    const errorBody = await error.response.json();

    // Handle specific status codes
    if (error.response.status === 404) {
      showNotification("Resource not found");
    } else if (error.response.status === 422) {
      // Validation errors
      handleValidationErrors(errorBody.errors);
    }
  } else if (isTimeoutError(error)) {
    // TimeoutError - no response (request timed out)
    showNotification("Request timed out. Please try again.");
  } else {
    // Network error or other error
    showNotification("Network error. Please check your connection.");
  }
}
```

### Consuming Error Response Bodies

**IMPORTANT**: Always consume or cancel error response bodies to prevent resource leaks:

```ts
try {
  await apiClient.get("api/v1/data").json();
} catch (error) {
  if (isHTTPError(error)) {
    // Option 1: Read the error response body
    const errorJson = await error.response.json();
    console.error("Error details:", errorJson);

    // Option 2: Cancel the body if you don't need it
    // await error.response.body?.cancel();
  }
}
```

## Request Cancellation

Use AbortController for cancelling requests:

```ts
const controller = new AbortController();
const { signal } = controller;

// Start request
const promise = apiClient.get("api/v1/large-data", { signal }).json();

// Cancel after 5 seconds
setTimeout(() => {
  controller.abort();
}, 5000);

try {
  const data = await promise;
} catch (error) {
  if (error.name === "AbortError") {
    console.log("Request was cancelled");
  }
}
```

## Progress Tracking

### Download Progress

```ts
const response = await apiClient.get("api/v1/large-file", {
  onDownloadProgress: (progress, chunk) => {
    const percent = Math.round(progress.percent * 100);
    console.log(
      `Downloaded: ${percent}% (${progress.transferredBytes} / ${progress.totalBytes} bytes)`
    );

    // Update UI progress bar
    updateProgressBar(percent);
  },
});
```

### Upload Progress

```ts
const formData = new FormData();
formData.append("file", fileInput.files[0]);

const response = await apiClient.post("api/v1/upload", {
  body: formData,
  onUploadProgress: (progress, chunk) => {
    const percent = Math.round(progress.percent * 100);
    console.log(
      `Uploaded: ${percent}% (${progress.transferredBytes} / ${progress.totalBytes} bytes)`
    );

    // Update UI progress bar
    updateProgressBar(percent);
  },
});
```

## Form Data Handling

### Sending multipart/form-data

```ts
const formData = new FormData();
formData.append("name", "Product Name");
formData.append("price", "99.99");
formData.append("image", fileInput.files[0]);

const product = await apiClient
  .post("api/v1/products", {
    body: formData,
    // Content-Type is automatically set to multipart/form-data
  })
  .json<Product>();
```

### Sending application/x-www-form-urlencoded

```ts
const params = new URLSearchParams();
params.set("username", "user@example.com");
params.set("password", "secret");

const response = await apiClient
  .post("api/v1/auth/login", {
    body: params,
    // Content-Type is automatically set to application/x-www-form-urlencoded
  })
  .json<AuthResponse>();
```

### Modifying FormData in Hooks

```ts
const response = await apiClient.post("api/v1/upload", {
  body: formData,
  hooks: {
    beforeRequest: [
      (request) => {
        const newFormData = new FormData();

        // Transform field names (e.g., to snake_case)
        for (const [key, value] of formData) {
          newFormData.set(toSnakeCase(key), value);
        }

        // IMPORTANT: Delete Content-Type to let Request regenerate it
        // with correct boundary
        request.headers.delete("content-type");

        return new Request(request, { body: newFormData });
      },
    ],
  },
});
```

## Context for Hooks

Pass arbitrary data to hooks without polluting the request:

```ts
const api = ky.create({
  hooks: {
    beforeRequest: [
      (request, options) => {
        const { skipAuth, businessId } = options.context;

        if (!skipAuth) {
          request.headers.set("Authorization", `Bearer ${getToken()}`);
        }

        if (businessId) {
          request.headers.set("X-Business-ID", businessId);
        }
      },
    ],
  },
});

// Use context
await api
  .get("api/v1/data", {
    context: {
      businessId: "123",
      skipAuth: false,
    },
  })
  .json();
```

## Production Best Practices

### 1. Centralized Error Handling

```ts
// src/api/errors.ts
import { HTTPError, TimeoutError } from "ky";
import { showNotification } from "@/lib/notifications";

export async function handleApiError(error: unknown): Promise<void> {
  if (error instanceof HTTPError) {
    const body = await error.response.json().catch(() => ({}));

    switch (error.response.status) {
      case 400:
        showNotification(body.message || "Invalid request", "error");
        break;
      case 401:
        // Handled by auth interceptor
        break;
      case 403:
        showNotification(
          "You do not have permission to perform this action",
          "error"
        );
        break;
      case 404:
        showNotification("Resource not found", "error");
        break;
      case 422:
        // Handle validation errors
        if (body.errors) {
          Object.values(body.errors).forEach((msg: any) => {
            showNotification(msg, "error");
          });
        }
        break;
      case 429:
        showNotification(
          "Too many requests. Please try again later",
          "warning"
        );
        break;
      case 500:
      case 502:
      case 503:
      case 504:
        showNotification("Server error. Please try again later", "error");
        break;
      default:
        showNotification(body.message || "An error occurred", "error");
    }
  } else if (error instanceof TimeoutError) {
    showNotification("Request timed out. Please try again", "error");
  } else {
    showNotification("Network error. Please check your connection", "error");
  }
}
```

### 2. Request/Response Logging

```ts
export const apiClient = ky.create({
  hooks: {
    beforeRequest: [
      (request, options, { retryCount }) => {
        if (retryCount === 0) {
          const requestId = crypto.randomUUID();
          request.headers.set("X-Request-ID", requestId);

          console.group(`[API Request] ${request.method} ${request.url}`);
          console.log("Request ID:", requestId);
          console.log("Headers:", Object.fromEntries(request.headers));
          console.groupEnd();

          // Store for response logging
          (request as any).startTime = Date.now();
          (request as any).requestId = requestId;
        }
      },
    ],
    afterResponse: [
      (request, options, response) => {
        const duration = Date.now() - (request as any).startTime;
        const requestId = (request as any).requestId;

        console.group(`[API Response] ${request.method} ${request.url}`);
        console.log("Request ID:", requestId);
        console.log("Status:", response.status);
        console.log("Duration:", `${duration}ms`);
        console.groupEnd();
      },
    ],
  },
});
```

### 3. Type-Safe API Wrappers

```ts
// src/api/users.ts
import { apiClient } from "./client";

export interface User {
  id: string;
  email: string;
  name: string;
  role: "admin" | "member";
}

export interface CreateUserInput {
  email: string;
  name: string;
  password: string;
}

export interface ListUsersParams {
  page?: number;
  limit?: number;
  search?: string;
}

export const usersApi = {
  async list(params?: ListUsersParams): Promise<User[]> {
    return apiClient
      .get("api/v1/users", { searchParams: params })
      .json<User[]>();
  },

  async get(id: string): Promise<User> {
    return apiClient.get(`api/v1/users/${id}`).json<User>();
  },

  async create(input: CreateUserInput): Promise<User> {
    return apiClient.post("api/v1/users", { json: input }).json<User>();
  },

  async update(id: string, input: Partial<User>): Promise<User> {
    return apiClient.patch(`api/v1/users/${id}`, { json: input }).json<User>();
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`api/v1/users/${id}`);
  },
};
```

### 4. Testing Strategies

```ts
// src/api/__tests__/client.test.ts
import { describe, it, expect, vi, beforeEach } from "vitest";
import { apiClient } from "../client";

// Mock fetch
global.fetch = vi.fn();

describe("apiClient", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("should add authorization header", async () => {
    (global.fetch as any).mockResolvedValueOnce({
      ok: true,
      json: async () => ({ data: "test" }),
    });

    await apiClient.get("test").json();

    const callArgs = (global.fetch as any).mock.calls[0];
    expect(callArgs[0].headers.get("Authorization")).toBe("Bearer token");
  });

  it("should retry on 429", async () => {
    (global.fetch as any)
      .mockResolvedValueOnce({
        ok: false,
        status: 429,
        headers: new Headers({ "Retry-After": "1" }),
      })
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: "success" }),
      });

    const result = await apiClient.get("test").json();

    expect(global.fetch).toHaveBeenCalledTimes(2);
    expect(result).toEqual({ data: "success" });
  });
});
```

## Common Patterns

### Pagination

```ts
export interface PaginatedResponse<T> {
  data: T[];
  meta: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

export async function fetchPaginated<T>(
  endpoint: string,
  page: number = 1,
  limit: number = 20
): Promise<PaginatedResponse<T>> {
  return apiClient
    .get(endpoint, {
      searchParams: { page, limit },
    })
    .json<PaginatedResponse<T>>();
}

// Usage
const orders = await fetchPaginated<Order>("api/v1/orders", 1, 20);
```

### File Downloads

```ts
export async function downloadFile(
  url: string,
  filename: string
): Promise<void> {
  const response = await apiClient.get(url, {
    onDownloadProgress: (progress) => {
      const percent = Math.round(progress.percent * 100);
      updateDownloadProgress(percent);
    },
  });

  const blob = await response.blob();
  const downloadUrl = window.URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = downloadUrl;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  window.URL.revokeObjectURL(downloadUrl);
}
```

### Debounced Search

```ts
import { debounce } from "lodash-es";

let searchController: AbortController | null = null;

export const debouncedSearch = debounce(async (query: string) => {
  // Cancel previous request
  if (searchController) {
    searchController.abort();
  }

  searchController = new AbortController();

  try {
    const results = await apiClient
      .get("api/v1/search", {
        searchParams: { q: query },
        signal: searchController.signal,
      })
      .json<SearchResult[]>();

    return results;
  } catch (error) {
    if (error.name === "AbortError") {
      // Request was cancelled, ignore
      return [];
    }
    throw error;
  }
}, 300);
```

## Security Considerations

1. **Never expose sensitive tokens in client code**
2. **Always use HTTPS in production** (validate with `prefixUrl`)
3. **Validate and sanitize user input before sending**
4. **Use HttpOnly cookies for refresh tokens** (handled by backend)
5. **Implement CSRF protection** if not using bearer tokens
6. **Set appropriate CORS headers** (backend responsibility)
7. **Consume error response bodies** to prevent leaks

## Performance Tips

1. **Use `ky.extend()` instead of `ky.create()`** when inheriting defaults
2. **Enable HTTP/2** for Node.js environments (see Ky docs)
3. **Use compression** (handled by server, Ky automatically accepts it)
4. **Implement request cancellation** for user-initiated actions
5. **Use progress tracking** for large uploads/downloads
6. **Cache responses** where appropriate (use custom cache layer)
7. **Optimize retry settings** based on your use case

## Debugging

```ts
// Enable verbose logging in development
if (import.meta.env.DEV) {
  apiClient = apiClient.extend({
    hooks: {
      beforeRequest: [
        (request) => {
          console.log("[Ky] Request:", {
            method: request.method,
            url: request.url,
            headers: Object.fromEntries(request.headers),
          });
        },
      ],
      afterResponse: [
        async (request, options, response) => {
          const body = await response.clone().text();
          console.log("[Ky] Response:", {
            status: response.status,
            headers: Object.fromEntries(response.headers),
            body: body.substring(0, 500), // First 500 chars
          });
        },
      ],
      beforeError: [
        (error) => {
          console.error("[Ky] Error:", error);
          return error;
        },
      ],
    },
  });
}
```

## Summary

Ky provides a powerful, type-safe HTTP client with:

- **Automatic retries** with exponential backoff and jitter
- **Comprehensive hook system** for request/response interception
- **Built-in timeout** management
- **Type-safe JSON** handling with generics
- **Error handling** with typed errors
- **Progress tracking** for uploads/downloads
- **Request cancellation** with AbortController

Use this guide as a reference when implementing HTTP requests in the Kyora Portal Web App. Always prefer creating domain-specific API client instances with appropriate defaults over inline configuration.
