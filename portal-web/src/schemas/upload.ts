import { z } from 'zod'

/**
 * Zod schema factories for file upload validation
 */

/**
 * Validate File instance
 */
export interface FileSchemaOptions {
  accept?: Array<string>
  maxSize?: number | string
  minSize?: number | string
  required?: boolean
}

export function fileSchema(options: FileSchemaOptions = {}) {
  let schema = z.instanceof(File, { message: 'invalid_file' })

  if (options.required === false) {
    schema = schema.optional() as any
  }

  return schema.refine(
    (file) => {
      // Validate file type
      if (options.accept && options.accept.length > 0) {
        const matches = options.accept.some((pattern) => {
          if (pattern.startsWith('.')) {
            return file.name.toLowerCase().endsWith(pattern.toLowerCase())
          }
          if (pattern.endsWith('/*')) {
            return file.type.startsWith(pattern.slice(0, -2))
          }
          return file.type === pattern
        })
        if (!matches) return false
      }

      // Validate file size
      const maxSize =
        typeof options.maxSize === 'string'
          ? parseSize(options.maxSize)
          : options.maxSize
      if (maxSize && file.size > maxSize) return false

      const minSize =
        typeof options.minSize === 'string'
          ? parseSize(options.minSize)
          : options.minSize
      if (minSize && file.size < minSize) return false

      return true
    },
    {
      message: 'invalid_file',
    },
  )
}

/**
 * Validate image file with dimension constraints
 */
export interface ImageSchemaOptions extends FileSchemaOptions {
  minWidth?: number
  maxWidth?: number
  minHeight?: number
  maxHeight?: number
}

export function imageSchema(options: ImageSchemaOptions = {}) {
  const baseOptions: FileSchemaOptions = {
    ...options,
    accept: [
      'image/jpeg',
      'image/jpg',
      'image/png',
      'image/webp',
      'image/heic',
      'image/heif',
    ],
  }

  return fileSchema(baseOptions)
}

/**
 * Validate AssetReference
 */
export function assetReferenceSchema() {
  return z.object({
    url: z.string().url(),
    assetId: z.string().startsWith('ast_'),
    originalUrl: z.string().url().optional(),
    thumbnailUrl: z.string().url().optional(),
    thumbnailOriginalUrl: z.string().url().optional(),
    metadata: z
      .object({
        altText: z.string().optional(),
        caption: z.string().optional(),
        width: z.number().positive().optional(),
        height: z.number().positive().optional(),
      })
      .optional(),
  })
}

/**
 * Validate array of AssetReferences
 */
export interface AssetReferenceArrayOptions {
  min?: number
  max?: number
  required?: boolean
}

export function assetReferenceArraySchema(
  options: AssetReferenceArrayOptions = {},
) {
  let schema = z.array(assetReferenceSchema())

  if (options.min !== undefined) {
    schema = schema.min(options.min, {
      message: `array_min_items`,
    })
  }

  if (options.max !== undefined) {
    schema = schema.max(options.max, {
      message: `array_max_items`,
    })
  }

  if (options.required === false) {
    return schema.optional()
  }

  return schema
}

/**
 * Validate either File or AssetReference (for update forms)
 */
export function fileOrAssetReference(options: FileSchemaOptions = {}) {
  return z.union([fileSchema(options), assetReferenceSchema()])
}

/**
 * Validate array of Files or AssetReferences
 */
export function fileOrAssetReferenceArray(
  options: FileSchemaOptions & AssetReferenceArrayOptions = {},
) {
  let schema = z.array(z.union([fileSchema(options), assetReferenceSchema()]))

  if (options.min !== undefined) {
    schema = schema.min(options.min, {
      message: `array_min_items`,
    })
  }

  if (options.max !== undefined) {
    schema = schema.max(options.max, {
      message: `array_max_items`,
    })
  }

  if (options.required === false) {
    return schema.optional()
  }

  return schema
}

/**
 * Helper: Parse size string to bytes
 */
function parseSize(size: string | number): number {
  if (typeof size === 'number') return size

  const units: Record<string, number> = {
    B: 1,
    KB: 1024,
    MB: 1024 * 1024,
    GB: 1024 * 1024 * 1024,
  }

  const match = size.toUpperCase().match(/^(\d+(?:\.\d+)?)\s*([KMGT]?B)$/)
  if (!match) {
    throw new Error(`Invalid size format: ${size}`)
  }

  const [, value, unit] = match
  return parseFloat(value) * (units[unit] || 1)
}

/**
 * Example schemas for common use cases
 */

// Business logo (single image, max 10MB)
export const businessLogoSchema = imageSchema({
  maxSize: '10MB',
  required: false,
})

// Product photos (1-10 images)
export const productPhotosSchema = assetReferenceArraySchema({
  min: 1,
  max: 10,
})

// Variant photos (max 10 images)
export const variantPhotosSchema = assetReferenceArraySchema({
  max: 10,
  required: false,
})

// Generic documents (PDF, DOC, etc.)
export const documentSchema = fileSchema({
  accept: [
    'application/pdf',
    '.pdf',
    'application/msword',
    '.doc',
    'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    '.docx',
  ],
  maxSize: '10MB',
})
