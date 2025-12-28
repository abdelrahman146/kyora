import { Languages, Globe, Check } from "lucide-react";
import { useLanguage } from "../../hooks/useLanguage";
import { useTranslation } from "react-i18next";

/**
 * Language configuration for supported locales
 * Add new languages here as they become available
 */
const LANGUAGE_CONFIG = {
  en: {
    code: "en",
    name: "English",
    nativeName: "English",
    flag: "🇬🇧",
    direction: "ltr",
  },
  ar: {
    code: "ar",
    name: "Arabic",
    nativeName: "العربية",
    flag: "🇸🇦",
    direction: "rtl",
  },
  // Future languages can be added here:
  // fr: {
  //   code: "fr",
  //   name: "French",
  //   nativeName: "Français",
  //   flag: "🇫🇷",
  //   direction: "ltr",
  // },
  // es: {
  //   code: "es",
  //   name: "Spanish",
  //   nativeName: "Español",
  //   flag: "🇪🇸",
  //   direction: "ltr",
  // },
} as const;

type LanguageSwitcherVariant = 
  | "dropdown" // Full dropdown with all languages
  | "toggle" // Simple toggle button (best for 2 languages)
  | "compact" // Compact dropdown in navbar
  | "iconOnly"; // Icon-only dropdown

interface LanguageSwitcherProps {
  variant?: LanguageSwitcherVariant;
  className?: string;
  showLabel?: boolean;
  showFlag?: boolean;
}

/**
 * Universal Language Switcher Component
 * 
 * A flexible, reusable language switcher that supports multiple design variants
 * and can easily be extended with new languages.
 * 
 * Features:
 * - Multiple design variants (dropdown, toggle, compact, iconOnly)
 * - Extensible language configuration
 * - RTL-aware positioning
 * - Accessible keyboard navigation
 * - Flag emoji support
 * - Native language names
 * 
 * Variants:
 * - `dropdown`: Full dropdown with language names and flags (default)
 * - `toggle`: Simple toggle button (best for 2 languages)
 * - `compact`: Compact version for navbars
 * - `iconOnly`: Shows only globe icon with dropdown
 * 
 * Usage:
 * ```tsx
 * // Full dropdown
 * <LanguageSwitcher variant="dropdown" />
 * 
 * // Simple toggle
 * <LanguageSwitcher variant="toggle" />
 * 
 * // Navbar compact version
 * <LanguageSwitcher variant="compact" />
 * 
 * // Icon only
 * <LanguageSwitcher variant="iconOnly" />
 * ```
 */
export function LanguageSwitcher({
  variant = "dropdown",
  className = "",
  showLabel = true,
  showFlag = true,
}: LanguageSwitcherProps) {
  const { language, changeLanguage, toggleLanguage, supportedLanguages } = useLanguage();
  const { t } = useTranslation();

  const currentLangConfig = LANGUAGE_CONFIG[language];
  const languages = supportedLanguages.map(code => LANGUAGE_CONFIG[code]);

  // Toggle variant - simple button to switch between 2 languages
  if (variant === "toggle") {
    const otherLanguage = languages.find(lang => lang.code !== language);
    
    return (
      <button
        onClick={toggleLanguage}
        className={`btn btn-ghost btn-sm gap-2 ${className}`}
        aria-label={t("common.changeLanguage")}
      >
        <Languages className="size-4" />
        {showLabel && (
          <span className="hidden sm:inline">
            {otherLanguage?.nativeName}
          </span>
        )}
      </button>
    );
  }

  // Icon only variant - shows just globe icon
  if (variant === "iconOnly") {
    return (
      <div className={`dropdown dropdown-end ${className}`}>
        <button
          tabIndex={0}
          className="btn btn-ghost btn-sm btn-circle"
          aria-label={t("common.changeLanguage")}
        >
          <Globe className="size-5" />
        </button>
        <ul
          tabIndex={0}
          className="dropdown-content menu bg-base-200 rounded-box z-50 mt-2 w-52 p-2 shadow-lg"
        >
          {languages.map((lang) => (
            <li key={lang.code}>
              <button
                onClick={() => {
                  changeLanguage(lang.code);
                }}
                className={`flex items-center gap-3 ${
                  language === lang.code ? "active" : ""
                }`}
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
    );
  }

  // Compact variant - for navbars and tight spaces
  if (variant === "compact") {
    return (
      <div className={`dropdown dropdown-end ${className}`}>
        <button
          tabIndex={0}
          className="btn btn-ghost btn-sm gap-2"
          aria-label={t("common.changeLanguage")}
        >
          {showFlag && (
            <span className="text-base" role="img" aria-label={currentLangConfig.name}>
              {currentLangConfig.flag}
            </span>
          )}
          <span className="uppercase text-xs font-semibold">
            {currentLangConfig.code}
          </span>
        </button>
        <ul
          tabIndex={0}
          className="dropdown-content menu bg-base-200 rounded-box z-50 mt-2 w-48 p-2 shadow-lg"
        >
          {languages.map((lang) => (
            <li key={lang.code}>
              <button
                onClick={() => {
                  changeLanguage(lang.code);
                }}
                className={`flex items-center gap-3 ${
                  language === lang.code ? "active" : ""
                }`}
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
    );
  }

  // Default dropdown variant - full featured
  return (
    <div className={`dropdown dropdown-end ${className}`}>
      <button
        tabIndex={0}
        className="btn btn-ghost gap-3"
        aria-label={t("common.changeLanguage")}
      >
        <Languages className="size-5" />
        {showFlag && (
          <span className="text-xl" role="img" aria-label={currentLangConfig.name}>
            {currentLangConfig.flag}
          </span>
        )}
        {showLabel && (
          <span className="font-medium">{currentLangConfig.nativeName}</span>
        )}
      </button>
      <ul
        tabIndex={0}
        className="dropdown-content menu bg-base-200 rounded-box z-50 mt-3 w-64 p-2 shadow-xl"
      >
        <li className="menu-title px-4 py-2">
          <span className="text-base-content/70 text-sm font-semibold">
            {t("common.selectLanguage")}
          </span>
        </li>
        {languages.map((lang) => (
          <li key={lang.code}>
            <button
              onClick={() => {
                changeLanguage(lang.code);
              }}
              className={`flex items-center gap-3 py-3 ${
                language === lang.code ? "active bg-primary/10" : ""
              }`}
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
  );
}
