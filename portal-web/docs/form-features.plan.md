# Form Features Implementation Plan

## Executive Summary

This document tracks the implementation status of the Kyora Portal Web form system, documenting what has been built and outlining future enhancements. The system is built on **TanStack Form v1** with a custom `useKyoraForm` composition layer that eliminates boilerplate while providing production-grade form handling.

This plan outlines the implementation of 4 critical advanced form features for the Kyora Portal Web form system. The approach prioritizes **optimal UX and best code quality** over quick delivery, ensuring all new components follow the established patterns of TanStack Form + Zod + i18n + daisyUI.

---

## ✅ Phase 0: Core Form System (COMPLETED)

### Architecture Foundation

**Core Stack:**
- **TanStack Form v1**: Granular state management with Subscribe pattern
- **Zod**: Runtime validation with translation key error messages
- **react-i18next**: Automatic error translation via 'errors' namespace
- **daisyUI + Tailwind CSS v4**: Component classes with RTL support
- **useKyoraForm**: Composition layer providing pre-bound components

### Implemented Components

**Field Components (via useKyoraForm):**
- ✅ `TextField`: Text/email/tel inputs with auto error handling
- ✅ `PasswordField`: Password input with visibility toggle
- ✅ `TextareaField`: Multi-line text input with character counter
- ✅ `SelectField`: Dropdown select with search support
- ✅ `CheckboxField`: Checkbox with label and description
- ✅ `RadioField`: Radio button group with flexible layout
- ✅ `ToggleField`: Toggle/switch component

**Form Components:**
- ✅ `FormRoot`: Replaces `<form>` with auto submit handling
- ✅ `SubmitButton`: Submit button with loading state
- ✅ `ErrorInfo`: Field-level error display
- ✅ `FormError`: Form-level error display

### Core Features Implemented

**Validation System:**
- ✅ Field-level Zod schema validation
- ✅ Progressive validation modes (submit → blur → change)
- ✅ Cross-field validation support
- ✅ Server error injection with RFC7807 support
- ✅ Auto-translated error messages

**User Experience:**
- ✅ Automatic focus on first invalid field
- ✅ Granular subscriptions prevent unnecessary re-renders
- ✅ Zero boilerplate with pre-bound components
- ✅ Type-safe with full TypeScript inference
- ✅ Mobile-optimized touch targets (50px minimum)
- ✅ RTL-first design with logical properties

**Accessibility:**
- ✅ ARIA attributes auto-applied
- ✅ Screen reader support
- ✅ Keyboard navigation
- ✅ Focus management
- ✅ Error announcements

### Architecture Patterns

**Component Structure:**
```
useKyoraForm (composition layer)
  ↓ returns
form.AppForm (context provider)
  ↓ wraps
form.AppField (field context)
  ↓ renders
field.TextField / field.PasswordField / etc.
```

**Critical Rule:**
ALL components using form context (FormRoot, SubmitButton, FormError, Subscribe) MUST be inside `<form.AppForm>` wrapper.

---

## 🚧 Future Enhancements (PLANNED)

---

## Phase 1: Multi-Select Support in FormSelect (PLANNED)

### Current Status

**SelectField exists** with single-select and search support. Multi-select mode is NOT yet implemented.

### Architecture

**Goal:** Extend existing `SelectField` component to support multi-select mode with chip-based UI.

**Component Structure:**
```
SelectField (composition - existing)
  ↓ wraps
form.Field (TanStack Form binding)
  ↓ renders
Multi-select chips UI (NEW)
```

### Implementation Steps

#### Step 1.1: Enhance FormSelect Atomic Component

**File:** `portal-web/src/components/atoms/FormSelect.tsx`

**Changes:**
1. Add `multiple?: boolean` prop to `FormSelectProps`
2. Update internal state to handle `string[]` when `multiple=true`
3. Add multi-select chip display below select element:
   - Show selected values as removable chips/tags
   - Support keyboard removal (Backspace, Delete)
   - Support click removal
   - Ensure RTL layout (chips flow right-to-left)
4. Update ARIA attributes for multi-select mode:
   - `aria-multiselectable="true"`
   - `aria-selected` on options
   - Screen reader announcements for selections
5. Add daisyUI styling for chips:
   - Use `badge` component with `badge-neutral` variant
   - Add close button with `btn-xs btn-circle` styling

**Validation Patterns:**
```typescript
// Array validation for multi-select
const categorySchema = z.array(z.string())
  .min(1, 'errors.form.selectAtLeastOne')
  .max(5, 'errors.form.selectTooMany')

// Unique values validation (built-in for select)
const tagsSchema = z.array(z.string())
  .refine(arr => new Set(arr).size === arr.length, 'errors.form.duplicateSelection')
```

**Translation Keys:**
```typescript
// Add to portal-web/src/i18n/en/errors.json
{
  "form": {
    "selectAtLeastOne": "Please select at least one option",
    "selectTooMany": "You can select maximum {{max}} options",
    "duplicateSelection": "Duplicate selections are not allowed"
  }
}
```

#### Step 1.2: Update SelectField Composition Component

**File:** `portal-web/src/components/molecules/SelectField.tsx`

**API:**
```typescript
interface SelectFieldProps<T> {
  form: ReturnType<typeof useKyoraForm>;
  name: string;
  label: string;
  options: Array<{ value: string; label: string }>;
  placeholder?: string;
  multiple?: boolean; // NEW
  required?: boolean;
  disabled?: boolean;
  helperText?: string;
  validator?: (value: T) => ValidationError | undefined;
}
```

**Example Usage:**
```typescript
// Single select (existing behavior)
<form.SelectField
  name="category"
  label={t('customer.category')}
  options={categoryOptions}
  validator={val => categorySchema.safeParse(val)}
/>

// Multi-select (new behavior)
<form.SelectField
  name="tags"
  label={t('customer.tags')}
  options={tagOptions}
  multiple
  validator={val => z.array(z.string()).min(1).max(5).safeParse(val)}
/>
```

#### Step 1.3: UX Enhancements

**Keyboard Navigation:**
- Arrow keys to navigate options
- Space/Enter to select/deselect
- Backspace to remove last selected chip
- Tab to move to next field
- Escape to close dropdown

**Touch Optimization:**
- Large touch targets (50px minimum height)
- Smooth chip removal animations
- Swipe gesture to remove chips (optional, nice-to-have)

**Visual Feedback:**
- Highlight selected options in dropdown
- Animate chip addition/removal (fade in/out)
- Show selection count when dropdown closed ("3 selected")
- Max selections warning (yellow badge when approaching limit)

#### Step 1.4: Update useKyoraForm

**File:** `portal-web/src/hooks/useKyoraForm.ts`

Add `SelectField` to pre-bound components:
```typescript
return {
  FormRoot,
  TextField,
  PasswordField,
  SelectField, // NEW
  SubmitButton,
  Field: form.Field,
  Subscribe: form.Subscribe,
  useStore: form.useStore,
  // ... other methods
}
```

#### Step 1.5: Documentation

**Documentation:**
Update `FORM_SYSTEM.md`:
- Add SelectField API reference
- Add multi-select usage examples
- Add array validation patterns
- Add accessibility notes

---

## Phase 2: Date/Time Pickers (PLANNED)

### Architecture

**Goal:** Create production-grade date/time field components integrated with useKyoraForm.

**Component Structure:**
```
DateField / TimeField / DateTimeField (new field components)
  ↓ registered in useKyoraForm
form.AppField
  ↓ renders
field.DateField / field.TimeField / field.DateTimeField
  ↓ internally uses
react-day-picker (library)
```

### Implementation Steps

#### Step 2.1: Install Dependencies

```bash
cd portal-web
npm install react-day-picker date-fns
npm install -D @types/react-day-picker
```

**Rationale:**
- `react-day-picker`: Production-ready, accessible, customizable calendar
- `date-fns`: Lightweight date utilities (already used in backend)
- Full TypeScript support

#### Step 2.2: Create DatePicker Atomic Component

**File:** `portal-web/src/components/atoms/DatePicker.tsx`

**Features:**
1. **Base Calendar UI:**
   - Popup calendar on input click
   - daisyUI styled input (`input input-bordered`)
   - Calendar styled with daisyUI colors (primary for selected date)
   - RTL support for Arabic (calendar flows right-to-left)
   - Consistent UI Design with Other UI elements and form elements.

2. **Input Formatting:**
   - Locale-aware display (en: MM/DD/YYYY, ar: DD/MM/YYYY)
   - User can type or select from calendar
   - Invalid input shows error state (red border)
   - Clear button (×) to reset value

3. **Validation Props:**
   ```typescript
   interface DatePickerProps {
     value?: Date;
     onChange: (date?: Date) => void;
     minDate?: Date;
     maxDate?: Date;
     disabledDates?: Date[];
     required?: boolean;
     disabled?: boolean;
     error?: string;
     placeholder?: string;
     // ... other atomic props
   }
   ```

4. **Keyboard Navigation:**
   - Arrow keys to navigate calendar days
   - Page Up/Down for months
   - Home/End for week start/end
   - Enter to select
   - Escape to close popup
   - Tab to move to next field

5. **Accessibility:**
   - `role="dialog"` for calendar popup
   - `aria-label` with current month/year
   - `aria-selected` on active date
   - Screen reader announcements for date selection
   - Focus trap when popup open
   - Return focus to input on close

6. **Mobile Optimization:**
   - Full-screen modal on mobile (< 768px)
   - Large touch targets (48px per day)
   - Swipe gestures for month navigation
   - Bottom action bar with "Today" and "Clear" buttons

**Styling (daisyUI):**
```typescript
// Input
className="input input-bordered w-full"

// Calendar popup
className="dropdown-content bg-base-100 rounded-box shadow-lg p-4"

// Selected date
className="bg-primary text-primary-content"

// Today
className="border border-primary"

// Disabled dates
className="text-base-300 cursor-not-allowed"
```

#### Step 2.3: Create TimePicker Atomic Component

**File:** `portal-web/src/components/atoms/TimePicker.tsx`

**Features:**
1. **Time Input UI:**
   - Two numeric inputs (hours, minutes) with colon separator
   - AM/PM toggle for 12-hour format (locale-aware)
   - Scroll wheel picker on mobile (native `<input type="time">` fallback)
   - daisyUI input styling

2. **Validation:**
   - Hours: 0-23 (24h) or 1-12 (12h)
   - Minutes: 0-59
   - Auto-advance to minutes after valid hours
   - Visual error states for invalid input

3. **Keyboard Navigation:**
   - Arrow keys to increment/decrement
   - Tab to move between hours/minutes/AM-PM
   - Type to set value directly

4. **Time Format Props:**
   ```typescript
   interface TimePickerProps {
     value?: Date; // Use Date for consistency, only time portion matters
     onChange: (time?: Date) => void;
     format?: '12h' | '24h'; // Default from locale
     step?: number; // Minute increments (default 1)
     minTime?: Date;
     maxTime?: Date;
     disabled?: boolean;
     error?: string;
   }
   ```

#### Step 2.4: Create DateTimeField Composition Component

**File:** `portal-web/src/components/molecules/DateTimeField.tsx`

**API:**
```typescript
interface DateTimeFieldProps {
  form: ReturnType<typeof useKyoraForm>;
  name: string;
  label: string;
  mode?: 'date' | 'time' | 'datetime'; // Control which pickers to show
  minDate?: Date;
  maxDate?: Date;
  required?: boolean;
  disabled?: boolean;
  helperText?: string;
  validator?: (value: Date | undefined) => ValidationError | undefined;
}
```

**Layout:**
- Date-only: Single DatePicker
- Time-only: Single TimePicker  
- DateTime: DatePicker + TimePicker side-by-side (mobile: stacked)

**Example Usage:**
```typescript
// Date only
<form.DateTimeField
  name="birthdate"
  label={t('customer.birthdate')}
  mode="date"
  maxDate={new Date()} // Can't be future
  validator={val => birthdateSchema.safeParse(val)}
/>

// Date + Time
<form.DateTimeField
  name="appointmentTime"
  label={t('order.deliveryTime')}
  mode="datetime"
  minDate={new Date()} // Can't be past
  validator={val => deliveryTimeSchema.safeParse(val)}
/>
```

#### Step 2.5: Create Composition Wrappers

**Files:**
- `portal-web/src/components/molecules/DateField.tsx` (thin wrapper for mode="date")
- `portal-web/src/components/molecules/TimeField.tsx` (thin wrapper for mode="time")

**Update useKyoraForm:**
```typescript
return {
  // ... existing
  DateField,
  TimeField,
  DateTimeField,
  // ... rest
}
```

#### Step 2.3: Validation Patterns

**Common Validators:**
```typescript
// Date range validation
const birthdateSchema = z.date()
  .max(new Date(), 'errors.form.dateCannotBeFuture')
  .refine(
    date => differenceInYears(new Date(), date) >= 18,
    'errors.form.mustBe18OrOlder'
  )

// Date comparison (after another field)
const endDateSchema = z.date()
  .refine(
    (endDate, ctx) => {
      const startDate = ctx.form.getFieldValue('startDate')
      return !startDate || endDate >= startDate
    },
    'errors.form.endDateMustBeAfterStartDate'
  )

// Business hours validation
const appointmentTimeSchema = z.date()
  .refine(
    date => {
      const hours = date.getHours()
      return hours >= 9 && hours < 17 // 9 AM - 5 PM
    },
    'errors.form.outsideBusinessHours'
  )
```

**Translation Keys:**
```typescript
// Add to portal-web/src/i18n/en/errors.json
{
  "form": {
    "dateCannotBeFuture": "Date cannot be in the future",
    "dateCannotBePast": "Date cannot be in the past",
    "mustBe18OrOlder": "Must be 18 years or older",
    "endDateMustBeAfterStartDate": "End date must be after start date",
    "outsideBusinessHours": "Please select a time between 9 AM and 5 PM",
    "invalidDate": "Invalid date format",
    "invalidTime": "Invalid time format"
  }
}
```

#### Step 2.4: Testing & Documentation

**Testing:**
- Date selection and formatting
- Time input validation
- Min/max date enforcement
- Keyboard navigation
- Screen reader announcements
- Mobile full-screen modal
- RTL layout (Arabic)
- Locale-aware formatting

**Documentation:**
Update `FORM_SYSTEM.md`:
- Add field components to API reference
- Add date/time validation patterns
- Add usage examples
- Add accessibility notes

---

## Phase 3: Wizard/Stepper Pattern (PLANNED)

### Architecture

**Goal:** Create multi-step form pattern with step-aware validation and progress tracking.

**Component Structure:**
```
useWizardForm (custom hook)
  ↓ wraps
useKyoraForm (form management)
  ↓ with
Stepper (progress indicator component)
```

### Implementation Steps

#### Step 3.1: Create Stepper Component

**File:** `portal-web/src/components/molecules/Stepper.tsx`

**Features:**
1. **Visual Progress:**
   - Numbered steps with labels
   - Active step highlighted (primary color)
   - Completed steps (checkmark icon, primary color)
   - Future steps (gray, disabled)
   - Progress line connecting steps

2. **Responsive Layout:**
   - Desktop: Horizontal steps across top
   - Mobile: Vertical steps on left side OR compact progress bar (Which ever provides the best mobile and small screen sizes experience)

3. **Interactivity:**
   - Click completed steps to navigate back
   - Cannot click future steps (validation required)
   - Keyboard navigation (Tab, Enter)

**API:**
```typescript
interface Step {
  id: string;
  label: string;
  description?: string;
  icon?: React.ComponentType;
}

interface StepperProps {
  steps: Step[];
  currentStep: number;
  completedSteps: Set<number>;
  onStepClick: (stepIndex: number) => void;
  orientation?: 'horizontal' | 'vertical';
}
```

**daisyUI Styling:**
```typescript
// Use steps component
<ul className="steps steps-horizontal">
  <li className="step step-primary">Register</li>
  <li className="step step-primary">Choose plan</li>
  <li className="step">Purchase</li>
  <li className="step">Receive Product</li>
</ul>
```

#### Step 3.2: Create useWizardForm Hook

**File:** `portal-web/src/hooks/useWizardForm.ts`

**Features:**
1. **Step Management:**
   - Initialize with step definitions
   - Track current step index
   - Track completed steps (Set)
   - Navigate forward/backward
   - Jump to specific step (if completed)

2. **Step-Aware Validation:**
   - Validate only current step fields on "Next"
   - Skip validation on "Back"
   - Block navigation to future steps if current invalid
   - Validate all steps on final submit

3. **Progress Persistence:**
   - Serialize step state to URL search params (`?step=<step_name>`)
   - Restore on page load
   - Optionally persist form values to localStorage
   - Clear on final submission or explicit reset

4. **TanStack Router Integration:**
   - Use `useNavigate` and `useSearch` for step param
   - Update URL on step change
   - Deep-linkable (share URL to specific step)

**API:**
```typescript
interface WizardStep {
  id: string;
  label: string;
  fields: string[]; // Field names to validate for this step
  validate?: (values: Record<string, any>) => ValidationError | undefined;
}

interface UseWizardFormOptions {
  steps: WizardStep[];
  onSubmit: (values: Record<string, any>) => Promise<void>;
  initialStep?: number;
  persistProgress?: boolean; // Save to localStorage
}

function useWizardForm(options: UseWizardFormOptions) {
  const form = useKyoraForm({ onSubmit: options.onSubmit })
  
  return {
    ...form, // All useKyoraForm methods
    
    // Wizard-specific
    currentStep: string,
    totalSteps: number,
    completedSteps: Set<string>,
    isFirstStep: boolean,
    isLastStep: boolean,
    canGoNext: boolean, // Current step valid
    canGoBack: boolean, // can be globally disabled
    
    nextStep: () => Promise<void>, // Validates + navigates
    prevStep: () => void,
    goToStep: (index: number) => void,
    resetWizard: () => void,
  }
}
```

**Example Usage:**
```typescript
const OnboardingWizard = () => {
  const wizard = useWizardForm({
    steps: [
      { id: 'email', label: 'Email', fields: ['email'] },
      { id: 'business', label: 'Business', fields: ['businessName', 'country'] },
      { id: 'confirm', label: 'Confirm', fields: [] },
    ],
    onSubmit: async (values) => {
      await api.onboarding.complete(values)
    },
    persistProgress: true,
  })

  return (
    <div>
      <Stepper
        steps={wizard.steps}
        currentStep={wizard.currentStep}
        completedSteps={wizard.completedSteps}
        onStepClick={wizard.goToStep}
      />
      
      <wizard.FormRoot>
        {wizard.currentStep === 0 && (
          <wizard.TextField name="email" label="Email" />
        )}
        
        {wizard.currentStep === 1 && (
          <>
            <wizard.TextField name="businessName" label="Business Name" />
            <wizard.SelectField name="country" label="Country" options={countries} />
          </>
        )}
        
        {wizard.currentStep === 2 && (
          <ConfirmationView values={wizard.useStore(s => s.values)} />
        )}
        
        <div className="flex gap-4">
          {!wizard.isFirstStep && (
            <button onClick={wizard.prevStep}>Back</button>
          )}
          
          {!wizard.isLastStep ? (
            <button onClick={wizard.nextStep} disabled={!wizard.canGoNext}>
              Next
            </button>
          ) : (
            <wizard.SubmitButton>Complete</wizard.SubmitButton>
          )}
        </div>
      </wizard.FormRoot>
    </div>
  )
}
```

#### Step 3.3: Accessibility & UX

**Keyboard Navigation:**
- Tab through steps
- Enter to activate clickable steps
- Arrow keys to navigate between fields within step

**Screen Reader Support:**
- Announce current step ("Step 2 of 4: Business Information")
- Announce step completion
- Announce validation errors specific to current step
- `aria-current="step"` on active step

**Visual Feedback:**
- Smooth step transitions (slide animation) -- slide direction depends on page direction and screen size (right/left/up/down).
- Loading state during validation
- Success checkmarks for completed steps
- Error indicators on steps with invalid fields

**Mobile Optimization:**
- Compact stepper (dots instead of labels)
- Full-width "Next"/"Back" buttons at bottom
- Sticky stepper at top during scroll

#### Step 3.3: Testing & Documentation

**Testing:**
- Step navigation (next/back/jump)
- Step-level validation
- URL persistence and deep-linking
- Progress persistence
- Form reset on completion
- Accessibility (keyboard, screen reader)

**Documentation:**
Update `FORM_SYSTEM.md`:
- Add useWizardForm API reference
- Add Stepper component usage
- Add multi-step examples
- Add persistence patterns

---

## Phase 4: Form Arrays (Repeating Fields) (PLANNED)

### Architecture

**Goal:** Support dynamic field arrays with add/remove/reorder operations.

**Component Structure:**
```
FieldArray (new component)
  ↓ uses
form.AppField with mode="array"
  ↓ renders
Repeated field items
  ↓ with
Array operations UI
```

### Implementation Steps

#### Step 4.1: Create FieldArray Component

**File:** `portal-web/src/components/molecules/FieldArray.tsx`

**Features:**
1. **Array Operations:**
   - Add new item (+ button)
   - Remove item (× button on each item)
   - Reorder via drag-and-drop (optional, see Step 4.2)
   - Validate each item individually
   - Validate array-level (min/max items, unique values)

2. **Render Function Pattern:**
   ```typescript
   <FieldArray
     name="phoneNumbers"
     label="Phone Numbers"
     minItems={1}
     maxItems={5}
     addButtonLabel="Add Phone"
     render={(item, index, remove) => (
       <div className="flex gap-2">
         <form.TextField
           name={`phoneNumbers.${index}.number`}
           label={`Phone ${index + 1}`}
         />
         <button onClick={remove}>Remove</button>
       </div>
     )}
   />
   ```

3. **Item Animations:**
   - Fade in when added
   - Fade out when removed
   - Smooth position change when reordered

4. **Empty State:**
   - Show placeholder when array empty
   - Prominent "Add First Item" button

**API:**
```typescript
interface FieldArrayProps<T> {
  form: ReturnType<typeof useKyoraForm>;
  name: string;
  label?: string;
  minItems?: number;
  maxItems?: number;
  defaultValue?: () => T; // Factory for new items
  addButtonLabel?: string;
  emptyMessage?: string;
  reorderable?: boolean; // Enable drag-and-drop
  validator?: (items: T[]) => ValidationError | undefined;
  render: (
    item: T,
    index: number,
    operations: {
      remove: () => void;
      moveUp: () => void;
      moveDown: () => void;
    }
  ) => React.ReactNode;
}
```

#### Step 4.2: Implement Drag-and-Drop Reordering

**Library:** `@dnd-kit/core` + `@dnd-kit/sortable`

```bash
npm install @dnd-kit/core @dnd-kit/sortable @dnd-kit/utilities
```

**Features:**
1. **Visual Feedback:**
   - Drag handle icon (⋮⋮) at start of each item
   - Ghost item follows cursor during drag
   - Drop indicator (blue line) shows insertion point
   - Dragged item opacity reduced

2. **Keyboard Reordering:**
   - Space to grab item
   - Arrow keys to move up/down
   - Space to drop
   - Escape to cancel

3. **Touch Support:**
   - Long-press to initiate drag (prevent conflict with scrolling)
   - Visual feedback during touch drag

4. **Accessibility:**
   - `role="list"` on array container
   - `role="listitem"` on each item
   - `aria-grabbed` during drag
   - Screen reader announcement: "Item moved from position X to Y"

**Implementation:**
```typescript
import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
} from '@dnd-kit/core'
import {
  SortableContext,
  sortableKeyboardCoordinates,
  verticalListSortingStrategy,
  useSortable,
} from '@dnd-kit/sortable'

const FieldArray = <T,>(props: FieldArrayProps<T>) => {
  // Implementation
}
  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    })
  )

  const handleDragEnd = (event) => {
    const { active, over } = event
    if (active.id !== over.id) {
      form.swapValues(props.name, active.id, over.id)
    }
  }

  return (
    <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
      <SortableContext items={items} strategy={verticalListSortingStrategy}>
        {items.map((item, index) => (
          <SortableItem key={item.id} id={item.id} index={index} {...props} />
        ))}
      </SortableContext>
    </DndContext>
  )
}
```

#### Step 4.3: Array Validation Patterns

**Common Validators:**
```typescript
// Min/Max items
const phoneNumbersSchema = z.array(phoneSchema)
  .min(1, 'errors.form.atLeastOnePhone')
  .max(5, 'errors.form.tooManyPhones')

// Unique values (e.g., email list)
const emailsSchema = z.array(z.string().email())
  .refine(
    emails => new Set(emails).size === emails.length,
    'errors.form.duplicateEmails'
  )

// Cross-item validation (e.g., date ranges don't overlap)
const timeRangesSchema = z.array(z.object({
  start: z.date(),
  end: z.date()
}))
  .refine(
    ranges => {
      for (let i = 0; i < ranges.length; i++) {
        for (let j = i + 1; j < ranges.length; j++) {
          if (ranges[i].end > ranges[j].start && ranges[i].start < ranges[j].end) {
            return false // Overlapping ranges
          }
        }
      }
      return true
    },
    'errors.form.overlappingTimeRanges'
  )
```

**Translation Keys:**
```typescript
// Add to portal-web/src/i18n/en/errors.json
{
  "form": {
    "atLeastOnePhone": "Please add at least one phone number",
    "tooManyPhones": "Maximum {{max}} phone numbers allowed",
    "duplicateEmails": "Email addresses must be unique",
    "overlappingTimeRanges": "Time ranges cannot overlap",
    "removeItemConfirm": "Are you sure you want to remove this item?"
  }
}
```

#### Step 4.4: UX Enhancements

**Add/Remove Animations:**
```css
/* Fade in new items */
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}

.field-array-item-enter {
  animation: fadeIn 0.3s ease-out;
}

/* Fade out removed items */
@keyframes fadeOut {
  from { opacity: 1; transform: scale(1); }
  to { opacity: 0; transform: scale(0.9); }
}

.field-array-item-exit {
  animation: fadeOut 0.2s ease-in;
}
```

**Confirmation Dialogs:**
- Show confirmation before removing item (optional, configurable)
- "Are you sure you want to remove this phone number?"

**Empty State Design:**
```typescript
<div className="border-2 border-dashed border-base-300 rounded-lg p-8 text-center">
  <p className="text-base-content/60 mb-4">
    {props.emptyMessage || 'No items added yet'}
  </p>
  <button className="btn btn-primary btn-outline" onClick={addItem}>
    {props.addButtonLabel || 'Add Item'}
  </button>
</div>
```

**Max Items Reached:**
```typescript
{items.length >= maxItems && (
  <div className="alert alert-warning">
    <span>Maximum {maxItems} items allowed</span>
  </div>
)}
```

#### Step 4.5: Update useKyoraForm

**File:** `portal-web/src/hooks/useKyoraForm.ts`

Add `FieldArray` to pre-bound components:
```typescript
return {
  // ... existing
  FieldArray,
  // ... rest
}
```

Ensure array operations exposed:
```typescript
// TanStack Form provides these
form.pushValue(name, value) // Add item
form.removeValue(name, index) // Remove item
form.swapValues(name, indexA, indexB) // Reorder
```

#### Step 4.6: Testing & Documentation (2 days)

**Tests:**
- Add/remove items
- Drag-and-drop reordering (mouse, touch, keyboard)
- Array-level validation (min/max, unique)
- Item-level validation
- Animations
- Accessibility (screen reader, keyboard)
- Focus management after remove

**Documentation:**
Update `FORM_SYSTEM.md`:
- Add FieldArray API reference
- Add drag-and-drop usage
- Add array validation patterns
- Add examples (contact list, order items, time ranges)

---

## Success Metrics

### Code Quality Metrics
- [ ] All components follow atomic → composition pattern
- [ ] Full TypeScript type safety (no `any`)
- [ ] Zero accessibility violations (axe-core audit)
- [ ] All error messages translatable (i18n)
- [ ] RTL support verified with Arabic locale

### UX Metrics
- [ ] Keyboard navigation works for all components
- [ ] Touch targets ≥ 50px on mobile
- [ ] Form completion time reduced vs. manual entry
- [ ] Zero user complaints about date/select pickers

### Performance Metrics
- [ ] No performance regression (< 100ms render time per component)
- [ ] No unnecessary re-renders (verified with React DevTools)
- [ ] Lazy loading for heavy components (react-day-picker)

### Documentation Metrics
- [ ] API documentation complete in FORM_SYSTEM.md
- [ ] At least 3 usage examples per component
- [ ] Migration guide for existing forms
- [ ] Accessibility guidelines documented

---

## Risk Assessment

### Technical Risks

**Risk 1: react-day-picker Bundle Size**
- **Likelihood:** High
- **Impact:** Medium (adds ~50KB to bundle)
- **Mitigation:** Lazy load DatePicker component using `React.lazy()`, only load when date field rendered

---

## Future Enhancements (Out of Scope)

These features are intentionally excluded from this plan but could be added later:

1. **Rich Text Editor Field:** WYSIWYG editor for description fields (consider Tiptap or Lexical)
2. **File Upload Field:** Drag-and-drop file uploads with preview (integrate with existing blob service)
3. **Autocomplete Field:** Searchable select with async data loading (integrate with TanStack Query)
4. **Color Picker Field:** Visual color selection for branding/themes
5. **Signature Field:** Capture digital signatures (useful for contracts/agreements)
6. **Conditional Fields:** Show/hide fields based on other field values (can be done with `form.Subscribe` already)
7. **Form Templates:** Pre-built form layouts for common use cases (customer intake, order forms, etc.)

---

## Conclusion

This plan provides a comprehensive roadmap for implementing 4 critical advanced form features following the established TanStack Form + Zod + i18n + daisyUI patterns. The phased approach allows for iterative development with continuous validation of architecture decisions.

**Key Principles:**
- ✅ Optimal UX over fast delivery
- ✅ Best code quality and maintainability
- ✅ Follow existing atomic → composition pattern
- ✅ Accessibility from Day 1, not an afterthought
- ✅ RTL-first design for Arabic users
- ✅ Mobile-optimized with large touch targets
- ✅ Progressive enhancement (start simple, add features)
