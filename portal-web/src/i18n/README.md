# Internationalization (i18n) Setup

This directory contains the internationalization configuration for the Kyora Portal Web App.

## Structure

```
i18n/
├── config.ts              # i18next configuration
├── locales/
│   ├── ar/
│   │   └── common.json   # Arabic translations
│   └── en/
│       └── common.json   # English translations
└── README.md             # This file
```

## Supported Languages

- **Arabic (ar)**: Primary language, RTL (Right-to-Left)
- **English (en)**: Secondary language, LTR (Left-to-Right), fallback language

## Usage

### In Components

```tsx
import { useTranslation } from 'react-i18next';

function MyComponent() {
  const { t } = useTranslation();
  
  return (
    <div>
      <button>{t('save')}</button>
      <p>{t('loading')}</p>
    </div>
  );
}
```

### Language Switching

Use the custom `useLanguage` hook:

```tsx
import { useLanguage } from '@/hooks/useLanguage';

function LanguageSwitcher() {
  const { currentLanguage, changeLanguage, isRTL } = useLanguage();
  
  return (
    <button onClick={() => changeLanguage(currentLanguage === 'ar' ? 'en' : 'ar')}>
      Switch to {currentLanguage === 'ar' ? 'English' : 'العربية'}
    </button>
  );
}
```

## Features

### Cookie Persistence
- Language selection is persisted in a cookie named `kyora_language`
- Cookie expires after 365 days

### Browser Language Detection
- If no saved language is found, the app detects the browser's language
- Falls back to English if the browser language is not supported

### Automatic Direction Handling
- The `useLanguage` hook automatically updates `document.dir` to `rtl` or `ltr`
- The `document.lang` attribute is also updated for accessibility

### Fallback Language
- English is always the fallback language if a translation key is missing

## Adding New Translations

1. Add the key to both `locales/en/common.json` and `locales/ar/common.json`
2. Use the translation in your component with `t('yourKey')`

Example:
```json
// locales/en/common.json
{
  "welcomeMessage": "Welcome to Kyora"
}

// locales/ar/common.json
{
  "welcomeMessage": "مرحبًا بك في كيورا"
}
```

## Adding New Namespaces

To organize translations by feature:

1. Create new JSON files in each locale folder:
   ```
   locales/en/auth.json
   locales/ar/auth.json
   ```

2. Update `config.ts`:
   ```ts
   import enAuth from './locales/en/auth.json';
   import arAuth from './locales/ar/auth.json';
   
   const resources = {
     en: {
       common: enCommon,
       auth: enAuth,  // Add new namespace
     },
     ar: {
       common: arCommon,
       auth: arAuth,  // Add new namespace
     },
   };
   ```

3. Use in components:
   ```tsx
   const { t } = useTranslation('auth');
   // or
   const { t } = useTranslation(['auth', 'common']);
   ```

## Best Practices

1. **Use nested keys** for organization:
   ```json
   {
     "auth": {
       "login": {
         "title": "Login",
         "button": "Sign In"
       }
     }
   }
   ```

2. **Use interpolation** for dynamic content:
   ```json
   {
     "welcome": "Welcome, {{name}}!"
   }
   ```
   ```tsx
   t('welcome', { name: 'Ahmed' })
   ```

3. **Test RTL layout** when adding new UI components to ensure proper mirroring

4. **Always provide translations** for both languages before deploying
