/**
 * Date Formatting Utilities
 *
 * Provides consistent date and time formatting across the application.
 * Uses Intl.DateTimeFormat for localized and timezone-aware formatting.
 *
 * Features:
 * - Relative time (e.g., "2 hours ago", "3 days ago")
 * - Short date format (e.g., "Jan 3, 2026")
 * - Full date format (e.g., "January 3, 2026")
 * - Date with time (e.g., "Jan 3, 2026 2:30 PM")
 * - Time only (e.g., "2:30 PM")
 * - Auto-detects user's locale or uses provided locale
 */

/**
 * Get the default locale for date formatting
 * Falls back to 'en-US' if not available
 */
function getDefaultLocale(): string {
  if (typeof navigator !== 'undefined' && navigator.language) {
    return navigator.language
  }
  return 'en-US'
}

/**
 * Format date as short date string (e.g., "Jan 3, 2026")
 */
export function formatDateShort(
  date: string | Date | null | undefined,
  locale?: string,
): string {
  if (!date) return '—'

  const dateObj = typeof date === 'string' ? new Date(date) : date

  if (isNaN(dateObj.getTime())) return '—'

  return new Intl.DateTimeFormat(locale || getDefaultLocale(), {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  }).format(dateObj)
}

/**
 * Format date as full date string (e.g., "January 3, 2026")
 */
export function formatDateLong(
  date: string | Date | null | undefined,
  locale?: string,
): string {
  if (!date) return '—'

  const dateObj = typeof date === 'string' ? new Date(date) : date

  if (isNaN(dateObj.getTime())) return '—'

  return new Intl.DateTimeFormat(locale || getDefaultLocale(), {
    month: 'long',
    day: 'numeric',
    year: 'numeric',
  }).format(dateObj)
}

/**
 * Format date with time (e.g., "Jan 3, 2026 2:30 PM")
 */
export function formatDateTime(
  date: string | Date | null | undefined,
  locale?: string,
): string {
  if (!date) return '—'

  const dateObj = typeof date === 'string' ? new Date(date) : date

  if (isNaN(dateObj.getTime())) return '—'

  return new Intl.DateTimeFormat(locale || getDefaultLocale(), {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
    hour12: true,
  }).format(dateObj)
}

/**
 * Format time only (e.g., "2:30 PM")
 */
export function formatTime(
  date: string | Date | null | undefined,
  locale?: string,
): string {
  if (!date) return '—'

  const dateObj = typeof date === 'string' ? new Date(date) : date

  if (isNaN(dateObj.getTime())) return '—'

  return new Intl.DateTimeFormat(locale || getDefaultLocale(), {
    hour: 'numeric',
    minute: '2-digit',
    hour12: true,
  }).format(dateObj)
}

/**
 * Format date as relative time (e.g., "2 hours ago", "3 days ago")
 */
export function formatRelativeTime(
  date: string | Date | null | undefined,
  locale?: string,
): string {
  if (!date) return '—'

  const dateObj = typeof date === 'string' ? new Date(date) : date

  if (isNaN(dateObj.getTime())) return '—'

  const now = new Date()
  const diffMs = now.getTime() - dateObj.getTime()
  const diffSecs = Math.floor(diffMs / 1000)
  const diffMins = Math.floor(diffSecs / 60)
  const diffHours = Math.floor(diffMins / 60)
  const diffDays = Math.floor(diffHours / 24)
  const diffMonths = Math.floor(diffDays / 30)
  const diffYears = Math.floor(diffDays / 365)

  const rtf = new Intl.RelativeTimeFormat(locale || getDefaultLocale(), {
    numeric: 'auto',
  })

  if (diffYears > 0) return rtf.format(-diffYears, 'year')
  if (diffMonths > 0) return rtf.format(-diffMonths, 'month')
  if (diffDays > 0) return rtf.format(-diffDays, 'day')
  if (diffHours > 0) return rtf.format(-diffHours, 'hour')
  if (diffMins > 0) return rtf.format(-diffMins, 'minute')
  return rtf.format(-diffSecs, 'second')
}
