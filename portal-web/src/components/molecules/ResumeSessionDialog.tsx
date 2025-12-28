import { useTranslation } from 'react-i18next';
import { Dialog } from '@/components/atoms/Dialog';

interface ResumeSessionDialogProps {
  open: boolean;
  onResume: () => void | Promise<void>;
  onStartFresh: () => void | Promise<void>;
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

  // Format stage for display
  const stageLabel = stage
    ? t(`onboarding:stages.${stage}`, { defaultValue: stage })
    : t('common:unknown');

  return (
    <Dialog
      open={open}
      title={t('onboarding:resumeSession.title')}
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
            {t('common:startFresh')}
          </button>
          <button
            onClick={() => void onResume()}
            className="btn btn-primary"
            disabled={isLoading}
          >
            {isLoading && <span className="loading loading-spinner loading-sm"></span>}
            {t('common:continue')}
          </button>
        </>
      }
    >
      <div>
        <p className="mb-4">
          {t('onboarding:resumeSession.message')}
        </p>
        
        {email && (
          <div className="rounded-lg bg-base-200 p-4 space-y-2">
            <div className="flex items-center justify-between">
              <span className="text-sm font-semibold text-base-content/70">
                {t('common:email')}
              </span>
              <span className="text-sm text-base-content">
                {email}
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm font-semibold text-base-content/70">
                {t('onboarding:stage')}
              </span>
              <span className="text-sm text-base-content">
                {stageLabel}
              </span>
            </div>
          </div>
        )}
      </div>
    </Dialog>
  );
}
