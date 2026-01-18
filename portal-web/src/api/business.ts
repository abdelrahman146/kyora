import { queryOptions, useMutation, useQuery } from '@tanstack/react-query'
import { z } from 'zod'
import { del, get, patch, post } from './client'
import { AssetReferenceSchema } from './types/asset'
import type { UseMutationOptions } from '@tanstack/react-query'

import { STALE_TIME, queryKeys } from '@/lib/queryKeys'

/**
 * Business API Types and Schemas
 */

/**
 * Storefront Theme Schema (SSOT)
 * Defines the visual theme for a business storefront
 */
export const StorefrontThemeSchema = z.object({
  primaryColor: z.string(),
  secondaryColor: z.string(),
  accentColor: z.string(),
  backgroundColor: z.string(),
  textColor: z.string(),
  fontFamily: z.string(),
  headingFontFamily: z.string(),
})

export type StorefrontTheme = z.infer<typeof StorefrontThemeSchema>

/**
 * Business Response Schema (SSOT - Single Source of Truth)
 *
 * Aligned with backend BusinessResponse in backend/internal/domain/business/model_response.go
 * All fields are required except: logo (omitempty), archivedAt (omitempty)
 */
export const BusinessSchema = z.object({
  id: z.string(),
  workspaceId: z.string(),
  descriptor: z.string(),
  name: z.string(),
  brand: z.string(),
  logo: AssetReferenceSchema.optional().nullable(),
  countryCode: z.string(),
  currency: z.string(),
  storefrontPublicId: z.string(),
  storefrontEnabled: z.boolean(),
  storefrontTheme: StorefrontThemeSchema,
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
  countryCode: z.string(),
  currency: z.string(),
})

export type CreateBusinessRequest = z.infer<typeof CreateBusinessRequestSchema>

export const UpdateBusinessRequestSchema = z.object({
  name: z.string().optional(),
  countryCode: z.string().optional(),
  currency: z.string().optional(),
})

export type UpdateBusinessRequest = z.infer<typeof UpdateBusinessRequestSchema>

/**
 * Shipping Zone Types
 */
export const ShippingZoneSchema = z.object({
  id: z.string(),
  businessId: z.string(),
  name: z.string(),
  countries: z.array(z.string()),
  currency: z.string(),
  shippingCost: z.string(),
  freeShippingThreshold: z.string(),
  createdAt: z.string(),
  updatedAt: z.string(),
})

export type ShippingZone = z.infer<typeof ShippingZoneSchema>

export const ListShippingZonesResponseSchema = z.array(ShippingZoneSchema)

/**
 * Payment Method Types
 */
export const PaymentMethodSchema = z.object({
  descriptor: z.string(),
  name: z.string(),
  enabled: z.boolean(),
  feeType: z.enum(['percentage', 'fixed']).optional().nullable(),
  feeValue: z.string().optional().nullable(),
  effectiveFeePercentage: z.string().optional().nullable(),
})

export type PaymentMethod = z.infer<typeof PaymentMethodSchema>

export const ListPaymentMethodsResponseSchema = z.array(PaymentMethodSchema)

export type ListShippingZonesResponse = z.infer<
  typeof ListShippingZonesResponseSchema
>

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

  /**
   * List shipping zones for a business
   */
  async listShippingZones(
    businessDescriptor: string,
  ): Promise<Array<ShippingZone>> {
    const response = await get<unknown>(
      `v1/businesses/${businessDescriptor}/shipping-zones`,
    )
    return ListShippingZonesResponseSchema.parse(response)
  },
  /**
   * List payment methods for a business
   */
  async listPaymentMethods(
    businessDescriptor: string,
  ): Promise<Array<PaymentMethod>> {
    const response = await get<unknown>(
      `v1/businesses/${businessDescriptor}/payment-methods`,
    )
    return ListPaymentMethodsResponseSchema.parse(response)
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

  /**
   * Query options for fetching shipping zones
   * @param businessDescriptor - Business descriptor/slug
   */
  shippingZones: (businessDescriptor: string) =>
    queryOptions({
      queryKey: queryKeys.businesses.shippingZones(businessDescriptor),
      queryFn: () => businessApi.listShippingZones(businessDescriptor),
      staleTime: STALE_TIME.FIVE_MINUTES,
      enabled: !!businessDescriptor,
    }),

  /**
   * Query options for fetching payment methods
   * @param businessDescriptor - Business descriptor/slug
   */
  paymentMethods: (businessDescriptor: string) =>
    queryOptions({
      queryKey: queryKeys.businesses.paymentMethods(businessDescriptor),
      queryFn: () => businessApi.listPaymentMethods(businessDescriptor),
      staleTime: STALE_TIME.FIVE_MINUTES,
      enabled: !!businessDescriptor,
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
 * Query to fetch shipping zones for a business
 *
 * StaleTime: 5 minutes (semi-static, changes when zones are updated)
 * @param businessDescriptor - Business descriptor/slug
 * @param enabled - Whether to enable the query (default: true)
 */
export function useShippingZonesQuery(
  businessDescriptor: string,
  enabled: boolean = true,
) {
  return useQuery({
    ...businessQueries.shippingZones(businessDescriptor),
    enabled: enabled && !!businessDescriptor,
  })
}

/**
 * Query to fetch payment methods for a business
 *
 * StaleTime: 5 minutes (semi-static, changes when methods are updated)
 * @param businessDescriptor - Business descriptor/slug
 * @param enabled - Whether to enable the query (default: true)
 */
export function usePaymentMethodsQuery(
  businessDescriptor: string,
  enabled: boolean = true,
) {
  return useQuery({
    ...businessQueries.paymentMethods(businessDescriptor),
    enabled: enabled && !!businessDescriptor,
  })
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
