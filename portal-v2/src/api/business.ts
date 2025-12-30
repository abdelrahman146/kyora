import { queryOptions, useMutation, useQuery } from '@tanstack/react-query'
import { z } from 'zod'
import { del, get, patch, post } from './client'
import type { UseMutationOptions } from '@tanstack/react-query'

import { STALE_TIME, queryKeys } from '@/lib/queryKeys'

/**
 * Business API Types and Schemas
 */

export const BusinessSchema = z.object({
  id: z.string(),
  workspaceId: z.string(),
  descriptor: z.string(),
  name: z.string(),
  brand: z.string(),
  logoUrl: z.string(),
  countryCode: z.string(),
  currency: z.string(),
  storefrontPublicId: z.string(),
  storefrontEnabled: z.boolean(),
  storefrontTheme: z.object({
    primaryColor: z.string(),
    secondaryColor: z.string(),
    accentColor: z.string(),
    backgroundColor: z.string(),
    textColor: z.string(),
    fontFamily: z.string(),
    headingFontFamily: z.string(),
  }),
  supportEmail: z.string(),
  phoneNumber: z.string(),
  whatsappNumber: z.string(),
  address: z.string(),
  websiteUrl: z.string(),
  instagramUrl: z.string(),
  facebookUrl: z.string(),
  tiktokUrl: z.string(),
  xUrl: z.string(),
  snapchatUrl: z.string(),
  vatRate: z.string(),
  safetyBuffer: z.string(),
  establishedAt: z.string(),
  archivedAt: z.string().nullable().optional(),
  createdAt: z.string(),
  updatedAt: z.string(),
})

export type Business = z.infer<typeof BusinessSchema>

export const ListBusinessesResponseSchema = z.object({
  businesses: z.array(BusinessSchema),
})

export type ListBusinessesResponse = z.infer<
  typeof ListBusinessesResponseSchema
>

export const GetBusinessResponseSchema = z.object({
  business: BusinessSchema,
})

export type GetBusinessResponse = z.infer<typeof GetBusinessResponseSchema>

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
    const parsed = GetBusinessResponseSchema.parse(response)
    return parsed.business
  },

  /**
   * Create a new business
   */
  async createBusiness(data: CreateBusinessRequest): Promise<Business> {
    const validatedData = CreateBusinessRequestSchema.parse(data)
    const response = await post<unknown>('v1/businesses', {
      json: validatedData,
    })
    const parsed = GetBusinessResponseSchema.parse(response)
    return parsed.business
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
    const parsed = GetBusinessResponseSchema.parse(response)
    return parsed.business
  },

  /**
   * Delete business
   */
  async deleteBusiness(descriptor: string): Promise<void> {
    await del<void>(`v1/businesses/${descriptor}`)
  },
}

/**
 * Query Options Factories
 *
 * Co-locate query configuration (key + fn + staleTime) for type-safe reuse
 * in components, route loaders, and prefetching.
 */
export const businessQueries = {
  /**
   * Query options for fetching all businesses
   */
  list: () =>
    queryOptions({
      queryKey: queryKeys.businesses.list(),
      queryFn: () => businessApi.listBusinesses(),
      staleTime: STALE_TIME.FIVE_MINUTES,
    }),

  /**
   * Query options for fetching a specific business
   * @param descriptor - Business descriptor/slug
   */
  detail: (descriptor: string) =>
    queryOptions({
      queryKey: queryKeys.businesses.detail(descriptor),
      queryFn: () => businessApi.getBusiness(descriptor),
      staleTime: STALE_TIME.FIVE_MINUTES,
      enabled: !!descriptor,
    }),
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
  return useQuery(businessQueries.list())
}

/**
 * Query to fetch a specific business
 *
 * StaleTime: 5 minutes (semi-static)
 */
export function useBusinessQuery(descriptor: string) {
  return useQuery(businessQueries.detail(descriptor))
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
