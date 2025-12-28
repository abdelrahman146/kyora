import { create } from "zustand";
import { persist } from "zustand/middleware";
import { metadataApi } from "@/api/metadata";
import type { CountryMetadata } from "@/api/types/metadata";

type LoadStatus = "idle" | "loading" | "loaded" | "error";

interface MetadataState {
  countries: CountryMetadata[];
  status: LoadStatus;
  error: string | null;
  loadedAt: number | null;

  loadCountries: (opts?: { force?: boolean }) => Promise<void>;
  clear: () => void;
}

const ONE_DAY_MS = 24 * 60 * 60 * 1000;

export const useMetadataStore = create<MetadataState>()(
  persist(
    (set, get) => ({
      countries: [],
      status: "idle",
      error: null,
      loadedAt: null,

      loadCountries: async (opts) => {
        const { force = false } = opts ?? {};
        const { countries, status, loadedAt } = get();

        if (!force) {
          if (status === "loading") return;
          if (
            countries.length > 0 &&
            loadedAt &&
            Date.now() - loadedAt < ONE_DAY_MS
          ) {
            return;
          }
        }

        set({ status: "loading", error: null });
        try {
          const { countries: fetched } = await metadataApi.listCountries();
          set({
            countries: fetched,
            status: "loaded",
            loadedAt: Date.now(),
            error: null,
          });
        } catch (err) {
          const message =
            err instanceof Error ? err.message : "Failed to load metadata";
          set({ status: "error", error: message });
        }
      },

      clear: () =>
        set({ countries: [], status: "idle", error: null, loadedAt: null }),
    }),
    {
      name: "kyora-metadata-store",
      partialize: (state) => ({
        countries: state.countries,
        loadedAt: state.loadedAt,
      }),
    }
  )
);
