import { useEffect, useRef, useState } from 'react';

import { useTranslation } from 'react-i18next';
import { Calendar, X } from 'lucide-react';
import { DayPicker } from 'react-day-picker';
import 'react-day-picker/dist/style.css';
import { format } from 'date-fns';
import { ar } from 'date-fns/locale';

import type { InputHTMLAttributes } from 'react';
import type { DateRange } from 'react-day-picker';

import { cn } from '@/lib/utils';
import { useLanguage } from '@/hooks/useLanguage';

type Size = 'sm' | 'md' | 'lg';

interface DateRangePickerProps
  extends Omit<
    InputHTMLAttributes<HTMLInputElement>,
    'value' | 'onChange' | 'size'
  > {
  value?: DateRange;
  onChange: (range?: DateRange) => void;
  minDate?: Date;
  maxDate?: Date;
  disabledDates?: Array<Date>;
  error?: string;
  size?: Size;
  numberOfMonths?: number;
}

export function DateRangePicker({
  value,
  onChange,
  minDate,
  maxDate,
  disabledDates = [],
  error,
  size = 'md',
  numberOfMonths = 2,
  placeholder,
  disabled,
  required,
  ...props
}: DateRangePickerProps) {
  const { t } = useTranslation();
  const { isRTL } = useLanguage();
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
      }
    };

    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside);
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [isOpen]);

  // Format date range for display
  const formatDateRange = (range?: DateRange): string => {
    if (!range?.from) return '';
    if (!range.to) {
      return format(range.from, 'PP', { locale: isRTL ? ar : undefined });
    }
    const fromStr = format(range.from, 'PP', { locale: isRTL ? ar : undefined });
    const toStr = format(range.to, 'PP', { locale: isRTL ? ar : undefined });
    return `${fromStr} - ${toStr}`;
  };

  const handleClear = (e: React.MouseEvent) => {
    e.stopPropagation();
    onChange(undefined);
  };

  const handleSelect = (range: DateRange | undefined) => {
    onChange(range);
    // Close dropdown when both dates are selected
    if (range?.from && range.to) {
      setIsOpen(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Escape') {
      setIsOpen(false);
      inputRef.current?.focus();
    } else if (e.key === 'Enter' && !isOpen) {
      e.preventDefault();
      setIsOpen(true);
    }
  };

  const sizeClasses: Record<Size, string> = {
    sm: 'input-sm',
    md: 'input-md',
    lg: 'input-lg',
  };

  return (
    <div className="relative w-full" ref={dropdownRef}>
      <div className="relative">
        <input
          ref={inputRef}
          type="text"
          readOnly
          value={formatDateRange(value)}
          placeholder={placeholder || t('common.selectDateRange')}
          disabled={disabled}
          required={required}
          onClick={() => !disabled && setIsOpen(!isOpen)}
          onKeyDown={handleKeyDown}
          className={cn(
            'input input-bordered w-full',
            sizeClasses[size],
            error && 'input-error',
            disabled && 'input-disabled',
            'pe-20 cursor-pointer'
          )}
          aria-invalid={!!error}
          aria-describedby={error ? `${props.id}-error` : undefined}
          {...props}
        />

        <div className="absolute inset-y-0 end-0 flex items-center gap-1 pe-3">
          {value?.from && !disabled && (
            <button
              type="button"
              onClick={handleClear}
              className="btn btn-ghost btn-xs btn-circle"
              aria-label={t('common.clear')}
              tabIndex={-1}
            >
              <X size={16} />
            </button>
          )}
          <Calendar
            size={18}
            className={cn(
              'text-base-content/50',
              disabled && 'text-base-content/30'
            )}
            aria-hidden="true"
          />
        </div>
      </div>

      {isOpen && (
        <div
          className={cn(
            'dropdown-content bg-base-100 rounded-box shadow-lg border border-base-300 p-4 z-50',
            'absolute mt-2',
            isRTL ? 'end-0' : 'start-0'
          )}
          role="dialog"
          aria-label={t('common.selectDateRange')}
          onKeyDown={(e) => {
            if (e.key === 'Escape') {
              setIsOpen(false);
              inputRef.current?.focus();
            }
          }}
        >
          <DayPicker
            mode="range"
            selected={value}
            onSelect={handleSelect}
            disabled={[
              ...(minDate ? [{ before: minDate }] : []),
              ...(maxDate ? [{ after: maxDate }] : []),
              ...(disabledDates.length > 0
                ? disabledDates.map((date) => ({
                    from: date,
                    to: date,
                  }))
                : []),
            ]}
            numberOfMonths={numberOfMonths}
            locale={isRTL ? ar : undefined}
            dir={isRTL ? 'rtl' : 'ltr'}
            className={cn('rdp-custom')}
            classNames={{
              months: 'flex gap-4 flex-col sm:flex-row',
              month: 'space-y-4',
              caption: 'flex justify-center pt-1 relative items-center',
              caption_label: 'text-base font-medium',
              nav: 'space-x-1 flex items-center',
              nav_button: cn(
                'btn btn-ghost btn-sm btn-circle',
                'hover:bg-base-200'
              ),
              nav_button_previous: cn(
                'absolute',
                isRTL ? 'end-1' : 'start-1'
              ),
              nav_button_next: cn('absolute', isRTL ? 'start-1' : 'end-1'),
              table: 'w-full border-collapse space-y-1',
              head_row: 'flex',
              head_cell: 'text-base-content/60 rounded-md w-9 font-normal text-sm',
              row: 'flex w-full mt-2',
              cell: cn(
                'relative p-0 text-center text-sm focus-within:relative focus-within:z-20',
                'h-9 w-9'
              ),
              day: cn(
                'btn btn-ghost btn-sm h-9 w-9 p-0 font-normal',
                'hover:bg-base-200'
              ),
              day_range_start: 'bg-primary text-primary-content hover:bg-primary/90',
              day_range_end: 'bg-primary text-primary-content hover:bg-primary/90',
              day_selected:
                'bg-primary text-primary-content hover:bg-primary/90 focus:bg-primary/90',
              day_today: 'border border-primary',
              day_outside: 'text-base-content/30 opacity-50',
              day_disabled: 'text-base-content/30 opacity-50 cursor-not-allowed',
              day_range_middle:
                'bg-primary/20 text-primary-content hover:bg-primary/30',
              day_hidden: 'invisible',
            }}
          />

          <div className="flex gap-2 mt-4 pt-4 border-t border-base-300">
            <button
              type="button"
              className="btn btn-sm btn-ghost flex-1"
              onClick={() => {
                onChange(undefined);
                setIsOpen(false);
              }}
            >
              {t('common.clear')}
            </button>
            <button
              type="button"
              className="btn btn-sm btn-primary flex-1"
              onClick={() => setIsOpen(false)}
              disabled={!value?.from || !value.to}
            >
              {t('common.apply')}
            </button>
          </div>
        </div>
      )}

      {error && (
        <p className="text-error text-sm mt-1" id={`${props.id}-error`}>
          {error}
        </p>
      )}
    </div>
  );
}
