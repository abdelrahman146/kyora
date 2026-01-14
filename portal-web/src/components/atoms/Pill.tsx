import type { ReactNode } from 'react'

import { cn } from '@/lib/utils'

export type PillVariant =
  | 'primary'
  | 'secondary'
  | 'success'
  | 'error'
  | 'warning'
  | 'info'
  | 'neutral'
  | 'ghost'

export type PillSize = 'sm' | 'md' | 'lg'

export interface PillProps {
  /** Main text content */
  children: ReactNode
  /** Optional icon displayed before the text */
  icon?: ReactNode
  /** Optional secondary text (e.g., timestamp) */
  secondary?: ReactNode
  /** Visual variant */
  variant?: PillVariant
  /** Size variant */
  size?: PillSize
  /** Additional CSS classes */
  className?: string
}

const variantClasses: Record<PillVariant, string> = {
  primary: 'bg-primary/15 text-primary border-primary/20',
  secondary: 'bg-secondary/15 text-secondary border-secondary/20',
  success: 'bg-success/15 text-success border-success/20',
  error: 'bg-error/15 text-error border-error/20',
  warning: 'bg-warning/15 text-warning border-warning/20',
  info: 'bg-info/15 text-info border-info/20',
  neutral: 'bg-base-200 text-base-content border-base-300',
  ghost: 'bg-transparent text-base-content/70 border-base-300',
}

const sizeClasses: Record<PillSize, string> = {
  sm: 'px-2 py-0.5 text-[11px] gap-1',
  md: 'px-2.5 py-1 text-xs gap-1.5',
  lg: 'px-3 py-1.5 text-sm gap-2',
}

const iconSizeClasses: Record<PillSize, string> = {
  sm: '[&>svg]:size-3',
  md: '[&>svg]:size-3.5',
  lg: '[&>svg]:size-4',
}

/**
 * Pill Component
 *
 * A compact status indicator with optional icon and secondary text.
 * Designed for mobile-first, responsive layouts.
 *
 * @example
 * ```tsx
 * <Pill icon={<CheckCircle2 />} variant="success">
 *   Ready
 * </Pill>
 *
 * <Pill icon={<Clock />} variant="info" secondary="2 min ago">
 *   Updated
 * </Pill>
 * ```
 */
export function Pill({
  children,
  icon,
  secondary,
  variant = 'neutral',
  size = 'md',
  className,
}: PillProps) {
  return (
    <span
      className={cn(
        // Base styles
        'inline-flex items-center rounded-full border font-medium',
        // Prevent text wrapping on small screens
        'whitespace-nowrap',
        // Variant and size
        variantClasses[variant],
        sizeClasses[size],
        iconSizeClasses[size],
        className,
      )}
    >
      {icon}
      <span>{children}</span>
      {secondary && <span className="opacity-70 font-normal">{secondary}</span>}
    </span>
  )
}
