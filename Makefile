### Kyora Monorepo Makefile
### - Keep targets simple and discoverable: `make help`

.DEFAULT_GOAL := help

ROOT_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
BACKEND_DIR := $(ROOT_DIR)/backend
PORTAL_DIR := $(ROOT_DIR)/portal-web
COMPOSE_FILE := $(ROOT_DIR)/docker-compose.dev.yml

GO ?= go
NPM ?= npm
AIR ?= air

SWAG_VERSION ?= v1.16.4
SWAG := $(GO) run github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION)

DC := docker compose -f $(COMPOSE_FILE)

PORTAL_PORT ?= 3000
SEED_SIZE ?= large
SEED_CLEAN ?= --clean

.PHONY: help
help: ## Show all targets
	@awk 'BEGIN {FS = ":.*##"; printf "\nKyora — Make targets\n\n"} \
		/^[a-zA-Z0-9_.-]+:.*##/ {printf "  %-22s %s\n", $$1, $$2} \
		END {printf "\nExamples:\n  make infra.up\n  make dev\n  make test\n  make portal.check\n\n"}' $(MAKEFILE_LIST)

## ----------------------------
## Infra (Docker Compose)
## ----------------------------

.PHONY: infra.up
infra.up: ## Start local infra (postgres, memcached, stripe-mock)
	@$(DC) up -d --remove-orphans

.PHONY: infra.down
infra.down: ## Stop local infra
	@$(DC) down

.PHONY: infra.reset
infra.reset: ## Stop infra and delete volumes (DANGEROUS)
	@$(DC) down -v

.PHONY: infra.ps
infra.ps: ## Show infra container status
	@$(DC) ps

.PHONY: infra.logs
infra.logs: ## Tail infra logs
	@$(DC) logs -f --tail=100

.PHONY: infra.restart
infra.restart: ## Restart local infra
	@$(DC) down
	@$(DC) up -d --remove-orphans

.PHONY: db.psql
db.psql: ## Open a psql session inside the postgres container
	@$(DC) exec -it postgres psql -U postgres -d kyora

## ----------------------------
## Dev
## ----------------------------

.PHONY: dev.ip
dev.ip: ## Print your current LAN IP (useful for phone testing)
	@IFACE=$$(route -n get default 2>/dev/null | awk '/interface:/{print $$2}' | head -n 1); \
	if [ -z "$$IFACE" ]; then IFACE=en0; fi; \
	IP=$$(ipconfig getifaddr $$IFACE 2>/dev/null || ipconfig getifaddr en0 2>/dev/null || ipconfig getifaddr en1 2>/dev/null); \
	if [ -z "$$IP" ]; then echo "Could not determine LAN IP"; exit 1; fi; \
	echo $$IP

.PHONY: dev.server
dev.server: ## Run backend development server with live reload (air)
	@echo "Starting backend development server"
	@cd $(BACKEND_DIR) && rm -rf tmp && $(AIR) server

.PHONY: dev.portal
dev.portal: ## Run portal dev server (LAN-friendly)
	@echo "Starting portal development server"
	@LAN_IP=$$($(MAKE) -s dev.ip) || true; \
	if [ -n "$$LAN_IP" ]; then echo "Portal URL (LAN): http://$$LAN_IP:$(PORTAL_PORT)"; fi
	@cd $(PORTAL_DIR) && \
	VITE_DEV_PORT=$(PORTAL_PORT) \
	VITE_DISABLE_TANSTACK_DEVTOOLS=$${VITE_DISABLE_TANSTACK_DEVTOOLS:-$${DISABLE_DEVTOOLS:-}} \
	$(NPM) run dev

.PHONY: dev
dev: ## Start backend + portal (parallel)
	@echo "Starting backend + portal"
	@$(MAKE) -j2 dev.server dev.portal

.PHONY: dev.infra
dev.infra: infra.up ## Start infra then start backend + portal
	@$(MAKE) dev

## ----------------------------
## Backend
## ----------------------------

.PHONY: backend.deps
backend.deps: ## Download backend dependencies
	@cd $(BACKEND_DIR) && $(GO) mod download

.PHONY: backend.tidy
backend.tidy: ## Tidy backend go.mod/go.sum
	@cd $(BACKEND_DIR) && $(GO) mod tidy

.PHONY: backend.fmt
backend.fmt: ## Format backend (gofmt via `go fmt`)
	@cd $(BACKEND_DIR) && $(GO) fmt ./...

.PHONY: backend.vet
backend.vet: ## Vet backend (basic static checks)
	@cd $(BACKEND_DIR) && $(GO) vet ./...

## ----------------------------
## Portal Web
## ----------------------------

.PHONY: portal.install
portal.install: ## Install portal dependencies
	@cd $(PORTAL_DIR) && $(NPM) install

.PHONY: portal.lint
portal.lint: ## Lint portal (eslint)
	@cd $(PORTAL_DIR) && $(NPM) run lint

.PHONY: portal.type-check
portal.type-check: ## Type-check portal (tsc --noEmit)
	@cd $(PORTAL_DIR) && $(NPM) run type-check

.PHONY: portal.format
portal.format: ## Format portal (prettier --write)
	@cd $(PORTAL_DIR) && $(NPM) run format -- --write .

.PHONY: portal.check
portal.check: ## Portal checks (lint + type-check)
	@cd $(PORTAL_DIR) && $(NPM) run lint && $(NPM) run type-check

.PHONY: portal.test
portal.test: ## Run portal tests (vitest)
	@cd $(PORTAL_DIR) && $(NPM) run test

.PHONY: portal.build
portal.build: ## Build portal for production
	@cd $(PORTAL_DIR) && $(NPM) run build

.PHONY: portal.preview
portal.preview: ## Preview production portal build locally
	@cd $(PORTAL_DIR) && $(NPM) run preview

.PHONY: clean.portal
clean.portal: ## Remove portal build output (dist)
	@rm -rf $(PORTAL_DIR)/dist

## ----------------------------
## OpenAPI (Backend)
## ----------------------------

.PHONY: openapi
openapi: ## Generate backend OpenAPI docs (Swaggo)
	@echo "Generating backend OpenAPI docs (Swaggo)..."
	@cd $(BACKEND_DIR) && rm -rf docs
	@cd $(BACKEND_DIR) && $(SWAG) init \
		-g main.go \
		-o ./docs \
		--parseDependency \
		--parseInternal \
		--parseDepth 2 \
		--outputTypes json,yaml
	@echo "OpenAPI generated: backend/docs/swagger.json and backend/docs/swagger.yaml"

.PHONY: openapi.check
openapi.check: openapi ## Verify OpenAPI output is committed (fails on git diff)
	@command -v git >/dev/null 2>&1 || { echo "git is required for openapi.check"; exit 1; }
	@git diff --exit-code -- backend/docs
	@echo "OpenAPI is up-to-date"

.PHONY: openapi.verify
openapi.verify: openapi.check ## Alias for openapi.check

## ----------------------------
## Backend Testing
## ----------------------------

.PHONY: test
test: ## Run all backend tests (verbose)
	@echo "Running all backend tests..."
	@cd $(BACKEND_DIR) && $(GO) test ./... -v

.PHONY: test.unit
test.unit: ## Run backend unit tests only
	@echo "Running backend unit tests..."
	@cd $(BACKEND_DIR) && $(GO) test ./internal/domain/... ./internal/platform/... -v

.PHONY: test.e2e
test.e2e: ## Run backend E2E tests only
	@echo "Running backend E2E tests..."
	@cd $(BACKEND_DIR) && $(GO) test ./internal/tests/e2e -v -timeout=120s

.PHONY: test.quick
test.quick: ## Run all backend tests (no verbose)
	@echo "Running backend tests (no verbose)..."
	@cd $(BACKEND_DIR) && $(GO) test ./...

## ----------------------------
## Coverage
## ----------------------------

.PHONY: test.coverage
test.coverage: ## Run backend tests with coverage summary
	@echo "Running backend tests with coverage..."
	@cd $(BACKEND_DIR) && $(GO) test ./... -cover -coverprofile=coverage.out
	@echo "\nCoverage summary:"
	@cd $(BACKEND_DIR) && $(GO) tool cover -func=coverage.out | tail -1

.PHONY: test.coverage.html
test.coverage.html: ## Generate backend HTML coverage report
	@echo "Generating HTML coverage report..."
	@cd $(BACKEND_DIR) && $(GO) test ./... -cover -coverprofile=coverage.out
	@cd $(BACKEND_DIR) && $(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: backend/coverage.html"

.PHONY: test.coverage.view
test.coverage.view: test.coverage.html ## Generate and open backend HTML coverage report
	@echo "Opening coverage report in browser..."
	@open backend/coverage.html 2>/dev/null || xdg-open backend/coverage.html 2>/dev/null || echo "Please open backend/coverage.html manually"

.PHONY: test.e2e.coverage
test.e2e.coverage: ## Run backend E2E tests with coverage summary
	@echo "Running backend E2E tests with coverage..."
	@cd $(BACKEND_DIR) && $(GO) test ./internal/tests/e2e -v -timeout=120s -cover -coverprofile=e2e_coverage.out
	@echo "\nE2E Coverage summary:"
	@cd $(BACKEND_DIR) && $(GO) tool cover -func=e2e_coverage.out | tail -1

## ----------------------------
## Seed
## ----------------------------

.PHONY: seed
seed: ## Seed backend database (vars: SEED_SIZE=large|small, SEED_CLEAN=--clean|)
	@echo "Seeding backend database with initial data..."
	@cd $(BACKEND_DIR) && STRIPE_BASE_URL="http://localhost:12111" $(GO) run . seed $(SEED_CLEAN) --size $(SEED_SIZE)
	@echo "Database seeding completed"

## ----------------------------
## Clean
## ----------------------------

.PHONY: clean.coverage
clean.coverage: ## Remove backend coverage report files
	@echo "Cleaning backend coverage reports..."
	@cd $(BACKEND_DIR) && rm -f coverage.out coverage.html e2e_coverage.out
	@echo "Coverage reports cleaned"

.PHONY: clean.backend
clean.backend: ## Remove backend build artifacts (tmp, build-errors.log)
	@echo "Cleaning backend build artifacts..."
	@cd $(BACKEND_DIR) && rm -rf tmp build-errors.log
	@echo "Backend artifacts cleaned"

.PHONY: clean
clean: clean.backend clean.coverage clean.portal ## Clean common build artifacts

## ----------------------------
## Monorepo meta
## ----------------------------

.PHONY: install
install: backend.deps portal.install ## Install all dependencies (Go + portal)

.PHONY: fmt
fmt: backend.fmt portal.format ## Format all (Go + portal)

.PHONY: lint
lint: backend.vet portal.lint ## Lint-ish checks (go vet + eslint)

.PHONY: check
check: lint test.quick portal.type-check ## Fast checks (lint + backend tests + portal type-check)

.PHONY: doctor
doctor: ## Check local tooling (go, docker, npm, air)
	@set -e; \
	missing=0; \
	for bin in "$(GO)" docker "$(NPM)" "$(AIR)"; do \
		if ! command -v $$bin >/dev/null 2>&1; then echo "Missing: $$bin"; missing=1; fi; \
	done; \
	if [ $$missing -ne 0 ]; then exit 1; fi; \
	echo "go: $$($(GO) version)"; \
	echo "docker: $$(docker --version)"; \
	echo "npm: $$(npm --version)"; \
	echo "air: $$($(AIR) -v 2>/dev/null || $(AIR) --version 2>/dev/null || echo \"(ok)\")"