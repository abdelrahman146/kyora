# Kyora Monorepo Structure

This document explains the monorepo organization and guidelines for adding new projects.

## Overview

Kyora uses a monorepo structure to house multiple related projects including backend services, frontend applications, and mobile apps. This approach provides:

- **Unified versioning** - All projects are versioned together
- **Shared tooling** - Common development tools and CI/CD pipelines
- **Code sharing** - Easy sharing of types, utilities, and documentation
- **Atomic changes** - Cross-project changes can be made in a single commit

## Current Structure

```
kyora/
├── backend/              # Go backend API server
│   ├── cmd/             # CLI commands
│   ├── internal/        # Internal packages
│   ├── main.go          # Entry point
│   ├── go.mod           # Go dependencies
│   └── README.md        # Backend documentation
├── .github/             # GitHub Actions workflows
├── .vscode/             # VS Code workspace settings
├── Makefile             # Root-level build commands
├── README.md            # Main repository documentation
├── STRUCTURE.md         # This file
└── LICENSE              # Project license
```

## Planned Projects

### Frontend (Coming Soon)

Web application built with modern frontend framework.

**Proposed location:** `frontend/` or `web/`

**Suggested tech stack:**
- React/Next.js or Vue/Nuxt.js
- TypeScript
- Tailwind CSS with daisyUI
- REST API client

### Mobile (Coming Soon)

Native mobile applications for iOS and Android.

**Proposed location:** `mobile/`

**Suggested tech stack:**
- React Native or Flutter
- Shared business logic with backend
- Mobile-optimized UI

### Additional Services (Future)

Microservices for specific functionalities.

**Proposed location:** `services/`

Examples:
- `services/notifications/` - Push notification service
- `services/analytics/` - Analytics processing service
- `services/exports/` - Report generation service

## Adding a New Project

When adding a new project to the monorepo, follow these guidelines:

### 1. Directory Structure

Create a top-level directory for your project:

```bash
mkdir projectname
```

### 2. Project Documentation

Each project must have its own `README.md` with:
- Overview and description
- Tech stack
- Setup instructions
- Development workflow
- Testing guidelines
- Deployment process

### 3. Makefile Integration

Add project-specific commands to the root `Makefile`:

```makefile
# Project Development
.PHONY: dev.project
dev.project:
	@echo "Starting project development server..."
	@cd projectname && npm run dev

# Project Testing
.PHONY: test.project
test.project:
	@echo "Running project tests..."
	@cd projectname && npm test
```

Update the `help` target to include new commands.

### 4. Dependencies

Each project manages its own dependencies:
- Go projects: `go.mod` and `go.sum`
- Node.js projects: `package.json` and `package-lock.json`
- Python projects: `requirements.txt` or `Pipfile`

### 5. Configuration

Store project-specific configuration within the project directory:
- Use `.env.example` files for environment variables
- Document all required configuration in project README
- Add sensitive config files to `.gitignore`

### 6. CI/CD

Add GitHub Actions workflows in `.github/workflows/`:
- `projectname-ci.yml` - Build and test
- `projectname-deploy.yml` - Deployment pipeline

### 7. Shared Code

For code shared between projects:

**Option A:** Create a shared package directory
```
shared/
├── types/          # Shared type definitions
├── utils/          # Common utilities
└── constants/      # Shared constants
```

**Option B:** Use language-specific package management
- Go: Internal packages
- TypeScript: Shared npm workspace packages
- Python: Local pip packages

## Makefile Conventions

The root Makefile groups commands by project:

```makefile
# Backend commands (existing)
.PHONY: dev.server test test.unit test.e2e
dev.server: ...
test: ...

# Frontend commands (example)
.PHONY: dev.web test.web build.web
dev.web: ...
test.web: ...

# Mobile commands (example)
.PHONY: dev.mobile test.mobile build.mobile
dev.mobile: ...
test.mobile: ...

# Global commands
.PHONY: test.all install.all clean.all
test.all:
	@make test
	@make test.web
	@make test.mobile

install.all:
	@cd backend && go mod download
	@cd frontend && npm install
	@cd mobile && npm install

clean.all:
	@cd backend && rm -rf tmp coverage.out
	@cd frontend && rm -rf node_modules dist
	@cd mobile && rm -rf node_modules build
```

## Git Workflow

### Branch Strategy

- `main` - Production-ready code
- `develop` - Integration branch
- `feature/*` - Feature branches
- `fix/*` - Bug fix branches
- `release/*` - Release preparation

### Commit Messages

Use conventional commits format:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types:**
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation changes
- `style` - Code style changes
- `refactor` - Code refactoring
- `test` - Adding tests
- `chore` - Maintenance tasks

**Scopes:**
- `backend` - Backend changes
- `frontend` - Frontend changes
- `mobile` - Mobile changes
- `ci` - CI/CD changes
- `docs` - Documentation changes

**Examples:**
```
feat(backend): add order export feature
fix(frontend): resolve login redirect issue
docs(mobile): update setup instructions
chore(ci): update deployment workflow
```

## IDE Configuration

### VS Code

The `.vscode/` directory contains shared workspace settings:
- Extensions recommendations
- Debugger configurations
- Task definitions
- Settings overrides

Each project can have additional project-specific VS Code settings.

## Dependencies Management

### Backend (Go)

```bash
# Update dependencies
cd backend && go get -u ./...
cd backend && go mod tidy
```

### Frontend (Node.js)

```bash
# Update dependencies
cd frontend && npm update
cd frontend && npm audit fix
```

### Version Pinning

- Pin major versions to avoid breaking changes
- Document any version constraints
- Test dependency updates in CI before merging

## Testing Strategy

### Unit Tests

Each project maintains its own unit tests within the project directory.

### Integration Tests

Integration tests that span multiple projects should be in a dedicated directory:
```
tests/
├── integration/
│   ├── api-web/         # Backend + Frontend integration
│   └── api-mobile/      # Backend + Mobile integration
└── e2e/                 # End-to-end tests
```

### Running All Tests

```bash
make test.all
```

## Documentation

### Project-Level Documentation

Each project directory contains:
- `README.md` - Main documentation
- `docs/` - Additional documentation (optional)

### Repository-Level Documentation

Root directory contains:
- `README.md` - Repository overview
- `STRUCTURE.md` - This file
- `CONTRIBUTING.md` - Contribution guidelines (optional)
- `CHANGELOG.md` - Version history (optional)

## Best Practices

1. **Keep projects independent** - Projects should be loosely coupled
2. **Document everything** - Good documentation is crucial in a monorepo
3. **Consistent tooling** - Use similar development tools across projects
4. **Clear ownership** - Define which team owns which project
5. **Automated testing** - Comprehensive CI/CD for all projects
6. **Version control** - Tag releases with project prefix (e.g., `backend-v1.0.0`)
7. **Code reviews** - Review changes across all affected projects
8. **Performance monitoring** - Track build times and test execution

## Resources

- [Monorepo Best Practices](https://monorepo.tools/)
- [Google's Monorepo Approach](https://research.google/pubs/pub45424/)
- [Conventional Commits](https://www.conventionalcommits.org/)

## Questions?

For questions about the monorepo structure, please open an issue or contact the maintainers.
