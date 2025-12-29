Kyora — AI coding agent quickstart

## What is Kyora?

Kyora is a business management assistant specifically designed for solo social media entrepreneurs and small teams who sell products on social media platforms. It's built for business owners who are excellent at creating and selling their products but may not have expertise in accounting or business management. Kyora handles all the complex financial tracking and business operations behind the scenes, presenting everything in a simple, easy-to-understand way so owners can focus on what they do best: creating and selling.

### The Problem Kyora Solves

Social media business owners face a unique challenge: they're passionate about their products and skilled at selling on platforms like Instagram, Facebook, TikTok, and WhatsApp, but they struggle with:

- **Financial Confusion**: Not knowing if they're actually making profit or just generating revenue
- **Order Chaos**: Tracking orders scattered across DMs, comments, and multiple platforms
- **Inventory Blindness**: Not knowing what's in stock, what's selling, or when to reorder
- **Cash Flow Mystery**: Unclear about their actual financial position and available cash
- **Tax Anxiety**: Worried about missing important financial records for tax time
- **Growth Paralysis**: Unable to make data-driven decisions because they lack proper business insights

Kyora becomes their silent business partner that handles all this complexity automatically, without requiring them to become accounting experts or spend hours on administrative tasks.

### Core Value Proposition

**Simplicity First**: Every feature is designed to be intuitive and require zero accounting knowledge. Complex financial concepts are translated into simple, actionable insights.

**Automatic Heavy Lifting**: Kyora automatically tracks revenue recognition, calculates profitability, monitors inventory levels, and maintains complete financial records without manual bookkeeping.

**Social Media Native**: Understanding that orders come through DMs, comments, and chat apps, Kyora makes it easy to quickly log orders from any source and keep everything organized in one place.

**Peace of Mind**: Owners can sleep well knowing their business finances are properly tracked, they have records for tax purposes, and they truly understand their financial position.

### Key Features

- **Simple Order Management**: Quickly add orders from any social media platform, track their status from payment to delivery, and automatically recognize revenue.
- **Effortless Customer Tracking**: Automatically build a customer database from orders, see purchase history, and identify your best customers without manual data entry.
- **Intelligent Inventory**: Know what's in stock, get alerts when items are running low, and understand which products are your best sellers.
- **Clear Financial Picture**: See your actual profit (not just revenue), understand your cash flow, and get simple reports that show how your business is really doing.
- **Automated Accounting**: Revenue recognition, expense tracking, and financial reporting happen automatically in the background.
- **Business Insights in Plain English**: No confusing charts or accounting jargon—just clear answers like "You made $X profit this month" and "Your best-selling product is Y."

### Multi-Tenancy Model

Kyora uses a workspace-based multi-tenancy architecture where:

- Each **workspace** represents one business owner or small team with their own subscription, users, and data.
- **Users** belong to a single workspace and have role-based permissions (admin/member for team collaboration).
- All data is strictly isolated by workspace to ensure complete privacy and security.
- Billing and subscriptions are managed at the workspace level, with affordable pricing tiers designed for solo entrepreneurs and small teams.

### Target Users

- **Solo social media sellers**: Individuals selling handmade goods, fashion, beauty products, food, or services through Instagram, Facebook, TikTok, WhatsApp, etc.
- **Small social commerce teams**: 2-5 person teams managing a social media-based business together
- **Side hustlers**: People running a business alongside their day job who need dead-simple management tools
- **Product creators**: Artisans, designers, bakers, makers who want to focus on creation, not administration
- **Non-technical entrepreneurs**: Business owners who are intimidated by complex software and just want something that works

### Key Business Flows

1. **Quick Onboarding**: Sign up with email or Google, verify email, name your workspace, and start adding orders immediately—no complex setup required.
2. **Simple Order Entry**: Add an order in seconds (customer name, product, price, payment received) → Kyora automatically tracks it through delivery → automatically recognizes revenue and updates inventory.
3. **Automatic Financial Tracking**: As orders are added → revenue is recognized → inventory is adjusted → profit is calculated → financial position is updated in real-time.
4. **Instant Business Insights**: Open the dashboard → see profit this month, total revenue, best customers, top products—all in plain language with simple visuals.
5. **Team Collaboration (Optional)**: Invite a helper or partner → assign them as admin or member → they can help manage orders while you focus on production.
6. **Subscription Management**: Start with a free or basic plan → as business grows, upgrade to handle more orders → billing happens automatically through Stripe.

Big picture

## Monorepo Structure

Kyora is organized as a **monorepo** to support multiple projects (backend, frontend, mobile, etc.):

```
kyora/
├── backend/              # Go backend API server (current focus)
│   ├── cmd/             # CLI commands
│   ├── internal/        # Internal packages
│   ├── main.go          # Entry point
│   ├── go.mod           # Go dependencies
│   ├── .air.toml        # Hot reload config
│   └── .kyora.yaml      # Backend configuration
├── storefront-web/       # Customer-facing storefront (React)
│   ├── src/             # App code
│   ├── public/
│   ├── package.json
│   └── DESIGN_SYSTEM.md
├── Makefile             # Root-level build commands
├── README.md            # Monorepo overview
└── STRUCTURE.md         # Monorepo guidelines
```

**Important**:

- All Go backend code is in the `backend/` directory. When referencing files or paths, always include the `backend/` prefix.
- All storefront web code is in the `storefront-web/` directory. When referencing files or paths, always include the `storefront-web/` prefix.
- All portal web code is in the `portal-web/` directory. When referencing files or paths, always include the `portal-web/` prefix.

**Project-Specific Reference Instructions**

Before working on any task, always refer to the relevant instruction files based on the project you're working on:

- **When working on `backend/` (Go API)**:

  - Always refer to `.github/instructions/backend.instructions.md` for architecture, patterns, and conventions
  - Refer to `.github/instructions/resend.instructions.md` when working with email functionality
  - Refer to `.github/instructions/stripe.instructions.md` when working with billing/payments

- **When working on `portal-web/` (React business dashboard)**:

  - Always refer to `.github/instructions/portal_web.instructions.md` for architecture, patterns, and conventions
  - Refer to `.github/instructions/ky.instructions.md` when making HTTP requests
  - Refer to `.github/instructions/react-router.instructions.md` when working with routing
  - Refer to `.github/instructions/branding.instructions.md` for design system and styling guidelines
  - Refer to `.github/instructions/daisyui.instructions.md` when working with UI components

- **When working on `storefront-web/` (Customer-facing storefront)**:
  - Refer to `storefront-web/DESIGN_SYSTEM.md` for component and styling guidelines
  - Refer to `.github/instructions/branding.instructions.md` for brand consistency
  - Refer to `.github/instructions/react-router.instructions.md` when working with routing
  - Refer to `.github/instructions/daisyui.instructions.md` when working with UI components

**General Instructions**

- the code is maintained by a single developer so we should always aim for simplicity and clarity in the code we write.
- the code is still under heavy development so we should always write code that is flexible and easy to change as requirements evolve and we can do breaking changes as needed.
- **never leave any TODOs or FIXMEs in the code**, we should always address them before finalizing the code.
- **never leave deprecated code in the codebase**, we should always remove it immediately to keep the code clean and maintainable. When creating a replacement for existing functionality, delete the old deprecated code completely - don't just mark it as deprecated.
- we should always follow the SOLID principles and best practices when writing code.
- we should always follow the existing code style and conventions used in the project to maintain consistency across the codebase.
- we should always write clear and concise comments and documentation for the code we write to ensure that it's easy to understand for other developers.
- we should always consider performance implications when writing code and optimize for efficiency where necessary.
- we should always consider security implications when writing code and ensure that the code is secure and follows best security practices.
- we should always consider scalability implications when writing code and ensure that the code can handle increased load and scale as needed.
- we should always consider maintainability implications when writing code and ensure that the code is easy to maintain and extend in the future
- the code should be clear and human maintainable and concise and follow best practices.
- when we have a generic sharable functionality we should add it in utils package either in helpers package or create its own package if its big enough.
- whenever you find inefficient or duplicate code across domains we should refactor it into a shared utility function in the utils package or helpers package.
- whenever you find a code that doesn't follow best practices or has potential bugs we should fix it to follow best practices and avoid potential issues in the future.
- the code output should be always secure, robust, and production-ready and follows best practices and standards and should be 100% complete.
- we should alwasys look for smart simple and elegant solutions to complex problems leaving the code clean and maintainble and very easy to fix and extend in the future.
- everytime we write code we should always think about the future and how this code will be used and maintained in the future and we should always write code that is easy to understand and maintain in the future.
- whenever we introduce new design pattern or archiectural changes we need to make sure it aligns with the overall architecture and design principles of the project and we should document the changes properly to explain the reasoning behind them.
- always do cleanups for no longer needed code or functions.
- never ever brief any implementation. always provide complete and thorough implementations.
- never ever settle on examples or partial implementations. always provide complete and thorough implementations.
- always aim for high-quality code that is secure, robust, maintainable, and production-ready

## Project-specific notes

### Backend (Go)

- Follow the rules in `.github/instructions/backend.instructions.md` for architecture, patterns, and multi-tenancy scoping.
- Prefer domain services for business logic and keep HTTP handlers thin.

### Storefront Web (React)

- Tech stack: React 19, React Router v7, Tailwind CSS v4 (CSS-first), daisyUI v5, TanStack Query, Zustand, i18next.
- Styling rules: follow `.github/instructions/branding.instructions.md`, `storefront-web/DESIGN_SYSTEM.md`, and the theme/tokens in `storefront-web/src/index.css`.
- RTL-first: never assume left/right; prefer logical properties and Tailwind `start-*` / `end-*` utilities.
