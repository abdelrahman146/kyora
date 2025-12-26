# Language Management System

## Overview

Kyora Portal Web implements a comprehensive language management system with:
- **Cookie persistence** - User language preference saved in cookies
- **Browser language detection** - Automatic detection with Arabic priority
- **Global hook** - `useLanguage()` hook for consistent language management
- **Automatic direction switching** - RTL/LTR handling for Arabic/English
- **Type-safe** - Full TypeScript support with proper types

## Language Detection Priority

The system follows this priority order when determining the initial language:

1. **Cookie value** (`kyora_language`) - User's saved preference
2. **Browser language** - If browser is set to Arabic (`ar`)
3. **English fallback** - Default if neither above applies

## Architecture

### Core Components

#### 1. `useLanguage()` Hook
**Location**: `src/hooks/useLanguage.ts`

The central language management hook that should be used throughout the application.

**API**:
```tsx
const {
  language,        // Current language: "ar" | "en"
  currentLanguage, // Alias for language
  isRTL,          // true if Arabic, false if English
  isArabic,       // true if current language is Arabic
  isEnglish,      // true if current language is English
  changeLanguage, // (lang: "ar" | "en") => void
  toggleLanguage, // () => void - Switch between ar/en
  supportedLanguages, // ["en", "ar"]
} = useLanguage();
```

**Features**:
- Automatically checks cookie and browser language on mount
- Updates cookie when language changes (1 year expiry)
- Updates document `dir` and `lang` attributes
- Triggers i18next language change
- Type-safe language values

#### 2. i18n Initialization
**Location**: `src/i18n/init.ts`

Initializes i18next with language detection logic.

**Flow**:
1. Checks `kyora_language` cookie
2. Falls back to browser language (if Arabic)
3. Falls back to English
4. Sets initial document attributes (`dir`, `lang`)
5. Initializes i18next with detected language
6. Sets up language change listener

#### 3. Cookie Management
**Location**: `src/lib/cookies.ts`

Utilities for reading/writing cookies:
- `getCookie(name)` - Read cookie value
- `setCookie(name, value, days)` - Write cookie with expiry
- `deleteCookie(name)` - Remove cookie

Cookie name: `kyora_language`
Cookie expiry: 365 days
Cookie flags: `SameSite=Lax`, `Secure` (in production)

## Usage Guide

### Basic Language Switching

```tsx
import { useLanguage } from "@/hooks/useLanguage";

function MyComponent() {
  const { language, toggleLanguage, isRTL } = useLanguage();

  return (
    <div>
      <button onClick={toggleLanguage}>
        {language === "ar" ? "English" : "عربي"}
      </button>
      
      {/* Use isRTL for conditional positioning */}
      <Toast position={isRTL ? "top-right" : "top-left"} />
    </div>
  );
}
```

### Change to Specific Language

```tsx
function LanguagePicker() {
  const { changeLanguage, language } = useLanguage();

  return (
    <select 
      value={language} 
      onChange={(e) => changeLanguage(e.target.value as "ar" | "en")}
    >
      <option value="ar">العربية</option>
      <option value="en">English</option>
    </select>
  );
}
```

### Conditional Rendering Based on Language

```tsx
function WelcomeMessage() {
  const { isArabic, isEnglish } = useLanguage();

  return (
    <div>
      {isArabic && <p>مرحباً بك في كيورا</p>}
      {isEnglish && <p>Welcome to Kyora</p>}
    </div>
  );
}
```

### Use with Toast Notifications

```tsx
import toast from "react-hot-toast";
import { useLanguage } from "@/hooks/useLanguage";

function LoginForm() {
  const { isRTL } = useLanguage();

  const handleSubmit = async () => {
    try {
      await login();
      toast.success("Login successful!", {
        position: isRTL ? "top-right" : "top-left",
      });
    } catch (error) {
      toast.error("Login failed", {
        position: isRTL ? "top-right" : "top-left",
      });
    }
  };
}
```

## Implementation Details

### How It Works

1. **App Initialization**:
   - `main.tsx` imports `./i18n/init` before rendering
   - `init.ts` runs language detection and sets initial state
   - Document attributes (`dir`, `lang`) are set immediately
   - i18next is initialized with detected language

2. **Component Mount**:
   - Component calls `useLanguage()`
   - Hook checks cookie and browser language
   - If different from current i18n language, triggers change
   - Returns current state and change functions

3. **Language Change**:
   - User clicks language switcher
   - Calls `changeLanguage("en")` or `toggleLanguage()`
   - Hook updates i18next language
   - Hook saves preference to cookie (365 days)
   - Hook updates document `dir` and `lang` attributes
   - i18next triggers re-render of all components using `useTranslation()`
   - Layout switches between RTL/LTR automatically

4. **Cookie Persistence**:
   - Language preference saved in `kyora_language` cookie
   - Cookie expires in 1 year
   - Cookie is read on every app load
   - User sees their preferred language immediately

### Browser Language Detection Logic

```typescript
function getBrowserLanguage(): "ar" | "en" {
  // Check primary language
  const primaryLang = navigator.language.split("-")[0];
  if (primaryLang === "ar") return "ar";

  // Check all preferred languages
  for (const lang of navigator.languages) {
    if (lang.startsWith("ar")) return "ar";
  }

  // Fallback to English
  return "en";
}
```

**Examples**:
- Browser: `ar-SA` → Detects `ar` ✅
- Browser: `en-US` → Falls back to `en` ✅
- Browser: `fr-FR` → Falls back to `en` ✅
- Browser Languages: `[en-US, ar-EG]` → Detects `ar` ✅

## Updated Components

All language switchers have been updated to use `useLanguage()`:

### 1. Login Page (`src/routes/login.tsx`)
```tsx
const { language, toggleLanguage, isRTL } = useLanguage();

// Toast positioning
<Toaster position={isRTL ? "top-right" : "top-left"} />

// Language switcher button
<button onClick={toggleLanguage}>
  {language === "ar" ? "English" : "العربية"}
</button>
```

### 2. Dashboard Page (`src/routes/dashboard.tsx`)
```tsx
const { language, toggleLanguage } = useLanguage();

<button onClick={toggleLanguage}>
  {language === "ar" ? "English" : "العربية"}
</button>
```

### 3. Home Page (`src/App.tsx`)
```tsx
const { language, toggleLanguage } = useLanguage();

<button onClick={toggleLanguage}>
  {language === 'ar' ? 'English' : 'عربي'}
</button>
```

## Testing

### Test Scenarios

#### Test 1: Cookie Persistence
1. Open app (language should be English by default for most browsers)
2. Switch to Arabic
3. Refresh page → Should load in Arabic ✅
4. Close browser, reopen → Should still be Arabic ✅

#### Test 2: Browser Language Detection (No Cookie)
1. Clear cookies
2. Set browser language to Arabic (`ar-SA`)
3. Open app → Should load in Arabic ✅
4. Set browser language to English
5. Clear cookies, reload → Should load in English ✅

#### Test 3: Language Switching
1. Open app in Arabic
2. Click "English" → Text changes to English, layout becomes LTR ✅
3. Click "عربي" → Text changes to Arabic, layout becomes RTL ✅
4. Check cookie → Should be updated ✅

#### Test 4: RTL/LTR Layout
1. Open app in Arabic
2. Check document attributes: `dir="rtl"`, `lang="ar"` ✅
3. Switch to English
4. Check document attributes: `dir="ltr"`, `lang="en"` ✅

#### Test 5: Toast Positioning
1. Open login page in Arabic
2. Trigger error → Toast appears top-right ✅
3. Switch to English
4. Trigger error → Toast appears top-left ✅

### Manual Testing Commands

```bash
# Check cookie value (in browser console)
document.cookie

# Check document attributes
document.documentElement.dir   // "rtl" or "ltr"
document.documentElement.lang  // "ar" or "en"

# Check current language (in React DevTools)
// Find component using useLanguage
// Check value of `language`

# Clear cookie and test detection
document.cookie = "kyora_language=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;"
location.reload()
```

## Migration Guide

### For Existing Code

**Before** (Old pattern - DO NOT USE):
```tsx
import { useTranslation } from "react-i18next";

function Component() {
  const { i18n } = useTranslation();

  const toggleLang = () => {
    const newLang = i18n.language === "ar" ? "en" : "ar";
    void i18n.changeLanguage(newLang);
  };

  return (
    <button onClick={toggleLang}>
      {i18n.language === "ar" ? "English" : "عربي"}
    </button>
  );
}
```

**After** (New pattern - USE THIS):
```tsx
import { useLanguage } from "@/hooks/useLanguage";

function Component() {
  const { language, toggleLanguage } = useLanguage();

  return (
    <button onClick={toggleLanguage}>
      {language === "ar" ? "English" : "عربي"}
    </button>
  );
}
```

## Benefits

✅ **Consistent behavior** - All language switching uses same logic
✅ **Cookie persistence** - User preference saved automatically
✅ **Browser detection** - Works without manual setup
✅ **Type safety** - TypeScript ensures correct usage
✅ **DRY principle** - No duplicate language switching logic
✅ **Maintainable** - Single source of truth for language management
✅ **Testable** - Clear separation of concerns

## Configuration

### Cookie Settings

To modify cookie behavior, edit `src/hooks/useLanguage.ts`:

```typescript
const LANGUAGE_COOKIE = "kyora_language"; // Cookie name
const cookieExpiry = 365; // Days (in setCookie call)
```

### Supported Languages

To add more languages, edit:

1. `src/hooks/useLanguage.ts`:
```typescript
const SUPPORTED_LANGUAGES = ["en", "ar", "fr"] as const; // Add "fr"
```

2. `src/i18n/init.ts`:
```typescript
// Add French translations
import fr from "./locales/fr/translation.json";

resources: {
  ar: { ... },
  en: { ... },
  fr: { ... }, // Add this
}
```

## Troubleshooting

### Language not persisting after page refresh
- Check browser console for cookie errors
- Verify cookie is being set: `document.cookie`
- Check if browser blocks cookies (privacy settings)

### Wrong language detected on first load
- Check browser language settings
- Verify cookie is cleared
- Check `detectLanguage()` logic in `init.ts`

### Direction not changing when switching languages
- Check if `dir` attribute is being updated: `document.documentElement.dir`
- Verify language change event is firing
- Check if `updateDocumentDirection()` is called

### TypeScript errors with useLanguage
- Ensure you're using the correct return type
- Check that language value is typed as `"ar" | "en"`
- Verify imports are correct

## Future Enhancements

Potential improvements:

- [ ] Add more languages (French, Spanish, etc.)
- [ ] Language preference in user profile (sync across devices)
- [ ] Automatic language detection based on location/IP
- [ ] Language-specific date/number formatting
- [ ] Translation management UI for non-developers
- [ ] Lazy loading of translation files
- [ ] Translation cache optimization

## Summary

The language management system is now:
- ✅ Cookie-based with automatic persistence
- ✅ Browser language detection with Arabic priority
- ✅ Managed by a single global hook (`useLanguage`)
- ✅ Type-safe and consistent across all components
- ✅ Automatically updates document direction (RTL/LTR)
- ✅ Integrated with all existing language switchers

All components should use `useLanguage()` hook instead of directly accessing i18next for language management.
