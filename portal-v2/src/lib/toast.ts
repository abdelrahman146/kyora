/**
 * Toast Notification Helper
 *
 * RTL-aware toast notifications with i18n support.
 * 4-second auto-dismiss, positioned based on document direction.
 */

import toast from 'react-hot-toast'
import { translateErrorAsync } from './translateError'
import type { TFunction } from 'i18next'

/**
 * Get toast position based on RTL direction
 */
function getToastPosition(): 'top-left' | 'top-right' {
  const dir = document.documentElement.getAttribute('dir')
  return dir === 'rtl' ? 'top-left' : 'top-right'
}

/**
 * Show success toast with translated message
 */
export function showSuccessToast(message: string) {
  toast.success(message, {
    duration: 4000,
    position: getToastPosition(),
  })
}

/**
 * Show error toast with translated message
 */
export function showErrorToast(message: string) {
  toast.error(message, {
    duration: 4000,
    position: getToastPosition(),
  })
}

/**
 * Show error toast from unknown error (parses and translates)
 */
export async function showErrorFromException(error: unknown, t: TFunction) {
  const message = await translateErrorAsync(error, t)
  showErrorToast(message)
}

/**
 * Show info toast
 */
export function showInfoToast(message: string) {
  toast(message, {
    duration: 4000,
    position: getToastPosition(),
  })
}

/**
 * Show loading toast (manual dismiss required)
 */
export function showLoadingToast(message: string) {
  return toast.loading(message, {
    position: getToastPosition(),
  })
}

/**
 * Dismiss specific toast
 */
export function dismissToast(toastId: string) {
  toast.dismiss(toastId)
}

/**
 * Dismiss all toasts
 */
export function dismissAllToasts() {
  toast.dismiss()
}
