# Kyora Monorepo Structure

This document describes how the Kyora monorepo is organized and the conventions for adding new projects.

## Current Structure

```
kyora/
├── backend/                 # Go backend API server
├── storefront-web/          # Customer-facing storefront (React)
├── docker-compose.dev.yml   # Local infra for development
├── Makefile                 # Root-level convenience commands
├── README.md
└── STRUCTURE.md             # This file
```

## Projects

### backend/

- Go 1.25.3 monolith.
- Architecture: `backend/internal/platform` (infrastructure) + `backend/internal/domain/*` (business modules).
- Entry point: `backend/main.go` → Cobra CLI (`backend/cmd/root.go`).
- Run server locally: `make dev.server` (root) or `cd backend && air server`.

### storefront-web/

- Customer-facing storefront (mobile-first) built with React + TypeScript + Vite.
- Styling: Tailwind CSS v4 (CSS-first) + daisyUI v5.
- RTL-first: i18next sets `html[dir]` and `html[lang]` in `storefront-web/src/App.tsx`.
- Design tokens + theme live in `storefront-web/src/index.css`.

## Adding a New Project

When adding a new top-level project (e.g., `portal-web/`, `mobile/`, `services/*`):

1. Create a top-level folder and keep the project self-contained.
2. Add a project `README.md` with setup, run, test, and deployment notes.
3. Integrate root `Makefile` targets for common workflows (dev/test/build) when helpful.
4. Keep configuration inside the project (`.env.example`, docs) and never commit secrets.

## Shared Conventions

- Prefer simple, explicit folder names at the repo root.
- Keep project dependencies isolated (Go uses `go.mod`, Node uses `package.json`).
- Repository-level docs live at root; project-level docs live inside each project.

## Local Infrastructure

- Start local dependencies with `docker compose -f docker-compose.dev.yml up -d`.
- Prefer making projects resilient when infra is missing (clear error messages, no panics).
