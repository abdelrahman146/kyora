# Development
.PHONY: dev.server
dev.server:
	@rm -rf tmp
	@air server

# Testing
.PHONY: test
test:
	@echo "Running all tests..."
	@go test ./... -v

.PHONY: test.unit
test.unit:
	@echo "Running unit tests..."
	@go test ./internal/domain/... ./internal/platform/... -v

.PHONY: test.e2e
test.e2e:
	@echo "Running E2E tests..."
	@go test ./internal/tests/e2e -v -timeout=120s

.PHONY: test.quick
test.quick:
	@echo "Running tests (no verbose)..."
	@go test ./...

# Coverage
.PHONY: test.coverage
test.coverage:
	@echo "Running tests with coverage..."
	@go test ./... -cover -coverprofile=coverage.out
	@echo "\nCoverage summary:"
	@go tool cover -func=coverage.out | tail -1

.PHONY: test.coverage.html
test.coverage.html:
	@echo "Generating HTML coverage report..."
	@go test ./... -cover -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@echo "Open coverage.html in your browser to view the report"

.PHONY: test.coverage.view
test.coverage.view: test.coverage.html
	@echo "Opening coverage report in browser..."
	@open coverage.html 2>/dev/null || xdg-open coverage.html 2>/dev/null || echo "Please open coverage.html manually"

.PHONY: test.e2e.coverage
test.e2e.coverage:
	@echo "Running E2E tests with coverage..."
	@go test ./internal/tests/e2e -v -timeout=120s -cover -coverprofile=e2e_coverage.out
	@echo "\nE2E Coverage summary:"
	@go tool cover -func=e2e_coverage.out | tail -1

# Clean
.PHONY: clean.coverage
clean.coverage:
	@echo "Cleaning coverage reports..."
	@rm -f coverage.out coverage.html e2e_coverage.out
	@echo "Coverage reports cleaned"

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  dev.server           - Run development server with live reload"
	@echo ""
	@echo "Testing:"
	@echo "  test                 - Run all tests (verbose)"
	@echo "  test.unit            - Run unit tests only"
	@echo "  test.e2e             - Run E2E tests only"
	@echo "  test.quick           - Run all tests (no verbose)"
	@echo ""
	@echo "Coverage:"
	@echo "  test.coverage        - Run tests with coverage report"
	@echo "  test.coverage.html   - Generate HTML coverage report"
	@echo "  test.coverage.view   - Generate and open HTML coverage in browser"
	@echo "  test.e2e.coverage    - Run E2E tests with coverage"
	@echo ""
	@echo "Clean:"
	@echo "  clean.coverage       - Remove coverage report files"