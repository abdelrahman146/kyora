# BottomSheet Quick Reference

## Import

```tsx
import { BottomSheet } from "@/components/molecules/BottomSheet";
```

## Minimal Example

```tsx
<BottomSheet
  isOpen={isOpen}
  onClose={() => { setIsOpen(false); }}
  title="My Drawer"
>
  <p>Content here</p>
</BottomSheet>
```

## Common Patterns

### Filter Drawer (with footer actions)
```tsx
<BottomSheet
  isOpen={isOpen}
  onClose={() => { setIsOpen(false); }}
  title="Filters"
  footer={
    <div className="flex gap-2">
      <button onClick={handleReset} className="btn btn-ghost flex-1">Reset</button>
      <button onClick={handleApply} className="btn btn-primary flex-1">Apply</button>
    </div>
  }
>
  {/* Filter content */}
</BottomSheet>
```

### Navigation Menu (left side)
```tsx
<BottomSheet
  isOpen={isMenuOpen}
  onClose={() => { setIsMenuOpen(false); }}
  title="Menu"
  side="start"
  size="sm"
>
  <nav>...</nav>
</BottomSheet>
```

### Shopping Cart (right side)
```tsx
<BottomSheet
  isOpen={isCartOpen}
  onClose={() => { setIsCartOpen(false); }}
  title="Cart"
  side="end"
  size="md"
  footer={<button className="btn btn-primary w-full">Checkout</button>}
>
  {/* Cart items */}
</BottomSheet>
```

### Full-Width Panel
```tsx
<BottomSheet
  isOpen={isOpen}
  onClose={() => { setIsOpen(false); }}
  title="Settings"
  size="full"
  contentClassName="max-w-4xl mx-auto"
>
  {/* Settings sections */}
</BottomSheet>
```

### Custom Header
```tsx
<BottomSheet
  isOpen={isOpen}
  onClose={() => { setIsOpen(false); }}
  header={
    <div className="flex items-center gap-3">
      <Avatar />
      <UserInfo />
    </div>
  }
>
  {/* Profile content */}
</BottomSheet>
```

### Confirmation Dialog
```tsx
<BottomSheet
  isOpen={isOpen}
  onClose={() => { setIsOpen(false); }}
  title="Confirm"
  size="sm"
  closeOnOverlayClick={false}
  footer={
    <div className="flex gap-2">
      <button onClick={onCancel} className="btn btn-ghost flex-1">Cancel</button>
      <button onClick={onConfirm} className="btn btn-error flex-1">Delete</button>
    </div>
  }
>
  <p>Are you sure?</p>
</BottomSheet>
```

## Props Quick Reference

| Prop | Type | Default | Use When |
|------|------|---------|----------|
| `isOpen` | `boolean` | - | Always (required) |
| `onClose` | `() => void` | - | Always (required) |
| `children` | `ReactNode` | - | Always (required) |
| `title` | `string` | - | Simple header with title |
| `header` | `ReactNode` | - | Custom header content |
| `footer` | `ReactNode` | - | Action buttons needed |
| `side` | `'start' \| 'end'` | `'end'` | Change desktop position |
| `size` | `'sm' \| 'md' \| 'lg' \| 'xl' \| 'full'` | `'md'` | Adjust width |
| `closeOnOverlayClick` | `boolean` | `true` | Prevent accidental close |
| `closeOnEscape` | `boolean` | `true` | Prevent accidental close |
| `showHeader` | `boolean` | `true` | Hide header completely |
| `showCloseButton` | `boolean` | `true` | Hide close button |
| `className` | `string` | - | Style drawer container |
| `contentClassName` | `string` | - | Style content area |
| `footerClassName` | `string` | - | Style footer area |

## Sizes

- `sm`: 384px - Navigation, small forms
- `md`: 448px - Filters, medium forms (default)
- `lg`: 512px - Details panels
- `xl`: 576px - Complex forms
- `full`: 100% - Settings, large content

## Behavior

- **Mobile**: Always bottom sheet, full width, max 85% height
- **Desktop**: Side drawer, configurable width, full height
- **RTL**: Automatic - `side="end"` goes to logical end (right in LTR, left in RTL)
- **Animations**: 300ms slide transition
- **Body Scroll**: Locked when open, restored when closed
- **Focus**: Trapped in drawer, restored to trigger on close
- **Keyboard**: Escape to close (configurable)

## RTL Notes

Use logical positions:
- `side="start"` = Left in LTR, Right in RTL
- `side="end"` = Right in LTR, Left in RTL

Don't use `left` or `right` in custom styles - use `start` and `end` classes.

## Accessibility

- Uses `role="dialog"` and `aria-modal="true"`
- Title auto-generates `aria-labelledby`
- Can override with `ariaLabel` prop
- Close button has proper `aria-label`
- Escape key support
- Focus management

## Tips

1. **Always use `useCallback`** for `onClose` to prevent re-renders
2. **Stable state** - manage `isOpen` in parent component
3. **Footer actions** - primary button on right/end
4. **Loading states** - show within drawer content
5. **Errors** - display inline, don't close drawer
6. **Long content** - drawer automatically scrolls
7. **Stacking** - avoid multiple drawers, use nested modals instead

## Complete Documentation

See [BottomSheet.md](./BottomSheet.md) for:
- Detailed API reference
- Multiple examples
- Best practices
- Troubleshooting
- Browser compatibility
- Performance tips
