/**
 * SocialMediaInputs Component
 *
 * An expandable section for entering social media handles in customer forms.
 * Features progressive disclosure to avoid overwhelming users with many optional fields.
 *
 * Features:
 * - Collapsible section (collapsed by default in create, expanded when data exists in edit)
 * - Shows count of filled fields in header
 * - Mobile-first responsive grid (1-2 columns)
 * - RTL/LTR support with proper icon alignment
 * - Platform-specific icons and placeholders
 * - All fields optional
 * - Clean validation and error display
 *
 * @example
 * ```tsx
 * <SocialMediaInputs
 *   instagramUsername={field.value}
 *   onInstagramChange={field.onChange}
 *   errors={{ instagramUsername: "Invalid username" }}
 *   disabled={isSubmitting}
 * />
 * ```
 */

import { useState, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { ChevronDown, ChevronUp } from "lucide-react";
import { cn } from "@/lib/utils";
import { FormInput } from "@/components/atoms/FormInput";
import { SocialIcon, SOCIAL_COLOR_CLASSES } from "@/components/icons/social";

export interface SocialMediaInputsProps {
  // Instagram
  instagramUsername?: string;
  onInstagramChange?: (value: string) => void;
  instagramError?: string;

  // Facebook
  facebookUsername?: string;
  onFacebookChange?: (value: string) => void;
  facebookError?: string;

  // TikTok
  tiktokUsername?: string;
  onTiktokChange?: (value: string) => void;
  tiktokError?: string;

  // Snapchat
  snapchatUsername?: string;
  onSnapchatChange?: (value: string) => void;
  snapchatError?: string;

  // X (Twitter)
  xUsername?: string;
  onXChange?: (value: string) => void;
  xError?: string;

  // WhatsApp
  whatsappNumber?: string;
  onWhatsappChange?: (value: string) => void;
  whatsappError?: string;

  // State
  disabled?: boolean;
  defaultExpanded?: boolean;
}

export function SocialMediaInputs({
  instagramUsername = "",
  onInstagramChange,
  instagramError,
  facebookUsername = "",
  onFacebookChange,
  facebookError,
  tiktokUsername = "",
  onTiktokChange,
  tiktokError,
  snapchatUsername = "",
  onSnapchatChange,
  snapchatError,
  xUsername = "",
  onXChange,
  xError,
  whatsappNumber = "",
  onWhatsappChange,
  whatsappError,
  disabled = false,
  defaultExpanded = false,
}: SocialMediaInputsProps) {
  const { t } = useTranslation();

  // Count filled fields for summary
  const filledCount = useMemo(() => {
    let count = 0;
    if (instagramUsername.trim()) count++;
    if (facebookUsername.trim()) count++;
    if (tiktokUsername.trim()) count++;
    if (snapchatUsername.trim()) count++;
    if (xUsername.trim()) count++;
    if (whatsappNumber.trim()) count++;
    return count;
  }, [
    instagramUsername,
    facebookUsername,
    tiktokUsername,
    snapchatUsername,
    xUsername,
    whatsappNumber,
  ]);

  // Auto-expand if data exists, otherwise respect defaultExpanded
  const [isExpanded, setIsExpanded] = useState(defaultExpanded || filledCount > 0);

  return (
    <div className="space-y-3">
      {/* Header - Collapsible */}
      <button
        type="button"
        onClick={() => {
          setIsExpanded(!isExpanded);
        }}
        className={cn(
          "w-full flex items-center justify-between",
          "px-4 py-3 rounded-lg",
          "bg-base-200/50 hover:bg-base-200",
          "transition-colors duration-200",
          "focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2",
          disabled && "opacity-60 cursor-not-allowed"
        )}
        disabled={disabled}
        aria-expanded={isExpanded}
        aria-controls="social-media-inputs"
      >
        <div className="flex items-center gap-3">
          <span className="font-medium text-base-content">
            {t("customers.form.social_media_section")}
          </span>
          {filledCount > 0 && (
            <span className="badge badge-primary badge-sm">
              {filledCount}
            </span>
          )}
        </div>
        {isExpanded ? (
          <ChevronUp size={20} className="text-base-content/60" />
        ) : (
          <ChevronDown size={20} className="text-base-content/60" />
        )}
      </button>

      {/* Content - Expandable */}
      {isExpanded && (
        <div
          id="social-media-inputs"
          className={cn(
            "grid gap-3",
            "grid-cols-1 sm:grid-cols-2",
            "animate-in fade-in slide-in-from-top-2 duration-200"
          )}
        >
          {/* Instagram */}
          <FormInput
            label={t("customers.form.instagram")}
            placeholder={t("customers.form.instagram_placeholder")}
            value={instagramUsername}
            onChange={(e) => {
              onInstagramChange?.(e.target.value);
            }}
            error={instagramError}
            disabled={disabled}
            autoComplete="off"
            dir="ltr"
            startIcon={
              <SocialIcon 
                platform="instagram" 
                className={cn("w-4 h-4", SOCIAL_COLOR_CLASSES.instagram)} 
              />
            }
          />

          {/* Facebook */}
          <FormInput
            label={t("customers.form.facebook")}
            placeholder={t("customers.form.facebook_placeholder")}
            value={facebookUsername}
            onChange={(e) => {
              onFacebookChange?.(e.target.value);
            }}
            error={facebookError}
            disabled={disabled}
            autoComplete="off"
            dir="ltr"
            startIcon={
              <SocialIcon 
                platform="facebook" 
                className={cn("w-4 h-4", SOCIAL_COLOR_CLASSES.facebook)} 
              />
            }
          />

          {/* TikTok */}
          <FormInput
            label={t("customers.form.tiktok")}
            placeholder={t("customers.form.tiktok_placeholder")}
            value={tiktokUsername}
            onChange={(e) => {
              onTiktokChange?.(e.target.value);
            }}
            error={tiktokError}
            disabled={disabled}
            autoComplete="off"
            dir="ltr"
            startIcon={
              <SocialIcon 
                platform="tiktok" 
                className="w-4 h-4"
              />
            }
          />

          {/* Snapchat */}
          <FormInput
            label={t("customers.form.snapchat")}
            placeholder={t("customers.form.snapchat_placeholder")}
            value={snapchatUsername}
            onChange={(e) => {
              onSnapchatChange?.(e.target.value);
            }}
            error={snapchatError}
            disabled={disabled}
            autoComplete="off"
            dir="ltr"
            startIcon={
              <SocialIcon 
                platform="snapchat" 
                className={cn("w-4 h-4", SOCIAL_COLOR_CLASSES.snapchat)} 
              />
            }
          />

          {/* X (Twitter) */}
          <FormInput
            label={t("customers.form.x")}
            placeholder={t("customers.form.x_placeholder")}
            value={xUsername}
            onChange={(e) => {
              onXChange?.(e.target.value);
            }}
            error={xError}
            disabled={disabled}
            autoComplete="off"
            dir="ltr"
            startIcon={
              <SocialIcon 
                platform="x" 
                className="w-4 h-4"
              />
            }
          />

          {/* WhatsApp */}
          <FormInput
            label={t("customers.form.whatsapp")}
            placeholder={t("customers.form.whatsapp_placeholder")}
            value={whatsappNumber}
            onChange={(e) => {
              onWhatsappChange?.(e.target.value);
            }}
            error={whatsappError}
            disabled={disabled}
            autoComplete="off"
            inputMode="tel"
            dir="ltr"
            startIcon={
              <SocialIcon 
                platform="whatsapp" 
                className={cn("w-4 h-4", SOCIAL_COLOR_CLASSES.whatsapp)} 
              />
            }
          />
        </div>
      )}

      {/* Helper text when collapsed */}
      {!isExpanded && filledCount === 0 && (
        <p className="text-sm text-base-content/60 px-4">
          {t("customers.form.social_media_hint")}
        </p>
      )}
    </div>
  );
}
