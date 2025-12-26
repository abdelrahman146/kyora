import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { Business } from "../api/types/business";

/**
 * Business Store
 *
 * Manages global business state including:
 * - List of available businesses
 * - Currently selected business and selectedBusinessId
 * - Sidebar UI state (collapsed on desktop, open/closed on mobile)
 */

interface BusinessState {
  // Data
  businesses: Business[];
  selectedBusiness: Business | null;
  selectedBusinessId: string | null;

  // UI State
  isSidebarCollapsed: boolean; // Desktop: collapsed (icon-only) state
  isSidebarOpen: boolean; // Mobile: drawer open state

  // Actions
  setBusinesses: (businesses: Business[]) => void;
  setSelectedBusiness: (business: Business | null) => void;
  setSelectedBusinessId: (id: string | null) => void;
  toggleSidebar: () => void;
  setSidebarCollapsed: (collapsed: boolean) => void;
  openSidebar: () => void;
  closeSidebar: () => void;
}

export const useBusinessStore = create<BusinessState>()(
  persist(
    (set) => ({
      // Initial State
      businesses: [],
      selectedBusiness: null,
      selectedBusinessId: null,
      isSidebarCollapsed: false,
      isSidebarOpen: false,

      // Actions
      setBusinesses: (businesses) => set({ businesses }),

      setSelectedBusiness: (business) =>
        set({
          selectedBusiness: business,
          selectedBusinessId: business?.id ?? null,
        }),

      setSelectedBusinessId: (id) => set({ selectedBusinessId: id }),

      // Desktop: toggle collapsed state (full width ↔ icon-only)
      toggleSidebar: () =>
        set((state) => ({ isSidebarCollapsed: !state.isSidebarCollapsed })),

      setSidebarCollapsed: (collapsed) =>
        set({ isSidebarCollapsed: collapsed }),

      // Mobile: open/close drawer
      openSidebar: () => {
        set({ isSidebarOpen: true });
      },
      closeSidebar: () => {
        set({ isSidebarOpen: false });
      },
    }),
    {
      name: "kyora-business-store",
      // Only persist selected business ID and desktop sidebar state
      partialize: (state) => ({
        selectedBusiness: state.selectedBusiness,
        selectedBusinessId: state.selectedBusinessId,
        isSidebarCollapsed: state.isSidebarCollapsed,
      }),
    }
  )
);
