import { useStore } from '@tanstack/react-store'
import type { User } from '@/api/types/auth'
import {
  authStore,
  clearAuth,
  login,
  logout,
  logoutAll,
  setUser,
} from '@/stores/authStore'

/**
 * useAuth Hook
 *
 * Provides React bindings for authStore with convenient helpers.
 * Use this hook to access authentication state and actions in components.
 *
 * @example
 * ```tsx
 * function MyComponent() {
 *   const { user, isAuthenticated, login, logout } = useAuth()
 *
 *   if (!isAuthenticated) {
 *     return <button onClick={() => login(email, password)}>Login</button>
 *   }
 *
 *   return (
 *     <div>
 *       <p>Welcome, {user?.firstName}!</p>
 *       <button onClick={logout}>Logout</button>
 *     </div>
 *   )
 * }
 * ```
 */
export function useAuth() {
  const state = useStore(authStore)

  return {
    // State
    user: state.user,
    isAuthenticated: state.isAuthenticated,
    isLoading: state.isLoading,

    // Actions
    login,
    logout,
    logoutAll,
    setUser,
    clearAuth,
  }
}

/**
 * useUser Hook
 *
 * Shorthand to access just the user object.
 * Returns null if not authenticated.
 */
export function useUser(): User | null {
  const state = useStore(authStore)
  return state.user
}

/**
 * useIsAuthenticated Hook
 *
 * Shorthand to check authentication status.
 */
export function useIsAuthenticated(): boolean {
  const state = useStore(authStore)
  return state.isAuthenticated
}
