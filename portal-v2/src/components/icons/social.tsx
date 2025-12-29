/**
 * Social Media Icons
 *
 * Centralized social media platform icons using react-icons/si (Simple Icons).
 * These icons are reused across the application (SocialMediaLink, SocialMediaInputs, etc.)
 *
 * Features:
 * - Official brand icons from Simple Icons
 * - Consistent sizing and styling
 * - Type-safe platform enumeration
 * - DRY principle - single source of truth
 *
 * @example
 * ```tsx
 * import { SocialIcon } from "@/components/icons/social";
 * <SocialIcon platform="instagram" className="w-5 h-5" />
 * ```
 */

import {
  SiFacebook,
  SiInstagram,
  SiSnapchat,
  SiTiktok,
  SiX,
} from 'react-icons/si'
import { MessageCircle } from 'lucide-react'
import type { ComponentType, SVGProps } from 'react'
import type { SocialPlatform } from './socialConstants'

// Re-export the type for convenience
export type { SocialPlatform }

export interface SocialIconProps extends SVGProps<SVGSVGElement> {
  platform: SocialPlatform
}

// Map of platform to icon component
const SOCIAL_ICON_MAP: Record<
  SocialPlatform,
  ComponentType<SVGProps<SVGSVGElement>>
> = {
  instagram: SiInstagram,
  facebook: SiFacebook,
  tiktok: SiTiktok,
  snapchat: SiSnapchat,
  x: SiX,
  whatsapp: MessageCircle, // Using lucide-react's MessageCircle as it's cleaner than WhatsApp logo
}

/**
 * SocialIcon Component
 *
 * Renders the appropriate icon for a given social media platform.
 * All icons from Simple Icons maintain consistent sizing and styling.
 *
 * @example
 * ```tsx
 * <SocialIcon platform="instagram" className="w-5 h-5 text-[#E4405F]" />
 * ```
 */
export function SocialIcon({
  platform,
  className,
  ...props
}: SocialIconProps) {
  const Icon = SOCIAL_ICON_MAP[platform]
  return <Icon className={className} {...props} />
}
