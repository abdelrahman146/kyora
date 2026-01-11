---
name: i18n-parity-portal-web
description: Ensure portal-web Arabic/English translation files stay in sync (same namespaces + same key paths) and that i18n init imports match disk. Use when adding/changing UI copy.
---

## What this skill does

- Verifies `portal-web/src/i18n/ar/*.json` and `portal-web/src/i18n/en/*.json` have the same set of namespaces (filenames).
- Verifies each namespace JSON has the same nested key paths between `ar` and `en`.
- Verifies `portal-web/src/i18n/init.ts` imports match what exists on disk.

## When to use

- You added/renamed translation keys in `portal-web`.
- You added a new translation namespace JSON file.
- You suspect missing Arabic/English strings are causing runtime issues.

## How to run (manual)

From the repo root:

- `node .github/skills/i18n-parity-portal-web/scripts/check-portal-i18n-parity.mjs`

Exit code:
- `0` = all good
- `1` = mismatches found

## How to fix failures

- If a namespace is missing: add the missing JSON file under both locales.
- If keys are missing: add the missing keys to the locale file.
- If `init.ts` imports are out of date: update imports and the `resources` map to include the namespace.

## Notes

- This skill intentionally checks **structure parity**, not whether translations are “good”.
- Keep namespace additions aligned with `portal-web/src/i18n/init.ts`.
