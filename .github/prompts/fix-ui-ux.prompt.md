---
description: Fix UI/UX issues in portal-web
agent: agent
tools: ["vscode", "execute", "read", "edit", "search", "web", "agent", "todo"]
model: Claude Opus 4.5 (copilot)
---

# Fix UI/UX Issue

You are fixing a UI/UX issue in the Kyora portal-web.

## Issue Description

${input:issue:Describe the UI/UX problem (e.g., "Button alignment broken in RTL", "Form looks cramped on mobile")}

## Instructions

Read design system rules first:

- [design-tokens.instructions.md](../instructions/design-tokens.instructions.md) for colors, spacing, typography
- [ui-implementation.instructions.md](../instructions/ui-implementation.instructions.md) for RTL, daisyUI, icons, accessibility

## UI/UX Standards

1. **RTL-First**: Design works perfectly in both RTL (Arabic) and LTR (English)
2. **Responsive**: Mobile-first, tablet, desktop breakpoints
3. **Accessibility**: ARIA labels, keyboard navigation, screen reader support
4. **Consistency**: Use design tokens (colors, spacing, typography)
5. **daisyUI**: Use semantic component classes, not arbitrary Tailwind
6. **Icons**: Lucide React icons with proper sizing

## Common UI/UX Issues

### RTL Issues

- Use `start`/`end` instead of `left`/`right` in Tailwind classes
- Use `ms-*`/`me-*` instead of `ml-*`/`mr-*`
- Use `ps-*`/`pe-*` instead of `pl-*`/`pr-*`
- Icons should flip in RTL where directional
- Text alignment: `text-start` not `text-left`

### Responsive Issues

- Use mobile-first breakpoints: base → `sm:` → `md:` → `lg:`
- Test on all screen sizes
- Stack components vertically on mobile
- Use responsive grid/flexbox

### Accessibility Issues

- Add `aria-label` to icon-only buttons
- Ensure keyboard navigation works
- Maintain color contrast ratios
- Add focus states to interactive elements

### Spacing Issues

- Use design token spacing: `space-xs`, `space-sm`, `space-md`, `space-lg`, `space-xl`
- Consistent padding/margin throughout app
- Use gap utilities for flexbox/grid spacing

### Color Issues

- Use design token colors: `primary`, `secondary`, `accent`, `neutral`
- Use daisyUI semantic colors: `btn-primary`, `alert-error`, etc.
- Maintain accessibility contrast ratios

## Workflow

1. Locate the component with UI/UX issue
2. Identify the root cause (spacing, colors, RTL, responsive, accessibility)
3. Apply fix using design tokens and daisyUI classes
4. Test in both RTL and LTR modes
5. Test on different screen sizes (mobile, tablet, desktop)
6. Verify accessibility (keyboard navigation, ARIA labels)

## Done

- UI/UX issue fixed
- Works in both RTL and LTR
- Responsive across all breakpoints
- Accessibility verified
- Design tokens used consistently
- Code is production-ready
