---
title: "portal-web README.md uses incorrect script name 'npm run start'"
date: 2026-01-18
status: open
impact: low
area: portal-web
tags: [documentation, developer-experience]
---

# Drift: portal-web README Uses Incorrect Script Name

## Summary

The `portal-web/README.md` file instructs developers to use `npm run start`, but the actual script in `package.json` is `npm run dev`.

## Current State

**File:** `portal-web/README.md` (lines 6-10)

```markdown
To run this application:

```bash
npm install
npm run start
```
```

**Reality:** `portal-web/package.json` has no `start` script. The dev server script is:

```json
{
  "scripts": {
    "dev": "vite dev",
    "build": "vite build",
    "preview": "vite preview"
  }
}
```

## Expected State

README should use the correct script name:

```markdown
To run this application:

```bash
npm install
npm run dev
```
```

## Impact

- **Developer Experience**: New developers following README will get error "missing script: start"
- **Documentation Accuracy**: README doesn't match actual implementation
- **Onboarding Friction**: Slows down new team members

## Why This Matters

The README is the first touchpoint for developers. Incorrect commands create immediate friction and reduce trust in documentation quality.

## Suggested Fix

Update `portal-web/README.md`:

```diff
 To run this application:
 
 ```bash
 npm install
-npm run start
+npm run dev
 ```
```

## Files to Change

1. `portal-web/README.md` (line 9)

## Related

- Package.json scripts: `portal-web/package.json`
- SSOT: `.github/instructions/portal-web-development.instructions.md` (correctly documents `npm run dev`)

## Priority

**Low** - Documentation issue only, doesn't affect functionality. However, should be fixed to improve developer experience.
