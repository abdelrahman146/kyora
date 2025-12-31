import { createRouter } from '@tanstack/react-router'
import { routeTree } from './routeTree.gen'
import type { QueryClient } from '@tanstack/react-query'
import type { User } from '@/api/types/auth'
import { authStore } from '@/stores/authStore'

// Import the generated route tree

/**
 * Router Context Type
 *
 * Provides auth state and queryClient to all routes via context.
 */
export interface RouterContext {
  auth: {
    user: User | null
    isAuthenticated: boolean
    isLoading: boolean
  }
  queryClient: QueryClient
}

/**
 * Create a new router instance
 *
 * Configured with:
 * - defaultPreload: 'intent' for hover-based prefetching
 * - Context with auth state from authStore and queryClient
 */
export const getRouter = (queryClient: QueryClient) => {
  const router = createRouter({
    routeTree,
    context: {
      auth: authStore.state,
      queryClient,
    } as RouterContext,
    defaultPreload: 'intent',
  })

  return router
}
