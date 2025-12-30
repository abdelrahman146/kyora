# Plan: Align portal-v2 with TanStack Best Practices

Transform portal-v2 into a production-grade TanStack application by systematically adopting all recommended patterns from the best practice guide. This phased approach prioritizes foundational patterns first, then layers advanced optimizations, ensuring each phase delivers measurable improvements to code quality, performance, and maintainability.

## Phase 1: Core Query Architecture (Foundation)

1. **Create queryOptions factories for all domains** in `api/` directory—replace scattered query definitions with co-located `queryOptions` exports (e.g., `customerQueries.list()`, `businessQueries.details()`) that bundle queryKey + queryFn + staleTime + enabled logic for type-safe reuse across components and route loaders

2. **Refactor all custom hooks** (`api/customer.ts`, `api/business.ts`, `api/user.ts`) to consume queryOptions factories—replace inline `useQuery({ queryKey: ..., queryFn: ... })` calls with `useQuery(customerQueries.list(...))` pattern

3. **Add queryOptions to route loaders** for critical data paths—update `routes/business/$businessDescriptor.tsx`, `routes/business/$businessDescriptor/customers/index.tsx` and other major routes to use `queryClient.ensureQueryData(customerQueries.list(...))` in loader functions

## Phase 2: Pagination & Data Loading UX

4. **Eliminate manual data fetching** in `routes/business/$businessDescriptor/customers/index.tsx`—remove `useState<Customer[]>`, manual `fetchCustomers` callback, and replace with proper `useQuery` consumption of `customerQueries.list()` with pagination params from search state

5. **Implement keepPreviousData pattern** across all paginated queries—add `placeholderData: keepPreviousData` to customer lists, order lists, and any other paginated data to eliminate loading flicker during page transitions

6. **Add route-level prefetching** for all major feature routes—create loaders that prefetch non-critical data without awaiting (e.g., metadata, filter options) while ensuring critical data with `ensureQueryData`

## Phase 3: Form Performance & Validation

7. **Refactor all TanStack Forms to use selectors**—update `components/organisms/customers/AddCustomerSheet.tsx`, business settings forms, and onboarding forms to replace `form.useStore()` with `form.useStore((state) => state.values.fieldName)` for precise subscriptions

8. **Implement async validation** for uniqueness checks—add async validators to customer email fields, business slugs, and other unique identifiers using Zod `.refine()` with API calls

9. **Standardize form-mutation separation**—ensure all form submission handlers delegate to mutations without side effects in forms, moving success callbacks, invalidations, and navigation to mutation `onSuccess` handlers

## Phase 4: Optimistic Updates

10. **Implement safe optimistic updates for create operations**—add `onMutate` with cache snapshot, `onError` rollback, and `onSettled` cleanup to customer creation, order creation, and other high-frequency mutations

11. **Implement safe optimistic updates for update operations**—add optimistic cache updates to edit customer, update business settings, and status change mutations with proper rollback logic

12. **Implement safe optimistic updates for delete operations**—add immediate cache removal with snapshot/restore for customer deletion, order cancellation, and similar destructive actions

## Phase 5: Router & Error Boundaries

13. **Add global error and notFound components**—create `components/templates/ErrorBoundary.tsx` and `components/templates/NotFound.tsx`, configure in `router.tsx` as router-wide defaults

14. **Add route-specific error components** where UX diverges—implement custom `errorComponent` in critical routes like `routes/business/$businessDescriptor.tsx` to provide contextual recovery options

15. **Ensure search param consistency** between URL state and query keys—audit all routes to verify that search params drive query keys completely (no hidden state), fixing `routes/business/$businessDescriptor/customers/index.tsx` filter discrepancies

## Phase 6: Store Optimizations

16. **Refactor all store subscriptions to use selectors**—update `hooks/useAuth.ts`, business store consumers, and metadata store usage to subscribe with `useStore(authStore, (state) => state.isAuthenticated)` pattern

17. **Remove metadata store duplication**—eliminate `stores/metadataStore.ts` and rely solely on TanStack Query's `gcTime` for caching countries/currencies, removing redundant TTL logic

18. **Add Derived stores for computed values**—convert computed functions in `stores/businessStore.ts` to proper `Derived` instances for automatic memoization

## Phase 7: i18n Enhancements

19. **Restructure translation namespaces by feature**—split `i18n/en/translation.ts` and `i18n/ar/translation.ts` into feature-based files (customers.ts, orders.ts, billing.ts, analytics.ts) and shared common.ts

20. **Implement lazy-loading for translation namespaces**—configure i18next backend plugin to load feature translations on demand when routes activate

21. **Include language in query keys** for localized API responses—audit API endpoints that return localized data (business categories, product descriptions) and add `lng` to their query keys, ensuring cache invalidation on language switch

22. **Add Accept-Language header to API client**—update `api/client.ts` to inject current language into all requests via `beforeRequest` hook

## Phase 8: Advanced Invalidation

23. **Refine business-scoped invalidation** in `lib/queryInvalidation.ts`—replace broad `invalidateBusinessScopedQueries` with targeted invalidations based on actual data dependencies (e.g., only invalidate orders/customers when switching businesses, not metadata)

24. **Move invalidation logic into mutation definitions**—centralize all `queryClient.invalidateQueries()` calls in mutation `onSuccess` handlers in `api/*.ts` files, removing ad-hoc invalidations from component callbacks

25. **Implement selective invalidation strategies**—use `queryClient.invalidateQueries({ queryKey: customerQueries.list._def, exact: false })` for list refreshes while preserving detail queries, add fine-grained invalidation for related data (e.g., updating customer should only invalidate that customer's orders)

## Further Considerations

### 1. Optimistic Update Scope

Implement optimistic updates incrementally, starting with customer creation (high-value, low-complexity) before tackling order mutations (complex state transitions). Accept that some mutations (bulk operations, complex workflows) may remain pessimistic for safety.

### 2. Translation Strategy Trade-offs

Lazy-loading namespaces adds complexity but significantly reduces initial bundle size. For portal-v2's current scope (~4 namespaces), upfront loading is acceptable, but plan to implement lazy loading before adding more features (aim to switch when exceeding ~10 namespaces).

### 4. Migration Strategy

Each phase builds on previous foundations—don't skip ahead. Phases 1-3 are mandatory for production readiness. Phases 4-8 are high-value quality improvements. Phases 9-10 are polish and long-term maintainability investments.

### 5. Performance Monitoring

After implementing form selectors (Phase 3) and store selectors (Phase 6), use React DevTools Profiler to measure render reduction. Target: 60%+ reduction in unnecessary re-renders for form-heavy pages.

## Gap Analysis Summary

### Current State Assessment (70/100)

**Strengths:**

- ✅ Excellent type safety with TypeScript throughout
- ✅ Well-structured query key factory with hierarchical keys
- ✅ Documented staleTime strategy per domain
- ✅ Proper authentication guards with route protection
- ✅ Clean separation of concerns in API layer
- ✅ Robust error parsing and handling
- ✅ Smart request deduplication and token refresh
- ✅ Proper Atomic Design component hierarchy

**Critical Gaps:**

- ❌ Not using queryOptions pattern (scattered query definitions)
- ❌ No placeholderData/keepPreviousData for pagination UX
- ❌ Manual data fetching in customer list route
- ❌ Forms don't use selectors (performance issues)
- ❌ No optimistic updates (slower perceived performance)
- ❌ Missing global error/notFound router components
- ❌ Store subscriptions without selectors
- ❌ Metadata store duplicates Query cache
- ❌ Language not in query keys for localized data
- ❌ Broad invalidation strategies

**Implementation Priorities:**

**High Priority (Phases 1-3):**

- queryOptions factories for type-safe reuse
- keepPreviousData for smooth pagination
- Form selectors for performance
- Route-level data loading

**Medium Priority (Phases 4-6):**

- Safe optimistic updates
- Router error boundaries
- Store selector refactoring
- Remove metadata store duplication

**Low Priority (Phases 7-10):**

- Feature-based translation namespaces
- Lazy-loading translations
- Route-driven language
- Comprehensive testing
- Documentation and ADRs

## Implementation Notes

### Phase Dependencies

- Phase 2 depends on Phase 1 (needs queryOptions for keepPreviousData)
- Phase 4 depends on Phase 1 (optimistic updates need queryOptions)
- Phase 8 depends on Phase 1 (invalidation uses queryOptions)
- Phase 9 depends on Phase 7 (route language needs namespace structure)

### Breaking Changes

- Phase 1: Refactoring hooks to use queryOptions (internal API change)
- Phase 6: Removing metadataStore (consumer code must migrate to Query)
- Phase 7: Splitting translation files (import paths change)

### Success Metrics

- **Phase 1-3 Complete:** 85/100 on best practices alignment
- **Phase 1-6 Complete:** 92/100 on best practices alignment
- **Phase 1-7 Complete:** 98/100 on best practices alignment (100% may be unrealistic for active development)
- **Performance:** 60%+ reduction in unnecessary re-renders
- **Bundle Size:** 20-30% reduction with lazy-loaded translations
- **User Experience:** Zero loading flicker during pagination
- **Code Quality:** 100% TypeScript strict mode compliance
- **Test Coverage:** 80%+ for critical flows
