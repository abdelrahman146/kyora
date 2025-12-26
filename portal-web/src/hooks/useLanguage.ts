import { useCallback, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { setCookie, getCookie } from "@/lib/cookies";

const LANGUAGE_COOKIE = "kyora_language";
const SUPPORTED_LANGUAGES = ["en", "ar"] as const;
type SupportedLanguage = (typeof SUPPORTED_LANGUAGES)[number];

function getBrowserLanguage(): SupportedLanguage {
  const browserLang = navigator.language.split("-")[0];
  return SUPPORTED_LANGUAGES.includes(browserLang as SupportedLanguage)
    ? (browserLang as SupportedLanguage)
    : "en";
}

function getInitialLanguage(): SupportedLanguage {
  const savedLanguage = getCookie(LANGUAGE_COOKIE);

  if (
    savedLanguage &&
    SUPPORTED_LANGUAGES.includes(savedLanguage as SupportedLanguage)
  ) {
    return savedLanguage as SupportedLanguage;
  }

  return getBrowserLanguage();
}

function updateDocumentDirection(language: SupportedLanguage): void {
  const isRTL = language === "ar";
  document.documentElement.dir = isRTL ? "rtl" : "ltr";
  document.documentElement.lang = language;
}

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
      setCookie(LANGUAGE_COOKIE, language);
      updateDocumentDirection(language);
    },
    [i18n]
  );

  const currentLanguage = i18n.language as SupportedLanguage;
  const isRTL = currentLanguage === "ar";

  return {
    currentLanguage,
    isRTL,
    changeLanguage,
    supportedLanguages: SUPPORTED_LANGUAGES,
  };
}
