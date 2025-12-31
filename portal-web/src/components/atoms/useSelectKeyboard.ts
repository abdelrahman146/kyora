/**
 * useSelectKeyboard Hook
 *
 * Manages keyboard navigation for select components.
 * Handles arrow keys, Enter, Escape, Home, End, and Tab.
 */

import { useCallback, useState } from 'react'
import type { FormSelectOption } from './FormSelect'

interface UseSelectKeyboardProps<T> {
  isOpen: boolean
  setIsOpen: (open: boolean) => void
  filteredOptions: Array<FormSelectOption<T>>
  onSelectOption: (value: T) => void
  onClose: () => void
  disabled?: boolean
}

interface UseSelectKeyboardReturn {
  focusedIndex: number
  setFocusedIndex: (index: number) => void
  handleKeyDown: (e: React.KeyboardEvent) => void
}

export function useSelectKeyboard<T extends string>({
  isOpen,
  setIsOpen,
  filteredOptions,
  onSelectOption,
  onClose,
  disabled = false,
}: UseSelectKeyboardProps<T>): UseSelectKeyboardReturn {
  const [focusedIndex, setFocusedIndex] = useState(-1)

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
      if (disabled) return

      switch (e.key) {
        case 'Enter':
        case ' ':
          if (!isOpen) {
            setIsOpen(true)
          } else if (
            focusedIndex >= 0 &&
            focusedIndex < filteredOptions.length
          ) {
            const option = filteredOptions[focusedIndex]
            if (!option.disabled) {
              onSelectOption(option.value)
            }
          }
          e.preventDefault()
          break

        case 'Escape':
          if (isOpen) {
            onClose()
            e.preventDefault()
          }
          break

        case 'ArrowDown':
          if (!isOpen) {
            setIsOpen(true)
          } else {
            setFocusedIndex((prev) => {
              const nextIndex = prev + 1
              return nextIndex < filteredOptions.length ? nextIndex : prev
            })
          }
          e.preventDefault()
          break

        case 'ArrowUp':
          if (isOpen) {
            setFocusedIndex((prev) => (prev > 0 ? prev - 1 : 0))
            e.preventDefault()
          }
          break

        case 'Home':
          if (isOpen) {
            setFocusedIndex(0)
            e.preventDefault()
          }
          break

        case 'End':
          if (isOpen) {
            setFocusedIndex(filteredOptions.length - 1)
            e.preventDefault()
          }
          break

        case 'Tab':
          if (isOpen) {
            onClose()
          }
          break
      }
    },
    [
      disabled,
      isOpen,
      focusedIndex,
      filteredOptions,
      onSelectOption,
      onClose,
      setIsOpen,
    ],
  )

  return {
    focusedIndex,
    setFocusedIndex,
    handleKeyDown,
  }
}
