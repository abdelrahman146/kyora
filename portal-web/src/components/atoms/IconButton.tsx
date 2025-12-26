import { forwardRef, type ButtonHTMLAttributes } from "react";
import { Loader2, type LucideIcon } from "lucide-react";
import { cn } from "@/lib/utils";

export interface IconButtonProps
  extends ButtonHTMLAttributes<HTMLButtonElement> {
  icon: LucideIcon;
  size?: "sm" | "md" | "lg";
  variant?: "ghost" | "outline" | "primary";
  loading?: boolean;
}

/**
 * IconButton Component
 *
 * Button with only an icon, commonly used in headers and toolbars.
 * Provides consistent touch target size (44x44px minimum).
 *
 * @example
 * ```tsx
 * <IconButton icon={Menu} onClick={toggleSidebar} aria-label="Toggle menu" />
 * <IconButton icon={Bell} variant="outline" aria-label="Notifications" />
 * ```
 */
export const IconButton = forwardRef<HTMLButtonElement, IconButtonProps>(
  (
    {
      icon: Icon,
      size = "md",
      variant = "ghost",
      loading = false,
      disabled,
      className,
      ...props
    },
    ref
  ) => {
    const sizeClasses = {
      sm: "h-8 w-8",
      md: "h-10 w-10",
      lg: "h-12 w-12",
    };

    const iconSizeMap = {
      sm: 16,
      md: 20,
      lg: 24,
    };

    const variantClasses = {
      ghost: "btn-ghost",
      outline: "btn-outline",
      primary: "btn-primary",
    };

    return (
      <button
        ref={ref}
        disabled={disabled ?? loading}
        className={cn(
          "btn btn-square",
          sizeClasses[size],
          variantClasses[variant],
          "active:scale-95 transition-transform",
          className
        )}
        {...props}
      >
        {loading ? (
          <Loader2 size={iconSizeMap[size]} className="animate-spin" />
        ) : (
          <Icon size={iconSizeMap[size]} />
        )}
      </button>
    );
  }
);

IconButton.displayName = "IconButton";
