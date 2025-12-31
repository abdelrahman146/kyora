# Portal Web Form System Documentation

## Overview

The portal-web uses a sophisticated form management system built on **TanStack Form v1** with a custom `useKyoraForm` composition layer that eliminates boilerplate while providing production-grade form handling.

### Key Architecture

**TanStack Form Composition Pattern:**
- `useKyoraForm` returns a form instance with `form.AppForm` (provides form context) and `form.AppField` (provides field context)
- Components registered in `fieldComponents` are accessed via `field.TextField`, `field.PasswordField`, etc.
- Components registered in `formComponents` (FormRoot, SubmitButton, FormError) require form context from `form.AppForm`

### Key Features

- ✅ **Zero Boilerplate**: Pre-bound components eliminate 75% of manual wiring
- ✅ **Progressive Validation**: Smart revalidation logic (submit → blur modes)
- ✅ **Auto-Translation**: Zod error keys automatically translated via i18n
- ✅ **Focus Management**: Automatic focus on first invalid field
- ✅ **Type-Safe**: Full TypeScript support with inferred types
- ✅ **Server Errors**: RFC7807 problem details integration
- ✅ **Performance**: Granular subscriptions prevent unnecessary re-renders

### ⚠️ Critical Rule

**ALL components that use form context (FormRoot, SubmitButton, FormError, Subscribe) MUST be inside `<form.AppForm>`.**

If you see: `Error: formContext only works when within a formComponent passed to createFormHook`, you have a component using form context placed outside `<form.AppForm>`.

## Quick Start

###  Basic Form

**IMPORTANT:** Always wrap forms in `<form.AppForm>` and use `<form.AppField>` with `field.ComponentName` pattern.

```tsx
import { useKyoraForm } from '@/lib/form'
import { z } from 'zod'

function LoginForm() {
  const { t } = useTranslation()
  
  const form = useKyoraForm({
    defaultValues: {
      email: '',
      password: '',
    },
    onSubmit: async ({ value }) => {
      await api.login(value)
    },
  })

  return (
    <form.AppForm>
      <form.FormRoot className="space-y-4">
        <form.FormError />
        
        <form.AppField
          name="email"
          validators={{
            onBlur: z.string().email('invalid_email'),
          }}
        >
          {(field) => (
            <field.TextField
              type="email"
              label={t('auth.email')}
              placeholder={t('auth.email_placeholder')}
              autoComplete="email"
            />
          )}
        </form.AppField>

        <form.AppField
          name="password"
          validators={{
            onBlur: z.string().min(8, 'password_too_short'),
          }}
        >
          {(field) => (
            <field.PasswordField
              label={t('auth.password')}
              autoComplete="current-password"
            />
          )}
        </form.AppField>

        <form.SubmitButton variant="primary">
          {t('auth.login')}
        </form.SubmitButton>
      </form.FormRoot>
    </form.AppForm>
  )
}
```

**Key points:**
1. `<form.AppForm>` wraps everything (provides form context)
2. `<form.AppField>` instead of `<form.Field>` (provides field context)
3. Use `{(field) => <field.TextField />}` pattern (components from `fieldComponents`)
4. FormRoot, SubmitButton, FormError must be inside `<form.AppForm>`

## API Reference

### `useKyoraForm(config)`

Returns a form instance with pre-bound components.

**Config:**
```typescript
{
  defaultValues: Record<string, any>
  validators?: Record<string, { onBlur?: ZodSchema }>
  onSubmit: (data: { value: T }) => void | Promise<void>
}
```

**Returns:**
```typescript
{
  // Pre-bound components
  FormRoot: Component         // Replaces <form>
  TextField: Component        // Text/email/tel inputs
  PasswordField: Component    // Password with toggle
  TextareaField: Component    // Multi-line text input
  SelectField: Component      // Dropdown with search
  CheckboxField: Component    // Checkbox with label
  RadioField: Component       // Radio button group
  ToggleField: Component      // Toggle/switch
  SubmitButton: Component     // Submit button
  ErrorInfo: Component        // Field error display
  FormError: Component        // Form-level errors
  
  // TanStack Form primitives  
  Field: Component            // Custom fields
  Subscribe: Component        // Granular subscriptions
  
  // Form methods
  handleSubmit: () => void
  setFieldValue: (name, value) => void
  reset: () => void
  // ... all TanStack Form methods
}
```

### Pre-bound Components

#### `<form.TextField>`

Standard text input with automatic error handling.

```tsx
<form.TextField
  name="email"                    // Field name (required)
  type="email"                    // Input type
  label="Email"                   // Field label
  placeholder="Enter email"       // Placeholder text
  autoComplete="email"            // Autocomplete hint
  inputMode="email"               // Mobile keyboard
  required                        // Visual indicator
  disabled                        // Disable input
/>
```

**Auto-handled:**
- Value binding
- Change handlers
- Blur handlers
- Error display (translated)
- Aria attributes

#### `<form.PasswordField>`

Password input with visibility toggle.

```tsx
<form.PasswordField
  name="password"
  label="Password"
  autoComplete="current-password"
  required
/>
```

**Features:**
- Eye icon toggle
- Translated labels
- All TextField features

#### `<form.TextareaField>`

Multi-line text input with character counter.

```tsx
<form.TextareaField
  name="description"
  label="Description"
  placeholder="Enter description"
  rows={4}                        // Number of visible rows
  maxLength={500}                 // Character limit
  showCount                       // Show character counter
  required
/>
```

**Auto-handled:**
- Value binding
- Change/blur handlers
- Error display (translated)
- Character counter

#### `<form.SelectField>`

Dropdown select with search and multi-select support.

**Single Select:**
```tsx
<form.AppField
  name="country"
  validators={{
    onBlur: z.string().min(1, 'required'),
  }}
>
  {(field) => (
    <field.SelectField
      label={t('customer.country')}
      options={[
        { value: 'us', label: 'United States' },
        { value: 'uk', label: 'United Kingdom' },
        { value: 'eg', label: 'Egypt' },
      ]}
      searchable                      // Enable search
      clearable                       // Show clear button
      required
    />
  )}
</form.AppField>
```

**Multi-Select with Chip UI:**
```tsx
<form.AppField
  name="tags"
  validators={{
    onBlur: z.array(z.string())
      .min(1, 'select_at_least_one')
      .max(5, 'select_too_many'),
  }}
>
  {(field) => (
    <field.SelectField
      label={t('customer.tags')}
      options={[
        { value: 'vip', label: 'VIP Customer', icon: <Star /> },
        { value: 'wholesale', label: 'Wholesale', icon: <Package /> },
        { value: 'repeat', label: 'Repeat Buyer', icon: <RefreshCw /> },
      ]}
      multiSelect                     // Enable multi-select mode
      searchable                      // Enable search
      clearable                       // Clear all selections
      required
    />
  )}
</form.AppField>
```

**Validation Patterns:**
```typescript
// Minimum selections
z.array(z.string()).min(1, 'select_at_least_one')

// Maximum selections
z.array(z.string()).max(5, 'select_too_many')

// Min and max
z.array(z.string()).min(1).max(5)

// Unique values (no duplicates)
z.array(z.string()).refine(
  (arr) => new Set(arr).size === arr.length,
  { message: 'duplicate_selection' }
)

// Custom validation
z.array(z.string()).refine(
  (arr) => arr.every((v) => validValues.includes(v)),
  { message: 'invalid_selection' }
)
```

**Features:**
- Search/filtering with real-time results
- Multi-select with chip-based UI
- Keyboard navigation (Arrow keys, Space/Enter, Backspace to remove last chip)
- Chip removal (click X button or Backspace/Delete keys)
- Clear all selections button
- RTL support (chips flow right-to-left in Arabic)
- Mobile bottom sheet / Desktop dropdown
- Touch-optimized (50px minimum height)
- Screen reader accessible

**Translation Keys Used:**
- `common.selected_count`: "{{count}} selected"
- `common.remove`: "Remove {{item}}"
- `common.clear_selection`: "Clear selection"
- `common.search_placeholder_generic`: "Search..."
- `common.no_options_found`: "No options found"
- `errors.validation.select_at_least_one`: "Please select at least one option."
- `errors.validation.select_too_many`: "You can select at most {{max}} options."
- `errors.validation.duplicate_selection`: "Duplicate selections are not allowed."
- `errors.validation.array_min_items`: "Please select at least {{min}} item(s)."
- `errors.validation.array_max_items`: "You can select at most {{max}} item(s)."

#### `<form.CheckboxField>`

Checkbox with label and description.

```tsx
<form.CheckboxField
  name="acceptTerms"
  label="Accept terms"
  description="I agree to the terms and conditions"
  required
/>
```

**Auto-handled:**
- Boolean value binding
- Error display
- Accessibility

#### `<form.RadioField>`

Radio button group with flexible layout.

```tsx
<form.RadioField
  name="plan"
  label="Select a plan"
  options={[
    { value: 'free', label: 'Free', description: '$0/month' },
    { value: 'pro', label: 'Pro', description: '$10/month' },
  ]}
  orientation="vertical"          // vertical | horizontal
  variant="primary"
  required
/>
```

**Features:**
- Multiple layout options
- Option descriptions
- Keyboard navigation

#### `<form.ToggleField>`

Toggle/switch component.

```tsx
<form.ToggleField
  name="notifications"
  label="Enable notifications"
  description="Receive email updates"
  size="md"
  variant="primary"
/>
```

**Auto-handled:**
- Boolean value binding
- Toggle state
- Error display

#### `<form.SubmitButton>`

Submit button with loading state.

```tsx
<form.SubmitButton
  variant="primary"               // btn-primary
  form="my-form-id"               // External form
  className="w-full"              // Additional classes
>
  Submit
</form.SubmitButton>
```

**Auto-handled:**
- Disabled during submission
- Loading state
- Form submission

#### `<form.Field>`

For custom controls not covered by pre-bound components.

```tsx
<form.Field name="customField">
  {(field) => (
    <CustomControl
      value={field.state.value}
      onChange={(value) => field.handleChange(value)}
      onBlur={field.handleBlur}
    />
  )}
</form.Field>
```

#### `<form.Subscribe>`

Granular subscriptions to prevent re-renders.

```tsx
<form.Subscribe selector={(state) => state.values}>
  {(values) => (
    <SocialMediaInputs
      instagram={values.instagram}
      onChange={(value) => form.setFieldValue('instagram', value)}
    />
  )}
</form.Subscribe>
```

## Validation

### Field-Level Validation

Use Zod schemas with translation keys as error messages:

```tsx
const form = useKyoraForm({
  validators: {
    email: { 
      onBlur: z.string()
        .min(1, 'required')                    // errors.required
        .email('invalid_email')                // errors.invalid_email
    },
    password: { 
      onBlur: z.string()
        .min(8, 'password_too_short')          // errors.password_too_short
    },
    confirmPassword: {
      onBlur: z.string()
        .min(1, 'required')
    },
  },
  // ...
})
```

### Cross-Field Validation

For fields that depend on each other:

```tsx
const form = useKyoraForm({
  validators: {
    phoneCode: {
      onChange: ({ value, fieldApi }) => {
        // Access other fields
        const phoneNumber = fieldApi.form.getFieldValue('phoneNumber')
        if (phoneNumber && !value) {
          return 'Phone code required when number provided'
        }
        return undefined
      },
    },
  },
  // ...
})
```

### Validation Timing

The form uses **progressive validation**:

| State | Behavior |
|-------|----------|
| Initial | No validation |
| First Submit | Validate all fields |
| After Submit | Validate on blur |
| Submit Again | Validate on submit |

This provides optimal UX: users aren't bothered until they try to submit.

## Error Handling

### Translation Flow

```
Zod Error Key → i18n errors namespace → Displayed Message
```

Example:
```tsx
// Validator
z.string().email('invalid_email')

// Translation (en/errors.json)
{
  "invalid_email": "Please enter a valid email address"
}

// Displayed
"Please enter a valid email address"
```

### Server Errors

Inject server errors using `createServerErrorValidator`:

```tsx
import { createServerErrorValidator } from '@/lib/form'

const form = useKyoraForm({
  validators: {
    email: { 
      onBlur: z.string().email('invalid_email'),
      onServer: createServerErrorValidator(),  // Injects server errors
    },
  },
  onSubmit: async ({ value }) => {
    try {
      await api.register(value)
    } catch (error) {
      // If RFC7807 error with field-level details
      // errors will automatically appear on fields
    }
  },
})
```

### Form-Level Errors

For errors that don't belong to a specific field:

```tsx
<form.FormRoot>
  <form.FormError />  {/* Shows form-level errors */}
  
  <form.TextField name="email" />
  {/* ... */}
</form.FormRoot>
```

## Advanced Patterns

### ⚠️ CRITICAL: Form Context and Component Placement

**The `form.AppForm` wrapper is REQUIRED for all form components that use form context.**

Components that use `useFormContext()` internally (FormRoot, SubmitButton, FormError) **MUST** be inside `<form.AppForm>`:

```tsx
// ❌ WRONG - SubmitButton outside form.AppForm
function MyForm() {
  const form = useKyoraForm({ /* ... */ })
  
  return (
    <>
      <form.AppForm>
        <form.FormRoot>
          {/* fields */}
        </form.FormRoot>
      </form.AppForm>
      <form.SubmitButton>Submit</form.SubmitButton>  {/* ❌ Error: formContext not available */}
    </>
  )
}

// ✅ CORRECT - All form components inside form.AppForm
function MyForm() {
  const form = useKyoraForm({ /* ... */ })
  
  return (
    <form.AppForm>
      <form.FormRoot>
        {/* fields */}
      </form.FormRoot>
      <form.SubmitButton>Submit</form.SubmitButton>  {/* ✅ Has form context */}
    </form.AppForm>
  )
}
```

**Error symptom:** If you see `Error: formContext only works when within a formComponent passed to createFormHook`, it means a component using `useFormContext()` is placed outside `<form.AppForm>`.

### External Form Submission

For modals/sheets with footer buttons outside the form:

```tsx
function AddCustomerSheet() {
  const formId = useId()
  const form = useKyoraForm({ /* ... */ })
  
  // ✅ CORRECT - Wrap ENTIRE component in form.AppForm
  return (
    <form.AppForm>
      <BottomSheet
        footer={
          <div>
            <button onClick={onClose}>Cancel</button>
            <form.SubmitButton form={formId}>  {/* ✅ Has access to form context */}
              Submit
            </form.SubmitButton>
          </div>
        }
      >
        <form.FormRoot id={formId}>
          {/* fields */}
        </form.FormRoot>
      </BottomSheet>
    </form.AppForm>
  )
}

// ❌ WRONG - form.AppForm only wraps FormRoot
function AddCustomerSheetWrong() {
  const formId = useId()
  const form = useKyoraForm({ /* ... */ })
  
  const footer = (  // ❌ Defined outside form.AppForm
    <div>
      <form.SubmitButton form={formId}>  {/* ❌ No form context! */}
        Submit
      </form.SubmitButton>
    </div>
  )
  
  return (
    <BottomSheet footer={footer}>
      <form.AppForm>  {/* ❌ form.AppForm in wrong place */}
        <form.FormRoot id={formId}>
          {/* fields */}
        </form.FormRoot>
      </form.AppForm>
    </BottomSheet>
  )
}
```

**Key principle:** `<form.AppForm>` must be the outermost wrapper that contains ALL components using form context (FormRoot, SubmitButton, FormError, Subscribe).

### Dependent Fields

Auto-update fields based on others:

```tsx
function CustomerForm() {
  const [selectedCountry, setSelectedCountry] = useState('US')
  
  const form = useKyoraForm({ /* ... */ })
  
  // Auto-link country to phone code
  useEffect(() => {
    const country = countries.find(c => c.code === selectedCountry)
    if (country?.phonePrefix) {
      form.setFieldValue('phoneCode', country.phonePrefix)
    }
  }, [selectedCountry, countries])
  
  return (
    <form.FormRoot>
      <form.Field name="countryCode">
        {(field) => (
          <CountrySelect
            value={field.state.value}
            onChange={(value) => {
              field.handleChange(value)
              setSelectedCountry(value)  // Trigger effect
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
  )
}
```

### Uncontrolled Components

For components that manage their own state:

```tsx
<form.Subscribe selector={(state) => state.values}>
  {(values) => (
    <SocialMediaInputs
      instagramUsername={values.instagramUsername}
      onInstagramChange={(value) =>
        form.setFieldValue('instagramUsername', value)
      }
      facebookUsername={values.facebookUsername}
      onFacebookChange={(value) =>
        form.setFieldValue('facebookUsername', value)
      }
    />
  )}
</form.Subscribe>
```

**Why Subscribe?**
- Reads form state reactively
- Only re-renders when selected state changes
- Prevents unnecessary re-renders of parent

## Performance

### Subscription Pattern

**❌ Bad:** Causes re-render on every state change
```tsx
const values = form.useStore((state) => state.values)
```

**✅ Good:** Only re-renders when values change
```tsx
<form.Subscribe selector={(state) => state.values}>
  {(values) => <div>{values.email}</div>}
</form.Subscribe>
```

### Field Isolation

Each `form.Field` is isolated - changing one field doesn't re-render others.

### Validation Debouncing

Validation runs immediately on blur, but you can debounce onChange:

```tsx
const form = useKyoraForm({
  validators: {
    email: {
      onBlur: z.string().email('invalid_email'),
      onChange: debounce(
        z.string().email('invalid_email'),
        300
      ),
    },
  },
})
```

## Migration Guide

### From Raw TanStack Form

**Before:**
```tsx
const form = useForm({
  defaultValues: { email: '' },
  validators: { onBlur: z.object({ email: z.string().email() }) },
  onSubmit: async ({ value }) => { /* ... */ },
})

return (
  <form onSubmit={(e) => { e.preventDefault(); form.handleSubmit() }}>
    <form.Field name="email">
      {(field) => (
        <div>
          <label>Email</label>
          <input
            value={field.state.value}
            onChange={(e) => field.handleChange(e.target.value)}
            onBlur={field.handleBlur}
          />
          {field.state.meta.errors.length > 0 && (
            <span>{t('errors:' + field.state.meta.errors[0])}</span>
          )}
        </div>
      )}
    </form.Field>
    
    <button type="submit" disabled={form.state.isSubmitting}>
      Submit
    </button>
  </form>
)
```

**After (27 → 7 lines):**
```tsx
const form = useKyoraForm({
  defaultValues: { email: '' },
  validators: {
    email: { onBlur: z.string().email('invalid_email') },
  },
  onSubmit: async ({ value }) => { /* ... */ },
})

return (
  <form.FormRoot>
    <form.TextField
      name="email"
      label="Email"
      type="email"
      required
    />
    
    <form.SubmitButton variant="primary">
      Submit
    </form.SubmitButton>
  </form.FormRoot>
)
```

**Benefits:**
- 74% less code
- Auto error translation
- Auto focus management
- No manual event handlers
- Type-safe
- Consistent UX

### Checklist

- [ ] Replace `useForm` with `useKyoraForm`
- [ ] Remove `{ t: tErrors }` translation import
- [ ] Convert validators from `{ onBlur: schema }` to field-level
- [ ] Replace `<form>` with `<form.FormRoot>`
- [ ] Replace standard inputs with `<form.TextField>`
- [ ] Replace password inputs with `<form.PasswordField>`
- [ ] Use `<form.Field>` for custom controls
- [ ] Replace submit buttons with `<form.SubmitButton>`
- [ ] Remove manual error display code
- [ ] Remove `form.state.isSubmitting` (use mutation.isPending)
- [ ] Wrap uncontrolled components in `<form.Subscribe>`
- [ ] Test validation timing
- [ ] Test error display
- [ ] Test focus management

## Examples

### Complete Login Form

See [src/components/organisms/LoginForm.tsx](../src/components/organisms/LoginForm.tsx)

### Password Reset

See [src/routes/auth/reset-password.tsx](../src/routes/auth/reset-password.tsx)

### Multi-Step Form

See [src/routes/onboarding/verify.tsx](../src/routes/onboarding/verify.tsx)

### Complex Form with Dependencies

See [src/routes/onboarding/business.tsx](../src/routes/onboarding/business.tsx)

### Multi-Select Form with Validation

Complete example with array validation, chip UI, and keyboard support:

```tsx
import { useKyoraForm } from '@/lib/form'
import { z } from 'zod'
import { Star, Package, RefreshCw } from 'lucide-react'

function CustomerTagsForm() {
  const { t } = useTranslation()
  
  const form = useKyoraForm({
    defaultValues: {
      customerName: '',
      tags: [], // Array for multi-select
      priority: '',
    },
    onSubmit: async ({ value }) => {
      await api.updateCustomer(customerId, value)
      toast.success(t('customer.updated'))
    },
  })

  const tagOptions = [
    { value: 'vip', label: 'VIP Customer', icon: <Star className="w-4 h-4" /> },
    { value: 'wholesale', label: 'Wholesale', icon: <Package className="w-4 h-4" /> },
    { value: 'repeat', label: 'Repeat Buyer', icon: <RefreshCw className="w-4 h-4" /> },
    { value: 'new', label: 'New Customer' },
    { value: 'discount', label: 'Discount Eligible' },
  ]

  const priorityOptions = [
    { value: 'high', label: 'High Priority' },
    { value: 'medium', label: 'Medium Priority' },
    { value: 'low', label: 'Low Priority' },
  ]

  return (
    <form.AppForm>
      <form.FormRoot className="space-y-6 max-w-2xl">
        <form.FormError />
        
        {/* Customer Name */}
        <form.AppField
          name="customerName"
          validators={{
            onBlur: z.string().min(2, 'min_length').max(100, 'max_length'),
          }}
        >
          {(field) => (
            <field.TextField
              label={t('customer.name')}
              placeholder={t('customer.name_placeholder')}
              required
            />
          )}
        </form.AppField>

        {/* Multi-Select Tags with Validation */}
        <form.AppField
          name="tags"
          validators={{
            onBlur: z.array(z.string())
              .min(1, 'select_at_least_one')
              .max(3, 'select_too_many')
              .refine(
                (arr) => new Set(arr).size === arr.length,
                { message: 'duplicate_selection' }
              ),
          }}
        >
          {(field) => (
            <field.SelectField
              label={t('customer.tags')}
              helperText={t('customer.tags_helper')}
              options={tagOptions}
              multiSelect
              searchable
              clearable
              required
            />
          )}
        </form.AppField>

        {/* Single Select Priority */}
        <form.AppField
          name="priority"
          validators={{
            onBlur: z.string().min(1, 'required'),
          }}
        >
          {(field) => (
            <field.SelectField
              label={t('customer.priority')}
              options={priorityOptions}
              searchable
              clearable
              required
            />
          )}
        </form.AppField>

        {/* Submit Button */}
        <form.SubmitButton variant="primary" size="lg">
          {t('customer.save_changes')}
        </form.SubmitButton>
      </form.FormRoot>
    </form.AppForm>
  )
}
```

**Key Features Demonstrated:**
- Multi-select with chip UI (tags field)
- Array validation (min 1, max 3, unique values)
- Single select with search (priority field)
- Icons in select options
- Keyboard navigation (Backspace to remove last tag)
- Translated error messages
- Helper text for guidance
- RTL layout support

## Troubleshooting

### Errors Not Translating

**Problem:** Errors show as keys like "invalid_email"

**Solution:** Ensure error keys exist in `src/i18n/*/errors.json`

### Form Not Submitting

**Problem:** Submit button does nothing

**Solution:** Check `<form.FormRoot>` has correct `id` matching `<form.SubmitButton form="...">`

### Field Not Validating

**Problem:** No error shown despite invalid input

**Solution:** Ensure validator is in `validators` object, not `onBlur` top-level

### Too Many Re-renders

**Problem:** Component re-renders on every keystroke

**Solution:** Use `<form.Subscribe>` instead of `form.useStore` for derived state

### Server Errors Not Showing

**Problem:** API errors don't appear on fields

**Solution:** Add `onServer: createServerErrorValidator()` to validators

## Architecture

### Component Hierarchy

```
useKyoraForm (composition layer)
  ├── TanStack Form (state management)
  ├── Zod (validation)
  ├── i18n (translation)
  └── Pre-bound Components
      ├── FormRoot
      ├── TextField
      ├── PasswordField
      ├── SubmitButton
      ├── ErrorInfo
      └── FormError
```

### File Structure

```
src/lib/form/
├── index.ts                    # Public API
├── useKyoraForm.ts             # Main hook
├── createFormHook.ts           # Composition factory
├── revalidateLogic.ts          # Progressive validation
├── useFocusOnError.ts          # Auto-focus management
├── createServerErrorValidator.ts  # Server error injection
└── components/
    ├── FormRoot.tsx
    ├── TextField.tsx
    ├── PasswordField.tsx
    ├── SubmitButton.tsx
    ├── ErrorInfo.tsx
    └── FormError.tsx
```

### Design Decisions

**Why TanStack Form?**
- Type-safe
- Framework-agnostic
- Granular subscriptions
- Powerful validation
- Battle-tested

**Why Zod?**
- Type inference
- Composable schemas
- Rich validation primitives
- Error customization

**Why Composition Layer?**
- Eliminates boilerplate
- Enforces consistency
- Easy to extend
- Opt-in complexity

**Why Translation Keys?**
- Single source of truth
- No manual translation calls
- Consistent UX
- Easy to maintain

## Contributing

When adding new form patterns:

1. Check if existing components cover the use case
2. If custom component needed, use `<form.Field>`
3. Document new patterns in this file
4. Add example to examples section
5. Update migration checklist if needed

## Future Enhancements

- [ ] Async validation support
- [ ] File upload component
- [ ] Multi-select support in FormSelect
- [ ] Date/time pickers
- [ ] Rich text editor
- [ ] Form arrays (repeating fields)
- [ ] Wizard/stepper pattern
- [ ] Auto-save draft
- [ ] Optimistic UI updates
