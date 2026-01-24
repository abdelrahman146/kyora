---
description: "Kyora inventory SSOT (backend + portal-web): products, variants, categories, stock semantics, search/sort, summary, RBAC"
applyTo: "backend/internal/domain/inventory/**,portal-web/src/features/inventory/**"
---

# Kyora Inventory SSOT (Backend + Portal Web)

This file is the **single source of truth** for inventory behavior that is **implemented today** across:

- Backend: `backend/internal/domain/inventory/**` + route wiring in `backend/internal/server/routes.go`
- Portal Web: `portal-web/src/api/inventory.ts` + `portal-web/src/routes/business/$businessDescriptor/inventory/**` and related sheets/forms

If you change inventory behavior, you must keep backend + portal-web consistent.

## Non-negotiables

- **Business-scoped always:** inventory data is scoped under `/v1/businesses/:businessDescriptor/...` and must never leak across businesses.
- **RBAC enforced on every route:** inventory endpoints are guarded by `role.ResourceInventory`.
- **Backend is the API contract:** portal-web must follow backend JSON shapes and semantics.
- **Errors are RFC7807 ProblemDetails:** rely on ProblemDetails typing and translation behavior (see `.github/instructions/frontend/_general/http-client.instructions.md`).

## Backend: API surface (authoritative)

All routes are under:

- `/v1/businesses/:businessDescriptor/inventory`

### Products

- `GET /products`
  - Query:
    - Pagination: `page`, `pageSize`
    - Sorting: `orderBy` (repeatable). Use `-field` for descending.
    - Filtering: `categoryId`, `stockStatus` (`in_stock|low_stock|out_of_stock`)
    - Search: `search` (normalized via `list.NormalizeSearchTerm`)
  - Response: `list.ListResponse<Product>` (camelCase list metadata)

- `GET /products/:productId` → returns product **including `variants`**
- `POST /products` → create product
- `POST /products/with-variants` → atomic create product + variants
- `PATCH /products/:productId` → update product (renames variants if name changes)
- `DELETE /products/:productId` → deletes product (variants cascade)
- `GET /products/:productId/variants` → list variants for a product

### Variants

- `GET /variants` → list variants across the business
- `GET /variants/:variantId`
- `POST /variants` → create variant (SKU can be auto-generated)
- `PATCH /variants/:variantId` → updates + normalization
- `DELETE /variants/:variantId`

### Categories

- `GET /categories` → list categories (no pagination)
- `GET /categories/:categoryId`
- `POST /categories` → create category
- `PATCH /categories/:categoryId`
- `DELETE /categories/:categoryId`

### Summary / insights

- `GET /summary?topLimit=N`
  - Returns computed metrics for the business:
    - `productsCount`, `variantsCount`, `categoriesCount`
    - `lowStockVariantsCount`, `outOfStockVariantsCount`
    - `totalStockUnits`, `inventoryValue`
    - `topProductsByInventoryValue` (array of Products)

- `GET /top-products?limit=N`
  - Returns an array of `{ product, inventoryValue }` ordered by inventory value DESC.

## Backend: JSON shapes (what clients must assume)

### List response metadata is camelCase

Inventory list endpoints use the shared list response:

- `items`
- `page`
- `pageSize`
- `totalCount`
- `totalPages`
- `hasMore`

This is verified by e2e tests such as:

- `backend/internal/tests/e2e/inventory_products_test.go`

### Inventory model JSON is also camelCase

Backend responses for inventory models use camelCase keys (examples verified by e2e):

- Product includes `businessId`, `categoryId`, and `variants[]`.
- Variant includes `productId`, `stockQuantity`, `stockQuantityAlert`.

Portal-web inventory typings currently contain **snake_case drift** (e.g. `business_id`, `page_size`). When modifying portal-web inventory code, align to backend’s camelCase shapes (match how `portal-web/src/api/order.ts` models list responses).

## Backend: stock semantics (variant-level)

Backend stock status filters operate on **variants**, not aggregated products:

- **Out of stock:** variant `stock_quantity == 0`
- **Low stock:** variant `stock_quantity <= stock_alert` (and can include zero depending on the scope/query)

Summary metrics are computed from variants:

- `lowStockVariantsCount` counts variants with `stock_quantity <= stock_alert`.
- `outOfStockVariantsCount` counts variants with `stock_quantity == 0`.

Portal-web may compute a product-level stock label (see `portal-web/src/features/inventory/utils/inventoryUtils.ts`) but backend filtering is variant-driven.

## Backend: search + ordering (important)

### Search implementation

`GET /products` supports search across:

- Product search vector
- Variant search vector
- Category search vector
- SKU search via trigram index (`variants.sku`)

This is backed by generated TSVectors + GIN indexes and a trigram GIN index for SKU.

### Ordering rules

- When `search` is provided, the default ordering is rank-based (best match first).
- `orderBy` supports ordinary columns (e.g., `name`) and computed ordering for:
  - `variantsCount`
  - `costPrice`
  - `stock`

Those computed sorts use an aggregation join (LATERAL) and are sensitive to query shape; don’t “simplify” this in service code unless you also update tests and validate query plans.

## Backend: creation/update rules (validation + normalization)

These behaviors are enforced in the service layer and verified by e2e tests:

- `CreateProductWithVariants` is **transactional**: product + variants are created atomically.
- **SKU auto-generation:** when creating a variant with an empty SKU, backend generates one.
- **Variant naming convention:** `Variant.Name = Product.Name + " - " + Variant.Code`.
- `UpdateProduct` renames all variants if the product name changes.
- `UpdateVariant` normalizes:
  - `code`: trimmed (e.g. `"  blue  " → "blue"`)
  - `sku`: trimmed
  - `currency`: uppercased (e.g. `"egp" → "EGP"`)

Categories normalize `descriptor` to a lowercase, trimmed slug-like value (verified by e2e tests).

## Backend: error mapping (ProblemDetails)

The shared response layer maps common DB errors:

- Record not found (GORM) → `404 Not Found` ProblemDetails
- Unique constraint violation → `409 Conflict` ProblemDetails
- Everything else → `500 Internal Error` ProblemDetails

Domain errors in `backend/internal/domain/inventory/errors.go` use the same ProblemDetails format.

## Backend: tests are the truth

When changing inventory behavior, update/add e2e tests under:

- `backend/internal/tests/e2e/inventory_*_test.go`

These tests codify real expectations like:

- List response uses camelCase metadata
- Product delete cascades variants
- SKU auto-generation
- Variant renaming on product rename
- RBAC: view allowed for `user`, manage forbidden

## Portal Web: current implementation (how it works today)

### Where the feature lives

**File placement SSOT:** `.github/instructions/frontend/projects/portal-web/code-structure.instructions.md`

- The items below are _current code locations_, not a requirement.
- Any new/refactored inventory UI must live under `portal-web/src/features/inventory/**`.

- API client + hooks: `portal-web/src/api/inventory.ts`
- Route wrappers: `portal-web/src/routes/business/$businessDescriptor/inventory/**`
- Feature components: `portal-web/src/features/inventory/components/**`
- Feature utils: `portal-web/src/features/inventory/utils/inventoryUtils.ts`
- Query keys (staleTime + invalidation): `portal-web/src/lib/queryKeys.ts`

### Routing + URL-driven state

Inventory list uses TanStack Router search params for:

- `search`, `page`, `pageSize`
- `sortBy`, `sortOrder` (translated to backend `orderBy` with `-` prefix)
- `categoryId`, `stockStatus`

### Data-fetching patterns

- Fetch via `get/post/patch/delVoid` in `portal-web/src/api/client.ts` (ky client with token refresh).
- Use TanStack Query keys under `queryKeys.inventory.*`.
- After mutations (create/update/delete), invalidate broadly with `queryKeys.inventory.all`.

### API Contract Alignment (Critical)

**Portal-web API types MUST match backend JSON exactly (camelCase).**

- `ListResponse` fields: `pageSize`, `totalCount`, `totalPages`, `hasMore` (NOT snake_case)
- `CreateVariantRequest.productId` (NOT `product_id`)
- `InventorySummaryResponse` fields: `productsCount`, `variantsCount`, `categoriesCount`, `lowStockVariantsCount`, `outOfStockVariantsCount`, `totalStockUnits`, `inventoryValue`, `topProductsByInventoryValue` (all camelCase)
- All Variant/Product fields use camelCase to match backend DTOs

This is non-negotiable and verified by backend E2E tests. See `.github/instructions/responses-dtos-swagger.instructions.md` for API contract standards.

## Extension checklist (when adding/changing inventory behavior)

Backend:

- Add/extend handler in `backend/internal/domain/inventory/handler_http.go`.
- Keep business scoping and ownership checks inside the service layer.
- If you add new list filters/sorts, update:
  - service query building
  - indexes/scopes in `storage.go` if needed
  - e2e coverage

Portal Web:

- Add API method + query/mutation hook in `portal-web/src/api/inventory.ts`.
- Update the route search schema and URL wiring in the inventory route.
- Keep UI Arabic/RTL-first and use existing UI primitives (see `.github/instructions/kyora/ux-strategy.instructions.md` and `.github/instructions/frontend/_general/ui-patterns.instructions.md`).
