---
name: brd-writer
description: "Writes a Kyora BRD (Business Requirements Document) from stakeholder requirements, optimized for Kyora’s mobile-first, Arabic/RTL-first DM-commerce customers. Use when you want a ready-to-build spec stored under /brds with a clear status workflow."
allowed-tools: ["read", "search", "edit"]
---

# Skill: BRD Writer (Kyora)

## When to use

- A stakeholder gives a feature/idea/problem in simple terms.
- You need a high-signal BRD that engineering can pick up later.

## Inputs you must collect

Ask short, non-technical questions (only what’s necessary):

1. Who is the user (seller, team member)?
2. Where does the workflow start (WhatsApp/Instagram/etc.)?
3. What is the #1 outcome the user wants?
4. What is the biggest confusion/pain today?
5. Any must-not-break constraints (payments, inventory, privacy)?
6. What does “done” look like for the stakeholder?

## Output requirements

1. Create a new file under `brds/`:

- `brds/BRD-YYYY-MM-DD-<short-slug>.md`

2. Use the template:

- `brds/BRD_TEMPLATE.md`

3. Keep language plain and customer-friendly:

- Mobile-first
- Arabic/RTL-first
- Avoid accounting jargon; use Kyora’s plain-language terms

4. Status rules:

- Start with `status: draft`
- Move to `status: ready` when reviewed
- Set `status: completed` when shipped and validated

## Quality bar

- Requirements must be testable and unambiguous.
- Every page/surface must define: primary action + empty/loading/error states.
- Include edge cases and failure handling.
- Include KPIs/events only if they are meaningful and measurable.

## Handoff hint

At the end of the BRD, include a short section titled **Handoff Notes for Engineering** with:

- Suggested owners (backend/portal-web/tests)
- Risky areas
- Dependencies
