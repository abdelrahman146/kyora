import { useStore } from '@tanstack/react-store'
import { authStore, login, logout, logoutAll } from '@/stores/authStore'

/**
 * useAuth Hook
 *
 * Provides access to authentication state and actions from authStore.
 *
 * @example
 * ```tsx
 * const { user, isAuthenticated, login, logout } = useAuth();
 *
 * const handleLogin = async () => {
 *   await login('email@example.com', 'password');
 * };
 * ```
 */
export function useAuth() {
  const state = useStore(authStore)

  return {
    user: state.user,
    isAuthenticated: state.isAuthenticated,
    isLoading: state.isLoading,
    login,
    logout,
    logoutAll,
  }
}
