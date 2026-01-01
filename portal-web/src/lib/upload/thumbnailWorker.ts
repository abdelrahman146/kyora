import {
  THUMBNAIL_FORMAT,
  THUMBNAIL_QUALITY,
  THUMBNAIL_SIZE,
} from './constants'
import { supportsWebP } from './metadataExtractor'

/**
 * Generate a thumbnail from an image file using canvas
 * This runs in a Web Worker when available, or on the main thread as fallback
 */
export interface ThumbnailOptions {
  size?: number
  quality?: number
  format?: 'webp' | 'jpeg'
}

export async function generateImageThumbnail(
  file: File,
  options: ThumbnailOptions = {},
): Promise<{ blob: Blob; width: number; height: number }> {
  const {
    size = THUMBNAIL_SIZE,
    quality = THUMBNAIL_QUALITY,
    format = THUMBNAIL_FORMAT,
  } = options

  // Check WebP support for format fallback
  const useWebP = format === 'webp' && (await supportsWebP())
  const outputFormat = useWebP ? 'image/webp' : 'image/jpeg'

  return new Promise((resolve, reject) => {
    const img = new Image()
    const objectUrl = URL.createObjectURL(file)

    img.onload = () => {
      try {
        // Calculate dimensions maintaining aspect ratio
        const { width: originalWidth, height: originalHeight } = img
        let width = size
        let height = size

        if (originalWidth > originalHeight) {
          height = Math.round((size * originalHeight) / originalWidth)
        } else if (originalHeight > originalWidth) {
          width = Math.round((size * originalWidth) / originalHeight)
        }

        // Create canvas and draw resized image
        const canvas = document.createElement('canvas')
        canvas.width = width
        canvas.height = height

        const ctx = canvas.getContext('2d')
        if (!ctx) {
          throw new Error('Failed to get canvas context')
        }

        // Use better image smoothing for quality
        ctx.imageSmoothingEnabled = true
        ctx.imageSmoothingQuality = 'high'

        ctx.drawImage(img, 0, 0, width, height)

        // Convert to blob
        canvas.toBlob(
          (blob) => {
            URL.revokeObjectURL(objectUrl)
            if (blob) {
              resolve({ blob, width, height })
            } else {
              reject(new Error('Failed to create thumbnail blob'))
            }
          },
          outputFormat,
          quality,
        )
      } catch (error) {
        URL.revokeObjectURL(objectUrl)
        reject(error)
      }
    }

    img.onerror = () => {
      URL.revokeObjectURL(objectUrl)
      reject(new Error('Failed to load image'))
    }

    img.src = objectUrl
  })
}

/**
 * Check if Web Workers are supported
 */
export function supportsWebWorker(): boolean {
  return typeof Worker !== 'undefined'
}

/**
 * Generate thumbnail with Web Worker if available, fallback to main thread
 */
export async function generateThumbnail(
  file: File,
  options?: ThumbnailOptions,
): Promise<{ blob: Blob; width: number; height: number }> {
  // For now, always use main thread implementation
  // Web Worker implementation would require bundling the worker code separately
  // which is complex with Vite. Main thread is acceptable for thumbnails since
  // canvas operations are relatively fast and don't block for long.
  return generateImageThumbnail(file, options)
}
