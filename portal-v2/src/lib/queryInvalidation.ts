import { isBusinessScopedQuery } from './queryKeys'
import type { QueryClient } from '@tanstack/react-query'

/**
 * Query Invalidation Utilities
 *
 * Helpers for invalidating queries based on business context.
 */

/**
 * Invalidate all business-scoped queries
 *
 * Called when switching businesses to clear all data from the previous business.
 * Only invalidates queries marked as business-scoped in queryKeys factory.
 *
 * @param queryClient - TanStack Query client instance
 * @param businessDescriptor - Optional business descriptor to invalidate specific business queries
 *
 * @example
 * ```ts
 * // Clear all business-scoped queries (any business)
 * invalidateBusinessScopedQueries(queryClient)
 *
 * // Clear queries for specific business
 * invalidateBusinessScopedQueries(queryClient, 'my-shop')
 * ```
 */
export function invalidateBusinessScopedQueries(
  queryClient: QueryClient,
  businessDescriptor?: string,
): void {
  const cache = queryClient.getQueryCache()
  const queries = cache.getAll()

  for (const query of queries) {
    const queryKey = query.queryKey

    // Check if this query is business-scoped
    if (!isBusinessScopedQuery(queryKey)) continue

    // If businessDescriptor provided, only invalidate queries for that business
    if (businessDescriptor && !queryKey.includes(businessDescriptor)) continue

    // Invalidate the query
    void queryClient.invalidateQueries({ queryKey })
  }
}

/**
 * Remove all business-scoped queries from cache
 *
 * More aggressive than invalidate - completely removes queries.
 * Use when user logs out or business is deleted.
 *
 * @param queryClient - TanStack Query client instance
 */
export function removeBusinessScopedQueries(queryClient: QueryClient): void {
  const cache = queryClient.getQueryCache()
  const queries = cache.getAll()

  for (const query of queries) {
    if (isBusinessScopedQuery(query.queryKey)) {
      queryClient.removeQueries({ queryKey: query.queryKey })
    }
  }
}
