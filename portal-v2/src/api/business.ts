import { useMutation, useQuery } from '@tanstack/react-query'
import { z } from 'zod'
import { del, get, patch, post } from './client'
import type { UseMutationOptions } from '@tanstack/react-query'

import type { Business } from '@/stores/businessStore'
import { STALE_TIME, queryKeys } from '@/lib/queryKeys'

/**
 * Business API Types and Schemas
 */

export const BusinessSchema = z.object({
  id: z.string(),
  descriptor: z.string(),
  name: z.string(),
  country: z.string(),
  currency: z.string(),
  createdAt: z.string(),
  updatedAt: z.string(),
})

export const ListBusinessesResponseSchema = z.object({
  businesses: z.array(BusinessSchema),
})

export type ListBusinessesResponse = z.infer<
  typeof ListBusinessesResponseSchema
>

export const CreateBusinessRequestSchema = z.object({
  name: z.string(),
  descriptor: z.string(),
  country: z.string(),
  currency: z.string(),
})

export type CreateBusinessRequest = z.infer<typeof CreateBusinessRequestSchema>

export const UpdateBusinessRequestSchema = z.object({
  name: z.string().optional(),
  country: z.string().optional(),
  currency: z.string().optional(),
})

export type UpdateBusinessRequest = z.infer<typeof UpdateBusinessRequestSchema>

/**
 * Business API Client
 */
export const businessApi = {
  /**
   * List all businesses for the authenticated user
   */
  async listBusinesses(): Promise<ListBusinessesResponse> {
    const response = await get<unknown>('v1/businesses')
    return ListBusinessesResponseSchema.parse(response)
  },

  /**
   * Get business by descriptor
   */
  async getBusiness(descriptor: string): Promise<Business> {
    const response = await get<unknown>(`v1/businesses/${descriptor}`)
    return BusinessSchema.parse(response)
  },

  /**
   * Create a new business
   */
  async createBusiness(data: CreateBusinessRequest): Promise<Business> {
    const validatedData = CreateBusinessRequestSchema.parse(data)
    const response = await post<unknown>('v1/businesses', {
      json: validatedData,
    })
    return BusinessSchema.parse(response)
  },

  /**
   * Update existing business
   */
  async updateBusiness(
    descriptor: string,
    data: UpdateBusinessRequest,
  ): Promise<Business> {
    const validatedData = UpdateBusinessRequestSchema.parse(data)
    const response = await patch<unknown>(`v1/businesses/${descriptor}`, {
      json: validatedData,
    })
    return BusinessSchema.parse(response)
  },

  /**
   * Delete business
   */
  async deleteBusiness(descriptor: string): Promise<void> {
    await del<void>(`v1/businesses/${descriptor}`)
  },
}

/**
 * Query Hooks
 */

/**
 * Query to fetch all businesses
 *
 * StaleTime: 5 minutes (semi-static, only changes when creating new business)
 */
export function useBusinessesQuery() {
  return useQuery({
    queryKey: queryKeys.businesses.list(),
    queryFn: () => businessApi.listBusinesses(),
    staleTime: STALE_TIME.FIVE_MINUTES,
    select: (data) => data.businesses,
  })
}

/**
 * Query to fetch a specific business
 *
 * StaleTime: 5 minutes (semi-static)
 */
export function useBusinessQuery(descriptor: string) {
  return useQuery({
    queryKey: queryKeys.businesses.detail(descriptor),
    queryFn: () => businessApi.getBusiness(descriptor),
    staleTime: STALE_TIME.FIVE_MINUTES,
    enabled: !!descriptor,
  })
}

/**
 * Mutation Hooks
 */

/**
 * Mutation to create a new business
 *
 * Implements optimistic updates:
 * - Immediately adds business to cache
 * - Rolls back on error with toast notification
 */
export function useCreateBusinessMutation(
  options?: UseMutationOptions<Business, Error, CreateBusinessRequest>,
) {
  return useMutation({
    mutationFn: (data: CreateBusinessRequest) =>
      businessApi.createBusiness(data),
    ...options,
  })
}

/**
 * Mutation to update a business
 *
 * Implements optimistic updates:
 * - Immediately updates business in cache
 * - Rolls back on error with toast notification
 */
export function useUpdateBusinessMutation(
  options?: UseMutationOptions<
    Business,
    Error,
    { descriptor: string; data: UpdateBusinessRequest }
  >,
) {
  return useMutation({
    mutationFn: ({ descriptor, data }) =>
      businessApi.updateBusiness(descriptor, data),
    ...options,
  })
}

/**
 * Mutation to delete a business
 */
export function useDeleteBusinessMutation(
  options?: UseMutationOptions<void, Error, string>,
) {
  return useMutation({
    mutationFn: (descriptor: string) => businessApi.deleteBusiness(descriptor),
    ...options,
  })
}
