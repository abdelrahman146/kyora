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

import { type AnchorHTMLAttributes, useMemo } from "react";
import { cn } from "@/lib/utils";
import { MessageCircle, type LucideIcon } from "lucide-react";

// Instagram icon (custom SVG - lucide brand icons deprecated)
const InstagramIcon = ({ className }: { className?: string }) => (
  <svg
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
    className={className}
    xmlns="http://www.w3.org/2000/svg"
  >
    <rect width="20" height="20" x="2" y="2" rx="5" ry="5" />
    <path d="M16 11.37A4 4 0 1 1 12.63 8 4 4 0 0 1 16 11.37z" />
    <line x1="17.5" x2="17.51" y1="6.5" y2="6.5" />
  </svg>
);

// Facebook icon (custom SVG - lucide brand icons deprecated)
const FacebookIcon = ({ className }: { className?: string }) => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
    className={className}
    xmlns="http://www.w3.org/2000/svg"
  >
    <path d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z" />
  </svg>
);

// X/Twitter icon (custom SVG - lucide brand icons deprecated)
const XIcon = ({ className }: { className?: string }) => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
    className={className}
    xmlns="http://www.w3.org/2000/svg"
  >
    <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z" />
  </svg>
);

// TikTok icon (not in lucide-react, using custom SVG)
const TikTokIcon = ({ className }: { className?: string }) => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
    className={className}
    xmlns="http://www.w3.org/2000/svg"
  >
    <path d="M19.59 6.69a4.83 4.83 0 0 1-3.77-4.25V2h-3.45v13.67a2.89 2.89 0 0 1-5.2 1.74 2.89 2.89 0 0 1 2.31-4.64 2.93 2.93 0 0 1 .88.13V9.4a6.84 6.84 0 0 0-1-.05A6.33 6.33 0 0 0 5 20.1a6.34 6.34 0 0 0 10.86-4.43v-7a8.16 8.16 0 0 0 4.77 1.52v-3.4a4.85 4.85 0 0 1-1-.1z" />
  </svg>
);

// Snapchat icon (not in lucide-react, using custom SVG)
const SnapchatIcon = ({ className }: { className?: string }) => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
    className={className}
    xmlns="http://www.w3.org/2000/svg"
  >
    <path d="M12.206.793c.99 0 4.347.276 5.93 3.821.529 1.193.403 3.219.299 4.847l-.003.06c-.012.18-.022.345-.03.51.075.045.203.09.401.09.3-.016.659-.12 1.033-.301.165-.088.344-.104.464-.104.182 0 .359.029.509.09.45.149.734.479.734.838.015.449-.39.839-1.213 1.168-.089.029-.209.075-.344.119-.45.135-1.139.36-1.333.81-.09.224-.061.524.12.868l.015.015c.06.136 1.526 3.475 4.791 4.014.255.044.435.27.42.509-.014.18-.134.389-.344.524-.569.345-1.553.434-2.378.434-.53 0-1.018-.075-1.318-.149-.195-.044-.375-.074-.509-.074-.09 0-.149.014-.209.074-.134.119-.119.449-.119.524.016.405-.014 1.021-.239 1.544-.36.853-1.122 1.528-2.068 1.873-.225.091-.539.21-1.018.21-.12 0-.255-.015-.375-.045-1.952-.434-2.79-1.648-3.373-2.691-.195-.345-.39-.705-.614-1.006-.464-.629-1.228-1.066-2.134-1.215-.195-.03-.359-.061-.509-.091-.404-.104-.824-.21-1.169-.21-.314 0-.627.091-.912.21l-.015.015c-.195.119-.404.239-.584.239-.209 0-.405-.134-.524-.27-.195-.24-.269-.57-.254-.899.029-.42.27-.779.614-.899 1.17-.389 2.531-1.515 3.306-2.544.195-.27.359-.539.479-.838.09-.241.119-.479.119-.689 0-.404-.165-.749-.479-1.006-.195-.165-.404-.3-.614-.419l-.015-.015c-.569-.315-1.049-.644-1.439-1.035-.524-.509-.824-1.095-.824-1.693 0-.449.24-.839.644-1.006.195-.075.42-.119.644-.119.405 0 .883.12 1.438.359.404.18.823.39 1.303.39.3 0 .569-.045.824-.12.194-.074.389-.149.584-.224l.06-.03c.405-.18.75-.329 1.169-.329.434 0 .869.15 1.243.42.509.359.779.914.779 1.559 0 .195-.015.389-.06.584-.045.195-.09.375-.09.554 0 .434.255.823.644 1.036.404.21.869.24 1.273.24.314 0 .614-.03.854-.09.27-.045.495-.09.645-.09.194 0 .344.03.479.09.404.149.689.524.689.974 0 .404-.195.794-.569 1.065-.225.165-.495.3-.795.435-.255.12-.539.21-.824.3-.404.12-.853.21-1.273.21-.27 0-.524-.03-.779-.09-.404-.09-.854-.195-1.348-.195-.644 0-1.303.165-1.952.479-.524.254-1.048.614-1.438 1.065-.404.464-.659 1.006-.659 1.575 0 .584.165 1.215.614 1.828.569.779 1.528 1.438 2.827 1.873.885.3 1.903.42 2.84.42.45 0 .868-.045 1.228-.119.404-.09.779-.195 1.093-.344.195-.091.405-.195.614-.299.479-.24.869-.45 1.318-.45.195 0 .374.045.524.12.404.195.644.614.644 1.094 0 .629-.344 1.288-.869 1.693-.404.314-.929.524-1.498.644-.254.045-.524.075-.794.075-.689 0-1.378-.195-2.009-.524-.404-.21-.779-.479-1.063-.779-.225-.24-.405-.524-.539-.824-.06-.135-.105-.27-.15-.405-.044-.135-.089-.27-.149-.405-.045-.09-.09-.18-.149-.27-.135-.18-.254-.345-.404-.495-.149-.135-.344-.255-.569-.344-.27-.12-.584-.18-.884-.18-.345 0-.704.075-1.048.195-.3.105-.599.24-.854.405-.12.075-.255.165-.375.255-.3.24-.614.479-.989.689-.51.3-1.108.479-1.723.479-.405 0-.824-.074-1.213-.225-.854-.314-1.558-.959-1.967-1.828-.224-.479-.329-1.005-.329-1.529 0-.254.029-.509.074-.748.03-.18.075-.345.12-.509.029-.105.06-.225.104-.345.105-.314.195-.629.195-.929 0-.18-.03-.345-.09-.494-.074-.18-.194-.345-.344-.495-.254-.255-.614-.404-1.003-.479l-.06-.015c-.314-.075-.644-.149-.989-.149-.569 0-1.168.12-1.768.345-.374.15-.734.329-1.048.554-.18.12-.344.27-.479.419-.299.315-.479.704-.479 1.109 0 .719.359 1.438.989 2.061.405.404.929.734 1.498 1.034l.06.029c.194.105.404.225.584.375.27.225.404.524.404.854 0 .254-.045.524-.149.779-.12.285-.284.569-.479.839-.794 1.034-2.151 2.16-3.321 2.549-.344.119-.584.479-.614.899-.015.329.06.659.254.899.119.136.315.27.524.27.18 0 .389-.12.584-.239l.015-.015c.285-.119.598-.21.912-.21.345 0 .765.106 1.169.21.15.03.314.061.509.091.906.149 1.67.586 2.134 1.215.224.301.419.661.614 1.006.584 1.043 1.421 2.257 3.373 2.691.12.03.255.045.375.045.479 0 .793-.119 1.018-.21.946-.345 1.708-1.02 2.068-1.873.225-.523.255-1.139.239-1.544 0-.075-.015-.405.119-.524.06-.06.119-.074.209-.074.134 0 .314.03.509.074.3.074.788.149 1.318.149.825 0 1.809-.089 2.378-.434.21-.135.33-.344.344-.524.015-.239-.165-.465-.42-.509-3.265-.539-4.731-3.878-4.791-4.014l-.015-.015c-.181-.344-.21-.644-.12-.868.194-.45.883-.675 1.333-.81.135-.044.255-.09.344-.119.823-.329 1.228-.719 1.213-1.168 0-.359-.284-.689-.734-.838-.15-.061-.327-.09-.509-.09-.12 0-.299.016-.464.104-.374.181-.733.285-1.033.301-.198 0-.326-.045-.401-.09.008-.165.018-.33.03-.51l.003-.06c.104-1.628.23-3.654-.299-4.847-1.583-3.545-4.94-3.821-5.93-3.821z" />
  </svg>
);

export type SocialPlatform = "instagram" | "facebook" | "tiktok" | "snapchat" | "x" | "whatsapp";

export interface SocialMediaLinkProps extends Omit<AnchorHTMLAttributes<HTMLAnchorElement>, "href"> {
  platform: SocialPlatform;
  username: string;
  size?: "sm" | "md" | "lg";
  variant?: "default" | "minimal";
}

interface PlatformConfig {
  icon: LucideIcon | typeof TikTokIcon | typeof SnapchatIcon | typeof InstagramIcon | typeof FacebookIcon | typeof XIcon;
  label: string;
  color: string;
  hoverColor: string;
  bgColor: string;
  hoverBgColor: string;
  getUrl: (username: string) => string;
}

const platformConfigs: Record<SocialPlatform, PlatformConfig> = {
  instagram: {
    icon: InstagramIcon,
    label: "Instagram",
    color: "text-[#E4405F]",
    hoverColor: "hover:text-[#C13584]",
    bgColor: "bg-[#E4405F]/10",
    hoverBgColor: "hover:bg-[#E4405F]/20",
    getUrl: (username) => `https://instagram.com/${username.replace("@", "")}`,
  },
  facebook: {
    icon: FacebookIcon,
    label: "Facebook",
    color: "text-[#1877F2]",
    hoverColor: "hover:text-[#0C63D4]",
    bgColor: "bg-[#1877F2]/10",
    hoverBgColor: "hover:bg-[#1877F2]/20",
    getUrl: (username) => `https://facebook.com/${username.replace("@", "")}`,
  },
  tiktok: {
    icon: TikTokIcon,
    label: "TikTok",
    color: "text-[#000000] dark:text-[#FFFFFF]",
    hoverColor: "hover:text-[#EE1D52]",
    bgColor: "bg-[#000000]/10 dark:bg-[#FFFFFF]/10",
    hoverBgColor: "hover:bg-[#EE1D52]/20",
    getUrl: (username) => `https://tiktok.com/@${username.replace("@", "")}`,
  },
  snapchat: {
    icon: SnapchatIcon,
    label: "Snapchat",
    color: "text-[#FFFC00]",
    hoverColor: "hover:text-[#FFFC00]",
    bgColor: "bg-[#FFFC00]/20",
    hoverBgColor: "hover:bg-[#FFFC00]/30",
    getUrl: (username) => `https://snapchat.com/add/${username.replace("@", "")}`,
  },
  x: {
    icon: XIcon,
    label: "X (Twitter)",
    color: "text-[#000000] dark:text-[#FFFFFF]",
    hoverColor: "hover:text-[#1DA1F2]",
    bgColor: "bg-[#000000]/10 dark:bg-[#FFFFFF]/10",
    hoverBgColor: "hover:bg-[#1DA1F2]/20",
    getUrl: (username) => `https://x.com/${username.replace("@", "")}`,
  },
  whatsapp: {
    icon: MessageCircle,
    label: "WhatsApp",
    color: "text-[#25D366]",
    hoverColor: "hover:text-[#128C7E]",
    bgColor: "bg-[#25D366]/10",
    hoverBgColor: "hover:bg-[#25D366]/20",
    getUrl: (username) => {
      // Remove all non-digit characters for WhatsApp
      const cleanNumber = username.replace(/\D/g, "");
      return `https://wa.me/${cleanNumber}`;
    },
  },
};

export function SocialMediaLink({
  platform,
  username,
  size = "md",
  variant = "default",
  className,
  ...props
}: SocialMediaLinkProps) {
  const config = platformConfigs[platform];
  const Icon = config.icon;

  const url = useMemo(() => config.getUrl(username), [config, username]);

  const sizeClasses = {
    sm: {
      container: "min-h-[44px] px-3 py-2 gap-2",
      icon: "w-4 h-4",
      text: "text-sm",
    },
    md: {
      container: "min-h-[48px] px-4 py-2.5 gap-2.5",
      icon: "w-5 h-5",
      text: "text-base",
    },
    lg: {
      container: "min-h-[52px] px-5 py-3 gap-3",
      icon: "w-6 h-6",
      text: "text-lg",
    },
  };

  const variantClasses =
    variant === "minimal"
      ? cn(
          // Minimal variant - transparent bg, only icon color
          "bg-transparent",
          config.color,
          config.hoverColor,
          "hover:bg-base-200"
        )
      : cn(
          // Default variant - colored bg
          config.bgColor,
          config.hoverBgColor,
          config.color,
          config.hoverColor
        );

  return (
    <a
      href={url}
      target="_blank"
      rel="noopener noreferrer"
      aria-label={`${config.label}: ${username}`}
      className={cn(
        // Base styles
        "inline-flex items-center justify-start",
        "rounded-lg",
        "font-medium",
        "transition-all duration-200",
        "focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2",
        "active:scale-[0.98]",
        // RTL/LTR support - use logical properties
        "text-start",
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
      <Icon
        className={cn(
          "shrink-0",
          sizeClasses[size].icon,
          "transition-transform duration-200 group-hover:scale-110"
        )}
        aria-hidden="true"
      />

      {/* Username - with proper LTR direction for social handles */}
      <span
        className={cn(
          "truncate font-medium",
          sizeClasses[size].text,
          // Force LTR for usernames/handles even in RTL layout
          "direction-ltr"
        )}
        dir="ltr"
      >
        {username.startsWith("@") ? username : `@${username}`}
      </span>
    </a>
  );
}
