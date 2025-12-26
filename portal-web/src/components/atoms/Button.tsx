import { forwardRef, type ButtonHTMLAttributes } from 'react';
import { Loader2 } from 'lucide-react';
import { cn } from '@/lib/utils';

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'ghost' | 'outline';
  size?: 'sm' | 'md' | 'lg';
  loading?: boolean;
  fullWidth?: boolean;
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  (
    {
      variant = 'primary',
      size = 'md',
      loading = false,
      fullWidth = false,
      disabled,
      children,
      className,
      ...props
    },
    ref
  ) => {
    const baseClasses =
      'btn inline-flex items-center justify-center gap-2 font-semibold transition-all active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed';

    const variantClasses = {
      primary: 'btn-primary',
      secondary: 'bg-primary-50 text-primary-700 hover:bg-primary-100',
      ghost: 'btn-ghost',
      outline: 'btn-outline',
    };

    const sizeClasses = {
      sm: 'btn-sm h-10 px-4 text-sm',
      md: 'btn-md h-[52px] px-6 text-base',
      lg: 'btn-lg h-14 px-8 text-lg',
    };

    return (
      <button
        ref={ref}
        disabled={disabled || loading}
        className={cn(
          baseClasses,
          variantClasses[variant],
          sizeClasses[size],
          fullWidth && 'btn-block w-full',
          'rounded-xl',
          className
        )}
        {...props}
      >
        {loading && <Loader2 className="animate-spin" size={18} />}
        {children}
      </button>
    );
  }
);

Button.displayName = 'Button';
