# API Quick Reference

## Import Paths

```typescript
// API Client
import apiClient, { setTokens, clearTokens, getAccessToken, hasValidToken } from "@/api/client";
import { get, post, put, patch, del } from "@/api/client";

// Authentication Service
import { authApi } from "@/api/auth";

// Types & Schemas
import type { User, LoginRequest, LoginResponse } from "@/api/types";
import { LoginRequestSchema, UserSchema } from "@/api/types";

// Error Handling
import { parseProblemDetails, parseValidationErrors, isHTTPError } from "@/lib/errorParser";
```

## Common Patterns

### Login

```typescript
const response = await authApi.login({ email, password });
// Tokens are automatically stored
```

### Logout

```typescript
await authApi.logoutCurrent();
// Tokens are automatically cleared
```

### Authenticated Request

```typescript
// Method 1: Using client directly
const data = await apiClient.get("v1/users").json();

// Method 2: Using typed helpers
const data = await get<User[]>("v1/users");
```

### Error Handling

```typescript
try {
  await authApi.login(credentials);
} catch (error) {
  const message = await parseProblemDetails(error);
  toast.error(message);
}
```

### Form Validation

```typescript
const { register, handleSubmit } = useForm({
  resolver: zodResolver(LoginRequestSchema),
});
```

## All Auth Methods

```typescript
authApi.login(credentials)
authApi.refreshToken(request)
authApi.logout(request)
authApi.logoutAll(request)
authApi.logoutCurrent()        // Helper: uses stored token
authApi.logoutAllCurrent()     // Helper: uses stored token
authApi.forgotPassword(request)
authApi.resetPassword(request)
authApi.loginWithGoogle(request)
authApi.requestEmailVerification(request)
authApi.verifyEmail(request)
```

## Key Features

✅ JWT Bearer token authentication  
✅ Automatic token refresh on 401  
✅ Retry logic with exponential backoff  
✅ Zod schema validation  
✅ User-friendly error messages  
✅ Type-safe API calls  
✅ Request/response logging in dev mode  
✅ ProblemDetails (RFC 7807) parsing  

## Next Steps

1. Test login/logout flows
2. Add more domain services (orders, customers, etc.)
3. Integrate with TanStack Query for caching
4. Build auth context/provider
5. Create protected route components
