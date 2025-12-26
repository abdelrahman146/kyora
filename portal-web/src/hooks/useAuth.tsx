import { useContext } from "react";
import { AuthContext, type AuthContextType } from "../contexts/auth/AuthContext";

/**
 * Hook to access authentication context
 *
 * Must be used within an AuthProvider.
 *
 * @example
 * ```tsx
 * function MyComponent() {
 *   const { user, login, logout, isAuthenticated } = useAuth();
 *
 *   if (!isAuthenticated) {
 *     return <LoginForm onLogin={login} />;
 *   }
 *
 *   return <div>Hello, {user.firstName}!</div>;
 * }
 * ```
 *
 * @throws {Error} If used outside of AuthProvider
 */
export function useAuth(): AuthContextType {
  const context = useContext(AuthContext);

  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }

  return context;
}
