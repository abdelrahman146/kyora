import { apiClient } from './client'
import type {
  AssetReference,
  CompleteMultipartUploadRequest,
  GenerateUploadURLsRequest,
  GenerateUploadURLsResponse,
  UploadDescriptor,
  UploadProgressEvent,
  UploadResult,
} from '@/types/asset'
import { generateThumbnail } from '@/lib/upload/thumbnailWorker'
import { generateVideoThumbnail } from '@/lib/upload/videoThumbnail'
import { extractMetadata } from '@/lib/upload/metadataExtractor'
import {
  MAX_CONCURRENT_UPLOADS,
  MAX_RETRIES,
  RETRY_DELAYS,
} from '@/lib/upload/constants'

/**
 * Generate upload URLs from backend
 */
export async function generateUploadURLs(
  businessDescriptor: string,
  files: Array<{ fileName: string; contentType: string; sizeBytes: number }>,
): Promise<GenerateUploadURLsResponse> {
  const response = await apiClient
    .post(`v1/businesses/${businessDescriptor}/assets/uploads`, {
      json: { files } satisfies GenerateUploadURLsRequest,
    })
    .json<GenerateUploadURLsResponse>()

  return response
}

/**
 * Complete multipart upload
 */
export async function completeMultipartUpload(
  businessDescriptor: string,
  assetId: string,
  parts: Array<{ partNumber: number; etag: string }>,
): Promise<void> {
  // Strip quotes from ETags if present
  const sanitizedParts = parts.map((part) => ({
    ...part,
    etag: part.etag.replace(/"/g, ''),
  }))

  await apiClient.post(
    `v1/businesses/${businessDescriptor}/assets/uploads/${assetId}/complete`,
    {
      json: { parts: sanitizedParts } satisfies CompleteMultipartUploadRequest,
    },
  )
}

/**
 * Upload file bytes to storage (S3 multipart or local POST)
 */
export async function uploadToStorage(
  descriptor: UploadDescriptor,
  file: File,
  onProgress?: (percent: number) => void,
  signal?: AbortSignal,
): Promise<Array<{ partNumber: number; etag: string }> | null> {
  if (
    descriptor.method === 'PUT' &&
    descriptor.partUrls &&
    descriptor.partSize
  ) {
    // S3 Multipart upload
    return uploadMultipart(descriptor, file, onProgress, signal)
  } else if (descriptor.method === 'POST' && descriptor.url) {
    // Local provider upload
    await uploadLocal(descriptor, file, onProgress, signal)
    return null
  } else {
    throw new Error('Invalid upload descriptor')
  }
}

/**
 * Upload file using S3 multipart
 */
async function uploadMultipart(
  descriptor: UploadDescriptor,
  file: File,
  onProgress?: (percent: number) => void,
  signal?: AbortSignal,
): Promise<Array<{ partNumber: number; etag: string }>> {
  const { partUrls, partSize } = descriptor
  if (!partUrls || !partSize) {
    throw new Error('Missing partUrls or partSize for multipart upload')
  }

  const parts: Array<{ partNumber: number; etag: string }> = []
  const totalParts = partUrls.length

  for (let i = 0; i < totalParts; i++) {
    const partNumber = i + 1 // 1-based
    const start = i * partSize
    const end = Math.min(start + partSize, file.size)
    const chunk = file.slice(start, end)

    // Upload part with retry
    const etag = await uploadPartWithRetry(partUrls[i], chunk, signal)
    parts.push({ partNumber, etag })

    // Report progress
    if (onProgress) {
      onProgress(Math.round(((i + 1) / totalParts) * 100))
    }

    // Check if aborted
    if (signal?.aborted) {
      throw new Error('Upload cancelled')
    }
  }

  return parts
}

/**
 * Upload a single part with exponential backoff retry
 */
async function uploadPartWithRetry(
  url: string,
  chunk: Blob,
  signal?: AbortSignal,
  retryCount = 0,
): Promise<string> {
  try {
    const response = await fetch(url, {
      method: 'PUT',
      body: chunk,
      signal,
    })

    if (!response.ok) {
      throw new Error(`Part upload failed: ${response.status}`)
    }

    const etag = response.headers.get('ETag')
    if (!etag) {
      throw new Error('Missing ETag in response. Check CORS configuration.')
    }

    return etag
  } catch (error) {
    // Retry on network errors
    if (retryCount < MAX_RETRIES && !signal?.aborted) {
      await new Promise((resolve) =>
        setTimeout(resolve, RETRY_DELAYS[retryCount] || 4000),
      )
      return uploadPartWithRetry(url, chunk, signal, retryCount + 1)
    }
    throw error
  }
}

/**
 * Upload file using local provider
 */
async function uploadLocal(
  descriptor: UploadDescriptor,
  file: File,
  onProgress?: (percent: number) => void,
  signal?: AbortSignal,
): Promise<void> {
  if (!descriptor.url) {
    throw new Error('Missing URL for local upload')
  }

  const headers = descriptor.headers || {}

  const response = await fetch(descriptor.url, {
    method: 'POST',
    body: file,
    headers,
    signal,
  })

  if (!response.ok) {
    throw new Error(`Local upload failed: ${response.status}`)
  }

  if (onProgress) {
    onProgress(100)
  }
}

/**
 * Upload a single file with thumbnail (if needed)
 */
export async function uploadFileWithThumbnail(
  businessDescriptor: string,
  file: File,
  onProgress?: (event: UploadProgressEvent) => void,
  signal?: AbortSignal,
): Promise<AssetReference> {
  // Step 1: Generate upload URLs
  const { uploads } = await generateUploadURLs(businessDescriptor, [
    {
      fileName: file.name,
      contentType: file.type,
      sizeBytes: file.size,
    },
  ])

  const descriptor = uploads[0]
  // Step 2: Upload main file
  onProgress?.({ stage: 'main', percent: 0 })
  const parts = await uploadToStorage(
    descriptor,
    file,
    (percent) => onProgress?.({ stage: 'main', percent }),
    signal,
  )

  // Step 3: Complete multipart for main file if needed
  if (parts && descriptor.method === 'PUT') {
    onProgress?.({ stage: 'completing', percent: 0 })
    await completeMultipartUpload(businessDescriptor, descriptor.assetId, parts)
  }

  // Step 4: Handle thumbnail if descriptor includes it
  let thumbnailUrl: string | undefined
  if (descriptor.thumbnail) {
    onProgress?.({ stage: 'thumbnail-generation', percent: 0 })

    let thumbnailBlob: Blob | null = null

    // Generate thumbnail based on file type
    if (file.type.startsWith('image/')) {
      const result = await generateThumbnail(file)
      thumbnailBlob = result.blob
    } else if (file.type.startsWith('video/')) {
      thumbnailBlob = await generateVideoThumbnail(file)
    }

    if (thumbnailBlob) {
      // Upload thumbnail
      onProgress?.({ stage: 'thumbnail-upload', percent: 0 })

      const thumbnailFile = new File([thumbnailBlob], `thumb_${file.name}`, {
        type: thumbnailBlob.type,
      })

      const thumbnailParts = await uploadToStorage(
        descriptor.thumbnail,
        thumbnailFile,
        (percent) => onProgress?.({ stage: 'thumbnail-upload', percent }),
        signal,
      )

      // Complete multipart for thumbnail if needed
      if (thumbnailParts && descriptor.thumbnail.method === 'PUT') {
        await completeMultipartUpload(
          businessDescriptor,
          descriptor.thumbnail.assetId,
          thumbnailParts,
        )
      }

      thumbnailUrl =
        descriptor.thumbnail.cdnUrl || descriptor.thumbnail.publicUrl
    }
  }

  // Step 5: Extract metadata
  const metadata = await extractMetadata(file)

  // Step 6: Build AssetReference
  const assetReference: AssetReference = {
    url: descriptor.cdnUrl || descriptor.publicUrl,
    assetId: descriptor.assetId,
    originalUrl:
      descriptor.publicUrl !== descriptor.cdnUrl
        ? descriptor.publicUrl
        : undefined,
    thumbnailUrl,
    metadata,
  }

  return assetReference
}

/**
 * Upload queue manager for concurrent uploads
 */
export class UploadQueue {
  private queue: Array<{
    file: File
    resolve: (ref: AssetReference) => void
    reject: (error: Error) => void
  }> = []
  private active = 0
  private businessDescriptor: string
  private onProgress?: (file: File, event: UploadProgressEvent) => void
  private abortControllers = new Map<File, AbortController>()

  constructor(
    businessDescriptor: string,
    onProgress?: (file: File, event: UploadProgressEvent) => void,
  ) {
    this.businessDescriptor = businessDescriptor
    this.onProgress = onProgress
  }

  /**
   * Add file to upload queue
   */
  enqueue(file: File): Promise<AssetReference> {
    return new Promise((resolve, reject) => {
      this.queue.push({ file, resolve, reject })
      this.processQueue()
    })
  }

  /**
   * Cancel upload for a specific file
   */
  cancel(file: File): void {
    const controller = this.abortControllers.get(file)
    if (controller) {
      controller.abort()
      this.abortControllers.delete(file)
    }

    // Remove from queue if not started
    const index = this.queue.findIndex((item) => item.file === file)
    if (index !== -1) {
      const item = this.queue.splice(index, 1)[0]
      item.reject(new Error('Upload cancelled'))
    }
  }

  /**
   * Process upload queue with concurrency limit
   */
  private processQueue(): void {
    while (this.active < MAX_CONCURRENT_UPLOADS && this.queue.length > 0) {
      const item = this.queue.shift()
      if (!item) break

      this.active++

      const controller = new AbortController()
      this.abortControllers.set(item.file, controller)

      uploadFileWithThumbnail(
        this.businessDescriptor,
        item.file,
        (event) => this.onProgress?.(item.file, event),
        controller.signal,
      )
        .then((ref) => {
          item.resolve(ref)
        })
        .catch((error) => {
          item.reject(error)
        })
        .finally(() => {
          this.active--
          this.abortControllers.delete(item.file)
          this.processQueue()
        })
    }
  }
}

/**
 * Upload multiple files with concurrency control
 */
export async function uploadFilesWithConcurrency(
  businessDescriptor: string,
  files: Array<File>,
  onProgress?: (file: File, event: UploadProgressEvent) => void,
): Promise<UploadResult> {
  const queue = new UploadQueue(businessDescriptor, onProgress)

  const results = await Promise.allSettled(
    files.map((file) => queue.enqueue(file)),
  )

  const success: Array<AssetReference> = []
  const failed: Array<{ file: File; error: string }> = []

  results.forEach((result, index) => {
    if (result.status === 'fulfilled') {
      success.push(result.value)
    } else {
      failed.push({
        file: files[index],
        error: result.reason?.message || 'Upload failed',
      })
    }
  })

  return { success, failed }
}
