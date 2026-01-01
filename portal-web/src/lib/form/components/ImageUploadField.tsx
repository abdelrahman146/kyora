/**
 * ImageUploadField Component - Specialized Image Upload
 *
 * Specialized variant of FileUploadField for image uploads only.
 * Forces image MIME types, shows grid layout with thumbnails, and supports
 * single/multiple modes with drag-and-drop reordering.
 *
 * Features:
 * - Image-only validation (JPEG, PNG, WebP, HEIC)
 * - Grid layout (2 cols mobile, 3 tablet, 4 desktop)
 * - Primary image indicator (first image)
 * - Drag-and-drop reordering with @dnd-kit
 * - Automatic thumbnail generation
 * - Mobile camera support
 *
 * Usage:
 * ```tsx
 * // Single image (business logo)
 * <form.Field name="logo">
 *   {(field) => (
 *     <field.ImageUploadField
 *       single
 *       label="Business Logo"
 *       maxSize="10MB"
 *     />
 *   )}
 * </form.Field>
 *
 * // Multiple images (product photos)
 * <form.Field name="photos">
 *   {(field) => (
 *     <field.ImageUploadField
 *       label="Product Photos"
 *       maxFiles={10}
 *       reorderable
 *     />
 *   )}
 * </form.Field>
 * ```
 */

import { FileUploadField } from './FileUploadField'
import type { FileUploadFieldProps } from './FileUploadField'

export interface ImageUploadFieldProps extends Omit<
  FileUploadFieldProps,
  'accept' | 'multiple'
> {
  /** Single image mode (default: multiple) */
  single?: boolean
}

export function ImageUploadField(props: ImageUploadFieldProps) {
  const { single = false, maxFiles = single ? 1 : 10, ...rest } = props

  return (
    <FileUploadField
      accept="image/jpeg,image/jpg,image/png,image/webp,image/heic,image/heif"
      multiple={!single}
      maxFiles={maxFiles}
      {...rest}
    />
  )
}
