# Language Switcher Component

Universal language switcher component with multiple design variants.

## Usage

```tsx
import { LanguageSwitcher } from "@/components/molecules/LanguageSwitcher";

// Full dropdown (default)
<LanguageSwitcher />

// Compact version for navbars
<LanguageSwitcher variant="compact" />

// Icon only
<LanguageSwitcher variant="iconOnly" />

// Simple toggle (best for 2 languages)
<LanguageSwitcher variant="toggle" />
```

## Variants

### Dropdown (Default)

Full-featured dropdown with language names, native names, and flags. Best for settings pages or dedicated language selection areas.

```tsx
<LanguageSwitcher variant="dropdown" />
```

### Compact

Compact version perfect for navbars and tight spaces. Shows current language code and flag.

```tsx
<LanguageSwitcher variant="compact" />
```

### Icon Only

Shows only a globe icon. Dropdown opens on click with full language list.

```tsx
<LanguageSwitcher variant="iconOnly" />
```

### Toggle

Simple toggle button to switch between 2 languages. Best when you have exactly 2 supported languages.

```tsx
<LanguageSwitcher variant="toggle" />
```

## Props

| Prop | Type | Default | Description |
| ------ | ------ | ------- | ----------- |
| `variant` | `"dropdown" \| "toggle" \| "compact" \| "iconOnly"` | `"dropdown"` | Design variant to use |
| `className` | `string` | `""` | Additional CSS classes |
| `showLabel` | `boolean` | `true` | Show language name label |
| `showFlag` | `boolean` | `true` | Show flag emoji |

## Adding New Languages

To add a new language:

1. Add the language configuration in `LanguageSwitcher.tsx`:

   ```tsx
   const LANGUAGE_CONFIG = {
     // ... existing languages
     fr: {
       code: "fr",
       name: "French",
       nativeName: "Français",
       flag: "🇫🇷",
       direction: "ltr",
     },
   } as const;
   ```

1. Update `useLanguage` hook to include the new language in `SUPPORTED_LANGUAGES`:

   ```tsx
   const SUPPORTED_LANGUAGES = ["en", "ar", "fr"] as const;
   ```

1. Add translation files in `src/i18n/locales/fr/`

1. Update i18n config in `src/i18n/config.ts`:

   ```tsx
   import frCommon from "./locales/fr/common.json";

   const resources = {
     // ... existing languages
     fr: {
       common: frCommon,
       // ... other namespaces
     },
   };
   ```

## Features

- ✅ Multiple design variants for different use cases
- ✅ RTL-aware dropdown positioning
- ✅ Accessible keyboard navigation
- ✅ Flag emoji support
- ✅ Native language names
- ✅ Persistent language selection (cookie-based)
- ✅ Automatic document direction update
- ✅ TypeScript type-safe
- ✅ Mobile-responsive
- ✅ DaisyUI styled

## Examples

### In Navigation Header

```tsx
<header>
  <Logo />
  <nav>
    {/* navigation links */}
  </nav>
  <LanguageSwitcher variant="compact" />
</header>
```

### In Settings Page

```tsx
<div className="card">
  <h2>Language Preferences</h2>
  <LanguageSwitcher variant="dropdown" />
</div>
```

### In Mobile Menu

```tsx
<div className="drawer-side">
  <ul className="menu">
    {/* menu items */}
    <li>
      <LanguageSwitcher variant="iconOnly" />
    </li>
  </ul>
</div>
```
