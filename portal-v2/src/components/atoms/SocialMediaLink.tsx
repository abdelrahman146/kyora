/**
 * SocialMediaLink Component
 *
 * A production-grade, mobile-first social media link button with proper RTL/LTR support.
 * Displays platform icon, username, and links to the user's profile.
 *
 * Features:
 * - Mobile-first design with large touch targets (min 44px)
 * - RTL/LTR support using logical properties
 * - Platform-specific brand colors and icons
 * - Accessible with proper ARIA labels
 * - Smooth hover/active states
 * - Opens in new tab with security attributes
 *
 * @example
 * ```tsx
 * <SocialMediaLink platform="instagram" username="john_doe" />
 * <SocialMediaLink platform="whatsapp" username="+971501234567" />
 * ```
 */

import { useMemo } from 'react'
import { cn } from '../../lib/utils'
import { SocialIcon  } from '../icons/social'
import {
  SOCIAL_BG_CLASSES,
  SOCIAL_COLOR_CLASSES,
  SOCIAL_HOVER_BG_CLASSES,
  SOCIAL_HOVER_COLOR_CLASSES,
  SOCIAL_LABELS,
} from '../icons/socialConstants'
import type {SocialPlatform} from '../icons/social';
import type { AnchorHTMLAttributes } from 'react'

export interface SocialMediaLinkProps
  extends Omit<AnchorHTMLAttributes<HTMLAnchorElement>, 'href'> {
  platform: SocialPlatform
  username: string
  size?: 'sm' | 'md' | 'lg'
  variant?: 'default' | 'minimal'
}

interface PlatformConfig {
  label: string
  getUrl: (username: string) => string
}

const platformConfigs: Record<SocialPlatform, PlatformConfig> = {
  instagram: {
    label: SOCIAL_LABELS.instagram,
    getUrl: (username) =>
      `https://instagram.com/${username.replace('@', '')}`,
  },
  facebook: {
    label: SOCIAL_LABELS.facebook,
    getUrl: (username) => `https://facebook.com/${username.replace('@', '')}`,
  },
  tiktok: {
    label: SOCIAL_LABELS.tiktok,
    getUrl: (username) => `https://tiktok.com/@${username.replace('@', '')}`,
  },
  snapchat: {
    label: SOCIAL_LABELS.snapchat,
    getUrl: (username) =>
      `https://snapchat.com/add/${username.replace('@', '')}`,
  },
  x: {
    label: SOCIAL_LABELS.x,
    getUrl: (username) => `https://x.com/${username.replace('@', '')}`,
  },
  whatsapp: {
    label: SOCIAL_LABELS.whatsapp,
    getUrl: (username) => {
      // Remove all non-digit characters for WhatsApp
      const cleanNumber = username.replace(/\D/g, '')
      return `https://wa.me/${cleanNumber}`
    },
  },
}

export function SocialMediaLink({
  platform,
  username,
  size = 'md',
  variant = 'default',
  className,
  ...props
}: SocialMediaLinkProps) {
  const config = platformConfigs[platform]
  const url = useMemo(() => config.getUrl(username), [config, username])

  const sizeClasses = {
    sm: {
      container: 'min-h-[44px] px-3 py-2 gap-2',
      icon: 'w-4 h-4',
      text: 'text-sm',
    },
    md: {
      container: 'min-h-[48px] px-4 py-2.5 gap-2.5',
      icon: 'w-5 h-5',
      text: 'text-base',
    },
    lg: {
      container: 'min-h-[52px] px-5 py-3 gap-3',
      icon: 'w-6 h-6',
      text: 'text-lg',
    },
  }

  const variantClasses =
    variant === 'minimal'
      ? cn(
          // Minimal variant - transparent bg, only icon color
          'bg-transparent',
          SOCIAL_COLOR_CLASSES[platform],
          SOCIAL_HOVER_COLOR_CLASSES[platform],
          'hover:bg-base-200'
        )
      : cn(
          // Default variant - colored bg
          SOCIAL_BG_CLASSES[platform],
          SOCIAL_HOVER_BG_CLASSES[platform],
          SOCIAL_COLOR_CLASSES[platform],
          SOCIAL_HOVER_COLOR_CLASSES[platform]
        )

  return (
    <a
      href={url}
      target="_blank"
      rel="noopener noreferrer"
      aria-label={`${config.label}: ${username}`}
      dir="ltr"
      className={cn(
        // Base styles
        'inline-flex items-center justify-start',
        'rounded-lg',
        'font-medium',
        'transition-all duration-200',
        'focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2',
        'active:scale-[0.98]',
        // Size
        sizeClasses[size].container,
        // Variant
        variantClasses,
        // Custom
        className
      )}
      {...props}
    >
      {/* Icon */}
      <SocialIcon
        platform={platform}
        className={cn(
          'shrink-0',
          sizeClasses[size].icon,
          'transition-transform duration-200 group-hover:scale-110'
        )}
        aria-hidden="true"
      />

      {/* Username - always displays LTR since social handles are always in English */}
      <span className={cn('truncate font-medium', sizeClasses[size].text)}>
        {username.startsWith('@') ? username : `@${username}`}
      </span>
    </a>
  )
}
