---
description: i18n & Translations SSOT (portal-web + storefront-web + backend-facing rules)
applyTo: "portal-web/**"
---

# i18n & Translations (SSOT)

**SSOT hierarchy**

- Parent: `.github/copilot-instructions.md`
- Portal architecture: `.github/instructions/portal-web-architecture.instructions.md`
- Portal HTTP + Query error-to-toast: `.github/instructions/http-tanstack-query.instructions.md`
- Backend error contract: `.github/instructions/errors-handling.instructions.md`

This doc exists because translation mistakes are a common failure mode (wrong namespace, missing locale keys, duplicated keys). Follow it strictly.

---

## 0) Non-negotiables

1. **No hard-coded user-facing strings** in UI.

- Any text shown to users must be translated via i18n keys (Arabic-first, English fallback).
- Exceptions should be extremely rare (e.g. purely technical debug logs).

2. **No duplicated translation keys.**

- Do not duplicate the same meaning/string across multiple namespaces or multiple places.
- “Shared” strings must live in one shared place (usually the `common` namespace).

3. **No duplicate keys inside a translation file/object.**

- Duplicate keys are considered invalid even if JSON parsing would “last-write-wins”.

4. **Every key must exist in every locale (en + ar).**

- If a key exists only in one locale, it’s a bug and must be tracked as drift.
- Do not rely on `fallbackLng` to “cover” missing keys.

5. **No legacy translation buckets.**

- Do not create or keep keys in catch-all buckets (e.g. a `translation.json`).
- The target state is: only explicit feature namespaces (see below), and every key lives in exactly one namespace.

---

## 1) Portal-web current i18n reality (what exists today)

### 1.1 Where translations live

Portal-web i18n is initialized in `portal-web/src/i18n/init.ts`.

Loaded namespaces include:

- `common` → `portal-web/src/i18n/{en,ar}/common.json`
- `errors` → `portal-web/src/i18n/{en,ar}/errors.json`
- `onboarding` → `portal-web/src/i18n/{en,ar}/onboarding.json`
- `upload` → `portal-web/src/i18n/{en,ar}/upload.json`
- `analytics` → `portal-web/src/i18n/{en,ar}/analytics.json`
- `inventory` → `portal-web/src/i18n/{en,ar}/inventory.json`
- `orders` → `portal-web/src/i18n/{en,ar}/orders.json`
- `auth` → `portal-web/src/i18n/{en,ar}/auth.json`
- `dashboard` → `portal-web/src/i18n/{en,ar}/dashboard.json`
- `customers` → `portal-web/src/i18n/{en,ar}/customers.json`
- `pagination` → `portal-web/src/i18n/{en,ar}/pagination.json`
- `home` → `portal-web/src/i18n/{en,ar}/home.json`
- `accounting` → `portal-web/src/i18n/{en,ar}/accounting.json`
- `reports` → `portal-web/src/i18n/{en,ar}/reports.json`

Default namespace is `common`.

**File structure reality:** Each namespace requires exactly 2 files:

- `portal-web/src/i18n/en/{namespace}.json` (English translations)
- `portal-web/src/i18n/ar/{namespace}.json` (Arabic translations)

When adding a new namespace:

1. Create both `en/{namespace}.json` and `ar/{namespace}.json`
2. Add the namespace to the resources object in `portal-web/src/i18n/init.ts`
3. Ensure both files have matching key structures

### 1.2 SSOT rule for code: always use an explicit namespace

To prevent “wrong file/wrong key” mistakes:

- Do not call `useTranslation()` without a namespace in new code.
- Always call `useTranslation('<namespace>')` and then use bare keys within that namespace.

Examples:

- ✅ `const { t } = useTranslation('common'); t('save')`
- ✅ `const { t: tErrors } = useTranslation('errors'); tErrors('http.404')`
- ❌ `useTranslation(); t('common.save')` (ambiguous + relies on a competing structure)
- ❌ `t('save', { ns: 'common' })` in UI code (allowed only in shared helpers)

Rationale: portal-web currently has multiple translation patterns in the codebase; explicit namespaces stop new code from making it worse.

---

## 2) Choosing where a new key should live (single correct approach)

### 2.1 Namespace responsibilities

Use these rules when adding new keys:

- `errors`: user-facing error messages, including:
  - HTTP status fallbacks (`http.400`, `http.404`, `http.500`, etc.)
  - Backend error codes mapped via `extensions.code` (e.g., `backend.account.invalid_credentials`)
  - Network errors (`network.timeout`, `network.connection`)
  - Validation errors (`validation.required`, `validation.invalid_email`, etc.)
  - Form errors (`form.minItemsRequired`, etc.)
  - Generic fallbacks (`generic.unexpected`, `generic.load_failed`)
- `common`: truly shared UI primitives:
  - Action buttons: `save`, `cancel`, `delete`, `edit`, `create`, `update`, `submit`, `confirm`, `back`, `next`, `add`, `remove`
  - States: `loading`, `no_results`, `empty_state`, `error_loading`
  - Page titles under `pages` object (see section 3.1.1)
  - Reusable component text that appears across multiple features
  - Generic labels that don't belong to a specific domain
- `onboarding`, `inventory`, `orders`, `analytics`, `upload`, `customers`, `accounting`, `reports`: domain-specific strings for those screens/components.
  - Keys within these namespaces should only contain strings specific to that feature
  - If a string is used in multiple features, move it to `common`
- **Do not use `translation`**. If there is no suitable namespace yet, create a new one and add it to `portal-web/src/i18n/init.ts`.

### 2.2 No duplication across namespaces

If a string is used in more than one feature:

- Put it in `common` and reference it from all features.
- Do not copy it into each feature namespace.

### 2.3 Translation key structure (how to organize keys within files)

Keys within a namespace should be organized hierarchically:

**Flat structure** for simple keys (preferred for most cases):

```json
{
  "save": "Save",
  "cancel": "Cancel",
  "loading": "Loading..."
}
```

**Nested structure** for grouped keys (use when keys naturally group together):

```json
{
  "pages": {
    "dashboard": "Dashboard",
    "inventory": "Inventory",
    "customers": "Customers"
  },
  "validation": {
    "required": "This field is required",
    "invalid_email": "Please enter a valid email"
  }
}
```

**Naming conventions:**

- Use `snake_case` for all keys: `customer_details`, `add_first_customer`, `recurring_expenses`
- Avoid camelCase or PascalCase: ❌ `customerDetails`, ❌ `CustomerDetails`
- Use descriptive names that indicate the context: ✅ `no_customers_message`, ❌ `message1`
- For backend error codes, preserve the exact code structure: `backend.account.invalid_credentials` maps to `extensions.code = "account.invalid_credentials"`

---

## 3) How translations must be called (one correct way)

### 3.1 In routes/components

- Always use `useTranslation('<namespace>')`.
- Prefer aliasing when multiple namespaces are needed:

```ts
const { t: tCommon } = useTranslation("common");
const { t: tOrders } = useTranslation("orders");
```

- Avoid passing `ns` per call in UI code.

#### 3.1.1 Route page titles (SSOT)

Route-level page titles must be **centralized** under `common.pages` so they can be reused everywhere (Header title, breadcrumbs, etc.).

**Current implementation:**

- Routes set `staticData.titleKey` pointing to keys in `common.pages`
- `DashboardLayout` reads `staticData.titleKey` from the matched route
- `DashboardLayout` calls `t(titleKey)` with the `common` namespace

**Rules:**

- In route files, set `staticData.titleKey` to a key like `pages.inventory` (no namespace prefix).
- Add the corresponding translations in:
  - `portal-web/src/i18n/en/common.json` under `pages`
  - `portal-web/src/i18n/ar/common.json` under `pages`
- Do not point route titles at feature namespaces (e.g. `inventory.title`, `customers.title`). Feature `title` keys can still exist for within-feature UI, but route titles must use `common.pages.*`.

**Exception for nested feature pages:**

- Some routes use feature-specific titles for sub-pages (e.g., `accounting:header.capital`, `accounting:header.assets`)
- This is allowed for deep feature navigation where the title is very specific to that feature context
- However, top-level feature routes (dashboard, inventory, customers, orders, etc.) must always use `common.pages.*`

**Example:**

```tsx
// ✅ Correct: Top-level route using common.pages
// portal-web/src/routes/business/$businessDescriptor/inventory/index.tsx
export const Route = createFileRoute(
  "/business/$businessDescriptor/inventory/",
)({
  staticData: {
    titleKey: "pages.inventory", // Resolved from common.pages.inventory
  },
  // ...
});

// ✅ Also correct: Nested feature page using feature namespace
// portal-web/src/routes/business/$businessDescriptor/accounting/capital.tsx
export const Route = createFileRoute(
  "/business/$businessDescriptor/accounting/capital",
)({
  staticData: {
    titleKey: "accounting:header.capital", // Deep feature-specific title
  },
  // ...
});

// ❌ Wrong: Top-level route using feature namespace
export const Route = createFileRoute(
  "/business/$businessDescriptor/inventory/",
)({
  staticData: {
    titleKey: "inventory.title", // Should be pages.inventory
  },
  // ...
});
```

**Current page titles in `common.pages`:**

- `dashboard`, `inventory`, `customers`, `customer_details`, `orders`, `analytics`, `accounting`, `expenses`, `recurring_expenses`, `capital`, `assets`, `billing`, `team`, `settings`, `reports`, `reports_health`, `reports_profit`, `reports_cashflow`

### 3.2 In shared utilities (toasts, error helpers)

Some shared utilities cannot easily know the namespace at compile time.

Allowed patterns:

- `t('key', { ns: 'errors' })` in helper utilities.
- Prefer using existing helpers:
  - `portal-web/src/lib/translateError.ts`
  - `portal-web/src/lib/toast.ts`

---

## 4) Error messages + translations (frontend/backend interaction)

### 4.1 Backend error translation workflow (SSOT)

Backend returns RFC 7807 Problem JSON with:

- `status`: HTTP status code
- `title`: Short error title
- `detail`: Human-readable error message (fallback)
- `extensions.code`: **Machine-readable error code** (required for all errors)

**Portal error translation pipeline:**

1. **Backend emits error** with `extensions.code`:

   ```json
   {
     "status": 401,
     "title": "Unauthorized",
     "detail": "Invalid email or password",
     "extensions": {
       "code": "account.invalid_credentials"
     }
   }
   ```

2. **Portal receives error** via ky HTTP client (see `.github/instructions/ky.instructions.md`)

3. **`parseProblemDetails(error)`** extracts the error code:
   - Priority 1: `extensions.code` (if present) → maps to `errors.backend.<code>`
   - Priority 2: HTTP status code → maps to `errors.http.<status>`
   - Priority 3: Generic fallback → `errors.generic.unexpected`

4. **`translateError(errorResult, t)`** looks up the translation key:
   - Key pattern: `backend.account.invalid_credentials`
   - Namespace: `errors`
   - Full lookup: `t('backend.account.invalid_credentials', { ns: 'errors' })`

5. **Global error handler** (in `portal-web/src/main.tsx`) displays the translated error as a toast

**Implementation files:**

- Error parsing: `portal-web/src/lib/errorParser.ts`
- Error translation: `portal-web/src/lib/translateError.ts`
- Toast display: `portal-web/src/lib/toast.ts`
- Global handler: `portal-web/src/main.tsx` (QueryCache/MutationCache `onError`)

**Translation key structure in `errors.json`:**

```json
{
  "backend": {
    "account.invalid_credentials": "Invalid email or password. Please try again.",
    "account.user_already_exists": "An account with this email already exists.",
    "onboarding.session_expired": "Session has expired. Please start onboarding again.",
    "billing.feature_not_available": "The '{{details.feature}}' feature is not available on your current plan.",
    "order.insufficient_stock": "Insufficient stock for product (variant: {{details.variantId}}). Requested: {{details.requestedQuantity}}, available: {{details.availableQuantity}}.",
    ...
  },
  "http": {
    "400": "Invalid request. Please check your input and try again.",
    "401": "You are not authorized. Please log in again.",
    "403": "You don't have permission to perform this action.",
    "404": "The requested resource was not found.",
    ...
  },
  "generic": {
    "unexpected": "An unexpected error occurred. Please try again.",
    "load_failed": "Failed to load data. Please try again."
  }
}
```

**Interpolation parameters:**

- Backend can include additional context in `extensions` beyond `code`
- These are passed to translation via `{{details.fieldName}}` syntax
- Example: `billing.feature_not_available` uses `{{details.feature}}`

### 4.2 Backend must NOT return "translation keys"

Backend is not locale-aware today.

- Backend returns Problem JSON (`application/problem+json`) with `title/detail` strings.
- Portal maps errors to translation keys via `parseProblemDetails()`.

If you need richer translated error UX:

- Prefer backend emitting stable machine-readable codes (e.g. `extensions.code`) and keep the portal mapping in one place (`portal-web/src/lib/errorParser.ts`).
- Do not implement endpoint-by-endpoint URL string checks long-term.

### 4.3 Portal must show translated, meaningful errors

- Default to translated toasts for user actions (mutations) via global handlers (see `.github/instructions/http-tanstack-query.instructions.md`).
- Only catch errors locally when the UI needs status-specific behavior.

**Error display patterns:**

✅ **Correct: Rely on global handler (default)**

```tsx
const mutation = useMutation({
  mutationFn: (data) => api.create(data),
  // No onError needed - global handler shows translated toast automatically
  onSuccess: () => {
    showSuccessToast(t("created_successfully"));
  },
});
```

✅ **Correct: Opt out for inline errors**

```tsx
const mutation = useMutation({
  mutationFn: (data) => api.create(data),
  meta: { errorToast: "off" }, // Suppress global toast
  onError: async (error) => {
    const translated = await translateErrorAsync(error, t);
    setFormError(translated); // Show inline instead
  },
});
```

❌ **Wrong: Show raw backend detail**

```tsx
catch (error) {
  toast.error(error.message);  // Never show raw error.message
}
```

---

## 5) Language detection and direction handling (SSOT)

### 5.1 Language detection priority (one way only)

Language is detected in this exact order (implemented in `portal-web/src/i18n/init.ts`):

1. **Cookie** (`kyora_language`) — user preference (highest priority)
2. **Browser language** — `navigator.language` / `navigator.languages`
   - If any language starts with `ar`, use Arabic
   - Otherwise, use English
3. **Fallback** — English (`en`)

**Implementation:**

```typescript
// portal-web/src/i18n/init.ts
function detectLanguage(): SupportedLanguage {
  // 1. Check cookie first
  const savedLanguage = getCookie("kyora_language");
  if (savedLanguage === "ar" || savedLanguage === "en") {
    return savedLanguage;
  }

  // 2. Check browser language
  const browserLang = navigator.language.split("-")[0];
  if (browserLang === "ar") return "ar";

  // Check all preferred languages
  for (const lang of navigator.languages) {
    if (lang.startsWith("ar")) return "ar";
  }

  // 3. Fallback to English
  return "en";
}
```

### 5.2 Document attributes (single source of truth)

The current language and text direction are set on `document.documentElement` attributes:

- `document.documentElement.lang` — Current language code (`'en'` or `'ar'`)
- `document.documentElement.dir` — Text direction (`'ltr'` or `'rtl'`)

**Where these are set:**

1. **Initial load** — `portal-web/src/i18n/init.ts` sets both attributes after detecting language
2. **Language change** — `i18n.on('languageChanged')` listener updates both attributes
3. **`useLanguage` hook** — `portal-web/src/hooks/useLanguage.ts` updates both when `changeLanguage()` is called

**Rule: Always read from `document.documentElement` or use `useLanguage` hook**

✅ **Correct: Read from document**

```tsx
const isRTL = document.documentElement.dir === "rtl";
const currentLang = document.documentElement.lang;
```

✅ **Correct: Use `useLanguage` hook (PREFERRED)**

```tsx
const { isRTL, language, isArabic } = useLanguage();
```

❌ **Wrong: Check i18n.language directly with .startsWith() (CRITICAL ANTI-PATTERN)**

```tsx
// NEVER do this - duplicates language detection logic
const { i18n } = useTranslation();
const isArabic = i18n.language.toLowerCase().startsWith("ar");
```

**Why this is wrong:**

- Duplicates language detection logic across components
- Harder to maintain (changes require updating multiple files)
- Violates single source of truth principle
- `useLanguage` already provides this functionality

**Pattern to follow when refactoring:**

```tsx
// Before (WRONG)
import { useTranslation } from "react-i18next";

export function MyComponent() {
  const { i18n } = useTranslation();
  const isArabic = i18n.language.toLowerCase().startsWith("ar");
  const countryName = isArabic ? country.nameAr : country.name;
  // ...
}

// After (CORRECT)
import { useTranslation } from "react-i18next";
import { useLanguage } from "@/hooks/useLanguage";

export function MyComponent() {
  const { t } = useTranslation("namespace"); // namespace-specific translation
  const { isArabic } = useLanguage(); // centralized language state
  const countryName = isArabic ? country.nameAr : country.name;
  // ...
}
```

### 5.3 The `useLanguage` hook (preferred way to access language)

**Location:** `portal-web/src/hooks/useLanguage.ts`

**Purpose:** Centralized language management with cookie persistence and document attribute updates

**Exports:**

```tsx
const {
  language, // Current language: 'en' | 'ar'
  currentLanguage, // Alias for language
  isRTL, // true if Arabic, false if English
  isArabic, // true if language is 'ar'
  isEnglish, // true if language is 'en'
  changeLanguage, // (lang: 'en' | 'ar') => void
  toggleLanguage, // () => void (switches between en/ar)
  supportedLanguages, // ['en', 'ar']
} = useLanguage();
```

**Usage examples:**

```tsx
// Check if RTL
const { isRTL } = useLanguage();
if (isRTL) {
  // Apply RTL-specific styling
}

// Check specific language
const { isArabic } = useLanguage();
const countryName = isArabic ? country.arName : country.enName;

// Change language
const { changeLanguage } = useLanguage();
changeLanguage("ar"); // Switch to Arabic

// Toggle between languages
const { toggleLanguage } = useLanguage();
toggleLanguage(); // Switch between en/ar
```

**What it does internally:**

1. Reads initial language from cookie or browser
2. Updates `document.documentElement.lang` and `document.documentElement.dir`
3. Persists language changes to `kyora_language` cookie (365 days expiry)
4. Synchronizes with i18next `changeLanguage()`

### 5.4 Rules for language/direction detection

**Do:**

- ✅ Use `useLanguage` hook in components (PREFERRED)
- ✅ Read `document.documentElement.dir` for one-off checks (rare cases)
- ✅ Use CSS logical properties when possible (`margin-inline-start`, `padding-inline-end`)

**Don't:**

- ❌ Check `i18n.language` directly with `.startsWith()` or `.toLowerCase()` (duplicates logic)
- ❌ Destructure `i18n` from `useTranslation()` just to check language (use `useLanguage` instead)
- ❌ Create computed language variables like `const language = isArabic ? 'ar' : 'en'` when `useLanguage` already provides this
- Duplicate language detection logic across components
- Set `document.documentElement.lang` or `.dir` outside of i18n init or `useLanguage`
- Create multiple ways to determine language or direction

---

## 6) Drift rules (what to log)

Log an item in `DRIFT_TODO.md` when you find any of these:

- A key exists in `en` but not in `ar` (or vice versa).
- A file duplicates the same key (or duplicates the same content across namespaces).
- UI code uses `useTranslation()` without a namespace (new code), or uses `t('common.*')`/`t('errors.*')` patterns that rely on conflicting structures.
- A route uses `staticData.titleKey` outside `common.pages.*` (exception: deep feature navigation like `accounting:header.capital` is allowed, but top-level routes must use `pages.*`).
- Components check `i18n.language.toLowerCase().startsWith('ar')` instead of using `useLanguage` hook.
- Components set `document.documentElement.lang` or `.dir` outside of i18n init or `useLanguage` hook.
- Storefront or portal diverges in translation system in a way that causes duplication or missing keys.
- Backend error codes in `extensions.code` exist but have no corresponding translation in `errors.backend.*`.
- Translation keys use camelCase or PascalCase instead of snake_case.

---

## 7) Why agents fail "first time" (root cause)

Portal-web is namespace-only.

If an agent or contributor adds a catch-all bucket (or uses `useTranslation()` without specifying a namespace), it usually causes:

- wrong namespace lookups
- duplicated meaning across namespaces
- missing locale parity

SSOT fix: always use explicit namespaces and keep each meaning in exactly one namespace.
