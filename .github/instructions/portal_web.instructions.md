---
description: Kyora Portal Web App - AI Agent Master Instructions
applyTo: "portal-web/**"
---

# Kyora Portal Web App - Master Instructions

This is the **comprehensive, authoritative instruction file** for the AI Agent maintaining the Kyora Portal Web App. All patterns, conventions, and implementation details are consolidated here.

---

## 1. Project Philosophy & Core Requirements

### 1.1 Target Audience & Value Proposition

**Users**: Busy social media business owners and entrepreneurs who:

- Sell products on Instagram, Facebook, TikTok, WhatsApp
- Need simple, powerful tools to manage orders, inventory, and finances
- Are often on mobile, switching between devices frequently
- May not have accounting or technical expertise

**Value**: Kyora provides professional business management without complexity. The Portal must reflect reliability, speed, and simplicity.

### 1.2 Design Philosophy

**Mobile-First & RTL-First**:

- Arabic is the primary language (RTL layout)
- English is fully supported (LTR layout)
- Every UI decision prioritizes mobile usability
- Minimum touch target: 44x44px
- Use logical CSS properties: `start-*`, `end-*`, `ms-*`, `me-*`, `ps-*`, `pe-*`
- **Never** use `left`, `right`, `margin-left`, `padding-right`, `float-left`

**Zero-State & Feedback**:

- Loading: Skeleton screens (pulse animations) matching content shape
- Actions: Toast notifications for every Create/Update/Delete
- Errors: Error boundaries for sections, never crash the app
- Forms: Bottom Sheets on mobile, modals/side drawers on desktop
- All toasts: RTL-aware positioning (top-right for Arabic, top-left for English)

**Production Standards**:

- No `// TODO` or `// FIXME` - implement complete features
- TypeScript `strict: true` - no `any` types
- Zod schemas for all data validation
- Atomic Design pattern for components
- Clean Architecture: UI, State, API separation

---

## 2. Technology Stack

| Category      | Technology              | Version Notes                       |
| ------------- | ----------------------- | ----------------------------------- |
| Framework     | React Router v7         | Framework mode with loaders/actions |
| Language      | TypeScript              | Strict mode enabled                 |
| Build Tool    | Vite                    | Latest                              |
| UI Library    | DaisyUI v5              | Based on Tailwind CSS v4            |
| Styling       | Tailwind CSS v4         | CSS-first configuration             |
| Forms         | React Hook Form         | With zod resolver                   |
| Validation    | Zod                     | All requests/responses              |
| i18n          | i18next + react-i18next | Arabic primary, English fallback    |
| HTTP Client   | ky                      | Production-grade with retry/refresh |
| Icons         | lucide-react            | Consistent icon set                 |
| Date Handling | date-fns                | With Arabic locale support          |
| Toast         | react-hot-toast         | RTL-aware positioning               |
| State         | Zustand                 | UI state only (not server state)    |

---

## 3. Architecture & Folder Structure

```
portal-web/src/
├── api/                      # API layer (ky-based)
│   ├── client.ts            # Centralized HTTP client with auth/retry/refresh
│   ├── auth.ts              # Authentication endpoints
│   ├── user.ts              # User endpoints
│   ├── business.ts          # Business endpoints
│   ├── onboarding.ts        # Onboarding flow endpoints
│   └── types/               # Zod schemas + TypeScript types
│       ├── auth.ts
│       ├── user.ts
│       └── index.ts
├── assets/                   # Static assets
├── components/               # Atomic Design structure
│   ├── atoms/               # Button, Input, Badge, Avatar, Modal, Skeleton
│   ├── molecules/           # FormInput, BottomSheet, LanguageSwitcher, Toast
│   ├── organisms/           # LoginForm, DataTable, Sidebar, Header
│   ├── templates/           # OnboardingLayout, DashboardLayout
│   └── routing/             # RequireAuth, route guards
├── contexts/                 # React contexts (use sparingly)
│   ├── AuthContext.tsx
│   └── OnboardingContext.tsx
├── hooks/                    # Custom hooks
│   ├── useAuth.tsx          # Authentication state/actions
│   ├── useLanguage.ts       # Language switching + RTL/LTR
│   └── useMediaQuery.ts     # Responsive breakpoints
├── i18n/                     # Internationalization
│   ├── init.ts              # i18next initialization
│   ├── config.ts            # i18next configuration
│   └── locales/
│       ├── ar/              # Arabic translations (primary)
│       │   ├── common.json
│       │   ├── auth.json
│       │   ├── onboarding.json
│       │   └── errors.json
│       └── en/              # English translations (fallback)
│           ├── common.json
│           ├── auth.json
│           ├── onboarding.json
│           └── errors.json
├── lib/                      # Utilities
│   ├── cookies.ts           # Cookie helpers (Secure flag in production)
│   ├── errorParser.ts       # ProblemDetails (RFC 7807) parser
│   ├── translateError.ts    # Error translation helper
│   ├── phone.ts             # Phone number utilities
│   └── utils.ts             # General utilities
├── routes/                   # React Router v7 routes
│   ├── login.tsx
│   ├── dashboard.tsx
│   ├── onboarding/          # Multi-step onboarding flow
│   │   ├── index.tsx        # Layout + guards
│   │   ├── plan.tsx         # Step 1: Plan selection
│   │   ├── verify.tsx       # Step 2: Email verification
│   │   ├── business.tsx     # Step 3: Business setup
│   │   ├── payment.tsx      # Step 4: Stripe checkout
│   │   └── complete.tsx     # Step 5: Finalization
│   └── dashboard/           # Protected dashboard routes
├── schemas/                  # Zod schemas for forms
│   └── auth.ts
├── stores/                   # Zustand stores (UI state only)
│   ├── businessStore.ts
│   └── metadataStore.ts
└── types/                    # Global TypeScript types
    └── index.ts
```

---

## 4. UI/UX & Design System

### 4.1 Colors (DaisyUI + Branding)

Follow `.github/instructions/branding.instructions.md` strictly:

- Primary: `#0D9488` (Teal-600)
- Secondary: `#EAB308` (Gold)
- Base colors: `base-100`, `base-200`, `base-300`
- Semantic: `success`, `error`, `warning`, `info`
- Use DaisyUI tokens: `bg-primary`, `text-primary-content`, etc.

### 4.2 Typography

- Font: `IBM Plex Sans Arabic` (Google Fonts)
- Fallback: `Almarai`
- Scale: Mobile-first
  - Display: 32px/1.2/Bold
  - H1: 24px/1.3/Bold
  - H2: 20px/1.3/SemiBold
  - Body: 16px/1.5/Regular
  - Caption: 12px/1.4/Medium

### 4.3 Spacing (4px baseline)

- Gap-XS: 4px, Gap-S: 8px, Gap-M: 16px, Gap-L: 24px, Gap-XL: 32px
- Safe area padding: 16px left/right
- Consistent use of `p-4`, `gap-2`, `space-y-4`

### 4.4 Border Radius

- `rounded-sm`: 4px (Tags, Checkboxes)
- `rounded-md`: 8px (Inputs, Inner cards)
- `rounded-lg`: 12px (Cards, Modals)
- `rounded-xl`: 16px (Bottom Sheets, Buttons)
- `rounded-full`: Pills, Avatars

### 4.5 Shadows

- `shadow-sm`: `0 1px 2px 0 rgb(0 0 0 / 0.05)` (Cards)
- `shadow-float`: `0 10px 15px -3px rgb(0 0 0 / 0.1)` (Floating elements)

### 4.6 Layouts

**App Shell (Desktop)**:

- Left Sidebar: Collapsible, logo at top, navigation items, user profile at bottom
- Top Header: Business Switcher, Search, Notifications, User Menu
- Main Content: Full height, scrollable, safe area padding

**App Shell (Mobile)**:

- Bottom Navigation: 4-5 primary items, icon + label
- Top Header: Business Switcher (or current context), actions
- Hamburger Menu: For secondary navigation ("More")
- Main Content: Full screen, safe area padding

**Business Switcher**:

- Located in header (both mobile and desktop)
- Dropdown with business list
- Shows current business name + avatar
- Changing business re-fetches data for new context

### 4.7 Interactions

**Bottom Sheets** (use `BottomSheet` component):

- Mobile: Slides up from bottom, 85% max height
- Desktop: Side drawer (left/right), configurable width
- Use for: Create, Edit, Filter actions
- Sizes: `sm` (384px), `md` (448px), `lg` (512px), `xl` (576px), `full`

**Modals** (use `Modal` component):

- Mobile: Bottom sheet behavior
- Desktop: Centered modal
- Use for: Confirmations, simple forms, alerts
- Sizes: `sm`, `md`, `lg`, `xl`, `full`

**Toast Notifications** (react-hot-toast):

- Position: Top-right (RTL), Top-left (LTR)
- Use `useLanguage()` hook to determine position
- Types: `success`, `error`, `info`, `warning`
- Auto-dismiss: 4 seconds

---

## 5. Authentication System

### 5.1 Strategy

- **Access Token**: In-memory (cleared on page refresh)
- **Refresh Token**: Secure cookie (`kyora_refresh_token`, 365 days, `SameSite=Lax`, `Secure` in production)
- **Auto-Refresh**: 401 detection → call `/v1/auth/refresh` → retry request
- **Token Storage**: Never use localStorage for tokens

### 5.2 Login Flow

```
1. User submits credentials
2. POST /v1/auth/login → { token, refreshToken, user }
3. setTokens(token, refreshToken) → saves to memory + cookie
4. Check user.hasCompletedOnboarding
5. Redirect to /dashboard OR /onboarding/plan
```

### 5.3 Session Restoration

```
1. App mounts
2. Check refresh token in cookie
3. If exists: POST /v1/auth/refresh → get new access token
4. Fetch user profile: GET /v1/users/me
5. Set user in AuthContext
6. User is logged in
```

### 5.4 API Client (`src/api/client.ts`)

```typescript
// Import and use
import apiClient, {
  setTokens,
  clearTokens,
  getAccessToken,
  hasValidToken,
} from "@/api/client";

// After login
setTokens(accessToken, refreshToken);

// Check auth
if (hasValidToken()) {
  /* user is logged in */
}

// Make authenticated request (automatic Bearer token)
const orders = await apiClient.get("v1/orders").json();

// On 401:
// - Client detects 401
// - Calls /v1/auth/refresh with cookie
// - Updates access token
// - Retries original request
// - If refresh fails → clearTokens() → navigate('/login')
```

### 5.5 useAuth Hook

```typescript
import { useAuth } from "@/hooks/useAuth";

const {
  user, // User | null
  isAuthenticated, // boolean
  isLoading, // boolean
  login, // (credentials) => Promise<void>
  logout, // () => Promise<void>
  logoutAll, // () => Promise<void> (all devices)
} = useAuth();
```

### 5.6 Route Guard

```typescript
import { RequireAuth } from "@/components/routing/RequireAuth";

// Protect routes
<Route
  element={
    <RequireAuth>
      <AppLayout />
    </RequireAuth>
  }
>
  <Route path="/dashboard" element={<Dashboard />} />
  <Route path="/orders" element={<Orders />} />
</Route>;
```

---

## 6. Onboarding Flow

### 6.1 Backend Session-Based Flow

```
POST /v1/onboarding/start → { sessionToken, stage, isPaid }
POST /v1/onboarding/email/otp
POST /v1/onboarding/email/verify
POST /v1/onboarding/oauth/google
POST /v1/onboarding/business
POST /v1/onboarding/payment/start → { checkoutUrl }
[User completes Stripe checkout]
[Webhook updates stage]
POST /v1/onboarding/complete → { user, token, refreshToken }
```

### 6.2 State Management

- Use `OnboardingContext` (sessionStorage persistence)
- URL-based navigation: `/onboarding/plan`, `/onboarding/verify`, etc.
- Session recovery across page refreshes
- Auto-cleanup on completion

### 6.3 Steps

1. **Plan Selection** (`/onboarding/plan`)

   - User selects billing plan
   - User enters email
   - POST /v1/onboarding/start → sessionToken

2. **Email Verification** (`/onboarding/verify`)

   - Option A: Email OTP → POST /v1/onboarding/email/otp → verify code
   - Option B: Google OAuth → redirect to Google → POST /v1/onboarding/oauth/google
   - User provides first name, last name, password

3. **Business Setup** (`/onboarding/business`)

   - User enters business name
   - Auto-generated descriptor (validated format)
   - User selects country + currency

4. **Payment** (`/onboarding/payment`) - Paid plans only

   - POST /v1/onboarding/payment/start → checkoutUrl
   - Redirect to Stripe Checkout
   - User completes payment
   - Return with `?status=success` or `?status=cancelled`

5. **Completion** (`/onboarding/complete`)
   - POST /v1/onboarding/complete
   - Backend creates workspace, user, business, subscription
   - Returns JWT tokens
   - Auto-login → redirect to /dashboard

### 6.4 Onboarding Layout

- Minimal layout (no sidebar)
- Progress indicator (percentage)
- Language switcher (icon-only variant)
- Mobile-first, RTL-ready

---

## 7. Internationalization (i18n)

### 7.1 Language Detection Priority

1. Cookie (`kyora_language`) - User's saved preference
2. Browser language (if Arabic: `navigator.languages` includes `ar*`)
3. English fallback

### 7.2 useLanguage Hook

```typescript
import { useLanguage } from "@/hooks/useLanguage";

const {
  language, // "ar" | "en"
  isRTL, // boolean
  isArabic, // boolean
  isEnglish, // boolean
  changeLanguage, // (lang: "ar" | "en") => void
  toggleLanguage, // () => void
} = useLanguage();

// Usage
<Toast position={isRTL ? "top-right" : "top-left"} />;
```

### 7.3 Translation Usage

```typescript
import { useTranslation } from "react-i18next";

// Single namespace
const { t } = useTranslation();
t("common.save"); // or t("common:save")

// Multiple namespaces
const { t } = useTranslation(["auth", "common"]);
t("auth:login.title");
t("common:cancel");

// With interpolation
t("auth:welcome", { name: user.firstName });
```

### 7.4 Translation File Structure

```
locales/
├── ar/
│   ├── common.json       # Shared terms (save, cancel, loading, etc.)
│   ├── auth.json         # Authentication (login, register, etc.)
│   ├── onboarding.json   # Onboarding flow
│   ├── errors.json       # Error messages
│   └── dashboard.json    # Dashboard-specific
└── en/
    └── (same structure)
```

### 7.5 Language Switcher Component

```typescript
import { LanguageSwitcher } from "@/components/molecules/LanguageSwitcher";

// Variants:
<LanguageSwitcher variant="dropdown" />   // Full dropdown (settings pages)
<LanguageSwitcher variant="compact" />    // Navbar (flag + code)
<LanguageSwitcher variant="iconOnly" />   // Minimal (onboarding)
<LanguageSwitcher variant="toggle" />     // Quick switch (auth pages)
```

### 7.6 RTL Layout Rules

- Use logical properties: `start`, `end`, `ms`, `me`, `ps`, `pe`
- Use Tailwind utilities: `start-0`, `end-4`, `text-start`, `text-end`
- Icons: Apply `rtl-mirror` class for directional icons
- Toast position: Conditional based on `isRTL`
- Document attributes: `dir="rtl"` and `lang="ar"` automatically set

---

## 8. Error Handling

### 8.1 Backend Error Format (RFC 7807 ProblemDetails)

```json
{
  "type": "https://example.com/errors/validation-failed",
  "title": "Validation Failed",
  "status": 422,
  "detail": "Invalid email format",
  "instance": "/v1/auth/login",
  "errors": {
    "email": "errors.validation.invalid_email"
  }
}
```

### 8.2 Error Parsing

```typescript
import { parseProblemDetails, parseValidationErrors } from "@/lib/errorParser";

// Parse error
const errorResult = await parseProblemDetails(error);
// Returns: { key: "errors.http.401", params?: {...}, fallback?: "..." }

// Extract field errors
const fieldErrors = await parseValidationErrors(error);
// Returns: { email: "errors.validation.invalid_email" } or null
```

### 8.3 Error Translation

```typescript
import { translateErrorAsync } from "@/lib/translateError";
import { useTranslation } from "react-i18next";

const { t } = useTranslation();

try {
  await authApi.login(credentials);
} catch (error) {
  // Translate error to user's language
  const message = await translateErrorAsync(error, t);
  toast.error(message); // Shows localized error message
}
```

### 8.4 Error Translation Keys

Located in `locales/{lang}/errors.json`:

```json
{
  "errors": {
    "http": {
      "400": "Invalid request",
      "401": "Unauthorized",
      "404": "Not found",
      "422": "Validation error",
      "500": "Server error"
    },
    "network": {
      "timeout": "Request timeout",
      "connection": "Network error"
    },
    "validation": {
      "required": "This field is required",
      "invalid_email": "Invalid email format"
    },
    "auth": {
      "invalid_credentials": "Invalid email or password"
    }
  }
}
```

### 8.5 Form Validation Errors

```typescript
// In form submission
try {
  await api.updateProfile(data);
} catch (error) {
  const validationErrors = await parseValidationErrors(error);

  if (validationErrors) {
    // Set field-specific errors
    Object.entries(validationErrors).forEach(([field, key]) => {
      setError(field, { message: t(key, { defaultValue: key }) });
    });
  } else {
    // Generic error
    const message = await translateErrorAsync(error, t);
    toast.error(message);
  }
}
```

---

## 9. API Layer

### 9.1 Client Configuration

- Base URL: `VITE_API_BASE_URL` environment variable
- Timeout: 30 seconds
- Retry: 2 attempts, exponential backoff (max 3s between retries)
- Retry on: 408, 413, 429, 500, 502, 503, 504
- Automatic Bearer token attachment
- Automatic 401 handling with refresh
- Request deduplication for identical concurrent requests

### 9.2 Making Requests

```typescript
// Option 1: Direct client
import apiClient from "@/api/client";
const data = await apiClient.get("v1/users").json();

// Option 2: Typed helpers (with deduplication)
import { get, post, put, patch, del } from "@/api/client";
const users = await get<User[]>("v1/users");
const newUser = await post<User>("v1/users", { json: userData });
```

### 9.3 Domain Services

Create typed API services for each domain:

```typescript
// src/api/users.ts
import { get, post, patch, del } from "./client";
import type { User, UpdateUserRequest } from "./types";

export const usersApi = {
  getCurrent: () => get<User>("v1/users/me"),
  update: (id: string, data: UpdateUserRequest) =>
    patch<User>(`v1/users/${id}`, { json: data }),
  delete: (id: string) => del<void>(`v1/users/${id}`),
};
```

### 9.4 Type Safety with Zod

```typescript
// src/api/types/user.ts
import { z } from "zod";

export const UserSchema = z.object({
  id: z.string().uuid(),
  email: z.string().email(),
  firstName: z.string(),
  lastName: z.string(),
  role: z.enum(["admin", "member"]),
  workspaceId: z.string().uuid(),
  createdAt: z.string().datetime(),
  updatedAt: z.string().datetime(),
});

export type User = z.infer<typeof UserSchema>;

// Validate responses
const response = await apiClient.get("v1/users/me").json();
const user = UserSchema.parse(response); // Throws if invalid
```

---

## 10. Component Patterns

### 10.1 Atomic Design Levels

**Atoms** (Basic building blocks):

- `Button`: Primary, secondary, ghost, outline, loading states
- `Input`: Text, email, password, with icons, error states
- `Badge`: Status, count, color variants
- `Avatar`: User initials, images, online status
- `Modal`: Bottom sheet on mobile, centered on desktop
- `Skeleton`: Pulse animations for loading states

**Molecules** (Combinations):

- `FormInput`: Input + label + error message
- `BottomSheet`: Responsive drawer (bottom on mobile, side on desktop)
- `LanguageSwitcher`: Language picker with multiple variants
- `ResendCountdownButton`: OTP resend with countdown
- `PhoneCodeSelect`: Country code picker
- `CountrySelect`: Country picker with flags

**Organisms** (Complex sections):

- `LoginForm`: Email/password form with validation
- `Header`: Business switcher, search, notifications, user menu
- `Sidebar`: Navigation menu with collapsible sections
- `FilterDrawer`: Filter panel with apply/reset actions
- `DataTable`: Sortable, paginated table with actions

**Templates** (Page layouts):

- `OnboardingLayout`: Minimal layout with progress indicator
- `DashboardLayout`: Sidebar + header + main content
- `AuthLayout`: Split screen with branding

### 10.2 Component Composition Example

```typescript
// Good: Composable atoms
<Button variant="primary" size="lg" disabled={isLoading}>
  {isLoading && <Spinner size="sm" />}
  {t("common.save")}
</Button>

// Better: Purpose-built molecule
<FormInput
  name="email"
  label={t("auth.email")}
  type="email"
  icon={<Mail />}
  error={errors.email?.message}
  {...register("email")}
/>

// Best: Context-aware organism
<LoginForm
  onSubmit={handleLogin}
  isLoading={isLoading}
  error={error}
/>
```

### 10.3 BottomSheet Usage

```typescript
import { BottomSheet } from "@/components/molecules/BottomSheet";

<BottomSheet
  isOpen={isOpen}
  onClose={() => setIsOpen(false)}
  title={t("filters.title")}
  size="md"
  side="end"
  footer={
    <div className="flex gap-2">
      <button onClick={handleReset} className="btn btn-ghost flex-1">
        {t("common.reset")}
      </button>
      <button onClick={handleApply} className="btn btn-primary flex-1">
        {t("common.apply")}
      </button>
    </div>
  }
>
  {/* Filter content */}
</BottomSheet>;
```

### 10.4 Modal Usage

```typescript
import { Modal } from "@/components/atoms/Modal";

<Modal
  isOpen={isConfirmOpen}
  onClose={() => setIsConfirmOpen(false)}
  title={t("common.confirm")}
  size="sm"
  closeOnBackdropClick={false}
  footer={
    <>
      <button onClick={() => setIsConfirmOpen(false)} className="btn btn-ghost">
        {t("common.cancel")}
      </button>
      <button onClick={handleDelete} className="btn btn-error">
        {t("common.delete")}
      </button>
    </>
  }
>
  <p>{t("messages.confirm_delete")}</p>
</Modal>;
```

---

## 11. Forms & Validation

### 11.1 React Hook Form + Zod

```typescript
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

const schema = z.object({
  email: z.string().email("errors.validation.invalid_email"),
  password: z.string().min(8, "errors.validation.password_too_short"),
});

type FormData = z.infer<typeof schema>;

const {
  register,
  handleSubmit,
  formState: { errors, isSubmitting },
  setError,
} = useForm<FormData>({
  resolver: zodResolver(schema),
});

const onSubmit = async (data: FormData) => {
  try {
    await authApi.login(data);
  } catch (error) {
    const message = await translateErrorAsync(error, t);
    toast.error(message);
  }
};
```

### 11.2 Form Input Pattern

```typescript
<FormInput
  label={t("auth.email")}
  type="email"
  icon={<Mail />}
  error={t(errors.email?.message as string)} // Translate error key
  disabled={isSubmitting}
  {...register("email")}
/>
```

### 11.3 Submit Button Pattern

```typescript
<Button
  type="submit"
  variant="primary"
  disabled={isSubmitting}
  className="w-full"
>
  {isSubmitting ? (
    <>
      <Spinner size="sm" />
      {t("auth.logging_in")}
    </>
  ) : (
    t("auth.login")
  )}
</Button>
```

---

## 12. Best Practices

### 12.1 Code Quality

✅ **DO**:

- Use Zod for all validation (forms, API responses)
- Use `useTranslation()` for all user-facing text
- Use `useLanguage()` for RTL/LTR detection
- Use `translateErrorAsync()` for all error messages
- Use logical properties (`start`, `end`) for layouts
- Use DaisyUI tokens for colors (`bg-primary`, `text-base-content`)
- Use atomic components for consistency
- Use `BottomSheet` for mobile-first drawers
- Use `Modal` for confirmations and simple forms
- Write complete, production-ready code (no TODOs)

❌ **DON'T**:

- Use `any` type
- Use hardcoded strings (use i18n)
- Use physical properties (`left`, `right`, `margin-left`)
- Use Tailwind color classes directly (use DaisyUI tokens)
- Use localStorage for tokens
- Use inline styles (use Tailwind classes)
- Leave incomplete features (implement fully)
- Write verbose comments (code should be self-documenting)

### 12.2 Performance

- Lazy load routes: `const Dashboard = lazy(() => import("./routes/dashboard"));`
- Debounce search inputs: `useDebouncedValue(searchTerm, 300)`
- Use skeleton screens instead of spinners
- Optimize images: WebP format, lazy loading
- Code split large components
- Memoize expensive calculations: `useMemo`, `useCallback`

### 12.3 Accessibility

- Use semantic HTML: `<nav>`, `<main>`, `<article>`, `<section>`
- Add ARIA labels: `aria-label`, `aria-describedby`, `aria-live`
- Ensure keyboard navigation: `tabIndex`, `onKeyDown`
- Focus management: Trap focus in modals/drawers
- Color contrast: Minimum WCAG AA (4.5:1 for normal text)
- Screen reader testing: VoiceOver (macOS/iOS), NVDA (Windows)

### 12.4 Security

- Never expose tokens in localStorage
- Use `Secure` flag for cookies in production
- Validate all user inputs with Zod
- Sanitize data before rendering (use React's built-in escaping)
- Use CSP headers (backend responsibility)
- Implement rate limiting (backend responsibility)
- Never log sensitive data (passwords, tokens, PII)

---

## 13. Testing Strategy

### 13.1 Manual Testing Checklist

- [ ] Login/logout flow
- [ ] Language switching (Arabic ↔ English)
- [ ] RTL layout correctness
- [ ] Mobile responsiveness (375px, 768px, 1440px)
- [ ] Toast notifications (position, duration, RTL)
- [ ] Form validation (inline errors, submission)
- [ ] Error handling (network errors, 401, 404, 500)
- [ ] Session restoration (refresh page while logged in)
- [ ] Onboarding flow (all steps, payment, completion)
- [ ] Business switching
- [ ] Keyboard navigation
- [ ] Screen reader compatibility

### 13.2 Browser Compatibility

✅ Supported:

- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+
- Mobile browsers (iOS Safari, Chrome Mobile)

### 13.3 Automated Testing (Future)

- Unit tests: Vitest + React Testing Library
- Integration tests: Playwright
- E2E tests: Playwright with backend running
- Visual regression: Percy or Chromatic

---

## 14. Deployment & Environment

### 14.1 Environment Variables

```bash
# .env.production
VITE_API_BASE_URL=https://api.kyora.app
```

### 14.2 Build

```bash
npm run build      # Production build
npm run preview    # Preview production build locally
```

### 14.3 Deployment Checklist

- [ ] Update API base URL in `.env.production`
- [ ] Test production build locally with `npm run preview`
- [ ] Verify all translations are complete (Arabic + English)
- [ ] Test authentication flow (login, refresh, logout)
- [ ] Test onboarding flow (all steps, Stripe integration)
- [ ] Verify Stripe public key is correct for production
- [ ] Test mobile responsiveness
- [ ] Test RTL layout in Arabic
- [ ] Verify error handling and user-friendly messages
- [ ] Check browser console for errors/warnings
- [ ] Test with screen reader
- [ ] Verify analytics/monitoring integration (if any)

---

## 15. Common Patterns & Examples

### 15.1 Protected Route

```typescript
// src/routes/dashboard.tsx
import { RequireAuth } from "@/components/routing/RequireAuth";

<Route
  element={
    <RequireAuth>
      <DashboardLayout />
    </RequireAuth>
  }
>
  <Route path="/dashboard" element={<Dashboard />} />
  <Route path="/orders" element={<Orders />} />
</Route>;
```

### 15.2 Data Fetching with React Router v7

```typescript
// src/routes/orders/index.tsx
import type { Route } from "./+types/index";

export async function loader({ request, params }: Route.LoaderArgs) {
  const orders = await ordersApi.list({ businessId: params.businessId });
  return { orders };
}

export default function Orders({ loaderData }: Route.ComponentProps) {
  const { orders } = loaderData;

  return (
    <div>
      {orders.map((order) => (
        <OrderCard key={order.id} order={order} />
      ))}
    </div>
  );
}
```

### 15.3 Form Submission with Action

```typescript
export async function action({ request }: Route.ActionArgs) {
  const formData = await request.formData();
  const email = formData.get("email") as string;
  const password = formData.get("password") as string;

  try {
    await authApi.login({ email, password });
    return redirect("/dashboard");
  } catch (error) {
    return { error: await parseProblemDetails(error) };
  }
}
```

### 15.4 RTL-Aware Component

```typescript
import { useLanguage } from "@/hooks/useLanguage";

function MyComponent() {
  const { isRTL } = useLanguage();

  return (
    <div className={`flex ${isRTL ? "flex-row-reverse" : "flex-row"}`}>
      <button className="ms-auto">{/* Always at the end */}</button>
    </div>
  );
}
```

### 15.5 Language Switcher Placement

```typescript
// Login page (auth layout)
<LanguageSwitcher variant="toggle" showLabel />

// Dashboard header (navbar)
<LanguageSwitcher variant="compact" />

// Onboarding layout (minimal)
<LanguageSwitcher variant="iconOnly" />

// User menu dropdown (mobile)
<LanguageSwitcher variant="toggle" showLabel={false} />
```

---

## 16. Troubleshooting

### 16.1 Common Issues

**Issue**: User logged out on page refresh

- **Expected**: Access token is in memory (security)
- **Solution**: Automatic session restoration via refresh token cookie

**Issue**: Infinite login redirect loop

- **Cause**: `/login` route wrapped in `RequireAuth`
- **Solution**: Remove `RequireAuth` from auth routes

**Issue**: Toast not showing in correct position

- **Cause**: Not using `useLanguage()` for RTL detection
- **Solution**: `<Toaster position={isRTL ? "top-right" : "top-left"} />`

**Issue**: Layout not mirroring in Arabic

- **Cause**: Using physical properties (`left`, `right`)
- **Solution**: Use logical properties (`start`, `end`, `ms`, `me`)

**Issue**: Translations not updating

- **Cause**: Translation keys changed but not in both languages
- **Solution**: Update both `ar/*.json` and `en/*.json`

**Issue**: API requests not authenticated

- **Cause**: Not using `apiClient`
- **Solution**: Always use `apiClient` or typed helpers (`get`, `post`, etc.)

---

## 17. Quick Reference

### 17.1 Key Imports

```typescript
// Authentication
import { useAuth } from "@/hooks/useAuth";
import { RequireAuth } from "@/components/routing/RequireAuth";

// API
import apiClient, { setTokens, clearTokens } from "@/api/client";
import { get, post, put, patch, del } from "@/api/client";

// i18n
import { useTranslation } from "react-i18next";
import { useLanguage } from "@/hooks/useLanguage";

// Error handling
import { translateErrorAsync } from "@/lib/translateError";
import { parseProblemDetails, parseValidationErrors } from "@/lib/errorParser";

// Components
import { Button } from "@/components/atoms";
import { FormInput } from "@/components/molecules";
import { BottomSheet } from "@/components/molecules/BottomSheet";
import { Modal } from "@/components/atoms/Modal";
import { LanguageSwitcher } from "@/components/molecules/LanguageSwitcher";
```

### 17.2 Essential Hooks

```typescript
const { user, isAuthenticated, login, logout } = useAuth();
const { t } = useTranslation();
const { language, isRTL, changeLanguage } = useLanguage();
const isMobile = useMediaQuery("(max-width: 768px)");
```

### 17.3 Common Patterns

```typescript
// Login
await login({ email, password });
navigate("/dashboard");

// Error handling
try {
  await api.doSomething();
} catch (error) {
  const message = await translateErrorAsync(error, t);
  toast.error(message);
}

// Form validation
const schema = z.object({ email: z.string().email() });
const { register, handleSubmit } = useForm({ resolver: zodResolver(schema) });

// RTL-aware positioning
<Toast position={isRTL ? "top-right" : "top-left"} />;
```

---

## 18. Future Considerations

- **TanStack Query**: For server state caching and optimistic updates
- **React Router Deferred Data**: For streaming large datasets
- **Service Worker**: For offline support and PWA capabilities
- **Web Push Notifications**: For real-time order updates
- **Biometric Authentication**: Fingerprint/Face ID for mobile
- **Dark Mode**: User-configurable theme switching
- **Advanced Analytics**: User behavior tracking and conversion funnels

---

**End of Master Instructions**

_This document is the single source of truth for the Kyora Portal Web App. All implementation details, patterns, and conventions are consolidated here. Refer to this file for all development decisions._
