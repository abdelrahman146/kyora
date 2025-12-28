/**
 * BottomSheet Component
 *
 * A versatile, production-grade drawer component that adapts to different screen sizes:
 * - Mobile: Bottom sheet with slide-up animation
 * - Desktop: Side drawer with slide-in animation
 *
 * Features:
 * - Fully responsive with mobile/desktop variants
 * - RTL/LTR support with automatic direction detection
 * - Smooth animations using CSS transitions and requestAnimationFrame
 * - Accessibility (ARIA attributes, keyboard navigation, focus management)
 * - Customizable positioning (left/right for desktop, bottom for mobile)
 * - Optional header, footer, and overlay
 * - Prevents body scroll when open
 * - Escape key to close
 * - Click outside to close (optional)
 * - Flexible content area with scroll support
 * - Production-ready with TypeScript types
 *
 * @example
 * // Basic usage
 * <BottomSheet
 *   isOpen={isOpen}
 *   onClose={() => setIsOpen(false)}
 *   title="Settings"
 * >
 *   <p>Your content here</p>
 * </BottomSheet>
 *
 * @example
 * // With custom footer
 * <BottomSheet
 *   isOpen={isOpen}
 *   onClose={() => setIsOpen(false)}
 *   title="Filters"
 *   footer={
 *     <div className="flex gap-2">
 *       <button className="btn btn-ghost flex-1">Reset</button>
 *       <button className="btn btn-primary flex-1">Apply</button>
 *     </div>
 *   }
 * >
 *   <FilterContent />
 * </BottomSheet>
 *
 * @example
 * // Left-side drawer on desktop
 * <BottomSheet
 *   isOpen={isOpen}
 *   onClose={() => setIsOpen(false)}
 *   title="Navigation"
 *   side="start"
 * >
 *   <Navigation />
 * </BottomSheet>
 */

import { useEffect, useState, useRef, useId, type ReactNode } from "react";
import { X } from "lucide-react";
import { useMediaQuery } from "../../hooks/useMediaQuery";
import { useTranslation } from "react-i18next";

export interface BottomSheetProps {
  /** Controls the open/closed state of the drawer */
  isOpen: boolean;

  /** Callback fired when the drawer should close */
  onClose: () => void;

  /** Title displayed in the header */
  title?: string;

  /** Main content of the drawer */
  children: ReactNode;

  /** Optional footer content (buttons, actions, etc.) */
  footer?: ReactNode;

  /** 
   * Side where the drawer appears on desktop
   * - 'start': Left in LTR, Right in RTL
   * - 'end': Right in LTR, Left in RTL
   * @default 'end'
   */
  side?: "start" | "end";

  /**
   * Maximum width of the drawer on desktop
   * @default 'md' (28rem / 448px)
   */
  size?: "sm" | "md" | "lg" | "xl" | "full";

  /**
   * Whether clicking the overlay closes the drawer
   * @default true
   */
  closeOnOverlayClick?: boolean;

  /**
   * Whether pressing Escape closes the drawer
   * @default true
   */
  closeOnEscape?: boolean;

  /**
   * Whether to show the header
   * @default true
   */
  showHeader?: boolean;

  /**
   * Whether to show the close button in header
   * @default true
   */
  showCloseButton?: boolean;

  /**
   * Custom header content (overrides title)
   */
  header?: ReactNode;

  /**
   * Additional CSS classes for the drawer container
   */
  className?: string;

  /**
   * Additional CSS classes for the content area
   */
  contentClassName?: string;

  /**
   * Additional CSS classes for the footer area
   */
  footerClassName?: string;

  /**
   * ARIA label for accessibility
   */
  ariaLabel?: string;

  /**
   * ID for aria-labelledby (references the title element)
   */
  ariaLabelledBy?: string;
}

const sizeClasses = {
  sm: "max-w-sm",
  md: "max-w-md",
  lg: "max-w-lg",
  xl: "max-w-xl",
  full: "max-w-full",
};

export function BottomSheet({
  isOpen,
  onClose,
  title,
  children,
  footer,
  side = "end",
  size = "md",
  closeOnOverlayClick = true,
  closeOnEscape = true,
  showHeader = true,
  showCloseButton = true,
  header,
  className = "",
  contentClassName = "",
  footerClassName = "",
  ariaLabel,
  ariaLabelledBy,
}: BottomSheetProps) {
  const { t } = useTranslation();
  const isMobile = useMediaQuery("(max-width: 768px)");
  const [isAnimating, setIsAnimating] = useState(false);
  const drawerRef = useRef<HTMLDivElement>(null);
  const previousActiveElement = useRef<HTMLElement | null>(null);
  const generatedId = useId();

  // Generate unique ID for aria-labelledby if title is provided
  const titleId = ariaLabelledBy ?? (title ? `drawer-title-${generatedId}` : undefined);

  // Handle open/close animation with requestAnimationFrame to avoid cascading renders
  useEffect(() => {
    if (isOpen) {
      // Store the currently focused element to restore focus later
      previousActiveElement.current = document.activeElement as HTMLElement;

      requestAnimationFrame(() => {
        setIsAnimating(true);
        // Focus the drawer after animation starts
        requestAnimationFrame(() => {
          drawerRef.current?.focus();
        });
      });
    } else {
      requestAnimationFrame(() => {
        setIsAnimating(false);
      });

      // Restore focus to the previously focused element
      if (previousActiveElement.current) {
        previousActiveElement.current.focus();
        previousActiveElement.current = null;
      }
    }
  }, [isOpen]);

  // Prevent body scroll when drawer is open
  useEffect(() => {
    if (isOpen) {
      const scrollBarWidth = window.innerWidth - document.documentElement.clientWidth;
      document.body.style.overflow = "hidden";
      document.body.style.paddingRight = `${String(scrollBarWidth)}px`;
    } else {
      document.body.style.overflow = "";
      document.body.style.paddingRight = "";
    }

    return () => {
      document.body.style.overflow = "";
      document.body.style.paddingRight = "";
    };
  }, [isOpen]);

  // Handle Escape key
  useEffect(() => {
    if (!closeOnEscape || !isOpen) return;

    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        e.preventDefault();
        onClose();
      }
    };

    document.addEventListener("keydown", handleEscape);
    return () => {
      document.removeEventListener("keydown", handleEscape);
    };
  }, [isOpen, onClose, closeOnEscape]);

  // Handle overlay click
  const handleOverlayClick = () => {
    if (closeOnOverlayClick) {
      onClose();
    }
  };

  // Don't render anything if not open
  if (!isOpen) return null;

  // Determine drawer positioning classes
  const getDrawerClasses = () => {
    if (isMobile) {
      // Mobile: Bottom sheet
      return `bottom-0 start-0 end-0 rounded-t-2xl max-h-[85vh] ${
        isAnimating ? "translate-y-0" : "translate-y-full"
      }`;
    }

    // Desktop: Side drawer
    const sideClass = side === "start" ? "start-0" : "end-0";
    const translateClass = isAnimating
      ? "translate-x-0"
      : side === "start"
      ? "ltr:-translate-x-full rtl:translate-x-full"
      : "ltr:translate-x-full rtl:-translate-x-full";

    return `top-0 ${sideClass} h-full w-full ${sizeClasses[size]} ${translateClass}`;
  };

  return (
    <>
      {/* Backdrop Overlay */}
      <div
        className={`fixed inset-0 bg-black/50 z-40 transition-opacity duration-300 ${
          isAnimating ? "opacity-100" : "opacity-0"
        }`}
        onClick={handleOverlayClick}
        aria-hidden="true"
      />

      {/* Drawer Container */}
      <div
        ref={drawerRef}
        className={`fixed z-50 bg-base-100 shadow-xl transition-transform duration-300 ease-in-out overflow-y-auto ${getDrawerClasses()} ${className}`}
        role="dialog"
        aria-modal="true"
        aria-label={ariaLabel}
        aria-labelledby={titleId}
        tabIndex={-1}
      >
        {/* Header */}
        {showHeader && (header ?? title) && (
          <div className="sticky top-0 bg-base-100 border-b border-base-300 px-4 py-4 flex items-center justify-between z-10">
            {header ?? (
              <h2 id={titleId} className="text-lg font-semibold text-base-content">
                {title}
              </h2>
            )}
            {showCloseButton && (
              <button
                type="button"
                onClick={onClose}
                className="btn btn-ghost btn-sm btn-circle"
                aria-label={t("common.close")}
              >
                <X size={20} />
              </button>
            )}
          </div>
        )}

        {/* Content Area */}
        <div className={`px-4 py-6 ${contentClassName}`}>{children}</div>

        {/* Footer */}
        {footer && (
          <div className={`sticky bottom-0 bg-base-100 border-t border-base-300 px-4 py-4 ${footerClassName}`}>
            {footer}
          </div>
        )}
      </div>
    </>
  );
}
