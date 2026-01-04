---
description: Debug and fix issues in portal-web (React dashboard)
agent: agent
tools: ["vscode", "execute", "read", "edit", "search", "web", "agent", "todo"]
---

# Debug Portal Web Issue

You are debugging an issue in the Kyora portal-web (React TanStack dashboard).

## Issue Description

${input:issue:Describe the bug you're experiencing (e.g., "Form submission not working", "Component not rendering")}

## Instructions

Read relevant architecture rules first:

- [portal-web-architecture.instructions.md](../instructions/portal-web-architecture.instructions.md) for tech stack patterns
- [portal-web-development.instructions.md](../instructions/portal-web-development.instructions.md) for workflow

## Debugging Workflow

1. **Locate the Problem**

   - Check browser console for errors
   - Search for error messages in codebase
   - Find the relevant component/route/API call
   - Read surrounding code for context

2. **Identify Root Cause**

   - Check for: Type errors, null/undefined, API errors, state issues
   - Verify TanStack Query cache invalidation
   - Check form validation (Zod schema)
   - Review HTTP requests (network tab)
   - Verify business scoping: business-owned requests must include the correct `businessDescriptor` in both UI route params and API paths (`v1/businesses/${businessDescriptor}/...`)
   - Check i18n keys exist for both locales

3. **Implement Fix**

   - Fix the root cause (not just symptoms)
   - Add null checks where needed
   - Fix TypeScript types if mismatched
   - Improve error handling
   - Update Zod schemas if validation failing

4. **Test the Fix**

   - Run tests: `cd portal-web && npm run test`
   - Run type check: `cd portal-web && npm run type-check`
   - Run lint: `cd portal-web && npm run lint`
   - Test manually: `npm run dev`
   - Test in both LTR (English) and RTL (Arabic)
   - Test responsive design (mobile/tablet/desktop)
   - Verify accessibility (keyboard navigation)

5. **Prevent Regression**
   - Add test case covering the bug
   - Update types/schemas if needed

## Common Portal Web Issues

- **API Errors**: Check Ky client setup, auth tokens, request/response types
- **Form Issues**: Check Zod validation, field components, form state
- **Routing Issues**: Check TanStack Router config, route parameters
- **State Issues**: Check TanStack Store (auth/business/metadata/onboarding), TanStack Query cache
- **UI Issues**: Check daisyUI classes, RTL support, responsive classes
- **i18n Issues**: Check translation keys exist in both `src/i18n/en/**` and `src/i18n/ar/**`
- **Type Errors**: Check Zod schema matches API response

## Done

- Root cause identified and fixed
- Tests pass
- No TypeScript errors
- Works in both RTL and LTR
- Responsive design verified
- Code is production-ready
