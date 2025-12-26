import { Navigate, useLocation } from "react-router";
import { useAuth } from "../../hooks/useAuth";

/**
 * RequireAuth Component Props
 */
interface RequireAuthProps {
  /** Children to render if authenticated */
  children: React.ReactNode;

  /** Optional redirect path (defaults to /login) */
  redirectTo?: string;
}

/**
 * Route Guard Component for Protected Routes
 *
 * Renders children if user is authenticated, otherwise redirects to login.
 * Preserves the intended destination in location state for post-login redirect.
 * Shows loading skeleton while checking authentication status.
 *
 * @example
 * ```tsx
 * // Protect a single route
 * <Route
 *   path="/dashboard"
 *   element={
 *     <RequireAuth>
 *       <Dashboard />
 *     </RequireAuth>
 *   }
 * />
 *
 * // Protect multiple routes with layout
 * <Route
 *   element={
 *     <RequireAuth>
 *       <AppLayout />
 *     </RequireAuth>
 *   }
 * >
 *   <Route path="/dashboard" element={<Dashboard />} />
 *   <Route path="/orders" element={<Orders />} />
 * </Route>
 * ```
 */
export function RequireAuth({
  children,
  redirectTo = "/login",
}: RequireAuthProps) {
  const { isAuthenticated, isLoading } = useAuth();
  const location = useLocation();

  // Show loading skeleton while checking auth status
  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="flex flex-col items-center gap-4">
          {/* Loading spinner */}
          <span className="loading loading-spinner loading-lg text-primary"></span>

          {/* Loading text */}
          <p className="text-base-content/60 text-sm">Loading...</p>
        </div>
      </div>
    );
  }

  // Redirect to login if not authenticated
  if (!isAuthenticated) {
    // Save the attempted URL for post-login redirect
    return <Navigate to={redirectTo} state={{ from: location }} replace />;
  }

  // User is authenticated - render protected content
  return <>{children}</>;
}
