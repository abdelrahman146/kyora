import { AlertCircle, CheckCircle2, Clock, Loader2 } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import type { ReactNode } from 'react'

import type { PillVariant } from '@/components'
import { Pill } from '@/components'
import { formatRelativeTime } from '@/lib/formatDate'

type PreviewState = 'idle' | 'loading' | 'success' | 'error' | 'stale'

interface PreviewStatusPillProps {
  state: PreviewState
  lastUpdated?: Date | null
  message?: string | null
}

const stateConfig: Record<
  PreviewState,
  { labelKey: string; variant: PillVariant; icon: ReactNode }
> = {
  idle: {
    labelKey: 'preview_status_idle',
    variant: 'ghost',
    icon: <Clock className="text-base-content/70" />,
  },
  loading: {
    labelKey: 'preview_status_running',
    variant: 'warning',
    icon: <Loader2 className="animate-spin" />,
  },
  success: {
    labelKey: 'preview_status_ready',
    variant: 'success',
    icon: <CheckCircle2 />,
  },
  error: {
    labelKey: 'preview_status_error',
    variant: 'error',
    icon: <AlertCircle />,
  },
  stale: {
    labelKey: 'preview_status_stale',
    variant: 'info',
    icon: <Clock />,
  },
}

export function PreviewStatusPill({
  state,
  lastUpdated,
  message,
}: PreviewStatusPillProps) {
  const { t: tOrders } = useTranslation('orders')

  const config = stateConfig[state]
  const showTimestamp = lastUpdated && state !== 'loading'

  return (
    <div className="flex flex-col items-start gap-1">
      <Pill
        icon={config.icon}
        variant={config.variant}
        size="sm"
        secondary={
          showTimestamp
            ? tOrders('preview_last_updated', {
                time: formatRelativeTime(lastUpdated),
              })
            : undefined
        }
      >
        {tOrders(config.labelKey)}
      </Pill>
      {message && state === 'error' && (
        <span className="text-xs text-error ps-1">{message}</span>
      )}
    </div>
  )
}
