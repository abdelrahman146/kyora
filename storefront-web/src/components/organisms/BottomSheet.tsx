import { memo, useEffect, useRef } from 'react';
import type { ReactNode } from 'react';
import { X } from 'lucide-react';
import { useTranslation } from 'react-i18next';

interface BottomSheetProps {
  isOpen: boolean;
  onClose: () => void;
  title?: string;
  children: ReactNode;
  footer?: ReactNode;
}

/**
 * BottomSheet Organism - Mobile-optimized drawer
 * Memoized to prevent unnecessary re-renders
 * Uses useEffect cleanup for scroll lock and keyboard events to prevent memory leaks
 * Follows KDS 3.2 & 4.5 with slide-up animation and backdrop blur
 */
export const BottomSheet = memo<BottomSheetProps>(function BottomSheet({ 
  isOpen, 
  onClose, 
  title, 
  children,
  footer,
}) {
  const { t } = useTranslation();
  const sheetRef = useRef<HTMLDivElement>(null);
  const backdropRef = useRef<HTMLDivElement>(null);

  // Prevent body scroll when sheet is open - with cleanup to prevent memory leaks
  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }

    // Cleanup function to restore scroll
    return () => {
      document.body.style.overflow = '';
    };
  }, [isOpen]);

  // Handle keyboard events - with cleanup to prevent memory leaks
  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isOpen) {
        onClose();
      }
    };

    window.addEventListener('keydown', handleEscape);
    
    // Cleanup function to remove event listener
    return () => {
      window.removeEventListener('keydown', handleEscape);
    };
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-end justify-center">
      {/* Backdrop */}
      <div
        ref={backdropRef}
        className="absolute inset-0 bg-black/40 backdrop-blur-sm"
        onClick={onClose}
        aria-hidden="true"
      />

      {/* Sheet */}
      <div
        ref={sheetRef}
        className="relative w-full max-w-5xl bg-white rounded-t-xl shadow-float max-h-[90vh] flex flex-col animate-slide-up"
        role="dialog"
        aria-modal="true"
        aria-label={title || t('sheet')}
      >
        {/* Handle */}
        <div className="flex justify-center pt-3 pb-2">
          <div className="w-12 h-1 bg-neutral-300 rounded-full" />
        </div>

        {/* Header */}
        {title && (
          <div className="flex items-center justify-between px-4 pb-3 border-b border-neutral-200">
            <h2 className="text-xl font-bold text-neutral-900">{title}</h2>
            <button
              type="button"
              onClick={onClose}
              className="btn btn-ghost btn-square btn-sm active-scale focus-ring"
              aria-label={t('close')}
            >
              <X className="w-5 h-5" />
            </button>
          </div>
        )}

        {/* Content (Scrollable) */}
        <div className="flex-1 overflow-y-auto px-4 py-4 scrollbar-hide">
          {children}
        </div>

        {/* Footer (Sticky) */}
        {footer && (
          <div className="border-t border-neutral-200 px-4 py-4 safe-bottom bg-white">
            {footer}
          </div>
        )}
      </div>
    </div>
  );
});
