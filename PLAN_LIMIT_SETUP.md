# Plan Limit Enforcement for Workspace Invitations

The workspace invitation system now includes **automatic enforcement** of the `MaxTeamMembers` plan limit. This prevents workspaces from inviting more users than their subscription plan allows.

## How It Works

When an admin attempts to invite a new user to the workspace:

1. **Authentication & Authorization** - User must be authenticated and have admin permissions
2. **Active Subscription Check** - `EnforceActiveSubscription` middleware verifies workspace has an active subscription
3. **Plan Limit Check** - `EnforcePlanWorkspaceLimits` middleware automatically:
   - Retrieves the workspace from context
   - Counts current number of users in the workspace
   - Compares against the plan's `MaxTeamMembers` limit
   - Blocks the invitation if limit would be exceeded
4. **Invitation Creation** - If within limits, invitation is created and email sent

## Implementation Details

### Billing Middleware

Location: `internal/domain/billing/middleware_http.go`

Two middleware functions handle plan limit enforcement:

```go
// For business-scoped limits
func EnforcePlanBusinessLimits(feature schema.Field, enforceFunc EnforcePlanLimitFunc) gin.HandlerFunc

// For workspace-scoped limits  
func EnforcePlanWorkspaceLimits(feature schema.Field, enforceFunc EnforcePlanLimitFunc) gin.HandlerFunc
```

The `EnforcePlanWorkspaceLimits` middleware:
- Gets workspace from context (set by `EnforceWorkspaceMembership`)
- Calls the provided counting function with workspace ID
- Checks against subscription plan limits before allowing the action

### Account Service

Location: `internal/domain/account/service.go`

Two methods handle user counting:

```go
// Direct counting by workspace ID
func (s *Service) CountWorkspaceUsers(ctx context.Context, workspaceID string) (int64, error)

// Wrapper that matches EnforcePlanLimitFunc signature
func (s *Service) CountWorkspaceUsersForPlanLimit(ctx context.Context, actor *User, id string) (int64, error)
```

### Account Handler

Location: `internal/domain/account/handler_http.go`

The handler accepts optional middleware via variadic parameters:

```go
func RegisterRoutes(router *gin.Engine, service *Service, enforcePlanLimitMiddleware ...gin.HandlerFunc)
```

The middleware is applied only to the POST invitation endpoint.

### Server Registration

Location: `internal/server/server.go`

The account routes are registered with two middleware functions:

```go
account.RegisterRoutes(r, accountSvc,
    billing.EnforceActiveSubscription(billingSvc),
    billing.EnforcePlanWorkspaceLimits(billing.PlanSchema.MaxTeamMembers, accountSvc.CountWorkspaceUsersForPlanLimit),
)
```

## Error Response

When a workspace attempts to invite a user beyond their plan limit:

**HTTP Status**: `403 Forbidden`

**Response Body**:

```json
{
  "type": "https://kyora.io/errors/plan-limit-reached",
  "title": "Plan Limit Reached",
  "status": 403,
  "detail": "Your current plan allows a maximum of 5 team members",
  "instance": "/v1/workspaces/wrk_xxx/invitations"
}
```

## Testing Plan Limits

To test the plan limit enforcement:

1. **Check current plan limits**:

```bash
curl -X GET http://localhost:8080/v1/billing/subscription \
  -H "Cookie: jwt=YOUR_JWT_TOKEN"
```

2. **Check current team size**:

```bash
curl -X GET http://localhost:8080/v1/workspaces/{workspaceId}/invitations \
  -H "Cookie: jwt=YOUR_JWT_TOKEN"
```

3. **Attempt to invite beyond limit**:

```bash
# If your plan allows 5 members and you have 5, this should fail
curl -X POST http://localhost:8080/v1/workspaces/{workspaceId}/invitations \
  -H "Cookie: jwt=YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "role": "user"
  }'
```

## Architecture Benefits

1. **No Import Cycles** - Uses existing middleware pattern with function parameters
2. **Separation of Concerns** - Business limits use `EnforcePlanBusinessLimits`, workspace limits use `EnforcePlanWorkspaceLimits`
3. **Reusable Pattern** - The same middleware can enforce other workspace-level limits
4. **Consistent Errors** - Uses the same RFC 7807 problem details as other billing limits
5. **Clean Dependencies** - Account service only needs to provide a counting function

## Middleware Chain

The complete middleware chain for the POST invitation endpoint:

1. `auth.EnforceAuthentication` - Validates JWT token
2. `account.EnforceValidActor` - Loads user from database
3. `account.EnforceWorkspaceMembership` - Verifies user belongs to workspace (sets workspace in context)
4. `account.EnforceActorPermissions` - Checks user has admin permissions
5. `billing.EnforceActiveSubscription` - Verifies active subscription (sets subscription in context)
6. `billing.EnforcePlanWorkspaceLimits` - Checks MaxTeamMembers limit
7. Handler executes - Creates invitation

## Future Enhancements

The same pattern can be applied to enforce other limits:

**Workspace-level limits** (using `EnforcePlanWorkspaceLimits`):
- Maximum pending invitations
- Storage quotas per workspace
- API rate limits per workspace

**Business-level limits** (using `EnforcePlanBusinessLimits`):
- Maximum businesses per account
- Maximum orders per month per business
- Revenue tracking limits

Just create a counting function and apply the appropriate middleware to the relevant endpoint.


