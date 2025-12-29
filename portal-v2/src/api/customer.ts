import { useMutation, useQuery } from '@tanstack/react-query'
import { z } from 'zod'
import { del, get, patch, post } from './client'
import type { UseMutationOptions } from '@tanstack/react-query'

import type { PaginatedResponse } from './types/common'
import { STALE_TIME, queryKeys } from '@/lib/queryKeys'

/**
 * Customer API Types and Schemas
 */

export const CustomerSchema = z.object({
  id: z.string(),
  businessId: z.string(),
  fullName: z.string(),
  email: z.string().email().optional().nullable(),
  phoneNumber: z.string().optional().nullable(),
  phonePrefix: z.string().optional().nullable(),
  address: z.string().optional().nullable(),
  city: z.string().optional().nullable(),
  country: z.string().optional().nullable(),
  instagramHandle: z.string().optional().nullable(),
  facebookHandle: z.string().optional().nullable(),
  notes: z.string().optional().nullable(),
  totalOrders: z.number().default(0),
  totalSpent: z.number().default(0),
  createdAt: z.string(),
  updatedAt: z.string(),
})

export type Customer = z.infer<typeof CustomerSchema>

export const ListCustomersResponseSchema = z.object({
  customers: z.array(CustomerSchema),
  pagination: z.object({
    page: z.number(),
    limit: z.number(),
    total: z.number(),
    totalPages: z.number(),
  }),
})

export type ListCustomersResponse = z.infer<typeof ListCustomersResponseSchema>

export const CreateCustomerRequestSchema = z.object({
  fullName: z.string().min(1, 'الاسم مطلوب'),
  email: z.string().email('البريد الإلكتروني غير صالح').optional(),
  phoneNumber: z.string().optional(),
  phonePrefix: z.string().optional(),
  address: z.string().optional(),
  city: z.string().optional(),
  country: z.string().optional(),
  instagramHandle: z.string().optional(),
  facebookHandle: z.string().optional(),
  notes: z.string().optional(),
})

export type CreateCustomerRequest = z.infer<typeof CreateCustomerRequestSchema>

export const UpdateCustomerRequestSchema = CreateCustomerRequestSchema.partial()

export type UpdateCustomerRequest = z.infer<typeof UpdateCustomerRequestSchema>

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
      limit?: number
    },
  ): Promise<ListCustomersResponse> {
    const searchParams = new URLSearchParams()
    if (params?.search) searchParams.set('search', params.search)
    if (params?.page) searchParams.set('page', params.page.toString())
    if (params?.limit) searchParams.set('limit', params.limit.toString())

    const query = searchParams.toString() ? `?${searchParams.toString()}` : ''
    const response = await get<unknown>(
      `v1/businesses/${businessDescriptor}/customers${query}`,
    )
    return ListCustomersResponseSchema.parse(response)
  },

  /**
   * Get customer by ID
   */
  async getCustomer(
    businessDescriptor: string,
    customerId: string,
  ): Promise<Customer> {
    const response = await get<unknown>(
      `v1/businesses/${businessDescriptor}/customers/${customerId}`,
    )
    return CustomerSchema.parse(response)
  },

  /**
   * Create a new customer
   */
  async createCustomer(
    businessDescriptor: string,
    data: CreateCustomerRequest,
  ): Promise<Customer> {
    const validatedData = CreateCustomerRequestSchema.parse(data)
    const response = await post<unknown>(
      `v1/businesses/${businessDescriptor}/customers`,
      { json: validatedData },
    )
    return CustomerSchema.parse(response)
  },

  /**
   * Update existing customer
   */
  async updateCustomer(
    businessDescriptor: string,
    customerId: string,
    data: UpdateCustomerRequest,
  ): Promise<Customer> {
    const validatedData = UpdateCustomerRequestSchema.parse(data)
    const response = await patch<unknown>(
      `v1/businesses/${businessDescriptor}/customers/${customerId}`,
      { json: validatedData },
    )
    return CustomerSchema.parse(response)
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
    limit?: number
  },
) {
  return useQuery({
    queryKey: queryKeys.customers.list(businessDescriptor, params),
    queryFn: () => customerApi.listCustomers(businessDescriptor, params),
    staleTime: STALE_TIME.THIRTY_SECONDS,
    enabled: !!businessDescriptor,
  })
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
  return useQuery({
    queryKey: queryKeys.customers.detail(businessDescriptor, customerId),
    queryFn: () => customerApi.getCustomer(businessDescriptor, customerId),
    staleTime: STALE_TIME.THIRTY_SECONDS,
    enabled: !!businessDescriptor && !!customerId,
  })
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
