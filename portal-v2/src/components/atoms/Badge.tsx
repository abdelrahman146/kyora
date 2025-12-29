import type {ReactNode} from 'react';
import { cn } from '@/lib/utils'

export interface BadgeProps {
  children: ReactNode
  variant?: 'primary' | 'secondary' | 'success' | 'error' | 'warning' | 'info' | 'neutral' | 'ghost'
  size?: 'xs' | 'sm' | 'md' | 'lg'
  outline?: boolean
  className?: string
}

/**
 * Badge Component
 *
 * Small status/label indicator using daisyUI badge classes.
 */
export function Badge({
  children,
  variant = 'neutral',
  size = 'md',
  outline = false,
  className,
}: BadgeProps) {
  const baseClasses = 'badge'

  const variantClasses = {
    primary: 'badge-primary',
    secondary: 'badge-secondary',
    success: 'badge-success',
    error: 'badge-error',
    warning: 'badge-warning',
    info: 'badge-info',
    neutral: 'badge-neutral',
    ghost: 'badge-ghost',
  }

  const sizeClasses = {
    xs: 'badge-xs',
    sm: 'badge-sm',
    md: 'badge-md',
    lg: 'badge-lg',
  }

  return (
    <span
      className={cn(
        baseClasses,
        variantClasses[variant],
        sizeClasses[size],
        outline && 'badge-outline',
        className,
      )}
    >
      {children}
    </span>
  )
}
