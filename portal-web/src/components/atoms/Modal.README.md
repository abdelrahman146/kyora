# Modal Component

Production-grade, reusable modal component with mobile-first design.

## Features

- 🎨 **Mobile-First Design**: Bottom sheet on mobile, centered modal on desktop
- 📱 **Fully Responsive**: Adapts to all screen sizes with smooth transitions
- ♿ **Accessible**: Focus trap, keyboard navigation (Escape to close), ARIA attributes
- 🎭 **Portal-based**: Renders at the end of document body for proper stacking
- 🔒 **Scroll Lock**: Prevents body scroll when modal is open
- 🎯 **Flexible Sizing**: Multiple size options (sm, md, lg, xl, full)
- 🌐 **RTL Support**: Works seamlessly with RTL languages using logical properties
- 🎬 **Smooth Animations**: CSS-based transitions for performance
- 🎨 **DaisyUI Themed**: Inherits theme colors and styling

## Basic Usage

```tsx
import { useState } from "react";
import { Modal } from "@/components/atoms";

function MyComponent() {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <>
      <button onClick={() => setIsOpen(true)} className="btn btn-primary">
        Open Modal
      </button>

      <Modal
        isOpen={isOpen}
        onClose={() => setIsOpen(false)}
        title="My Modal Title"
      >
        <p>Modal content goes here...</p>
      </Modal>
    </>
  );
}
```

## Props

| Prop                    | Type                                  | Default | Description                                       |
| ----------------------- | ------------------------------------- | ------- | ------------------------------------------------- |
| `isOpen`                | `boolean`                             | -       | **Required**. Whether the modal is open           |
| `onClose`               | `() => void`                          | -       | **Required**. Callback when modal should close    |
| `title`                 | `ReactNode`                           | -       | Modal title shown in header                       |
| `children`              | `ReactNode`                           | -       | **Required**. Modal body content                  |
| `footer`                | `ReactNode`                           | -       | Footer content (typically action buttons)         |
| `size`                  | `"sm" \| "md" \| "lg" \| "xl" \| "full"` | `"md"` | Modal width                                       |
| `closeOnBackdropClick`  | `boolean`                             | `true`  | Allow closing by clicking backdrop                |
| `closeOnEscape`         | `boolean`                             | `true`  | Allow closing with Escape key                     |
| `showCloseButton`       | `boolean`                             | `true`  | Show X button in top-right                        |
| `className`             | `string`                              | -       | Additional CSS classes for container              |
| `contentClassName`      | `string`                              | -       | Additional CSS classes for modal box              |
| `scrollable`            | `boolean`                             | `true`  | Enable scrolling in modal content                 |
| `zIndex`                | `number`                              | `50`    | Custom z-index value                              |

## Examples

### Confirmation Dialog

```tsx
<Modal
  isOpen={isDeleteModalOpen}
  onClose={() => setIsDeleteModalOpen(false)}
  title="Delete Item"
  size="sm"
  footer={
    <>
      <button onClick={() => setIsDeleteModalOpen(false)} className="btn btn-ghost">
        Cancel
      </button>
      <button onClick={handleDelete} className="btn btn-error">
        Delete
      </button>
    </>
  }
>
  <p>Are you sure you want to delete this item? This action cannot be undone.</p>
</Modal>
```

### Form Modal

```tsx
<Modal
  isOpen={isFormOpen}
  onClose={() => setIsFormOpen(false)}
  title="Add New Item"
  size="lg"
  footer={
    <>
      <button onClick={() => setIsFormOpen(false)} className="btn btn-ghost">
        Cancel
      </button>
      <button onClick={handleSubmit} className="btn btn-primary">
        Save
      </button>
    </>
  }
>
  <form className="space-y-4">
    <div className="form-control">
      <label className="label">
        <span className="label-text">Name</span>
      </label>
      <input type="text" className="input input-bordered" />
    </div>
    <div className="form-control">
      <label className="label">
        <span className="label-text">Description</span>
      </label>
      <textarea className="textarea textarea-bordered" />
    </div>
  </form>
</Modal>
```

### Large Content Modal

```tsx
<Modal
  isOpen={isDetailsOpen}
  onClose={() => setIsDetailsOpen(false)}
  title="Item Details"
  size="xl"
  scrollable={true}
  footer={
    <button onClick={() => setIsDetailsOpen(false)} className="btn btn-primary">
      Close
    </button>
  }
>
  <div className="prose max-w-none">
    {/* Long content that will scroll */}
    <h2>Section 1</h2>
    <p>Long content...</p>
    <h2>Section 2</h2>
    <p>More content...</p>
  </div>
</Modal>
```

### Full-Screen Modal (Mobile)

```tsx
<Modal
  isOpen={isFullScreenOpen}
  onClose={() => setIsFullScreenOpen(false)}
  title="Full Details"
  size="full"
  closeOnBackdropClick={false}
>
  <div className="space-y-6">
    {/* Full content */}
  </div>
</Modal>
```

### Alert/Info Modal (No Footer)

```tsx
<Modal
  isOpen={isAlertOpen}
  onClose={() => setIsAlertOpen(false)}
  title="Important Notice"
  size="md"
>
  <div className="alert alert-info">
    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" className="stroke-current shrink-0 w-6 h-6">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
    </svg>
    <span>Please note this important information.</span>
  </div>
</Modal>
```

### Custom Styling

```tsx
<Modal
  isOpen={isCustomOpen}
  onClose={() => setCustomOpen(false)}
  title="Custom Styled Modal"
  className="bg-linear-to-br from-primary/10 to-secondary/10"
  contentClassName="border-2 border-primary"
>
  <p>Custom styled modal content</p>
</Modal>
```

## Accessibility

The Modal component follows WAI-ARIA best practices:

- **Role**: `dialog` with `aria-modal="true"`
- **Labeling**: `aria-labelledby` links to the title
- **Focus Management**: Automatically focuses first focusable element
- **Keyboard Navigation**:
  - `Escape` closes the modal (if `closeOnEscape={true}`)
  - Focus is trapped within the modal
- **Screen Readers**: Proper semantic HTML and ARIA attributes

## Mobile Behavior

On mobile devices (< md breakpoint):

- Modal appears as a **bottom sheet** sliding up from the bottom
- Rounded only on top corners for native app feel
- Height limited to 90vh to allow peek at content behind
- Smooth slide-in animation

On desktop (≥ md breakpoint):

- Modal is **centered** on screen
- Fully rounded corners
- Fade-in animation
- Max height of 85vh

## Best Practices

1. **Always provide a way to close**: Either enable `closeOnBackdropClick`, `closeOnEscape`, or `showCloseButton`
2. **Use appropriate size**: Choose the size that fits your content without excessive scrolling
3. **Clear actions**: Put primary actions in the footer, typically aligned to the end
4. **Keep content focused**: Modals should be for focused tasks, not full pages
5. **Avoid nested modals**: Don't open a modal from within a modal
6. **Test with keyboard**: Ensure all actions can be performed with keyboard only
7. **Test on mobile**: Verify the bottom sheet behavior works well

## Performance

- **Portal rendering**: Modal is rendered outside the React tree to avoid z-index issues
- **CSS animations**: Uses CSS transitions for better performance than JS animations
- **Scroll lock**: Efficiently prevents body scroll without layout shifts
- **Event cleanup**: All event listeners are properly cleaned up on unmount

## Browser Support

Works in all modern browsers that support:

- CSS `position: fixed`
- CSS transitions
- React Portals
- `document.body.style` manipulation
