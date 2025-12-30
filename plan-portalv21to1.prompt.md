Plan: 1:1 Portal-web to Portal-v2 Migration
Copy ALL content and functionality from portal-web to portal-v2 file-by-file, adapting only for TanStack Router/Form/Store while keeping portal-v2's route structure and TanStack Store.

Steps
Migrate missing API files — Copy api/address.ts, api/user.ts, api/types/business.ts, api/types/customer.ts to portal-v2, wrapping functions with TanStack Query hooks following existing patterns in customer.ts.

Migrate missing lib utilities — Copy lib/phone.ts, lib/sessionStorage.ts to portal-v2 as-is (no stack changes needed).

Migrate customer organisms — Copy organisms/customers/AddCustomerSheet.tsx, EditCustomerSheet.tsx, AddressSheet.tsx to portal-v2, refactoring React Hook Form → TanStack Form.

Overwrite components with portal-web versions — For EVERY existing component (Pagination, Header, Sidebar, BottomNav, DashboardLayout, Modal, Table, etc.), copy the full portal-web implementation and adapt only: useNavigate/<Link> → TanStack Router equivalents, form handling → TanStack Form, keeping all translations, features, and UI identical.

Overwrite route files with portal-web content — Copy exact content from home.tsx → index.tsx, dashboard.tsx → business/$businessDescriptor.tsx, customers.tsx → customers/index.tsx, etc., adapting route params from :param to $param syntax.

Verify and sync TanStack Stores — Compare AuthContext with authStore, OnboardingContext with onboardingStore, and businessStore with businessStore — ensure identical state shape and all methods exist, optimized for TanStack Store patterns.

Sync translation files 100% — Diff and merge ALL translation content from ar and en/ into portal-v2's i18n folders, ensuring every key exists and every component uses t() calls.

Migrate types and add barrel exports — Copy types/index.ts, create index.ts barrel exports in atoms/, molecules/, organisms/, templates/.

Delete divergent portal-v2 components — Remove CustomerForm.tsx (replaced by AddCustomerSheet/EditCustomerSheet), and any simplified/stub implementations that don't match portal-web quality.

Final verification pass — Review every TSX file in portal-v2 to confirm: no hardcoded text (all use t()), all features from portal-web exist, all props/APIs match, RTL support intact, mobile-first responsive design preserved.
