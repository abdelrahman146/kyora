import { Check, Globe, Languages } from 'lucide-react'
import { useTranslation } from 'react-i18next'

import { useLanguage } from '@/hooks/useLanguage'

const LANGUAGE_CONFIG = {
  en: {
    code: 'en',
    name: 'English',
    nativeName: 'English',
    flag: '🇬🇧',
    direction: 'ltr',
  },
  ar: {
    code: 'ar',
    name: 'Arabic',
    nativeName: 'العربية',
    flag: '🇸🇦',
    direction: 'rtl',
  },
} as const

type LanguageSwitcherVariant = 'dropdown' | 'toggle' | 'compact' | 'iconOnly'

export interface LanguageSwitcherProps {
  variant?: LanguageSwitcherVariant
  className?: string
  showLabel?: boolean
  showFlag?: boolean
}

export function LanguageSwitcher({
  variant = 'dropdown',
  className = '',
  showLabel = true,
  showFlag = true,
}: LanguageSwitcherProps) {
  const { language, changeLanguage, toggleLanguage, supportedLanguages } =
    useLanguage()
  const { t: tCommon } = useTranslation('common')

  const currentLangConfig = LANGUAGE_CONFIG[language]
  const languages = supportedLanguages.map((code) => LANGUAGE_CONFIG[code])

  if (variant === 'toggle') {
    const otherLanguage = languages.find((lang) => lang.code !== language)

    return (
      <button
        onClick={toggleLanguage}
        className={`btn btn-ghost btn-sm gap-2 ${className}`}
        aria-label={tCommon('changeLanguage')}
        type="button"
      >
        <Languages className="size-4" />
        {showLabel && (
          <span className="hidden sm:inline">{otherLanguage?.nativeName}</span>
        )}
      </button>
    )
  }

  if (variant === 'iconOnly') {
    return (
      <div className={`dropdown dropdown-end ${className}`}>
        <button
          tabIndex={0}
          className="btn btn-ghost btn-sm btn-circle"
          aria-label={tCommon('changeLanguage')}
          type="button"
        >
          <Globe className="size-5" />
        </button>
        <ul
          tabIndex={0}
          className="dropdown-content menu bg-base-200 rounded-box z-50 mt-2 w-52 p-2"
        >
          {languages.map((lang) => (
            <li key={lang.code}>
              <button
                onClick={() => {
                  changeLanguage(lang.code)
                }}
                className={`flex items-center gap-3 ${
                  language === lang.code ? 'active' : ''
                }`}
                type="button"
              >
                <span className="text-xl" role="img" aria-label={lang.name}>
                  {lang.flag}
                </span>
                <span className="flex-1">{lang.nativeName}</span>
                {language === lang.code && <Check className="size-4" />}
              </button>
            </li>
          ))}
        </ul>
      </div>
    )
  }

  if (variant === 'compact') {
    return (
      <div className={`dropdown dropdown-end ${className}`}>
        <button
          tabIndex={0}
          className="btn btn-ghost btn-sm gap-2"
          aria-label={tCommon('changeLanguage')}
          type="button"
        >
          {showFlag && (
            <span
              className="text-base"
              role="img"
              aria-label={currentLangConfig.name}
            >
              {currentLangConfig.flag}
            </span>
          )}
          <span className="uppercase text-xs font-semibold">
            {currentLangConfig.code}
          </span>
        </button>
        <ul
          tabIndex={0}
          className="dropdown-content menu bg-base-200 rounded-box z-50 mt-2 w-48 p-2"
        >
          {languages.map((lang) => (
            <li key={lang.code}>
              <button
                onClick={() => {
                  changeLanguage(lang.code)
                }}
                className={`flex items-center gap-3 ${
                  language === lang.code ? 'active' : ''
                }`}
                type="button"
              >
                <span className="text-base" role="img" aria-label={lang.name}>
                  {lang.flag}
                </span>
                <span className="flex-1 text-sm">{lang.nativeName}</span>
                {language === lang.code && <Check className="size-4" />}
              </button>
            </li>
          ))}
        </ul>
      </div>
    )
  }

  return (
    <div className={`dropdown dropdown-end ${className}`}>
      <button
        tabIndex={0}
        className="btn btn-ghost gap-3"
        aria-label={tCommon('changeLanguage')}
        type="button"
      >
        <Languages className="size-5" />
        {showFlag && (
          <span
            className="text-xl"
            role="img"
            aria-label={currentLangConfig.name}
          >
            {currentLangConfig.flag}
          </span>
        )}
        {showLabel && (
          <span className="font-medium">{currentLangConfig.nativeName}</span>
        )}
      </button>
      <ul
        tabIndex={0}
        className="dropdown-content menu bg-base-200 rounded-box z-50 mt-3 w-64 p-2"
      >
        <li className="menu-title px-4 py-2">
          <span className="text-base-content/70 text-sm font-semibold">
            {tCommon('selectLanguage')}
          </span>
        </li>
        {languages.map((lang) => (
          <li key={lang.code}>
            <button
              onClick={() => {
                changeLanguage(lang.code)
              }}
              className={`flex items-center gap-3 py-3 ${
                language === lang.code ? 'active bg-primary/10' : ''
              }`}
              type="button"
            >
              <span className="text-2xl" role="img" aria-label={lang.name}>
                {lang.flag}
              </span>
              <div className="flex-1 text-start">
                <div className="font-medium">{lang.nativeName}</div>
                <div className="text-xs text-base-content/60">{lang.name}</div>
              </div>
              {language === lang.code && (
                <Check className="size-5 text-primary" />
              )}
            </button>
          </li>
        ))}
      </ul>
    </div>
  )
}
