import { create } from "zustand";
import { persist } from "zustand/middleware";

export interface CartLineItem {
  variantId: string;
  productId: string;
  productName: string;
  variantName: string;
  unitPrice: string; // decimal string
  currency: string;
  photoUrl?: string;
  quantity: number;
  note?: string;
}

export interface CartState {
  itemsByVariantId: Record<string, CartLineItem>;
  isCartOpen: boolean;

  addItem: (item: Omit<CartLineItem, "quantity">, quantity?: number) => void;
  removeItem: (variantId: string, quantity?: number) => void;
  setQuantity: (variantId: string, quantity: number) => void;
  updateNote: (variantId: string, note: string) => void;
  clearCart: () => void;

  openCart: () => void;
  closeCart: () => void;
}

const STORAGE_KEY = "kyora_storefront_cart_v1";

export const useCartStore = create<CartState>()(
  persist(
    (set) => ({
      itemsByVariantId: {},
      isCartOpen: false,

      addItem: (item, quantity = 1) => {
        if (!item?.variantId) return;
        if (quantity <= 0) return;

        set((state) => {
          const existing = state.itemsByVariantId[item.variantId];
          const nextQty = (existing?.quantity || 0) + quantity;

          return {
            itemsByVariantId: {
              ...state.itemsByVariantId,
              [item.variantId]: {
                ...(existing || { ...item, quantity: 0 }),
                ...item,
                quantity: nextQty,
              },
            },
          };
        });
      },

      removeItem: (variantId, quantity = 1) => {
        if (!variantId) return;
        if (quantity <= 0) return;

        set((state) => {
          const existing = state.itemsByVariantId[variantId];
          if (!existing) return state;

          const nextQty = existing.quantity - quantity;
          if (nextQty > 0) {
            return {
              itemsByVariantId: {
                ...state.itemsByVariantId,
                [variantId]: { ...existing, quantity: nextQty },
              },
            };
          }

          const next = { ...state.itemsByVariantId };
          delete next[variantId];
          return { itemsByVariantId: next };
        });
      },

      setQuantity: (variantId, quantity) => {
        if (!variantId) return;
        if (!Number.isFinite(quantity)) return;

        set((state) => {
          const existing = state.itemsByVariantId[variantId];
          if (!existing) return state;

          if (quantity <= 0) {
            const next = { ...state.itemsByVariantId };
            delete next[variantId];
            return { itemsByVariantId: next };
          }

          return {
            itemsByVariantId: {
              ...state.itemsByVariantId,
              [variantId]: { ...existing, quantity },
            },
          };
        });
      },

      updateNote: (variantId, note) => {
        if (!variantId) return;
        const trimmed = (note || "").slice(0, 500);
        set((state) => {
          const existing = state.itemsByVariantId[variantId];
          if (!existing) return state;
          return {
            itemsByVariantId: {
              ...state.itemsByVariantId,
              [variantId]: { ...existing, note: trimmed },
            },
          };
        });
      },

      clearCart: () => set({ itemsByVariantId: {}, isCartOpen: false }),

      openCart: () => set({ isCartOpen: true }),
      closeCart: () => set({ isCartOpen: false }),
    }),
    {
      name: STORAGE_KEY,
      version: 1,
      partialize: (state) => ({
        itemsByVariantId: state.itemsByVariantId,
      }),
    }
  )
);

export function cartItemsArray(
  state: Pick<CartState, "itemsByVariantId">
): CartLineItem[] {
  return Object.values(state.itemsByVariantId).sort((a, b) =>
    a.productName.localeCompare(b.productName)
  );
}

export function cartTotalQuantity(
  state: Pick<CartState, "itemsByVariantId">
): number {
  return Object.values(state.itemsByVariantId).reduce(
    (acc, x) => acc + (x.quantity || 0),
    0
  );
}
