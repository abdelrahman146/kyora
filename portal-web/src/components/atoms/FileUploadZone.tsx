import { useRef, useState } from 'react'
import { CloudUpload } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { cn } from '@/lib/utils'

export interface FileUploadZoneProps {
  onFilesSelected: (files: Array<File>) => void
  accept?: string
  multiple?: boolean
  disabled?: boolean
  maxFiles?: number
  className?: string
  children?: React.ReactNode
  /** When true, indicates the max file count has been reached */
  isMaxReached?: boolean
  /** Whether uploads are currently in progress */
  isUploading?: boolean
  /** Optional label to display while uploading */
  uploadingLabel?: string
}

/**
 * FileUploadZone - Drag-and-drop file upload zone
 *
 * Features:
 * - Drag-and-drop file upload
 * - Click to select files
 * - Visual feedback on drag-over
 * - Mobile-optimized layout
 * - Keyboard accessible
 * - RTL support
 *
 * @example
 * <FileUploadZone
 *   accept="image/*"
 *   multiple
 *   maxFiles={10}
 *   onFilesSelected={(files) => console.log(files)}
 * />
 */
export function FileUploadZone({
  onFilesSelected,
  accept,
  multiple = false,
  disabled = false,
  maxFiles,
  className,
  children,
  isMaxReached = false,
  isUploading = false,
  uploadingLabel,
}: FileUploadZoneProps) {
  const [isDragActive, setIsDragActive] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)
  const { t } = useTranslation('upload')
  const isMaxFilesReached = Boolean(maxFiles && isMaxReached)
  const isZoneDisabled = disabled || isMaxFilesReached

  const handleDragEnter = (e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    if (!isZoneDisabled) {
      setIsDragActive(true)
    }
  }

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragActive(false)
  }

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
  }

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragActive(false)

    if (isZoneDisabled) return

    const files = Array.from(e.dataTransfer.files)
    if (files.length > 0) {
      onFilesSelected(files)
    }
  }

  const handleClick = () => {
    if (!isZoneDisabled) {
      fileInputRef.current?.click()
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if ((e.key === ' ' || e.key === 'Enter') && !isZoneDisabled) {
      e.preventDefault()
      handleClick()
    }
  }

  const handleFileInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (isZoneDisabled) return

    const files = Array.from(e.target.files || [])
    if (files.length > 0) {
      onFilesSelected(files)
    }
    // Reset input so same file can be selected again
    e.target.value = ''
  }

  return (
    <div
      role="button"
      tabIndex={isZoneDisabled ? -1 : 0}
      onClick={handleClick}
      onKeyDown={handleKeyDown}
      onDragEnter={handleDragEnter}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
      aria-label={
        isMaxFilesReached
          ? t('maximumFilesReached')
          : isUploading
            ? (uploadingLabel ?? t('uploading'))
            : t('selectFiles')
      }
      aria-disabled={isZoneDisabled}
      className={cn(
        'relative flex flex-col items-center justify-center',
        'min-h-[200px] p-8',
        'border-2 border-dashed rounded-lg',
        'transition-all duration-200',
        !isZoneDisabled && 'cursor-pointer',
        isDragActive &&
          !isZoneDisabled &&
          'border-primary bg-primary/10 scale-[1.02]',
        !isDragActive &&
          !isZoneDisabled &&
          'border-base-300 hover:border-primary/50 hover:bg-base-200/50',
        isZoneDisabled &&
          'border-base-300/50 bg-base-200/30 cursor-not-allowed opacity-60',
        className,
      )}
    >
      <input
        ref={fileInputRef}
        type="file"
        accept={accept}
        multiple={multiple}
        onChange={handleFileInputChange}
        disabled={isZoneDisabled}
        className="hidden"
        aria-hidden="true"
      />

      {children || (
        <>
          <CloudUpload
            className={cn(
              'w-12 h-12 mb-4 md:w-16 md:h-16',
              isDragActive && !disabled
                ? 'text-primary'
                : 'text-base-content/40',
            )}
            aria-hidden="true"
          />

          <p className="text-base md:text-lg font-medium text-base-content/80 mb-2 text-center">
            {isDragActive && !disabled ? (
              t('dropFilesHere')
            ) : isMaxFilesReached ? (
              t('maximumFilesReached')
            ) : isUploading ? (
              (uploadingLabel ?? t('waitingForUploads'))
            ) : (
              <>
                <span className="hidden md:inline">
                  {t('clickToUploadOrDragDrop')}
                </span>
                <span className="inline md:hidden">{t('tapToUpload')}</span>
              </>
            )}
          </p>

          {accept && !isMaxFilesReached && (
            <p
              id="upload-requirements"
              className="text-sm text-base-content/60 text-center"
            >
              {accept.includes('image') && t('imageFormatsUpTo10MB')}
              {accept.includes('video') && t('videoFormatsUpTo100MB')}
              {!accept.includes('image') &&
                !accept.includes('video') &&
                t('acceptedFileTypes')}
            </p>
          )}

          {maxFiles && !isMaxFilesReached && (
            <p className="text-xs text-base-content/50 mt-2">
              {t('maximumFiles', { count: maxFiles })}
            </p>
          )}
        </>
      )}
    </div>
  )
}
