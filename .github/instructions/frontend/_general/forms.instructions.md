---
description: Frontend forms system core - TanStack Form, field patterns, submission, error display (reusable across portal-web, storefront-web)
applyTo: "portal-web/**,storefront-web/**"
---

# Frontend Forms System

TanStack Form-based system with zero-boilerplate field composition.

**Cross-refs:**

- Validation: `./forms-validation.instructions.md`
- UI components: `./ui-patterns.instructions.md`
- HTTP: `./http-client.instructions.md`

---

## 1. Core Pattern

### useKyoraForm Hook

```tsx
import { useKyoraForm } from "@/lib/form";

const form = useKyoraForm({
  defaultValues: {
    email: "",
    password: "",
  },
  onSubmit: async ({ value }) => {
    await api.login(value);
  },
});
```

### Form Structure

```tsx
<form.AppForm>
  <form.FormRoot className="space-y-4">
    <form.FormError /> {/* Form-level errors */}
    {/* Fields */}
    <form.SubmitButton variant="primary">Submit</form.SubmitButton>
  </form.FormRoot>
</form.AppForm>
```

**Critical:** ALL form components must be inside `<form.AppForm>`

---

## 2. Field Pattern (Mandatory)

**NEVER use components directly. ALWAYS use AppField pattern:**

```tsx
// ❌ WRONG
<TextField name="email" label="Email" />

// ✅ CORRECT
<form.AppField name="email" validators={{ onBlur: z.string().email() }}>
  {(field) => (
    <field.TextField
      label="Email"
      type="email"
      required
    />
  )}
</form.AppField>
```

**Why:** Ensures automatic value binding, error handling, validation timing, focus management.

---

## 3. Field Components

### TextField

```tsx
<form.AppField
  name="email"
  validators={{ onBlur: z.string().email("validation.invalid_email") }}
>
  {(field) => (
    <field.TextField
      type="email" // text | email | url | tel | search
      label="Email"
      placeholder="you@example.com"
      autoComplete="email"
      inputMode="email"
      autoCapitalize="none"
      autoCorrect="off"
      spellCheck={false}
      enterKeyHint="next"
      required
      disabled
    />
  )}
</form.AppField>
```

### PasswordField

```tsx
<form.AppField name="password" validators={{ onBlur: z.string().min(8) }}>
  {(field) => (
    <field.PasswordField
      label="Password"
      autoComplete="current-password"
      required
    />
  )}
</form.AppField>
```

**Features:** Eye icon toggle, translated labels

### TextareaField

```tsx
<form.AppField name="description" validators={{ onBlur: z.string().max(500) }}>
  {(field) => (
    <field.TextareaField
      label="Description"
      rows={4}
      maxLength={500}
      showCount // Show "45/500"
      required
    />
  )}
</form.AppField>
```

### SelectField (Single)

```tsx
<form.AppField name="country" validators={{ onBlur: z.string().min(1) }}>
  {(field) => (
    <field.SelectField
      label="Country"
      options={[
        { value: "ae", label: "UAE" },
        { value: "sa", label: "Saudi Arabia" },
      ]}
      searchable
      clearable
      required
    />
  )}
</form.AppField>
```

### SelectField (Multi-Select)

```tsx
<form.AppField
  name="tags"
  validators={{
    onBlur: z.array(z.string()).min(1, "validation.select_at_least_one"),
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

**Features:** Chip UI, keyboard nav (Arrow keys, Backspace removes last chip)

### CheckboxField

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

### RadioField

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

### ToggleField

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

### DateField

```tsx
<form.AppField
  name="birthdate"
  validators={{
    onBlur: z.date().max(new Date(), "validation.date_cannot_be_future"),
  }}
>
  {(field) => (
    <field.DateField
      label="Birth Date"
      minAge={18}
      maxDate={new Date()}
      clearable
      required
    />
  )}
</form.AppField>
```

### TimeField

```tsx
<form.AppField name="appointmentTime">
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

### DateTimeField

```tsx
<form.AppField
  name="eventDateTime"
  validators={{
    onBlur: z.date().min(new Date(), "validation.must_be_future"),
  }}
>
  {(field) => (
    <field.DateTimeField
      mode="datetime" // 'date' | 'time' | 'datetime'
      label="Event Date & Time"
      datePickerProps={{ minDate: new Date() }}
      timePickerProps={{ minuteStep: 30 }}
      required
    />
  )}
</form.AppField>
```

### DateRangeField

```tsx
<form.AppField
  name="reportDateRange"
  validators={{
    onBlur: z.custom<DateRange>((val) => {
      const range = val as DateRange;
      if (!range?.from || !range?.to) return "validation.date_range_required";
      return undefined;
    }),
  }}
>
  {(field) => (
    <field.DateRangeField
      label="Report Period"
      numberOfMonths={2}
      minDate={new Date("2020-01-01")}
      maxDate={new Date()}
      clearable
      required
    />
  )}
</form.AppField>
```

---

## 4. Validation Timing

```tsx
validators={{
  onBlur: z.string().email('validation.invalid_email'),     // Default - best UX
  onChange: /* Use for real-time (username, password strength) */
  onChangeAsync: /* Use for API checks with debounce */
}}
```

**Default:** `onBlur` (validates after user leaves field)

---

## 5. Mobile Keyboard UX

### Email

```tsx
<field.TextField
  type="email"
  autoComplete="email"
  inputMode="email"
  autoCapitalize="none"
  autoCorrect="off"
  spellCheck={false}
  enterKeyHint="next"
/>
```

### Phone

```tsx
<field.TextField
  type="tel"
  autoComplete="tel"
  inputMode="tel"
  dir="ltr" // Keep LTR in RTL UI
  enterKeyHint="next"
/>
```

### Quantity / Numeric Code

```tsx
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

---

## 6. Form Submission

### Basic

```tsx
const form = useKyoraForm({
  defaultValues: { email: "", password: "" },
  onSubmit: async ({ value }) => {
    await api.login(value);
  },
});
```

### With Mutation

```tsx
const mutation = useMutation({
  mutationFn: createCustomer,
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ["customers"] });
    toast.success(t("customer.created"));
  },
});

const form = useKyoraForm({
  defaultValues: { name: "", email: "" },
  onSubmit: async ({ value }) => {
    await mutation.mutateAsync(value);
  },
});
```

### With Server Errors

```tsx
import { createServerErrorValidator } from "@/lib/form";

const form = useKyoraForm({
  validators: {
    email: {
      onBlur: z.string().email("validation.invalid_email"),
      onServer: createServerErrorValidator(), // Injects RFC7807 errors
    },
  },
  onSubmit: async ({ value }) => {
    try {
      await api.register(value);
    } catch (error) {
      // Field-level errors automatically injected
      // Form-level errors shown in <form.FormError />
    }
  },
});
```

---

## 7. Bottom Sheet Pattern

```tsx
function MySheet({ isOpen, onClose }) {
  const form = useKyoraForm({ /* ... */ });
  const formId = useId();

  return (
    <form.AppForm>  {/* MUST wrap everything */}
      <BottomSheet
        isOpen={isOpen}
        onClose={onClose}
        title="Add Customer"
        footer={  {/* Footer MUST be here */}
          <div className="flex gap-2">
            <Button
              type="button"     {/* MUST be type="button" */}
              variant="ghost"
              onClick={onClose}
              disabled={isSubmitting}
            >
              Cancel
            </Button>
            <form.SubmitButton
              form={formId}     {/* MUST match FormRoot id */}
              variant="primary"
              disabled={isSubmitting}
            >
              Save
            </form.SubmitButton>
          </div>
        }
      >
        <form.FormRoot id={formId} className="space-y-4">
          {/* Fields */}
        </form.FormRoot>
      </BottomSheet>
    </form.AppForm>
  );
}
```

**Rules:**

- `form.AppForm` wraps entire BottomSheet (including footer)
- Cancel button has `type="button"` to prevent submission
- `form.SubmitButton` has `form={formId}` matching `FormRoot` id
- Both buttons disabled while `isSubmitting`

---

## 8. Translation Keys

**ALWAYS use translation keys, NEVER hardcoded messages:**

```tsx
// ❌ WRONG
z.string().min(8, "Password must be at least 8 characters");

// ✅ CORRECT
z.string().min(8, "validation.password_min_length");
```

**Translation key requirements:**

- Must be prefixed with `validation.`
- Must exist in `i18n/en/errors.json` under `validation.*`
- Must have Arabic translation in `i18n/ar/errors.json`

**Common keys:**

```typescript
"validation.required";
"validation.invalid_email";
"validation.invalid_phone";
"validation.min_length"; // interpolates {{min}}
"validation.max_length"; // interpolates {{max}}
"validation.positive_number";
"validation.invalid_date";
```

---

## 9. Performance

### Subscription Pattern

```tsx
// ❌ Bad: Re-renders on every state change
const values = form.useStore((state) => state.values);

// ✅ Good: Only re-renders when values change
<form.Subscribe selector={(state) => state.values}>
  {(values) => <div>{values.email}</div>}
</form.Subscribe>;
```

### Validation Debouncing

```tsx
validators={{
  onChangeAsyncDebounceMs: 500,
  onChangeAsync: async ({ value }) => {
    // Only fires after 500ms of no changes
    return await checkAvailability(value);
  },
}}
```

---

## 10. Common Patterns

### Password Confirmation

```tsx
<form.AppField
  name="confirmPassword"
  validators={{
    onChangeListenTo: ["password"],
    onBlur: ({ value, fieldApi }) => {
      if (value !== fieldApi.form.getFieldValue("password")) {
        return "validation.passwords_must_match";
      }
      return undefined;
    },
  }}
>
  {(field) => <field.PasswordField label="Confirm Password" />}
</form.AppField>
```

### Username Availability

```tsx
<form.AppField
  name="username"
  validators={{
    onChangeAsync: async ({ value }) => {
      if (value.length < 3) return "validation.username_too_short";
      const exists = await checkUsername(value);
      return exists ? "validation.username_taken" : undefined;
    },
    onChangeAsyncDebounceMs: 500,
  }}
>
  {(field) => <field.TextField label="Username" />}
</form.AppField>
```

### Date Range (Start < End)

```tsx
<form.AppField
  name="endDate"
  validators={{
    onChangeListenTo: ["startDate"],
    onBlur: ({ value, fieldApi }) => {
      const startDate = fieldApi.form.getFieldValue("startDate");
      if (startDate && value && value <= startDate) {
        return "validation.end_date_must_be_after_start_date";
      }
      return undefined;
    },
  }}
>
  {(field) => <field.DateField label="End Date" />}
</form.AppField>
```

---

## 11. Troubleshooting

### "formContext only works when..."

**Cause:** Component using form context is outside `<form.AppForm>`

**Solution:** Wrap ENTIRE component in `<form.AppForm>`, including footer

### Form not submitting

**Cause:** `<form.FormRoot>` id doesn't match `<form.SubmitButton form="...">`

**Solution:** Use `useId()` and pass same value to both

### Field not validating

**Cause:** Not using AppField pattern

**Solution:** Use `<form.AppField>` with validators prop

---

## Agent Validation

Before completing form task:

- ☑ `<form.AppForm>` wraps everything (including BottomSheet footer)
- ☑ **ZERO** direct component usage - all fields use `<form.AppField>` pattern
- ☑ FormRoot has `id` prop matching SubmitButton's `form` prop
- ☑ All validators use `onBlur` by default
- ☑ **ZERO** hardcoded error messages - all use `validation.*` keys
- ☑ All validation keys exist in `i18n/*/errors.json`
- ☑ Cancel button has `type="button"`
- ☑ All labels use `t()` function
- ☑ RTL: No `left`/`right` classes
- ☑ LTR-only fields (phone, codes) have `dir="ltr"`

---

## Resources

- TanStack Form: https://tanstack.com/form/latest
- Zod Docs: https://zod.dev
- Implementation: `portal-web/src/lib/form/`
