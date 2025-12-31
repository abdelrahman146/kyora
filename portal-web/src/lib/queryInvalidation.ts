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
 * **When to use:**
 * - After mutations that affect multiple business-scoped resources (e.g., deleting a business)
 * - When explicitly "refreshing" all data for the current business (user-triggered)
 *
 * **When NOT to use:**
 * - In route loaders/beforeLoad (causes immediate refetch of data you just loaded)
 * - After every navigation (wastes bandwidth, violates staleTime strategy)
 * - Instead: let TanStack Query's staleTime handle refetching, or invalidate specific queries
 *
 * **Prefer targeted invalidation:**
 * - After creating an order: `invalidateQueries({ queryKey: ['businesses', descriptor, 'orders'] })`
 * - After updating a customer: `invalidateQueries({ queryKey: ['businesses', descriptor, 'customers', customerId] })`
 *
 * @param queryClient - TanStack Query client instance
 * @param businessDescriptor - Optional business descriptor to invalidate specific business queries
 *
 * @example
 * ```ts
 * // In a mutation handler after deleting a business
 * onSuccess: () => {
 *   invalidateBusinessScopedQueries(queryClient)
 * }
 *
 * // For specific business refresh (e.g., user clicks "Refresh" button)
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
