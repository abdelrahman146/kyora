import type { ReactNode } from 'react'
import { cn } from '@/lib/utils'

export interface AvatarProps {
  src?: string
  alt?: string
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl'
  placeholder?: ReactNode
  fallback?: string
  shape?: 'circle' | 'square'
  online?: boolean
  offline?: boolean
  className?: string
}

/**
 * Avatar Component
 *
 * User avatar with placeholder support and online/offline indicators.
 */
export function Avatar({
  src,
  alt,
  size = 'md',
  placeholder,
  fallback,
  shape = 'circle',
  online,
  offline,
  className,
}: AvatarProps) {
  const sizeClasses = {
    xs: 'w-6',
    sm: 'w-8',
    md: 'w-12',
    lg: 'w-16',
    xl: 'w-24',
  }

  const shapeClass = shape === 'circle' ? 'rounded-full' : 'rounded-lg'

  const onlineIndicator = online ?? offline

  return (
    <div className={cn('avatar', onlineIndicator && 'online', className)}>
      {online && (
        <div className="indicator-item badge badge-success badge-xs"></div>
      )}
      {offline && (
        <div className="indicator-item badge badge-error badge-xs"></div>
      )}
      <div className={cn(sizeClasses[size], shapeClass)}>
        {src ? (
          <img src={src} alt={alt ?? 'Avatar'} />
        ) : (
          <div className="flex h-full w-full items-center justify-center bg-neutral text-neutral-content">
            {fallback ? (
              <span>{fallback.slice(0, 2).toUpperCase()}</span>
            ) : (
              (placeholder ?? <span>?</span>)
            )}
          </div>
        )}
      </div>
    </div>
  )
}
