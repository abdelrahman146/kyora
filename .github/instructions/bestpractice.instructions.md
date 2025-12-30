---
description: end-to-end guide for building client-side applications
applyTo: "portal-v2/**,storefront-web/**"
---

Below is a pragmatic, end-to-end guide for building **client-side SPAs** with **TanStack Query + Router + Form + Store**. It is opinionated toward: (1) type-safety, (2) predictable data flow, (3) minimal re-renders, and (4) operational reliability (errors, retries, offline-ish behavior).

---

## 1) Mental model: separate “where state lives”

A high-quality SPA is mostly about putting each category of state in the right place:

1. **Server state** (remote, cached, deduped, invalidated): **TanStack Query**
2. **URL state** (navigational state, filters, pagination, identity): **TanStack Router** search params + route params
3. **Form state** (input buffers, validation, submission): **TanStack Form**
4. **Client/UI state** (theme, panels, ephemeral UI toggles): **TanStack Store** (sparingly)

A useful rule: if the source of truth is the backend, **do not mirror it in a client store**—keep it in Query’s cache and derive UI from queries/mutations.

---

## 2) Baseline architecture and folder conventions

A scalable TanStack SPA usually benefits from **feature-first structure** with a small shared “platform” layer:

```
src/
  app/                 # app bootstrap, providers, router, query client
  routes/              # router tree + route components (file-based if desired)
  features/
    users/
      api/             # fetchers, DTOs, queryOptions factories
      components/
      routes/
      forms/
  shared/
    api/               # http client, auth token handling, interceptors
    ui/                # design system components
    utils/             # typed helpers, logging, guards
```

**Design principle:** keep fetchers + queryOptions near the feature that owns them, but reusable via exports.

---

## 3) TanStack Query best practices

### 3.1 Query keys: consistency or chaos

TanStack Query caches by **queryKey**, which must be a **top-level array** and should be JSON-serializable and unique to the data. ([tanstack.com][1])

**Do**

- Use **query key factories** per domain (users, projects, etc.).
- Include _all identity inputs_ (route params, filters, pagination) in the key.
- Prefer stable primitives; if you include objects, ensure stable shape and serialization.

**Don’t**

- Don’t hand-write keys ad hoc across components (drift is guaranteed).
- Don’t omit “filters” from the key and then wonder why caching is wrong.

**Common mistake**

- `['users', { page: 1, sort }]` where `{...}` is created inline and shape varies. This creates confusing cache fragmentation and invalidations.

---

### 3.2 Centralize query definitions with `queryOptions`

In v5, a best practice is to co-locate `queryKey` + `queryFn` (and default behaviors) using `queryOptions`, which improves reuse and TypeScript inference. ([tanstack.com][2])

**Pattern: queryOptions factory**

```ts
import { queryOptions } from "@tanstack/react-query";

export const userQueries = {
  list: (filters: { page: number; q?: string }) =>
    queryOptions({
      queryKey: ["users", "list", filters],
      queryFn: () => api.users.list(filters),
      staleTime: 60_000,
    }),

  detail: (id: string) =>
    queryOptions({
      queryKey: ["users", "detail", id],
      queryFn: () => api.users.byId(id),
      staleTime: 5 * 60_000,
    }),
};
```

**Do**

- Treat these factories as the canonical API for that domain’s reads.
- Use them for `useQuery`, prefetching, route loaders, and invalidations.

**Don’t**

- Don’t duplicate queryFn logic across components and loaders.

---

### 3.3 Staleness, refetching, and “important defaults”

Unconfigured queries can refetch more than you expect; Query’s “recommended” approach to reducing excessive refetching is setting `staleTime`, and you can tune refetch triggers like window focus, mount, and reconnect. ([tanstack.com][3])

**Do**

- Set `staleTime` intentionally for each query category:

  - “Rarely changes” (permissions, config): minutes/hours
  - “User-facing but tolerates staleness” (lists): tens of seconds to minutes
  - “Must be fresh” (trading-like): short staleTime + explicit refetch rules

**Don’t**

- Don’t globally set everything to `staleTime: 0` and then fight the refetch behavior everywhere.

**Pitfall**

- Disabling `refetchOnWindowFocus` globally without compensating: users return to stale UIs and lose trust. (If you disable it, you need another freshness strategy.) ([tanstack.com][4])

---

### 3.4 Invalidation strategy: targeted, not “nuke from orbit”

Use `queryClient.invalidateQueries` to mark queries stale and refetch when appropriate. ([tanstack.com][5])

**Do**

- Invalidate _just the affected slice_ (e.g., `['users', 'detail', id]` and relevant lists).
- Prefer invalidation over manual cache rewriting unless you truly need immediate UI changes.

**Don’t**

- Don’t invalidate `['users']` broadly if you can invalidate `['users','detail',id]` + one list key.

---

### 3.5 Pagination and “keep previous data” (v5 approach)

For paginated/laged queries, v5 uses `placeholderData` (or `keepPreviousData`) to avoid jarring loading state resets between pages. ([tanstack.com][6])

**Do**

- Use `placeholderData: keepPreviousData` for page transitions where UX continuity matters.

**Don’t**

- Don’t show a full-screen spinner for every small page change; it makes SPAs feel slow.

---

### 3.6 Optimistic updates: do it safely or don’t do it

TanStack Query supports optimistic updates via `onMutate` (cache update + rollback) or “via the UI” (local optimistic UI without cache writes). ([tanstack.com][7])

**Do**

- Cancel in-flight queries, snapshot prior cache, write optimistic cache, rollback on error, then invalidate.
- Keep the optimistic scope small (single entity or single list row).

**Don’t**

- Don’t optimistically update multiple unrelated queries unless you have a clear rollback story.

**Common pitfall**

- Forgetting rollback: users see “saved” states that later revert unpredictably.

---

## 4) TanStack Router best practices

### 4.1 Route-level data loading and Query integration

Router loaders can preload data in parallel and render quickly (often via suspense-style patterns). ([tanstack.com][8])
Query’s prefetching guide explicitly discusses router-level integration choices: block rendering until data is ready vs start prefetching without awaiting. ([tanstack.com][9])

**Recommended pattern (SPA):**

- **Critical data**: loader `await queryClient.ensureQueryData(...)`
- **Non-critical data**: loader starts `prefetchQuery(...)` without awaiting

**Do**

- Tie loader inputs (params/search) to query keys; if search params identify the data, they should identify the queryKey. ([tanstack.com][8])

**Don’t**

- Don’t fetch in both loader _and_ component with different keys—double requests and divergent caches.

---

### 4.2 Search params: validate, type, and treat as untrusted input

TanStack Router supports search-param tooling and emphasizes validation patterns for route-bound search. ([tanstack.com][10])

**Do**

- Validate search params at the route boundary.
- Normalize defaults (pageIndex=0, sort=…) consistently so URLs remain stable.

**Don’t**

- Don’t treat raw `URLSearchParams` as typed truth.

**Common mistakes**

- Putting complex filter objects into search params without canonical ordering/normalization → cache key chaos.
- Using search params to drive UI but forgetting they must also drive query keys.

---

### 4.3 Error, pending, and not-found handling (make it explicit)

Router supports `pendingComponent`, `errorComponent`, and `notFoundComponent`, and it’s strongly recommended to provide not-found handling at least at the root (or router default) to avoid undesirable defaults. ([tanstack.com][11])

**Do**

- Provide router-wide defaults (pending/error/not-found) for consistency.
- Use route-specific overrides where UX differs.

**Don’t**

- Don’t rely on framework defaults; you’ll ship ugly or misleading states.

---

### 4.4 Authentication & guards: `beforeLoad` + redirect

For authenticated routes, Router supports throwing a `redirect()` from `beforeLoad`, and if `beforeLoad` throws, children won’t attempt to load. ([tanstack.com][12])

**Do**

- Put auth gating in a parent layout route’s `beforeLoad`.
- Use a single “session source of truth” (often a Query `session` query) and have `beforeLoad` ensure it.

**Don’t**

- Don’t sprinkle auth checks inside components; it creates flicker and inconsistent navigation rules.

---

### 4.5 Code splitting: route-level lazy boundaries

TanStack Router supports code splitting of route components and related elements (pending/error/notFound/loader), including automatic splitting guidance. ([tanstack.com][13])

**Do**

- Split by route for large apps; keep shared layout chunks stable.
- Ensure error/pending UI is lightweight and doesn’t pull huge bundles.

**Don’t**

- Don’t code split so aggressively that navigation becomes death-by-100-chunks.

---

## 5) TanStack Form best practices

### 5.1 Validation: Standard Schema and adapters

TanStack Form supports validators following the **Standard Schema** specification (Zod/Valibot/etc.). ([tanstack.com][14])

**Do**

- Use schema validation as the single source of truth for form shape and constraints.
- Distinguish:

  - **Sync validation**: basic constraints (required, min/max)
  - **Async validation**: uniqueness checks, server-derived rules (debounced)

**Don’t**

- Don’t run async validation on every keystroke without debounce/cancellation.

**Pitfall**

- Assuming validation returns transformed values (docs explicitly warn it does not). ([tanstack.com][14])
  If you need transforms, do them in submission handling or a mapping layer.

---

### 5.2 Performance: selectors are not optional in practice

TanStack Form uses a store internally; docs recommend using `useStore` with a selector to avoid unnecessary re-renders—omitting selectors can re-render on any store change. ([tanstack.com][15])

**Do**

- Select only the state you render (`values.firstName`, `errorMap`, etc.). ([tanstack.com][15])

**Don’t**

- Don’t subscribe to the entire form store in every field component.

---

### 5.3 Submission pattern: Form + Mutation

**Recommended approach**

- `useMutation` owns the submission side effects (network, invalidation, optimistic updates).
- The form owns “input buffering,” validation state, and disabled/loading UX.

**Do**

- Disable submit while mutation is pending.
- On success:

  - close modal / navigate
  - invalidate the relevant queries (targeted)

- On error:

  - map server validation errors into field errors (if applicable)

**Don’t**

- Don’t put navigation and invalidation logic deep inside field components.

---

## 6) TanStack Store best practices (use it, but keep it small)

TanStack Store is a framework-agnostic store/signals implementation with adapters. ([tanstack.com][16])
It supports derived/computed state via `Derived`. ([tanstack.com][17])

**What belongs in Store**

- UI preferences (theme, density, sidebar collapsed)
- Client-only state that is not worth encoding in URL
- Cross-route ephemeral UI coordination (rare)

**What does not belong in Store**

- Backend entities already managed by Query
- Cached lists, details, and server truth (you will create sync bugs)

**Do**

- Use derived state for computed values rather than re-computing everywhere. ([tanstack.com][17])
- Subscribe with selectors so components only re-render for relevant changes (same principle as Form).

**Don’t**

- Don’t build a “god store” that becomes your shadow backend cache.

---

## 7) Cross-library design patterns that scale

### Pattern A: “Route owns critical data”

- Router loader ensures critical query data exists
- Component reads with `useQuery(userQueries.detail(id))` and renders immediately from cache

**Benefits**

- No “loading flash” after navigation
- Predictable waterfall control

---

### Pattern B: “URL is the state machine for lists”

For list pages:

- Filters/pagination live in **search params**
- Query key includes those params
- UI updates via `router.navigate({ search: ... })`

**Benefits**

- Deep links work
- Back/forward works
- Cache aligns with navigation

---

### Pattern C: “Mutation consequences are centralized”

Define mutation helpers that:

- perform the mutation
- update or invalidate affected queries (targeted)
- optionally do optimistic updates safely

**Benefits**

- Side effects don’t sprawl across components

---

### Pattern D: “Forms are leaf nodes; mutations are roots”

- Form components stay dumb: validate + submit values
- Feature-level hook owns mutation + navigation + invalidation policy

**Benefits**

- Easy to test
- Easy to reuse forms (modal vs page)

---

## 8) Common mistakes and pitfalls (checklist)

### Query pitfalls

- Inconsistent query keys across features → duplicated cache entries and invalidations (solve with queryOptions factories). ([tanstack.com][2])
- Over-invalidation (“invalidate everything”) → performance collapse and UI flicker. ([tanstack.com][5])
- Forgetting `staleTime` tuning → excessive background refetches. ([tanstack.com][3])
- Misusing optimistic updates without rollback. ([tanstack.com][7])

### Router pitfalls

- Unvalidated search params → type bugs + cache mismatches. ([tanstack.com][10])
- No not-found strategy → users hit ugly defaults. ([tanstack.com][11])
- Auth checks inside components → flicker and leaky access control (use `beforeLoad` + redirect). ([tanstack.com][12])

### Form pitfalls

- Subscribing to the whole form store everywhere → re-render storms (use selectors). ([tanstack.com][15])
- Assuming validation transforms values → mismatches at submit time. ([tanstack.com][14])

### Store pitfalls

- Treating Store as a second Query cache → eventual consistency bugs and wasted effort.

---

## 9) Optimized approaches (what “good” looks like)

**Performance**

- Route-level prefetch for critical screens (loader + ensureQueryData). ([tanstack.com][9])
- Paginated lists use placeholderData/keepPreviousData for smooth transitions. ([tanstack.com][18])
- Form and Store subscriptions use selectors to minimize re-renders. ([tanstack.com][15])
- Router code splitting by route, with lightweight pending/error components. ([tanstack.com][13])

**Reliability**

- Explicit error boundaries and reset strategy when using suspense-style query loading. ([tanstack.com][19])
- Auth gating in `beforeLoad` with redirects. ([tanstack.com][12])

**Maintainability**

- QueryOptions factories + typed route search params become your “contract surface”
- Mutations centralize invalidation and optimistic logic

---

## 10) “Do and Don’t” quick reference

### Do

- Do model server state with Query, URL state with Router, forms with Form, and only true UI state with Store.
- Do centralize query definitions via `queryOptions`. ([tanstack.com][2])
- Do validate search params and align them with query keys. ([tanstack.com][10])
- Do provide router-wide pending/error/not-found defaults. ([tanstack.com][11])
- Do use selectors for Form store subscriptions. ([tanstack.com][15])

### Don’t

- Don’t mirror backend entities into Store “because it’s convenient.”
- Don’t invalidate huge swaths of cache when you can invalidate precisely. ([tanstack.com][5])
- Don’t treat URL/search params as trusted or typed without validation. ([tanstack.com][8])
- Don’t optimistically update without a rollback plan. ([tanstack.com][7])

---

[1]: https://tanstack.com/query/v5/docs/react/guides/query-keys?utm_source=chatgpt.com "Query Keys | TanStack Query React Docs"
[2]: https://tanstack.com/query/v5/docs/react/guides/query-options?utm_source=chatgpt.com "Query Options | TanStack Query React Docs"
[3]: https://tanstack.com/query/v5/docs/react/guides/important-defaults?utm_source=chatgpt.com "Important Defaults | TanStack Query React Docs"
[4]: https://tanstack.com/query/v5/docs/react/guides/window-focus-refetching?utm_source=chatgpt.com "Window Focus Refetching | TanStack Query React Docs"
[5]: https://tanstack.com/query/v5/docs/react/guides/query-invalidation?utm_source=chatgpt.com "Query Invalidation | TanStack Query React Docs"
[6]: https://tanstack.com/query/v5/docs/react/guides/migrating-to-v5?utm_source=chatgpt.com "Migrating to TanStack Query v5"
[7]: https://tanstack.com/query/v5/docs/react/guides/optimistic-updates?utm_source=chatgpt.com "Optimistic Updates | TanStack Query React Docs"
[8]: https://tanstack.com/router/v1/docs/framework/react/guide/data-loading?utm_source=chatgpt.com "Data Loading | TanStack Router React Docs"
[9]: https://tanstack.com/query/v5/docs/react/guides/prefetching?utm_source=chatgpt.com "Prefetching & Router Integration"
[10]: https://tanstack.com/router/v1/docs/framework/react/guide/search-params?utm_source=chatgpt.com "Search Params | TanStack Router React Docs"
[11]: https://tanstack.com/router/v1/docs/framework/react/guide/not-found-errors?utm_source=chatgpt.com "Not Found Errors | TanStack Router React Docs"
[12]: https://tanstack.com/router/v1/docs/framework/react/guide/authenticated-routes?utm_source=chatgpt.com "Authenticated Routes | TanStack Router React Docs"
[13]: https://tanstack.com/router/v1/docs/framework/react/guide/code-splitting?utm_source=chatgpt.com "Code Splitting | TanStack Router React Docs"
[14]: https://tanstack.com/form/v1/docs/framework/react/guides/validation?utm_source=chatgpt.com "Form and Field Validation | TanStack Form React Docs"
[15]: https://tanstack.com/form/v1/docs/framework/react/guides/basic-concepts?utm_source=chatgpt.com "Basic Concepts and Terminology | TanStack Form React ..."
[16]: https://tanstack.com/store/latest/docs?utm_source=chatgpt.com "Overview | TanStack Store Docs"
[17]: https://tanstack.com/store/latest/docs/reference/classes/Derived?utm_source=chatgpt.com "Derived | TanStack Store Docs"
[18]: https://tanstack.com/query/v5/docs/react/guides/paginated-queries?utm_source=chatgpt.com "Paginated / Lagged Queries | TanStack Query React Docs"
[19]: https://tanstack.com/query/v5/docs/react/guides/suspense?utm_source=chatgpt.com "Suspense | TanStack Query React Docs"

## 1) Staying DRY without overengineering (but still optimized)

The practical target is **“one place to change things”** for conventions, without building an abstraction layer that becomes its own framework.

### 1.1 A simple rule set that prevents both duplication and overdesign

**Rule A: DRY the _policy_, not the implementation.**
Centralize _decisions_ (key naming, cache/invalidation rules, route conventions, translation conventions, theme token mapping). Allow small amounts of repetition in leaf components.

**Rule B: Promote patterns only after the 2nd or 3rd copy.**
If you have only one usage, keep it local. If you see the same shape emerging repeatedly, extract it.

**Rule C: Prefer “configuration objects” over “wrapper functions”.**
For TanStack Query, this means using `queryOptions` (co-located key + fn + defaults) rather than building a custom `useApiQuery` abstraction. ([tanstack.com][1])

### 1.2 The minimal “platform layer” that keeps projects consistent

Keep these as small, stable modules:

- `app/queryClient.ts` – QueryClient defaults (staleTime policies, retry policy)
- `app/router.tsx` – router setup, root pending/error/notFound conventions (Router supports defaults like `defaultPendingMs`) ([tanstack.com][2])
- `app/i18n.ts` – i18next init, language detection strategy, namespace loading helpers (init once) ([i18next.com][3])
- `app/theme.ts` – daisyUI theme application + token resolution for charts
- `shared/http.ts` – fetch client and headers (including Accept-Language when needed)

Everything else should be **feature-owned** (queries/mutations/routes/components/translations per feature).

### 1.3 “Feature contract” pattern (high DRY, low complexity)

For each feature, export exactly these (names vary):

- `api.ts` (fetchers)
- `queries.ts` (queryOptions factories)
- `mutations.ts` (mutation hooks + invalidation policy)
- `route.ts` (route search schema defaults, mapping URL → filters)
- `i18n/` (namespace resource files)

This avoids:

- giant global “api layer”
- query keys scattered across UI
- ad hoc invalidation logic everywhere

#### Example: queryOptions + typed invalidation in one place

Using `queryOptions` keeps the key and function co-located and reusable across components and route loaders. ([tanstack.com][1])

```ts
// features/users/queries.ts
import { queryOptions } from "@tanstack/react-query";

export const usersQ = {
  list: (filters: { page: number; q?: string }) =>
    queryOptions({
      queryKey: ["users", "list", filters], // query keys must be top-level arrays and uniquely identify data :contentReference[oaicite:4]{index=4}
      queryFn: () => api.users.list(filters),
      staleTime: 60_000,
    }),

  detail: (id: string) =>
    queryOptions({
      queryKey: ["users", "detail", id],
      queryFn: () => api.users.byId(id),
      staleTime: 5 * 60_000,
    }),
};
```

### 1.4 Optimization without cleverness (high ROI, low cognitive load)

**Do**

- Use **URL search params** as the canonical state for list filters/pagination.
- Use **route loaders** for critical prefetch (avoid “loading flash” after navigation).
- Keep **stable references**: memoize chart options, large derived objects, and callbacks.
- Use **selectors** for store subscriptions (both TanStack Store and TanStack Form benefit from selecting minimal slices).

**Don’t**

- Don’t introduce a generic “Repository” or “Service” layer unless you truly have cross-cutting needs (auditing, offline queueing, multi-tenant routing).
- Don’t build a “mega store” that mirrors server entities—Query already _is_ your server-state store.

---

## 2) i18next / react-i18next best practices (SPA + TanStack Router/Query)

### 2.1 Initialization and Suspense (avoid “double init” and flicker)

- **Initialize i18next once**. The i18next docs explicitly advise not to call init multiple times; use `changeLanguage` for switching. ([i18next.com][3])
- In React, `useTranslation` can **trigger Suspense while translations load** (unless you disable it with `useSuspense: false`). ([react.i18next.com][4])

**Recommended default**

- Keep `useSuspense: true` (clean UX), and ensure your router/root has a consistent pending UI.
- If you _must_ avoid Suspense fallback during language change, you need a deliberate “preload then switch” flow (see 2.4).

### 2.2 Namespace strategy that stays DRY

Use **namespaces per feature** (and optionally a shared `common`):

- `common` (buttons, nav, generic toasts)
- `users`, `billing`, `reports` etc.

This gives you:

- smaller bundles
- simpler ownership
- fewer merge conflicts in translation files

### 2.3 Key design rules (keeps translation maintainable)

**Do**

- Use keys that encode meaning and scope: `users.list.title`, `users.form.email.label`
- Prefer full sentences in resources rather than concatenating fragments in code.

**Don’t**

- Don’t overuse interpolation. i18next explicitly recommends using interpolation sparingly; avoid it when the translated content can be fully self-contained. ([i18next.com][5])

### 2.4 Route-driven language (best UX) + preloading namespaces

A clean SPA approach is to treat language as **URL state**:

- `/:lang/...` (path param), or
- `?lang=xx` (search param)

**Flow**

1. Router validates `lang`
2. In a parent route `beforeLoad`/loader:

   - preload required namespaces for that route
   - then call `i18n.changeLanguage(lang)` (or do it first if you’re OK with Suspense fallback)

**Why preload?**
Because `useTranslation` can suspend when not ready. ([react.i18next.com][4])

### 2.5 Query cache correctness in multilingual apps

If your API responses are language-dependent (common when the server localizes strings):

- Include language in:

  - the **request header** (e.g., `Accept-Language`)
  - and the **queryKey** (or otherwise segment caches by language)

If you do not, you risk serving cached data in the wrong language after a switch.

### 2.6 Practical “Do / Don’t” list for i18n

**Do**

- Keep i18n init in `app/i18n.ts` (single instance). ([i18next.com][3])
- Use namespaces per feature.
- Keep interpolation minimal and structured. ([i18next.com][5])
- Preload namespaces for a route before switching language when you want no flicker. ([react.i18next.com][4])

**Don’t**

- Don’t store translated strings in TanStack Store or Query cache—store keys + values and translate at render time (except for server-localized payloads).

---

## 3) Branded, themed Chart.js in a TanStack + daisyUI SPA

### 3.1 The key principle: derive Chart.js styling from **daisyUI semantic tokens**

daisyUI explicitly encourages using semantic color utilities like `bg-primary` instead of hard-coded color utilities, and these semantic colors are driven by theme CSS variables. ([daisyUI][6])

So the DRY approach is:

- daisyUI theme is the single “design system”
- Chart.js colors/fonts are **computed from the active theme**, not duplicated hex codes

### 3.2 A robust token resolver (future-proof, avoids coupling to internal CSS variable names)

Instead of relying on exact variable names, you can resolve colors by rendering a hidden element with a daisyUI class and reading computed styles.

```ts
// app/themeTokens.ts
export function resolveBgColorFromClass(className: string): string {
  const el = document.createElement("div");
  el.className = className;
  el.style.position = "absolute";
  el.style.left = "-9999px";
  document.body.appendChild(el);

  const color = getComputedStyle(el).backgroundColor;
  document.body.removeChild(el);
  return color;
}

export function getChartTokens() {
  return {
    primary: resolveBgColorFromClass("bg-primary"),
    secondary: resolveBgColorFromClass("bg-secondary"),
    accent: resolveBgColorFromClass("bg-accent"),
    base: resolveBgColorFromClass("bg-base-100"),
    text: (() => {
      const el = document.createElement("div");
      el.className = "text-base-content";
      el.style.position = "absolute";
      el.style.left = "-9999px";
      document.body.appendChild(el);
      const color = getComputedStyle(el).color;
      document.body.removeChild(el);
      return color;
    })(),
  };
}
```

This stays DRY because:

- brand/theme changes happen in daisyUI theme config
- charts automatically follow without additional work

### 3.3 Applying tokens to Chart.js (colors, fonts, background)

Chart.js supports configuring background/border colors and text colors. ([chartjs.org][7])
It also supports global font defaults via `Chart.defaults.font`. ([chartjs.org][8])
Plugins are the most efficient customization mechanism when you need custom behavior (like canvas background). ([chartjs.org][9])

**Base options builder**

```ts
// app/chartTheme.ts
import type { ChartOptions } from "chart.js";

export function buildThemedOptions(
  tokens: ReturnType<typeof import("./themeTokens").getChartTokens>
): ChartOptions<"line"> {
  return {
    responsive: true,
    maintainAspectRatio: false,
    color: tokens.text, // default text color :contentReference[oaicite:16]{index=16}
    plugins: {
      legend: { labels: { color: tokens.text } },
      tooltip: { titleColor: tokens.text, bodyColor: tokens.text },
    },
    scales: {
      x: { ticks: { color: tokens.text }, grid: { color: "rgba(0,0,0,0.08)" } },
      y: { ticks: { color: tokens.text }, grid: { color: "rgba(0,0,0,0.08)" } },
    },
  };
}
```

**Canvas background plugin** (so the chart matches `bg-base-100`)

```ts
// app/chartPlugins.ts
import type { Plugin } from "chart.js";

export function canvasBackgroundPlugin(backgroundColor: string): Plugin {
  return {
    id: "canvasBackground",
    beforeDraw(chart) {
      const { ctx, chartArea } = chart;
      if (!chartArea) return;
      ctx.save();
      ctx.fillStyle = backgroundColor;
      ctx.fillRect(
        chartArea.left,
        chartArea.top,
        chartArea.right - chartArea.left,
        chartArea.bottom - chartArea.top
      );
      ctx.restore();
    },
  };
}
```

### 3.4 Performance guidelines for Chart.js in SPAs

Chart.js has explicit performance guidance:

- Use **decimation** for large line datasets (best results, reduces memory and draw cost). ([chartjs.org][10])
- Disable parsing when you can supply data in the expected internal format (`parsing: false`)—Chart.js documents that parsing can be disabled and comes with data format requirements. ([chartjs.org][11])
- Disable animations when rendering large datasets or frequent updates (animations can be disabled by setting the animation node to `false`). ([chartjs.org][12])

**Example performance switches**

```ts
const options = {
  animation: false, // disable animations :contentReference[oaicite:20]{index=20}
  parsing: false, // only if your data is already in the required format :contentReference[oaicite:21]{index=21}
  plugins: {
    decimation: { enabled: true, algorithm: "min-max" }, // for line charts :contentReference[oaicite:22]{index=22}
  },
};
```

### 3.5 Theme selection with TanStack Store (simple, DRY, no overreach)

TanStack Store is designed for lightweight client state, and `useStore` supports selectors. ([tanstack.com][13])

Use it for:

- `theme` (daisyUI `data-theme`)
- optional `brand` variant if you white-label

```ts
// app/themeStore.ts
import { Store } from "@tanstack/store";

export const themeStore = new Store<{ theme: string }>({ theme: "corporate" });
```

```ts
// app/applyTheme.ts
export function applyDaisyTheme(theme: string) {
  document.documentElement.setAttribute("data-theme", theme); // daisyUI theme mechanism :contentReference[oaicite:24]{index=24}
}
```

Then in React:

- subscribe to theme
- apply `data-theme`
- rebuild chart tokens/options

This keeps it DRY and avoids an overengineered “theme system” in TypeScript.

### 3.6 Branding guidelines for Chart.js + daisyUI

**Do**

- Use **semantic series roles**:

  - primary series = `primary`
  - comparison series = `secondary`
  - highlight/outlier = `accent`

- Keep gridlines low-contrast and let the data carry emphasis.
- Limit categorical palettes (e.g., 6–8 distinct colors max); use tooltips/labels for clarity.

**Don’t**

- Don’t hard-code hex values in chart components; resolve from theme.
- Don’t mix “brand palette” and “UI palette” unless you have a token mapping layer.

---

## 4) Putting it together (DRY, optimized, not overbuilt)

A good “low ceremony” integration point is a `ChartFrame` component that:

- reads `theme` from TanStack Store (selector)
- applies daisyUI theme attribute
- resolves tokens once per theme change
- memoizes Chart.js options

Separately, keep i18n route-level behavior in a parent route loader and keep Query’s cache language-correct by keying on language when necessary.

---

[1]: https://tanstack.com/query/v5/docs/react/guides/query-options?utm_source=chatgpt.com "Query Options | TanStack Query React Docs"
[2]: https://tanstack.com/router/v1/docs/framework/react/api/router/RouterOptionsType?utm_source=chatgpt.com "RouterOptions | TanStack Router React Docs"
[3]: https://www.i18next.com/overview/api?utm_source=chatgpt.com "API | i18next documentation"
[4]: https://react.i18next.com/latest/usetranslation-hook?utm_source=chatgpt.com "useTranslation (hook)"
[5]: https://www.i18next.com/principles/best-practices?utm_source=chatgpt.com "Best Practices"
[6]: https://daisyui.com/docs/colors/?lang=en&utm_source=chatgpt.com "Colors — daisyUI Tailwind CSS Component UI Library"
[7]: https://www.chartjs.org/docs/latest/general/colors.html?utm_source=chatgpt.com "Colors"
[8]: https://www.chartjs.org/docs/latest/general/fonts.html?utm_source=chatgpt.com "Fonts"
[9]: https://www.chartjs.org/docs/latest/developers/plugins.html?utm_source=chatgpt.com "Plugins"
[10]: https://www.chartjs.org/docs/latest/general/performance.html?utm_source=chatgpt.com "Performance"
[11]: https://www.chartjs.org/docs/latest/api/interfaces/ParsingOptions.html?utm_source=chatgpt.com "Interface: ParsingOptions"
[12]: https://www.chartjs.org/docs/latest/configuration/animations.html?utm_source=chatgpt.com "Animations"
[13]: https://tanstack.com/store/latest/docs/quick-start?utm_source=chatgpt.com "Quick Start | TanStack Store Docs"
