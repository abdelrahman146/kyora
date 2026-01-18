/**
 * AdvisorPanel Component
 *
 * Friendly advisor-style panel for displaying contextual insights and recommendations.
 * Designed to feel helpful rather than alarming, with a conversational tone.
 *
 * Features:
 * - Lightbulb icon for helpful suggestions
 * - Soft, non-alarming colors
 * - Multiple insights support
 * - Clean, card-based design
 * - RTL-compatible
 * - Mobile-first responsive
 *
 * Usage:
 * - Desktop: Sticky sidebar alongside main content
 * - Mobile: Inline card after hero section
 *
 * @example
 * ```tsx
 * <AdvisorPanel
 *   title="Financial Advisor"
 *   insights={[
 *     { type: 'positive', message: 'Your profit margin is healthy.' },
 *     { type: 'suggestion', message: 'Consider reducing expenses.' }
 *   ]}
 * />
 * ```
 */
import {
  AlertCircle,
  Info,
  Lightbulb,
  Sparkles,
  TrendingUp,
} from 'lucide-react'
import type { LucideIcon } from 'lucide-react'

import { cn } from '@/lib/utils'

export interface AdvisorInsight {
  /** Insight type determines icon and subtle styling */
  type: 'positive' | 'suggestion' | 'alert' | 'info'
  /** The insight message */
  message: string
  /** Optional custom icon */
  icon?: LucideIcon
}

export interface AdvisorPanelProps {
  /** Panel title */
  title: string
  /** List of insights to display */
  insights: Array<AdvisorInsight>
  /** Additional CSS classes */
  className?: string
}

const insightConfig = {
  positive: {
    icon: TrendingUp,
    iconColor: 'text-success',
    dotColor: 'bg-success',
  },
  suggestion: {
    icon: Lightbulb,
    iconColor: 'text-warning',
    dotColor: 'bg-warning',
  },
  alert: {
    icon: AlertCircle,
    iconColor: 'text-error',
    dotColor: 'bg-error',
  },
  info: {
    icon: Info,
    iconColor: 'text-info',
    dotColor: 'bg-info',
  },
} as const

export function AdvisorPanel({
  title,
  insights,
  className,
}: AdvisorPanelProps) {
  if (insights.length === 0) return null

  return (
    <div
      className={cn(
        'card bg-gradient-to-br from-primary/5 to-secondary/5 border border-base-300',
        className,
      )}
      role="complementary"
      aria-label={title}
    >
      {/* Header */}
      <div className="card-body p-4 sm:p-6">
        <div className="flex items-center gap-3 mb-4">
          <div className="flex h-10 w-10 items-center justify-center rounded-full bg-primary/10">
            <Sparkles className="h-5 w-5 text-primary" aria-hidden="true" />
          </div>
          <h3 className="text-lg font-semibold text-base-content">{title}</h3>
        </div>

        {/* Insights List */}
        <div className="space-y-4">
          {insights.map((insight, index) => {
            const config = insightConfig[insight.type]
            const Icon = insight.icon ?? config.icon

            return (
              <div
                key={index}
                className="flex items-start gap-3 rounded-lg bg-base-100 p-3 transition-shadow hover:shadow-sm"
              >
                {/* Icon with colored dot */}
                <div className="relative shrink-0">
                  <Icon
                    className={cn('h-5 w-5', config.iconColor)}
                    aria-hidden="true"
                  />
                  <span
                    className={cn(
                      'absolute -end-1 -top-1 h-2 w-2 rounded-full',
                      config.dotColor,
                    )}
                    aria-hidden="true"
                  />
                </div>

                {/* Message */}
                <p className="flex-1 text-sm leading-relaxed text-base-content/80">
                  {insight.message}
                </p>
              </div>
            )
          })}
        </div>
      </div>
    </div>
  )
}

AdvisorPanel.displayName = 'AdvisorPanel'
