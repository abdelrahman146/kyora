import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import ar from "./locales/ar/translation.json";
import en from "./locales/en/translation.json";
import arErrors from "./locales/ar/errors.json";
import enErrors from "./locales/en/errors.json";

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
  lng: "ar",
  fallbackLng: "en",
  defaultNS: "translation",
  interpolation: {
    escapeValue: false,
  },
});

i18n.on("languageChanged", (lng) => {
  document.documentElement.dir = lng === "ar" ? "rtl" : "ltr";
  document.documentElement.lang = lng;
});

document.documentElement.dir = i18n.language === "ar" ? "rtl" : "ltr";
document.documentElement.lang = i18n.language;

export default i18n;
