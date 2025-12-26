import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { Business } from "../api/types/business";

/**
 * Business Store
 *
 * Manages global business state including:
 * - List of available businesses
 * - Currently selected business
 * - Sidebar collapsed state (UI state)
 */

interface BusinessState {
  // Data
  businesses: Business[];
  selectedBusiness: Business | null;

  // UI State
  isSidebarCollapsed: boolean;

  // Actions
  setBusinesses: (businesses: Business[]) => void;
  setSelectedBusiness: (business: Business | null) => void;
  toggleSidebar: () => void;
  setSidebarCollapsed: (collapsed: boolean) => void;
}

export const useBusinessStore = create<BusinessState>()(
  persist(
    (set) => ({
      // Initial State
      businesses: [],
      selectedBusiness: null,
      isSidebarCollapsed: false,

      // Actions
      setBusinesses: (businesses) => set({ businesses }),

      setSelectedBusiness: (business) => set({ selectedBusiness: business }),

      toggleSidebar: () =>
        set((state) => ({ isSidebarCollapsed: !state.isSidebarCollapsed })),

      setSidebarCollapsed: (collapsed) =>
        set({ isSidebarCollapsed: collapsed }),
    }),
    {
      name: "kyora-business-store",
      // Only persist selected business ID and sidebar state
      partialize: (state) => ({
        selectedBusiness: state.selectedBusiness,
        isSidebarCollapsed: state.isSidebarCollapsed,
      }),
    }
  )
);
