---
description: Kyora UX Strategy — Core principles, interaction patterns, information architecture, emotional design, decision framework
applyTo: "portal-web/**,storefront-web/**,mobile-web/**"
---

# Kyora UX Strategy

**Purpose**: How to design for Kyora's target customer — mobile-first, Arabic-first, confidence-building, plain language.

**When to use**: Designing new features, scoping UX flows, resolving design conflicts, writing UI copy.

---

## 1. Core UX Principles

### 1.1 Mobile-First (90% of Usage)

**What it means**:

- Design for 375px viewport first, progressively enhance for larger screens
- One-column layouts by default
- Primary actions at thumb-friendly positions (bottom, center)
- Touch targets ≥44px (prefer 48px+)
- Full-width CTAs on mobile

**Why it matters**: 80%+ of users interact on smartphones; desktop is secondary.

**Implementation rules**:

```tsx
// ✅ Correct — Mobile default, desktop enhancement
<button className="btn btn-primary w-full md:w-auto">Submit</button>

// ❌ Wrong — Desktop default, mobile override
<button className="btn btn-primary w-auto max-md:w-full">Submit</button>
```

**Decision test**: "Would this work on a phone while standing in a busy market?" If no, redesign.

---

### 1.2 Arabic-First (Not Arabic-Added)

**What it means**:

- RTL layout is the default (not an afterthought)
- Arabic phrasing drives English translation (not the reverse)
- Icons, spacing, alignment all respect RTL reading flow
- Mixed-direction content (phone numbers, IDs) explicitly marked

**Why it matters**: Users must feel the UI is "native", not translated.

**Implementation rules**:

- Use logical properties (`start`/`end`, not `left`/`right`)
- Rotate directional arrows 180° in RTL (`isRTL ? 'rotate-180' : ''`)
- Mark LTR content explicitly: `<span dir="ltr">{phone}</span>`

**See**: `kyora/design-system.instructions.md` (section 7: RTL-First Rules)

**Decision test**: "Does this feel like it was built for Arabic speakers, or adapted from English?" If adapted, rethink.

---

### 1.3 Plain Language (No Accounting Jargon)

**What it means**:

- Use everyday money words: "Profit", "Cash in hand", "Money in/out", "Best seller"
- Avoid accounting terms: "Ledger", "Accrual", "EBITDA", "COGS", "Depreciation"
- Keep sentences short, action-oriented, verb-first

**Why it matters**: Users fear complexity and feeling stupid. Plain language builds confidence.

**Approved vocabulary**:

| Context   | Use                     | Don't Use             |
| --------- | ----------------------- | --------------------- |
| Financial | Profit, Cash in hand    | Revenue, Net income   |
| Inventory | Low stock, Out of stock | Reorder point, SKU    |
| Actions   | Add order, Save         | Create transaction    |
| Status    | Pending, Paid           | Processed, Reconciled |
| Guidance  | What to do next         | Recommended actions   |

**See**: `kyora/brand-key.instructions.md` (section 7: Tone Guidelines)

**Decision test**: "Would a high school graduate with no business training understand this?" If no, simplify.

---

### 1.4 Clarity Over Complexity

**What it means**:

- Simple flows (fewer steps, fewer decisions)
- No hidden features (obvious where to find things)
- Obvious next steps ("What to do next" prompts)
- Progressive disclosure (show essentials, hide advanced)

**Why it matters**: Users are time-poor, avoid complexity, fear mistakes.

**Flow design rules**:

- **Max 3 steps** for common tasks (add order, create product)
- **Single primary action** per screen (mobile)
- **Clear exits** (Cancel/Back always visible)
- **Confirm destructive actions** (Delete, Archive)

**Decision test**: "Can a user complete this task in under 60 seconds on their first try?" If no, simplify flow.

---

### 1.5 Confidence-Building (Not Intimidating)

**What it means**:

- Show progress, confirm actions, explain impact
- Celebrate small wins ("Order added!", "Profit this month: 1,200 SAR")
- Forgiving inputs (auto-format, smart defaults, undo)
- No judgment language ("Fix this" not "Error: Invalid input")

**Why it matters**: Users fear making mistakes, feeling stupid, being exposed as disorganized.

**Messaging guidelines**:

| Situation         | Wrong tone                     | Right tone                                 |
| ----------------- | ------------------------------ | ------------------------------------------ |
| Validation error  | "Invalid input"                | "Phone number must be 10 digits"           |
| Empty state       | "No data"                      | "Add your first order to start"            |
| Success action    | "Transaction completed"        | "Order saved!"                             |
| Low stock warning | "Stock level below threshold"  | "3 items left — restock soon?"             |
| Delete confirm    | "This action cannot be undone" | "Are you sure? This will remove the order" |

**Decision test**: "Does this message make the user feel smart or stupid?" If stupid, rewrite.

---

## 2. Interaction Patterns

### 2.1 Touch-Optimized

**Rules**:

- Touch targets ≥44×44px (prefer 48px+)
- Swipe gestures for lists (future: swipe to delete, mark as paid)
- Bottom sheets for mobile forms (not modals)
- One-handed use: Primary actions at thumb-reach (bottom), secondary at top

**Common mistakes to avoid**:

- ❌ Small icon-only buttons (<40px)
- ❌ Desktop-style dropdowns on mobile (use action sheets)
- ❌ Hover-dependent interactions (no hover on touch)

**Implementation**:

```tsx
// ✅ Correct — Adequate touch target
<button className="btn btn-ghost min-h-12 min-w-12">
  <Trash2 size={20} />
</button>

// ❌ Wrong — Too small for touch
<button className="p-1">
  <Trash2 size={16} />
</button>
```

---

### 2.2 Forgiving Inputs

**What it means**:

- Auto-format phone numbers, prices, dates
- Smart defaults (today's date, last-used customer)
- Undo actions (especially delete)
- Inline validation (show errors on blur, not on every keystroke)

**Examples**:

```tsx
// Auto-format phone number
<input type="tel" inputMode="numeric" autoComplete="tel" />
// Displays: 05XX XXX XXXX (formatted automatically)

// Smart default date
<input type="date" defaultValue={new Date().toISOString().split('T')[0]} />

// Undo delete (toast with undo button)
toast.success('Order deleted', {
  action: {
    label: 'Undo',
    onClick: () => restoreOrder(id)
  }
})
```

---

### 2.3 Keyboard UX

**Rules**:

- Use proper `inputMode`: `numeric` for prices/quantities, `tel` for phones, `email` for email
- Use `autoComplete` attributes (speeds up form filling)
- Avoid layouts that break when keyboard opens (keep submit button reachable)

**Implementation**:

```tsx
// ✅ Correct — Optimized keyboard
<input
  type="text"
  inputMode="numeric"
  autoComplete="tel"
  placeholder="05XX XXX XXXX"
/>

// ❌ Wrong — Generic keyboard
<input type="text" placeholder="Phone" />
```

---

## 3. Information Architecture

### 3.1 Progressive Disclosure

**What it means**:

- Show essentials first, hide advanced features
- "What to do next" prompts for empty states
- Contextual help (inline hints, tooltips, not separate help docs)

**Hierarchy**:

1. **Primary info**: Always visible (order status, profit, stock count)
2. **Secondary info**: Visible but de-emphasized (timestamps, secondary actions)
3. **Advanced features**: Hidden behind "Advanced" toggle or separate screen

**Example**:

```tsx
// Product form — basic fields visible, advanced collapsed
<form>
  <form.TextField name="name" label="Product Name" required />
  <form.TextField name="price" label="Price" required />

  <Accordion>
    <AccordionItem title="Advanced Options">
      <form.TextField name="sku" label="SKU (optional)" />
      <form.TextField name="barcode" label="Barcode (optional)" />
    </AccordionItem>
  </Accordion>
</form>
```

---

### 3.2 Scannable Content

**What it means**:

- Visual hierarchy (headings, spacing, grouping)
- Icons + labels (not labels alone)
- Short sentences (<15 words)
- Bullet points over paragraphs

**Layout patterns**:

- **Card-based lists**: Easier to scan than tables on mobile
- **Section dividers**: Clear visual breaks between groups
- **Icon + text**: Faster recognition than text alone

**Example**:

```tsx
// ✅ Scannable — Card with icon, title, meta
<div className="card bg-base-100 border border-base-300">
  <div className="card-body">
    <div className="flex items-center gap-3">
      <Package size={24} className="text-primary" />
      <div>
        <h3 className="font-semibold">Product Name</h3>
        <p className="text-sm text-base-content/70">12 in stock · 1,200 SAR</p>
      </div>
    </div>
  </div>
</div>

// ❌ Not scannable — Wall of text
<div>
  <p>Product Name - Stock: 12 - Price: 1,200 SAR</p>
</div>
```

---

### 3.3 Contextual Help

**What it means**:

- Inline hints (below input fields)
- Tooltips (hover/tap for details)
- Examples (show sample input)
- No separate help docs (help is embedded)

**Implementation**:

```tsx
// Inline hint
<form.TextField
  name="price"
  label="Price"
  hint="Enter the selling price (excluding delivery)"
/>

// Tooltip
<label className="flex items-center gap-2">
  Profit Margin
  <HelpCircle size={16} className="text-base-content/50" title="Profit divided by revenue" />
</label>

// Example
<form.TextField
  name="phone"
  label="Phone Number"
  placeholder="05XX XXX XXXX"
/>
```

---

## 4. Emotional Design

### 4.1 Reassurance (Peace of Mind)

**What it means**:

- "You're not missing things" messaging
- "It's handled" confirmations
- "You can focus on selling" prompts

**Examples**:

- **Empty state**: "No orders yet — add your first order to start tracking"
- **Success toast**: "Order saved! You can view it anytime in Orders."
- **Low stock alert**: "3 items left — time to restock?"

---

### 4.2 Empowerment (Confidence)

**What it means**:

- "Your business feels professional" pride
- "You're organized and in control" self-expression
- Celebrate progress (milestones, achievements)

**Examples**:

- **Dashboard summary**: "You made 1,200 SAR profit this month 🎉"
- **First order**: "Great! Your first order is saved. Keep going!"
- **Inventory healthy**: "All products in stock — you're ready to sell."

---

### 4.3 No Judgment (Support Without Shame)

**What it means**:

- Never make users feel "wrong" or "stupid"
- Guide them to fix (not blame them)
- "Let's fix this" tone (not "Error: You did this wrong")

**Examples**:

| Situation      | Wrong (judgmental)        | Right (supportive)                                   |
| -------------- | ------------------------- | ---------------------------------------------------- |
| Missing field  | "Error: Name is required" | "Add a product name to continue"                     |
| Low stock      | "Critical stock level"    | "3 items left — restock soon?"                       |
| No profit      | "Revenue too low"         | "No profit yet — keep tracking orders to see trends" |
| Delete confirm | "This will destroy data"  | "Are you sure? This will remove the order"           |

---

## 5. Content Strategy

### 5.1 Action-Oriented Copy

**What it means**:

- Verb-first CTAs: "Add order", "View profit", "Mark as paid"
- Not noun-based: "Order creation", "Profitability view", "Payment status update"

**Button copy guidelines**:

| Context     | Good         | Bad                   |
| ----------- | ------------ | --------------------- |
| Create      | Add Order    | Create New Order      |
| Save        | Save         | Submit Changes        |
| Delete      | Delete       | Remove Item           |
| View        | View Details | See More Info         |
| Mark status | Mark as Paid | Update Payment Status |

---

### 5.2 Outcome-Focused Messaging

**What it means**:

- Tell users what they'll get (not just what the feature does)
- "Know your best seller" (not "View sales analytics")
- "See what's working" (not "Generate reports")

**Examples**:

- **Analytics section**: "Know your best sellers and top customers"
- **Inventory section**: "Never run out of stock"
- **Accounting section**: "See your profit and cash in hand"

---

### 5.3 User Vocabulary (Not Business Jargon)

**What it means**:

- Use words from jobs-to-be-done (what users say, not what business textbooks say)
- "Cash in hand" (not "Liquid assets")
- "Best seller" (not "Top SKU by revenue")

**See**: `kyora/target-customer.instructions.md` (section 9: Quotes from User Research)

---

## 6. Loading, Empty, and Error States

### 6.1 Loading States

**Rules**:

- Prefer **skeletons** to spinners for page-level loading (stable layouts)
- Use **spinners** for inline actions (button loading, form submit)
- Keep layout stable (avoid jumping content)

**Implementation**:

```tsx
// Page-level loading (skeleton)
{
  isLoading ? (
    <div className="space-y-3">
      {Array.from({ length: 5 }).map((_, i) => (
        <div key={i} className="skeleton h-32 rounded-box" />
      ))}
    </div>
  ) : (
    <ProductList products={products} />
  );
}

// Button loading (spinner)
<button disabled={isPending} className="btn btn-primary">
  {isPending ? <Loader2 className="animate-spin" size={20} /> : "Save"}
</button>;
```

---

### 6.2 Empty States

**Rules**:

- Always include: **what it means** + **clear next action**
- Use friendly, encouraging tone
- Show icon (size 48–64px)

**Implementation**:

```tsx
<div className="flex flex-col items-center justify-center py-12 text-center">
  <Package size={48} className="text-base-content/30 mb-4" />
  <h3 className="text-lg font-semibold mb-2">No products yet</h3>
  <p className="text-base-content/70 mb-6">
    Add your first product to start tracking inventory
  </p>
  <button onClick={onAdd} className="btn btn-primary gap-2">
    <Plus size={20} />
    Add Product
  </button>
</div>
```

---

### 6.3 Error States

**Rules**:

- Explain **what happened** + **how to fix**
- No stack traces, no technical jargon
- Provide action (Retry, Contact Support)

**Implementation**:

```tsx
<div className="card bg-base-100 border border-error/50">
  <div className="card-body text-center">
    <AlertCircle size={48} className="text-error mx-auto mb-4" />
    <h3 className="text-lg font-semibold mb-2">Couldn't load orders</h3>
    <p className="text-base-content/70 mb-4">
      Check your internet connection and try again
    </p>
    <button onClick={retry} className="btn btn-primary">
      Retry
    </button>
  </div>
</div>
```

---

## 7. Visual Style Constraints

### 7.1 No Shadows, No Gradients

**Why**: Maintains minimal, calm aesthetic; prevents visual clutter.

**Rule**: Use **borders** and **spacing** for separation (not shadows).

**See**: `kyora/design-system.instructions.md` (section 9: Borders & Shadows)

---

### 7.2 daisyUI Semantic Classes

**Why**: Consistent theming, easier maintenance.

**Rule**: Use daisyUI classes for components (`.btn`, `.card`, `.badge`), Tailwind for spacing/layout only.

**See**: `kyora/design-system.instructions.md` (section 6: Component Patterns)

---

## 8. Responsive Patterns

### Mobile vs Desktop

| Component    | Mobile (<768px)         | Desktop (≥768px)  |
| ------------ | ----------------------- | ----------------- |
| Modal        | BottomSheet (slides up) | Dialog (centered) |
| Navigation   | Bottom bar              | Sidebar           |
| Form buttons | Full width (`w-full`)   | Auto width        |
| Page layout  | Single column           | Multi-column      |
| Lists        | Cards (stacked)         | Table or cards    |

**Implementation**:

```tsx
// Conditional rendering based on screen size
const isMobile = useMediaQuery("(max-width: 768px)");

{
  isMobile ? (
    <BottomSheet isOpen={isOpen} onClose={onClose}>
      <Form />
    </BottomSheet>
  ) : (
    <Dialog isOpen={isOpen} onClose={onClose}>
      <Form />
    </Dialog>
  );
}
```

---

## 9. Accessibility Requirements

### ARIA Attributes

- Icon-only buttons: `aria-label="Action description"`
- Form inputs: `aria-invalid`, `aria-describedby`, `aria-required`
- Modals: `role="dialog"`, `aria-modal="true"`, `aria-labelledby`
- Decorative icons: `aria-hidden="true"`

### Focus Management

- All interactive elements must have visible focus ring
- Modals/overlays must trap focus
- Keyboard navigation: Tab, Escape, Enter/Space, Arrow keys

**See**: `kyora/design-system.instructions.md` (section 12: Touch Targets & Accessibility)

---

## 10. Decision-Making Framework

When in conflict (simplicity vs features, mobile vs desktop, etc.), use this hierarchy:

### Hierarchy of Concerns

1. **Does it serve the ICP?** (Solo entrepreneurs, Middle East, DM-first)
2. **Does it reduce complexity?** (Fewer steps, fewer decisions)
3. **Does it build confidence?** (Peace of mind, not intimidation)
4. **Is it mobile-first?** (Phone-friendly, not desktop-centric)
5. **Is it Arabic-first?** (Native RTL, not translated afterthought)

### Conflict Resolution Examples

**Scenario 1**: Feature request — "Add advanced filtering (10+ filters)"

- **Conflict**: Power vs simplicity
- **Decision**: Add 3 core filters (status, date, customer), hide rest behind "Advanced"
- **Rationale**: ICP rarely needs complex filtering; keep UI clean

**Scenario 2**: UI layout — "Desktop users want multi-column dashboard"

- **Conflict**: Desktop optimization vs mobile-first
- **Decision**: Single column default, multi-column on large screens (`lg:grid-cols-2`)
- **Rationale**: Mobile-first principle (80%+ users on mobile)

**Scenario 3**: Copy — "Use 'Revenue' instead of 'Money in'"

- **Conflict**: Professional terminology vs plain language
- **Decision**: Keep "Money in"
- **Rationale**: Plain language principle (users fear accounting jargon)

---

## 11. Portal Web UX Checklist

Before completing any portal-web UI work:

- ☑ **Mobile-first**: One-column layout, full-width CTAs, touch targets ≥44px
- ☑ **Arabic-first**: No `left`/`right` classes, directional icons rotate, `dir="ltr"` for phone/IDs
- ☑ **Plain language**: No jargon ("Profit" not "EBITDA", "Cash in hand" not "Accrual")
- ☑ **Clarity**: Simple flow (≤3 steps), obvious next action, single primary CTA
- ☑ **Confidence-building**: Success confirmations, progress indicators, no judgment language
- ☑ **Loading states**: Skeletons for pages, spinners for buttons, stable layouts
- ☑ **Empty states**: Icon + message + CTA ("Add your first order to start")
- ☑ **Error states**: What happened + how to fix + retry button
- ☑ **Visual style**: Borders (not shadows), daisyUI classes, spacing for separation
- ☑ **Accessibility**: Focus visible, aria-labels, keyboard navigation

---

## Related Documentation

- **Brand Key**: `.github/instructions/kyora/brand-key.instructions.md` — Voice, tone, positioning
- **Target Customer**: `.github/instructions/kyora/target-customer.instructions.md` — Who we're designing for
- **Design System**: `.github/instructions/kyora/design-system.instructions.md` — Visual patterns, components
- **Business Model**: `.github/instructions/kyora/business-model.instructions.md` — Market context, ICP
- **Forms**: `.github/instructions/forms.instructions.md` — Form-specific patterns
- **i18n**: `.github/instructions/frontend/_general/i18n.instructions.md` — Translation rules
