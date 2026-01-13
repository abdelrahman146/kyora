---
description: i18n & Translations SSOT (portal-web + storefront-web + backend-facing rules)
applyTo: "portal-web/src/i18n/**,portal-web/src/routes/**,portal-web/src/components/**,portal-web/src/lib/**,storefront-web/src/i18n/**,storefront-web/src/pages/**,storefront-web/src/components/**,backend/internal/platform/types/problem/**,backend/internal/platform/response/**,backend/internal/domain/**/handler_http.go"
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

Default namespace is `common`.

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

- `errors`: user-facing error messages, including HTTP status fallbacks and auth/onboarding errors.
- `common`: truly shared UI primitives (Save/Cancel/Delete, generic empty states, reusable component text).
- `onboarding`, `inventory`, `orders`, `analytics`, `upload`: domain-specific strings for those screens/components.
- **Do not use `translation`**. If there is no suitable namespace yet, create a new one (e.g. `auth`, `dashboard`, `customers`) and add it to `portal-web/src/i18n/init.ts`.

### 2.2 No duplication across namespaces

If a string is used in more than one feature:

- Put it in `common` and reference it from all features.
- Do not copy it into each feature namespace.

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

### 3.2 In shared utilities (toasts, error helpers)

Some shared utilities cannot easily know the namespace at compile time.

Allowed patterns:

- `t('key', { ns: 'errors' })` in helper utilities.
- Prefer using existing helpers:
  - `portal-web/src/lib/translateError.ts`
  - `portal-web/src/lib/toast.ts`

---

## 4) Error messages + translations (frontend/backed interaction)

### 4.1 Backend must NOT return “translation keys”

Backend is not locale-aware today.

- Backend returns Problem JSON (`application/problem+json`) with `title/detail` strings.
- Portal maps errors to translation keys via `parseProblemDetails()`.

If you need richer translated error UX:

- Prefer backend emitting stable machine-readable codes (e.g. `extensions.code`) and keep the portal mapping in one place (`portal-web/src/lib/errorParser.ts`).
- Do not implement endpoint-by-endpoint URL string checks long-term.

### 4.2 Portal must show translated, meaningful errors

- Default to translated toasts for user actions (mutations) via global handlers (see `.github/instructions/http-tanstack-query.instructions.md`).
- Only catch errors locally when the UI needs status-specific behavior.

---

## 5) Drift rules (what to log)

Log an item in `DRIFT_TODO.md` when you find any of these:

- A key exists in `en` but not in `ar` (or vice versa).
- A file duplicates the same key (or duplicates the same content across namespaces).
- UI code uses `useTranslation()` without a namespace (new code), or uses `t('common.*')`/`t('errors.*')` patterns that rely on conflicting structures.
- Storefront or portal diverges in translation system in a way that causes duplication or missing keys.

---

## 6) Why agents fail “first time” (root cause)

Portal-web is namespace-only.

If an agent or contributor adds a catch-all bucket (or uses `useTranslation()` without specifying a namespace), it usually causes:

- wrong namespace lookups
- duplicated meaning across namespaces
- missing locale parity

SSOT fix: always use explicit namespaces and keep each meaning in exactly one namespace.
