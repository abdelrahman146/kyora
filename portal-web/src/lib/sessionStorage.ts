/**
 * LocalStorage utility for managing onboarding session tokens.
 *
 * Security considerations:
 * - Only stores the session token (not sensitive user data)
 * - Token is validated on backend with every request
 * - Sessions expire after 24 hours on backend
 * - Should only be used over HTTPS in production
 */

const SESSION_TOKEN_KEY = "kyora_onboarding_session";

export const sessionStorage = {
  /**
   * Retrieves the stored session token.
   * @returns The session token or null if not found
   */
  getToken(): string | null {
    try {
      return localStorage.getItem(SESSION_TOKEN_KEY);
    } catch {
      return null;
    }
  },

  /**
   * Stores a session token.
   * @param token - The session token to store
   */
  setToken(token: string): void {
    try {
      localStorage.setItem(SESSION_TOKEN_KEY, token);
    } catch {
      // Silent fail - storage might not be available
    }
  },

  /**
   * Removes the stored session token.
   */
  clearToken(): void {
    try {
      localStorage.removeItem(SESSION_TOKEN_KEY);
    } catch {
      // Silent fail
    }
  },

  /**
   * Checks if a session token exists.
   * @returns true if a token is stored
   */
  hasToken(): boolean {
    return this.getToken() !== null;
  },
};
