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

#### `<form.DateField>`

Date picker with calendar popup.

```tsx
<form.AppField
  name="birthdate"
  validators={{
    onBlur: z.date()
      .max(new Date(), 'date_cannot_be_future')
      .refine(
        (date) => {
          const age = new Date().getFullYear() - date.getFullYear()
          return age >= 18
        },
        { message: 'must_be_18_or_older' }
      ),
  }}
>
  {(field) => (
    <field.DateField
      label="Birth Date"
      minAge={18}                     // Min age validation (years)
      maxDate={new Date()}            // Max date allowed
      disableWeekends                 // Disable Sat/Sun
      clearable                       // Show clear button
      required
    />
  )}
</form.AppField>
```

**Features:**
- Calendar popup with month/year navigation
- RTL support (Arabic/English locales)
- Date format: dd/MM/yyyy (Arabic) or MM/dd/yyyy (English)
- Keyboard navigation (Arrow keys, Page Up/Down, Home/End, Enter, Escape)
- Mobile: Full-screen modal
- Desktop: Dropdown popup
- Min/max date validation
- Disabled dates support
- Clear button
- Touch-optimized (50px minimum height)

**Validation Patterns:**
```typescript
// Basic date required
z.date()

// Date cannot be in future
z.date().max(new Date(), 'date_cannot_be_future')

// Minimum age (18+)
z.date().refine(
  (date) => {
    const age = new Date().getFullYear() - date.getFullYear()
    return age >= 18
  },
  { message: 'must_be_18_or_older' }
)

// Date range
z.date()
  .min(new Date('2020-01-01'))
  .max(new Date('2025-12-31'))

// Custom business logic
z.date().refine(
  (date) => {
    const day = date.getDay()
    return day !== 0 && day !== 6 // No weekends
  },
  { message: 'weekdays_only' }
)
```

**Translation Keys Used:**
- `common.date.selectDate`: "Select date"
- `common.date.today`: "Today"
- `common.clear`: "Clear"
- `errors.validation.invalid_date`: "Please enter a valid date."
- `errors.validation.date_cannot_be_future`: "Date cannot be in the future."
- `errors.validation.must_be_18_or_older`: "You must be at least 18 years old."

#### `<form.TimeField>`

Time picker with hour/minute controls.

```tsx
<form.AppField
  name="appointmentTime"
  validators={{
    onBlur: z.date()
      .refine(
        (date) => {
          const hours = date.getHours()
          return hours >= 9 && hours < 17 // 9 AM - 5 PM
        },
        { message: 'outside_business_hours' }
      ),
  }}
>
  {(field) => (
    <field.TimeField
      label="Appointment Time"
      use24Hour={false}               // 12-hour format with AM/PM
      minuteStep={15}                 // 15-minute increments
      clearable
      required
    />
  )}
</form.AppField>
```

**Features:**
- Two numeric inputs (hours, minutes)
- AM/PM toggle for 12-hour format
- 24-hour format support (locale-aware)
- Arrow buttons for increment/decrement
- Keyboard navigation (Arrow Up/Down, Tab)
- Auto-advance from hours to minutes
- Minute step intervals (1, 5, 15, 30)
- Touch-optimized (50px minimum height)

**Validation Patterns:**
```typescript
// Basic time required
z.date()

// Business hours only (9 AM - 5 PM)
z.date().refine(
  (date) => {
    const hours = date.getHours()
    return hours >= 9 && hours < 17
  },
  { message: 'outside_business_hours' }
)

// Specific time range
z.date().refine(
  (date) => {
    const hours = date.getHours()
    const minutes = date.getMinutes()
    const totalMinutes = hours * 60 + minutes
    return totalMinutes >= 540 && totalMinutes <= 1020 // 9:00 AM - 5:00 PM
  },
  { message: 'outside_hours' }
)

// 15-minute intervals only
z.date().refine(
  (date) => date.getMinutes() % 15 === 0,
  { message: 'must_be_15_min_intervals' }
)
```

**Translation Keys Used:**
- `common.date.hours`: "Hours"
- `common.date.minutes`: "Minutes"
- `common.date.period`: "Period"
- `common.date.incrementHours`: "Increment hours"
- `common.date.decrementHours`: "Decrement hours"
- `common.date.incrementMinutes`: "Increment minutes"
- `common.date.decrementMinutes`: "Decrement minutes"
- `errors.validation.invalid_time`: "Please enter a valid time."
- `errors.validation.outside_business_hours`: "Time must be within business hours (9 AM - 5 PM)."

#### `<form.DateTimeField>`

Combined date and time picker with flexible modes.

```tsx
// Date and Time mode
<form.AppField
  name="eventDateTime"
  validators={{
    onBlur: z.date()
      .min(new Date(), 'must_be_future')
      .refine(
        (date) => {
          const hours = date.getHours()
          return hours >= 9 && hours < 17
        },
        { message: 'outside_business_hours' }
      ),
  }}
>
  {(field) => (
    <field.DateTimeField
      mode="datetime"                 // 'date' | 'time' | 'datetime'
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

// Date only mode (equivalent to DateField)
<form.AppField name="startDate">
  {(field) => (
    <field.DateTimeField
      mode="date"
      label="Start Date"
    />
  )}
</form.AppField>

// Time only mode (equivalent to TimeField)
<form.AppField name="meetingTime">
  {(field) => (
    <field.DateTimeField
      mode="time"
      label="Meeting Time"
    />
  )}
</form.AppField>
```

**Features:**
- Three modes: date, time, datetime
- Responsive layout: side-by-side on desktop, stacked on mobile
- Independent date and time props
- Preserves date when changing time and vice versa
- All DateField and TimeField features combined

**Validation Patterns:**
```typescript
// End date must be after start date
const form = useKyoraForm({
  defaultValues: {
    startDate: null as Date | null,
    endDate: null as Date | null,
  },
})

<form.AppField
  name="endDate"
  validators={{
    onChange: z.date(),
    onBlurListenTo: ['startDate'],
    onBlur: ({ value, fieldApi }) => {
      const startDate = fieldApi.form.getFieldValue('startDate')
      if (startDate && value && value <= startDate) {
        return 'end_date_must_be_after_start_date'
      }
      return undefined
    },
  }}
>
  {(field) => (
    <field.DateTimeField
      mode="datetime"
      label="End Date & Time"
    />
  )}
</form.AppField>
```

**Translation Keys Used:**
- All DateField and TimeField translation keys
- `errors.validation.end_date_must_be_after_start_date`: "End date must be after start date."

#### `<form.DateRangeField>`

Date range picker for selecting start and end dates with dual-calendar interface.

```tsx
<form.AppField
  name="reportDateRange"
  validators={{
    onBlur: z.custom<DateRange>((val) => {
      const range = val as DateRange
      if (!range?.from || !range?.to) {
        return 'date_range_required'
      }
      if (range.to < range.from) {
        return 'date_range_invalid'
      }
      return undefined
    }),
  }}
>
  {(field) => (
    <field.DateRangeField
      label="Report Period"
      placeholder="Select date range"
      minDate={new Date('2020-01-01')}
      maxDate={new Date()}
      numberOfMonths={2}              // Show 2 months side-by-side (default: 2)
      disabledDates={[               // Disable specific dates
        new Date('2024-12-25'),
        new Date('2024-01-01'),
      ]}
      clearable
      required
    />
  )}
</form.AppField>
```

**Features:**
- Dual-calendar popup (configurable with `numberOfMonths`)
- Date range selection with visual highlighting
- Auto-close when both dates selected
- Clear button to reset selection
- Format: "MM/DD/YYYY - MM/DD/YYYY" (locale-aware)
- RTL support (Arabic/English locales)
- Keyboard navigation (Arrow keys navigate calendar, Escape closes, Enter opens)
- Mobile: Action bar with Clear/Apply buttons
- Desktop: Dropdown popup with side-by-side months
- Min/max date constraints
- Disabled dates support
- Touch-optimized (50px minimum height)

**DateRange Type:**
```typescript
import type { DateRange } from 'react-day-picker'

interface DateRange {
  from?: Date
  to?: Date
}

// Usage
const form = useKyoraForm({
  defaultValues: {
    dateRange: undefined as DateRange | undefined,
  },
})
```

**Common Use Cases:**

**1. Booking System:**
```tsx
<form.AppField
  name="bookingRange"
  validators={{
    onBlur: z.custom<DateRange>((val) => {
      const range = val as DateRange
      if (!range?.from || !range?.to) return 'date_range_required'
      
      // Minimum 1 night
      const days = Math.ceil((range.to.getTime() - range.from.getTime()) / (1000 * 60 * 60 * 24))
      if (days < 1) return 'minimum_one_night'
      
      // Maximum 30 days
      if (days > 30) return 'maximum_30_days'
      
      return undefined
    }),
  }}
>
  {(field) => (
    <field.DateRangeField
      label="Check-in / Check-out"
      minDate={new Date()}
      numberOfMonths={2}
    />
  )}
</form.AppField>
```

**2. Analytics Date Filter:**
```tsx
<form.AppField
  name="analyticsRange"
  validators={{
    onBlur: z.custom<DateRange>((val) => {
      const range = val as DateRange
      if (!range?.from || !range?.to) return 'date_range_required'
      
      // Cannot be future dates
      const now = new Date()
      if (range.to > now) return 'date_cannot_be_future'
      
      return undefined
    }),
  }}
>
  {(field) => (
    <field.DateRangeField
      label="Analysis Period"
      maxDate={new Date()}
      numberOfMonths={2}
    />
  )}
</form.AppField>
```

**3. Business Report with Weekday Validation:**
```tsx
<form.AppField
  name="reportRange"
  validators={{
    onBlur: z.custom<DateRange>((val) => {
      const range = val as DateRange
      if (!range?.from || !range?.to) return 'date_range_required'
      
      // Start and end must be weekdays
      if (range.from.getDay() === 0 || range.from.getDay() === 6) {
        return 'start_must_be_weekday'
      }
      if (range.to.getDay() === 0 || range.to.getDay() === 6) {
        return 'end_must_be_weekday'
      }
      
      return undefined
    }),
  }}
>
  {(field) => (
    <field.DateRangeField
      label="Report Period (Weekdays Only)"
      numberOfMonths={1}
      disabledDates={getWeekendDates()}  // Helper to disable all weekends
    />
  )}
</form.AppField>
```

**Validation Patterns:**
```typescript
// Basic date range required
z.custom<DateRange>((val) => {
  const range = val as DateRange
  if (!range?.from || !range?.to) return 'date_range_required'
  return undefined
})

// End date must be after start date
z.custom<DateRange>((val) => {
  const range = val as DateRange
  if (!range?.from || !range?.to) return 'date_range_required'
  if (range.to < range.from) return 'date_range_invalid'
  return undefined
})

// Maximum range duration (e.g., 90 days)
z.custom<DateRange>((val) => {
  const range = val as DateRange
  if (!range?.from || !range?.to) return 'date_range_required'
  
  const days = Math.ceil((range.to.getTime() - range.from.getTime()) / (1000 * 60 * 60 * 24))
  if (days > 90) return 'maximum_90_days'
  
  return undefined
})

// Minimum range duration (e.g., 7 days)
z.custom<DateRange>((val) => {
  const range = val as DateRange
  if (!range?.from || !range?.to) return 'date_range_required'
  
  const days = Math.ceil((range.to.getTime() - range.from.getTime()) / (1000 * 60 * 60 * 24))
  if (days < 7) return 'minimum_7_days'
  
  return undefined
})

// Business days only (no weekends)
z.custom<DateRange>((val) => {
  const range = val as DateRange
  if (!range?.from || !range?.to) return 'date_range_required'
  
  const isWeekday = (date: Date) => {
    const day = date.getDay()
    return day !== 0 && day !== 6
  }
  
  if (!isWeekday(range.from) || !isWeekday(range.to)) {
    return 'weekdays_only'
  }
  
  return undefined
})

// Range cannot include today
z.custom<DateRange>((val) => {
  const range = val as DateRange
  if (!range?.from || !range?.to) return 'date_range_required'
  
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  
  if (range.from <= today && range.to >= today) {
    return 'cannot_include_today'
  }
  
  return undefined
})
```

**Translation Keys Used:**
- `common.date.selectDateRange`: "Select date range"
- `common.clear`: "Clear"
- `common.apply`: "Apply"
- `errors.validation.invalid_date_range`: "Please enter a valid date range."
- `errors.validation.date_range_required`: "Both start and end dates are required."
- `errors.validation.date_range_invalid`: "End date must be after start date."
- `errors.validation.minimum_one_night`: "Booking must be at least 1 night."
- `errors.validation.maximum_30_days`: "Booking cannot exceed 30 days."
- `errors.validation.maximum_90_days`: "Date range cannot exceed 90 days."
- `errors.validation.minimum_7_days`: "Date range must be at least 7 days."
- `errors.validation.weekdays_only`: "Dates must be weekdays only."
- `errors.validation.cannot_include_today`: "Date range cannot include today."

**Accessibility:**
- `role="dialog"` on calendar popup
- `aria-label` with translated "Select date range"
- `aria-invalid` and `aria-describedby` for error states
- Focus trap within calendar
- Keyboard navigation with screen reader announcements

#### `<form.FieldArray>`

Dynamic array field for managing lists of repeating items with drag-and-drop reordering, add/remove operations, and array-specific validations.

```tsx
<form.AppField
  name="phoneNumbers"
  validators={{
    onChange: ({ value }) => {
      // Validate min/max items
      const minMaxError = validateArrayLength(value, {
        min: 1,
        max: 5,
        minErrorKey: 'form.atLeastOnePhone',
        maxErrorKey: 'form.tooManyPhones',
      })
      if (minMaxError) return minMaxError
      
      // Validate uniqueness
      return validateUniqueValues(value, {
        errorKey: 'form.duplicateEmails',
      })
    },
  }}
>
  {(field) => (
    <field.FieldArray
      label="Phone Numbers"
      addButtonLabel="Add Phone Number"
      emptyMessage="No phone numbers yet"
      minItems={1}
      maxItems={5}
      reorderable
      defaultValue={{ number: '', type: 'mobile' }}
      render={(item, operations, index) => (
        <div className="flex gap-2">
          <input
            type="tel"
            value={item.number}
            onChange={(e) => {
              const updated = [...field.state.value]
              updated[index] = { ...item, number: e.target.value }
              field.handleChange(updated)
            }}
            placeholder="Phone number"
            className="input flex-1"
          />
          <select
            value={item.type}
            onChange={(e) => {
              const updated = [...field.state.value]
              updated[index] = { ...item, type: e.target.value }
              field.handleChange(updated)
            }}
            className="select"
          >
            <option value="mobile">Mobile</option>
            <option value="home">Home</option>
            <option value="work">Work</option>
          </select>
        </div>
      )}
    />
  )}
</form.AppField>
```

**Props:**
- `label` (string): Field label
- `addButtonLabel` (string): Label for add button (default: "Add Item")
- `emptyMessage` (string): Message when array is empty
- `minItems` (number): Minimum number of items
- `maxItems` (number): Maximum number of items
- `reorderable` (boolean): Enable drag-and-drop reordering
- `defaultValue` (T): Default value for new items
- `render` (function): Render function for each item
  - `item`: Current item data
  - `operations`: Object with `remove()`, `moveUp()`, `moveDown()` methods
  - `index`: Item index

**Features:**
- Drag-and-drop reordering (via @dnd-kit)
- Add/remove operations with keyboard support
- Move up/down buttons (accessible alternative to drag-and-drop)
- Empty state with add button
- Max items warning (disables add button)
- Array validation helpers (min/max items, unique values, cross-item validation)
- Smooth animations (300ms transitions)
- RTL support (drag handles positioned correctly)
- Touch-optimized (50px minimum touch targets)
- Keyboard navigation (Tab, Space/Enter for drag, Arrow keys for reorder)
- Screen reader accessible (live regions, ARIA labels)

**Array Validation Utilities:**

The form system provides comprehensive validation utilities for array fields:

```typescript
import {
  validateArrayLength,
  validateUniqueValues,
  validateArrayItems,
  validateNoOverlap,
  validateArrayConditionally,
  validateArrayAnd,
  validateArrayOr,
  validateArrayCount,
} from '@/lib/form'

// 1. Min/Max Items Validation
validators: {
  onChange: ({ value }) => validateArrayLength(value, {
    min: 1,
    max: 10,
    minErrorKey: 'form.atLeastOnePhone',
    maxErrorKey: 'form.tooManyPhones',
  })
}

// 2. Unique Values Validation
validators: {
  onChange: ({ value }) => {
    // Simple string array
    return validateUniqueValues(value, {
      errorKey: 'form.duplicateEmails',
    })
    
    // Object array with extractor
    return validateUniqueValues(value, {
      extractor: (item) => item.email,
      errorKey: 'form.duplicateEmails',
    })
  }
}

// 3. Per-Item Validation
validators: {
  onChange: ({ value }) => validateArrayItems(value, (item, index) => {
    if (!item.name) return 'form.nameRequired'
    if (item.price < 0) return 'form.priceNegative'
    return undefined
  })
}

// 4. Overlap Validation (e.g., time ranges)
validators: {
  onChange: ({ value }) => validateNoOverlap(value, {
    extractor: (item) => ({ start: item.startTime, end: item.endTime }),
    errorKey: 'form.overlappingTimeRanges',
  })
}

// 5. Conditional Array Validation
validators: {
  onChange: ({ value, fieldApi }) => validateArrayConditionally(
    value,
    () => fieldApi.form.getFieldValue('hasPhoneSupport'),
    (array) => validateArrayLength(array, {
      min: 1,
      minErrorKey: 'form.atLeastOnePhone',
    })
  )
}

// 6. Combine Validators (AND logic)
validators: {
  onChange: ({ value }) => validateArrayAnd(value, [
    (array) => validateArrayLength(array, { min: 1, max: 10 }),
    (array) => validateUniqueValues(array),
    (array) => validateArrayItems(array, (item) => 
      item.email?.includes('@') ? undefined : 'form.invalidEmail'
    ),
  ])
}

// 7. Combine Validators (OR logic - at least one must pass)
validators: {
  onChange: ({ value }) => validateArrayOr(value, [
    (array) => validateArrayLength(array, { min: 1 }),
    (array) => array.some(item => item.isPrimary) ? undefined : 'form.noPrimary',
  ])
}

// 8. Count Validation (e.g., exactly one primary)
validators: {
  onChange: ({ value }) => validateArrayCount(value, {
    extractor: (item) => item.isPrimary,
    matchValue: true,
    exactCount: 1,
    errorKey: 'form.exactlyOnePrimary',
  })
}
```

**Common Use Cases:**

**1. Contact List with Multiple Phones:**
```tsx
<form.AppField
  name="phoneNumbers"
  validators={{
    onChange: ({ value }) => validateArrayAnd(value, [
      (array) => validateArrayLength(array, {
        min: 1,
        minErrorKey: 'form.atLeastOnePhone',
      }),
      (array) => validateUniqueValues(array, {
        extractor: (item) => item.number,
        errorKey: 'form.duplicatePhones',
      }),
    ]),
  }}
>
  {(field) => (
    <field.FieldArray
      label="Phone Numbers"
      addButtonLabel="Add Phone"
      minItems={1}
      maxItems={3}
      defaultValue={{ number: '', type: 'mobile', isPrimary: false }}
      render={(item, operations, index) => (
        <div className="flex gap-2">
          <input
            type="tel"
            value={item.number}
            onChange={(e) => {
              const updated = [...field.state.value]
              updated[index] = { ...item, number: e.target.value }
              field.handleChange(updated)
            }}
            placeholder="+1 (555) 123-4567"
            className="input flex-1"
          />
          <label className="flex items-center gap-2">
            <input
              type="checkbox"
              checked={item.isPrimary}
              onChange={(e) => {
                const updated = field.state.value.map((phone, i) => ({
                  ...phone,
                  isPrimary: i === index ? e.target.checked : false,
                }))
                field.handleChange(updated)
              }}
            />
            Primary
          </label>
        </div>
      )}
    />
  )}
</form.AppField>
```

**2. Order Items with Quantity and Price:**
```tsx
<form.AppField
  name="orderItems"
  validators={{
    onChange: ({ value }) => validateArrayAnd(value, [
      (array) => validateArrayLength(array, { min: 1 }),
      (array) => validateArrayItems(array, (item) => {
        if (!item.productId) return 'form.productRequired'
        if (item.quantity <= 0) return 'form.quantityInvalid'
        if (item.price < 0) return 'form.priceInvalid'
        return undefined
      }),
    ]),
  }}
>
  {(field) => (
    <field.FieldArray
      label="Order Items"
      addButtonLabel="Add Item"
      reorderable
      defaultValue={{ productId: '', quantity: 1, price: 0 }}
      render={(item, operations, index) => (
        <div className="grid grid-cols-3 gap-2">
          <select
            value={item.productId}
            onChange={(e) => {
              const updated = [...field.state.value]
              updated[index] = { ...item, productId: e.target.value }
              field.handleChange(updated)
            }}
            className="select"
          >
            <option value="">Select Product</option>
            {products.map((p) => (
              <option key={p.id} value={p.id}>{p.name}</option>
            ))}
          </select>
          <input
            type="number"
            min="1"
            value={item.quantity}
            onChange={(e) => {
              const updated = [...field.state.value]
              updated[index] = { ...item, quantity: parseInt(e.target.value) }
              field.handleChange(updated)
            }}
            placeholder="Qty"
            className="input"
          />
          <input
            type="number"
            min="0"
            step="0.01"
            value={item.price}
            onChange={(e) => {
              const updated = [...field.state.value]
              updated[index] = { ...item, price: parseFloat(e.target.value) }
              field.handleChange(updated)
            }}
            placeholder="Price"
            className="input"
          />
        </div>
      )}
    />
  )}
</form.AppField>
```

**3. Time Range Scheduling with Overlap Validation:**
```tsx
<form.AppField
  name="timeRanges"
  validators={{
    onChange: ({ value }) => validateArrayAnd(value, [
      (array) => validateArrayLength(array, { min: 1, max: 10 }),
      (array) => validateNoOverlap(array, {
        extractor: (item) => ({
          start: new Date(`2024-01-01T${item.startTime}`),
          end: new Date(`2024-01-01T${item.endTime}`),
        }),
        errorKey: 'form.overlappingTimeRanges',
      }),
      (array) => validateArrayItems(array, (item) => {
        const start = new Date(`2024-01-01T${item.startTime}`)
        const end = new Date(`2024-01-01T${item.endTime}`)
        if (end <= start) return 'form.endBeforeStart'
        return undefined
      }),
    ]),
  }}
>
  {(field) => (
    <field.FieldArray
      label="Available Time Slots"
      addButtonLabel="Add Time Slot"
      reorderable
      defaultValue={{ startTime: '09:00', endTime: '17:00' }}
      render={(item, operations, index) => (
        <div className="flex gap-2">
          <input
            type="time"
            value={item.startTime}
            onChange={(e) => {
              const updated = [...field.state.value]
              updated[index] = { ...item, startTime: e.target.value }
              field.handleChange(updated)
            }}
            className="input"
          />
          <span className="flex items-center">to</span>
          <input
            type="time"
            value={item.endTime}
            onChange={(e) => {
              const updated = [...field.state.value]
              updated[index] = { ...item, endTime: e.target.value }
              field.handleChange(updated)
            }}
            className="input"
          />
        </div>
      )}
    />
  )}
</form.AppField>
```

**Drag-and-Drop Behavior:**
- **Desktop**: Grab handle appears on hover (GripVertical icon)
- **Mobile**: Tap and hold (300ms) to start dragging
- **Keyboard**: Tab to handle → Space/Enter to grab → Arrow Up/Down to move → Space/Enter to drop
- **Visual Feedback**: Item dims and transforms during drag
- **Smooth Animations**: 300ms easing transitions
- **RTL Support**: Handle positioned correctly for Arabic layout

**Translation Keys Used:**
- `common.array.addItem`: "Add Item"
- `common.array.removeItem`: "Remove Item"
- `common.array.moveUp`: "Move Up"
- `common.array.moveDown`: "Move Down"
- `common.array.noItems`: "No items yet"
- `common.array.maxItemsReached`: "Maximum items reached"
- `common.array.dragToReorder`: "Drag to reorder"
- `common.array.item`: "Item {{index}}"
- `errors.form.minItemsRequired`: "At least {{min}} item(s) required."
- `errors.form.maxItemsExceeded`: "Maximum {{max}} item(s) allowed."
- `errors.form.duplicateValues`: "Duplicate values are not allowed."
- `errors.form.overlappingRanges`: "Items have overlapping ranges."
- `errors.form.invalidItem`: "Item {{index}} is invalid."

**Accessibility:**
- `role="list"` on container
- `role="listitem"` on each item
- `aria-label` on drag handle: "Drag to reorder"
- `aria-label` on remove button: "Remove item {{index}}"
- `aria-live="polite"` for add/remove announcements
- Keyboard navigation (Tab, Space/Enter, Arrow keys)
- Focus management (focus returns to appropriate element after operations)
- Screen reader announcements for all operations

#### `<form.FileUploadField>`

Generic file upload field with drag-and-drop, progress tracking, mobile camera support, and automatic thumbnail generation. Handles both file creation (File[]) and asset updates (AssetReference[]) seamlessly.

**Basic File Upload:**
```tsx
<form.AppField
  name="documents"
  validators={{
    onBlur: fileSchema({
      maxSize: '5MB',
      maxFiles: 3,
    }),
  }}
>
  {(field) => (
    <field.FileUploadField
      label={t('documents.upload')}
      accept=".pdf,.doc,.docx"
      maxFiles={3}
      maxSize="5MB"
      multiple
    />
  )}
</form.AppField>
```

**Single Business Logo (with Update Mode):**
```tsx
import { BusinessContext } from '@/contexts/BusinessContext'

function BusinessSettings() {
  const { business } = useBusiness()
  const form = useKyoraForm({
    defaultValues: {
      logo: business.logo || null, // AssetReference | null
      name: business.name,
    },
    onSubmit: async ({ value }) => {
      await api.updateBusiness(value)
    },
  })

  return (
    <BusinessContext.Provider value={business.descriptor}>
      <form.AppForm>
        <form.FormRoot>
          <form.AppField
            name="logo"
            validators={{
              onBlur: businessLogoSchema(), // Single image, 2MB
            }}
          >
            {(field) => (
              <field.FileUploadField
                label={t('business.logo')}
                accept="image/*"
                maxFiles={1}
                maxSize="2MB"
              />
            )}
          </form.AppField>
          
          <form.SubmitButton>
            {t('common.save')}
          </form.SubmitButton>
        </form.FormRoot>
      </form.AppForm>
    </BusinessContext.Provider>
  )
}
```

**Multiple Product Photos (Reorderable):**
```tsx
<form.AppField
  name="photos"
  validators={{
    onBlur: productPhotosSchema(), // 2-10 images, 10MB each
  }}
>
  {(field) => (
    <field.FileUploadField
      label={t('product.photos')}
      description={t('product.photos_description')}
      maxFiles={10}
      maxSize="10MB"
      reorderable
      multiple
      required
    />
  )}
</form.AppField>
```

**Props:**
```ts
interface FileUploadFieldProps {
  label?: string                   // Field label
  description?: string             // Helper text
  accept?: string                  // MIME types or extensions
  maxFiles?: number                // Maximum number of files
  maxSize?: string                 // Max file size (e.g., "10MB")
  multiple?: boolean               // Allow multiple files
  reorderable?: boolean            // Enable drag-drop reordering
  required?: boolean               // Visual indicator
  disabled?: boolean               // Disable field
  onUploadComplete?: (refs: AssetReference[]) => void // Upload callback
}
```

**Features:**
- **Optimistic UI**: Shows preview immediately, uploads in background
- **Mode Detection**: Automatically detects File[] vs AssetReference[] mode
- **Progress Tracking**: Real-time upload progress with cancel support
- **Mobile Optimized**: Camera/gallery buttons (50px touch targets)
- **Thumbnail Generation**: Automatic image thumbnails (300x300px WebP)
- **Video Support**: FFmpeg-based video thumbnail extraction
- **Concurrent Uploads**: Max 3 simultaneous uploads with queuing
- **Error Recovery**: Retry failed uploads, clear validation errors
- **Drag-Drop**: Native file drop zone with visual feedback
- **Accessible**: ARIA labels, keyboard navigation, screen reader support

**Upload Flow:**
1. User selects/drops files → Client-side validation
2. Optimistic preview shown immediately
3. Background upload starts (progress tracked)
4. Thumbnail generated and uploaded
5. Field value updates with AssetReference[]
6. Form can be submitted anytime

**Validation Example:**
```tsx
// Using preset schema
validators={{
  onBlur: productPhotosSchema(), // 2-10 images, 10MB
}}

// Custom validation
validators={{
  onBlur: fileSchema({
    accept: ['image/jpeg', 'image/png', 'application/pdf'],
    maxSize: '5MB',
    maxFiles: 5,
    minFiles: 1,
  }),
}}

// Array validation
validators={{
  onBlur: assetReferenceArraySchema({
    min: 2,
    max: 10,
  }),
}}
```

**Business Context Requirement:**
The upload system requires `businessDescriptor` from context to generate upload URLs:

```tsx
import { BusinessContext } from '@/contexts/BusinessContext'

// Provide context at component root
<BusinessContext.Provider value={business.descriptor}>
  <form.AppForm>
    {/* FileUploadField components */}
  </form.AppForm>
</BusinessContext.Provider>
```

#### `<form.ImageUploadField>`

Specialized variant of FileUploadField for image-only uploads with smart defaults.

**Single Image Upload:**
```tsx
<form.AppField
  name="avatar"
  validators={{
    onBlur: imageSchema({
      maxSize: '2MB',
    }),
  }}
>
  {(field) => (
    <field.ImageUploadField
      label={t('user.avatar')}
      single
      maxSize="2MB"
    />
  )}
</form.AppField>
```

**Multiple Images (Gallery):**
```tsx
<form.AppField
  name="gallery"
  validators={{
    onBlur: imageSchema({
      maxSize: '10MB',
      maxFiles: 20,
    }),
  }}
>
  {(field) => (
    <field.ImageUploadField
      label={t('product.gallery')}
      maxFiles={20}
      reorderable
    />
  )}
</form.AppField>
```

**Props:**
```ts
interface ImageUploadFieldProps {
  single?: boolean                 // Single image (maxFiles=1)
  // ... all FileUploadFieldProps except 'accept' and 'multiple'
}
```

**Features:**
- Pre-configured for images: `image/jpeg,image/jpg,image/png,image/webp,image/heic,image/heif`
- Default maxFiles: 10 (or 1 if single=true)
- Automatic WebP thumbnail generation
- Mobile camera/gallery access
- All FileUploadField features

**Preset Schemas:**
```tsx
// Business logo (single, 2MB)
import { businessLogoSchema } from '@/schemas/upload'

// Product photos (2-10 images, 10MB each)
import { productPhotosSchema } from '@/schemas/upload'

// Variant photos (1-5 images, 10MB each)
import { variantPhotosSchema } from '@/schemas/upload'
```

**Within FieldArray (Nested Upload):**
```tsx
<form.AppField name="variants">
  {(field) => (
    <field.FieldArray>
      {({ fields, pushItem, removeItem }) => (
        <>
          {fields.map((item, i) => (
            <div key={item.id}>
              <form.AppField name={`variants[${i}].name`}>
                {(field) => <field.TextField label="Variant Name" />}
              </form.AppField>
              
              <form.AppField
                name={`variants[${i}].photos`}
                validators={{
                  onBlur: variantPhotosSchema(), // 1-5 images
                }}
              >
                {(field) => (
                  <field.ImageUploadField
                    label="Variant Photos"
                    maxFiles={5}
                  />
                )}
              </form.AppField>
              
              <button onClick={() => removeItem(i)}>
                Remove Variant
              </button>
            </div>
          ))}
          
          <button onClick={() => pushItem({ name: '', photos: [] })}>
            Add Variant
          </button>
        </>
      )}
    </field.FieldArray>
  )}
</form.AppField>
```

**BottomSheet Integration:**
```tsx
import { BottomSheet } from '@/components/molecules/BottomSheet'

function AddProductModal({ isOpen, onClose }) {
  const form = useKyoraForm({
    defaultValues: {
      name: '',
      photos: [],
    },
    onSubmit: async ({ value }) => {
      await api.createProduct(value)
      onClose()
    },
  })

  return (
    <BottomSheet
      isOpen={isOpen}
      onClose={onClose}
      title={t('product.add')}
    >
      <form.AppForm>
        <form.FormRoot className="p-4 space-y-4">
          <form.AppField name="name">
            {(field) => (
              <field.TextField
                label={t('product.name')}
                required
              />
            )}
          </form.AppField>
          
          <form.AppField
            name="photos"
            validators={{
              onBlur: productPhotosSchema(),
            }}
          >
            {(field) => (
              <field.ImageUploadField
                label={t('product.photos')}
                maxFiles={10}
                reorderable
              />
            )}
          </form.AppField>
          
          {/* Sticky footer with progress */}
          <div className="sticky bottom-0 bg-base-100 p-4 border-t">
            <form.Subscribe selector={(state) => state.isSubmitting}>
              {(isSubmitting) => (
                <form.SubmitButton className="w-full" disabled={isSubmitting}>
                  {isSubmitting ? t('common.uploading') : t('common.save')}
                </form.SubmitButton>
              )}
            </form.Subscribe>
          </div>
        </form.FormRoot>
      </form.AppForm>
    </BottomSheet>
  )
}
```

**Translation Keys:**
```yaml
# Upload namespace (upload.json)
upload:
  dropzone_title: "Drop files here"
  dropzone_or: "or"
  choose_files: "Choose Files"
  choose_file: "Choose File"
  take_picture: "Take Picture"
  choose_from_gallery: "Choose from Gallery"
  uploading: "Uploading..."
  uploading_thumbnail: "Uploading thumbnail..."
  upload_failed: "Upload failed"
  retry: "Retry"
  remove: "Remove"
  drag_to_reorder: "Drag to reorder"
  max_files_reached: "Maximum {{max}} files allowed"
  filesUploading_one: "{{count}} file uploading"
  filesUploading_other: "{{count}} files uploading"
  
# Validation errors (errors.validation)
errors:
  validation:
    file_too_large: "File size exceeds {{maxSize}}"
    invalid_file_type: "Invalid file type. Accepted: {{accept}}"
    too_many_files: "Maximum {{max}} files allowed"
    min_files_required: "At least {{min}} file(s) required"
```

**Troubleshooting:**

**Issue:** "businessDescriptor is required"
```tsx
// ❌ Missing context
<form.AppForm>
  <form.AppField name="photos">
    {(field) => <field.ImageUploadField />}
  </form.AppField>
</form.AppForm>

// ✅ Provide context
import { BusinessContext } from '@/contexts/BusinessContext'

<BusinessContext.Provider value={business.descriptor}>
  <form.AppForm>
    <form.AppField name="photos">
      {(field) => <field.ImageUploadField />}
    </form.AppField>
  </form.AppForm>
</BusinessContext.Provider>
```

**Issue:** FFmpeg not loading for video thumbnails
```bash
# Ensure FFmpeg WASM files are in public/
cp node_modules/@ffmpeg/core/dist/ffmpeg-core.js public/
cp node_modules/@ffmpeg/core/dist/ffmpeg-core.wasm public/
```

**Issue:** Thumbnails not generating
```tsx
// Check browser support (WebP fallback to JPEG)
// Thumbnails: 300x300px, 0.8 quality, WebP/JPEG format
// Configurable in /lib/upload/constants.ts
```

**Issue:** Upload progress not showing
```tsx
// Upload progress is tracked automatically via useFileUpload
// Ensure field value updates trigger re-render
<form.Subscribe selector={(state) => state.values.photos}>
  {(photos) => <div>Uploaded: {photos?.length || 0}</div>}
</form.Subscribe>
```

**Issue:** Mobile camera not working
```tsx
// Ensure accept includes image types
<field.ImageUploadField
  accept="image/*"  // Enables camera/gallery
  maxFiles={1}
/>

// For generic files, camera won't appear (correct behavior)
<field.FileUploadField
  accept=".pdf,.doc"  // No camera for documents
/>
```

**Best Practices:**
1. **Always provide BusinessContext** at component root
2. **Use preset schemas** (businessLogoSchema, productPhotosSchema) for consistency
3. **Enable reorderable** for multiple images (better UX)
4. **Set appropriate maxSize** based on file type (2MB logos, 10MB products)
5. **Show upload progress** in BottomSheet/Modal footers
6. **Handle failed uploads** with retry button (automatic)
7. **Test mobile experience** with real devices (camera/gallery)
8. **Validate both files and references** depending on mode

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
