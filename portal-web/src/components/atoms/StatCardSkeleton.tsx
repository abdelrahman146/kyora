import { cn } from '@/lib/utils'

export interface StatCardSkeletonProps {
  variant?: 'simple' | 'complex'
  className?: string
}

export const StatCardSkeleton = ({
  variant = 'simple',
  className,
}: StatCardSkeletonProps) => {
  return (
    <div
      className={cn(
        'card rounded-box bg-base-100 border border-base-300 animate-pulse',
        className,
      )}
      role="status"
      aria-label="Loading statistics..."
    >
      <div className="card-body p-4">
        {variant === 'simple' ? (
          <SimpleStatSkeleton />
        ) : (
          <ComplexStatSkeleton />
        )}
      </div>
    </div>
  )
}

function SimpleStatSkeleton() {
  return (
    <div className="flex items-start justify-between gap-4">
      <div className="flex flex-col gap-2 flex-1">
        <div className="h-3 w-20 rounded bg-base-300" />
        <div className="h-8 w-32 rounded bg-base-300" />
      </div>
      <div className="h-10 w-10 rounded-full bg-base-300" />
    </div>
  )
}

function ComplexStatSkeleton() {
  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-start justify-between gap-4">
        <div className="flex flex-col gap-2 flex-1">
          <div className="h-3 w-24 rounded bg-base-300" />
          <div className="h-9 w-40 rounded bg-base-300" />
        </div>
        <div className="h-6 w-16 rounded-full bg-base-300" />
      </div>
      <div className="flex items-center gap-4 pt-2 border-t border-base-300">
        <div className="h-3 w-28 rounded bg-base-300" />
        <div className="h-3 w-24 rounded bg-base-300" />
      </div>
    </div>
  )
}
