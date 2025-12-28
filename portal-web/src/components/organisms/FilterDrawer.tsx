/**
 * FilterDrawer Component
 *
 * A bottom sheet/drawer component for filtering lists on mobile and desktop.
 *
 * Features:
 * - Mobile: Swipeable bottom sheet
 * - Desktop: Side drawer or modal
 * - Accessible (keyboard navigation, ARIA attributes)
 * - RTL-compatible
 * - Responsive design
 */

import { useEffect } from "react";
import type { ReactNode } from "react";
import { X } from "lucide-react";
import { useMediaQuery } from "../../hooks/useMediaQuery";

export interface FilterDrawerProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  children: ReactNode;
  onApply?: () => void;
  onReset?: () => void;
  applyLabel?: string;
  resetLabel?: string;
}

export function FilterDrawer({
  isOpen,
  onClose,
  title,
  children,
  onApply,
  onReset,
  applyLabel = "Apply Filters",
  resetLabel = "Reset",
}: FilterDrawerProps) {
  const isMobile = useMediaQuery("(max-width: 768px)");

  // Prevent body scroll when drawer is open
  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = "hidden";
    } else {
      document.body.style.overflow = "";
    }

    return () => {
      document.body.style.overflow = "";
    };
  }, [isOpen]);

  // Close on Escape key
  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === "Escape" && isOpen) {
        onClose();
      }
    };

    document.addEventListener("keydown", handleEscape);
    return () => {
      document.removeEventListener("keydown", handleEscape);
    };
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  return (
    <>
      {/* Backdrop */}
      <div
        className="fixed inset-0 bg-black/50 z-40 transition-opacity"
        onClick={onClose}
        aria-hidden="true"
      />

      {/* Drawer */}
      <div
        className={`fixed z-50 bg-base-100 shadow-xl transition-transform duration-300 translate-y-0 ${
          isMobile
            ? "bottom-0 start-0 end-0 rounded-t-2xl max-h-[85vh] overflow-y-auto"
            : "top-0 end-0 h-full w-full max-w-md overflow-y-auto"
        }`}
        role="dialog"
        aria-modal="true"
        aria-labelledby="filter-drawer-title"
      >
        {/* Header */}
        <div className="sticky top-0 bg-base-100 border-b border-base-300 px-4 py-4 flex items-center justify-between z-10">
          <h2 id="filter-drawer-title" className="text-lg font-semibold">
            {title}
          </h2>
          <button
            type="button"
            onClick={onClose}
            className="btn btn-ghost btn-sm btn-circle"
            aria-label="Close filters"
          >
            <X size={20} />
          </button>
        </div>

        {/* Content */}
        <div className="px-4 py-6 space-y-6">{children}</div>

        {/* Footer Actions */}
        {(onApply !== undefined || onReset !== undefined) && (
          <div className="sticky bottom-0 bg-base-100 border-t border-base-300 px-4 py-4 flex gap-2">
            {onReset && (
              <button
                type="button"
                onClick={onReset}
                className="btn btn-ghost flex-1"
              >
                {resetLabel}
              </button>
            )}
            {onApply && (
              <button
                type="button"
                onClick={() => {
                  onApply();
                  onClose();
                }}
                className="btn btn-primary flex-1"
              >
                {applyLabel}
              </button>
            )}
          </div>
        )}
      </div>
    </>
  );
}
