# Kyora Portal Web App

The business management portal for Kyora - built for social media entrepreneurs.

## Tech Stack

- **Framework**: React 19 with TypeScript
- **Build Tool**: Vite 7
- **Routing**: React Router v7 (Framework mode)
- **UI Library**: daisyUI 5 + Tailwind CSS
- **Forms**: React Hook Form + Zod validation
- **i18n**: i18next + react-i18next
- **HTTP Client**: ky
- **Icons**: lucide-react
- **Date**: date-fns

## Project Structure

\`\`\`
src/
├── api/              # API client and endpoints
├── assets/           # Static assets
├── components/       # Atomic Design components
│   ├── atoms/       # Basic building blocks
│   ├── molecules/   # Component combinations
│   ├── organisms/   # Complex sections
│   └── templates/   # Page layouts
├── hooks/           # Custom React hooks
├── i18n/            # Translations (ar, en)
├── lib/             # Utilities and helpers
├── routes/          # React Router routes
│   ├── _auth/      # Auth pages
│   ├── _app/       # Main app pages
│   └── onboarding/ # Onboarding flow
├── stores/          # State management (Zustand)
└── types/           # TypeScript definitions
\`\`\`

## Getting Started

### Prerequisites

- Node.js 20.19+ or 22.12+
- npm or pnpm

### Installation

\`\`\`bash
npm install
\`\`\`

### Environment Variables

Copy \`.env.example\` to \`.env\` and configure:

\`\`\`bash
cp .env.example .env
\`\`\`

### Development

\`\`\`bash
npm run dev
\`\`\`

The app will be available at \`http://localhost:3000\`

### Build

\`\`\`bash
npm run build
\`\`\`

### Preview Production Build

\`\`\`bash
npm run preview
\`\`\`

## Code Quality

### Linting

\`\`\`bash
npm run lint
\`\`\`

### Type Checking

\`\`\`bash
npm run type-check
\`\`\`

## Design Principles

- **RTL-First**: Arabic is the primary language. Always use logical CSS properties.
- **Mobile-First**: Every design decision prioritizes mobile usability.
- **Strict TypeScript**: No \`any\` types, strict mode enabled.
- **No Technical Debt**: No TODO or FIXME comments allowed.

## Key Features

- Multi-language support (Arabic, English)
- RTL/LTR layout switching
- Strict TypeScript configuration
- Path aliases (@/components, @/api, etc.)
- ESLint with strict rules
- daisyUI component library
- API client with interceptors

## Contributing

Follow the patterns established in \`.github/instructions/portal_web.instructions.md\`
