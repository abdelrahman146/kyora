/**
 * SocialMediaHandles Component
 *
 * Elegant, minimal social media display with clean icon-based design.
 * Mobile-first, RTL-ready, accessible.
 *
 * Features:
 * - Clean icon circles with platform colors
 * - Minimal text, maximum clarity
 * - Hover states with subtle animations
 * - Mobile-optimized spacing
 * - RTL/LTR support
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
import {
  FaFacebook,
  FaInstagram,
  FaSnapchat,
  FaTiktok,
  FaWhatsapp,
  FaXTwitter,
} from 'react-icons/fa6'
import { cn } from '../../lib/utils'
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
  icon: React.ComponentType<{ className?: string }>
  label: string
  url: string
  colorClass: string
  hoverClass: string
}

const platformConfigs = {
  instagram: {
    icon: FaInstagram,
    label: 'Instagram',
    getUrl: (username: string) =>
      `https://instagram.com/${username.replace('@', '')}`,
    colorClass: 'bg-[#E4405F]/10 text-[#E4405F]',
    hoverClass: 'hover:bg-[#E4405F]/20',
  },
  facebook: {
    icon: FaFacebook,
    label: 'Facebook',
    getUrl: (username: string) =>
      `https://facebook.com/${username.replace('@', '')}`,
    colorClass: 'bg-[#1877F2]/10 text-[#1877F2]',
    hoverClass: 'hover:bg-[#1877F2]/20',
  },
  tiktok: {
    icon: FaTiktok,
    label: 'TikTok',
    getUrl: (username: string) =>
      `https://tiktok.com/@${username.replace('@', '')}`,
    colorClass: 'bg-base-content/10 text-base-content',
    hoverClass: 'hover:bg-base-content/20',
  },
  snapchat: {
    icon: FaSnapchat,
    label: 'Snapchat',
    getUrl: (username: string) =>
      `https://snapchat.com/add/${username.replace('@', '')}`,
    colorClass: 'bg-[#FFFC00]/20 text-base-content',
    hoverClass: 'hover:bg-[#FFFC00]/30',
  },
  x: {
    icon: FaXTwitter,
    label: 'X',
    getUrl: (username: string) => `https://x.com/${username.replace('@', '')}`,
    colorClass: 'bg-base-content/10 text-base-content',
    hoverClass: 'hover:bg-base-content/20',
  },
  whatsapp: {
    icon: FaWhatsapp,
    label: 'WhatsApp',
    getUrl: (username: string) => {
      const cleanNumber = username.replace(/\D/g, '')
      return `https://wa.me/${cleanNumber}`
    },
    colorClass: 'bg-[#25D366]/10 text-[#25D366]',
    hoverClass: 'hover:bg-[#25D366]/20',
  },
}

export function SocialMediaHandles({
  instagramUsername,
  facebookUsername,
  tiktokUsername,
  snapchatUsername,
  xUsername,
  whatsappNumber,
  size = 'md',
  variant: _variant = 'default',
  className,
}: SocialMediaHandlesProps) {
  // Filter and map social handles
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
        (handle): handle is { platform: SocialPlatform; username: string } =>
          handle.username != null && handle.username.trim() !== '',
      )
      .map(({ platform, username }) => {
        const config = platformConfigs[platform]
        return {
          platform,
          username: username.trim(),
          icon: config.icon,
          label: config.label,
          url: config.getUrl(username.trim()),
          colorClass: config.colorClass,
          hoverClass: config.hoverClass,
        }
      })
  }, [
    instagramUsername,
    facebookUsername,
    tiktokUsername,
    snapchatUsername,
    xUsername,
    whatsappNumber,
  ])

  if (handles.length === 0) {
    return null
  }

  const sizeClasses = {
    sm: 'size-10',
    md: 'size-12',
    lg: 'size-14',
  }

  const iconSizeClasses = {
    sm: 'text-lg',
    md: 'text-xl',
    lg: 'text-2xl',
  }

  return (
    <div className={cn('w-full', className)}>
      <div className="flex flex-wrap gap-3">
        {handles.map(
          ({
            platform,
            username,
            icon: Icon,
            label,
            url,
            colorClass,
            hoverClass,
          }) => (
            <a
              key={platform}
              href={url}
              target="_blank"
              rel="noopener noreferrer"
              className={cn(
                'group flex items-center gap-3 px-4 py-2.5 rounded-lg border border-base-300',
                'transition-all duration-200',
                'hover:border-base-content/20',
                'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2',
              )}
              aria-label={`${label}: ${username}`}
            >
              <div
                className={cn(
                  sizeClasses[size],
                  'rounded-full flex items-center justify-center flex-shrink-0',
                  'transition-colors duration-200',
                  colorClass,
                  hoverClass,
                )}
              >
                <Icon className={iconSizeClasses[size]} />
              </div>
              <div className="flex-1 min-w-0">
                <div className="text-xs text-base-content/60 leading-none mb-1">
                  {label}
                </div>
                <div className="font-medium text-sm truncate" dir="ltr">
                  {username}
                </div>
              </div>
            </a>
          ),
        )}
      </div>
    </div>
  )
}
