import { clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'
import type { ClassValue } from 'clsx'

/**
 * Utility function to merge Tailwind CSS classes
 *
 * Combines clsx for conditional classes and tailwind-merge to prevent conflicts.
 */
export function cn(...inputs: Array<ClassValue>) {
  return twMerge(clsx(inputs))
}

export function formatCountdownDuration(totalSeconds: number): string {
  if (!Number.isFinite(totalSeconds)) return '0s'

  const secondsInt = Math.max(0, Math.floor(totalSeconds))
  const hours = Math.floor(secondsInt / 3600)
  const minutes = Math.floor((secondsInt % 3600) / 60)
  const seconds = secondsInt % 60

  if (hours > 0) {
    return `${String(hours)}hr:${String(minutes).padStart(2, '0')}m:${String(
      seconds,
    ).padStart(2, '0')}s`
  }

  if (minutes > 0) {
    return `${String(minutes)}m:${String(seconds).padStart(2, '0')}s`
  }

  return `${String(seconds)}s`
}
