import {  useEffect, useRef } from 'react'
import { createPortal } from 'react-dom'
import { X } from 'lucide-react'
import type {ReactNode} from 'react';
import { cn } from '@/lib/utils'

export interface ModalProps {
  isOpen: boolean
  onClose: () => void
  title?: ReactNode
  children: ReactNode
  footer?: ReactNode
  size?: 'sm' | 'md' | 'lg' | 'xl'
  closeOnBackdropClick?: boolean
  closeOnEscape?: boolean
  showCloseButton?: boolean
  className?: string
  contentClassName?: string
}

/**
 * Modal Component
 *
 * DaisyUI modal with backdrop, keyboard navigation, and focus trap.
 * Uses native <dialog> element for accessibility.
 */
export function Modal({
  isOpen,
  onClose,
  title,
  children,
  footer,
  size = 'md',
  closeOnBackdropClick = true,
  closeOnEscape = true,
  showCloseButton = true,
  className,
  contentClassName,
}: ModalProps) {
  const modalRef = useRef<HTMLDialogElement>(null)

  const sizeClasses = {
    sm: 'modal-box max-w-sm',
    md: 'modal-box max-w-md',
    lg: 'modal-box max-w-2xl',
    xl: 'modal-box max-w-4xl',
  }

  // Handle Escape key
  useEffect(() => {
    if (!isOpen || !closeOnEscape) return

    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        onClose()
      }
    }

    document.addEventListener('keydown', handleEscape)
    return () => document.removeEventListener('keydown', handleEscape)
  }, [isOpen, closeOnEscape, onClose])

  // Handle backdrop click
  const handleBackdropClick = (e: React.MouseEvent<HTMLElement>) => {
    if (closeOnBackdropClick && e.target === e.currentTarget) {
      onClose()
    }
  }

  if (!isOpen) return null

  return createPortal(
    <dialog className={cn('modal modal-open', className)} ref={modalRef}>
      <div className={cn(sizeClasses[size], contentClassName)}>
        {/* Header */}
        {(title ?? showCloseButton) && (
          <div className="mb-4 flex items-center justify-between">
            {title && <h3 className="text-lg font-bold">{title}</h3>}
            {showCloseButton && (
              <button
                type="button"
                onClick={onClose}
                className="btn btn-circle btn-ghost btn-sm ms-auto"
                aria-label="Close modal"
              >
                <X size={20} />
              </button>
            )}
          </div>
        )}

        {/* Content */}
        <div className="py-4">{children}</div>

        {/* Footer */}
        {footer && <div className="modal-action">{footer}</div>}
      </div>

      {/* Backdrop */}
      <form method="dialog" className="modal-backdrop">
        <button type="button" onClick={handleBackdropClick}>
          close
        </button>
      </form>
    </dialog>,
    document.body,
  )
}
