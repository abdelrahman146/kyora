/**
 * AsOfDatePicker - Date selector button for financial reports
 *
 * Displays the current "As of" date and opens a calendar picker directly when clicked.
 * Updates the URL search param to reflect the selected date.
 *
 * Features:
 * - Shows current date in localized format
 * - Opens calendar picker directly using DatePicker in buttonMode
 * - Modal on desktop, bottom sheet on mobile (no clipping)
 * - Updates URL with YYYY-MM-DD format
 * - Max date is today (can't view future reports)
 */
import { useNavigate } from '@tanstack/react-router'
import { format } from 'date-fns'
import { Calendar } from 'lucide-react'
import { useCallback } from 'react'
import { useTranslation } from 'react-i18next'

import { DatePicker } from '@/components/form'
import { formatDateShort } from '@/lib/formatDate'

export interface AsOfDatePickerProps {
  /** Current date value from search params */
  asOf?: string
  /** Route path to navigate to with the new date */
  routeTo: string
  /** Business descriptor for route params */
  businessDescriptor: string
  /** Additional class names */
  className?: string
}

export function AsOfDatePicker({
  asOf,
  routeTo,
  businessDescriptor,
  className,
}: AsOfDatePickerProps) {
  const { t } = useTranslation('reports')
  const navigate = useNavigate()

  // Parse current date from search param or use today
  const currentDate = asOf ? new Date(asOf) : new Date()

  // Handle date selection
  const handleDateChange = useCallback(
    (date: Date | undefined) => {
      if (!date) return

      // Format date as YYYY-MM-DD for URL
      const formattedDate = format(date, 'yyyy-MM-dd')

      // Navigate to the same page with new date
      void navigate({
        to: routeTo,
        params: { businessDescriptor },
        search: { asOf: formattedDate },
      })
    },
    [navigate, routeTo, businessDescriptor],
  )

  return (
    <DatePicker
      buttonMode
      value={currentDate}
      onChange={handleDateChange}
      maxDate={new Date()}
      clearable={false}
      label={t('hub.as_of', { date: formatDateShort(currentDate) })}
      buttonContent={
        <>
          <Calendar className="h-4 w-4" aria-hidden="true" />
          <span className="whitespace-nowrap">
            {t('hub.as_of', { date: formatDateShort(currentDate) })}
          </span>
        </>
      }
      className={className}
    />
  )
}

AsOfDatePicker.displayName = 'AsOfDatePicker'
