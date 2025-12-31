/**
 * SocialMediaHandles Component
 *
 * A production-grade molecule component that displays all social media handles for a customer.
 * Automatically filters out empty values and provides a clean, responsive grid layout.
 *
 * Features:
 * - Mobile-first responsive grid (1-2 columns based on screen size)
 * - RTL/LTR support with logical properties
 * - DRY principle - reuses SocialMediaLink atom
 * - Empty state handling
 * - Accessible with proper semantic HTML
 * - Smooth animations and transitions
 *
 * @example
 * ```tsx
 * <SocialMediaHandles
 *   instagramUsername="john_doe"
 *   facebookUsername="johndoe"
 *   whatsappNumber="+971501234567"
 * />
 * ```
 */

import { useMemo } from 'react'
import { cn } from '../../lib/utils'
import { SocialMediaLink } from '../atoms/SocialMediaLink'
import type { SocialPlatform } from '../icons/social'

export interface SocialMediaHandlesProps {
  instagramUsername?: string | null
  facebookUsername?: string | null
  tiktokUsername?: string | null
  snapchatUsername?: string | null
  xUsername?: string | null
  whatsappNumber?: string | null
  size?: 'sm' | 'md' | 'lg'
  variant?: 'default' | 'minimal'
  className?: string
}

interface SocialHandle {
  platform: SocialPlatform
  username: string
}

export function SocialMediaHandles({
  instagramUsername,
  facebookUsername,
  tiktokUsername,
  snapchatUsername,
  xUsername,
  whatsappNumber,
  size = 'md',
  variant = 'default',
  className,
}: SocialMediaHandlesProps) {
  // Filter and map social handles - memoized for performance
  const handles = useMemo<Array<SocialHandle>>(() => {
    const allHandles: Array<{
      platform: SocialPlatform
      username: string | null | undefined
    }> = [
      { platform: 'instagram', username: instagramUsername },
      { platform: 'facebook', username: facebookUsername },
      { platform: 'tiktok', username: tiktokUsername },
      { platform: 'snapchat', username: snapchatUsername },
      { platform: 'x', username: xUsername },
      { platform: 'whatsapp', username: whatsappNumber },
    ]

    return allHandles
      .filter(
        (handle): handle is { platform: SocialPlatform; username: string } => {
          // Filter out null, undefined, and empty strings
          return handle.username != null && handle.username.trim() !== ''
        },
      )
      .map(({ platform, username }) => ({
        platform,
        username: username.trim(),
      }))
  }, [
    instagramUsername,
    facebookUsername,
    tiktokUsername,
    snapchatUsername,
    xUsername,
    whatsappNumber,
  ])

  // If no handles, don't render anything
  if (handles.length === 0) {
    return null
  }

  return (
    <div className={cn('w-full', className)}>
      {/* Grid layout - responsive and RTL-aware */}
      <div
        className={cn(
          'grid gap-3',
          // Mobile: 1 column, Tablet+: 2 columns
          'grid-cols-1 sm:grid-cols-2',
          // Ensure proper alignment in RTL
          'items-start',
        )}
      >
        {handles.map(({ platform, username }) => (
          <SocialMediaLink
            key={platform}
            platform={platform}
            username={username}
            size={size}
            variant={variant}
          />
        ))}
      </div>
    </div>
  )
}
