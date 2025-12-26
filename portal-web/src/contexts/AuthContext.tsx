import {
  useState,
  useEffect,
  useCallback,
  type ReactNode,
} from "react";
import { authApi } from "../api/auth";
import { userApi } from "../api/user";
import {
  setTokens,
  clearTokens,
  getRefreshToken,
  hasValidToken,
} from "../api/client";
import { AuthContext, type AuthContextType } from "./auth/AuthContext";
import type { User, LoginRequest } from "../api/types/auth";

/**
 * Authentication Provider Props
 */
interface AuthProviderProps {
  children: ReactNode;
}

/**
 * Authentication Provider Component
 *
 * Manages authentication state, token storage, and auth operations.
 *
 * Token Storage Strategy:
 * - Access Token: In-memory (cleared on page refresh for security)
 * - Refresh Token: Secure cookie (persistent across sessions)
 *
 * On mount, attempts to restore session using refresh token from cookie.
 *
 * @example
 * ```tsx
 * <AuthProvider>
 *   <App />
 * </AuthProvider>
 * ```
 */
export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  /**
   * Restore session on mount using refresh token
   * If refresh token exists in cookie, attempt to get new access token
   */
  const restoreSession = useCallback(async () => {
    setIsLoading(true);

    try {
      const refreshToken = getRefreshToken();

      if (!refreshToken) {
        // No refresh token found - user needs to login
        setIsLoading(false);
        return;
      }

      // Attempt to refresh access token using stored refresh token
      const response = await authApi.refreshToken({ refreshToken });

      // Store new tokens
      setTokens(response.token, response.refreshToken);

      // Fetch user profile using new access token
      const user = await userApi.getCurrentUser();
      setUser(user);
    } catch (error) {
      // Refresh failed - clear invalid tokens
      clearTokens();
      setUser(null);

      console.warn("[Auth] Session restoration failed:", error);
    } finally {
      setIsLoading(false);
    }
  }, []);

  /**
   * Initialize auth state on mount
   */
  useEffect(() => {
    void restoreSession();
  }, [restoreSession]);

  /**
   * Login with email and password
   */
  const login = useCallback(async (credentials: LoginRequest) => {
    setIsLoading(true);

    try {
      const response = await authApi.login(credentials);

      // Store tokens
      setTokens(response.token, response.refreshToken);

      // Set user
      setUser(response.user);
    } catch (error: unknown) {
      // Clear any partial state on error
      clearTokens();
      setUser(null);
      throw error; // Re-throw for component error handling
    } finally {
      setIsLoading(false);
    }
  }, []);

  /**
   * Logout current session
   * Revokes the current refresh token on the backend
   */
  const logout = useCallback(async () => {
    try {
      // Call backend to revoke current refresh token
      await authApi.logoutCurrent();
    } catch (error) {
      console.warn("[Auth] Logout API call failed:", error);
      // Continue with local cleanup even if API call fails
    } finally {
      // Clear tokens and user state
      clearTokens();
      setUser(null);
    }
  }, []);

  /**
   * Logout all sessions
   * Revokes all refresh tokens for this user on the backend
   */
  const logoutAll = useCallback(async () => {
    try {
      // Call backend to revoke all refresh tokens
      await authApi.logoutAllCurrent();
    } catch (error) {
      console.warn("[Auth] Logout all API call failed:", error);
      // Continue with local cleanup even if API call fails
    } finally {
      // Clear tokens and user state
      clearTokens();
      setUser(null);
    }
  }, []);

  const isAuthenticated = user !== null && hasValidToken();

  const value: AuthContextType = {
    user,
    isLoading,
    isAuthenticated,
    login,
    logout,
    logoutAll,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
