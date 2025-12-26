import { useCallback, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { setCookie, getCookie } from "@/lib/cookies";

const LANGUAGE_COOKIE = "kyora_language";
const SUPPORTED_LANGUAGES = ["en", "ar"] as const;
type SupportedLanguage = (typeof SUPPORTED_LANGUAGES)[number];

/**
 * Get browser language preference
 * Returns 'ar' if any Arabic locale is detected, otherwise 'en'
 */
function getBrowserLanguage(): SupportedLanguage {
  // Check primary language
  const primaryLang = navigator.language.split("-")[0];
  if (primaryLang === "ar") {
    return "ar";
  }

  // Check all preferred languages
  const languages =
    navigator.languages.length > 0 ? navigator.languages : [navigator.language];
  for (const lang of languages) {
    const code = lang.split("-")[0];
    if (code === "ar") {
      return "ar";
    }
  }

  // Fallback to English
  return "en";
}

/**
 * Get initial language from cookie or browser
 * Priority: 1. Cookie, 2. Browser (if Arabic), 3. English fallback
 */
function getInitialLanguage(): SupportedLanguage {
  // Check cookie first
  const savedLanguage = getCookie(LANGUAGE_COOKIE);
  if (
    savedLanguage &&
    SUPPORTED_LANGUAGES.includes(savedLanguage as SupportedLanguage)
  ) {
    return savedLanguage as SupportedLanguage;
  }

  // Check browser language
  return getBrowserLanguage();
}

function updateDocumentDirection(language: SupportedLanguage): void {
  const isRTL = language === "ar";
  document.documentElement.dir = isRTL ? "rtl" : "ltr";
  document.documentElement.lang = language;
}

/**
 * Language Management Hook
 *
 * Provides centralized language management with:
 * - Cookie persistence for user preference
 * - Browser language detection
 * - Automatic dir/lang attribute updates
 * - Type-safe language switching
 *
 * Priority order:
 * 1. Cookie value (user preference)
 * 2. Browser language (if Arabic)
 * 3. Fallback to English
 *
 * Usage:
 * ```tsx
 * const { language, changeLanguage, toggleLanguage, isRTL } = useLanguage();
 *
 * // Switch language
 * changeLanguage('en');
 *
 * // Toggle between Arabic and English
 * toggleLanguage();
 *
 * // Check if current language is RTL
 * if (isRTL) { ... }
 * ```
 */
export function useLanguage() {
  const { i18n } = useTranslation();

  useEffect(() => {
    const initialLanguage = getInitialLanguage();

    if (i18n.language !== initialLanguage) {
      void i18n.changeLanguage(initialLanguage);
    }

    updateDocumentDirection(initialLanguage);
  }, [i18n]);

  const changeLanguage = useCallback(
    (language: SupportedLanguage) => {
      void i18n.changeLanguage(language);
      setCookie(LANGUAGE_COOKIE, language, 365);
      updateDocumentDirection(language);
    },
    [i18n]
  );

  const toggleLanguage = useCallback(() => {
    const newLang = i18n.language === "ar" ? "en" : "ar";
    changeLanguage(newLang);
  }, [i18n.language, changeLanguage]);

  const currentLanguage = i18n.language as SupportedLanguage;
  const isRTL = currentLanguage === "ar";

  return {
    /** Current language code */
    language: currentLanguage,
    /** Current language code (alias) */
    currentLanguage,
    /** Whether current language is RTL */
    isRTL,
    /** Whether current language is Arabic */
    isArabic: currentLanguage === "ar",
    /** Whether current language is English */
    isEnglish: currentLanguage === "en",
    /** Change language to specified locale */
    changeLanguage,
    /** Toggle between Arabic and English */
    toggleLanguage,
    /** Array of supported language codes */
    supportedLanguages: SUPPORTED_LANGUAGES,
  };
}
