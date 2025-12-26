# Kyora Storefront - Design System Implementation

This document outlines the complete revamp of the `storefront-web` application according to the **Kyora Design System (KDS)**.

## ✨ What's New

### 🎨 Design System Integration

- **Color Palette**: Implemented the complete KDS color system with primary (teal), secondary (gold), and semantic colors
- **Typography**: IBM Plex Sans Arabic as the primary font with proper fallbacks
- **Spacing**: 4px baseline grid system with logical spacing tokens
- **Border Radius**: Consistent radius values (sm: 4px, md: 8px, lg: 12px, xl: 16px)
- **Shadows**: Elevation system with sm and float shadow variants

### 🌍 RTL-First Design

- **Direction Toggle**: Automatic RTL/LTR switching based on language
- **Logical Properties**: All spacing uses `inline-start/end` instead of left/right
- **Icon Mirroring**: Directional icons (arrows, chevrons) are automatically mirrored in RTL
- **Arabic Typography**: Optimized for Arabic text rendering with proper line heights

### 📱 Mobile-First Components

#### 1. **StorefrontHeader** (`components/StorefrontHeader.tsx`)
- Sticky header with scroll shadow effect
- Left: Hamburger menu (RTL start)
- Center: Logo and brand name
- Right: Search + Cart with badge (RTL end)
- Minimum 44px touch targets

#### 2. **ProductCard** (`components/ProductCard.tsx`)
- 2-column grid layout (strict mobile)
- 1:1 aspect ratio images with object-cover
- Floating "+" button on image corner
- Quantity controls appear when item is in cart
- Title: max 2 lines with line-clamp
- Price: bold, primary-700 color
- Stock indicators (out of stock, low stock)

#### 3. **ProductCardSkeleton** (`components/ProductCard.tsx`)
- Pulsing gray blocks matching card shape
- Used during loading states
- Prevents layout shift

#### 4. **ProductList** (`components/ProductListNew.tsx`)
- Category tabs with horizontal scroll
- Active state with primary color
- Empty states with friendly messaging
- Grid always 2 columns on mobile

#### 5. **StickyCartBar** (`components/StickyCartBar.tsx`)
- Appears only when cart has items
- Floats at bottom above safe area
- Shows item count and total
- High contrast (primary-900 bg, white text)
- Proper touch target (56px height)

#### 6. **BottomSheet** (`components/BottomSheet.tsx`)
- Slide-up panel from bottom
- Gray pill handle for discoverability
- Scrollable content area
- Optional sticky footer for actions
- Backdrop blur effect
- Keyboard accessible (Escape to close)

#### 7. **FloatingWhatsAppButton** (`components/FloatingWhatsAppButton.tsx`)
- Fixed bottom-right (bottom-left in RTL)
- WhatsApp brand green (#25D366)
- Pre-filled message with store name
- High visibility with shadow
- Smooth hover/active animations

#### 8. **EmptyState** (`components/EmptyState.tsx`)
- Friendly empty state illustrations
- Used for empty categories, cart, etc.
- Optional action button
- Centered layout with proper spacing

### 🎯 Key Features

#### Accessibility (A11y)
- **ARIA Labels**: All icon buttons have proper aria-labels
- **Keyboard Navigation**: Full keyboard support (Tab, Enter, Escape)
- **Focus Management**: Visible focus rings (never hidden)
- **Screen Readers**: Semantic HTML and ARIA attributes
- **Touch Targets**: Minimum 44x44px for all interactive elements
- **Color Contrast**: WCAG AA compliant (primary-600 on white)

#### Performance
- **Lazy Loading**: Images load on demand
- **Skeleton States**: No white screens during loading
- **Optimized Re-renders**: Memoized expensive computations
- **Small Bundle**: Only necessary dependencies loaded

#### PWA Features
- **Offline-First**: Service worker ready
- **Add to Home Screen**: Proper manifest.json
- **iOS Safe Areas**: Safe-top/safe-bottom classes
- **Tap Highlight**: Disabled for native feel
- **Smooth Scrolling**: Hardware-accelerated animations

## 🛠️ Technical Stack

- **React 19**: Latest React with concurrent features
- **TypeScript**: Fully typed, no `any` types
- **Tailwind CSS v4**: Latest version with design tokens
- **DaisyUI v5**: Component library following KDS
- **Zustand**: Lightweight state management
- **TanStack Query**: Server state management
- **Lucide React**: Icon library (700+ icons)
- **i18next**: Internationalization (Arabic/English)

## 📂 File Structure

```
src/
├── components/
│   ├── BottomSheet.tsx          # Modal drawer from bottom
│   ├── CartDrawer.tsx            # Cart with checkout form
│   ├── EmptyState.tsx            # Friendly empty states
│   ├── FloatingWhatsAppButton.tsx # WhatsApp FAB
│   ├── ImageTile.tsx             # Optimized image component
│   ├── LanguageSwitcher.tsx      # AR/EN toggle
│   ├── ProductCard.tsx           # Product card + skeleton
│   ├── ProductListNew.tsx        # Grid with categories
│   ├── StickyCartBar.tsx         # Bottom cart summary
│   └── StorefrontHeader.tsx      # Sticky header
├── pages/
│   └── StorefrontPage.tsx        # Main storefront page
├── index.css                     # KDS design tokens
└── App.tsx                       # RTL direction handler
```

## 🎨 Design Tokens

### Colors
```css
Primary (Brand Teal):
  --primary-50: #F0FDFA
  --primary-600: #0D9488 (Main)
  --primary-900: #134E4A

Secondary (Gold/Sand):
  --secondary-500: #EAB308

Neutral:
  --neutral-0: #FFFFFF
  --neutral-50: #F8FAFC (Background)
  --neutral-900: #0F172A (Text)

Semantic:
  --success: #10B981
  --error: #EF4444
  --warning: #F59E0B
```

### Typography
```
Font: IBM Plex Sans Arabic
Display: 32px / Bold
H1: 24px / Bold
H2: 20px / SemiBold
H3: 18px / Medium
Body-L: 16px / Regular
Body-M: 14px / Regular
Caption: 12px / Medium
```

### Spacing (4px baseline)
```
--gap-xs: 4px
--gap-s: 8px
--gap-m: 16px
--gap-l: 24px
--gap-xl: 32px
```

## 🚀 Usage

### Development
```bash
npm run dev
```

### Build
```bash
npm run build
```

### Preview
```bash
npm run preview
```

## ♿ Accessibility Checklist

- [x] Semantic HTML5 elements
- [x] ARIA labels for icon buttons
- [x] Keyboard navigation (Tab, Enter, Escape)
- [x] Focus visible (custom focus-ring class)
- [x] Color contrast WCAG AA
- [x] Touch targets >= 44px
- [x] Alt text for images
- [x] Screen reader friendly
- [x] Skip to content link (future)
- [x] No seizure-inducing animations

## 🌐 Internationalization

### Supported Languages
- **Arabic (ar)**: Primary language, RTL layout
- **English (en)**: Secondary language, LTR layout

### Adding New Translations
Edit `src/i18n/translations.ts`:
```typescript
export const translations = {
  en: {
    common: {
      yourKey: "Your English Text"
    }
  },
  ar: {
    common: {
      yourKey: "النص بالعربية"
    }
  }
}
```

## 📱 Responsive Breakpoints

- **Mobile**: < 640px (base, 2 columns)
- **Tablet**: 640px - 1024px (optional, not implemented yet)
- **Desktop**: >= 1024px (max-w-5xl container)

## 🎭 Animation Guidelines

- **Press Effects**: `active:scale-95` on buttons
- **Hover States**: Smooth color transitions (200ms)
- **Slide-Up**: Bottom sheet animation (300ms cubic-bezier)
- **Skeleton Pulse**: 2s ease-in-out infinite
- **Fade-In**: Page transitions (optional)

## 🔧 Customization

### Changing Brand Colors
Edit `src/index.css`:
```css
@plugin "daisyui/theme" {
  name: "kyora";
  --color-primary: oklch(...); /* Your color */
}
```

### Adjusting Grid Columns
Edit `src/components/ProductListNew.tsx`:
```tsx
// Change from 2 to 3 columns
<div className="grid grid-cols-3 gap-3">
```

## 🐛 Known Issues & Future Enhancements

### Current Limitations
- Menu drawer is placeholder (needs category navigation)
- Search functionality not implemented yet
- Product detail page (single product view) pending
- Order tracking page pending
- PWA manifest needs refinement

### Planned Features
- [ ] Product quick view modal
- [ ] Advanced filtering (price, stock)
- [ ] Product search with autocomplete
- [ ] Wishlist/favorites
- [ ] Share product via social media
- [ ] Multi-currency support
- [ ] Customer reviews/ratings
- [ ] Order tracking page
- [ ] Push notifications

## 📚 Resources

- [Kyora Design System](../.github/instructions/branding.instructions.md)
- [DaisyUI v5 Docs](../.github/instructions/daisyui.instructions.md)
- [Tailwind CSS v4](https://tailwindcss.com/)
- [Lucide Icons](https://lucide.dev/)
- [IBM Plex Sans Arabic](https://fonts.google.com/specimen/IBM+Plex+Sans+Arabic)

## 🤝 Contributing

When adding new components:
1. Follow KDS specifications strictly
2. Use design tokens (no magic numbers)
3. Support RTL from day one
4. Add TypeScript types
5. Include accessibility attributes
6. Test on mobile devices
7. Document in this README

## 📄 License

See main repository LICENSE file.

---

**Built with ❤️ for Middle East entrepreneurs**
