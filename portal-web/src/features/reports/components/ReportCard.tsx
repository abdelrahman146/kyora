/**
 * ReportCard Component
 *
 * Clickable card for the Reports Hub page displaying:
 * - Icon and title
 * - Key metric (large, prominent)
 * - Secondary metrics (smaller, supporting)
 * - "View Details" CTA
 *
 * Entire card is clickable with proper touch targets (min 44px).
 * Supports RTL layout and mobile-first responsive design.
 *
 * @example
 * ```tsx
 * <ReportCard
 *   title="Business Health"
 *   icon={Heart}
 *   keyMetric={{ label: "What Your Business Is Worth", value: "$25,000" }}
 *   secondaryValues={[
 *     { label: "Cash", value: "$8,000" },
 *     { label: "Inventory", value: "$15,000" }
 *   ]}
 *   href="/business/$businessDescriptor/reports/health"
 *   params={{ businessDescriptor: "my-business" }}
 * />
 * ```
 */
import { Link } from '@tanstack/react-router'
import { ChevronLeft, ChevronRight } from 'lucide-react'
import { useTranslation } from 'react-i18next'

import type { LucideIcon } from 'lucide-react'

import { Skeleton } from '@/components/atoms/Skeleton'
import { useLanguage } from '@/hooks/useLanguage'
import { cn } from '@/lib/utils'

export interface ReportCardProps {
  /** Card title displayed next to icon */
  title: string
  /** Icon displayed in the header */
  icon: LucideIcon
  /** Primary metric displayed prominently */
  keyMetric: {
    label: string
    value: string
  }
  /** Optional secondary metrics displayed below key metric */
  secondaryValues?: Array<{
    label: string
    value: string
  }>
  /** Navigation target path */
  href: string
  /** Route params for dynamic segments */
  params?: Record<string, string>
  /** Search params to append to URL */
  searchParams?: Record<string, string>
  /** Additional CSS classes */
  className?: string
}

export function ReportCard({
  title,
  icon: Icon,
  keyMetric,
  secondaryValues = [],
  href,
  params,
  searchParams,
  className,
}: ReportCardProps) {
  const { t } = useTranslation('reports')
  const { isRTL } = useLanguage()
  const ChevronIcon = isRTL ? ChevronLeft : ChevronRight

  return (
    <Link
      to={href}
      params={params}
      search={searchParams}
      className={cn(
        'block rounded-box border border-base-300 p-4',
        'transition-colors hover:bg-base-200/50 active:scale-[0.98]',
        'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2',
        className,
      )}
    >
      {/* Header: Icon + Title */}
      <div className="flex items-center gap-2 mb-3">
        <Icon className="h-5 w-5 text-primary" aria-hidden="true" />
        <h3 className="text-base font-semibold text-base-content">{title}</h3>
      </div>

      {/* Key Metric */}
      <div className="mb-2">
        <p className="text-xs text-base-content/60 mb-1">{keyMetric.label}</p>
        <p className="text-2xl font-bold text-base-content">
          {keyMetric.value}
        </p>
      </div>

      {/* Secondary Values */}
      {secondaryValues.length > 0 && (
        <div className="flex flex-wrap gap-x-4 gap-y-1 mb-3 text-sm text-base-content/70">
          {secondaryValues.map((item, index) => (
            <span key={index}>
              {item.label}: <span className="font-medium">{item.value}</span>
            </span>
          ))}
        </div>
      )}

      {/* CTA */}
      <div className="flex items-center justify-end gap-1 text-sm font-medium text-primary">
        {t('cards.view_details')}
        <ChevronIcon className="h-4 w-4" aria-hidden="true" />
      </div>
    </Link>
  )
}

ReportCard.displayName = 'ReportCard'

export interface ReportCardSkeletonProps {
  /** Additional CSS classes */
  className?: string
}

export function ReportCardSkeleton({ className }: ReportCardSkeletonProps) {
  return (
    <div className={cn('rounded-box border border-base-300 p-4', className)}>
      {/* Header skeleton */}
      <div className="flex items-center gap-2 mb-3">
        <Skeleton className="h-5 w-5 rounded" />
        <Skeleton className="h-5 w-32" />
      </div>

      {/* Key metric skeleton */}
      <div className="mb-2">
        <Skeleton className="h-3 w-24 mb-2" />
        <Skeleton className="h-8 w-40" />
      </div>

      {/* Secondary values skeleton */}
      <div className="flex gap-4 mb-3">
        <Skeleton className="h-4 w-20" />
        <Skeleton className="h-4 w-24" />
      </div>

      {/* CTA skeleton */}
      <div className="flex justify-end">
        <Skeleton className="h-4 w-24" />
      </div>
    </div>
  )
}

ReportCardSkeleton.displayName = 'ReportCardSkeleton'
