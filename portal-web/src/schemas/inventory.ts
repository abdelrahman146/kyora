import { z } from 'zod'

/**
 * Asset Reference Schema (matches backend type)
 */
export const assetReferenceSchema = z.object({
  url: z.string().min(1, 'validation.required'),
  assetId: z.string().min(1, 'validation.required'),
  originalUrl: z.string().optional(),
  thumbnailUrl: z.string().optional(),
  thumbnailOriginalUrl: z.string().optional(),
  metadata: z
    .object({
      altText: z.string().optional(),
      caption: z.string().optional(),
      width: z.number().optional(),
      height: z.number().optional(),
    })
    .optional(),
})

export type AssetReferenceFormData = z.infer<typeof assetReferenceSchema>

/**
 * Product Variant Schema for creation
 */
export const createVariantSchema = z.object({
  code: z.string().min(1, 'validation.required'),
  sku: z.string().optional(),
  photos: z
    .array(assetReferenceSchema)
    .max(10, 'validation.max_photos')
    .optional(),
  costPrice: z
    .string()
    .min(1, 'validation.required')
    .refine(
      (val) => {
        const num = parseFloat(val)
        return !isNaN(num) && num >= 0
      },
      { message: 'validation.positive_number' },
    ),
  salePrice: z
    .string()
    .min(1, 'validation.required')
    .refine(
      (val) => {
        const num = parseFloat(val)
        return !isNaN(num) && num >= 0
      },
      { message: 'validation.positive_number' },
    ),
  stockQuantity: z.coerce
    .number({ message: 'validation.required' })
    .int({ message: 'validation.integer' })
    .min(0, { message: 'validation.min_zero' }),
  stockQuantityAlert: z.coerce
    .number({ message: 'validation.required' })
    .int({ message: 'validation.integer' })
    .min(0, { message: 'validation.min_zero' }),
})

export type CreateVariantFormData = z.infer<typeof createVariantSchema>

/**
 * Single Variant Product Schema (combines product + variant)
 */
export const singleVariantProductSchema = z.object({
  name: z.string().min(1, 'validation.required'),
  description: z.string().optional(),
  photos: z
    .array(assetReferenceSchema)
    .max(10, 'validation.max_photos')
    .optional(),
  categoryId: z.string().min(1, 'validation.required'),
  // Variant fields (embedded seamlessly)
  code: z.string().min(1, 'validation.required'),
  sku: z.string().optional(),
  variantPhotos: z
    .array(assetReferenceSchema)
    .max(10, 'validation.max_photos')
    .optional(),
  costPrice: z
    .string()
    .min(1, 'validation.required')
    .refine(
      (val) => {
        const num = parseFloat(val)
        return !isNaN(num) && num >= 0
      },
      { message: 'validation.positive_number' },
    ),
  salePrice: z
    .string()
    .min(1, 'validation.required')
    .refine(
      (val) => {
        const num = parseFloat(val)
        return !isNaN(num) && num >= 0
      },
      { message: 'validation.positive_number' },
    ),
  stockQuantity: z.coerce
    .number({ message: 'validation.required' })
    .int({ message: 'validation.integer' })
    .min(0, { message: 'validation.min_zero' }),
  stockQuantityAlert: z.coerce
    .number({ message: 'validation.required' })
    .int({ message: 'validation.integer' })
    .min(0, { message: 'validation.min_zero' }),
})

export type SingleVariantProductFormData = z.infer<
  typeof singleVariantProductSchema
>

/**
 * Multi Variant Product Schema (product + array of variants)
 */
export const multiVariantProductSchema = z.object({
  name: z.string().min(1, 'validation.required'),
  description: z.string().optional(),
  photos: z
    .array(assetReferenceSchema)
    .max(10, 'validation.max_photos')
    .optional(),
  categoryId: z.string().min(1, 'validation.required'),
  variants: z
    .array(createVariantSchema)
    .min(1, 'validation.at_least_one_variant')
    .max(50, 'validation.max_variants'),
})

export type MultiVariantProductFormData = z.infer<
  typeof multiVariantProductSchema
>

/**
 * Update Product Schema (for editing product info only)
 */
export const updateProductSchema = z.object({
  name: z.string().min(1, 'validation.required').optional(),
  description: z.string().optional(),
  photos: z
    .array(assetReferenceSchema)
    .max(10, 'validation.max_photos')
    .optional(),
  categoryId: z.string().min(1, 'validation.required').optional(),
})

export type UpdateProductFormData = z.infer<typeof updateProductSchema>

/**
 * Update Variant Schema
 */
export const updateVariantSchema = z.object({
  code: z.string().min(1, 'validation.required').optional(),
  sku: z.string().optional(),
  photos: z
    .array(assetReferenceSchema)
    .max(10, 'validation.max_photos')
    .optional(),
  costPrice: z
    .string()
    .refine(
      (val) => {
        const num = parseFloat(val)
        return !isNaN(num) && num >= 0
      },
      { message: 'validation.positive_number' },
    )
    .optional(),
  salePrice: z
    .string()
    .refine(
      (val) => {
        const num = parseFloat(val)
        return !isNaN(num) && num >= 0
      },
      { message: 'validation.positive_number' },
    )
    .optional(),
  stockQuantity: z.coerce
    .number({ message: 'validation.required' })
    .int({ message: 'validation.integer' })
    .min(0, { message: 'validation.min_zero' })
    .optional(),
  stockQuantityAlert: z.coerce
    .number({ message: 'validation.required' })
    .int({ message: 'validation.integer' })
    .min(0, { message: 'validation.min_zero' })
    .optional(),
})

export type UpdateVariantFormData = z.infer<typeof updateVariantSchema>
