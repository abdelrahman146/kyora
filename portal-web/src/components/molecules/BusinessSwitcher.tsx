import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { Check, ChevronDown, Building2, Loader2 } from "lucide-react";
import { Avatar } from "../atoms/Avatar";
import { useBusinessStore } from "../../stores/businessStore";
import { businessApi } from "../../api/business";
import { cn } from "@/lib/utils";

/**
 * BusinessSwitcher Component
 *
 * Dropdown in the header that allows users to switch between businesses.
 * Fetches businesses on mount and updates global selectedBusinessId state.
 *
 * Features:
 * - Lists all available businesses
 * - Shows current selected business
 * - Updates global state on selection
 * - Loading and error states
 * - RTL support
 *
 * @example
 * ```tsx
 * <BusinessSwitcher />
 * ```
 */
export function BusinessSwitcher() {
  const { t } = useTranslation();
  const { businesses, selectedBusiness, setBusinesses, setSelectedBusiness } =
    useBusinessStore();
  const [isLoading, setIsLoading] = useState(false);
  const [isOpen, setIsOpen] = useState(false);

  // Fetch businesses on mount if not already loaded
  useEffect(() => {
    const loadBusinesses = async () => {
      if (businesses.length > 0) return;

      setIsLoading(true);
      try {
        const fetchedBusinesses = await businessApi.listBusinesses();
        setBusinesses(fetchedBusinesses);

        // Auto-select first business if none selected
        if (!selectedBusiness && fetchedBusinesses.length > 0) {
          setSelectedBusiness(fetchedBusinesses[0]);
        }
      } catch (error) {
        console.error("Failed to load businesses:", error);
      } finally {
        setIsLoading(false);
      }
    };

    void loadBusinesses();
  }, [businesses.length, selectedBusiness, setBusinesses, setSelectedBusiness]);

  const handleSelectBusiness = (businessId: string) => {
    const business = businesses.find((b) => b.id === businessId);
    if (business) {
      setSelectedBusiness(business);
      setIsOpen(false);
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center gap-2 px-3 py-2 bg-base-200 rounded-lg">
        <Loader2 size={16} className="animate-spin text-base-content/60" />
        <span className="text-sm text-base-content/60">
          {t("common.loading")}
        </span>
      </div>
    );
  }

  if (businesses.length === 0) {
    return null;
  }

  return (
    <div className="dropdown dropdown-end">
      {/* Trigger Button */}
      <button
        type="button"
        tabIndex={0}
        onClick={() => { setIsOpen(!isOpen); }}
        className="btn btn-ghost gap-2 h-auto min-h-0 py-2 px-3 hover:bg-base-200"
      >
        <Avatar
          src={selectedBusiness?.logoUrl}
          fallback={selectedBusiness?.name}
          size="sm"
          shape="square"
        />
        <div className="flex flex-col items-start">
          <span className="text-sm font-semibold text-base-content">
            {selectedBusiness?.name ?? t("dashboard.select_business")}
          </span>
          {selectedBusiness?.brand && (
            <span className="text-xs text-base-content/60">
              {selectedBusiness.brand}
            </span>
          )}
        </div>
        <ChevronDown
          size={16}
          className={cn(
            "transition-transform text-base-content/60",
            isOpen && "rotate-180"
          )}
        />
      </button>

      {/* Dropdown Menu */}
      {isOpen && (
        <ul
          tabIndex={0}
          className="dropdown-content menu p-2 shadow-lg bg-base-100 rounded-lg w-64 mt-2 border border-base-300"
        >
          {businesses.length === 0 ? (
            <li className="px-4 py-3 text-center text-sm text-base-content/60">
              {t("dashboard.no_businesses")}
            </li>
          ) : (
            businesses.map((business) => (
              <li key={business.id}>
                <button
                  type="button"
                  onClick={() => { handleSelectBusiness(business.id); }}
                  className={cn(
                    "flex items-center gap-3 px-3 py-2 text-start",
                    selectedBusiness?.id === business.id &&
                      "bg-primary-50 text-primary-700"
                  )}
                >
                  <Avatar
                    src={business.logoUrl}
                    fallback={business.name}
                    size="sm"
                    shape="square"
                  />
                  <div className="flex-1 min-w-0">
                    <div className="font-medium text-sm truncate">
                      {business.name}
                    </div>
                    {business.brand && (
                      <div className="text-xs text-base-content/60 truncate">
                        {business.brand}
                      </div>
                    )}
                  </div>
                  {selectedBusiness?.id === business.id && (
                    <Check size={16} className="text-primary-600 flex-shrink-0" />
                  )}
                </button>
              </li>
            ))
          )}

          {/* Divider and Add Business Option */}
          {businesses.length > 0 && (
            <>
              <div className="divider my-1"></div>
              <li>
                <button
                  type="button"
                  className="flex items-center gap-2 px-3 py-2 text-primary-600 hover:bg-primary-50"
                >
                  <Building2 size={16} />
                  <span className="text-sm font-medium">
                    {t("dashboard.add_business")}
                  </span>
                </button>
              </li>
            </>
          )}
        </ul>
      )}
    </div>
  );
}
