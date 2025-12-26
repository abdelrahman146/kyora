import { type HTMLAttributes } from 'react';
import { cn } from '@/lib/utils';

export interface SkeletonProps extends HTMLAttributes<HTMLDivElement> {
  variant?: 'text' | 'circular' | 'rectangular';
  width?: string | number;
  height?: string | number;
}

export const Skeleton = ({
  variant = 'rectangular',
  width,
  height,
  className,
  style,
  ...props
}: SkeletonProps) => {
  const baseClasses = 'animate-pulse bg-base-300';

  const variantClasses = {
    text: 'rounded h-4',
    circular: 'rounded-full',
    rectangular: 'rounded-md',
  };

  const inlineStyles: React.CSSProperties = {
    width: width ? (typeof width === 'number' ? `${String(width)}px` : width) : undefined,
    height: height ? (typeof height === 'number' ? `${String(height)}px` : height) : undefined,
    ...style,
  };

  return (
    <div
      className={cn(baseClasses, variantClasses[variant], className)}
      style={inlineStyles}
      {...props}
    />
  );
};

Skeleton.displayName = 'Skeleton';

export interface SkeletonTextProps {
  lines?: number;
  className?: string;
}

export const SkeletonText = ({ lines = 3, className }: SkeletonTextProps) => {
  return (
    <div className={cn('space-y-2', className)}>
      {Array.from({ length: lines }).map((_, index) => (
        <Skeleton
          key={index}
          variant="text"
          width={index === lines - 1 ? '70%' : '100%'}
        />
      ))}
    </div>
  );
};

SkeletonText.displayName = 'SkeletonText';
