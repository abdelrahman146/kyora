import { type ReactNode } from 'react';
import { X } from 'lucide-react';

interface DialogProps {
  /**
   * Controls dialog visibility
   */
  open: boolean;
  
  /**
   * Handler called when dialog should close
   */
  onClose?: () => void;
  
  /**
   * Dialog title
   */
  title?: string;
  
  /**
   * Dialog description
   */
  description?: string;
  
  /**
   * Dialog content
   */
  children: ReactNode;
  
  /**
   * Footer actions (buttons, etc.)
   */
  footer?: ReactNode;
  
  /**
   * Size variant
   * @default 'md'
   */
  size?: 'sm' | 'md' | 'lg' | 'xl' | 'full';
  
  /**
   * Whether to show close button
   * @default true
   */
  showCloseButton?: boolean;
  
  /**
   * Whether clicking backdrop closes dialog
   * @default true
   */
  closeOnBackdrop?: boolean;
  
  /**
   * Whether to center content vertically
   * @default false
   */
  centered?: boolean;
  
  /**
   * Additional CSS classes for modal box
   */
  className?: string;
}

/**
 * Generic Dialog Component
 * 
 * Built with daisyUI 5 modal, follows Kyora design system.
 * Mobile-first, responsive, RTL-ready, accessible.
 * 
 * @example
 * ```tsx
 * <Dialog
 *   open={isOpen}
 *   onClose={() => setIsOpen(false)}
 *   title="Confirm Action"
 *   description="Are you sure you want to proceed?"
 *   footer={
 *     <>
 *       <button onClick={onCancel} className="btn btn-ghost">
 *         Cancel
 *       </button>
 *       <button onClick={onConfirm} className="btn btn-primary">
 *         Confirm
 *       </button>
 *     </>
 *   }
 * >
 *   <p>This action cannot be undone.</p>
 * </Dialog>
 * ```
 */
export function Dialog({
  open,
  onClose,
  title,
  description,
  children,
  footer,
  size = 'md',
  showCloseButton = true,
  closeOnBackdrop = true,
  centered = false,
  className = '',
}: DialogProps) {
  if (!open) return null;

  const sizeClasses = {
    sm: 'max-w-sm',
    md: 'max-w-md',
    lg: 'max-w-lg',
    xl: 'max-w-xl',
    full: 'max-w-full',
  };

  const handleBackdropClick = (e: React.MouseEvent) => {
    if (closeOnBackdrop && onClose && e.target === e.currentTarget) {
      onClose();
    }
  };

  return (
    <dialog className="modal modal-open" onClick={handleBackdropClick}>
      <div
        className={`modal-box ${sizeClasses[size]} ${
          centered ? 'flex flex-col justify-center' : ''
        } ${className}`}
      >
        {/* Header */}
        {(title ?? showCloseButton) && (
          <div className="flex items-start justify-between gap-4 mb-4">
            <div className="flex-1 min-w-0">
              {title && (
                <h3 className="text-lg font-bold text-base-content line-clamp-2">
                  {title}
                </h3>
              )}
              {description && (
                <p className="text-sm text-base-content/70 mt-1">
                  {description}
                </p>
              )}
            </div>
            
            {showCloseButton && onClose && (
              <button
                onClick={onClose}
                className="btn btn-sm btn-circle btn-ghost shrink-0"
                aria-label="Close dialog"
              >
                <X className="w-5 h-5" />
              </button>
            )}
          </div>
        )}

        {/* Content */}
        <div className="overflow-y-auto max-h-[60vh]">{children}</div>

        {/* Footer */}
        {footer && (
          <div className="modal-action mt-6">
            <div className="flex items-center justify-end gap-2 w-full flex-wrap">
              {footer}
            </div>
          </div>
        )}
      </div>
    </dialog>
  );
}
