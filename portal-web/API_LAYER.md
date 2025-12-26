# API Layer Documentation

This document describes the robust HTTP layer implementation for the Kyora Portal Web application.

## Overview

The API layer provides a production-ready HTTP client with the following features:

- ✅ JWT Bearer token authentication with automatic refresh
- ✅ Zod schema validation for all requests and responses
- ✅ Comprehensive error handling with user-friendly messages
- ✅ Retry logic with exponential backoff
- ✅ Type-safe API methods
- ✅ RFC 7807 ProblemDetails error parsing

## Architecture

```
src/api/
├── client.ts          # Centralized ky HTTP client with auth & retry logic
├── auth.ts            # Authentication service with all auth endpoints
├── types/
│   ├── auth.ts        # Zod schemas for auth types
│   └── index.ts       # Re-exports all types
└── ...                # Other domain-specific API services

src/lib/
└── errorParser.ts     # ProblemDetails error parser utility
```

## Core Components

### 1. API Client (`src/api/client.ts`)

The centralized HTTP client built on top of `ky` with the following features:

#### Features

- **Authentication**: Automatically adds JWT Bearer token to all requests
- **Token Refresh**: Detects 401 responses and automatically refreshes the access token
- **Retry Logic**: Exponential backoff for failed requests (408, 413, 429, 500, 502, 503, 504)
- **Error Handling**: Parses backend ProblemDetails format into user-friendly messages
- **Request Logging**: Logs all requests/responses in development mode

#### Token Management

```typescript
import { setTokens, clearTokens, getAccessToken, hasValidToken } from "@/api/client";

// After successful login
setTokens(accessToken, refreshToken);

// Check if user is authenticated
if (hasValidToken()) {
  // User is logged in
}

// On logout
clearTokens();
```

#### Making API Calls

```typescript
import apiClient from "@/api/client";

// Using the client directly
const user = await apiClient.get("v1/users/me").json();

// Using typed helper methods
import { get, post, put, patch, del } from "@/api/client";

const users = await get<User[]>("v1/users");
const newUser = await post<User>("v1/users", { json: userData });
```

### 2. Authentication Service (`src/api/auth.ts`)

Provides type-safe methods for all authentication endpoints.

#### Available Methods

```typescript
import { authApi } from "@/api/auth";

// Login with email & password
const response = await authApi.login({
  email: "user@example.com",
  password: "password123",
});
// Returns: { token, refreshToken, user }

// Refresh access token
const tokens = await authApi.refreshToken({
  refreshToken: "...",
});
// Returns: { token, refreshToken }

// Logout current session
await authApi.logoutCurrent();

// Logout all devices
await authApi.logoutAllCurrent();

// Forgot password
await authApi.forgotPassword({
  email: "user@example.com",
});

// Reset password
await authApi.resetPassword({
  token: "reset-token",
  password: "newPassword123",
});

// Login with Google OAuth
const response = await authApi.loginWithGoogle({
  code: "google-oauth-code",
});

// Request email verification
await authApi.requestEmailVerification({
  email: "user@example.com",
});

// Verify email
await authApi.verifyEmail({
  token: "verification-token",
});
```

### 3. Type Schemas (`src/api/types/auth.ts`)

All request/response types are validated using Zod schemas.

#### Available Schemas

```typescript
import {
  LoginRequestSchema,
  LoginResponseSchema,
  UserSchema,
  // ... other schemas
} from "@/api/types/auth";

// Validate request data before sending
const validatedData = LoginRequestSchema.parse({
  email: "user@example.com",
  password: "password123",
});

// Validate response data after receiving
const validatedResponse = LoginResponseSchema.parse(apiResponse);
```

#### Type Exports

```typescript
import type {
  User,
  LoginRequest,
  LoginResponse,
  RefreshRequest,
  RefreshResponse,
  ProblemDetails,
  // ... other types
} from "@/api/types/auth";
```

### 4. Error Parser (`src/lib/errorParser.ts`)

Parses backend RFC 7807 ProblemDetails errors into user-friendly messages.

#### Usage

```typescript
import { parseProblemDetails, parseValidationErrors, isHTTPError } from "@/lib/errorParser";

try {
  await authApi.login(credentials);
} catch (error) {
  // Get user-friendly error message
  const message = await parseProblemDetails(error);
  toast.error(message);

  // Extract validation errors (if any)
  const validationErrors = await parseValidationErrors(error);
  if (validationErrors) {
    // validationErrors = { email: "Invalid format", password: "Too short" }
    Object.entries(validationErrors).forEach(([field, message]) => {
      setError(field, { message });
    });
  }

  // Check for specific HTTP status
  if (isHTTPError(error)) {
    if (error.response.status === 429) {
      toast.error("Too many requests. Please wait a moment.");
    }
  }
}
```

## Usage Examples

### Complete Login Flow

```typescript
import { authApi } from "@/api/auth";
import { parseProblemDetails } from "@/lib/errorParser";
import { useNavigate } from "react-router-dom";

async function handleLogin(email: string, password: string) {
  try {
    // Login (automatically stores tokens)
    const response = await authApi.login({ email, password });

    // Response is validated and typed
    console.log(`Welcome, ${response.user.firstName}!`);

    // Redirect to dashboard
    navigate("/dashboard");
  } catch (error) {
    // Parse backend error into user-friendly message
    const errorMessage = await parseProblemDetails(error);
    toast.error(errorMessage);
  }
}
```

### Automatic Token Refresh

The API client automatically handles token refresh:

```typescript
// Make any authenticated request
const orders = await apiClient.get("v1/orders").json();

// If the access token is expired (401 response):
// 1. Client automatically calls /v1/auth/refresh with refresh token
// 2. Updates tokens in memory
// 3. Retries the original request with new token
// 4. Returns the result seamlessly
//
// If refresh fails:
// 1. Clears tokens
// 2. Redirects to /login
```

### Form Validation with Zod

```typescript
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoginRequestSchema } from "@/api/types/auth";

const {
  register,
  handleSubmit,
  formState: { errors },
} = useForm({
  resolver: zodResolver(LoginRequestSchema),
});

const onSubmit = async (data) => {
  // data is already validated by Zod schema
  await authApi.login(data);
};
```

### Protected Route Component

```typescript
import { hasValidToken } from "@/api/client";
import { Navigate } from "react-router-dom";

function ProtectedRoute({ children }) {
  if (!hasValidToken()) {
    return <Navigate to="/login" replace />;
  }

  return children;
}
```

## Configuration

### Environment Variables

```bash
# .env
VITE_API_BASE_URL=http://localhost:8080
```

### Client Options

You can customize the API client in `src/api/client.ts`:

```typescript
export const apiClient = ky.create({
  prefixUrl: API_BASE_URL,
  timeout: 30000, // 30 seconds
  retry: {
    limit: 2, // Retry failed requests 2 times
    statusCodes: [408, 413, 429, 500, 502, 503, 504],
    backoffLimit: 3000, // Max 3 seconds between retries
  },
});
```

## Error Handling Best Practices

### 1. Always Use Try-Catch

```typescript
try {
  await authApi.login(credentials);
} catch (error) {
  const message = await parseProblemDetails(error);
  toast.error(message);
}
```

### 2. Handle Specific Error Cases

```typescript
import { isHTTPError } from "@/lib/errorParser";

try {
  await authApi.login(credentials);
} catch (error) {
  if (isHTTPError(error)) {
    if (error.response.status === 401) {
      toast.error("Invalid email or password");
    } else if (error.response.status === 429) {
      toast.error("Too many login attempts. Please wait.");
    } else {
      toast.error(await parseProblemDetails(error));
    }
  }
}
```

### 3. Extract Validation Errors

```typescript
import { parseValidationErrors } from "@/lib/errorParser";

try {
  await authApi.resetPassword(data);
} catch (error) {
  const validationErrors = await parseValidationErrors(error);
  if (validationErrors) {
    // Display field-specific errors
    Object.entries(validationErrors).forEach(([field, message]) => {
      setError(field, { message });
    });
  } else {
    // Generic error
    toast.error(await parseProblemDetails(error));
  }
}
```

## Testing

### Mocking API Calls

```typescript
import { vi } from "vitest";
import * as authApi from "@/api/auth";

// Mock successful login
vi.spyOn(authApi, "login").mockResolvedValue({
  token: "mock-token",
  refreshToken: "mock-refresh-token",
  user: {
    id: "1",
    email: "test@example.com",
    firstName: "Test",
    lastName: "User",
    // ... other fields
  },
});

// Mock failed login
vi.spyOn(authApi, "login").mockRejectedValue(
  new HTTPError(
    new Response(
      JSON.stringify({
        detail: "Invalid credentials",
        status: 401,
      }),
      { status: 401 }
    )
  )
);
```

## Security Considerations

1. **Token Storage**: Access tokens are stored in memory (not localStorage) to prevent XSS attacks
2. **Refresh Tokens**: Should be stored in httpOnly cookies (backend responsibility)
3. **HTTPS Only**: Always use HTTPS in production
4. **Token Expiry**: Access tokens should have short expiry (e.g., 15 minutes)
5. **Refresh Rotation**: Refresh tokens are rotated on each refresh request

## Performance Tips

1. **Parallel Requests**: Use `Promise.all()` for independent requests
2. **Request Cancellation**: Use AbortController for user-initiated cancellations
3. **Caching**: Consider using TanStack Query for automatic caching and revalidation
4. **Debouncing**: Debounce search/filter requests to reduce API calls

## Next Steps

1. **Add more domain services**: Create `orders.ts`, `customers.ts`, `inventory.ts`, etc.
2. **Implement TanStack Query**: Add query hooks for data fetching and caching
3. **Add request cancellation**: Implement AbortController for long-running requests
4. **Create API mocks**: Set up MSW (Mock Service Worker) for development and testing
5. **Add performance monitoring**: Track API request times and error rates

## Support

For questions or issues, refer to:
- [ky documentation](https://github.com/sindresorhus/ky)
- [Zod documentation](https://zod.dev/)
- [RFC 7807 ProblemDetails](https://datatracker.ietf.org/doc/html/rfc7807)
