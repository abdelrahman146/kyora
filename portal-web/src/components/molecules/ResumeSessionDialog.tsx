import { useTranslation } from 'react-i18next';

interface ResumeSessionDialogProps {
  open: boolean;
  onResume: () => void;
  onStartFresh: () => void;
  email?: string;
  stage?: string;
  isLoading?: boolean;
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
  const { t } = useTranslation(['onboarding', 'common']);

  if (!open) return null;

  // Format stage for display
  const stageLabel = stage
    ? t(`onboarding:stages.${stage}`, { defaultValue: stage })
    : t('common:unknown');

  return (
    <dialog className="modal modal-open">
      <div className="modal-box">
        <h3 className="text-lg font-bold">
          {t('onboarding:resumeSession.title')}
        </h3>
        
        <div className="py-4">
          <p className="mb-2">
            {t('onboarding:resumeSession.message')}
          </p>
          
          {email && (
            <div className="mt-4 rounded-lg bg-base-200 p-3">
              <p className="text-sm">
                <span className="font-semibold">{t('common:email')}:</span>{' '}
                {email}
              </p>
              <p className="text-sm">
                <span className="font-semibold">{t('onboarding:stage')}:</span>{' '}
                {stageLabel}
              </p>
            </div>
          )}
        </div>

        <div className="modal-action">
          <button
            onClick={onStartFresh}
            className="btn btn-ghost"
            disabled={isLoading}
          >
            {t('common:startFresh')}
          </button>
          <button
            onClick={onResume}
            className="btn btn-primary"
            disabled={isLoading}
          >
            {isLoading && <span className="loading loading-spinner loading-sm"></span>}
            {t('common:continue')}
          </button>
        </div>
      </div>
    </dialog>
  );
}
