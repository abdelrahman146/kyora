import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { MessageSquarePlus, Paperclip } from 'lucide-react'
import type { CustomerNote } from '@/api/customer'
import type { AssetReference } from '@/api/types/asset'
import { BottomSheet } from '@/components/molecules/BottomSheet'
import { formatDateShort } from '@/lib/formatDate'

export interface CustomerNotesProps {
  notes: Array<CustomerNote>
  onAddNote?: (
    content: string,
    attachments: Array<AssetReference>,
  ) => Promise<void>
  isAddingNote?: boolean
  emptyMessage?: string
}

/**
 * CustomerNotes Component
 *
 * Displays a list of customer notes with attachments and supports adding new notes.
 * Reusable for both customer and order details pages.
 *
 * Features:
 * - Display notes with text, timestamp, and attachments
 * - Image/video thumbnails for AssetReference attachments
 * - Generic file representation for non-image/video files
 * - Add new note via BottomSheet with textarea and file upload
 * - Mobile-first, RTL-ready
 * - Empty state with CTA
 */
export function CustomerNotes({
  notes,
  onAddNote,
  isAddingNote = false,
  emptyMessage,
}: CustomerNotesProps) {
  const { t } = useTranslation('customers')
  const [isAddNoteOpen, setIsAddNoteOpen] = useState(false)
  const [noteContent, setNoteContent] = useState('')
  const [noteAttachments, setNoteAttachments] = useState<Array<AssetReference>>(
    [],
  )

  const handleAddNote = async () => {
    if (!noteContent.trim() || !onAddNote) return

    try {
      await onAddNote(noteContent.trim(), noteAttachments)
      setNoteContent('')
      setNoteAttachments([])
      setIsAddNoteOpen(false)
    } catch (error) {
      // Error handling is done in the parent component
    }
  }

  // Helper functions for future attachment support
  // const getFileIcon = (contentType: string | undefined) => { ... }
  // const getFileName = (asset: AssetReference): string => { ... }
  // const isImageOrVideo = (asset: AssetReference): boolean => { ... }

  if (notes.length === 0 && !onAddNote) {
    return (
      <div className="card bg-base-100 border border-base-300">
        <div className="card-body p-4">
          <h3 className="text-sm font-semibold text-base-content/60 uppercase tracking-wide mb-4">
            {t('details.notes')}
          </h3>
          <div className="text-center py-12 text-base-content/60">
            <div className="size-16 rounded-full bg-base-200 flex items-center justify-center mx-auto mb-4">
              <MessageSquarePlus size={32} className="opacity-40" />
            </div>
            <p className="font-medium">
              {emptyMessage || t('details.no_notes')}
            </p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <>
      <div className="card bg-base-100 border border-base-300">
        <div className="card-body p-4">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-sm font-semibold text-base-content/60 uppercase tracking-wide">
              {t('details.notes')}
            </h3>
            {onAddNote && (
              <button
                type="button"
                className="btn btn-primary btn-sm gap-2"
                onClick={() => setIsAddNoteOpen(true)}
                disabled={isAddingNote}
              >
                <MessageSquarePlus size={16} />
                <span className="hidden sm:inline">
                  {t('details.add_note')}
                </span>
                <span className="sm:hidden">+</span>
              </button>
            )}
          </div>

          {notes.length === 0 ? (
            <div className="text-center py-12 text-base-content/60">
              <div className="size-16 rounded-full bg-base-200 flex items-center justify-center mx-auto mb-4">
                <MessageSquarePlus size={32} className="opacity-40" />
              </div>
              <p className="font-medium mb-4">
                {emptyMessage || t('details.no_notes')}
              </p>
              {onAddNote && (
                <button
                  type="button"
                  className="btn btn-outline btn-sm gap-2"
                  onClick={() => setIsAddNoteOpen(true)}
                  disabled={isAddingNote}
                >
                  <MessageSquarePlus size={16} />
                  <span className="hidden sm:inline">
                    {t('details.add_note')}
                  </span>
                  <span className="sm:hidden">+</span>
                </button>
              )}
            </div>
          ) : (
            <div className="space-y-3">
              {notes.map((note) => (
                <div
                  key={note.id}
                  className="p-4 bg-base-200 rounded-lg border border-base-300"
                >
                  <p className="text-sm whitespace-pre-wrap mb-3">
                    {note.content}
                  </p>

                  {/* Attachments - Placeholder for future implementation */}
                  {/* Once backend supports attachments, uncomment and implement */}
                  {/* {note.attachments && note.attachments.length > 0 && (
                    <div className="space-y-2 mb-3">
                      <div className="flex items-center gap-1.5 text-xs text-base-content/60">
                        <Paperclip size={12} />
                        <span>{t('details.attachments')}</span>
                      </div>
                      <div className="grid grid-cols-2 sm:grid-cols-3 gap-2">
                        {note.attachments.map((attachment, idx) => {
                          const isMedia = isImageOrVideo(attachment)
                          const thumbnailUrl = getThumbnailUrl(attachment)

                          return isMedia && thumbnailUrl ? (
                            <a
                              key={idx}
                              href={attachment.url}
                              target="_blank"
                              rel="noopener noreferrer"
                              className="relative aspect-square rounded-lg overflow-hidden border border-base-300 hover:opacity-80 transition-opacity"
                            >
                              <img
                                src={thumbnailUrl}
                                alt={attachment.metadata?.altText || `Attachment ${idx + 1}`}
                                className="w-full h-full object-cover"
                              />
                            </a>
                          ) : (
                            <a
                              key={idx}
                              href={attachment.url}
                              target="_blank"
                              rel="noopener noreferrer"
                              className="flex items-center gap-2 p-2 rounded-lg border border-base-300 bg-base-100 hover:bg-base-200 transition-colors"
                            >
                              {getFileIcon(attachment.metadata?.caption)}
                              <span className="text-xs truncate flex-1">
                                {getFileName(attachment)}
                              </span>
                            </a>
                          )
                        })}
                      </div>
                    </div>
                  )} */}

                  <div className="text-xs text-base-content/60">
                    {formatDateShort(note.createdAt)}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Add Note Bottom Sheet */}
      {onAddNote && (
        <BottomSheet
          isOpen={isAddNoteOpen}
          onClose={() => {
            setIsAddNoteOpen(false)
            setNoteContent('')
            setNoteAttachments([])
          }}
          title={t('details.add_note')}
          size="lg"
          footer={
            <div className="flex gap-2 justify-end">
              <button
                type="button"
                className="btn btn-ghost"
                onClick={() => {
                  setIsAddNoteOpen(false)
                  setNoteContent('')
                  setNoteAttachments([])
                }}
                disabled={isAddingNote}
              >
                {t('common.cancel')}
              </button>
              <button
                type="button"
                className="btn btn-primary"
                onClick={() => void handleAddNote()}
                disabled={!noteContent.trim() || isAddingNote}
              >
                {isAddingNote && (
                  <span className="loading loading-spinner loading-sm" />
                )}
                {t('common.save')}
              </button>
            </div>
          }
        >
          <div className="space-y-4">
            <div className="form-control">
              <label className="label">
                <span className="label-text">{t('details.note_content')}</span>
              </label>
              <textarea
                className="textarea textarea-bordered h-32 resize-none"
                placeholder={t('details.note_placeholder')}
                value={noteContent}
                onChange={(e) => setNoteContent(e.target.value)}
                disabled={isAddingNote}
              />
            </div>

            {/* File Upload - Placeholder for future implementation */}
            {/* Once backend supports attachments, uncomment and implement FileUploadForm */}
            {/* <div className="form-control">
              <label className="label">
                <span className="label-text">{t('details.attachments_optional')}</span>
              </label>
              <FileUploadForm
                businessDescriptor={businessDescriptor}
                onUploadComplete={(assets) => setNoteAttachments(assets)}
                maxFiles={5}
                accept="image/*,video/*,application/pdf,.doc,.docx,.xls,.xlsx"
              />
            </div> */}

            <div className="alert alert-info">
              <Paperclip size={16} />
              <span className="text-sm">
                {t('details.attachments_coming_soon')}
              </span>
            </div>
          </div>
        </BottomSheet>
      )}
    </>
  )
}
