---
title: "GORM Model Embedded in Domain Models Exposed to API Responses"
date: 2026-01-18
priority: medium
category: consistency
status: open
domain: backend
---

# GORM Model Embedded in Domain Models Exposed to API Responses

## Summary

Several domain models embed `gorm.Model` directly, which leaks GORM internals (PascalCase timestamps like `CreatedAt`, `UpdatedAt`, `DeletedAt`) into JSON API responses. This violates the backend-core instruction that states: "do not expose storage/ORM implementation details (GORM internals) to clients" and creates inconsistency in response casing (should be camelCase).

## Current State

Multiple domain models embed `gorm.Model` and are returned directly in HTTP responses:

### Locations Found

1. **backend/internal/domain/billing/model.go** (lines 149, 260):
   ```go
   type Plan struct {
       gorm.Model
       // ... other fields
   }
   
   type Subscription struct {
       gorm.Model
       // ... other fields
   }
   ```

2. **backend/internal/domain/account/model.go** (lines 23, 74, 106, 223):
   ```go
   type Workspace struct {
       gorm.Model
       // ... other fields
   }
   
   type User struct {
       gorm.Model
       // ... other fields
   }
   
   type UserInvitation struct {
       gorm.Model
       // ... other fields
   }
   
   type Session struct {
       gorm.Model
       // ... other fields
   }
   ```

## Expected State

Per backend-core.instructions.md:

> **All JSON fields returned to clients must be `camelCase`.**
> 
> - ✅ `createdAt`, `updatedAt`, `deletedAt`
> - ❌ `CreatedAt`, `UpdatedAt`, `DeletedAt`
> 
> Rule: do not expose GORM's embedded `gorm.Model` fields directly to JSON.

And from responses-dtos-swagger.instructions.md:

> **Policy:** handlers should return explicit response DTOs.

## Impact

- **Medium**: API responses include PascalCase timestamp fields (e.g., `CreatedAt`) instead of camelCase (`createdAt`).
- Portal-web may be modeling responses with wrong casing or have inconsistent type definitions.
- Swagger/OpenAPI generation includes GORM internal fields with wrong casing.
- Breaks backend's stated JSON naming convention standard.

## Affected Endpoints

Any endpoint that returns these models directly:

- Billing: Plan and Subscription endpoints
- Account: User, Workspace, Invitation, Session endpoints

## Suggested Fix

### Option 1: Create Response DTOs (Recommended)

Create explicit response types that map from domain models:

```go
// backend/internal/domain/account/model_response.go
type UserResponse struct {
    ID        string    `json:"id"`
    Email     string    `json:"email"`
    FirstName string    `json:"firstName"`
    // ... other fields
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}

func ToUserResponse(u *User) UserResponse {
    return UserResponse{
        ID:        u.ID,
        Email:     u.Email,
        FirstName: u.FirstName,
        // ...
        CreatedAt: u.CreatedAt,
        UpdatedAt: u.UpdatedAt,
    }
}
```

Update handlers to return `ToUserResponse(user)` instead of `user`.

### Option 2: Custom JSON Tags on Models

Add explicit JSON tags to domain models and don't embed `gorm.Model`:

```go
type User struct {
    ID        string    `gorm:"primaryKey" json:"id"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
    
    Email     string `json:"email"`
    // ...
}
```

**Note**: Option 1 is preferred as it better separates storage concerns from API contracts and aligns with the DTO pattern mentioned in responses-dtos-swagger.instructions.md.

## Related Instructions

- `.github/instructions/backend-core.instructions.md` (JSON casing standards)
- `.github/instructions/responses-dtos-swagger.instructions.md` (DTO layer guidance)

## Related Issues

None yet identified.

## Notes

- This is a **consistency drift**, not a functional bug (endpoints work, but violate established patterns).
- Should be addressed before portal-web implements more features that depend on these response shapes.
- After fixing, regenerate Swagger docs with `make openapi`.
