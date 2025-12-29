import { ChevronLeft, ChevronRight } from 'lucide-react'
import { cn } from '@/lib/utils'

export interface PaginationProps {
  currentPage: number
  totalPages: number
  onPageChange: (page: number) => void
  showPageNumbers?: boolean
  className?: string
}

/**
 * Pagination Component
 *
 * Page navigation controls with prev/next buttons.
 * RTL-aware using logical properties.
 */
export function Pagination({
  currentPage,
  totalPages,
  onPageChange,
  showPageNumbers = true,
  className,
}: PaginationProps) {
  const canGoPrevious = currentPage > 1
  const canGoNext = currentPage < totalPages

  return (
    <div className={cn('flex items-center justify-center gap-2', className)}>
      <div className="join">
        <button
          className="btn join-item btn-sm"
          onClick={() => onPageChange(currentPage - 1)}
          disabled={!canGoPrevious}
          aria-label="Previous page"
        >
          <ChevronRight size={16} />
          السابق
        </button>

        {showPageNumbers && (
          <button className="btn join-item btn-sm pointer-events-none">
            صفحة {currentPage} من {totalPages}
          </button>
        )}

        <button
          className="btn join-item btn-sm"
          onClick={() => onPageChange(currentPage + 1)}
          disabled={!canGoNext}
          aria-label="Next page"
        >
          التالي
          <ChevronLeft size={16} />
        </button>
      </div>
    </div>
  )
}
