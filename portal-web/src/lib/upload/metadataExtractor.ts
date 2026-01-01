import type { AssetMetadata } from '@/types/asset'

/**
 * Extract metadata from an image file
 */
export async function extractImageMetadata(file: File): Promise<AssetMetadata> {
  return new Promise((resolve, reject) => {
    const img = new Image()
    const objectUrl = URL.createObjectURL(file)

    img.onload = () => {
      URL.revokeObjectURL(objectUrl)
      resolve({
        width: img.naturalWidth,
        height: img.naturalHeight,
      })
    }

    img.onerror = () => {
      URL.revokeObjectURL(objectUrl)
      reject(new Error('Failed to load image'))
    }

    img.src = objectUrl
  })
}

/**
 * Extract basic file metadata
 */
export function extractFileMetadata(): AssetMetadata {
  return {
    // Caption and altText should be added by user
    altText: undefined,
    caption: undefined,
  }
}

/**
 * Extract metadata from a file based on its type
 */
export async function extractMetadata(file: File): Promise<AssetMetadata> {
  if (file.type.startsWith('image/')) {
    try {
      return await extractImageMetadata(file)
    } catch (error) {
      console.warn('Failed to extract image metadata:', error)
      return extractFileMetadata()
    }
  }

  // For non-images, return basic metadata
  return extractFileMetadata()
}

/**
 * Check if browser supports WebP
 */
let webpSupported: boolean | null = null

export async function supportsWebP(): Promise<boolean> {
  if (webpSupported !== null) return webpSupported

  return new Promise((resolve) => {
    const img = new Image()
    img.onload = () => {
      webpSupported = img.width > 0 && img.height > 0
      resolve(webpSupported)
    }
    img.onerror = () => {
      webpSupported = false
      resolve(false)
    }
    // Minimal WebP image (1x1 transparent pixel)
    img.src =
      'data:image/webp;base64,UklGRhoAAABXRUJQVlA4TA0AAAAvAAAAEAcQERGIiP4HAA=='
  })
}
