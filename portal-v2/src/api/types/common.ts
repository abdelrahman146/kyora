/**
 * Common API Types
 *
 * Shared types used across multiple API endpoints.
 */

/**
 * Pagination metadata
 */
export interface PaginationMeta {
  page: number
  limit: number
  total: number
  totalPages: number
}

/**
 * Generic paginated response
 */
export interface PaginatedResponse<T> {
  data: Array<T>
  pagination: PaginationMeta
}

/**
 * Sort direction
 */
export type SortDirection = 'asc' | 'desc'

/**
 * Filter operator types
 */
export type FilterOperator =
  | 'eq'
  | 'ne'
  | 'gt'
  | 'gte'
  | 'lt'
  | 'lte'
  | 'contains'
  | 'startsWith'
  | 'endsWith'

/**
 * Generic filter
 */
export interface Filter {
  field: string
  operator: FilterOperator
  value: string | number | boolean
}

/**
 * Generic sort
 */
export interface Sort {
  field: string
  direction: SortDirection
}

/**
 * Generic list query params
 */
export interface ListQueryParams {
  page?: number
  limit?: number
  search?: string
  filters?: Array<Filter>
  sort?: Array<Sort>
}
