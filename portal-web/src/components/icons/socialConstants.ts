/**
 * Social Media Constants
 *
 * Centralized constants for social media platforms including colors, labels, and type definitions.
 * Separated from the icon components to comply with React Fast Refresh requirements.
 *
 * @example
 * ```tsx
 * import { SOCIAL_COLORS, SOCIAL_LABELS } from "@/components/icons/socialConstants";
 * const color = SOCIAL_COLORS["instagram"];
 * const label = SOCIAL_LABELS["instagram"];
 * ```
 */

export type SocialPlatform =
  | 'instagram'
  | 'facebook'
  | 'tiktok'
  | 'snapchat'
  | 'x'
  | 'whatsapp'

// Platform brand colors (official colors)
export const SOCIAL_COLORS: Record<SocialPlatform, string> = {
  instagram: '#E4405F',
  facebook: '#1877F2',
  tiktok: '#000000',
  snapchat: '#FFFC00',
  x: '#000000',
  whatsapp: '#25D366',
}

// Tailwind CSS classes for colors
export const SOCIAL_COLOR_CLASSES: Record<SocialPlatform, string> = {
  instagram: 'text-[#E4405F]',
  facebook: 'text-[#1877F2]',
  tiktok: 'text-[#000000] dark:text-[#FFFFFF]',
  snapchat: 'text-[#FFFC00]',
  x: 'text-[#000000] dark:text-[#FFFFFF]',
  whatsapp: 'text-[#25D366]',
}

// Tailwind CSS classes for hover colors
export const SOCIAL_HOVER_COLOR_CLASSES: Record<SocialPlatform, string> = {
  instagram: 'hover:text-[#C13584]',
  facebook: 'hover:text-[#0C63D4]',
  tiktok: 'hover:text-[#EE1D52]',
  snapchat: 'hover:text-[#FFFC00]',
  x: 'hover:text-[#1DA1F2]',
  whatsapp: 'hover:text-[#128C7E]',
}

// Background colors
export const SOCIAL_BG_CLASSES: Record<SocialPlatform, string> = {
  instagram: 'bg-[#E4405F]/10',
  facebook: 'bg-[#1877F2]/10',
  tiktok: 'bg-[#000000]/10 dark:bg-[#FFFFFF]/10',
  snapchat: 'bg-[#FFFC00]/20',
  x: 'bg-[#000000]/10 dark:bg-[#FFFFFF]/10',
  whatsapp: 'bg-[#25D366]/10',
}

// Hover background colors
export const SOCIAL_HOVER_BG_CLASSES: Record<SocialPlatform, string> = {
  instagram: 'hover:bg-[#E4405F]/20',
  facebook: 'hover:bg-[#1877F2]/20',
  tiktok: 'hover:bg-[#EE1D52]/20',
  snapchat: 'hover:bg-[#FFFC00]/30',
  x: 'hover:bg-[#1DA1F2]/20',
  whatsapp: 'hover:bg-[#25D366]/20',
}

// Platform labels for accessibility and display
export const SOCIAL_LABELS: Record<SocialPlatform, string> = {
  instagram: 'Instagram',
  facebook: 'Facebook',
  tiktok: 'TikTok',
  snapchat: 'Snapchat',
  x: 'X (Twitter)',
  whatsapp: 'WhatsApp',
}
