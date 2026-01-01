import { GripVertical } from 'lucide-react'
import { cn } from '@/lib/utils'

export interface DragHandleProps {
  className?: string
  'aria-label'?: string
}

/**
 * DragHandle - Touch-optimized drag handle for reorderable lists
 *
 * Features:
 * - 50px touch target (mobile-first)
 * - GripVertical icon
 * - Cursor: grab/grabbing
 * - Keyboard accessible
 * - RTL support
 *
 * @example
 * <DragHandle aria-label="Drag to reorder" />
 */
export function DragHandle({
  className,
  'aria-label': ariaLabel,
}: DragHandleProps) {
  return (
    <div
      className={cn(
        'flex items-center justify-center',
        'w-[50px] h-[50px]',
        'cursor-grab active:cursor-grabbing',
        'text-base-content/40 hover:text-base-content/70',
        'transition-colors duration-200',
        'touch-manipulation',
        className,
      )}
      aria-label={ariaLabel || 'Drag to reorder'}
      role="button"
      tabIndex={0}
    >
      <GripVertical className="w-5 h-5" aria-hidden="true" />
    </div>
  )
}
