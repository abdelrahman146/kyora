import { useCallback, useRef, useState } from 'react'
import type {
  AssetReference,
  UploadProgressEvent,
  UploadState,
} from '@/types/asset'
import { uploadFilesWithConcurrency } from '@/api/assets'

export interface UseFileUploadOptions {
  businessDescriptor: string
  onSuccess?: (references: Array<AssetReference>) => void
  onError?: (error: Error) => void
  onProgress?: (states: Map<File, UploadState>) => void
}

export interface UseFileUploadReturn {
  /** Current upload states for all files */
  uploadStates: Map<File, UploadState>
  /** Start uploading files */
  upload: (files: Array<File>) => Promise<Array<AssetReference>>
  /** Cancel upload for a specific file */
  cancelUpload: (file: File) => void
  /** Cancel all uploads */
  cancelAll: () => void
  /** Check if any uploads are in progress */
  isUploading: boolean
  /** Get overall upload progress (0-100) */
  overallProgress: number
}

/**
 * Hook for managing file uploads with progress tracking
 */
export function useFileUpload(
  options: UseFileUploadOptions,
): UseFileUploadReturn {
  const { businessDescriptor, onSuccess, onError, onProgress } = options

  const [uploadStates, setUploadStates] = useState<Map<File, UploadState>>(
    new Map(),
  )
  const abortControllersRef = useRef<Map<File, AbortController>>(new Map())

  // Update upload state for a file
  const updateState = useCallback(
    (file: File, update: Partial<UploadState>) => {
      setUploadStates((prev) => {
        const newStates = new Map(prev)
        const current = newStates.get(file) || {
          file,
          status: 'pending',
          progress: 0,
        }
        newStates.set(file, { ...current, ...update })
        onProgress?.(newStates)
        return newStates
      })
    },
    [onProgress],
  )

  // Upload files
  const upload = useCallback(
    async (files: Array<File>): Promise<Array<AssetReference>> => {
      // Initialize states
      files.forEach((file) => {
        updateState(file, {
          file,
          status: 'pending',
          progress: 0,
          error: undefined,
        })
      })

      try {
        const result = await uploadFilesWithConcurrency(
          businessDescriptor,
          files,
          (file, event: UploadProgressEvent) => {
            // Map progress event to upload state
            const statusMap: Record<
              UploadProgressEvent['stage'],
              UploadState['status']
            > = {
              main: 'uploading',
              'thumbnail-generation': 'generating-thumbnail',
              'thumbnail-upload': 'uploading-thumbnail',
              completing: 'completing',
            }

            updateState(file, {
              status: statusMap[event.stage],
              progress: event.percent,
            })
          },
        )

        // Update success states
        result.success.forEach((ref, index) => {
          const file = files.find(
            (f) => f.name === ref.metadata?.altText || files[index],
          )
          if (file) {
            updateState(file, {
              status: 'success',
              progress: 100,
              assetReference: ref,
            })
          }
        })

        // Update error states
        result.failed.forEach(({ file, error }) => {
          updateState(file, {
            status: 'error',
            error,
          })
        })

        onSuccess?.(result.success)

        if (result.failed.length > 0) {
          throw new Error(`${result.failed.length} file(s) failed to upload`)
        }

        return result.success
      } catch (error) {
        onError?.(error as Error)
        throw error
      }
    },
    [businessDescriptor, updateState, onSuccess, onError],
  )

  // Cancel upload for a specific file
  const cancelUpload = useCallback((file: File) => {
    const controller = abortControllersRef.current.get(file)
    if (controller) {
      controller.abort()
      abortControllersRef.current.delete(file)
    }

    setUploadStates((prev) => {
      const newStates = new Map(prev)
      newStates.delete(file)
      return newStates
    })
  }, [])

  // Cancel all uploads
  const cancelAll = useCallback(() => {
    abortControllersRef.current.forEach((controller) => controller.abort())
    abortControllersRef.current.clear()
    setUploadStates(new Map())
  }, [])

  // Calculate overall progress
  const isUploading = Array.from(uploadStates.values()).some(
    (state) =>
      state.status !== 'success' &&
      state.status !== 'error' &&
      state.status !== 'pending',
  )

  const overallProgress =
    Array.from(uploadStates.values()).reduce((sum, state) => {
      return sum + (state.progress || 0)
    }, 0) / Math.max(uploadStates.size, 1)

  return {
    uploadStates,
    upload,
    cancelUpload,
    cancelAll,
    isUploading,
    overallProgress,
  }
}
