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
    } catch (error) {
      console.error("Failed to read session token:", error);
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
    } catch (error) {
      console.error("Failed to save session token:", error);
    }
  },

  /**
   * Removes the stored session token.
   */
  clearToken(): void {
    try {
      localStorage.removeItem(SESSION_TOKEN_KEY);
    } catch (error) {
      console.error("Failed to clear session token:", error);
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
