import { useId, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { MessageSquarePlus, Trash2 } from 'lucide-react'
import { z } from 'zod'
import { useKyoraForm } from '@/lib/form'
import { BottomSheet } from '@/components/molecules/BottomSheet'
import { formatDateShort } from '@/lib/formatDate'

export interface Note {
  id: string
  content: string
  CreatedAt: string
}

export interface NotesProps {
  notes: Array<Note>
  onAddNote?: (content: string) => Promise<void>
  onDeleteNote?: (noteId: string) => Promise<void>
  isAddingNote?: boolean
  isDeletingNote?: boolean
  deletingNoteId?: string | null
  emptyMessage?: string
  showAddButton?: boolean
  className?: string
  maxLength?: number
}

/**
 * Notes Component
 *
 * Production-grade, reusable notes system combining list + add form.
 * Designed for customer notes, order notes, and any entity that needs note-taking.
 *
 * Features:
 * - Mobile-first BottomSheet for adding notes
 * - Kyora form system with TextareaField (proper validation)
 * - Clean timeline/feed design (supports future attachments)
 * - RTL-first with logical properties
 * - Accessible with ARIA labels
 * - Consistent sizing with portal design system
 * - Empty state with CTA
 * - Optimistic UI with loading states
 */
export function Notes({
  notes,
  onAddNote,
  onDeleteNote,
  isAddingNote = false,
  isDeletingNote = false,
  deletingNoteId = null,
  emptyMessage,
  showAddButton = true,
  className = '',
  maxLength = 1000,
}: NotesProps) {
  const { t } = useTranslation('common')
  const [isAddSheetOpen, setIsAddSheetOpen] = useState(false)
  const formId = useId()

  // Form for adding notes
  const form = useKyoraForm({
    defaultValues: {
      content: '',
    },
    onSubmit: async ({ value }) => {
      if (!onAddNote) return
      await onAddNote(value.content.trim())
      form.reset()
      setIsAddSheetOpen(false)
    },
  })

  // Validation schema
  const noteSchema = z.object({
    content: z
      .string()
      .min(1, t('notes.validation.required'))
      .max(maxLength, t('notes.validation.max_length', { max: maxLength })),
  })

  // Empty state
  if (notes.length === 0) {
    return (
      <>
        <div
          className={`card bg-base-100 border border-base-300 ${className}`}
          role="region"
          aria-label={t('notes.section_label')}
        >
          <div className="card-body p-4 sm:p-6">
            <div className="text-center py-8 sm:py-12">
              <div className="size-16 sm:size-20 rounded-full bg-base-200 flex items-center justify-center mx-auto mb-3 sm:mb-4">
                <MessageSquarePlus
                  size={32}
                  className="text-base-content/30 sm:size-10"
                  aria-hidden="true"
                />
              </div>
              <h3 className="text-base sm:text-lg font-semibold text-base-content mb-1 sm:mb-2">
                {emptyMessage || t('notes.no_notes')}
              </h3>
              <p className="text-sm text-base-content/60 mb-4 sm:mb-6">
                {t('notes.empty_description')}
              </p>
              {onAddNote && showAddButton && (
                <button
                  type="button"
                  className="btn btn-primary btn-sm gap-2"
                  onClick={() => setIsAddSheetOpen(true)}
                  disabled={isAddingNote}
                  aria-label={t('notes.add_first_note')}
                >
                  {isAddingNote ? (
                    <span
                      className="loading loading-spinner loading-xs"
                      aria-hidden="true"
                    />
                  ) : (
                    <MessageSquarePlus size={16} aria-hidden="true" />
                  )}
                  <span>{t('notes.add_first_note')}</span>
                </button>
              )}
            </div>
          </div>
        </div>

        {/* Add Note BottomSheet */}
        {onAddNote && (
          <form.AppForm>
            <BottomSheet
              isOpen={isAddSheetOpen}
              onClose={() => {
                setIsAddSheetOpen(false)
                form.reset()
              }}
              title={t('notes.add_note')}
              size="md"
              footer={
                <div className="flex gap-3">
                  <button
                    type="button"
                    className="btn btn-ghost btn-sm flex-1"
                    onClick={() => {
                      setIsAddSheetOpen(false)
                      form.reset()
                    }}
                    disabled={isAddingNote}
                  >
                    {t('cancel')}
                  </button>

                  <form.Subscribe selector={(state) => state.values.content}>
                    {(contentValue) => (
                      <form.SubmitButton
                        variant="primary"
                        size="sm"
                        className="flex-1"
                        disabled={isAddingNote || !contentValue.trim()}
                        form={`add-note-form-${formId}`}
                      >
                        {isAddingNote ? (
                          <>
                            <span
                              className="loading loading-spinner loading-xs"
                              aria-hidden="true"
                            />
                            <span>{t('notes.adding')}</span>
                          </>
                        ) : (
                          <>
                            <MessageSquarePlus size={16} aria-hidden="true" />
                            <span>{t('notes.add')}</span>
                          </>
                        )}
                      </form.SubmitButton>
                    )}
                  </form.Subscribe>
                </div>
              }
            >
              <form.FormRoot
                id={`add-note-form-${formId}`}
                className="space-y-4"
              >
                <form.FormError />

                <form.AppField
                  name="content"
                  validators={{
                    onChange: noteSchema.shape.content,
                  }}
                >
                  {(field) => (
                    <field.TextareaField
                      label={t('notes.content_label')}
                      placeholder={t('notes.content_placeholder')}
                      rows={6}
                      maxLength={maxLength}
                      showCounter
                      autoFocus
                    />
                  )}
                </form.AppField>
              </form.FormRoot>
            </BottomSheet>
          </form.AppForm>
        )}
      </>
    )
  }

  // Notes list - Clean feed/timeline design
  return (
    <>
      <div
        className={`card bg-base-100 border border-base-300 ${className}`}
        role="region"
        aria-label={t('notes.section_label')}
      >
        <div className="card-body p-4 sm:p-6">
          {/* Header */}
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-base font-semibold text-base-content">
              {t('notes.notes')} ({notes.length})
            </h3>
            {onAddNote && showAddButton && (
              <button
                type="button"
                className="btn btn-primary btn-sm gap-2"
                onClick={() => setIsAddSheetOpen(true)}
                disabled={isAddingNote}
                aria-label={t('notes.add_note')}
              >
                {isAddingNote ? (
                  <span
                    className="loading loading-spinner loading-xs"
                    aria-hidden="true"
                  />
                ) : (
                  <MessageSquarePlus size={16} aria-hidden="true" />
                )}
                <span className="hidden sm:inline">{t('notes.add_note')}</span>
                <span className="sm:hidden">{t('notes.add')}</span>
              </button>
            )}
          </div>

          {/* Notes Timeline/Feed */}
          <div className="space-y-3">
            {notes.map((note, index) => {
              const isDeleting = isDeletingNote && deletingNoteId === note.id
              const isLast = index === notes.length - 1

              return (
                <div
                  key={note.id}
                  className={`group relative ${!isLast ? 'pb-3 border-b border-base-300' : ''} ${
                    isDeleting ? 'opacity-50' : ''
                  }`}
                  role="article"
                  aria-label={t('notes.note_label')}
                >
                  {/* Note Header: Timestamp + Delete */}
                  <div className="flex items-center justify-between gap-2 mb-2">
                    <time
                      className="text-xs text-base-content/60 font-medium"
                      dateTime={note.CreatedAt}
                    >
                      {formatDateShort(note.CreatedAt)}
                    </time>

                    {/* Delete Button - Always visible on touch devices, hover on desktop */}
                    {onDeleteNote && (
                      <button
                        type="button"
                        className="btn btn-ghost btn-xs btn-square opacity-100 sm:opacity-0 sm:group-hover:opacity-100 transition-opacity"
                        onClick={() => onDeleteNote(note.id)}
                        disabled={isDeleting}
                        aria-label={t('notes.delete_note')}
                        title={t('notes.delete_note')}
                      >
                        {isDeleting ? (
                          <span
                            className="loading loading-spinner loading-xs"
                            aria-hidden="true"
                          />
                        ) : (
                          <Trash2
                            size={14}
                            className="text-error"
                            aria-hidden="true"
                          />
                        )}
                      </button>
                    )}
                  </div>

                  {/* Note Content */}
                  <div className="text-sm text-base-content whitespace-pre-wrap break-words leading-relaxed">
                    {note.content}
                  </div>

                  {/* Future: Attachments would go here */}
                  {/* <div className="flex flex-wrap gap-2 mt-3">
                    {note.attachments?.map(attachment => (
                      <AttachmentPreview key={attachment.id} {...attachment} />
                    ))}
                  </div> */}
                </div>
              )
            })}
          </div>
        </div>
      </div>

      {/* Add Note BottomSheet */}
      {onAddNote && (
        <form.AppForm>
          <BottomSheet
            isOpen={isAddSheetOpen}
            onClose={() => {
              setIsAddSheetOpen(false)
              form.reset()
            }}
            title={t('notes.add_note')}
            size="md"
            footer={
              <div className="flex gap-3">
                <button
                  type="button"
                  className="btn btn-ghost btn-sm flex-1"
                  onClick={() => {
                    setIsAddSheetOpen(false)
                    form.reset()
                  }}
                  disabled={isAddingNote}
                >
                  {t('cancel')}
                </button>

                <form.Subscribe selector={(state) => state.values.content}>
                  {(contentValue) => (
                    <form.SubmitButton
                      variant="primary"
                      size="sm"
                      className="flex-1"
                      disabled={isAddingNote || !contentValue.trim()}
                      form={`add-note-form-${formId}`}
                    >
                      {isAddingNote ? (
                        <>
                          <span
                            className="loading loading-spinner loading-xs"
                            aria-hidden="true"
                          />
                          <span>{t('notes.adding')}</span>
                        </>
                      ) : (
                        <>
                          <MessageSquarePlus size={16} aria-hidden="true" />
                          <span>{t('notes.add')}</span>
                        </>
                      )}
                    </form.SubmitButton>
                  )}
                </form.Subscribe>
              </div>
            }
          >
            <form.FormRoot id={`add-note-form-${formId}`} className="space-y-4">
              <form.FormError />

              <form.AppField
                name="content"
                validators={{
                  onChange: noteSchema.shape.content,
                }}
              >
                {(field) => (
                  <field.TextareaField
                    label={t('notes.content_label')}
                    placeholder={t('notes.content_placeholder')}
                    rows={6}
                    maxLength={maxLength}
                    showCounter
                    autoFocus
                  />
                )}
              </form.AppField>
            </form.FormRoot>
          </BottomSheet>
        </form.AppForm>
      )}
    </>
  )
}
