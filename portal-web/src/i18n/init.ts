import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'
import arCommon from './ar/common.json'
import arErrors from './ar/errors.json'
import arOnboarding from './ar/onboarding.json'
import arUpload from './ar/upload.json'
import arAnalytics from './ar/analytics.json'
import arInventory from './ar/inventory.json'
import arOrders from './ar/orders.json'
import arAuth from './ar/auth.json'
import arDashboard from './ar/dashboard.json'
import arCustomers from './ar/customers.json'
import arPagination from './ar/pagination.json'
import arHome from './ar/home.json'
import enCommon from './en/common.json'
import enErrors from './en/errors.json'
import enOnboarding from './en/onboarding.json'
import enUpload from './en/upload.json'
import enAnalytics from './en/analytics.json'
import enInventory from './en/inventory.json'
import enOrders from './en/orders.json'
import enAuth from './en/auth.json'
import enDashboard from './en/dashboard.json'
import enCustomers from './en/customers.json'
import enPagination from './en/pagination.json'
import enHome from './en/home.json'
import { getCookie } from '@/lib/cookies'

/**
 * Detect initial language from cookie or browser
 * Priority: 1. Cookie, 2. Browser (if Arabic), 3. English fallback
 */
function detectLanguage(): 'ar' | 'en' {
  // 1. Check cookie for saved preference
  const savedLanguage = getCookie('kyora_language')
  if (savedLanguage === 'ar' || savedLanguage === 'en') {
    return savedLanguage
  }

  // 2. Check browser language
  const primaryLang = navigator.language.split('-')[0]
  if (primaryLang === 'ar') {
    return 'ar'
  }

  // Check all preferred languages
  const languages =
    navigator.languages.length > 0 ? navigator.languages : [navigator.language]
  for (const lang of languages) {
    const code = lang.split('-')[0]
    if (code === 'ar') {
      return 'ar'
    }
  }

  // 3. Default to English
  return 'en'
}

// Detect language from cookie or browser
const detectedLanguage = detectLanguage()

// Initialize document attributes before i18n loads
document.documentElement.lang = detectedLanguage
document.documentElement.dir = detectedLanguage === 'ar' ? 'rtl' : 'ltr'

// Initialize i18next with detected language
void i18n.use(initReactI18next).init({
  resources: {
    ar: {
      common: arCommon,
      errors: arErrors,
      onboarding: arOnboarding,
      upload: arUpload,
      analytics: arAnalytics,
      inventory: arInventory,
      orders: arOrders,
      auth: arAuth,
      dashboard: arDashboard,
      customers: arCustomers,
      pagination: arPagination,
      home: arHome,
    },
    en: {
      common: enCommon,
      errors: enErrors,
      onboarding: enOnboarding,
      upload: enUpload,
      analytics: enAnalytics,
      inventory: enInventory,
      orders: enOrders,
      auth: enAuth,
      dashboard: enDashboard,
      customers: enCustomers,
      pagination: enPagination,
      home: enHome,
    },
  },
  lng: detectedLanguage,
  fallbackLng: 'en',
  defaultNS: 'common',
  ns: [
    'common',
    'errors',
    'onboarding',
    'upload',
    'analytics',
    'inventory',
    'orders',
    'auth',
    'dashboard',
    'customers',
    'pagination',
    'home',
  ],
  interpolation: {
    escapeValue: false,
  },
})

// Listen for language changes and update document attributes
i18n.on('languageChanged', (lng) => {
  document.documentElement.dir = lng === 'ar' ? 'rtl' : 'ltr'
  document.documentElement.lang = lng
})

export default i18n
