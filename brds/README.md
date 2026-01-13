# Kyora BRDs (Business Requirements Documents)

This folder stores product-ready BRDs that engineering can pick up later.

## Status workflow

Each BRD uses YAML frontmatter with a `status` field.

Allowed statuses:

- `draft` — early, still gathering requirements
- `ready` — reviewed and ready for engineering planning
- `in-progress` — being implemented
- `completed` — shipped and validated
- `paused` — deprioritized or blocked

Update `status: completed` when the work is done.

## Naming

Use:

- `brds/BRD-YYYY-MM-DD-<short-slug>.md`

Example:

- `brds/BRD-2026-01-13-whatsapp-order-confirmation.md`

## Template

Start from [BRD_TEMPLATE.md](BRD_TEMPLATE.md).
