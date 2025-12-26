import { type HTMLAttributes } from 'react';
import { cn } from '@/lib/utils';

export interface BadgeProps extends HTMLAttributes<HTMLSpanElement> {
  variant?: 'default' | 'primary' | 'secondary' | 'success' | 'warning' | 'error' | 'info';
  size?: 'sm' | 'md' | 'lg';
}

export const Badge = ({
  variant = 'default',
  size = 'md',
  className,
  children,
  ...props
}: BadgeProps) => {
  const baseClasses = 'badge inline-flex items-center justify-center font-medium';

  const variantClasses = {
    default: 'badge-neutral',
    primary: 'badge-primary',
    secondary: 'badge-secondary',
    success: 'badge-success',
    warning: 'badge-warning',
    error: 'badge-error',
    info: 'badge-info',
  };

  const sizeClasses = {
    sm: 'badge-sm text-xs px-2 py-0.5',
    md: 'badge-md text-sm px-3 py-1',
    lg: 'badge-lg text-base px-4 py-1.5',
  };

  return (
    <span
      className={cn(
        baseClasses,
        variantClasses[variant],
        sizeClasses[size],
        'rounded-sm',
        className
      )}
      {...props}
    >
      {children}
    </span>
  );
};

Badge.displayName = 'Badge';
