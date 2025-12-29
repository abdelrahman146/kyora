/**
 * InfiniteScroll Component
 *
 * A component that triggers a callback when the user scrolls near the bottom.
 * Used for infinite scrolling lists on mobile.
 *
 * Features:
 * - Intersection Observer API for performance
 * - Configurable threshold
 * - Loading indicator
 * - End of list detection
 */

import { useEffect, useRef } from 'react'
import type { ReactNode } from 'react'

export interface InfiniteScrollProps {
  children: ReactNode
  hasMore: boolean
  isLoading: boolean
  onLoadMore: () => void
  threshold?: number // pixels from bottom to trigger load
  loadingMessage?: string
  endMessage?: string
}

export function InfiniteScroll({
  children,
  hasMore,
  isLoading,
  onLoadMore,
  threshold = 200,
  loadingMessage = 'Loading more...',
  endMessage = 'No more items',
}: InfiniteScrollProps) {
  const observerTarget = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!hasMore || isLoading) return

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting) {
          onLoadMore()
        }
      },
      {
        root: null,
        rootMargin: `${String(threshold)}px`,
        threshold: 0,
      }
    )

    const currentTarget = observerTarget.current
    if (currentTarget) {
      observer.observe(currentTarget)
    }

    return () => {
      if (currentTarget) {
        observer.unobserve(currentTarget)
      }
    }
  }, [hasMore, isLoading, onLoadMore, threshold])

  return (
    <>
      {children}

      {/* Observer Target */}
      <div ref={observerTarget} className="py-4 text-center">
        {isLoading && (
          <div className="flex items-center justify-center gap-2">
            <span className="loading loading-spinner loading-sm"></span>
            <span className="text-sm text-base-content/60">
              {loadingMessage}
            </span>
          </div>
        )}
        {!isLoading && !hasMore && (
          <div className="text-sm text-base-content/40">{endMessage}</div>
        )}
      </div>
    </>
  )
}
