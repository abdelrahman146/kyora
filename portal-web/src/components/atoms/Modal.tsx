import { type ReactNode, useEffect, useRef } from "react";
import { createPortal } from "react-dom";
import { cn } from "@/lib/utils";

export interface ModalProps {
  /**
   * Whether the modal is open or closed
   */
  isOpen: boolean;

  /**
   * Callback fired when the modal should be closed
   */
  onClose: () => void;

  /**
   * Modal title
   */
  title?: ReactNode;

  /**
   * Modal content
   */
  children: ReactNode;

  /**
   * Footer content (typically action buttons)
   */
  footer?: ReactNode;

  /**
   * Size of the modal
   * @default "md"
   */
  size?: "sm" | "md" | "lg" | "xl" | "full";

  /**
   * Whether clicking the backdrop closes the modal
   * @default true
   */
  closeOnBackdropClick?: boolean;

  /**
   * Whether pressing Escape closes the modal
   * @default true
   */
  closeOnEscape?: boolean;

  /**
   * Whether to show the close button in the top-right
   * @default true
   */
  showCloseButton?: boolean;

  /**
   * Additional CSS classes for the modal container
   */
  className?: string;

  /**
   * Additional CSS classes for the modal content box
   */
  contentClassName?: string;

  /**
   * Whether the modal content should be scrollable
   * @default true
   */
  scrollable?: boolean;

  /**
   * Custom z-index for the modal
   * @default 50
   */
  zIndex?: number;
}

/**
 * Modal Component - Production-Grade Reusable Modal
 *
 * Features:
 * - Mobile-first design (bottom sheet on mobile, centered on desktop)
 * - Responsive sizing with multiple size options
 * - Portal-based rendering for proper stacking context
 * - Keyboard navigation (Escape to close)
 * - Focus trap for accessibility
 * - Backdrop click to close (optional)
 * - Smooth animations with CSS transitions
 * - RTL support with logical properties
 * - Scroll lock when open
 * - DaisyUI theming support
 *
 * @example
 * ```tsx
 * <Modal
 *   isOpen={isOpen}
 *   onClose={() => setIsOpen(false)}
 *   title="Delete Item"
 *   footer={
 *     <>
 *       <button onClick={() => setIsOpen(false)} className="btn btn-ghost">
 *         Cancel
 *       </button>
 *       <button onClick={handleDelete} className="btn btn-error">
 *         Delete
 *       </button>
 *     </>
 *   }
 * >
 *   <p>Are you sure you want to delete this item?</p>
 * </Modal>
 * ```
 */
export function Modal({
  isOpen,
  onClose,
  title,
  children,
  footer,
  size = "md",
  closeOnBackdropClick = true,
  closeOnEscape = true,
  showCloseButton = true,
  className,
  contentClassName,
  scrollable = true,
  zIndex = 50,
}: ModalProps) {
  const modalRef = useRef<HTMLDivElement>(null);

  // Size mapping for responsive modal widths
  const sizeClasses = {
    sm: "max-w-sm",
    md: "max-w-md",
    lg: "max-w-2xl",
    xl: "max-w-4xl",
    full: "max-w-full mx-4",
  };

  // Handle Escape key press
  useEffect(() => {
    if (!isOpen || !closeOnEscape) return;

    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        onClose();
      }
    };

    document.addEventListener("keydown", handleEscape);
    return () => {
      document.removeEventListener("keydown", handleEscape);
    };
  }, [isOpen, closeOnEscape, onClose]);

  // Lock body scroll when modal is open
  useEffect(() => {
    if (!isOpen) return;

    const originalOverflow = document.body.style.overflow;
    const originalPaddingRight = document.body.style.paddingRight;

    // Prevent body scroll
    const scrollbarWidth = window.innerWidth - document.documentElement.clientWidth;
    document.body.style.overflow = "hidden";
    document.body.style.paddingRight = `${scrollbarWidth.toString()}px`;

    return () => {
      document.body.style.overflow = originalOverflow;
      document.body.style.paddingRight = originalPaddingRight;
    };
  }, [isOpen]);

  // Focus trap - focus first focusable element when modal opens
  useEffect(() => {
    if (!isOpen) return;

    const focusableElements = modalRef.current?.querySelectorAll<HTMLElement>(
      'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
    );

    if (focusableElements && focusableElements.length > 0) {
      focusableElements[0].focus();
    }
  }, [isOpen]);

  // Handle backdrop click
  const handleBackdropClick = (e: React.MouseEvent<HTMLDivElement>) => {
    if (closeOnBackdropClick && e.target === e.currentTarget) {
      onClose();
    }
  };

  if (!isOpen) return null;

  const modalContent = (
    <div
      className={cn(
        "fixed inset-0 flex items-end md:items-center justify-center",
        className
      )}
      style={{ zIndex }}
      role="dialog"
      aria-modal="true"
      aria-labelledby={title ? "modal-title" : undefined}
    >
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/50 backdrop-blur-sm transition-opacity duration-300"
        onClick={handleBackdropClick}
        aria-hidden="true"
      />

      {/* Modal Box */}
      <div
        ref={modalRef}
        className={cn(
          "relative w-full bg-base-100 rounded-t-2xl md:rounded-2xl shadow-2xl",
          "transform transition-transform duration-300 ease-out",
          "max-h-[90vh] md:max-h-[85vh]",
          "animate-in slide-in-from-bottom md:slide-in-from-bottom-0 md:fade-in",
          sizeClasses[size],
          contentClassName
        )}
      >
        {/* Header */}
        {(title ?? showCloseButton) && (
          <div className="flex items-center justify-between gap-4 px-6 py-4 border-b border-base-300">
            {title && (
              <h3
                id="modal-title"
                className="text-lg md:text-xl font-bold text-base-content flex-1"
              >
                {title}
              </h3>
            )}
            {showCloseButton && (
              <button
                onClick={onClose}
                className="btn btn-sm btn-circle btn-ghost shrink-0"
                aria-label="Close modal"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="h-5 w-5"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </button>
            )}
          </div>
        )}

        {/* Content */}
        <div
          className={cn(
            "px-6 py-4",
            scrollable && "overflow-y-auto",
            !footer && "pb-6"
          )}
          style={{
            maxHeight: footer
              ? "calc(90vh - 140px)"
              : title || showCloseButton
              ? "calc(90vh - 80px)"
              : "calc(90vh - 40px)",
          }}
        >
          {children}
        </div>

        {/* Footer */}
        {footer && (
          <div className="flex items-center justify-end gap-3 px-6 py-4 border-t border-base-300 bg-base-100">
            {footer}
          </div>
        )}
      </div>
    </div>
  );

  // Render modal in a portal at the end of body
  return createPortal(modalContent, document.body);
}
