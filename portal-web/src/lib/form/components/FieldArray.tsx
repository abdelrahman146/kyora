/**
 * FieldArray Component - Form Composition Layer for Dynamic Array Fields
 *
 * Pre-bound component for managing dynamic array fields with add/remove/reorder operations.
 * Supports drag-and-drop reordering, validation, and animations.
 *
 * Usage within form:
 * ```tsx
 * <form.Field name="phoneNumbers" mode="array">
 *   {(field) => (
 *     <field.FieldArray
 *       label="Phone Numbers"
 *       minItems={1}
 *       maxItems={5}
 *       addButtonLabel="Add Phone"
 *       defaultValue={() => ({ number: '', type: 'mobile' })}
 *       render={(item, index, operations) => (
 *         <div className="flex gap-2">
 *           <form.Field name={`phoneNumbers[${index}].number`}>
 *             {(subField) => (
 *               <subField.TextField label={`Phone ${index + 1}`} />
 *             )}
 *           </form.Field>
 *           <button onClick={operations.remove}>Remove</button>
 *         </div>
 *       )}
 *     />
 *   )}
 * </form.Field>
 * ```
 */

import { useCallback, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import {
  DndContext,
  DragOverlay,
  KeyboardSensor,
  PointerSensor,
  TouchSensor,
  closestCenter,
  useSensor,
  useSensors,
} from '@dnd-kit/core'
import {
  SortableContext,
  sortableKeyboardCoordinates,
  useSortable,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import { AlertCircle, GripVertical, Plus, Trash2 } from 'lucide-react'
import { useFieldContext } from '../contexts'
import type { DragEndEvent, DragStartEvent } from '@dnd-kit/core'
import type { ArrayItemOperations, FieldArrayProps } from '../types'
import { cn } from '@/lib/utils'

/**
 * Sortable Item Wrapper
 */
interface SortableItemProps<T> {
  id: string
  index: number
  item: T
  isReorderable: boolean
  operations: ArrayItemOperations
  render: (
    item: T,
    index: number,
    operations: ArrayItemOperations,
  ) => React.ReactNode
  isDragging?: boolean
}

function SortableItem<T>({
  id,
  index,
  item,
  isReorderable,
  operations,
  render,
  isDragging,
}: SortableItemProps<T>) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging: isSortableDragging,
  } = useSortable({ id, disabled: !isReorderable })

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging || isSortableDragging ? 0.5 : 1,
  }

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={cn(
        'group relative rounded-lg border border-base-300 bg-base-100 transition-all duration-200',
        (isDragging || isSortableDragging) &&
          'shadow-lg ring-2 ring-primary/20',
        'animate-in fade-in slide-in-from-top-2 duration-300',
      )}
      role="listitem"
    >
      <div className="flex items-start gap-3 p-4">
        {/* Drag Handle */}
        {isReorderable && (
          <button
            type="button"
            className={cn(
              'btn btn-ghost btn-sm btn-square shrink-0 cursor-grab active:cursor-grabbing',
              'text-base-content/40 hover:text-base-content/70',
              'focus:outline-none focus:ring-2 focus:ring-primary/20',
            )}
            aria-label={`Drag to reorder item ${index + 1}`}
            {...attributes}
            {...listeners}
          >
            <GripVertical size={18} aria-hidden="true" />
          </button>
        )}

        {/* Item Content */}
        <div className="flex-1 min-w-0">{render(item, index, operations)}</div>

        {/* Remove Button */}
        <button
          type="button"
          onClick={operations.remove}
          className={cn(
            'btn btn-ghost btn-sm btn-square shrink-0 text-error',
            'opacity-0 group-hover:opacity-100 transition-opacity',
            'focus:opacity-100 focus:outline-none focus:ring-2 focus:ring-error/20',
          )}
          aria-label={`Remove item ${index + 1}`}
        >
          <Trash2 size={18} aria-hidden="true" />
        </button>
      </div>
    </div>
  )
}

/**
 * FieldArray Component
 */
export function FieldArray<T = any>(props: FieldArrayProps<T>) {
  const {
    label,
    minItems = 0,
    maxItems,
    defaultValue,
    addButtonLabel,
    emptyMessage,
    reorderable = true,
    helperText,
    render,
    className,
  } = props

  const field = useFieldContext<Array<T>>()
  const { t } = useTranslation(['common', 'errors'])

  // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
  const items = field.state.value || []
  const [activeId, setActiveId] = useState<string | null>(null)

  // Configure drag-and-drop sensors
  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8, // Prevent accidental drags
      },
    }),
    useSensor(TouchSensor, {
      activationConstraint: {
        delay: 250, // Long-press on mobile
        tolerance: 5,
      },
    }),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    }),
  )

  // Extract error from field state and translate
  const error = useMemo(() => {
    const errors = field.state.meta.errors
    if (errors.length === 0) return undefined

    const firstError = errors[0]
    if (typeof firstError === 'string') {
      return t(firstError, { ns: 'errors' })
    }

    if (
      typeof firstError === 'object' &&
      firstError &&
      'message' in firstError
    ) {
      const errorObj = firstError as { message: string; code?: number }
      return t(errorObj.message, { ns: 'errors' })
    }

    return undefined
  }, [field.state.meta.errors, t])

  const showError = field.state.meta.isTouched && error

  // Array operations
  const addItem = useCallback(() => {
    const newItem = defaultValue ? defaultValue() : ({} as T)
    const newItems = [...items, newItem]
    field.handleChange(newItems)
  }, [items, defaultValue, field])

  const removeItem = useCallback(
    (index: number) => {
      const newItems = items.filter((_, i) => i !== index)
      field.handleChange(newItems)
    },
    [items, field],
  )

  const moveItem = useCallback(
    (fromIndex: number, toIndex: number) => {
      const newItems = [...items]
      const [removed] = newItems.splice(fromIndex, 1)
      newItems.splice(toIndex, 0, removed)
      field.handleChange(newItems)
    },
    [items, field],
  )

  // Drag handlers
  const handleDragStart = (event: DragStartEvent) => {
    setActiveId(String(event.active.id))
  }

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event

    if (over && active.id !== over.id) {
      const oldIndex = items.findIndex((_, i) => String(i) === active.id)
      const newIndex = items.findIndex((_, i) => String(i) === over.id)

      if (oldIndex !== -1 && newIndex !== -1) {
        moveItem(oldIndex, newIndex)
      }
    }

    setActiveId(null)
  }

  const handleDragCancel = () => {
    setActiveId(null)
  }

  // Create operations for each item
  const createOperations = useCallback(
    (index: number): ArrayItemOperations => ({
      remove: () => {
        if (minItems > 0 && items.length <= minItems) {
          // Show error or confirmation
          return
        }
        removeItem(index)
      },
      moveUp: () => {
        if (index > 0) {
          moveItem(index, index - 1)
        }
      },
      moveDown: () => {
        if (index < items.length - 1) {
          moveItem(index, index + 1)
        }
      },
    }),
    [items.length, minItems, removeItem, moveItem],
  )

  const canAddMore = !maxItems || items.length < maxItems
  const activeItem = activeId ? items[parseInt(activeId)] : null

  return (
    <div className={cn('form-control', className)}>
      {/* Label */}
      {label && (
        <label className="label">
          <span className="label-text text-base-content/70 font-medium">
            {label}
            {minItems > 0 && <span className="text-error ms-1">*</span>}
          </span>
          {maxItems && (
            <span className="label-text-alt text-base-content/60">
              {items.length} / {maxItems}
            </span>
          )}
        </label>
      )}

      {/* Helper Text */}
      {helperText && !showError && (
        <div className="label mb-2">
          <span className="label-text-alt text-base-content/60">
            {helperText}
          </span>
        </div>
      )}

      {/* Error Message */}
      {showError && (
        <div className="alert alert-error mb-4" role="alert">
          <AlertCircle size={20} aria-hidden="true" />
          <span>{error}</span>
        </div>
      )}

      {/* Array Items */}
      {items.length > 0 ? (
        <DndContext
          sensors={sensors}
          collisionDetection={closestCenter}
          onDragStart={handleDragStart}
          onDragEnd={handleDragEnd}
          onDragCancel={handleDragCancel}
        >
          <SortableContext
            items={items.map((_, i) => String(i))}
            strategy={verticalListSortingStrategy}
          >
            <ul className="space-y-3 mb-4" role="list" aria-label={label}>
              {items.map((item, index) => (
                <SortableItem
                  key={index}
                  id={String(index)}
                  index={index}
                  item={item}
                  isReorderable={reorderable && items.length > 1}
                  operations={createOperations(index)}
                  render={render}
                  isDragging={String(index) === activeId}
                />
              ))}
            </ul>
          </SortableContext>

          {/* Drag Overlay */}
          <DragOverlay>
            {activeItem && activeId ? (
              <div className="rounded-lg border-2 border-primary bg-base-100 shadow-xl opacity-90">
                <div className="flex items-start gap-3 p-4">
                  <div className="flex-1 min-w-0">
                    {render(
                      activeItem,
                      parseInt(activeId),
                      createOperations(parseInt(activeId)),
                    )}
                  </div>
                </div>
              </div>
            ) : null}
          </DragOverlay>
        </DndContext>
      ) : (
        /* Empty State */
        <div className="border-2 border-dashed border-base-300 rounded-lg p-8 text-center mb-4 transition-all duration-200 hover:border-base-content/20">
          <p className="text-base-content/60 mb-4">
            {emptyMessage || t('common.array.noItems', 'No items added yet')}
          </p>
        </div>
      )}

      {/* Add Button */}
      <div className="flex gap-2">
        <button
          type="button"
          onClick={addItem}
          disabled={!canAddMore}
          className={cn(
            'btn btn-outline btn-primary gap-2',
            items.length === 0 && 'btn-lg',
          )}
          aria-label={addButtonLabel || t('common.array.addItem', 'Add Item')}
        >
          <Plus size={18} aria-hidden="true" />
          <span>{addButtonLabel || t('common.array.addItem', 'Add Item')}</span>
        </button>

        {/* Max Items Warning */}
        {maxItems && items.length >= maxItems && (
          <div className="alert alert-warning py-2 px-4 text-sm">
            <span>
              {t(
                'common.array.maxItemsReached',
                'Maximum {{max}} items allowed',
                { max: maxItems },
              )}
            </span>
          </div>
        )}
      </div>
    </div>
  )
}
