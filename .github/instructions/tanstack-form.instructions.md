---
description: TanStack Form - React Form Management Library
applyTo: "portal-web/**,storefront-web/**"
---

# TanStack Form - React Form Management Library

## Overview

TanStack Form is a powerful, flexible form management library for React that provides complete control over validation, error handling, and form state without being opinionated about markup. It uses TanStack Store for reactive state management.

## Core Philosophy

- **Framework Agnostic Markup**: TanStack Form doesn't dictate your HTML structure
- **Type Safety First**: Full TypeScript support with strong type inference
- **Flexibility Over Convention**: Highly customizable validation, error handling, and submission
- **No Unnecessary Re-renders**: Uses signals/store pattern to prevent performance issues
- **Headless by Design**: Bring your own UI components

---

## Installation & Setup

```bash
npm install @tanstack/react-form
# For schema validation (optional)
npm install zod @tanstack/zod-form-adapter
```

Basic form setup:

```tsx
import { useForm } from "@tanstack/react-form";

function App() {
  const form = useForm({
    defaultValues: {
      firstName: "",
      lastName: "",
    },
    onSubmit: async ({ value }) => {
      console.log(value);
    },
  });

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault();
        form.handleSubmit();
      }}
    >
      {/* Fields */}
    </form>
  );
}
```

---

## Validation

### When Validation Runs

Control exactly when validation occurs using validator callbacks:

- `onChange`: Validate on every keystroke
- `onBlur`: Validate when field loses focus
- `onSubmit`: Validate only on form submission
- `onMount`: Validate when component mounts

**Example - Multiple Validation Timings:**

```tsx
<form.Field
  name="age"
  validators={{
    onChange: ({ value }) => (value < 13 ? "Must be 13 or older" : undefined),
    onBlur: ({ value }) => (value < 0 ? "Invalid value" : undefined),
  }}
>
  {(field) => (
    <>
      <input
        value={field.state.value}
        onChange={(e) => field.handleChange(e.target.valueAsNumber)}
        onBlur={field.handleBlur}
      />
      {!field.state.meta.isValid && (
        <em>{field.state.meta.errors.join(", ")}</em>
      )}
    </>
  )}
</form.Field>
```

### Field-Level vs Form-Level Validation

**Field-Level** (recommended for most cases):

```tsx
<form.Field
  name="email"
  validators={{
    onChange: ({ value }) =>
      !value.includes("@") ? "Invalid email" : undefined,
  }}
/>
```

**Form-Level** (for cross-field validation):

```tsx
const form = useForm({
  validators: {
    onChange({ value }) {
      if (value.age < 13) {
        return "Must be 13 or older to sign";
      }
      return undefined;
    },
  },
});

// Access form-level errors
const formErrorMap = useStore(form.store, (state) => state.errorMap);
```

**Setting Field Errors from Form Validators:**

```tsx
const form = useForm({
  validators: {
    onSubmitAsync: async ({ value }) => {
      const hasErrors = await verifyDataOnServer(value);
      if (hasErrors) {
        return {
          form: "Invalid data", // Optional form-level error
          fields: {
            age: "Must be 13 or older",
            "socials[0].url": "Invalid URL", // Nested field
            "details.email": "Email required",
          },
        };
      }
      return null;
    },
  },
});
```

> **Important**: Field-specific validation will overwrite form-level field errors.

### Async Validation

**Built-in Debouncing:**

```tsx
<form.Field
  name="username"
  asyncDebounceMs={500} // Debounce all async validators by 500ms
  validators={{
    onChangeAsyncDebounceMs: 1500, // Override for specific validator
    onChangeAsync: async ({ value }) => {
      await new Promise((resolve) => setTimeout(resolve, 1000));
      return value.length < 3 ? "Too short" : undefined;
    },
    onBlurAsync: async ({ value }) => {
      const exists = await checkUsernameExists(value);
      return exists ? "Username taken" : undefined;
    },
  }}
/>
```

**Async Validation Behavior:**

- Sync validation runs first by default
- Async only runs if sync succeeds
- Set `asyncAlways: true` to run async regardless

### Standard Schema Validation (Zod, Valibot, ArkType, Effect/Schema)

TanStack Form supports Standard Schema libraries natively:

```tsx
import { z } from 'zod'

const userSchema = z.object({
  age: z.number().gte(13, 'Must be 13 or older'),
})

// Form-level schema
const form = useForm({
  validators: {
    onChange: userSchema,
  },
})

// Field-level schema
<form.Field
  name="age"
  validators={{
    onChange: z.number().gte(13, 'Must be 13 or older'),
    onChangeAsync: z.number().refine(
      async (value) => {
        const currentAge = await fetchCurrentAge()
        return value >= currentAge
      },
      { message: 'Can only increase age' }
    ),
  }}
/>

// Combining schema with custom logic
<form.Field
  name="age"
  validators={{
    onChangeAsync: async ({ value, fieldApi }) => {
      const errors = fieldApi.parseValueWithSchema(
        z.number().gte(13, 'Must be 13 or older')
      )
      if (errors) return errors
      // Continue with additional validation
    },
  }}
/>
```

> **Note**: Schemas validate but don't transform values. See submission handling for transformations.

### Preventing Invalid Form Submission

```tsx
<form.Subscribe
  selector={(state) => [state.canSubmit, state.isSubmitting]}
  children={([canSubmit, isSubmitting]) => (
    <button type="submit" disabled={!canSubmit}>
      {isSubmitting ? "..." : "Submit"}
    </button>
  )}
/>
```

Combine `canSubmit` with `isPristine` to prevent submission before any changes:

```tsx
disabled={!canSubmit || isPristine}
```

---

## Dynamic Validation

Change validation rules based on form state using `onDynamic`:

```tsx
import { revalidateLogic } from "@tanstack/react-form";

const form = useForm({
  validationLogic: revalidateLogic(), // Required to enable onDynamic
  validators: {
    onDynamic: ({ value }) => {
      if (!value.firstName) {
        return { firstName: "First name required" };
      }
      return undefined;
    },
  },
});
```

**Revalidation Modes:**

```tsx
revalidateLogic({
  mode: "submit", // Before first submission (default: 'submit')
  modeAfterSubmission: "blur", // After submission (default: 'change')
});
```

Available modes: `'change'`, `'blur'`, `'submit'`

**Accessing Dynamic Errors:**

```tsx
<p>{form.state.errorMap.onDynamic?.firstName}</p>
```

**With Fields:**

```tsx
<form.Field
  name="age"
  validators={{
    onDynamic: ({ value }) => (value > 18 ? undefined : "Must be over 18"),
  }}
>
  {(field) => (
    <div>
      <input
        type="number"
        onChange={(e) => field.handleChange(e.target.valueAsNumber)}
        onBlur={field.handleBlur}
      />
      <p>{field.state.meta.errorMap.onDynamic}</p>
    </div>
  )}
</form.Field>
```

**Async Dynamic Validation:**

```tsx
validators: {
  onDynamicAsyncDebounceMs: 500,
  onDynamicAsync: async ({ value }) => {
    const isValid = await validateUsername(value.username)
    return isValid ? undefined : { username: 'Already taken' }
  },
}
```

**With Standard Schemas:**

```tsx
import { z } from 'zod'

validators: {
  onDynamic: z.object({
    firstName: z.string().min(1, 'Required'),
    lastName: z.string().min(1, 'Required'),
  }),
}
```

---

## Custom Error Types

TanStack Form supports any error type (strings, numbers, booleans, objects, arrays):

### String Errors (Most Common)

```tsx
validators: {
  onChange: ({ value }) =>
    value < 13 ? 'Must be 13 or older' : undefined,
}

// Display
{field.state.meta.errors.map((error, i) => (
  <div key={i}>{error}</div>
))}
```

### Number Errors

```tsx
validators: {
  onChange: ({ value }) =>
    value < 18 ? 18 - value : undefined, // Returns number or undefined
}

// Display - TypeScript knows error is a number
<div>You need {field.state.meta.errors[0]} more years</div>
```

### Boolean Errors

```tsx
validators: {
  onChange: ({ value }) =>
    !value ? true : undefined,
}

// Display
{field.state.meta.errors[0] === true && (
  <div>You must accept terms</div>
)}
```

### Object Errors

```tsx
validators: {
  onChange: ({ value }) => {
    if (!value.includes('@')) {
      return {
        message: 'Invalid email',
        severity: 'error',
        code: 1001,
      }
    }
    return undefined
  },
}

// Display - TypeScript knows the error shape
{typeof field.state.meta.errors[0] === 'object' && (
  <div className={field.state.meta.errors[0].severity}>
    {field.state.meta.errors[0].message}
    <small>(Code: {field.state.meta.errors[0].code})</small>
  </div>
)}
```

### Array Errors

```tsx
validators: {
  onChange: ({ value }) => {
    const errors = []
    if (value.length < 8) errors.push('Too short')
    if (!/[A-Z]/.test(value)) errors.push('Missing uppercase')
    if (!/[0-9]/.test(value)) errors.push('Missing number')
    return errors.length ? errors : undefined
  },
}

// Display
{Array.isArray(field.state.meta.errors) && (
  <ul>
    {field.state.meta.errors.map((err, i) => (
      <li key={i}>{err}</li>
    ))}
  </ul>
)}
```

### Accessing Errors by Source

Use `errorMap` to access errors by validation timing:

```tsx
{
  field.state.meta.errorMap.onChange && (
    <div>{field.state.meta.errorMap.onChange}</div>
  );
}

{
  field.state.meta.errorMap.onBlur && (
    <div>{field.state.meta.errorMap.onBlur}</div>
  );
}
```

### The `disableErrorFlat` Prop

By default, errors from all validators are flattened into a single `errors` array. Use `disableErrorFlat` to preserve error sources:

```tsx
<form.Field
  name="email"
  disableErrorFlat
  validators={{
    onChange: ({ value }) =>
      !value.includes("@") ? "Invalid format" : undefined,
    onBlur: ({ value }) =>
      !value.endsWith(".com") ? "Only .com allowed" : undefined,
  }}
>
  {(field) => (
    <>
      {field.state.meta.errorMap.onChange && (
        <div className="real-time">{field.state.meta.errorMap.onChange}</div>
      )}
      {field.state.meta.errorMap.onBlur && (
        <div className="blur-feedback">{field.state.meta.errorMap.onBlur}</div>
      )}
    </>
  )}
</form.Field>
```

### Type Safety

`errorMap` keys are strongly typed to match validator return types:

```tsx
<form.Field
  name="password"
  validators={{
    onChange: ({ value }): string | undefined =>
      value.length < 8 ? 'Too short' : undefined,
    onBlur: ({ value }): { message: string, level: string } | undefined =>
      !/[A-Z]/.test(value) ? { message: 'Missing uppercase', level: 'warning' } : undefined,
  }}
>
  {(field) => {
    // TypeScript knows the exact types
    const onChangeError: string | undefined = field.state.meta.errorMap.onChange
    const onBlurError: { message: string, level: string } | undefined = field.state.meta.errorMap.onBlur

    return (/* ... */)
  }}
</form.Field>
```

---

## Arrays and Nested Fields

### Basic Array Usage

```tsx
function App() {
  const form = useForm({
    defaultValues: {
      people: [],
    },
  });

  return (
    <form.Field name="people" mode="array">
      {(field) => (
        <div>
          {field.state.value.map((_, i) => (
            <form.Field key={i} name={`people[${i}].name`}>
              {(subField) => (
                <input
                  value={subField.state.value}
                  onChange={(e) => subField.handleChange(e.target.value)}
                />
              )}
            </form.Field>
          ))}
          <button
            onClick={() => field.pushValue({ name: "", age: 0 })}
            type="button"
          >
            Add person
          </button>
        </div>
      )}
    </form.Field>
  );
}
```

### Array Field Methods

- `pushValue(value)`: Add item to end
- `insertValue(index, value)`: Insert at index
- `removeValue(index)`: Remove item
- `swapValues(indexA, indexB)`: Swap positions
- `moveValue(from, to)`: Move item

---

## Linked Fields

Revalidate one field when another changes:

```tsx
<form.Field name="password">
  {(field) => (
    <input
      value={field.state.value}
      onChange={(e) => field.handleChange(e.target.value)}
    />
  )}
</form.Field>

<form.Field
  name="confirm_password"
  validators={{
    onChangeListenTo: ['password'], // Revalidate when password changes
    onChange: ({ value, fieldApi }) => {
      if (value !== fieldApi.form.getFieldValue('password')) {
        return 'Passwords do not match'
      }
      return undefined
    },
  }}
>
  {(field) => (
    <div>
      <input
        value={field.state.value}
        onChange={(e) => field.handleChange(e.target.value)}
      />
      {field.state.meta.errors.map((err) => (
        <div key={err}>{err}</div>
      ))}
    </div>
  )}
</form.Field>
```

Also supports `onBlurListenTo` for blur-triggered revalidation.

---

## Reactivity

TanStack Form doesn't cause re-renders by default. Subscribe to values using:

### `useStore` Hook

Perfect for accessing values in component logic:

```tsx
import { useStore } from "@tanstack/react-store";

const firstName = useStore(form.store, (state) => state.values.firstName);
const errors = useStore(form.store, (state) => state.errorMap);
```

> **Warning**: `useStore` causes component re-renders. Always use a selector to limit subscriptions.

### `form.Subscribe` Component

Best for UI reactivity (doesn't trigger component-level re-renders):

```tsx
<form.Subscribe
  selector={(state) => state.values.firstName}
  children={(firstName) => (
    <form.Field>
      {(field) => (
        <input
          name="lastName"
          value={field.state.lastName}
          onChange={field.handleChange}
        />
      )}
    </form.Field>
  )}
/>
```

**Rule of Thumb**: Use `form.Subscribe` for UI, `useStore` for logic.

---

## Listeners (Side Effects)

React to form/field events with listeners (onChange, onBlur, onMount, onSubmit):

### Field Listeners

```tsx
<form.Field
  name="country"
  listeners={{
    onChange: ({ value }) => {
      console.log(`Country changed to: ${value}`);
      form.setFieldValue("province", ""); // Reset dependent field
    },
    onChangeDebounceMs: 500, // Debounce for 500ms
  }}
>
  {(field) => (
    <input
      value={field.state.value}
      onChange={(e) => field.handleChange(e.target.value)}
    />
  )}
</form.Field>
```

### Form Listeners

```tsx
const form = useForm({
  listeners: {
    onMount: ({ formApi }) => {
      loggingService("mount", formApi.state.values);
    },
    onChange: ({ formApi, fieldApi }) => {
      // fieldApi is the field that triggered the change
      if (formApi.state.isValid) {
        formApi.handleSubmit(); // Autosave
      }
    },
    onChangeDebounceMs: 500,
  },
});
```

**Available Events:**

- `onMount`: Form/field mounted
- `onChange`: Value changed
- `onBlur`: Field blurred
- `onSubmit`: Form submitted

---

## Async Initial Values

Integrate with TanStack Query for async initial values:

```tsx
import { useForm } from "@tanstack/react-form";
import { useQuery } from "@tanstack/react-query";

export default function App() {
  const { data, isLoading } = useQuery({
    queryKey: ["data"],
    queryFn: async () => {
      await new Promise((resolve) => setTimeout(resolve, 1000));
      return { firstName: "John", lastName: "Doe" };
    },
  });

  const form = useForm({
    defaultValues: {
      firstName: data?.firstName ?? "",
      lastName: data?.lastName ?? "",
    },
    onSubmit: async ({ value }) => {
      console.log(value);
    },
  });

  if (isLoading) return <p>Loading...</p>;

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault();
        form.handleSubmit();
      }}
    >
      {/* Fields */}
    </form>
  );
}
```

---

## Focus Management

TanStack Form doesn't manage focus automatically (philosophy: no markup opinions), but you can implement it easily:

### React DOM

```tsx
const form = useForm({
  onSubmitInvalid() {
    const invalidInput = document.querySelector(
      '[aria-invalid="true"]'
    ) as HTMLInputElement
    invalidInput?.focus()
  },
})

// In field
<input
  aria-invalid={!field.state.meta.isValid && field.state.meta.isTouched}
  // ...
/>
```

### React Native

```tsx
import { useRef } from 'react'
import { TextInput } from 'react-native'

const fields = useRef([] as Array<{ input: TextInput; name: string }>)

const form = useForm({
  onSubmitInvalid({ formApi }) {
    const errorMap = formApi.state.errorMap.onChange
    const inputs = fields.current

    let firstInput
    for (const input of inputs) {
      if (errorMap[input.name]) {
        firstInput = input.input
        break
      }
    }
    firstInput?.focus()
  },
})

// In field
<TextInput
  ref={(input) => {
    fields.current[0] = { input, name: field.name }
  }}
  // ...
/>
```

---

## Form Composition

### Custom Form Hooks

Create reusable form hooks with pre-bound components:

```tsx
import { createFormHookContexts, createFormHook } from "@tanstack/react-form";

// 1. Create contexts
export const { fieldContext, formContext, useFieldContext, useFormContext } =
  createFormHookContexts();

// 2. Create custom field components
function TextField({ label }: { label: string }) {
  const field = useFieldContext<string>();
  return (
    <label>
      <span>{label}</span>
      <input
        value={field.state.value}
        onChange={(e) => field.handleChange(e.target.value)}
      />
    </label>
  );
}

function SubscribeButton({ label }: { label: string }) {
  const form = useFormContext();
  return (
    <form.Subscribe selector={(state) => state.isSubmitting}>
      {(isSubmitting) => (
        <button type="submit" disabled={isSubmitting}>
          {label}
        </button>
      )}
    </form.Subscribe>
  );
}

// 3. Create custom form hook
const { useAppForm, withForm } = createFormHook({
  fieldContext,
  formContext,
  fieldComponents: {
    TextField,
  },
  formComponents: {
    SubscribeButton,
  },
});

// 4. Use in app
function App() {
  const form = useAppForm({
    defaultValues: {
      firstName: "John",
      lastName: "Doe",
    },
  });

  return (
    <form.AppForm>
      <form.AppField
        name="firstName"
        children={(field) => <field.TextField label="First Name" />}
      />
      <form.SubscribeButton label="Submit" />
    </form.AppForm>
  );
}
```

### Breaking Forms into Smaller Pieces

Use `withForm` for large forms:

```tsx
const ChildForm = withForm({
  defaultValues: {
    firstName: "John",
    lastName: "Doe",
  },
  props: {
    title: "Child Form",
  },
  render: function Render({ form, title }) {
    return (
      <div>
        <p>{title}</p>
        <form.AppField
          name="firstName"
          children={(field) => <field.TextField label="First Name" />}
        />
        <form.AppForm>
          <form.SubscribeButton label="Submit" />
        </form.AppForm>
      </div>
    );
  },
});

function App() {
  const form = useAppForm({
    defaultValues: {
      firstName: "John",
      lastName: "Doe",
    },
  });

  return <ChildForm form={form} title="Testing" />;
}
```

> **Note**: Use `function Render() {}` syntax (not arrow functions) to avoid ESLint hook warnings.

### Reusable Field Groups

Create reusable groups of related fields:

```tsx
const FieldGroupPasswordFields = withFieldGroup({
  defaultValues: {
    password: '',
    confirm_password: '',
  },
  props: {
    title: 'Password',
  },
  render: function Render({ group, title }) {
    // Access reactive values
    const password = useStore(group.store, (state) => state.values.password)

    return (
      <div>
        <h2>{title}</h2>
        <group.AppField name="password">
          {(field) => <field.TextField label="Password" />}
        </group.AppField>
        <group.AppField
          name="confirm_password"
          validators={{
            onChangeListenTo: ['password'],
            onChange: ({ value, fieldApi }) => {
              if (value !== group.getFieldValue('password')) {
                return 'Passwords do not match'
              }
              return undefined
            },
          }}
        >
          {(field) => (
            <div>
              <field.TextField label="Confirm Password" />
              <field.ErrorInfo />
            </div>
          )}
        </group.AppField>
      </div>
    )
  },
})

// Use in form
<FieldGroupPasswordFields
  form={form}
  fields="account_data" // Or map: { password: 'pwd', confirm_password: 'confirmPwd' }
  title="Passwords"
/>
```

**Field Mapping:**

```tsx
// Nested fields
fields="account_data"

// Top-level fields (custom mapping)
fields={{
  password: 'password',
  confirm_password: 'confirm_password',
}}

// Array fields
fields={`linked_accounts[${i}]`}
```

### Tree-Shaking with React.lazy

```tsx
// src/hooks/form-context.ts
import { createFormHookContexts } from '@tanstack/react-form'

export const { fieldContext, useFieldContext, formContext, useFormContext } =
  createFormHookContexts()

// src/components/text-field.tsx
import { useFieldContext } from '../hooks/form-context'

export default function TextField({ label }: { label: string }) {
  const field = useFieldContext<string>()
  return (/* ... */)
}

// src/hooks/form.ts
import { lazy } from 'react'
import { createFormHook } from '@tanstack/react-form'

const TextField = lazy(() => import('../components/text-field'))

const { useAppForm } = createFormHook({
  fieldContext,
  formContext,
  fieldComponents: {
    TextField,
  },
})

// src/App.tsx
import { Suspense } from 'react'

export default function App() {
  return (
    <Suspense fallback={<p>Loading...</p>}>
      <PeoplePage />
    </Suspense>
  )
}
```

---

## Best Practices

### 1. Validation Strategy

- **Use `onBlur` for most fields** - Better UX than `onChange` spam
- **Use `onChange` for real-time requirements** - Username availability, password strength
- **Use `onSubmit` for expensive validations** - API calls, complex cross-field checks
- **Combine validators** - Different rules at different times

```tsx
validators={{
  onBlur: ({ value }) => value.length < 3 ? 'Too short' : undefined,
  onChangeAsync: async ({ value }) => {
    const exists = await checkUsername(value)
    return exists ? 'Taken' : undefined
  },
  onChangeAsyncDebounceMs: 500,
}}
```

### 2. Error Display

- **Show errors only after interaction**: Use `field.state.meta.isTouched`
- **Use `errorMap` for specific error sources**: Different styling for different validators
- **Translate error keys**: Store validation keys, translate in UI

```tsx
{
  !field.state.meta.isValid && field.state.meta.isTouched && (
    <em role="alert">{translateError(field.state.meta.errors[0])}</em>
  );
}
```

### 3. Performance

- **Always use selectors with `useStore`**: Avoid unnecessary re-renders
- **Use `form.Subscribe` for UI**: More performant than `useStore` for rendering
- **Debounce async validation**: Prevent API spam
- **Lazy-load large forms**: Use `React.lazy` + `Suspense`

### 4. Type Safety

- **Define form types explicitly**: Better autocomplete and error detection
- **Use Standard Schemas**: Zod/Valibot for runtime + compile-time safety
- **Type error shapes**: Custom error types are fully typed

```tsx
type FormValues = {
  firstName: string;
  lastName: string;
  age: number;
};

const form = useForm<FormValues>({
  defaultValues: {
    firstName: "",
    lastName: "",
    age: 0,
  },
});
```

### 5. Form Submission

- **Prevent default**: Always `e.preventDefault()` and `e.stopPropagation()`
- **Check `canSubmit`**: Disable button when invalid
- **Show loading state**: Use `isSubmitting` for button state
- **Handle errors gracefully**: Use `onSubmitInvalid` for focus management

```tsx
<form
  onSubmit={(e) => {
    e.preventDefault();
    e.stopPropagation();
    form.handleSubmit();
  }}
>
  <form.Subscribe
    selector={(state) => [state.canSubmit, state.isSubmitting]}
    children={([canSubmit, isSubmitting]) => (
      <button type="submit" disabled={!canSubmit}>
        {isSubmitting ? "Submitting..." : "Submit"}
      </button>
    )}
  />
</form>
```

---

## Common Patterns

### Password Confirmation

```tsx
<form.Field
  name="confirm_password"
  validators={{
    onChangeListenTo: ["password"],
    onChange: ({ value, fieldApi }) => {
      if (value !== fieldApi.form.getFieldValue("password")) {
        return "Passwords must match";
      }
      return undefined;
    },
  }}
/>
```

### Autosave

```tsx
const form = useForm({
  listeners: {
    onChange: ({ formApi }) => {
      if (formApi.state.isValid) {
        formApi.handleSubmit();
      }
    },
    onChangeDebounceMs: 1000, // Debounce autosave
  },
});
```

### Conditional Fields

```tsx
<form.Subscribe
  selector={(state) => state.values.hasAddress}
  children={(hasAddress) => (
    hasAddress && (
      <form.Field name="address">
        {(field) => (/* address input */)}
      </form.Field>
    )
  )}
/>
```

### Multi-Step Forms

```tsx
const [step, setStep] = useState(0);

const form = useForm({
  onSubmit: ({ value }) => {
    if (step < 2) {
      setStep(step + 1);
    } else {
      submitForm(value);
    }
  },
});

return (
  <>
    {step === 0 && <Step1Fields />}
    {step === 1 && <Step2Fields />}
    {step === 2 && <Step3Fields />}
  </>
);
```

---

## API Quick Reference

### useForm Options

```tsx
{
  defaultValues: {},
  validators: {
    onChange: (state) => {},
    onBlur: (state) => {},
    onSubmit: (state) => {},
    onDynamic: (state) => {}, // Requires validationLogic
  },
  onSubmit: async ({ value }) => {},
  onSubmitInvalid: ({ formApi }) => {},
  listeners: {
    onChange: ({ formApi, fieldApi }) => {},
    onBlur: ({ formApi, fieldApi }) => {},
    onMount: ({ formApi }) => {},
    onSubmit: ({ formApi }) => {},
  },
  validationLogic: revalidateLogic(), // For onDynamic
}
```

### Field Options

```tsx
<form.Field
  name="fieldName"
  mode="value" // or "array"
  validators={{
    onChange: ({ value, fieldApi }) => {},
    onBlur: ({ value, fieldApi }) => {},
    onChangeAsync: async ({ value, fieldApi }) => {},
    onBlurAsync: async ({ value, fieldApi }) => {},
    onChangeListenTo: ["otherField"],
    onBlurListenTo: ["otherField"],
  }}
  asyncDebounceMs={500}
  onChangeAsyncDebounceMs={1000}
  disableErrorFlat={false}
  listeners={{
    onChange: ({ value }) => {},
    onBlur: ({ value }) => {},
  }}
/>
```

### Field State

```tsx
field.state.value; // Current value
field.state.meta.errors; // Array of all errors
field.state.meta.errorMap; // Errors by source (onChange, onBlur, etc.)
field.state.meta.isValid; // Is field valid
field.state.meta.isTouched; // Has field been touched
field.state.meta.isDirty; // Has value changed from default
```

### Form State

```tsx
form.state.values; // All form values
form.state.errors; // All form errors
form.state.errorMap; // Form-level errors by source
form.state.isValid; // Is entire form valid
form.state.isSubmitting; // Is form currently submitting
form.state.canSubmit; // Can form be submitted (valid + touched)
form.state.isPristine; // No fields have been touched
form.state.isDirty; // Any field has changed
```

---

## Troubleshooting

### Issue: Fields not revalidating when another field changes

**Solution**: Use `onChangeListenTo` or `onBlurListenTo`

```tsx
validators={{
  onChangeListenTo: ['dependentField'],
  onChange: ({ value, fieldApi }) => {
    // Will rerun when dependentField changes
  },
}}
```

### Issue: Form causing too many re-renders

**Solution**: Use `form.Subscribe` instead of `useStore`, or add selectors

```tsx
// ❌ Bad - triggers re-render on any state change
const formState = useStore(form.store)

// ✅ Good - only re-renders when firstName changes
const firstName = useStore(form.store, (state) => state.values.firstName)

// ✅ Better - no component re-render
<form.Subscribe
  selector={(state) => state.values.firstName}
  children={(firstName) => <div>{firstName}</div>}
/>
```

### Issue: Async validation firing too often

**Solution**: Add debouncing

```tsx
<form.Field
  name="username"
  asyncDebounceMs={500} // Global debounce
  validators={{
    onChangeAsyncDebounceMs: 1000, // Specific debounce
    onChangeAsync: async ({ value }) => {
      // Only fires after 1000ms of no changes
    },
  }}
/>
```

### Issue: Form errors not clearing on field change

**Solution**: This is expected with `onBlur` validation. Use `onChange` or clear manually:

```tsx
listeners={{
  onChange: ({ fieldApi }) => {
    // Clear errors on change if needed
    fieldApi.setMeta({ errors: [] })
  },
}}
```

### Issue: TypeScript errors with error types

**Solution**: Use `errorMap` for type-safe access or cast errors array

```tsx
// Type-safe
const onChangeError = field.state.meta.errorMap.onChange;

// With casting
const error = field.state.meta.errors[0] as string;
```

---

## Integration Examples

### With TanStack Query

```tsx
import { useForm } from '@tanstack/react-form'
import { useQuery, useMutation } from '@tanstack/react-query'

function App() {
  // Fetch initial data
  const { data, isLoading } = useQuery({
    queryKey: ['user'],
    queryFn: fetchUser,
  })

  // Submit mutation
  const mutation = useMutation({
    mutationFn: updateUser,
  })

  const form = useForm({
    defaultValues: {
      firstName: data?.firstName ?? '',
      lastName: data?.lastName ?? '',
    },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync(value)
    },
  })

  if (isLoading) return <p>Loading...</p>

  return (/* form */)
}
```

### With Zod

```tsx
import { useForm } from "@tanstack/react-form";
import { zodValidator } from "@tanstack/zod-form-adapter";
import { z } from "zod";

const schema = z.object({
  email: z.string().email("Invalid email"),
  password: z.string().min(8, "Must be 8+ characters"),
});

const form = useForm({
  validatorAdapter: zodValidator(), // Enable Zod support
  validators: {
    onChange: schema,
  },
});
```

---

## Resources

- [Official Documentation](https://tanstack.com/form/latest)
- [GitHub Repository](https://github.com/tanstack/form)
- [Examples](https://tanstack.com/form/latest/docs/framework/react/examples)
- [Discord Community](https://tlinz.com/discord)
