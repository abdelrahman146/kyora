import { cn } from '../../lib/utils'

export interface LogoProps {
  /**
   * Whether to show text alongside icon
   * @default true
   */
  showText?: boolean

  /**
   * Size variant
   * @default 'md'
   */
  size?: 'sm' | 'md' | 'lg'

  /**
   * Additional CSS classes
   */
  className?: string
}

/**
 * Kyora Logo Component
 *
 * Displays the Kyora brand icon (K) with optional text.
 * Responsive, supports multiple sizes.
 *
 * @example
 * ```tsx
 * <Logo />
 * <Logo size="lg" />
 * <Logo showText={false} />
 * ```
 */
export function Logo({ showText = true, size = 'md', className }: LogoProps) {
  const sizeClasses = {
    sm: 'h-8',
    md: 'h-10',
    lg: 'h-12',
  }

  const iconSizeClasses = {
    sm: 'text-base',
    md: 'text-lg',
    lg: 'text-xl',
  }

  const textSizeClasses = {
    sm: 'text-lg',
    md: 'text-xl',
    lg: 'text-2xl',
  }

  return (
    <div className={cn('flex items-center gap-2', className)}>
      {/* Icon */}
      <div
        className={cn(
          'flex aspect-square items-center justify-center rounded-lg bg-primary text-primary-content font-bold',
          sizeClasses[size],
          iconSizeClasses[size],
        )}
      >
        K
      </div>

      {/* Text */}
      {showText && (
        <span
          className={cn('font-bold text-base-content', textSizeClasses[size])}
        >
          Kyora
        </span>
      )}
    </div>
  )
}
