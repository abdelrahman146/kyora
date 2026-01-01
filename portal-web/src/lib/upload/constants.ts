/**
 * Upload system configuration constants
 */

/** Maximum number of concurrent file uploads */
export const MAX_CONCURRENT_UPLOADS = 3

/** Thumbnail dimensions (square) */
export const THUMBNAIL_SIZE = 300

/** Thumbnail JPEG quality (0.0 - 1.0) */
export const THUMBNAIL_QUALITY = 0.8

/** Preferred thumbnail format (WebP for better compression, JPEG fallback) */
export const THUMBNAIL_FORMAT = 'webp'

/** Retry configuration */
export const RETRY_DELAYS = [1000, 2000, 4000] // 1s, 2s, 4s
export const MAX_RETRIES = 3

/** Video thumbnail extraction time (seconds) */
export const VIDEO_THUMBNAIL_TIME = 1.0

/** Default chunk size for multipart uploads (5MB) */
export const DEFAULT_CHUNK_SIZE = 5 * 1024 * 1024

/** Session storage key for upload queue persistence */
export const UPLOAD_QUEUE_STORAGE_KEY = 'kyora_upload_queue'
