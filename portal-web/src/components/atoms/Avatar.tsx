import { forwardRef, type ImgHTMLAttributes } from "react";
import { cn } from "@/lib/utils";
import { Building2 } from "lucide-react";

export interface AvatarProps extends ImgHTMLAttributes<HTMLImageElement> {
  size?: "xs" | "sm" | "md" | "lg" | "xl";
  fallback?: string;
  shape?: "circle" | "square";
}

/**
 * Avatar Component
 *
 * Displays user or business avatar with automatic fallback.
 *
 * @example
 * ```tsx
 * <Avatar src={user.avatarUrl} alt={user.name} fallback="AB" />
 * <Avatar src={business.logoUrl} fallback={business.name[0]} shape="square" />
 * ```
 */
export const Avatar = forwardRef<HTMLImageElement, AvatarProps>(
  (
    {
      size = "md",
      fallback,
      shape = "circle",
      src,
      alt,
      className,
      ...props
    },
    ref
  ) => {
    const sizeClasses = {
      xs: "w-6 h-6 text-xs",
      sm: "w-8 h-8 text-sm",
      md: "w-10 h-10 text-base",
      lg: "w-12 h-12 text-lg",
      xl: "w-16 h-16 text-xl",
    };

    const shapeClasses = {
      circle: "rounded-full",
      square: "rounded-lg",
    };

    // If no src or src fails to load, show fallback
    if (!src) {
      return (
        <div
          className={cn(
            "flex items-center justify-center bg-primary-100 text-primary-700 font-semibold",
            sizeClasses[size],
            shapeClasses[shape],
            className
          )}
        >
          {fallback ? (
            <span>{fallback.charAt(0).toUpperCase()}</span>
          ) : (
            <Building2 size={16} />
          )}
        </div>
      );
    }

    return (
      <img
        ref={ref}
        src={src}
        alt={alt}
        className={cn(
          "object-cover",
          sizeClasses[size],
          shapeClasses[shape],
          className
        )}
        onError={(e) => {
          // Hide image and show fallback on error
          e.currentTarget.style.display = "none";
        }}
        {...props}
      />
    );
  }
);

Avatar.displayName = "Avatar";
