import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import { getCookie, setCookie } from "../utils/cookie";
import { translations, type SupportedLanguage } from "./translations";

const LANG_COOKIE = "kyora_sf_lang";

function normalizeLang(input: string | undefined): SupportedLanguage {
  const v = (input || "").trim().toLowerCase();
  if (v.startsWith("ar")) return "ar";
  if (v.startsWith("en")) return "en";
  return "ar";
}

function applyDocumentDirection(lang: SupportedLanguage): void {
  const dir = lang === "ar" ? "rtl" : "ltr";
  document.documentElement.dir = dir;
  document.documentElement.lang = lang;
}

export function initI18n(): void {
  if (i18n.isInitialized) return;

  const initial = normalizeLang(getCookie(LANG_COOKIE) || navigator.language);

  i18n
    .use(initReactI18next)
    .init({
      resources: translations,
      lng: initial,
      fallbackLng: "ar",
      defaultNS: "common",
      interpolation: { escapeValue: false },
    })
    .catch(() => {
      // i18n init errors should not block the storefront
    });

  applyDocumentDirection(initial);

  i18n.on("languageChanged", (lng) => {
    const lang = normalizeLang(lng);
    setCookie(LANG_COOKIE, lang);
    applyDocumentDirection(lang);
  });
}

export function supportedLanguages(): Array<{
  code: SupportedLanguage;
  label: string;
}> {
  return [
    { code: "ar", label: "العربية" },
    { code: "en", label: "English" },
  ];
}

export function cookieLanguage(): SupportedLanguage {
  return normalizeLang(getCookie(LANG_COOKIE));
}
