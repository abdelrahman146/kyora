import { z } from 'zod'

/**
 * Asset API Types and Schemas
 *
 * Shared schemas for asset references used across all resources
 * (business, product, variant, order, etc.)
 */

/**
 * Asset Metadata Schema
 *
 * Optional metadata for assets including alt text, captions, and dimensions
 */
export const AssetMetadataSchema = z.object({
  altText: z.string().optional(),
  caption: z.string().optional(),
  width: z.number().optional(),
  height: z.number().optional(),
})

export type AssetMetadata = z.infer<typeof AssetMetadataSchema>

/**
 * Asset Reference Schema
 *
 * Complete asset reference including URLs, IDs, and metadata.
 * Used across all resources that support asset uploads.
 */
export const AssetReferenceSchema = z.object({
  url: z.string(),
  assetId: z.string().optional(),
  originalUrl: z.string().optional(),
  thumbnailUrl: z.string().optional(),
  thumbnailOriginalUrl: z.string().optional(),
  metadata: AssetMetadataSchema.optional(),
})

export type AssetReference = z.infer<typeof AssetReferenceSchema>
