# Onboarding Flow Refactoring - Backend-Driven Session Management

## Completed Work

### Backend (100% Complete)
✅ Added GET /v1/onboarding/session endpoint
✅ Service method: GetSession(ctx, sessionToken) 
✅ Error handling: ErrSessionTokenRequired
✅ Comprehensive E2E tests in onboarding_get_session_test.go
✅ Route registered in routes.go
✅ Backend compiles successfully

### Frontend API Layer (100% Complete)
✅ GetSessionResponseSchema added to types/onboarding.ts
✅ onboardingApi.getSession() method implemented
✅ TypeScript types exported
✅ All imports updated

### Frontend Context (100% Complete - New File Ready)
✅ New OnboardingContext.tsx created with:
  - loadSession(sessionToken) to restore from backend
  - Removed all sessionStorage dependencies
  - Backend is single source of truth
  - updateFromSession() helper for sync
  - All state derived from GetSessionResponse

## Remaining Work

### 1. Replace OnboardingContext.tsx
```bash
cd portal-web/src/contexts
mv OnboardingContext.tsx.new OnboardingContext.tsx
```

### 2. Update Onboarding Pages
Each page needs to call loadSession() on mount to restore state:

**plan.tsx, email.tsx, verify.tsx, business.tsx, payment.tsx, complete.tsx:**
```typescript
useEffect(() => {
  const params = new URLSearchParams(window.location.search);
  const token = params.get('sessionToken') || sessionToken;
  if (token && !stage) {
    loadSession(token);
  }
}, [loadSession, sessionToken, stage]);
```

### 3. Update Navigation
When navigating between steps, pass sessionToken as query param:
```typescript
navigate(`/onboarding/business?sessionToken=${sessionToken}`);
```

### 4. Remove setBusinessDetails calls
No longer needed - state updates automatically via API responses

### 5. Test E2E Flow
```bash
cd backend
go test ./internal/tests/e2e -run TestOnboardingGetSession -v
```

### 6. Run Quality Checks
```bash
cd portal-web
npm run type-check
npm run lint
```

## Benefits of New Architecture

1. **Simplified State Management**: Backend is single source of truth
2. **Resume Capability**: Users can refresh/return and continue where they left off
3. **No Client Storage**: More secure, no sync issues
4. **Easier Debugging**: All state visible in backend database
5. **Better UX**: Seamless flow even with browser refresh

## API Usage Pattern

```typescript
// Start new session
const response = await onboardingApi.startSession({ email, planDescriptor });
startSession(response, email, plan);

// Resume existing session (on page load/refresh)
const token = getSessionTokenFromURL();
await loadSession(token);

// After any operation, optionally refresh state
await loadSession(sessionToken);
```

## Session Token Flow

1. User starts: POST /start → receives sessionToken
2. sessionToken passed as query param in all navigation
3. On any page load: GET /session?sessionToken=xxx → restores full state
4. All mutations (verify, business, payment) update backend session
5. Frontend stays in sync by reading responses or calling getSession()

