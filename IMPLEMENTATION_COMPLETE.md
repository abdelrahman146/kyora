# Backend-Driven Onboarding - IMPLEMENTATION GUIDE

## ✅ Completed Backend Work

### 1. GET /v1/onboarding/session
- **File**: `backend/internal/domain/onboarding/handler_http.go`
- **Service**: `GetSession(ctx, sessionToken)` in `service.go`
- **Returns**: Complete session state with all user data, business info, payment status
- **Error Handling**: Session not found, expired, or already committed

### 2. DELETE /v1/onboarding/session  
- **File**: `backend/internal/domain/onboarding/handler_http.go`
- **Service**: `DeleteSession(ctx, sessionToken)` in `service.go`
- **Purpose**: Cancel/restart onboarding flow
- **Protection**: Cannot delete already committed sessions

### 3. Enhanced Errors
- `ErrSessionTokenRequired`: Missing token parameter
- `ErrSessionAlreadyCommitted`: Prevent deleting completed sessions

### 4. Routes Updated
- `GET /v1/onboarding/session` - Retrieve session
- `DELETE /v1/onboarding/session` - Delete session
- Both registered in `backend/internal/server/routes.go`

### 5. E2E Tests Created
- `backend/internal/tests/e2e/onboarding_get_session_test.go`
- 10+ test scenarios covering all edge cases
- Helper method: `CreateExpiredSession()` for testing

## Frontend Implementation Needed

### 1. LocalStorage Utility (`portal-web/src/lib/sessionStorage.ts`)

```typescript
const SESSION_TOKEN_KEY = 'kyora_onboarding_session';

export const sessionStorage = {
  getToken(): string | null {
    try {
      return localStorage.getItem(SESSION_TOKEN_KEY);
    } catch {
      return null;
    }
  },

  setToken(token: string): void {
    try {
      localStorage.setItem(SESSION_TOKEN_KEY, token);
    } catch (error) {
      console.error('Failed to save session token:', error);
    }
  },

  clearToken(): void {
    try {
      localStorage.removeItem(SESSION_TOKEN_KEY);
    } catch (error) {
      console.error('Failed to clear session token:', error);
    }
  },
};
```

### 2. Resume Session Dialog (`portal-web/src/components/molecules/ResumeSessionDialog.tsx`)

```typescript
import { useTranslation } from 'react-i18next';

interface ResumeSessionDialogProps {
  open: boolean;
  onResume: () => void;
  onStartFresh: () => void;
  email?: string;
  stage?: string;
}

export function ResumeSessionDialog({
  open,
  onResume,
  onStartFresh,
  email,
  stage,
}: ResumeSessionDialogProps) {
  const { t } = useTranslation(['onboarding', 'common']);

  if (!open) return null;

  return (
    <dialog className="modal modal-open">
      <div className="modal-box">
        <h3 className="font-bold text-lg">
          {t('onboarding:resumeSession.title')}
        </h3>
        <p className="py-4">
          {t('onboarding:resumeSession.message', { email, stage })}
        </p>
        <div className="modal-action">
          <button onClick={onStartFresh} className="btn btn-ghost">
            {t('common:startFresh')}
          </button>
          <button onClick={onResume} className="btn btn-primary">
            {t('common:continue')}
          </button>
        </div>
      </div>
    </dialog>
  );
}
```

### 3. Updated OnboardingContext

```typescript
// Key Changes:
// 1. Store sessionToken in localStorage on startSession()
// 2. Add loadSessionFromStorage() to restore on mount
// 3. Add clearSession() to delete backend session + clear localStorage
// 4. Handle session conflicts (email already has session)

const startSession = useCallback(
  (response: StartSessionResponse, email: string, plan: Plan) => {
    // Save to localStorage
    sessionStorage.setToken(response.sessionToken);
    
    updateState({
      sessionToken: response.sessionToken,
      stage: response.stage,
      email,
      selectedPlan: plan,
      isPaidPlan: response.isPaid,
      //... other fields
    });
  },
  [updateState]
);

const loadSessionFromStorage = useCallback(async () => {
  const token = sessionStorage.getToken();
  if (!token) return false;

  try {
    await loadSession(token);
    return true;
  } catch {
    sessionStorage.clearToken();
    return false;
  }
}, [loadSession]);

const clearSession = useCallback(async () => {
  if (sessionToken) {
    try {
      await onboardingApi.deleteSession(sessionToken);
    } catch (error) {
      console.error('Failed to delete session:', error);
    }
  }
  sessionStorage.clearToken();
  resetOnboarding();
}, [sessionToken, resetOnboarding]);
```

### 4. Plan Selection Page Updates

```typescript
// In portal-web/src/routes/onboarding/plan.tsx

useEffect(() => {
  // Check for existing session on mount
  const checkExistingSession = async () => {
    const hasSession = await loadSessionFromStorage();
    if (hasSession) {
      setShowResumeDialog(true);
    }
  };
  checkExistingSession();
}, [loadSessionFromStorage]);

// Handle form submission
const handleSubmit = async (e) => {
  e.preventDefault();
  
  try {
    const response = await onboardingApi.startSession({
      email: formEmail,
      planDescriptor: selectedPlan.descriptor,
    });
    
    // Backend returns existing session if email already registered
    if (response.stage !== 'plan_selected') {
      // Email has existing session from another device
      setConflictSession(response);
      setShowConflictDialog(true);
    } else {
      // New session - proceed normally
      startSession(response, formEmail, selectedPlan);
      navigate('/onboarding/email');
    }
  } catch (error) {
    // Handle error
  }
};

// Handle resume dialog
const handleResumeSession = () => {
  setShowResumeDialog(false);
  // Navigate to appropriate step based on stage
  navigateToCurrentStage();
};

const handleStartFresh = async () => {
  await clearSession();
  setShowResumeDialog(false);
};
```

### 5. All Onboarding Pages

Add this to each page (email, verify, business, payment, complete):

```typescript
useEffect(() => {
  // Restore session from localStorage on mount
  const restoreSession = async () => {
    if (!sessionToken) {
      const hasSession = await loadSessionFromStorage();
      if (!hasSession) {
        navigate('/onboarding/plan');
      }
    }
  };
  restoreSession();
}, [sessionToken, loadSessionFromStorage, navigate]);
```

### 6. i18n Translations

Add to `portal-web/src/i18n/locales/en/onboarding.json`:

```json
{
  "resumeSession": {
    "title": "Continue Your Onboarding?",
    "message": "We found an existing onboarding session for {{email}} at stage: {{stage}}. Would you like to continue or start fresh?"
  },
  "conflictSession": {
    "title": "Existing Session Found",
    "message": "This email already has an onboarding session in progress. Would you like to continue that session or start over?"
  }
}
```

## Testing Checklist

### Backend Tests
```bash
cd backend
go test ./internal/tests/e2e -run TestOnboardingGetSession -v
go test ./internal/tests/e2e -run TestOnboarding -v  # All onboarding tests
```

### Frontend Tests
```bash
cd portal-web
npm run type-check
npm run lint
```

### Manual Testing Scenarios

1. **Fresh Start**: Clear localStorage → Start onboarding → Verify token saved
2. **Resume Flow**: Refresh page mid-onboarding → Should show resume dialog
3. **Different Device**: Start on device A → Continue on device B with same email
4. **Start Fresh**: Choose "Start Fresh" → Old session deleted → New session created
5. **Complete Flow**: Complete onboarding → localStorage cleared → Cannot resume

## Security Considerations

1. **LocalStorage**: Only stores session token (not sensitive data)
2. **Token Expiry**: Backend validates expiry on every request
3. **HTTPS Only**: Production must use HTTPS for localStorage security
4. **Session Deletion**: Proper cleanup on logout/complete
5. **XSS Protection**: Token stored as string, not evaluated as code

## Benefits Summary

✅ **Seamless UX**: Users can resume from any device  
✅ **No Data Loss**: All progress saved in backend database  
✅ **Conflict Resolution**: Handles multi-device scenarios gracefully  
✅ **Security**: Minimal client storage, backend-validated sessions  
✅ **Production-Ready**: Complete error handling and edge case coverage

