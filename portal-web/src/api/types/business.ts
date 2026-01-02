import { z } from 'zod'
import {
  
  AssetMetadataSchema,
  
  AssetReferenceSchema
} from './asset'
import type {AssetMetadata, AssetReference} from './asset';

/**
 * Business API types based on backend swagger.json
 */

// Re-export asset types for backward compatibility
export { AssetMetadataSchema, AssetReferenceSchema }
export type { AssetMetadata, AssetReference }

// Storefront Theme Schema
export const StorefrontThemeSchema = z.object({
  primaryColor: z.string().optional(),
  secondaryColor: z.string().optional(),
  accentColor: z.string().optional(),
  fontFamily: z.string().optional(),
})

export type StorefrontTheme = z.infer<typeof StorefrontThemeSchema>

// Business Response Schema
export const BusinessSchema = z.object({
  id: z.string(),
  workspaceId: z.string(),
  name: z.string(),
  descriptor: z.string(),
  brand: z.string().optional(),
  logo: AssetReferenceSchema.optional().nullable(),
  phoneNumber: z.string().optional(),
  whatsappNumber: z.string().optional(),
  supportEmail: z.string().optional(),
  websiteUrl: z.string().optional(),
  facebookUrl: z.string().optional(),
  instagramUrl: z.string().optional(),
  xUrl: z.string().optional(),
  tiktokUrl: z.string().optional(),
  snapchatUrl: z.string().optional(),
  address: z.string().optional(),
  countryCode: z.string().optional(),
  currency: z.string().optional(),
  vatRate: z.string().optional(),
  safetyBuffer: z.string().optional(),
  establishedAt: z.string().optional(),
  storefrontEnabled: z.boolean().optional(),
  storefrontPublicId: z.string().optional(),
  storefrontTheme: StorefrontThemeSchema.optional(),
  createdAt: z.string(),
  updatedAt: z.string(),
  archivedAt: z.string().optional().nullable(),
})

export type Business = z.infer<typeof BusinessSchema>

// List Businesses Response Schema
export const ListBusinessesResponseSchema = z.object({
  businesses: z.array(BusinessSchema),
})

export type ListBusinessesResponse = z.infer<
  typeof ListBusinessesResponseSchema
>

// Create Business Input Schema
export const CreateBusinessInputSchema = z.object({
  name: z.string().min(1, 'Business name is required'),
  descriptor: z.string().min(1, 'Business descriptor is required'),
  brand: z.string().optional(),
  countryCode: z.string().optional(),
  currency: z.string().optional(),
})

export type CreateBusinessInput = z.infer<typeof CreateBusinessInputSchema>
