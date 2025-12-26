# Localized Error Handling Guide

## Overview

The error handling system now returns **translation keys** instead of hardcoded English strings. This enables full localization support for Arabic (primary) and English.

## Key Components

### 1. `errorParser.ts` - Returns Translation Keys

```typescript
import { parseProblemDetails, type ErrorResult } from "@/lib/errorParser";

// ErrorResult structure
interface ErrorResult {
  key: string; // i18n translation key (e.g., "errors.http.401")
  params?: Record<string, string | number>; // Interpolation params
  fallback?: string; // Backend error message (if provided)
}
```

### 2. `translateError.ts` - Helper for Translation

```typescript
import { translateError, translateErrorAsync } from "@/lib/translateError";
import { useTranslation } from "react-i18next";

const { t } = useTranslation();
const errorResult = await parseProblemDetails(error);
const message = translateError(errorResult, t);
```

### 3. Translation Files

- **English**: `src/i18n/locales/en/errors.json`
- **Arabic**: `src/i18n/locales/ar/errors.json`

## Usage Examples

### Example 1: Login Form with Error Toast

```tsx
import { useTranslation } from "react-i18next";
import { authApi } from "@/api/auth";
import { translateErrorAsync } from "@/lib/translateError";
import toast from "react-hot-toast";

function LoginForm() {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      await authApi.login({ email, password });
      toast.success(t("auth.login.success"));
      navigate("/dashboard");
    } catch (error) {
      // Translate error and show toast
      const message = await translateErrorAsync(error, t);
      toast.error(message);
    } finally {
      setLoading(false);
    }
  };

  return <form onSubmit={handleSubmit}>{/* form fields */}</form>;
}
```

### Example 2: Inline Error Display with React Hook Form

```tsx
import { useTranslation } from "react-i18next";
import { parseProblemDetails } from "@/lib/errorParser";
import { translateError } from "@/lib/translateError";

function RegisterForm() {
  const { t } = useTranslation();
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const onSubmit = async (data: RegisterFormData) => {
    try {
      await authApi.register(data);
    } catch (error) {
      // Parse and translate error
      const errorResult = await parseProblemDetails(error);
      const message = translateError(errorResult, t);
      setErrorMessage(message);
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      {errorMessage && (
        <div className="alert alert-error">{errorMessage}</div>
      )}
      {/* form fields */}
    </form>
  );
}
```

### Example 3: Global Error Handler with Context

```tsx
import { useTranslation } from "react-i18next";
import { translateErrorAsync } from "@/lib/translateError";

function useGlobalErrorHandler() {
  const { t } = useTranslation();

  const handleError = async (error: unknown) => {
    const message = await translateErrorAsync(error, t);

    // Log to error tracking service
    console.error("[Error]", message, error);

    // Show toast notification
    toast.error(message);
  };

  return handleError;
}

// Usage in components
function MyComponent() {
  const handleError = useGlobalErrorHandler();

  const doSomething = async () => {
    try {
      await someApiCall();
    } catch (error) {
      handleError(error);
    }
  };
}
```

### Example 4: Validation Errors (Field-Level)

```tsx
import { parseValidationErrors } from "@/lib/errorParser";
import { useTranslation } from "react-i18next";

function ProfileForm() {
  const { t } = useTranslation();
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});

  const onSubmit = async (data: ProfileData) => {
    try {
      await profileApi.update(data);
    } catch (error) {
      // Extract field-level validation errors
      const validationErrors = await parseValidationErrors(error);

      if (validationErrors) {
        // Translate each field error
        const translated = Object.entries(validationErrors).reduce(
          (acc, [field, key]) => ({
            ...acc,
            [field]: t(key as any, { defaultValue: key }),
          }),
          {}
        );
        setFieldErrors(translated);
      } else {
        // Show generic error
        const message = await translateErrorAsync(error, t);
        toast.error(message);
      }
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <Input
        name="email"
        error={fieldErrors.email}
        // ...
      />
    </form>
  );
}
```

### Example 5: Custom Error Keys from Backend

If your backend returns specific error types, you can map them to translation keys:

```tsx
import { parseProblemDetails } from "@/lib/errorParser";
import { translateError } from "@/lib/translateError";

function PasswordResetForm() {
  const { t } = useTranslation();

  const onSubmit = async (data: ResetPasswordData) => {
    try {
      await authApi.resetPassword(data);
    } catch (error) {
      const errorResult = await parseProblemDetails(error);

      // Override key for specific error types
      if (errorResult.key === "errors.http.401") {
        // Check if backend provided a specific error type
        const isExpiredToken = errorResult.fallback?.includes("expired");

        if (isExpiredToken) {
          toast.error(t("auth.reset_token_expired"));
          return;
        }
      }

      // Use default translation
      const message = translateError(errorResult, t);
      toast.error(message);
    }
  };
}
```

## Translation Key Structure

### HTTP Status Codes

```
errors.http.400  -> Invalid request
errors.http.401  -> Unauthorized
errors.http.403  -> Forbidden
errors.http.404  -> Not found
errors.http.422  -> Validation error
errors.http.429  -> Too many requests
errors.http.500  -> Server error
errors.http.502  -> Bad gateway
errors.http.503  -> Service unavailable
errors.http.504  -> Gateway timeout
errors.http.4xx  -> Generic client error (with status param)
errors.http.5xx  -> Generic server error (with status param)
```

### Network Errors

```
errors.network.timeout     -> Request timeout
errors.network.connection  -> Network/connection error
```

### Generic Errors

```
errors.generic.unexpected  -> Unexpected error
errors.generic.message     -> Generic error with custom message
```

### Authentication Errors (Custom)

```
errors.auth.invalid_credentials  -> Wrong email/password
errors.auth.session_expired      -> Session expired
errors.auth.email_not_verified   -> Email not verified
errors.auth.account_locked       -> Account locked
```

### Validation Errors (Custom)

```
errors.validation.required           -> Required field
errors.validation.invalid_email      -> Invalid email format
errors.validation.invalid_password   -> Password too short
errors.validation.password_mismatch  -> Passwords don't match
```

## Adding New Error Keys

### 1. Add to English translations (`en/errors.json`)

```json
{
  "errors": {
    "business": {
      "not_found": "Business not found",
      "access_denied": "You don't have access to this business"
    }
  }
}
```

### 2. Add to Arabic translations (`ar/errors.json`)

```json
{
  "errors": {
    "business": {
      "not_found": "لم يتم العثور على النشاط التجاري",
      "access_denied": "ليس لديك حق الوصول إلى هذا النشاط التجاري"
    }
  }
}
```

### 3. Use in code

```typescript
// Option 1: Direct translation key
toast.error(t("errors.business.not_found"));

// Option 2: Return from backend as ProblemDetails.detail
// Backend: return &problem.Problem{Detail: "errors.business.not_found"}
// Frontend will automatically use it as translation key
```

## Best Practices

1. **Always use `translateError()` or `translateErrorAsync()`** - Don't try to translate error results manually
2. **Provide fallbacks** - Backend error messages serve as fallbacks if translation key is missing
3. **Use specific keys when possible** - Custom auth/validation keys are better than generic HTTP status keys
4. **Keep keys namespaced** - Use dot notation: `errors.domain.specific_error`
5. **Test both languages** - Always verify error messages in Arabic and English
6. **Update both translation files** - When adding new keys, update both `en/errors.json` and `ar/errors.json`

## Migration from Old Code

If you have code using the old `parseProblemDetails()` that returned strings:

```typescript
// OLD (returns string)
const message = await parseProblemDetails(error);
toast.error(message);

// NEW (returns ErrorResult → translate)
const { t } = useTranslation();
const message = await translateErrorAsync(error, t);
toast.error(message);

// OR
const errorResult = await parseProblemDetails(error);
const message = translateError(errorResult, t);
toast.error(message);
```
