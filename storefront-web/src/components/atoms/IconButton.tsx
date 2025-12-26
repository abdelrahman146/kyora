import { memo, type ButtonHTMLAttributes, type ReactNode } from 'react';

interface IconButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  icon: ReactNode;
  label: string;
  size?: 'sm' | 'md' | 'lg';
  variant?: 'ghost' | 'primary' | 'secondary';
}

/**
 * IconButton Atom - Icon-only button with accessibility
 * Memoized to prevent unnecessary re-renders
 */
export const IconButton = memo<IconButtonProps>(function IconButton({
  icon,
  label,
  size = 'md',
  variant = 'ghost',
  className = '',
  ...props
}) {
  const sizeClasses = {
    sm: 'btn-sm',
    md: 'btn-md',
    lg: 'btn-lg',
  };

  const variantClasses = {
    ghost: 'btn-ghost',
    primary: 'btn-primary',
    secondary: 'btn-secondary',
  };

  const classes = [
    'btn btn-square active-scale focus-ring',
    sizeClasses[size],
    variantClasses[variant],
    className,
  ]
    .filter(Boolean)
    .join(' ');

  return (
    <button
      type="button"
      className={classes}
      aria-label={label}
      {...props}
    >
      {icon}
    </button>
  );
});
