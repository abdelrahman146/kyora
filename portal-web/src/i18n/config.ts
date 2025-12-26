import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import enCommon from "./locales/en/common.json";
import arCommon from "./locales/ar/common.json";

const resources = {
  en: {
    common: enCommon,
  },
  ar: {
    common: arCommon,
  },
};

void i18n.use(initReactI18next).init({
  resources,
  fallbackLng: "en",
  defaultNS: "common",
  ns: ["common"],
  interpolation: {
    escapeValue: false,
  },
  react: {
    useSuspense: false,
  },
});

export default i18n;
