/**
 * FileUploadField Component - Form Composition Layer
 *
 * Generic file upload field with optimistic UI, progress tracking, and drag-and-drop reordering.
 * Handles both File[] (creation) and AssetReference[] (update) modes seamlessly.
 *
 * Features:
 * - Optimistic UI (show preview before upload completes)
 * - Real-time upload progress tracking
 * - Drag-and-drop file selection and reordering
 * - Automatic thumbnail generation for images/videos
 * - Mobile camera/gallery support
 * - Retry failed uploads
 * - RTL support
 * - Accessible keyboard navigation
 *
 * Usage:
 * ```tsx
 * // Creation mode (File[] → AssetReference[])
 * <form.Field name="documents">
 *   {(field) => (
 *     <field.FileUploadField
 *       accept=".pdf,.doc"
 *       maxFiles={5}
 *       maxSize="10MB"
 *     />
 *   )}
 * </form.Field>
 *
 * // Update mode (AssetReference[] with add/remove)
 * <form.Field name="documents" defaultValue={existingRefs}>
 *   {(field) => (
 *     <field.FileUploadField accept=".pdf" maxFiles={5} />
 *   )}
 * </form.Field>
 * ```
 */

import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useState,
} from 'react'
import { Camera, Image as ImageIcon } from 'lucide-react'
import { useTranslation } from 'react-i18next'

import { useFieldContext } from '../contexts'
import type { AssetReference } from '@/types/asset'
import {
  FilePreview,
  FileThumbnail,
  FileUploadZone,
  UploadProgress,
} from '@/components/atoms'
import { useFileUpload } from '@/lib/upload'
import { useObjectURLs } from '@/lib/upload/filePreviewManager'
import { validateFiles } from '@/lib/upload/fileValidator'
import { getThumbnailUrl } from '@/lib/assetUrl'
import { useMediaQuery } from '@/hooks'

// Business Context for getting businessDescriptor from route/layout
export const BusinessContext = createContext<{
  businessDescriptor: string
} | null>(null)

export interface FileUploadFieldProps {
  /** Accepted file types (MIME types or extensions) */
  accept?: string
  /** Maximum number of files */
  maxFiles?: number
  /** Maximum file size (e.g., "10MB") */
  maxSize?: string
  /** Allow multiple file selection */
  multiple?: boolean
  /** Enable drag-and-drop reordering */
  reorderable?: boolean
  /** Field label */
  label?: string
  /** Helper text */
  hint?: string
  /** Disabled state */
  disabled?: boolean
  /** Visual required indicator */
  required?: boolean
  /** Callback when upload completes */
  onUploadComplete?: (references: Array<AssetReference>) => void
}

function isAssetReference(value: unknown): value is AssetReference {
  return (
    typeof value === 'object' &&
    value !== null &&
    'url' in value &&
    'assetId' in value &&
    typeof (value as AssetReference).url === 'string'
  )
}

function isFile(value: unknown): value is File {
  return value instanceof File
}

export function FileUploadField(props: FileUploadFieldProps) {
  const {
    accept,
    maxFiles,
    maxSize,
    multiple = false,
    label,
    hint,
    disabled = false,
    onUploadComplete,
  } = props

  const field = useFieldContext<
    File | Array<File> | AssetReference | Array<AssetReference>
  >()
  const { t } = useTranslation(['upload', 'errors'])
  const isMobile = useMediaQuery('(max-width: 768px)')

  // Get business descriptor from context
  const businessContext = useContext(BusinessContext)
  const businessDescriptor = businessContext?.businessDescriptor || 'default'

  // Separate existing references and pending uploads
  const [pendingFiles, setPendingFiles] = useState<Array<File>>([])

  // Parse field value into existing references and files
  const { existingReferences, newFiles } = useMemo(() => {
    const value = field.state.value
    const refs: Array<AssetReference> = []
    const files: Array<File> = []

    if (Array.isArray(value)) {
      value.forEach((item) => {
        if (isAssetReference(item)) {
          refs.push(item)
        } else if (isFile(item)) {
          files.push(item)
        }
      })
    } else if (isAssetReference(value)) {
      refs.push(value)
    } else if (isFile(value)) {
      files.push(value)
    }

    return { existingReferences: refs, newFiles: files }
  }, [field.state.value])

  // Upload hook
  const { uploadStates, upload, cancelUpload, isUploading } = useFileUpload({
    businessDescriptor,
    onSuccess: (references) => {
      // Merge uploaded references into field value
      const currentValue = field.state.value
      let newValue: Array<AssetReference> | AssetReference

      if (Array.isArray(currentValue)) {
        newValue = [...existingReferences, ...references]
      } else {
        newValue = references[0] || existingReferences[0]
      }

      field.handleChange(newValue as any)
      setPendingFiles([])
      onUploadComplete?.(references)
    },
    onError: (error) => {
      console.error('Upload failed:', error)
      // Show error toast with translated message
      const errorMessage = error.message.includes('404')
        ? t('uploadFailed', { ns: 'upload' }) +
          ': ' +
          t('invalidPath', { ns: 'errors' })
        : t('uploadFailed', { ns: 'upload' })

      // Use toast library if available, otherwise console
      if (typeof window !== 'undefined' && 'toast' in window) {
        // @ts-ignore - window.toast is dynamically added by react-hot-toast
        window.toast?.error?.(errorMessage)
      }
    },
  })

  // Object URLs for file previews
  const objectURLs = useObjectURLs([...newFiles, ...pendingFiles])

  // Handle file selection
  const handleFilesSelected = useCallback(
    (files: Array<File>) => {
      if (disabled) return

      // Validate files
      const validationResults = validateFiles(files, {
        accept: accept?.split(','),
        maxSize,
        maxFiles: maxFiles
          ? maxFiles - existingReferences.length - newFiles.length
          : undefined,
        checkDuplicates: true,
      })

      // Filter valid files
      const validFiles = files.filter((file) => {
        const result = validationResults.get(file)
        if (!result?.valid) {
          console.warn(
            `File ${file.name} validation failed:`,
            result?.errorMessage,
          )
          return false
        }
        return true
      })

      if (validFiles.length === 0) return

      // Add to pending and start upload
      setPendingFiles((prev) => [...prev, ...validFiles])

      upload(validFiles).catch((error) => {
        console.error('Upload error:', error)
      })
    },
    [
      disabled,
      accept,
      maxSize,
      maxFiles,
      existingReferences.length,
      newFiles.length,
      upload,
    ],
  )

  // Handle remove
  const handleRemove = useCallback(
    (item: AssetReference | File) => {
      if (isAssetReference(item)) {
        // Remove from field value
        const currentValue = field.state.value
        if (Array.isArray(currentValue)) {
          const newValue = currentValue.filter((v) => {
            if (isAssetReference(v)) {
              return v.assetId !== item.assetId
            }
            return true
          })
          field.handleChange(newValue as any)
        } else if (
          isAssetReference(currentValue) &&
          currentValue.assetId === item.assetId
        ) {
          field.handleChange(undefined as any)
        }
      } else if (isFile(item)) {
        // Cancel upload and remove from pending
        cancelUpload(item)
        setPendingFiles((prev) => prev.filter((f) => f !== item))
      }
    },
    [field, cancelUpload],
  )

  // Drag-and-drop reordering
  const [draggedIndex, setDraggedIndex] = useState<number | null>(null)

  const handleDragStart = useCallback((index: number) => {
    setDraggedIndex(index)
  }, [])

  const handleDragEnd = useCallback(() => {
    setDraggedIndex(null)
  }, [])

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault()
  }, [])

  const handleDrop = useCallback(
    (targetIndex: number) => {
      if (draggedIndex === null || draggedIndex === targetIndex) return

      const currentValue = field.state.value
      if (!Array.isArray(currentValue)) return

      // Only reorder existing references (not pending files)
      const refs = currentValue.filter(isAssetReference)
      if (draggedIndex >= refs.length || targetIndex >= refs.length) return

      // Reorder array
      const newRefs = [...refs]
      const [removed] = newRefs.splice(draggedIndex, 1)
      newRefs.splice(targetIndex, 0, removed)

      // Update field with reordered refs
      field.handleChange(newRefs as any)
      setDraggedIndex(null)
    },
    [draggedIndex, field],
  )

  // Extract error
  const error = useMemo(() => {
    const errors = field.state.meta.errors
    if (errors.length === 0) return undefined
    const firstError = errors[0]
    if (typeof firstError === 'string') {
      return t(firstError, { ns: 'errors' })
    }
    return undefined
  }, [field.state.meta.errors, t])

  const showError = field.state.meta.isTouched && error
  const totalFiles =
    existingReferences.length + newFiles.length + pendingFiles.length
  const isMaxReached = maxFiles ? totalFiles >= maxFiles : false

  // Calculate overall upload progress
  const overallProgress = useMemo(() => {
    const states = Array.from(uploadStates.values())
    if (states.length === 0) return 0
    const total = states.reduce((sum, state) => sum + state.progress, 0)
    return Math.round(total / states.length)
  }, [uploadStates])

  return (
    <div className="form-control w-full">
      {label && (
        <label className="label">
          <span className="label-text text-base-content/70 font-medium">
            {label}
            {props.required && <span className="text-error ms-1">*</span>}
          </span>
        </label>
      )}

      {/* Upload zone */}
      {!isMaxReached && (
        <div className="mb-4">
          <FileUploadZone
            accept={accept}
            multiple={multiple}
            disabled={disabled || isUploading}
            maxFiles={maxFiles}
            onFilesSelected={handleFilesSelected}
          >
            {isMobile && accept?.includes('image') && (
              <div className="flex gap-2 mt-4">
                <button
                  type="button"
                  className="btn btn-sm btn-outline"
                  onClick={(e) => {
                    e.stopPropagation()
                    const input = document.createElement('input')
                    input.type = 'file'
                    input.accept = accept
                    input.capture = 'environment'
                    input.onchange = (event) => {
                      const files = Array.from(
                        (event.target as HTMLInputElement).files || [],
                      )
                      handleFilesSelected(files)
                    }
                    input.click()
                  }}
                >
                  <Camera className="w-4 h-4" />
                  {t('takePicture', { ns: 'upload' })}
                </button>
                <button
                  type="button"
                  className="btn btn-sm btn-outline"
                  onClick={(e) => {
                    e.stopPropagation()
                    const input = document.createElement('input')
                    input.type = 'file'
                    input.accept = accept
                    input.multiple = multiple
                    input.onchange = (event) => {
                      const files = Array.from(
                        (event.target as HTMLInputElement).files || [],
                      )
                      handleFilesSelected(files)
                    }
                    input.click()
                  }}
                >
                  <ImageIcon className="w-4 h-4" />
                  {t('chooseFromGallery', { ns: 'upload' })}
                </button>
              </div>
            )}
          </FileUploadZone>
        </div>
      )}

      {/* File previews grid */}
      {totalFiles > 0 && (
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          {/* Existing references with drag-and-drop reordering */}
          {existingReferences.map((ref, index) => (
            <FilePreview
              key={ref.assetId}
              src={getThumbnailUrl(ref) || ''}
              alt={ref.metadata?.altText || 'Uploaded file'}
              onRemove={() => handleRemove(ref)}
              disabled={disabled}
              draggable={props.reorderable !== false}
              onDragStart={(e) => {
                e.dataTransfer.effectAllowed = 'move'
                handleDragStart(index)
              }}
              onDragEnd={handleDragEnd}
              onDragOver={handleDragOver}
              onDrop={(e) => {
                e.preventDefault()
                handleDrop(index)
              }}
            />
          ))}

          {/* Pending uploads with progress */}
          {pendingFiles.map((file) => {
            const uploadState = uploadStates.get(file)
            const objectURL = objectURLs.get(file)

            // For images, show preview
            if (file.type.startsWith('image/') && objectURL) {
              return (
                <FilePreview
                  key={file.name}
                  src={objectURL}
                  alt={file.name}
                  isLoading={uploadState?.status !== 'success'}
                  progress={uploadState?.progress || 0}
                  error={uploadState?.error}
                  onRemove={() => handleRemove(file)}
                  onRetry={() => upload([file])}
                  disabled={disabled}
                />
              )
            }

            // For non-images, show file thumbnail with progress
            return (
              <div key={file.name} className="relative">
                <FileThumbnail
                  fileName={file.name}
                  fileSize={file.size}
                  fileType={file.type}
                />
                {uploadState && uploadState.status !== 'success' && (
                  <div className="absolute inset-0 bg-base-300/80 flex items-center justify-center">
                    <UploadProgress
                      fileName={file.name}
                      fileSize={file.size}
                      progress={uploadState.progress}
                      status={t(`upload.${uploadState.status}`)}
                      onCancel={() => handleRemove(file)}
                    />
                  </div>
                )}
              </div>
            )
          })}
        </div>
      )}

      {/* Overall progress indicator */}
      {isUploading && (
        <div className="mt-4">
          <div className="alert alert-info">
            <span>
              {t('filesUploading', {
                count: pendingFiles.length,
                ns: 'upload',
              })}{' '}
              ({overallProgress}%)
            </span>
          </div>
        </div>
      )}

      {/* Helper text */}
      {hint && !showError && (
        <label className="label">
          <span className="label-text-alt text-base-content/60">{hint}</span>
        </label>
      )}

      {/* Error message */}
      {showError && (
        <label className="label">
          <span className="label-text-alt text-error" role="alert">
            {error}
          </span>
        </label>
      )}

      {/* Max files warning */}
      {isMaxReached && (
        <div className="alert alert-warning mt-2">
          <span>{t('maxFilesReached', { max: maxFiles, ns: 'upload' })}</span>
        </div>
      )}
    </div>
  )
}
