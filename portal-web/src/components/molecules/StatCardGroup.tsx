import type { ReactNode } from 'react'
import { cn } from '@/lib/utils'

export interface StatCardGroupProps {
  children: ReactNode
  className?: string
  cols?: 1 | 2 | 3 | 4
}

/**
 * StatCardGroup Component
 *
 * Responsive grid wrapper for multiple StatCard components.
 * Automatically adapts to screen size:
 * - Mobile: 1 column
 * - Tablet: 2 columns
 * - Desktop: 3-4 columns (configurable)
 */
export const StatCardGroup = ({
  children,
  className,
  cols = 3,
}: StatCardGroupProps) => {
  const colClasses = {
    1: 'grid-cols-1',
    2: 'grid-cols-1 md:grid-cols-2',
    3: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3',
    4: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-4',
  }

  return (
    <div className={cn('grid gap-4', colClasses[cols], className)}>
      {children}
    </div>
  )
}
