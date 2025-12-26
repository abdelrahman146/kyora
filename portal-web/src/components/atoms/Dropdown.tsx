import { useState, useEffect, useRef, type ReactNode } from "react";
import { cn } from "@/lib/utils";

export interface DropdownProps {
  /** The trigger button element */
  trigger: ReactNode;
  /** The dropdown content */
  children: ReactNode;
  /** Additional classes for the dropdown container */
  className?: string;
  /** Additional classes for the content */
  contentClassName?: string;
  /** Alignment of dropdown */
  align?: "start" | "end";
  /** Width of dropdown content */
  width?: string;
}

/**
 * Dropdown Component
 *
 * A reusable dropdown component with proper click-outside detection and state management.
 *
 * Features:
 * - Click outside to close
 * - Keyboard support (Escape to close)
 * - Proper z-index stacking
 * - Customizable alignment
 * - RTL support
 * - Accessible
 *
 * @example
 * ```tsx
 * <Dropdown
 *   trigger={
 *     <button className="btn">Open Menu</button>
 *   }
 * >
 *   <ul className="menu">
 *     <li><a>Item 1</a></li>
 *     <li><a>Item 2</a></li>
 *   </ul>
 * </Dropdown>
 * ```
 */
export function Dropdown({
  trigger,
  children,
  className,
  contentClassName,
  align = "end",
  width = "auto",
}: DropdownProps) {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
      }
    };

    if (isOpen) {
      document.addEventListener("mousedown", handleClickOutside);
    }

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [isOpen]);

  // Close dropdown on Escape key
  useEffect(() => {
    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === "Escape" && isOpen) {
        setIsOpen(false);
      }
    };

    if (isOpen) {
      document.addEventListener("keydown", handleEscape);
    }

    return () => {
      document.removeEventListener("keydown", handleEscape);
    };
  }, [isOpen]);

  const handleToggle = () => {
    setIsOpen(!isOpen);
  };

  const handleClose = () => {
    setIsOpen(false);
  };

  return (
    <div ref={dropdownRef} className={cn("relative", className)}>
      {/* Trigger */}
      <div onClick={handleToggle} role="button" tabIndex={0}>
        {trigger}
      </div>

      {/* Dropdown Content */}
      {isOpen && (
        <div
          className={cn(
            "absolute top-full mt-2 z-50",
            "bg-base-100 rounded-lg shadow-lg border border-base-300",
            "animate-in fade-in slide-in-from-top-2 duration-200",
            align === "end" && "end-0",
            align === "start" && "start-0",
            contentClassName
          )}
          style={{ width }}
          onClick={handleClose}
        >
          {children}
        </div>
      )}
    </div>
  );
}
