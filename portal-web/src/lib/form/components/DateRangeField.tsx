import { useTranslation } from 'react-i18next';
import type { DateRange } from 'react-day-picker';

import { DateRangePicker } from '@/components/atoms/DateRangePicker';
import { useFieldContext } from '@/lib/form';

interface DateRangeFieldProps {
  label?: string;
  placeholder?: string;
  minDate?: Date;
  maxDate?: Date;
  disabledDates?: Array<Date>;
  required?: boolean;
  disabled?: boolean;
  helperText?: string;
  size?: 'sm' | 'md' | 'lg';
  numberOfMonths?: number;
}

export function DateRangeField({
  label,
  placeholder,
  minDate,
  maxDate,
  disabledDates,
  required,
  disabled,
  helperText,
  size = 'md',
  numberOfMonths = 2,
}: DateRangeFieldProps) {
  const field = useFieldContext<DateRange | undefined>();
  const { t } = useTranslation();

  const errorMessage = field.state.meta.errors[0];

  return (
    <div className="w-full">
      {label && (
        <label className="label" htmlFor={field.name}>
          <span className="label-text">
            {label}
            {required && <span className="text-error ms-1">*</span>}
          </span>
        </label>
      )}

      <DateRangePicker
        id={field.name}
        value={field.state.value}
        onChange={(range: DateRange | undefined) => field.handleChange(range)}
        onBlur={field.handleBlur}
        minDate={minDate}
        maxDate={maxDate}
        disabledDates={disabledDates}
        placeholder={placeholder}
        required={required}
        disabled={disabled}
        error={errorMessage}
        size={size}
        numberOfMonths={numberOfMonths}
      />

      {helperText && !errorMessage && (
        <p className="text-base-content/60 text-sm mt-1">{helperText}</p>
      )}

      {errorMessage && (
        <p className="text-error text-sm mt-1" role="alert">
          {t(`errors.${errorMessage}`, { defaultValue: errorMessage })}
        </p>
      )}
    </div>
  );
}
