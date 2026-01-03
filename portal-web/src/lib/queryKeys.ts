/**
 * Centralized Query Key Factory
 *
 * Provides type-safe, hierarchical query keys for TanStack Query with documented staleTime strategy.
 * Keys are nested to enable efficient invalidation (e.g., invalidate all customer queries at once).
 *
 * @example
 * ```ts
 * // Invalidate all customer queries
 * queryClient.invalidateQueries({ queryKey: queryKeys.customers.all })
 *
 * // Invalidate specific customer detail
 * queryClient.invalidateQueries({ queryKey: queryKeys.customers.detail(customerId) })
 *
 * // Invalidate all business-scoped queries on business switch
 * invalidateBusinessScopedQueries(queryClient, businessDescriptor)
 * ```
 */

/**
 * StaleTime Strategy (how long data is considered fresh):
 *
 * - User profile: 5 minutes (rarely changes during session)
 * - Businesses list: 5 minutes (semi-static, only changes when creating new business)
 * - Selected business: 5 minutes (semi-static)
 * - Customers list: 30 seconds (business-critical, moderate update frequency)
 * - Customer details: 30 seconds (business-critical)
 * - Orders: 15 seconds (high priority, frequently updated)
 * - Inventory: 1 minute (moderate priority, moderate update frequency)
 * - Analytics: 5 minutes (computed data, acceptable staleness)
 * - Metadata (countries/currencies): 24 hours (static reference data)
 * - Onboarding session: 0 (always fresh, time-sensitive)
 */

export const STALE_TIME = {
  NEVER: 0,
  FIFTEEN_SECONDS: 15 * 1000,
  THIRTY_SECONDS: 30 * 1000,
  ONE_MINUTE: 60 * 1000,
  FIVE_MINUTES: 5 * 60 * 1000,
  ONE_HOUR: 60 * 60 * 1000,
  TWENTY_FOUR_HOURS: 24 * 60 * 60 * 1000,
} as const

// Legacy export for backwards compatibility
export const staleTime = {
  user: STALE_TIME.FIVE_MINUTES,
  businesses: STALE_TIME.FIVE_MINUTES,
  customers: STALE_TIME.THIRTY_SECONDS,
  orders: STALE_TIME.FIFTEEN_SECONDS,
  inventory: STALE_TIME.ONE_MINUTE,
  analytics: STALE_TIME.FIVE_MINUTES,
  metadata: STALE_TIME.TWENTY_FOUR_HOURS,
} as const

/**
 * User queries
 * StaleTime: 5 minutes
 */
export const user = {
  all: ['user'] as const,
  current: () => [...user.all, 'current'] as const,
  profile: () => [...user.all, 'profile'] as const,
  businesses: () => [...user.all, 'businesses'] as const,
} as const

/**
 * Business queries
 * StaleTime: 5 minutes (semi-static data)
 */
export const businesses = {
  all: ['businesses'] as const,
  list: () => [...businesses.all, 'list'] as const,
  detail: (descriptor: string) =>
    [...businesses.all, 'detail', descriptor] as const,
} as const

// Legacy export for backwards compatibility
export const business = {
  all: ['business'] as const,
  detail: (businessId: string) => [...business.all, businessId] as const,
  metadata: (businessId: string) =>
    [...business.all, businessId, 'metadata'] as const,
} as const

/**
 * Customer queries (business-scoped)
 * StaleTime: 30 seconds (business-critical data)
 * Invalidated on business switch
 */
export const customers = {
  all: ['customers'] as const,
  businessScoped: true,
  lists: () => [...customers.all, 'list'] as const,
  list: (
    businessDescriptor: string,
    filters?: {
      search?: string
      page?: number
      pageSize?: number
      orderBy?: Array<string>
    },
  ) => [...customers.lists(), businessDescriptor, filters] as const,
  details: () => [...customers.all, 'detail'] as const,
  detail: (businessDescriptor: string, customerId: string) =>
    [...customers.details(), businessDescriptor, customerId] as const,
} as const

/**
 * Address queries (business-scoped, nested under customer)
 * StaleTime: 30 seconds (business-critical data)
 * Invalidated on business switch
 */
export const addresses = {
  all: ['addresses'] as const,
  businessScoped: true,
  lists: () => [...addresses.all, 'list'] as const,
  list: (businessDescriptor: string, customerId: string) =>
    [...addresses.lists(), businessDescriptor, customerId] as const,
  details: () => [...addresses.all, 'detail'] as const,
  detail: (businessDescriptor: string, customerId: string, addressId: string) =>
    [
      ...addresses.details(),
      businessDescriptor,
      customerId,
      addressId,
    ] as const,
} as const

/**
 * Order queries (business-scoped)
 * StaleTime: 15 seconds (high priority, frequently updated)
 * Invalidated on business switch
 */
export const orders = {
  all: ['orders'] as const,
  businessScoped: true,
  lists: () => [...orders.all, 'list'] as const,
  list: (
    businessDescriptor: string,
    filters?: {
      status?: string
      customerId?: string
      dateFrom?: string
      dateTo?: string
      page?: number
      limit?: number
    },
  ) => [...orders.lists(), businessDescriptor, filters] as const,
  details: () => [...orders.all, 'detail'] as const,
  detail: (businessDescriptor: string, orderId: string) =>
    [...orders.details(), businessDescriptor, orderId] as const,
} as const

/**
 * Inventory queries (business-scoped)
 * StaleTime: 1 minute (moderate priority)
 * Invalidated on business switch
 */
export const inventory = {
  all: ['inventory'] as const,
  businessScoped: true,
  lists: () => [...inventory.all, 'list'] as const,
  list: (
    businessDescriptor: string,
    filters?: {
      search?: string
      lowStock?: boolean
      categoryId?: string
      stockStatus?: string
      page?: number
      limit?: number
      orderBy?: Array<string>
    },
  ) => [...inventory.lists(), businessDescriptor, filters] as const,
  details: () => [...inventory.all, 'detail'] as const,
  detail: (businessDescriptor: string, productId: string) =>
    [...inventory.details(), businessDescriptor, productId] as const,
} as const

/**
 * Analytics queries (business-scoped)
 * StaleTime: 5 minutes (computed data, acceptable staleness)
 * Invalidated on business switch
 */
export const analytics = {
  all: ['analytics'] as const,
  businessScoped: true,
  dashboard: (
    businessDescriptor: string,
    dateRange?: { from: string; to: string },
  ) => [...analytics.all, 'dashboard', businessDescriptor, dateRange] as const,
  revenue: (
    businessDescriptor: string,
    dateRange?: { from: string; to: string },
  ) => [...analytics.all, 'revenue', businessDescriptor, dateRange] as const,
  topProducts: (
    businessDescriptor: string,
    dateRange?: { from: string; to: string },
  ) =>
    [...analytics.all, 'topProducts', businessDescriptor, dateRange] as const,
  topCustomers: (
    businessDescriptor: string,
    dateRange?: { from: string; to: string },
  ) =>
    [...analytics.all, 'topCustomers', businessDescriptor, dateRange] as const,
} as const

/**
 * Metadata query keys (global, not business-scoped)
 * StaleTime: 24 hours (static reference data)
 */
export const metadata = {
  all: ['metadata'] as const,
  countries: () => [...metadata.all, 'countries'] as const,
  currencies: () => [...metadata.all, 'currencies'] as const,
  plans: () => [...metadata.all, 'plans'] as const,
} as const

/**
 * Onboarding queries (session-specific)
 * StaleTime: 0 (always fresh, time-sensitive)
 */
export const onboarding = {
  all: ['onboarding'] as const,
  session: (sessionToken: string | null) =>
    [...onboarding.all, 'session', sessionToken] as const,
  plans: () => [...onboarding.all, 'plans'] as const,
} as const

/**
 * Export all query key factories
 */
export const queryKeys = {
  user,
  businesses,
  business, // Legacy compatibility
  customers,
  addresses,
  orders,
  inventory,
  analytics,
  metadata,
  onboarding,
  staleTime,
} as const

/**
 * Helper to identify if a query key belongs to business-scoped data
 *
 * Used by invalidateBusinessScopedQueries to determine which queries
 * should be cleared when switching businesses.
 */
export function isBusinessScopedQuery(
  queryKey: ReadonlyArray<unknown>,
): boolean {
  if (!Array.isArray(queryKey) || queryKey.length === 0) return false

  const rootKey = queryKey[0]
  return (
    rootKey === 'customers' ||
    rootKey === 'addresses' ||
    rootKey === 'orders' ||
    rootKey === 'inventory' ||
    rootKey === 'analytics'
  )
}
