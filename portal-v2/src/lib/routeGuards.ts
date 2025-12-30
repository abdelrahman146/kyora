import { redirect } from '@tanstack/react-router'
import { authStore } from '@/stores/authStore'

/**
 * Route Guard: Require Authentication
 *
 * Use this in `beforeLoad` hooks to protect routes that require authentication.
 * Redirects to login if user is not authenticated, preserving the intended destination.
 *
 * @example
 * ```tsx
 * export const Route = createFileRoute('/business/$businessDescriptor/')({
 *   beforeLoad: requireAuth,
 *   component: DashboardPage,
 * })
 * ```
 */
export function requireAuth() {
  const { isAuthenticated } = authStore.state

  if (!isAuthenticated) {
    // Redirect to login with return URL
    throw redirect({
      to: '/auth/login',
      search: {
        redirect: window.location.pathname,
      },
    })
  }
}

/**
 * Route Guard: Redirect if Authenticated
 *
 * Use this in `beforeLoad` hooks for auth pages (login, register, etc.)
 * Redirects to home if user is already authenticated.
 *
 * @example
 * ```tsx
 * export const Route = createFileRoute('/auth/login')({
 *   beforeLoad: redirectIfAuthenticated,
 *   component: LoginPage,
 * })
 * ```
 */
export function redirectIfAuthenticated() {
  const { isAuthenticated } = authStore.state

  if (isAuthenticated) {
    // Redirect to home if already logged in
    throw redirect({
      to: '/',
    })
  }
}
