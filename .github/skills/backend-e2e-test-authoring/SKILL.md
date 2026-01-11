---
name: backend-e2e-test-authoring
description: Write or extend Kyora backend E2E tests using the existing container-based harness, helpers, and HTTP client patterns. Use when adding/changing backend behavior.
---

## What this skill does

A practical workflow for adding E2E coverage without reinventing infrastructure.

## Where tests live

- E2E suite entrypoint: `backend/internal/tests/e2e/main_test.go`
- Tests: `backend/internal/tests/e2e/*_test.go`
- Shared infra + helpers: `backend/internal/tests/testutils/` and domain `*_helpers_test.go`

## Workflow

1. Find the closest existing test
   - Prefer extending an existing domain test file or reusing its helpers.

2. Reuse helpers for setup
   - Look for `<domain>_helpers_test.go` in `backend/internal/tests/e2e/`.
   - Avoid reaching into DB directly unless existing helpers already do.

3. Make API calls through the existing patterns
   - Keep HTTP calls consistent with the suite’s conventions (base URL, auth headers, etc.).

4. Assert tenant isolation
   - When relevant, create a second workspace/business and ensure cross-access is denied.

5. Keep runtime bounded
   - Don’t spin up extra containers; use the shared `TestMain` environment.

## References (SSOT)

- Testing guidelines: `.github/instructions/backend-testing.instructions.md`
- Backend patterns: `.github/instructions/backend-core.instructions.md`

## Manual run command

- `make test.e2e`

(Run manually; it boots the container-based E2E environment.)
