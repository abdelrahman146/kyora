import { memo, useMemo, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { Globe } from 'lucide-react';
import { supportedLanguages } from '../../i18n/init';

/**
 * LanguageSwitcher Molecule - Language selection dropdown
 * Memoized to prevent unnecessary re-renders
 * Optimized with useMemo and useCallback
 */
export const LanguageSwitcher = memo(function LanguageSwitcher() {
  const { i18n, t } = useTranslation();
  
  // Memoize supported languages
  const langs = useMemo(() => supportedLanguages(), []);
  
  // Memoize current language
  const current = useMemo(
    () => langs.find((l) => l.code === i18n.language) || langs[0],
    [langs, i18n.language]
  );

  // Memoize language change handler
  const handleLanguageChange = useCallback(
    (code: string) => {
      i18n.changeLanguage(code);
    },
    [i18n]
  );

  return (
    <div className="dropdown dropdown-end">
      <button 
        type="button" 
        tabIndex={0}
        className="btn btn-sm btn-ghost gap-2 active-scale focus-ring" 
        aria-label={t('language')}
      >
        <Globe className="w-4 h-4" strokeWidth={2} />
        <span className="font-medium text-sm">{current.label}</span>
      </button>

      <ul 
        tabIndex={0}
        className="dropdown-content menu bg-base-100 rounded-xl shadow-lg border border-base-200 p-2 w-32 z-50 mt-1"
      >
        {langs.map((l) => (
          <li key={l.code}>
            <button
              type="button"
              className={`text-sm active-scale focus-ring rounded-lg ${
                l.code === i18n.language 
                  ? 'bg-primary text-primary-content font-semibold' 
                  : 'hover:bg-base-200'
              }`}
              onClick={() => handleLanguageChange(l.code)}
            >
              {l.label}
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
});
