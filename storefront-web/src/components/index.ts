// Atomic Design Component Structure
// This file provides convenient exports for all components organized by atomic design principles

// Atoms - Basic building blocks
export * from "./atoms";

// Molecules - Simple combinations of atoms
export * from "./molecules";

// Organisms - Complex combinations forming distinct sections
export * from "./organisms";

// Named exports for backward compatibility
export { ErrorBoundary } from "./atoms";
export { EmptyState } from "./atoms";
export { ImageTile } from "./atoms";
export { Button } from "./atoms";
export { IconButton } from "./atoms";
export { Badge } from "./atoms";

export { ProductCard, ProductCardSkeleton } from "./molecules";
export { StickyCartBar } from "./molecules";
export { FloatingWhatsAppButton } from "./molecules";
export { LanguageSwitcher } from "./molecules";
export { CartItem } from "./molecules";
export { FormInput } from "./molecules";
export { FormSelect } from "./molecules";
export { OrderSummary } from "./molecules";

export { StorefrontHeader } from "./organisms";
export { ProductListNew } from "./organisms";
export { BottomSheet } from "./organisms";
export { CartDrawer } from "./organisms";
export { WhatsAppButton } from "./organisms";
