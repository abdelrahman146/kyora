# Authentication & Token Management Guide

## Overview

Kyora Portal Web implements **secure JWT-based authentication** with automatic token refresh, persistent sessions, and comprehensive security measures.

## Architecture

### Token Strategy

**Access Token** (In-Memory):
- Stored in memory (JavaScript variable)
- Short-lived (typically 15-60 minutes)
- Cleared on page refresh/tab close
- Used for API authentication (Bearer token)
- **More secure** - not accessible to scripts after page reload

**Refresh Token** (Secure Cookie):
- Stored in browser cookie with `Secure` flag (HTTPS-only in production)
- Long-lived (typically 30-90 days)
- Survives page refresh and browser restart
- Used only to obtain new access tokens
- `SameSite=Lax` protection against CSRF

### Why This Approach?

1. **Security**: Access tokens in memory are cleared on refresh, limiting XSS attack window
2. **Convenience**: Refresh token in cookie maintains persistent login
3. **Balance**: Users don't need to re-login frequently, but sensitive operations are protected
4. **Best Practice**: Aligns with OWASP recommendations for SPA token storage

## Components

### 1. Token Management (`src/api/client.ts`)

```typescript
import { getAccessToken, setTokens, clearTokens } from "@/api/client";

// Get current access token
const token = getAccessToken(); // Returns null if not set

// Store tokens after login
setTokens(accessToken, refreshToken);

// Clear tokens on logout
clearTokens();
```

**Automatic Features**:
- Attaches `Authorization: Bearer <token>` header to all requests
- Detects 401 Unauthorized responses
- Automatically calls `/v1/auth/refresh` with refresh token
- Retries original request with new access token
- Deduplicates simultaneous refresh requests
- Redirects to `/login` if refresh fails

### 2. Authentication Context (`src/hooks/useAuth.tsx`)

```typescript
import { useAuth } from "@/hooks/useAuth";

function MyComponent() {
  const { user, login, logout, isAuthenticated, isLoading } = useAuth();

  // user: User | null
  // isAuthenticated: boolean
  // isLoading: boolean (true during initialization)
  // login: (credentials) => Promise<void>
  // logout: () => Promise<void>
  // logoutAll: () => Promise<void>
}
```

**Features**:
- Automatically restores session on mount using refresh token
- Manages user state across the application
- Handles login/logout operations
- Provides loading states for UI feedback

### 3. Route Guard (`src/components/routing/RequireAuth.tsx`)

```typescript
import { RequireAuth } from "@/components/routing/RequireAuth";

// Protect a single route
<Route
  path="/dashboard"
  element={
    <RequireAuth>
      <Dashboard />
    </RequireAuth>
  }
/>

// Protect multiple routes with layout
<Route
  element={
    <RequireAuth>
      <AppLayout />
    </RequireAuth>
  }
>
  <Route path="/dashboard" element={<Dashboard />} />
  <Route path="/orders" element={<Orders />} />
  <Route path="/settings" element={<Settings />} />
</Route>
```

**Features**:
- Shows loading spinner while checking auth status
- Redirects to `/login` if not authenticated
- Preserves intended destination in location state
- Supports custom redirect paths

## Usage Examples

### Example 1: App Setup with AuthProvider

```tsx
// src/main.tsx or src/App.tsx
import { AuthProvider } from "@/hooks/useAuth";

function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />

          {/* Protected routes */}
          <Route
            element={
              <RequireAuth>
                <AppLayout />
              </RequireAuth>
            }
          >
            <Route path="/dashboard" element={<Dashboard />} />
            <Route path="/orders" element={<Orders />} />
          </Route>
        </Routes>
      </Router>
    </AuthProvider>
  );
}
```

### Example 2: Login Form

```tsx
import { useAuth } from "@/hooks/useAuth";
import { useNavigate, useLocation } from "react-router";
import { translateErrorAsync } from "@/lib/translateError";
import { useTranslation } from "react-i18next";
import toast from "react-hot-toast";

function LoginForm() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const { t } = useTranslation();
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      await login({ email, password });

      // Redirect to intended destination or dashboard
      const from = location.state?.from?.pathname || "/dashboard";
      navigate(from, { replace: true });

      toast.success(t("auth.login.success"));
    } catch (error) {
      const message = await translateErrorAsync(error, t);
      toast.error(message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      {/* Form fields */}
      <button type="submit" disabled={isLoading}>
        {isLoading ? "Logging in..." : "Login"}
      </button>
    </form>
  );
}
```

### Example 3: Protected Dashboard

```tsx
import { useAuth } from "@/hooks/useAuth";

function Dashboard() {
  const { user } = useAuth();

  return (
    <div>
      <h1>Welcome, {user?.firstName}!</h1>
      <p>Email: {user?.email}</p>
      <p>Role: {user?.role}</p>
    </div>
  );
}

// Usage (automatically protected by RequireAuth wrapper)
<Route
  path="/dashboard"
  element={
    <RequireAuth>
      <Dashboard />
    </RequireAuth>
  }
/>
```

### Example 4: Logout Button

```tsx
import { useAuth } from "@/hooks/useAuth";
import { useNavigate } from "react-router";
import toast from "react-hot-toast";

function LogoutButton() {
  const { logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = async () => {
    try {
      await logout();
      navigate("/login");
      toast.success("Logged out successfully");
    } catch (error) {
      console.error("Logout failed:", error);
      // Still navigate to login even if API call fails
      navigate("/login");
    }
  };

  return (
    <button onClick={handleLogout} className="btn btn-ghost">
      Logout
    </button>
  );
}
```

### Example 5: Show User Info in Header

```tsx
import { useAuth } from "@/hooks/useAuth";

function AppHeader() {
  const { user, isAuthenticated, logout } = useAuth();

  if (!isAuthenticated) {
    return null;
  }

  return (
    <header className="navbar bg-base-100">
      <div className="flex-1">
        <a className="btn btn-ghost text-xl">Kyora</a>
      </div>

      <div className="flex-none gap-2">
        {/* User dropdown */}
        <div className="dropdown dropdown-end">
          <div tabIndex={0} role="button" className="avatar btn btn-circle btn-ghost">
            <div className="w-10 rounded-full">
              <span className="flex h-full items-center justify-center bg-primary text-primary-content">
                {user?.firstName?.[0] ?? "?"}
              </span>
            </div>
          </div>

          <ul className="menu dropdown-content menu-sm z-[1] mt-3 w-52 rounded-box bg-base-100 p-2 shadow">
            <li>
              <a className="justify-between">
                {user?.firstName} {user?.lastName}
                <span className="badge">{user?.role}</span>
              </a>
            </li>
            <li>
              <a>Settings</a>
            </li>
            <li>
              <a onClick={logout}>Logout</a>
            </li>
          </ul>
        </div>
      </div>
    </header>
  );
}
```

### Example 6: Conditional Rendering Based on Auth

```tsx
import { useAuth } from "@/hooks/useAuth";
import { Navigate } from "react-router";

function LandingPage() {
  const { isAuthenticated, isLoading } = useAuth();

  // Show loading while checking auth
  if (isLoading) {
    return <div>Loading...</div>;
  }

  // Redirect to dashboard if already authenticated
  if (isAuthenticated) {
    return <Navigate to="/dashboard" replace />;
  }

  // Show landing page for unauthenticated users
  return (
    <div>
      <h1>Welcome to Kyora</h1>
      <a href="/login">Login</a>
      <a href="/register">Sign Up</a>
    </div>
  );
}
```

### Example 7: Logout All Sessions

```tsx
import { useAuth } from "@/hooks/useAuth";

function SecuritySettings() {
  const { logoutAll } = useAuth();
  const navigate = useNavigate();

  const handleLogoutAll = async () => {
    if (!confirm("This will log you out from all devices. Continue?")) {
      return;
    }

    try {
      await logoutAll();
      navigate("/login");
      toast.success("Logged out from all devices");
    } catch (error) {
      toast.error("Failed to logout from all devices");
    }
  };

  return (
    <div>
      <h2>Active Sessions</h2>
      <button onClick={handleLogoutAll} className="btn btn-error">
        Logout All Devices
      </button>
    </div>
  );
}
```

## Security Considerations

### ✅ What We Do Right

1. **Access Token in Memory**: Not accessible after page refresh, limiting XSS attack window
2. **Secure Cookie Flag**: HTTPS-only transmission in production
3. **SameSite=Lax**: Protects against CSRF attacks
4. **Automatic Token Refresh**: Minimizes time using expired tokens
5. **Deduplication**: Prevents multiple simultaneous refresh requests
6. **Graceful Failure**: Redirects to login if refresh fails
7. **No Local Storage**: Avoiding the most vulnerable storage option

### ⚠️ Limitations & Future Improvements

1. **Not True HttpOnly**: Client-side JavaScript can read the refresh token from cookies
   - **Solution**: Backend should set refresh token via `Set-Cookie` header with `HttpOnly` flag
   - **Current**: We set cookie client-side, which is readable by JavaScript

2. **Refresh Token Rotation**: Not implemented yet
   - **Best Practice**: Issue new refresh token on each refresh
   - **Current**: Reuses same refresh token

3. **Token Expiry Tracking**: No client-side expiry checking
   - **Enhancement**: Could decode JWT to check expiry and proactively refresh

4. **Session Fingerprinting**: No device/browser fingerprinting
   - **Enhancement**: Track session metadata for security monitoring

### Migration to True HttpOnly Cookies

If backend starts setting refresh token via HttpOnly cookie:

```typescript
// Remove this line from setTokens():
setCookie(REFRESH_TOKEN_COOKIE_NAME, refresh, 365);

// Backend will set it via:
// Set-Cookie: kyora_refresh_token=<token>; HttpOnly; Secure; SameSite=Lax; Path=/; Max-Age=7776000

// Keep getRefreshToken() as is - it will read from cookie
// But client-side JavaScript cannot access HttpOnly cookies via document.cookie
```

## Troubleshooting

### Issue: User logged out on page refresh

**Cause**: Access token is in memory (by design)

**Solution**: This is expected behavior. The app should automatically restore session using refresh token from cookie.

**Check**:
1. Verify refresh token exists in cookies (DevTools → Application → Cookies)
2. Check console for session restoration errors
3. Verify backend `/v1/auth/refresh` endpoint is working

### Issue: Infinite login redirect loop

**Cause**: `RequireAuth` redirects to `/login`, which redirects back

**Solution**: Don't wrap login route in `RequireAuth`

```tsx
// ❌ WRONG
<Route path="/login" element={<RequireAuth><Login /></RequireAuth>} />

// ✅ CORRECT
<Route path="/login" element={<Login />} />
```

### Issue: 401 errors after period of inactivity

**Cause**: Refresh token expired

**Solution**: Expected behavior. User needs to login again after refresh token expiry (typically 30-90 days)

### Issue: Token not attached to requests

**Cause**: `apiClient` not used for API calls

**Solution**: Always use `apiClient` or `authApi` for authenticated requests

```tsx
// ❌ WRONG
fetch("/v1/orders");

// ✅ CORRECT
import { apiClient } from "@/api/client";
apiClient.get("v1/orders");
```

## Testing

### Manual Testing Checklist

- [ ] Login with valid credentials → Redirects to dashboard
- [ ] Login with invalid credentials → Shows error message
- [ ] Page refresh while logged in → Stays logged in (session restored)
- [ ] Logout → Redirects to login, tokens cleared
- [ ] Access protected route while logged out → Redirects to login
- [ ] Login after redirect → Returns to intended destination
- [ ] Access token expires → Auto-refreshes and retries request
- [ ] Refresh token expires → Redirects to login
- [ ] Open multiple tabs → All tabs share auth state
- [ ] Close all tabs and reopen → Session restored from cookie

### Automated Testing

```typescript
// Example: Testing useAuth hook
import { renderHook, waitFor } from "@testing-library/react";
import { AuthProvider, useAuth } from "@/hooks/useAuth";

test("login sets user and tokens", async () => {
  const { result } = renderHook(() => useAuth(), {
    wrapper: AuthProvider,
  });

  await waitFor(() => {
    expect(result.current.isLoading).toBe(false);
  });

  await result.current.login({
    email: "test@example.com",
    password: "password",
  });

  expect(result.current.isAuthenticated).toBe(true);
  expect(result.current.user).not.toBeNull();
});
```

## API Reference

See individual file documentation:
- [client.ts](../src/api/client.ts) - Token management and API interceptors
- [useAuth.tsx](../src/hooks/useAuth.tsx) - Authentication context and hook
- [RequireAuth.tsx](../src/components/routing/RequireAuth.tsx) - Route guard component
- [auth.ts](../src/api/auth.ts) - Authentication API methods
