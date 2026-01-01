import type { FileValidationResult } from '@/types/asset'
import { FILE_CATEGORY_LIMITS } from '@/types/asset'

/**
 * Parse human-readable file size (e.g., "10MB", "5GB") to bytes
 */
export function parseSize(size: string | number): number {
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
 * Format bytes to human-readable size
 */
export function formatSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  if (bytes < 1024 * 1024 * 1024)
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
  return `${(bytes / (1024 * 1024 * 1024)).toFixed(1)} GB`
}

/**
 * Check if a MIME type matches an accept pattern
 * Supports wildcards like "image/*"
 */
export function matchesMimeType(mimeType: string, pattern: string): boolean {
  if (pattern === mimeType) return true
  if (pattern.endsWith('/*')) {
    const prefix = pattern.slice(0, -2)
    return mimeType.startsWith(prefix)
  }
  return false
}

/**
 * Validate file against accept patterns
 */
function validateFileType(
  file: File,
  accept?: Array<string>,
): FileValidationResult {
  if (!accept || accept.length === 0) {
    return { valid: true }
  }

  const matches = accept.some((pattern) => {
    // Check MIME type
    if (matchesMimeType(file.type, pattern)) return true

    // Check file extension
    if (pattern.startsWith('.')) {
      return file.name.toLowerCase().endsWith(pattern.toLowerCase())
    }

    return false
  })

  if (!matches) {
    return {
      valid: false,
      error: 'invalid-type',
      errorMessage: `File type ${file.type} is not allowed. Accepted types: ${accept.join(', ')}`,
    }
  }

  return { valid: true }
}

/**
 * Validate file size
 */
function validateFileSize(
  file: File,
  maxSize?: number,
  minSize?: number,
): FileValidationResult {
  if (maxSize !== undefined && file.size > maxSize) {
    return {
      valid: false,
      error: 'too-large',
      errorMessage: `File size ${formatSize(file.size)} exceeds maximum ${formatSize(maxSize)}`,
    }
  }

  if (minSize !== undefined && file.size < minSize) {
    return {
      valid: false,
      error: 'too-small',
      errorMessage: `File size ${formatSize(file.size)} is below minimum ${formatSize(minSize)}`,
    }
  }

  return { valid: true }
}

/**
 * Check for duplicate filenames in a batch
 */
export function hasDuplicateFilenames(files: Array<File>): boolean {
  const names = new Set<string>()
  for (const file of files) {
    if (names.has(file.name)) return true
    names.add(file.name)
  }
  return false
}

/**
 * Get file category based on MIME type
 */
export function getFileCategory(
  mimeType: string,
): keyof typeof FILE_CATEGORY_LIMITS | null {
  for (const [category, config] of Object.entries(FILE_CATEGORY_LIMITS)) {
    if ((config.types as unknown as Array<string>).includes(mimeType)) {
      return category as keyof typeof FILE_CATEGORY_LIMITS
    }
  }
  return null
}

/**
 * Get maximum allowed size for a file based on its MIME type
 */
export function getMaxSizeForFile(file: File): number {
  const category = getFileCategory(file.type)
  if (!category) {
    // Default to 10MB for unknown types
    return 10 * 1024 * 1024
  }
  return FILE_CATEGORY_LIMITS[category].maxSize
}

/**
 * Check if a file type supports thumbnail generation
 */
export function supportsThumbnail(mimeType: string): boolean {
  const category = getFileCategory(mimeType)
  return category ? FILE_CATEGORY_LIMITS[category].supportsThumbnail : false
}

/**
 * Validate a single file
 */
export interface ValidateFileOptions {
  accept?: Array<string>
  maxSize?: number | string
  minSize?: number | string
}

export function validateFile(
  file: File,
  options: ValidateFileOptions = {},
): FileValidationResult {
  const { accept, maxSize, minSize } = options

  // Parse size strings to numbers
  const maxSizeBytes = maxSize !== undefined ? parseSize(maxSize) : undefined
  const minSizeBytes = minSize !== undefined ? parseSize(minSize) : undefined

  // Validate type
  const typeResult = validateFileType(file, accept)
  if (!typeResult.valid) return typeResult

  // Validate size
  const sizeResult = validateFileSize(file, maxSizeBytes, minSizeBytes)
  if (!sizeResult.valid) return sizeResult

  return { valid: true }
}

/**
 * Validate multiple files in a batch
 */
export interface ValidateFilesOptions extends ValidateFileOptions {
  maxFiles?: number
  checkDuplicates?: boolean
}

export function validateFiles(
  files: Array<File>,
  options: ValidateFilesOptions = {},
): Map<File, FileValidationResult> {
  const results = new Map<File, FileValidationResult>()
  const { maxFiles, checkDuplicates = true } = options

  // Check max files limit
  if (maxFiles !== undefined && files.length > maxFiles) {
    const error: FileValidationResult = {
      valid: false,
      error: 'unknown',
      errorMessage: `Maximum ${maxFiles} files allowed, but ${files.length} were provided`,
    }
    files.forEach((file) => results.set(file, error))
    return results
  }

  // Check for duplicates
  if (checkDuplicates && hasDuplicateFilenames(files)) {
    const names = new Set<string>()
    files.forEach((file) => {
      if (names.has(file.name)) {
        results.set(file, {
          valid: false,
          error: 'duplicate',
          errorMessage: `Duplicate filename: ${file.name}`,
        })
      } else {
        names.add(file.name)
      }
    })
  }

  // Validate each file individually
  files.forEach((file) => {
    if (!results.has(file)) {
      results.set(file, validateFile(file, options))
    }
  })

  return results
}
