import { memo, type SelectHTMLAttributes } from 'react';
import { type FieldError } from 'react-hook-form';

interface FormSelectProps extends Omit<SelectHTMLAttributes<HTMLSelectElement>, 'size'> {
  label?: string;
  error?: FieldError;
  size?: 'xs' | 'sm' | 'md' | 'lg';
  options: Array<{ value: string; label: string }>;
}

export const FormSelect = memo<FormSelectProps>(function FormSelect({
  label,
  error,
  size = 'sm',
  options,
  className = '',
  ...props
}) {
  const sizeClass = size === 'xs' ? 'select-xs' : size === 'sm' ? 'select-sm' : size === 'lg' ? 'select-lg' : 'select-md';
  
  return (
    <div className="form-control w-full">
      {label && (
        <label className="label pb-1">
          <span className="label-text text-sm font-medium text-neutral-700">{label}</span>
        </label>
      )}
      <select
        className={`select ${sizeClass} w-full border-neutral-200 bg-white focus:border-primary-500 ${
          error ? 'border-error focus:border-error' : ''
        } ${className}`}
        aria-invalid={error ? 'true' : 'false'}
        {...props}
      >
        {options.map((opt) => (
          <option key={opt.value} value={opt.value}>
            {opt.label}
          </option>
        ))}
      </select>
      {error && (
        <label className="label pt-1">
          <span className="label-text-alt text-error">{error.message}</span>
        </label>
      )}
    </div>
  );
});
