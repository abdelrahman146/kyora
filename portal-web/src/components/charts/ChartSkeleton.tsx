import { cn } from '@/lib/utils'

export interface ChartSkeletonProps {
  variant?: 'line' | 'bar' | 'pie' | 'doughnut' | 'mixed'
  height?: number
  className?: string
}

export const ChartSkeleton = ({
  variant = 'line',
  height = 320,
  className,
}: ChartSkeletonProps) => {
  return (
    <div
      className={cn('relative w-full animate-pulse', className)}
      style={{ height }}
      role="status"
      aria-label="Loading chart..."
    >
      <div className="flex h-full w-full flex-col gap-2">
        {/* Legend skeleton */}
        <div className="flex items-center justify-center gap-4">
          <div className="h-3 w-20 rounded bg-base-300" />
          <div className="h-3 w-20 rounded bg-base-300" />
        </div>

        {/* Chart area */}
        <div className="relative flex-1">
          {variant === 'line' && <LineChartSkeleton />}
          {variant === 'bar' && <BarChartSkeleton />}
          {variant === 'pie' && <PieChartSkeleton />}
          {variant === 'doughnut' && <PieChartSkeleton />}
          {variant === 'mixed' && <MixedChartSkeleton />}
        </div>

        {/* X-axis labels skeleton */}
        {(variant === 'line' || variant === 'bar' || variant === 'mixed') && (
          <div className="flex justify-between px-2">
            <div className="h-2 w-12 rounded bg-base-300" />
            <div className="h-2 w-12 rounded bg-base-300" />
            <div className="h-2 w-12 rounded bg-base-300" />
            <div className="h-2 w-12 rounded bg-base-300" />
          </div>
        )}
      </div>
    </div>
  )
}

function LineChartSkeleton() {
  return (
    <div className="absolute inset-0 flex items-end justify-around px-4 pb-4">
      {[60, 75, 55, 85, 70, 90, 65, 80].map((height, i) => (
        <div key={i} className="flex w-full flex-col items-center justify-end">
          <div
            className="w-1 rounded-full bg-base-300"
            style={{ height: `${height}%` }}
          />
        </div>
      ))}
    </div>
  )
}

function BarChartSkeleton() {
  return (
    <div className="absolute inset-0 flex items-end justify-around gap-2 px-4 pb-4">
      {[60, 75, 55, 85, 70, 90].map((height, i) => (
        <div
          key={i}
          className="w-full rounded-t bg-base-300"
          style={{ height: `${height}%` }}
        />
      ))}
    </div>
  )
}

function PieChartSkeleton() {
  return (
    <div className="absolute inset-0 flex items-center justify-center">
      <div className="h-48 w-48 rounded-full bg-base-300" />
    </div>
  )
}

function MixedChartSkeleton() {
  return (
    <div className="absolute inset-0 flex items-end justify-around gap-1 px-4 pb-4">
      {[
        { type: 'bar', height: 60 },
        { type: 'line', height: 75 },
        { type: 'bar', height: 55 },
        { type: 'line', height: 85 },
        { type: 'bar', height: 70 },
        { type: 'line', height: 90 },
      ].map((item, i) => (
        <div key={i} className="flex w-full flex-col items-center justify-end">
          <div
            className={cn(
              'bg-base-300',
              item.type === 'bar' ? 'w-full rounded-t' : 'w-1 rounded-full',
            )}
            style={{ height: `${item.height}%` }}
          />
        </div>
      ))}
    </div>
  )
}
