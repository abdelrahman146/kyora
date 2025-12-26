---
description: Kyora Portal Web App - AI Agent Instructions
applyTo: "portal-web/**"
---

Kyora Portal Web App - AI Agent Instructions

This document serves as the Master Instruction File for the AI Agent responsible for creating and maintaining the Kyora Portal Web App. This application is the client-side command center for business owners using the Kyora platform.

1. Project Philosophy & Core Requirements

1.1. Target Audience

Users: Busy social media business owners and entrepreneurs.

Context: Users are often on the go, switching between mobile and desktop devices. They need quick access to data, fast interactions, and zero friction.

Value Proposition: Kyora provides deep analysis and management tools. The Portal must reflect reliability, professionalism, and speed.

1.2. Design Philosophy

Mobile-First & Thumb-Friendly: Every design decision must prioritize mobile usability. Interactive elements must be easily reachable.

Forms: Prefer Bottom Sheet Modals (Drawers) over centered modals on mobile.

Navigation: Accessible, thumb-friendly menus.

Touch Targets: Minimum 44px for all tappable elements.

Arabic First (RTL First):

Arabic is the primary language. The application must be designed natively for RTL (Right-To-Left).

English (LTR) is a second-class citizen but must be fully supported.

Layouts: Use logical CSS properties (e.g., ms- (margin-start), pe- (padding-end), start-, end-) instead of physical properties (left, right) to ensure automatic mirroring.

Zero-State & Feedback:

Loading: Use skeletons (pulse animations) matching the content shape, not generic spinners.

Feedback: Every action (Create, Update, Delete) must have a visual confirmation (Toast notification).

Errors: Fail gracefully. Use Error Boundaries for sections. Never crash the entire app. Provide "Retry" buttons.

1.3. Production Standards

No Technical Debt: Do not write // TODO or // FIXME. Implement complete features.

Strict Types: TypeScript strict: true. No any types. Zod schemas for all data.

Atomic Design: Follow the folder structure and composition patterns used in storefront-web.

Clean Architecture: Separation of concerns (UI, State, API).

2. Technology Stack

Framework: React Router v7 (using Framework mode with loader and action patterns).

Language: TypeScript.

Build Tool: Vite.

UI Library: DaisyUI (based on Tailwind CSS).

Styling: Tailwind CSS.

Form Handling: React Hook Form.

Validation: Zod.

Internationalization: i18next with react-i18next.

Icons: lucide-react (or standard icon set defined in branding).

Date Handling: date-fns (ensure locale support for Arabic).

HTTP Client: Native fetch wrapper or ky with interceptors for Auth.

3. Architecture & Folder Structure

Adhere to the Atomic Design principle. The structure should mirror the storefront-web approach but adapted for a dashboard/portal context.

src/
├── api/ # API client, types, and endpoints (generated or manually typed from swagger.json)
│ ├── auth.ts
│ ├── business.ts
│ └── ...
├── assets/ # Static assets (images, fonts)
├── components/ # Atomic Design Components
│ ├── atoms/ # Basic building blocks (Button, Input, Badge, Avatar, Skeleton)
│ ├── molecules/ # Combinations (FormInput, SearchBar, Toast, ModalHeader)
│ ├── organisms/ # Complex sections (LoginForm, DataTable, Sidebar, BottomSheet)
│ └── templates/ # Page layouts (DashboardLayout, AuthLayout, SettingsLayout)
├── hooks/ # Custom React hooks (useAuth, useToast, useMediaQuery)
├── i18n/ # Translation files (locales/ar, locales/en)
├── lib/ # Utilities (validation, formatting, constants)
├── routes/ # React Router v7 Route definitions (File-system routing preferred)
│ ├── \_auth/ # Auth layout routes (Login, Register)
│ ├── \_app/ # Main app layout (Sidebar, Header)
│ │ ├── dashboard/
│ │ ├── inventory/
│ │ ├── orders/
│ │ └── settings/
│ └── onboarding/ # Onboarding specific routes
├── stores/ # Client-side state (Zustand) - strictly for UI state (e.g., toggle sidebar)
└── types/ # Global TypeScript definitions

4. UI/UX & Design System

4.1. Theming (DaisyUI)

Colors: Refer strictly to .github/instructions/branding.instructions.md for primary, secondary, and accent colors.

Dark/Light Mode: Support both if required by branding, otherwise stick to the defined theme.

Consistency: Use Tailwind utility classes for margins and padding to ensure consistent spacing (e.g., p-4, gap-2).

4.2. Layouts

App Shell:

Desktop: Left Sidebar (collapsible), Top Header (Search, Business Switcher, Profile).

Mobile: Bottom Navigation Bar (Top 4-5 items), Hamburger menu for "More". Top Header for context switching.

Business Switcher:

Located in the Header.

Dropdown allowing users to switch between businesses in their workspace.

Changing business updates the global context/URL and re-fetches data.

4.3. Interactions

Bottom Sheets (Drawers):

Use for all "Create", "Edit", and "Filter" actions on Mobile.

On Desktop, these can adapt to Side Drawers (Slide-overs) or Center Modals depending on content density.

Toasts:

Position: Bottom-center (Mobile), Top-start (Desktop RTL) / Top-end (Desktop LTR).

Types: Success (Green), Error (Red), Info (Blue), Warning (Yellow).

5. Feature Implementation Guidelines

5.1. Authentication (Security & Flow)

Strategy: JWT (Bearer Token) + Refresh Token.

Storage:

access_token: In-memory (or short-lived storage if needed, but prefer strict security).

refresh_token: HttpOnly Cookie (if backend supports) or secure Cookie set by client.

Login Flow:

User enters credentials.

POST /v1/auth/login.

Save tokens.

Redirect to Dashboard or Onboarding (based on user state).

Interceptor Logic:

On 401 Unauthorized response:

Pause outgoing requests.

Call /v1/auth/refresh using the stored cookie.

If successful: Update access_token, retry original failed requests.

If failed: Clear session, Redirect to Login.

Google OAuth:

Implement using the backend's OAuth callback flow.

Button: "Continue with Google".

5.2. Onboarding Journey

Critical Requirement: Smooth, linear, and encouraging.

Backend Integration: Read swagger.json -> Onboarding endpoints.

Flow:

Welcome: Introduction to Kyora.

Workspace Setup: Naming the workspace.

Business Setup: Creating the first business entity.

Plan Selection: Displaying available subscription tiers.

Completion: Success animation and redirection to Dashboard.

State: Persist onboarding step in URL or backend state to allow resuming.

5.3. Workspace & Team Management (RBAC)

Role-Based Access Control:

Fetch user permissions from the backend on load.

UI Guarding: Hide buttons/links (e.g., "Billing", "Invite Member") if the user lacks the required permission.

Route Guarding: Redirect users from forbidden routes.

Business Switching:

Users may belong to multiple businesses within a workspace.

Switching business acts as a global filter for all subsequent data fetches.

5.4. Resource Modules (Inventory, Orders, Customers)

Standard Layout:

Header: Title + "Add New" Button.

Toolbar: Search Input (Debounced) + Filter Button (Opens Sheet) + Sort Dropdown.

Data View:

Mobile: Card List View (Info stacked vertically).

Desktop: Table View (Sortable columns).

Pagination: Infinite Scroll for mobile lists, standard pagination for desktop tables.

Forms:

Use react-hook-form + zod resolver.

Validate inputs locally before submission.

Show inline error messages (translated).

Disable submit button while isSubmitting.

5.5. Subscription & Billing (Stripe)

Integration: Use Stripe Elements or redirect to Stripe Checkout based on backend implementation.

Features:

Plan Upgrade/Downgrade: Visual comparison of features.

Invoices: List view with download PDF capability.

Payment Methods: Add/Remove cards.

Error Handling: Specific handling for payment failures (SCA requirements, declined cards).

6. Coding Conventions

6.1. React Router v7 Patterns

Use Loaders for data fetching:

export async function loader({ request, params }: Route.LoaderArgs) {
const data = await api.getOrders(params.businessId);
return data;
}

Use Actions for mutations:

export async function action({ request }: Route.ActionArgs) {
const formData = await request.formData();
// validation and api call
return redirect('/orders');
}

Use useLoaderData to access data in components.

Use useSubmit or <Form> for interactions.

6.2. Internationalization (i18next)

Keys: Use nested keys for organization (e.g., auth.login.title, common.save).

Interpolation: Use {{value}} for dynamic data.

Hook: const { t, i18n } = useTranslation();

Direction: document.dir must update when language changes.

6.3. API & Data

Centralized Client: Create a singleton API client instance.

Types: All API responses must be typed via Zod schemas or TypeScript interfaces derived from Swagger.

Error Parsing: Helper function to parse backend error responses (Standard ProblemDetails format) into user-friendly messages.

6.4. Branding

Typography: Use the font family specified in branding.instructions.md.

Radius: Consistent border-radius (e.g., rounded-box from DaisyUI).

Shadows: Soft shadows for depth.

7. Implementation Checklist for Agent

Analyze Swagger: Map every screen to a specific API endpoint.

Setup Routes: Define the route hierarchy in React Router v7.

Build Atoms: Create base components (Button, Input) with DaisyUI/Tailwind.

Implement Auth: Login/Register/Refresh logic.

Build Shell: Sidebar, Header, Business Switcher.

Implement CRUD: Build one module (e.g., Customers) fully to establish the pattern, then replicate for Inventory/Orders.

Polish UX: Add loading skeletons, toasts, and verify RTL layout.

Final Review: Check against copilot-instructions.md constraints.
