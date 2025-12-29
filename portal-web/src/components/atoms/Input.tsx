import { forwardRef, useId, type InputHTMLAttributes } from 'react';
import { cn } from '@/lib/utils';

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  helperText?: string;
  startIcon?: React.ReactNode;
  endIcon?: React.ReactNode;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  (
    {
      label,
      error,
      helperText,
      startIcon,
      endIcon,
      className,
      id,
      ...props
    },
    ref
  ) => {
    const generatedId = useId();
    const inputId = id ?? generatedId;

    return (
      <div className="w-full">
        {label && (
          <label
            htmlFor={inputId}
            className="label block mb-2 text-sm font-medium text-base-content"
          >
            <span className="label-text">{label}</span>
          </label>
        )}
        <div className="relative">
          {startIcon && (
            <div className="absolute inset-y-0 start-0 z-10 flex items-center ps-3 pointer-events-none text-base-content/50">
              <span aria-hidden="true">{startIcon}</span>
            </div>
          )}
          <input
            ref={ref}
            id={inputId}
            className={cn(
              'input input-bordered w-full h-[50px] relative z-0',
              'bg-base-100 text-base-content',
              'text-start placeholder:text-base-content/40',
              'focus:outline-none focus:border-primary focus:ring-2 focus:ring-primary/20',
              'transition-all duration-200',
              error ? 'input-error border-error focus:border-error focus:ring-error/20' : '',
              startIcon ? 'ps-10' : '',
              endIcon ? 'pe-10' : '',
              className
            )}
            aria-invalid={error ? 'true' : 'false'}
            aria-describedby={
              error ? `${inputId}-error` : helperText ? `${inputId}-helper` : undefined
            }
            {...props}
          />
          {endIcon && (
            <div className="absolute inset-y-0 end-0 z-10 flex items-center pe-3 pointer-events-none text-base-content/50">
              <span aria-hidden="true">{endIcon}</span>
            </div>
          )}
        </div>
        {error && (
          <p
            id={`${inputId}-error`}
            className="mt-2 text-sm text-error"
            role="alert"
          >
            {error}
          </p>
        )}
        {!error && helperText && (
          <p
            id={`${inputId}-helper`}
            className="mt-2 text-sm text-neutral-500"
          >
            {helperText}
          </p>
        )}
      </div>
    );
  }
);

Input.displayName = 'Input';
