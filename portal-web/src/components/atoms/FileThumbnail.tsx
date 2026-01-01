import { Archive, File, FileText, Music, Video } from 'lucide-react'
import { cn } from '@/lib/utils'
import { formatSize } from '@/lib/upload/fileValidator'

export interface FileThumbnailProps {
  fileName: string
  fileSize: number
  fileType: string
  className?: string
}

/**
 * Get appropriate icon for file type
 */
function getFileIcon(mimeType: string) {
  if (
    mimeType.startsWith('application/pdf') ||
    mimeType.includes('document') ||
    mimeType.includes('text')
  ) {
    return FileText
  }
  if (mimeType.startsWith('video/')) {
    return Video
  }
  if (mimeType.startsWith('audio/')) {
    return Music
  }
  if (
    mimeType.includes('zip') ||
    mimeType.includes('tar') ||
    mimeType.includes('rar') ||
    mimeType.includes('7z')
  ) {
    return Archive
  }
  return File
}

/**
 * FileThumbnail - Generic file thumbnail for non-image files
 *
 * Features:
 * - Type-specific icons (document, video, audio, archive)
 * - File name and size display
 * - Consistent sizing with image previews
 * - RTL support
 *
 * @example
 * <FileThumbnail
 *   fileName="report.pdf"
 *   fileSize={1024000}
 *   fileType="application/pdf"
 * />
 */
export function FileThumbnail({
  fileName,
  fileSize,
  fileType,
  className,
}: FileThumbnailProps) {
  const Icon = getFileIcon(fileType)

  return (
    <div
      className={cn(
        'flex flex-col items-center justify-center',
        'aspect-square rounded-lg',
        'border-2 border-base-300',
        'bg-base-200 p-4',
        className,
      )}
    >
      <Icon
        className="w-12 h-12 md:w-16 md:h-16 text-base-content/40 mb-2"
        aria-hidden="true"
      />

      <p
        className="text-sm font-medium text-base-content text-center truncate w-full px-2"
        title={fileName}
      >
        {fileName}
      </p>

      <p className="text-xs text-base-content/60 mt-1">
        {formatSize(fileSize)}
      </p>
    </div>
  )
}
