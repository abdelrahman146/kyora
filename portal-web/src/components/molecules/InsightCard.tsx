/**
 * InsightCard Component
 *
 * Displays contextual insights with icon, variant styling, and optional action link.
 * Used in financial reports to show actionable suggestions based on metrics.
 *
 * Variants:
 * - success: Green background, CheckCircle icon (positive metrics)
 * - warning: Yellow background, AlertTriangle icon (concerning metrics)
 * - error: Red background, AlertCircle icon (critical issues)
 * - info: Blue background, Info icon (tips and information)
 *
 * @example
 * ```tsx
 * <InsightCard
 *   variant="success"
 *   message="Your business is in a healthy position."
 * />
 *
 * <InsightCard
 *   variant="warning"
 *   message="Your cash position is low."
 *   link={{
 *     label: "View Expenses",
 *     href: "/business/$businessDescriptor/accounting/expenses",
 *     params: { businessDescriptor: "my-business" }
 *   }}
 * />
 * ```
 */
import { Link } from '@tanstack/react-router'
import {
  AlertCircle,
  AlertTriangle,
  CheckCircle,
  ChevronLeft,
  ChevronRight,
  Info,
} from 'lucide-react'
import type { LucideIcon } from 'lucide-react'

import { cn } from '@/lib/utils'
import { useLanguage } from '@/hooks/useLanguage'

export interface InsightCardProps {
  /** Visual style and default icon selection */
  variant: 'success' | 'warning' | 'error' | 'info'
  /** The insight message to display */
  message: string
  /** Optional custom icon (overrides default variant icon) */
  icon?: LucideIcon
  /** Optional navigation link */
  link?: {
    label: string
    href: string
    params?: Record<string, string>
  }
  /** Additional CSS classes */
  className?: string
}

const variantConfig = {
  success: {
    bg: 'bg-success/10',
    border: 'border-success/20',
    text: 'text-success',
    defaultIcon: CheckCircle,
  },
  warning: {
    bg: 'bg-warning/10',
    border: 'border-warning/20',
    text: 'text-warning',
    defaultIcon: AlertTriangle,
  },
  error: {
    bg: 'bg-error/10',
    border: 'border-error/20',
    text: 'text-error',
    defaultIcon: AlertCircle,
  },
  info: {
    bg: 'bg-info/10',
    border: 'border-info/20',
    text: 'text-info',
    defaultIcon: Info,
  },
} as const

export function InsightCard({
  variant,
  message,
  icon,
  link,
  className,
}: InsightCardProps) {
  const { isRTL } = useLanguage()
  const config = variantConfig[variant]
  const Icon = icon ?? config.defaultIcon
  const ChevronIcon = isRTL ? ChevronLeft : ChevronRight

  return (
    <div
      className={cn(
        'rounded-box border p-4',
        config.bg,
        config.border,
        className,
      )}
      role="note"
      aria-label={message}
    >
      <div className="flex items-start gap-3">
        <Icon
          className={cn('h-5 w-5 shrink-0 mt-0.5', config.text)}
          aria-hidden="true"
        />
        <div className="flex-1 space-y-2">
          <p className="text-sm text-base-content">{message}</p>
          {link && (
            <Link
              to={link.href}
              params={link.params}
              className={cn(
                'inline-flex items-center gap-1 text-sm font-medium',
                config.text,
                'hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2',
              )}
            >
              {link.label}
              <ChevronIcon className="h-4 w-4" aria-hidden="true" />
            </Link>
          )}
        </div>
      </div>
    </div>
  )
}

InsightCard.displayName = 'InsightCard'
