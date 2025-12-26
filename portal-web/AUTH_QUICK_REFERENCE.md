# Authentication Quick Reference

## Setup (One-time)

```tsx
// src/main.tsx
import { AuthProvider } from "@/hooks/useAuth";

<AuthProvider>
  <App />
</AuthProvider>
```

## Protect Routes

```tsx
import { RequireAuth } from "@/components/routing/RequireAuth";

// Single route
<Route
  path="/dashboard"
  element={
    <RequireAuth>
      <Dashboard />
    </RequireAuth>
  }
/>

// Layout with nested routes
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
```

## Use in Components

```tsx
import { useAuth } from "@/hooks/useAuth";

function MyComponent() {
  const { user, login, logout, isAuthenticated, isLoading } = useAuth();

  // Login
  const handleLogin = async () => {
    try {
      await login({ email, password });
      navigate("/dashboard");
    } catch (error) {
      // Handle error
    }
  };

  // Logout
  const handleLogout = async () => {
    await logout();
    navigate("/login");
  };

  // Check auth
  if (isLoading) return <div>Loading...</div>;
  if (!isAuthenticated) return <div>Please login</div>;

  // Show user data
  return <div>Hello, {user?.firstName}!</div>;
}
```

## Token Management (Automatic)

```typescript
// Already handled by apiClient - no manual intervention needed!

// ✅ Automatic: Bearer token attached to requests
// ✅ Automatic: 401 detection and token refresh
// ✅ Automatic: Request retry with new token
// ✅ Automatic: Redirect to login if refresh fails

// Manual token access (rarely needed)
import { getAccessToken, getRefreshToken } from "@/api/client";

const token = getAccessToken(); // or null
const refresh = getRefreshToken(); // or null
```

## Common Patterns

### Login Form
```tsx
const { login } = useAuth();
const location = useLocation();

const handleSubmit = async (credentials) => {
  await login(credentials);
  const from = location.state?.from?.pathname || "/dashboard";
  navigate(from, { replace: true });
};
```

### Logout Button
```tsx
const { logout } = useAuth();

<button onClick={logout}>Logout</button>
```

### User Avatar
```tsx
const { user } = useAuth();

<div className="avatar">
  <span>{user?.firstName?.[0]}</span>
</div>
```

### Conditional Redirect
```tsx
const { isAuthenticated, isLoading } = useAuth();

if (isLoading) return <Loading />;
if (isAuthenticated) return <Navigate to="/dashboard" />;

return <LandingPage />;
```

## API Methods

```typescript
const { user, isAuthenticated, isLoading, login, logout, logoutAll } = useAuth();
```

| Method | Type | Description |
|--------|------|-------------|
| `user` | `User \| null` | Current user or null |
| `isAuthenticated` | `boolean` | Has valid token and user |
| `isLoading` | `boolean` | Checking auth status |
| `login(credentials)` | `Promise<void>` | Login with email/password |
| `logout()` | `Promise<void>` | Logout current session |
| `logoutAll()` | `Promise<void>` | Logout all sessions |

## User Object

```typescript
interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  role: "admin" | "user";
  workspaceId: string;
  createdAt: string;
  updatedAt: string;
}
```

## Error Handling

```tsx
import { translateErrorAsync } from "@/lib/translateError";
import { useTranslation } from "react-i18next";

const { t } = useTranslation();

try {
  await login(credentials);
} catch (error) {
  const message = await translateErrorAsync(error, t);
  toast.error(message); // Localized error message
}
```

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Logged out on refresh | Expected - access token is in memory. Session auto-restores via refresh token. |
| Infinite redirect loop | Don't wrap `/login` in `RequireAuth` |
| Token not attached | Use `apiClient` for all API calls |
| 401 errors persist | Check refresh token in cookies (DevTools → Application → Cookies) |

## Full Documentation

- [AUTHENTICATION.md](./AUTHENTICATION.md) - Complete guide with examples
- [AUTHENTICATION_SUMMARY.md](./AUTHENTICATION_SUMMARY.md) - Implementation summary
