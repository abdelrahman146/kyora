Plan: Strict 1:1 Portal-web → Portal-v2 Migration (No Legacy URLs)
We’ll deliver a true 1:1 UX for all routes/components/files by overwriting portal-v2 with portal-web implementations, adapting only for TanStack Router/Form/Store, and keeping portal-v2’s file-based route structure + TanStack Store. We will NOT add any legacy/backward-compatible routes or redirects; portal-v2 URLs are the source of truth, and all navigation/UI will be updated to match portal-web UX within that structure.

Steps

1. Lock i18n conventions across portal-v2 (eliminate raw keys everywhere)

Align init.ts with portal-web key style: nested keys in translation (e.g. t('auth.welcome_back')).
Replace all invalid namespace usage (t('auth:...')) in auth/onboarding/business routes under routes.
Overwrite LanguageSwitcher to portal-web behavior/variants: LanguageSwitcher.tsx → LanguageSwitcher.tsx.

2. Enforce a single URL system: portal-v2 route structure only (no legacy, no redirects)

Remove any “old path” references in UI and navigation; update all internal links to portal-v2 file routes.
Ensure layouts/nav components (Header/Sidebar/BottomNav) generate portal-v2 URLs consistently.

3. Auth routes 1:1 UX migration (within portal-v2 URLs)

Overwrite each auth screen to match portal-web UX states and visuals, adapting only navigation + form stack:
login.tsx → login.tsx
forgot-password.tsx → forgot-password.tsx
reset-password.tsx → reset-password.tsx
oauth-callback.tsx → oauth-callback.tsx
Ensure LanguageSwitcher is present on login and any auth pages where portal-web shows it.
Overwrite shared layouts + navigation components to portal-web parity (TanStack Router-native)

4. Overwrite these portal-v2 components with portal-web implementations and adapt only navigation primitives:
   DashboardLayout.tsx → DashboardLayout.tsx
   Header.tsx → Header.tsx
   Sidebar.tsx → Sidebar.tsx
   BottomNav.tsx → BottomNav.tsx

5. Overwrite feature routes to portal-web UX (within portal-v2 URLs), then reconcile TanStack Router search schemas

Home: home.tsx → index.tsx
Dashboard: dashboard.tsx → index.tsx
Customers list/detail: overwrite to portal-web UX while keeping portal-v2 paths:
customers.tsx → index.tsx
customers.$customerId.tsx → $customerId.tsx
Final verification pass: every TS/TSX in portal-v2 is 1:1 UX-complete and production-grade

6. Verify no hardcoded text, every user string uses t(), all portal-web UX states exist, RTL matches, shared components are reused, and no simplified/stub implementations remain under src.
   Further Considerations
   Confirm portal-v2 URLs are canonical; remove any legacy references in UI/docs.
   Standardize “TanStack Router navigation pattern” to avoid Link typing pitfalls (consistent wrapper or helper).
   Ensure LanguageSwitcher parity everywhere portal-web shows it (auth + header + onboarding layouts)

7. Final verification pass — Review every TSX file in portal-v2 to confirm: no hardcoded text (all use t()), all features from portal-web exist, all props/APIs match, RTL support intact, mobile-first responsive design preserved.
