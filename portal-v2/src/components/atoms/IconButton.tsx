import { forwardRef } from 'react'
import { Loader2 } from 'lucide-react'
import { cn } from '../../lib/utils'
import type { LucideIcon } from 'lucide-react'
import type { ButtonHTMLAttributes } from 'react'

export interface IconButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  /**
   * Icon component from lucide-react
   */
  icon: LucideIcon

  /**
   * Size variant
   * @default 'md'
   */
  size?: 'sm' | 'md' | 'lg'

  /**
   * Style variant
   * @default 'ghost'
   */
  variant?: 'ghost' | 'outline' | 'primary'

  /**
   * Loading state
   */
  loading?: boolean

  /**
   * Additional CSS classes
   */
  className?: string
}

/**
 * Icon-Only Button Component
 *
 * Optimized for icon-only actions (close, settings, etc.).
 * Built with daisyUI 5 btn and btn-square utilities.
 *
 * @example
 * ```tsx
 * <IconButton icon={Settings} variant="ghost" size="md" />
 * <IconButton icon={X} variant="outline" onClick={onClose} />
 * <IconButton icon={Save} variant="primary" loading={isSaving} />
 * ```
 */
export const IconButton = forwardRef<HTMLButtonElement, IconButtonProps>(
  (
    {
      icon: Icon,
      size = 'md',
      variant = 'ghost',
      loading = false,
      className,
      disabled,
      ...props
    },
    ref
  ) => {
    const sizeClasses = {
      sm: 'h-8 w-8 btn-sm',
      md: 'h-10 w-10',
      lg: 'h-12 w-12 btn-lg',
    }

    const variantClasses = {
      ghost: 'btn-ghost',
      outline: 'btn-outline',
      primary: 'btn-primary',
    }

    return (
      <button
        ref={ref}
        disabled={disabled || loading}
        className={cn(
          'btn btn-square',
          sizeClasses[size],
          variantClasses[variant],
          className
        )}
        {...props}
      >
        {loading ? (
          <Loader2 className="h-5 w-5 animate-spin" />
        ) : (
          <Icon className="h-5 w-5" />
        )}
      </button>
    )
  }
)

IconButton.displayName = 'IconButton'
