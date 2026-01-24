---
description: Portal Web development workflow - Vite config, dev server, build process, testing, deployment (portal-web only)
applyTo: "portal-web/**"
---

# Portal Web Development

Development workflow for portal-web.

**Cross-refs:**

- General architecture: `../../_general/architecture.instructions.md`
- General testing: `../../_general/testing.instructions.md`

---

## 1. Local Development

### Start Dev Server

```bash
# From project root
make dev.portal

# Or directly
cd portal-web
npm run dev

# Custom port
PORTAL_PORT=3001 make dev.portal
```

**Dev server features:**

- Hot Module Replacement (HMR)
- Fast Refresh (preserves React state)
- Auto HTTPS (self-signed cert)
- API proxy to backend

---

## 2. Vite Configuration

### vite.config.ts

```typescript
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { TanStackRouterVite } from "@tanstack/router-vite-plugin";
import path from "path";

export default defineConfig({
  plugins: [
    react(),
    TanStackRouterVite(), // Auto-generates routeTree.gen.ts
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    port: 3000,
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: "dist",
    sourcemap: true,
    rollupOptions: {
      output: {
        manualChunks: {
          "react-vendor": ["react", "react-dom"],
          tanstack: [
            "@tanstack/react-router",
            "@tanstack/react-query",
            "@tanstack/react-form",
          ],
          chart: ["chart.js", "react-chartjs-2"],
        },
      },
    },
  },
});
```

---

## 3. TanStack Router Route Generation

TanStack Router uses **file-based routing**. Route tree is auto-generated.

### Route Generation

```bash
# Auto-generated on save (dev mode)
npm run dev

# Manual generation
npm run build:routes
```

**Output:** `src/routeTree.gen.ts`

**Important:**

- Never edit `routeTree.gen.ts` manually
- Commit it to version control
- Regenerate after adding/removing route files

---

## 4. Testing Workflow

### Run All Tests

```bash
# From project root
make portal.check

# Or directly
cd portal-web
npm run test
```

### Unit Tests (Vitest)

```bash
npm run test:unit

# Watch mode
npm run test:unit:watch

# Coverage
npm run test:coverage
```

### E2E Tests (Playwright)

```bash
npm run test:e2e

# Headed mode (see browser)
npm run test:e2e:headed

# Specific test file
npm run test:e2e -- tests/orders.spec.ts
```

---

## 5. Type Checking

### TypeScript Check

```bash
# From project root
make portal.check

# Or directly
npm run typecheck
```

**tsconfig.json:**

```json
{
  "compilerOptions": {
    "target": "ES2022",
    "lib": ["ES2022", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "moduleResolution": "bundler",
    "jsx": "react-jsx",
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "skipLibCheck": true,
    "esModuleInterop": true,
    "resolveJsonModule": true,
    "paths": {
      "@/*": ["./src/*"]
    }
  },
  "include": ["src"],
  "exclude": ["node_modules", "dist"]
}
```

---

## 6. Linting & Formatting

### ESLint

```bash
npm run lint

# Auto-fix
npm run lint:fix
```

**eslint.config.js:**

```javascript
import js from "@eslint/js";
import react from "eslint-plugin-react";
import reactHooks from "eslint-plugin-react-hooks";
import tseslint from "@typescript-eslint/eslint-plugin";

export default [
  js.configs.recommended,
  {
    files: ["**/*.{ts,tsx}"],
    plugins: {
      "@typescript-eslint": tseslint,
      react: react,
      "react-hooks": reactHooks,
    },
    rules: {
      "react-hooks/rules-of-hooks": "error",
      "react-hooks/exhaustive-deps": "warn",
      "@typescript-eslint/no-unused-vars": [
        "error",
        { argsIgnorePattern: "^_" },
      ],
    },
  },
];
```

### Prettier

```bash
npm run format

# Check only
npm run format:check
```

**prettier.config.js:**

```javascript
export default {
  semi: true,
  singleQuote: true,
  trailingComma: "es5",
  printWidth: 100,
  tabWidth: 2,
};
```

---

## 7. Build Process

### Production Build

```bash
# From project root
make portal.build

# Or directly
npm run build
```

**Output:** `dist/`

### Preview Production Build

```bash
make portal.preview

# Or directly
npm run preview
```

Serves production build locally at `http://localhost:4173`

---

## 8. Environment Variables

### .env Files

```
.env                  # Defaults (committed)
.env.local            # Local overrides (NOT committed)
.env.production       # Production (NOT committed)
```

### Example .env

```bash
VITE_API_BASE_URL=http://localhost:8080
VITE_APP_ENV=development
```

### Usage

```typescript
const apiUrl = import.meta.env.VITE_API_BASE_URL;
```

**Rules:**

- Prefix all vars with `VITE_`
- Never commit secrets (use `.env.local`)
- Public vars only (exposed to client)

---

## 9. Debugging

### React DevTools

Install browser extension:

- Chrome: https://chrome.google.com/webstore/detail/react-developer-tools
- Firefox: https://addons.mozilla.org/en-US/firefox/addon/react-devtools/

### TanStack Query DevTools

```tsx
// src/main.tsx
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";

<QueryClientProvider client={queryClient}>
  <App />
  <ReactQueryDevtools initialIsOpen={false} />
</QueryClientProvider>;
```

### TanStack Router DevTools

```tsx
// src/router.tsx
import { TanStackRouterDevtools } from "@tanstack/router-devtools";

<Router>
  <App />
  <TanStackRouterDevtools position="bottom-right" />
</Router>;
```

---

## 10. Common Tasks

### Add New Dependency

```bash
npm install <package>

# Development dependency
npm install -D <package>
```

**Remember:** New dependencies require PO gate (see AGENTS.md).

### Add New Route

1. Create route file:

   ```tsx
   // src/routes/business/$businessDescriptor/new-feature/index.tsx
   export const Route = createFileRoute(
     "/business/$businessDescriptor/new-feature/",
   )({
     component: NewFeaturePage,
   });
   ```

2. Route tree auto-regenerates (dev mode)

3. Add to sidebar (if needed):
   ```tsx
   // src/features/dashboard-layout/components/Sidebar.tsx
   <SidebarLink to="/business/$businessDescriptor/new-feature">
     New Feature
   </SidebarLink>
   ```

### Add New Translation

1. Add keys to namespace files:

   ```typescript
   // src/i18n/ar/new-feature.ts
   export default {
     title: "الميزة الجديدة",
     description: "الوصف",
   };
   ```

2. Import in i18n config:

   ```typescript
   // src/i18n/init.ts
   import newFeatureAr from "./ar/new-feature";
   import newFeatureEn from "./en/new-feature";
   ```

3. Use in component:
   ```tsx
   const { t } = useTranslation(["newFeature"]);
   t("newFeature:title");
   ```

---

## 11. Performance Optimization

### Code Splitting

Vite automatically splits code by routes. For manual splits:

```tsx
import { lazy } from "react";

const HeavyComponent = lazy(() => import("./HeavyComponent"));

<Suspense fallback={<Loading />}>
  <HeavyComponent />
</Suspense>;
```

### Bundle Analysis

```bash
npm run build -- --mode analyze
```

Uses `rollup-plugin-visualizer` to generate bundle report.

---

## 12. CI/CD

### GitHub Actions (Example)

```yaml
name: Portal CI

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: "20"
      - run: cd portal-web && npm ci
      - run: cd portal-web && npm run typecheck
      - run: cd portal-web && npm run lint
      - run: cd portal-web && npm run test
      - run: cd portal-web && npm run build
```

---

## Agent Validation

Before completing development task:

- ☑ Dev server runs (`make dev.portal`)
- ☑ Type check passes (`make portal.check`)
- ☑ Linting passes (`npm run lint`)
- ☑ Tests pass (`npm run test`)
- ☑ Production build succeeds (`make portal.build`)
- ☑ No console errors/warnings
- ☑ Route tree regenerated if routes changed
- ☑ New dependencies justified (PO gate)

---

## Resources

- Vite Docs: https://vitejs.dev/
- TanStack Router: https://tanstack.com/router
- Vitest: https://vitest.dev/
- Playwright: https://playwright.dev/
