import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import enCommon from "./locales/en/common.json";
import arCommon from "./locales/ar/common.json";
import enOnboarding from "./locales/en/onboarding.json";
import arOnboarding from "./locales/ar/onboarding.json";

const resources = {
  en: {
    common: enCommon,
    onboarding: enOnboarding,
  },
  ar: {
    common: arCommon,
    onboarding: arOnboarding,
  },
};

void i18n.use(initReactI18next).init({
  resources,
  fallbackLng: "en",
  defaultNS: "common",
  ns: ["common", "onboarding"],
  interpolation: {
    escapeValue: false,
  },
  react: {
    useSuspense: false,
  },
});

export default i18n;
