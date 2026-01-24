---
description: Frontend forms validation - Zod schemas, async validation, cross-field, array validation (reusable across portal-web, storefront-web)
applyTo: "portal-web/**,storefront-web/**"
---

# Forms Validation

Zod-based validation patterns for TanStack Form.

**Cross-refs:**

- Core forms: `./forms.instructions.md`
- HTTP client: `./http-client.instructions.md`

---

## 1. Basic Zod Schemas

```typescript
import { z } from "zod";

// Email
z.string().email("validation.invalid_email");

// Password
z.string().min(8, "validation.password_min_length");

// Number
z.number().positive("validation.positive_number");

// Integer
z.number().int("validation.integer_only");

// URL
z.string().url("validation.invalid_url");

// Phone (basic)
z.string().regex(/^\+?[0-9]{10,15}$/, "validation.invalid_phone");

// Date
z.date()
  .min(new Date("2020-01-01"), "validation.date_too_old")
  .max(new Date(), "validation.date_cannot_be_future");
```

---

## 2. Custom Validators

### Min/Max with Interpolation

```typescript
// Length
z.string()
  .min(3, "validation.min_length") // Interpolates {{min}}
  .max(100, "validation.max_length"); // Interpolates {{max}}

// Number range
z.number()
  .min(0, "validation.min_value") // Interpolates {{min}}
  .max(1000, "validation.max_value"); // Interpolates {{max}}
```

### Refine with Custom Logic

```typescript
z.string().refine((val) => val.trim().length > 0, {
  message: "validation.required_non_empty",
});

z.number().refine((val) => val % 5 === 0, {
  message: "validation.must_be_multiple_of_5",
});
```

### Transform + Validate

```typescript
z.string()
  .transform((val) => val.trim().toLowerCase())
  .refine((val) => val.length >= 3, { message: "validation.min_length" });
```

---

## 3. Cross-Field Validation

### Password Confirmation

```tsx
<form.AppField
  name="confirmPassword"
  validators={{
    onChangeListenTo: ["password"],
    onBlur: ({ value, fieldApi }) => {
      const password = fieldApi.form.getFieldValue("password");
      if (value !== password) {
        return "validation.passwords_must_match";
      }
      return undefined;
    },
  }}
/>
```

### Date Range

```tsx
<form.AppField
  name="endDate"
  validators={{
    onChangeListenTo: ["startDate"],
    onBlur: ({ value, fieldApi }) => {
      const startDate = fieldApi.form.getFieldValue("startDate");
      if (!startDate || !value) return undefined;

      if (value <= startDate) {
        return "validation.end_date_must_be_after_start_date";
      }

      // Max 90 days
      const days = Math.ceil(
        (value.getTime() - startDate.getTime()) / (1000 * 60 * 60 * 24),
      );
      if (days > 90) {
        return "validation.date_range_max_90_days";
      }

      return undefined;
    },
  }}
/>
```

### Conditional Required

```tsx
<form.AppField
  name="phoneNumber"
  validators={{
    onChangeListenTo: ["hasPhoneSupport"],
    onChange: ({ value, fieldApi }) => {
      const requirePhone = fieldApi.form.getFieldValue("hasPhoneSupport");
      if (requirePhone && !value) {
        return "validation.phone_required";
      }
      return undefined;
    },
  }}
/>
```

---

## 4. Async Validation

### Username Availability

```tsx
<form.AppField
  name="username"
  validators={{
    onChange: z.string().min(3, "validation.username_too_short"),
    onChangeAsync: async ({ value }) => {
      if (value.length < 3) return undefined; // Skip if too short

      const exists = await checkUsernameExists(value);
      return exists ? "validation.username_taken" : undefined;
    },
    onChangeAsyncDebounceMs: 500, // Debounce API calls
  }}
>
  {(field) => (
    <field.TextField
      label="Username"
      suffix={field.state.meta.isValidating ? <Spinner /> : null}
    />
  )}
</form.AppField>
```

### Email Availability

```tsx
<form.AppField
  name="email"
  validators={{
    onBlur: z.string().email("validation.invalid_email"),
    onChangeAsync: async ({ value }) => {
      if (!z.string().email().safeParse(value).success) {
        return undefined; // Skip if invalid format
      }

      const exists = await checkEmailExists(value);
      return exists ? "validation.email_already_registered" : undefined;
    },
    onChangeAsyncDebounceMs: 500,
  }}
/>
```

### Promo Code Validation

```tsx
<form.AppField
  name="promoCode"
  validators={{
    onChangeAsync: async ({ value }) => {
      if (!value || value.length < 4) return undefined;

      try {
        const result = await validatePromoCode(value);
        if (!result.valid) {
          return result.error || "validation.promo_code_invalid";
        }
        return undefined;
      } catch {
        return "validation.promo_code_check_failed";
      }
    },
    onChangeAsyncDebounceMs: 500,
  }}
/>
```

---

## 5. Array Validation

### Min/Max Items

```typescript
z.array(z.string()).min(1, "validation.select_at_least_one");
z.array(z.string()).max(5, "validation.select_too_many");
```

### Unique Values

```typescript
z.array(z.string()).refine((arr) => new Set(arr).size === arr.length, {
  message: "validation.duplicate_values",
});
```

### Array Validation Utilities

```typescript
import {
  validateArrayLength,
  validateUniqueValues,
  validateArrayItems,
  validateNoOverlap,
  validateArrayAnd,
  validateArrayOr,
  validateArrayCount,
} from '@/lib/form';

// Min/max items
<form.AppField
  name="phoneNumbers"
  validators={{
    onChange: ({ value }) => validateArrayLength(value, { min: 1, max: 5 })
  }}
/>

// Unique values
<form.AppField
  name="emails"
  validators={{
    onChange: ({ value }) => validateUniqueValues(value, {
      extractor: (item) => item.email,
      errorKey: 'validation.duplicate_emails',
    })
  }}
/>

// Per-item validation
<form.AppField
  name="addresses"
  validators={{
    onChange: ({ value }) => validateArrayItems(value, (item, index) => {
      if (!item.street) return 'validation.street_required';
      if (!item.city) return 'validation.city_required';
      return undefined;
    })
  }}
/>

// No overlap (time ranges)
<form.AppField
  name="timeSlots"
  validators={{
    onChange: ({ value }) => validateNoOverlap(value, {
      extractor: (item) => ({ start: item.startTime, end: item.endTime }),
      errorKey: 'validation.time_slots_overlap',
    })
  }}
/>

// Combine validators (AND)
<form.AppField
  name="tags"
  validators={{
    onChange: ({ value }) => validateArrayAnd(value, [
      (arr) => validateArrayLength(arr, { min: 1, max: 10 }),
      (arr) => validateUniqueValues(arr, { extractor: (item) => item }),
    ])
  }}
/>

// Count validation (exactly one primary)
<form.AppField
  name="addresses"
  validators={{
    onChange: ({ value }) => validateArrayCount(value, {
      extractor: (item) => item.isPrimary,
      matchValue: true,
      exactCount: 1,
      errorKey: 'validation.exactly_one_primary_address',
    })
  }}
/>
```

---

## 6. Complex Schemas

### Nested Objects

```typescript
const addressSchema = z.object({
  street: z.string().min(1, "validation.required"),
  city: z.string().min(1, "validation.required"),
  country: z.string().min(2, "validation.country_code_required"),
  postalCode: z.string().optional(),
});

const customerSchema = z.object({
  name: z.string().min(2, "validation.min_length"),
  email: z.string().email("validation.invalid_email"),
  phone: z.string().regex(/^\+?[0-9]{10,15}$/, "validation.invalid_phone"),
  address: addressSchema,
});
```

### Array of Objects

```typescript
const phoneSchema = z.object({
  number: z.string().regex(/^\+?[0-9]{10,15}$/, "validation.invalid_phone"),
  type: z.enum(["mobile", "work", "home"]),
  isPrimary: z.boolean(),
});

const customerSchema = z.object({
  name: z.string().min(2, "validation.min_length"),
  phones: z
    .array(phoneSchema)
    .min(1, "validation.at_least_one_phone")
    .refine((phones) => phones.filter((p) => p.isPrimary).length === 1, {
      message: "validation.exactly_one_primary_phone",
    }),
});
```

### Discriminated Unions

```typescript
const paymentMethodSchema = z.discriminatedUnion("type", [
  z.object({
    type: z.literal("cash"),
  }),
  z.object({
    type: z.literal("card"),
    cardNumber: z.string().min(16, "validation.card_number_invalid"),
    expiryDate: z.string().regex(/^\d{2}\/\d{2}$/, "validation.expiry_invalid"),
  }),
  z.object({
    type: z.literal("bank_transfer"),
    accountNumber: z.string().min(1, "validation.required"),
    bankName: z.string().min(1, "validation.required"),
  }),
]);
```

---

## 7. Dependent Fields

### Auto-Calculate

```typescript
const form = useKyoraForm({
  defaultValues: {
    quantity: 1,
    price: 0,
    total: 0,
  },
});

// Watch quantity and price, update total
form.useStore((state) => {
  const { quantity, price } = state.values;
  const newTotal = quantity * price;
  if (state.values.total !== newTotal) {
    form.setFieldValue("total", newTotal);
  }
});
```

### Conditional Fields

```tsx
<form.Subscribe selector={(state) => state.values.hasCustomAddress}>
  {(hasCustomAddress) =>
    hasCustomAddress && (
      <form.AppField name="customAddress">
        {(field) => <field.TextareaField label="Custom Address" />}
      </form.AppField>
    )
  }
</form.Subscribe>
```

---

## 8. Server-Side Validation

### Injecting Backend Errors

```typescript
import { createServerErrorValidator } from "@/lib/form";

const form = useKyoraForm({
  validators: {
    email: {
      onBlur: z.string().email("validation.invalid_email"),
      onServer: createServerErrorValidator(), // Injects RFC7807 errors
    },
    password: {
      onBlur: z.string().min(8, "validation.password_min_length"),
      onServer: createServerErrorValidator(),
    },
  },
  onSubmit: async ({ value }) => {
    try {
      await api.register(value);
    } catch (error) {
      // Field-level errors automatically appear on fields
      // Form-level errors shown in <form.FormError />
    }
  },
});
```

**How it works:**

1. Backend returns RFC7807 with field details
2. `createServerErrorValidator()` parses and injects errors
3. Errors appear on corresponding fields
4. Form-level errors shown in `<form.FormError />`

---

## 9. Validation Helpers

### Email

```typescript
export const emailSchema = z
  .string()
  .email("validation.invalid_email")
  .min(1, "validation.required");
```

### Phone

```typescript
export const phoneSchema = z
  .string()
  .regex(/^\+?[0-9]{10,15}$/, "validation.invalid_phone")
  .min(1, "validation.required");
```

### Password

```typescript
export const passwordSchema = z
  .string()
  .min(8, "validation.password_min_length")
  .regex(/[A-Z]/, "validation.password_uppercase_required")
  .regex(/[a-z]/, "validation.password_lowercase_required")
  .regex(/[0-9]/, "validation.password_number_required");
```

### URL

```typescript
export const urlSchema = z
  .string()
  .url("validation.invalid_url")
  .optional()
  .or(z.literal("")); // Allow empty
```

### Age

```typescript
export const birthDateSchema = z
  .date()
  .max(new Date(), "validation.date_cannot_be_future")
  .refine(
    (date) => {
      const age = Math.floor(
        (Date.now() - date.getTime()) / (365.25 * 24 * 60 * 60 * 1000),
      );
      return age >= 18;
    },
    { message: "validation.must_be_18_or_older" },
  );
```

---

## 10. Validation Patterns Reference

### Required Field

```typescript
z.string().min(1, "validation.required");
```

### Optional Field

```typescript
z.string().optional();
// or
z.string().nullable();
```

### Min/Max Length

```typescript
z.string().min(3, "validation.min_length").max(100, "validation.max_length");
```

### Min/Max Value

```typescript
z.number().min(0, "validation.min_value").max(1000, "validation.max_value");
```

### Regex Pattern

```typescript
z.string().regex(/^[A-Z0-9]+$/, "validation.alphanumeric_uppercase_only");
```

### Enum

```typescript
z.enum(["draft", "pending", "completed"], {
  errorMap: () => ({ message: "validation.invalid_status" }),
});
```

### Array Min/Max

```typescript
z.array(z.string())
  .min(1, "validation.select_at_least_one")
  .max(10, "validation.select_too_many");
```

### Object Shape

```typescript
z.object({
  name: z.string().min(1),
  age: z.number().min(0),
}).strict(); // Disallow extra keys
```

---

## Agent Validation

Before completing validation task:

- ☑ **ZERO** hardcoded error messages - all use `validation.*` keys
- ☑ All validation keys exist in `i18n/en/errors.json`
- ☑ All validation keys exist in `i18n/ar/errors.json`
- ☑ Async validators use debounce (500ms default)
- ☑ Cross-field validators use `onChangeListenTo` or `onBlurListenTo`
- ☑ Array validators use utility functions from `lib/form`
- ☑ Server errors use `createServerErrorValidator()`
- ☑ Password validators enforce min length + complexity

---

## Resources

- Zod Docs: https://zod.dev
- TanStack Form: https://tanstack.com/form/latest
- Implementation: `portal-web/src/lib/form/`
