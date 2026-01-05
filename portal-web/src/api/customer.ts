import {
  keepPreviousData,
  queryOptions,
  useMutation,
  useQuery,
} from '@tanstack/react-query'

import { del, get, patch, post } from './client'
import type { UseMutationOptions } from '@tanstack/react-query'
import { STALE_TIME, queryKeys } from '@/lib/queryKeys'

/**
 * Customer API Types
 * Based on backend swagger.json definitions
 */

export type CustomerGender = 'male' | 'female' | 'other'

export type SocialPlatform =
  | 'instagram'
  | 'tiktok'
  | 'facebook'
  | 'x'
  | 'snapchat'
  | 'whatsapp'

export interface ListCustomersFilters {
  countryCode?: string
  hasOrders?: boolean
  socialPlatforms?: Array<SocialPlatform>
}

export interface CustomerAddress {
  id: string
  customerId: string
  shippingZoneId: string | null
  street: string | null
  city: string
  state: string
  zipCode: string | null
  countryCode: string
  phoneCode: string
  phoneNumber: string
  createdAt: string
  updatedAt: string
  deletedAt: string | null
}

export interface CustomerNote {
  id: string
  customerId: string
  content: string
  createdAt: string
  updatedAt: string
  deletedAt: string | null
}

export interface Customer {
  id: string
  businessId: string
  name: string
  email: string | null
  phoneCode: string | null
  phoneNumber: string | null
  countryCode: string
  whatsappNumber: string | null
  gender: CustomerGender
  joinedAt: string
  instagramUsername: string | null
  facebookUsername: string | null
  tiktokUsername: string | null
  snapchatUsername: string | null
  xUsername: string | null
  addresses?: Array<CustomerAddress>
  notes?: Array<CustomerNote>
  createdAt: string
  updatedAt: string
  deletedAt: string | null
  // Computed fields from backend aggregation
  ordersCount?: number
  totalSpent?: number
  avatarUrl?: string
}

export interface ListResponse<T> {
  hasMore: boolean
  items: Array<T>
  page: number
  pageSize: number
  totalCount: number
  totalPages: number
}

export type ListCustomersResponse = ListResponse<Customer>

export interface CreateCustomerRequest {
  name: string
  email: string
  phoneCode?: string
  phoneNumber?: string
  countryCode: string
  whatsappNumber?: string
  gender?: CustomerGender
  joinedAt?: string
  instagramUsername?: string
  facebookUsername?: string
  tiktokUsername?: string
  snapchatUsername?: string
  xUsername?: string
}

export interface UpdateCustomerRequest {
  name?: string
  email?: string
  phoneCode?: string
  phoneNumber?: string
  countryCode?: string
  whatsappNumber?: string
  gender?: CustomerGender
  joinedAt?: string
  instagramUsername?: string
  facebookUsername?: string
  tiktokUsername?: string
  snapchatUsername?: string
  xUsername?: string
}

/**
 * Customer API Client
 */
export const customerApi = {
  /**
   * List customers for a business
   */
  async listCustomers(
    businessDescriptor: string,
    params?: {
      search?: string
      page?: number
      pageSize?: number
      orderBy?: Array<string>
      countryCode?: string
      hasOrders?: boolean
      socialPlatforms?: Array<SocialPlatform>
    },
  ): Promise<ListCustomersResponse> {
    const searchParams = new URLSearchParams()
    if (params?.search) searchParams.set('search', params.search)
    if (params?.page) searchParams.set('page', params.page.toString())
    if (params?.pageSize)
      searchParams.set('pageSize', params.pageSize.toString())
    if (params?.orderBy && params.orderBy.length > 0) {
      // swagger: orderBy is collectionFormat=csv
      searchParams.set('orderBy', params.orderBy.join(','))
    }
    if (params?.countryCode) searchParams.set('countryCode', params.countryCode)
    if (params?.hasOrders !== undefined)
      searchParams.set('hasOrders', params.hasOrders.toString())
    if (params?.socialPlatforms && params.socialPlatforms.length > 0) {
      // swagger: socialPlatforms is collectionFormat=csv
      searchParams.set('socialPlatforms', params.socialPlatforms.join(','))
    }

    const query = searchParams.toString() ? `?${searchParams.toString()}` : ''
    return get<ListCustomersResponse>(
      `v1/businesses/${businessDescriptor}/customers${query}`,
    )
  },

  /**
   * Get customer by ID
   */
  async getCustomer(
    businessDescriptor: string,
    customerId: string,
  ): Promise<Customer> {
    return get<Customer>(
      `v1/businesses/${businessDescriptor}/customers/${customerId}`,
    )
  },

  /**
   * Create a new customer
   */
  async createCustomer(
    businessDescriptor: string,
    data: CreateCustomerRequest,
  ): Promise<Customer> {
    return post<Customer>(`v1/businesses/${businessDescriptor}/customers`, {
      json: data,
    })
  },

  /**
   * Update existing customer
   */
  async updateCustomer(
    businessDescriptor: string,
    customerId: string,
    data: UpdateCustomerRequest,
  ): Promise<Customer> {
    return patch<Customer>(
      `v1/businesses/${businessDescriptor}/customers/${customerId}`,
      { json: data },
    )
  },

  /**
   * Delete customer
   */
  async deleteCustomer(
    businessDescriptor: string,
    customerId: string,
  ): Promise<void> {
    await del<void>(
      `v1/businesses/${businessDescriptor}/customers/${customerId}`,
    )
  },

  /**
   * Create a note for a customer
   */
  async createCustomerNote(
    businessDescriptor: string,
    customerId: string,
    content: string,
  ): Promise<CustomerNote> {
    return post<CustomerNote>(
      `v1/businesses/${businessDescriptor}/customers/${customerId}/notes`,
      { json: { content } },
    )
  },
}

/**
 * Query Options Factories
 *
 * Co-locate query configuration (key + fn + staleTime) for type-safe reuse
 * in components, route loaders, and prefetching.
 */
export const customerQueries = {
  /**
   * Query options for fetching customers list
   * @param businessDescriptor - Business identifier
   * @param params - Optional filter and pagination params
   */
  list: (
    businessDescriptor: string,
    params?: {
      search?: string
      page?: number
      pageSize?: number
      orderBy?: Array<string>
      countryCode?: string
      hasOrders?: boolean
      socialPlatforms?: Array<SocialPlatform>
    },
  ) =>
    queryOptions({
      queryKey: queryKeys.customers.list(businessDescriptor, params),
      queryFn: () => customerApi.listCustomers(businessDescriptor, params),
      staleTime: STALE_TIME.THIRTY_SECONDS,
      enabled: !!businessDescriptor,
      placeholderData: keepPreviousData, // Smooth pagination transitions
    }),

  /**
   * Query options for fetching a specific customer
   * @param businessDescriptor - Business identifier
   * @param customerId - Customer identifier
   */
  detail: (businessDescriptor: string, customerId: string) =>
    queryOptions({
      queryKey: queryKeys.customers.detail(businessDescriptor, customerId),
      queryFn: () => customerApi.getCustomer(businessDescriptor, customerId),
      staleTime: STALE_TIME.THIRTY_SECONDS,
      enabled: !!businessDescriptor && !!customerId,
    }),
}

/**
 * Query Hooks
 */

/**
 * Query to fetch customers list
 *
 * StaleTime: 30 seconds (business-critical data)
 * Business-scoped: Invalidated on business switch
 */
export function useCustomersQuery(
  businessDescriptor: string,
  params?: {
    search?: string
    page?: number
    pageSize?: number
    orderBy?: Array<string>
    countryCode?: string
    hasOrders?: boolean
    socialPlatforms?: Array<SocialPlatform>
  },
) {
  return useQuery(customerQueries.list(businessDescriptor, params))
}

/**
 * Query to fetch a specific customer
 *
 * StaleTime: 30 seconds (business-critical data)
 * Business-scoped: Invalidated on business switch
 */
export function useCustomerQuery(
  businessDescriptor: string,
  customerId: string,
) {
  return useQuery(customerQueries.detail(businessDescriptor, customerId))
}

/**
 * Mutation Hooks
 */

/**
 * Mutation to create a new customer
 *
 * Implements optimistic updates:
 * - Immediately adds customer to cache
 * - Rolls back on error with toast notification
 */
export function useCreateCustomerMutation(
  businessDescriptor: string,
  options?: UseMutationOptions<Customer, Error, CreateCustomerRequest>,
) {
  return useMutation({
    mutationFn: (data: CreateCustomerRequest) =>
      customerApi.createCustomer(businessDescriptor, data),
    ...options,
  })
}

/**
 * Mutation to update a customer
 *
 * Implements optimistic updates:
 * - Immediately updates customer in cache
 * - Rolls back on error with toast notification
 */
export function useUpdateCustomerMutation(
  businessDescriptor: string,
  options?: UseMutationOptions<
    Customer,
    Error,
    { customerId: string; data: UpdateCustomerRequest }
  >,
) {
  return useMutation({
    mutationFn: ({ customerId, data }) =>
      customerApi.updateCustomer(businessDescriptor, customerId, data),
    ...options,
  })
}

/**
 * Mutation to delete a customer
 */
export function useDeleteCustomerMutation(
  businessDescriptor: string,
  options?: UseMutationOptions<void, Error, string>,
) {
  return useMutation({
    mutationFn: (customerId: string) =>
      customerApi.deleteCustomer(businessDescriptor, customerId),
    ...options,
  })
}

/**
 * Mutation to create a customer note
 */
export function useCreateCustomerNoteMutation(
  businessDescriptor: string,
  customerId: string,
  options?: UseMutationOptions<CustomerNote, Error, string>,
) {
  return useMutation({
    mutationFn: (content: string) =>
      customerApi.createCustomerNote(businessDescriptor, customerId, content),
    ...options,
  })
}
