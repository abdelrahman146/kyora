---
description: "Tenant Isolation SSOT — workspace/business hierarchy, scoping, cross-tenant prevention, RBAC enforcement"
applyTo: "backend/**,portal-web/**"
---

# Tenant Isolation — Single Source of Truth (SSOT)

**SSOT Hierarchy:**

- Parent: `.github/copilot-instructions.md`
- Related: `backend/_general/architecture.instructions.md`, `platform/errors-handling.instructions.md`, all `domain/*.instructions.md`

**When to Read:**

- Implementing any data query (backend or frontend)
- Adding new API endpoints
- Working with workspace/business context
- Implementing RBAC checks
- Debugging cross-tenant access issues

---

## 1) Tenant Hierarchy

Kyora uses **two-level tenant isolation**:

```
Workspace (top-level tenant)
  └─ Business (sub-tenant, multiple per workspace)
      └─ Domain data (customers, orders, inventory, etc.)
```

### Rules:

1. **Workspace is the auth boundary**: JWTs carry `workspaceId`, all authenticated requests must not access other workspaces
2. **Business is the data boundary**: Most domain resources (orders, customers, inventory) belong to a business
3. **Workspace-level resources**: Users, invitations, workspace settings, subscription/billing
4. **Business-level resources**: Everything else (orders, customers, inventory, accounting, analytics)

---

## 2) Backend Scoping Enforcement

### 2.1 Repository Scopes

**MUST use scopes for all queries:**

```go
// ✅ Correct: workspace-scoped
func (r *UserRepository) GetByID(ctx context.Context, workspaceID, userID string) (*User, error) {
    var user User
    err := r.db.WithContext(ctx).
        Scopes(ScopeWorkspaceID(workspaceID)).
        First(&user, "id = ?", userID).Error
    return &user, err
}

// ✅ Correct: business-scoped
func (r *OrderRepository) GetByID(ctx context.Context, businessID, orderID string) (*Order, error) {
    var order Order
    err := r.db.WithContext(ctx).
        Scopes(ScopeBusinessID(businessID)).
        First(&order, "id = ?", orderID).Error
    return &order, err
}

// 🚫 Wrong: no scoping
func (r *OrderRepository) GetByID(ctx context.Context, orderID string) (*Order, error) {
    var order Order
    err := r.db.WithContext(ctx).First(&order, "id = ?", orderID).Error // Cross-tenant leak!
    return &order, err
}
```

**Scope helpers** (defined in `backend/internal/platform/database/scopes.go`):

- `ScopeWorkspaceID(workspaceID string)` — adds `WHERE workspace_id = ?`
- `ScopeBusinessID(businessID string)` — adds `WHERE business_id = ?`
- `ScopeArchived(includeArchived bool)` — filters archived records

### 2.2 Service Layer Validation

Services **MUST validate tenant context** before calling repositories:

```go
func (s *OrderService) GetOrderByID(ctx context.Context, businessID, orderID string) (*Order, error) {
    // Validate inputs
    if businessID == "" || orderID == "" {
        return nil, domain.NewValidationError("business_id and order_id required")
    }

    // Repository automatically scopes by businessID
    order, err := s.repo.GetByID(ctx, businessID, orderID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, domain.NewNotFoundError("order", orderID)
        }
        return nil, domain.WrapError(err, "failed to get order")
    }

    return order, nil
}
```

### 2.3 Middleware Enforcement

**Auth middleware** (`backend/internal/platform/middleware/auth.go`) extracts and validates workspace:

- Parses JWT from `Authorization: Bearer <token>`
- Validates token signature and expiration
- Extracts `workspaceId` claim
- Attaches to request context via `request.WithWorkspaceID(ctx, workspaceID)`

**RBAC middleware** enforces permissions after auth:

- Extracts user role from request context
- Checks resource + action permissions
- Returns **403 Forbidden** if insufficient permissions

### 2.4 Route Parameter Extraction

**Handler functions MUST extract businessID from routes:**

```go
// Route: /v1/businesses/:businessDescriptor/orders/:orderId
func (h *OrderHandler) GetOrder(c *gin.Context) {
    // Extract workspace from auth middleware
    workspaceID := request.GetWorkspaceID(c.Request.Context())

    // Extract business descriptor from route
    businessDescriptor := c.Param("businessDescriptor")

    // Resolve business (validates workspace ownership)
    business, err := h.businessService.GetByDescriptor(c.Request.Context(), workspaceID, businessDescriptor)
    if err != nil {
        response.Error(c, err)
        return
    }

    // Extract order ID
    orderID := c.Param("orderId")

    // Fetch order (scoped to business)
    order, err := h.orderService.GetOrderByID(c.Request.Context(), business.ID, orderID)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.Success(c, http.StatusOK, order)
}
```

**Never trust workspace/business IDs from URL params** for workspace-scoped routes. Always use authenticated context.

---

## 3) Frontend Scoping Enforcement

### 3.1 API Client Context

**TanStack Router context** provides workspace/business IDs:

```typescript
// portal-web/src/api/client.ts extracts from route context
const workspaceId = routerContext.workspaceId;
const businessId = routerContext.businessId;
```

### 3.2 Route Guards

**Routes enforce tenant context** (see `portal-web/src/routes/`):

```typescript
// Workspace-level route
export const Route = createFileRoute("/_auth/workspace")({
  beforeLoad: ({ context }) => {
    if (!context.workspaceId) {
      throw redirect({ to: "/auth/login" });
    }
  },
});

// Business-level route
export const Route = createFileRoute(
  "/_auth/workspace/$businessDescriptor/orders",
)({
  beforeLoad: async ({ context, params }) => {
    const business = await context.queryClient.ensureQueryData(
      businessQueries.byDescriptor(params.businessDescriptor),
    );
    return { business };
  },
});
```

### 3.3 Query Scoping

**All API queries MUST include tenant context:**

```typescript
// ✅ Correct: business-scoped query
export const ordersQueries = {
  list: (businessId: string, filters?: OrderFilters) =>
    queryOptions({
      queryKey: ["orders", "list", businessId, filters],
      queryFn: () => ordersApi.list(businessId, filters),
    }),
};

// 🚫 Wrong: no scoping
export const ordersQueries = {
  list: (
    filters?: OrderFilters, // Missing businessId!
  ) =>
    queryOptions({
      queryKey: ["orders", "list", filters],
      queryFn: () => ordersApi.list(filters),
    }),
};
```

---

## 4) Cross-Tenant Prevention

### 4.1 Error Responses

**MUST return 404 for cross-tenant access attempts** (not 403):

```go
// ✅ Correct: 404 for missing or unauthorized resource
order, err := s.repo.GetByID(ctx, businessID, orderID)
if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domain.NewNotFoundError("order", orderID) // 404
    }
    return nil, domain.WrapError(err, "failed to get order")
}
```

**Why 404, not 403?**

- Prevents information leakage (attacker can't tell if resource exists in another tenant)
- Consistent with REST semantics (resource doesn't exist _in this context_)

### 4.2 Testing Cross-Tenant Access

**E2E tests MUST verify tenant isolation:**

```go
func TestOrderCrossTenantAccess(t *testing.T) {
    // Setup two workspaces
    workspace1 := createTestWorkspace(t, "workspace1")
    workspace2 := createTestWorkspace(t, "workspace2")

    business1 := createTestBusiness(t, workspace1.ID, "business1")
    business2 := createTestBusiness(t, workspace2.ID, "business2")

    order1 := createTestOrder(t, business1.ID)

    // User from workspace2 tries to access order from workspace1
    token := getAuthToken(t, workspace2.ID)

    resp := makeRequest(t, "GET",
        fmt.Sprintf("/v1/businesses/%s/orders/%s", business1.Descriptor, order1.ID),
        token)

    assert.Equal(t, http.StatusNotFound, resp.StatusCode) // Must be 404
}
```

---

## 5) RBAC Integration

**Role-Based Access Control** is enforced at the business level:

### 5.1 Permission Model

```
User
  └─ WorkspaceMembership (role: owner, admin, member)
      └─ BusinessMembership (role per business)
          └─ Permissions (resource + action)
```

### 5.2 Backend RBAC Middleware

```go
// Route with RBAC protection
ordersGroup.GET("/:orderId",
    middleware.RequireAuth(),
    middleware.RequireBusinessAccess(), // Validates workspace → business ownership
    middleware.RequirePermission(role.ResourceOrders, role.ActionView),
    handlers.GetOrder)
```

### 5.3 Frontend Permission Checks

```typescript
// Check if user can manage orders
const canManageOrders = usePermission('orders', 'manage');

// Conditionally render UI
{canManageOrders && <Button onClick={handleDelete}>Delete Order</Button>}
```

**Permission helpers** (see `portal-web/src/lib/permissions.ts`):

- `usePermission(resource, action)` — checks if current user has permission
- `hasPermission(user, resource, action)` — programmatic check

---

## 6) Data Scoping Checklist

**Before implementing any data access:**

- [ ] **Backend Repository**: Does the query use `ScopeWorkspaceID` or `ScopeBusinessID`?
- [ ] **Backend Service**: Does it validate tenant IDs before calling repo?
- [ ] **Backend Handler**: Does it extract workspace from auth context (not URL params)?
- [ ] **Frontend Query**: Does the query key include `businessId` or `workspaceId`?
- [ ] **Frontend Route**: Does the route guard enforce tenant context?
- [ ] **Error Handling**: Does cross-tenant access return 404 (not 403)?
- [ ] **Tests**: Are cross-tenant access attempts covered?

---

## 7) Common Pitfalls

### 🚫 Don't: Trust URL parameters for workspace

```go
// Wrong: workspace from URL can be forged
workspaceID := c.Param("workspaceId") // Attacker can inject any workspace ID!
```

```go
// Correct: workspace from authenticated context
workspaceID := request.GetWorkspaceID(c.Request.Context())
```

### 🚫 Don't: Skip scoping on "admin" routes

```go
// Wrong: even admin endpoints must scope by workspace
func (r *UserRepository) ListAll(ctx context.Context) ([]*User, error) {
    var users []*User
    err := r.db.Find(&users).Error // No scoping = cross-tenant leak!
    return users, err
}
```

```go
// Correct: admin can list users, but only in their workspace
func (r *UserRepository) ListByWorkspace(ctx context.Context, workspaceID string) ([]*User, error) {
    var users []*User
    err := r.db.Scopes(ScopeWorkspaceID(workspaceID)).Find(&users).Error
    return users, err
}
```

### 🚫 Don't: Use business ID from JWT

```go
// Wrong: JWT only carries workspaceId, not businessId
businessID := claims["businessId"] // This doesn't exist!
```

```go
// Correct: business ID comes from route param, validated against workspace
businessDescriptor := c.Param("businessDescriptor")
business, err := h.businessService.GetByDescriptor(ctx, workspaceID, businessDescriptor)
```

---

## 8) References

**Backend:**

- `backend/internal/platform/database/scopes.go` — Scope helpers
- `backend/internal/platform/middleware/auth.go` — Auth middleware
- `backend/internal/platform/middleware/rbac.go` — RBAC middleware
- `backend/internal/platform/request/context.go` — Context helpers

**Frontend:**

- `portal-web/src/api/client.ts` — API client with tenant context
- `portal-web/src/lib/permissions.ts` — Permission helpers
- `portal-web/src/routes/**` — Route guards

**Related Instructions:**

- `backend/_general/architecture.instructions.md` — Backend domain structure
- `platform/errors-handling.instructions.md` — Error response patterns
- All `domain/*.instructions.md` — Domain-specific scoping rules
