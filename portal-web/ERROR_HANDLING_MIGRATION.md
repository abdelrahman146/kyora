# Error Handling Migration Notice

**IMPORTANT: Error handling is now fully localized!**

## What Changed

`parseProblemDetails()` now returns `ErrorResult` with translation keys instead of hardcoded English strings.

## Quick Migration

```typescript
// ❌ OLD (will break - type error)
const message = await parseProblemDetails(error);
toast.error(message);

// ✅ NEW (recommended)
import { translateErrorAsync } from "@/lib/translateError";
import { useTranslation } from "react-i18next";

const { t } = useTranslation();
const message = await translateErrorAsync(error, t);
toast.error(message);

// ✅ NEW (alternative)
import { parseProblemDetails } from "@/lib/errorParser";
import { translateError } from "@/lib/translateError";

const errorResult = await parseProblemDetails(error);
const message = translateError(errorResult, t);
toast.error(message);
```

## Full Documentation

See [LOCALIZED_ERROR_HANDLING.md](./LOCALIZED_ERROR_HANDLING.md) for:
- Complete usage examples
- Translation key structure
- Best practices
- Adding custom error keys
