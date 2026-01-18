---
title: "GORM Model Embedded in Domain Models Exposed to API Responses"
date: 2026-01-18
priority: medium
category: consistency
status: resolved
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

### Create Response DTOs (Recommended)

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

## Resolution

**Status:** ✅ Resolved  
**Date:** 2026-01-19  
**Approach Taken:** Option 1 (Updated code to match instructions)

### Harmonization Summary

Fixed the GORM model exposure issue in the onboarding domain where the `CompleteOnboarding` endpoint was returning a raw `User` model wrapped in `gin.H`. This violated the established pattern of using explicit response DTOs with camelCase JSON fields.

### Pattern Applied

All handlers now use explicit response DTOs created via `To<Model>Response()` converter functions:
- Account domain: `UserResponse`, `WorkspaceResponse`, `UserInvitationResponse`, `LoginResponse`
- Billing domain: `PlanResponse`, `SubscriptionResponse`
- All responses use camelCase JSON field names (`createdAt`, `updatedAt`, not `CreatedAt`, `UpdatedAt`)

### Files Changed

1. **backend/internal/domain/onboarding/handler_http.go**
   - Line 261: Changed `gin.H{"user": user, ...}` to `gin.H{"user": account.ToUserResponse(user), ...}` in `CompleteOnboarding` handler
   - Lines 35-50: Removed unused `completeResponse` struct that directly embedded raw User model
   - Line 242: Updated Swagger annotation from `completeResponse` to `account.LoginResponse`

2. **backend/docs/swagger.json** - Regenerated
3. **backend/docs/swagger.yaml** - Regenerated

### Migration Completeness

- Total instances found: 2 (onboarding handler + unused response type)
- Instances harmonized: 2
- Remaining drift: 0 (verified via grep search)
- Account & Billing domains already had proper DTOs in place

### Validation Results

**Testing:**
- ✅ E2E tests: 12/12 Onboarding tests passed
- ✅ Full suite: 74 tests passed (100%)
- ✅ Response casing: Verified LoginResponse returns camelCase fields
- ✅ No PascalCase GORM fields leak into responses
- ✅ OpenAPI/Swagger regenerated correctly

### Instruction Files Updated

**1. backend-core.instructions.md**
   - Enhanced "Responses and errors (RFC7807)" section with explicit DTO pattern examples
   - Added code examples showing ✅ CORRECT vs ❌ WRONG patterns
   - Added three anti-patterns to prevent recurrence:
     - Returning raw GORM models directly
     - Wrapping raw models in `gin.H{}`
     - Embedding GORM models in response structs
   - Added reference implementation pointer to Account domain model_response.go

**2. responses-dtos-swagger.instructions.md**
   - Rewrote section 2.1 "Don't return GORM models directly" to be CRITICAL/MANDATORY
   - Expanded with detailed anti-pattern examples showing the exact mistakes made
   - Added "Correct Pattern" section with full working examples
   - Made it clear: "ALWAYS create explicit response DTO", "ALWAYS use To*Response converters", "NEVER wrap raw models in gin.H{}"

### Prevention Measures

These drift-specific instruction updates ensure recurrence is unlikely:

1. **Explicit Mandatory Rules:**
   - Instruction files now state "MANDATORY: Never return GORM models directly"
   - Pattern is explicitly enforced in multiple places

2. **Anti-Pattern Documentation:**
   - Three specific anti-patterns documented with code examples
   - Shows the exact drift that occurred (wrapping raw user model in gin.H)
   - Clear correction with working code samples

3. **Reference Implementations:**
   - Backend-core now points to Account and Billing domains as ground-truth examples
   - Makes it clear where to look for the correct pattern
   - New developers/agents will find these references immediately

4. **Root Cause Prevention:**
   - Instructions explain WHY this pattern matters: camelCase standard, Swagger/portal consistency, abstraction violation
   - Covers all three problem areas: code, tests, and documentation

Pattern is now codified in instruction files to prevent future drift. Any new handler returning a model will immediately violate documented anti-patterns.
