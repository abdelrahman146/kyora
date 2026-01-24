---
description: "Adding new apps/services to Kyora monorepo — scaffolding, Makefile integration, CI/CD, AGENTS.md"
applyTo: "**/AGENTS.md,Makefile,package.json,go.mod"
---

# Adding Projects to Kyora Monorepo — Single Source of Truth (SSOT)

**SSOT Hierarchy:**

- Parent: `.github/copilot-instructions.md`, `AGENTS.md`
- Related: `monorepo/structure.instructions.md`, `monorepo/workflows.instructions.md`

**When to Read:**

- Adding a new backend service
- Adding a new frontend app
- Creating shared libraries
- Integrating new apps into CI/CD

---

## 1) Scaffolding New Backend Apps

### Directory Structure

```
<app-name>/                 # e.g., admin-api/, webhook-worker/
├── AGENTS.md               # App-specific agent manifest
├── cmd/                    # CLI commands
│   ├── root.go
│   └── server.go
├── docs/                   # Swagger/OpenAPI (if HTTP service)
├── internal/
│   ├── domain/             # Business logic
│   ├── platform/           # Infrastructure
│   ├── server/             # HTTP server (if applicable)
│   └── tests/              # Tests
├── go.mod
├── go.sum
└── main.go
```

### Steps

1. **Create directory and initialize Go module:**

```bash
mkdir <app-name>
cd <app-name>
go mod init github.com/kyora/<app-name>
```

2. **Create AGENTS.md:**

```markdown
# <App Name>

For monorepo structure and SSOT references, see [root AGENTS.md](../AGENTS.md).

## Overview

[Brief description of what this app does]

## Tech Stack

- Go 1.22+
- [Other dependencies]

## Commands

See root Makefile for all commands. App-specific:

- `make dev.<app>` — Run app in dev mode
- `make test.<app>` — Run app tests
- `make build.<app>` — Build app binary

## SSOT References

- [backend/_general/architecture.instructions.md] — Backend patterns
- [backend/_general/go-patterns.instructions.md] — Go patterns
- [backend/_general/testing.instructions.md] — Testing guidelines
```

3. **Add Makefile targets** (in root `Makefile`):

```makefile
# <App Name> targets
.PHONY: dev.<app>
dev.<app>: ## Run <app> in dev mode
	@cd <app-name> && go run main.go

.PHONY: test.<app>
test.<app>: ## Run <app> tests
	@cd <app-name> && go test ./...

.PHONY: build.<app>
build.<app>: ## Build <app> binary
	@cd <app-name> && go build -o bin/<app> main.go
```

4. **Add docker-compose service** (if needed):

```yaml
# In docker-compose.dev.yml
services:
  <app-name>:
    build:
      context: ./<app-name>
      dockerfile: Dockerfile.dev
    ports:
      - "<port>:<port>"
    environment:
      - DATABASE_URL=postgres://kyora:kyora_dev_password@postgres:5432/kyora_dev?sslmode=disable
      - MEMCACHED_SERVERS=memcached:11211
    depends_on:
      - postgres
      - memcached
    volumes:
      - ./<app-name>:/app
    command: go run main.go
```

5. **Add GitHub Actions workflow** (if separate deployment):

```yaml
# .github/workflows/<app-name>-ci.yml
name: <App Name> CI

on:
  push:
    branches: [main]
    paths:
      - "<app-name>/**"
      - ".github/workflows/<app-name>-ci.yml"
  pull_request:
    branches: [main]
    paths:
      - "<app-name>/**"

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"
      - name: Test
        run: make test.<app>
      - name: Build
        run: make build.<app>
```

6. **Update root AGENTS.md:**

```markdown
## Project Structure

...

├── <app-name>/ # [Brief description]
│ ├── AGENTS.md # App-specific manifest
│ └── ...

...

## SSOT Entry Points

Backend:

- ...
- [<app-name>/AGENTS.md](<app-name>/AGENTS.md) — <App> specifics
```

---

## 2) Scaffolding New Frontend Apps

### Directory Structure

```
<app-name>/                 # e.g., storefront-web/, mobile-web/
├── AGENTS.md               # App-specific agent manifest
├── public/                 # Static assets
├── src/
│   ├── api/                # API client modules
│   ├── components/         # Shared UI components
│   ├── features/           # Feature modules
│   ├── routes/             # Page routes
│   ├── i18n/               # Translations
│   ├── lib/                # Utilities
│   ├── stores/             # Global state
│   ├── types/              # TypeScript types
│   ├── main.tsx
│   └── router.tsx
├── index.html
├── package.json
├── vite.config.ts
├── tsconfig.json
└── eslint.config.js
```

### Steps

1. **Create directory and initialize npm:**

```bash
mkdir <app-name>
cd <app-name>
npm init -y
```

2. **Install dependencies:**

```bash
npm install react react-dom
npm install -D vite @vitejs/plugin-react typescript
npm install @tanstack/react-router @tanstack/react-query
npm install tailwindcss daisyui
npm install react-i18next i18next
```

3. **Create `package.json` scripts:**

```json
{
  "name": "<app-name>",
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "preview": "vite preview",
    "typecheck": "tsc --noEmit",
    "lint": "eslint . --ext ts,tsx --report-unused-disable-directives --max-warnings 0"
  }
}
```

4. **Create AGENTS.md:**

```markdown
# <App Name>

For monorepo structure and SSOT references, see [root AGENTS.md](../AGENTS.md).

## Overview

[Brief description of what this app does]

## Tech Stack

- React + Vite
- TanStack (Router, Query, Form, Store)
- Tailwind + daisyUI
- i18n (Arabic/RTL-first)

## Commands

See root Makefile for all commands. App-specific:

- `make dev.<app>` — Run app in dev mode
- `make <app>.check` — Lint + type check
- `make <app>.build` — Production build

## SSOT References

- [frontend/_general/architecture.instructions.md] — Frontend patterns
- [frontend/_general/ui-patterns.instructions.md] — UI/RTL guidelines
- [frontend/_general/i18n.instructions.md] — i18n patterns
```

5. **Add Makefile targets** (in root `Makefile`):

```makefile
# <App Name> targets
.PHONY: dev.<app>
dev.<app>: ## Run <app> in dev mode
	@cd <app-name> && npm run dev

.PHONY: <app>.check
<app>.check: ## Lint + type check <app>
	@cd <app-name> && npm run lint && npm run typecheck

.PHONY: <app>.build
<app>.build: ## Build <app> for production
	@cd <app-name> && npm run build

.PHONY: <app>.preview
<app>.preview: ## Preview <app> production build
	@cd <app-name> && npm run preview

.PHONY: <app>.install
<app>.install: ## Install <app> dependencies
	@cd <app-name> && npm install
```

6. **Add GitHub Actions workflow:**

```yaml
# .github/workflows/<app-name>-ci.yml
name: <App Name> CI

on:
  push:
    branches: [main]
    paths:
      - "<app-name>/**"
      - ".github/workflows/<app-name>-ci.yml"
  pull_request:
    branches: [main]
    paths:
      - "<app-name>/**"

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: "20"
      - name: Install
        run: make <app>.install
      - name: Check
        run: make <app>.check
      - name: Build
        run: make <app>.build
```

7. **Update root AGENTS.md:**

```markdown
## Project Structure

...

├── <app-name>/ # [Brief description]
│ ├── AGENTS.md # App-specific manifest
│ └── ...

...

## SSOT Entry Points

Frontend:

- ...
- [<app-name>/AGENTS.md](<app-name>/AGENTS.md) — <App> specifics
```

---

## 3) Creating Shared Libraries

### When to Create Shared Libraries

**Create a shared library when:**

- Code is reused across 3+ apps
- Logic is truly generic (no app-specific dependencies)
- Library has a clear, minimal API

**Don't create shared libraries for:**

- App-specific utilities (keep in `lib/` or `platform/`)
- Code used in only 1-2 apps (duplication is acceptable)
- Rapidly changing code (wait until stable)

### Shared Library Structure

```
shared/
├── ts-utils/               # Shared TypeScript utilities
│   ├── package.json
│   ├── src/
│   │   ├── index.ts
│   │   ├── types.ts
│   │   └── utils.ts
│   └── tsconfig.json
│
├── go-utils/               # Shared Go utilities
│   ├── go.mod
│   └── ...
│
└── design-tokens/          # Shared design tokens
    ├── package.json
    ├── tokens.json
    └── build.js
```

### Steps (TypeScript Example)

1. **Create shared library:**

```bash
mkdir -p shared/ts-utils
cd shared/ts-utils
npm init -y
```

2. **Add to workspaces** (if using npm workspaces):

```json
// Root package.json
{
  "workspaces": ["portal-web", "storefront-web", "shared/ts-utils"]
}
```

3. **Link in consuming apps:**

```json
// portal-web/package.json
{
  "dependencies": {
    "@kyora/ts-utils": "workspace:*"
  }
}
```

4. **Document in root AGENTS.md:**

```markdown
## Shared Libraries

- `shared/ts-utils/` — Shared TypeScript utilities
- `shared/go-utils/` — Shared Go utilities
- `shared/design-tokens/` — Shared design tokens
```

---

## 4) Makefile Integration Patterns

### Target Naming

- **Dev**: `dev.<app>` — Run app in dev mode
- **Test**: `test.<app>` — Run app tests
- **Build**: `build.<app>` — Build app binary/bundle
- **Check**: `<app>.check` — Lint + type check (frontend)
- **Install**: `<app>.install` — Install dependencies (frontend)

### Help Text

Always include help text for new targets:

```makefile
.PHONY: dev.<app>
dev.<app>: ## Run <app> in dev mode
	@cd <app-name> && ...
```

### Validation Targets

Add validation targets for CI:

```makefile
.PHONY: validate
validate: test.<app> <app>.check ## Validate <app> (tests + checks)
```

---

## 5) CI/CD Integration

### Workflow File Naming

- **Pattern**: `<app-name>-ci.yml`
- **Location**: `.github/workflows/`

### Trigger Patterns

Always scope triggers to changed paths:

```yaml
on:
  push:
    branches: [main]
    paths:
      - "<app-name>/**"
      - ".github/workflows/<app-name>-ci.yml"
      - "shared/**" # If using shared libraries
```

### Steps

**Standard CI steps:**

1. Checkout
2. Setup language runtime
3. Install dependencies
4. Run validation (tests, lint, type check)
5. Build production artifact
6. (Optional) Deploy

**Example:**

```yaml
jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: "20"
      - run: make <app>.install
      - run: make <app>.check
      - run: make <app>.build
```

---

## 6) App Integration Checklist

**Before considering an app "fully integrated":**

- [ ] Directory created with standard structure
- [ ] AGENTS.md created (references root AGENTS.md)
- [ ] Makefile targets added (`dev`, `test`, `build`, `check`, `install`)
- [ ] Makefile help text added
- [ ] Docker-compose service added (if needed)
- [ ] GitHub Actions workflow added
- [ ] Root AGENTS.md updated with app reference
- [ ] Root README.md updated (if applicable)
- [ ] Dependencies documented
- [ ] Development workflow documented

---

## 7) References

**Monorepo patterns:**

- [monorepo/structure.instructions.md] — Directory structure
- [monorepo/workflows.instructions.md] — Makefile, docker-compose, CI/CD

**Backend patterns:**

- [backend/_general/architecture.instructions.md] — Backend structure
- [backend/_general/go-patterns.instructions.md] — Go patterns

**Frontend patterns:**

- [frontend/_general/architecture.instructions.md] — Frontend structure
- [frontend/_general/ui-patterns.instructions.md] — UI patterns

**Root:**

- [AGENTS.md] — Root agent manifest
- [Makefile] — Command definitions
- [docker-compose.dev.yml] — Infrastructure config
