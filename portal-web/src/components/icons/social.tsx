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
 * import { SocialIcon, SOCIAL_COLORS } from "@/components/icons/social";
 * <SocialIcon platform="instagram" className="w-5 h-5" />
 * ```
 */

import {
  SiInstagram,
  SiFacebook,
  SiTiktok,
  SiSnapchat,
  SiX,
} from "react-icons/si";
import { MessageCircle } from "lucide-react";
import { type ComponentType, type SVGProps } from "react";

export type SocialPlatform = "instagram" | "facebook" | "tiktok" | "snapchat" | "x" | "whatsapp";

export interface SocialIconProps extends SVGProps<SVGSVGElement> {
  platform: SocialPlatform;
}

// Map of platform to icon component
const SOCIAL_ICON_MAP: Record<SocialPlatform, ComponentType<SVGProps<SVGSVGElement>>> = {
  instagram: SiInstagram,
  facebook: SiFacebook,
  tiktok: SiTiktok,
  snapchat: SiSnapchat,
  x: SiX,
  whatsapp: MessageCircle, // Using lucide-react's MessageCircle as it's cleaner than WhatsApp logo
};

// Platform brand colors (official colors)
export const SOCIAL_COLORS: Record<SocialPlatform, string> = {
  instagram: "#E4405F",
  facebook: "#1877F2",
  tiktok: "#000000",
  snapchat: "#FFFC00",
  x: "#000000",
  whatsapp: "#25D366",
};

// Tailwind CSS classes for colors
export const SOCIAL_COLOR_CLASSES: Record<SocialPlatform, string> = {
  instagram: "text-[#E4405F]",
  facebook: "text-[#1877F2]",
  tiktok: "text-[#000000] dark:text-[#FFFFFF]",
  snapchat: "text-[#FFFC00]",
  x: "text-[#000000] dark:text-[#FFFFFF]",
  whatsapp: "text-[#25D366]",
};

// Tailwind CSS classes for hover colors
export const SOCIAL_HOVER_COLOR_CLASSES: Record<SocialPlatform, string> = {
  instagram: "hover:text-[#C13584]",
  facebook: "hover:text-[#0C63D4]",
  tiktok: "hover:text-[#EE1D52]",
  snapchat: "hover:text-[#FFFC00]",
  x: "hover:text-[#1DA1F2]",
  whatsapp: "hover:text-[#128C7E]",
};

// Background colors
export const SOCIAL_BG_CLASSES: Record<SocialPlatform, string> = {
  instagram: "bg-[#E4405F]/10",
  facebook: "bg-[#1877F2]/10",
  tiktok: "bg-[#000000]/10 dark:bg-[#FFFFFF]/10",
  snapchat: "bg-[#FFFC00]/20",
  x: "bg-[#000000]/10 dark:bg-[#FFFFFF]/10",
  whatsapp: "bg-[#25D366]/10",
};

// Hover background colors
export const SOCIAL_HOVER_BG_CLASSES: Record<SocialPlatform, string> = {
  instagram: "hover:bg-[#E4405F]/20",
  facebook: "hover:bg-[#1877F2]/20",
  tiktok: "hover:bg-[#EE1D52]/20",
  snapchat: "hover:bg-[#FFFC00]/30",
  x: "hover:bg-[#1DA1F2]/20",
  whatsapp: "hover:bg-[#25D366]/20",
};

// Platform labels for accessibility and display
export const SOCIAL_LABELS: Record<SocialPlatform, string> = {
  instagram: "Instagram",
  facebook: "Facebook",
  tiktok: "TikTok",
  snapchat: "Snapchat",
  x: "X (Twitter)",
  whatsapp: "WhatsApp",
};

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
export function SocialIcon({ platform, className, ...props }: SocialIconProps) {
  const Icon = SOCIAL_ICON_MAP[platform];
  return <Icon className={className} {...props} />;
}
