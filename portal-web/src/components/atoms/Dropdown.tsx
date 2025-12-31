import { useCallback, useEffect, useRef, useState } from 'react'
import { createPortal } from 'react-dom'
import { cn } from '../../lib/utils'
import { useMediaQuery } from '../../hooks/useMediaQuery'
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
   * Width of dropdown (desktop only)
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

  /**
   * Title for mobile modal (optional)
   */
  mobileTitle?: string

  /**
   * Whether to show close button on mobile modal
   * @default true
   */
  showMobileClose?: boolean

  /**
   * Callback when dropdown opens
   */
  onOpen?: () => void

  /**
   * Callback when dropdown closes
   */
  onClose?: () => void
}

/**
 * Mobile-First Dropdown Component
 *
 * - On mobile (< 768px): Renders as a bottom sheet modal with smooth animations
 * - On desktop: Renders as a traditional dropdown menu
 *
 * Features:
 * - Full accessibility (ARIA attributes, keyboard navigation, focus management)
 * - RTL/LTR support with logical properties
 * - Body scroll lock on mobile modal
 * - Click outside to close
 * - Escape key to close
 * - Smooth animations
 * - Focus trap on mobile
 * - Production-ready edge case handling
 *
 * @example
 * ```tsx
 * <Dropdown
 *   trigger={<button className="btn">Menu</button>}
 *   align="end"
 *   mobileTitle="Options"
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
  mobileTitle,
  showMobileClose = true,
  onOpen,
  onClose,
}: DropdownProps) {
  const [isOpen, setIsOpen] = useState(false)
  const [isAnimating, setIsAnimating] = useState(false)
  const dropdownRef = useRef<HTMLDivElement>(null)
  const modalContentRef = useRef<HTMLDivElement>(null)
  const triggerRef = useRef<HTMLElement | null>(null)
  const isMobile = useMediaQuery('(max-width: 768px)')

  // Store original trigger element for focus restoration
  useEffect(() => {
    if (dropdownRef.current) {
      triggerRef.current = dropdownRef.current.querySelector(
        '[role="button"]',
      ) as HTMLElement
    }
  }, [])

  // Body scroll lock for mobile modal
  useEffect(() => {
    if (isOpen && isMobile) {
      const originalStyle = window.getComputedStyle(document.body).overflow
      document.body.style.overflow = 'hidden'
      
      return () => {
        document.body.style.overflow = originalStyle
      }
    }
  }, [isOpen, isMobile])

  // Focus trap for mobile modal
  useEffect(() => {
    if (isOpen && isMobile && modalContentRef.current) {
      const modalElement = modalContentRef.current
      const focusableElements = modalElement.querySelectorAll<HTMLElement>(
        'a[href], button:not([disabled]), textarea:not([disabled]), input:not([disabled]), select:not([disabled]), [tabindex]:not([tabindex="-1"])',
      )
      
      if (focusableElements.length === 0) return

      const firstElement = focusableElements[0]
      const lastElement = focusableElements[focusableElements.length - 1]

      // Focus first element
      firstElement.focus()

      const handleTabKey = (e: KeyboardEvent) => {
        if (e.key !== 'Tab') return

        if (e.shiftKey) {
          // Shift + Tab
          if (document.activeElement === firstElement) {
            e.preventDefault()
            lastElement.focus()
          }
        } else {
          // Tab
          if (document.activeElement === lastElement) {
            e.preventDefault()
            firstElement.focus()
          }
        }
      }

      modalElement.addEventListener('keydown', handleTabKey)
      return () => modalElement.removeEventListener('keydown', handleTabKey)
    }
  }, [isOpen, isMobile])

  // Click outside to close (desktop only)
  useEffect(() => {
    if (!isMobile && isOpen) {
      const handleClickOutside = (event: MouseEvent) => {
        if (
          dropdownRef.current &&
          !dropdownRef.current.contains(event.target as Node)
        ) {
          handleClose()
        }
      }

      // Use capture phase to handle events before they bubble
      document.addEventListener('mousedown', handleClickOutside, true)
      return () => {
        document.removeEventListener('mousedown', handleClickOutside, true)
      }
    }
  }, [isOpen, isMobile])

  // Escape key to close
  useEffect(() => {
    if (isOpen) {
      const handleEscape = (event: KeyboardEvent) => {
        if (event.key === 'Escape') {
          event.preventDefault()
          event.stopPropagation()
          handleClose()
        }
      }

      document.addEventListener('keydown', handleEscape, true)
      return () => {
        document.removeEventListener('keydown', handleEscape, true)
      }
    }
  }, [isOpen])

  const handleOpen = useCallback(() => {
    setIsOpen(true)
    setIsAnimating(true)
    onOpen?.()
  }, [onOpen])

  const handleClose = useCallback(() => {
    setIsAnimating(false)
    
    // Wait for animation to complete before removing from DOM
    setTimeout(() => {
      setIsOpen(false)
      onClose?.()
      
      // Restore focus to trigger
      triggerRef.current?.focus()
    }, 200)
  }, [onClose])

  const handleToggle = () => {
    if (isOpen) {
      handleClose()
    } else {
      handleOpen()
    }
  }

  const handleKeyDown = (event: React.KeyboardEvent) => {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault()
      handleToggle()
    }
  }

  const handleContentClick = (event: React.MouseEvent) => {
    // Allow clicks inside content to propagate (e.g., menu items)
    // Close only if clicking directly on menu items or explicitly closing
    const target = event.target as HTMLElement
    
    // Close if clicking on a link or button inside the dropdown
    if (
      target.tagName === 'A' ||
      target.tagName === 'BUTTON' ||
      target.closest('a') ||
      target.closest('button')
    ) {
      handleClose()
    }
  }

  // Mobile Modal Rendering
  const renderMobileModal = () => {
    if (!isOpen) return null

    return createPortal(
      <div
        className={cn(
          'fixed inset-0 z-50 flex items-end justify-center',
          'transition-opacity duration-200',
          isAnimating ? 'opacity-100' : 'opacity-0',
        )}
        role="dialog"
        aria-modal="true"
        aria-label={mobileTitle || 'Menu'}
      >
        {/* Backdrop */}
        <div
          className="absolute inset-0 bg-base-content/50 backdrop-blur-sm"
          onClick={handleClose}
          aria-hidden="true"
        />

        {/* Modal Content */}
        <div
          ref={modalContentRef}
          className={cn(
            'relative w-full max-h-[85vh]',
            'bg-base-100 rounded-t-xl shadow-xl',
            'overflow-hidden flex flex-col',
            'transition-transform duration-200 ease-out',
            isAnimating
              ? 'translate-y-0'
              : 'translate-y-full',
            contentClassName,
          )}
          onClick={handleContentClick}
        >
          {/* Drag Handle */}
          <div className="flex justify-center py-2 px-4 shrink-0">
            <div className="w-12 h-1 bg-base-300 rounded-full" aria-hidden="true" />
          </div>

          {/* Header (if title provided) */}
          {(mobileTitle || showMobileClose) && (
            <div className="flex items-center justify-between px-4 pb-3 shrink-0 border-b border-base-300">
              {mobileTitle && (
                <h2 className="text-lg font-semibold text-base-content">
                  {mobileTitle}
                </h2>
              )}
              {showMobileClose && (
                <button
                  type="button"
                  onClick={handleClose}
                  className="btn btn-ghost btn-sm btn-square ms-auto"
                  aria-label="Close"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="20"
                    height="20"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  >
                    <line x1="18" y1="6" x2="6" y2="18" />
                    <line x1="6" y1="6" x2="18" y2="18" />
                  </svg>
                </button>
              )}
            </div>
          )}

          {/* Content */}
          <div className="overflow-y-auto flex-1 overscroll-contain">
            {children}
          </div>
        </div>
      </div>,
      document.body,
    )
  }

  // Desktop Dropdown Rendering
  const renderDesktopDropdown = () => {
    if (!isOpen) return null

    return (
      <div
        className={cn(
          'absolute top-full mt-2 z-50',
          'bg-base-100 rounded-lg shadow-lg border border-base-300',
          'transition-all duration-200 origin-top',
          isAnimating
            ? 'opacity-100 scale-100'
            : 'opacity-0 scale-95',
          align === 'end' && 'end-0',
          align === 'start' && 'start-0',
          contentClassName,
        )}
        style={{ width, minWidth: '160px' }}
        onClick={handleContentClick}
        role="menu"
        aria-orientation="vertical"
      >
        {children}
      </div>
    )
  }

  return (
    <div ref={dropdownRef} className={cn('relative', className)}>
      {/* Trigger */}
      <div
        onClick={handleToggle}
        onKeyDown={handleKeyDown}
        role="button"
        tabIndex={0}
        aria-haspopup="true"
        aria-expanded={isOpen}
        aria-label="Open menu"
      >
        {trigger}
      </div>

      {/* Render based on screen size */}
      {isMobile ? renderMobileModal() : renderDesktopDropdown()}
    </div>
  )
}
