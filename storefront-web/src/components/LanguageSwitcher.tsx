import { useTranslation } from 'react-i18next';
import { ChevronDownIcon, LanguageIcon } from '@heroicons/react/24/outline';
import { supportedLanguages } from '../i18n/init';

export function LanguageSwitcher() {
  const { i18n, t } = useTranslation();
  const langs = supportedLanguages();
  const current = langs.find((l) => l.code === i18n.language) || langs[0];

  return (
    <div className="dropdown dropdown-end">
      <button type="button" className="btn btn-ghost btn-sm gap-2" aria-label={t('language')}>
        <LanguageIcon className="h-5 w-5" />
        <span className="font-semibold">{current.code.toUpperCase()}</span>
        <ChevronDownIcon className="h-4 w-4 opacity-70" />
      </button>

      <ul className="dropdown-content menu bg-base-100 rounded-box shadow p-2 w-44">
        {langs.map((l) => (
          <li key={l.code}>
            <button
              type="button"
              className={l.code === i18n.language ? 'menu-active' : undefined}
              onClick={() => i18n.changeLanguage(l.code)}
            >
              {l.label}
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}
