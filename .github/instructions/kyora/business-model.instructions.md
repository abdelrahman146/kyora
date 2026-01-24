---
description: Kyora Business Model — Market, value proposition, revenue model, competitive advantage
applyTo: "**/*"
---

# Kyora Business Model

**Purpose**: Context for business decisions, feature prioritization, pricing discussions, and strategic planning.

**When to use**: Evaluating new features, scoping MVP vs premium capabilities, discussing market fit, planning roadmap.

---

## 1. Market Context

### Geographic Market

- **Primary**: Middle East (Saudi Arabia, UAE, Egypt, Jordan, Kuwait, Bahrain, Oman)
- **Secondary**: MENA region expansion (North Africa, Levant)
- **Future**: Other emerging markets with similar DM-commerce patterns

### Market Characteristics

- **Social commerce dominance**: Instagram, WhatsApp, TikTok are primary sales channels (not e-commerce websites)
- **DM-first selling**: Entrepreneurs take orders via direct messages, not shopping carts
- **Mobile-heavy**: 80%+ of transactions happen on mobile devices
- **Arabic-first culture**: Primary language for business communication
- **Cash-based economy**: Many transactions still cash-on-delivery (COD)
- **Informal economy**: Micro-businesses often start unregistered, formalize later

### Market Dynamics

- **Growing**: Social commerce growing 30%+ annually in Middle East
- **Underserved**: Existing tools (POS, accounting software) not built for DM workflows
- **Low digitization**: Many sellers still use WhatsApp notes, screenshots, paper notebooks
- **Trust-based**: Seller reputation and personal relationships drive sales

---

## 2. Customer Segment (ICP)

### Primary Segment: Solo Micro-Entrepreneurs

**Profile:**

- Solo entrepreneurs, side hustlers, home-based makers
- 25–45 years old, often women
- Selling handmade goods, resold products, food, beauty, fashion
- 10–100 orders/month (growing)
- Revenue: $500–$5,000/month

**Behavior:**

- Sell via Instagram/WhatsApp DMs (no website)
- Manage inventory manually (or in their head)
- Track finances on paper, notes apps, or not at all
- Fear "proper business tools" (too complex, too expensive)

### Secondary Segment: Micro-Teams (2–5 people)

**Profile:**

- Small teams (often family members or friends)
- Slightly higher volume (100–500 orders/month)
- Revenue: $5,000–$20,000/month
- Need collaboration (multiple people taking orders)

**Behavior:**

- Still DM-first, but need coordination
- May have part-time help (delivery, production)
- Want to "feel professional" without big investments

### Adjacent Segments (Future)

- **Brick-and-mortar with DM sales**: Shops that also sell via Instagram
- **Event-based sellers**: Sell at markets/pop-ups, manage orders via DMs
- **Service providers**: Tutors, consultants, beauty services (appointment-based)

---

## 3. Value Proposition

### Job-to-be-Done (JTBD)

**When** I'm running a small business through Instagram/WhatsApp DMs,  
**I want to** know I'm not losing money or missing orders,  
**So I can** feel in control and focus on selling (not admin chaos).

### Core Value Delivered

| User need                     | Kyora solution                                      |
| ----------------------------- | --------------------------------------------------- |
| Track orders without chaos    | DM-native order capture, status tracking            |
| Know if I'm making money      | Automatic profit calculation, cash in hand summary  |
| Avoid running out of stock    | Low stock alerts, inventory tracking                |
| Feel professional             | Clean mobile UI, tax-ready records                  |
| Not need accounting knowledge | Plain language, automatic bookkeeping in background |

### Key Differentiators (vs competitors)

1. **DM-native**: Built for Instagram/WhatsApp workflows (not adapted from e-commerce)
2. **Arabic-first**: Native Arabic UI (not translated afterthought)
3. **Mobile-first**: Designed for phone use (not desktop-first with mobile view)
4. **Plain language**: "Profit", "Cash in hand" (not "EBITDA", "Accrual")
5. **Automatic**: Background bookkeeping (no manual categorization)

**See also**: `kyora/brand-key.instructions.md` (section 9: Discriminator)

---

## 4. Revenue Model

### Freemium SaaS

**Free Plan (Starter):**

- Core features: Orders, inventory, customers, basic accounting
- Limits: 50 orders/month, 50 products, 1 user
- Goal: Let users experience value before paying

**Paid Plans (Growth, Scale):**

- Unlock: Higher limits, multi-user collaboration, advanced analytics, integrations
- Pricing: Tiered based on order volume and team size
- Billing: Monthly/annual subscriptions (Stripe)

### Pricing Philosophy

- **Value-based**: Price scales with business size (fairness)
- **Transparent**: No hidden fees, clear limits
- **Accessible**: Free plan is genuinely useful (not crippled trial)
- **Plan-gated features**: Premium capabilities (analytics, multi-business, permissions) behind paid plans

**Implementation**: See `.github/instructions/billing.instructions.md` for plan limits and feature gates.

---

## 5. Business Model Canvas

### Key Activities

- **Product development**: Mobile-first UI, DM workflow automation
- **Automation**: Background bookkeeping, profit calculation, stock alerts
- **Customer support**: Onboarding help, Arabic-language support
- **Marketing**: Content marketing, social media, influencer partnerships

### Key Resources

- **Platform**: Mobile-first web app (React + Go API)
- **Team**: Product, engineering, customer success
- **Brand**: "Quiet expert" positioning, Arabic-first credibility
- **Data**: Transaction patterns, inventory trends (for future AI/ML features)

### Key Partners

- **Payment processors**: Stripe (international), local payment gateways (future)
- **Social platforms**: Instagram, WhatsApp (where customers operate)
- **Accounting firms**: Tax preparation partners (future referral program)

### Customer Relationships

- **Self-service**: Mobile app, onboarding flow, help docs
- **Proactive support**: Low stock alerts, "What to do next" prompts
- **Community**: User groups, social media engagement (future)

### Channels

- **Direct**: Kyora website, app download
- **Organic**: Instagram/WhatsApp seller communities, word-of-mouth
- **Paid**: Social media ads targeting DM sellers
- **Partnerships**: Influencer collaborations, seller associations

### Cost Structure

- **R&D**: Engineering, product design
- **Infrastructure**: Cloud hosting (AWS/GCP), database, CDN
- **Third-party services**: Stripe fees, email (Resend), storage (S3)
- **Marketing**: Ads, content creation
- **Support**: Customer success team

---

## 6. Competitive Advantage

### Sustainable Differentiators

| Advantage             | Why it's hard to copy                                         |
| --------------------- | ------------------------------------------------------------- |
| Arabic-first design   | Requires native Arabic UX thinking (not just translation)     |
| DM workflow depth     | Requires deep understanding of Instagram/WhatsApp selling     |
| Plain language AI     | Requires brand voice consistency across all features          |
| Mobile-first polish   | Most competitors are desktop-first with bolted-on mobile      |
| Automatic bookkeeping | Requires accounting domain expertise + simple UX (hard to do) |

### Network Effects (Future)

- **Data**: More users → better insights (e.g., "Best-selling products in your category")
- **Community**: More sellers → stronger word-of-mouth, user-generated content
- **Integrations**: More volume → leverage for partnerships (e.g., shipping, payments)

---

## 7. Unit Economics (High-Level)

**Customer Acquisition Cost (CAC):**

- Target: <$50 (organic + paid mix)
- Payback period: <6 months (for paid users)

**Lifetime Value (LTV):**

- Average subscription: $20–$50/month
- Retention target: 80%+ annually
- LTV goal: $500+ over 2 years

**Gross Margin:**

- SaaS standard: 70–80% (low infrastructure costs)

**Note**: These are directional targets, not hardcoded rules. Focus on product-market fit first, unit economics optimization second.

---

## 8. Go-to-Market Strategy

### Phase 1: Early Adopters (Current)

- **Target**: Power users in Saudi Arabia (Instagram sellers, 50–200 orders/month)
- **Channel**: Direct outreach, Instagram ads, influencer partnerships
- **Goal**: Product-market fit validation, feature refinement

### Phase 2: Growth (Next 12 months)

- **Target**: Broader Middle East (UAE, Egypt, Kuwait)
- **Channel**: Content marketing, SEO, paid social
- **Goal**: 10,000+ active users, strong retention metrics

### Phase 3: Scale (Future)

- **Target**: Regional dominance, adjacent segments (brick-and-mortar, services)
- **Channel**: Partnerships, API integrations, community-led growth
- **Goal**: Category leadership in Arabic social commerce tools

---

## 9. Product Strategy Implications

### Feature Prioritization Framework

**High Priority** (must-haves for ICP):

- DM order capture (fast, mobile-first)
- Inventory tracking (low stock alerts)
- Profit calculation (automatic, plain language)
- Customer management (notes, addresses)
- Arabic/RTL support (native, not translated)

**Medium Priority** (valuable for retention):

- Analytics (best sellers, top customers)
- Multi-user collaboration (for micro-teams)
- Export/reporting (tax preparation)
- Integrations (future: shipping, payments)

**Low Priority** (nice-to-haves, not differentiating):

- Advanced accounting (beyond profit/cash flow)
- E-commerce website builder (not ICP need)
- CRM features (marketing automation, email campaigns)

### Build vs Buy Decisions

| Capability        | Decision | Rationale                                         |
| ----------------- | -------- | ------------------------------------------------- |
| Auth/identity     | Buy      | Commodity (use proven libraries)                  |
| Payments          | Buy      | Stripe (reliable, feature-rich)                   |
| Email             | Buy      | Resend (transactional emails)                     |
| Storage           | Buy      | S3 (cheap, scalable)                              |
| Bookkeeping logic | Build    | Core differentiator (automatic, plain language)   |
| Order workflow    | Build    | DM-native flow is unique, not available off-shelf |
| Analytics         | Build    | Tailored to seller needs (best seller, low stock) |

---

## 10. Success Metrics

### North Star Metric

**"Active sellers who feel in control"**

- Proxy: Users who log in 3+ times/week
- Goal: 70%+ of users are active weekly

### Key Product Metrics

- **Activation**: Time to first order created (<5 minutes)
- **Engagement**: Orders created per week (target: 10+)
- **Retention**: 7-day, 30-day, 90-day retention (target: 80%, 60%, 50%)
- **Conversion**: Free → Paid (target: 10–15%)

### Business Metrics

- **MRR**: Monthly Recurring Revenue
- **Churn**: Monthly churn rate (target: <5%)
- **NPS**: Net Promoter Score (target: 50+)

---

## Decision-Making Guidance

When evaluating new features or pivots:

1. **ICP test**: Does this serve solo micro-entrepreneurs in the Middle East?
2. **JTBD test**: Does this help users "feel in control" without complexity?
3. **Differentiation test**: Does this strengthen our DM-first, Arabic-first, mobile-first advantage?
4. **Unit economics test**: Does this improve activation, engagement, or retention?

If the answer to 3+ questions is "No", deprioritize the feature.

---

## Related Documentation

- **Brand Key**: `.github/instructions/kyora/brand-key.instructions.md` — Why Kyora exists
- **Target Customer**: `.github/instructions/kyora/target-customer.instructions.md` — Who we serve
- **UX Strategy**: `.github/instructions/kyora/ux-strategy.instructions.md` — How to design for ICP
- **Billing**: `.github/instructions/billing.instructions.md` — Plan limits and feature gates
