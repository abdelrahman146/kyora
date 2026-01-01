/**
 * Asset management types matching backend AssetReference contract
 */

/**
 * AssetReference - JSONB stored on domain entities (Business.Logo, Product.Photos, etc.)
 * Source of truth: backend/internal/platform/types/asset/asset.go
 */
export interface AssetReference {
  /** Required: CDN or public URL for the asset */
  url: string
  /** Required for garbage collection and future migrations */
  assetId: string
  /** Optional: Fallback S3 URL */
  originalUrl?: string
  /** Optional: Thumbnail CDN/public URL (for images/videos) */
  thumbnailUrl?: string
  /** Optional: Fallback S3 URL for thumbnail */
  thumbnailOriginalUrl?: string
  /** Optional: Additional metadata */
  metadata?: AssetMetadata
}

/**
 * Asset metadata stored within AssetReference
 */
export interface AssetMetadata {
  /** Alt text for accessibility */
  altText?: string
  /** Caption or description */
  caption?: string
  /** Image/video width in pixels */
  width?: number
  /** Image/video height in pixels */
  height?: number
}

/**
 * Upload descriptor returned from backend when generating upload URLs
 */
export interface UploadDescriptor {
  /** Asset ID (ast_...) */
  assetId: string
  /** HTTP method: "PUT" for S3 multipart, "POST" for local */
  method: 'PUT' | 'POST'
  /** Public/CDN URL for the uploaded asset */
  publicUrl: string
  /** CDN URL (same as publicUrl in most cases) */
  cdnUrl: string

  // S3 Multipart fields (when method === 'PUT')
  /** Size of each part in bytes */
  partSize?: number
  /** Total number of parts */
  totalParts?: number
  /** Pre-signed URLs for each part (1-based indexing) */
  partUrls?: Array<string>
  /** S3 upload ID */
  uploadId?: string

  // Local provider fields (when method === 'POST')
  /** Single upload URL for local provider */
  url?: string
  /** Headers to include in local upload request */
  headers?: Record<string, string>

  // Optional thumbnail descriptor
  /** Nested upload descriptor for thumbnail (if backend decides it's needed) */
  thumbnail?: UploadDescriptor
}

/**
 * Response from generateUploadURLs API
 */
export interface GenerateUploadURLsResponse {
  uploads: Array<UploadDescriptor>
}

/**
 * Request to generateUploadURLs API
 */
export interface GenerateUploadURLsRequest {
  files: Array<{
    fileName: string
    contentType: string
    sizeBytes: number
  }>
}

/**
 * Request to complete multipart upload
 */
export interface CompleteMultipartUploadRequest {
  parts: Array<{
    /** Part number (1-based) */
    partNumber: number
    /** ETag from part upload response (without quotes) */
    etag: string
  }>
}

/**
 * Upload state machine status
 */
export type UploadStatus =
  | 'pending' // Queued but not started
  | 'uploading' // Uploading main file
  | 'generating-thumbnail' // Generating thumbnail in worker/ffmpeg
  | 'uploading-thumbnail' // Uploading thumbnail bytes
  | 'completing' // Completing multipart upload
  | 'success' // Upload completed successfully
  | 'error' // Upload failed

/**
 * Upload state for a single file
 */
export interface UploadState {
  /** The file being uploaded */
  file: File
  /** Current status */
  status: UploadStatus
  /** Upload progress (0-100) */
  progress: number
  /** Error message if status === 'error' */
  error?: string
  /** Resulting AssetReference (available after success or optimistically) */
  assetReference?: AssetReference
  /** Descriptor from backend */
  descriptor?: UploadDescriptor
  /** Cancel controller for aborting upload */
  abortController?: AbortController
}

/**
 * Progress event emitted during upload
 */
export interface UploadProgressEvent {
  /** Current stage of upload */
  stage: 'main' | 'thumbnail-generation' | 'thumbnail-upload' | 'completing'
  /** Progress percentage (0-100) */
  percent: number
  /** Optional status message */
  message?: string
}

/**
 * Result of batch upload operation
 */
export interface UploadResult {
  /** Successfully uploaded files */
  success: Array<AssetReference>
  /** Failed uploads with error details */
  failed: Array<{
    file: File
    error: string
  }>
}

/**
 * File validation error types
 */
export type FileValidationError =
  | 'invalid-type' // MIME type not in accept list
  | 'too-large' // File exceeds maxSize
  | 'too-small' // File below minSize
  | 'duplicate' // Duplicate filename in batch
  | 'dimensions' // Image dimensions out of range
  | 'unknown' // Unknown validation error

/**
 * File validation result
 */
export interface FileValidationResult {
  valid: boolean
  error?: FileValidationError
  errorMessage?: string
}

/**
 * Backend file category limits (matching asset_upload.instructions.md)
 */
export const FILE_CATEGORY_LIMITS = {
  images: {
    types: [
      'image/jpeg',
      'image/jpg',
      'image/png',
      'image/webp',
      'image/heic',
      'image/heif',
      'image/gif',
    ],
    maxSize: 10 * 1024 * 1024, // 10 MB
    supportsThumbnail: true,
  },
  videos: {
    types: [
      'video/mp4',
      'video/quicktime',
      'video/x-msvideo',
      'video/x-matroska',
      'video/webm',
    ],
    maxSize: 100 * 1024 * 1024, // 100 MB
    supportsThumbnail: true,
  },
  audio: {
    types: ['audio/mpeg', 'audio/wav', 'audio/ogg', 'audio/mp4', 'audio/aac'],
    maxSize: 20 * 1024 * 1024, // 20 MB
    supportsThumbnail: false,
  },
  documents: {
    types: [
      'application/pdf',
      'application/msword',
      'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
      'text/plain',
      'application/rtf',
      'application/vnd.oasis.opendocument.text',
    ],
    maxSize: 10 * 1024 * 1024, // 10 MB
    supportsThumbnail: false,
  },
  compressed: {
    types: [
      'application/zip',
      'application/x-tar',
      'application/gzip',
      'application/x-rar-compressed',
      'application/x-7z-compressed',
    ],
    maxSize: 50 * 1024 * 1024, // 50 MB
    supportsThumbnail: false,
  },
} as const
