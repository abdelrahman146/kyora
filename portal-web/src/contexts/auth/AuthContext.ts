import { createContext } from "react";
import type { User, LoginRequest } from "../../api/types/auth";

/**
 * Authentication Context Interface
 */
export interface AuthContextType {
  /** Current authenticated user (null if not authenticated) */
  user: User | null;

  /** Whether the auth state is currently loading (checking tokens, etc.) */
  isLoading: boolean;

  /** Whether user is authenticated (has valid tokens and user data) */
  isAuthenticated: boolean;

  /** Login with email and password */
  login: (credentials: LoginRequest) => Promise<void>;

  /**
   * Set authenticated session directly.
   * Useful when an endpoint returns user + tokens (e.g. onboarding completion).
   */
  setSession: (session: { user: User; token: string; refreshToken: string }) => void;

  /** Logout current session (revoke refresh token) */
  logout: () => Promise<void>;

  /** Logout all sessions (revoke all refresh tokens for this user) */
  logoutAll: () => Promise<void>;
}

export const AuthContext = createContext<AuthContextType | undefined>(
  undefined
);
