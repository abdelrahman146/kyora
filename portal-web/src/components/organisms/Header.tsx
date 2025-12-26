import { useTranslation } from "react-i18next";
import { Menu } from "lucide-react";
import { IconButton } from "../atoms/IconButton";
import { BusinessSwitcher } from "../molecules/BusinessSwitcher";
import { UserMenu } from "../molecules/UserMenu";
import { useBusinessStore } from "../../stores/businessStore";
import { useLanguage } from "../../hooks/useLanguage";
import { cn } from "@/lib/utils";

interface HeaderProps {
  title?: string;
}

/**
 * Header Component
 *
 * Top header bar containing business switcher, page title, and user menu.
 *
 * Features:
 * - Business switcher dropdown
 * - Page title (responsive)
 * - User profile menu
 * - Mobile hamburger menu
 * - RTL support
 * - Responsive padding based on sidebar state
 */
export function Header({ title }: HeaderProps) {
  const { t } = useTranslation();
  const { toggleLanguage } = useLanguage();
  const { isSidebarCollapsed, toggleSidebar } = useBusinessStore();

  return (
    <header
      className={cn(
        "fixed top-0 end-0 h-16 bg-base-100 border-b border-base-300 z-30 transition-all duration-300",
        "md:start-64",
        isSidebarCollapsed && "md:start-20",
        "start-0"
      )}
    >
      <div className="h-full flex items-center justify-between gap-4 px-4">
        {/* Left Section: Mobile Menu + Title */}
        <div className="flex items-center gap-3">
          {/* Mobile Menu Toggle */}
          <IconButton
            icon={Menu}
            size="md"
            variant="ghost"
            onClick={toggleSidebar}
            aria-label={t("dashboard.toggle_menu")}
            className="md:hidden"
          />

          {/* Page Title */}
          {title && (
            <h1 className="text-lg font-bold text-base-content hidden sm:block">
              {title}
            </h1>
          )}
        </div>

        {/* Right Section: Business Switcher + Language Toggle + User Menu */}
        <div className="flex items-center gap-2">
          {/* Business Switcher */}
          <BusinessSwitcher />

          {/* Language Toggle */}
          <button
            type="button"
            onClick={toggleLanguage}
            className="btn btn-ghost btn-sm"
          >
            <span className="text-sm font-medium">عربي / EN</span>
          </button>

          {/* User Menu */}
          <UserMenu />
        </div>
      </div>
    </header>
  );
}
