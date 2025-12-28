/**
 * FilterDrawer Component
 *
 * A specialized drawer component for filtering lists.
 * Built on top of the generic BottomSheet component.
 *
 * Features:
 * - Mobile: Bottom sheet with slide-up animation
 * - Desktop: Side drawer with slide-in animation
 * - Apply and Reset action buttons
 * - Accessible (keyboard navigation, ARIA attributes)
 * - RTL-compatible
 * - Responsive design
 */

import type { ReactNode } from "react";
import { BottomSheet } from "../molecules/BottomSheet";

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
  // Build footer with action buttons if provided
  const footer =
    onApply !== undefined || onReset !== undefined ? (
      <div className="flex gap-2">
        {onReset && (
          <button type="button" onClick={onReset} className="btn btn-ghost flex-1">
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
    ) : undefined;

  return (
    <BottomSheet
      isOpen={isOpen}
      onClose={onClose}
      title={title}
      footer={footer}
      side="end"
      size="md"
      contentClassName="space-y-6"
      ariaLabel="Filter options"
    >
      {children}
    </BottomSheet>
  );
}
