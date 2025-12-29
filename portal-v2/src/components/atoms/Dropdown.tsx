import { useEffect, useRef, useState } from 'react'
import { cn } from '../../lib/utils'
import type { ReactNode } from 'react'

export interface DropdownProps {
  /**
   * Trigger element (button, icon, etc.)
   */
  trigger: ReactNode

  /**
   * Dropdown content
   */
  children: ReactNode

  /**
   * Alignment of dropdown relative to trigger
   * @default 'start'
   */
  align?: 'start' | 'end'

  /**
   * Width of dropdown
   * @default '200px'
   */
  width?: string

  /**
   * Additional CSS classes for dropdown container
   */
  className?: string

  /**
   * Additional CSS classes for content
   */
  contentClassName?: string
}

/**
 * Generic Dropdown Component
 *
 * Handles click-outside detection, keyboard navigation (Escape to close),
 * RTL-aware alignment, and accessible interactions.
 *
 * @example
 * ```tsx
 * <Dropdown
 *   trigger={<button className="btn">Menu</button>}
 *   align="end"
 * >
 *   <ul className="menu p-2">
 *     <li><a>Profile</a></li>
 *     <li><a>Settings</a></li>
 *     <li><a>Logout</a></li>
 *   </ul>
 * </Dropdown>
 * ```
 */
export function Dropdown({
  trigger,
  children,
  align = 'start',
  width = '200px',
  className,
  contentClassName,
}: DropdownProps) {
  const [isOpen, setIsOpen] = useState(false)
  const dropdownRef = useRef<HTMLDivElement>(null)

  // Click outside to close
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false)
      }
    }

    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside)
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside)
    }
  }, [isOpen])

  // Escape key to close
  useEffect(() => {
    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        setIsOpen(false)
      }
    }

    if (isOpen) {
      document.addEventListener('keydown', handleEscape)
    }

    return () => {
      document.removeEventListener('keydown', handleEscape)
    }
  }, [isOpen])

  const handleToggle = () => {
    setIsOpen(!isOpen)
  }

  const handleClose = () => {
    setIsOpen(false)
  }

  return (
    <div ref={dropdownRef} className={cn('relative', className)}>
      {/* Trigger */}
      <div onClick={handleToggle} role="button" tabIndex={0}>
        {trigger}
      </div>

      {/* Dropdown Content */}
      {isOpen && (
        <div
          className={cn(
            'absolute top-full mt-2 z-50',
            'bg-base-100 rounded-lg shadow-lg border border-base-300',
            'animate-in fade-in slide-in-from-top-2 duration-200',
            align === 'end' && 'end-0',
            align === 'start' && 'start-0',
            contentClassName
          )}
          style={{ width }}
          onClick={handleClose}
        >
          {children}
        </div>
      )}
    </div>
  )
}
