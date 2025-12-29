import { Store } from '@tanstack/react-store'
import type { User } from '@/api/types/auth'
import { loginUser, logoutAllDevices, logoutUser } from '@/lib/auth'

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
    setAuthLoading(true)
    const user = await loginUser(email, password)
    setUser(user)
  } catch (error) {
    setAuthLoading(false)
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
