---
description: Kyora Brand Key — Foundation for all voice, messaging, and positioning decisions
applyTo: "**/*"
---

# Brand Key Model: Kyora

**Purpose**: Single source of truth for brand strategy, voice, positioning, and messaging across all Kyora projects.

**When to use**: All content decisions (UI copy, marketing, docs), UX tone, competitive positioning, feature naming.

---

## 1. Brand Root Strength

**Founder reality Kyora is built on:**

Social commerce entrepreneurs in the Middle East run real businesses through DMs, but they are forced to manage operations and money with tools that feel complex, scattered, and stressful. Kyora exists to make that business feel controlled, clear, and "handled" without needing business or accounting knowledge.

**Why this matters**: This truth drives every product decision — if a feature requires "accounting knowledge", it violates our brand root.

---

## 2. Competitive Environment

**What Kyora competes against (in the user's mind):**

| Category         | Examples                                     | User perception                                      |
| ---------------- | -------------------------------------------- | ---------------------------------------------------- |
| Manual tracking  | WhatsApp notes, Instagram screenshots        | "Messy, but I understand it"                         |
| Spreadsheets     | Excel, Google Sheets                         | "Intimidating, error-prone, feels like homework"     |
| Generic tools    | POS apps, basic invoicing/inventory apps     | "Doesn't fit how I actually sell (DMs)"              |
| "Proper" systems | ERPs, accounting software (QuickBooks, Xero) | "Too complex, not mobile, not Arabic, too expensive" |

**Category tension:**

Entrepreneurs want to feel professional and in control, but fear complexity and making mistakes.

**Decision guidance**: When designing features, ask: "Does this feel more like a spreadsheet or like Instagram?" Choose the latter.

---

## 3. Target Audience

**Primary user profile:**

- **Who**: Solo entrepreneurs, side hustlers, home-based makers, micro-teams (2–5 people)
- **Where**: Middle East (Arabic-speaking, mobile-heavy)
- **How they sell**: Instagram, WhatsApp, TikTok, Facebook DMs
- **Tech literacy**: Low to moderate
- **Context**: Time-poor, juggling production + selling + delivery + family

**What they avoid**: "Finance" tools (feel like they're for bigger businesses), anything that requires setup or training.

**See also**: `kyora/target-customer.instructions.md` for deeper psychographic profile.

---

## 4. Consumer Insight

**Human truth Kyora should be built around:**

> "I'm good at selling and making products, but I don't know how to track everything. I worry I'm losing money or messing up, and I don't want to feel stupid trying to 'do business properly.' I just want someone to quietly keep things organized and tell me what matters."

**What this means for product**:

- Don't make them "learn business" → Automate the business mechanics
- Don't make them feel stupid → Use plain language, no judgment
- Tell them what matters → Proactive insights, "What to do next"

---

## 5. Benefits

### Functional Benefits (what Kyora does for them)

- Captures orders fast from any channel and tracks status clearly
- Shows stock and alerts before items run out
- Automatically keeps the business record (money in/out, profit, cash in hand)
- Turns daily activity into clear, usable answers: best seller, top customers, what's working
- Keeps records ready for taxes/reporting without forcing the user to "learn accounting"

### Emotional Benefits (how it makes them feel)

- **Peace of mind**: "I'm not missing things."
- **Confidence**: "I know where my business stands."
- **Relief**: "It's handled, and I can focus on selling."
- **Pride**: "My business feels professional."

### Self-Expressive Benefits (what it says about them)

- "I run my business smartly, even if I'm small."
- "I'm organized and in control."

**Feature test**: Ask "What emotional benefit does this deliver?" If the answer is "None", reconsider the feature.

---

## 6. Values and Personality

### Core Brand Values

| Value                          | What it means in practice                                          |
| ------------------------------ | ------------------------------------------------------------------ |
| Simplicity over sophistication | Fewer fields, fewer steps, fewer decisions                         |
| Clarity over complexity        | Plain language, visible status, obvious next steps                 |
| Support without judgment       | Never make users feel "wrong" — guide them to fix                  |
| Trust and responsibility       | The business is "safe" with Kyora — no data loss, accurate records |
| Local reality respect          | DM-first, mobile-first, Arabic-first — not "adapted" but native    |

### Brand Personality (how Kyora should feel)

- **Calm, discreet, dependable** — like a reliable assistant who works quietly in the background
- **Practical, clear, encouraging** — "Here's what you need to do" vs "Here's a wall of options"
- **Never preachy, never technical** — no lectures about "best practices" or "proper accounting"
- **"Quiet expert" energy** — knows more than you, but never makes you feel small about it

**Wrong tones to avoid**:

- ❌ Corporate/stuffy: "Please ensure all fields are populated prior to submission"
- ❌ Overly casual: "Yo! Let's add that product 🔥"
- ❌ Technical: "Configure your chart of accounts"
- ❌ Patronizing: "Don't worry, it's easy!"

**Right tone examples**:

- ✅ "Add your first product to get started"
- ✅ "3 items are low in stock"
- ✅ "You made 1,200 SAR profit this month"

---

## 7. Tone Guidelines

### Use These Words

| Category  | Preferred terms                                     |
| --------- | --------------------------------------------------- |
| Financial | Profit, Cash in hand, Money in, Money out           |
| Inventory | Low stock, Best seller, Out of stock                |
| Actions   | Add, Save, Mark as paid, What to do next            |
| Status    | Pending, Completed, Paid, Shipped                   |
| Guidance  | Here's what's working, You're all set, Almost there |

### Avoid These Words

| Category          | Forbidden terms (and why)                                      |
| ----------------- | -------------------------------------------------------------- |
| Accounting jargon | Ledger, Accrual, EBITDA, COGS, Depreciation (intimidating)     |
| Technical terms   | Configure, Initialize, Deploy, Integrate (sounds like IT work) |
| Formal business   | Revenue, Expenditure, Assets, Liabilities (feels corporate)    |
| Complexity cues   | Advanced settings, Custom fields, Workflows (scary)            |

### Sentence Structure

- **Keep sentences short**: Prefer action + outcome. "Save your changes" not "Please ensure changes are saved."
- **Verb-first CTAs**: "Add order", "View profit", "Mark as paid" (not "Order creation", "Profitability view")
- **Active voice**: "You made 1,200 SAR" not "1,200 SAR was made"

### Arabic/RTL Mindset

- UI language should feel **native**, not translated
- Phrasing should respect Arabic sentence flow and cultural norms
- Numbers, dates, currency should follow locale conventions
- See: `.github/instructions/frontend/_general/i18n.instructions.md`

---

## 8. Reasons to Believe

**Proof points Kyora should consistently communicate (product-led):**

| Proof point                      | User-facing message example                           |
| -------------------------------- | ----------------------------------------------------- |
| DM-native order flow             | "Take orders from Instagram, WhatsApp, or TikTok DMs" |
| Automatic background bookkeeping | "Kyora keeps your business record automatically"      |
| Mobile-first speed               | "Add an order in under 30 seconds"                    |
| Clear business truth summaries   | "See your profit, cash in hand, and best sellers"     |
| Multi-business support           | "Run multiple businesses from one account"            |
| Tax-ready records                | "Your records are ready when you need them"           |

**Implementation note**: These should appear in onboarding, empty states, and marketing — not as features, but as reassurances.

---

## 9. Discriminator

**What Kyora does differently (in one line):**

> Kyora is the only Arabic-first, mobile-first platform built specifically for DM-driven commerce that runs the business behind the scenes and speaks to you in simple, everyday money language.

**Competitive positioning table:**

| Dimension      | Generic POS/inventory   | Accounting software   | Kyora                             |
| -------------- | ----------------------- | --------------------- | --------------------------------- |
| Channel        | In-store, e-commerce    | Any                   | **DM-first (Instagram/WhatsApp)** |
| Language       | English-first           | English-first         | **Arabic-first**                  |
| Device         | Desktop or tablet       | Desktop               | **Mobile-first**                  |
| User knowledge | Some business knowledge | Accounting knowledge  | **No business knowledge**         |
| Automation     | Manual entry            | Manual categorization | **Automatic bookkeeping**         |

---

## 10. Brand Essence

**The most compact expression of Kyora:**

**"Business clarity, without complexity."**

Use this as the north star for all product and content decisions:

- Does this feature add **clarity** or **complexity**?
- Does this copy make the business **clearer** or more **confusing**?
- Does this design feel **simple** or **complicated**?

If the answer leans toward complexity, rethink it.

---

## Decision-Making Framework

When in doubt about voice, features, or UX:

1. **Ask**: Does this make the user feel smart or stupid?
2. **Test**: Can a non-technical user understand this in under 5 seconds?
3. **Filter**: Would our user say "This is for people like me" or "This is for big businesses"?
4. **Choose**: Simplicity over sophistication, always.

---

## Related Documentation

- **Target Customer**: `.github/instructions/kyora/target-customer.instructions.md` — Deep user profile
- **UX Strategy**: `.github/instructions/kyora/ux-strategy.instructions.md` — How brand translates to interaction design
- **Business Model**: `.github/instructions/kyora/business-model.instructions.md` — Market positioning
- **i18n/Copy**: `.github/instructions/frontend/_general/i18n.instructions.md` — Translation rules
