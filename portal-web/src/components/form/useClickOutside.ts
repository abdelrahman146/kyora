/**
 * useClickOutside Hook
 *
 * Detects clicks outside a referenced element and triggers callback.
 * Handles both mouse and touch events with capture phase for better detection.
 */

import { useEffect, useRef } from 'react'

interface UseClickOutsideProps {
  isActive: boolean
  onClickOutside: () => void
}

export function useClickOutside<T extends HTMLElement = HTMLDivElement>({
  isActive,
  onClickOutside,
}: UseClickOutsideProps) {
  const ref = useRef<T>(null)

  useEffect(() => {
    if (!isActive) return

    const handleClickOutside = (event: MouseEvent | TouchEvent) => {
      if (ref.current && !ref.current.contains(event.target as Node)) {
        onClickOutside()
      }
    }

    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        onClickOutside()
      }
    }

    // Use capture phase for better click-outside detection
    document.addEventListener('mousedown', handleClickOutside, true)
    document.addEventListener('touchstart', handleClickOutside, true)
    document.addEventListener('keydown', handleEscape)

    return () => {
      document.removeEventListener('mousedown', handleClickOutside, true)
      document.removeEventListener('touchstart', handleClickOutside, true)
      document.removeEventListener('keydown', handleEscape)
    }
  }, [isActive, onClickOutside])

  return ref
}
