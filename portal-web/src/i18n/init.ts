import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import ar from "./locales/ar/translation.json";
import en from "./locales/en/translation.json";
import arErrors from "./locales/ar/errors.json";
import enErrors from "./locales/en/errors.json";
import { getCookie } from "../lib/cookies";

/**
 * Detect initial language from cookie or browser
 * Priority: 1. Cookie, 2. Browser (if Arabic), 3. English fallback
 */
function detectLanguage(): "ar" | "en" {
  // 1. Check cookie for saved preference
  const savedLanguage = getCookie("kyora_language");
  if (savedLanguage === "ar" || savedLanguage === "en") {
    return savedLanguage;
  }

  // 2. Check browser language
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

  // 3. Default to English
  return "en";
}

// Detect language from cookie or browser
const detectedLanguage = detectLanguage();

// Initialize document attributes before i18n loads
document.documentElement.lang = detectedLanguage;
document.documentElement.dir = detectedLanguage === "ar" ? "rtl" : "ltr";

// Initialize i18next with detected language
void i18n.use(initReactI18next).init({
  resources: {
    ar: {
      translation: ar,
      errors: arErrors,
    },
    en: {
      translation: en,
      errors: enErrors,
    },
  },
  lng: detectedLanguage,
  fallbackLng: "en",
  defaultNS: "translation",
  interpolation: {
    escapeValue: false,
  },
});

// Listen for language changes and update document attributes
i18n.on("languageChanged", (lng) => {
  document.documentElement.dir = lng === "ar" ? "rtl" : "ltr";
  document.documentElement.lang = lng;
});

export default i18n;
