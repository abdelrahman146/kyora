import { useTranslation } from 'react-i18next'
import { Dialog } from './Dialog'

interface ResumeSessionDialogProps {
  open: boolean
  onResume: () => void | Promise<void>
  onStartFresh: () => void | Promise<void>
  email?: string
  stage?: string
  isLoading?: boolean
}

/**
 * Dialog shown when an existing onboarding session is found in localStorage.
 * Allows user to continue where they left off or start fresh.
 */
export function ResumeSessionDialog({
  open,
  onResume,
  onStartFresh,
  email,
  stage,
  isLoading = false,
}: ResumeSessionDialogProps) {
  const { t: tOnboarding } = useTranslation('onboarding')
  const { t: tCommon } = useTranslation('common')

  // Format stage for display
  const stageLabel = stage
    ? tOnboarding(`stages.${stage}`, { defaultValue: stage })
    : tCommon('unknown')

  return (
    <Dialog
      open={open}
      title={tOnboarding('resumeSession.title')}
      size="md"
      showCloseButton={false}
      closeOnBackdrop={false}
      footer={
        <>
          <button
            onClick={() => void onStartFresh()}
            className="btn btn-ghost"
            disabled={isLoading}
          >
            {tCommon('startFresh')}
          </button>
          <button
            onClick={() => void onResume()}
            className="btn btn-primary"
            disabled={isLoading}
          >
            {isLoading && (
              <span className="loading loading-spinner loading-sm"></span>
            )}
            {tCommon('continue')}
          </button>
        </>
      }
    >
      <div>
        <p className="mb-4">{tOnboarding('resumeSession.message')}</p>

        {email && (
          <div className="rounded-lg bg-base-200 p-4 space-y-2">
            <div className="flex items-center justify-between">
              <span className="text-sm font-semibold text-base-content/70">
                {tCommon('email')}
              </span>
              <span className="text-sm text-base-content">{email}</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm font-semibold text-base-content/70">
                {tOnboarding('stage')}
              </span>
              <span className="text-sm text-base-content">{stageLabel}</span>
            </div>
          </div>
        )}
      </div>
    </Dialog>
  )
}
