# Backend Development
.PHONY: dev.server
dev.server:
	@echo "Starting backend development server..."
	@cd backend && rm -rf tmp && air server

# Backend OpenAPI
.PHONY: openapi
openapi:
	@echo "Generating backend OpenAPI docs (Swaggo)..."
	@cd backend && rm -rf docs
	@cd backend && go run github.com/swaggo/swag/cmd/swag@v1.16.4 init \
		-g main.go \
		-o ./docs \
		--parseDependency \
		--parseInternal \
		--parseDepth 2 \
		--outputTypes json,yaml
	@echo "OpenAPI generated: backend/docs/swagger.json and backend/docs/swagger.yaml"

# Backend Testing
.PHONY: test
test:
	@echo "Running all backend tests..."
	@cd backend && go test ./... -v

.PHONY: test.unit
test.unit:
	@echo "Running backend unit tests..."
	@cd backend && go test ./internal/domain/... ./internal/platform/... -v

.PHONY: test.e2e
test.e2e:
	@echo "Running backend E2E tests..."
	@cd backend && go test ./internal/tests/e2e -v -timeout=120s

.PHONY: test.quick
test.quick:
	@echo "Running backend tests (no verbose)..."
	@cd backend && go test ./...

# Backend Coverage
.PHONY: test.coverage
test.coverage:
	@echo "Running backend tests with coverage..."
	@cd backend && go test ./... -cover -coverprofile=coverage.out
	@echo "\nCoverage summary:"
	@cd backend && go tool cover -func=coverage.out | tail -1

.PHONY: test.coverage.html
test.coverage.html:
	@echo "Generating HTML coverage report..."
	@cd backend && go test ./... -cover -coverprofile=coverage.out
	@cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: backend/coverage.html"
	@echo "Open backend/coverage.html in your browser to view the report"

.PHONY: test.coverage.view
test.coverage.view: test.coverage.html
	@echo "Opening coverage report in browser..."
	@open backend/coverage.html 2>/dev/null || xdg-open backend/coverage.html 2>/dev/null || echo "Please open backend/coverage.html manually"

.PHONY: test.e2e.coverage
test.e2e.coverage:
	@echo "Running backend E2E tests with coverage..."
	@cd backend && go test ./internal/tests/e2e -v -timeout=120s -cover -coverprofile=e2e_coverage.out
	@echo "\nE2E Coverage summary:"
	@cd backend && go tool cover -func=e2e_coverage.out | tail -1

# Backend Clean
.PHONY: clean.coverage
clean.coverage:
	@echo "Cleaning backend coverage reports..."
	@cd backend && rm -f coverage.out coverage.html e2e_coverage.out
	@echo "Coverage reports cleaned"

.PHONY: clean.backend
clean.backend:
	@echo "Cleaning backend build artifacts..."
	@cd backend && rm -rf tmp build-errors.log
	@echo "Backend artifacts cleaned"

seed:
	@echo "Seeding backend database with initial data..."
	@cd backend && STRIPE_BASE_URL="http://localhost:12111" go run . seed --clean --size large
	@echo "Database seeding completed"

# Help
.PHONY: help
help:
	@echo "Kyora Monorepo - Available targets:"
	@echo ""
	@echo "Backend Development:"
	@echo "  dev.server           - Run backend development server with live reload"
	@echo ""
	@echo "Backend OpenAPI:"
	@echo "  openapi              - Generate backend OpenAPI docs (Swaggo)"
	@echo "  openapi.verify       - Generate and verify OpenAPI covers real routes (E2E)"
	@echo ""
	@echo "Backend Testing:"
	@echo "  test                 - Run all backend tests (verbose)"
	@echo "  test.unit            - Run backend unit tests only"
	@echo "  test.e2e             - Run backend E2E tests only"
	@echo "  test.quick           - Run all backend tests (no verbose)"
	@echo ""
	@echo "Backend Coverage:"
	@echo "  test.coverage        - Run backend tests with coverage report"
	@echo "  test.coverage.html   - Generate HTML coverage report"
	@echo "  test.coverage.view   - Generate and open HTML coverage in browser"
	@echo "  test.e2e.coverage    - Run backend E2E tests with coverage"
	@echo ""
	@echo "Clean:"
	@echo "  clean.coverage       - Remove backend coverage report files"
	@echo "  clean.backend        - Remove backend build artifacts"
	@echo ""
	@echo "General:"
	@echo "  help                 - Show this help message"