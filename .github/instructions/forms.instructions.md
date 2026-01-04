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

## ⚠️ CRITICAL RULES

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

### 2. Field Pattern

**Always use `<form.AppField>` + `{(field) => <field.ComponentName />}` pattern:**

```tsx
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

### 3. Validation Timing (Default: onBlur)

```tsx
validators={{
  onBlur: z.string().email('invalid_email'),  // Default - best UX
  onChange: /* Use for real-time (username, password strength) */
  onChangeAsync: /* Use for API checks with debounce */
}}
```

### 4. Translation Keys

**Always use translation keys, never hardcoded messages:**

```tsx
// ❌ WRONG
z.string().min(8, "Password must be at least 8 characters");

// ✅ CORRECT
z.string().min(8, "validation.password_min_length");
// Translation: src/i18n/*/errors.json → { "validation": { "password_min_length": "..." } }
```

**Important (Prevents Missing-Key Bugs):**

- Validation keys must be prefixed with `validation.` (e.g. `validation.required`, `validation.invalid_email`).
- The app translates these via `src/lib/translateValidationError.ts`.
- Keys without the `validation.` prefix are treated as already-translated strings.

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
      { message: "outside_business_hours" }
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
    (range.to.getTime() - range.from.getTime()) / (1000 * 60 * 60 * 24)
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

**CRITICAL:** `form.AppForm` must wrap the ENTIRE component, including footer:

```tsx
function AddCustomerSheet() {
  const formId = useId();
  const form = useKyoraForm({
    /* ... */
  });

  // ✅ CORRECT - AppForm wraps everything
  return (
    <form.AppForm>
      <BottomSheet
        footer={
          <div>
            <button onClick={onClose}>Cancel</button>
            <form.SubmitButton form={formId}>Submit</form.SubmitButton>
          </div>
        }
      >
        <form.FormRoot id={formId}>{/* fields */}</form.FormRoot>
      </BottomSheet>
    </form.AppForm>
  );
}
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

**Solution:** Add key to errors.json with translated message

### Issue: Form not submitting

**Cause:** `<form.FormRoot>` missing correct `id` matching `<form.SubmitButton form="...">`

**Solution:** Use `useId()` and pass same value to both

### Issue: Field not validating

**Cause:** Validator in wrong place (top-level `onBlur` instead of `validators.email.onBlur`)

**Solution:** Structure validators as `validators: { fieldName: { onBlur: ... } }`

### Issue: "businessDescriptor is required"

**Cause:** FileUploadField/ImageUploadField used without BusinessContext

**Solution:** Wrap component in `<BusinessContext.Provider value={business.descriptor}>`

---

## Agent Validation Checklist

Before completing form task:

- ☑ All form components inside `<form.AppForm>`
- ☑ All fields use `<form.AppField>` + `{(field) => <field.Component />}` pattern
- ☑ Validators use `onBlur` by default (or documented exception)
- ☑ All error messages are translation keys, not hardcoded strings
- ☑ Translation keys exist in `src/i18n/*/errors.json`
- ☑ Cross-field validation uses `onChangeListenTo` or `onBlurListenTo`
- ☑ File uploads wrapped in `<BusinessContext.Provider>`
- ☑ No `any` types in validators or handlers
- ☑ RTL: No `left`/`right` classes (use `start`/`end`)
- ☑ Accessibility: `aria-label` on icon-only buttons

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
