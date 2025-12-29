import type { User } from '@/api/types/auth'
import { authApi } from '@/api/auth'
import { clearTokens, getRefreshToken, setTokens } from '@/api/client'

/**
 * Session Restoration
 *
 * Attempts to restore user session using refresh token from cookie.
 * If successful, returns the authenticated user. Otherwise returns null.
 *
 * @returns User object if session restored, null otherwise
 */
export async function restoreSession(): Promise<User | null> {
  try {
    const refreshToken = getRefreshToken()

    if (!refreshToken) {
      // No refresh token found - user needs to login
      return null
    }

    // Attempt to refresh access token using stored refresh token
    const response = await authApi.refreshToken({ refreshToken })

    // Store new tokens
    setTokens(response.token, response.refreshToken)

    // TODO: Fetch user profile - will be implemented when user API is available
    // For now, return null to force re-login
    // const user = await userApi.getCurrentUser();
    // return user;

    return null
  } catch {
    // Refresh failed - clear invalid tokens
    clearTokens()
    return null
  }
}

/**
 * Login Helper
 *
 * Authenticates user with email and password, stores tokens, and returns user object.
 *
 * @param email - User email
 * @param password - User password
 * @returns User object
 */
export async function loginUser(
  email: string,
  password: string,
): Promise<User> {
  const response = await authApi.login({ email, password })
  setTokens(response.token, response.refreshToken)
  return response.user
}

/**
 * Logout Helper
 *
 * Logs out the current user by revoking refresh token.
 */
export async function logoutUser(): Promise<void> {
  await authApi.logoutCurrent()
  clearTokens()
}

/**
 * Logout All Devices Helper
 *
 * Logs out the user from all devices by revoking all refresh tokens.
 */
export async function logoutAllDevices(): Promise<void> {
  await authApi.logoutAllCurrent()
  clearTokens()
}

// Re-export token management functions for convenience
export { setTokens, clearTokens, getRefreshToken } from '@/api/client'
export { hasValidToken, getAccessToken } from '@/api/client'
