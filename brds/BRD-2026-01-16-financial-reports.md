---
status: draft
owner: product
created_at: 2026-01-16
updated_at: 2026-01-16
stakeholders:
  - name: Product Lead
    role: Approver
  - name: Engineering Manager
    role: Implementation Lead
  - name: UX Designer
    role: Design Reviewer
areas:
  - portal-web
  - backend
kpis:
  - name: Financial Report Page Views
    definition: Number of times users view financial report pages per week
    baseline: 0 (new feature)
    target: >60% of active businesses visit reports weekly
  - name: Time to Financial Clarity
    definition: Time from login to viewing financial health summary
    baseline: N/A
    target: <10 seconds (2-3 taps from dashboard)
  - name: User Confidence Score
    definition: Survey-based confidence rating on financial health understanding
    baseline: N/A
    target: >4.2/5.0
  - name: Report Engagement Depth
    definition: % of users who view all 3 report types in a session
    baseline: N/A
    target: >40%
---

# BRD: Financial Reports — Business Health at a Glance

## 1) Problem (in plain language)

### What is the user trying to do?

Social media commerce entrepreneurs want to understand their business's financial health **without** studying accounting or understanding complex financial terms. They need answers to simple questions:

- "How much money do I actually have?"
- "Am I making profit or losing money?"
- "Can I safely take money out of the business?"
- "Where is my money going?"
- "What's my business worth?"

### What is painful/confusing today?

- **Accounting is intimidating**: Balance sheets, P&L statements, cash flow reports are designed for accountants, not sellers.
- **No clear picture**: Without Kyora's reports, users either don't know their financial status, or they maintain messy spreadsheets that are error-prone and time-consuming.
- **Fear of taking money out**: Users don't know how much they can safely withdraw without hurting the business.
- **Jargon overload**: Terms like "retained earnings", "COGS", "liabilities", and "equity" mean nothing to most users.

### Why does it matter for Kyora's customers?

- **Peace of mind**: This is Kyora's core promise—financial visibility without the stress.
- **Decision support**: Users need to know if they can afford new inventory, hire help, or invest in marketing.
- **Subscription driver**: Many users will subscribe to Kyora specifically for this financial clarity feature.
- **Mobile-first necessity**: Users check their financial status from their phones while on the go—waiting for deliveries, between customers, etc.

## 2) Customer & Context

### Primary user persona(s)

1. **Solo Seller (Primary)**: Manages everything alone. Needs quick, glanceable financial status.
2. **Side Hustler**: Has a day job, runs business in spare time. Needs clarity without time investment.
3. **Micro-Team Owner (2-5 people)**: Wants to share financial visibility with partners/family.

### Where does this happen?

- **On the go**: Between customer chats, during commute, while waiting.
- **End of day/week**: Quick check on business health.
- **Before big decisions**: "Can I afford to order more inventory?" "Should I take money out?"

### Device context (mobile-first)

- **Primary**: Mobile phones (iPhone, Android)
- **Secondary**: Desktop/tablet for deeper analysis
- **Network conditions**: Variable—reports should feel instant on 4G

### Language direction: RTL/Arabic-first requirements

- All content must be fully RTL-compatible
- Numbers should be displayed with proper formatting for Arabic locale
- Charts must support RTL text labels and legends
- Date formats should respect locale (Arabic months, day names)
- Currency formatting must respect business currency settings

## 3) Goals (what success looks like)

- **Goal 1: Instant Financial Clarity**: User understands business health within 5 seconds of viewing a report.
- **Goal 2: Jargon-Free Experience**: Zero accounting terminology visible to users—everything in plain language.
- **Goal 3: Mobile-First Excellence**: Reports are beautiful, fast, and fully functional on mobile.
- **Goal 4: Actionable Insights**: Each report suggests "what to do next" based on the numbers.
- **Goal 5: Confidence Building**: Users feel confident about their financial decisions after viewing reports.

## 4) Non-goals (explicitly out of scope)

- **Non-goal 1: Advanced Accounting Features**: No double-entry bookkeeping, ledgers, or journal entries.
- **Non-goal 2: Tax Filing Automation**: We won't file taxes or generate tax forms (future consideration).
- **Non-goal 3: Multi-Currency Reports**: Reports are single-currency (business currency) only.
- **Non-goal 4: Historical Comparison**: "Compare to last month/year" is deferred to a future iteration.
- **Non-goal 5: PDF Export**: Report export/printing is out of scope for this iteration.
- **Non-goal 6: Forecasting/Projections**: AI-based forecasts are future scope.

## 5) User journey (happy path)

### Journey A: Quick Health Check (most common)

1. User taps "Reports" in bottom navigation (or sidebar on desktop).
2. User lands on **Reports Hub** showing 3 report cards with key metrics visible immediately.
3. User sees their **Safe to Draw** amount prominently—feels reassured or alerted.
4. User taps "Business Health" card to see full snapshot.
5. User understands their financial position in plain language, closes app feeling informed.

### Journey B: Understanding Profit

1. User wonders "Am I actually making money?"
2. User navigates to Reports → taps "Profit & Earnings" card.
3. User sees clear breakdown: Money In → Product Costs → Running Costs → **What You Keep**.
4. User understands where money goes, identifies highest expense category.
5. User makes informed decision about expenses or pricing.

### Journey C: Cash Visibility

1. User asks "Where did my cash go?"
2. User navigates to Reports → taps "Cash Movement" card.
3. User sees simple flow: Cash Start → Cash In → Cash Out → **Cash Now**.
4. User identifies if business is cash-healthy or cash-strapped.
5. User decides whether to delay a purchase or take a withdrawal.

## 6) Edge cases & failure handling

### Case: No data yet (new business)

- Expected behavior: Show empty state with friendly guidance
- What the user sees: "Your financial picture will appear here once you start recording orders and expenses. Get started by creating your first order!"
- Action button: "Create First Order"

### Case: Negative safe-to-draw amount

- Expected behavior: Show warning, explain clearly
- What the user sees: "⚠️ Caution: Your business expenses currently exceed available profits. Consider reducing expenses or waiting before withdrawing money."
- Visual: Error/warning color treatment on the amount

### Case: Negative cash balance (approximated)

- Expected behavior: Show alert with explanation
- What the user sees: "💸 Your cash position is estimated to be negative. This usually means expenses have exceeded income. Review your recent expenses."
- Visual: Error color, link to expenses page

### Case: Loading large report data

- Expected behavior: Skeleton loaders with subtle animation
- What the user sees: Smooth skeleton that matches report layout
- Duration: Target <2s for initial load

### Case: Report fetch fails (network error)

- Expected behavior: Show retry option
- What the user sees: "Couldn't load your report. Check your connection and try again."
- Action: "Retry" button

### Case: Partial data (e.g., no expenses recorded)

- Expected behavior: Show available data, note what's missing
- What the user sees: Normal report with a note: "ℹ️ You haven't recorded any expenses yet. Your profit shown is before deducting running costs."
- Action: Link to add expenses

## 7) UX / IA (mobile-first)

### Information Architecture

```
Reports (Hub Page)
├── Business Health (Financial Position)
├── Profit & Earnings (P&L)
└── Cash Movement (Cash Flow)
```

### Navigation Entry Points

1. **Bottom Navigation** (mobile): "Reports" icon (new item, requires nav update)
2. **Sidebar** (desktop): "Reports" menu item under business menu
3. **Dashboard**: Optional quick-link card showing key metric

---

### Page 1: Reports Hub (`/business/:descriptor/reports`)

**Purpose**: Quick overview of all financial reports with key metrics visible at a glance—user should get value without drilling in.

**Primary action**: Tap a report card to view full details.

**Secondary actions**:
- Change report date (as-of date picker)
- Quick links to related sections (Expenses, Capital)

**Content (what must be shown)**:

**Hero Section: Safe to Draw**
- Large, prominent "Safe to Draw" amount
- Subtitle: "The amount you can safely take out"
- Color coding: Green (positive), Yellow (low), Red (negative)

**Report Cards Grid (3 cards)**:

1. **Business Health Card**
   - Title: "Business Health" / "صحة عملك"
   - Key metric: "What Your Business Is Worth: [TotalEquity]"
   - Secondary: "Cash: [CashOnHand] | Inventory: [InventoryValue]"
   - CTA: "View Details →"

2. **Profit & Earnings Card**
   - Title: "Profit & Earnings" / "أرباحك"
   - Key metric: "What You Keep: [NetProfit]"
   - Secondary: "Revenue: [Revenue] | Costs: [COGS + Expenses]"
   - Trend indicator: ↑/↓ based on profit (future: vs last period)
   - CTA: "View Details →"

3. **Cash Movement Card**
   - Title: "Cash Movement" / "حركة النقد"
   - Key metric: "Cash Now: [CashAtEnd]"
   - Secondary: "In: [TotalCashIn] | Out: [TotalCashOut]"
   - CTA: "View Details →"

**As-Of Date Picker**:
- Default: Today
- Format: "As of [date]" / "حتى تاريخ [date]"
- Mobile: Bottom sheet date picker
- Desktop: Dropdown date picker

**Empty state**:
- Illustration: Friendly graphic of a shop/store
- Headline: "Your financial picture starts here"
- Body: "Once you start recording orders and expenses, you'll see a complete view of your business finances."
- CTA: "Create Your First Order"

**Loading state**:
- 3 skeleton cards matching card layout
- Skeleton for Safe to Draw hero section
- Subtle pulse animation

**Error state**:
- Icon: Exclamation or refresh icon
- Headline: "Couldn't load your reports"
- Body: "We're having trouble connecting. Please check your internet and try again."
- CTA: "Retry"

**i18n keys needed**:
- `reports.hub.title`: "Reports" / "التقارير"
- `reports.hub.safe_to_draw`: "Safe to Draw" / "المبلغ الآمن للسحب"
- `reports.hub.safe_to_draw_subtitle`: "The amount you can safely take out" / "المبلغ الذي يمكنك سحبه بأمان"
- `reports.hub.as_of`: "As of {date}" / "حتى تاريخ {date}"
- `reports.cards.business_health`: "Business Health" / "صحة عملك"
- `reports.cards.profit_earnings`: "Profit & Earnings" / "أرباحك"
- `reports.cards.cash_movement`: "Cash Movement" / "حركة النقد"
- `reports.cards.view_details`: "View Details" / "عرض التفاصيل"
- `reports.metrics.business_worth`: "What Your Business Is Worth" / "قيمة عملك"
- `reports.metrics.what_you_keep`: "What You Keep" / "ما تحتفظ به"
- `reports.metrics.cash_now`: "Cash Now" / "النقد الحالي"
- `reports.empty.title`: "Your financial picture starts here" / "صورتك المالية تبدأ هنا"
- `reports.empty.body`: "Once you start recording orders and expenses, you'll see a complete view of your business finances." / "بمجرد أن تبدأ بتسجيل الطلبات والمصروفات، ستظهر لك صورة كاملة عن ماليات عملك."
- `reports.error.title`: "Couldn't load your reports" / "تعذر تحميل تقاريرك"
- `reports.error.body`: "We're having trouble connecting. Please check your internet and try again." / "نواجه مشكلة في الاتصال. يرجى التحقق من الإنترنت والمحاولة مجدداً."
- `reports.error.retry`: "Retry" / "إعادة المحاولة"

---

### Page 2: Business Health (`/business/:descriptor/reports/health`)

**Purpose**: Show what the business owns, owes, and is worth in plain language. (Balance Sheet translated for non-accountants)

**Primary action**: Understand financial position at a glance.

**Secondary actions**:
- Navigate to related sections (Assets, Expenses, Capital)
- Change as-of date

**Content (what must be shown)**:

**Header**:
- Title: "Business Health" / "صحة عملك"
- Subtitle: "A snapshot of what your business owns and is worth"
- Date: "As of [date]"
- Back navigation

**Hero Section: Business Worth (Total Equity)**:
- Large metric: "[TotalEquity]"
- Label: "What Your Business Is Worth" / "قيمة عملك"
- Explanation tooltip: "This is your business value: everything you own minus what you owe."

**Section 1: What You Own (Assets)**

Visual: Horizontal stacked bar or segmented progress bar

| Item | Amount | Visual |
|------|--------|--------|
| Cash on Hand | [CashOnHand] | Green segment |
| Inventory Value | [InventoryValue] | Blue segment |
| Fixed Assets | [FixedAssets] | Purple segment |
| **Total** | [TotalAssets] | Bold total |

Subsection "Cash on Hand" explainer:
- "Your available cash = (Money from sales + Owner investment) - (Expenses + Withdrawals + Asset purchases + Inventory)"

**Section 2: What You Owe (Liabilities)**

- Display: "[TotalLiabilities]" (currently always 0)
- Note: "Kyora doesn't track loans or credit yet. This will be available in a future update."

**Section 3: Owner's Stake (Equity)**

| Item | Amount | Meaning |
|------|--------|---------|
| Money You Put In | [OwnerInvestment] | Capital invested |
| Money You Took Out | [OwnerDraws] | Withdrawals |
| Profit Kept in Business | [RetainedEarnings] | Accumulated profit |
| **Your Business Value** | [TotalEquity] | Net worth |

**Insight Card**:
- If TotalEquity > 0: "✅ Your business is in a healthy position. You own more than you owe."
- If TotalEquity < 0: "⚠️ Your business value is negative. This means expenses and withdrawals exceed profits. Consider reducing costs or injecting more capital."
- If CashOnHand < 0: "💸 Your cash position is estimated as negative. Review recent expenses or add capital."

**Quick Links**:
- "View Assets" → Assets page
- "View Capital" → Capital page (investments/withdrawals)

**Empty state**: Same as hub, with context about needing data.

**Loading state**: Skeleton matching sections above.

**Error state**: Retry pattern.

**i18n keys needed**:
- `reports.health.title`: "Business Health" / "صحة عملك"
- `reports.health.subtitle`: "A snapshot of what your business owns and is worth" / "لمحة عما يملكه عملك وقيمته"
- `reports.health.business_worth`: "What Your Business Is Worth" / "قيمة عملك"
- `reports.health.what_you_own`: "What You Own" / "ما تملكه"
- `reports.health.what_you_owe`: "What You Owe" / "ما عليك"
- `reports.health.owners_stake`: "Owner's Stake" / "حصة المالك"
- `reports.health.cash_on_hand`: "Cash on Hand" / "النقد المتوفر"
- `reports.health.inventory_value`: "Inventory Value" / "قيمة المخزون"
- `reports.health.fixed_assets`: "Equipment & Assets" / "المعدات والأصول"
- `reports.health.total_assets`: "Total You Own" / "إجمالي ما تملكه"
- `reports.health.total_liabilities`: "Total You Owe" / "إجمالي ما عليك"
- `reports.health.money_put_in`: "Money You Put In" / "المال الذي وضعته"
- `reports.health.money_took_out`: "Money You Took Out" / "المال الذي سحبته"
- `reports.health.profit_kept`: "Profit Kept in Business" / "الأرباح المحتفظ بها"
- `reports.health.liabilities_note`: "Kyora doesn't track loans or credit yet." / "كيورا لا تتتبع القروض أو الائتمان بعد."
- `reports.health.insight_healthy`: "Your business is in a healthy position." / "عملك في وضع صحي."
- `reports.health.insight_negative`: "Your business value is negative." / "قيمة عملك سالبة."
- `reports.health.insight_cash_negative`: "Your cash position is estimated as negative." / "وضعك النقدي المقدر سالب."

---

### Page 3: Profit & Earnings (`/business/:descriptor/reports/profit`)

**Purpose**: Show where money comes from and where it goes, ending with what the user keeps. (P&L Statement in plain language)

**Primary action**: Understand profit breakdown.

**Secondary actions**:
- View expenses by category
- Navigate to related sections

**Content (what must be shown)**:

**Header**:
- Title: "Profit & Earnings" / "أرباحك"
- Subtitle: "Where your money comes from and where it goes"
- Date: "As of [date]"

**Hero Section: Net Profit**:
- Large metric: "[NetProfit]"
- Label: "What You Keep" / "ما تحتفظ به"
- Color: Green if positive, Red if negative
- Explanation: "This is your final profit after all costs and expenses."

**Visual: Profit Waterfall (Sankey-style or stepped bar)**

```
Revenue [█████████████████] [amount]
   ↓
- Product Costs (COGS) [███████] [amount]
   ↓
= Gross Profit [██████████] [amount]
   ↓
- Running Costs (Expenses) [████] [amount]
   ↓
= What You Keep [██████] [amount]
```

**Section 1: Money Coming In**:
- Revenue: [Revenue]
- Source label: "From your sales"

**Section 2: Product Costs**:
- COGS: [COGS]
- Explanation: "The cost of making or buying what you sold"

**Section 3: Gross Profit**:
- Gross Profit: [GrossProfit]
- Explanation: "What's left after product costs"
- Margin indicator: "Gross Margin: [GrossProfit/Revenue * 100]%"

**Section 4: Running Costs (Expenses by Category)**:

Chart: Horizontal bar chart or doughnut showing category breakdown

| Category | Amount | % of Total |
|----------|--------|------------|
| [Category 1] | [Amount] | [%] |
| [Category 2] | [Amount] | [%] |
| ... | ... | ... |
| **Total Running Costs** | [TotalExpenses] | 100% |

**Section 5: Final Profit**:
- Net Profit: [NetProfit]
- Profit Margin: "[NetProfit/Revenue * 100]%"

**Insight Card**:
- If NetProfit > 0: "🎉 You're profitable! For every [currency] of sales, you keep [margin]%."
- If NetProfit < 0: "📉 You're currently losing money. Your expenses exceed your gross profit by [amount]."
- Highest expense category: "💡 Your biggest expense is [category] at [amount] ([%])."

**Quick Links**:
- "View All Expenses" → Expenses page
- "View Orders" → Orders page

**Empty state**: Guide to create orders and record expenses.

**Loading state**: Skeleton waterfall + sections.

**Error state**: Standard retry.

**i18n keys needed**:
- `reports.profit.title`: "Profit & Earnings" / "أرباحك"
- `reports.profit.subtitle`: "Where your money comes from and where it goes" / "من أين يأتي مالك وأين يذهب"
- `reports.profit.what_you_keep`: "What You Keep" / "ما تحتفظ به"
- `reports.profit.money_in`: "Money Coming In" / "الأموال الواردة"
- `reports.profit.from_sales`: "From your sales" / "من مبيعاتك"
- `reports.profit.product_costs`: "Product Costs" / "تكلفة المنتجات"
- `reports.profit.product_costs_explanation`: "The cost of making or buying what you sold" / "تكلفة صنع أو شراء ما بعته"
- `reports.profit.gross_profit`: "Gross Profit" / "الربح الإجمالي"
- `reports.profit.gross_profit_explanation`: "What's left after product costs" / "ما يتبقى بعد تكلفة المنتجات"
- `reports.profit.gross_margin`: "Gross Margin" / "هامش الربح الإجمالي"
- `reports.profit.running_costs`: "Running Costs" / "مصاريف التشغيل"
- `reports.profit.total_running_costs`: "Total Running Costs" / "إجمالي مصاريف التشغيل"
- `reports.profit.final_profit`: "Final Profit" / "الربح النهائي"
- `reports.profit.profit_margin`: "Profit Margin" / "هامش الربح"
- `reports.profit.insight_profitable`: "You're profitable!" / "أنت تحقق ربحاً!"
- `reports.profit.insight_losing`: "You're currently losing money." / "أنت تخسر حالياً."
- `reports.profit.insight_biggest_expense`: "Your biggest expense is {category}." / "أكبر مصروفاتك هي {category}."

---

### Page 4: Cash Movement (`/business/:descriptor/reports/cashflow`)

**Purpose**: Show how cash flows through the business—where it comes from and where it goes. (Cash Flow Statement simplified)

**Primary action**: Understand cash in/out dynamics.

**Secondary actions**:
- View investments/withdrawals
- View expenses

**Content (what must be shown)**:

**Header**:
- Title: "Cash Movement" / "حركة النقد"
- Subtitle: "How cash flows through your business"
- Date: "As of [date]"

**Hero Section: Cash Now (Ending Balance)**:
- Large metric: "[CashAtEnd]"
- Label: "Cash Now" / "النقد الحالي"
- Color: Green if positive, Red if negative
- Explanation: "Your current cash runway"

**Visual: Cash Flow Diagram**

Simple flow visualization showing:
```
Start: [CashAtStart]
         ↓
    ┌─────────────────┐
    │  Money Coming In │
    │  ═══════════════ │
    │  + From Sales    │  [CashFromCustomers]
    │  + From Owner    │  [CashFromOwner]
    │  ───────────────│
    │  Total In        │  [TotalCashIn]
    └────────┬────────┘
             ↓
    ┌─────────────────┐
    │  Money Going Out │
    │  ═══════════════ │
    │  - Inventory     │  [InventoryPurchases]
    │  - Running Costs │  [OperatingExpenses]
    │  - Equipment     │  [BusinessInvestments]
    │  - To Owner      │  [OwnerDraws]
    │  ───────────────│
    │  Total Out       │  [TotalCashOut]
    └────────┬────────┘
             ↓
    Net Change: [NetCashFlow] (+ or -)
             ↓
End: Cash Now [CashAtEnd]
```

**Section 1: Cash Coming In**:

| Source | Amount |
|--------|--------|
| From Customers (Sales) | [CashFromCustomers] |
| From Owner (Investment) | [CashFromOwner] |
| **Total Cash In** | [TotalCashIn] |

**Section 2: Cash Going Out**:

| Destination | Amount |
|-------------|--------|
| Inventory Purchases | [InventoryPurchases] |
| Running Costs (Expenses) | [OperatingExpenses] |
| Equipment & Assets | [BusinessInvestments] |
| To Owner (Withdrawals) | [OwnerDraws] |
| **Total Cash Out** | [TotalCashOut] |

**Section 3: Net Change**:
- Net Cash Flow: [NetCashFlow]
- Indicator: Green arrow up if positive, Red arrow down if negative
- Explanation: "Your cash [increased/decreased] by [amount]"

**Insight Card**:
- If NetCashFlow > 0: "✅ Cash Healthy: More cash came in than went out."
- If NetCashFlow < 0: "⚠️ Cash Alert: You spent more than you received. Monitor your cash carefully."
- If OwnerDraws > CashFromCustomers: "💡 Tip: You withdrew more than your sales revenue. Consider aligning withdrawals with income."

**Quick Links**:
- "View Capital" → Capital page
- "View Expenses" → Expenses page

**Empty state**: Same pattern.

**Loading state**: Skeleton flow diagram.

**Error state**: Standard retry.

**i18n keys needed**:
- `reports.cashflow.title`: "Cash Movement" / "حركة النقد"
- `reports.cashflow.subtitle`: "How cash flows through your business" / "كيف ينتقل النقد في عملك"
- `reports.cashflow.cash_now`: "Cash Now" / "النقد الحالي"
- `reports.cashflow.cash_runway`: "Your current cash runway" / "المدى النقدي الحالي"
- `reports.cashflow.cash_start`: "Cash at Start" / "النقد في البداية"
- `reports.cashflow.cash_end`: "Cash at End" / "النقد في النهاية"
- `reports.cashflow.money_in`: "Money Coming In" / "الأموال الواردة"
- `reports.cashflow.money_out`: "Money Going Out" / "الأموال الصادرة"
- `reports.cashflow.from_customers`: "From Customers" / "من العملاء"
- `reports.cashflow.from_owner`: "From Owner" / "من المالك"
- `reports.cashflow.total_in`: "Total Cash In" / "إجمالي النقد الوارد"
- `reports.cashflow.inventory_purchases`: "Inventory Purchases" / "مشتريات المخزون"
- `reports.cashflow.running_costs`: "Running Costs" / "مصاريف التشغيل"
- `reports.cashflow.equipment_assets`: "Equipment & Assets" / "المعدات والأصول"
- `reports.cashflow.to_owner`: "To Owner" / "للمالك"
- `reports.cashflow.total_out`: "Total Cash Out" / "إجمالي النقد الصادر"
- `reports.cashflow.net_change`: "Net Change" / "صافي التغير"
- `reports.cashflow.cash_increased`: "Your cash increased by {amount}" / "ازداد نقدك بمقدار {amount}"
- `reports.cashflow.cash_decreased`: "Your cash decreased by {amount}" / "انخفض نقدك بمقدار {amount}"
- `reports.cashflow.insight_healthy`: "Cash Healthy: More cash came in than went out." / "النقد صحي: دخل نقد أكثر مما خرج."
- `reports.cashflow.insight_alert`: "Cash Alert: You spent more than you received." / "تنبيه نقدي: أنفقت أكثر مما استلمت."
- `reports.cashflow.insight_withdrawal_tip`: "Tip: You withdrew more than your sales revenue." / "نصيحة: سحبت أكثر من إيرادات مبيعاتك."

---

### Copy principles

- **Use plain language**: "What You Keep" not "Net Profit", "Running Costs" not "Operating Expenses", "Money Coming In" not "Revenue".
- **Prefer actionable CTAs**: "View Details", "View Expenses", "Retry"
- **Explain complex concepts inline**: Tooltip or small explainer text for concepts users might not understand.
- **Celebrate wins**: Positive language when metrics are good ("You're profitable!", "Cash Healthy")
- **Gentle warnings**: Non-scary language for negative metrics ("Caution", "Alert", not "DANGER" or "CRITICAL")

## 8) Functional requirements

### Core Functionality

- **FR-1**: Reports Hub page must load and display 3 report cards with key metrics within 3 seconds on 4G mobile.
- **FR-2**: Each report card must show the most important metric prominently without requiring drill-in.
- **FR-3**: Safe to Draw amount must be visible on the hub page as a hero metric.
- **FR-4**: All reports must support an "as-of" date parameter, defaulting to today.
- **FR-5**: All monetary values must be formatted according to business currency settings.
- **FR-6**: All date values must be formatted according to user locale (Arabic/English).
- **FR-7**: All report pages must support full RTL layout and Arabic text.

### Navigation

- **FR-8**: Reports must be accessible from bottom navigation on mobile (requires nav update).
- **FR-9**: Reports must be accessible from sidebar on desktop.
- **FR-10**: Each report page must have clear back navigation to the Reports Hub.
- **FR-11**: Report pages must have quick links to related sections (Expenses, Capital, Assets).

### Data Display

- **FR-12**: Business Health page must display: Total Assets breakdown, Total Liabilities (with note), Equity breakdown, and insight card.
- **FR-13**: Profit & Earnings page must display: Revenue, COGS, Gross Profit, Expenses by Category chart, Net Profit, and insight card.
- **FR-14**: Cash Movement page must display: Cash flow visualization, Cash In breakdown, Cash Out breakdown, Net Cash Flow, and insight card.
- **FR-15**: Each report must include contextual insight cards with actionable suggestions.

### States

- **FR-16**: All report pages must have a skeleton loading state that matches final layout.
- **FR-17**: All report pages must have an empty state for businesses with no data.
- **FR-18**: All report pages must have an error state with retry functionality.
- **FR-19**: Negative values must be clearly indicated with color coding (red) and appropriate iconography.

### Charts & Visualizations

- **FR-20**: Expenses by Category must be displayed as a horizontal bar chart or doughnut chart.
- **FR-21**: Assets breakdown must be displayed as a segmented/stacked bar visualization.
- **FR-22**: All charts must use the existing Chart.js infrastructure with RTL support.
- **FR-23**: All charts must be touch-friendly on mobile with proper tap targets.

### Performance

- **FR-24**: Report data must be cached client-side to avoid redundant API calls within a session.
- **FR-25**: Initial page load must show skeleton within 100ms, data within 3s.

## 9) Data & permissions

### Tenant scoping (workspace + business)

- All report data is scoped to a single business within a workspace.
- Reports are accessed via `/v1/businesses/:businessDescriptor/analytics/reports/*` endpoints.
- Business descriptor comes from URL route parameter.

### Roles (admin/member)

- Both admin and member roles can VIEW financial reports.
- Permission required: `role.ActionView` on `role.ResourceFinancialReports`.
- No write operations on reports (read-only).

### What must never leak across tenants

- Financial data must never be accessible across workspaces.
- Financial data must never be accessible across businesses within the same workspace.
- Backend enforces this via middleware chain: `EnforceAuthentication` → `EnforceValidActor` → `EnforceWorkspaceMembership` → `EnforceBusinessValidity`.

## 10) Analytics & KPIs

### Events to track

| Event Name | Properties | When Triggered |
|------------|------------|----------------|
| `reports.hub.viewed` | `businessId`, `hasData` | Reports Hub page loaded |
| `reports.health.viewed` | `businessId`, `asOf` | Business Health page loaded |
| `reports.profit.viewed` | `businessId`, `asOf` | Profit & Earnings page loaded |
| `reports.cashflow.viewed` | `businessId`, `asOf` | Cash Movement page loaded |
| `reports.date_changed` | `businessId`, `reportType`, `newDate` | User changes as-of date |
| `reports.quick_link.clicked` | `businessId`, `from`, `to` | User clicks a quick link |
| `reports.retry.clicked` | `businessId`, `reportType` | User clicks retry after error |

### KPI impact expectation

- **Weekly Active Report Viewers**: Expect 60%+ of active businesses to view reports at least once per week.
- **Report Depth**: Expect 40%+ of report viewers to view all 3 reports in a single session.
- **Time to Insight**: Expect users to understand their financial status within 10 seconds of landing on Reports Hub.
- **User Confidence**: Survey-based metric showing users feel confident about their financial understanding.

## 11) Rollout & risks

### Rollout plan

1. **Phase 1: Internal Testing** (1 week)
   - Deploy to staging environment
   - Internal team testing with test data
   - Fix critical bugs

2. **Phase 2: Soft Launch** (1 week)
   - Enable for 10% of businesses (feature flag)
   - Monitor performance and errors
   - Gather initial feedback

3. **Phase 3: Full Launch**
   - Enable for all businesses
   - Announce feature via in-app notification
   - Monitor KPIs

### Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Performance issues on large datasets | Medium | High | Optimize queries, add caching, test with large data |
| User confusion despite plain language | Medium | Medium | User testing before launch, iterate on copy |
| RTL/Arabic rendering issues | Low | Medium | Thorough RTL testing, Arabic native review |
| Backend calculation errors | Low | High | Comprehensive E2E tests, reconciliation checks |
| Mobile layout issues | Medium | Medium | Test on multiple devices, responsive design review |

### Mitigations

- Pre-launch performance testing with realistic data volumes
- Arabic-native user testing for copy and RTL layout
- Backend E2E tests for all financial calculations
- Device testing matrix (iPhone SE, Samsung, iPad, Desktop)

## 12) Open questions

- **Q1**: Should we add a "period comparison" feature (this month vs last month) in this iteration or defer?
  - Defer to future iteration to keep scope focused.

- **Q2**: Should reports be accessible to the `member` role or `admin` only?
  - Both can view. Financial transparency is valuable for teams.

- **Q3**: Should we add PDF/export functionality in this iteration?
  - Defer. Users can screenshot on mobile.

- **Q4**: How should we handle businesses with zero revenue but expenses recorded?
  - Show reports normally with appropriate messaging about no sales yet.

- **Q5**: Should the "as-of" date be editable, or always show today?
  - Editable for flexibility, but default to today.

## 13) Backend API enhancements (if needed)

The current backend API provides good foundation but may need these enhancements:

### Enhancement 1: Period-based financial reports (optional, deferred)

Currently, financial reports (`financial-position`, `profit-and-loss`, `cash-flow`) are inception-to-date (from beginning to `asOf`). For future period comparisons:

- Consider adding `from` parameter to enable period-based reports (e.g., last month only).
- This would enable "this month vs last month" comparisons.

**Status**: Defer to future iteration.

### Enhancement 2: Percentage calculations

Backend currently returns raw amounts. Portal needs to calculate percentages client-side:

- Gross Margin: `GrossProfit / Revenue * 100`
- Net Margin: `NetProfit / Revenue * 100`
- Expense category %: `CategoryAmount / TotalExpenses * 100`

**Recommendation**: Calculate client-side. Keep backend simple.

### Enhancement 3: Insight generation (future AI feature)

Future consideration: Backend could return generated insights based on financial data patterns.

**Status**: Out of scope for this iteration.

## 14) Acceptance criteria (definition of done)

### Functionality
- [ ] Reports Hub page loads and displays 3 report cards with key metrics
- [ ] Safe to Draw amount is prominently displayed on hub
- [ ] Business Health page displays full financial position data
- [ ] Profit & Earnings page displays P&L with expense breakdown chart
- [ ] Cash Movement page displays cash flow visualization
- [ ] As-of date picker works and updates all reports
- [ ] All quick links navigate to correct sections

### Mobile & RTL
- [ ] All pages are fully functional on mobile (iPhone SE, Samsung)
- [ ] All pages render correctly in RTL mode
- [ ] Arabic translations are complete and natural
- [ ] Touch targets are minimum 44x44px
- [ ] Charts render correctly in RTL

### States
- [ ] Skeleton loading states match final layout on all pages
- [ ] Empty states guide user to take action
- [ ] Error states allow retry
- [ ] Negative values are clearly indicated with red color

### Performance
- [ ] Initial load < 3 seconds on 4G mobile
- [ ] Skeleton appears within 100ms

### Quality
- [ ] No accounting jargon visible to users
- [ ] All monetary values formatted correctly
- [ ] All dates formatted according to locale
- [ ] Insight cards provide actionable guidance

### Navigation
- [ ] Reports accessible from mobile bottom nav
- [ ] Reports accessible from desktop sidebar
- [ ] Back navigation works on all report pages

### Testing
- [ ] Manual testing on iPhone, Android, Desktop
- [ ] RTL/Arabic testing by native speaker
- [ ] Empty state testing (new business with no data)
- [ ] Negative values testing
