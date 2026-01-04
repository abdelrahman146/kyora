---
description: Debug and fix issues in the backend (Go API)
agent: agent
tools: ["vscode", "execute", "read", "edit", "search", "web", "agent", "todo"]
---

# Debug Backend Issue

You are debugging an issue in the Kyora backend (Go monolith).

## Issue Description

${input:issue:Describe the bug or error you're experiencing (e.g., "Getting 500 error when creating new order")}

## Instructions

Read relevant architecture rules first:

- [backend-core.instructions.md](../instructions/backend-core.instructions.md) for architecture patterns
- [backend-testing.instructions.md](../instructions/backend-testing.instructions.md) if debugging tests

## Debugging Workflow

1. **Locate the Problem**

   - Search for error messages or stack traces
   - Find the relevant handler/service/repository
   - Read surrounding code for context

2. **Identify Root Cause**

   - Check for: SQL errors, nil pointer dereferences, missing validations
   - Verify workspace_id filtering (multi-tenancy)
   - Check error handling and logging
   - Review database queries and indexes

3. **Implement Fix**

   - Fix the root cause (not just symptoms)
   - Add defensive checks where needed
   - Improve error messages for debugging
   - Add validation if missing

4. **Test the Fix**

   - Run affected tests: `cd backend && go test ./...`
   - Test manually with `make run` or `make dev`
   - Verify fix doesn't break other features

5. **Prevent Regression**
   - Add test case covering the bug
   - Update validation rules if needed
   - Document fix in code if complex

## Common Backend Issues

- **500 Errors**: Check logs, nil pointer, SQL syntax, missing joins
- **401/403**: Check auth middleware, workspace permissions
- **400 Errors**: Check request validation, required fields
- **Slow Queries**: Check indexes, N+1 queries, missing `workspace_id` filter
- **Panics**: Check nil dereferences, type assertions, array bounds

## Done

- Root cause identified and fixed
- Tests pass
- No new errors introduced
- Code is production-ready
