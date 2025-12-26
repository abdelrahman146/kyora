import { memo, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import type { PublicBusiness } from '../../api/types';

interface BrandHeaderProps {
  business: PublicBusiness;
}

/**
 * BrandHeader Atom - Centered brand logo and name display
 * Shown at the top of the page content, not in the sticky header
 * Memoized to prevent unnecessary re-renders
 */
export const BrandHeader = memo<BrandHeaderProps>(function BrandHeader({ business }) {
  const { t } = useTranslation();

  const businessName = useMemo(
    () => business.name || t('store'),
    [business.name, t]
  );

  const logoUrl = useMemo(() => business.logoUrl || '', [business.logoUrl]);

  const avatarLetter = useMemo(
    () => (businessName.trim().charAt(0) || '?').toUpperCase(),
    [businessName]
  );

  return (
    <div className="flex flex-col items-center justify-center py-8 md:py-10">
      {logoUrl ? (
        <div className="avatar">
          <div className="w-20 h-20 md:w-24 md:h-24 rounded-full ring-2 ring-base-300/50">
            <img src={logoUrl} alt={businessName} className="object-cover" />
          </div>
        </div>
      ) : (
        <div className="avatar placeholder">
          <div className="w-20 h-20 md:w-24 md:h-24 rounded-full bg-primary text-primary-content ring-2 ring-base-300/50">
            <span className="text-3xl md:text-4xl font-bold">{avatarLetter}</span>
          </div>
        </div>
      )}
      <h1 className="font-bold text-2xl md:text-3xl mt-4 text-center px-4 text-neutral-900">
        {businessName}
      </h1>
    </div>
  );
});
