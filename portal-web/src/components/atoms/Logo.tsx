import { type HTMLAttributes } from "react";
import { cn } from "@/lib/utils";

export interface LogoProps extends HTMLAttributes<HTMLDivElement> {
  size?: "sm" | "md" | "lg";
  showText?: boolean;
}

/**
 * Logo Component
 *
 * Displays the Kyora brand logo with optional text.
 * RTL-aware and responsive.
 *
 * @example
 * ```tsx
 * <Logo size="md" showText />
 * <Logo size="sm" />
 * ```
 */
export function Logo({ size = "md", showText = true, className }: LogoProps) {
  const sizeClasses = {
    sm: "h-6",
    md: "h-8",
    lg: "h-10",
  };

  const textSizeClasses = {
    sm: "text-lg",
    md: "text-xl",
    lg: "text-2xl",
  };

  return (
    <div className={cn("flex items-center gap-2", className)}>
      {/* Logo Icon - Using a simple geometric shape for now */}
      <div
        className={cn(
          "bg-primary-600 rounded-lg flex items-center justify-center aspect-square",
          sizeClasses[size]
        )}
      >
        <span className="text-white font-bold text-sm">K</span>
      </div>

      {/* Logo Text */}
      {showText && (
        <span
          className={cn(
            "font-bold text-primary-900",
            textSizeClasses[size]
          )}
        >
          Kyora
        </span>
      )}
    </div>
  );
}
