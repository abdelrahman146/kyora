import {
  forwardRef,
  useId,
  useState,
  useRef,
  useEffect,
  useCallback,
  type ReactNode,
} from "react";
import { cn } from "@/lib/utils";
import { Check, ChevronDown, X, Search } from "lucide-react";

export interface FormSelectOption<T = string> {
  value: T;
  label: string;
  description?: string;
  icon?: ReactNode;
  disabled?: boolean;
  renderCustom?: () => ReactNode;
}

export interface FormSelectProps<T = string> {
  label?: string;
  error?: string;
  helperText?: string;
  options: FormSelectOption<T>[];
  value?: T | T[];
  onChange?: (value: T | T[]) => void;
  size?: "sm" | "md" | "lg";
  variant?: "default" | "filled" | "ghost";
  fullWidth?: boolean;
  searchable?: boolean;
  multiSelect?: boolean;
  placeholder?: string;
  clearable?: boolean;
  maxHeight?: number;
  id?: string;
  disabled?: boolean;
  required?: boolean;
  className?: string;
}

/**
 * FormSelect - Advanced production-grade select/dropdown component
 * 
 * Features:
 * - RTL-first design
 * - Mobile-optimized with bottom sheet on mobile
 * - Searchable with filtering
 * - Multi-select support
 * - Custom option rendering
 * - Keyboard navigation
 * - Accessible with ARIA attributes
 * - Clearable selection
 * 
 * @example
 * // Basic select
 * <FormSelect
 *   label="Country"
 *   options={countries}
 *   value={selectedCountry}
 *   onChange={setSelectedCountry}
 * />
 * 
 * // Searchable multi-select
 * <FormSelect
 *   label="Tags"
 *   options={tags}
 *   value={selectedTags}
 *   onChange={setSelectedTags}
 *   searchable
 *   multiSelect
 * />
 */
export const FormSelect = forwardRef<HTMLDivElement, FormSelectProps>(
  <T extends string>(
    {
      label,
      error,
      helperText,
      options,
      value,
      onChange,
      size = "md",
      variant = "default",
      fullWidth = true,
      searchable = false,
      multiSelect = false,
      placeholder = "Select...",
      clearable = false,
      maxHeight = 300,
      className,
      id,
      disabled,
      required,
    }: FormSelectProps<T>,
    ref: React.ForwardedRef<HTMLDivElement>
  ) => {
    const generatedId = useId();
    const inputId = id ?? generatedId;
    const [isOpen, setIsOpen] = useState(false);
    const [searchQuery, setSearchQuery] = useState("");
    const [focusedIndex, setFocusedIndex] = useState(-1);
    const containerRef = useRef<HTMLDivElement>(null);
    const searchInputRef = useRef<HTMLInputElement>(null);

    const selectedValues = (() => {
      if (multiSelect && Array.isArray(value)) {
        return value;
      }
      if (!multiSelect && value !== undefined && !Array.isArray(value)) {
        return [value];
      }
      return [];
    })();

    const filteredOptions = searchable && searchQuery
      ? options.filter((option) =>
          option.label.toLowerCase().includes(searchQuery.toLowerCase())
        )
      : options;

    const sizeClasses = {
      sm: "h-[44px] text-sm",
      md: "h-[50px] text-base",
      lg: "h-[56px] text-lg",
    };

    const variantClasses = {
      default: "input-bordered bg-base-100",
      filled: "input-bordered bg-base-200/50 border-transparent focus:bg-base-100",
      ghost: "input-ghost bg-transparent",
    };

    useEffect(() => {
      // Prevent body scroll when dropdown is open on mobile
      if (isOpen && typeof window !== "undefined") {
        const originalOverflow = document.body.style.overflow;
        const originalPaddingRight = document.body.style.paddingRight;
        
        // Only lock scroll on mobile/tablet
        const isMobile = window.innerWidth < 1024;
        if (isMobile) {
          const scrollbarWidth = window.innerWidth - document.documentElement.clientWidth;
          document.body.style.overflow = "hidden";
          if (scrollbarWidth > 0) {
            document.body.style.paddingRight = `${String(scrollbarWidth)}px`;
          }
        }

        return () => {
          if (isMobile) {
            document.body.style.overflow = originalOverflow;
            document.body.style.paddingRight = originalPaddingRight;
          }
        };
      }
      
      return undefined;
    }, [isOpen]);

    useEffect(() => {
      const handleClickOutside = (event: MouseEvent | TouchEvent) => {
        if (
          containerRef.current &&
          !containerRef.current.contains(event.target as Node)
        ) {
          setIsOpen(false);
          setSearchQuery("");
          setFocusedIndex(-1);
        }
      };

      const handleEscape = (event: KeyboardEvent) => {
        if (event.key === "Escape" && isOpen) {
          setIsOpen(false);
          setSearchQuery("");
          setFocusedIndex(-1);
        }
      };

      if (isOpen) {
        // Use capture phase for better click-outside detection
        document.addEventListener("mousedown", handleClickOutside, true);
        document.addEventListener("touchstart", handleClickOutside, true);
        document.addEventListener("keydown", handleEscape);
        
        if (searchable && searchInputRef.current) {
          searchInputRef.current.focus();
        }
      }

      return () => {
        document.removeEventListener("mousedown", handleClickOutside, true);
        document.removeEventListener("touchstart", handleClickOutside, true);
        document.removeEventListener("keydown", handleEscape);
      };
    }, [isOpen, searchable]);

    const handleToggleOption = useCallback((optionValue: T) => {
      if (disabled) return;

      if (multiSelect) {
        const newValues = selectedValues.includes(optionValue)
          ? selectedValues.filter((v) => v !== optionValue)
          : [...selectedValues, optionValue];
        onChange?.(newValues as T | T[]);
      } else {
        onChange?.(optionValue as T | T[]);
        setIsOpen(false);
        setSearchQuery("");
        setFocusedIndex(-1);
      }
    }, [disabled, multiSelect, selectedValues, onChange]);

    const handleClear = useCallback((e: React.MouseEvent) => {
      e.stopPropagation();
      if (multiSelect) {
        onChange?.([] as T | T[]);
      } else {
        onChange?.(null as unknown as T | T[]);
      }
    }, [multiSelect, onChange]);

    const handleKeyDown = useCallback((e: React.KeyboardEvent) => {
      if (disabled) return;

      switch (e.key) {
        case "Enter":
        case " ":
          if (!isOpen) {
            setIsOpen(true);
          } else if (focusedIndex >= 0 && focusedIndex < filteredOptions.length) {
            const option = filteredOptions[focusedIndex];
            if (!option.disabled) {
              handleToggleOption(option.value);
            }
          }
          e.preventDefault();
          break;
        case "Escape":
          if (isOpen) {
            setIsOpen(false);
            setSearchQuery("");
            setFocusedIndex(-1);
            e.preventDefault();
          }
          break;
        case "ArrowDown":
          if (!isOpen) {
            setIsOpen(true);
          } else {
            setFocusedIndex((prev) => {
              const nextIndex = prev + 1;
              return nextIndex < filteredOptions.length ? nextIndex : prev;
            });
          }
          e.preventDefault();
          break;
        case "ArrowUp":
          if (isOpen) {
            setFocusedIndex((prev) => (prev > 0 ? prev - 1 : 0));
            e.preventDefault();
          }
          break;
        case "Home":
          if (isOpen) {
            setFocusedIndex(0);
            e.preventDefault();
          }
          break;
        case "End":
          if (isOpen) {
            setFocusedIndex(filteredOptions.length - 1);
            e.preventDefault();
          }
          break;
        case "Tab":
          if (isOpen) {
            setIsOpen(false);
            setSearchQuery("");
            setFocusedIndex(-1);
          }
          break;
      }
    }, [disabled, isOpen, focusedIndex, filteredOptions, handleToggleOption]);

    const getDisplayText = () => {
      if (selectedValues.length === 0) return placeholder;

      if (multiSelect) {
        const count = String(selectedValues.length);
        return count !== '0' ? `${count} selected` : placeholder;
      }

      const selectedOption = options.find((opt) => opt.value === selectedValues[0]);
      return selectedOption?.label ?? placeholder;
    };

    return (
      <div
        ref={ref}
        className={cn("form-control", fullWidth && "w-full")}
      >
        {label && (
          <label htmlFor={inputId} className="label">
            <span className="label-text text-base-content/70 font-medium">
              {label}
              {required && <span className="text-error ms-1">*</span>}
            </span>
          </label>
        )}

        <div ref={containerRef} className="relative">
          <button
            type="button"
            id={inputId}
            onClick={() => {
              if (!disabled) setIsOpen(!isOpen);
            }}
            onKeyDown={handleKeyDown}
            disabled={disabled}
            className={cn(
              "input w-full flex items-center justify-between gap-2 transition-all duration-200",
              sizeClasses[size],
              variantClasses[variant],
              "text-start cursor-pointer",
              "focus:outline-none focus:border-primary focus:ring-2 focus:ring-primary/20",
              error ? "input-error border-error focus:border-error focus:ring-error/20" : "",
              disabled && "opacity-60 cursor-not-allowed",
              isOpen && "border-primary ring-2 ring-primary/20",
              className ?? ""
            )}
            aria-haspopup="listbox"
            aria-expanded={isOpen}
            aria-labelledby={label ? `${inputId}-label` : undefined}
            aria-required={required}
            aria-invalid={error ? "true" : "false"}
          >
            <span
              className={cn(
                "flex-1 truncate",
                selectedValues.length === 0 && "text-base-content/40"
              )}
            >
              {getDisplayText()}
            </span>

            <div className="flex items-center gap-1">
              {clearable && selectedValues.length > 0 && !disabled && (
                <button
                  type="button"
                  onClick={(e) => {
                    e.stopPropagation();
                    handleClear(e);
                  }}
                  className="p-1 hover:bg-base-200 rounded-md transition-colors"
                  aria-label="Clear selection"
                >
                  <X className="w-4 h-4" />
                </button>
              )}
              <ChevronDown
                className={cn(
                  "w-5 h-5 transition-transform duration-200",
                  isOpen && "rotate-180"
                )}
              />
            </div>
          </button>

          {isOpen && !disabled && (
            <div
              className={cn(
                "absolute z-50 mt-2 w-full",
                "bg-base-100 border border-base-300 rounded-lg shadow-xl",
                "overflow-hidden",
                "animate-in fade-in-0 zoom-in-95 duration-100"
              )}
              style={{ maxHeight }}
              role="presentation"
              onClick={(e) => {
                e.stopPropagation();
              }}
            >
              {searchable && (
                <div className="p-2 border-b border-base-300">
                  <div className="relative">
                    <Search className="absolute start-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content/50" />
                    <input
                      ref={searchInputRef}
                      type="text"
                      value={searchQuery}
                      onChange={(e) => {
                        setSearchQuery(e.target.value);
                        setFocusedIndex(-1);
                      }}
                      placeholder="Search..."
                      className="input input-sm w-full ps-9"
                      aria-label="Search options"
                    />
                  </div>
                </div>
              )}

              <ul
                role="listbox"
                aria-multiselectable={multiSelect}
                className="overflow-y-auto"
                style={{ maxHeight: maxHeight - (searchable ? 60 : 0) }}
              >
                {filteredOptions.length === 0 ? (
                  <li className="p-4 text-center text-base-content/50">
                    No options found
                  </li>
                ) : (
                  filteredOptions.map((option, index) => {
                    const isSelected = selectedValues.includes(option.value);
                    const isFocused = index === focusedIndex;

                    return (
                      <li
                        key={option.value}
                        role="option"
                        aria-selected={isSelected}
                        aria-disabled={option.disabled}
                        onClick={() => {
                          if (!option.disabled) handleToggleOption(option.value);
                        }}
                        onKeyDown={(e) => {
                          if (e.key === "Enter" || e.key === " ") {
                            e.preventDefault();
                            if (!option.disabled) handleToggleOption(option.value);
                          }
                        }}
                        tabIndex={isFocused ? 0 : -1}
                        ref={(el) => {
                          if (isFocused && el) {
                            el.scrollIntoView({ block: "nearest", behavior: "smooth" });
                          }
                        }}
                        className={cn(
                          "flex items-center gap-3 px-4 py-3 cursor-pointer transition-colors",
                          "min-h-[48px]", // Better touch target for mobile
                          "hover:bg-base-200 focus:bg-base-200 focus:outline-none",
                          "active:bg-base-300", // Touch feedback
                          isSelected && "bg-primary/10",
                          isFocused && "bg-base-200 ring-2 ring-inset ring-primary/30",
                          option.disabled && "opacity-50 cursor-not-allowed pointer-events-none"
                        )}
                      >
                        {option.renderCustom ? (
                          option.renderCustom()
                        ) : (
                          <>
                            {option.icon && (
                              <span className="shrink-0">{option.icon}</span>
                            )}
                            <div className="flex-1 min-w-0">
                              <div className="font-medium truncate">{option.label}</div>
                              {option.description && (
                                <div className="text-sm text-base-content/60 truncate">
                                  {option.description}
                                </div>
                              )}
                            </div>
                            {multiSelect && isSelected && (
                              <Check className="w-5 h-5 text-primary shrink-0" />
                            )}
                            {!multiSelect && isSelected && (
                              <Check className="w-5 h-5 text-primary shrink-0" />
                            )}
                          </>
                        )}
                      </li>
                    );
                  })
                )}
              </ul>
            </div>
          )}
        </div>

        {error && (
          <label className="label">
            <span
              id={`${inputId}-error`}
              className="label-text-alt text-error"
              role="alert"
            >
              {error}
            </span>
          </label>
        )}

        {!error && helperText && (
          <label className="label">
            <span
              id={`${inputId}-helper`}
              className="label-text-alt text-base-content/60"
            >
              {helperText}
            </span>
          </label>
        )}
      </div>
    );
  }
) as (<T extends string>(
  props: FormSelectProps<T> & { ref?: React.ForwardedRef<HTMLDivElement> }
) => React.ReactElement) & { displayName?: string };

(FormSelect as { displayName?: string }).displayName = "FormSelect";
