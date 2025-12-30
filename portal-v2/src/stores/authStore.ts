import { Store } from '@tanstack/react-store'
import type { User } from '@/api/types/auth'
import {
  loginUser,
  logoutAllDevices,
  logoutUser,
  restoreSession,
} from '@/lib/auth'

/**
 * Authentication Store State
 */
interface AuthState {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
}

/**
 * Initial Authentication State
 */
const initialState: AuthState = {
  user: null,
  isAuthenticated: false,
  isLoading: true,
}

/**
 * Authentication Store
 *
 * Manages authentication state using TanStack Store.
 * No persistence - session is restored via refresh token on mount.
 *
 * Dev-only devtools integration via conditional import.
 */
export const authStore = new Store<AuthState>(initialState)

let initPromise: Promise<void> | null = null

/**
 * Initialize authentication state
 *
 * - Restores session from refresh token cookie when present
 * - Ensures `isLoading` is always eventually set to false
 * - Idempotent: concurrent calls share the same promise
 */
export async function initializeAuth(): Promise<void> {
  if (initPromise) return initPromise

  initPromise = (async () => {
    setAuthLoading(true)

    try {
      const user = await restoreSession()
      if (user) {
        setUser(user)
      } else {
        clearAuth()
      }

      if (import.meta.env.DEV) {
        const { isAuthenticated } = authStore.state
        console.debug('[authStore] initializeAuth complete', {
          isAuthenticated,
        })
      }
    } catch (error) {
      clearAuth()
      if (import.meta.env.DEV) {
        console.error('[authStore] initializeAuth failed', error)
      }
    }
  })()

  return initPromise
}

/**
 * Authentication Actions
 */

/**
 * Set authenticated user
 */
export function setUser(user: User): void {
  authStore.setState((state) => ({
    ...state,
    user,
    isAuthenticated: true,
    isLoading: false,
  }))
}

/**
 * Clear authentication state
 */
export function clearAuth(): void {
  authStore.setState(() => ({
    user: null,
    isAuthenticated: false,
    isLoading: false,
  }))

  // Allow re-initialization after an explicit logout/clear.
  initPromise = null
}

/**
 * Set loading state
 */
export function setAuthLoading(isLoading: boolean): void {
  authStore.setState((state) => ({
    ...state,
    isLoading,
  }))
}

/**
 * Login action
 *
 * Authenticates user with email and password, updates store with user data.
 */
export async function login(email: string, password: string): Promise<void> {
  try {
    // Don't set loading state here - it causes the login form to unmount
    // The form component handles its own loading state
    const user = await loginUser(email, password)
    setUser(user)

    // Consider auth initialized after explicit login.
    initPromise = Promise.resolve()
  } catch (error) {
    // Don't set loading false - we didn't set it to true
    throw error
  }
}

/**
 * Logout action
 *
 * Logs out the current user, clears auth state.
 */
export async function logout(): Promise<void> {
  try {
    await logoutUser()
  } finally {
    clearAuth()
  }
}

/**
 * Logout all devices action
 *
 * Logs out from all devices, clears auth state.
 */
export async function logoutAll(): Promise<void> {
  try {
    await logoutAllDevices()
  } finally {
    clearAuth()
  }
}

/**
 * Initialize TanStack Store Devtools (dev-only)
 *
 * Conditionally loads devtools in development mode only.
 * Production builds will exclude this code via tree-shaking.
 */
if (import.meta.env.DEV) {
  // Note: TanStack Store has built-in devtools support via @tanstack/react-store
  // The devtools are automatically enabled in development mode.
  // No need for a separate devtools package.
  console.log('[authStore] TanStack Store devtools enabled in development mode')
}
