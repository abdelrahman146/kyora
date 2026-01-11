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

import { useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { ChevronDown, ChevronUp } from 'lucide-react'
import { cn } from '@/lib/utils'
import { FormInput } from '@/components/form/FormInput'
import { SocialIcon } from '@/components/icons/social'
import { SOCIAL_COLOR_CLASSES } from '@/components/icons/socialConstants'

export interface SocialMediaInputsProps {
  // Instagram
  instagramUsername?: string
  onInstagramChange?: (value: string) => void
  instagramError?: string

  // Facebook
  facebookUsername?: string
  onFacebookChange?: (value: string) => void
  facebookError?: string

  // TikTok
  tiktokUsername?: string
  onTiktokChange?: (value: string) => void
  tiktokError?: string

  // Snapchat
  snapchatUsername?: string
  onSnapchatChange?: (value: string) => void
  snapchatError?: string

  // X (Twitter)
  xUsername?: string
  onXChange?: (value: string) => void
  xError?: string

  // WhatsApp
  whatsappNumber?: string
  onWhatsappChange?: (value: string) => void
  whatsappError?: string

  // State
  disabled?: boolean
  defaultExpanded?: boolean
}

export function SocialMediaInputs({
  instagramUsername = '',
  onInstagramChange,
  instagramError,
  facebookUsername = '',
  onFacebookChange,
  facebookError,
  tiktokUsername = '',
  onTiktokChange,
  tiktokError,
  snapchatUsername = '',
  onSnapchatChange,
  snapchatError,
  xUsername = '',
  onXChange,
  xError,
  whatsappNumber = '',
  onWhatsappChange,
  whatsappError,
  disabled = false,
  defaultExpanded = false,
}: SocialMediaInputsProps) {
  const { t: tCustomers } = useTranslation('customers')

  // Count filled fields for summary
  const filledCount = useMemo(() => {
    let count = 0
    if (instagramUsername.trim()) count++
    if (facebookUsername.trim()) count++
    if (tiktokUsername.trim()) count++
    if (snapchatUsername.trim()) count++
    if (xUsername.trim()) count++
    if (whatsappNumber.trim()) count++
    return count
  }, [
    instagramUsername,
    facebookUsername,
    tiktokUsername,
    snapchatUsername,
    xUsername,
    whatsappNumber,
  ])

  // Auto-expand if data exists, otherwise respect defaultExpanded
  const [isExpanded, setIsExpanded] = useState(
    defaultExpanded || filledCount > 0,
  )

  return (
    <div className="space-y-3">
      {/* Header - Collapsible */}
      <button
        type="button"
        onClick={() => {
          setIsExpanded(!isExpanded)
        }}
        className={cn(
          'w-full flex items-center justify-between',
          'px-4 py-3 rounded-lg',
          'bg-base-200/50 hover:bg-base-200',
          'transition-colors duration-200',
          'focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2',
          disabled && 'opacity-60 cursor-not-allowed',
        )}
        disabled={disabled}
        aria-expanded={isExpanded}
        aria-controls="social-media-inputs"
      >
        <div className="flex items-center gap-3">
          <span className="font-medium text-base-content">
            {tCustomers('form.social_media_section')}
          </span>
          {filledCount > 0 && (
            <span className="badge badge-primary badge-sm">{filledCount}</span>
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
            'grid gap-3',
            'grid-cols-1 sm:grid-cols-2',
            'animate-in fade-in slide-in-from-top-2 duration-200',
          )}
        >
          {/* Instagram */}
          <FormInput
            label={tCustomers('form.instagram')}
            placeholder={tCustomers('form.instagram_placeholder')}
            value={instagramUsername}
            onChange={(e) => {
              onInstagramChange?.(e.target.value)
            }}
            error={instagramError}
            disabled={disabled}
            autoComplete="off"
            dir="ltr"
            startIcon={
              <SocialIcon
                platform="instagram"
                className={cn('w-4 h-4', SOCIAL_COLOR_CLASSES.instagram)}
              />
            }
          />

          {/* Facebook */}
          <FormInput
            label={tCustomers('form.facebook')}
            placeholder={tCustomers('form.facebook_placeholder')}
            value={facebookUsername}
            onChange={(e) => {
              onFacebookChange?.(e.target.value)
            }}
            error={facebookError}
            disabled={disabled}
            autoComplete="off"
            dir="ltr"
            startIcon={
              <SocialIcon
                platform="facebook"
                className={cn('w-4 h-4', SOCIAL_COLOR_CLASSES.facebook)}
              />
            }
          />

          {/* TikTok */}
          <FormInput
            label={tCustomers('form.tiktok')}
            placeholder={tCustomers('form.tiktok_placeholder')}
            value={tiktokUsername}
            onChange={(e) => {
              onTiktokChange?.(e.target.value)
            }}
            error={tiktokError}
            disabled={disabled}
            autoComplete="off"
            dir="ltr"
            startIcon={<SocialIcon platform="tiktok" className="w-4 h-4" />}
          />

          {/* Snapchat */}
          <FormInput
            label={tCustomers('form.snapchat')}
            placeholder={tCustomers('form.snapchat_placeholder')}
            value={snapchatUsername}
            onChange={(e) => {
              onSnapchatChange?.(e.target.value)
            }}
            error={snapchatError}
            disabled={disabled}
            autoComplete="off"
            dir="ltr"
            startIcon={
              <SocialIcon
                platform="snapchat"
                className={cn('w-4 h-4', SOCIAL_COLOR_CLASSES.snapchat)}
              />
            }
          />

          {/* X (Twitter) */}
          <FormInput
            label={tCustomers('form.x')}
            placeholder={tCustomers('form.x_placeholder')}
            value={xUsername}
            onChange={(e) => {
              onXChange?.(e.target.value)
            }}
            error={xError}
            disabled={disabled}
            autoComplete="off"
            dir="ltr"
            startIcon={<SocialIcon platform="x" className="w-4 h-4" />}
          />

          {/* WhatsApp */}
          <FormInput
            label={tCustomers('form.whatsapp')}
            placeholder={tCustomers('form.whatsapp_placeholder')}
            value={whatsappNumber}
            onChange={(e) => {
              onWhatsappChange?.(e.target.value)
            }}
            error={whatsappError}
            disabled={disabled}
            autoComplete="off"
            inputMode="tel"
            dir="ltr"
            startIcon={
              <SocialIcon
                platform="whatsapp"
                className={cn('w-4 h-4', SOCIAL_COLOR_CLASSES.whatsapp)}
              />
            }
          />
        </div>
      )}

      {/* Helper text when collapsed */}
      {!isExpanded && filledCount === 0 && (
        <p className="text-sm text-base-content/60 px-4">
          {tCustomers('form.social_media_hint')}
        </p>
      )}
    </div>
  )
}
