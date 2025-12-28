/**
 * SearchInput Component
 *
 * A debounced search input component for filtering lists.
 *
 * Features:
 * - Automatic debouncing (default 300ms)
 * - Clear button when input has value
 * - Loading indicator during debounce
 * - Accessible keyboard navigation
 * - RTL-compatible
 * - Mobile-friendly (min 44px touch target)
 */

import { useState, useEffect, useRef } from "react";
import { Search, X } from "lucide-react";

export interface SearchInputProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  debounceMs?: number;
  disabled?: boolean;
  className?: string;
}

export function SearchInput({
  value,
  onChange,
  placeholder = "Search...",
  debounceMs = 300,
  disabled = false,
  className = "",
}: SearchInputProps) {
  const [localValue, setLocalValue] = useState(value);
  const [isDebouncePending, setIsDebouncePending] = useState(false);
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);

  // Sync external value changes
  useEffect(() => {
    setLocalValue(value);
  }, [value]);

  // Debounced onChange
  useEffect(() => {
    if (localValue === value) {
      setIsDebouncePending(false);
      return;
    }

    setIsDebouncePending(true);

    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }

    timeoutRef.current = setTimeout(() => {
      onChange(localValue);
      setIsDebouncePending(false);
    }, debounceMs);

    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, [localValue, value, onChange, debounceMs]);

  const handleClear = () => {
    setLocalValue("");
    onChange("");
  };

  return (
    <div className={`relative ${className}`}>
      {/* Search Icon */}
      <Search
        size={20}
        className="absolute top-1/2 -translate-y-1/2 start-3 text-base-content/50 pointer-events-none"
      />

      {/* Input */}
      <input
        type="text"
        value={localValue}
        onChange={(e) => {
          setLocalValue(e.target.value);
        }}
        placeholder={placeholder}
        disabled={disabled}
        className="input input-bordered w-full ps-10 pe-10"
        aria-label={placeholder}
      />

      {/* Clear Button / Loading Indicator */}
      {localValue && (
        <button
          type="button"
          onClick={handleClear}
          disabled={disabled}
          className="btn btn-ghost btn-sm btn-circle absolute top-1/2 -translate-y-1/2 end-2"
          aria-label="Clear search"
        >
          {isDebouncePending ? (
            <span className="loading loading-spinner loading-sm"></span>
          ) : (
            <X size={18} />
          )}
        </button>
      )}
    </div>
  );
}
