import { memo, type InputHTMLAttributes } from 'react';
import { type FieldError } from 'react-hook-form';

interface FormInputProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'size'> {
  label?: string;
  error?: FieldError;
  size?: 'xs' | 'sm' | 'md' | 'lg';
}

export const FormInput = memo<FormInputProps>(function FormInput({
  label,
  error,
  size = 'sm',
  className = '',
  ...props
}) {
  const sizeClass = size === 'xs' ? 'input-xs' : size === 'sm' ? 'input-sm' : size === 'lg' ? 'input-lg' : 'input-md';
  
  return (
    <div className="form-control w-full">
      {label && (
        <label className="label pb-1">
          <span className="label-text text-sm font-medium text-neutral-700">{label}</span>
        </label>
      )}
      <input
        className={`input ${sizeClass} w-full border-neutral-200 bg-white focus:border-primary-500 ${
          error ? 'border-error focus:border-error' : ''
        } ${className}`}
        aria-invalid={error ? 'true' : 'false'}
        {...props}
      />
      {error && (
        <label className="label pt-1">
          <span className="label-text-alt text-error">{error.message}</span>
        </label>
      )}
    </div>
  );
});
