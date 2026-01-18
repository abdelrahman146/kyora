---
description: Portal Web Form System - Complete Reference
applyTo: "portal-web/**,storefront-web/**"
---

# Portal Web Form System

**SSOT Hierarchy:**

- Parent: copilot-instructions.md
- Peers: ui-implementation.instructions.md, ky.instructions.md
- Required Reading: design-tokens.instructions.md

**When to Read:**

- Creating/modifying forms in portal-web or storefront-web
- Form validation patterns
- Field component implementation
- File upload flows

---

## ⚠️ CRITICAL RULES (MANDATORY - NO EXCEPTIONS)

### 1. Form Context Requirement

**ALL components that use form context (FormRoot, SubmitButton, FormError, Subscribe) MUST be inside `<form.AppForm>`.**

```tsx
// ❌ WRONG
<>
  <form.AppForm>
    <form.FormRoot>{/* fields */}</form.FormRoot>
  </form.AppForm>
  <form.SubmitButton>Submit</form.SubmitButton>  {/* Error! */}
</>

// ✅ CORRECT
<form.AppForm>
  <form.FormRoot>{/* fields */}</form.FormRoot>
  <form.SubmitButton>Submit</form.SubmitButton>
</form.AppForm>
```

**Error symptom:** `Error: formContext only works when within a formComponent passed to createFormHook`

**Bottom Sheet Pattern (STRICT):**

```tsx
// ✅ CORRECT - AppForm wraps EVERYTHING including footer
function MySheet() {
  const form = useKyoraForm({...});
  const formId = useId(); // or hardcoded stable id

  return (
    <form.AppForm>
      <BottomSheet
        footer={
          <div className="flex gap-2">
            <Button variant="ghost" onClick={onClose}>
              {t('common:cancel')}
            </Button>
            <form.SubmitButton form={formId} variant="primary">
              {t('common:save')}
            </form.SubmitButton>
          </div>
        }
      >
        <form.FormRoot id={formId}>
          {/* fields */}
        </form.FormRoot>
      </BottomSheet>
    </form.AppForm>
  );
}
```

### 2. Field Pattern (MANDATORY)

**NEVER use form components directly. ALWAYS use `<form.AppField>` + `{(field) => <field.ComponentName />}` pattern:**

```tsx
// ❌ WRONG - Using component directly
<TextField name="email" label={t('auth.email')} />
<PriceInput name="amount" />

// ✅ CORRECT - Using AppField pattern
<form.AppField name="email" validators={{...}}>
  {(field) => (
    <field.TextField
      label={t('auth.email')}
      type="email"
      required
    />
  )}
</form.AppField>
```

**Why this matters:**

- Ensures automatic value binding
- Ensures error handling works
- Ensures validation timing works
- Ensures focus management works

### 3. Field Component Selection (STRICT RULES)

**Always choose the correct field component for the data type:**

| Data Type            | Component                           | Example                         |
| -------------------- | ----------------------------------- | ------------------------------- |
| Money/Price/Currency | `<field.PriceField>`                | amount, price, cost, value, fee |
| Quantity/Count       | `<field.QuantityField>`             | quantity, stock, count          |
| Email                | `<field.TextField type="email">`    | email                           |
| Phone                | `<field.TextField type="tel">`      | phone, mobile                   |
| Password             | `<field.PasswordField>`             | password, confirmPassword       |
| Long text            | `<field.TextareaField>`             | description, note, comment      |
| Single choice        | `<field.SelectField>`               | category, status, country       |
| Multiple choice      | `<field.SelectField multiSelect>`   | tags, permissions               |
| Yes/No               | `<field.ToggleField>`               | enabled, active, isRecurring    |
| Date                 | `<field.DateField>`                 | birthdate, startDate, dueDate   |
| Time                 | `<field.TimeField>`                 | appointmentTime, openingTime    |
| Date + Time          | `<field.DateTimeField>`             | eventDateTime, scheduledAt      |
| Date range           | `<field.DateRangeField>`            | reportPeriod, dateRange         |
| File(s)              | `<field.FileUploadField>`           | documents, attachments          |
| Image(s)             | `<field.ImageUploadField>`          | photos, gallery, avatar         |
| Customer             | `<field.CustomerSelectField>`       | customerId                      |
| Product variant      | `<field.ProductVariantSelectField>` | variantId                       |
| Address              | `<field.AddressSelectField>`        | shippingAddressId               |

**Money fields (CRITICAL):**

```tsx
// ❌ WRONG - Never use TextField for money
<field.TextField name="amount" type="number" />
<field.TextField name="price" inputMode="decimal" />

// ✅ CORRECT - Always use PriceField for money
<field.PriceField
  label={t('form.amount')}
  currencyCode={currency}
  placeholder="0.00"
  required
/>
```

**PriceField automatically:**

- Sets `inputMode="decimal"` for mobile numeric keyboard
- Sets `dir="ltr"` to keep numbers left-to-right in RTL
- Shows currency code as prefix
- Prevents invalid characters
- Formats decimal places correctly (max 2)
- Handles comma/dot decimal separator

### 4. Validation Timing (Default: onBlur)

```tsx
validators={{
  onBlur: z.string().email('validation.invalid_email'),  // Default - best UX
  onChange: /* Use for real-time (username, password strength) */
  onChangeAsync: /* Use for API checks with debounce */
}}
```

### 5. Translation Keys (MANDATORY)

**ALWAYS use translation keys from `src/i18n/*/errors.json`, NEVER hardcoded messages:**

```tsx
// ❌ WRONG - Hardcoded English message
z.string().min(8, "Password must be at least 8 characters");
z.string().email("Please enter a valid email");
z.number().positive("Must be positive");

// ✅ CORRECT - Translation keys
z.string().min(8, "validation.password_min_length");
z.string().email("validation.invalid_email");
z.number().positive("validation.positive_number");
```

**Translation key requirements:**

- Validation keys **MUST** be prefixed with `validation.` (e.g. `validation.required`, `validation.invalid_email`).
- Keys **MUST** exist in `src/i18n/en/errors.json` under `validation.*`.
- Keys **MUST** exist in `src/i18n/ar/errors.json` with Arabic translation.
- The app translates these via `src/lib/translateValidationError.ts`.
- Keys without the `validation.` prefix are treated as already-translated strings.

**Common validation keys available:**

```typescript
// Required fields
"validation.required";

// Format validation
"validation.invalid_email";
"validation.invalid_phone";
"validation.invalid_format";

// Length validation
"validation.min_length"; // interpolates {{min}}
"validation.max_length"; // interpolates {{max}}
"validation.password_min_length";

// Number validation
"validation.positive_number";
"validation.min_zero";
"validation.min_value"; // interpolates {{min}}
"validation.max_value"; // interpolates {{max}}
"validation.invalid_number";
"validation.integer_only";

// Date validation
"validation.invalid_date";
"validation.date_cannot_be_future";
"validation.date_cannot_be_past";
"validation.end_date_must_be_after_start_date";

// Selection validation
"validation.select_at_least_one";
"validation.select_too_many"; // interpolates {{max}}
```

**If you need a new validation message:**

1. Add it to `portal-web/src/i18n/en/errors.json` under `validation.*`
2. Add it to `portal-web/src/i18n/ar/errors.json` with Arabic translation
3. Use the key in your validator

### 6. Bottom Sheet Structure (STRICT)

**ALL bottom sheets with forms MUST follow this structure:**

```tsx
function MySheet({ isOpen, onClose }) {
  const { t } = useTranslation();
  const form = useKyoraForm({ /* ... */ });
  const formId = useId(); // or "my-form-id"

  return (
    <form.AppForm>  {/* MUST wrap everything */}
      <BottomSheet
        isOpen={isOpen}
        onClose={onClose}
        title={t('title')}
        footer={  {/* Footer MUST be here, not below FormRoot */}
          <div className="flex gap-2">
            <Button
              type="button"  {/* MUST be type="button" */}
              variant="ghost"
              className="flex-1"
              onClick={onClose}
              disabled={isSubmitting}
            >
              {t('common:cancel')}
            </Button>
            <form.SubmitButton
              form={formId}  {/* MUST match FormRoot id */}
              variant="primary"
              className="flex-1"
              disabled={isSubmitting}
            >
              {isSubmitting ? t('common:saving') : t('common:save')}
            </form.SubmitButton>
          </div>
        }
      >
        <form.FormRoot id={formId} className="space-y-4">
          {/* Fields go here */}
        </form.FormRoot>
      </BottomSheet>
    </form.AppForm>
  );
}
```

**Rules:**

- `form.AppForm` MUST wrap the entire BottomSheet (including footer)
- `footer` prop MUST be used for buttons (not after FormRoot)
- Cancel button MUST have `type="button"` to prevent form submission
- `form.SubmitButton` MUST have `form={formId}` matching `FormRoot` id
- Both buttons MUST be disabled while `isSubmitting`
- FormRoot MUST have the `id` prop matching SubmitButton's `form` prop

---

## Quick Start

### Basic Form (3 Steps)

```tsx
import { useKyoraForm } from "@/lib/form";
import { z } from "zod";

function LoginForm() {
  const { t } = useTranslation();

  // 1. Create form
  const form = useKyoraForm({
    defaultValues: {
      email: "",
      password: "",
    },
    onSubmit: async ({ value }) => {
      await api.login(value);
    },
  });

  // 2. Wrap in AppForm + FormRoot
  return (
    <form.AppForm>
      <form.FormRoot className="space-y-4">
        <form.FormError />

        {/* 3. Add fields with AppField pattern */}
        <form.AppField
          name="email"
          validators={{
            onBlur: z.string().email("invalid_email"),
          }}
        >
          {(field) => (
            <field.TextField
              type="email"
              label={t("auth.email")}
              autoComplete="email"
            />
          )}
        </form.AppField>

        <form.AppField
          name="password"
          validators={{
            onBlur: z.string().min(8, "password_too_short"),
          }}
        >
          {(field) => (
            <field.PasswordField
              label={t("auth.password")}
              autoComplete="current-password"
            />
          )}
        </form.AppField>

        <form.SubmitButton variant="primary">
          {t("auth.login")}
        </form.SubmitButton>
      </form.FormRoot>
    </form.AppForm>
  );
}
```

---

## Field Components Reference

## Mobile Keyboard UX (Required)

Kyora Portal is mobile-first. Field configuration must produce the correct keyboard and behave well in Arabic/RTL.

### Defaults

- Use `autoComplete` whenever possible (helps speed and reduces errors).
- Prefer `enterKeyHint="next"` for multi-field forms and `enterKeyHint="done"` for the last field.
- For identifiers/numbers (phone, order id, IBAN, codes): set `dir="ltr"` and use an appropriate `inputMode`.
- Avoid `type="number"` for currency; use `<field.PriceField />` (it uses `inputMode="decimal"` and `dir="ltr"`).

### Common Patterns

```tsx
// Email
<field.TextField
  type="email"
  autoComplete="email"
  inputMode="email"
  autoCapitalize="none"
  autoCorrect="off"
  spellCheck={false}
  enterKeyHint="next"
/>

// Phone (keep LTR even in Arabic UI)
<field.TextField
  type="tel"
  autoComplete="tel"
  inputMode="tel"
  dir="ltr"
  enterKeyHint="next"
/>

// Quantity / numeric code (use text + inputMode, not type=number)
<field.TextField
  type="text"
  inputMode="numeric"
  pattern="[0-9]*"
  dir="ltr"
  autoCapitalize="none"
  autoCorrect="off"
  spellCheck={false}
  enterKeyHint="done"
/>
```

### TextField

Standard text input (email, tel, text, url, search).

```tsx
<form.AppField name="email" validators={{ onBlur: z.string().email() }}>
  {(field) => (
    <field.TextField
      type="email" // text | email | url | tel | search
      label={t("auth.email")}
      placeholder={t("auth.email_placeholder")}
      autoComplete="email"
      inputMode="email"
      autoCapitalize="none"
      autoCorrect="off"
      spellCheck={false}
      enterKeyHint="next"
      required // Visual indicator
      disabled // Disable input
    />
  )}
</form.AppField>
```

**Auto-handled:** Value binding, change/blur handlers, error display, aria attributes

---

### PasswordField

Password with visibility toggle.

```tsx
<form.AppField name="password" validators={{ onBlur: z.string().min(8) }}>
  {(field) => (
    <field.PasswordField
      label={t("auth.password")}
      autoComplete="current-password"
      required
    />
  )}
</form.AppField>
```

**Features:** Eye icon toggle, translated labels, all TextField features

---

### TextareaField

Multi-line text with character counter.

```tsx
<form.AppField name="description" validators={{ onBlur: z.string().max(500) }}>
  {(field) => (
    <field.TextareaField
      label={t("common.description")}
      rows={4}
      maxLength={500}
      showCount // Show "45/500"
      required
    />
  )}
</form.AppField>
```

---

### SelectField (Single + Multi-Select)

Dropdown with search and multi-select chip UI.

**Single Select:**

```tsx
<form.AppField name="country" validators={{ onBlur: z.string().min(1) }}>
  {(field) => (
    <field.SelectField
      label={t("common.country")}
      options={[
        { value: "ae", label: t("countries.ae") },
        { value: "sa", label: t("countries.sa") },
      ]}
      searchable
      clearable
      required
    />
  )}
</form.AppField>
```

**Multi-Select with Chips:**

```tsx
<form.AppField
  name="tags"
  validators={{
    onBlur: z
      .array(z.string())
      .min(1, "select_at_least_one")
      .max(5, "select_too_many"),
  }}
>
  {(field) => (
    <field.SelectField
      label="Tags"
      options={[
        { value: "vip", label: "VIP", icon: <Star /> },
        { value: "wholesale", label: "Wholesale" },
      ]}
      multiSelect // Enable chips
      searchable
      clearable
      required
    />
  )}
</form.AppField>
```

**Validation Patterns:**

```typescript
// Min selections
z.array(z.string()).min(1, "select_at_least_one");

// Max selections
z.array(z.string()).max(5, "select_too_many");

// Unique values
z.array(z.string()).refine((arr) => new Set(arr).size === arr.length, {
  message: "duplicate_selection",
});
```

**Features:** Search/filter, chip UI, keyboard nav (Arrow keys, Backspace removes last chip), RTL support, mobile bottom sheet

---

### CheckboxField

Checkbox with label.

```tsx
<form.AppField name="acceptTerms">
  {(field) => (
    <field.CheckboxField
      label="Accept terms"
      description="I agree to the terms and conditions"
      required
    />
  )}
</form.AppField>
```

---

### RadioField

Radio button group.

```tsx
<form.AppField name="plan">
  {(field) => (
    <field.RadioField
      label="Select Plan"
      options={[
        { value: "free", label: "Free", description: "$0/month" },
        { value: "pro", label: "Pro", description: "$10/month" },
      ]}
      orientation="vertical" // vertical | horizontal
      required
    />
  )}
</form.AppField>
```

---

### ToggleField

Toggle switch.

```tsx
<form.AppField name="notifications">
  {(field) => (
    <field.ToggleField
      label="Enable notifications"
      description="Receive email updates"
    />
  )}
</form.AppField>
```

---

### DateField

Date picker with calendar.

```tsx
<form.AppField
  name="birthdate"
  validators={{
    onBlur: z
      .date()
      .max(new Date(), "date_cannot_be_future")
      .refine((date) => new Date().getFullYear() - date.getFullYear() >= 18, {
        message: "must_be_18_or_older",
      }),
  }}
>
  {(field) => (
    <field.DateField
      label="Birth Date"
      minAge={18}
      maxDate={new Date()}
      disableWeekends
      clearable
      required
    />
  )}
</form.AppField>
```

**Features:** Calendar popup, RTL support, month/year nav, keyboard nav, mobile full-screen modal

---

### TimeField

Time picker with hour/minute controls.

```tsx
<form.AppField
  name="appointmentTime"
  validators={{
    onBlur: z.date().refine(
      (date) => {
        const hours = date.getHours();
        return hours >= 9 && hours < 17;
      },
      { message: "outside_business_hours" },
    ),
  }}
>
  {(field) => (
    <field.TimeField
      label="Appointment Time"
      use24Hour={false} // 12-hour with AM/PM
      minuteStep={15} // 15-min increments
      clearable
      required
    />
  )}
</form.AppField>
```

**Features:** Numeric inputs, AM/PM toggle, arrow buttons, keyboard nav, auto-advance

---

### DateTimeField

Combined date + time picker.

```tsx
<form.AppField
  name="eventDateTime"
  validators={{
    onBlur: z.date().min(new Date(), "must_be_future"),
  }}
>
  {(field) => (
    <field.DateTimeField
      mode="datetime" // 'date' | 'time' | 'datetime'
      label="Event Date & Time"
      datePickerProps={{
        minDate: new Date(),
        disableWeekends: true,
      }}
      timePickerProps={{
        minuteStep: 30,
        use24Hour: false,
      }}
      required
    />
  )}
</form.AppField>
```

---

### DateRangeField

Dual-calendar date range picker.

```tsx
import type { DateRange } from 'react-day-picker'

<form.AppField
  name="reportDateRange"
  validators={{
    onBlur: z.custom<DateRange>((val) => {
      const range = val as DateRange
      if (!range?.from || !range?.to) return 'date_range_required'
      if (range.to < range.from) return 'date_range_invalid'
      return undefined
    }),
  }}
>
  {(field) => (
    <field.DateRangeField
      label="Report Period"
      numberOfMonths={2}              // Show 2 months side-by-side
      minDate={new Date('2020-01-01')}
      maxDate={new Date()}
      disabledDates={[...]}
      clearable
      required
    />
  )}
</form.AppField>
```

**Validation Patterns:**

```typescript
// Maximum range (90 days)
z.custom<DateRange>((val) => {
  const range = val as DateRange;
  if (!range?.from || !range?.to) return "date_range_required";
  const days = Math.ceil(
    (range.to.getTime() - range.from.getTime()) / (1000 * 60 * 60 * 24),
  );
  if (days > 90) return "maximum_90_days";
  return undefined;
});

// Weekdays only
z.custom<DateRange>((val) => {
  const range = val as DateRange;
  const isWeekday = (date: Date) => ![0, 6].includes(date.getDay());
  if (!isWeekday(range.from) || !isWeekday(range.to)) return "weekdays_only";
  return undefined;
});
```

---

### FieldArray

Dynamic array with drag-drop reordering.

```tsx
<form.AppField
  name="phoneNumbers"
  validators={{
    onChange: ({ value }) => {
      return validateArrayAnd(value, [
        (array) => validateArrayLength(array, { min: 1, max: 5 }),
        (array) =>
          validateUniqueValues(array, {
            extractor: (item) => item.number,
            errorKey: "form.duplicatePhones",
          }),
      ]);
    },
  }}
>
  {(field) => (
    <field.FieldArray
      label="Phone Numbers"
      addButtonLabel="Add Phone"
      minItems={1}
      maxItems={5}
      reorderable
      defaultValue={{ number: "", type: "mobile" }}
      render={(item, operations, index) => (
        <div className="flex gap-2">
          <input
            type="tel"
            value={item.number}
            onChange={(e) => {
              const updated = [...field.state.value];
              updated[index] = { ...item, number: e.target.value };
              field.handleChange(updated);
            }}
            className="input flex-1"
          />
          <select
            value={item.type}
            onChange={(e) => {
              const updated = [...field.state.value];
              updated[index] = { ...item, type: e.target.value };
              field.handleChange(updated);
            }}
            className="select"
          >
            <option value="mobile">Mobile</option>
            <option value="work">Work</option>
          </select>
        </div>
      )}
    />
  )}
</form.AppField>
```

**Array Validation Utilities:**

```typescript
import {
  validateArrayLength,
  validateUniqueValues,
  validateArrayItems,
  validateNoOverlap,
  validateArrayAnd,
  validateArrayOr,
  validateArrayCount,
} from "@/lib/form";

// Min/max items
validateArrayLength(value, { min: 1, max: 10 });

// Unique values
validateUniqueValues(value, {
  extractor: (item) => item.email,
  errorKey: "form.duplicateEmails",
});

// Per-item validation
validateArrayItems(value, (item, index) => {
  if (!item.name) return "form.nameRequired";
  return undefined;
});

// Overlap validation (time ranges)
validateNoOverlap(value, {
  extractor: (item) => ({ start: item.startTime, end: item.endTime }),
  errorKey: "form.overlappingTimeRanges",
});

// Combine validators (AND)
validateArrayAnd(value, [validator1, validator2, validator3]);

// Combine validators (OR - at least one must pass)
validateArrayOr(value, [validator1, validator2]);

// Count validation (exactly one primary)
validateArrayCount(value, {
  extractor: (item) => item.isPrimary,
  matchValue: true,
  exactCount: 1,
  errorKey: "form.exactlyOnePrimary",
});
```

**Features:** Drag-drop reordering (via @dnd-kit), add/remove, move up/down buttons, empty state, max items warning, smooth animations, RTL support, keyboard nav, touch-optimized

---

### FileUploadField

Generic file upload with drag-drop and progress.

```tsx
import { BusinessContext } from "@/contexts/BusinessContext";

<BusinessContext.Provider value={business.descriptor}>
  <form.AppForm>
    <form.AppField
      name="documents"
      validators={{
        onBlur: fileSchema({
          maxSize: "5MB",
          maxFiles: 3,
        }),
      }}
    >
      {(field) => (
        <field.FileUploadField
          label="Documents"
          accept=".pdf,.doc,.docx"
          maxFiles={3}
          maxSize="5MB"
          multiple
        />
      )}
    </form.AppField>
  </form.AppForm>
</BusinessContext.Provider>;
```

**Features:** Optimistic UI, progress tracking, mobile camera, thumbnail generation, concurrent uploads, error recovery, drag-drop, accessible

**Business Context Required:** Upload system needs `businessDescriptor` from context

---

### ImageUploadField

Specialized image-only upload.

```tsx
// Single image
<form.AppField
  name="avatar"
  validators={{
    onBlur: imageSchema({ maxSize: '2MB' }),
  }}
>
  {(field) => (
    <field.ImageUploadField
      label="Avatar"
      single
      maxSize="2MB"
    />
  )}
</form.AppField>

// Multiple images (gallery)
<form.AppField
  name="gallery"
  validators={{
    onBlur: imageSchema({ maxSize: '10MB', maxFiles: 20 }),
  }}
>
  {(field) => (
    <field.ImageUploadField
      label="Product Gallery"
      maxFiles={20}
      reorderable
    />
  )}
</form.AppField>
```

**Preset Schemas:**

```typescript
import {
  businessLogoSchema,
  productPhotosSchema,
  variantPhotosSchema,
} from "@/schemas/upload";

// Business logo: single, 2MB
businessLogoSchema();

// Product photos: 2-10 images, 10MB each
productPhotosSchema();

// Variant photos: 1-5 images, 10MB each
variantPhotosSchema();
```

---

## Validation Patterns

### Timing Strategy

| Timing             | Use Case           | Example                                      |
| ------------------ | ------------------ | -------------------------------------------- |
| `onBlur`           | Default (best UX)  | Email, password, text fields                 |
| `onChange`         | Real-time feedback | Username availability, password strength     |
| `onChangeAsync`    | API checks         | Username exists, promo code validation       |
| `onChangeListenTo` | Cross-field        | Password confirmation, end date > start date |

### Common Patterns

**Email:**

```tsx
validators={{
  onBlur: z.string().email('invalid_email'),
}}
```

**Password:**

```tsx
validators={{
  onBlur: z.string().min(8, 'password_too_short'),
}}
```

**Password Confirmation:**

```tsx
<form.AppField
  name="confirmPassword"
  validators={{
    onChangeListenTo: ["password"],
    onBlur: ({ value, fieldApi }) => {
      if (value !== fieldApi.form.getFieldValue("password")) {
        return "passwords_must_match";
      }
      return undefined;
    },
  }}
/>
```

**Username Availability (Async + Debounce):**

```tsx
validators={{
  onChangeAsync: async ({ value }) => {
    if (value.length < 3) return 'username_too_short'
    const exists = await checkUsernameExists(value)
    return exists ? 'username_taken' : undefined
  },
  onChangeAsyncDebounceMs: 500,
}}
```

**Date Range (Start < End):**

```tsx
<form.AppField
  name="endDate"
  validators={{
    onChangeListenTo: ["startDate"],
    onBlur: ({ value, fieldApi }) => {
      const startDate = fieldApi.form.getFieldValue("startDate");
      if (startDate && value && value <= startDate) {
        return "end_date_must_be_after_start_date";
      }
      return undefined;
    },
  }}
/>
```

**Conditional Validation:**

```tsx
validators={{
  onChange: ({ value, fieldApi }) => {
    const requirePhone = fieldApi.form.getFieldValue('hasPhoneSupport')
    if (requirePhone && !value) {
      return 'phone_required'
    }
    return undefined
  },
}}
```

---

## Advanced Patterns

### Modal/Sheet with External Submit

**CRITICAL:** `form.AppForm` must wrap the ENTIRE component, including footer.

**Strict pattern (MUST follow exactly):**

```tsx
function AddCustomerSheet({ isOpen, onClose }) {
  const { t } = useTranslation();
  const formId = useId(); // or "add-customer-form"
  const form = useKyoraForm({
    defaultValues: {
      /* ... */
    },
    onSubmit: async ({ value }) => {
      await api.createCustomer(value);
      onClose();
    },
  });

  const isSubmitting = createMutation.isPending;

  // ✅ CORRECT - AppForm wraps everything including footer
  return (
    <form.AppForm>
      <BottomSheet
        isOpen={isOpen}
        onClose={onClose}
        title={t("customers:add_customer")}
        closeOnOverlayClick={!isSubmitting}
        closeOnEscape={!isSubmitting}
        footer={
          <div className="flex gap-2">
            <Button
              type="button"
              variant="ghost"
              className="flex-1"
              onClick={onClose}
              disabled={isSubmitting}
            >
              {t("common:cancel")}
            </Button>
            <form.SubmitButton
              form={formId}
              variant="primary"
              className="flex-1"
              disabled={isSubmitting}
            >
              {isSubmitting ? t("common:saving") : t("common:save")}
            </form.SubmitButton>
          </div>
        }
      >
        <form.FormRoot id={formId} className="space-y-4">
          {/* fields go here */}
        </form.FormRoot>
      </BottomSheet>
    </form.AppForm>
  );
}
```

**Common mistakes to avoid:**

```tsx
// ❌ WRONG - AppForm doesn't wrap footer
<>
  <form.AppForm>
    <BottomSheet>
      <form.FormRoot id={formId}>{/* fields */}</form.FormRoot>
    </BottomSheet>
  </form.AppForm>
  <div>
    <button onClick={onClose}>Cancel</button>
    <form.SubmitButton form={formId}>Submit</form.SubmitButton>
  </div>
</>

// ❌ WRONG - footer not used, buttons after FormRoot
<form.AppForm>
  <BottomSheet>
    <form.FormRoot id={formId}>{/* fields */}</form.FormRoot>
    <div>
      <button onClick={onClose}>Cancel</button>
      <form.SubmitButton form={formId}>Submit</form.SubmitButton>
    </div>
  </BottomSheet>
</form.AppForm>

// ❌ WRONG - cancel button missing type="button"
<Button onClick={onClose}>{t('cancel')}</Button>
// This triggers form submission in some browsers!

// ❌ WRONG - form id mismatch
<form.FormRoot id="my-form">{/* fields */}</form.FormRoot>
...
<form.SubmitButton form="different-form">Submit</form.SubmitButton>
```

### Dependent Fields (Auto-Update)

```tsx
function CustomerForm() {
  const [selectedCountry, setSelectedCountry] = useState("US");
  const form = useKyoraForm({
    /* ... */
  });

  // Auto-link country to phone code
  useEffect(() => {
    const country = countries.find((c) => c.code === selectedCountry);
    if (country?.phonePrefix) {
      form.setFieldValue("phoneCode", country.phonePrefix);
    }
  }, [selectedCountry]);

  return (
    <form.FormRoot>
      <form.Field name="countryCode">
        {(field) => (
          <CountrySelect
            value={field.state.value}
            onChange={(value) => {
              field.handleChange(value);
              setSelectedCountry(value);
            }}
          />
        )}
      </form.Field>

      <form.Field name="phoneCode">
        {(field) => (
          <PhoneCodeSelect
            value={field.state.value}
            onChange={field.handleChange}
            countryCode={selectedCountry}
          />
        )}
      </form.Field>
    </form.FormRoot>
  );
}
```

### Uncontrolled Components

For components that manage their own state:

```tsx
<form.Subscribe selector={(state) => state.values}>
  {(values) => (
    <SocialMediaInputs
      instagram={values.instagram}
      onInstagramChange={(value) => form.setFieldValue("instagram", value)}
      facebook={values.facebook}
      onFacebookChange={(value) => form.setFieldValue("facebook", value)}
    />
  )}
</form.Subscribe>
```

---

## Server Errors

Inject server errors using `createServerErrorValidator`:

```tsx
import { createServerErrorValidator } from "@/lib/form";

const form = useKyoraForm({
  validators: {
    email: {
      onBlur: z.string().email("invalid_email"),
      onServer: createServerErrorValidator(), // Injects RFC7807 errors
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

1. API returns RFC7807 error with field details
2. Server error validator injects errors into form
3. Errors appear on corresponding fields
4. Form-level errors shown in `<form.FormError />`

---

## Performance

### Subscription Pattern

**❌ Bad:** Re-renders on every state change

```tsx
const values = form.useStore((state) => state.values);
```

**✅ Good:** Only re-renders when values change

```tsx
<form.Subscribe selector={(state) => state.values}>
  {(values) => <div>{values.email}</div>}
</form.Subscribe>
```

### Field Isolation

Each `form.Field` is isolated. Changing one field doesn't re-render others.

### Validation Debouncing

```tsx
validators={{
  onChangeAsyncDebounceMs: 500,  // Debounce async validation
  onChangeAsync: async ({ value }) => {
    // Only fires after 500ms of no changes
  },
}}
```

---

## Troubleshooting

### Issue: "formContext only works when..."

**Cause:** Component using `useFormContext()` is outside `<form.AppForm>`

**Solution:** Wrap ENTIRE component tree in `<form.AppForm>`, including footer buttons

### Issue: Errors not translating

**Cause:** Translation key missing in `src/i18n/*/errors.json`

**Solution:** Add key to errors.json with translated message under `validation.*`

### Issue: Form not submitting

**Cause:** `<form.FormRoot>` missing correct `id` matching `<form.SubmitButton form="...">`

**Solution:** Use `useId()` and pass same value to both

### Issue: Field not validating

**Cause:** Validator in wrong place or not using AppField pattern

**Solution:** Use `<form.AppField>` with validators prop, structure as `validators: { onBlur: ... }`

### Issue: "businessDescriptor is required"

**Cause:** FileUploadField/ImageUploadField used without BusinessContext

**Solution:** Wrap component in `<BusinessContext.Provider value={business.descriptor}>`

### Issue: Money field showing wrong keyboard on mobile

**Cause:** Using TextField instead of PriceField

**Solution:** Use `<field.PriceField>` which automatically sets `inputMode=\"decimal\"` and `dir=\"ltr\"`

### Issue: Cancel button submitting form

**Cause:** Missing `type=\"button\"` attribute

**Solution:** Always add `type=\"button\"` to cancel/secondary buttons in forms

### Issue: Form state not updating

**Cause:** Using components directly instead of AppField pattern

**Solution:** Replace `<TextField />` with `<form.AppField>{(field) => <field.TextField />}</form.AppField>`

---

## Common Pitfalls & How to Avoid Them

### Pitfall 1: Using form components directly

```tsx
// ❌ WRONG - Component used directly
<TextField
  name=\"email\"
  value={email}
  onChange={(e) => setEmail(e.target.value)}
  label={t('auth.email')}
/>

// ❌ WRONG - PriceInput used directly
<PriceInput
  name=\"amount\"
  value={amount}
  onChange={(e) => setAmount(e.target.value)}
/>

// ✅ CORRECT - Always use AppField pattern
<form.AppField name=\"email\" validators={{ onBlur: z.string().email('validation.invalid_email') }}>
  {(field) => (
    <field.TextField
      label={t('auth.email')}
      type=\"email\"
    />
  )}
</form.AppField>

<form.AppField name=\"amount\" validators={{ onChange: ({ value }) => { /* ... */ } }}>
  {(field) => (
    <field.PriceField
      label={t('form.amount')}
      currencyCode={currency}
    />
  )}
</form.AppField>
```

**Why:** Direct usage bypasses form context, validation, error handling, and focus management.

### Pitfall 2: Wrong component for money fields

```tsx
// ❌ WRONG - TextField with type=\"number\"
<form.AppField name=\"amount\">
  {(field) => (
    <field.TextField
      type=\"number\"
      step=\"0.01\"
      min=\"0\"
    />
  )}
</form.AppField>

// ❌ WRONG - TextField with inputMode
<form.AppField name=\"price\">
  {(field) => (
    <field.TextField
      inputMode=\"decimal\"
      pattern=\"[0-9]*\"
    />
  )}
</form.AppField>

// ✅ CORRECT - PriceField for all money values
<form.AppField name=\"amount\">
  {(field) => (
    <field.PriceField
      label={t('form.amount')}
      currencyCode={currency}
      placeholder=\"0.00\"
    />
  )}
</form.AppField>
```

**Why:** PriceField handles decimal formatting, currency display, mobile keyboard, RTL, and validation automatically.

**Money field detection:**

- If the field name contains: `amount`, `price`, `cost`, `fee`, `value`, `total`, `subtotal`, `vat`, `discount`, `shipping`
- If the label refers to money/currency
- If the value represents money

**Always use PriceField.**

### Pitfall 3: Hardcoded validation messages

```tsx
// ❌ WRONG - English hardcoded
validators={{
  onBlur: z.string().min(8, \"Password must be at least 8 characters\"),
}}

validators={{
  onChange: ({ value }) => {
    if (!value) return \"This field is required\";
    if (value.length < 3) return \"Must be at least 3 characters\";
    return undefined;
  },
}}

// ✅ CORRECT - Translation keys
validators={{
  onBlur: z.string().min(8, \"validation.password_min_length\"),
}}

validators={{
  onChange: ({ value }) => {
    if (!value) return \"validation.required\";
    if (value.length < 3) return \"validation.min_length\"; // with interpolation
    return undefined;
  },
}}
```

**Why:** Hardcoded messages break i18n, ignore RTL, and create maintenance burden.

### Pitfall 4: BottomSheet structure mistakes

```tsx
// ❌ WRONG - AppForm doesn't wrap footer
<form.AppForm>
  <BottomSheet>
    <form.FormRoot id=\"my-form\">
      {/* fields */}
    </form.FormRoot>
  </BottomSheet>
</form.AppForm>
<div>
  <Button onClick={onClose}>Cancel</Button>
  <form.SubmitButton form=\"my-form\">Save</form.SubmitButton>
</div>

// ❌ WRONG - Footer buttons after FormRoot
<form.AppForm>
  <BottomSheet>
    <form.FormRoot id=\"my-form\">
      {/* fields */}
    </form.FormRoot>
    <div className=\"flex gap-2\">
      <Button onClick={onClose}>Cancel</Button>
      <form.SubmitButton form=\"my-form\">Save</form.SubmitButton>
    </div>
  </BottomSheet>
</form.AppForm>

// ✅ CORRECT - AppForm wraps everything, footer is a prop
<form.AppForm>
  <BottomSheet
    footer={
      <div className=\"flex gap-2\">
        <Button type=\"button\" variant=\"ghost\" onClick={onClose}>
          {t('common:cancel')}
        </Button>
        <form.SubmitButton form=\"my-form\" variant=\"primary\">
          {t('common:save')}
        </form.SubmitButton>
      </div>
    }
  >
    <form.FormRoot id=\"my-form\">
      {/* fields */}
    </form.FormRoot>
  </BottomSheet>
</form.AppForm>
```

**Why:** SubmitButton needs form context, which only exists inside AppForm. Footer must be inside AppForm scope.

### Pitfall 5: Cancel button submitting form

```tsx
// ❌ WRONG - Missing type=\"button\"
<Button onClick={onClose}>Cancel</Button>

// ❌ WRONG - Using form.SubmitButton
<form.SubmitButton onClick={onClose}>Cancel</form.SubmitButton>

// ✅ CORRECT - type=\"button\" prevents submission
<Button type=\"button\" onClick={onClose}>
  {t('common:cancel')}
</Button>
```

**Why:** Buttons inside forms default to `type=\"submit\"` and will trigger form submission when clicked.

### Pitfall 6: Missing validation keys

```tsx
// ❌ WRONG - Key doesn't exist
validators={{
  onBlur: z.string().email('validation.email_format_invalid'),
}}
// Result: \"validation.email_format_invalid\" shown to user (not translated)

// ✅ CORRECT - Use existing key
validators={{
  onBlur: z.string().email('validation.invalid_email'),
}}
// Result: \"Please enter a valid email address\" (EN) or \"الرجاء إدخال بريد إلكتروني صحيح\" (AR)
```

**How to check:**

1. Search `src/i18n/en/errors.json` for the key
2. Verify it exists under `validation.*`
3. Verify Arabic translation exists in `src/i18n/ar/errors.json`

### Pitfall 7: Wrong validation timing

```tsx
// ❌ WRONG - Using onChange for all fields (too aggressive)
<form.AppField
  name=\"email\"
  validators={{
    onChange: z.string().email('validation.invalid_email'),
  }}
>
  {(field) => <field.TextField type=\"email\" />}
</form.AppField>
// Result: Error shows while user is still typing

// ✅ CORRECT - Use onBlur by default
<form.AppField
  name=\"email\"
  validators={{
    onBlur: z.string().email('validation.invalid_email'),
  }}
>
  {(field) => <field.TextField type=\"email\" />}
</form.AppField>
// Result: Error shows after user leaves field

// ✅ CORRECT - Use onChange only when needed (real-time feedback)
<form.AppField
  name=\"username\"
  validators={{
    onChange: z.string().min(3, 'validation.min_length'),
    onChangeAsync: async ({ value }) => {
      const exists = await checkUsername(value);
      return exists ? 'validation.username_taken' : undefined;
    },
    onChangeAsyncDebounceMs: 500,
  }}
>
  {(field) => <field.TextField />}
</form.AppField>
```

**When to use each:**

- `onBlur`: Default for most fields (best UX)
- `onChange`: Real-time feedback (username availability, password strength)
- `onChangeAsync`: API validation with debounce

### Pitfall 8: Form ID mismatch

```tsx
// ❌ WRONG - IDs don't match
<form.FormRoot id=\"create-form\">
  {/* fields */}
</form.FormRoot>
...
<form.SubmitButton form=\"submit-form\">Save</form.SubmitButton>

// ✅ CORRECT - Matching IDs
const formId = useId(); // or \"create-expense-form\"
...
<form.FormRoot id={formId}>
  {/* fields */}
</form.FormRoot>
...
<form.SubmitButton form={formId}>Save</form.SubmitButton>
```

**Why:** SubmitButton needs to know which form to submit via the `form` attribute.

---

## Agent Validation Checklist (MANDATORY BEFORE COMPLETING TASK)

**Before marking a form task as complete, verify ALL of these:**

### Structure & Context

- ☑ `<form.AppForm>` wraps EVERYTHING (including BottomSheet footer if present)
- ☑ All form components (FormRoot, SubmitButton, FormError, Subscribe) are inside `<form.AppForm>`
- ☑ FormRoot has an `id` prop
- ☑ SubmitButton has `form` prop matching FormRoot's `id`
- ☑ No form components are used outside `<form.AppForm>`

### Field Pattern

- ☑ **ZERO direct component usage** - every field uses `<form.AppField>` + `{(field) => <field.Component />}` pattern
- ☑ **NO** `<TextField .../>`, `<SelectField .../>`, `<PriceInput .../>`, etc. used directly
- ☑ All fields use the render prop pattern: `{(field) => <field.ComponentName />}`

### Field Component Selection

- ☑ Money/price fields use `<field.PriceField>` (NOT TextField with type=\"number\")
- ☑ Quantity fields use `<field.QuantityField>`
- ☑ Email fields use `<field.TextField type=\"email\">`
- ☑ Password fields use `<field.PasswordField>`
- ☑ Long text uses `<field.TextareaField>` (NOT TextField)
- ☑ Single select uses `<field.SelectField>`
- ☑ Multi-select uses `<field.SelectField multiSelect>`
- ☑ Dates use `<field.DateField>`
- ☑ Times use `<field.TimeField>`
- ☑ Images use `<field.ImageUploadField>`
- ☑ Files use `<field.FileUploadField>`
- ☑ Toggles/switches use `<field.ToggleField>`
- ☑ Radio groups use `<field.RadioField>`

### Validation

- ☑ All validators use `onBlur` by default (exceptions must be documented)
- ☑ **ZERO hardcoded error messages** - all use `validation.*` keys
- ☑ All validation keys exist in `src/i18n/en/errors.json` under `validation.*`
- ☑ All validation keys exist in `src/i18n/ar/errors.json` with Arabic translation
- ☑ Cross-field validation uses `onChangeListenTo` or `onBlurListenTo`
- ☑ Async validation uses `onChangeAsync` with `onChangeAsyncDebounceMs`

### Bottom Sheet (if applicable)

- ☑ `form.AppForm` wraps the entire BottomSheet (including footer)
- ☑ Footer buttons are in the `footer` prop (NOT after FormRoot)
- ☑ Cancel button has `type=\"button\"` to prevent form submission
- ☑ Cancel button is disabled during submission
- ☑ Submit button is disabled during submission
- ☑ Both buttons use `className=\"flex-1\"` for equal width
- ☑ Submit button shows loading state (e.g., `{isSubmitting ? t('saving') : t('save')}`)

### Translation & i18n

- ☑ All labels use `t()` function, no hardcoded text
- ☑ All validation messages use translation keys
- ☑ Button labels use translation keys
- ☑ Placeholder text uses translation keys
- ☑ Helper text uses translation keys

### TypeScript & Types

- ☑ No `any` types in validators or handlers
- ☑ Form values are properly typed
- ☑ Field names match form value types
- ☑ Validator functions have proper return types

### Accessibility

- ☑ All fields have labels (via `label` prop)
- ☑ Required fields marked with `required` prop
- ☑ Icon-only buttons have `aria-label`
- ☑ Error messages are associated with fields

### RTL Support

- ☑ No `left`/`right` classes used (use `start`/`end` or `ms`/`me`)
- ☑ LTR-only fields (phone, codes, IDs) have `dir=\"ltr\"`
- ☑ Money fields automatically set `dir=\"ltr\"` (via PriceField)

### File Uploads (if applicable)

- ☑ FileUploadField/ImageUploadField wrapped in `<BusinessContext.Provider>`
- ☑ Business descriptor is passed to context
- ☑ File validation schema used (e.g., `imageSchema()`, `fileSchema()`)

### Performance

- ☑ Avoid `form.useStore()` in components (use `<form.Subscribe>` instead)
- ☑ Async validators have debounce configured
- ☑ Large lists use field array optimizations

### Common Mistakes to Check

- ☑ NOT using `<TextField name=\"amount\" type=\"number\">` for money
- ☑ NOT using `<input>` or other native elements directly
- ☑ NOT hardcoding English error messages
- ☑ NOT placing SubmitButton outside form.AppForm
- ☑ NOT forgetting `type=\"button\"` on cancel buttons
- ☑ NOT using wrong field component for data type

**If ANY checkbox is unchecked, the form is NOT complete. Fix it before finishing.**

---

## See Also

- **UI Components:** `.github/instructions/ui-implementation.instructions.md` → daisyUI usage, RTL rules
- **HTTP Requests:** `.github/instructions/ky.instructions.md` → Form submission patterns
- **File Uploads:** `.github/instructions/asset_upload.instructions.md` → Backend contract
- **Design Tokens:** `.github/instructions/design-tokens.instructions.md` → Colors, typography

---

## Resources

- TanStack Form Docs: https://tanstack.com/form/latest
- Zod Docs: https://zod.dev
- Implementation: `portal-web/src/lib/form/`
- Examples: `portal-web/src/routes/onboarding/`, `portal-web/src/routes/auth/`
