/**
 * Centralized Query Key Factory
 *
 * Provides consistent query key generation with business scoping and tiered staleTime configuration.
 * All query keys follow the pattern: [scope, resource, ...filters]
 *
 * Tiered staleTime Strategy (defined in main.tsx QueryClient defaults):
 * - User Profile: 5 minutes (changes infrequently)
 * - Businesses: 5 minutes (changes infrequently)
 * - Customers: 30 seconds (moderate change frequency)
 * - Orders: 15 seconds (frequent changes)
 * - Inventory: 1 minute (moderate change frequency)
 * - Analytics: 5 minutes (calculated data, can be cached longer)
 * - Metadata: 24 hours (rarely changes)
 *
 * @example
 * ```tsx
 * // User-scoped query
 * const { data } = useQuery({
 *   queryKey: queryKeys.user.profile(),
 *   queryFn: () => userApi.getProfile(),
 *   staleTime: queryKeys.staleTime.user,
 * });
 *
 * // Business-scoped query
 * const { data } = useQuery({
 *   queryKey: queryKeys.customers.list(businessId, filters),
 *   queryFn: () => customerApi.list(businessId, filters),
 *   staleTime: queryKeys.staleTime.customers,
 * });
 * ```
 */

/**
 * Tiered staleTime configuration (in milliseconds)
 * Used consistently across all queries
 */
export const staleTime = {
  user: 1000 * 60 * 5, // 5 minutes
  businesses: 1000 * 60 * 5, // 5 minutes
  customers: 1000 * 30, // 30 seconds
  orders: 1000 * 15, // 15 seconds
  inventory: 1000 * 60, // 1 minute
  analytics: 1000 * 60 * 5, // 5 minutes
  metadata: 1000 * 60 * 60 * 24, // 24 hours
} as const

/**
 * User-scoped query keys
 */
export const user = {
  all: ['user'] as const,
  profile: () => [...user.all, 'profile'] as const,
  businesses: () => [...user.all, 'businesses'] as const,
} as const

/**
 * Business-scoped query keys
 * Use businessScoped() helper to automatically prefix with business ID
 */
export const business = {
  all: ['business'] as const,
  detail: (businessId: string) => [...business.all, businessId] as const,
  metadata: (businessId: string) =>
    [...business.all, businessId, 'metadata'] as const,
} as const

/**
 * Customer-scoped query keys
 */
export const customers = {
  all: ['customers'] as const,
  lists: () => [...customers.all, 'list'] as const,
  list: (
    businessId: string,
    filters?: {
      search?: string
      page?: number
      limit?: number
    },
  ) => [...customers.lists(), businessId, filters] as const,
  detail: (businessId: string, customerId: string) =>
    [...customers.all, businessId, 'detail', customerId] as const,
} as const

/**
 * Order-scoped query keys
 */
export const orders = {
  all: ['orders'] as const,
  lists: () => [...orders.all, 'list'] as const,
  list: (
    businessId: string,
    filters?: {
      status?: string
      customerId?: string
      dateFrom?: string
      dateTo?: string
      page?: number
      limit?: number
    },
  ) => [...orders.lists(), businessId, filters] as const,
  detail: (businessId: string, orderId: string) =>
    [...orders.all, businessId, 'detail', orderId] as const,
} as const

/**
 * Inventory-scoped query keys
 */
export const inventory = {
  all: ['inventory'] as const,
  lists: () => [...inventory.all, 'list'] as const,
  list: (
    businessId: string,
    filters?: {
      search?: string
      lowStock?: boolean
      page?: number
      limit?: number
    },
  ) => [...inventory.lists(), businessId, filters] as const,
  detail: (businessId: string, productId: string) =>
    [...inventory.all, businessId, 'detail', productId] as const,
} as const

/**
 * Analytics-scoped query keys
 */
export const analytics = {
  all: ['analytics'] as const,
  dashboard: (businessId: string, dateRange?: { from: string; to: string }) =>
    [...analytics.all, businessId, 'dashboard', dateRange] as const,
  revenue: (businessId: string, dateRange?: { from: string; to: string }) =>
    [...analytics.all, businessId, 'revenue', dateRange] as const,
  topProducts: (businessId: string, dateRange?: { from: string; to: string }) =>
    [...analytics.all, businessId, 'topProducts', dateRange] as const,
  topCustomers: (
    businessId: string,
    dateRange?: { from: string; to: string },
  ) => [...analytics.all, businessId, 'topCustomers', dateRange] as const,
} as const

/**
 * Metadata query keys (plans, currencies, etc.)
 */
export const metadata = {
  all: ['metadata'] as const,
  plans: () => [...metadata.all, 'plans'] as const,
  currencies: () => [...metadata.all, 'currencies'] as const,
  countries: () => [...metadata.all, 'countries'] as const,
} as const

/**
 * Onboarding query keys
 */
export const onboarding = {
  all: ['onboarding'] as const,
  session: () => [...onboarding.all, 'session'] as const,
} as const

/**
 * Helper to invalidate all queries scoped to a specific business
 * Use this when switching businesses or when business data changes significantly
 *
 * @example
 * ```tsx
 * // Invalidate all business-scoped queries when switching business
 * await queryClient.invalidateQueries({
 *   predicate: (query) => businessScoped.includes(query.queryKey, businessId)
 * });
 * ```
 */
export const businessScoped = {
  /**
   * Check if a query key is scoped to a specific business
   */
  includes: (queryKey: ReadonlyArray<unknown>, businessId: string): boolean => {
    // Business-scoped keys include the businessId as the second element
    // e.g., ['customers', 'list', 'business-123', ...]
    return (
      queryKey.length >= 2 &&
      typeof queryKey[1] === 'object' &&
      queryKey[1] !== null &&
      'businessId' in queryKey[1] &&
      queryKey[1].businessId === businessId
    )
  },

  /**
   * Get all business-scoped resources for a specific business
   * Useful for bulk invalidation
   */
  all: (businessId: string) => [businessId] as const,
}

/**
 * Export all query key factories
 */
export const queryKeys = {
  user,
  business,
  customers,
  orders,
  inventory,
  analytics,
  metadata,
  onboarding,
  businessScoped,
  staleTime,
} as const
