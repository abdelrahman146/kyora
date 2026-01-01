import type { AssetReference } from '@/api/types/asset'

/**
 * Get the thumbnail URL from an AssetReference with intelligent fallback logic
 *
 * Fallback chain:
 * 1. thumbnailUrl (preferred thumbnail)
 * 2. thumbnailOriginalUrl (original thumbnail)
 * 3. url (main photo as fallback)
 * 4. originalUrl (original main photo)
 * 5. undefined (no valid URL found)
 *
 * @param asset - The asset reference object (can be null/undefined)
 * @returns The best available thumbnail URL or undefined if no valid URL found
 *
 * @example
 * ```tsx
 * <Avatar src={getThumbnailUrl(business.logo)} />
 * ```
 */
export function getThumbnailUrl(
  asset: AssetReference | null | undefined,
): string | undefined {
  if (!asset) {
    return undefined
  }

  return (
    asset.thumbnailUrl ||
    asset.thumbnailOriginalUrl ||
    asset.url ||
    asset.originalUrl ||
    undefined
  )
}

/**
 * Get the full photo URL from an AssetReference with intelligent fallback logic
 *
 * Fallback chain:
 * 1. url (main photo)
 * 2. originalUrl (original main photo)
 * 3. thumbnailUrl (thumbnail as fallback)
 * 4. thumbnailOriginalUrl (original thumbnail)
 * 5. undefined (no valid URL found)
 *
 * @param asset - The asset reference object (can be null/undefined)
 * @returns The best available photo URL or undefined if no valid URL found
 *
 * @example
 * ```tsx
 * <img src={getPhotoUrl(product.image)} alt="Product" />
 * ```
 */
export function getPhotoUrl(
  asset: AssetReference | null | undefined,
): string | undefined {
  if (!asset) {
    return undefined
  }

  return (
    asset.url ||
    asset.originalUrl ||
    asset.thumbnailUrl ||
    asset.thumbnailOriginalUrl ||
    undefined
  )
}
